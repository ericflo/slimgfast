package slimgfast

import (
	"github.com/golang/groupcache"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// ProxyFetcher fetches images from an HTTP server.
type ProxyFetcher struct {
	ProxyUrlPrefix string
	url            string
}

// ParseURL looks at the proxy URL prefix, appends the path of the request URL,
// and arrives at a URL to the image we want to fetch.
func (f *ProxyFetcher) ParseURL(rawUrl string) error {
	parsedUrl, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return err
	}
	f.url = f.ProxyUrlPrefix + parsedUrl.Path
	log.Println(f.url)
	return nil
}

// Fetch makes an HTTP GET request to fetch the image data from the URL
// determined by ParseURL
func (f *ProxyFetcher) Fetch(req *ImageRequest, dest groupcache.Sink) error {
	resp, err := http.Get(f.url)
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
