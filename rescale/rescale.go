package rescale

import (
	"encoding/json"
	"reflect"

	"github.com/disintegration/imaging"
	"github.com/swim-services/swim_porter/porterror"
	"github.com/swim-services/swim_porter/resource"
	"github.com/swim-services/swim_porter/utils"
)

type RescaleOptions struct {
	ShowCredits bool
	Algorithm   imaging.ResampleFilter
}

type rescaler struct {
	in    *utils.MapFS
	scale int
}

func Rescale(in []byte, scale int, opts RescaleOptions) ([]byte, error) {
	zipMap, err := utils.Unzip(in)
	if err != nil {
		return []byte{}, err
	}
	fs := utils.NewMapFS(zipMap)
	err = RescaleRaw(fs, scale, opts)
	if err != nil {
		return []byte{}, err
	}
	out, err := utils.Zip(fs.RawMap())
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}

func RescaleRaw(in *utils.MapFS, scale int, opts RescaleOptions) error {
	if reflect.DeepEqual(opts.Algorithm, imaging.ResampleFilter{}) {
		opts.Algorithm = imaging.NearestNeighbor
	}
	p := &rescaler{in: in, scale: scale}
	err := p.doRescale(opts)
	if err != nil {
		return porterror.Wrap(err)
	}
	return nil
}

func (p *rescaler) doRescale(opts RescaleOptions) error {
	if err := p.manifest(opts.ShowCredits); err != nil {
		return porterror.Wrap(err)
	}
	if err := p.rescaleDir("textures/blocks", opts.Algorithm); err != nil {
		return porterror.Wrap(err)
	}
	if err := p.rescaleDir("textures/items", opts.Algorithm); err != nil {
		return porterror.Wrap(err)
	}
	if err := p.rescaleDir("textures/entity", opts.Algorithm); err != nil {
		return porterror.Wrap(err)
	}
	return nil
}

func (p *rescaler) manifest(showCredits bool) error {
	bedrockManifestOrig, err := p.in.Read("manifest.json")
	if err != nil {
		return porterror.Wrap(porterror.ErrManifestNotFound)
	}
	bedrockManifest, err := resource.UnmarshalJSON(bedrockManifestOrig)
	if err != nil {
		return porterror.Wrap(err)
	}
	utils.ChangeUUID(&bedrockManifest)
	if showCredits {
		bedrockManifest.Header.Description += "\n§aRescaled by §dSwim Auto Rescale §f| §bdiscord.gg/swim"
	}
	bedrockManifestBytes, err := json.Marshal(bedrockManifest)
	if err != nil {
		return porterror.Wrap(err)
	}
	p.in.Write(bedrockManifestBytes, "manifest.json")
	return nil
}
