package particlefix

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"image"
	"image/draw"

	"github.com/disintegration/imaging"
	"github.com/gameparrot/fastpng"
	"github.com/swim-services/swim_porter/internal"
	"github.com/swim-services/swim_porter/porterror"
	"github.com/swim-services/swim_porter/resource"
	"github.com/swim-services/swim_porter/utils"
)

const particleGridSize = 16
const vanillaParticleSize = 8

//go:embed particles.png
var vanillaParticleData []byte

var vanillaParticles image.Image

func init() {
	var err error
	vanillaParticles, err = fastpng.Decode(bytes.NewBuffer(vanillaParticleData))
	if err != nil {
		panic(err)
	}
}

type particlefixer struct {
	in *utils.MapFS
}

func FixParticles(in []byte) ([]byte, error) {
	zipMap, err := utils.Unzip(in)
	if err != nil {
		return []byte{}, err
	}
	fs := utils.NewMapFS(zipMap)
	err = FixParticlesRaw(utils.NewMapFS(zipMap))
	if err != nil {
		return []byte{}, err
	}
	out, err := utils.Zip(fs.RawMap())
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}

func DoFix(in *utils.MapFS) error {
	particleData, err := in.Read("textures/particle/particles.png")
	if err != nil {
		return err
	}
	particleImg, err := fastpng.Decode(bytes.NewBuffer(particleData))
	if err != nil {
		return porterror.Wrap(err).WithMessage("read image textures/particle/particles.png")
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

	internal.WritePng(outImg, "textures/particle/particles.png", in)
	return nil
}

func FixParticlesRaw(in *utils.MapFS) error {
	p := &particlefixer{in: in}
	err := p.doParticleFix()
	if err != nil {
		return err
	}
	return nil
}

func (p *particlefixer) doParticleFix() error {
	if err := p.manifest(); err != nil {
		return err
	}
	DoFix(p.in)
	return nil
}

func (p *particlefixer) manifest() error {
	bedrockManifestOrig, err := p.in.Read("manifest.json")
	if err != nil {
		return porterror.Wrap(porterror.ErrManifestNotFound)
	}
	bedrockManifest, err := resource.UnmarshalJSON(bedrockManifestOrig)
	if err != nil {
		return err
	}
	utils.ChangeUUID(&bedrockManifest)
	bedrockManifest.Header.Description += "\n§aParticles fixed by §dSwim Particle Fixer §f| §bdiscord.gg/swim"
	bedrockManifestBytes, err := json.Marshal(bedrockManifest)
	if err != nil {
		return err
	}
	p.in.Write(bedrockManifestBytes, "manifest.json")
	return nil
}
