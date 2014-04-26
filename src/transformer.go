package slimgfast

import (
	"image"
)

type Transformer interface {
	Transform(req *ImageRequest, image image.Image) (image.Image, error)
}
