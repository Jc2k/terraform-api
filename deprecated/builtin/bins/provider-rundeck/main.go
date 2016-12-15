package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/rundeck"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: rundeck.Provider,
	})
}
