package wafopenapi_test

import (
	"context"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/wafopenapi"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// TestUnitAliCloudWafDefenseAssetSchema asserts the defense asset schema invariants that the API
// contract makes load-bearing. The authoritative online (zhenyuan) DefenseAsset definitions fix:
//
//   - asset_type is a Create enum (AI, Client_Injection, Apisec, Account_Protector, SSO) and is
//     not accepted by ModifyDefenseAsset, so a change must replace the resource;
//   - instance_id is required on every DefenseAsset call and binds the asset to its WAF instance,
//     so a change must replace the resource;
//   - resource and conditions are required on Create and not accepted by Modify, so they are
//     immutable;
//   - asset_config is required on both Create and Modify;
//   - id is server-assigned and never part of the user contract.
//
// The asset_type RequiresReplace guard is the regression for the previous round's failure, where
// asset_type was not configured with RequiresReplace and a change was planned as an in-place
// update — which then drove ModifyDefenseAsset with a type the asset never had.
func TestUnitAliCloudWafDefenseAssetSchema(t *testing.T) {
	r := wafopenapi.NewWafDefenseAssetResource()

	var metadata fwresource.MetadataResponse
	r.Metadata(context.Background(), fwresource.MetadataRequest{ProviderTypeName: "alicloud"}, &metadata)
	if metadata.TypeName != "alicloud_waf_defense_asset" {
		t.Fatalf("TypeName = %q, want alicloud_waf_defense_asset", metadata.TypeName)
	}

	var schemaResp fwresource.SchemaResponse
	r.Schema(context.Background(), fwresource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building schema: %v", schemaResp.Diagnostics)
	}

	objectType, ok := schemaResp.Schema.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatal("resource schema does not describe an object")
	}

	// asset_type is the identity-defining attribute the API validates: changing it must replace the
	// resource. This is the regression guard for the previous round's failure, where asset_type was
	// not configured with RequiresReplace and a change was planned as an in-place update.
	assetType, ok := schemaResp.Schema.Attributes["asset_type"].(schema.StringAttribute)
	if !ok {
		t.Fatal("asset_type attribute is missing or is not a string attribute")
	}
	if !assetType.Required || len(assetType.PlanModifiers) != 1 {
		t.Fatalf("asset_type must stay Required with one plan modifier, got Required=%v modifiers=%d", assetType.Required, len(assetType.PlanModifiers))
	}
	object := func(assetType string) tftypes.Value {
		return tftypes.NewValue(objectType, map[string]tftypes.Value{
			"id":           tftypes.NewValue(tftypes.String, "asset-123"),
			"instance_id":  tftypes.NewValue(tftypes.String, "waf_v2_public_cn-jyi4sglpj03"),
			"asset_type":   tftypes.NewValue(tftypes.String, assetType),
			"resource":     tftypes.NewValue(tftypes.String, "res-123"),
			"conditions":   tftypes.NewValue(tftypes.String, `[{"Domain":"example.com","UrlPath":"/","HttpMethod":"GET"}]`),
			"asset_config": tftypes.NewValue(tftypes.String, `{}`),
			"asset_name":   tftypes.NewValue(tftypes.String, nil),
			"description":  tftypes.NewValue(tftypes.String, nil),
			"asset_status": tftypes.NewValue(tftypes.String, "1"),
			"timeouts":     tftypes.NewValue(objectType.AttributeTypes["timeouts"], nil),
		})
	}
	var modifierResp planmodifier.StringResponse
	assetType.PlanModifiers[0].PlanModifyString(context.Background(), planmodifier.StringRequest{
		State:      tfsdk.State{Schema: schemaResp.Schema, Raw: object("AI")},
		Plan:       tfsdk.Plan{Schema: schemaResp.Schema, Raw: object("Apisec")},
		StateValue: types.StringValue("AI"),
		PlanValue:  types.StringValue("Apisec"),
	}, &modifierResp)
	if !modifierResp.RequiresReplace {
		t.Fatal("asset_type plan modifier did not request replacement on change")
	}

	// instance_id is required on every DefenseAsset call, so a change must also replace the resource.
	instanceId, ok := schemaResp.Schema.Attributes["instance_id"].(schema.StringAttribute)
	if !ok {
		t.Fatal("instance_id attribute is missing or is not a string attribute")
	}
	if !instanceId.Required || len(instanceId.PlanModifiers) != 1 {
		t.Fatalf("instance_id must stay Required with one plan modifier, got Required=%v modifiers=%d", instanceId.Required, len(instanceId.PlanModifiers))
	}
	var instanceModifierResp planmodifier.StringResponse
	instanceId.PlanModifiers[0].PlanModifyString(context.Background(), planmodifier.StringRequest{
		State:      tfsdk.State{Schema: schemaResp.Schema, Raw: object("AI")},
		Plan:       tfsdk.Plan{Schema: schemaResp.Schema, Raw: object("AI")},
		StateValue: types.StringValue("waf_v2_public_cn-jyi4sglpj03"),
		PlanValue:  types.StringValue("waf_v2_public_cn-jyi4sglpj04"),
	}, &instanceModifierResp)
	if !instanceModifierResp.RequiresReplace {
		t.Fatal("instance_id plan modifier did not request replacement on change")
	}

	// resource and conditions are Create-only (not accepted by ModifyDefenseAsset): both must be
	// Required and request replacement on change.
	for _, name := range []string{"resource", "conditions"} {
		attr, ok := schemaResp.Schema.Attributes[name].(schema.StringAttribute)
		if !ok {
			t.Fatalf("%s attribute is missing or is not a string attribute", name)
		}
		if !attr.Required {
			t.Fatalf("%s must be Required, got Required=%v", name, attr.Required)
		}
		var hasReplace bool
		for _, m := range attr.PlanModifiers {
			var resp planmodifier.StringResponse
			m.PlanModifyString(context.Background(), planmodifier.StringRequest{
				State:      tfsdk.State{Schema: schemaResp.Schema, Raw: object("AI")},
				Plan:       tfsdk.Plan{Schema: schemaResp.Schema, Raw: object("AI")},
				StateValue: types.StringValue("v1"),
				PlanValue:  types.StringValue("v2"),
			}, &resp)
			if resp.RequiresReplace {
				hasReplace = true
			}
		}
		if !hasReplace {
			t.Fatalf("%s must have a RequiresReplace plan modifier", name)
		}
	}

	// asset_config is required on both Create and Modify.
	assetConfig, ok := schemaResp.Schema.Attributes["asset_config"].(schema.StringAttribute)
	if !ok {
		t.Fatal("asset_config attribute is missing or is not a string attribute")
	}
	if !assetConfig.Required {
		t.Fatalf("asset_config must be Required, got Required=%v", assetConfig.Required)
	}

	// id is server-assigned and must not be part of the user contract.
	idAttr, ok := schemaResp.Schema.Attributes["id"].(schema.StringAttribute)
	if !ok {
		t.Fatal("id attribute is missing or is not a string attribute")
	}
	if !idAttr.Computed || idAttr.Required || idAttr.Optional {
		t.Fatalf("id must stay Computed-only, got Computed=%v Required=%v Optional=%v", idAttr.Computed, idAttr.Required, idAttr.Optional)
	}
}
