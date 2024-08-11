package k0s

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/k0sproject/version"

	"github.com/Mirantis/launchpad/pkg/host/network"
	"github.com/k0sproject/dig"
)

const (
	RoleController    = "controller"
	ServiceController = "k0scontroller"
	RoleWorker        = "worker"
	ServiceWorker     = "k0sworker"

	DefaultK0sConfigPath = "/etc/k0s/k0sd.yaml"
	DefaultK0sDataPath   = "/var/lib/k0s"
	DefaultK0sBinaryPath = "/usr/bin/k0s"
	DefaultK0sTokenPath  = "/etc/k0s/k0s-token"
)

var (
	JoinTokenExpireDuration = time.Duration(10) * time.Minute
)

type k0sversion struct {
	K0s          string `json:"k0s"`
	RunC         string `json:"runc"`
	ContainerD   string `json:"containerd"`
	Kubernetes   string `json:"kubernetes"`
	Kine         string `json:"kine"`
	Etcd         string `json:"etcd"`
	Konnectivity string `json:"konnectivity"`
}

var (
	ErrK0sBinaryNotFound = errors.New("K0s binaries not installed onto the machine")
	ErrK0sNotRunning     = errors.New("K0s not running on the machine")
)

type k0sstatus struct {
	Version       *version.Version `json:"Version"`
	Pid           int              `json:"Pid"`
	PPid          int              `json:"PPid"`
	Role          string           `json:"Role"`
	SysInit       string           `json:"SysInit"`
	StubFile      string           `json:"StubFile"`
	Workloads     bool             `json:"Workloads"`
	Args          []string         `json:"Args"`
	ClusterConfig dig.Mapping      `json:"ClusterConfig"`
	K0sVars       dig.Mapping      `json:"K0sVars"`
}

func k0sErrorAnalyze(e string, err error) error {
	if strings.Contains(e, "can't get \"status\"") {
		return fmt.Errorf("%w; %s", ErrK0sNotRunning, err.Error())
	}
	if strings.Contains(e, "command not found") {
		return fmt.Errorf("%w; %s", ErrK0sBinaryNotFound, err.Error())
	}

	return err
}

// DownloadURL url to use for the k0s binary for the config version
func DownloadK0sURL(v version.Version, arch string) string {
	// https://github.com/k0sproject/k0s/releases/download/v1.30.0%2Bk0s.0/k0s-v1.30.0+k0s.0-amd64
	// v1.30.0%2Bk0s.0/k0s-v1.30.0+k0s.0-amd64
	return fmt.Sprintf("%[1]s/%[2]s/k0s-%[2]s-%[3]s", K0sReleaseLinkBase, v.String(), arch)
}

func (c *Component) CollectClusterSans(ctx context.Context) []string {

	chs, cerr := c.GetControllerHosts(ctx)
	if cerr != nil {
		slog.WarnContext(ctx, "error getting controller host list")
		return []string{}
	}
	if len(chs) == 0 {
		slog.WarnContext(ctx, "no controlers, so no sans can be collected")
		return []string{}
	}

	var sans []string

	for _, ch := range chs {
		chn := network.HostGetNetwork(ch)
		cn, nerr := chn.Network(ctx)
		if nerr != nil {
			slog.WarnContext(ctx, fmt.Sprintf("%s: host doesn't have any network features, so can't provide any san information", ch.Id()))
			continue
		}
		sans = append(sans, cn.PublicAddress)
		if cn.PrivateAddress != "" {
			sans = append(sans, cn.PrivateAddress)
		}
	}

	return sans
}
