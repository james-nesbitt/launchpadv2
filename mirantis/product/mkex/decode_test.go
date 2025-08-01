package mkex_test

import (
	"testing"

	"github.com/Mirantis/launchpad/mirantis/product/mkex"
	"github.com/Mirantis/launchpad/pkg/component"
	"gopkg.in/yaml.v3"
)

func Test_DecodeSanity(t *testing.T) {
	ctx := t.Context()
	y := `
role: manager
`
	var yn yaml.Node

	if err := yaml.Unmarshal([]byte(y), &yn); err != nil {
		t.Fatal("failed to unmarshal")
	}

	c, cerr := component.DecodeComponent(mkex.ComponentType, "test", yn.Decode)
	if cerr != nil {
		t.Fatalf("error decoding: %s", cerr.Error())
	}

	if verr := c.Validate(ctx); verr != nil {
		t.Fatalf("invalid component created: %s", verr.Error())
	}
}
