package blend

import (
	"errors"
	"image"
	"image/color"
	"image/draw"

	"github.com/disintegration/imaging"
)

func Blend(in image.Image, fadePercent int) (image.Image, error) {
	if fadePercent >= 50 || fadePercent < 1 {
		return nil, errors.New("fade percent must be between 1 and 49")
	}
	bounds := in.Bounds()
	newImg := image.NewRGBA(bounds)
	draw.Draw(newImg, bounds, in, image.Point{}, draw.Src)
	width := bounds.Dx()
	blendAmt := int(float32(width) * float32(fadePercent) / 100.0)
	for x := 0; x < blendAmt; x++ {
		for y := 0; y < bounds.Dy(); y++ {
			newImg.Set(width-x-1, y, blendColors(in.At(blendAmt-x, y), in.At(width-x-1, y), float64(x)/float64(blendAmt)))
			newImg.Set(x, y, blendColors(in.At(width-(blendAmt-x)-1, y), in.At(x, y), float64(x)/float64(blendAmt)))
		}
	}
	return imaging.Crop(newImg, image.Rect(blendAmt/2, 0, bounds.Dx()-(blendAmt/2), bounds.Dy())), nil
}

func blendColors(c1, c2 color.Color, alpha float64) color.RGBA {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()

	r := uint32(float64(r1)*(1-alpha) + float64(r2)*alpha)
	g := uint32(float64(g1)*(1-alpha) + float64(g2)*alpha)
	b := uint32(float64(b1)*(1-alpha) + float64(b2)*alpha)

	return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: 255}
}
