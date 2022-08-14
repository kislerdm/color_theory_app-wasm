package main

import (
	"reflect"
	"testing"
)

func TestOutput_generateUI(t *testing.T) {
	type fields struct {
		Name   string
		IsWarm bool
	}
	tests := []struct {
		name     string
		fields   fields
		wantHtml string
	}{
		{
			name: "happy path: cold",
			fields: fields{
				Name:   "foo",
				IsWarm: false,
			},
			wantHtml: `<div><label for="output_name" id="output_label">Color Name:</label><output name="color_name" id="output_name"> foo</output></div><div><label for="output_type" id="output_label">Color Type:</label><output name="color_type" id="output_type"> Cool</output></div>`,
		},
		{
			name: "happy path: warm",
			fields: fields{
				Name:   "bar",
				IsWarm: true,
			},
			wantHtml: `<div><label for="output_name" id="output_label">Color Name:</label><output name="color_name" id="output_name"> bar</output></div><div><label for="output_type" id="output_label">Color Type:</label><output name="color_type" id="output_type"> Warm</output></div>`,
		},
		{
			name: "unhappy path: name not found",
			fields: fields{
				Name:   "",
				IsWarm: true,
			},
			wantHtml: `<div><label for="output_name" id="output_label">Color Name:</label><output name="color_name" id="output_name"> Not found</output></div><div><label for="output_type" id="output_label">Color Type:</label><output name="color_type" id="output_type"> Warm</output></div>`,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				o := &Output{
					Name:   tt.fields.Name,
					IsWarm: tt.fields.IsWarm,
				}
				if gotHtml := o.generateUI(); gotHtml != tt.wantHtml {
					t.Errorf("generateUI() = %v, want %v", gotHtml, tt.wantHtml)
				}
			},
		)
	}
}

func Test_generateOutput(t *testing.T) {
	type args struct {
		r float64
		g float64
		b float64
	}
	tests := []struct {
		name    string
		args    args
		want    Output
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{187, 58, 58},
			want: Output{
				Name:   "Well Read",
				IsWarm: true,
			},
			wantErr: false,
		},
		{
			name:    "unhappy path: wrong input",
			args:    args{-1, 58, 58},
			want:    Output{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := generateOutput(tt.args.r, tt.args.g, tt.args.b)
				if (err != nil) != tt.wantErr {
					t.Errorf("generateOutput() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("generateOutput() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
