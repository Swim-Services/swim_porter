package port

import (
	"bytes"
	"image"
	"image/png"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/swim-services/swim_porter/internal"
	"github.com/swim-services/swim_porter/particlefix"
	"github.com/swim-services/swim_porter/porterror"
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
	particlefix.DoFix(p.out)
	if err := p.grassSide(); err != nil {
		return err
	}
	if err := p.water(false); err != nil {
		return err
	}
	if err := p.water(true); err != nil {
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
				return porterror.Wrap(err)
			}
			p.out.Write(grassSideTGA, "textures/blocks/grass_side.tga")
			p.out.Delete("textures/blocks/grass_side.png")
			snowSideTGA, err := internal.SideOverlayTGA(grassSideOverlay, dirt)
			if err != nil {
				return porterror.Wrap(err)
			}
			p.out.Write(snowSideTGA, "textures/blocks/grass_side_snow.tga")
			p.out.Delete("textures/blocks/grass_side_snowed.png")
		}
	}
	return nil
}

func (p *porter) itemsFix() {
	internal.ParallelMap(p.out.Dir("textures/items"), func(path string, data []byte) {
		if strings.HasSuffix(strings.ToLower(path), ".png") {
			if img, err := png.Decode(bytes.NewReader(data)); err == nil {
				internal.WritePng(imageTransparencyFix(img, 127), "textures/items"+path, p.out)
			}
		}
	})
}

func (p *porter) fire() {
	p.out.Rename("textures/blocks/fire_layer_1.png", "textures/blocks/fire_1.png")
	p.out.Rename("textures/blocks/fire_layer_0.png", "textures/blocks/fire_0.png")
}

func (p *porter) water(flow bool) error {
	var waterType string
	if flow {
		waterType = "_flow"
	} else {
		waterType = "_still"
	}
	if data, err := p.out.Read("textures/blocks/water" + waterType + ".png"); err == nil {
		waterImg, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			return porterror.Wrap(err).WithMessage("read image textures/blocks/water%s.png", waterType)
		}
		greyWater := internal.AlphaMult(imaging.Grayscale(waterImg), 2)
		if err := internal.WritePng(greyWater, "textures/blocks/water"+waterType+"_grey.png", p.out); err != nil {
			return porterror.Wrap(err)
		}
	}
	return nil
}

func (p *porter) painting() {
	p.out.Rename("textures/painting/paintings_kristoffer_zetterstrand.png", "textures/painting/kz.png")
}

func imageTransparencyFix(raw image.Image, cutoff uint8) *image.NRGBA {
	newImg := imaging.Clone(raw)
	for i := 3; i < len(newImg.Pix); i += 4 {
		if newImg.Pix[i] <= cutoff {
			newImg.Pix[i-3] = 255
			newImg.Pix[i-2] = 255
			newImg.Pix[i-1] = 255
			newImg.Pix[i] = 0
		}
	}
	return newImg
}
