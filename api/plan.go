package api

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/xanzy/terraform-api/api/tfpb"
	"github.com/xanzy/terraform-api/terraform"
	"golang.org/x/net/context"
)

// Plan implements the TerraformServer interface
func (s *Server) Plan(c context.Context, req *tfpb.PlanRequest) (*tfpb.PlanResponse, error) {
	resp := &tfpb.PlanResponse{
		Actions: make(map[string]tfpb.ResourceAction),
	}

	oldState, err := terraform.ReadState(bytes.NewReader(req.State))
	if err != nil {
		return nil, fmt.Errorf("Error reading state: %v", err)
	}

	hooks := []terraform.Hook{&PlanHook{
		resp: resp,
	}}

	ctx, err := s.newContext(req.Config, req.Destroy, nil, req.State, req.Parallelism, hooks)
	if err != nil {
		return nil, err
	}

	if err := validateContext(ctx); err != nil {
		return nil, fmt.Errorf("Error validating context: %v", err)
	}

	if req.Refresh {
		_, err := ctx.Refresh()
		if err != nil {
			return nil, fmt.Errorf("Error refreshing state: %v", err)
		}
	}

	plan, err := ctx.Plan()
	if err != nil {
		return nil, fmt.Errorf("Error running plan: %v", err)
	}

	var b bytes.Buffer
	err = terraform.WritePlan(plan, &b)
	if err != nil {
		return nil, fmt.Errorf("Error writing plan: %v", err)
	}
	resp.Plan = b.Bytes()

	resp.Diff, err = json.Marshal(plan.Diff)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling diff: %v", err)
	}

	resp.State, err = json.Marshal(plan.State)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling refreshed state: %v", err)
	}

	// Check if we need to update the state serial
	plan.State.IncrementSerialMaybe(oldState)
	resp.Serial = plan.State.Serial

	return resp, nil
}
