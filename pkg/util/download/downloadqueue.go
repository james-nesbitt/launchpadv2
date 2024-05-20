package download

import (
	"context"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/adrg/xdg"
)

func NewQueueDownload(c *http.Client) *QueueDownload {
	if c == nil {
		c = http.DefaultClient
	}
	return &QueueDownload{
		c: c,
	}
}

type QueueDownload struct {
	c  *http.Client
	mu sync.Mutex
}

func (qd *QueueDownload) Download(ctx context.Context, url string) (io.ReadCloser, string, error) {
	irs, igerr := qd.c.Get(url)
	if igerr != nil {
		return nil, "", igerr
	}
	defer irs.Body.Close()

	cd := irs.Header.Get("Content-Disposition")
	_, params, _ := mime.ParseMediaType(cd)
	fs := filepath.Join("launchpad", "k0s", params["filename"])

	// only download one at a time (prevents duplicate downloads, and lowers bandwidth)
	qd.mu.Lock()
	defer qd.mu.Unlock()

	cfs, cfserr := xdg.SearchCacheFile(fs)
	if cfserr != nil {
		cfs, _ = xdg.CacheFile(fs)
		cf, cferr := os.Create(cfs)
		if cferr != nil {
			return nil, cfs, cferr
		}

		if _, err := io.Copy(cf, irs.Body); err != nil {
			return nil, cfs, err
		}
		if err := cf.Close(); err != nil {
			return nil, cfs, err
		}
	}

	cf, cferr := os.Open(cfs)
	if cferr != nil {
		return nil, cfs, cferr
	}

	return cf, cfs, nil
}
