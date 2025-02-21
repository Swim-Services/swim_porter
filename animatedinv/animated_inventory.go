package animatedinv

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"strconv"

	"image/gif"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/resource"
	"github.com/swim-services/swim_porter/utils"
)

const out_x, out_y = 352, 332

//go:embed assets.zip
var assetsZip []byte

//go:embed slots.png
var slotsPng []byte

var assets map[string][]byte
var slots image.Image

func init() {
	zipMap, err := utils.Unzip(assetsZip)
	if err != nil {
		panic(err)
	}
	assets = zipMap
	slotsImg, err := png.Decode(bytes.NewReader(slotsPng))
	if err != nil {
		panic(err)
	}
	slots = slotsImg
}

type invmaker struct {
	in        *gif.GIF
	static    image.Image
	out       *utils.MapFS
	lastFrame *image.NRGBA
	splitGif  []*image.NRGBA
	name      string
	overlay   bool
	addSlots  bool
}

func MakeAnimated(in *gif.GIF, name string, addSlots bool) ([]byte, error) {
	maker := &invmaker{in: in, out: utils.NewMapFS(cloneMap(assets)), name: name, addSlots: addSlots}
	maker.splitGif = SplitAnimatedGIF(maker.in)
	if err := maker.manifest(); err != nil {
		return []byte{}, err
	}
	if err := maker.icon(); err != nil {
		return []byte{}, err
	}
	if err := maker.createUi(); err != nil {
		return []byte{}, err
	}
	if err := maker.makeGifInventory(); err != nil {
		return []byte{}, err
	}
	zipped, err := utils.Zip(maker.out.RawMap())
	if err != nil {
		return []byte{}, err
	}
	return zipped, nil
}

func MakeOverlay(in image.Image, name string, addSlots bool) ([]byte, error) {
	maker := &invmaker{static: in, out: utils.NewMapFS(cloneMap(assets)), name: name, addSlots: addSlots, overlay: true}
	if err := maker.manifest(); err != nil {
		return []byte{}, err
	}
	if err := maker.icon(); err != nil {
		return []byte{}, err
	}
	if err := maker.createUi(); err != nil {
		return []byte{}, err
	}
	if err := maker.makeOverlay(); err != nil {
		return []byte{}, err
	}
	zipped, err := utils.Zip(maker.out.RawMap())
	if err != nil {
		return []byte{}, err
	}
	return zipped, nil
}

func (p *invmaker) makeGifInventory() error {
	out := image.NewRGBA(image.Rect(0, 0, max(out_x, out_x*int(math.Ceil(float64(len(p.in.Image))/40))), out_y*min(40, len(p.in.Image))))
	for i, frame := range p.splitGif {
		nrgbaImg := (frame)
		startx := (i / 40) * out_x
		starty := (i % 40) * out_y
		draw.Draw(out, image.Rect(startx, starty, startx+out_x, starty+out_y), imaging.Resize(nrgbaImg, out_x, out_y, imaging.Box), image.Point{}, draw.Over)
	}
	var overlay image.Image
	if p.addSlots {
		overlay = slots
	} else {
		overlay = image.NewRGBA(image.Rect(0, 0, out_x, out_y))
	}
	blankWriter := bytes.NewBuffer([]byte{})
	if err := png.Encode(blankWriter, overlay); err != nil {
		return err
	}
	p.out.Write(blankWriter.Bytes(), "textures/animated_ui/inventory_bg/inventory_overlay.png")
	writer := bytes.NewBuffer([]byte{})
	if err := png.Encode(writer, out); err != nil {
		return err
	}
	p.out.Write(writer.Bytes(), "textures/animated_ui/inventory_bg/inventory_vertical_flipbook.png")
	return nil
}

func (p *invmaker) makeOverlay() error {
	out := imaging.Resize(p.static, out_x, out_y, imaging.Box)
	if p.addSlots {
		draw.Draw(out, out.Bounds(), slots, image.Point{}, draw.Over)
	}
	writer := bytes.NewBuffer([]byte{})
	if err := png.Encode(writer, out); err != nil {
		return err
	}
	p.out.Write(writer.Bytes(), "textures/animated_ui/inventory_bg/inventory_overlay.png")
	return nil
}

func (p *invmaker) manifest() error {
	bedrockManifest := resource.Manifest{
		FormatVersion: 1,
		Header: resource.Header{
			Name:               p.name + " Inventory",
			Description:        p.name + "\n§aInventory by §dSwim Inventory Maker §f| §bdiscord.gg/swim",
			UUID:               uuid.New(),
			Version:            [3]int{2, 0, 0},
			MinimumGameVersion: [3]int{1, 12, 1},
		},
		Modules: []resource.Module{{
			UUID:        uuid.New().String(),
			Description: p.name + "\n§aInventory by §dSwim Inventory Maker §f| §bdiscord.gg/swim",
			Type:        "resources",
			Version:     [3]int{2, 0, 0},
		}},
	}
	bedrockManifestBytes, err := json.Marshal(bedrockManifest)
	if err != nil {
		return err
	}
	p.out.Write(bedrockManifestBytes, "manifest.json")
	return nil
}

func (p *invmaker) icon() error {
	var img image.Image
	if p.overlay {
		img = p.static
	} else {
		if len(p.in.Image) < 1 {
			return nil
		}
		img = p.in.Image[0]
	}
	writer := bytes.NewBuffer([]byte{})
	if err := png.Encode(writer, img); err != nil {
		return err
	}
	p.out.Write(writer.Bytes(), "pack_icon.png")
	return nil
}

func (p *invmaker) createUi() error {
	if !p.overlay {
		uiFile, err := p.out.Read("ui/_global_variables.json")
		if err != nil {
			return err
		}
		uiFile = bytes.ReplaceAll(uiFile, []byte("num"), []byte(strconv.Itoa(len(p.in.Image))))
		uiFile = bytes.ReplaceAll(uiFile, []byte("fum_frames"), []byte(fmt.Sprintf("%f", arrayAvg(p.in.Delay)/100)))
		p.out.Write(uiFile, "ui/_global_variables.json")
	}
	baseFile, err := p.out.Read("uidx/animated_ui/inventory_bg_base.uidx")
	if err != nil {
		return err
	}
	if !p.overlay {
		baseFile = append(bytes.TrimSuffix(baseFile, []byte("\n")), []byte(",\n")...)
		for i := range p.in.Image {
			next := i + 2
			after := ","
			if i+1 == len(p.in.Image) {
				next = 1
				after = ""
			}
			newStr := fmt.Sprintf(`"%d@CrisXolt_anm_inv_bg_base.inventory_uv_base": { "$uv_frame": [ %d, %d], "$max_inventory_uv_frames": "($total_inventory_frames = '%d_frames')", "$next_frame": "@CrisXolt_anm_inv_bg_base.%d" }%s
`, i+1, out_x*(i/40), out_y*(i%40), i+1, next, after)
			baseFile = append(baseFile, []byte(newStr)...)
		}
	}
	baseFile = append(baseFile, byte('}'))
	p.out.Write(baseFile, "uidx/animated_ui/inventory_bg_base.uidx")

	return nil
}

func arrayAvg(in []int) float64 {
	if len(in) == 0 {
		return 0
	}
	sum := 0
	for _, cur := range in {
		sum += cur
	}
	return (float64(sum)) / (float64(len(in)))
}
