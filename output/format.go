package output

import (
	"fmt"
	"io"
	"text/template"

	"github.com/mitchellh/pointerstructure"
)

// Config configures the behaviour of the output formatter.
type Config struct {
	// Format must contain the name of the formatter that should be used.
	Format string
	// Formatters can be set to configure user-defined formatters. If empty,
	// the built-in formatters are used.
	Formatters FormatterMap
	// Template configures the template for template-based formatters.
	Template string
	// TemplateConfig is optional configuration for the underlying template
	// struct for template-based formatter.
	// See: https://golang.org/pkg/text/template/#Template
	TemplateConfig struct {
		// Funcs configures additional template funcs.
		// See: https://golang.org/pkg/text/template/#Template.Funcs
		Funcs template.FuncMap
		// Options configures additional template options.
		// See: https://golang.org/pkg/text/template/#Template.Option
		Options []string
	}
	// TemplateItems configures the behaviour of template based formatters. If
	// true the template applies only to the items of a slice that should be
	// formatted. This allows for omission of the `range` loop in templates.
	// This config only applies to slices. Has no effect on non-slice types.
	TemplateItems bool
	// TrailingNewline ensures that the formatted output always ends with a
	// trailing newline character if it is not empty. Does not add a newline if
	// the formatted output is already terminated by one.
	TrailingNewline bool
	// JSONPointer can be fill with a jsonpointer expression to select a nested
	// field within an object before passing it on to the formatter. If empty
	// the original object is passed as is.
	// See RFC: https://datatracker.ietf.org/doc/html/rfc6901
	JSONPointer string
}

// Format formats v using the given config and writes the result to w. Returns
// any error that may occur during formatting. On errors nothing is written to
// w.
func Format(w io.Writer, v interface{}, config *Config) error {
	buf, err := FormatBytes(v, config)
	if err != nil {
		return err
	}

	_, err = w.Write(buf)
	return err
}

// FormatBytes formats v using the given config and returns the formatted
// bytes. Returns any error that may occur during formatting.
func FormatBytes(v interface{}, config *Config) ([]byte, error) {
	formatters := config.Formatters
	if len(formatters) == 0 {
		formatters = Formatters
	}

	f, ok := formatters[config.Format]
	if !ok {
		return nil, fmt.Errorf("no formatter for format %q", config.Format)
	}

	if config.JSONPointer != "" {
		pv, err := pointerstructure.Get(v, config.JSONPointer)
		if err != nil {
			return nil, err
		}

		v = pv
	}

	buf, err := f.Format(v, config)
	if err != nil {
		return nil, err
	}

	// Append a trailing newline if requested, but only if the formatter did
	// not already do that for us.
	if config.TrailingNewline && len(buf) > 0 && buf[len(buf)-1] != '\n' {
		buf = append(buf, '\n')
	}

	return buf, nil
}

// FormatString formats v using the given config and returns the formatted
// string. Returns any error that may occur during formatting.
func FormatString(v interface{}, config *Config) (string, error) {
	buf, err := FormatBytes(v, config)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
