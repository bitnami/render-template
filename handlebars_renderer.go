package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/aymerick/raymond"
)

func orHelper(a interface{}, b interface{}, options *raymond.Options) interface{} {
	aStr := raymond.Str(a)
	bStr := raymond.Str(b)

	if aStr != "" {
		return aStr
	} else if bStr != "" {
		return bStr
	}
	return ""
}

func quoteHelper(s string) raymond.SafeString {
	return raymond.SafeString(fmt.Sprintf("%q", s))
}

func jsonEscapeHelper(s string) raymond.SafeString {
	b, err := json.Marshal(s)
	// raymond does not allow returning errors from helpers
	// TODO: Setup recovering from the panic when calling Exec
	if err != nil {
		panic(err)
	}
	return raymond.SafeString(b)
}

func init() {
	raymond.RegisterHelper("or", orHelper)
	raymond.RegisterHelper("quote", quoteHelper)
	raymond.RegisterHelper("json_escape", jsonEscapeHelper)
}

type renderer interface {
	RenderTemplate(in io.Reader, data *templateData) (string, error)
}

type handlerbarsRenderer struct {
}

func newHandlerbarsRenderer() renderer {
	return &handlerbarsRenderer{}
}

func (h *handlerbarsRenderer) RenderTemplate(in io.Reader, data *templateData) (string, error) {
	tpl, err := h.parseTemplate(in)
	if err != nil {
		return "", fmt.Errorf("cannot parse template: %v", err)
	}
	parsedData := h.convertTemplateData(data)

	return tpl.Exec(parsedData)
}

func (h *handlerbarsRenderer) parseTemplate(in io.Reader) (*raymond.Template, error) {
	b, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}
	return raymond.Parse(string(b))
}

// strings chars are escaped as they are considerd unsafe for HTML. We need to mark them as safe
func (h *handlerbarsRenderer) convertTemplateData(data *templateData) map[string]raymond.SafeString {
	rawData := data.Data()
	res := make(map[string]raymond.SafeString, len(rawData))
	for key, value := range rawData {
		res[key] = raymond.SafeString(value)
	}
	return res
}
