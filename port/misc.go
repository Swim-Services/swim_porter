package port

import (
	"bytes"
	_ "embed"
	"image"
	"image/draw"
	"image/png"

	"github.com/swim-services/swim_porter/port/internal"
	"github.com/swim-services/swim_porter/port/utils"
)

//go:embed assets.zip
var assets []byte

var assetsMapFS *utils.MapFS

func init() {
	var err error
	assetsMap, err := utils.Unzip(assets)
	if err != nil {
		panic(err)
	}
	assetsMapFS = utils.NewMapFS(assetsMap)
}

func (p *porter) misc() error {
	if err := p.title(); err != nil {
		return err
	}
	p.panorama()
	p.sounds()
	p.armor()
	return nil
}

func (p *porter) panorama() {
	p.out.InterCopyDir(p.in, "assets/minecraft/textures/gui/title/background", "textures/ui")
}

func (p *porter) sounds() {
	p.out.InterCopyDir(p.in, "assets/minecraft/sounds", "sounds")
}

func (p *porter) armor() {
	for old, new := range utils.ARMOR_MAP {
		p.out.Rename("textures/models/armor/"+old+".png", "textures/models/armor/"+new+".png")
	}
}

func (p *porter) title() error {
	if data, err := p.out.Read("textures/gui/title/minecraft.png"); err == nil {
		titleImg, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			return err
		}
		start := internal.LowAlpha(titleImg, titleImg.Bounds().Min)
		start2 := internal.LowNoAlpha(titleImg, image.Point{titleImg.Bounds().Min.X, start}, titleImg.Bounds().Max)
		left := internal.RightAlpha(titleImg, titleImg.Bounds().Min, image.Point{titleImg.Bounds().Max.X, start})
		left2 := internal.RightAlpha(titleImg, image.Point{X: titleImg.Bounds().Min.X, Y: start}, image.Point{titleImg.Bounds().Max.X, start * 2})
		newImg := image.NewRGBA(image.Rect(0, 0, left+left2+1, start))
		draw.Draw(newImg, newImg.Bounds(), titleImg, image.Point{}, draw.Src)
		draw.Draw(newImg, newImg.Bounds().Add(image.Point{X: left}), titleImg, image.Point{Y: start2}, draw.Src)
		if err := internal.WritePng(newImg, "textures/ui/title.png", p.out); err != nil {
			return err
		}
	}
	return nil
}
