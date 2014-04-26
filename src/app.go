package slimgfast

import (
	"fmt"
	"github.com/golang/groupcache"
	"net/http"
)

const RESIZED_IMAGE_SOURCE_NAME string = "slimgfast_resized_image_source"

type App struct {
	SizeCounter *SizeCounter
	cache       *groupcache.Group
	workerGroup *WorkerGroup
}

func NewApp(sizeCounter *SizeCounter, fetcher Fetcher, numWorkers int, cacheMegabytes int64) *App {
	workerGroup := &WorkerGroup{NumWorkers: numWorkers}
	app := &App{
		SizeCounter: sizeCounter,
		workerGroup: workerGroup,
	}
	imageSource := NewImageSource(fetcher)
	app.cache = groupcache.NewGroup(RESIZED_IMAGE_SOURCE_NAME, cacheMegabytes<<20, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			req, err := ImageRequestFromCacheKey(key)
			if err != nil {
				return err
			}
			resizedData, err := workerGroup.Resize(imageSource, req)
			if err != nil {
				return err
			}
			dest.SetBytes(resizedData)
			return nil
		}))
	return app
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := ImageRequestFromURLString(r.URL.String())
	if err != nil {
		handleError(http.StatusNotFound, err.Error(), w, r)
		return
	}

	if size, err := req.Size(); err != nil {
		// We don't care to capture stats about requests with no size
	} else {
		app.SizeCounter.CountSize(size)
	}

	var resizedData []byte
	imgSink := groupcache.AllocatingByteSliceSink(&resizedData)
	cacheKey, err := req.CacheKey()
	if err != nil {
		handleError(http.StatusNotFound, err.Error(), w, r)
		return
	}
	err = app.cache.Get(nil, cacheKey, imgSink)
	if err != nil {
		handleError(http.StatusNotFound, err.Error(), w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "image/jpeg")
	w.Write(resizedData)
}

func (app *App) Start() {
	app.workerGroup.Start()
}

func (app *App) Close() {
	app.workerGroup.Close()
}

func handleError(status int, content string, w http.ResponseWriter, r *http.Request) {
	// TODO: Generate an error image in the correct dimensions
	w.WriteHeader(status)
	fmt.Fprint(w, content)
}
