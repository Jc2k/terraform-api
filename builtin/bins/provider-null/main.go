package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/null"
	"github.com/xanzy/terraform-api/plugin"
	"github.com/xanzy/terraform-api/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return null.Provider()
		},
	})
}
