package int32planmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// RequiresReplace mirrors stringplanmodifier.RequiresReplace for Int32
// attributes. It flags replacement only on an update whose planned value differs
// from the prior state value, so an unconfigured (computed) attribute never
// forces a replacement.
func RequiresReplace() planmodifier.Int32 {
	return int32RequiresReplaceModifier{}
}

type int32RequiresReplaceModifier struct{}

func (m int32RequiresReplaceModifier) Description(_ context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

func (m int32RequiresReplaceModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m int32RequiresReplaceModifier) PlanModifyInt32(ctx context.Context, req planmodifier.Int32Request, resp *planmodifier.Int32Response) {
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
