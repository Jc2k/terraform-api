package api

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/xanzy/terraform-api/api/tfpb"
	"github.com/xanzy/terraform-api/terraform"
	"golang.org/x/net/context"
)

// Refresh implements the TerraformServer interface
func (s *Server) Refresh(c context.Context, req *tfpb.RefreshRequest) (*tfpb.RefreshResponse, error) {
	oldState, err := terraform.ReadState(bytes.NewReader(req.State))
	if err != nil {
		return nil, fmt.Errorf("Error reading state: %v", err)
	}

	resp := &tfpb.RefreshResponse{}

	ctx, err := s.newContext(req.Config, false, nil, req.State, req.Parallelism, nil)
	if err != nil {
		return nil, err
	}

	if err := validateContext(ctx); err != nil {
		return nil, fmt.Errorf("Error validating context: %v", err)
	}

	newState, err := ctx.Refresh()
	if err != nil {
		return nil, fmt.Errorf("Error refreshing state: %v", err)
	}

	// Check if we need to update the state serial
	newState.IncrementSerialMaybe(oldState)
	resp.Serial = newState.Serial

	resp.State, err = json.Marshal(newState)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling refreshed state: %v", err)
	}

	return resp, nil
}
