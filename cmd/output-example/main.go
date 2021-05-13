// Example application for demonstating the output package.
package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/martinohmann/exit"
	"github.com/martinohmann/exp/cli"
	"github.com/martinohmann/exp/output"
	"github.com/mgutz/ansi"
	"github.com/spf13/pflag"
)

func main() {
	cli.Run(run)
}

func run() error {
	// Register custom format func.
	output.RegisterFormatFunc("xml", func(v interface{}, config *output.Config) ([]byte, error) {
		return xml.Marshal(v)
	})

	config := &output.Config{
		Format:   "json",
		Template: `{{color "cyan"}}{{.}}{{color "reset"}}`,
		TemplateConfig: output.TemplateConfig{
			// Custom template func definitions.
			Funcs: template.FuncMap{
				"color": ansi.ColorCode,
			},
		},
	}

	args, err := parseArgs(config)
	if err != nil {
		return exit.Error(exit.CodeUsage, err)
	}

	rc, err := readCloser(args)
	if err != nil {
		return err
	}
	defer rc.Close()

	var obj interface{}

	if err := decode(rc, &obj); err != nil {
		return err
	}

	// This is doing the actual work in this example.
	return output.Format(os.Stdout, obj, config)
}

func parseArgs(config *output.Config) ([]string, error) {
	fs := pflag.NewFlagSet("output-example", pflag.ContinueOnError)

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: output-example [<file>] [flags]

Accepts json on stdin or from a file and prints it to stdout in a format dictated by the provided flags.

This is an example app for playing around with the github.com/martinohmann/exp/output package.

Flags:`)
		fs.PrintDefaults()
	}

	fs.StringVarP(&config.Format, "output", "o", config.Format, fmt.Sprintf("output format, valid values: '%s'", strings.Join(output.FormatterNames(), "', '")))
	fs.StringVarP(&config.Template, "template", "t", config.Template, "output template. ignored unless output format is 'gotemplate'")
	fs.StringVarP(&config.JSONPointer, "jsonpointer", "j", config.JSONPointer, "json pointer for filtering the data before formatting, e.g. '/foo/0/bar'")
	fs.BoolVar(&config.TemplateItems, "items", config.TemplateItems, "if true, the template applies to the items if the input is a slice. ignored unless output format is 'gotemplate'")
	fs.BoolVar(&config.TrailingNewline, "newline", config.TrailingNewline, "ensure output ends with a trailing newline")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	return fs.Args(), nil
}

func readCloser(args []string) (io.ReadCloser, error) {
	if len(args) == 0 || args[0] == "-" {
		return io.NopCloser(bufio.NewReader(os.Stdin)), nil
	}

	return os.Open(args[0])
}

func decode(r io.Reader, obj interface{}) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(buf, obj)
}
