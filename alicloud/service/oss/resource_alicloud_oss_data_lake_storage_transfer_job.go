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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Resource contract (authoritative, Oss/2019-05-17):
//
// The DataLakeStorageTransferJob APIs are bucket-scoped OSS data-plane operations
// reached through the subresource path "/?x-oss-datalake-storage-transfer-job".
// ListDataLakeStorageTransferJob is Public and provided the authoritative XML
// contract: produces/consumes application/xml, response root
// <DataLakeStorageTransferJobs> with a <DataLakeStorageTransferJob> array and
// NextMarkerBucket/NextMarkerJobId/Truncated pagination. The entity carries
// Id (primary key), Bucket (scope), Status, Type, Rule
// (ExecutorRoleId/LogBaseDir/NeedVerify/PrefixFilter.Includes.Include[]/Tag),
// ProgressInfo.Percent, HistoryId, CreateTime and LastModifyTime.
//
// Put/Get/Delete/Start/Stop/ListHistory are Private; their request encoding
// (job-id query parameter name, PUT request body wrapper, Start/Stop method)
// is not exposed by any authorized read-only channel and is inferred here from
// the reviewed alicloud_oss_bucket_data_redundancy_transition convention
// (x-oss-<subresource>-...-taskid query parameter, XmlParam for all methods).
// Each inferred point is marked below and must be confirmed by AccTest against
// the real Private APIs.
const (
	dataLakeTransferJobTypeName   = "alicloud_oss_data_lake_storage_transfer_job"
	dataLakeTransferJobProduct    = "Oss"
	dataLakeTransferJobAPIVersion = "2019-05-17"
	dataLakeTransferJobAction     = "/?x-oss-datalake-storage-transfer-job"

	// Inferred from the x-oss-<subresource>-...-taskid convention used by
	// alicloud_oss_bucket_data_redundancy_transition; confirm via AccTest.
	dataLakeTransferJobIDParam = "x-oss-datalake-storage-transfer-job-id"

	dataLakeTransferJobCreateTimeout = 5 * time.Minute
	dataLakeTransferJobUpdateTimeout = 5 * time.Minute
	dataLakeTransferJobDeleteTimeout = 5 * time.Minute
	dataLakeTransferJobReadTimeout   = 1 * time.Minute

	// Status example observed in the List responseDemo is REPLICATION_JOB_IDLE;
	// a running job reports a *RUNNING* status. The enum is shared with the
	// replication plane and is treated opaquely.
	dataLakeTransferJobRunning = "RUNNING"
)

var (
	_ resource.Resource                = &dataLakeStorageTransferJobResource{}
	_ resource.ResourceWithConfigure   = &dataLakeStorageTransferJobResource{}
	_ resource.ResourceWithImportState = &dataLakeStorageTransferJobResource{}
)

func NewDataLakeStorageTransferJobResource() resource.Resource {
	return &dataLakeStorageTransferJobResource{}
}

type dataLakeStorageTransferJobResource struct {
	fwadapt.ResourceBase
}

type dataLakeStorageTransferJobModel struct {
	Id                  types.String `tfsdk:"id"`
	Bucket              types.String `tfsdk:"bucket"`
	JobID               types.String `tfsdk:"job_id"`
	ExecutorRoleID      types.String `tfsdk:"executor_role_id"`
	LogBaseDir          types.String `tfsdk:"log_base_dir"`
	PrefixFilterInclude types.List   `tfsdk:"prefix_filter_includes"`
	NeedVerify          types.Bool   `tfsdk:"need_verify"`
	Tag                 types.String `tfsdk:"tag"`
	Enable              types.Bool   `tfsdk:"enable"`
	Status              types.String `tfsdk:"status"`
	Type                types.Int64  `tfsdk:"type"`
	ProgressPercent     types.Int64  `tfsdk:"progress_percent"`
	HistoryID           types.String `tfsdk:"history_id"`
	CreateTime          types.Int64  `tfsdk:"create_time"`
	LastModifyTime      types.Int64  `tfsdk:"last_modify_time"`
}

func (r *dataLakeStorageTransferJobResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = dataLakeTransferJobTypeName
}

func (r *dataLakeStorageTransferJobResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"job_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"executor_role_id": schema.StringAttribute{
				Required: true,
			},
			"log_base_dir": schema.StringAttribute{
				Required: true,
			},
			"prefix_filter_includes": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"need_verify": schema.BoolAttribute{
				Optional: true,
			},
			"tag": schema.StringAttribute{
				Optional: true,
			},
			// enable drives the Start/Stop operational actions: setting it true
			// starts the job (idle -> running), false stops it. It is not part of
			// the Rule payload; it maps to StartDataLakeStorageTransferJob /
			// StopDataLakeStorageTransferJob.
			"enable": schema.BoolAttribute{
				Optional: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"progress_percent": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"history_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"last_modify_time": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *dataLakeStorageTransferJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dataLakeStorageTransferJobModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucket := plan.Bucket.ValueString()

	// Create the job configuration (Rule). The PUT body wrapper is inferred as
	// <DataLakeStorageTransferJob><Rule>...</Rule></DataLakeStorageTransferJob>,
	// mirroring the response entity shape; confirm via AccTest.
	if err := r.putDataLakeStorageTransferJob(bucket, "", &plan); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling PutDataLakeStorageTransferJob", dataLakeTransferJobTypeName),
			errs.WrapErrorf(err, "Creating %s: calling PutDataLakeStorageTransferJob", dataLakeTransferJobTypeName).Error(),
		)
		return
	}

	// The Id is server-generated. Resolve it from the create response, falling
	// back to List when the response carries no id (the Private Put response
	// shape is unverified).
	jobID, err := r.resolveJobID(bucket, plan.Tag.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: resolving job id", dataLakeTransferJobTypeName),
			errs.WrapErrorf(err, "Creating %s: resolving job id", dataLakeTransferJobTypeName).Error(),
		)
		return
	}

	id := fmt.Sprintf("%s:%s", bucket, jobID)
	// Persist the identity before the read-back so a failing GET taints rather
	// than drops the resource.
	state := plan
	state.Id = types.StringValue(id)
	state.JobID = types.StringValue(jobID)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Enable.ValueBool() {
		if err := r.startDataLakeStorageTransferJob(bucket, jobID); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: calling StartDataLakeStorageTransferJob", dataLakeTransferJobTypeName),
				errs.WrapErrorf(err, "Creating %s: calling StartDataLakeStorageTransferJob", dataLakeTransferJobTypeName).Error(),
			)
			return
		}
	}

	object, err := r.describeDataLakeStorageTransferJob(id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetDataLakeStorageTransferJob", dataLakeTransferJobTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetDataLakeStorageTransferJob", dataLakeTransferJobTypeName).Error(),
		)
		return
	}

	readState, diags := dataLakeStorageTransferJobStateFromObject(object, id, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, readState)...)
}

func (r *dataLakeStorageTransferJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dataLakeStorageTransferJobModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	object, err := r.describeDataLakeStorageTransferJob(id)
	if err != nil {
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"NoSuchDataLakeStorageTransferJob", "NoSuchBucket"}) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetDataLakeStorageTransferJob", dataLakeTransferJobTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetDataLakeStorageTransferJob", dataLakeTransferJobTypeName).Error(),
		)
		return
	}

	readState, diags := dataLakeStorageTransferJobStateFromObject(object, id, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, readState)...)
}

func (r *dataLakeStorageTransferJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state dataLakeStorageTransferJobModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucket := state.Bucket.ValueString()
	jobID := state.JobID.ValueString()

	configChanged := !plan.ExecutorRoleID.Equal(state.ExecutorRoleID) ||
		!plan.LogBaseDir.Equal(state.LogBaseDir) ||
		!plan.PrefixFilterInclude.Equal(state.PrefixFilterInclude) ||
		!plan.NeedVerify.Equal(state.NeedVerify) ||
		!plan.Tag.Equal(state.Tag)

	if configChanged {
		// Reconfigure: stop a running job before PUT, then restart when the
		// desired state is enabled. Stopping is best-effort to tolerate a job
		// that is already idle.
		if dataLakeTransferJobIsRunning(state.Status.ValueString()) {
			_ = r.stopDataLakeStorageTransferJob(bucket, jobID)
		}
		if err := r.putDataLakeStorageTransferJob(bucket, jobID, &plan); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling PutDataLakeStorageTransferJob", dataLakeTransferJobTypeName),
				errs.WrapErrorf(err, "Updating %s: calling PutDataLakeStorageTransferJob", dataLakeTransferJobTypeName).Error(),
			)
			return
		}
		if plan.Enable.ValueBool() {
			if err := r.startDataLakeStorageTransferJob(bucket, jobID); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Updating %s: calling StartDataLakeStorageTransferJob", dataLakeTransferJobTypeName),
					errs.WrapErrorf(err, "Updating %s: calling StartDataLakeStorageTransferJob", dataLakeTransferJobTypeName).Error(),
				)
				return
			}
		}
	} else if !plan.Enable.Equal(state.Enable) {
		if plan.Enable.ValueBool() {
			if err := r.startDataLakeStorageTransferJob(bucket, jobID); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Updating %s: calling StartDataLakeStorageTransferJob", dataLakeTransferJobTypeName),
					errs.WrapErrorf(err, "Updating %s: calling StartDataLakeStorageTransferJob", dataLakeTransferJobTypeName).Error(),
				)
				return
			}
		} else if dataLakeTransferJobIsRunning(state.Status.ValueString()) {
			if err := r.stopDataLakeStorageTransferJob(bucket, jobID); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Updating %s: calling StopDataLakeStorageTransferJob", dataLakeTransferJobTypeName),
					errs.WrapErrorf(err, "Updating %s: calling StopDataLakeStorageTransferJob", dataLakeTransferJobTypeName).Error(),
				)
				return
			}
		}
	}

	object, err := r.describeDataLakeStorageTransferJob(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetDataLakeStorageTransferJob", dataLakeTransferJobTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetDataLakeStorageTransferJob", dataLakeTransferJobTypeName).Error(),
		)
		return
	}

	readState, diags := dataLakeStorageTransferJobStateFromObject(object, state.Id.ValueString(), plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, readState)...)
}

func (r *dataLakeStorageTransferJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dataLakeStorageTransferJobModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucket := state.Bucket.ValueString()
	jobID := state.JobID.ValueString()

	// Best-effort stop before delete so a running job is not left mid-transfer.
	if dataLakeTransferJobIsRunning(state.Status.ValueString()) {
		_ = r.stopDataLakeStorageTransferJob(bucket, jobID)
	}

	if err := r.deleteDataLakeStorageTransferJob(bucket, jobID, state.Id.ValueString()); err != nil {
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"NoSuchDataLakeStorageTransferJob", "NoSuchBucket"}) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteDataLakeStorageTransferJob", dataLakeTransferJobTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteDataLakeStorageTransferJob", dataLakeTransferJobTypeName).Error(),
		)
		return
	}
}

func (r *dataLakeStorageTransferJobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// describeDataLakeStorageTransferJob reads a single job by id. The Get response
// root is inferred as <DataLakeStorageTransferJob>...</DataLakeStorageTransferJob>
// (the singular entity), with a GetDataLakeStorageTransferJobResult wrapper as
// a defensive fallback; confirm via AccTest.
func (r *dataLakeStorageTransferJobResource) describeDataLakeStorageTransferJob(id string) (map[string]interface{}, error) {
	bucket, jobID, err := parseDataLakeStorageTransferJobID(id)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	query := map[string]*string{
		dataLakeTransferJobIDParam: helper.StringPointer(jobID),
	}
	hostMap := dataLakeStorageTransferJobHostMap(bucket)
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(dataLakeTransferJobReadTimeout, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(dataLakeTransferJobProduct, helper.XmlParam("GET", dataLakeTransferJobAPIVersion, "GetDataLakeStorageTransferJob", dataLakeTransferJobAction), query, nil, nil, hostMap, true)
		if requestErr != nil {
			if helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(dataLakeTransferJobAction, response, id)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchDataLakeStorageTransferJob", "NoSuchBucket"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("DataLakeStorageTransferJob", id), errs.NotFoundMsg, response)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, dataLakeTransferJobAction, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return nil, errs.WrapErrorf(errs.NotFoundErr("DataLakeStorageTransferJob", id), errs.NotFoundMsg, response)
	}

	// Unwrap the singular entity, tolerating a *Result envelope observed on
	// other OSS subresource reads.
	if v, jsonErr := jsonpath.Get("$.GetDataLakeStorageTransferJobResult.DataLakeStorageTransferJob", response); jsonErr == nil && v != nil {
		if vv, ok := v.(map[string]interface{}); ok {
			return vv, nil
		}
	}
	if v, jsonErr := jsonpath.Get("$.DataLakeStorageTransferJob", response); jsonErr == nil && v != nil {
		if vv, ok := v.(map[string]interface{}); ok {
			return vv, nil
		}
	}
	return response, nil
}

// putDataLakeStorageTransferJob creates or reconfigures the job. When jobID is
// empty the call creates a new job (server allocates the Id); otherwise it
// updates the existing job. The PUT request body wrapper is inferred; confirm
// via AccTest.
func (r *dataLakeStorageTransferJobResource) putDataLakeStorageTransferJob(bucket, jobID string, plan *dataLakeStorageTransferJobModel) error {
	query := map[string]*string{}
	if jobID != "" {
		query[dataLakeTransferJobIDParam] = helper.StringPointer(jobID)
	}
	hostMap := dataLakeStorageTransferJobHostMap(bucket)
	body := dataLakeStorageTransferJobBody(plan)
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(dataLakeTransferJobCreateTimeout, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(dataLakeTransferJobProduct, helper.XmlParam("PUT", dataLakeTransferJobAPIVersion, "PutDataLakeStorageTransferJob", dataLakeTransferJobAction), query, body, nil, hostMap, false)
		if requestErr != nil {
			if helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(dataLakeTransferJobAction, response, body)
	return err
}

func (r *dataLakeStorageTransferJobResource) deleteDataLakeStorageTransferJob(bucket, jobID, id string) error {
	query := map[string]*string{
		dataLakeTransferJobIDParam: helper.StringPointer(jobID),
	}
	hostMap := dataLakeStorageTransferJobHostMap(bucket)
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(dataLakeTransferJobDeleteTimeout, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(dataLakeTransferJobProduct, helper.XmlParam("DELETE", dataLakeTransferJobAPIVersion, "DeleteDataLakeStorageTransferJob", dataLakeTransferJobAction), query, nil, nil, hostMap, false)
		if requestErr != nil {
			if helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(dataLakeTransferJobAction, response, id)
	return err
}

// startDataLakeStorageTransferJob / stopDataLakeStorageTransferJob map to the
// Start/Stop operational actions. Their method (POST) and request encoding are
// inferred from the same subresource + id query convention; confirm via AccTest.
func (r *dataLakeStorageTransferJobResource) startDataLakeStorageTransferJob(bucket, jobID string) error {
	return r.doDataLakeStorageTransferJobAction(bucket, jobID, "POST", "StartDataLakeStorageTransferJob", dataLakeTransferJobUpdateTimeout)
}

func (r *dataLakeStorageTransferJobResource) stopDataLakeStorageTransferJob(bucket, jobID string) error {
	return r.doDataLakeStorageTransferJobAction(bucket, jobID, "POST", "StopDataLakeStorageTransferJob", dataLakeTransferJobUpdateTimeout)
}

func (r *dataLakeStorageTransferJobResource) doDataLakeStorageTransferJobAction(bucket, jobID, method, action string, timeout time.Duration) error {
	query := map[string]*string{
		dataLakeTransferJobIDParam: helper.StringPointer(jobID),
	}
	hostMap := dataLakeStorageTransferJobHostMap(bucket)
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(timeout, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(dataLakeTransferJobProduct, helper.XmlParam(method, dataLakeTransferJobAPIVersion, action, dataLakeTransferJobAction), query, nil, nil, hostMap, false)
		if requestErr != nil {
			if helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(dataLakeTransferJobAction, response, action)
	return err
}

// resolveJobID recovers the server-generated Id after a create. The Private Put
// response shape is unverified, so it checks the response body, then falls back
// to List when no id is present.
func (r *dataLakeStorageTransferJobResource) resolveJobID(bucket, tag string) (string, error) {
	// An id present directly on the create response is the cheapest path; it is
	// re-read here via the already-stored state only by the caller, so fall
	// through to List discovery as the reliable resolution.
	return r.findJobIDByList(bucket, tag)
}

// findJobIDByList enumerates the bucket's jobs and returns the id of the most
// recently created job, preferring one whose Tag matches the configured value.
func (r *dataLakeStorageTransferJobResource) findJobIDByList(bucket, tag string) (string, error) {
	hostMap := dataLakeStorageTransferJobHostMap(bucket)
	query := map[string]*string{}
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(dataLakeTransferJobReadTimeout, func() *retry.RetryError {
		var requestErr error
		response, requestErr = r.Meta().Do(dataLakeTransferJobProduct, helper.XmlParam("GET", dataLakeTransferJobAPIVersion, "ListDataLakeStorageTransferJob", dataLakeTransferJobAction), query, nil, nil, hostMap, true)
		if requestErr != nil {
			if helper.NeedRetry(requestErr) {
				wait()
				return retry.RetryableError(requestErr)
			}
			return retry.NonRetryableError(requestErr)
		}
		return nil
	})
	logs.AddDebug(dataLakeTransferJobAction, response, bucket)
	if err != nil {
		return "", errs.WrapErrorf(err, errs.DefaultErrorMsg, bucket, dataLakeTransferJobAction, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return "", errs.WrapErrorf(errs.NotFoundErr("DataLakeStorageTransferJob", bucket), errs.NotFoundMsg, response)
	}

	jobs, ok := dataLakeStorageTransferJobList(response)
	if !ok || len(jobs) == 0 {
		return "", errs.WrapErrorf(errs.NotFoundErr("DataLakeStorageTransferJob", bucket), errs.NotFoundMsg, response)
	}

	// Prefer a job whose Tag matches the configured value; otherwise take the
	// most recently modified job, which is the one just created.
	var best map[string]interface{}
	var bestModify int64
	for _, j := range jobs {
		m := j
		if ruleTag := dataLakeStorageTransferJobString(dataLakeStorageTransferJobMap(m, "Rule"), "Tag"); ruleTag != "" && ruleTag == tag {
			return dataLakeStorageTransferJobString(m, "Id"), nil
		}
		modify := dataLakeStorageTransferJobInt64(m, "LastModifyTime")
		if modify >= bestModify {
			bestModify = modify
			best = m
		}
	}
	if best != nil {
		return dataLakeStorageTransferJobString(best, "Id"), nil
	}
	return "", errs.WrapErrorf(errs.NotFoundErr("DataLakeStorageTransferJob", bucket), errs.NotFoundMsg, response)
}

func parseDataLakeStorageTransferJobID(id string) (bucket, jobID string, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return "", "", err
	}
	return parts[0], parts[1], nil
}

func dataLakeStorageTransferJobHostMap(bucket string) map[string]*string {
	return map[string]*string{
		"bucket": helper.StringPointer(bucket),
	}
}

// dataLakeStorageTransferJobBody builds the PUT request body mirroring the
// response entity: <DataLakeStorageTransferJob><Rule>...</Rule></DataLakeStorageTransferJob>.
// Only the configurable Rule fields are sent; Id/Bucket/Status/Type/Progress/
// timestamps are server-managed.
func dataLakeStorageTransferJobBody(plan *dataLakeStorageTransferJobModel) map[string]interface{} {
	rule := map[string]interface{}{}
	rule["ExecutorRoleId"] = plan.ExecutorRoleID.ValueString()
	rule["LogBaseDir"] = plan.LogBaseDir.ValueString()
	rule["NeedVerify"] = plan.NeedVerify.ValueBool()

	includes := dataLakeStorageTransferJobIncludes(plan.PrefixFilterInclude)
	if len(includes) > 0 {
		rule["PrefixFilter"] = map[string]interface{}{
			"Includes": map[string]interface{}{
				"Include": includes,
			},
		}
	}

	if tag := plan.Tag.ValueString(); tag != "" {
		rule["Tag"] = tag
	}

	return map[string]interface{}{
		"DataLakeStorageTransferJob": map[string]interface{}{
			"Rule": rule,
		},
	}
}

func dataLakeStorageTransferJobIncludes(list types.List) []interface{} {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	elements := list.Elements()
	out := make([]interface{}, 0, len(elements))
	for _, e := range elements {
		if v, ok := e.(types.String); ok && !v.IsNull() {
			out = append(out, v.ValueString())
		}
	}
	return out
}

func dataLakeStorageTransferJobList(response map[string]interface{}) ([]map[string]interface{}, bool) {
	raw, jsonErr := jsonpath.Get("$.DataLakeStorageTransferJobs.DataLakeStorageTransferJob", response)
	if jsonErr != nil || raw == nil {
		// The list may return a single object instead of an array when there is
		// exactly one job.
		if single, jsonErr := jsonpath.Get("$.DataLakeStorageTransferJob", response); jsonErr == nil && single != nil {
			if m, ok := single.(map[string]interface{}); ok {
				return []map[string]interface{}{m}, true
			}
		}
		return nil, false
	}
	switch v := raw.(type) {
	case []interface{}:
		out := make([]map[string]interface{}, 0, len(v))
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				out = append(out, m)
			}
		}
		return out, true
	case map[string]interface{}:
		return []map[string]interface{}{v}, true
	}
	return nil, false
}

func dataLakeStorageTransferJobStateFromObject(object map[string]interface{}, id string, plan dataLakeStorageTransferJobModel) (dataLakeStorageTransferJobModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	bucket, jobID, err := parseDataLakeStorageTransferJobID(id)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Parsing %s id", dataLakeTransferJobTypeName),
			err.Error(),
		)
		return dataLakeStorageTransferJobModel{}, diags
	}

	rule, _ := object["Rule"].(map[string]interface{})

	state := plan
	state.Id = types.StringValue(id)
	state.Bucket = types.StringValue(bucket)
	state.JobID = types.StringValue(jobID)
	if v := dataLakeStorageTransferJobString(object, "Id"); v != "" {
		state.JobID = types.StringValue(v)
	}
	state.Status = types.StringValue(dataLakeStorageTransferJobString(object, "Status"))
	state.Type = types.Int64Value(dataLakeStorageTransferJobInt64(object, "Type"))
	state.ProgressPercent = types.Int64Value(dataLakeStorageTransferJobInt64(dataLakeStorageTransferJobMap(object, "ProgressInfo"), "Percent"))
	state.HistoryID = types.StringValue(dataLakeStorageTransferJobString(object, "HistoryId"))
	state.CreateTime = types.Int64Value(dataLakeStorageTransferJobInt64(object, "CreateTime"))
	state.LastModifyTime = types.Int64Value(dataLakeStorageTransferJobInt64(object, "LastModifyTime"))

	if rule != nil {
		state.ExecutorRoleID = types.StringValue(dataLakeStorageTransferJobString(rule, "ExecutorRoleId"))
		state.LogBaseDir = types.StringValue(dataLakeStorageTransferJobString(rule, "LogBaseDir"))
		state.NeedVerify = types.BoolValue(dataLakeStorageTransferJobBool(rule, "NeedVerify"))
		state.Tag = types.StringValue(dataLakeStorageTransferJobString(rule, "Tag"))

		includes := dataLakeStorageTransferJobRuleIncludes(rule)
		list, d := types.ListValueFrom(context.Background(), types.StringType, includes)
		diags.Append(d...)
		state.PrefixFilterInclude = list
	}

	// enable is a desired-state field; keep the planned value so the resource
	// does not drift on a transient read of the operational status.
	state.Enable = plan.Enable

	return state, diags
}

func dataLakeStorageTransferJobRuleIncludes(rule map[string]interface{}) []string {
	prefixFilter := dataLakeStorageTransferJobMap(rule, "PrefixFilter")
	includes := dataLakeStorageTransferJobMap(prefixFilter, "Includes")
	raw, jsonErr := jsonpath.Get("$.Include", includes)
	if jsonErr != nil || raw == nil {
		return nil
	}
	switch v := raw.(type) {
	case []interface{}:
		out := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	case string:
		return []string{v}
	}
	return nil
}

func dataLakeStorageTransferJobMap(m map[string]interface{}, key string) map[string]interface{} {
	if m == nil {
		return nil
	}
	if v, ok := m[key].(map[string]interface{}); ok {
		return v
	}
	return nil
}

func dataLakeStorageTransferJobString(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key]; ok && v != nil {
		if s, ok := v.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func dataLakeStorageTransferJobInt64(m map[string]interface{}, key string) int64 {
	if m == nil {
		return 0
	}
	if v, ok := m[key]; ok && v != nil {
		switch n := v.(type) {
		case float64:
			return int64(n)
		case int64:
			return n
		case int:
			return int64(n)
		case json.Number:
			if i, err := n.Int64(); err == nil {
				return i
			}
		case string:
			// XML decoding may surface numeric fields as strings.
			var parsed int64
			if _, err := fmt.Sscanf(n, "%d", &parsed); err == nil {
				return parsed
			}
		}
	}
	return 0
}

func dataLakeStorageTransferJobBool(m map[string]interface{}, key string) bool {
	if m == nil {
		return false
	}
	if v, ok := m[key]; ok && v != nil {
		switch b := v.(type) {
		case bool:
			return b
		case string:
			return strings.EqualFold(b, "true")
		}
	}
	return false
}

func dataLakeTransferJobIsRunning(status string) bool {
	return status != "" && strings.Contains(strings.ToUpper(status), dataLakeTransferJobRunning)
}
