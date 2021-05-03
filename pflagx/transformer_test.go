package pflagx

import (
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func TestRegisterTransformerFunc(t *testing.T) {
	t.Run("nil fn panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		RegisterTransformerFunc(fs, "the-flag", nil)
	})

	t.Run("panics if flags are already parsed", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		require.NoError(t, fs.Parse(nil))

		RegisterTransformerFunc(fs, "the-flag", strings.ToUpper)
	})

	t.Run("panics if flag was not defined", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

		RegisterTransformerFunc(fs, "the-flag", strings.ToUpper)
	})

	t.Run("registers TransformerFunc", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		theFlag := fs.String("the-flag", "", "the usage")

		RegisterTransformerFunc(fs, "the-flag", strings.ToUpper)

		require.NoError(t, fs.Parse([]string{"--the-flag", "valid-value"}))
		require.Equal(t, "VALID-VALUE", *theFlag)
	})
}
