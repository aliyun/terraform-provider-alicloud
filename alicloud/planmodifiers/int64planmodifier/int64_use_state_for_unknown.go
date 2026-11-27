package int64planmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// UseStateForUnknown mirrors stringplanmodifier.UseStateForUnknown for Int64
// attributes. It copies a known prior state value into the planned value, so an
// unconfigured, computed attribute keeps its current value instead of being
// shown as "(known after apply)" on every update.
//
// On an Optional+Computed+ForceNew attribute, place UseStateForUnknown before
// RequiresReplace: the framework marks an unconfigured computed attribute
// unknown on update, UseStateForUnknown resolves it back to the prior value, and
// RequiresReplace then fires only on a genuine configuration change.
func UseStateForUnknown() planmodifier.Int64 {
	return int64UseStateForUnknownModifier{}
}

type int64UseStateForUnknownModifier struct{}

func (m int64UseStateForUnknownModifier) Description(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

func (m int64UseStateForUnknownModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m int64UseStateForUnknownModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	// Do nothing if there is no state (resource is being created).
	if req.State.Raw.IsNull() {
		return
	}

	// Do nothing if there is a known planned value.
	if !req.PlanValue.IsUnknown() {
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.ConfigValue.IsUnknown() {
		return
	}

	resp.PlanValue = req.StateValue
}
