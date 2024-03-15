package port

import (
	"bytes"
	"image"
	"image/png"

	"github.com/disintegration/imaging"
	"github.com/swim-services/swim_porter/port/internal"
)

func (p *porter) entityFixes() error {
	if err := p.fixZombie(); err != nil {
		return err
	}
	return nil
}

func (p *porter) fixZombie() error {
	if data, err := p.out.Read("textures/entity/zombie/zombie.png"); err == nil {
		zombieImg, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			return err
		}
		bounds := zombieImg.Bounds()
		newZombie := imaging.Crop(zombieImg, image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Min.Y+(bounds.Dy()/2)))
		if err := internal.WritePng(newZombie, "textures/entity/zombie/zombie.png", p.out); err != nil {
			return err
		}
	}
	return nil
}
