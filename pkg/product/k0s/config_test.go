package k0s_test

import (
	"net/http"
	"testing"

	"github.com/Mirantis/launchpad/pkg/product/k0s"
	"github.com/k0sproject/version"
)

func Test_ConfigURL(t *testing.T) {
	arch := "amd64"
	vs := "v1.30.0+k0s.0"
	v, verr := version.NewVersion(vs)
	if verr != verr {
		t.Errorf("err with version %s: %s", vs, verr.Error())
	}

	config := k0s.Config{
		Version: *v,
	}

	url := config.DownloadURL(arch)

	r, err := http.Get(url)
	if err != nil {
		t.Errorf("couldn't download version:")
	}
	r.Body.Close()
}
