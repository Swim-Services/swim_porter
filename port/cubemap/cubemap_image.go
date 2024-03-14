package cubemap

import (
	"image"
	"image/draw"
	"math"
)

var cubemapNumber = [6]int{0, 2, 1, 3, 4, 5}

func CubemapFromImage(in image.Image) [6]image.Image {
	img := image.NewRGBA(in.Bounds())
	draw.Draw(img, img.Bounds(), in, image.Point{}, draw.Src)
	var imgs [6]image.Image
	for i := 0; i < 6; i++ {
		newImg := renderFace(img, i, math.Pi, i, 2048)
		imgs[cubemapNumber[i]] = newImg
	}
	return imgs
}

type cube struct {
	x, y, z float64
}

var orientations = []func(out *cube, x, y float64){
	func(out *cube, x, y float64) {
		out.x = -1
		out.y = -x
		out.z = -y
	},
	func(out *cube, x, y float64) {
		out.x = 1
		out.y = x
		out.z = -y
	},
	func(out *cube, x, y float64) {
		out.x = x
		out.y = -1
		out.z = -y
	},
	func(out *cube, x, y float64) {
		out.x = -x
		out.y = 1
		out.z = -y
	},
	func(out *cube, x, y float64) {
		out.x = -y
		out.y = -x
		out.z = 1
	},
	func(out *cube, x, y float64) {
		out.x = y
		out.y = -x
		out.z = -1
	},
}

func clamp(x, min, max float64) float64 {
	return math.Min(max, math.Max(x, min))
}

func mod(x, n float64) float64 {
	return math.Mod((math.Mod(x, n))+n, n)
}

func copyPixelNearest(read, write *image.RGBA) func(float64, float64, int) {
	width, height := read.Bounds().Dx(), read.Bounds().Dy()
	readIndex := func(x, y int) int {
		return 4 * (y*width + x)
	}

	return func(xFrom, yFrom float64, to int) {
		nearest := readIndex(
			int(clamp(math.Round(xFrom), 0, float64(width-1))),
			int(clamp(math.Round(yFrom), 0, float64(height-1))),
		)

		for channel := 0; channel < 3; channel++ {
			write.Pix[to+channel] = read.Pix[nearest+channel]
		}
	}
}

func renderFace(readData *image.RGBA, face int, rotation float64, num int, maxWidth float64) *image.RGBA {
	num += 1
	faceWidth := math.Min(maxWidth, float64(readData.Bounds().Dx())/4)
	faceHeight := faceWidth

	cube := &cube{}
	orientation := orientations[face]

	writeData := image.NewRGBA(image.Rect(0, 0, int(faceWidth), int(faceHeight)))

	copyPixel := copyPixelNearest(readData, writeData)

	for x := 0; x < int(faceWidth); x++ {
		for y := 0; y < int(faceHeight); y++ {
			to := 4 * (y*int(faceWidth) + x)

			writeData.Pix[to+3] = 255
			orientation(cube, (2*(float64(x)+0.5)/faceWidth - 1), (2*(float64(y)+0.5)/faceHeight - 1))

			r := math.Sqrt(cube.x*cube.x + cube.y*cube.y + cube.z*cube.z)
			lon := mod(math.Atan2(cube.y, cube.x)+rotation, 2*math.Pi)
			lat := math.Acos(cube.z / r)

			copyPixel(float64(readData.Bounds().Dx())*lon/math.Pi/2-0.5, float64(readData.Bounds().Dy())*lat/math.Pi-0.5, to)
		}
	}

	return (writeData)
}
