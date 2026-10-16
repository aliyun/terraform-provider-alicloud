package paiflow

import (
	"context"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const paiRunTypeName = "alicloud_pai_run"

// PAIFlow ROA API constants shared across the lifecycle methods.
const (
	paiFlowProduct = "PAIFlow"
	paiFlowVersion = "2021-02-02"

	paiRunBasePath  = "/api/v1/pipelineruns"
	paiRunStartPath = "/api/v1/pipelineruns/%s/start"
)

var (
	_ resource.Resource                     = &paiRunResource{}
	_ resource.ResourceWithImportState      = &paiRunResource{}
	_ resource.ResourceWithConfigValidators = &paiRunResource{}
)

type paiRunResource struct {
	fwadapt.ResourceBase
}

type paiRunResourceModel struct {
	Id               types.String `tfsdk:"id"`
	WorkspaceId      types.String `tfsdk:"workspace_id"`
	PipelineId       types.String `tfsdk:"pipeline_id"`
	PipelineManifest types.String `tfsdk:"pipeline_manifest"`
	Name             types.String `tfsdk:"name"`
	Arguments        types.String `tfsdk:"arguments"`
	Options          types.String `tfsdk:"options"`
	SourceType       types.String `tfsdk:"source_type"`
	SourceId         types.String `tfsdk:"source_id"`
	Accessibility    types.String `tfsdk:"accessibility"`
	Status           types.String `tfsdk:"status"`
	Manifest         types.String `tfsdk:"manifest"`
	PipelineRunUri   types.String `tfsdk:"pipeline_run_uri"`
	StartedAt        types.String `tfsdk:"started_at"`
	FinishedAt       types.String `tfsdk:"finished_at"`
	Duration         types.String `tfsdk:"duration"`
	NodeId           types.String `tfsdk:"node_id"`
	UserId           types.String `tfsdk:"user_id"`
	ParentUserId     types.String `tfsdk:"parent_user_id"`
	Message          types.String `tfsdk:"message"`
	GmtCreateTime    types.String `tfsdk:"gmt_create_time"`
	GmtModifiedTime  types.String `tfsdk:"gmt_modified_time"`
}

func (r *paiRunResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = paiRunTypeName
}

func (r *paiRunResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI Flow Pipeline Run resource. " +
			"A pipeline run is a single execution of a pipeline. " +
			"For information about PAI Flow Pipeline Run and how to use it, see [PAI Flow API](https://www.alibabacloud.com/help/en/pai/developer-reference/api-paiflow-2021-02-02-createpipelinerun).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `pipeline_run_id`.",
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
				Computed:    true,
				Description: "The ID of the pipeline to run. Exactly one of `pipeline_id` and `pipeline_manifest` must be set.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pipeline_manifest": schema.StringAttribute{
				Optional:    true,
				Description: "The manifest of the pipeline to run, in YAML. Exactly one of `pipeline_id` and `pipeline_manifest` must be set.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the pipeline run. If omitted, PAI Flow generates one.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"arguments": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The run arguments in YAML.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"options": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The run options in JSON.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The source type of the pipeline run. Valid values: `UNKNOWN`, `SDK`, `DESIGNER`, `M6`. Default: `UNKNOWN`.",
				Validators: []validator.String{
					stringvalidator.OneOf("UNKNOWN", "SDK", "DESIGNER", "M6"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The source ID of the pipeline run.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"accessibility": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The accessibility of the pipeline run. Valid values: `PUBLIC`, `PRIVATE`. Default: `PUBLIC`.",
				Validators: []validator.String{
					stringvalidator.OneOf("PUBLIC", "PRIVATE"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the pipeline run, e.g. `Running`, `Succeeded`, `Failed`, `Stopped`.",
			},
			"manifest": schema.StringAttribute{
				Computed:    true,
				Description: "The manifest of the pipeline run.",
			},
			"pipeline_run_uri": schema.StringAttribute{
				Computed:    true,
				Description: "The URI of the pipeline run.",
			},
			"started_at": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the pipeline run started.",
			},
			"finished_at": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the pipeline run finished.",
			},
			"duration": schema.StringAttribute{
				Computed:    true,
				Description: "The duration of the pipeline run, in milliseconds.",
			},
			"node_id": schema.StringAttribute{
				Computed:    true,
				Description: "The node ID of the pipeline run.",
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "The user ID that started the pipeline run.",
			},
			"parent_user_id": schema.StringAttribute{
				Computed:    true,
				Description: "The parent user ID of the pipeline run.",
			},
			"message": schema.StringAttribute{
				Computed:    true,
				Description: "The message of the pipeline run, populated on failure.",
			},
			"gmt_create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the pipeline run.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gmt_modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "The modification time of the pipeline run.",
			},
		},
	}
}

// ConfigValidators enforces the cross-attribute constraint at plan time, before
// any API call: exactly one of pipeline_id and pipeline_manifest must be set.
func (r *paiRunResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("pipeline_id"),
			path.MatchRoot("pipeline_manifest"),
		),
	}
}

func (r *paiRunResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan paiRunResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"WorkspaceId": plan.WorkspaceId.ValueString(),
	}
	if !plan.PipelineId.IsNull() {
		body["PipelineId"] = plan.PipelineId.ValueString()
	}
	if !plan.PipelineManifest.IsNull() {
		body["PipelineManifest"] = plan.PipelineManifest.ValueString()
	}
	if !plan.Name.IsNull() {
		body["Name"] = plan.Name.ValueString()
	}
	if !plan.Arguments.IsNull() {
		body["Arguments"] = plan.Arguments.ValueString()
	}
	if !plan.Options.IsNull() {
		body["Options"] = plan.Options.ValueString()
	}
	if !plan.SourceType.IsNull() {
		body["SourceType"] = plan.SourceType.ValueString()
	}
	if !plan.SourceId.IsNull() {
		body["SourceId"] = plan.SourceId.ValueString()
	}
	if !plan.Accessibility.IsNull() {
		body["Accessibility"] = plan.Accessibility.ValueString()
	}

	response, err := r.Client().RoaPost(paiFlowProduct, paiFlowVersion, paiRunBasePath, nil, nil, body, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreatePipelineRun", paiRunTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreatePipelineRun", paiRunTypeName).Error(),
		)
		return
	}

	pipelineRunId, _ := response["PipelineRunId"].(string)
	if pipelineRunId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreatePipelineRun", paiRunTypeName),
			fmt.Sprintf("The response carried no PipelineRunId: %v", response),
		)
		return
	}

	// Start the run. PAIFlow separates creating a run from starting it; the
	// resource models a run that executes, so issue Start after Create.
	startPath := fmt.Sprintf(paiRunStartPath, pipelineRunId)
	if _, err := r.Client().RoaPut(paiFlowProduct, paiFlowVersion, startPath, nil, nil, nil, true); err != nil {
		// A start error does not always mean the run cannot run: a freshly
		// created run may already be starting. Only fail when the run cannot
		// be observed at all afterwards.
		object, getErr := getPaiRun(r.Client(), pipelineRunId)
		if getErr != nil || object == nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: calling StartPipelineRun", paiRunTypeName),
				errs.WrapErrorf(err, "Creating %s: calling StartPipelineRun for %s", paiRunTypeName, pipelineRunId).Error(),
			)
			return
		}
	}

	object, err := getPaiRun(r.Client(), pipelineRunId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetPipelineRun", paiRunTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetPipelineRun", paiRunTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreatePipelineRun", paiRunTypeName),
			fmt.Sprintf("The pipeline run %s is not visible via GetPipelineRun right after creation", pipelineRunId),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, paiRunState(object, plan))...)
}

func (r *paiRunResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state paiRunResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getPaiRun(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetPipelineRun", paiRunTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetPipelineRun", paiRunTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, paiRunState(object, state))...)
}

func (r *paiRunResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Pipeline runs are immutable: every writable argument is ForceNew, so the
	// only plan changes that reach Update touch computed attributes. There is no
	// UpdatePipelineRun API, so just reflect the server-side state.
	var plan, state paiRunResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getPaiRun(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetPipelineRun", paiRunTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetPipelineRun", paiRunTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, paiRunState(object, plan))...)
}

func (r *paiRunResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state paiRunResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pipelineRunId := state.Id.ValueString()
	action := fmt.Sprintf("/api/v1/pipelineruns/%s", pipelineRunId)
	if _, err := r.Client().RoaDelete(paiFlowProduct, paiFlowVersion, action, nil, nil, nil, true); err != nil {
		// The run is gone: a not-found on delete is a successful delete.
		if isPipelineRunNotExist(err) || errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeletePipelineRun", paiRunTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeletePipelineRun for %s", paiRunTypeName, pipelineRunId).Error(),
		)
		return
	}
}

func (r *paiRunResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getPaiRun queries GetPipelineRun by run ID and returns the run object, or nil
// when the run no longer exists.
func getPaiRun(client *connectivity.AliyunClient, pipelineRunId string) (map[string]interface{}, error) {
	action := fmt.Sprintf("/api/v1/pipelineruns/%s", pipelineRunId)
	response, err := client.RoaGet(paiFlowProduct, paiFlowVersion, action, nil, nil, nil)
	if err != nil {
		if isPipelineRunNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return response, nil
}

// paiRunState maps a GetPipelineRun response onto the resource model. Empty
// strings become null so absent attributes read back as absent rather than "".
// pipeline_manifest is a pure input that the API does not echo back separately,
// so it is carried forward from the plan/state.
func paiRunState(object map[string]interface{}, prior paiRunResourceModel) paiRunResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}

	manifest := str("Manifest")
	// pipeline_manifest is not part of the GetPipelineRun response; preserve
	// the value the user supplied so a refresh does not drift into a recreate.
	pipelineManifest := prior.PipelineManifest

	return paiRunResourceModel{
		Id:               str("PipelineRunId"),
		WorkspaceId:      str("WorkspaceId"),
		PipelineId:       str("PipelineId"),
		PipelineManifest: pipelineManifest,
		Name:             str("Name"),
		Arguments:        str("Arguments"),
		Options:          str("Options"),
		SourceType:       str("SourceType"),
		SourceId:         str("SourceId"),
		Accessibility:    str("Accessibility"),
		Status:           str("Status"),
		Manifest:         manifest,
		PipelineRunUri:   str("PipelineRunUri"),
		StartedAt:        str("StartedAt"),
		FinishedAt:       str("FinishedAt"),
		Duration:         str("Duration"),
		NodeId:           str("NodeId"),
		UserId:           str("UserId"),
		ParentUserId:     str("ParentUserId"),
		Message:          str("Message"),
		GmtCreateTime:    str("GmtCreateTime"),
		GmtModifiedTime:  str("GmtModifiedTime"),
	}
}

// isPipelineRunNotExist reports whether err is a PAIFlow run-not-found error.
// The Get endpoint answers a missing run with code 200200004; the Delete
// endpoint answers with HTTP 404. Both arrive on a *tea.SDKError from the ROA
// client; the substring fallback covers wrapped forms.
func isPipelineRunNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok {
		if code := tea.StringValue(e.Code); code == "200200004" || code == "404" {
			return true
		}
		if tea.IntValue(e.StatusCode) == 404 {
			return true
		}
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "200200004") ||
		strings.Contains(msg, "not exist") ||
		strings.Contains(msg, "not found")
}
