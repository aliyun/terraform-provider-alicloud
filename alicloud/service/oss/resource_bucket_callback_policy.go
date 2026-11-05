package oss

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	bucketCallbackPolicyTypeName = "alicloud_oss_bucket_callback_policy"
	bucketCallbackPolicyAction   = "/?policy&comp=callback"

	bucketCallbackPolicyProduct    = "Oss"
	bucketCallbackPolicyAPIVersion = "2019-05-17"

	// bucketCallbackPolicyTimeout bounds the retry envelope of the policy writes
	// (Put and Delete); the read has its own one-minute budget, matching the
	// other OSS bucket sub-resources.
	bucketCallbackPolicyTimeout = 5 * time.Minute

	// notFoundCodes are returned by OSS when the bucket or its callback policy
	// is absent; both are treated as "the resource no longer exists".
	notFoundCodes = "NoSuchBucket,BucketCallbackPolicyNotExist"
)

var notFoundCodeList = strings.Split(notFoundCodes, ",")

var (
	_ resource.Resource                = &bucketCallbackPolicyResource{}
	_ resource.ResourceWithConfigure   = &bucketCallbackPolicyResource{}
	_ resource.ResourceWithImportState = &bucketCallbackPolicyResource{}
)

func NewBucketCallbackPolicyResource() resource.Resource {
	return &bucketCallbackPolicyResource{}
}

type bucketCallbackPolicyResource struct {
	fwadapt.ResourceBase
}

// policyItemModel is the per-item shape of the policy_item list block.
type policyItemModel struct {
	Callback    types.String `tfsdk:"callback"`
	PolicyName  types.String `tfsdk:"policy_name"`
	CallbackVar types.String `tfsdk:"callback_var"`
}

type bucketCallbackPolicyResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Bucket     types.String `tfsdk:"bucket"`
	PolicyItem types.List   `tfsdk:"policy_item"`
}

func (r *bucketCallbackPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = bucketCallbackPolicyTypeName
}

func (r *bucketCallbackPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bucket": schema.StringAttribute{
				Required:    true,
				Description: "The name of the bucket.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy_item": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The container that holds the callback policy configuration items.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"callback": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Callback address. It is a Base64-encoded JSON string.",
						},
						"policy_name": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Name of the callback policy configuration item.",
						},
						"callback_var": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Callback parameters. It is a Base64-encoded JSON string.",
						},
					},
				},
			},
		},
	}
}

func (r *bucketCallbackPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan bucketCallbackPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucket := plan.Bucket.ValueString()
	var items []policyItemModel
	resp.Diagnostics.Append(plan.PolicyItem.ElementsAs(ctx, &items, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.putBucketCallbackPolicy(bucket, items); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling PutBucketCallbackPolicy", bucketCallbackPolicyTypeName),
			errs.WrapErrorf(err, "Creating %s: calling PutBucketCallbackPolicy", bucketCallbackPolicyTypeName).Error(),
		)
		return
	}

	// The resource id is the owning bucket name; the callback policy is a
	// singleton sub-resource of the bucket.
	id := bucket
	state := bucketCallbackPolicyResourceModel{
		Id:     types.StringValue(id),
		Bucket: plan.Bucket,
	}
	// Persist the configured items verbatim before the read-back: a failing GET
	// must not leave Terraform without a resource to taint and destroy, and the
	// API performs no server-side transformation of these fields.
	state.PolicyItem = plan.PolicyItem
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := r.describeBucketCallbackPolicy(id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBucketCallbackPolicy", bucketCallbackPolicyTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetBucketCallbackPolicy (%s)", bucketCallbackPolicyTypeName, id).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: after create", bucketCallbackPolicyTypeName),
			fmt.Sprintf("The resource %s is not visible via GetBucketCallbackPolicy right after creation", id),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, bucketCallbackPolicyState(ctx, object, &plan))...)
}

func (r *bucketCallbackPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state bucketCallbackPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	object, err := r.describeBucketCallbackPolicy(id)
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBucketCallbackPolicy", bucketCallbackPolicyTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetBucketCallbackPolicy", bucketCallbackPolicyTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, bucketCallbackPolicyState(ctx, object, &state))...)
}

func (r *bucketCallbackPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state bucketCallbackPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// bucket is ForceNew, so the owning bucket never changes here. The policy
	// is a singleton on the bucket; any change to policy_item is a full Put.
	if !plan.PolicyItem.Equal(state.PolicyItem) {
		var items []policyItemModel
		resp.Diagnostics.Append(plan.PolicyItem.ElementsAs(ctx, &items, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if err := r.putBucketCallbackPolicy(state.Bucket.ValueString(), items); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling PutBucketCallbackPolicy", bucketCallbackPolicyTypeName),
				errs.WrapErrorf(err, "Updating %s: calling PutBucketCallbackPolicy", bucketCallbackPolicyTypeName).Error(),
			)
			return
		}
	}

	object, err := r.describeBucketCallbackPolicy(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBucketCallbackPolicy", bucketCallbackPolicyTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetBucketCallbackPolicy", bucketCallbackPolicyTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, bucketCallbackPolicyState(ctx, object, &plan))...)
}

func (r *bucketCallbackPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state bucketCallbackPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.deleteBucketCallbackPolicy(state.Bucket.ValueString())
	if err != nil {
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, notFoundCodeList) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteBucketCallbackPolicy", bucketCallbackPolicyTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteBucketCallbackPolicy", bucketCallbackPolicyTypeName).Error(),
		)
		return
	}
}

func (r *bucketCallbackPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// describeBucketCallbackPolicy fetches the callback policy attached to a
// bucket. The id is the resource id: the bucket name. The response wraps the
// policy in a BucketCallbackPolicy root element, which is extracted here.
func (r *bucketCallbackPolicyResource) describeBucketCallbackPolicy(id string) (map[string]interface{}, error) {
	query := map[string]*string{}
	hostMap := map[string]*string{
		"bucket": helper.StringPointer(id),
	}
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(1*time.Minute, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(bucketCallbackPolicyProduct, helper.XmlParam("GET", bucketCallbackPolicyAPIVersion, "GetBucketCallbackPolicy", bucketCallbackPolicyAction), query, nil, nil, hostMap, false)
		if requestErr != nil {
			if helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(bucketCallbackPolicyAction, response, id)
	if err != nil {
		if errs.IsExpectedErrors(err, notFoundCodeList) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("BucketCallbackPolicy", id), errs.NotFoundMsg, response)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, bucketCallbackPolicyAction, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return nil, errs.WrapErrorf(errs.NotFoundErr("BucketCallbackPolicy", id), errs.NotFoundMsg, response)
	}

	// OSS wraps the policy in a BucketCallbackPolicy root element; extract it
	// and fall back to the raw response when no wrapper is present.
	if v, jsonErr := jsonpath.Get("$.BucketCallbackPolicy", response); jsonErr == nil && v != nil {
		if vv, ok := v.(map[string]interface{}); ok {
			response = vv
		}
	}
	return response, nil
}

func (r *bucketCallbackPolicyResource) putBucketCallbackPolicy(bucket string, items []policyItemModel) error {
	query := map[string]*string{}
	hostMap := map[string]*string{
		"bucket": helper.StringPointer(bucket),
	}
	body := map[string]interface{}{
		"BucketCallbackPolicy": map[string]interface{}{
			"PolicyItem": policyItemToWire(items),
		},
	}
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(bucketCallbackPolicyTimeout, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(bucketCallbackPolicyProduct, helper.XmlParam("PUT", bucketCallbackPolicyAPIVersion, "PutBucketCallbackPolicy", bucketCallbackPolicyAction), query, body, nil, hostMap, false)
		if requestErr != nil {
			if helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(bucketCallbackPolicyAction, response, body)
	return err
}

func (r *bucketCallbackPolicyResource) deleteBucketCallbackPolicy(bucket string) error {
	query := map[string]*string{}
	hostMap := map[string]*string{
		"bucket": helper.StringPointer(bucket),
	}
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(bucketCallbackPolicyTimeout, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(bucketCallbackPolicyProduct, helper.XmlParam("DELETE", bucketCallbackPolicyAPIVersion, "DeleteBucketCallbackPolicy", bucketCallbackPolicyAction), query, nil, nil, hostMap, false)
		if requestErr != nil {
			if helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(bucketCallbackPolicyAction, response, bucket)
	return err
}

// bucketCallbackPolicyState maps a GetBucketCallbackPolicy document onto the
// resource model. The bucket is taken from the response (falling back to the
// prior id) and policy_item is rebuilt from the PolicyItem list, preserving the
// configured values for fields the response omits.
func bucketCallbackPolicyState(ctx context.Context, object map[string]interface{}, prior *bucketCallbackPolicyResourceModel) bucketCallbackPolicyResourceModel {
	bucket := stringFromResponse(object, "Bucket")
	if bucket.IsNull() && prior != nil {
		bucket = prior.Bucket
	}

	items := policyItemListFromResponse(object["PolicyItem"])
	if prior != nil {
		items = preservePlannedPolicyItems(ctx, items, prior.PolicyItem)
	}

	return bucketCallbackPolicyResourceModel{
		Id:         types.StringValue(bucket.ValueString()),
		Bucket:     bucket,
		PolicyItem: items,
	}
}

// policyItemToWire builds the wire-shaped (PascalCase) list of maps sent in the
// PutBucketCallbackPolicy body, skipping unset fields so the server receives
// only what the user configured.
func policyItemToWire(items []policyItemModel) []interface{} {
	out := make([]interface{}, 0, len(items))
	for _, it := range items {
		m := map[string]interface{}{}
		if !it.Callback.IsNull() && !it.Callback.IsUnknown() {
			m["Callback"] = it.Callback.ValueString()
		}
		if !it.PolicyName.IsNull() && !it.PolicyName.IsUnknown() {
			m["PolicyName"] = it.PolicyName.ValueString()
		}
		if !it.CallbackVar.IsNull() && !it.CallbackVar.IsUnknown() {
			m["CallbackVar"] = it.CallbackVar.ValueString()
		}
		out = append(out, m)
	}
	return out
}

// policyItemListFromResponse rebuilds the policy_item list from the Get response.
// OSS may return a single item as a bare map rather than a one-element list, so
// both shapes are accepted; the response order is preserved.
func policyItemListFromResponse(raw interface{}) types.List {
	elemType := types.ObjectType{AttrTypes: policyItemAttrTypes()}
	if raw == nil {
		return types.ListNull(elemType)
	}

	switch v := raw.(type) {
	case []interface{}:
		items := make([]policyItemModel, 0, len(v))
		for _, e := range v {
			items = append(items, policyItemModelFromResponse(e))
		}
		return listFromModels(elemType, items)
	case map[string]interface{}:
		// OSS decodes a single PolicyItem as a bare object; normalize to a list.
		return listFromModels(elemType, []policyItemModel{policyItemModelFromResponse(v)})
	default:
		return types.ListNull(elemType)
	}
}

func policyItemModelFromResponse(e interface{}) policyItemModel {
	m, _ := e.(map[string]interface{})
	return policyItemModel{
		Callback:    stringFromResponse(m, "Callback"),
		PolicyName:  stringFromResponse(m, "PolicyName"),
		CallbackVar: stringFromResponse(m, "CallbackVar"),
	}
}

// preservePlannedPolicyItems keeps the configured value for any field the
// response left absent while taking the server-returned values where present,
// so a canonical re-encoding does not surface as a spurious diff.
func preservePlannedPolicyItems(ctx context.Context, responseList types.List, planned types.List) types.List {
	if planned.IsNull() || planned.IsUnknown() {
		return responseList
	}
	var plannedItems []policyItemModel
	if diags := planned.ElementsAs(ctx, &plannedItems, false); diags.HasError() {
		return responseList
	}
	var responseItems []policyItemModel
	if diags := responseList.ElementsAs(ctx, &responseItems, false); diags.HasError() {
		return responseList
	}

	// Match by policy_name when both sides carry it; otherwise fall back to
	// index, so a reordering by the server does not drop configured values.
	priorByName := make(map[string]policyItemModel, len(plannedItems))
	for _, p := range plannedItems {
		if !p.PolicyName.IsNull() && !p.PolicyName.IsUnknown() {
			priorByName[p.PolicyName.ValueString()] = p
		}
	}
	for i := range responseItems {
		// Clamp to planned length when the server returned more items than
		// configured; an unconfigured extra item is still recorded as-is.
		var prior policyItemModel
		if i < len(plannedItems) {
			prior = plannedItems[i]
		}
		if name := responseItems[i].PolicyName; !name.IsNull() && !name.IsUnknown() {
			if p, ok := priorByName[name.ValueString()]; ok {
				prior = p
			}
		}
		responseItems[i].Callback = firstKnown(responseItems[i].Callback, prior.Callback)
		responseItems[i].PolicyName = firstKnown(responseItems[i].PolicyName, prior.PolicyName)
		responseItems[i].CallbackVar = firstKnown(responseItems[i].CallbackVar, prior.CallbackVar)
	}
	return listFromModels(types.ObjectType{AttrTypes: policyItemAttrTypes()}, responseItems)
}

func listFromModels(elemType attr.Type, items []policyItemModel) types.List {
	values := make([]attr.Value, 0, len(items))
	for _, it := range items {
		obj, diags := types.ObjectValue(policyItemAttrTypes(), map[string]attr.Value{
			"callback":     it.Callback,
			"policy_name":  it.PolicyName,
			"callback_var": it.CallbackVar,
		})
		if diags.HasError() {
			continue
		}
		values = append(values, obj)
	}
	list, diags := types.ListValue(elemType, values)
	if diags.HasError() {
		return types.ListNull(elemType)
	}
	return list
}

// stringFromResponse returns a present string (even when empty) as a known
// value, and an absent field as null. This mirrors the v1 d.Set behavior where
// an API-returned "" is a set value, not absence.
func stringFromResponse(object map[string]interface{}, key string) types.String {
	if object == nil {
		return types.StringNull()
	}
	v, ok := object[key]
	if !ok || v == nil {
		return types.StringNull()
	}
	if s, ok := v.(string); ok {
		return types.StringValue(s)
	}
	return types.StringNull()
}

// firstKnown returns value when it is known (non-null, non-unknown), otherwise
// fallback. Used to keep the planned value where the response omitted a field.
func firstKnown(value, fallback types.String) types.String {
	if !value.IsNull() && !value.IsUnknown() {
		return value
	}
	return fallback
}

func policyItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"callback":     types.StringType,
		"policy_name":  types.StringType,
		"callback_var": types.StringType,
	}
}
