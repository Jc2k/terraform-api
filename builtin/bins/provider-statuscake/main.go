package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/statuscake"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: statuscake.Provider,
	})
}
