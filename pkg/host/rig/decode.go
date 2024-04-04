package righost

import "github.com/Mirantis/launchpad/pkg/host"

func init() {
	host.RegisterDecoder("rig", RigHostDecode)
}

func RigHostDecode(d func(interface{}) error) (host.Host, error) {
	var rc Config
	if err := d(&rc); err != nil {
		return nil, err
	}

	rh := NewHost(&rc)
	return &rh, nil
}
