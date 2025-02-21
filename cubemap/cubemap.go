package cubemap

import (
	"image"

	"github.com/disintegration/imaging"
)

func BuildCubemap(img image.Image) [6]image.Image {
	cubemap := img.(interface {
		SubImage(r image.Rectangle) image.Image
	})
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	x := width / 3
	y := height / 2
	var images [6]image.Image

	images[0] = cubemap.SubImage(image.Rect(height, 0, x+height, y))
	images[1] = cubemap.SubImage(image.Rect(0, y, x, y*2))
	images[2] = cubemap.SubImage(image.Rect(x, y, x*2, y*2))
	images[3] = cubemap.SubImage(image.Rect(height, y, x+height, y*2))
	images[4] = imaging.FlipV(imaging.FlipH(cubemap.SubImage(image.Rect(x, 0, x*2, y))))
	images[5] = imaging.FlipV(imaging.FlipH(cubemap.SubImage(image.Rect(0, 0, x, y))))

	return images
}
