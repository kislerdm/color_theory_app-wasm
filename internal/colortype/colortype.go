//go:build !gen
// +build !gen

package colortype

import (
	_ "embed"
	"log"
)

var m *Model

func init() {
	if modelDef == nil {
		log.Fatalln("wrong model generated")
	}
	m = &Model{
		trees:           modelDef,
		BinaryThreshold: 0.5,
	}
}

// FindColorTypeByRGB predicts the color type based on its RGB code.
func FindColorTypeByRGB(r, g, b float64) (bool, error) {
	p, err := m.Predict(
		SparceMatrix{
			{
				"r": r,
				"g": g,
				"b": b,
			},
		},
	)

	if err != nil {
		return false, err
	}

	return p[0], nil
}
