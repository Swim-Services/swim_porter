package recolor

import (
	"image"
	"image/color"
)

type Algorithm interface {
	SetColor(c color.RGBA)
	RecolorImage(in image.Image, fileName string) (image.Image, error)
	DefaultList() []string
}

func FromString(s string) (Algorithm, bool) {
	switch s {
	case "hue":
		return &Hue{}, true
	case "hue_v2":
		return &HueV2{}, true
	case "tint":
		return &Tint{}, true
	case "gray_tint":
		return &GrayTint{}, true
	}
	return nil, false
}
