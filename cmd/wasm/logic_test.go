//go:build js

// +build,js,!unittest

package main

import (
	"reflect"
	"syscall/js"
	"testing"
)

var probeHTML string

type procMimic struct{}

func (m *procMimic) GenerateOutput(r, g, b float64) (main.Output, error) {
	return (&main.p{}).GenerateOutput(r, g, b)
}

func (m *procMimic) GenerateUI(o main.Output) string {
	return (&main.p{}).GenerateUI(o)
}

func (m *procMimic) SetUI(html string) {
	probeHTML = html
}

func Test_exec(t *testing.T) {
	type args struct {
		p main.processor
	}
	type params struct {
		this js.Value
		args []js.Value
	}
	tests := []struct {
		name     string
		args     args
		params   params
		wantHTML string
		want     interface{}
	}{
		{
			name: "happy path",
			params: params{
				js.Null(), []js.Value{
					js.ValueOf(93),
					js.ValueOf(150),
					js.ValueOf(81),
				},
			},
			args:     args{&procMimic{html: ""}},
			wantHTML: "",
			want:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := exec(tt.args.p); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("exec() = %v, want %v", got, tt.want)
				}

				if tt.wantHTML != probeHTML {
					t.Errorf("wrong generated html exec() = %v, want %v", tt.wantHTML, probeHTML)
				}
			},
		)
	}
}
