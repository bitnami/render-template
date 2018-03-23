package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

/*
func getCliArguments(port int, state string, host string, timeout int) []string {
	return []string{
		"--state", state,
		"--host", host,
		"--timeout", fmt.Sprintf("%d", timeout),
		fmt.Sprintf("%d", port),
	}
}

func testViaStruct(port int, state string, host string, timeout int, t *testing.T) error {
	cmd := NewWaitForPortCmd()
	cmd.State = state
	cmd.Host = host
	cmd.Timeout = timeout
	cmd.Args.Port = port

	return cmd.Execute([]string{})
}

func testViaCli(port int, state string, host string, timeout int, t *testing.T) error {
	cliArgs := getCliArguments(port, state, host, timeout)
	res := RunTool(cliArgs...)
	if !res.Success() {
		return fmt.Errorf("%s", res.Stderr())
	}
	return nil
}
*/

// setenv modify the environment with the new provided map and
// returns a callback to set it back to the original state
func setenv(env map[string]string) (restoreCn func()) {
	toRestore := make(map[string]string)
	toUnset := []string{}
	for k, v := range env {
		oldValue, ok := os.LookupEnv(k)
		if ok {
			toRestore[k] = oldValue
		} else {
			toUnset = append(toUnset, k)
		}
		os.Setenv(k, v)
	}
	return func() {
		for k, v := range toRestore {
			os.Setenv(k, v)
		}
		for _, k := range toUnset {
			os.Unsetenv(k)
		}
	}
}
func replaceEnv(env map[string]string) (restoreCn func()) {
	toRestore := make(map[string]string)
	for _, envLine := range os.Environ() {
		data := strings.SplitN(envLine, "=", 2)
		if len(data) != 2 {
			continue
		}
		key := data[0]
		value := data[1]
		toRestore[key] = value
	}
	os.Clearenv()
	for k, v := range env {
		os.Setenv(k, v)
	}
	return func() {
		for k, v := range toRestore {
			os.Setenv(k, v)
		}
	}
}

func mergeMaps(args ...map[string]string) map[string]string {
	new := make(map[string]string)
	for _, m := range args {
		for k, v := range m {
			new[k] = v
		}
	}
	return new
}
func writeDataFile(file string, data map[string]string) error {
	fh, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("cannot open file %q: %v", file, err)
	}
	defer fh.Close()
	for k, v := range data {
		if _, err := io.WriteString(fh, fmt.Sprintf("%s=%s\n", k, v)); err != nil {
			return fmt.Errorf("cannot write data to file: %v", err)
		}

	}
	return nil
}
