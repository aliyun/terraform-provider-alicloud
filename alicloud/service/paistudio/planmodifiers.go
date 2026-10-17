package paistudio

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// mapRequiresReplace returns a planmodifier.Map that flags the resource for
// replacement whenever a known, non-null tags value differs from the prior
// state. Tags are Create-only on CreateResourceGroup (UpdateResourceGroup has
// no Tag parameter), so any change forces a new resource. It mirrors the
// built-in stringplanmodifier.RequiresReplace semantics for the Map type.
func mapRequiresReplace() planmodifier.Map {
	return mapRequiresReplaceModifier{}
}

type mapRequiresReplaceModifier struct{}

func (m mapRequiresReplaceModifier) Description(_ context.Context) string {
	return "Flags the resource for replacement when the value changes."
}

func (m mapRequiresReplaceModifier) MarkdownDescription(_ context.Context) string {
	return "Flags the resource for replacement when the value changes."
}

func (m mapRequiresReplaceModifier) PlanModifyMap(_ context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	// Unknown values cannot be compared; let them be replanned.
	if req.PlanValue.IsUnknown() {
		return
	}
	// Null planned value (attribute removed) carries no replace intent here;
	// a removed tags map is handled as a no-op rather than a destructive
	// replacement, matching Create-only semantics.
	if req.PlanValue.IsNull() {
		return
	}
	if req.PlanValue.Equal(req.StateValue) {
		return
	}
	resp.RequiresReplace = true
}
