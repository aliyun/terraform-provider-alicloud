package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const securityGroupSnapshotPolicyTypeName = "alicloud_security_group_snapshot_policy"

// createDescribeMaxAttempts bounds the post-Create refresh retries that cover
// the eventual-consistency window in which a just-created policy is not yet
// visible via Describe. createDescribeRetryInterval is the wait between
// attempts. A real transport/API error is never retried.
const (
	createDescribeMaxAttempts   = 5
	createDescribeRetryInterval = 2 * time.Second
)

var (
	_ resource.Resource                = &securityGroupSnapshotPolicyResource{}
	_ resource.ResourceWithImportState = &securityGroupSnapshotPolicyResource{}
)

func NewSecurityGroupSnapshotPolicyResource() resource.Resource {
	return &securityGroupSnapshotPolicyResource{}
}

type securityGroupSnapshotPolicyResource struct {
	fwadapt.ResourceBase
}

type securityGroupSnapshotPolicyResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	SecurityGroupSnapshotPolicyId   types.String `tfsdk:"security_group_snapshot_policy_id"`
	SecurityGroupSnapshotPolicyName types.String `tfsdk:"security_group_snapshot_policy_name"`
	Status                          types.String `tfsdk:"status"`
	RegionId                        types.String `tfsdk:"region_id"`
}

func (r *securityGroupSnapshotPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = securityGroupSnapshotPolicyTypeName
}

func (r *securityGroupSnapshotPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a ECS Security Group Snapshot Policy resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `security_group_snapshot_policy_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"security_group_snapshot_policy_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the security group snapshot policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"security_group_snapshot_policy_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the security group snapshot policy.",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The status of the security group snapshot policy. Valid values: `enable` and `disable`.",
				Validators: []validator.String{
					stringvalidator.OneOf("enable", "disable"),
				},
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the security group snapshot policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *securityGroupSnapshotPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan securityGroupSnapshotPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"RegionId": r.Client().RegionId,
	}
	if v := plan.SecurityGroupSnapshotPolicyName.ValueString(); v != "" && !plan.SecurityGroupSnapshotPolicyName.IsUnknown() {
		request["SecurityGroupSnapshotPolicyName"] = v
	}
	if v := plan.Status.ValueString(); v != "" && !plan.Status.IsUnknown() {
		request["Status"] = v
	}

	response, err := r.Client().RpcPost("Ecs", "2016-03-14", "CreateSecurityGroupSnapshotPolicy", nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateSecurityGroupSnapshotPolicy", securityGroupSnapshotPolicyTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateSecurityGroupSnapshotPolicy", securityGroupSnapshotPolicyTypeName).Error(),
		)
		return
	}

	id, ok := extractSecurityGroupSnapshotPolicyId(response)
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateSecurityGroupSnapshotPolicy", securityGroupSnapshotPolicyTypeName),
			fmt.Sprintf("The response carried no SecurityGroupSnapshotPolicyId: %v", response),
		)
		return
	}

	// Persist the recoverable id before the post-Create refresh. If the Describe
	// that follows errors, or the resource is not yet visible (eventual
	// consistency), the created cloud resource is still tracked: the next Read
	// recovers from this id instead of leaking a resource Terraform cannot
	// address (R4.2).
	plan.Id = types.StringValue(id)
	plan.SecurityGroupSnapshotPolicyId = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	// A freshly created policy may take a moment to become visible via Describe;
	// retry for a bounded window rather than failing the create at the first
	// not-yet-visible read. A real error or a malformed response is surfaced
	// immediately (fail-closed); only an empty-but-well-formed list is retried.
	object, err := describeSecurityGroupSnapshotPolicyVisible(func() (map[string]interface{}, error) {
		return describeSecurityGroupSnapshotPolicy(r.Client(), id)
	}, createDescribeMaxAttempts, createDescribeRetryInterval)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeSecurityGroupSnapshotPolicies", securityGroupSnapshotPolicyTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeSecurityGroupSnapshotPolicies", securityGroupSnapshotPolicyTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateSecurityGroupSnapshotPolicy", securityGroupSnapshotPolicyTypeName),
			fmt.Sprintf("The security group snapshot policy %s is not visible via DescribeSecurityGroupSnapshotPolicies right after creation", id),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, securityGroupSnapshotPolicyState(object, r.Client().RegionId))...)
}

func (r *securityGroupSnapshotPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state securityGroupSnapshotPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := describeSecurityGroupSnapshotPolicy(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeSecurityGroupSnapshotPolicies", securityGroupSnapshotPolicyTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeSecurityGroupSnapshotPolicies", securityGroupSnapshotPolicyTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, securityGroupSnapshotPolicyState(object, r.Client().RegionId))...)
}

func (r *securityGroupSnapshotPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state securityGroupSnapshotPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"RegionId":                      r.Client().RegionId,
		"SecurityGroupSnapshotPolicyId": state.Id.ValueString(),
	}

	// security_group_snapshot_policy_name and status are Optional+Computed
	// without UseStateForUnknown: when the user leaves one unset its plan value
	// is unknown, and ValueString() would return "" — sending an empty string
	// that clears the field. Guard with IsUnknown so an unset field is left
	// untouched (the subsequent Read back-fills it); an explicitly set value
	// still drives the modify.
	changed := false
	if !plan.SecurityGroupSnapshotPolicyName.IsUnknown() && !plan.SecurityGroupSnapshotPolicyName.Equal(state.SecurityGroupSnapshotPolicyName) {
		request["SecurityGroupSnapshotPolicyName"] = plan.SecurityGroupSnapshotPolicyName.ValueString()
		changed = true
	}
	if !plan.Status.IsUnknown() && !plan.Status.Equal(state.Status) {
		request["Status"] = plan.Status.ValueString()
		changed = true
	}

	if changed {
		if _, err := r.Client().RpcPost("Ecs", "2016-03-14", "ModifySecurityGroupSnapshotPolicy", nil, request, true); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling ModifySecurityGroupSnapshotPolicy", securityGroupSnapshotPolicyTypeName),
				errs.WrapErrorf(err, "Updating %s: calling ModifySecurityGroupSnapshotPolicy", securityGroupSnapshotPolicyTypeName).Error(),
			)
			return
		}
	}

	object, err := describeSecurityGroupSnapshotPolicy(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeSecurityGroupSnapshotPolicies", securityGroupSnapshotPolicyTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeSecurityGroupSnapshotPolicies", securityGroupSnapshotPolicyTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, securityGroupSnapshotPolicyState(object, r.Client().RegionId))...)
}

func (r *securityGroupSnapshotPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state securityGroupSnapshotPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.Client().RpcPost("Ecs", "2016-03-14", "DeleteSecurityGroupSnapshotPolicy", nil, map[string]interface{}{
		"RegionId":                      r.Client().RegionId,
		"SecurityGroupSnapshotPolicyId": state.Id.ValueString(),
	}, true)
	if err != nil && !errs.NotFoundError(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteSecurityGroupSnapshotPolicy", securityGroupSnapshotPolicyTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteSecurityGroupSnapshotPolicy", securityGroupSnapshotPolicyTypeName).Error(),
		)
		return
	}
}

func (r *securityGroupSnapshotPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// describeSecurityGroupSnapshotPolicy queries DescribeSecurityGroupSnapshotPolicies
// by id and returns the matching policy object, or nil when the policy no longer
// exists.
//
// A confirmed NotFound is distinguished from a malformed response: the API
// returning a NotFound error, or a well-formed but empty list, yields (nil, nil)
// so Read/Update can drop the resource. A missing root, an unrecognized type, an
// unknown inner shape, or a non-empty list that carries an element this helper
// cannot parse (a non-map item, or a map whose id is missing, nil or the wrong
// type) yields an error so the caller keeps the state instead of treating an
// uncertain response as a deletion (R4.4): only a proven absence — an empty
// list, or a list of well-formed policies none of which matches — clears the
// id.
func describeSecurityGroupSnapshotPolicy(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	ids, _ := json.Marshal([]string{id})
	request := map[string]interface{}{
		"RegionId":                       client.RegionId,
		"SecurityGroupSnapshotPolicyIds": string(ids),
	}
	response, err := client.RpcPost("Ecs", "2016-03-14", "DescribeSecurityGroupSnapshotPolicies", nil, request, true)
	if err != nil {
		if errs.NotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	items, ok := securityGroupSnapshotPolicyItems(response)
	if !ok {
		return nil, fmt.Errorf("DescribeSecurityGroupSnapshotPolicies returned an unrecognized response shape for %s: %v", id, response)
	}

	// An empty but well-formed list is a confirmed NotFound.
	if len(items) == 0 {
		return nil, nil
	}

	// Only a legitimate empty list (above) or a list of well-formed policies
	// none of which matches the requested id is a confirmed absence. A non-empty
	// list that carries an element this helper cannot parse — a non-map item, or
	// a map whose id is missing, nil or the wrong type — is not trustworthy: that
	// element may be the requested policy, so it is reported as an error and the
	// caller keeps the state instead of dropping the resource (R4.4). This stops a
	// malformed-but-non-empty Describe from masquerading as a NotFound.
	var malformed []interface{}
	for _, raw := range items {
		item, ok := raw.(map[string]interface{})
		if !ok {
			malformed = append(malformed, raw)
			continue
		}
		pid, parseable := policyIdOfOk(item)
		if !parseable {
			malformed = append(malformed, raw)
			continue
		}
		if pid == id {
			return item, nil
		}
	}
	if len(malformed) > 0 {
		return nil, fmt.Errorf("DescribeSecurityGroupSnapshotPolicies(%s): response carried %d malformed policy element(s); cannot confirm absence: %v", id, len(malformed), malformed)
	}
	// A well-formed, non-empty list of parseable policies, none of which matches:
	// the API returned real policies, none is the requested one, so it is absent.
	return nil, nil
}

// describeSecurityGroupSnapshotPolicyVisible retries a Describe across the
// eventual-consistency window in which a just-created policy is not yet
// visible. It returns as soon as describe yields a real error, a malformed
// response (error) or a visible object; only an empty-but-well-formed list
// (nil, nil) is retried, up to maxAttempts.
func describeSecurityGroupSnapshotPolicyVisible(describe func() (map[string]interface{}, error), maxAttempts int, interval time.Duration) (map[string]interface{}, error) {
	var object map[string]interface{}
	var err error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		object, err = describe()
		if err != nil || object != nil {
			return object, err
		}
		if attempt < maxAttempts-1 {
			time.Sleep(interval)
		}
	}
	return nil, nil
}

// securityGroupSnapshotPolicyItems extracts the policy array from a Describe
// response in either of the two POP list shapes: a top-level array, or an
// inner object's array. The second return is false when the root is absent or
// has a type this helper does not model, which the caller treats as a
// malformed response rather than a NotFound.
func securityGroupSnapshotPolicyItems(response map[string]interface{}) ([]interface{}, bool) {
	if sets, ok := response["SecurityGroupSnapshotPolicySets"].([]interface{}); ok {
		return sets, true
	}
	if inner, ok := response["SecurityGroupSnapshotPolicySets"].(map[string]interface{}); ok {
		if arr, ok := inner["SecurityGroupSnapshotPolicy"].([]interface{}); ok {
			return arr, true
		}
		// Inner object present but no recognizable array: malformed.
		return nil, false
	}
	// Root missing entirely: malformed (the API always returns the root for a
	// list call; its absence means the response is not what this code models).
	return nil, false
}

// securityGroupSnapshotPolicyIdOf reads a policy id without fmt.Sprintf
// coercion. It returns (id, true) only for a non-empty string or a JSON number;
// a missing, nil or wrong-type (map/slice/bool/...) id returns ("", false).
//
// RpcPost parses response bodies through tea-utils ReadAsJSON, which decodes
// with json.Decoder.UseNumber(), so a JSON number arrives as json.Number (a
// named string type), not float64. json.Number is normalized via String() —
// preserving the API's literal text — so a numeric id stays the same string
// across Create, Describe and state. The float64 branch is kept for callers
// that build a response map without going through RpcPost (for example, direct
// unit tests of this helper), so the contract is uniform whichever numeric
// representation a caller hands in.
//
// Create's extract, Describe's policyIdOfOk and the state writer share this one
// contract so a numeric id is consistent across Create->Read->state and a
// malformed id is rejected everywhere — not just in Describe — instead of the
// previous asymmetry where Create's default branch coerced a map/slice/bool
// into a fmt.Sprintf pseudo-id (D1), the state writer only accepted a string id,
// collapsing a numeric id to null and overwriting the id persisted by Create
// (D2), and neither branch recognized the json.Number RpcPost actually delivers
// (D3).
func securityGroupSnapshotPolicyIdOf(raw interface{}) (string, bool) {
	if raw == nil {
		return "", false
	}
	switch v := raw.(type) {
	case string:
		if v == "" {
			return "", false
		}
		return v, true
	case json.Number:
		s := v.String()
		if s == "" {
			return "", false
		}
		return s, true
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), true
	default:
		return "", false
	}
}

// extractSecurityGroupSnapshotPolicyId reads the Create response's id via the
// shared securityGroupSnapshotPolicyIdOf contract: a missing, nil or wrong-type
// (map/slice/bool/...) id is reported as absent (ok=false) so the caller fails
// the create rather than chasing a fmt.Sprintf pseudo-id such as "<nil>" or
// "map[x:sp-1]" (D1).
func extractSecurityGroupSnapshotPolicyId(response map[string]interface{}) (string, bool) {
	return securityGroupSnapshotPolicyIdOf(response["SecurityGroupSnapshotPolicyId"])
}

// policyIdOfOk reads a policy's id via the shared securityGroupSnapshotPolicyIdOf
// contract so Describe's matching uses the exact same string / float64
// normalization rule as Create's extract and the state writer. A missing, nil or
// wrong-type id returns ("", false) so the caller treats the element as
// malformed — it cannot be reliably matched — rather than coercing it into a
// pseudo-id that never matches and silently skipping it (which let a malformed
// element fall through to a NotFound and drop the resource's state).
func policyIdOfOk(item map[string]interface{}) (string, bool) {
	return securityGroupSnapshotPolicyIdOf(item["SecurityGroupSnapshotPolicyId"])
}

// securityGroupSnapshotPolicyState maps the policy object of a Describe response onto
// the resource model. Empty strings become null: optional fields the API omits read
// back as absent rather than "". The id field uses the shared
// securityGroupSnapshotPolicyIdOf contract (string or float64 normalized via
// FormatFloat) so a numeric id returned by Describe does not collapse to null and
// overwrite the id persisted by Create (D2); other fields stay string-only, since
// a numeric name/status would be a malformed response rather than a valid value.
func securityGroupSnapshotPolicyState(object map[string]interface{}, regionId string) securityGroupSnapshotPolicyResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	idStr := func() types.String {
		if id, ok := securityGroupSnapshotPolicyIdOf(object["SecurityGroupSnapshotPolicyId"]); ok {
			return types.StringValue(id)
		}
		return types.StringNull()
	}
	regionIdStr := str("RegionId")
	if regionIdStr.IsNull() && regionId != "" {
		regionIdStr = types.StringValue(regionId)
	}
	return securityGroupSnapshotPolicyResourceModel{
		Id:                              idStr(),
		SecurityGroupSnapshotPolicyId:   idStr(),
		SecurityGroupSnapshotPolicyName: str("SecurityGroupSnapshotPolicyName"),
		Status:                          str("Status"),
		RegionId:                        regionIdStr,
	}
}
