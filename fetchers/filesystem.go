package fetchers

import (
	"github.com/ericflo/slimgfast"
	"github.com/golang/groupcache"
	"io/ioutil"
	"log"
	"net/url"
	"path"
)

// FilesystemFetcher fetches images from the filesystem.
type FilesystemFetcher struct {
	PathPrefix string
}

// parseURL looks at the base directory, appends the path of the request URL,
// and arrives at a path to the image on the filesystem.
func parseFSUrl(f *FilesystemFetcher, rawUrl string) (string, error) {
	parsedUrl, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return "", err
	}
	filePath := path.Clean(f.PathPrefix + parsedUrl.Path)
	log.Println(filePath)
	return filePath, nil
}

// Fetch opens and reads in the image data from the file requested by the user.
func (f *FilesystemFetcher) Fetch(req *slimgfast.ImageRequest, dest groupcache.Sink) error {
	filePath, err := parseFSUrl(f, req.Url)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	dest.SetBytes(data)
	return nil
}
