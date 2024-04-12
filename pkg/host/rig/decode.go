package righost

import "github.com/Mirantis/launchpad/pkg/host"

func init() {
	host.RegisterDecoder("rig", Decode)
}

func Decode(id string, d func(interface{}) error) (host.Host, error) {
	var rc Config
	if err := d(&rc); err != nil {
		return nil, err
	}
	rc.Id = id

	rh := NewHost(&rc)
	return &rh, nil
}
