//go:build js
// +build js

package main

import (
	"fmt"
	"syscall/js"

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

func run(args []js.Value) (html string, err error) {
	if len(args) < 3 {
		return "", fmt.Errorf("no r, g, b input provided")
	}

	r := args[0].Float()
	g := args[1].Float()
	b := args[2].Float()

	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return "", fmt.Errorf("wrong RGB input")
	}

	isWarm, err := colortype.FindColorTypeByRGB(r, g, b)
	if err != nil {
		return "", err
	}

	colorName := colorname.FindColorNameByRGB(r, g, b)

	return ui(colorName, isWarm), nil
}

func main() {
	js.Global().Set(
		"start", js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				htmlString, err := run(args)
				if err != nil {
					return map[string]interface{}{"error": err.Error()}
				}

				js.Global().Get("document").Call("getElementById", "color_output").Set("innerHTML", htmlString)

				return nil
			},
		),
	)
	select {}
}
