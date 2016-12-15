package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/vsphere"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: vsphere.Provider,
	})
}
