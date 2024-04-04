package decode

import (
	"gopkg.in/yaml.v3"
)

func DecodeTestYaml(hyb []byte) func(interface{}) error {
	dyn := make(chan *yaml.Node)

	go func() {
		dhc := dummyContainer{yn: dyn}
		yaml.Unmarshal(hyb, &dhc)
	}()

	py := <-dyn
	return py.Decode
}

type dummyContainer struct {
	yn chan *yaml.Node
}

func (dhc *dummyContainer) UnmarshalYAML(yn *yaml.Node) error {
	dhc.yn <- yn
	return nil
}
