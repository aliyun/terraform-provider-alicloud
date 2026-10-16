package paiflow

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const pipelineRunTypeName = "alicloud_pai_flow_pipeline_run"

// Status values returned by GetPipelineRun. Only "Succeeded" is a success
// terminal; the rest are either transient or failure terminals.
const (
	statusSucceeded  = "Succeeded"
	statusFailed     = "Failed"
	statusTerminated = "Terminated"
	statusSkipped    = "Skipped"
)

// defaultCreateTimeout bounds how long Create blocks waiting for the run to
// settle. A pipeline run is a one-shot execution, so Create waits until the
// run reaches a terminal status rather than returning while still running.
const defaultCreateTimeout = 30 * time.Minute

// pollInterval is the gap between GetPipelineRun polls while waiting.
const pollInterval = 5 * time.Second

var (
	_ resource.Resource                = &pipelineRunResource{}
	_ resource.ResourceWithImportState = &pipelineRunResource{}
)

func NewPipelineRunResource() resource.Resource {
	return &pipelineRunResource{}
}

type pipelineRunResource struct {
	fwadapt.ResourceBase
}

type pipelineRunResourceModel struct {
	Id                types.String `tfsdk:"id"`
	PipelineRunId     types.String `tfsdk:"pipeline_run_id"`
	WorkspaceId       types.String `tfsdk:"workspace_id"`
	PipelineId        types.String `tfsdk:"pipeline_id"`
	Manifest          types.String `tfsdk:"manifest"`
	Arguments         types.String `tfsdk:"arguments"`
	Name              types.String `tfsdk:"name"`
	NoConfirmRequired types.Bool   `tfsdk:"no_confirm_required"`
	Status            types.String `tfsdk:"status"`
	GmtCreateTime     types.String `tfsdk:"gmt_create_time"`
	GmtModifiedTime   types.String `tfsdk:"gmt_modified_time"`
}

func (r *pipelineRunResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = pipelineRunTypeName
}

func (r *pipelineRunResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAIFlow Pipeline Run resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `pipeline_run_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pipeline_run_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the pipeline run. It is the unique key of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the workspace.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"pipeline_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the pipeline to run. Required unless `manifest` is supplied as an inline pipeline definition.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"manifest": schema.StringAttribute{
				Optional:    true,
				Description: "The manifest of the pipeline run, in YAML. An inline pipeline definition used when `pipeline_id` is not supplied.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"arguments": schema.StringAttribute{
				Optional:    true,
				Description: "The run arguments, as a JSON string.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the pipeline run. Updatable via UpdatePipelineRun.",
			},
			"no_confirm_required": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to skip the run confirmation. Create-only; not read back by GetPipelineRun.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the pipeline run. Valid values: `Initialized`, `Running`, `Succeeded`, `Failed`, `Suspended`, `Terminated`, `Unknown`, `Skipped`, `Terminating`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gmt_create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the pipeline run was created.",
			},
			"gmt_modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the pipeline run was last modified.",
			},
		},
	}
}

func (r *pipelineRunResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan pipelineRunResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"WorkspaceId": plan.WorkspaceId.ValueString(),
	}
	if !plan.PipelineId.IsNull() {
		request["PipelineId"] = plan.PipelineId.ValueString()
	}
	if !plan.Manifest.IsNull() {
		request["Manifest"] = plan.Manifest.ValueString()
	}
	if !plan.Arguments.IsNull() {
		request["Arguments"] = plan.Arguments.ValueString()
	}
	if !plan.Name.IsNull() {
		request["Name"] = plan.Name.ValueString()
	}
	if !plan.NoConfirmRequired.IsNull() {
		request["NoConfirmRequired"] = plan.NoConfirmRequired.ValueBool()
	}

	response, err := r.Client().RoaPost("PAIFlow", "2021-02-02", "/api/v1/pipelineruns", nil, nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreatePipelineRun", pipelineRunTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreatePipelineRun", pipelineRunTypeName).Error(),
		)
		return
	}

	pipelineRunId, _ := response["PipelineRunId"].(string)
	if pipelineRunId == "" {
		// Some Roa responses wrap the id in a stringified form.
		if v := response["PipelineRunId"]; v != nil {
			pipelineRunId = fmt.Sprint(v)
		}
	}
	if pipelineRunId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreatePipelineRun", pipelineRunTypeName),
			fmt.Sprintf("The response carried no PipelineRunId: %v", response),
		)
		return
	}

	// A pipeline run is a one-shot execution: Create waits until it settles.
	if err := waitForPipelineRunTerminal(ctx, r.Client(), pipelineRunId, defaultCreateTimeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for pipeline run %s to settle", pipelineRunTypeName, pipelineRunId),
			err.Error(),
		)
		return
	}

	object, err := getPipelineRun(r.Client(), pipelineRunId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetPipelineRun", pipelineRunTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetPipelineRun", pipelineRunTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreatePipelineRun", pipelineRunTypeName),
			fmt.Sprintf("The pipeline run %s is not visible via GetPipelineRun right after creation", pipelineRunId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, pipelineRunState(object, &plan))...)
}

func (r *pipelineRunResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state pipelineRunResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.PipelineRunId.ValueString()
	if id == "" {
		id = state.Id.ValueString()
	}

	object, err := getPipelineRun(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetPipelineRun", pipelineRunTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetPipelineRun", pipelineRunTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, pipelineRunState(object, &state))...)
}

func (r *pipelineRunResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state pipelineRunResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// UpdatePipelineRun only updates Name; only call the API when it changed.
	id := state.PipelineRunId.ValueString()
	if id == "" {
		id = state.Id.ValueString()
	}

	if !plan.Name.Equal(state.Name) {
		action := fmt.Sprintf("/api/v1/pipelineruns/%s", id)
		body := map[string]interface{}{}
		if !plan.Name.IsNull() {
			body["Name"] = plan.Name.ValueString()
		}
		if _, err := r.Client().RoaPut("PAIFlow", "2021-02-02", action, nil, nil, body, true); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling UpdatePipelineRun", pipelineRunTypeName),
				errs.WrapErrorf(err, "Updating %s: calling UpdatePipelineRun", pipelineRunTypeName).Error(),
			)
			return
		}
	}

	// Always re-read so the Computed-only fields survive: the Update plan does
	// not carry them, so setting the plan directly would drop them.
	object, err := getPipelineRun(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetPipelineRun", pipelineRunTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetPipelineRun", pipelineRunTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, pipelineRunState(object, &plan))...)
}

func (r *pipelineRunResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state pipelineRunResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.PipelineRunId.ValueString()
	if id == "" {
		id = state.Id.ValueString()
	}
	action := fmt.Sprintf("/api/v1/pipelineruns/%s", id)

	_, err := r.Client().RoaDelete("PAIFlow", "2021-02-02", action, nil, nil, nil, true)
	if err != nil && !isPipelineRunNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeletePipelineRun", pipelineRunTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeletePipelineRun", pipelineRunTypeName).Error(),
		)
		return
	}
}

func (r *pipelineRunResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getPipelineRun queries GetPipelineRun by id and returns the response object,
// or nil when the run no longer exists.
func getPipelineRun(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	action := fmt.Sprintf("/api/v1/pipelineruns/%s", id)
	response, err := client.RoaGet("PAIFlow", "2021-02-02", action, nil, nil, nil)
	if err != nil {
		if isPipelineRunNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return response, nil
}

// waitForPipelineRunTerminal polls GetPipelineRun until the run reaches a
// terminal status. A pipeline run is an execution record: the managed state
// is the run's existence, so reaching any terminal (Succeeded, Failed,
// Terminated, Skipped) completes Create — the outcome is reflected in the
// Computed status attribute, not in apply success. Only a timeout or the run
// disappearing mid-flight is an error.
func waitForPipelineRunTerminal(ctx context.Context, client *connectivity.AliyunClient, id string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		object, err := getPipelineRun(client, id)
		if err != nil {
			return errs.WrapErrorf(err, "waiting for pipeline run %s", id)
		}
		if object == nil {
			return fmt.Errorf("pipeline run %s disappeared while waiting", id)
		}
		if isPipelineRunTerminal(strVal(object["Status"])) {
			return nil
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("waiting for pipeline run %s to settle: timeout after %s (last status %q)", id, timeout, strVal(object["Status"]))
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(pollInterval):
		}
	}
}

// isPipelineRunTerminal reports whether status is a settled end state.
func isPipelineRunTerminal(status string) bool {
	switch status {
	case statusSucceeded, statusFailed, statusTerminated, statusSkipped:
		return true
	}
	return false
}

// pipelineRunState maps a GetPipelineRun response onto the resource model.
// Attributes whose echoed value is not round-trip-equal to the configured value
// are carried over from the prior/configured state rather than the Get response,
// so the resource stays idempotent and Framework's apply-consistency check holds.
//   - manifest / arguments / no_confirm_required: not echoed by GetPipelineRun.
//   - name: echoed, but PAIFlow appends a uniqueness suffix to the pipeline run
//     name on create, so the returned Name differs from the configured value.
//
// pipeline_run_id is the authoritative identity.
func pipelineRunState(object map[string]interface{}, prior *pipelineRunResourceModel) pipelineRunResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	id := str("PipelineRunId")
	m := pipelineRunResourceModel{
		Id:              id,
		PipelineRunId:   id,
		WorkspaceId:     str("WorkspaceId"),
		PipelineId:      str("PipelineId"),
		Status:          str("Status"),
		GmtCreateTime:   str("GmtCreateTime"),
		GmtModifiedTime: str("GmtModifiedTime"),
	}
	if prior != nil {
		m.Manifest = prior.Manifest
		m.Arguments = prior.Arguments
		m.NoConfirmRequired = prior.NoConfirmRequired
		m.Name = prior.Name
	}
	return m
}

// strVal returns a string from a loosely-typed response value.
func strVal(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}

// isPipelineRunNotExist reports whether err is the pipeline-run-not-found
// error. The code arrives on a *tea.SDKError from the Roa call; PAIFlow uses
// the numeric code 200200004 for not-found (same as pipelines), and the
// substring fallback covers wrapped forms and HTTP 404.
func isPipelineRunNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil {
		switch *e.Code {
		case "200200004", "404":
			return true
		}
	}
	msg := err.Error()
	return strings.Contains(msg, "200200004") ||
		strings.Contains(msg, "not found") ||
		strings.Contains(msg, "NotExists") ||
		strings.Contains(msg, "does not exist")
}
