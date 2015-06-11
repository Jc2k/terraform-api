package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/tls"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: tls.Provider,
	})
}
