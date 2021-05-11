package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/ghodss/yaml"
)

// Formatter can format values.
type Formatter interface {
	// Format takes v and returns the formatted bytes. The Formatter may change
	// behaviour depending on the passed in config.
	Format(v interface{}, config *Config) ([]byte, error)
}

// FormatFunc is a func that implements the Formatter interface.
type FormatFunc func(v interface{}, config *Config) ([]byte, error)

// Format implements the Formatter interface.
func (f FormatFunc) Format(v interface{}, config *Config) ([]byte, error) {
	return f(v, config)
}

// FormatterMap maps a user-defined name to a formatter.
type FormatterMap map[string]Formatter

// Register registers a Formatter. Panics if a formatter with the same name
// already exists.
func (m FormatterMap) Register(name string, f Formatter) {
	if _, exists := m[name]; exists {
		panic(fmt.Sprintf("formatter with name %q already registered", name))
	}

	m[name] = f
}

// Names returns a sorted slice of formatter names. This is useful to present
// allowed values in command line flags.
func (m FormatterMap) Names() []string {
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}

// Formatters is a map of built-in formatters which can be extended to make
// more formatters globally available. These are used as a fallback if the user
// does not provide a custom map as part of the config to Format* funcs.
var Formatters = FormatterMap{
	"json": FormatFunc(func(v interface{}, config *Config) ([]byte, error) {
		return json.MarshalIndent(v, "", "  ")
	}),
	"yaml": FormatFunc(func(v interface{}, config *Config) ([]byte, error) {
		return yaml.Marshal(v)
	}),
	"gotemplate": FormatFunc(func(v interface{}, config *Config) ([]byte, error) {
		var buf bytes.Buffer

		if err := formatTemplate(&buf, v, config); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}),
}
