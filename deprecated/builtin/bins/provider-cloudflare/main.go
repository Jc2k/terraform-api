package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/cloudflare"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cloudflare.Provider,
	})
}
