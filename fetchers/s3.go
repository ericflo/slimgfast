package fetchers

import (
	"errors"
	"github.com/golang/groupcache"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"strings"
)

// S3Fetcher fetches images from an S3 bucket.
type S3Fetcher struct {
	Auth   aws.Auth
	Region aws.Region
}

// parseURL looks at the request URL path to determine the AWS bucket and
// filename.
func parseS3Url(f *S3Fetcher, urlPath string) (string, string, error) {
	pathSegments := strings.Split(urlPath, "/")
	if len(pathSegments) < 3 {
		return "", "", errors.New("Url needs to be /BUCKET/filename.jpg")
	}
	bucketname := pathSegments[1]
	filename := strings.Join(pathSegments[2:], "/")
	return bucketname, filename, nil
}

// Fetch grabs the image data from the bucket and filename requested by the
// user.
func (f *S3Fetcher) Fetch(urlPath string, dest groupcache.Sink) error {
	bucketname, filename, err := parseS3Url(f, urlPath)
	if err != nil {
		return err
	}
	conn := s3.New(f.Auth, f.Region)
	bucket := conn.Bucket(bucketname)
	if data, err := bucket.Get(filename); err != nil {
		return err
	} else {
		dest.SetBytes(data)
	}
	return nil
}
