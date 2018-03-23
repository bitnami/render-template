package main

import (
	"io"
	"reflect"
	"testing"

	"github.com/aymerick/raymond"
)

func Test_handlerbarsRenderer_RenderTemplate(t *testing.T) {
	type args struct {
		in   io.Reader
		data *templateData
	}
	tests := []struct {
		name    string
		h       *handlerbarsRenderer
		args    args
		want    string
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.RenderTemplate(tt.args.in, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("handlerbarsRenderer.RenderTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handlerbarsRenderer.RenderTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handlerbarsRenderer_parseTemplate(t *testing.T) {
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name    string
		h       *handlerbarsRenderer
		args    args
		want    *raymond.Template
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.parseTemplate(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("handlerbarsRenderer.parseTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handlerbarsRenderer.parseTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}
