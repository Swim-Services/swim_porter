package recolor

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"maps"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gameparrot/fastpng"
	"github.com/swim-services/swim_porter/internal"
	"github.com/swim-services/swim_porter/porterror"
	"github.com/swim-services/swim_porter/resource"
	"github.com/swim-services/swim_porter/utils"
)

type RecolorOptions struct {
	ShowCredits bool
	NewColor    color.RGBA
	Alg         Algorithm
}

type recolorer struct {
	in *utils.MapFS
}

func Recolor(in []byte, opts RecolorOptions) ([]byte, error) {
	zipMap, err := utils.Unzip(in)
	if err != nil {
		return []byte{}, err
	}
	fs := utils.NewMapFS(zipMap)
	err = RecolorRaw(fs, opts)
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
	if opts.Alg == nil {
		opts.Alg = &HueV2{}
	}
	p := &recolorer{in: in}
	err := p.doRecolor(opts)
	if err != nil {
		return porterror.Wrap(err)
	}
	return nil
}

func (p *recolorer) doRecolor(opts RecolorOptions) error {
	if err := p.manifest(opts.ShowCredits); err != nil {
		return porterror.Wrap(err)
	}
	opts.Alg.SetColor(opts.NewColor)
	list := opts.Alg.DefaultList()
	internal.ParallelMap(maps.Clone(p.in.RawMap()), func(file string, data []byte) {
		name := filepath.Base(file)
		ext := path.Ext(name)
		if ext != ".png" {
			return
		}
		nameNoExt := name[:strings.LastIndex(name, ext)]
		if !slices.Contains(list, nameNoExt) {
			return
		}
		img, err := fastpng.Decode(bytes.NewReader(data))
		if err != nil {
			return // ignore invalid images
		}
		var newImg image.Image
		err = nil
		if newImg, err = opts.Alg.RecolorImage(img, file); err != nil {
			return
		}
		if err := internal.WritePng(newImg, file, p.in); err != nil {
			return
		}
	})
	return nil
}

func (p *recolorer) manifest(showCredits bool) error {
	bedrockManifestOrig, err := p.in.Read("manifest.json")
	if err != nil {
		return porterror.Wrap(porterror.ErrManifestNotFound)
	}
	bedrockManifest, err := resource.UnmarshalJSON(bedrockManifestOrig)
	if err != nil {
		return porterror.Wrap(err)
	}
	utils.ChangeUUID(&bedrockManifest)
	bedrockManifest.Header.Name += "§r§b Recolor"
	if showCredits {
		bedrockManifest.Header.Description += "\n§aRecolored by §dSwim Auto Recolor §f| §bdiscord.gg/swim"
	}
	bedrockManifestBytes, err := json.Marshal(bedrockManifest)
	if err != nil {
		return porterror.Wrap(err)
	}
	p.in.Write(bedrockManifestBytes, "manifest.json")
	return nil
}
