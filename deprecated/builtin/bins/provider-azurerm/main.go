package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/azurerm"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: azurerm.Provider,
	})
}
