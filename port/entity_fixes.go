package port

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"

	"github.com/disintegration/imaging"
	"github.com/gameparrot/tga"
	"github.com/swim-services/swim_porter/port/internal"
)

func (p *porter) entityFixes() error {
	if err := p.fixZombie(); err != nil {
		return err
	}
	if err := p.fixSheep(); err != nil {
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

func (p *porter) fixSheep() error {
	if data, err := p.out.Read("textures/entity/sheep/sheep.png"); err == nil {
		if datafur, err := p.out.Read("textures/entity/sheep/sheep_fur.png"); err == nil {
			sheepImg, err := png.Decode(bytes.NewReader(data))
			if err != nil {
				return err
			}
			sheepFurImg, err := png.Decode(bytes.NewReader(datafur))
			if err != nil {
				return err
			}
			bounds := sheepImg.Bounds()
			newImg := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dx()))
			if sheepFurImg.Bounds().Dx() != bounds.Dx() {
				sheepFurImg = imaging.Resize(sheepFurImg, bounds.Dx(), bounds.Dy(), imaging.NearestNeighbor)
			}
			draw.Draw(newImg, bounds, sheepImg, image.Point{}, draw.Src)
			draw.Draw(newImg, bounds.Add(image.Point{Y: bounds.Dy()}), sheepFurImg, image.Point{}, draw.Src)
			newBounds := newImg.Bounds()
			for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
				for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
					curColor := newImg.RGBAAt(x, y)
					max := max(curColor.R, curColor.G, curColor.B)
					min := min(curColor.R, curColor.G, curColor.B)
					if int(curColor.R)+int(curColor.G)+int(curColor.B) < 700 && max-min > 20 {
						curColor.A = 1
						newImg.Set(x, y, curColor)
					}
				}
			}
			p.out.Delete("textures/entity/sheep/sheep.png")
			p.out.Delete("textures/entity/sheep/sheep_fur.png")
			writer := bytes.NewBuffer([]byte{})
			if err := tga.Encode(writer, newImg); err != nil {
				return err
			}
			p.out.Write(writer.Bytes(), "textures/entity/sheep/sheep.tga")
		}
	}
	return nil
}
