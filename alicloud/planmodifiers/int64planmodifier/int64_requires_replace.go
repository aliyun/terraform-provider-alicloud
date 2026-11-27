package int64planmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// RequiresReplace mirrors stringplanmodifier.RequiresReplace for Int64
// attributes. It flags replacement only on an update whose planned value differs
// from the prior state value, so an unconfigured (computed) attribute never
// forces a replacement.
func RequiresReplace() planmodifier.Int64 {
	return int64RequiresReplaceModifier{}
}

type int64RequiresReplaceModifier struct{}

func (m int64RequiresReplaceModifier) Description(_ context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

func (m int64RequiresReplaceModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m int64RequiresReplaceModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
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
