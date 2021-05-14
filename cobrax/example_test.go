package cobrax_test

import (
	"fmt"
	"os"

	"github.com/martinohmann/exp/cobrax"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ExampleHookViper() {
	v := viper.New()

	// Configure viper as needed.
	v.SetEnvPrefix("example")

	rootCmd := &cobra.Command{
		Use: "example",
		// Install the hook which binds viper to the command and sets flags not
		// explicitly set by the user to values from the viper env and config
		// files.
		PersistentPreRunE: cobrax.HookViper(v),
		Run: func(cmd *cobra.Command, args []string) {
			// Do something.
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func ExampleExecute() {
	v := viper.New()

	// Configure viper as needed.
	v.SetEnvPrefix("example")

	rootCmd := &cobra.Command{
		Use: "example",
		Run: func(cmd *cobra.Command, args []string) {
			// Do something.
		},
	}

	// Binds viper to the command and sets flags not explicitly set by the user
	// to values from the viper env and config files. This is an alternative to
	// cmd.Execute().
	if err := cobrax.Execute(rootCmd, v); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
