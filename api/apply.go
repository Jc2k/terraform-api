package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xanzy/terraform-api/api/tfpb"
	"github.com/xanzy/terraform-api/terraform"
)

// Apply implements the TerraformServer interface
func (s *Server) Apply(req *tfpb.ApplyRequest, stream tfpb.Terraform_ApplyServer) error {
	resp := &tfpb.ApplyResponse{
		States: make(map[string]tfpb.ResourceState),
	}

	if req.Destroy && req.Plan != nil {
		return errors.New("Destroy can't be called with a plan")
	}

	oldState, err := terraform.ReadState(bytes.NewReader(req.State))
	if err != nil {
		return fmt.Errorf("Error reading state: %v", err)
	}

	hooks := []terraform.Hook{&ApplyHook{
		oldState: oldState,
		stream:   stream,
		resp:     resp,
	}}

	ctx, err := s.newContext(req.Config, req.Destroy, req.Plan, req.State, req.Parallelism, hooks)
	if err != nil {
		return err
	}

	if err := validateContext(ctx); err != nil {
		return fmt.Errorf("Error validating context: %v", err)
	}

	if req.Plan == nil {
		if req.Refresh {
			_, err := ctx.Refresh()
			if err != nil {
				return fmt.Errorf("Error refreshing state: %v", err)
			}
		}

		_, err := ctx.Plan()
		if err != nil {
			return fmt.Errorf("Error creating plan: %v", err)
		}
	}

	newState, err := ctx.Apply()
	parseErrors(&resp.Errors, err)

	// Make sure we send the last known state back to the client
	if newState != nil {
		// Check if we need to update the state serial
		newState.IncrementSerialMaybe(oldState)
		resp.Serial = newState.Serial

		resp.State, err = json.Marshal(newState)
		if err != nil {
			return fmt.Errorf("Error marshalling final state: %v", err)
		}

		err := stream.Send(resp)
		if err != nil {
			return fmt.Errorf("Error sending final state: %v", err)
		}
	}

	return nil
}
