package consul

import (
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/xanzy/terraform-api/config"
	"github.com/xanzy/terraform-api/helper/schema"
	"github.com/xanzy/terraform-api/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"consul": testAccProvider,
	}

	// Use the demo address for the acceptance tests
	testAccProvider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		conf := consulapi.DefaultConfig()
		conf.Address = "demo.consul.io:80"
		return consulapi.NewClient(conf)
	}
}

func TestResourceProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestResourceProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func TestResourceProvider_Configure(t *testing.T) {
	rp := Provider()

	raw := map[string]interface{}{
		"address":    "demo.consul.io:80",
		"datacenter": "nyc3",
		"scheme":     "https",
	}

	rawConfig, err := config.NewRawConfig(raw)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = rp.Configure(terraform.NewResourceConfig(rawConfig))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}
