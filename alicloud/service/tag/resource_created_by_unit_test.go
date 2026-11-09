package tag

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestUnitCreatedByTimeoutsSchema(t *testing.T) {
	r := &createdByResource{}
	resp := resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	block, ok := resp.Schema.Blocks["timeouts"].(schema.SingleNestedBlock)
	if !ok {
		t.Fatalf("timeouts block has type %T", resp.Schema.Blocks["timeouts"])
	}
	if _, ok := block.CustomType.(timeouts.Type); !ok {
		t.Fatalf("timeouts CustomType = %T, want timeouts.Type", block.CustomType)
	}
	if len(block.Attributes) != 3 {
		t.Fatalf("timeouts has %d attributes, want 3", len(block.Attributes))
	}
	for _, operation := range []string{"create", "update", "delete"} {
		if _, ok := block.Attributes[operation]; !ok {
			t.Errorf("timeouts is missing %q", operation)
		}
	}
}

func TestUnitCreatedByTimeoutsValidation(t *testing.T) {
	r := &createdByResource{}
	resp := resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	block := resp.Schema.Blocks["timeouts"].(schema.SingleNestedBlock)

	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{name: "positive", value: "10m"},
		{name: "malformed", value: "soon", wantError: true},
		{name: "zero", value: "0s"},
		{name: "negative", value: "-1m"},
	}
	for _, operation := range []string{"create", "update", "delete"} {
		attribute := block.Attributes[operation].(schema.StringAttribute)
		for _, tt := range tests {
			t.Run(operation+"/"+tt.name, func(t *testing.T) {
				hasError := false
				for _, attributeValidator := range attribute.Validators {
					validatorResp := validator.StringResponse{}
					attributeValidator.ValidateString(context.Background(), validator.StringRequest{
						Path:        path.Root("timeouts").AtName(operation),
						ConfigValue: types.StringValue(tt.value),
					}, &validatorResp)
					hasError = hasError || validatorResp.Diagnostics.HasError()
				}
				if hasError != tt.wantError {
					t.Fatalf("validation error = %v, want %v", hasError, tt.wantError)
				}
			})
		}
	}
}

func TestUnitCreatedByTimeoutsStateCompatibility(t *testing.T) {
	ctx := context.Background()
	r := &createdByResource{}
	schemaResp := resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	timeoutsType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"create": tftypes.String,
		"update": tftypes.String,
		"delete": tftypes.String,
	}}
	stateType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"id":          tftypes.String,
		"open_status": tftypes.Bool,
		"timeouts":    timeoutsType,
	}}
	decode := func(t *testing.T, rawTimeouts tftypes.Value) createdByResourceModel {
		t.Helper()
		state := tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":          tftypes.NewValue(tftypes.String, "cn-shanghai"),
			"open_status": tftypes.NewValue(tftypes.Bool, true),
			"timeouts":    rawTimeouts,
		})}
		var model createdByResourceModel
		if diagnostics := state.Get(ctx, &model); diagnostics.HasError() {
			t.Fatalf("decoding legacy-shape state: %v", diagnostics)
		}
		return model
	}

	t.Run("null", func(t *testing.T) {
		model := decode(t, tftypes.NewValue(timeoutsType, nil))
		if !model.Timeouts.IsNull() {
			t.Fatal("null timeouts did not remain null")
		}
	})
	t.Run("configured", func(t *testing.T) {
		model := decode(t, tftypes.NewValue(timeoutsType, map[string]tftypes.Value{
			"create": tftypes.NewValue(tftypes.String, "10m"),
			"update": tftypes.NewValue(tftypes.String, nil),
			"delete": tftypes.NewValue(tftypes.String, "15m"),
		}))
		createTimeout, diagnostics := model.Timeouts.Create(ctx, 5*time.Minute)
		if diagnostics.HasError() || createTimeout != 10*time.Minute {
			t.Fatalf("create timeout = %s, diagnostics = %v", createTimeout, diagnostics)
		}
		deleteTimeout, diagnostics := model.Timeouts.Delete(ctx, 5*time.Minute)
		if diagnostics.HasError() || deleteTimeout != 15*time.Minute {
			t.Fatalf("delete timeout = %s, diagnostics = %v", deleteTimeout, diagnostics)
		}
		updateTimeout, diagnostics := model.Timeouts.Update(ctx, 5*time.Minute)
		if diagnostics.HasError() || updateTimeout != 5*time.Minute {
			t.Fatalf("update timeout = %s, diagnostics = %v", updateTimeout, diagnostics)
		}
	})
}
