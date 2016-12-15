package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/dme"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dme.Provider,
	})
}
