package slimgfast

import (
	"github.com/golang/groupcache"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type ProxyFetcher struct {
	ProxyUrlPrefix string
	url            string
}

func (f *ProxyFetcher) ParseURL(rawUrl string) error {
	parsedUrl, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return err
	}
	f.url = f.ProxyUrlPrefix + parsedUrl.Path
	log.Println(f.url)
	return nil
}

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
