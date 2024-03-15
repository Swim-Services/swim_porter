package rescale

import (
	"bytes"
	"image"
	"image/png"
	"path"
	"path/filepath"

	"github.com/blezek/tga"
	"github.com/disintegration/imaging"
)

var filters = map[string]imaging.ResampleFilter{"nearest_neighbor": imaging.NearestNeighbor, "box": imaging.Box, "linear": imaging.Linear, "hermite": imaging.Hermite, "mitchellnetravali": imaging.MitchellNetravali, "catmull_rom": imaging.CatmullRom, "bspline": imaging.BSpline, "gaussian": imaging.Gaussian, "bartlett": imaging.Bartlett, "lanczos": imaging.Lanczos, "hann": imaging.Hann, "hamming": imaging.Hamming, "blackman": imaging.Blackman, "welch": imaging.Welch, "cosine": imaging.Cosine}

func (p *rescaler) rescaleDir(dir string, filter imaging.ResampleFilter) error {
	for file, data := range p.in.Dir(dir) {
		name := filepath.Base(file)
		ext := path.Ext(name)

		var img image.Image
		var err error
		switch ext {
		case ".png":
			img, err = png.Decode(bytes.NewReader(data))
		case ".tga":
			img, err = tga.Decode(bytes.NewReader(data))
		default:
			continue
		}
		if err != nil {
			continue // ignore invalid images
		}
		if img.Bounds().Dx() > 128 {
			continue
		}
		newImg := imaging.Resize(img, p.scale, p.scale, filter)
		writer := bytes.NewBuffer([]byte{})

		switch ext {
		case ".png":
			err = png.Encode(writer, newImg)
		case ".tga":
			err = tga.Encode(writer, newImg)
		}
		if err != nil {
			return err
		}
		p.in.Write(writer.Bytes(), dir+file)
	}
	return nil

}

func ParseAlgorithm(alg string) (imaging.ResampleFilter, bool) {
	parsedFilter, ok := filters[alg]
	return parsedFilter, ok
}

func GetAlgorithms() []string {
	filterStr := make([]string, len(filters))
	i := 0
	for k := range filters {
		filterStr[i] = k
		i++
	}
	return filterStr
}
