package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/xanzy/terraform-api/config"
	"github.com/xanzy/terraform-api/config/module"
	"github.com/xanzy/terraform-api/terraform"
)

// Server represents a Terraform gRPC server
type Server struct {
	providers    map[string]terraform.ResourceProviderFactory
	provisioners map[string]terraform.ResourceProvisionerFactory
}

// NewServer returns an initialized Terraform gRPC server
func NewServer(
	providers map[string]terraform.ResourceProviderFactory,
	provisioners map[string]terraform.ResourceProvisionerFactory) *Server {

	return &Server{
		providers:    providers,
		provisioners: provisioners,
	}
}

func (s *Server) newContext(
	rConf json.RawMessage,
	rDestroy bool,
	rPlan []byte,
	rState json.RawMessage,
	rParallelism int32,
	hooks []terraform.Hook) (*terraform.Context, error) {
	conf, err := config.LoadJSON(rConf)
	if err != nil {
		return nil, err
	}

	mod := module.NewTree("", conf)
	err = mod.Load(nil, module.GetModeNone)
	if err != nil {
		return nil, fmt.Errorf("Error loading module: %s", err)
	}

	ctxOpts := &terraform.ContextOpts{
		Destroy:      rDestroy,
		Hooks:        hooks,
		Module:       mod,
		Parallelism:  int(rParallelism),
		Providers:    s.providers,
		Provisioners: s.provisioners,
	}

	if rState != nil {
		b := bytes.NewBuffer(rState)
		state, err := terraform.ReadState(b)
		if err != nil {
			return nil, fmt.Errorf("Error reading state: %s", err)
		}
		ctxOpts.State = state
	}

	if rPlan != nil {
		b := bytes.NewBuffer(rPlan)
		plan, err := terraform.ReadPlan(b)
		if err != nil {
			return nil, fmt.Errorf("Error reading plan: %s", err)
		}

		return plan.Context(ctxOpts), nil
	}

	return terraform.NewContext(ctxOpts), nil
}

func validateContext(ctx *terraform.Context) error {
	_, es := ctx.Validate()
	if len(es) == 0 {
		return nil
	}

	return errors.New(multierror.ListFormatFunc(es))
}
