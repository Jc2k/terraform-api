package api

import (
	"encoding/json"
	"fmt"

	"github.com/xanzy/terraform-api/api/tfpb"
	"github.com/xanzy/terraform-api/terraform"
	"golang.org/x/net/context"
)

// State implements the TerraformServer interface
func (s *Server) State(c context.Context, req *tfpb.StateRequest) (*tfpb.StateResponse, error) {
	resp := &tfpb.StateResponse{}

	state, err := json.Marshal(terraform.NewState())
	if err != nil {
		return nil, fmt.Errorf("Error marshalling new state: %v", err)
	}

	resp.State = state

	return resp, nil
}
