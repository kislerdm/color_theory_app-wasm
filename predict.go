package app

import (
	_ "embed"
	"log"
)

//go:embed train/model.json
var modelCfg []byte

var m *Model

func init() {
	var err error
	m, err = LoadModelConfig(modelCfg)
	if err != nil {
		log.Fatalln("cannot load the model", err)
	}
}

// Input defines prediction input.
type Input struct {
	R, G, B float64
}

// Predict predicts the color type based on its RGB code.
func Predict(data Input) (bool, error) {
	r, err := m.Predict(
		SparceMatrix{
			{
				"r": data.R,
				"g": data.G,
				"b": data.B,
			},
		},
	)

	if err != nil {
		return false, err
	}

	return r[0], nil
}
