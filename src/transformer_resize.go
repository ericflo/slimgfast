package slimgfast

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
	// TODO: Inspect the req.Fit attribute and either crop or scale it,
	//       as necessary.
	return resized, nil
}
