package crosshairmaker

import (
	"encoding/json"
	"image"
	"image/color"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/resource"
	"github.com/swim-services/swim_porter/internal"
	"github.com/swim-services/swim_porter/utils"
)

func CrosshairPack(name string, img image.Image, scale float64, isColor bool) ([]byte, error) {
	fs := utils.NewMapFS(make(map[string][]byte))
	if err := internal.WritePng(img, "pack_icon.png", fs); err != nil {
		return []byte{}, err
	}

	if err := internal.WritePng(makeCrosshairFromImage(img, scale, isColor), "textures/ui/cross_hair.png", fs); err != nil {
		return []byte{}, err
	}
	var desc = name + "\n§aCrosshair made by §dSwim Crosshair Maker §f| §bdiscord.gg/swim"
	bedrockManifest := resource.Manifest{
		FormatVersion: 1,
		Header: resource.Header{
			Name:               name + " Crosshair",
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

func makeCrosshairFromImage(img image.Image, textureScale float64, isColor bool) image.Image {
	size := min(128, img.Bounds().Dx())
	textureSize := int(float64(size) * textureScale)
	offset := (size - textureSize) / 2
	out := image.NewNRGBA(image.Rect(0, 0, size, size))
	in := imaging.Resize(img, textureSize, textureSize, imaging.NearestNeighbor)
	for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x++ {
		for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
			if isColor {
				out.Set(x+offset, y+offset, in.At(x, y))
			} else if in.NRGBAAt(x, y).A > 10 {
				out.Set(x+offset, y+offset, color.White)

			}
		}
	}
	return out
}
