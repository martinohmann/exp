package pflagx

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// BindViper binds a *viper.Viper to a *pflag.FlagSet. Values of flags not
// explicitly set by the user will be filled with values obtained from viper's
// env or configuration.
//
// If v is nil the global viper instance is used.
//
// Preference rules are as follows with later ones taking higher precedence:
//
//   1. flag default value
//   2. value from config file(s)
//   3. env var value
//   4. flag value provided on the commandline
//
// Returns an error if the type of a value for viper is not compatible with the
// corresponding flag value type or if there are errors while reading the viper
// config. Nonexistent configuration files do not cause errors.
func BindViper(fs *pflag.FlagSet, v *viper.Viper) error {
	if v == nil {
		v = viper.GetViper()
	}

	if err := v.ReadInConfig(); err != nil {
		var notFoundErr viper.ConfigFileNotFoundError

		// Non-existent config files are fine and should not be treated as
		// errors.
		if !errors.As(err, &notFoundErr) {
			return err
		}
	}

	// Enable automatic configuration via environment variables. The
	// EnvKeyReplacer is required to correctly build environment variables for
	// flags that contain dashes or dots.
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	return bindFlags(fs, v)
}

// bindFlags binds environment variable and config values from v to flags that
// were not explicitly set by the user.
func bindFlags(flags *pflag.FlagSet, v *viper.Viper) (err error) {
	flags.VisitAll(func(f *pflag.Flag) {
		if err != nil {
			return
		}

		if err = bindFlag(flags, f, v); err != nil {
			err = fmt.Errorf("failed to set flag from env or config: %w", err)
		}
	})

	return err
}

func bindFlag(flags *pflag.FlagSet, f *pflag.Flag, v *viper.Viper) error {
	if f.Changed {
		// If the flag was explicitly provided on the commandline by the
		// user, do not attempt override it with a value from the viper env
		// or config.
		return nil
	}

	// Bind environment variable for flag before looking up the value.
	v.BindEnv(f.Name) // nolint: errcheck

	val := lookupValue(v, f.Name)
	if val == nil {
		// Key was not found in viper.
		return nil
	}

	// Key was found, attempt to set the flag using the string representation
	// of the value.
	return flags.Set(f.Name, stringify(val))
}

// lookupValue retrieves the value for key from v. Attempts to also look up
// nested config values if no value for key was found. For example if key is
// 'some-key' and it is not found, a lookup for 'some.key' will be attempted as
// well.
func lookupValue(v *viper.Viper, key string) interface{} {
	val := v.Get(key)
	if val != nil {
		return val
	}

	dottedKey := strings.ReplaceAll(key, "-", ".")

	return v.Get(dottedKey)
}

// stringify converts val to a string that can be parsed by *pflag.FlagSet.Set.
func stringify(val interface{}) string {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Map:
		return stringifyMap(val)
	case reflect.Slice:
		return stringifySlice(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// stringifySlice produces a string of the form val1,val2,val3 from val. Panics
// if val is not of slice type.
func stringifySlice(val interface{}) string {
	rval := reflect.ValueOf(val)
	n := rval.Len()
	s := make([]string, n)

	for i := 0; i < n; i++ {
		v := rval.Index(i).Interface()
		s[i] = fmt.Sprintf("%v", v)
	}

	return strings.Join(s, ",")
}

// stringifyMap produces a string of the form key1=val1,key2=val2,key3=val3
// from val. Panics if val is not of map type.
func stringifyMap(val interface{}) string {
	rval := reflect.ValueOf(val)
	s := make([]string, 0, rval.Len())

	for iter := rval.MapRange(); iter.Next(); {
		k := iter.Key().Interface()
		v := iter.Value().Interface()
		s = append(s, fmt.Sprintf("%v=%v", k, v))
	}

	return strings.Join(s, ",")
}
