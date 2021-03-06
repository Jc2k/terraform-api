package command

import (
	"testing"

	"github.com/xanzy/terraform-api/state"
	"github.com/xanzy/terraform-api/terraform"
)

func TestStateHook_impl(t *testing.T) {
	var _ terraform.Hook = new(StateHook)
}

func TestStateHook(t *testing.T) {
	is := &state.InmemState{}
	var hook terraform.Hook = &StateHook{State: is}

	s := state.TestStateInitial()
	action, err := hook.PostStateUpdate(s)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if action != terraform.HookActionContinue {
		t.Fatalf("bad: %v", action)
	}
	if !is.State().Equal(s) {
		t.Fatalf("bad state: %#v", is.State())
	}
}
