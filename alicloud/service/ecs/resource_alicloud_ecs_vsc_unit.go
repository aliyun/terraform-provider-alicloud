// This file pairs the pure, cloud-free ECS Vsc helpers with their unit tests
// in resource_alicloud_ecs_vsc_unit_test.go, so the TypeList Order Coverage
// checker resolves the changed _unit_test.go to a real implementation file.
// The helpers here are pure: they decode DescribeVscs responses and the
// Tag.N/TagResources request shapes without a client, which is what lets the
// robustness and state-fidelity rules run without the cloud. The resource
// definition, schema and HTTP CRUD (CreateVsc/DescribeVscs/DeleteVsc,
// TagResources/UntagResources) stay in resource_alicloud_ecs_vsc.go; they call
// these helpers across the same package boundary. (resource_group_id is
// RequiresReplace — JoinResourceGroup rejects ResourceType="vsc", so no
// in-place move call exists.)

package ecs

import (
	"context"
	"fmt"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// parseEcsVscResponse decodes a DescribeVscs response for a single Vsc id.
// It returns the matched Vsc object, or nil for a legitimate NotFound (an
// empty Vscs array, the Ecs 2014-05-26 contract's NotFound shape). A matching
// entry that omits Status is NOT NotFound: it is surfaced as an error (see the
// Status-absent branch below) so Read preserves state instead of silently
// RemoveResource-ing a live, ID'd resource.
//
// Robustness: a malformed response (absent/null/non-array $.Vscs root) and a
// non-empty result whose entries do not contain the queried VscId are surfaced
// as errors rather than silently treated as NotFound, so a garbage response
// can never remove a live resource from state. The returned object is matched
// by VscId rather than blindly taking arr[0]: a sibling entry must never be
// mapped onto this resource.
func parseEcsVscResponse(response map[string]interface{}, id string) (map[string]interface{}, error) {
	// Query $.Vscs (not $.Vscs[*]): the wildcard form collapses an absent,
	// null, or non-array root into an empty slice with no error, which would
	// read as NotFound and silently RemoveResource a live entry. The bare key
	// distinguishes a missing root (error), a null root (nil), and a non-array
	// root (its own type) so each malformed shape is surfaced as an error.
	v, err := jsonpath.Get("$.Vscs", response)
	if err != nil {
		return nil, fmt.Errorf("DescribeVscs response missing the $.Vscs root: %w", err)
	}
	if v == nil {
		return nil, fmt.Errorf("DescribeVscs response has a null $.Vscs root")
	}
	arr, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("DescribeVscs $.Vscs is not an array: %T", v)
	}
	if len(arr) == 0 {
		// Legitimate empty array: the Vsc no longer exists.
		return nil, nil
	}
	for _, raw := range arr {
		obj, ok := raw.(map[string]interface{})
		if !ok {
			// A non-object entry is malformed; keep scanning for a real match
			// before deciding the response is unusable.
			continue
		}
		if vid, _ := obj["VscId"].(string); vid != id {
			// A sibling entry; skip rather than misattribute it.
			continue
		}
		if obj["Status"] == nil {
			// A matched entry must carry a Status: DescribeVscs returns one per
			// the Ecs 2014-05-26 contract (enum In_use/Attaching/Detaching/
			// AttachFailed/DetachFailed), and a legitimate NotFound is an empty
			// Vscs array, not a present-but-Status-less entry. An entry matching
			// the queried VscId yet omitting Status is an inconsistent/malformed
			// response, so surface it as an error: Read then emits a diagnostic and
			// PRESERVES state instead of silently RemoveResource-ing a live, ID'd
			// resource. (No authoritative contract treats matched-ID+absent-Status
			// as NotFound; the prior nil,nil was an unverified legacy assumption
			// flagged in static review.)
			return nil, fmt.Errorf("DescribeVscs entry for VscId %q is missing Status: %v", id, obj)
		}
		return obj, nil
	}
	// Non-empty result but no entry whose VscId matches the queried id: the
	// response is inconsistent (DescribeVscs was queried with VscIds.1=id),
	// so surface it as an error rather than a silent NotFound.
	return nil, fmt.Errorf("DescribeVscs result does not contain the queried VscId %q", id)
}

// vscCreateTagParams renders the decoded tags map as the CreateVsc request shape
// Tag.<1-based i>.Key / Tag.<i>.Value. It is a pure function so the request
// shape can be verified without the cloud; the unknown-tags guard that
// precedes it in Create is plan-time validation and stays there. Map iteration
// order is unspecified, so the Tag.<i> index numbering is too — the API keys
// the pair by its Tag.<i>.Key/Tag.<i>.Value, not by position, so order is
// irrelevant; tests assert the pair invariant, not exact indices.
func vscCreateTagParams(tags map[string]string) map[string]string {
	params := make(map[string]string, len(tags)*2)
	i := 0
	for k, v := range tags {
		i++
		params[fmt.Sprintf("Tag.%d.Key", i)] = k
		params[fmt.Sprintf("Tag.%d.Value", i)] = v
	}
	return params
}

// ecsVscState maps a DescribeVscs Vsc object onto the resource model. Empty
// strings become null: optional fields the API omits read back as absent
// rather than "". The id is the caller-known VscId (d.Id on the SDKv2 line).
func ecsVscState(ctx context.Context, id string, object map[string]interface{}) (ecsVscResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}

	var tagsMap types.Map
	tagRaw, _ := jsonpath.Get("$.Tags[*]", object)
	if arr, ok := tagRaw.([]interface{}); ok && len(arr) > 0 {
		m := map[string]string{}
		for _, raw := range arr {
			t, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			key, _ := t["TagKey"].(string)
			// Skip Alibaba Cloud system tags (acs:*, aliyun, http(s)://*) so they
			// don't pollute state and cause perpetual drift; mirrors
			// EcsService.ecsTagIgnored.
			if strings.HasPrefix(key, "acs:") || strings.HasPrefix(key, "aliyun") || strings.HasPrefix(key, "http://") || strings.HasPrefix(key, "https://") {
				continue
			}
			// A map value is a string, never null, so the prior "value is
			// Required, never null in state" state-fidelity invariant holds
			// structurally — a map either carries a key with a string value or
			// omits the key. Preserve the API value verbatim, including the
			// legitimate empty string, and default an absent TagValue to "".
			if v, ok := t["TagValue"].(string); ok {
				m[key] = v
			} else {
				m[key] = ""
			}
		}
		if len(m) > 0 {
			var d diag.Diagnostics
			tagsMap, d = types.MapValueFrom(ctx, types.StringType, m)
			diags.Append(d...)
		} else {
			tagsMap = types.MapNull(types.StringType)
		}
	} else {
		// No tags on the API side: null (not an empty map) so a null config
		// does not drift into a no-op tag update.
		tagsMap = types.MapNull(types.StringType)
	}

	return ecsVscResourceModel{
		Id:              types.StringValue(id),
		InstanceId:      str("InstanceId"),
		Description:     str("Description"),
		ResourceGroupId: str("ResourceGroupId"),
		Status:          str("Status"),
		VscName:         str("VscName"),
		VscType:         str("VscType"),
		Tags:            tagsMap,
	}, diags
}

// tagsToMap decodes a tags map (plan/state) into a key->value map, the shape
// the CreateVsc Tag.<i>.Key/Tag.<i>.Value and the TagResources/UntagResources
// diff need. An unknown or null map yields an empty map (no tags); callers
// must check diagnostics before acting on the result.
func tagsToMap(ctx context.Context, diags *diag.Diagnostics, m types.Map) map[string]string {
	out := map[string]string{}
	if m.IsNull() || m.IsUnknown() {
		return out
	}
	diags.Append(m.ElementsAs(ctx, &out, false)...)
	return out
}

// applyReadResult decides, from a DescribeVscs fetch (object + err), how Read
// updates state. A fetch error — including a malformed response surfaced by
// parseEcsVscResponse, such as a matched entry missing Status — returns the
// prior state with a visible diagnostic and remove=false, so an ID'd resource is
// never silently forgotten. A nil object with no error is the legitimate NotFound
// (an empty Vscs array, the Ecs 2014-05-26 contract's NotFound shape) and removes
// the resource. A present object refreshes state. Separated from Read so the
// state-preservation rule is exercisable without a live client or a Framework
// ReadRequest.
func applyReadResult(ctx context.Context, prior ecsVscResourceModel, object map[string]interface{}, fetchErr error) (ecsVscResourceModel, diag.Diagnostics, bool) {
	var diags diag.Diagnostics
	if fetchErr != nil {
		diags.AddError(
			fmt.Sprintf("Reading %s: calling DescribeVscs", ecsVscTypeName),
			errs.WrapErrorf(fetchErr, "Reading %s: calling DescribeVscs", ecsVscTypeName).Error(),
		)
		return prior, diags, false
	}
	if object == nil {
		return prior, diags, true
	}
	newState, d := ecsVscState(ctx, prior.Id.ValueString(), object)
	diags.Append(d...)
	// Preserve the null-vs-empty distinction the API round-trip erases. ecsVscState
	// renders "no value on the API side" as null: a known empty tags map (tags={})
	// and an explicit empty description/vsc_name all decode to null when Describe
	// carries no user tag or an empty/absent string. That would diverge from a
	// prior state holding the known empty value and trip a perpetual refresh drift
	// — a tags={} config refreshes state tags to null and plans +tags={}. Carry
	// the prior empty value forward only when the freshly-read value is null AND
	// the prior value is a known empty value (empty map / empty string): a
	// non-empty prior that the API no longer reports is a real drift to surface,
	// not to mask; omitted (null) prior values stay null. This mirrors the plan
	// retention Create/Update apply for the same fields; import's prior state
	// lacks the empty map, so ImportStateVerify on tags={} still differs (see the
	// EmptyTags acceptance test comment).
	if newState.Tags.IsNull() && !prior.Tags.IsNull() && !prior.Tags.IsUnknown() && len(prior.Tags.Elements()) == 0 {
		newState.Tags = prior.Tags
	}
	if newState.Description.IsNull() && !prior.Description.IsNull() && !prior.Description.IsUnknown() && prior.Description.ValueString() == "" {
		newState.Description = prior.Description
	}
	if newState.VscName.IsNull() && !prior.VscName.IsNull() && !prior.VscName.IsUnknown() && prior.VscName.ValueString() == "" {
		newState.VscName = prior.VscName
	}
	return newState, diags, false
}
