package port

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"math"

	"github.com/disintegration/imaging"
	"github.com/swim-services/swim_porter/port/internal"
	"github.com/swim-services/swim_porter/port/recolor"
	"github.com/swim-services/swim_porter/port/utils"
)

func (p *porter) beds() error {
	headTop, err := p.tryLoadBed("bed_head_top")
	if err != nil {
		return nil
	}
	headSide, err := p.tryLoadBed("bed_head_side")
	if err != nil {
		return nil
	}
	headEnd, err := p.tryLoadBed("bed_head_end")
	if err != nil {
		return nil
	}
	feetTop, err := p.tryLoadBed("bed_feet_top")
	if err != nil {
		return nil
	}
	feetSide, err := p.tryLoadBed("bed_feet_side")
	if err != nil {
		return nil
	}
	feetEnd, err := p.tryLoadBed("bed_feet_end")
	if err != nil {
		return nil
	}
	newSize := headTop.Bounds().Dx() * 4
	newImg := image.NewRGBA(image.Rect(0, 0, newSize, newSize))

	baseDraw := int(math.Round(float64(newSize) / 10.6666666667))

	draw.Draw(newImg, headTop.Bounds().Add(image.Point{X: baseDraw, Y: baseDraw}), imaging.Rotate90(headTop), image.Point{}, draw.Src)
	draw.Draw(newImg, feetTop.Bounds().Add(image.Point{X: baseDraw, Y: baseDraw + feetTop.Bounds().Dx()}), imaging.Rotate90(feetTop), image.Point{}, draw.Over)

	feetHeight := int(math.Round(float64(feetSide.Bounds().Dy()) / 5.33333333333))

	feetStart := internal.LowNoAlpha(feetSide, feetSide.Bounds().Min, feetSide.Bounds().Max)

	draw.Draw(newImg, headSide.Bounds().Add(image.Point{X: baseDraw + headSide.Bounds().Dx(), Y: baseDraw}), imaging.Rotate90(headSide), image.Point{X: feetStart}, draw.Over)
	draw.Draw(newImg, headSide.Bounds().Add(image.Point{X: -feetHeight, Y: baseDraw}), imaging.FlipH(imaging.Rotate90(headSide)), image.Point{}, draw.Over)

	draw.Draw(newImg, feetSide.Bounds().Add(image.Point{X: baseDraw + feetSide.Bounds().Dx(), Y: baseDraw + feetTop.Bounds().Dx()}), imaging.Rotate90(feetSide), image.Point{X: feetStart}, draw.Over)
	draw.Draw(newImg, feetSide.Bounds().Add(image.Point{X: -feetHeight, Y: baseDraw + feetTop.Bounds().Dx()}), imaging.FlipH(imaging.Rotate90(feetSide)), image.Point{}, draw.Over)

	draw.Draw(newImg, feetEnd.Bounds().Add(image.Point{X: baseDraw + feetEnd.Bounds().Dx(), Y: -feetHeight}), imaging.Rotate180(feetEnd), image.Point{}, draw.Over)
	draw.Draw(newImg, headEnd.Bounds().Add(image.Point{X: baseDraw, Y: -feetHeight}), imaging.Rotate180(headEnd), image.Point{}, draw.Over)

	if woodPlanksData, err := p.out.Read("textures/blocks/planks_oak.png"); err == nil {
		planksImg, err := png.Decode(bytes.NewReader(woodPlanksData))
		planksImg = imaging.Resize(planksImg, headTop.Bounds().Dx(), headTop.Bounds().Dy(), imaging.NearestNeighbor)
		if err != nil {
			return err
		}
		draw.Draw(newImg, planksImg.Bounds().Add(image.Point{X: baseDraw + feetSide.Bounds().Dx() + feetSide.Bounds().Dx() - feetStart - feetHeight, Y: baseDraw}), planksImg, image.Point{}, draw.Over)
		draw.Draw(newImg, planksImg.Bounds().Add(image.Point{X: baseDraw + feetSide.Bounds().Dx() + feetSide.Bounds().Dx() - feetStart - feetHeight, Y: baseDraw + feetTop.Bounds().Dx()}), planksImg, image.Point{}, draw.Over)
	}

	endFoot0, endFoot1 := feet(feetEnd, feetHeight, 2)
	sideEndFoot, _ := feet(feetSide, feetHeight, 0)
	headFoot0, headFoot1 := feet(headEnd, feetHeight, 2)
	_, sideHeadFoot := feet(headEnd, feetHeight, 1)

	endFoot0Flip := imaging.FlipV(endFoot0)
	endFoot1Flip := imaging.FlipV(endFoot1)

	headFoot0Flip := imaging.FlipV(headFoot0)
	headFoot1Flip := imaging.FlipV(headFoot1)

	sideEndFoot90 := imaging.Rotate90(sideEndFoot)
	sideEndFoot270 := imaging.Rotate270(sideEndFoot)

	sideHeadFoot90 := imaging.Rotate90(sideHeadFoot)
	sideHeadFoot270 := imaging.Rotate270(sideHeadFoot)

	drawFoot(newImg, endFoot0Flip, feetHeight, baseDraw, 2, 2, headTop)
	drawFoot(newImg, imaging.FlipH(endFoot0Flip), feetHeight, baseDraw, 1, 2, headTop)

	drawFoot(newImg, endFoot1Flip, feetHeight, baseDraw, 6, 2, headTop)
	drawFoot(newImg, imaging.FlipH(endFoot1Flip), feetHeight, baseDraw, 5, 2, headTop)

	drawFoot(newImg, sideEndFoot270, feetHeight, baseDraw, 0, 3, headTop)
	drawFoot(newImg, sideEndFoot270, feetHeight, baseDraw, 4, 3, headTop)
	drawFoot(newImg, sideEndFoot90, feetHeight, baseDraw, 2, 3, headTop)
	drawFoot(newImg, sideEndFoot90, feetHeight, baseDraw, 6, 3, headTop)

	drawFoot(newImg, sideHeadFoot270, feetHeight, baseDraw, 0, 1, headTop)
	drawFoot(newImg, sideHeadFoot90, feetHeight, baseDraw, 2, 1, headTop)
	drawFoot(newImg, sideHeadFoot270, feetHeight, baseDraw, 4, 1, headTop)
	drawFoot(newImg, sideHeadFoot90, feetHeight, baseDraw, 6, 1, headTop)

	drawFoot(newImg, imaging.FlipH(headFoot0Flip), feetHeight, baseDraw, 5, 0, headTop)
	drawFoot(newImg, endFoot0Flip, feetHeight, baseDraw, 6, 0, headTop)

	drawFoot(newImg, imaging.FlipH(headFoot1Flip), feetHeight, baseDraw, 1, 0, headTop)
	drawFoot(newImg, endFoot1Flip, feetHeight, baseDraw, 2, 0, headTop)

	drawFoot(newImg, endFoot0, feetHeight, baseDraw, 3, 3, headTop)
	drawFoot(newImg, endFoot1, feetHeight, baseDraw, 7, 3, headTop)

	drawFoot(newImg, headFoot0, feetHeight, baseDraw, 7, 1, headTop)
	drawFoot(newImg, headFoot1, feetHeight, baseDraw, 3, 1, headTop)

	p.recolorBeds(newImg, baseDraw)

	return nil
}

func (p *porter) recolorBeds(redBed image.Image, baseDraw int) error {
	recolor1Start := image.Point{X: redBed.Bounds().Dx() / 32, Y: baseDraw + redBed.Bounds().Dy()/8}
	recolor1End := recolor1Start.Add(image.Point{X: redBed.Bounds().Dx()/4 + redBed.Bounds().Dx()/8, Y: redBed.Bounds().Dx()/4 + redBed.Bounds().Dx()/8})

	recolor2Start := image.Point{X: baseDraw + redBed.Bounds().Dx()/4, Y: redBed.Bounds().Dx() / 32}
	recolor2End := recolor2Start.Add(image.Point{X: redBed.Bounds().Dx() / 4, Y: baseDraw - redBed.Bounds().Dx()/32})

	internal.WritePng(redBed, "textures/entity/bed/red.png", p.out)
	for name, newColor := range utils.BEDS_MAP {
		internal.WritePng(recolor.GrayTintRange(recolor.GrayTintRange(redBed, newColor, recolor1Start, recolor1End, 3), newColor, recolor2Start, recolor2End, 3), "textures/entity/bed/"+name+".png", p.out)
	}
	return nil
}

func (p *porter) tryLoadBed(name string) (image.Image, error) {
	data, err := p.out.Read("textures/blocks/" + name + ".png")
	if err != nil {
		return nil, err
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return img, nil
}

func drawFoot(dst draw.Image, in image.Image, feetHeight int, baseDraw int, xPos int, yPos int, headTop image.Image) {
	draw.Draw(dst, in.Bounds().Add(image.Point{X: feetHeight * xPos, Y: baseDraw + headTop.Bounds().Dx()*2 + feetHeight*yPos}), in, image.Point{}, draw.Over)
}

func feet(img image.Image, feetHeight int, s int) (*image.NRGBA, *image.NRGBA) {
	if s == 0 {
		return imaging.Crop(img, image.Rect(img.Bounds().Min.X, img.Bounds().Max.Y-feetHeight, img.Bounds().Min.X+feetHeight, img.Bounds().Max.Y)), nil
	} else if s == 1 {
		return nil, imaging.Crop(img, image.Rect(img.Bounds().Max.X-feetHeight, img.Bounds().Max.Y-feetHeight, img.Bounds().Max.X, img.Bounds().Max.Y))
	} else {
		return imaging.Crop(img, image.Rect(img.Bounds().Min.X, img.Bounds().Max.Y-feetHeight, img.Bounds().Min.X+feetHeight, img.Bounds().Max.Y)), imaging.Crop(img, image.Rect(img.Bounds().Max.X-feetHeight, img.Bounds().Max.Y-feetHeight, img.Bounds().Max.X, img.Bounds().Max.Y))
	}
}
