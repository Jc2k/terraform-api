package api

import (
	"encoding/json"
	"sync"

	"github.com/xanzy/terraform-api/api/tfpb"
	"github.com/xanzy/terraform-api/terraform"
)

// ApplyHook is used during an apply request
type ApplyHook struct {
	terraform.NilHook
	sync.Mutex

	oldState *terraform.State
	stream   tfpb.Terraform_ApplyServer
	resp     *tfpb.ApplyResponse
}

// PreApply is called before a single resource is applied, it adds the new
// state to the ApplyResponse and sends it to the calling gRPC client
func (h *ApplyHook) PreApply(
	n *terraform.InstanceInfo,
	s *terraform.InstanceState,
	d *terraform.InstanceDiff) (terraform.HookAction, error) {
	h.Lock()
	defer h.Unlock()

	h.resp.States[n.HumanId()] = tfpb.ResourceState_StateRunning

	// Write the new state over the connected gRPC stream
	if err := h.stream.Send(h.resp); err != nil {
		return terraform.HookActionHalt, err
	}

	return terraform.HookActionContinue, nil
}

// PostApply is called after a single resource is applied, it adds the new
// state to the ApplyResponse and sends it to the calling gRPC client
func (h *ApplyHook) PostApply(
	n *terraform.InstanceInfo,
	s *terraform.InstanceState,
	err error) (terraform.HookAction, error) {
	h.Lock()
	defer h.Unlock()

	if err != nil {
		h.resp.States[n.HumanId()] = tfpb.ResourceState_StateError
	} else {
		h.resp.States[n.HumanId()] = tfpb.ResourceState_StateSuccess
	}

	// Write the new state over the connected gRPC stream
	if err := h.stream.Send(h.resp); err != nil {
		return terraform.HookActionHalt, err
	}

	return terraform.HookActionContinue, nil
}

// PostStateUpdate continuously updates the state in a ApplyResponse
// and sends the updated response to the calling gRPC client
func (h *ApplyHook) PostStateUpdate(newState *terraform.State) (terraform.HookAction, error) {
	h.Lock()
	defer h.Unlock()

	if newState != nil {
		// Check if we need to update the state serial
		newState.IncrementSerialMaybe(h.oldState)
		h.resp.Serial = newState.Serial

		state, err := json.Marshal(newState)
		if err != nil {
			return terraform.HookActionHalt, err
		}
		h.resp.State = state

		// Write the new state over the connected gRPC stream
		if err := h.stream.Send(h.resp); err != nil {
			return terraform.HookActionHalt, err
		}
	}

	// Continue forth
	return terraform.HookActionContinue, nil
}

// PlanHook is used during a plan request
type PlanHook struct {
	terraform.NilHook
	sync.Mutex

	resp *tfpb.PlanResponse
}

// PostDiff is triggered after each individual resource is diffed, and adds
// the required action for each resource to the PlanResponse
func (h *PlanHook) PostDiff(
	n *terraform.InstanceInfo,
	d *terraform.InstanceDiff) (terraform.HookAction, error) {
	h.Lock()
	defer h.Unlock()

	switch d.ChangeType() {
	case terraform.DiffCreate:
		h.resp.Actions[n.HumanId()] = tfpb.ResourceAction_ActionCreate
	case terraform.DiffUpdate:
		h.resp.Actions[n.HumanId()] = tfpb.ResourceAction_ActionUpdate
	case terraform.DiffDestroy:
		h.resp.Actions[n.HumanId()] = tfpb.ResourceAction_ActionDestroy
	case terraform.DiffDestroyCreate:
		h.resp.Actions[n.HumanId()] = tfpb.ResourceAction_ActionRecreate
	default:
		h.resp.Actions[n.HumanId()] = tfpb.ResourceAction_ActionNone
	}

	return terraform.HookActionContinue, nil
}
