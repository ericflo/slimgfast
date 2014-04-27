package slimgfast

import (
	"github.com/golang/groupcache"
	"io/ioutil"
	"log"
	"net/url"
	"path"
)

// FilesystemFetcher fetches images from the filesystem.
type FilesystemFetcher struct {
	PathPrefix string
	path       string
}

// ParseURL looks at the base directory, appends the path of the request URL,
// and arrives at a path to the image on the filesystem.
func (f *FilesystemFetcher) ParseURL(rawUrl string) error {
	parsedUrl, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return err
	}
	f.path = path.Clean(f.PathPrefix + parsedUrl.Path)
	log.Println(f.path)
	return nil
}

// Fetch opens and reads in the image data from the file determined by ParseURL.
func (f *FilesystemFetcher) Fetch(req *ImageRequest, dest groupcache.Sink) error {
	data, err := ioutil.ReadFile(f.path)
	if err != nil {
		return err
	}
	dest.SetBytes(data)
	return nil
}
