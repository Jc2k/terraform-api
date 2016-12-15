package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/digitalocean"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: digitalocean.Provider,
	})
}
