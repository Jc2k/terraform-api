package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/xanzy/terraform-api/api/tfpb"
	"github.com/xanzy/terraform-api/terraform"
	"golang.org/x/net/context"
)

// Taint implements the TerraformServer interface
func (s *Server) Taint(c context.Context, req *tfpb.TaintRequest) (*tfpb.TaintResponse, error) {
	resp := &tfpb.TaintResponse{}

	if req.Module == "" {
		req.Module = "root"
	} else {
		req.Module = "root." + req.Module
	}

	b := bytes.NewBuffer(req.State)
	state, err := terraform.ReadState(b)
	if err != nil {
		return nil, fmt.Errorf("Error reading state: %v", err)
	}

	// Get the proper module we want to taint
	modPath := strings.Split(req.Module, ".")
	mod := state.ModuleByPath(modPath)
	if mod == nil {
		return nil, fmt.Errorf("Module %s could not be found", req.Module)
	}

	// If there are no resources in this module, we can skip this all together
	if len(mod.Resources) > 0 {
		// Get the resource we're looking for and if it exists, then taint it
		rs, ok := mod.Resources[req.Resource]
		if ok {
			rs.Taint()
		}
	}

	resp.State, err = json.Marshal(state)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling tainted state: %v", err)
	}

	return resp, nil
}
