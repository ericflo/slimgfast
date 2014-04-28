// Package slimfast is a library that allows you to create a service which will
// resize, serve, and cache images.
//
// For a full guide, visit https://www.github.com/ericflo/slimgfast
package slimgfast

import (
	"fmt"
	"github.com/golang/groupcache"
	"net/http"
	"time"
)

const RESIZED_IMAGE_SOURCE_NAME string = "slimgfast_resized_image_source"

// App ties together the fetcher, the transformers, handles the lifecycle of an
// image request, and does the actual HTTP serving.
type App struct {
	sizeCounter *SizeCounter
	cache       *groupcache.Group
	workerGroup *WorkerGroup
}

// NewApp returns an App that is initialized and ready to be started.
func NewApp(fetcher Fetcher, transformers []Transformer, counterFilename string, numWorkers int, cacheMegabytes int64) (*App, error) {
	workerGroup := &WorkerGroup{
		NumWorkers:   numWorkers,
		Transformers: transformers,
	}
	// Create a counter to track image size requests
	sizeCounter, err := NewSizeCounter(counterFilename)
	if err != nil {
		return nil, err
	}

	app := &App{
		sizeCounter: sizeCounter,
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
	return app, err
}

// ServeHTTP is responsible for actually kicking off the image transformations
// and serving the image back to the user who requested it.
func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := ImageRequestFromURLString(r.URL.String())
	if err != nil {
		handleError(http.StatusNotFound, err.Error(), w, r)
		return
	}

	if size, err := req.Size(); err != nil {
		// We don't care to capture stats about requests with no size
	} else {
		app.sizeCounter.CountSize(size)
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

// Start starts the application worker group and size counter goroutines.
func (app *App) Start() {
	app.workerGroup.Start()
	// Should we un-hardcode this? Does anyone care?
	app.sizeCounter.Start(1 * time.Second)
}

// Close signals to the worker group and size counter goroutines to exit.
func (app *App) Close() {
	app.workerGroup.Close()
	app.sizeCounter.Close()
}

// handleError handles any errors that happen in the HTTP request/response
// cycle.
func handleError(status int, content string, w http.ResponseWriter, r *http.Request) {
	// TODO: Generate an error image in the correct dimensions
	w.WriteHeader(status)
	fmt.Fprint(w, content)
}
