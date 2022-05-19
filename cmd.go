package main

import (
	"fmt"
	"io"
	"os"
)

// RenderTemplateCmd allows rendering templates
type RenderTemplateCmd struct {
	DataFile string `short:"f" long:"data-file" description:"Properties file containing the replacements for the template" value-name:"DATA_FILE"`
	Args     struct {
		TemplateFile string `positional-arg-name:"template-file" description:"File containing the template to render. Its contents can be also passed through stdin"`
	} `positional-args:"yes"`
	OutWriter io.Writer
}

// NewRenderTemplateCmd returns a new RenderTemplateCmd
func NewRenderTemplateCmd() *RenderTemplateCmd {
	return &RenderTemplateCmd{OutWriter: os.Stdout}
}

func (c *RenderTemplateCmd) getTemplateData() (*templateData, error) {
	d := newTemplateData()

	if c.DataFile != "" {
		if err := d.LoadBatchDataFile(c.DataFile); err != nil {
			return nil, fmt.Errorf("cannot read template data file %q: %v", c.DataFile, err)
		}
	}
	return d, nil
}

// Execute performs rendering
func (c *RenderTemplateCmd) Execute(args []string) (err error) {
	var in io.Reader
	if c.Args.TemplateFile != "" {
		fh, err := os.Open(c.Args.TemplateFile)
		if err != nil {
			return fmt.Errorf("cannot open template file: %v", err)
		}
		defer fh.Close()
		in = fh
	} else if hasPipedStdin() {
		in = os.Stdin
	} else {
		return fmt.Errorf("you must provide a template file as an argument or through stdin")
	}

	data, err := c.getTemplateData()
	if err != nil {
		return err
	}
	renderer := newHandlerbarsRenderer()
	str, err := renderer.RenderTemplate(in, data)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(c.OutWriter, str)
	return err
}

func hasPipedStdin() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
