package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"

	"github.com/SectorLabs/terraform-provider-circleci/circleci"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: circleci.Provider,
	})
}
