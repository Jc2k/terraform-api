package terraform

import (
	"github.com/xanzy/terraform-api/helper/schema"
	"github.com/xanzy/terraform-api/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"terraform_remote_state": resourceRemoteState(),
		},
	}
}
