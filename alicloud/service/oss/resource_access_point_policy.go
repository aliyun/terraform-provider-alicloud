package oss

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	accessPointPolicyTypeName = "alicloud_oss_access_point_policy"
	accessPointPolicyAction   = "/?accessPointPolicy"

	accessPointPolicyProduct    = "Oss"
	accessPointPolicyAPIVersion = "2019-05-17"

	// accessPointPolicyTimeout bounds the retry envelope of the policy writes
	// (Put and Delete); the read has its own one-minute budget.
	accessPointPolicyTimeout = 5 * time.Minute
)

var (
	_ resource.Resource                = &accessPointPolicyResource{}
	_ resource.ResourceWithConfigure   = &accessPointPolicyResource{}
	_ resource.ResourceWithImportState = &accessPointPolicyResource{}
)

func NewAccessPointPolicyResource() resource.Resource {
	return &accessPointPolicyResource{}
}

type accessPointPolicyResource struct {
	fwadapt.ResourceBase
}

type accessPointPolicyResourceModel struct {
	Id              types.String `tfsdk:"id"`
	AccessPointName types.String `tfsdk:"access_point_name"`
	Bucket          types.String `tfsdk:"bucket"`
	Policy          types.String `tfsdk:"policy"`
}

func (r *accessPointPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = accessPointPolicyTypeName
}

func (r *accessPointPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_point_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"bucket": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					jsonTemplateEquivalentModifier{},
				},
			},
		},
	}
}

func (r *accessPointPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan accessPointPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucket := plan.Bucket.ValueString()
	accessPointName := plan.AccessPointName.ValueString()
	if err := r.putAccessPointPolicy(bucket, accessPointName, plan.Policy.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling PutAccessPointPolicy", accessPointPolicyTypeName),
			errs.WrapErrorf(err, "Creating %s: calling PutAccessPointPolicy", accessPointPolicyTypeName).Error(),
		)
		return
	}

	id := fmt.Sprintf("%s:%s", bucket, accessPointName)
	// Persist the resource identity before the read-back: the policy exists
	// remotely from here on, so a failing GET must not leave Terraform without
	// a resource to taint and destroy.
	state := accessPointPolicyResourceModel{
		Id:              types.StringValue(id),
		AccessPointName: plan.AccessPointName,
		Bucket:          plan.Bucket,
	}
	if !plan.Policy.IsUnknown() {
		state.Policy = plan.Policy
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := r.describeAccessPointPolicy(id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetAccessPointPolicy", accessPointPolicyTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetAccessPointPolicy", accessPointPolicyTypeName).Error(),
		)
		return
	}

	state = accessPointPolicyState(object)
	state.Id = types.StringValue(id)
	state.Bucket = plan.Bucket
	state.AccessPointName = plan.AccessPointName
	if !plan.Policy.IsUnknown() {
		// Write the configured document back verbatim: the API returns a
		// canonical re-encoding, and any difference from the known planned
		// value fails Terraform's apply consistency check.
		state.Policy = plan.Policy
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *accessPointPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state accessPointPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	object, err := r.describeAccessPointPolicy(id)
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetAccessPointPolicy", accessPointPolicyTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetAccessPointPolicy", accessPointPolicyTypeName).Error(),
		)
		return
	}

	model, err := accessPointPolicyStateFromId(object, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing resource id", accessPointPolicyTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *accessPointPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state accessPointPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Policy.Equal(state.Policy) {
		if err := r.putAccessPointPolicy(state.Bucket.ValueString(), state.AccessPointName.ValueString(), plan.Policy.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling PutAccessPointPolicy", accessPointPolicyTypeName),
				errs.WrapErrorf(err, "Updating %s: calling PutAccessPointPolicy", accessPointPolicyTypeName).Error(),
			)
			return
		}
	}

	object, err := r.describeAccessPointPolicy(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetAccessPointPolicy", accessPointPolicyTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetAccessPointPolicy", accessPointPolicyTypeName).Error(),
		)
		return
	}

	model := accessPointPolicyState(object)
	model.Id = state.Id
	model.Bucket = plan.Bucket
	model.AccessPointName = plan.AccessPointName
	if !plan.Policy.IsUnknown() {
		// See Create: keep the configured document so the applied state matches
		// the known planned value.
		model.Policy = plan.Policy
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *accessPointPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state accessPointPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.deleteAccessPointPolicy(state.Bucket.ValueString(), state.AccessPointName.ValueString(), state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"NoSuchAccessPointPolicy", "NoSuchAccessPoint"}) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteAccessPointPolicy", accessPointPolicyTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteAccessPointPolicy", accessPointPolicyTypeName).Error(),
		)
		return
	}
}

func (r *accessPointPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// describeAccessPointPolicy fetches the policy document attached to an access
// point. The id is the resource id: <bucket>:<access_point_name>.
func (r *accessPointPolicyResource) describeAccessPointPolicy(id string) (map[string]interface{}, error) {
	bucket, accessPointName, err := parseAccessPointPolicyId(id)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	query := accessPointPolicyQuery(accessPointName)
	hostMap := accessPointPolicyHostMap(bucket)
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(accessPointPolicyProduct, helper.XmlJsonParam("GET", accessPointPolicyAPIVersion, "GetAccessPointPolicy", accessPointPolicyAction), query, nil, nil, hostMap, true)
		if requestErr != nil {
			if helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(accessPointPolicyAction, response, id)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchAccessPointPolicy", "NoSuchAccessPoint"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("AccessPointPolicy", id), errs.NotFoundMsg, response)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, accessPointPolicyAction, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return nil, errs.WrapErrorf(errs.NotFoundErr("AccessPointPolicy", id), errs.NotFoundMsg, response)
	}

	// OSS access-point responses are wrapped in a *Result element; extract it
	// defensively and fall back to the raw response when no wrapper is present.
	if v, jsonErr := jsonpath.Get("$.GetAccessPointPolicyResult", response); jsonErr == nil && v != nil {
		if vv, ok := v.(map[string]interface{}); ok {
			response = vv
		}
	}
	return response, nil
}

func parseAccessPointPolicyId(id string) (bucket, accessPointName string, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid Resource Id %s. Expected bucket:access_point_name, got %s", id, id)
	}
	return parts[0], parts[1], nil
}

func accessPointPolicyQuery(accessPointName string) map[string]*string {
	return map[string]*string{
		"x-oss-access-point-name": helper.StringPointer(accessPointName),
	}
}

func accessPointPolicyHostMap(bucket string) map[string]*string {
	return map[string]*string{
		"bucket": helper.StringPointer(bucket),
	}
}

func (r *accessPointPolicyResource) putAccessPointPolicy(bucket, accessPointName, policy string) error {
	query := accessPointPolicyQuery(accessPointName)
	hostMap := accessPointPolicyHostMap(bucket)
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(accessPointPolicyTimeout, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(accessPointPolicyProduct, helper.JsonXmlParam("PUT", accessPointPolicyAPIVersion, "PutAccessPointPolicy", accessPointPolicyAction), query, policy, nil, hostMap, false)
		if requestErr != nil {
			if helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(accessPointPolicyAction, response, policy)
	return err
}

func (r *accessPointPolicyResource) deleteAccessPointPolicy(bucket, accessPointName, id string) error {
	query := accessPointPolicyQuery(accessPointName)
	hostMap := accessPointPolicyHostMap(bucket)
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(accessPointPolicyTimeout, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(accessPointPolicyProduct, helper.JsonXmlParam("DELETE", accessPointPolicyAPIVersion, "DeleteAccessPointPolicy", accessPointPolicyAction), query, nil, nil, hostMap, false)
		if requestErr != nil {
			// A policy operation can race the access point creation, as with the
			// flat alicloud_oss_access_point resource; retry that conflict.
			if errs.IsExpectedErrors(requestErr, []string{"AccessPointCreatingConflict"}) || helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(accessPointPolicyAction, response, id)
	return err
}

// accessPointPolicyState maps a GetAccessPointPolicy document onto the
// resource model. The policy is stored as the canonical JSON encoding of the
// response so semantically equal documents read back identically.
func accessPointPolicyState(object map[string]interface{}) accessPointPolicyResourceModel {
	return accessPointPolicyResourceModel{
		Policy: types.StringValue(convert.ConvertMapToJsonStringIgnoreError(object)),
	}
}

// accessPointPolicyStateFromId builds the model from a fetched policy plus the
// owning ids encoded in the resource id, which the policy document itself does
// not carry.
func accessPointPolicyStateFromId(object map[string]interface{}, id string) (accessPointPolicyResourceModel, error) {
	bucket, accessPointName, err := parseAccessPointPolicyId(id)
	if err != nil {
		return accessPointPolicyResourceModel{}, err
	}
	model := accessPointPolicyState(object)
	model.Id = types.StringValue(id)
	model.Bucket = types.StringValue(bucket)
	model.AccessPointName = types.StringValue(accessPointName)
	return model, nil
}

// jsonTemplateEquivalentModifier keeps the prior state value in the plan when
// the planned policy differs only in JSON formatting, so an equivalent
// document is not reported as a change.
type jsonTemplateEquivalentModifier struct{}

func (jsonTemplateEquivalentModifier) Description(context.Context) string {
	return "Keeps the current value when the planned policy is an equivalent JSON document."
}

func (m jsonTemplateEquivalentModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (jsonTemplateEquivalentModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsUnknown() || req.PlanValue.IsNull() || req.StateValue.IsUnknown() || req.StateValue.IsNull() {
		return
	}
	if equal, err := helper.CompareJsonTemplateAreEquivalent(req.StateValue.ValueString(), req.PlanValue.ValueString()); err == nil && equal {
		resp.PlanValue = req.StateValue
	}
}
