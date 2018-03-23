package main

import (
	"flag"
	"os"
	"testing"

	tu "github.com/bitnami/gonit/testutils"
	ca "github.com/juamedgod/cliassert"
	"github.com/stretchr/testify/require"
)

func Test_main(t *testing.T) {
	templateData1 := "hello, {{name}}.\nWelcome to {{company}}"
	expectedResult1 := "hello, user.\nWelcome to bitnami"
	tests := []struct {
		name    string
		tplText string
		tplFile string

		dataText       string
		env            map[string]string
		wantErr        bool
		stdin          string
		expectedErr    interface{}
		expectedResult string
	}{
		{
			name:           "Template from stdin",
			stdin:          templateData1,
			env:            map[string]string{"name": "user", "company": "bitnami"},
			expectedResult: expectedResult1,
		},
		{
			name:        "No template from file nor stdin",
			wantErr:     true,
			expectedErr: "you must provide a template file as an argument or through stdin",
		},

		{
			name:        "Invalid template file",
			wantErr:     true,
			tplFile:     "some_nonexistent_file.txt",
			expectedErr: "cannot open template file",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := tu.NewSandbox()
			defer sb.Cleanup()
			cb := replaceEnv(tt.env)
			defer cb()

			args := []string{}

			if len(tt.dataText) > 0 {
				dataFile, err := sb.Write("data.txt", tt.dataText)
				require.NoError(t, err)
				args = append(args, "-f", dataFile)
			}

			if tt.stdin == "" {
				if tt.tplText != "" {
					file, err := sb.Write("template.conf.tpl", tt.tplText)
					require.NoError(t, err)
					args = append(args, file)
				} else if tt.tplFile != "" {
					args = append(args, tt.tplFile)

				}
			}
			res := runTool(os.Args[0], args, tt.stdin)
			if tt.wantErr {
				if res.Success() {
					t.Errorf("the command was expected to fail but succeeded")
				} else if tt.expectedErr != nil {
					res.AssertErrorMatch(t, tt.expectedErr)
				}
			} else {
				res.AssertSuccessMatch(t, tt.expectedResult)
			}
		})
	}
}

func TestMain(m *testing.M) {
	if os.Getenv("BE_TOOL") == "1" {
		main()
		os.Exit(0)
		return
	}
	flag.Parse()
	c := m.Run()
	os.Exit(c)
}

func runTool(bin string, args []string, stdin string) ca.CmdResult {
	cmd := ca.NewCommand()
	if stdin != "" {
		cmd.SetStdin(stdin)
	}
	os.Setenv("BE_TOOL", "1")
	defer os.Unsetenv("BE_TOOL")
	return cmd.Exec(bin, args...)
}

func RunTool(args ...string) ca.CmdResult {
	return runTool(os.Args[0], args, "")
}
