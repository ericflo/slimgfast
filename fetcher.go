package slimgfast

import "github.com/golang/groupcache"

// Fetcher is the interface that is used to fetch images from some source,
// which could be the filesystem, a remote URL, or S3 -- but it could be from
// anywhere.
type Fetcher interface {
	Fetch(urlPath string, dest groupcache.Sink) error
}
