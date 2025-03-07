package recolor

import (
	"image"
	"image/color"
	"strings"

	"github.com/crazy3lf/colorconv"
	"github.com/swim-services/swim_porter/utils"
)

type HueV2 struct {
	h, s, v float64
}

func NewHueV2(color color.RGBA) *HueV2 {
	h, s, v := colorconv.RGBToHSV(color.R, color.G, color.B)
	return &HueV2{h: h, s: s, v: v}
}

func (h *HueV2) SetColor(color color.RGBA) {
	h.h, h.s, h.v = colorconv.RGBToHSV(color.R, color.G, color.B)
}

func (h *HueV2) RecolorImage(in image.Image, fileName string) (image.Image, error) {
	bounds := in.Bounds()
	out := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := in.At(x, y).RGBA()
			if a < 1024 {
				continue
			}
			_, s, v := colorconv.RGBToHSV(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			if s > 0.15 || strings.Contains(fileName, "overworld_cubemap") {
				s = min(1, s*h.s*1.2)
				v = min(1, v*h.v*1.2)
			}
			outR, outG, outB, err := colorconv.HSVToRGB(h.h, s, v)
			if err != nil {
				return nil, err
			}
			out.Set(x, y, color.RGBA{outR, outG, outB, uint8(a >> 8)})
		}
	}
	return out, nil
}

func (h *HueV2) DefaultList() []string {
	return append(utils.DEFAULT_RECOLOR_LIST, "diamond_ore")
}
