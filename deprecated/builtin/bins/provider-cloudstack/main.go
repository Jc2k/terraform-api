package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/cloudstack"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cloudstack.Provider,
	})
}
