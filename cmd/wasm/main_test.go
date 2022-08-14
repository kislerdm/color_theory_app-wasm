//go:build js
// +build js

package main

import (
	"syscall/js"
	"testing"
)

func Test_run(t *testing.T) {
	type args struct {
		args []js.Value
	}
	tests := []struct {
		name     string
		args     args
		wantHtml string
		wantErr  bool
	}{
		{
			name: "happy path: Black",
			args: args{
				[]js.Value{js.ValueOf(0), js.ValueOf(0), js.ValueOf(0)},
			},
			wantHtml: `<div><label for="output_name" id="output_label">Color Name:</label><output name="color_name" id="output_name"> Black</output></div><div><label for="output_type" id="output_label">Color Type:</label><output name="color_type" id="output_type"> Cool</output></div>`,
			wantErr:  false,
		},
		{
			name: "happy path: Red",
			args: args{
				[]js.Value{js.ValueOf(255), js.ValueOf(0), js.ValueOf(0)},
			},
			wantHtml: `<div><label for="output_name" id="output_label">Color Name:</label><output name="color_name" id="output_name"> Red</output></div><div><label for="output_type" id="output_label">Color Type:</label><output name="color_type" id="output_type"> Warm</output></div>`,
			wantErr:  false,
		},
		{
			name:    "unhappy path: not enough input args",
			args:    args{[]js.Value{js.ValueOf(255)}},
			wantErr: true,
		},
		{
			name:    "unhappy path: wrong input enough input args",
			args:    args{[]js.Value{js.ValueOf(-1), js.ValueOf(0), js.ValueOf(0)}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gotHtml, err := run(tt.args.args)
				if (err != nil) != tt.wantErr {
					t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if gotHtml != tt.wantHtml {
					t.Errorf("run() gotHtml = %v, want %v", gotHtml, tt.wantHtml)
				}
			},
		)
	}
}

func Test_ui(t *testing.T) {
	type args struct {
		colorName string
		isWarm    bool
	}
	tests := []struct {
		name     string
		args     args
		wantHtml string
	}{
		{
			name: "color not found",
			args: args{
				colorName: "",
				isWarm:    false,
			},
			wantHtml: `<div><label for="output_name" id="output_label">Color Name:</label><output name="color_name" id="output_name"> Not found</output></div><div><label for="output_type" id="output_label">Color Type:</label><output name="color_type" id="output_type"> Cool</output></div>`,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if gotHtml := ui(tt.args.colorName, tt.args.isWarm); gotHtml != tt.wantHtml {
					t.Errorf("ui() = %v, want %v", gotHtml, tt.wantHtml)
				}
			},
		)
	}
}
