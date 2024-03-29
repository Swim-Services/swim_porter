package port

import (
	"encoding/json"

	"github.com/swim-services/swim_porter/port/porterror"
	"github.com/swim-services/swim_porter/port/utils"

	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/resource"
)

type PortOptions struct {
	ShowCredits    bool
	SkyboxOverride string
}

type porter struct {
	in   *utils.MapFS
	out  *utils.MapFS
	name string
}

func Port(in []byte, name string, opts PortOptions) ([]byte, error) {
	zipMap, err := utils.Unzip(in)
	if err != nil {
		return []byte{}, err
	}
	outFS, err := PortRaw(utils.NewMapFS(zipMap), name, opts)
	if err != nil {
		return []byte{}, err
	}
	out, err := utils.Zip(outFS.RawMap())
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}

func PortRaw(in *utils.MapFS, name string, opts PortOptions) (*utils.MapFS, error) {
	p := &porter{in: in, out: utils.NewMapFS(make(map[string][]byte)), name: name}
	err := p.doPort(opts)
	if err != nil {
		return nil, err
	}
	return p.out, nil
}

func (p *porter) doPort(opts PortOptions) error {
	if err := p.manifest(opts.ShowCredits); err != nil {
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
	if err := p.environment(opts.SkyboxOverride); err != nil {
		return err
	}
	if err := p.misc(); err != nil {
		return err
	}
	p.itemsFix()
	if err := p.beds(); err != nil {
		return err
	}
	if err := p.entityFixes(); err != nil {
		return err
	}
	return nil
}

func (p *porter) manifest(showCredits bool) error {
	meta, err := p.in.Read("pack.mcmeta")
	if err != nil {
		return porterror.New("pack.mcmeta not found")
	}
	javaMeta, err := utils.PackMcmeta(meta)
	var desc string
	if err != nil {
		desc = p.name
	} else {
		desc = javaMeta.Pack.Description
	}
	if showCredits {
		desc += "\n§aPorted by §dSwim Auto Port §f| §bdiscord.gg/swim"
	}
	bedrockManifest := resource.Manifest{
		FormatVersion: 1,
		Header: resource.Header{
			Name:               p.name,
			Description:        desc,
			UUID:               uuid.New().String(),
			Version:            [3]int{2, 0, 0},
			MinimumGameVersion: [3]int{1, 12, 1},
		},
		Modules: []resource.Module{{
			UUID:        uuid.New().String(),
			Description: desc,
			Type:        "resources",
			Version:     [3]int{2, 0, 0},
		}},
	}
	bedrockManifestBytes, err := json.Marshal(bedrockManifest)
	if err != nil {
		return porterror.Wrap(err)
	}
	p.out.Write(bedrockManifestBytes, "/manifest.json")
	return nil
}

func (p *porter) icon() {
	if icon, err := p.in.Read("pack.png"); err == nil {
		p.out.Write(icon, "/pack_icon.png")
	}
}
