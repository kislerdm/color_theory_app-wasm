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

func FindColorNameByRGBv2(r, g, b float64) string {
	output := ""
	distance := distMaxThreshold

	for _, color := range colorsData {

		dR := color.R - r
		dG := color.G - g
		dB := color.B - b

		if dR == 0 && dG == 0 && dB == 0 {
			return color.Name
		}

		if d := sqrt(dR*dR + dG*dG + dB*dB); d < distance {
			output = color.Name
			distance = d
		}

	}

	return output
}

func sqrt(v float64) float64 {
	// from quake3 inverse sqrt algorithm
	// ref: https://medium.com/@adrien.za/fast-inverse-square-root-in-go-and-javascript-for-fun-6b891e74e5a8
	const magic64 = 0x5FE6EB50C7B537A9

	n2, th := v*0.5, float64(1.5)
	b := math.Float64bits(v)
	b = magic64 - (b >> 1)
	f := math.Float64frombits(b)
	f *= th - (n2 * f * f)
	return f
}
