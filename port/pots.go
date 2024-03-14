package port

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"math"

	"github.com/swim-services/swim_porter/port/internal"
	"github.com/swim-services/swim_porter/port/recolor"
	"github.com/swim-services/swim_porter/port/utils"

	"github.com/disintegration/imaging"
)

func (p *porter) pots() error {
	if err := p.tintPots(false); err != nil {
		return err
	}
	if err := p.tintPots(true); err != nil {
		return err
	}
	if err := p.potionEffectsUI(); err != nil {
		return err
	}
	return nil
}

func (p *porter) tintPots(splash bool) error {
	var data []byte
	var err error
	if splash {
		data, err = p.out.Read("textures/items/potion_bottle_splash.png")
	} else {
		data, err = p.out.Read("textures/items/potion_bottle_drinkable.png")
	}
	if err != nil {
		return nil
	}
	overlayBytes, err := p.out.Read("textures/items/potion_overlay.png")
	if err != nil {
		return nil
	}
	blank, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}
	overlay, err := png.Decode(bytes.NewReader(overlayBytes))
	if err != nil {
		return err
	}
	if overlay.Bounds() != blank.Bounds() {
		overlay = imaging.Resize(overlay, blank.Bounds().Dx(), blank.Bounds().Dy(), imaging.NearestNeighbor)
	}
	for potType, color := range utils.POTS_MAP {
		over := recolor.Tint(overlay, color)
		canvas := image.NewRGBA(blank.Bounds())
		draw.Draw(canvas, blank.Bounds(), blank, image.Point{0, 0}, draw.Src)
		draw.Draw(canvas, blank.Bounds(), over, image.Point{0, 0}, draw.Over)
		writer := bytes.NewBuffer([]byte{})
		if err := png.Encode(writer, canvas); err != nil {
			return err
		}
		imgBytes := writer.Bytes()
		if splash {
			p.out.Write(imgBytes, "textures/items/potion_bottle_splash_"+potType+".png")
		} else {
			p.out.Write(imgBytes, "textures/items/potion_bottle_"+potType+".png")
		}
	}
	return nil
}

func (p *porter) potionEffectsUI() error {
	if inv, err := p.out.Read("textures/gui/container/inventory.png"); err == nil {
		invImg, err := png.Decode(bytes.NewReader(inv))
		if err != nil {
			return err
		}
		bounds := invImg.Bounds()
		sin := float64(bounds.Dx()) / 4.41379310345
		epilsonFlat := int(math.Round(sin))
		startingY := bounds.Dy() - epilsonFlat
		hypo := float64(bounds.Dy()) / 14.2222222222
		cellChangeFactor := int(math.Round(hypo))
		effects := []string{"speed", "slowness", "haste", "mining_fatigue", "strength", "weakness", "poison", "regeneration", "invisibility", "saturation", "jump_boost", "nausea", "night_vision", "blindness", "resistance", "fire_resistance", "water_breathing", "wither", "absorption"}
		sheetRow := 0
		var x int
		for i := 0; i < 19; i++ {
			if i%8 == 0 && i != 0 {
				x = 0
				sheetRow = 1
				startingY = startingY + cellChangeFactor
			} else {
				x = sheetRow * cellChangeFactor
				sheetRow++
			}
			subImage := invImg.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(image.Rect(x, startingY, x+cellChangeFactor, startingY+cellChangeFactor))

			if err := internal.WritePng(subImage, "textures/ui/"+effects[i]+"_effect.png", p.out); err != nil {
				return err
			}
		}

	}
	return nil
}
