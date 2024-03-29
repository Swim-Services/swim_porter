package port

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"path"
	"slices"
	"strings"

	"github.com/swim-services/swim_porter/port/internal"
	"github.com/swim-services/swim_porter/port/porterror"
	"github.com/swim-services/swim_porter/port/utils"

	"github.com/disintegration/imaging"
)

func (p *porter) ui() error {
	p.font()
	p.gui()
	if err := p.guiFix(); err != nil {
		return err
	}
	if err := p.containerUI(); err != nil {
		return err
	}
	if err := p.crosshair(); err != nil {
		return err
	}
	if err := p.mobileButtons(); err != nil {
		return err
	}
	return nil
}

func (p *porter) font() {
	p.out.InterCopyDir(p.in, "assets/minecraft/mcpatcher/font", "font")
	p.out.InterCopyDir(p.in, "assets/minecraft/font", "font")
	p.out.InterCopyDir(p.in, "assets/minecraft/textures/font", "font")
	p.out.Rename("font/ascii.png", "font/default8.png")
}

func (p *porter) gui() {
	p.out.Rename("textures/gui/widgets.png", "textures/gui/gui.png")
}

func (p *porter) guiFix() error {
	if gui, err := p.out.Read("textures/gui/gui.png"); err == nil {
		guiImg, err := png.Decode(bytes.NewReader(gui))
		if err != nil {
			return porterror.Wrap(err)
		}
		if guiImg.Bounds().Dy() != 256 {
			if err := internal.WritePng(imaging.Resize(guiImg, 256, 256, imaging.NearestNeighbor), "textures/gui/gui.png", p.out); err != nil {
				return porterror.Wrap(err)
			}
		}
	}
	return nil
}

func (p *porter) crosshair() error {
	if data, err := p.out.Read("textures/gui/icons.png"); err == nil {
		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			return porterror.Wrap(err)
		}
		iconsImg := img.(interface {
			SubImage(r image.Rectangle) image.Image
		})
		crosshairSize := img.Bounds().Dx() / 16
		crosshairImg := iconsImg.SubImage(image.Rect(0, 0, crosshairSize, crosshairSize))
		canvas := image.NewRGBA(image.Rect(0, 0, crosshairSize, crosshairSize))
		draw.Draw(canvas, canvas.Bounds(), image.Black, image.Point{}, draw.Src)
		draw.Draw(canvas, canvas.Bounds(), crosshairImg, image.Point{}, draw.Over)
		if err := internal.WritePng(canvas, "textures/ui/cross_hair.png", p.out); err != nil {
			return porterror.Wrap(err)
		}
	}
	return nil
}

func (p *porter) mobileButtons() error {
	if data, err := p.out.Read("textures/gui/gui.png"); err == nil {
		gui, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			return porterror.Wrap(err)
		}
		mobileImgData, err := assetsMapFS.Read("mobile/mobile_buttons.png")
		if err != nil {
			return porterror.Wrap(err)
		}
		mobileImg, err := png.Decode(bytes.NewReader(mobileImgData))
		if err != nil {
			return porterror.Wrap(err)
		}
		canvas := image.NewRGBA(gui.Bounds())
		draw.Draw(canvas, gui.Bounds(), gui, image.Point{}, draw.Src)
		draw.Draw(canvas, gui.Bounds(), mobileImg, image.Point{}, draw.Over)
		if err := internal.WritePng(canvas, "textures/gui/gui.png", p.out); err != nil {
			return porterror.Wrap(err)
		}
	}
	return nil
}

func (p *porter) containerUI() error {
	_, check1 := p.out.Read("textures/gui/container/inventory.png")
	_, check2 := p.out.Read("textures/gui/container/double_normal.png")
	if check1 != nil && check2 != nil {
		return nil
	}
	p.out.InterCopyDir(assetsMapFS, "ui", "ui")
	p.out.InterCopyDir(assetsMapFS, "textures_uidx/uidx", "textures/uidx")
	p.out.InterCopyDir(p.in, "assets/minecraft/textures/gui/container", "assets/minecraft/textures/gui/container")
	p.out.InterCopyDir(assetsMapFS, "recipe_book", "assets/uidx/textures/gui/container/recipe_book")
	p.out.Copy("assets/minecraft/textures/gui/container/generic_54.png", "assets/uidx/textures/gui/container/generic_54.png")
	p.out.Copy("assets/minecraft/textures/gui/container/generic_54.png", "assets/uidx/textures/gui/container/ender_chest.png")
	if _, err := p.out.Read("assets/minecraft/textures/gui/container/creative_inventory"); err != nil {
		p.out.InterCopyDir(assetsMapFS, "container/creative_inventory", "assets/minecraft/textures/gui/container/creative_inventory")
	}
	if p.out.Copy("assets/minecraft/textures/gui/container/generic_54.png", "assets/uidx/textures/gui/container/small_chest.png") != nil {
		for path, data := range assetsMapFS.Dir("assets/container/") {
			if _, err := p.out.Read("assets/minecraft/textures/gui/container/" + path); err != nil {
				p.out.Write(data, "assets/minecraft/textures/gui/container/"+path)
			}
		}
	}
	globalVarData, err := p.out.Read("ui/_global_variables.json")
	if err != nil {
		return porterror.Wrap(err)
	}
	globalVars := string(globalVarData)
	for filePath, data := range p.out.Dir("assets/minecraft/textures/gui/container/") {
		if filePath == "blast_furnace.png" || filePath == "smoker.png" {
			filePath = "furnace.png"
		}
		if !strings.HasSuffix(filePath, ".png") {
			continue
		}
		fileExtension := path.Ext(filePath)
		fileName := filePath[:strings.LastIndex(filePath, fileExtension)]

		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			return porterror.Wrap(err)
		}
		width := img.Bounds().Dx()
		height := img.Bounds().Dy()
		dimension := height
		validDimensions := []int{256, 512, 1024, 2048, 4096, 8192}
		validWidth := slices.Contains(validDimensions, width)
		validHeight := slices.Contains(validDimensions, height)
		if height != width || !validHeight || !validWidth {
			dimension = utils.FindClosestDimension(height, width, validDimensions)
			if err := internal.WritePng(imaging.Resize(img, dimension, dimension, imaging.NearestNeighbor), filePath, p.out); err != nil {
				return porterror.Wrap(err)
			}
		}
		globalVars = strings.ReplaceAll(globalVars, fmt.Sprintf("\"$%s_resolution\": \"256x\",", fileName), fmt.Sprintf("\"$%s_resolution\": \"%dx\",", fileName, dimension))
	}
	p.out.Write([]byte(globalVars), "ui/_global_variables.json")
	return nil
}
