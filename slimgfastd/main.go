package main

import (
	"flag"
	"fmt"
	"github.com/ericflo/slimgfast"
	"github.com/ericflo/slimgfast/fetchers"
	"github.com/golang/groupcache"
	"log"
	"net/http"
	"os"
)

var COUNTER_FILENAME = *flag.String(
	"counter_filename",
	"/tmp/slimfast_sizes.json",
	"The file where we'll save statistical information about which sizes were requested",
)
var GROUPCACHE_HOSTS = *flag.String(
	"groupcache_hosts",
	"http://localhost:4401",
	"The URL prefix that you would like to assign to groupcache",
)
var PORT = *flag.String("port", "4400", "The port to serve images on")
var NUM_WORKERS = *flag.Int(
	"num_workers",
	4,
	"The number of worker goroutines to spawn",
)
var OUTPUT_CACHE_MB = int64(*flag.Int(
	"output_cache_mb",
	512,
	"The amount of cache to reserve for resized images",
))
var MAX_WIDTH = *flag.Int(
	"max_width",
	2048,
	"The max width of the resized image",
)
var MAX_HEIGHT = *flag.Int(
	"max_height",
	2048,
	"The max height of the resized image",
)

func parseFlags() slimgfast.Fetcher {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage of slimgfastd: slimgfastd [OPTIONS] COMMAND PREFIX\n")
		fmt.Fprintf(os.Stderr, "Available commands: proxy, filesystem\n")
		fmt.Fprintf(os.Stderr, "Note: PREFIX is the URL prefix for proxying, or the file path prefix for filesystem\n\n")
		fmt.Fprintf(os.Stderr, "Example: slimgfastd -num_workers 8 proxy http://i.imgur.com\n")
		fmt.Fprintf(os.Stderr, "Example: slimgfastd -output_cache_mb 128 filesystem /srv/project/static/images\n\n")
		fmt.Fprintf(os.Stderr, "Defaults:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n\n")
		os.Exit(1)
	}

	// Set up the fetcher
	flag.Parse()

	command := flag.Arg(0)
	if command == "" {
		flag.Usage()
	}

	var fetcher slimgfast.Fetcher
	if command == "proxy" || command == "filesystem" {
		prefix := flag.Arg(1)
		if prefix == "" {
			flag.Usage()
		}
		if command == "proxy" {
			fetcher = &fetchers.ProxyFetcher{ProxyUrlPrefix: prefix}
		} else {
			fetcher = &fetchers.FilesystemFetcher{PathPrefix: prefix}
		}
	} else {
		if command == "s3" {
			log.Println("Sorry, S3 hasn't been implemented yet in the daemon.")
		}
		flag.Usage()
	}
	return fetcher
}

func main() {
	fetcher := parseFlags()

	// Instantiate the transformers
	resizeTransformer := &slimgfast.TransformerResize{}
	transformers := []slimgfast.Transformer{resizeTransformer}

	// Create the app
	app, err := slimgfast.NewApp(
		fetcher,
		transformers,
		COUNTER_FILENAME,
		NUM_WORKERS,
		OUTPUT_CACHE_MB,
		MAX_WIDTH,
		MAX_HEIGHT,
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
