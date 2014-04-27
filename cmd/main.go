package main

import (
	"flag"
	"github.com/ericflo/slimgfast"
	"github.com/golang/groupcache"
	"net/http"
)

var COUNTER_FILENAME string = slimgfast.GetEnvString(
	"SLIMGFAST_COUNTER_FILENAME", "/tmp/sizes.json")
var GROUPCACHE_HOSTS string = slimgfast.GetEnvString(
	"SLIMGFAST_GROUPCACHE_HOSTS", "http://localhost:4401")
var PORT string = slimgfast.GetEnvString("SLIMGFAST_PORT", "4400")
var NUM_WORKERS int = slimgfast.GetEnvInt("SLIMGFAST_NUM_WORKERS", 4)
var THUMB_CACHE_MEGABYTES int64 = int64(slimgfast.GetEnvInt(
	"SLIMGFAST_THUMB_CACHE_MEGABYTES", 512))

func main() {
	// Set up the fetcher
	flag.Parse()
	prefix := flag.Arg(0)
	if prefix == "" {
		panic("Must pass the prefix to the command")
	}
	fetcher := &slimgfast.ProxyFetcher{ProxyUrlPrefix: prefix}
	//fetcher := &slimgfast.FilesystemFetcher{PathPrefix: prefix}

	// Instantiate the transformers
	resizeTransformer := &slimgfast.TransformerResize{}
	transformers := []slimgfast.Transformer{resizeTransformer}

	// Create the app
	app := slimgfast.NewApp(
		fetcher,
		transformers,
		COUNTER_FILENAME,
		NUM_WORKERS,
		THUMB_CACHE_MEGABYTES,
	)

	// Set up our groupcache pool
	peers := groupcache.NewHTTPPool(GROUPCACHE_HOSTS)
	go http.ListenAndServe(GROUPCACHE_HOSTS, http.HandlerFunc(peers.ServeHTTP))

	// Start the app
	app.Start()
	defer app.Close()

	// Start the HTTP server
	http.ListenAndServe(":"+PORT, app)
}
