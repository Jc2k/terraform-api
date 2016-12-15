package api

import (
	"github.com/hashicorp/go-multierror"
	"github.com/xanzy/terraform-api/api/tfpb"
	"github.com/xanzy/terraform-api/config/module"
	"golang.org/x/net/context"
)

// Validate implements the TerraformServer interface
func (s *Server) Validate(c context.Context, req *tfpb.ValidateRequest) (*tfpb.ValidateResponse, error) {
	resp := &tfpb.ValidateResponse{Valid: true}

	ctx, err := s.newContext(req.Config, false, nil, nil, 0, nil)
	if err != nil {
		return nil, err
	}

	if ws, es := ctx.Validate(); len(ws) > 0 || len(es) > 0 {
		resp.Warnings = ws

		if len(es) > 0 {
			resp.Valid = false
			parseErrors(&resp.Errors, es)
		}
	}

	return resp, nil
}

// parseErrors parses (and if needed converts) a whole range
// of possible values into a slice of strings
func parseErrors(s *[]string, v interface{}) {
	if v == nil {
		return
	}

	switch v := v.(type) {
	case []string:
		*s = v
	case []error:
		for _, e := range v {
			parseErrors(s, e)
		}
	case *multierror.Error:
		parseErrors(s, v.Errors)
	case *module.TreeError:
		parseErrors(s, v.Err)
	case error:
		parseErrors(s, v.Error())
	case string:
		*s = append(*s, v)
	}
}
