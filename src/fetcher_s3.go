package slimgfast

import (
	"errors"
	"github.com/golang/groupcache"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"net/url"
	"strings"
)

type S3Fetcher struct {
	Auth     aws.Auth
	Region   aws.Region
	bucket   string
	filename string
}

func (f *S3Fetcher) ParseURL(rawUrl string) error {
	parsedUrl, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return err
	}
	pathSegments := strings.Split(parsedUrl.Path, "/")
	if len(pathSegments) < 3 {
		return errors.New("Url needs to be /BUCKET/filename.jpg")
	}
	f.bucket = pathSegments[1]
	f.filename = strings.Join(pathSegments[2:], "/")
	return nil
}

func (f *S3Fetcher) Fetch(req *ImageRequest, dest groupcache.Sink) error {
	conn := s3.New(f.Auth, f.Region)
	bucket := conn.Bucket(f.bucket)
	if data, err := bucket.Get(f.filename); err != nil {
		return err
	} else {
		dest.SetBytes(data)
	}
	return nil
}
