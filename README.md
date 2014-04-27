# Slimgfast

Slimgfast is a library that allows you to create a scalable, efficient, dynamic
image origin server.

As an example, if you want an image to be resized to 640x480, you could hit the
url at http://127.0.0.1/my-image.jpg?w=640&h=480 and slimgfast will dynamically
resize your image to the correct dimensions, serve it, and cache it for later.

Slimgfast comes with an executable which supports a baseline default, but most
advanced users will want to use it as a library for maximum configurability.

GoDoc: http://godoc.org/github.com/ericflo/slimgfast/src

### Pronunciation

It's pronounced like "slimmage fast" :)

## Getting Started

The easiest way to get a copy of slimfast is to use "go get":

    go get github.com/ericflo/slimgfast

In your code you can now import the library:

```go
import github.com/ericflo/slimgfast/src

// For example:
fetcher := &slimgfast.ProxyFetcher{ProxyUrlPrefix: "http://i.imgur.com"}
```

To see what the default executable can do, navigate to the newly-downloaded
slimgfast directory and run:

    go run main.go http://i.imgur.com

Now to load an image we can do:

    open http://localhost:4400/EgLrnVL.jpg

With no arguments it will just proxy the original image which lives at
http://i.imgur.com/EgLrnVL.jpg.  With arguments it will resize it:

    open http://localhost:4400/EgLrnVL.jpg\?w=300\&h=300

## Using Slimgfast as a library

The steps for setting up a slimfast instance are fairly straightforward:

* Create a **fetcher** that will know how to read images from the upstream
  source
* Create a list of **transformers**, or potential operations that can be
  applied to the image (e.g. resize)
* Instantiate an **app struct**, which collects all the fetchers and
  transformers and handles the actual HTTP requests
* Spin up the **groupcache** library so it knows who its peers are
* Start the app and the http server

In fact, this is all that
[main.go](https://github.com/ericflo/slimgfast/blob/master/main.go) is doing.

## Creating your own Fetcher

Creating a Fetcher is straightforward, you only have to implement the Fetcher
inteface, which means implementing the following:

```go
Fetch(req *ImageRequest, dest groupcache.Sink) error
ParseURL(rawUrl string) error
```

Since it's really not all that much code, here's the body of the filesystem
fetcher as an example:

```go
type FilesystemFetcher struct {
    PathPrefix string
    path       string
}

func (f *FilesystemFetcher) ParseURL(rawUrl string) error {
    parsedUrl, err := url.ParseRequestURI(rawUrl)
    if err != nil {
        return err
    }
    f.path = path.Clean(f.PathPrefix + parsedUrl.Path)
    return nil
}

func (f *FilesystemFetcher) Fetch(req *ImageRequest, dest groupcache.Sink) error {
    data, err := ioutil.ReadFile(f.path)
    if err != nil {
        return err
    }
    dest.SetBytes(data)
    return nil
}
```

## Creating your own Transformer

Creating a Transformer is similarly straightforward to creating a Fetcher,
you have to implement the Transformer interface, which has only one method:

```go
Transform(req *ImageRequest, image image.Image) (image.Image, error)
```

So, it takes an image, and the request, and then returns the transformed image
(or an error.)  Here's the body of the resize transformer as an example:

```go
import (
    "github.com/nfnt/resize"
    "image"
)

type TransformerResize struct{}

func (t *TransformerResize) Transform(req *ImageRequest, image image.Image) (image.Image, error) {
    resized := resize.Resize(
        uint(req.Width),
        uint(req.Height),
        image,
        resize.Lanczos3,
    )
    return resized, nil
}
```

You could easily write a transformer that uses ImageMagick or epeg, if you want
either more power or more performance.  Or you could write a filter to change
the brightness, or the contrast, or turn it black and white, or any other
interesting image transformation.  Since you have access to the request,
you can parse the querystring with whatever semantics makes sense for your
needs.

## I want something with commercial support

You should check out http://imgix.com/, which is a well-run startup that offers
a similar but more advanced commercial image service.

## Status

Status: Very, very alpha.  It started as some code I'd written for work, but
it ended up being more or less a complete rewrite, and this version hasn't seen
any production traffic, ever.  I'll remove this warning when I'm more confident
in it and have run it in production.