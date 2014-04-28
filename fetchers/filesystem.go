package fetchers

import (
	"github.com/golang/groupcache"
	"io/ioutil"
	"path"
)

// FilesystemFetcher fetches images from the filesystem.
type FilesystemFetcher struct {
	PathPrefix string
}

// Fetch opens and reads in the image data from the file requested by the user.
func (f *FilesystemFetcher) Fetch(urlPath string, dest groupcache.Sink) error {
	filePath := path.Clean(f.PathPrefix + urlPath)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	dest.SetBytes(data)
	return nil
}
