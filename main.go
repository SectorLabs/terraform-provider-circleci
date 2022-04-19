package main

import (
	"github.com/SectorLabs/terraform-provider-circleci/circleci"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: circleci.Provider,
	})
}
