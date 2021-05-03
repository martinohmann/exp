package pflagx

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func TestRegisterValidatorFunc(t *testing.T) {
	t.Run("nil fn panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		RegisterValidatorFunc(fs, "the-flag", nil)
	})

	t.Run("panics if flags are already parsed", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		require.NoError(t, fs.Parse(nil))

		RegisterValidatorFunc(fs, "the-flag", func(s string) error { return nil })
	})

	t.Run("panics if flag was not defined", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

		RegisterValidatorFunc(fs, "the-flag", func(s string) error { return nil })
	})

	t.Run("registers ValidatorFunc", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		theFlag := fs.String("the-flag", "", "the usage")

		RegisterValidatorFunc(fs, "the-flag", AnyOf("valid-value"))

		require.NoError(t, fs.Parse([]string{"--the-flag", "valid-value"}))
		require.Equal(t, "valid-value", *theFlag)
	})

	t.Run("fails on Parse()", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.String("the-flag", "", "the usage")

		RegisterValidatorFunc(fs, "the-flag", AnyOf("valid-value", "other-value"))

		err := fs.Parse([]string{"--the-flag", "invalid-value"})
		require.EqualError(t, err, `invalid argument "invalid-value" for "--the-flag" flag: possible values: "valid-value", "other-value"`)
	})
}
