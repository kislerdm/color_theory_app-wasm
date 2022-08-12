// Package logic defines the logic to handle client input in web browser.
// go:build js && wasm
package logic

import (
	"fmt"
	"syscall/js"

	"github.com/kislerdm/color_theory_app-wasm/internal/colorname"
	"github.com/kislerdm/color_theory_app-wasm/internal/colotype"
)

func sendErr(msg string) map[string]interface{} {
	return map[string]interface{}{"error": msg}
}

type Output struct {
	Name   string
	IsWarm bool
}

type processor interface {
	GenerateOutput(r, g, b float64) (Output, error)
	GenerateUI(Output) string
	SetUI(html string)
}

type p struct{}

func (p *p) GenerateOutput(r, g, b float64) (Output, error) {
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

func (p *p) GenerateUI(o Output) string {
	html := ""

	if o.Name != "" {
		html += `<div><label for="output_name" id="output_label">Color Name:</label><output name="color_name" id="output_name"> ` +
			o.Name + `</output></div>\n`
	}
	t := "Cool"
	if o.IsWarm {
		t = "Warm"
	}
	html += `<div><label for="output_type" id="output_label">Color Type:</label><output name="color_type" id="output_type"> ` + t + `</output></div>\n`

	return html
}

func (p *p) SetUI(html string) {
	js.Global().Get("document").Call("getElementById", "color_output").Set(
		"innerHTML", html,
	)
}

func exec(p processor) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		if len(args) < 3 {
			return sendErr("no r, g, b input provided")
		}

		r := args[0].Float()
		g := args[1].Float()
		b := args[2].Float()

		o, err := p.GenerateOutput(r, g, b)
		if err != nil {
			return sendErr(err.Error())
		}

		p.SetUI(p.GenerateUI(o))

		return nil
	}
}

// Run executes the logic upon the input.
func Run() {
	js.Global().Set("start", js.FuncOf(exec(&p{})))
}
