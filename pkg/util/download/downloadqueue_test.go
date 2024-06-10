package download_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/Mirantis/launchpad/pkg/util/download"
)

func Test_DownloadQueue(t *testing.T) {
	ctx, c := context.WithTimeout(context.Background(), time.Second*60)
	defer c()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 1)
		w.Header().Add("content-disposition", `attachment; filename="filename.jpg"`)
		w.Write([]byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"))
	}))
	defer ts.Close()

	dq := download.NewQueueDownload(ts.Client())

	d1, fn, d1err := dq.Download(ctx, ts.URL)
	if d1err != nil {
		t.Fatalf("error occurred downloading: %s", d1err.Error())
	}
	if filepath.Base(fn) != "filename.jpg" {
		t.Fatalf("wrong filename returned: %s", fn)
	}
	defer d1.Close()
}
