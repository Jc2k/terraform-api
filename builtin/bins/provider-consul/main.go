package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/consul"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: consul.Provider,
	})
}
