package utils

import (
	"errors"
	"image/color"
	"strconv"
)

func ParseHex(hex string) (color.RGBA, error) {
	if len(hex) != 6 {
		return color.RGBA{}, errors.New("invalid hex")
	}
	r, err := strconv.ParseInt(hex[0:2], 16, 64)
	if err != nil {
		return color.RGBA{}, err
	}
	g, err := strconv.ParseInt(hex[2:4], 16, 64)
	if err != nil {
		return color.RGBA{}, err
	}
	b, err := strconv.ParseInt(hex[4:6], 16, 64)
	if err != nil {
		return color.RGBA{}, err
	}

	// Create a color.Color instance from the parsed components
	c := color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255, // Alpha channel, assuming fully opaque
	}
	return c, nil
}
