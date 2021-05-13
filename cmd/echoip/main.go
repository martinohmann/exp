package main

import (
	"context"
	"os"

	"github.com/martinohmann/exit"
	"github.com/martinohmann/exp/cli"
	"github.com/martinohmann/exp/http/echoip"
	"github.com/martinohmann/exp/jsonx"
	"github.com/spf13/pflag"
)

func main() {
	cli.Run(run)
}

func run() error {
	fs := pflag.NewFlagSet("echoip", pflag.ContinueOnError)

	ip := fs.IP("ip", nil, "IP address to lookup info for.")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return exit.Error(exit.CodeUsage, err)
	}

	client := echoip.NewClient("https://ifconfig.co/")

	resp, err := client.Lookup(context.Background(), &echoip.Options{IP: *ip})
	if err != nil {
		return err
	}

	return jsonx.WriteIndent(os.Stdout, resp, "", "  ")
}
