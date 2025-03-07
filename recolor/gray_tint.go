package recolor

import (
	"image"
	"image/color"
	"math"

	"github.com/disintegration/imaging"
	"github.com/swim-services/swim_porter/utils"
)

type GrayTint struct {
	color color.RGBA
}

func NewGrayTint(color color.RGBA) *GrayTint {
	return &GrayTint{color: color}
}

func (g *GrayTint) SetColor(color color.RGBA) {
	g.color = color
}

func (t *GrayTint) RecolorImage(in image.Image) (image.Image, error) {
	in = imaging.Grayscale(in)
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

			r = uint32(math.Min(float64(r)*(float64(t.color.R)/128.0), 255))
			g = uint32(math.Min(float64(g)*(float64(t.color.G)/128.0), 255))
			b = uint32(math.Min(float64(b)*(float64(t.color.B)/128.0), 255))

			dst.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}
	return dst, nil
}

func GrayTintRange(in image.Image, tint color.RGBA, start image.Point, end image.Point, mult float64) image.Image {
	dst := imaging.Clone(in)
	in = imaging.Grayscale(in)
	for y := start.Y; y < end.Y; y++ {
		for x := start.X; x < end.X; x++ {
			pixelColor := in.At(x, y)
			r, g, b, a := pixelColor.RGBA()
			if a < 1024 {
				continue
			}
			r >>= 8
			g >>= 8
			b >>= 8
			a >>= 8

			r = uint32(math.Min(float64(r)*(float64(tint.R)/(256/mult)), 255))
			g = uint32(math.Min(float64(g)*(float64(tint.G)/(256/mult)), 255))
			b = uint32(math.Min(float64(b)*(float64(tint.B)/(256/mult)), 255))

			dst.SetNRGBA(x, y, color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}
	return dst
}

func (g *GrayTint) DefaultList() []string {
	return utils.DEFAULT_RECOLOR_LIST
}
