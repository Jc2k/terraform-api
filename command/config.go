package command

import (
	"github.com/mitchellh/cli"
	"github.com/xanzy/terraform-api/terraform"
)

// Config is a structure used to configure many commands with Terraform
// configurations.
type Config struct {
	Hooks     []terraform.Hook
	Providers map[string]terraform.ResourceProviderFactory
	Ui        cli.Ui
}

func (c *Config) ContextOpts() *terraform.ContextOpts {
	hooks := make([]terraform.Hook, len(c.Hooks)+1)
	copy(hooks, c.Hooks)
	hooks[len(c.Hooks)] = &UiHook{Ui: c.Ui}

	return &terraform.ContextOpts{
		Hooks:     hooks,
		Providers: c.Providers,
	}
}
