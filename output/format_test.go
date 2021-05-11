package output

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatString(t *testing.T) {
	tests := []formatTestCase{
		{
			cfg: Config{Format: "none"},
			err: errors.New(`no formatter for format "none"`),
		},
		{
			name: "unserializable type",
			cfg:  Config{Format: "json"},
			v:    func() {},
			err:  errors.New(`json: unsupported type: func()`),
		},
		{
			cfg:  Config{Format: "json"},
			want: "null",
		},
		{
			cfg:  Config{Format: "json"},
			v:    map[string]interface{}{"foo": "bar"},
			want: "{\n  \"foo\": \"bar\"\n}",
		},
		{
			cfg:  Config{Format: "yaml"},
			v:    map[string]interface{}{"foo": "bar"},
			want: "foo: bar\n",
		},
		{
			name: "ensures trailing newline",
			cfg:  Config{Format: "json", TrailingNewline: true},
			v:    map[string]interface{}{"foo": "bar"},
			want: "{\n  \"foo\": \"bar\"\n}\n",
		},
		{
			name: "evaluates json pointer",
			cfg:  Config{Format: "json", JSONPointer: "/foo"},
			v:    map[string]interface{}{"foo": "bar"},
			want: "\"bar\"",
		},
		{
			name: "does not add trailing newline if there already is one",
			cfg:  Config{Format: "yaml", TrailingNewline: true},
			v:    map[string]interface{}{"foo": "bar"},
			want: "foo: bar\n",
		},
		{
			name: "gotemplate requires template",
			cfg:  Config{Format: "gotemplate"},
			v:    map[string]interface{}{"foo": "bar", "baz": 42},
			err:  errors.New("template must not be empty"),
		},
		{
			name: "gotemplate",
			cfg:  Config{Format: "gotemplate", Template: "{{.baz}}"},
			v:    map[string]interface{}{"foo": "bar", "baz": 42},
			want: "42",
		},
		{
			name: "gotemplate slice",
			cfg:  Config{Format: "gotemplate", Template: "{{range $i, $v := .}}{{$i}}: {{$v}}\n{{end}}"},
			v:    []interface{}{"foo", "bar", "baz", 42},
			want: "0: foo\n1: bar\n2: baz\n3: 42\n",
		},
		{
			name: "gotemplate template items slice",
			cfg:  Config{Format: "gotemplate", Template: "{{.}}", TemplateItems: true},
			v:    []interface{}{"foo", "bar", "baz", 42},
			want: "foo\nbar\nbaz\n42",
		},
		{
			name: "custom formatter",
			cfg: Config{
				Format: "customjson",
				Formatters: FormatterMap{
					"customjson": FormatFunc(func(v interface{}, config *Config) ([]byte, error) {
						buf, err := json.MarshalIndent(v, "prefix:", "  ")
						if err != nil {
							return nil, err
						}

						return append([]byte(`prefix:`), buf...), nil
					}),
				},
			},
			v:    map[string]interface{}{"foo": "bar"},
			want: "prefix:{\nprefix:  \"foo\": \"bar\"\nprefix:}",
		},
	}

	testFormat(t, tests, FormatString)
}

type formatTestCase struct {
	name string
	cfg  Config
	v    interface{}
	want string
	err  error
}

type formatTestFunc func(v interface{}, config *Config) (string, error)

func testFormat(t *testing.T, tests []formatTestCase, fn formatTestFunc) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := fn(test.v, &test.cfg)
			if test.err != nil {
				require.EqualError(t, err, test.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.want, got)
			}
		})
	}
}
