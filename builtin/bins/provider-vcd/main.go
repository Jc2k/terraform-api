package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/vcd"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: vcd.Provider,
	})
}
