package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/template"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: template.Provider,
	})
}
