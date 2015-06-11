package api

import (
	"encoding/json"
	"fmt"

	"github.com/xanzy/terraform-api/api/pb"
	"golang.org/x/net/context"
)

// Refresh implements the TerraformServer interface
func (s *Server) Refresh(c context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	resp := new(pb.RefreshResponse)

	ctx, err := s.newContext(req.Config, false, nil, req.State, req.Parallelism, nil)
	if err != nil {
		return nil, err
	}

	if err := validateContext(ctx); err != nil {
		return nil, fmt.Errorf("Error validating context: %v", err)
	}

	state, err := ctx.Refresh()
	if err != nil {
		return nil, fmt.Errorf("Error refreshing state: %v", err)
	}

	resp.State, err = json.Marshal(state)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling refreshed state: %v", err)
	}

	return resp, nil
}
