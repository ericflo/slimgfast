package slimgfast

import (
	"github.com/golang/groupcache"
)

const DEFAULT_IMAGE_SOURCE_NAME string = "slimgfast_image_source"
const DEFAULT_CACHE_SIZE_MB int64 = 128

type Fetcher interface {
	Fetch(req *ImageRequest, dest groupcache.Sink) error
	ParseURL(rawUrl string) error
}

type ImageSource struct {
	cache *groupcache.Group
}

func NewImageSource(fetcher Fetcher) *ImageSource {
	return NewImageSourceCustomCache(
		fetcher,
		DEFAULT_IMAGE_SOURCE_NAME,
		DEFAULT_CACHE_SIZE_MB,
	)
}

func NewImageSourceCustomCache(fetcher Fetcher, cacheName string, cacheMegabytes int64) *ImageSource {
	cache := groupcache.NewGroup(cacheName, cacheMegabytes<<20, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			req, err := ImageRequestFromCacheKey(key)
			if err != nil {
				return err
			}
			return fetcher.Fetch(req, dest)
		}))
	return &ImageSource{cache: cache}
}

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
