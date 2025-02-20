package port

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"

	"github.com/swim-services/swim_porter/port/cubemap"
	"github.com/swim-services/swim_porter/port/internal"
	"github.com/swim-services/swim_porter/port/porterror"
	"golang.org/x/sync/errgroup"
)

func (p *porter) environment(skyboxOverride string) error {
	p.destroyStages()
	p.colorMap()
	if err := p.sky(skyboxOverride); err != nil {
		return err
	}
	return nil
}

func (p *porter) destroyStages() {
	for i := 0; i < 10; i++ {
		if data, err := p.in.Read(fmt.Sprintf("assets/minecraft/textures/blocks/destroy_stage_%d.png", i)); err == nil {
			p.out.Write(data, fmt.Sprintf("textures/environment/destroy_stage_%d.png", i))
		}
	}
}

func (p *porter) colorMap() {
	p.out.InterCopyDir(p.in, "assets/minecraft/mcpatcher/colormap", "textures/colormap")
}

func (p *porter) sky(skyboxOverride string) error {
	skyboxes := []string{"cloud1", "cloud2", "starfield03", "starfield", "skybox", "skybox2"}
	var skyMap image.Image
	found := false
	if data, err := p.in.Read("assets/minecraft/mcpatcher/sky/world0/" + skyboxOverride + ".png"); err == nil && skyboxOverride != "" {
		if skyMap, err = png.Decode(bytes.NewReader(data)); err == nil {
			found = true
		}
	} else {
		for _, box := range skyboxes {
			if data, err := p.in.Read("assets/minecraft/mcpatcher/sky/world0/" + box + ".png"); err == nil {
				if skyMap, err = png.Decode(bytes.NewReader(data)); err == nil {
					found = true
					break
				}
			}
		}
	}
	if !found {
		return nil
	}

	cubemapImages := cubemap.BuildCubemap(skyMap)

	if p.opts.OffsetSky {
		totalWidth := 0
		for _, img := range cubemapImages {
			totalWidth += img.Bounds().Dx()
		}
		multAmt := max(4.5, min(8, float64(totalWidth)/1024))
		equi := cubemap.CubemapToEquirectangular(cubemapImages, multAmt)
		cubemapImages = cubemap.CubemapFromImage(equi, cubemap.CubemapImageOpts{VertOffset: 0.41, DivAmt: multAmt})
	}

	errs, _ := errgroup.WithContext(context.Background())
	for i, img := range cubemapImages {
		ii := i
		errs.Go(func() error {
			if err := internal.WritePng(img, fmt.Sprintf("textures/environment/overworld_cubemap/cubemap_%d.png", ii), p.out); err != nil {
				return porterror.Wrap(err)
			}
			return nil
		})
	}
	return errs.Wait()
}
