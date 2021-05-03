package pflagx

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// ValidatorFunc is a func that can be registered via RegisterValidatorFunc to
// validate flag values before they are set.
type ValidatorFunc func(val string) error

type validatedValue struct {
	pflag.Value
	fn ValidatorFunc
}

func newValidatedValue(value pflag.Value, fn ValidatorFunc) *validatedValue {
	return &validatedValue{
		Value: value,
		fn:    fn,
	}
}

// Set implements the pflag.Value interface.
func (f *validatedValue) Set(s string) error {
	if err := f.fn(s); err != nil {
		return err
	}

	return f.Value.Set(s)
}

// RegisterValidatorFunc registers a ValidatorFunc for the named flag on fs.
// The func is invoked before setting the flag value. Panics if fn is nil, fs
// does not contain a flag with name or if the FlagSet is already parsed.
func RegisterValidatorFunc(fs *pflag.FlagSet, name string, fn ValidatorFunc) {
	if fn == nil {
		panic("pflagx.RegisterValidatorFunc: nil validator func")
	}

	if fs.Parsed() {
		panic("pflagx.RegisterValidatorFunc: must be invoked before fs.Parse()")
	}

	flag := fs.Lookup(name)
	if flag == nil {
		panic(fmt.Sprintf("pflagx.RegisterValidatorFunc: flag %q not defined", name))
	}

	flag.Value = newValidatedValue(flag.Value, fn)
}

// AnyOf returns a ValidatorFunc that allows a flag to have any of the provided
// values.
func AnyOf(values ...string) ValidatorFunc {
	return func(val string) error {
		for _, v := range values {
			if v == val {
				return nil
			}
		}

		return fmt.Errorf(`possible values: "%s"`, strings.Join(values, `", "`))
	}
}
