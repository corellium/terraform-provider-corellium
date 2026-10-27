package main

import (
	"context"
	"flag"
	"log"

	"terraform-provider-corellium/corellium"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	 flag.Parse()

    opts := providerserver.ServeOpts{
        Address: "registry.terraform.io/corellium/corellium",
        Debug:   *corellium.Debug,
    }

    err := providerserver.Serve(context.Background(), corellium.New, opts)

    if err != nil {
        log.Fatal(err.Error())
    }
}
