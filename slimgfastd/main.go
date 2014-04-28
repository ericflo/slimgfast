package main

import (
	"flag"
	"github.com/ericflo/slimgfast"
	"github.com/ericflo/slimgfast/fetchers"
	"github.com/golang/groupcache"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var COUNTER_FILENAME = getEnvString(
	"SLIMGFAST_COUNTER_FILENAME", "/tmp/sizes.json")
var GROUPCACHE_HOSTS = getEnvString(
	"SLIMGFAST_GROUPCACHE_HOSTS", "http://localhost:4401")
var PORT = getEnvString("SLIMGFAST_PORT", "4400")
var NUM_WORKERS = getEnvInt("SLIMGFAST_NUM_WORKERS", 4)
var THUMB_CACHE_MEGABYTES = int64(
	getEnvInt("SLIMGFAST_THUMB_CACHE_MEGABYTES", 512))

func main() {
	// Set up the fetcher
	flag.Parse()
	prefix := flag.Arg(0)
	if prefix == "" {
		panic("Must pass the prefix to the command")
	}
	fetcher := &fetchers.ProxyFetcher{ProxyUrlPrefix: prefix}
	//fetcher := &fetchers.FilesystemFetcher{PathPrefix: prefix}

	// Instantiate the transformers
	resizeTransformer := &slimgfast.TransformerResize{}
	transformers := []slimgfast.Transformer{resizeTransformer}

	// Create the app
	app, err := slimgfast.NewApp(
		fetcher,
		transformers,
		COUNTER_FILENAME,
		NUM_WORKERS,
		THUMB_CACHE_MEGABYTES,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Set up our groupcache pool
	peers := groupcache.NewHTTPPool(GROUPCACHE_HOSTS)
	go http.ListenAndServe(GROUPCACHE_HOSTS, http.HandlerFunc(peers.ServeHTTP))

	// Start the app
	app.Start()
	defer app.Close()

	// Start the HTTP server
	if err = http.ListenAndServe(":"+PORT, app); err != nil {
		log.Fatal(err)
	}
}

// UTILITIES

// environ builds a full mapping of environment variables
func environ() map[string]string {
	_env := make(map[string]string)
	for _, item := range os.Environ() {
		splits := strings.SplitN(item, "=", 2)
		_env[splits[0]] = splits[1]
	}
	return _env
}

// getEnvString tries first to get a string from the environment, but falls
// back on a default provided value.
func getEnvString(key, def string) string {
	resp, ok := environ()[key]
	if !ok {
		return def
	}
	return resp
}

// getEnvInt tries first to get and parse an int from the environment, but
// falls back on a default provided value.
func getEnvInt(key string, def int) int {
	rawVal, ok := environ()[key]
	if !ok {
		return def
	}
	resp, err := strconv.Atoi(rawVal)
	if err != nil {
		return def
	}
	return resp
}
