package heroku

import (
	"os"
	"testing"

	"github.com/xanzy/terraform-api/helper/schema"
	"github.com/xanzy/terraform-api/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"heroku": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("HEROKU_EMAIL"); v == "" {
		t.Fatal("HEROKU_EMAIL must be set for acceptance tests")
	}

	if v := os.Getenv("HEROKU_API_KEY"); v == "" {
		t.Fatal("HEROKU_API_KEY must be set for acceptance tests")
	}
}
