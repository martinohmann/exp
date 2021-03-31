package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/martinohmann/exp/http/echoip"
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

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	if err := enc.Encode(resp); err != nil {
		log.Fatal(err)
	}
}
