package recolor

import (
	"image"
	"image/color"

	"github.com/swim-services/swim_porter/utils"
)

type Tint struct {
	color color.RGBA
}

func NewTint(color color.RGBA) *Tint {
	return &Tint{color: color}
}

func (t *Tint) SetColor(color color.RGBA) {
	t.color = color
}

func (t *Tint) RecolorImage(in image.Image, fileName string) (image.Image, error) {
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
			r = (r + uint32(t.color.R)) / 2
			g = (g + uint32(t.color.G)) / 2
			b = (b + uint32(t.color.B)) / 2

			dst.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}
	return dst, nil
}

func (t *Tint) DefaultList() []string {
	return utils.DEFAULT_RECOLOR_LIST
}
