package float32planmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// RequiresReplace mirrors stringplanmodifier.RequiresReplace for Float32
// attributes. It flags replacement only on an update whose planned value differs
// from the prior state value, so an unconfigured (computed) attribute never
// forces a replacement.
func RequiresReplace() planmodifier.Float32 {
	return float32RequiresReplaceModifier{}
}

type float32RequiresReplaceModifier struct{}

func (m float32RequiresReplaceModifier) Description(_ context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

func (m float32RequiresReplaceModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m float32RequiresReplaceModifier) PlanModifyFloat32(ctx context.Context, req planmodifier.Float32Request, resp *planmodifier.Float32Response) {
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
