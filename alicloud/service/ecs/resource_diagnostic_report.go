// Package ecs holds the Framework-native ECS resources. DiagnosticReport is the
// first one; it ports the legacy SDKv2 alicloud/resource_alicloud_ecs_diagnostic_report.go
// to the Plugin Framework, preserving the create-only schema, the write-only
// AdditionalOptions map, the dual NotFound paths and the post-delete poll.
package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const diagnosticReportTypeName = "alicloud_ecs_diagnostic_report"

// Mirror the SDKv2 timeouts: Create and Delete allow 5 minutes, the Describe
// read/poll allows 1 minute of throttling backoff.
const (
	diagnosticReportCreateTimeout   = 5 * time.Minute
	diagnosticReportDeleteTimeout   = 5 * time.Minute
	diagnosticReportDescribeTimeout = 1 * time.Minute
)

var (
	_ resource.Resource                = &diagnosticReportResource{}
	_ resource.ResourceWithImportState = &diagnosticReportResource{}
)

func NewDiagnosticReportResource() resource.Resource {
	return &diagnosticReportResource{}
}

type diagnosticReportResource struct {
	fwadapt.ResourceBase
}

type diagnosticReportResourceModel struct {
	Id                types.String `tfsdk:"id"`
	ResourceId        types.String `tfsdk:"resource_id"`
	MetricSetId       types.String `tfsdk:"metric_set_id"`
	StartTime         types.String `tfsdk:"start_time"`
	EndTime           types.String `tfsdk:"end_time"`
	AdditionalOptions types.Map    `tfsdk:"additional_options"`
}

func (r *diagnosticReportResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = diagnosticReportTypeName
}

func (r *diagnosticReportResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a ECS Diagnostic Report resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the Diagnostic Report, equivalent to the ReportId returned by CreateDiagnosticReport.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the resource to diagnose, such as an ECS instance ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"metric_set_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the diagnostic metric set. If it is empty, the default set dms-instancedefault of the ECS instance is used.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"start_time": schema.StringAttribute{
				Optional:    true,
				Description: "Start time. Only takes effect for diagnostic indicators that do not need to execute cloud assistant instructions in GuestOS.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"end_time": schema.StringAttribute{
				Optional:    true,
				Description: "End time. Only takes effect for diagnostic indicators that do not need to execute cloud assistant instructions in GuestOS.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// additional_options is write-only: the API accepts it on Create but
			// DescribeDiagnosticReportAttributes never returns it, so it is never
			// refreshed. It is kept out of import verification for the same reason.
			"additional_options": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Additional options for the diagnostic report, as a map of string key-value pairs sent as a JSON-encoded string.",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *diagnosticReportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan diagnosticReportResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"RegionId":   r.Client().RegionId,
		"ResourceId": plan.ResourceId.ValueString(),
	}
	if !plan.MetricSetId.IsNull() {
		request["MetricSetId"] = plan.MetricSetId.ValueString()
	}
	if !plan.StartTime.IsNull() {
		request["StartTime"] = plan.StartTime.ValueString()
	}
	if !plan.EndTime.IsNull() {
		request["EndTime"] = plan.EndTime.ValueString()
	}
	if !plan.AdditionalOptions.IsNull() && !plan.AdditionalOptions.IsUnknown() {
		// additional_options is a write-only map<string,string>. ElementsAs to a concrete
		// map[string]string — a map[string]interface{} target fails because the framework's
		// reflection cannot map tftypes.String into interface{}. The CreateDiagnosticReport
		// API accepts additional_options as a JSON "extension", so the marshaled map is
		// sent as a JSON string, mirroring the kms ephemeral_plaintext EncryptionContext
		// pattern and the provider-wide practice of pre-marshaling complex values before
		// RpcPost rather than relying on the tea-rpc gateway to encode a raw map.
		opts := make(map[string]string)
		resp.Diagnostics.Append(plan.AdditionalOptions.ElementsAs(ctx, &opts, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		encoded, merr := json.Marshal(opts)
		if merr != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: encoding additional_options", diagnosticReportTypeName),
				errs.WrapErrorf(merr, "Creating %s: encoding additional_options", diagnosticReportTypeName).Error(),
			)
			return
		}
		request["AdditionalOptions"] = string(encoded)
	}

	response, err := callRpcWithRetry(ctx, diagnosticReportCreateTimeout, func() (map[string]interface{}, error) {
		return r.Client().RpcPost("Ecs", "2014-05-26", "CreateDiagnosticReport", nil, request, true)
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDiagnosticReport", diagnosticReportTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateDiagnosticReport", diagnosticReportTypeName).Error(),
		)
		return
	}

	reportId, idErr := diagnosticReportId(response)
	if idErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDiagnosticReport", diagnosticReportTypeName),
			fmt.Sprintf("The response did not carry a valid ReportId: %s; raw response: %v", idErr.Error(), response),
		)
		return
	}

	// Persist the just-created ReportId (and the write-only additional_options carried
	// from plan) into state BEFORE refreshing. CreateDiagnosticReport has already
	// succeeded, so a transient Describe failure or brief invisibility must not lose
	// the resource; the plan fields are preserved here and only the echoed fields are
	// overwritten on a successful refresh below. The next refresh will retry.
	plan.Id = types.StringValue(reportId)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	object, derr := describeDiagnosticReport(ctx, r.Client(), reportId)
	if derr != nil {
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Creating %s: refreshing after create", diagnosticReportTypeName),
			errs.WrapErrorf(derr, "the report %s was created but could not be refreshed right away; the ID is retained in state and the next refresh will retry", reportId).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Creating %s: refreshing after create", diagnosticReportTypeName),
			fmt.Sprintf("The report %s was created but is not yet visible via DescribeDiagnosticReportAttributes; the ID is retained in state and the next refresh will retry", reportId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, diagnosticReportState(object, plan.AdditionalOptions))...)
}

func (r *diagnosticReportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state diagnosticReportResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := describeDiagnosticReport(ctx, r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeDiagnosticReportAttributes", diagnosticReportTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeDiagnosticReportAttributes", diagnosticReportTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, diagnosticReportState(object, state.AdditionalOptions))...)
}

// Update is unreachable in practice: every user-settable attribute carries
// RequiresReplace, so any plan change triggers replacement rather than an
// in-place update. The method exists to satisfy the Resource interface and copies
// the plan into state defensively.
func (r *diagnosticReportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan diagnosticReportResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *diagnosticReportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state diagnosticReportResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := callRpcWithRetry(ctx, diagnosticReportDeleteTimeout, func() (map[string]interface{}, error) {
		return r.Client().RpcPost("Ecs", "2014-05-26", "DeleteDiagnosticReports", nil, map[string]interface{}{
			"ReportIds.1": state.Id.ValueString(),
			"RegionId":    r.Client().RegionId,
		}, true)
	})
	if err != nil && !isReportNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteDiagnosticReports", diagnosticReportTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteDiagnosticReports", diagnosticReportTypeName).Error(),
		)
		return
	}
	// NotFound on delete is a no-op success.

	// DeleteDiagnosticReports may be asynchronous; poll DescribeDiagnosticReportAttributes
	// until the report disappears so a later refresh observes no stale row.
	if err := waitForDiagnosticReportDeletion(ctx, r.Client(), state.Id.ValueString(), diagnosticReportDeleteTimeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: waiting for deletion", diagnosticReportTypeName),
			errs.WrapErrorf(err, "Deleting %s: waiting for deletion of %s", diagnosticReportTypeName, state.Id.ValueString()).Error(),
		)
	}
}

func (r *diagnosticReportResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// describeDiagnosticReport fetches a DiagnosticReport by id. It distinguishes the two
// legal "no longer exists" paths — the 403 InvalidParameter.ReportNotExist error and the
// 200 envelope that carries only RequestId with no ReportId — both of which return a nil
// object, nil error. Any other error propagates. A response that is nil, carries a
// ReportId of the wrong type, is empty, or whose ReportId does not match the requested id
// is a bad response: it is returned as an error so Read/Delete retain the existing state
// rather than silently removing the resource or writing a null id. Other 4xx/5xx are not
// treated as not-found.
func describeDiagnosticReport(ctx context.Context, client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	response, err := callRpcWithRetry(ctx, diagnosticReportDescribeTimeout, func() (map[string]interface{}, error) {
		return client.RpcPost("Ecs", "2014-05-26", "DescribeDiagnosticReportAttributes", nil, map[string]interface{}{
			"ReportId": id,
			"RegionId": client.RegionId,
		}, true)
	})
	if err != nil {
		// Only the specific ReportNotExist code is a legal "gone" signal; other
		// 4xx/5xx (including unrelated 403/404) propagate so Read/Delete do not
		// silently drop a resource that may still exist.
		if isReportNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("DescribeDiagnosticReportAttributes returned an empty response for report %s", id)
	}
	// A present ReportId must be a non-empty string that matches the requested id.
	// A null ReportId ({"ReportId":null,...}) is NOT a legal "gone" envelope: it is
	// a malformed response and must surface as an error so Read/Delete retain the
	// existing state rather than silently removing the resource. A wrong type, an
	// empty string, or a mismatched id is likewise a bad response.
	raw, present := response["ReportId"]
	if present {
		if raw == nil {
			return nil, fmt.Errorf("DescribeDiagnosticReportAttributes returned a null ReportId for report %s", id)
		}
		s, ok := raw.(string)
		if !ok {
			return nil, fmt.Errorf("DescribeDiagnosticReportAttributes returned a ReportId of unexpected type %T for report %s", raw, id)
		}
		if s == "" {
			return nil, fmt.Errorf("DescribeDiagnosticReportAttributes returned an empty ReportId for report %s", id)
		}
		if s != id {
			return nil, fmt.Errorf("DescribeDiagnosticReportAttributes returned ReportId %q which does not match the requested %q", s, id)
		}
		return response, nil
	}
	// ReportId key is absent. The ONLY legal "report no longer exists" 200 envelope
	// carries only RequestId with no ReportId and no other business fields. An empty
	// response object ({}) or one carrying unrelated object fields (e.g.
	// {"ResourceId":"i-other"}) is malformed and must surface as an error so
	// Read/Delete retain the existing state rather than silently removing the
	// resource. Other 4xx/5xx are not reached here; they propagate as errors above.
	if isLegalGoneEnvelope(response) {
		return nil, nil
	}
	return nil, fmt.Errorf("DescribeDiagnosticReportAttributes returned a malformed response for report %s: %v", id, response)
}

// isLegalGoneEnvelope reports whether response is the legal 200 "report no longer
// exists" envelope: a non-nil map carrying only RequestId with no ReportId and no
// other business fields. The ECS API returns this minimal envelope when a report is
// gone; an empty object ({}) or one carrying unrelated business fields is malformed
// and must not be treated as gone, so Read/Delete retain the existing state.
func isLegalGoneEnvelope(response map[string]interface{}) bool {
	rid, ok := response["RequestId"]
	if !ok || rid == nil {
		return false
	}
	for k := range response {
		// The tea-rpc SDK injects transport metadata (notably "_headers") into the
		// response map alongside the JSON body fields. Such underscore-prefixed keys
		// are not business fields and are ignored here: the legal "report no longer
		// exists" envelope still carries only RequestId among business fields. Any
		// other business field (e.g. ResourceId) keeps this from being "gone".
		if k != "RequestId" && !strings.HasPrefix(k, "_") {
			return false
		}
	}
	return true
}

// DiagnosticReportExists reports whether a diagnostic report is present. It returns
// (true, nil) when the report exists and (false, nil) when the report is gone, covering
// the two legal "no longer exists" paths the ECS API exposes — the 403
// InvalidParameter.ReportNotExist error and the 200 envelope that carries only RequestId
// with no ReportId. Any other failure, including a transient RPC error or a malformed
// response, returns (false, err) so a caller cannot mistake a bad response for a
// successfully deleted report. It is the destroy verification used by the acceptance test.
func DiagnosticReportExists(client *connectivity.AliyunClient, id string) (bool, error) {
	object, err := describeDiagnosticReport(context.Background(), client, id)
	if err != nil {
		return false, err
	}
	return object != nil, nil
}

// waitForDiagnosticReportDeletion polls DescribeDiagnosticReportAttributes until the
// report is gone or the deadline expires, mirroring the SDKv2 post-delete poll.
func waitForDiagnosticReportDeletion(ctx context.Context, client *connectivity.AliyunClient, id string, timeout time.Duration) error {
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	deadline := time.Now().Add(timeout)
	var lastStatus interface{}
	for {
		object, err := describeDiagnosticReport(ctx, client, id)
		if err != nil {
			return err
		}
		if object == nil {
			return nil
		}
		lastStatus = object["Status"]
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("waiting for deletion of ecs diagnostic report %s timeout, last status: %v", id, lastStatus)
		}
		wait()
	}
}

// callRpcWithRetry invokes fn once, then retries transient errors reported by
// helper.NeedRetry with the same incremental backoff the SDKv2 line used, until ctx
// is done or timeout elapses. Non-retryable errors are returned on the first try.
func callRpcWithRetry(ctx context.Context, timeout time.Duration, fn func() (map[string]interface{}, error)) (map[string]interface{}, error) {
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	deadline := time.Now().Add(timeout)
	for {
		response, err := fn()
		if err == nil {
			return response, nil
		}
		if !helper.NeedRetry(err) {
			return response, err
		}
		if ctx.Err() != nil {
			return response, err
		}
		if time.Now().After(deadline) {
			return response, err
		}
		wait()
	}
}

// diagnosticReportState maps a DescribeDiagnosticReportAttributes response onto the
// resource model. additional_options is preserved from the caller (config or prior
// state) because the API never returns it. Empty strings become null.
func diagnosticReportState(object map[string]interface{}, additionalOptions types.Map) diagnosticReportResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	return diagnosticReportResourceModel{
		Id:                str("ReportId"),
		ResourceId:        str("ResourceId"),
		MetricSetId:       str("MetricSetId"),
		StartTime:         str("StartTime"),
		EndTime:           str("EndTime"),
		AdditionalOptions: additionalOptions,
	}
}

// diagnosticReportId extracts the ReportId from a CreateDiagnosticReport response. It
// accepts only a non-empty string; a missing key, nil, or any non-string type (number,
// map, slice) is rejected with an error so a malformed response cannot yield a
// pseudo-ID. This replaces the SDKv2 fmt.Sprint fallback, which would have coerced 42,
// {"k":"v"} or [] into a usable-looking but wrong ID.
func diagnosticReportId(response map[string]interface{}) (string, error) {
	raw, ok := response["ReportId"]
	if !ok || raw == nil {
		return "", fmt.Errorf("the response carried no ReportId")
	}
	s, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("the response ReportId is not a string: %T", raw)
	}
	if s == "" {
		return "", fmt.Errorf("the response ReportId is an empty string")
	}
	return s, nil
}

// isReportNotExist reports whether err is the ECS report-not-found error. The code
// arrives on a *tea.SDKError from RpcPost; the substring fallback covers wrapped
// forms, matching IsExpectedErrors' behavior on the SDKv2 line.
func isReportNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil && *e.Code == "InvalidParameter.ReportNotExist" {
		return true
	}
	return strings.Contains(err.Error(), "InvalidParameter.ReportNotExist")
}
