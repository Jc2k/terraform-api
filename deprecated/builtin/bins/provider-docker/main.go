package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/docker"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: docker.Provider,
	})
}
