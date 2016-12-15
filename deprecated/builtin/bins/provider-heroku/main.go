package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/heroku"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: heroku.Provider,
	})
}
