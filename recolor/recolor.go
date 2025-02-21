package recolor

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/png"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/swim-services/swim_porter/internal"
	"github.com/swim-services/swim_porter/porterror"
	"github.com/swim-services/swim_porter/resource"
	"github.com/swim-services/swim_porter/utils"
	stripjsoncomments "github.com/trapcodeio/go-strip-json-comments"
)

type RecolorOptions struct {
	ShowCredits bool
	NewColor    color.RGBA
	Alg         string
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
	if opts.Alg == "" {
		opts.Alg = "hue"
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
		var newImg image.Image
		err = nil
		switch opts.Alg {
		case "tint":
			newImg = Tint(img, opts.NewColor)
		case "hue":
			newImg, err = HueShift(img, float64(GetHue(int(opts.NewColor.R), int(opts.NewColor.G), int(opts.NewColor.B))))
		case "gray_tint":
			newImg = GrayTint(img, opts.NewColor)
		default:
			return errors.New("unknown algorithm")
		}
		if err != nil {
			return porterror.Wrap(err)
		}
		if err := internal.WritePng(newImg, file, p.in); err != nil {
			return porterror.Wrap(err)
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
	err = json.Unmarshal([]byte(stripjsoncomments.Strip(string(bedrockManifestOrig))), &bedrockManifest)
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
