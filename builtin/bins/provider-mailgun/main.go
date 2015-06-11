package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/mailgun"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: mailgun.Provider,
	})
}
