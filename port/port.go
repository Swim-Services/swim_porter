package port

import (
	"encoding/json"
	"sync"

	"github.com/swim-services/swim_porter/porterror"
	"github.com/swim-services/swim_porter/utils"

	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/resource"
)

type PortOptions struct {
	ShowCredits    bool
	SkyboxOverride string
	OffsetSky      bool
}

type porter struct {
	in   *utils.MapFS
	out  *utils.MapFS
	name string
	opts PortOptions
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
	p := &porter{in: in, out: utils.NewMapFS(make(map[string][]byte)), name: name, opts: opts}
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
	wg := sync.WaitGroup{}
	wg.Add(1)
	var skyErr error
	go func() {
		defer wg.Done()
		if err := p.environment(opts.SkyboxOverride); err != nil {
			if portError, ok := err.(*porterror.PortError); ok {
				skyErr = portError.WithMessage("port environment")
				return
			}
			skyErr = err
		}
	}()
	p.icon()
	if err := p.textures(); err != nil {
		if portError, ok := err.(*porterror.PortError); ok {
			return portError.WithMessage("port textures")
		}
		return err
	}
	if err := p.pots(); err != nil {
		if portError, ok := err.(*porterror.PortError); ok {
			return portError.WithMessage("port potions")
		}
		return err
	}
	if err := p.xp(); err != nil {
		if portError, ok := err.(*porterror.PortError); ok {
			return portError.WithMessage("port xp bar")
		}
		return err
	}
	if err := p.ui(); err != nil {
		if portError, ok := err.(*porterror.PortError); ok {
			return portError.WithMessage("port ui")
		}
		return err
	}
	if err := p.misc(); err != nil {
		if portError, ok := err.(*porterror.PortError); ok {
			return portError.WithMessage("port misc")
		}
		return err
	}
	p.itemsFix()
	if err := p.beds(); err != nil {
		if portError, ok := err.(*porterror.PortError); ok {
			return portError.WithMessage("port beds")
		}
		return err
	}
	if err := p.entityFixes(); err != nil {
		if portError, ok := err.(*porterror.PortError); ok {
			return portError.WithMessage("entity fixes")
		}
		return err
	}
	wg.Wait()
	return skyErr
}

func (p *porter) manifest(showCredits bool) error {
	meta, err := p.in.Read("pack.mcmeta")
	if err != nil {
		return porterror.Wrap(porterror.ErrMcmetaNotFound)
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
			UUID:               uuid.New(),
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
