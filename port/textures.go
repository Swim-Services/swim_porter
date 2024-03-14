package port

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"strings"

	"github.com/swim-services/swim_porter/port/internal"
)

func (p *porter) textures() error {
	for name, data := range p.in.Dir("assets/minecraft/textures") {
		name = strings.ReplaceAll(name, "/block/", "/blocks/")
		name = strings.ReplaceAll(name, "/item/", "/items/")
		p.out.Write(data, "textures"+name)
	}
	p.chestFix()
	p.fire()
	p.painting()
	if err := p.grassSide(); err != nil {
		return err
	}
	return nil
}

func (p *porter) chestFix() {
	p.out.Rename("textures/entity/chest/normal_double.png", "textures/entity/chest/double_normal.png")
}

func (p *porter) grassSide() error {
	if dirt, err := p.out.Read("textures/blocks/dirt.png"); err == nil {
		if grassSideOverlay, err := p.out.Read("textures/blocks/grass_side_overlay.png"); err == nil {
			grassSideTGA, err := internal.SideOverlayTGA(grassSideOverlay, dirt)
			if err != nil {
				return err
			}
			p.out.Write(grassSideTGA, "textures/blocks/grass_side.tga")
			p.out.Delete("textures/blocks/grass_side.png")
			snowSideTGA, err := internal.SideOverlayTGA(grassSideOverlay, dirt)
			if err != nil {
				return err
			}
			p.out.Write(snowSideTGA, "textures/blocks/grass_side_snow.tga")
			p.out.Delete("textures/blocks/grass_side_snowed.png")
		}
	}
	return nil
}

func (p *porter) itemsFix() {
	for path, data := range p.out.Dir("textures/items") {
		if strings.HasSuffix(strings.ToLower(path), ".png") {
			if img, err := png.Decode(bytes.NewReader(data)); err == nil {
				internal.WritePng(imageTransparencyFix(img), "textures/items"+path, p.out)
			}
		}
	}
}

func (p *porter) fire() {
	p.out.Rename("textures/blocks/fire_layer_1.png", "textures/blocks/fire_1.png")
	p.out.Rename("textures/blocks/fire_layer_0.png", "textures/blocks/fire_0.png")
}

func (p *porter) painting() {
	p.out.Rename("textures/painting/paintings_kristoffer_zetterstrand.png", "textures/painting/kz.png")
}

func imageTransparencyFix(raw image.Image) *image.RGBA {
	bounds := raw.Bounds()

	image := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := raw.At(x, y)
			_, _, _, alpha := rgba.RGBA()
			alpha >>= 8
			if alpha < 127 {
				rgba = color.NRGBA{0xff, 0xff, 0xff, 0x00} // Fully transparent white
			}

			image.Set(x, y, rgba)
		}
	}
	return image
}
