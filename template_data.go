package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mmikulicic/multierror"
)

// TemplateData provides different sources of data to use when rendering the templates
type templateData struct {
	data    map[string]string
	envData map[string]string
}

func newTemplateData() *templateData {
	return &templateData{
		data: make(map[string]string),
	}
}

func (d *templateData) EnvData() map[string]string {
	if len(d.envData) == 0 {
		d.loadEnvVars()
	}
	return d.envData
}
func (d *templateData) CustomData() map[string]string {
	return d.data
}
func (d *templateData) Data() map[string]string {
	envData := d.EnvData()
	customData := d.CustomData()
	res := make(map[string]string, len(envData)+len(customData))
	for k, v := range envData {
		res[k] = v
	}
	// custom data takes precedence over the environment
	for k, v := range customData {
		res[k] = v
	}
	return res
}

func (d *templateData) AddData(k string, v string) {
	d.data[k] = v
}

func (d *templateData) LoadBatchDataFile(file string) error {
	fh, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("cannot read data file %q: %v", file, err)
	}
	defer fh.Close()
	return d.LoadBatchData(fh)
}

func (d *templateData) LoadBatchData(in io.Reader) error {
	var errs error
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Allow comments and blank lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			errs = multierror.Append(errs, fmt.Errorf("malformed line %q: could not find '=' separator", truncateString(line, 50)))
			continue
		}
		d.AddData(kv[0], kv[1])
	}
	// If the scanner failed just return that error instead of mixing with possible malformed lines
	if scanner.Err() != nil {
		return fmt.Errorf("cannot read template data: %v", scanner.Err())
	}
	return errs
}
func (d *templateData) loadEnvVars() {
	envStrs := os.Environ()
	d.envData = make(map[string]string, len(envStrs))
	for _, envLine := range envStrs {
		kv := strings.SplitN(envLine, "=", 2)
		if len(kv) != 2 {
			continue
		}
		d.envData[kv[0]] = kv[1]
	}
}

func truncateString(str string, num int) string {
	if len(str) <= num {
		return str
	}
	if len(str) > 3 {
		num -= 3
	}
	return str[0:num] + "..."
}
