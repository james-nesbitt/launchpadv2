//go:build rig_integration
// +build rig_integration

package righost_test

import (
	"context"
	"testing"

	righost "github.com/Mirantis/launchpad/pkg/host/rig"
	"gopkg.in/yaml.v3"
)

var (
	rigHostsYaml = "" // FILL THIS WITH A TEST ARGUMENT
)

func Test_RigDecode_Integration(t *testing.T) {
	ctx := context.Background()
	nsh := map[string]yaml.Node{}

	if err := yaml.Unmarshal([]byte(rigHostsYaml), &nsh); err != nil {
		t.Errorf("Couldn't unmarshall hosts yaml: %s", err.Error())
	}

	for id, rhyn := range nsh {
		rh, err := righost.Decode(id, rhyn.Decode)
		if err != nil {
			t.Errorf("Rig host decode fail [%s] : %s", id, err.Error())
			continue
		}

		if rh.Id() != id {
			t.Errorf("Rig Host has wrong ID: %+v", rh)
		}

		var cmd string
		if rh.IsWindows() {
			cmd = "dir"
		} else {
			cmd = "ls -la"
		}

		if o, e, err := rh.Exec(ctx, cmd, nil); err != nil {
			t.Errorf("Rig host exec err [%s] : %s -> %+v", id, err.Error(), rh)
		} else {

			if len(o) == 0 {
				t.Error("Rig host exec output was empty")
			}
			if len(e) != 0 {
				t.Error("Rig host exec error output was not empty")
			}
		}
	}
}
