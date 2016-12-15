package null

import (
	"testing"

	"github.com/xanzy/terraform-api/helper/schema"
	"github.com/xanzy/terraform-api/terraform"
)

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}
