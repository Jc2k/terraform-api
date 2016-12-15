package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/terraform"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: terraform.Provider,
	})
}
