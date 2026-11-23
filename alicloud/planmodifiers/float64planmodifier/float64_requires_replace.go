package float64planmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// RequiresReplace mirrors stringplanmodifier.RequiresReplace for Float64
// attributes. It flags replacement only on an update whose planned value differs
// from the prior state value, so an unconfigured (computed) attribute never
// forces a replacement.
func RequiresReplace() planmodifier.Float64 {
	return float64RequiresReplaceModifier{}
}

type float64RequiresReplaceModifier struct{}

func (m float64RequiresReplaceModifier) Description(_ context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

func (m float64RequiresReplaceModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m float64RequiresReplaceModifier) PlanModifyFloat64(ctx context.Context, req planmodifier.Float64Request, resp *planmodifier.Float64Response) {
	// Do not replace on resource creation.
	if req.State.Raw.IsNull() {
		return
	}

	// Do not replace on resource destroy.
	if req.Plan.Raw.IsNull() {
		return
	}

	// Do not replace if the plan and state values are equal.
	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	resp.RequiresReplace = true
}
