//go:build js && !unittest
// +build js,!unittest

package main

import (
	"syscall/js"

	"github.com/kislerdm/color_theory_app-wasm/cmd/wasm/logic"
)

func main() {
	js.Global().Set(
		"start", js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				if len(args) < 3 {
					return map[string]interface{}{"error": "no r, g, b input provided"}
				}

				htmlString, err := logic.Run(
					args[0].Float(),
					args[1].Float(),
					args[2].Float(),
				)

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
