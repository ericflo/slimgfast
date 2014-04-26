package slimgfast

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

type ImageRequest struct {
	Url    string
	Width  int
	Height int
	Fit    string
}

func ImageRequestFromURLString(rawUrl string) (*ImageRequest, error) {
	req := ImageRequest{Url: rawUrl}
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	// Purposely ignore the strcon errors here, if they're empty strings we'll
	// notice that later and serve the image at its original size.
	req.Width, _ = strconv.Atoi(parsedUrl.Query().Get("w"))
	req.Height, _ = strconv.Atoi(parsedUrl.Query().Get("h"))
	req.Fit = parsedUrl.Query().Get("fit")
	return &req, nil
}

func ImageRequestFromCacheKey(cacheKey string) (*ImageRequest, error) {
	req := ImageRequest{}
	err := json.Unmarshal([]byte(cacheKey), &req)
	return &req, err
}

func (req *ImageRequest) CacheKey() (string, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (req *ImageRequest) Size() (*Size, error) {
	if req.Width == 0 || req.Width == -1 {
		return nil, errors.New("Image does not have a specified width")
	}
	if req.Height == 0 || req.Width == -1 {
		return nil, errors.New("Image does not have a specified height")
	}
	return &Size{Width: uint(req.Width), Height: uint(req.Height)}, nil
}
