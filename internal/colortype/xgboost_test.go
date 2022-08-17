package colortype

import (
	"reflect"
	"testing"
)

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
