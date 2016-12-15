package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/dnsimple"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dnsimple.Provider,
	})
}
