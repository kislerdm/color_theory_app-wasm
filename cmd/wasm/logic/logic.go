package logic

import (
	"errors"

	"github.com/kislerdm/color_theory_app-wasm/internal/colorname"
	"github.com/kislerdm/color_theory_app-wasm/internal/colortype"
)

func ui(colorName string, isWarm bool) (html string) {
	if colorName == "" {
		colorName = "Not found"
	}

	colorType := "Cool"
	if isWarm {
		colorType = "Warm"
	}

	html = `<div><label for="output_name" id="output_label">Color Name:</label><output name="color_name" id="output_name"> ` +
		colorName + `</output></div><div><label for="output_type" id="output_label">Color Type:</label><output name="color_type" id="output_type"> ` + colorType + `</output></div>`

	return
}

func Run(r, g, b float64) (html string, err error) {
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return "", errors.New("wrong RGB input")
	}

	isWarm, err := colortype.FindColorTypeByRGB(r, g, b)
	if err != nil {
		return "", err
	}

	colorName := colorname.FindColorNameByRGB(r, g, b)

	return ui(colorName, isWarm), nil
}
