package cubemap

import (
	"image"
	"math"

	"github.com/disintegration/imaging"
	"github.com/swim-services/swim_porter/internal"
)

func CubemapToEquirectangular(cubeMap [6]image.Image, multAmt float64) image.Image {
	var nrgbaImages [6]*image.NRGBA
	totalWidth := 0
	for i, img := range cubeMap {
		totalWidth += img.Bounds().Dx()
		nrgbaImages[i] = imaging.Clone(img)
	}
	outWidth := int(float64(totalWidth) / 6 * multAmt)
	outHeight := outWidth / 2
	outImg := image.NewNRGBA(image.Rect(0, 0, outWidth, outHeight))

	internal.Parallel(0, outHeight, func(c <-chan int) {
		for j := range c {
			v := 1 - (float64(j) / float64(outHeight))
			theta := v * math.Pi
			for i := 0; i < outWidth; i++ {
				u := (float64(i) / float64(outWidth))
				phi := u * 2 * math.Pi

				x := math.Sin(phi) * math.Sin(theta) * -1
				y := math.Cos(theta)
				z := math.Cos(phi) * math.Sin(theta) * -1

				a := max(math.Abs(x), math.Abs(y), math.Abs(z))

				var sourceImageInd int
				var xPixel, yPixel float64
				shouldFlip := false
				xa := x / a
				ya := y / a
				za := z / a
				if xa == 1 {
					sourceImageInd = 3
					xPixel = ((za + 1) / 2) - 1
					yPixel = (ya + 1) / 2
				} else if xa == -1 {
					sourceImageInd = 1
					xPixel = (za + 1) / 2
					yPixel = (ya + 1) / 2
				} else if ya == 1 {
					shouldFlip = true
					sourceImageInd = 5
					xPixel = (xa + 1) / 2
					yPixel = ((za + 1) / 2) - 1
				} else if ya == -1 {
					shouldFlip = true
					sourceImageInd = 4
					xPixel = (xa + 1) / 2
					yPixel = (za + 1) / 2
				} else if za == 1 {
					sourceImageInd = 2
					xPixel = (xa + 1) / 2
					yPixel = (ya + 1) / 2
				} else if za == -1 {
					sourceImageInd = 0
					xPixel = ((xa + 1) / 2) - 1
					yPixel = (ya + 1) / 2
				} else {
					continue
				}

				img := nrgbaImages[sourceImageInd]
				pixX := min(int(math.Floor(math.Abs(xPixel)*float64(img.Bounds().Dx()))), img.Bounds().Dx()-1)
				pixY := min(int(math.Floor(math.Abs(yPixel)*float64(img.Bounds().Dy()))), img.Bounds().Dy()-1)
				if shouldFlip {
					pixX = img.Bounds().Dx() - pixX - 1
					pixY = img.Bounds().Dy() - pixY - 1
				}
				col := img.NRGBAAt(pixX, pixY)
				outImg.SetNRGBA(i, j, col)
			}
		}
	})
	return outImg
}
