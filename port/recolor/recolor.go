package recolor

import (
	"image"
	"image/color"
)

func Tint(in image.Image, tint color.RGBA) image.Image {
	bounds := in.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixelColor := in.At(x, y)
			r, g, b, a := pixelColor.RGBA()
			if a < 1024 {
				continue
			}
			r >>= 8
			g >>= 8
			b >>= 8
			a >>= 8
			r = (r + uint32(tint.R)) / 2
			g = (r + uint32(tint.G)) / 2
			b = (r + uint32(tint.B)) / 2

			dst.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}
	return dst
}
