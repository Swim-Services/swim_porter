package port

import (
	"bytes"
	_ "embed"
	"image"
	"image/draw"
	"image/png"

	"github.com/disintegration/imaging"
	"github.com/swim-services/swim_porter/port/internal"
	"github.com/swim-services/swim_porter/port/utils"
)

const particleGridSize = 16
const vanillaParticleSize = 8

//go:embed assets.zip
var assets []byte

var assetsMapFS *utils.MapFS

var vanillaParticles image.Image

func init() {
	var err error
	assetsMap, err := utils.Unzip(assets)
	if err != nil {
		panic(err)
	}
	assetsMapFS = utils.NewMapFS(assetsMap)

	vanillaParticleData, err := assetsMapFS.Read("particle/particles.png")
	if err != nil {
		panic(err)
	}
	vanillaParticles, err = png.Decode(bytes.NewBuffer(vanillaParticleData))
	if err != nil {
		panic(err)
	}
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
		starta := internal.LowNoAlpha(titleImg, titleImg.Bounds().Min, titleImg.Bounds().Max)
		start := internal.LowAlpha(titleImg, image.Point{titleImg.Bounds().Min.X, starta})
		start2 := internal.LowNoAlpha(titleImg, image.Point{titleImg.Bounds().Min.X, start}, titleImg.Bounds().Max)
		left := internal.RightAlpha(titleImg, titleImg.Bounds().Min, image.Point{titleImg.Bounds().Max.X, start})
		left2 := internal.RightAlpha(titleImg, image.Point{X: titleImg.Bounds().Min.X, Y: start}, image.Point{titleImg.Bounds().Max.X, start * 2})
		newImg := image.NewRGBA(image.Rect(0, 0, left+left2+1, start))
		draw.Draw(newImg, newImg.Bounds(), titleImg, image.Point{}, draw.Src)
		draw.Draw(newImg, newImg.Bounds().Add(image.Point{X: left}), titleImg, image.Point{Y: start2}, draw.Src)
		internal.WritePng(newImg, "textures/ui/title.png", p.out)
	}
	return nil
}

func (p *porter) particlesFix() {
	particleData, err := p.out.Read("textures/particle/particles.png")
	if err != nil {
		return
	}
	particleImg, err := png.Decode(bytes.NewBuffer(particleData))
	if err != nil {
		return
	}

	outImg := imaging.Clone(particleImg)

	particleSizeX := particleImg.Bounds().Dx() / particleGridSize
	particleSizeY := particleImg.Bounds().Dy() / particleGridSize

	for x := 0; x < particleGridSize; x++ {
		for y := 0; y < particleGridSize; y++ {
			currentX := x * particleSizeX
			currentY := y * particleSizeY
			isBlank := true
		T:
			for xx := 0; xx < particleSizeX; xx++ {
				for yy := 0; yy < particleSizeY; yy++ {
					a := outImg.NRGBAAt(currentX+xx, currentY+yy).A
					if a > 5 {
						isBlank = false
						break T
					}
				}
			}
			if isBlank {
				vanillaParticle := imaging.Crop(vanillaParticles, image.Rect(x*vanillaParticleSize, y*vanillaParticleSize, (x+1)*vanillaParticleSize, (y+1)*vanillaParticleSize))
				vanillaParticle = imaging.Resize(vanillaParticle, particleSizeX, particleSizeY, imaging.NearestNeighbor)
				draw.Draw(outImg, image.Rect(currentX, currentY, currentX+particleSizeX, currentY+particleSizeY), vanillaParticle, image.Point{}, draw.Src)
			}
		}
	}

	internal.WritePng(outImg, "textures/particle/particles.png", p.out)
}
