// Package fwutil holds Framework schema scaffolding shared by product packages:
// the operation-timeout block generated resources expose. It lives outside the
// product packages so a second product does not carry a copy of the same
// declarations. Retry helpers live in alicloud/helper.
package fwutil

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type timeoutDurationValidator struct{}

// TimeoutTypes returns the timeouts attribute types for the operations a
// resource exposes.
func TimeoutTypes(operations ...string) map[string]attr.Type {
	attributeTypes := make(map[string]attr.Type, len(operations))
	for _, operation := range operations {
		attributeTypes[operation] = types.StringType
	}
	return attributeTypes
}

// TimeoutBlock returns the timeouts block schema for the operations a resource
// exposes. Only the operations passed in are accepted, so a resource without an
// update path does not advertise an update timeout.
func TimeoutBlock(operations ...string) schema.SingleNestedBlock {
	attributes := map[string]schema.Attribute{}
	for _, operation := range operations {
		attributes[operation] = schema.StringAttribute{
			Optional:    true,
			Description: "Maximum retry and wait duration for this operation, for example 10m.",
			Validators:  []validator.String{timeoutDurationValidator{}},
		}
	}
	return schema.SingleNestedBlock{Attributes: attributes}
}

func (timeoutDurationValidator) Description(context.Context) string {
	return "A positive Go duration, for example 10m or 1h."
}

func (v timeoutDurationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (timeoutDurationValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	duration, err := time.ParseDuration(req.ConfigValue.ValueString())
	if err != nil || duration <= 0 {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid operation timeout", "Use a positive Go duration, for example 10m or 1h.")
	}
}

// TimeoutDuration resolves the configured duration for operation, falling back
// to fallback when the block or the operation is not set.
func TimeoutDuration(value types.Object, operation string, fallback time.Duration) (time.Duration, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	if value.IsNull() {
		return fallback, diagnostics
	}
	if value.IsUnknown() {
		diagnostics.AddError("Unknown operation timeout", "Resolve the timeouts block before applying.")
		return 0, diagnostics
	}
	raw, ok := value.Attributes()[operation]
	if !ok || raw.IsNull() {
		return fallback, diagnostics
	}
	configured, ok := raw.(types.String)
	if !ok || configured.IsUnknown() {
		diagnostics.AddError("Invalid operation timeout", "Resolve the operation timeout to a duration string before applying.")
		return 0, diagnostics
	}
	duration, err := time.ParseDuration(configured.ValueString())
	if err != nil || duration <= 0 {
		diagnostics.AddError("Invalid operation timeout", "Use a positive Go duration, for example 10m or 1h.")
		return 0, diagnostics
	}
	return duration, diagnostics
}
