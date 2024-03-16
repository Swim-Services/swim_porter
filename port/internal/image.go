package internal

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"github.com/swim-services/swim_porter/port/utils"

	"github.com/disintegration/imaging"
	"github.com/gameparrot/tga"
)

func SideOverlayTGA(overlay, base []byte) ([]byte, error) {
	overlayImg, err := png.Decode(bytes.NewReader(overlay))
	if err != nil {
		return []byte{}, err
	}
	baseImg, err := png.Decode(bytes.NewReader(base))
	if err != nil {
		return []byte{}, err
	}
	bounds := baseImg.Bounds()

	canvas := image.NewNRGBA(bounds)
	draw.Draw(canvas, bounds, baseImg, image.Point{}, draw.Src)

	drawAlpha(canvas, 1)

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
	if err := png.Encode(writer, img); err != nil {
		return err
	}
	fs.Write(writer.Bytes(), path)
	return nil
}

func drawAlpha(in *image.NRGBA, alpha uint8) {
	bounds := in.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.X; y < bounds.Max.Y; y++ {
			rgba := in.At(x, y).(color.NRGBA)
			in.SetNRGBA(x, y, color.NRGBA{rgba.R, rgba.G, rgba.B, alpha})
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
