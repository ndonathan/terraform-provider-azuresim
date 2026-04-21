package main

import (
	"context"
	"log"

	"github.com/ndonathan/terraform-provider-azuresim/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version string = "dev"
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "ndonathan.github.io/ndonathan/azuresim",
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
