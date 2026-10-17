package paistudio

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// userVpcObjectModel mirrors the UserVpc component schema. Its tfsdk tags must
// match userVpcAttrTypes so types.Object reflection round-trips correctly.
type userVpcObjectModel struct {
	VpcId              types.String `tfsdk:"vpc_id"`
	SwitchId           types.String `tfsdk:"switch_id"`
	SecurityGroupId    types.String `tfsdk:"security_group_id"`
	RoleArn            types.String `tfsdk:"role_arn"`
	DefaultRoute       types.String `tfsdk:"default_route"`
	ExtendedCIDRs      types.List   `tfsdk:"extended_cidrs"`
	DefaultForwardInfo types.Object `tfsdk:"default_forward_info"`
}

// forwardInfoObjectModel mirrors the ForwardInfo component schema.
type forwardInfoObjectModel struct {
	EipAllocationId types.String `tfsdk:"eip_allocation_id"`
	NatGatewayId    types.String `tfsdk:"nat_gateway_id"`
}

// Attribute type maps for the nested object attributes, used to round-trip
// types.Object values through reflection via ObjectFrom / As.

var forwardInfoAttrTypes = map[string]attr.Type{
	"eip_allocation_id": types.StringType,
	"nat_gateway_id":    types.StringType,
}

var userVpcAttrTypes = map[string]attr.Type{
	"vpc_id":            types.StringType,
	"switch_id":         types.StringType,
	"security_group_id": types.StringType,
	"role_arn":          types.StringType,
	"default_route":     types.StringType,
	"extended_cidrs":    types.ListType{ElemType: types.StringType},
	"default_forward_info": types.ObjectType{
		AttrTypes: forwardInfoAttrTypes,
	},
}

// tagsMapToTagArray converts a configured tags map into the CreateResourceGroup
// Tag request shape: [{Key, Value}, ...].
func tagsMapToTagArray(m types.Map) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(m.Elements()))
	for k, v := range m.Elements() {
		if sv, ok := v.(types.String); ok && !sv.IsNull() && !sv.IsUnknown() {
			out = append(out, map[string]interface{}{
				"Key":   k,
				"Value": sv.ValueString(),
			})
		}
	}
	return out
}

// tagsArrayToMap converts the GetResourceGroup Tags response shape
// [{TagKey, TagValue}, ...] into a computed tags map.
func tagsArrayToMap(tags []interface{}) types.Map {
	goMap := make(map[string]string, len(tags))
	for _, raw := range tags {
		t, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		k, _ := t["TagKey"].(string)
		v, _ := t["TagValue"].(string)
		if k == "" {
			continue
		}
		goMap[k] = v
	}
	v, diags := types.MapValueFrom(context.Background(), types.StringType, goMap)
	if diags.HasError() {
		return types.MapNull(types.StringType)
	}
	return v
}

// userVpcToObject converts a configured user_vpc object into the API map shape
// (PascalCase keys), skipping null/unknown sub-values. Returns ok=false when
// the object could not be decoded.
func userVpcToObject(ctx context.Context, obj types.Object, diags *diag.Diagnostics) (map[string]interface{}, bool) {
	var uv userVpcObjectModel
	diags.Append(obj.As(ctx, &uv, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, false
	}
	out := map[string]interface{}{}
	if v := uv.VpcId; !v.IsNull() && !v.IsUnknown() {
		out["VpcId"] = v.ValueString()
	}
	if v := uv.SwitchId; !v.IsNull() && !v.IsUnknown() {
		out["SwitchId"] = v.ValueString()
	}
	if v := uv.SecurityGroupId; !v.IsNull() && !v.IsUnknown() {
		out["SecurityGroupId"] = v.ValueString()
	}
	if v := uv.RoleArn; !v.IsNull() && !v.IsUnknown() {
		out["RoleArn"] = v.ValueString()
	}
	if v := uv.DefaultRoute; !v.IsNull() && !v.IsUnknown() {
		out["DefaultRoute"] = v.ValueString()
	}
	if !uv.ExtendedCIDRs.IsNull() && !uv.ExtendedCIDRs.IsUnknown() {
		cidrs := make([]string, 0, len(uv.ExtendedCIDRs.Elements()))
		for _, e := range uv.ExtendedCIDRs.Elements() {
			if sv, ok := e.(types.String); ok && !sv.IsNull() && !sv.IsUnknown() {
				cidrs = append(cidrs, sv.ValueString())
			}
		}
		if len(cidrs) > 0 {
			out["ExtendedCIDRs"] = cidrs
		}
	}
	if !uv.DefaultForwardInfo.IsNull() && !uv.DefaultForwardInfo.IsUnknown() {
		if fi, ok := forwardInfoToObject(ctx, uv.DefaultForwardInfo, diags); ok {
			out["DefaultForwardInfo"] = fi
		}
		if diags.HasError() {
			return nil, false
		}
	}
	return out, true
}

func forwardInfoToObject(ctx context.Context, obj types.Object, diags *diag.Diagnostics) (map[string]interface{}, bool) {
	var fi forwardInfoObjectModel
	diags.Append(obj.As(ctx, &fi, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, false
	}
	out := map[string]interface{}{}
	if v := fi.EipAllocationId; !v.IsNull() && !v.IsUnknown() {
		out["EipAllocationId"] = v.ValueString()
	}
	if v := fi.NatGatewayId; !v.IsNull() && !v.IsUnknown() {
		out["NatGatewayId"] = v.ValueString()
	}
	return out, true
}

// userVpcFromObject converts the GetResourceGroup UserVpc response map into a
// computed types.Object matching the user_vpc attribute schema.
func userVpcFromObject(obj map[string]interface{}) types.Object {
	str := func(key string) types.String {
		if v, ok := obj[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	uv := userVpcObjectModel{
		VpcId:           str("VpcId"),
		SwitchId:        str("SwitchId"),
		SecurityGroupId: str("SecurityGroupId"),
		RoleArn:         str("RoleArn"),
		DefaultRoute:    str("DefaultRoute"),
		// Initialize the complex sub-values to typed nulls so the
		// types.ObjectFrom reflection below sees well-formed values when the
		// API omits them, instead of zero types.Object/List values whose
		// type does not match the attribute schema.
		ExtendedCIDRs:      types.ListNull(types.StringType),
		DefaultForwardInfo: types.ObjectNull(forwardInfoAttrTypes),
	}
	if raw, ok := obj["ExtendedCIDRs"].([]interface{}); ok && len(raw) > 0 {
		cidrs := make([]string, 0, len(raw))
		for _, c := range raw {
			if s, ok := c.(string); ok {
				cidrs = append(cidrs, s)
			}
		}
		if v, d := types.ListValueFrom(context.Background(), types.StringType, cidrs); !d.HasError() {
			uv.ExtendedCIDRs = v
		}
	}
	if fi, ok := obj["DefaultForwardInfo"].(map[string]interface{}); ok && len(fi) > 0 {
		uv.DefaultForwardInfo = forwardInfoFromObject(fi)
	}
	v, diags := types.ObjectValueFrom(context.Background(), userVpcAttrTypes, &uv)
	if diags.HasError() {
		return types.ObjectNull(userVpcAttrTypes)
	}
	return v
}

func forwardInfoFromObject(obj map[string]interface{}) types.Object {
	str := func(key string) types.String {
		if v, ok := obj[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	fi := forwardInfoObjectModel{
		EipAllocationId: str("EipAllocationId"),
		NatGatewayId:    str("NatGatewayId"),
	}
	v, diags := types.ObjectValueFrom(context.Background(), forwardInfoAttrTypes, &fi)
	if diags.HasError() {
		return types.ObjectNull(forwardInfoAttrTypes)
	}
	return v
}
