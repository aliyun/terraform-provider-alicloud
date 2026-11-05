package oss

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	bucketRequesterQoSInfoTypeName = "alicloud_oss_bucket_requester_qos_info"
	bucketRequesterQoSInfoAction   = "/?requesterQosInfo"

	bucketRequesterQoSInfoProduct    = "Oss"
	bucketRequesterQoSInfoAPIVersion = "2019-05-17"

	// bucketRequesterQoSInfoTimeout bounds the retry envelope of the QoS writes
	// (Put and Delete); the read has its own one-minute budget.
	bucketRequesterQoSInfoTimeout = 5 * time.Minute
)

var (
	_ resource.Resource                = &bucketRequesterQoSInfoResource{}
	_ resource.ResourceWithConfigure   = &bucketRequesterQoSInfoResource{}
	_ resource.ResourceWithImportState = &bucketRequesterQoSInfoResource{}
)

func NewBucketRequesterQoSInfoResource() resource.Resource {
	return &bucketRequesterQoSInfoResource{}
}

type bucketRequesterQoSInfoResource struct {
	fwadapt.ResourceBase
}

// bucketRequesterQoSInfoResourceModel is the Terraform state for a single
// requester QoS configuration attached to a bucket. The composite identity is
// <bucket>:<qos_requester>, matching the OSS API's bucket(host) +
// qosRequester(query) addressing.
type bucketRequesterQoSInfoResourceModel struct {
	Id types.String `tfsdk:"id"`

	// Identity — both parts force replacement: a different bucket or requester
	// is a different remote entity, not an in-place update.
	Bucket       types.String `tfsdk:"bucket"`
	QosRequester types.String `tfsdk:"qos_requester"`

	// QoS quota fields. PutBucketRequesterQoSInfo upserts the whole
	// QoSConfiguration, so the resource owns every field and they are required.
	ExtranetDownloadBandwidth types.Int64 `tfsdk:"extranet_download_bandwidth"`
	ExtranetQps               types.Int64 `tfsdk:"extranet_qps"`
	ExtranetUploadBandwidth   types.Int64 `tfsdk:"extranet_upload_bandwidth"`
	IntranetDownloadBandwidth types.Int64 `tfsdk:"intranet_download_bandwidth"`
	IntranetQps               types.Int64 `tfsdk:"intranet_qps"`
	IntranetUploadBandwidth   types.Int64 `tfsdk:"intranet_upload_bandwidth"`
	TotalDownloadBandwidth    types.Int64 `tfsdk:"total_download_bandwidth"`
	TotalQps                  types.Int64 `tfsdk:"total_qps"`
	TotalUploadBandwidth      types.Int64 `tfsdk:"total_upload_bandwidth"`
}

func (r *bucketRequesterQoSInfoResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = bucketRequesterQoSInfoTypeName
}

func (r *bucketRequesterQoSInfoResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bucket": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"qos_requester": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"extranet_download_bandwidth": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"extranet_qps": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"extranet_upload_bandwidth": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"intranet_download_bandwidth": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"intranet_qps": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"intranet_upload_bandwidth": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"total_download_bandwidth": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"total_qps": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"total_upload_bandwidth": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *bucketRequesterQoSInfoResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan bucketRequesterQoSInfoResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucket := plan.Bucket.ValueString()
	qosRequester := plan.QosRequester.ValueString()
	if err := r.putBucketRequesterQoSInfo(bucket, qosRequester, plan.qosConfigurationBody()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling PutBucketRequesterQoSInfo", bucketRequesterQoSInfoTypeName),
			errs.WrapErrorf(err, "Creating %s: calling PutBucketRequesterQoSInfo", bucketRequesterQoSInfoTypeName).Error(),
		)
		return
	}

	id := fmt.Sprintf("%s:%s", bucket, qosRequester)
	// Persist the resource identity before the read-back: the QoS config exists
	// remotely from here on, so a failing GET must not leave Terraform without a
	// resource to taint and destroy.
	state := plan
	state.Id = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := r.describeBucketRequesterQoSInfo(id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBucketRequesterQoSInfo", bucketRequesterQoSInfoTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetBucketRequesterQoSInfo", bucketRequesterQoSInfoTypeName).Error(),
		)
		return
	}

	refreshed := bucketRequesterQoSInfoState(object)
	refreshed.Id = types.StringValue(id)
	refreshed.Bucket = plan.Bucket
	refreshed.QosRequester = plan.QosRequester
	resp.Diagnostics.Append(resp.State.Set(ctx, refreshed)...)
}

func (r *bucketRequesterQoSInfoResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state bucketRequesterQoSInfoResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	object, err := r.describeBucketRequesterQoSInfo(id)
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBucketRequesterQoSInfo", bucketRequesterQoSInfoTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetBucketRequesterQoSInfo", bucketRequesterQoSInfoTypeName).Error(),
		)
		return
	}

	bucket, qosRequester, err := parseBucketRequesterQoSInfoId(id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing resource id", bucketRequesterQoSInfoTypeName),
			err.Error(),
		)
		return
	}

	model := bucketRequesterQoSInfoState(object)
	model.Id = types.StringValue(id)
	model.Bucket = types.StringValue(bucket)
	model.QosRequester = types.StringValue(qosRequester)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *bucketRequesterQoSInfoResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state bucketRequesterQoSInfoResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Bucket and qos_requester are ForceReplace, so an Update only fires when the
	// QoS quota fields change; the upsert Put replaces the whole configuration.
	if err := r.putBucketRequesterQoSInfo(state.Bucket.ValueString(), state.QosRequester.ValueString(), plan.qosConfigurationBody()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling PutBucketRequesterQoSInfo", bucketRequesterQoSInfoTypeName),
			errs.WrapErrorf(err, "Updating %s: calling PutBucketRequesterQoSInfo", bucketRequesterQoSInfoTypeName).Error(),
		)
		return
	}

	object, err := r.describeBucketRequesterQoSInfo(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBucketRequesterQoSInfo", bucketRequesterQoSInfoTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetBucketRequesterQoSInfo", bucketRequesterQoSInfoTypeName).Error(),
		)
		return
	}

	model := bucketRequesterQoSInfoState(object)
	model.Id = state.Id
	model.Bucket = plan.Bucket
	model.QosRequester = plan.QosRequester
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *bucketRequesterQoSInfoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state bucketRequesterQoSInfoResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.deleteBucketRequesterQoSInfo(state.Bucket.ValueString(), state.QosRequester.ValueString(), state.Id.ValueString())
	if err != nil {
		// A vanished requester config or a vanished bucket both count as deleted.
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, bucketRequesterQoSInfoNotFoundCodes) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteBucketRequesterQoSInfo", bucketRequesterQoSInfoTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteBucketRequesterQoSInfo", bucketRequesterQoSInfoTypeName).Error(),
		)
		return
	}
}

func (r *bucketRequesterQoSInfoResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// bucketRequesterQoSInfoNotFoundCodes are the service-side codes that indicate
// the requester QoS config (or its bucket) no longer exists; they are treated as
// a successful delete and as RemoveResource on read.
var bucketRequesterQoSInfoNotFoundCodes = []string{
	"NoSuchRequesterQoSInfo",
	"NoSuchBucket",
	"NoSuchRequesterQosConfiguration",
}

// parseBucketRequesterQoSInfoId splits <bucket>:<qos_requester>.
func parseBucketRequesterQoSInfoId(id string) (bucket, qosRequester string, err error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid Resource Id %s. Expected bucket:qos_requester, got %s", id, id)
	}
	return parts[0], parts[1], nil
}

func bucketRequesterQoSInfoQuery(qosRequester string) map[string]*string {
	return map[string]*string{
		"qosRequester": helper.StringPointer(qosRequester),
	}
}

func bucketRequesterQoSInfoHostMap(bucket string) map[string]*string {
	return map[string]*string{
		"bucket": helper.StringPointer(bucket),
	}
}

// qosConfigurationBody builds the JSON request body for PutBucketRequesterQoSInfo.
// The OSS gateway converts the JSON object into the XML QoSConfiguration element
// the API consumes.
func (m bucketRequesterQoSInfoResourceModel) qosConfigurationBody() string {
	qos := map[string]interface{}{
		"ExtranetDownloadBandwidth": m.ExtranetDownloadBandwidth.ValueInt64(),
		"ExtranetQps":               m.ExtranetQps.ValueInt64(),
		"ExtranetUploadBandwidth":   m.ExtranetUploadBandwidth.ValueInt64(),
		"IntranetDownloadBandwidth": m.IntranetDownloadBandwidth.ValueInt64(),
		"IntranetQps":               m.IntranetQps.ValueInt64(),
		"IntranetUploadBandwidth":   m.IntranetUploadBandwidth.ValueInt64(),
		"TotalDownloadBandwidth":    m.TotalDownloadBandwidth.ValueInt64(),
		"TotalQps":                  m.TotalQps.ValueInt64(),
		"TotalUploadBandwidth":      m.TotalUploadBandwidth.ValueInt64(),
	}
	body, err := json.Marshal(map[string]interface{}{"QoSConfiguration": qos})
	if err != nil {
		// Marshal of a flat int map cannot fail in practice; fall back to an empty
		// configuration rather than blocking the whole apply.
		logs.AddDebug(bucketRequesterQoSInfoAction, err, qos)
		return "{}"
	}
	return string(body)
}

func (r *bucketRequesterQoSInfoResource) putBucketRequesterQoSInfo(bucket, qosRequester, body string) error {
	query := bucketRequesterQoSInfoQuery(qosRequester)
	hostMap := bucketRequesterQoSInfoHostMap(bucket)
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(bucketRequesterQoSInfoTimeout, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(bucketRequesterQoSInfoProduct, helper.JsonXmlParam("PUT", bucketRequesterQoSInfoAPIVersion, "PutBucketRequesterQoSInfo", bucketRequesterQoSInfoAction), query, body, nil, hostMap, false)
		if requestErr != nil {
			if helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(bucketRequesterQoSInfoAction, response, body)
	return err
}

func (r *bucketRequesterQoSInfoResource) deleteBucketRequesterQoSInfo(bucket, qosRequester, id string) error {
	query := bucketRequesterQoSInfoQuery(qosRequester)
	hostMap := bucketRequesterQoSInfoHostMap(bucket)
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(bucketRequesterQoSInfoTimeout, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(bucketRequesterQoSInfoProduct, helper.JsonXmlParam("DELETE", bucketRequesterQoSInfoAPIVersion, "DeleteBucketRequesterQoSInfo", bucketRequesterQoSInfoAction), query, nil, nil, hostMap, false)
		if requestErr != nil {
			if helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(bucketRequesterQoSInfoAction, response, id)
	return err
}

// describeBucketRequesterQoSInfo fetches a requester QoS configuration. The id
// is the resource id: <bucket>:<qos_requester>.
func (r *bucketRequesterQoSInfoResource) describeBucketRequesterQoSInfo(id string) (map[string]interface{}, error) {
	bucket, qosRequester, err := parseBucketRequesterQoSInfoId(id)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	query := bucketRequesterQoSInfoQuery(qosRequester)
	hostMap := bucketRequesterQoSInfoHostMap(bucket)
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(bucketRequesterQoSInfoProduct, helper.XmlJsonParam("GET", bucketRequesterQoSInfoAPIVersion, "GetBucketRequesterQoSInfo", bucketRequesterQoSInfoAction), query, nil, nil, hostMap, true)
		if requestErr != nil {
			if helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(bucketRequesterQoSInfoAction, response, id)
	if err != nil {
		if errs.IsExpectedErrors(err, bucketRequesterQoSInfoNotFoundCodes) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("BucketRequesterQoSInfo", id), errs.NotFoundMsg, response)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, bucketRequesterQoSInfoAction, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return nil, errs.WrapErrorf(errs.NotFoundErr("BucketRequesterQoSInfo", id), errs.NotFoundMsg, response)
	}

	// OSS sub-resource responses are wrapped in a *Result element when served
	// through the OpenAPI gateway; extract it defensively and fall back to the
	// raw response when no wrapper is present.
	if v, jsonErr := jsonpath.Get("$.GetBucketRequesterQoSInfoResult", response); jsonErr == nil && v != nil {
		if vv, ok := v.(map[string]interface{}); ok {
			response = vv
		}
	}

	// The RequesterQoSInfo container holds the requester id and the QoS
	// configuration; fall back to the raw response when the container is absent.
	if v, jsonErr := jsonpath.Get("$.RequesterQoSInfo", response); jsonErr == nil && v != nil {
		if vv, ok := v.(map[string]interface{}); ok {
			return vv, nil
		}
	}
	return response, nil
}

// bucketRequesterQoSInfoState maps a fetched RequesterQoSInfo document onto the
// resource model. The requester id is read back from the API when present.
func bucketRequesterQoSInfoState(object map[string]interface{}) bucketRequesterQoSInfoResourceModel {
	qos := readQoSConfigurationMap(object["QoSConfiguration"])
	return bucketRequesterQoSInfoResourceModel{
		QosRequester:              stringOrNull(object["Requester"]),
		ExtranetDownloadBandwidth: int64OrNull(qos["ExtranetDownloadBandwidth"]),
		ExtranetQps:               int64OrNull(qos["ExtranetQps"]),
		ExtranetUploadBandwidth:   int64OrNull(qos["ExtranetUploadBandwidth"]),
		IntranetDownloadBandwidth: int64OrNull(qos["IntranetDownloadBandwidth"]),
		IntranetQps:               int64OrNull(qos["IntranetQps"]),
		IntranetUploadBandwidth:   int64OrNull(qos["IntranetUploadBandwidth"]),
		TotalDownloadBandwidth:    int64OrNull(qos["TotalDownloadBandwidth"]),
		TotalQps:                  int64OrNull(qos["TotalQps"]),
		TotalUploadBandwidth:      int64OrNull(qos["TotalUploadBandwidth"]),
	}
}

// readQoSConfigurationMap coerces the QoSConfiguration value to a map whether
// the gateway returned it as a JSON object or as a JSON string.
func readQoSConfigurationMap(raw interface{}) map[string]interface{} {
	if raw == nil {
		return nil
	}
	if m, ok := raw.(map[string]interface{}); ok {
		return m
	}
	if s, ok := raw.(string); ok && s != "" {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(s), &m); err == nil {
			return m
		}
	}
	return nil
}

func int64OrNull(raw interface{}) types.Int64 {
	if raw == nil {
		return types.Int64Null()
	}
	switch v := raw.(type) {
	case float64:
		return types.Int64Value(int64(v))
	case int64:
		return types.Int64Value(v)
	case int:
		return types.Int64Value(int64(v))
	case string:
		if v == "" {
			return types.Int64Null()
		}
		var n int64
		if _, err := fmt.Sscanf(v, "%d", &n); err == nil {
			return types.Int64Value(n)
		}
		return types.Int64Null()
	}
	return types.Int64Null()
}

func stringOrNull(raw interface{}) types.String {
	if raw == nil {
		return types.StringNull()
	}
	if s, ok := raw.(string); ok {
		return types.StringValue(s)
	}
	return types.StringNull()
}
