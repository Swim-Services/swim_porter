package cubemap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"

	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/resource"
	"github.com/swim-services/swim_porter/port/utils"
)

func SkyPack(name string, img [6]image.Image, action string) ([]byte, error) {
	fs := utils.NewMapFS(make(map[string][]byte))
	for i, image := range img {
		writer := bytes.NewBuffer([]byte{})
		err := png.Encode(writer, image)
		if err != nil {
			return []byte{}, err
		}
		rawBytes := writer.Bytes()
		if i == 0 {
			fs.Write(rawBytes, "pack_icon.png")
		}
		fs.Write(rawBytes, fmt.Sprintf("textures/environment/overworld_cubemap/cubemap_%d.png", i))
	}
	var desc = name + "\n§a" + action + " by §dSwim Auto Port §f| §bdiscord.gg/swim"
	bedrockManifest := resource.Manifest{
		FormatVersion: 1,
		Header: resource.Header{
			Name:               name + " Sky",
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
		return []byte{}, err
	}
	fs.Write(bedrockManifestBytes, "manifest.json")
	rawOut, err := utils.Zip(fs.RawMap())
	if err != nil {
		return []byte{}, err
	}
	return rawOut, nil
}
