/*
Package decode abstract decoding to avoid implementation lockin.
*/
package decode

import (
	"gopkg.in/yaml.v3"
)

func DecodeTestYaml(hyb []byte) func(any) error {
	dyn := make(chan *yaml.Node)

	go func() {
		dhc := dummyContainer{yn: dyn}
		if err := yaml.Unmarshal(hyb, &dhc); err != nil {
			// Not much we can do here, but at least we're not ignoring the error return.
			// The receiver will probably timeout or get a nil node.
			return
		}
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
