package main

import (
	"bytes"
	"testing"

	tu "github.com/bitnami/gonit/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tplTest struct {
	tplText        string
	env            map[string]string
	data           map[string]string
	expectedResult string
}

func TestRenderTemplateCmd_Execute(t *testing.T) {
	tests := []struct {
		name        string
		test        tplTest
		wantErr     bool
		expectedErr interface{}
	}{
		{
			name: "Simple template",
			test: tplTest{
				tplText:        `hello {{to}}{{#if EXTRA}}and EXTRA{{/if}}`,
				env:            map[string]string{"to": "world"},
				data:           map[string]string{},
				expectedResult: `hello world`,
			},
		},
		{
			name: "Simple template with extra data",
			test: tplTest{
				tplText:        `hello {{to}}{{#if EXTRA}} and {{EXTRA}}{{/if}}`,
				env:            map[string]string{"to": "world"},
				data:           map[string]string{"EXTRA": "bitnami"},
				expectedResult: `hello world and bitnami`,
			},
		},
		{
			name: "Multiline Template",
			test: tplTest{
				tplText:        "This is a {{how_much}} complex\nmultiline {{what}}\n{{empty_value}}\nCreated to test {{program}}\n",
				env:            map[string]string{"how_much": "very", "what": "example", "program": "render-template"},
				expectedResult: "This is a very complex\nmultiline example\n\nCreated to test render-template\n",
			},
		},
		{
			name: "Quote helper",
			test: tplTest{
				tplText:        `Quoted {{quote string}}`,
				env:            map[string]string{"string": `this " has some "quoted" words ' and chars`},
				expectedResult: `Quoted "this \" has some \"quoted\" words ' and chars"`,
			},
		},
		{
			name: "or helper",
			test: tplTest{
				tplText:        `{{#if (or A B)}}A-OR-B{{/if}};{{#if (or B A)}}B-OR-A{{/if}};{{#if (or B C)}}EMPTY{{/if}};{{#if (or D E)}}UNDEFINED{{/if}};{{#if (or D A)}}D-OR-A{{/if}}`,
				env:            map[string]string{"A": "yes", "B": "", "C": ""},
				expectedResult: `A-OR-B;B-OR-A;;;D-OR-A`,
			},
		},
		{
			name: "json_escape helper",
			test: tplTest{
				tplText:        `val1={{json_escape VAL1}};val2={{json_escape VAL2}}`,
				env:            map[string]string{"VAL1": `t''"`, "VAL2": `hello world`},
				expectedResult: `val1="t''\"";val2="hello world"`,
			},
		},
		{
			name: "Malformed template",
			test: tplTest{
				tplText: `{{if `,
			},
			wantErr:     true,
			expectedErr: "cannot parse template: Parse error on line 1",
		},
	}
	for _, tt := range tests {
		sb := tu.NewSandbox()
		defer sb.Cleanup()
		file, err := sb.Write("my.conf.tpl", tt.test.tplText)
		require.NoError(t, err)
		cmd := NewRenderTemplateCmd()
		cmd.Args.TemplateFile = file

		if len(tt.test.data) > 0 {
			dataFile := sb.TempFile("data.txt")
			writeDataFile(dataFile, tt.test.data)
			cmd.DataFile = dataFile
		}

		t.Run(tt.name, func(t *testing.T) {
			var err error
			cb := setenv(tt.test.env)
			defer cb()
			b := &bytes.Buffer{}
			cmd.OutWriter = b

			err = cmd.Execute([]string{})
			stdout := b.String()

			if (err != nil) != tt.wantErr {
				t.Errorf("RenderTemplateCmd.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if tt.expectedErr != nil {
					assert.Regexp(t, tt.expectedErr, err, "Expected error %v to match %v", err, tt.expectedErr)
				}
			} else {
				assert.Equal(t, tt.test.expectedResult, stdout)
			}
		})
	}
}

func TestRenderTemplateCmd_Execute_Errors(t *testing.T) {
	sb := tu.NewSandbox()
	defer sb.Cleanup()
	existentDataFile := sb.Touch("exitent-data.txt")
	nonExistentDataFile := sb.Normalize("non-exitent-data.txt")

	existentTemplate := sb.Touch("exitent-template.tpl")
	nonExistentTemplate := sb.Normalize("non-exitent-template.tpl")

	tests := []struct {
		name        string
		expectedErr interface{}
		tplFile     string
		dataFile    string
		// TODO: provide stdin
		stdin string
	}{
		{
			name:        "Template file does not exists",
			tplFile:     nonExistentTemplate,
			dataFile:    existentDataFile,
			expectedErr: "cannot open template file",
		},
		{
			name:        "Data file does not exists",
			dataFile:    nonExistentDataFile,
			tplFile:     existentTemplate,
			expectedErr: "cannot read template data file",
		},
	}
	for _, tt := range tests {
		cmd := NewRenderTemplateCmd()
		cmd.Args.TemplateFile = tt.tplFile
		cmd.DataFile = tt.dataFile
		t.Run(tt.name, func(t *testing.T) {

			err := cmd.Execute([]string{})

			if err == nil {
				t.Errorf("RenderTemplateCmd.Execute() was expected to fail")
			} else {
				if tt.expectedErr != nil {
					assert.Regexp(t, tt.expectedErr, err, "Expected error %v to match %v", err, tt.expectedErr)
				}
			}
		})
	}
}
