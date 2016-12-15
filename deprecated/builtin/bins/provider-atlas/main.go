package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/atlas"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: atlas.Provider,
	})
}
