package ga

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Object attribute types for the nested sets, used when constructing
// types.Set values from API responses and when building empty sets.
var (
	componentObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"service_component_id": types.StringType,
			"priority":             types.Int64Type,
			"timeout":              types.Int64Type,
			"config":               types.StringType,
			"fail_policy":          types.StringType,
			"component_name":       types.StringType,
		},
	}
	resourceObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"accelerator_id": types.StringType,
			"resource_type":  types.StringType,
			"resource_id":    types.StringType,
			"associate_id":   types.StringType,
		},
	}
)

// stringValue maps an API field to a types.String, with empty values read as
// null — matching the established ims Framework convention so Optional and
// Computed attributes do not surface "" where the API omits them.
func stringValue(v interface{}) types.String {
	switch x := v.(type) {
	case string:
		if x == "" {
			return types.StringNull()
		}
		return types.StringValue(x)
	case nil:
		return types.StringNull()
	default:
		s := fmt.Sprint(x)
		if s == "" {
			return types.StringNull()
		}
		return types.StringValue(s)
	}
}

// int64Value maps an API numeric field (float64 from JSON) to types.Int64.
func int64Value(v interface{}) types.Int64 {
	switch x := v.(type) {
	case float64:
		return types.Int64Value(int64(x))
	case int:
		return types.Int64Value(int64(x))
	case int64:
		return types.Int64Value(x)
	case string:
		if n, err := strconv.ParseInt(x, 10, 64); err == nil {
			return types.Int64Value(n)
		}
		return types.Int64Null()
	case nil:
		return types.Int64Null()
	default:
		return types.Int64Null()
	}
}

func componentSet(ctx context.Context, items []interface{}) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	elements := make([]attr.Value, 0, len(items))
	for _, raw := range items {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		obj, d := types.ObjectValue(componentObjectType.AttrTypes, map[string]attr.Value{
			"service_component_id": stringValue(item["ServiceComponentId"]),
			"priority":             int64Value(item["Priority"]),
			"timeout":              int64Value(item["Timeout"]),
			"config":               stringValue(item["Config"]),
			"fail_policy":          stringValue(item["FailPolicy"]),
			"component_name":       stringValue(item["ComponentName"]),
		})
		diags.Append(d...)
		elements = append(elements, obj)
	}
	set, d := types.SetValue(componentObjectType, elements)
	diags.Append(d...)
	return set, diags
}

func resourceSet(ctx context.Context, items []interface{}) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	elements := make([]attr.Value, 0, len(items))
	for _, raw := range items {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		obj, d := types.ObjectValue(resourceObjectType.AttrTypes, map[string]attr.Value{
			"accelerator_id": stringValue(item["AcceleratorId"]),
			"resource_type":  stringValue(item["ResourceType"]),
			"resource_id":    stringValue(item["ResourceId"]),
			"associate_id":   stringValue(item["AssociateId"]),
		})
		diags.Append(d...)
		elements = append(elements, obj)
	}
	set, d := types.SetValue(resourceObjectType, elements)
	diags.Append(d...)
	return set, diags
}

func componentsFromSet(ctx context.Context, s types.Set) []componentModel {
	var items []componentModel
	_ = s.ElementsAs(ctx, &items, false)
	return items
}

func resourcesFromSet(ctx context.Context, s types.Set) []extensionResourceModel {
	var items []extensionResourceModel
	_ = s.ElementsAs(ctx, &items, false)
	return items
}

func indexComponents(items []componentModel) map[string]componentModel {
	m := make(map[string]componentModel, len(items))
	for _, c := range items {
		m[c.ServiceComponentId.ValueString()] = c
	}
	return m
}

func indexResources(items []extensionResourceModel) map[string]extensionResourceModel {
	m := make(map[string]extensionResourceModel, len(items))
	for _, r := range items {
		m[r.AcceleratorId.ValueString()+"|"+r.ResourceType.ValueString()+"|"+r.ResourceId.ValueString()] = r
	}
	return m
}

func componentsEqual(a, b componentModel) bool {
	return a.ServiceComponentId.Equal(b.ServiceComponentId) &&
		a.Priority.Equal(b.Priority) &&
		a.Timeout.Equal(b.Timeout) &&
		a.Config.Equal(b.Config) &&
		a.FailPolicy.Equal(b.FailPolicy)
}

func tagsToMapFromState(ctx context.Context, m types.Map) map[string]string {
	out := make(map[string]string)
	_ = m.ElementsAs(ctx, &out, false)
	return out
}

func tagsMapToValue(ctx context.Context, tags map[string]string, diags *diag.Diagnostics) types.Map {
	if len(tags) == 0 {
		set, d := types.MapValue(types.StringType, map[string]attr.Value{})
		diags.Append(d...)
		return set
	}
	elements := make(map[string]attr.Value, len(tags))
	for k, v := range tags {
		elements[k] = types.StringValue(v)
	}
	set, d := types.MapValue(types.StringType, elements)
	diags.Append(d...)
	return set
}
