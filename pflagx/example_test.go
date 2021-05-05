package pflagx_test

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/martinohmann/exp/pflagx"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func ExampleRegisterValidatorFunc() {
	fs := pflag.NewFlagSet("pflagx", pflag.ContinueOnError)
	fs.String("my-flag", "", "flag usage")

	pflagx.RegisterValidatorFunc(fs, "my-flag", pflagx.AnyOf("one", "two"))

	fmt.Println(fs.Parse([]string{"--my-flag", "three"}))

	// Output:
	// invalid argument "three" for "--my-flag" flag: possible values: "one", "two"
}

func ExampleRegisterTransformerFunc() {
	fs := pflag.NewFlagSet("pflagx", pflag.ContinueOnError)
	val := fs.String("my-flag", "", "flag usage")

	pflagx.RegisterTransformerFunc(fs, "my-flag", strings.ToUpper)

	err := fs.Parse([]string{"--my-flag", "the-value"})
	if err != nil {
		panic(err)
	}

	fmt.Println(*val)

	// Output:
	// THE-VALUE
}

func ExampleFunc() {
	fs := pflag.NewFlagSet("pflagx", pflag.ContinueOnError)

	var val *big.Int

	pflagx.Func(fs, "big-int", "a big int", func(s string) error {
		z, ok := new(big.Int).SetString(s, 10)
		if !ok {
			return errors.New("not a big int")
		}

		val = z
		return nil
	})

	err := fs.Parse([]string{"--big-int", "12345678901234567890123456789012345678901234567890"})
	if err != nil {
		panic(err)
	}

	fmt.Println(val.String())

	// Output:
	// 12345678901234567890123456789012345678901234567890
}

func ExampleBindViper() {
	flags := pflag.NewFlagSet("snakes", pflag.ContinueOnError)
	val := flags.String("the-flag", "", "the flag usage")

	if err := flags.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	v := viper.New()

	// Configure viper as needed.
	v.SetEnvPrefix("example")

	if err := pflagx.BindViper(flags, v); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	fmt.Println(*val)
}
