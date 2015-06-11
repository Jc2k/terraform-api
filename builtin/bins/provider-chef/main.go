package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/chef"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: chef.Provider,
	})
}
