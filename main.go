package main

import (
	"context"
	"log"

	"github.com/ndonathan/terraform-provider-azuresim/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/ndonathan/azuresim",
	}

	err := providerserver.Serve(context.Background(), provider.New("0.1.0"), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
