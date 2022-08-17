package colorname

import (
	_ "embed"
	"math"
)

type color struct {
	Name    string
	R, G, B float64
}

type colors []color

type distancesMap map[float64]string

const distMaxThreshold = 50.

func (c colors) calculateAcceptableDistances(r, g, b float64) distancesMap {
	o := distancesMap{}
	for _, el := range c {
		distance := math.Sqrt(math.Pow(el.R-r, 2) + math.Pow(el.G-g, 2) + math.Pow(el.B-b, 2))
		if distance <= distMaxThreshold {
			o[distance] = el.Name
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
	for distance, name := range distances {
		if distance < distanceMin {
			distanceMin = distance
			output = name
		}
	}

	return
}

// FindColorNameByRGB returns Color name.
func FindColorNameByRGB(r, g, b float64) string {
	return colorsData.ReadColorNameByRGB(r, g, b)
}
