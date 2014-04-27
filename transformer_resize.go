package slimgfast

import (
	"github.com/nfnt/resize"
	"image"
)

// TransformerResize is the primary Transformer that will resize images to the
// proper size.
type TransformerResize struct{}

// Transform resizes the image as requested.
func (t *TransformerResize) Transform(req *ImageRequest, image image.Image) (image.Image, error) {
	resized := resize.Resize(
		uint(req.Width),
		uint(req.Height),
		image,
		resize.Lanczos3,
	)
	// TODO: Inspect the req.Fit attribute and either crop or scale it,
	//       as necessary.
	return resized, nil
}
