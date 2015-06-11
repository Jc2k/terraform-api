package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/postgresql"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: postgresql.Provider,
	})
}
