package port

import (
	"bytes"
	"fmt"
	"image"
	"image/png"

	"github.com/swedeachu/swim_porter/port/cubemap"
	"github.com/swedeachu/swim_porter/port/internal"
)

func (p *porter) environment() error {
	p.destroyStages()
	p.colorMap()
	if err := p.sky(); err != nil {
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

func (p *porter) sky() error {
	skyboxes := []string{"cloud1", "cloud2", "starfield03", "starfield", "skybox", "skybox2"}
	var skyMap image.Image
	found := false
	for _, box := range skyboxes {
		if data, err := p.in.Read("assets/minecraft/mcpatcher/sky/world0/" + box + ".png"); err == nil {
			if skyMap, err = png.Decode(bytes.NewReader(data)); err == nil {
				found = true
				break
			}
		}
	}
	if !found {
		return nil
	}
	cubemapImages := cubemap.BuildCubemap(skyMap)
	for i, img := range cubemapImages {
		if err := internal.WritePng(img, fmt.Sprintf("textures/environment/overworld_cubemap/cubemap_%d.png", i), p.out); err != nil {
			return err
		}
	}
	return nil
}
