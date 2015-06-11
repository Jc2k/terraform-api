package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/google"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: google.Provider,
	})
}
