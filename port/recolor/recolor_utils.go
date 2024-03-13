package recolor

import (
	"image"
	"image/color"
	"math"

	"github.com/crazy3lf/colorconv"
	"github.com/disintegration/imaging"
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
			g = (g + uint32(tint.G)) / 2
			b = (b + uint32(tint.B)) / 2

			dst.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}
	return dst
}

func GrayTint(in image.Image, tint color.RGBA) image.Image {
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

			r = uint32(math.Min(float64(r)*(float64(tint.R)/128.0), 255))
			g = uint32(math.Min(float64(g)*(float64(tint.G)/128.0), 255))
			b = uint32(math.Min(float64(b)*(float64(tint.B)/128.0), 255))

			dst.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}
	return dst
}

func HueShift(in image.Image, iHUE float64) (image.Image, error) {
	bounds := in.Bounds()
	out := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := in.At(x, y).RGBA()
			if a < 1024 {
				continue
			}
			_, s, v := colorconv.RGBToHSV(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			outR, outG, outB, err := colorconv.HSVToRGB(iHUE, s, v)
			if err != nil {
				return nil, err
			}
			out.Set(x, y, color.RGBA{outR, outG, outB, uint8(a >> 8)})
		}
	}
	return out, nil
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
