package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/openfaas/terraform-provider-openfaas/openfaas"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: openfaas.Provider})
}
