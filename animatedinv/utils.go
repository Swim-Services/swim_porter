package animatedinv

import (
	"image"
	"image/color"
	"image/gif"
	"iter"

	"github.com/disintegration/imaging"
	"github.com/swim-services/swim_porter/internal"
)

func cloneMap(in map[string][]byte) map[string][]byte {
	newMap := make(map[string][]byte)
	for key, value := range in {
		newVal := make([]byte, len(value))
		copy(newVal, value)
		newMap[key] = newVal
	}
	return newMap
}

func SplitAnimatedGIF(gifImg *gif.GIF) []*image.NRGBA {
	var out = make([]*image.NRGBA, len(gifImg.Image))

	imgWidth, imgHeight := getGifDimensions(gifImg)

	overpaintImage := image.NewNRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	for i, srcImg := range gifImg.Image {
		drawOverAsync(overpaintImage, srcImg)
		out[i] = imaging.Clone(overpaintImage)
	}
	return out
}

func AnimatedGIFIter(gifImg *gif.GIF) iter.Seq[*image.NRGBA] {
	imgWidth, imgHeight := getGifDimensions(gifImg)
	overpaintImage := image.NewNRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	return func(yield func(*image.NRGBA) bool) {
		for _, srcImg := range gifImg.Image {
			drawOverAsync(overpaintImage, srcImg)
			if !yield(overpaintImage) {
				return
			}
		}
	}
}
func getGifDimensions(gif *gif.GIF) (x, y int) {
	var lowestX int
	var lowestY int
	var highestX int
	var highestY int

	for _, img := range gif.Image {
		if img.Rect.Min.X < lowestX {
			lowestX = img.Rect.Min.X
		}
		if img.Rect.Min.Y < lowestY {
			lowestY = img.Rect.Min.Y
		}
		if img.Rect.Max.X > highestX {
			highestX = img.Rect.Max.X
		}
		if img.Rect.Max.Y > highestY {
			highestY = img.Rect.Max.Y
		}
	}

	return highestX - lowestX, highestY - lowestY
}

func drawOverAsync(dst *image.NRGBA, src *image.Paletted) {
	internal.Parallel(src.Bounds().Min.X, src.Bounds().Max.X, func(c <-chan int) {
		for x := range c {
			for y := src.Bounds().Min.Y; y < src.Bounds().Max.Y; y++ {
				col := src.RGBA64At(x, y)
				dst.SetRGBA64(x, y, blendColors(dst.RGBA64At(x, y), col, float64(col.A)/65535))
			}
		}
	})
}

func blendColors(c1, c2 color.RGBA64, alpha float64) color.RGBA64 {
	r := uint16(float64(c1.R)*(1-alpha) + float64(c2.R)*alpha)
	g := uint16(float64(c1.G)*(1-alpha) + float64(c2.G)*alpha)
	b := uint16(float64(c1.B)*(1-alpha) + float64(c2.B)*alpha)

	return color.RGBA64{R: r, G: g, B: b, A: 65535}
}
