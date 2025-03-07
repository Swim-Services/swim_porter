package recolor

import (
	"image"
	"image/color"
	"math"

	"github.com/crazy3lf/colorconv"
	"github.com/swim-services/swim_porter/utils"
)

type Hue struct {
	hue float64
}

func NewHue(color color.RGBA) *Hue {
	return &Hue{hue: float64(GetHue(int(color.R), int(color.G), int(color.B)))}
}

func (h *Hue) SetColor(color color.RGBA) {
	h.hue = float64(GetHue(int(color.R), int(color.G), int(color.B)))
}

func (h *Hue) RecolorImage(in image.Image, fileName string) (image.Image, error) {
	bounds := in.Bounds()
	out := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := in.At(x, y).RGBA()
			if a < 1024 {
				continue
			}
			_, s, v := colorconv.RGBToHSV(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			outR, outG, outB, err := colorconv.HSVToRGB(h.hue, s, v)
			if err != nil {
				return nil, err
			}
			out.Set(x, y, color.RGBA{outR, outG, outB, uint8(a >> 8)})
		}
	}
	return out, nil
}

func (h *Hue) DefaultList() []string {
	return append(utils.DEFAULT_RECOLOR_LIST, "diamond_ore")
}

func GetHue(red, green, blue int) int {
	min := math.Min(math.Min(float64(red), float64(green)), float64(blue))
	max := math.Max(math.Max(float64(red), float64(green)), float64(blue))

	if min == max {
		return 0
	}

	var hue float64
	if max == float64(red) {
		hue = (float64(green) - float64(blue)) / (max - min)
	} else if max == float64(green) {
		hue = 2 + (float64(blue)-float64(red))/(max-min)
	} else {
		hue = 4 + (float64(red)-float64(green))/(max-min)
	}

	hue = hue * 60
	if hue < 0 {
		hue = hue + 360
	}

	return int(math.Round(hue))
}
