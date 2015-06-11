package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/xanzy/terraform-api/api/pb"
	"github.com/xanzy/terraform-api/terraform"
	"golang.org/x/net/context"
)

// Taint implements the TerraformServer interface
func (s *Server) Taint(c context.Context, req *pb.TaintRequest) (*pb.TaintResponse, error) {
	resp := new(pb.TaintResponse)

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

	// If there are no resources in this module, it is an error
	if len(mod.Resources) == 0 {
		return nil, fmt.Errorf("Module %s has no resources", req.Module)
	}

	// Get the resource we're looking for
	rs, ok := mod.Resources[req.Resource]
	if !ok {
		return nil,
			fmt.Errorf("Resource %s couldn't be found in the module %s", req.Resource, req.Module)
	}

	// Taint the resource and save the updated state
	rs.Taint()

	resp.State, err = json.Marshal(state)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling tainted state: %v", err)
	}

	return resp, nil
}
