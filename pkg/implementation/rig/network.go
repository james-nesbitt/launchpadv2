package rig

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/Mirantis/launchpad/pkg/host/exec"
	"github.com/Mirantis/launchpad/pkg/host/network"
)

var (
	sbinPath = "PATH=/usr/local/sbin:/usr/sbin:/sbin:$PATH"
)

// Network retrieve the Network information for the machine.
func (p hostPlugin) Network(ctx context.Context) (network.Network, error) {
	n := network.Network{}

	if pi, err := p.resolvePrivateInterface(ctx); err != nil {
		return n, err
	} else {
		n.PrivateInterface = pi
	}

	if ip, err := p.resolvePublicIP(); err != nil {
		return n, err
	} else {
		n.PublicAddress = ip
	}

	if ip, err := p.resolveInternaIP(ctx, n.PrivateInterface, n.PublicAddress); err != nil {
		return n, err
	} else {
		n.PrivateAddress = ip
	}

	return n, nil
}

func (p hostPlugin) resolvePrivateInterface(ctx context.Context) (string, error) {
	cmd := fmt.Sprintf(`%s; (ip route list scope global | grep -P "\b(172|10|192\.168)\.") || (ip route list | grep -m1 default)`, sbinPath)
	o, e, err := p.Exec(ctx, cmd, nil, exec.ExecOptions{})
	if err != nil {
		return "", fmt.Errorf("could not detect private interface ;%w : %s", err, e)
	}

	re := regexp.MustCompile(`\bdev (\w+)`)
	match := re.FindSubmatch([]byte(o))
	if len(match) == 0 {
		return "", fmt.Errorf("can't find 'dev' in output")
	}
	return string(match[1]), nil
}

func (p hostPlugin) resolvePublicIP() (string, error) {
	return p.rig.ConnectionConfig.SSH.Address, nil
}

func (p hostPlugin) resolveInternaIP(ctx context.Context, privateInterface string, publicIP string) (string, error) {
	o, e, err := p.Exec(ctx, fmt.Sprintf("%s ip -o addr show dev %s scope global", sbinPath, privateInterface), nil, exec.ExecOptions{})
	if err != nil {
		return "", fmt.Errorf("%s: failed to find private interface: %s :: %s", p.hid(), err.Error(), e)
	}

	lines := strings.Split(o, "\n")
	for _, line := range lines {
		items := strings.Fields(line)
		if len(items) < 4 {
			//log.Debugf("not enough items in ip address line (%s), skipping...", items)
			continue
		}

		idx := strings.Index(items[3], "/")
		if idx == -1 {
			//log.Debugf("no CIDR mask in ip address line (%s), skipping...", items)
			continue
		}
		addr := items[3][:idx]

		if addr != publicIP {
			//log.Infof("%s: using %s as private IP", h, addr)
			if net.ParseIP(addr) != nil {
				return addr, nil
			}
		}
	}
	return publicIP, nil
}
