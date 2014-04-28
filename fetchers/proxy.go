package fetchers

import (
	"github.com/ericflo/slimgfast"
	"github.com/golang/groupcache"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// ProxyFetcher fetches images from an HTTP server.
type ProxyFetcher struct {
	ProxyUrlPrefix string
}

// parseURL looks at the proxy URL prefix, appends the path of the request URL,
// and arrives at a URL to the image we want to fetch.
func parseProxyUrl(f *ProxyFetcher, rawUrl string) (string, error) {
	parsedUrl, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return "", err
	}
	fullUrl := f.ProxyUrlPrefix + parsedUrl.Path
	log.Println(fullUrl)
	return fullUrl, nil
}

// Fetch makes an HTTP GET request to fetch the image data requested by the
// user.
func (f *ProxyFetcher) Fetch(req *slimgfast.ImageRequest, dest groupcache.Sink) error {
	fullUrl, err := parseProxyUrl(f, req.Url)
	if err != nil {
		return err
	}
	resp, err := http.Get(fullUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	dest.SetBytes(data)
	return nil
}
