// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"log"

	"terraform-provider-movies/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version string = "dev"
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "hashicorp.com/edu/movies",
		Debug:   false,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
