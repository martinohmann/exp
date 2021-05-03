package pflagx

import (
	"errors"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func TestFunc(t *testing.T) {
	t.Run("sets flag", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

		var val string

		Func(fs, "foo", "", func(v string) error {
			val = v
			return nil
		})

		require.NoError(t, fs.Parse([]string{"--foo", "bar"}))
		require.Equal(t, "bar", val)
	})

	t.Run("func error is parse error", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

		FuncP(fs, "foo", "f", "", func(v string) error {
			return errors.New("whoops")
		})

		require.EqualError(t, fs.Parse([]string{"-f", "bar"}), `invalid argument "bar" for "-f, --foo" flag: whoops`)
	})
}
