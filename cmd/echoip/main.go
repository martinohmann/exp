package main

import (
	"context"
	"log"
	"os"

	"github.com/martinohmann/exp/http/echoip"
	"github.com/martinohmann/exp/json"
	"github.com/spf13/pflag"
)

func main() {
	ip := pflag.IP("ip", nil, "IP address to lookup info for.")
	pflag.Parse()

	client := echoip.NewClient("https://ifconfig.co/")

	resp, err := client.Lookup(context.Background(), &echoip.Options{IP: *ip})
	if err != nil {
		log.Fatal(err)
	}

	err = json.WriteIndent(os.Stdout, resp, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
}
