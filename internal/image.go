package internal

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/swim-services/swim_porter/utils"

	"github.com/disintegration/imaging"
	"github.com/gameparrot/fastpng"
	"github.com/gameparrot/tga"
)

func SideOverlayTGA(overlay, base []byte) ([]byte, error) {
	overlayImg, err := fastpng.Decode(bytes.NewReader(overlay))
	if err != nil {
		return []byte{}, err
	}
	baseImg, err := fastpng.Decode(bytes.NewReader(base))
	if err != nil {
		return []byte{}, err
	}
	bounds := baseImg.Bounds()

	canvas := image.NewNRGBA(bounds)
	draw.Draw(canvas, bounds, baseImg, image.Point{}, draw.Src)

	DrawAlpha(canvas, 1)

	scaledOverlay := imaging.Resize(overlayImg, bounds.Dx(), bounds.Dy(), imaging.NearestNeighbor)
	draw.Draw(canvas, bounds, scaledOverlay, image.Point{}, draw.Over)

	writer := bytes.NewBuffer([]byte{})
	if err = tga.Encode(writer, canvas); err != nil {
		return []byte{}, err
	}
	return writer.Bytes(), nil
}

func WritePng(img image.Image, path string, fs *utils.MapFS) error {
	writer := bytes.NewBuffer([]byte{})
	if err := fastpng.Encode(writer, img); err != nil {
		return fmt.Errorf("write image %s: %w", path, err)
	}
	fs.Write(writer.Bytes(), path)
	return nil
}

func DrawAlpha(in *image.NRGBA, alpha uint8) {
	bounds := in.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.X; y < bounds.Max.Y; y++ {
			rgba := in.At(x, y).(color.NRGBA)
			if rgba.A == 0 {
				continue
			}
			in.SetNRGBA(x, y, color.NRGBA{rgba.R, rgba.G, rgba.B, alpha})
		}
	}
}

func DrawAlphaOver(bg *image.NRGBA, in *image.NRGBA, alpha uint8) {
	bounds := in.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.X; y < bounds.Max.Y; y++ {
			rgba := in.At(x, y).(color.NRGBA)
			if rgba.A == 0 {
				continue
			}
			bg.SetNRGBA(x, y, color.NRGBA{rgba.R, rgba.G, rgba.B, alpha})
		}
	}
}

func AlphaMult(in image.Image, mult int) *image.NRGBA {
	bounds := in.Bounds()
	dst := image.NewNRGBA(bounds)
	draw.Draw(dst, bounds, in, image.Point{}, draw.Src)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			imgColor := dst.NRGBAAt(x, y)
			newA := min(255, int(imgColor.A)*mult)
			imgColor.A = uint8(newA)
			dst.Set(x, y, imgColor)
		}
	}
	return dst
}

func LowAlpha(in image.Image, start image.Point) int {
	most := -1
	bounds := in.Bounds()
	for y := start.Y; y < bounds.Max.Y; y++ {
		allNoAlpha := true
		for x := start.X; x < bounds.Max.X; x++ {
			_, _, _, a := in.At(x, y).RGBA()
			if a != 0 {
				allNoAlpha = false
				break
			}
		}
		if y > most {
			most = y
		}
		if allNoAlpha {
			return most
		}
	}
	return most
}

func LowNoAlpha(in image.Image, start image.Point, to image.Point) int {
	bounds := in.Bounds()
	for y := start.Y; y < to.Y; y++ {
		for x := start.X; x < to.X; x++ {
			_, _, _, a := in.At(x, y).RGBA()
			if a != 0 {
				return y
			}
		}
	}
	return bounds.Min.X
}

func RightAlpha(in image.Image, start image.Point, to image.Point) int {
	bounds := in.Bounds()
	for x := to.X - 1; x >= start.X; x-- {
		for y := start.Y; y < to.Y; y++ {
			_, _, _, a := in.At(x, y).RGBA()
			if a != 0 {
				return x
			}
		}
	}
	return bounds.Min.X
}
