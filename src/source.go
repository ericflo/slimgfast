package slimgfast

import (
	"github.com/golang/groupcache"
)

const DEFAULT_IMAGE_SOURCE_NAME string = "slimgfast_image_source"
const DEFAULT_CACHE_SIZE_MB int64 = 128

// Fetcher is the interface that is used to fetch images from some source,
// which could be the filesystem, a remote URL, or S3 -- but it could be from
// anywhere.
type Fetcher interface {
	Fetch(req *ImageRequest, dest groupcache.Sink) error
	ParseURL(rawUrl string) error
}

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
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			req, err := ImageRequestFromCacheKey(key)
			if err != nil {
				return err
			}
			err = fetcher.ParseURL(req.Url)
			if err != nil {
				return err
			}
			return fetcher.Fetch(req, dest)
		}))
	return &ImageSource{cache: cache}
}

// GetImageData gets the image data the request asked for, either from cache or
// from the associated Fetcher.
func (src *ImageSource) GetImageData(req *ImageRequest) ([]byte, error) {
	var img []byte
	imgSink := groupcache.AllocatingByteSliceSink(&img)
	cacheKey, err := req.CacheKey()
	if err != nil {
		return nil, err
	}
	err = src.cache.Get(nil, cacheKey, imgSink)
	return img, err
}
