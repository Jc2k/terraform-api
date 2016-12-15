package main

import (
	"github.com/xanzy/terraform-api/builtin/provisioners/remote-exec"
	"github.com/xanzy/terraform-api/plugin"
	"github.com/xanzy/terraform-api/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: func() terraform.ResourceProvisioner {
			return new(remoteexec.ResourceProvisioner)
		},
	})
}
