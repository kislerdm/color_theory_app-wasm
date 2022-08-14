package colorname

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"io"
	"log"
	"math"
	"strconv"
)

//go:embed colors.csv
var colorsDataRaw []byte

type color struct {
	Name    string
	Hex     string
	R, G, B float64
}

type colors []color

type distancesMap map[string]float64

const distMaxThreshold = 50.

func (c colors) calculateAcceptableDistances(r, g, b float64) distancesMap {
	o := distancesMap{}
	for _, el := range c {
		distance := math.Sqrt(math.Pow(el.R-r, 2) + math.Pow(el.G-g, 2) + math.Pow(el.B-b, 2))
		if distance <= distMaxThreshold {
			o[el.Name] = distance
		}
	}
	return o
}

func (c colors) ReadColorNameByRGB(r, g, b float64) (output string) {
	distances := c.calculateAcceptableDistances(r, g, b)

	if len(distances) == 0 {
		return
	}

	distanceMin := distMaxThreshold
	for name, distance := range distances {
		if distance < distanceMin {
			distanceMin = distance
			output = name
		}
	}

	return
}

var colorsData colors

func mustStringToFloat(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return v
}

func init() {
	r := csv.NewReader(bytes.NewReader(colorsDataRaw))

	rowInd := 0
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
				Hex:  row[1],
				R:    mustStringToFloat(row[2]),
				G:    mustStringToFloat(row[3]),
				B:    mustStringToFloat(row[4]),
			},
		)
	}
}

// FindColorNameByRGB returns color name.
func FindColorNameByRGB(r, g, b float64) string {
	return colorsData.ReadColorNameByRGB(r, g, b)
}
