package recolor

import (
	"bytes"
	"encoding/json"
	"errors"
	"image/color"
	"image/png"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/sandertv/gophertunnel/minecraft/resource"
	"github.com/swedeachu/swim_porter/port/internal"
	"github.com/swedeachu/swim_porter/port/utils"
)

type RecolorOptions struct {
	ShowCredits bool
	NewColor    string
}

type recolorer struct {
	in    *utils.MapFS
	color color.RGBA
}

func Recolor(in []byte, opts RecolorOptions) ([]byte, error) {
	zipMap, err := utils.Unzip(in)
	if err != nil {
		return []byte{}, err
	}
	fs := utils.NewMapFS(zipMap)
	err = RecolorRaw(utils.NewMapFS(zipMap), opts)
	if err != nil {
		return []byte{}, err
	}
	out, err := utils.Zip(fs.RawMap())
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}

func RecolorRaw(in *utils.MapFS, opts RecolorOptions) error {
	p := &recolorer{in: in}
	err := p.doRecolor(opts)
	if err != nil {
		return err
	}
	return nil
}

func (p *recolorer) doRecolor(opts RecolorOptions) error {
	if err := p.manifest(opts.ShowCredits); err != nil {
		return err
	}
	color, err := utils.ParseHex(opts.NewColor)
	if err != nil {
		return err
	}
	for file, data := range p.in.RawMap() {
		name := filepath.Base(file)
		ext := path.Ext(name)
		if ext != ".png" {
			continue
		}
		nameNoExt := name[:strings.LastIndex(name, ext)]
		if !slices.Contains(utils.DEFAULT_RECOLOR_LIST, nameNoExt) {
			continue
		}
		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			continue // ignore invalid images
		}
		newImg, err := HueShift(img, float64(GetHue(int(color.R), int(color.G), int(color.B))))
		if err != nil {
			return err
		}
		if err := internal.WritePng(newImg, file, p.in); err != nil {
			return err
		}
	}
	return nil
}

func (p *recolorer) manifest(showCredits bool) error {
	bedrockManifestOrig, err := p.in.Read("manifest.json")
	if err != nil {
		return errors.New("manifest.json not found")
	}
	var bedrockManifest resource.Manifest
	err = json.Unmarshal(bedrockManifestOrig, &bedrockManifest)
	if err != nil {
		return err
	}
	bedrockManifest.Header.Name += "§r§b Recolor"
	if showCredits {
		bedrockManifest.Header.Description += "\n§aRecolored by §dSwim Auto Recolor §f| §bdiscord.gg/swim"
	}
	bedrockManifestBytes, err := json.Marshal(bedrockManifest)
	if err != nil {
		return err
	}
	p.in.Write(bedrockManifestBytes, "manifest.json")
	return nil
}
