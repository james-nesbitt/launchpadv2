package implementation

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Mirantis/launchpad/pkg/host/exec"
)

func NewRpmOSTreeExecClient(e exec.HostExecutor, opts exec.ExecOptions) RpmOSTreeExecClient {
	return RpmOSTreeExecClient{
		e:    e,
		opts: opts,
	}
}

type RpmOSTreeExecClient struct {
	e    exec.HostExecutor
	opts exec.ExecOptions
}

func (c RpmOSTreeExecClient) Status(ctx context.Context) ([]RpmOSTreeStatusDeployment, error) {
	cmd := rpmOSTreeCmd([]string{"status", "--verbose", "--json"})

	var sds struct {
		Deployments []RpmOSTreeStatusDeployment `json:"deployments"`
	}

	o, e, err := c.e.Exec(ctx, cmd, nil, c.opts)
	if err != nil {
		return sds.Deployments, fmt.Errorf("Error getting RPMOSTRee status: %s :: %s", err.Error(), e)
	}

	if err := json.Unmarshal([]byte(o), &sds); err != nil {
		return sds.Deployments, fmt.Errorf("Error parsing RPMOSTRee status: %s :: %s", err.Error(), e)
	}

	return sds.Deployments, nil
}

func rpmOSTreeCmd(cmd []string) string {
	return strings.Join(append([]string{"rpm-ostree"}, cmd...), " ")
}
