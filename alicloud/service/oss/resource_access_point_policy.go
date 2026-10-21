package oss

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
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

// accessPointPolicyClient is the OSS client primitive the resource calls; it
// makes the wire requests replaceable in unit tests.
type accessPointPolicyClient interface {
	Do(apiProductCode string, apiParams *openapi.Params, query map[string]*string, body interface{}, headers map[string]*string, hostMap map[string]*string, autoRetry bool) (map[string]interface{}, error)
}

type accessPointPolicyResource struct {
	fwadapt.ResourceBase

	client accessPointPolicyClient
}

func (r *accessPointPolicyResource) policyClient() accessPointPolicyClient {
	if r.client != nil {
		return r.client
	}
	return r.Client()
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

// Configure takes the provider client through the standard Framework hook. The
// Do-shaped seam is accepted here as well so lifecycle tests can drive the
// resource through its public factory with a wire-level fake.
func (r *accessPointPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client, ok := req.ProviderData.(*connectivity.AliyunClient); ok {
		r.ResourceBase.Configure(ctx, resource.ConfigureRequest{ProviderData: client}, resp)
		return
	}
	if api, ok := req.ProviderData.(accessPointPolicyClient); ok {
		r.client = api
		return
	}
	r.ResourceBase.Configure(ctx, req, resp)
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
	client := r.policyClient()
	if err := putAccessPointPolicy(client, bucket, accessPointName, plan.Policy.ValueString()); err != nil {
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

	object, err := describeAccessPointPolicy(client, id)
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
	object, err := describeAccessPointPolicy(r.policyClient(), id)
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

	client := r.policyClient()
	if !plan.Policy.Equal(state.Policy) {
		if err := putAccessPointPolicy(client, state.Bucket.ValueString(), state.AccessPointName.ValueString(), plan.Policy.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling PutAccessPointPolicy", accessPointPolicyTypeName),
				errs.WrapErrorf(err, "Updating %s: calling PutAccessPointPolicy", accessPointPolicyTypeName).Error(),
			)
			return
		}
	}

	object, err := describeAccessPointPolicy(client, state.Id.ValueString())
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

	err := deleteAccessPointPolicy(r.policyClient(), state.Bucket.ValueString(), state.AccessPointName.ValueString(), state.Id.ValueString())
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

// DescribeAccessPointPolicy fetches the policy document attached to an access
// point. The id is the resource id: <bucket>:<access_point_name>.
func DescribeAccessPointPolicy(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	return describeAccessPointPolicy(client, id)
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

func putAccessPointPolicy(client accessPointPolicyClient, bucket, accessPointName, policy string) error {
	query := accessPointPolicyQuery(accessPointName)
	hostMap := accessPointPolicyHostMap(bucket)
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(accessPointPolicyTimeout, func() *retry.RetryError {
		var requestErr error
		response, requestErr = client.Do("Oss", helper.JsonXmlParam("PUT", "2019-05-17", "PutAccessPointPolicy", accessPointPolicyAction), query, policy, nil, hostMap, false)
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

func describeAccessPointPolicy(client accessPointPolicyClient, id string) (map[string]interface{}, error) {
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
		response, requestErr = client.Do("Oss", helper.XmlJsonParam("GET", "2019-05-17", "GetAccessPointPolicy", accessPointPolicyAction), query, nil, nil, hostMap, true)
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

func deleteAccessPointPolicy(client accessPointPolicyClient, bucket, accessPointName, id string) error {
	query := accessPointPolicyQuery(accessPointName)
	hostMap := accessPointPolicyHostMap(bucket)
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(accessPointPolicyTimeout, func() *retry.RetryError {
		var requestErr error
		response, requestErr = client.Do("Oss", helper.JsonXmlParam("DELETE", "2019-05-17", "DeleteAccessPointPolicy", accessPointPolicyAction), query, nil, nil, hostMap, false)
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
// the planned policy differs from it only in JSON formatting (key order or
// whitespace), so an equivalent document is not reported as a change.
type jsonTemplateEquivalentModifier struct{}

func (jsonTemplateEquivalentModifier) Description(ctx context.Context) string {
	return "Suppresses changes when both policy documents are semantically equivalent JSON."
}

func (jsonTemplateEquivalentModifier) MarkdownDescription(ctx context.Context) string {
	return jsonTemplateEquivalentModifier{}.Description(ctx)
}

func (jsonTemplateEquivalentModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	if req.PlanValue.IsUnknown() || req.StateValue.IsUnknown() || req.StateValue.IsNull() {
		return
	}
	if req.PlanValue.Equal(req.StateValue) {
		return
	}
	equivalent, err := helper.CompareJsonTemplateAreEquivalent(req.StateValue.ValueString(), req.PlanValue.ValueString())
	if err == nil && equivalent {
		resp.PlanValue = req.StateValue
	}
}
