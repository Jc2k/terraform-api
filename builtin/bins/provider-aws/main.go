package main

import (
	"github.com/xanzy/terraform-api/builtin/providers/aws"
	"github.com/xanzy/terraform-api/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: aws.Provider,
	})
}
