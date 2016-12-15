package template

import (
	"github.com/xanzy/terraform-api/helper/schema"
	"github.com/xanzy/terraform-api/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"template_file":             resourceFile(),
			"template_cloudinit_config": resourceCloudinitConfig(),
		},
	}
}
