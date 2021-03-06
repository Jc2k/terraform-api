package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/azure"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: azure.Provider,
	})
}
