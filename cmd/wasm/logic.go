package main

import (
	"fmt"

	"github.com/kislerdm/color_theory_app-wasm/internal/colorname"
	"github.com/kislerdm/color_theory_app-wasm/internal/colotype"
)

type Output struct {
	Name   string
	IsWarm bool
}

func (o *Output) generateUI() (html string) {
	s := o.Name
	if o.Name == "" {
		s = "Not found"
	}
	html += `<div><label for="output_name" id="output_label">Color Name:</label><output name="color_name" id="output_name"> ` +
		s + `</output></div>`

	s = "Cool"
	if o.IsWarm {
		s = "Warm"
	}
	html += `<div><label for="output_type" id="output_label">Color Type:</label><output name="color_type" id="output_type"> ` + s + `</output></div>`

	return
}

func generateOutput(r, g, b float64) (Output, error) {
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return Output{}, fmt.Errorf("wrong RGB input")
	}

	t, err := colotype.FindColorTypeByRGB(r, g, b)
	if err != nil {
		return Output{}, err
	}

	return Output{
		Name:   colorname.FindColorNameByRGB(r, g, b),
		IsWarm: t,
	}, nil
}
