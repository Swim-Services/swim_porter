package port

import (
	"encoding/json"
	"errors"

	"github.com/swedeachu/swim_porter/port/utils"

	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/resource"
)

type porter struct {
	in   *utils.MapFS
	out  *utils.MapFS
	name string
}

func Port(in []byte, name string) ([]byte, error) {
	zipMap, err := utils.Unzip(in)
	if err != nil {
		return []byte{}, err
	}
	outFS, err := PortRaw(utils.NewMapFS(zipMap), name)
	if err != nil {
		return []byte{}, err
	}
	out, err := utils.Zip(outFS.RawMap())
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}

func PortRaw(in *utils.MapFS, name string) (*utils.MapFS, error) {
	p := &porter{in: in, out: utils.NewMapFS(make(map[string][]byte)), name: name}
	err := p.doPort()
	if err != nil {
		return nil, err
	}
	return p.out, nil
}

func (p *porter) doPort() error {
	if err := p.manifest(); err != nil {
		return err
	}
	p.icon()
	if err := p.textures(); err != nil {
		return err
	}
	if err := p.pots(); err != nil {
		return err
	}
	if err := p.xp(); err != nil {
		return err
	}
	if err := p.ui(); err != nil {
		return err
	}
	if err := p.environment(); err != nil {
		return err
	}
	p.misc()
	p.itemsFix()

	return nil
}

func (p *porter) manifest() error {
	meta, err := p.in.Read("pack.mcmeta")
	if err != nil {
		return errors.New("pack.mcmeta not found")
	}
	javaMeta, err := utils.PackMcmeta(meta)
	if err != nil {
		return err
	}
	bedrockManifest := resource.Manifest{
		FormatVersion: 1,
		Header: resource.Header{
			Name:               p.name,
			Description:        javaMeta.Pack.Description,
			UUID:               uuid.New().String(),
			Version:            [3]int{2, 0, 0},
			MinimumGameVersion: [3]int{1, 12, 1},
		},
		Modules: []resource.Module{{
			UUID:        uuid.New().String(),
			Description: javaMeta.Pack.Description,
			Type:        "resources",
			Version:     [3]int{2, 0, 0},
		}},
	}
	bedrockManifestBytes, err := json.Marshal(bedrockManifest)
	if err != nil {
		return err
	}
	p.out.Write(bedrockManifestBytes, "/manifest.json")
	return nil
}

func (p *porter) icon() {
	if icon, err := p.in.Read("pack.png"); err == nil {
		p.out.Write(icon, "/pack_icon.png")
	}
}
