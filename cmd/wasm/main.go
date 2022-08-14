//go:build js
// +build js

package main

import (
	"syscall/js"
)

func sendErr(msg string) map[string]interface{} {
	return map[string]interface{}{"error": msg}
}

func run(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return sendErr("no r, g, b input provided")
	}

	r := args[0].Float()
	g := args[1].Float()
	b := args[2].Float()

	o, err := generateOutput(r, g, b)
	if err != nil {
		return sendErr(err.Error())
	}

	js.Global().Get("document").Call("getElementById", "color_output").Set("innerHTML", o.generateUI())

	return nil
}

func main() {
	js.Global().Set("start", js.FuncOf(run))
	select {}
}
