package v21

import (
	"strings"

	"github.com/Mirantis/launchpad/pkg/component"
	"github.com/Mirantis/launchpad/pkg/config"
	"github.com/Mirantis/launchpad/pkg/project"
)

const (
	ID = "launchpad.mirantis.com/v2.1"
)

func init() {
	config.RegisterSpecDecoder(ID, Decode)
}

// Decode a project from the spec.
func Decode(cl *project.Project, d func(interface{}) error) error {
	var cs Spec

	if err := d(&cs); err != nil {
		return err
	}

	if cl.Components == nil {
		cl.Components = component.Components{}
	}

	for id, c := range cs.Components.components() {
		cl.Components[strings.ToUpper(id)] = c
	}

	return nil
}
