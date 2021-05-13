package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/martinohmann/exit"
	"github.com/spf13/pflag"
)

// Run runs fn and handles returned errors by exiting the program with a
// suitable exit code.
func Run(fn func() error) {
	exit.SetErrorHandler(func(err error) (code int, handled bool) {
		if errors.Is(err, pflag.ErrHelp) {
			return exit.CodeOK, true
		}

		fmt.Fprintln(os.Stderr, color.RedString("error:"), err)

		return 0, false
	})

	exit.Exit(fn())
}
