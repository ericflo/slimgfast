package fetchers

import (
	"errors"
	"fmt"
	"github.com/golang/groupcache"
	"io/ioutil"
	"net/http"
)

// ProxyFetcher fetches images from an HTTP server.
type ProxyFetcher struct {
	ProxyUrlPrefix string
}

// Fetch makes an HTTP GET request to fetch the image data requested by the
// user.
func (f *ProxyFetcher) Fetch(urlPath string, dest groupcache.Sink) error {
	fullUrl := f.ProxyUrlPrefix + urlPath
	resp, err := http.Get(fullUrl)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		errStr := fmt.Sprintf(
			"Got a bad status code back (expected 200, got %d)",
			resp.StatusCode,
		)
		return errors.New(errStr)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	dest.SetBytes(data)
	return nil
}
