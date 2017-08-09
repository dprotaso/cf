package main

import (
	"fmt"

	"github.com/dprotaso/cf"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	// Example of a required flag
	Endpoint string `short:"a" long:"api" description:"A name" required:"true"`
}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		panic(err)
	}

	creds := cf.PasswordCredentials(
		"cf",
		"",
		"")

	var (
		client *cf.Client
		err    error
		bytes  []byte
	)

	if client, err = cf.NewClient(opts.Endpoint, creds); err != nil {
		panic(err)
	}

	if bytes, err = client.Orgs(); err == nil {
		fmt.Printf("%s", string(bytes))
		return
	}

	panic(err)
}
