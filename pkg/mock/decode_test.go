package mock_test

import (
	"testing"

	"github.com/Mirantis/launchpad/pkg/mock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func Test_DecodeBasicYaml(t *testing.T) {
	y := "" // put more here when we have more in the decode spec
	var yn yaml.Node

	assert.Nil(t, yaml.Unmarshal([]byte(y), &yn), "basic unmarshal failure")

	cid := "test"
	c, derr := mock.DecodeComponent(cid, yn.Decode)
	assert.Nil(t, derr, "decode error occurred")

	assert.Equal(t, cid, c.Name(), "wrong name stored in mock component")
}
