package fontfix

import (
	"bytes"
	"encoding/json"
	"errors"
	"image/color"
	"image/png"

	"github.com/disintegration/imaging"
	"github.com/swim-services/swim_porter/internal"
	"github.com/swim-services/swim_porter/porterror"
	"github.com/swim-services/swim_porter/resource"
	"github.com/swim-services/swim_porter/utils"
)

type fontfixer struct {
	in *utils.MapFS
}

func FixFont(in []byte) ([]byte, error) {
	zipMap, err := utils.Unzip(in)
	if err != nil {
		return []byte{}, err
	}
	fs := utils.NewMapFS(zipMap)
	err = FixFontRaw(utils.NewMapFS(zipMap))
	if err != nil {
		return []byte{}, err
	}
	out, err := utils.Zip(fs.RawMap())
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}

func FixFontRaw(in *utils.MapFS) error {
	p := &fontfixer{in: in}
	err := p.doFontFix()
	if err != nil {
		return err
	}
	return nil
}

func (p *fontfixer) doFontFix() error {
	if err := p.manifest(); err != nil {
		return err
	}
	if data, err := p.in.Read("font/default8.png"); err == nil {
		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			return porterror.Wrap(err).WithMessage("read image font/default8.png")
		}
		bounds := img.Bounds()
		if bounds.Dx() == 128 || bounds.Dx() == 1152 {
			return errors.New("this pack's font does not need to be fixed")
		}
		newImg := imaging.Resize(img, 1152, 1152, imaging.NearestNeighbor)
		newBounds := newImg.Bounds()
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y += 72 {
			for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
				if y == 0 {
					continue
				}
				newImg.Set(x, y-1, color.Transparent)
			}
		}
		if err := internal.WritePng(newImg, "font/default8.png", p.in); err != nil {
			return err
		}
	} else {
		return errors.New("default8 font not found; this pack does not contain a font")
	}
	return nil
}

func (p *fontfixer) manifest() error {
	bedrockManifestOrig, err := p.in.Read("manifest.json")
	if err != nil {
		return porterror.Wrap(porterror.ErrManifestNotFound)
	}
	bedrockManifest, err := resource.UnmarshalJSON(bedrockManifestOrig)
	if err != nil {
		return err
	}
	utils.ChangeUUID(&bedrockManifest)
	bedrockManifest.Header.Description += "\n§aFont fixed by §dSwim Font Fixer §f| §bdiscord.gg/swim"
	bedrockManifestBytes, err := json.Marshal(bedrockManifest)
	if err != nil {
		return err
	}
	p.in.Write(bedrockManifestBytes, "manifest.json")
	return nil
}
