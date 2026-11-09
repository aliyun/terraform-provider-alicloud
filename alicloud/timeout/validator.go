package timeout

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type PositiveTimeoutValidator struct{}

func (PositiveTimeoutValidator) Description(context.Context) string {
	return "A positive Go duration, for example 10m or 1h."
}

func (v PositiveTimeoutValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (PositiveTimeoutValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	duration, err := time.ParseDuration(req.ConfigValue.ValueString())
	if err != nil || duration > 0 {
		return
	}
	resp.Diagnostics.AddAttributeError(req.Path, "Invalid operation timeout", "Use a positive Go duration, for example 10m or 1h.")
}
