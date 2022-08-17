//go:build gen
// +build gen

// colorname datagen
package main

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"io"
	"log"
	"os"
	"path"
	"text/template"
)

//go:embed colors.csv
var colorsDataRaw []byte

func getOutputPath() string {
	p, _ := os.Getwd()
	return path.Dir(p) + "/colors.go"
}

type color struct {
	Name, R, G, B string
}

func main() {
	r := csv.NewReader(bytes.NewReader(colorsDataRaw))

	rowInd := 0
	var colorsData []color
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		if rowInd == 0 {
			rowInd++
			continue
		}

		colorsData = append(
			colorsData, color{
				Name: row[0],
				R:    row[2],
				G:    row[3],
				B:    row[4],
			},
		)
	}

	if err := generateColorsGoFile(colorsData); err != nil {
		log.Fatalln(err)
	}
}

func generateColorsGoFile(v []color) error {
	const t = `{{- define "color" -}}{ Name: "{{ .Name }}", R: {{ .R }}, G: {{ .G }}, B: {{ .B }} },{{- end }}//go:build !gen
// +build !gen

package colorname

var colorsData = []color{
{{- range . }}
	{{ template "color" . }}
{{- end }}
}`
	f, err := os.Create(getOutputPath())
	if err != nil {
		return err
	}
	defer f.Close()

	return template.Must(template.New("colors").Parse(t)).Execute(f, v)
}
