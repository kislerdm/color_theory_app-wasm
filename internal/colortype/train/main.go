//go:build !unittest
// +build !unittest

// Model object generator
package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/kislerdm/color_theory_app-wasm/internal/colortype"
)

//go:embed model/model.json
var modelCfg []byte

func main() {
	var v colortype.XGB
	if err := json.NewDecoder(bytes.NewReader(modelCfg)).Decode(&v); err != nil {
		log.Fatalln(err)
	}

	if err := generateModelGoFile(&v); err != nil {
		log.Fatalln(err)
	}
}

func getOutputPath() string {
	p, _ := os.Getwd()
	return path.Dir(p) + "/model.go"
}

func generateModelGoFile(v *colortype.XGB) error {
	const t = `{{- define "node" -}}{
		ID:        {{ .ID }},
		{{- if .Children }}
		Depth:     {{ .Depth }},
		Feature:   "{{ .Feature }}",
		Threshold: {{ .Threshold }},
		Yes:       {{ .Yes }},
		No:        {{ .No }},
		Missing:   {{ .Missing }},
		Children:  []*node{
		{{- range .Children }}
			{{ template "node" . }}
		{{- end }}
	},
		{{- end }}
		{{- if not .Children }}
		Leaf:      {{ .Leaf }},
		{{- end }}
},
{{- end }}package colortype

// ModelDef defines the XGB model trees.
var ModelDef = XGB{
{{- range . }}
	{{ template "node" . }}
{{- end }}
}`
	log.Println(getOutputPath())
	f, err := os.Create(getOutputPath())
	if err != nil {
		return err
	}
	defer f.Close()

	return template.Must(template.New("xgbmodel").Parse(t)).Execute(f, &v)
}
