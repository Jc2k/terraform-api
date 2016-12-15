package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/openstack"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: openstack.Provider,
	})
}
