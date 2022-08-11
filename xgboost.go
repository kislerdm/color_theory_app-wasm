package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
)

type node struct {
	ID        int     `json:"nodeid"`
	Depth     int     `json:"depth,omitempty"`
	Feature   string  `json:"split,omitempty"`
	Threshold float64 `json:"split_condition,omitempty"`
	Yes       int     `json:"yes,omitempty"`
	No        int     `json:"no,omitempty"`
	Missing   int     `json:"missing,omitempty"`
	Leaf      float64 `json:"leaf,omitempty"`
	Children  []*node `json:"children,omitempty"`
}

func (t *node) getNodeByID(id int) *node {
	if t.ID == id {
		return t
	}
	for _, nodeChild := range t.Children {
		if nodeChild.ID == id {
			return nodeChild
		}
		if n := nodeChild.getNodeByID(id); n != nil {
			return n
		}
	}
	return nil
}

// XGB defines a model.
type XGB []node

// SparceVector data vector to predict.
type SparceVector map[string]float64

// SparceMatrix dataframe input to predict.
type SparceMatrix []SparceVector

// Prediction defines prediction output results.
type Prediction []float64

// Model defines a binary classifier.
type Model struct {
	trees           XGB
	BinaryThreshold float64
}

func (m *Model) predictRow(dataRow SparceVector) (float64, error) {
	var v float64
	for _, tree := range m.trees {
		idx := 0
		for {
			n := tree.getNodeByID(idx)
			if n == nil {
				return 0, fmt.Errorf("node %d not found", idx)
			}

			if len(n.Children) == 0 {
				v += n.Leaf
				break
			}

			v, ok := dataRow[n.Feature]
			if !ok {
				idx = n.Missing
				continue
			}

			idx = n.Yes
			if v >= n.Threshold {
				idx = n.No
			}
		}
	}

	return v, nil
}

func (m *Model) Predict(dataFrame SparceMatrix) ([]bool, error) {
	var o = make([]bool, len(dataFrame))

	for i, row := range dataFrame {
		p, err := m.predictRow(row)
		if err != nil {
			return nil, err
		}

		o[i] = false
		if sigmoid(p) >= m.BinaryThreshold {
			o[i] = true
		}
	}

	return o, nil
}

// LoadModelConfig loads model JSON configuration from bytes.
func LoadModelConfig(data []byte) (*Model, error) {
	var v XGB
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&v); err != nil {
		return nil, err
	}
	return &Model{
		trees:           v,
		BinaryThreshold: 0.5,
	}, nil
}

func sigmoid(x float64) float64 {
	return 1. / (1. + math.Exp(-x))
}
