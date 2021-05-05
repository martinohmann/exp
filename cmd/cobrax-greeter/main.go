// Example application for demonstating the cobrax package.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/martinohmann/exp/cobrax"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type options struct {
	verbose    bool
	message    string
	listenAddr string
}

func main() {
	v := viper.New()
	v.SetEnvPrefix("greeter")
	v.SetConfigName("greeter")
	v.AddConfigPath(".")

	opts := &options{
		message:    "Hello World!",
		listenAddr: "127.0.0.1:8080",
	}

	cmd := newRootCommand(opts)

	err := cobrax.Execute(cmd, v)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func newRootCommand(opts *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "greeter",
		Short:         "Prints greetings",
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.SilenceUsage = true

			if !opts.verbose {
				log.SetOutput(io.Discard)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(opts.message)
		},
	}

	cmd.PersistentFlags().BoolVar(&opts.verbose, "verbose", opts.verbose, "verbose output")
	cmd.PersistentFlags().StringVarP(&opts.message, "message", "m", opts.message, "a nice greeting message")

	cmd.AddCommand(newServeCommand(opts))

	return cmd
}

func newServeCommand(opts *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serves greetings",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Printf("listening on %s\n", opts.listenAddr)

			handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				defer func(start time.Time) {
					log.Printf("%s %s from %s took %v", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
				}(time.Now())

				fmt.Fprintln(rw, opts.message)
			})

			return http.ListenAndServe(opts.listenAddr, handler)
		},
	}

	cmd.Flags().StringVar(&opts.listenAddr, "listen-addr", opts.listenAddr, "address to listen on")

	return cmd
}
