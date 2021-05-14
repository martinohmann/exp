package pflagx

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runBindViperTestCases(t *testing.T) {
	type options struct {
		boolOpt        bool
		intOpt         int
		stringOpt      string
		stringSliceOpt []string
		intSliceOpt    []int
		stringToIntOpt map[string]int
	}

	testCases := []struct {
		name        string
		opts        options
		expected    options
		expectedErr error
		configureV  func(v *viper.Viper)
		env         map[string]string
		files       map[string]string
		args        []string
	}{
		{
			name: "options from ENV",
			configureV: func(v *viper.Viper) {
				v.SetEnvPrefix("snakes")
			},
			env:      map[string]string{"SNAKES_INT_OPT": "2"},
			opts:     options{intOpt: 42, stringOpt: "hello"},
			expected: options{intOpt: 2, stringOpt: "hello"},
		},
		{
			name:     "options from ENV, no envPrefix",
			env:      map[string]string{"INT_OPT": "2"},
			opts:     options{intOpt: 42, stringOpt: "hello"},
			expected: options{intOpt: 2, stringOpt: "hello"},
		},
		{
			name: "options from CLI > ENV",
			configureV: func(v *viper.Viper) {
				v.SetEnvPrefix("snakes")
			},
			args:     []string{"--int-opt", "13"},
			env:      map[string]string{"SNAKES_INT_OPT": "2"},
			opts:     options{intOpt: 42, stringOpt: "hello"},
			expected: options{intOpt: 13, stringOpt: "hello"},
		},
		{
			name: "options from config file",
			configureV: func(v *viper.Viper) {
				v.SetConfigName("snakes")
			},
			files: map[string]string{
				"snakes.json": `{"string-opt": "world", "bool-opt": true}`,
			},
			opts:     options{intOpt: 42, stringOpt: "hello"},
			expected: options{intOpt: 42, stringOpt: "world", boolOpt: true},
		},
		{
			name: "options from ENV > config file",
			configureV: func(v *viper.Viper) {
				v.SetConfigName("snakes")
				v.SetEnvPrefix("snakes")
			},
			env: map[string]string{"SNAKES_STRING_OPT": "all"},
			files: map[string]string{
				"snakes.json": `{"string-opt": "world", "bool-opt": true}`,
			},
			opts:     options{intOpt: 42, stringOpt: "hello"},
			expected: options{intOpt: 42, stringOpt: "all", boolOpt: true},
		},
		{
			name: "invalid config files cause errors",
			configureV: func(v *viper.Viper) {
				v.SetConfigName("snakes")
			},
			files: map[string]string{
				"snakes.yaml": "{invalid",
			},
			expectedErr: errors.New(`While parsing config: yaml: line 1: did not find expected ',' or '}'`),
		},
		{
			name: "nonexistent config files are ok",
			configureV: func(v *viper.Viper) {
				v.SetConfigName("snakes")
			},
			opts:     options{intOpt: 42},
			expected: options{intOpt: 42},
		},
		{
			name: "options from file, ENV and flags",
			configureV: func(v *viper.Viper) {
				v.SetConfigName("snakes")
				v.SetEnvPrefix("snakes")
			},
			args: []string{"--int-opt", "13"},
			env:  map[string]string{"SNAKES_STRING_OPT": "all", "SNAKES_INT_SLICE_OPT": "1,2"},
			files: map[string]string{
				"snakes.yaml": `---
string-opt: world
bool-opt: true
int-opt: 2`,
			},
			opts: options{
				intOpt:    42,
				stringOpt: "hello",
			},
			expected: options{
				boolOpt:     true,
				intOpt:      13,
				stringOpt:   "all",
				intSliceOpt: []int{1, 2},
			},
		},
		{
			name:        "options from ENV, invalid slice",
			env:         map[string]string{"INT_SLICE_OPT": "1,two"},
			expectedErr: errors.New(`failed to set flag from env or config: invalid argument "1,two" for "-I, --int-slice-opt" flag: strconv.Atoi: parsing "two": invalid syntax`),
		},
		{
			name:        "options from ENV, invalid map",
			env:         map[string]string{"STRING_TO_INT_OPT": "one=two"},
			expectedErr: errors.New(`failed to set flag from env or config: invalid argument "one=two" for "-M, --string-to-int-opt" flag: strconv.Atoi: parsing "two": invalid syntax`),
		},
		{
			name: "options from file, nested",
			configureV: func(v *viper.Viper) {
				v.SetConfigName("snakes")
			},
			files: map[string]string{
				"snakes.yaml": `---
string.opt: world
int:
  opt: 2
string-slice-opt: [foo, bar, 1]
int.slice:
  opt: [1, 2]
string-to-int-opt:
  foo: 5
  bar: 6`,
			},
			opts: options{
				intOpt:    42,
				stringOpt: "hello",
			},
			expected: options{
				intOpt:         2,
				stringOpt:      "world",
				stringSliceOpt: []string{"foo", "bar", "1"},
				intSliceOpt:    []int{1, 2},
				stringToIntOpt: map[string]int{"foo": 5, "bar": 6},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := viper.New()
			if tc.configureV != nil {
				tc.configureV(v)
			}

			opts := tc.opts

			for key, val := range tc.env {
				oldVal, exists := os.LookupEnv(key)
				if exists {
					defer os.Setenv(key, oldVal)
				} else {
					defer os.Unsetenv(key)
				}

				os.Setenv(key, val)
			}

			if len(tc.files) > 0 {
				dir := t.TempDir()

				v.AddConfigPath(dir)

				for name, content := range tc.files {
					err := ioutil.WriteFile(filepath.Join(dir, name), []byte(content), 0644)
					require.NoError(t, err)
				}
			}

			fs := pflag.NewFlagSet("snakes", pflag.ContinueOnError)

			fs.BoolVarP(&opts.boolOpt, "bool-opt", "b", opts.boolOpt, "A bool flag")
			fs.IntVarP(&opts.intOpt, "int-opt", "i", opts.intOpt, "An int flag")
			fs.StringVarP(&opts.stringOpt, "string-opt", "s", opts.stringOpt, "A string flag")
			fs.StringSliceVarP(&opts.stringSliceOpt, "string-slice-opt", "S", opts.stringSliceOpt, "A string slice flag")
			fs.IntSliceVarP(&opts.intSliceOpt, "int-slice-opt", "I", opts.intSliceOpt, "An int slice flag")
			fs.StringToIntVarP(&opts.stringToIntOpt, "string-to-int-opt", "M", opts.stringToIntOpt, "A string to int flag")

			require.NoError(t, fs.Parse(tc.args))

			err := BindViper(fs, v)
			if tc.expectedErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, opts)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			}
		})
	}
}

func TestBindViper(t *testing.T) {
	t.Run("falls back to global viper", func(t *testing.T) {
		fs := pflag.NewFlagSet("snakes", pflag.ContinueOnError)
		val := fs.Int("the-flag", 0, "some int flag")

		viper.Set("the-flag", 42)
		defer viper.Reset()

		require.NoError(t, fs.Parse([]string{}))
		require.NoError(t, BindViper(fs, nil))
		require.Equal(t, 42, *val)
	})

	runBindViperTestCases(t)
}
