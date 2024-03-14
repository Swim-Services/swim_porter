package port

import (
	_ "embed"

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

func (p *porter) misc() {
	p.panorama()
	p.sounds()
	p.armor()
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
