package main

import (
	"reflect"
	"regexp"
	"testing"

	tu "github.com/bitnami/gonit/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type templateDataTest struct {
	name           string
	inputEnv       map[string]string
	dataText       string
	expectedResult map[string]string
	wantErr        bool
	expectedErr    *regexp.Regexp
}

func Test_templateData_Data(t *testing.T) {
	env1 := map[string]string{"A": "Value1", "a": "another value", "INSTALLDIR": "/opt/bitnami"}
	tests := []templateDataTest{

		{
			name:           "Only env",
			inputEnv:       env1,
			expectedResult: env1,
		},
		{
			name:     "Env and data",
			inputEnv: env1,
			dataText: `
USER=bitnami
HOME=/home/bitnami
`,
			expectedResult: mergeMaps(env1, map[string]string{"USER": "bitnami", "HOME": "/home/bitnami"}),
		},
		{
			name: "Data with comments and malformed",
			dataText: `
A=Value a
#B=Value b

C=Value c

Some value

D=Value d

			`,
			expectedResult: map[string]string{"A": "Value a", "C": "Value c", "D": "Value d"},
			wantErr:        true,
			expectedErr:    regexp.MustCompile(`^malformed line "Some value": could not find '=' separator$`),
		},
		{
			name:           "Data overrides env",
			inputEnv:       map[string]string{"A": "value a", "B": "value b"},
			expectedResult: map[string]string{"A": "value a", "B": "new value b", "C": "value c"},
			dataText:       "B=new value b\nC=value c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			sb := tu.NewSandbox()
			defer sb.Cleanup()
			file, err := sb.Write("data.txt", tt.dataText)
			require.NoError(t, err)

			cb := replaceEnv(tt.inputEnv)
			defer cb()
			d := newTemplateData()

			err = d.LoadBatchDataFile(file)
			if tt.wantErr {
				if tt.expectedErr != nil {
					tu.AssertErrorMatch(t, err, tt.expectedErr)
				} else {
					assert.Error(t, err)
				}
			} else {
				assert.NoError(t, err)
			}
			res := d.Data()
			assert.True(t, reflect.DeepEqual(tt.expectedResult, res), "expected %v to be equal to %v", res, tt.expectedResult)
		})
	}
}
