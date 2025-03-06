package skyfix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"

	"github.com/gameparrot/fastpng"
	"github.com/swim-services/swim_porter/cubemap"
	"github.com/swim-services/swim_porter/internal"
	"github.com/swim-services/swim_porter/porterror"
	"github.com/swim-services/swim_porter/resource"
	"github.com/swim-services/swim_porter/utils"
	"golang.org/x/sync/errgroup"
)

type skyfixer struct {
	in *utils.MapFS
}

func FixSky(in []byte) ([]byte, error) {
	zipMap, err := utils.Unzip(in)
	if err != nil {
		return []byte{}, err
	}
	fs := utils.NewMapFS(zipMap)
	err = FixSkyRaw(utils.NewMapFS(zipMap))
	if err != nil {
		return []byte{}, err
	}
	out, err := utils.Zip(fs.RawMap())
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}

func FixSkyRaw(in *utils.MapFS) error {
	p := &skyfixer{in: in}
	err := p.doSkyFix()
	if err != nil {
		return err
	}
	return nil
}

func (p *skyfixer) doSkyFix() error {
	if err := p.manifest(); err != nil {
		return err
	}

	var cubeMap [6]image.Image

	for i := 0; i < 6; i++ {
		file, err := p.in.Read(fmt.Sprintf("textures/environment/overworld_cubemap/cubemap_%d.png", i))
		if err != nil {
			return porterror.New("pack does not contain sky")
		}
		pngImg, err := fastpng.Decode(bytes.NewReader(file))
		if err != nil {
			return porterror.Wrap(err).WithMessage("read cubemap_%d.png")
		}
		cubeMap[i] = pngImg
	}

	totalWidth := 0
	for _, img := range cubeMap {
		totalWidth += img.Bounds().Dx()
	}
	multAmt := max(4.5, min(8, float64(totalWidth)/1024))
	equi := cubemap.CubemapToEquirectangular(cubeMap, multAmt)
	cubeMap = cubemap.CubemapFromImage(equi, cubemap.CubemapImageOpts{VertOffset: 0.41, DivAmt: multAmt})

	errs, _ := errgroup.WithContext(context.Background())
	for i, img := range cubeMap {
		ii := i
		errs.Go(func() error {
			if err := internal.WritePng(img, fmt.Sprintf("textures/environment/overworld_cubemap/cubemap_%d.png", ii), p.in); err != nil {
				return porterror.Wrap(err)
			}
			return nil
		})
	}
	return errs.Wait()
}

func (p *skyfixer) manifest() error {
	bedrockManifestOrig, err := p.in.Read("manifest.json")
	if err != nil {
		return porterror.Wrap(porterror.ErrManifestNotFound)
	}
	bedrockManifest, err := resource.UnmarshalJSON(bedrockManifestOrig)
	if err != nil {
		return err
	}
	utils.ChangeUUID(&bedrockManifest)
	bedrockManifest.Header.Description += "\n§aSky fixed by §dSwim Sky Fixer §f| §bdiscord.gg/swim"
	bedrockManifestBytes, err := json.Marshal(bedrockManifest)
	if err != nil {
		return err
	}
	p.in.Write(bedrockManifestBytes, "manifest.json")
	return nil
}
