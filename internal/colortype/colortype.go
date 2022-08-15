package colortype

import (
	_ "embed"
	"log"
)

var m *Model

func init() {
	var err error
	m, err = LoadModelConfig()
	if err != nil {
		log.Fatalln("cannot load the model", err)
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
