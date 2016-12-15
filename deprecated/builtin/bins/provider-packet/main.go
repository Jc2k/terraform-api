package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/packet"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: packet.Provider,
	})
}
