package slimgfast

import (
	"image"
)

// Transformer is the interface that any image transformation operations must
// adhere to.  By implementing one Transform method which takes the
// ImageRequest and an image object and returns a new image object, we can
// achieve any effect we want.
type Transformer interface {
	Transform(req *ImageRequest, image image.Image) (image.Image, error)
}
