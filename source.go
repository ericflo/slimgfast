package slimgfast

import (
	"github.com/golang/groupcache"
	"net/url"
)

const DEFAULT_IMAGE_SOURCE_NAME = "slimgfast_image_source"
const DEFAULT_CACHE_SIZE_MB = int64(128)

// ImageSource is an abstraction over a fetcher which caches intelligently, and
// serves as the primary internal interface to fetchers.
type ImageSource struct {
	cache *groupcache.Group
}

// NewImageSource initializes and returns an *ImageSource with sane default
// values.
func NewImageSource(fetcher Fetcher) *ImageSource {
	return NewImageSourceCustomCache(
		fetcher,
		DEFAULT_IMAGE_SOURCE_NAME,
		DEFAULT_CACHE_SIZE_MB,
	)
}

// NewImageSourceCustomCache initializes and returns an *ImageSource with
// a custom groupcache name and a custom cache size.
func NewImageSourceCustomCache(fetcher Fetcher, cacheName string, cacheMegabytes int64) *ImageSource {
	cache := groupcache.NewGroup(cacheName, cacheMegabytes<<20, groupcache.GetterFunc(
		func(ctx groupcache.Context, urlPath string, dest groupcache.Sink) error {
			return fetcher.Fetch(urlPath, dest)
		}))
	return &ImageSource{cache: cache}
}

// GetImageData gets the image data the request asked for, either from cache or
// from the associated Fetcher.
func (src *ImageSource) GetImageData(req *ImageRequest) ([]byte, error) {
	var img []byte
	imgSink := groupcache.AllocatingByteSliceSink(&img)
	parsedUrl, err := url.ParseRequestURI(req.Url)
	if err != nil {
		return img, err
	}
	err = src.cache.Get(nil, parsedUrl.Path, imgSink)
	return img, err
}
