package app

import (
	"reflect"
	"testing"
)

func TestLoadModelConfig(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Model
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				data: []byte(`[{ "nodeid": 0, "depth": 0, "split": "b", "split_condition": 162, "yes": 1, "no": 2, "missing": 1 , "children": [
    { "nodeid": 1, "depth": 1, "split": "r", "split_condition": 86, "yes": 3, "no": 4, "missing": 3 , "children": [
      { "nodeid": 3, "depth": 2, "split": "g", "split_condition": 98, "yes": 7, "no": 8, "missing": 7 , "children": [
        { "nodeid": 7, "leaf": -0 }, 
        { "nodeid": 8, "leaf": -0.666666687 }
      ]}, 
      { "nodeid": 4, "depth": 2, "split": "b", "split_condition": 112, "yes": 9, "no": 10, "missing": 9 , "children": [
        { "nodeid": 9, "leaf": 0.70588237 }, 
        { "nodeid": 10, "leaf": 0.176470593 }
      ]}
    ]}, 
    { "nodeid": 2, "depth": 1, "split": "r", "split_condition": 247, "yes": 5, "no": 6, "missing": 5 , "children": [
      { "nodeid": 5, "leaf": -0.826086938 }, 
      { "nodeid": 6, "depth": 2, "split": "g", "split_condition": 221, "yes": 11, "no": 12, "missing": 11 , "children": [
        { "nodeid": 11, "leaf": -0.200000003 }, 
        { "nodeid": 12, "leaf": 0.25 }
      ]}
    ]}
  ]}]`),
			},
			want: &Model{
				trees: XGB{
					{
						ID:        0,
						Depth:     0,
						Feature:   "b",
						Threshold: 162,
						Yes:       1,
						No:        2,
						Missing:   1,
						Leaf:      0,
						Children: []*node{
							{
								ID:        1,
								Depth:     1,
								Feature:   "r",
								Threshold: 86,
								Yes:       3,
								No:        4,
								Missing:   3,
								Leaf:      0,
								Children: []*node{
									{
										ID:        3,
										Depth:     2,
										Feature:   "g",
										Threshold: 98,
										Yes:       7,
										No:        8,
										Missing:   7,
										Leaf:      0,
										Children: []*node{
											{
												ID:   7,
												Leaf: -0,
											},
											{
												ID:   8,
												Leaf: -0.666666687,
											},
										},
									},
									{
										ID:        4,
										Depth:     2,
										Feature:   "b",
										Threshold: 112,
										Yes:       9,
										No:        10,
										Missing:   9,
										Leaf:      0,
										Children: []*node{
											{
												ID:   9,
												Leaf: 0.70588237,
											},
											{
												ID:   10,
												Leaf: 0.176470593,
											},
										},
									},
								},
							},
							{
								ID:        2,
								Depth:     1,
								Feature:   "r",
								Threshold: 247,
								Yes:       5,
								No:        6,
								Missing:   5,
								Leaf:      0,
								Children: []*node{
									{
										ID:   5,
										Leaf: -0.826086938,
									},
									{
										ID:        6,
										Depth:     2,
										Feature:   "g",
										Threshold: 221,
										Yes:       11,
										No:        12,
										Missing:   11,
										Leaf:      0,
										Children: []*node{
											{
												ID:   11,
												Leaf: -0.200000003,
											},
											{
												ID:   12,
												Leaf: 0.25,
											},
										},
									},
								},
							},
						},
					},
				},
				BinaryThreshold: 0.5,
			},
			wantErr: false,
		},
		{
			name: "unhappy path: corrupt json",
			args: args{
				data: []byte(`[ "nodeid": 0, "depth": 0, "split": "b", "split_condition": 162, "yes": 1, "no": 2, "missing": 1 , "children": [
    { "nodeid": 1, "depth": 1, "split": "r", "split_condition": 86, "yes": 3, "no": 4, "missing": 3 , "children": [
      { "nodeid": 3, "depth": 2, "split": "g", "split_condition": 98, "yes": 7, "no": 8, "missing": 7 , "children": [
        { "nodeid": 7, "leaf": -0 }, 
        { "nodeid": 8, "leaf": -0.666666687 }
      ]}, 
      { "nodeid": 4, "depth": 2, "split": "b", "split_condition": 112, "yes": 9, "no": 10, "missing": 9 , "children": [
        { "nodeid": 9, "leaf": 0.70588237 }, 
        { "nodeid": 10, "leaf": 0.176470593 }
      ]}
    ]}, 
    { "nodeid": 2, "depth": 1, "split": "r", "split_condition": 247, "yes": 5, "no": 6, "missing": 5 , "children": [
      { "nodeid": 5, "leaf": -0.826086938 }, 
      { "nodeid": 6, "depth": 2, "split": "g", "split_condition": 221, "yes": 11, "no": 12, "missing": 11 , "children": [
        { "nodeid": 11, "leaf": -0.200000003 }, 
        { "nodeid": 12, "leaf": 0.25 }
      ]}
    ]}
  ]}]`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := LoadModelConfig(tt.args.data)
				if (err != nil) != tt.wantErr {
					t.Errorf("LoadModelConfig() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("LoadModelConfig() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestModel_Predict(t *testing.T) {
	type fields struct {
		trees           XGB
		BinaryThreshold float64
	}
	type args struct {
		dataFrame SparceMatrix
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []bool
		wantErr bool
	}{
		{
			name: "happy path",
			fields: fields{
				trees: XGB{
					{
						ID:        0,
						Depth:     0,
						Feature:   "b",
						Threshold: 162,
						Yes:       1,
						No:        2,
						Missing:   1,
						Leaf:      0,
						Children: []*node{
							{
								ID:        1,
								Depth:     1,
								Feature:   "r",
								Threshold: 86,
								Yes:       3,
								No:        4,
								Missing:   3,
								Leaf:      0,
								Children: []*node{
									{
										ID:        3,
										Depth:     2,
										Feature:   "g",
										Threshold: 98,
										Yes:       7,
										No:        8,
										Missing:   7,
										Leaf:      0,
										Children: []*node{
											{
												ID:   7,
												Leaf: -0,
											},
											{
												ID:   8,
												Leaf: -0.666666687,
											},
										},
									},
									{
										ID:        4,
										Depth:     2,
										Feature:   "b",
										Threshold: 112,
										Yes:       9,
										No:        10,
										Missing:   9,
										Leaf:      0,
										Children: []*node{
											{
												ID:   9,
												Leaf: 0.70588237,
											},
											{
												ID:   10,
												Leaf: 0.176470593,
											},
										},
									},
								},
							},
							{
								ID:        2,
								Depth:     1,
								Feature:   "r",
								Threshold: 247,
								Yes:       5,
								No:        6,
								Missing:   5,
								Leaf:      0,
								Children: []*node{
									{
										ID:   5,
										Leaf: -0.826086938,
									},
									{
										ID:        6,
										Depth:     2,
										Feature:   "g",
										Threshold: 221,
										Yes:       11,
										No:        12,
										Missing:   11,
										Leaf:      0,
										Children: []*node{
											{
												ID:   11,
												Leaf: -0.200000003,
											},
											{
												ID:   12,
												Leaf: 0.25,
											},
										},
									},
								},
							},
						},
					},
				},
				BinaryThreshold: 0.5,
			},
			args: args{
				dataFrame: SparceMatrix{
					{"r": 254., "g": 37., "b": 0.},
					{"r": 182, "g": 221, "b": 199},
				},
			},
			want:    []bool{true, false},
			wantErr: false,
		},
		{
			name: "unhappy path: node missing",
			fields: fields{
				trees: XGB{
					{
						ID:        0,
						Depth:     0,
						Feature:   "b",
						Threshold: 162,
						Yes:       1,
						No:        2,
						Missing:   1,
						Leaf:      0,
						Children: []*node{
							{
								ID:        1,
								Depth:     1,
								Feature:   "r",
								Threshold: 86,
								Yes:       3,
								No:        4,
								Missing:   3,
								Leaf:      0,
								Children: []*node{
									{
										ID:        3,
										Depth:     2,
										Feature:   "g",
										Threshold: 98,
										Yes:       7,
										No:        8,
										Missing:   7,
										Leaf:      0,
										Children: []*node{
											{
												ID:   7,
												Leaf: -0,
											},
											{
												ID:   8,
												Leaf: -0.666666687,
											},
										},
									},
								},
							},
							{
								ID:        2,
								Depth:     1,
								Feature:   "r",
								Threshold: 247,
								Yes:       5,
								No:        6,
								Missing:   5,
								Leaf:      0,
								Children: []*node{
									{
										ID:   5,
										Leaf: -0.826086938,
									},
									{
										ID:        6,
										Depth:     2,
										Feature:   "g",
										Threshold: 221,
										Yes:       11,
										No:        12,
										Missing:   11,
										Leaf:      0,
										Children: []*node{
											{
												ID:   11,
												Leaf: -0.200000003,
											},
											{
												ID:   12,
												Leaf: 0.25,
											},
										},
									},
								},
							},
						},
					},
				},
				BinaryThreshold: 0.5,
			},
			args: args{
				dataFrame: SparceMatrix{
					{"r": 254., "g": 37., "b": 0.},
					{"r": 182, "g": 221, "b": 199},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy path: missing feature - fall back",
			fields: fields{
				trees: XGB{
					{
						ID:        0,
						Depth:     0,
						Feature:   "foo",
						Threshold: 162,
						Yes:       1,
						No:        2,
						Missing:   1,
						Leaf:      0,
						Children: []*node{
							{
								ID:        1,
								Depth:     1,
								Feature:   "r",
								Threshold: 86,
								Yes:       3,
								No:        4,
								Missing:   3,
								Leaf:      0,
								Children: []*node{
									{
										ID:        3,
										Depth:     2,
										Feature:   "g",
										Threshold: 98,
										Yes:       7,
										No:        8,
										Missing:   7,
										Leaf:      0,
										Children: []*node{
											{
												ID:   7,
												Leaf: -0,
											},
											{
												ID:   8,
												Leaf: -0.666666687,
											},
										},
									},
									{
										ID:        4,
										Depth:     2,
										Feature:   "b",
										Threshold: 112,
										Yes:       9,
										No:        10,
										Missing:   9,
										Leaf:      0,
										Children: []*node{
											{
												ID:   9,
												Leaf: 0.70588237,
											},
											{
												ID:   10,
												Leaf: 0.176470593,
											},
										},
									},
								},
							},
							{
								ID:        2,
								Depth:     1,
								Feature:   "r",
								Threshold: 247,
								Yes:       5,
								No:        6,
								Missing:   5,
								Leaf:      0,
								Children: []*node{
									{
										ID:   5,
										Leaf: -0.826086938,
									},
									{
										ID:        6,
										Depth:     2,
										Feature:   "g",
										Threshold: 221,
										Yes:       11,
										No:        12,
										Missing:   11,
										Leaf:      0,
										Children: []*node{
											{
												ID:   11,
												Leaf: -0.200000003,
											},
											{
												ID:   12,
												Leaf: 0.25,
											},
										},
									},
								},
							},
						},
					},
				},
				BinaryThreshold: 0.5,
			},
			args: args{
				dataFrame: SparceMatrix{
					{"r": 254., "g": 37., "b": 0.},
					{"r": 182, "g": 221, "b": 199},
				},
			},
			want:    []bool{true, true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				m := &Model{
					trees:           tt.fields.trees,
					BinaryThreshold: tt.fields.BinaryThreshold,
				}
				got, err := m.Predict(tt.args.dataFrame)
				if (err != nil) != tt.wantErr {
					t.Errorf("Predict() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Predict() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
