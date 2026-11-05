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
	accessPointPolicyForObjectProcessTypeName = "alicloud_oss_access_point_policy_for_object_process"
	accessPointPolicyForObjectProcessAction   = "/?accessPointPolicyForObjectProcess"

	accessPointPolicyForObjectProcessProduct    = "Oss"
	accessPointPolicyForObjectProcessAPIVersion = "2019-05-17"

	// accessPointPolicyForObjectProcessTimeout bounds the retry envelope of the
	// policy writes (Put and Delete); the read has its own one-minute budget.
	accessPointPolicyForObjectProcessTimeout = 5 * time.Minute
)

var (
	_ resource.Resource                = &accessPointPolicyForObjectProcessResource{}
	_ resource.ResourceWithConfigure   = &accessPointPolicyForObjectProcessResource{}
	_ resource.ResourceWithImportState = &accessPointPolicyForObjectProcessResource{}
)

func NewAccessPointPolicyForObjectProcessResource() resource.Resource {
	return &accessPointPolicyForObjectProcessResource{}
}

type accessPointPolicyForObjectProcessResource struct {
	fwadapt.ResourceBase
}

type accessPointPolicyForObjectProcessResourceModel struct {
	Id              types.String `tfsdk:"id"`
	AccessPointName types.String `tfsdk:"access_point_name"`
	Bucket          types.String `tfsdk:"bucket"`
	Policy          types.String `tfsdk:"policy"`
}

func (r *accessPointPolicyForObjectProcessResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = accessPointPolicyForObjectProcessTypeName
}

func (r *accessPointPolicyForObjectProcessResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *accessPointPolicyForObjectProcessResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan accessPointPolicyForObjectProcessResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucket := plan.Bucket.ValueString()
	accessPointName := plan.AccessPointName.ValueString()
	if err := r.putAccessPointPolicyForObjectProcess(bucket, accessPointName, plan.Policy.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling PutAccessPointPolicyForObjectProcess", accessPointPolicyForObjectProcessTypeName),
			errs.WrapErrorf(err, "Creating %s: calling PutAccessPointPolicyForObjectProcess", accessPointPolicyForObjectProcessTypeName).Error(),
		)
		return
	}

	id := fmt.Sprintf("%s:%s", bucket, accessPointName)
	// Persist the resource identity before the read-back: the policy exists
	// remotely from here on, so a failing GET must not leave Terraform without
	// a resource to taint and destroy.
	state := accessPointPolicyForObjectProcessResourceModel{
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

	object, err := r.describeAccessPointPolicyForObjectProcess(id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetAccessPointPolicyForObjectProcess", accessPointPolicyForObjectProcessTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetAccessPointPolicyForObjectProcess", accessPointPolicyForObjectProcessTypeName).Error(),
		)
		return
	}

	state = accessPointPolicyForObjectProcessState(object)
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

func (r *accessPointPolicyForObjectProcessResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state accessPointPolicyForObjectProcessResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	object, err := r.describeAccessPointPolicyForObjectProcess(id)
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetAccessPointPolicyForObjectProcess", accessPointPolicyForObjectProcessTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetAccessPointPolicyForObjectProcess", accessPointPolicyForObjectProcessTypeName).Error(),
		)
		return
	}

	model, err := accessPointPolicyForObjectProcessStateFromId(object, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing resource id", accessPointPolicyForObjectProcessTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *accessPointPolicyForObjectProcessResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state accessPointPolicyForObjectProcessResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Policy.Equal(state.Policy) {
		if err := r.putAccessPointPolicyForObjectProcess(state.Bucket.ValueString(), state.AccessPointName.ValueString(), plan.Policy.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling PutAccessPointPolicyForObjectProcess", accessPointPolicyForObjectProcessTypeName),
				errs.WrapErrorf(err, "Updating %s: calling PutAccessPointPolicyForObjectProcess", accessPointPolicyForObjectProcessTypeName).Error(),
			)
			return
		}
	}

	object, err := r.describeAccessPointPolicyForObjectProcess(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetAccessPointPolicyForObjectProcess", accessPointPolicyForObjectProcessTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetAccessPointPolicyForObjectProcess", accessPointPolicyForObjectProcessTypeName).Error(),
		)
		return
	}

	model := accessPointPolicyForObjectProcessState(object)
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

func (r *accessPointPolicyForObjectProcessResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state accessPointPolicyForObjectProcessResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.deleteAccessPointPolicyForObjectProcess(state.Bucket.ValueString(), state.AccessPointName.ValueString(), state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"NoSuchAccessPointPolicyForObjectProcess", "NoSuchAccessPointForObjectProcess", "NoSuchAccessPoint"}) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteAccessPointPolicyForObjectProcess", accessPointPolicyForObjectProcessTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteAccessPointPolicyForObjectProcess", accessPointPolicyForObjectProcessTypeName).Error(),
		)
		return
	}
}

func (r *accessPointPolicyForObjectProcessResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// describeAccessPointPolicyForObjectProcess fetches the policy document attached
// to an Object FC access point. The id is the resource id: <bucket>:<access_point_name>.
func (r *accessPointPolicyForObjectProcessResource) describeAccessPointPolicyForObjectProcess(id string) (map[string]interface{}, error) {
	bucket, accessPointName, err := parseAccessPointPolicyForObjectProcessId(id)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	query := accessPointPolicyForObjectProcessQuery(accessPointName)
	hostMap := accessPointPolicyForObjectProcessHostMap(bucket)
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(accessPointPolicyForObjectProcessProduct, helper.XmlJsonParam("GET", accessPointPolicyForObjectProcessAPIVersion, "GetAccessPointPolicyForObjectProcess", accessPointPolicyForObjectProcessAction), query, nil, nil, hostMap, true)
		if requestErr != nil {
			if helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(accessPointPolicyForObjectProcessAction, response, id)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchAccessPointPolicyForObjectProcess", "NoSuchAccessPointForObjectProcess", "NoSuchAccessPoint"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("AccessPointPolicyForObjectProcess", id), errs.NotFoundMsg, response)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, accessPointPolicyForObjectProcessAction, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return nil, errs.WrapErrorf(errs.NotFoundErr("AccessPointPolicyForObjectProcess", id), errs.NotFoundMsg, response)
	}

	// OSS Object FC access-point policy responses are returned as the policy
	// document itself. Some responses are still wrapped in a *Result element;
	// extract it defensively and fall back to the raw response when no wrapper
	// is present.
	if v, jsonErr := jsonpath.Get("$.GetAccessPointPolicyForObjectProcessResult", response); jsonErr == nil && v != nil {
		if vv, ok := v.(map[string]interface{}); ok {
			response = vv
		}
	}
	return response, nil
}

func parseAccessPointPolicyForObjectProcessId(id string) (bucket, accessPointName string, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid Resource Id %s. Expected bucket:access_point_name, got %s", id, id)
	}
	return parts[0], parts[1], nil
}

func accessPointPolicyForObjectProcessQuery(accessPointName string) map[string]*string {
	return map[string]*string{
		"x-oss-access-point-for-object-process-name": helper.StringPointer(accessPointName),
	}
}

func accessPointPolicyForObjectProcessHostMap(bucket string) map[string]*string {
	return map[string]*string{
		"bucket": helper.StringPointer(bucket),
	}
}

func (r *accessPointPolicyForObjectProcessResource) putAccessPointPolicyForObjectProcess(bucket, accessPointName, policy string) error {
	query := accessPointPolicyForObjectProcessQuery(accessPointName)
	hostMap := accessPointPolicyForObjectProcessHostMap(bucket)
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(accessPointPolicyForObjectProcessTimeout, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(accessPointPolicyForObjectProcessProduct, helper.JsonXmlParam("PUT", accessPointPolicyForObjectProcessAPIVersion, "PutAccessPointPolicyForObjectProcess", accessPointPolicyForObjectProcessAction), query, policy, nil, hostMap, false)
		if requestErr != nil {
			if helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(accessPointPolicyForObjectProcessAction, response, policy)
	return err
}

func (r *accessPointPolicyForObjectProcessResource) deleteAccessPointPolicyForObjectProcess(bucket, accessPointName, id string) error {
	query := accessPointPolicyForObjectProcessQuery(accessPointName)
	hostMap := accessPointPolicyForObjectProcessHostMap(bucket)
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(accessPointPolicyForObjectProcessTimeout, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(accessPointPolicyForObjectProcessProduct, helper.JsonXmlParam("DELETE", accessPointPolicyForObjectProcessAPIVersion, "DeleteAccessPointPolicyForObjectProcess", accessPointPolicyForObjectProcessAction), query, nil, nil, hostMap, false)
		if requestErr != nil {
			// A policy operation can race the Object FC access point creation; retry
			// that conflict the same way the flat access point resource does.
			if errs.IsExpectedErrors(requestErr, []string{"AccessPointCreatingConflict", "AccessPointForObjectProcessCreatingConflict"}) || helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(accessPointPolicyForObjectProcessAction, response, id)
	return err
}

// accessPointPolicyForObjectProcessState maps a GetAccessPointPolicyForObjectProcess
// document onto the resource model. The policy is stored as the canonical JSON
// encoding of the response so semantically equal documents read back identically.
func accessPointPolicyForObjectProcessState(object map[string]interface{}) accessPointPolicyForObjectProcessResourceModel {
	return accessPointPolicyForObjectProcessResourceModel{
		Policy: types.StringValue(convert.ConvertMapToJsonStringIgnoreError(object)),
	}
}

// accessPointPolicyForObjectProcessStateFromId builds the model from a fetched policy
// plus the owning ids encoded in the resource id, which the policy document itself
// does not carry.
func accessPointPolicyForObjectProcessStateFromId(object map[string]interface{}, id string) (accessPointPolicyForObjectProcessResourceModel, error) {
	bucket, accessPointName, err := parseAccessPointPolicyForObjectProcessId(id)
	if err != nil {
		return accessPointPolicyForObjectProcessResourceModel{}, err
	}
	model := accessPointPolicyForObjectProcessState(object)
	model.Id = types.StringValue(id)
	model.Bucket = types.StringValue(bucket)
	model.AccessPointName = types.StringValue(accessPointName)
	return model, nil
}
