package animatedinv

import (
	"image"
	"image/draw"
	"image/gif"

	"github.com/disintegration/imaging"
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

func SplitAnimatedGIF(gif *gif.GIF) []*image.NRGBA {
	var out = make([]*image.NRGBA, len(gif.Image))

	imgWidth, imgHeight := getGifDimensions(gif)

	overpaintImage := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	draw.Draw(overpaintImage, overpaintImage.Bounds(), imaging.Clone(gif.Image[0]), image.Point{}, draw.Src)

	for i, srcImg := range gif.Image {
		draw.Draw(overpaintImage, overpaintImage.Bounds(), srcImg, image.Point{}, draw.Over)
		out[i] = imaging.Clone(overpaintImage)
	}
	return out
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
