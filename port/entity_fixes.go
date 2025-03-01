package port

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"github.com/disintegration/imaging"
	"github.com/gameparrot/tga"
	"github.com/swim-services/swim_porter/internal"
	"github.com/swim-services/swim_porter/porterror"
	"github.com/swim-services/swim_porter/recolor"
)

func (p *porter) entityFixes() error {
	if err := p.fixZombie(); err != nil {
		return err
	}
	if err := p.fixSheep(); err != nil {
		return err
	}
	if err := p.fixLeather(); err != nil {
		return err
	}
	return nil
}

func (p *porter) fixZombie() error {
	if data, err := p.out.Read("textures/entity/zombie/zombie.png"); err == nil {
		zombieImg, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			return porterror.Wrap(err).WithMessage("read image textures/entity/zombie/zombie.png")
		}
		bounds := zombieImg.Bounds()
		newZombie := imaging.Crop(zombieImg, image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Min.Y+(bounds.Dy()/2)))
		if err := internal.WritePng(newZombie, "textures/entity/zombie/zombie.png", p.out); err != nil {
			return porterror.Wrap(err)
		}
	}
	return nil
}

func (p *porter) fixSheep() error {
	if data, err := p.out.Read("textures/entity/sheep/sheep.png"); err == nil {
		if datafur, err := p.out.Read("textures/entity/sheep/sheep_fur.png"); err == nil {
			sheepImg, err := png.Decode(bytes.NewReader(data))
			if err != nil {
				return porterror.Wrap(err).WithMessage("read image textures/entity/sheep/sheep.png")
			}
			sheepFurImg, err := png.Decode(bytes.NewReader(datafur))
			if err != nil {
				return porterror.Wrap(err).WithMessage("read image textures/entity/sheep/sheep_fur.png")
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
				return porterror.Wrap(err)
			}
			p.out.Write(writer.Bytes(), "textures/entity/sheep/sheep.tga")
		}
	}
	return nil
}

func (p *porter) fixLeather() error {
	for i := 1; i <= 2; i++ {
		imgPath := fmt.Sprintf("textures/models/armor/cloth_%d.png", i)
		overlayPath := fmt.Sprintf("textures/models/armor/leather_layer_%d_overlay.png", i)
		if data, err := p.out.Read(imgPath); err == nil {
			img, err := png.Decode(bytes.NewReader(data))
			if err != nil {
				return porterror.Wrap(err).WithMessage("read image textures/models/armor/cloth_%d.png", i)
			}
			newImg := recolor.Tint(img, color.RGBA{R: 190, G: 120, B: 80})
			if err := internal.WritePng(newImg, imgPath, p.out); err != nil {
				return porterror.Wrap(err)
			}

			if overlayData, err := p.out.Read(overlayPath); err == nil {
				overlay, err := png.Decode(bytes.NewReader(overlayData))
				if err != nil {
					return porterror.Wrap(err).WithMessage("read image textures/models/armor/leather_layer_%d_overlay.png", i)
				}
				drawImg := imaging.Clone(img)
				if overlay.Bounds().Dx() != drawImg.Bounds().Dx() {
					overlay = imaging.Resize(overlay, drawImg.Bounds().Dx(), drawImg.Bounds().Dy(), imaging.NearestNeighbor)
				}
				internal.DrawAlphaOver(drawImg, imaging.Clone(overlay), 1)
				writer := bytes.NewBuffer([]byte{})
				if err := tga.Encode(writer, imageTransparencyFix(drawImg, 0)); err != nil {
					return porterror.Wrap(err)
				}
				p.out.Write(writer.Bytes(), fmt.Sprintf("textures/models/armor/leather_%d.tga", i))
				p.out.Delete(overlayPath)
			}
		}
	}

	return nil
}
