package pailangstudio

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const (
	knowledgeBaseJobTypeName = "alicloud_pailangstudio_knowledge_base_job"
	pailangStudioProductCode = "PAILangStudio"
	pailangStudioAPIVersion  = "2024-07-10"
	// idSeparator joins the three identity fields (WorkspaceId, KnowledgeBaseId,
	// KnowledgeBaseJobId) into the single Terraform id; none of them is alone
	// enough to address the resource.
	idSeparator = ":"
)

var (
	_ resource.Resource                = &knowledgeBaseJobResource{}
	_ resource.ResourceWithImportState = &knowledgeBaseJobResource{}
)

func NewKnowledgeBaseJobResource() resource.Resource {
	return &knowledgeBaseJobResource{}
}

type knowledgeBaseJobResource struct {
	fwadapt.ResourceBase
}

// knowledgeBaseJobModel mirrors the resource schema. The identity is composite
// (workspace:knowledgeBase:job); the scalar fields map one-to-one to the
// KnowledgeBaseJob API, and the three nested objects mirror the API structs.
type knowledgeBaseJobModel struct {
	Id                 types.String `tfsdk:"id"`
	WorkspaceId        types.String `tfsdk:"workspace_id"`
	KnowledgeBaseId    types.String `tfsdk:"knowledge_base_id"`
	KnowledgeBaseJobId types.String `tfsdk:"knowledge_base_job_id"`
	Description        types.String `tfsdk:"description"`
	Accessibility      types.String `tfsdk:"accessibility"`
	JobAction          types.String `tfsdk:"job_action"`
	MaxRunningTime     types.Int64  `tfsdk:"max_running_time_in_seconds"`
	ResourceId         types.String `tfsdk:"resource_id"`
	EcsSpecs           types.List   `tfsdk:"ecs_specs"`
	UserVpc            types.Object `tfsdk:"user_vpc"`
	EmbeddingConfig    types.Object `tfsdk:"embedding_config"`
	Status             types.String `tfsdk:"status"`
	GmtCreateTime      types.String `tfsdk:"gmt_create_time"`
	GmtModifiedTime    types.String `tfsdk:"gmt_modified_time"`
	GmtFinishTime      types.String `tfsdk:"gmt_finish_time"`
	Creator            types.String `tfsdk:"creator"`
	ErrorMessage       types.String `tfsdk:"error_message"`
	RegionId           types.String `tfsdk:"region_id"`
}

func (r *knowledgeBaseJobResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = knowledgeBaseJobTypeName
}

func (r *knowledgeBaseJobResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI LangStudio Knowledge Base Job resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The composite ID of the resource, formatted as `<workspace_id>:<knowledge_base_id>:<knowledge_base_job_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.StringAttribute{
				Required:    true,
				Description: "The workspace ID. Part of the resource identity; immutable after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"knowledge_base_id": schema.StringAttribute{
				Required:    true,
				Description: "The knowledge base ID. Part of the resource identity and the Create path; immutable after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"knowledge_base_job_id": schema.StringAttribute{
				Computed:    true,
				Description: "The knowledge base job ID, assigned by the server on Create.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the knowledge base job. The only field updatable via UpdateKnowledgeBaseJob.",
			},
			"accessibility": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The accessibility of the job. Valid values: `PUBLIC`, `PRIVATE`. Defaults to `PUBLIC`. Immutable after creation.",
				Validators: []validator.String{
					stringvalidator.OneOf("PUBLIC", "PRIVATE"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"job_action": schema.StringAttribute{
				Optional:    true,
				Description: "The job action. Immutable after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"max_running_time_in_seconds": schema.Int64Attribute{
				Optional:    true,
				Description: "The maximum running time of the job in seconds. Immutable after creation.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"resource_id": schema.StringAttribute{
				Optional:    true,
				Description: "The prepaid resource ID. Immutable after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ecs_specs": schema.ListNestedAttribute{
				Optional:    true,
				Description: "The ECS specifications. Immutable after creation.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The node type. Defaults to `Worker`.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"instance_type": schema.StringAttribute{
							Optional:    true,
							Description: "The instance type.",
						},
						"pod_count": schema.Int64Attribute{
							Optional:    true,
							Description: "The pod count.",
						},
					},
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"user_vpc": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "The user VPC configuration. Immutable after creation.",
				Attributes: map[string]schema.Attribute{
					"vpc_id": schema.StringAttribute{
						Optional:    true,
						Description: "The VPC ID.",
					},
					"v_switch_id": schema.StringAttribute{
						Optional:    true,
						Description: "The vSwitch ID.",
					},
					"security_group_id": schema.StringAttribute{
						Optional:    true,
						Description: "The security group ID.",
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"embedding_config": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "The embedding configuration. Immutable after creation.",
				Attributes: map[string]schema.Attribute{
					"batch_size": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: "The embedding batch size.",
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"concurrency": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: "The embedding concurrency. Valid values: 1 to 200.",
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
							int64validator.AtMost(200),
						},
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The job status. Valid values: `Running`, `Success`, `Failed`.",
			},
			"gmt_create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gmt_modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "The modification time.",
			},
			"gmt_finish_time": schema.StringAttribute{
				Computed:    true,
				Description: "The finish time.",
			},
			"creator": schema.StringAttribute{
				Computed:    true,
				Description: "The creator.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"error_message": schema.StringAttribute{
				Computed:    true,
				Description: "The job error message, populated when the job fails.",
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *knowledgeBaseJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan knowledgeBaseJobModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// KnowledgeBaseId is a path parameter on Create. The API document marks it
	// optional, but a path parameter is always required, so the schema enforces
	// it; the path below substitutes it verbatim and must not be dropped.
	body := map[string]interface{}{
		"WorkspaceId": plan.WorkspaceId.ValueString(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		body["Description"] = plan.Description.ValueString()
	}
	if !plan.Accessibility.IsNull() && !plan.Accessibility.IsUnknown() {
		body["Accessibility"] = plan.Accessibility.ValueString()
	}
	if !plan.JobAction.IsNull() && !plan.JobAction.IsUnknown() {
		body["JobAction"] = plan.JobAction.ValueString()
	}
	if !plan.MaxRunningTime.IsNull() && !plan.MaxRunningTime.IsUnknown() {
		body["MaxRunningTimeInSeconds"] = plan.MaxRunningTime.ValueInt64()
	}
	if !plan.ResourceId.IsNull() && !plan.ResourceId.IsUnknown() {
		body["ResourceId"] = plan.ResourceId.ValueString()
	}
	if ecsSpecs, diags := ecsSpecsRequest(ctx, &plan); !diags.HasError() {
		resp.Diagnostics.Append(diags...)
		if ecsSpecs != nil {
			body["EcsSpecs"] = ecsSpecs
		}
	} else {
		resp.Diagnostics.Append(diags...)
		return
	}
	if uv, diags := userVpcRequest(ctx, &plan); !diags.HasError() {
		resp.Diagnostics.Append(diags...)
		if uv != nil {
			body["UserVpc"] = uv
		}
	} else {
		resp.Diagnostics.Append(diags...)
		return
	}
	if ec, diags := embeddingConfigRequest(ctx, &plan); !diags.HasError() {
		resp.Diagnostics.Append(diags...)
		if ec != nil {
			body["EmbeddingConfig"] = ec
		}
	} else {
		resp.Diagnostics.Append(diags...)
		return
	}

	pathName := fmt.Sprintf("/api/v1/langstudio/knowledgebases/%s/knowledgebasejobs", plan.KnowledgeBaseId.ValueString())
	response, err := r.Client().RoaPost(pailangStudioProductCode, pailangStudioAPIVersion, pathName, map[string]*string{}, nil, body, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateKnowledgeBaseJob", knowledgeBaseJobTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateKnowledgeBaseJob", knowledgeBaseJobTypeName).Error(),
		)
		return
	}

	// CreateKnowledgeBaseJob returns the new job ID at the top level of the
	// response body, alongside RequestId; there is no envelope wrapper.
	jobId, _ := response["KnowledgeBaseJobId"].(string)
	if jobId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateKnowledgeBaseJob", knowledgeBaseJobTypeName),
			fmt.Sprintf("The response carried no KnowledgeBaseJobId: %v", response),
		)
		return
	}

	plan.Id = types.StringValue(compositeID(plan.WorkspaceId.ValueString(), plan.KnowledgeBaseId.ValueString(), jobId))
	plan.KnowledgeBaseJobId = types.StringValue(jobId)

	object, err := getKnowledgeBaseJob(ctx, r.Client(), plan.WorkspaceId.ValueString(), plan.KnowledgeBaseId.ValueString(), jobId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetKnowledgeBaseJob", knowledgeBaseJobTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetKnowledgeBaseJob", knowledgeBaseJobTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateKnowledgeBaseJob", knowledgeBaseJobTypeName),
			fmt.Sprintf("The job %s is not visible via GetKnowledgeBaseJob right after creation", jobId),
		)
		return
	}
	resp.Diagnostics.Append(setStateFromObject(ctx, &plan, object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *knowledgeBaseJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state knowledgeBaseJobModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId, knowledgeBaseId, jobId, ok := parseCompositeID(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing id", knowledgeBaseJobTypeName),
			fmt.Sprintf("The id %q is not in the expected `<workspace_id>:<knowledge_base_id>:<knowledge_base_job_id>` form", state.Id.ValueString()),
		)
		return
	}

	object, err := getKnowledgeBaseJob(ctx, r.Client(), workspaceId, knowledgeBaseId, jobId)
	if err != nil {
		// NotFoundError only treats HTTP 404 or an explicit NotFound message as
		// absence: 403, empty and malformed responses are surfaced, not silently
		// removed, so a transient or permission failure never drops state.
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetKnowledgeBaseJob", knowledgeBaseJobTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetKnowledgeBaseJob", knowledgeBaseJobTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	state.WorkspaceId = types.StringValue(workspaceId)
	state.KnowledgeBaseId = types.StringValue(knowledgeBaseId)
	state.KnowledgeBaseJobId = types.StringValue(jobId)
	resp.Diagnostics.Append(setStateFromObject(ctx, &state, object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *knowledgeBaseJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state knowledgeBaseJobModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, knowledgeBaseId, jobId, ok := parseCompositeID(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing id", knowledgeBaseJobTypeName),
			fmt.Sprintf("The id %q is not in the expected `<workspace_id>:<knowledge_base_id>:<knowledge_base_job_id>` form", state.Id.ValueString()),
		)
		return
	}

	// UpdateKnowledgeBaseJob only accepts Description (plus the required
	// WorkspaceId). It does not extend to VPC or other fields, so the diff is
	// constrained to Description. The PUT is skipped when Description is
	// unchanged, so an empty update never risks clearing Description server-side.
	if !plan.Description.Equal(state.Description) {
		body := map[string]interface{}{
			"WorkspaceId": plan.WorkspaceId.ValueString(),
		}
		if plan.Description.IsNull() {
			body["Description"] = ""
		} else {
			body["Description"] = plan.Description.ValueString()
		}
		pathName := fmt.Sprintf("/api/v1/langstudio/knowledgebases/%s/knowledgebasejobs/%s", knowledgeBaseId, jobId)
		if _, err := r.Client().RoaPut(pailangStudioProductCode, pailangStudioAPIVersion, pathName, map[string]*string{}, nil, body, true); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling UpdateKnowledgeBaseJob", knowledgeBaseJobTypeName),
				errs.WrapErrorf(err, "Updating %s: calling UpdateKnowledgeBaseJob", knowledgeBaseJobTypeName).Error(),
			)
			return
		}
	}

	object, err := getKnowledgeBaseJob(ctx, r.Client(), plan.WorkspaceId.ValueString(), knowledgeBaseId, jobId)
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetKnowledgeBaseJob", knowledgeBaseJobTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetKnowledgeBaseJob", knowledgeBaseJobTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(setStateFromObject(ctx, &plan, object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *knowledgeBaseJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state knowledgeBaseJobModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId, knowledgeBaseId, jobId, ok := parseCompositeID(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing id", knowledgeBaseJobTypeName),
			fmt.Sprintf("The id %q is not in the expected `<workspace_id>:<knowledge_base_id>:<knowledge_base_job_id>` form", state.Id.ValueString()),
		)
		return
	}

	pathName := fmt.Sprintf("/api/v1/langstudio/knowledgebases/%s/knowledgebasejobs/%s", knowledgeBaseId, jobId)
	query := map[string]*string{
		"WorkspaceId": tea.String(workspaceId),
	}
	_, err := r.Client().RoaDelete(pailangStudioProductCode, pailangStudioAPIVersion, pathName, query, nil, nil, true)
	if err != nil && !errs.NotFoundError(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteKnowledgeBaseJob", knowledgeBaseJobTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteKnowledgeBaseJob", knowledgeBaseJobTypeName).Error(),
		)
		return
	}
}

func (r *knowledgeBaseJobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	workspaceId, knowledgeBaseId, jobId, ok := parseCompositeID(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Importing %s: parsing id", knowledgeBaseJobTypeName),
			fmt.Sprintf("The id %q is not in the expected `<workspace_id>:<knowledge_base_id>:<knowledge_base_job_id>` form", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(compositeID(workspaceId, knowledgeBaseId, jobId)))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace_id"), types.StringValue(workspaceId))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("knowledge_base_id"), types.StringValue(knowledgeBaseId))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("knowledge_base_job_id"), types.StringValue(jobId))...)
}

// getKnowledgeBaseJob queries GetKnowledgeBaseJob by the three identity fields
// and returns the flat job object from the response body, or nil when the job
// no longer exists. The response carries the KnowledgeBaseJob fields at the top
// level of the body; there is no envelope to unwrap.
func getKnowledgeBaseJob(ctx context.Context, client *connectivity.AliyunClient, workspaceId, knowledgeBaseId, jobId string) (map[string]interface{}, error) {
	pathName := fmt.Sprintf("/api/v1/langstudio/knowledgebases/%s/knowledgebasejobs/%s", knowledgeBaseId, jobId)
	query := map[string]*string{
		"WorkspaceId": tea.String(workspaceId),
	}
	response, err := client.RoaGet(pailangStudioProductCode, pailangStudioAPIVersion, pathName, query, nil, nil)
	if err != nil {
		if errs.NotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	if response == nil {
		return nil, nil
	}
	return response, nil
}

func compositeID(workspaceId, knowledgeBaseId, jobId string) string {
	return workspaceId + idSeparator + knowledgeBaseId + idSeparator + jobId
}

func parseCompositeID(id string) (workspaceId, knowledgeBaseId, jobId string, ok bool) {
	parts := strings.SplitN(id, idSeparator, 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

// setStateFromObject maps the flat GetKnowledgeBaseJob response onto the model.
// The response field names are the PascalCase API names; empty strings become
// null so absent optional and computed fields read back as absent, not "".
func setStateFromObject(ctx context.Context, m *knowledgeBaseJobModel, obj map[string]interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// Identity (WorkspaceId/KnowledgeBaseId/KnowledgeBaseJobId/Id) is set
	// authoritatively by the callers from the composite id / create response.
	// The Get response may not echo path/query parameters back, so only adopt a
	// response value when it is present; never overwrite the caller-set identity
	// with null, and do not recompute the composite Id from response fields.
	if v := strVal(obj["WorkspaceId"]); !v.IsNull() {
		m.WorkspaceId = v
	}
	if v := strVal(obj["KnowledgeBaseId"]); !v.IsNull() {
		m.KnowledgeBaseId = v
	}
	if v := strVal(obj["KnowledgeBaseJobId"]); !v.IsNull() {
		m.KnowledgeBaseJobId = v
	}
	m.Description = strVal(obj["Description"])
	m.Accessibility = strVal(obj["Accessibility"])
	m.JobAction = strVal(obj["JobAction"])
	m.ResourceId = strVal(obj["ResourceId"])
	m.Status = strVal(obj["Status"])
	m.GmtCreateTime = strVal(obj["GmtCreateTime"])
	m.GmtModifiedTime = strVal(obj["GmtModifiedTime"])
	m.GmtFinishTime = strVal(obj["GmtFinishTime"])
	m.Creator = strVal(obj["Creator"])
	m.ErrorMessage = strVal(obj["ErrorMessage"])
	m.RegionId = strVal(obj["RegionId"])
	m.MaxRunningTime = intVal(obj["MaxRunningTimeInSeconds"])

	ecsSpecs, d := ecsSpecsState(ctx, obj["EcsSpecs"])
	diags.Append(d...)
	m.EcsSpecs = ecsSpecs

	uv, d := userVpcState(ctx, obj["UserVpc"])
	diags.Append(d...)
	m.UserVpc = uv

	ec, d := embeddingConfigState(ctx, obj["EmbeddingConfig"])
	diags.Append(d...)
	m.EmbeddingConfig = ec
	return diags
}

// --- nested request builders: model -> API body ---

func ecsSpecsRequest(ctx context.Context, m *knowledgeBaseJobModel) ([]interface{}, diag.Diagnostics) {
	if m.EcsSpecs.IsNull() || m.EcsSpecs.IsUnknown() {
		return nil, nil
	}
	var specs []ecsSpecModel
	diags := m.EcsSpecs.ElementsAs(ctx, &specs, false)
	if diags.HasError() {
		return nil, diags
	}
	out := make([]interface{}, 0, len(specs))
	for _, s := range specs {
		entry := map[string]interface{}{}
		// type is Optional+Computed (server defaults to Worker); when omitted the
		// plan value is unknown, so skip sending rather than emit an empty string.
		if !s.Type.IsNull() && !s.Type.IsUnknown() {
			entry["Type"] = s.Type.ValueString()
		}
		if !s.InstanceType.IsNull() {
			entry["InstanceType"] = s.InstanceType.ValueString()
		}
		if !s.PodCount.IsNull() {
			entry["PodCount"] = s.PodCount.ValueInt64()
		}
		out = append(out, entry)
	}
	return out, diags
}

func userVpcRequest(ctx context.Context, m *knowledgeBaseJobModel) (map[string]interface{}, diag.Diagnostics) {
	if m.UserVpc.IsNull() || m.UserVpc.IsUnknown() {
		return nil, nil
	}
	var uv userVpcModel
	diags := m.UserVpc.As(ctx, &uv, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}
	out := map[string]interface{}{}
	if !uv.VpcId.IsNull() {
		out["VpcId"] = uv.VpcId.ValueString()
	}
	if !uv.VSwitchId.IsNull() {
		out["VSwitchId"] = uv.VSwitchId.ValueString()
	}
	if !uv.SecurityGroupId.IsNull() {
		out["SecurityGroupId"] = uv.SecurityGroupId.ValueString()
	}
	return out, diags
}

func embeddingConfigRequest(ctx context.Context, m *knowledgeBaseJobModel) (map[string]interface{}, diag.Diagnostics) {
	if m.EmbeddingConfig.IsNull() || m.EmbeddingConfig.IsUnknown() {
		return nil, nil
	}
	var ec embeddingConfigModel
	diags := m.EmbeddingConfig.As(ctx, &ec, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}
	out := map[string]interface{}{}
	// batch_size and concurrency are Optional+Computed (server defaults apply);
	// when omitted the plan value is unknown, so skip sending rather than emit 0
	// (a Concurrency of 0 would also violate the 1-200 validator server-side).
	if !ec.BatchSize.IsNull() && !ec.BatchSize.IsUnknown() {
		out["BatchSize"] = ec.BatchSize.ValueInt64()
	}
	if !ec.Concurrency.IsNull() && !ec.Concurrency.IsUnknown() {
		out["Concurrency"] = ec.Concurrency.ValueInt64()
	}
	return out, diags
}

// --- nested state mappers: API body -> model ---

func ecsSpecsState(ctx context.Context, raw interface{}) (types.List, diag.Diagnostics) {
	if raw == nil {
		return types.ListNull(ecsSpecObjectType), nil
	}
	arr, ok := raw.([]interface{})
	if !ok || len(arr) == 0 {
		return types.ListNull(ecsSpecObjectType), nil
	}
	var diags diag.Diagnostics
	values := make([]attr.Value, 0, len(arr))
	for _, e := range arr {
		em, _ := e.(map[string]interface{})
		obj, d := types.ObjectValue(ecsSpecAttrTypes, map[string]attr.Value{
			"type":          strVal(em["Type"]),
			"instance_type": strVal(em["InstanceType"]),
			"pod_count":     intVal(em["PodCount"]),
		})
		diags.Append(d...)
		values = append(values, obj)
	}
	list, d := types.ListValue(ecsSpecObjectType, values)
	diags.Append(d...)
	return list, diags
}

func userVpcState(ctx context.Context, raw interface{}) (types.Object, diag.Diagnostics) {
	if raw == nil {
		return types.ObjectNull(userVpcAttrTypes), nil
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return types.ObjectNull(userVpcAttrTypes), nil
	}
	return types.ObjectValue(userVpcAttrTypes, map[string]attr.Value{
		"vpc_id":            strVal(m["VpcId"]),
		"v_switch_id":       strVal(m["VSwitchId"]),
		"security_group_id": strVal(m["SecurityGroupId"]),
	})
}

func embeddingConfigState(ctx context.Context, raw interface{}) (types.Object, diag.Diagnostics) {
	if raw == nil {
		return types.ObjectNull(embeddingConfigAttrTypes), nil
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return types.ObjectNull(embeddingConfigAttrTypes), nil
	}
	return types.ObjectValue(embeddingConfigAttrTypes, map[string]attr.Value{
		"batch_size":  intVal(m["BatchSize"]),
		"concurrency": intVal(m["Concurrency"]),
	})
}

// --- nested model types and attribute-type maps ---

type ecsSpecModel struct {
	Type         types.String `tfsdk:"type"`
	InstanceType types.String `tfsdk:"instance_type"`
	PodCount     types.Int64  `tfsdk:"pod_count"`
}

type userVpcModel struct {
	VpcId           types.String `tfsdk:"vpc_id"`
	VSwitchId       types.String `tfsdk:"v_switch_id"`
	SecurityGroupId types.String `tfsdk:"security_group_id"`
}

type embeddingConfigModel struct {
	BatchSize   types.Int64 `tfsdk:"batch_size"`
	Concurrency types.Int64 `tfsdk:"concurrency"`
}

var ecsSpecAttrTypes = map[string]attr.Type{
	"type":          types.StringType,
	"instance_type": types.StringType,
	"pod_count":     types.Int64Type,
}

var ecsSpecObjectType = types.ObjectType{AttrTypes: ecsSpecAttrTypes}

var userVpcAttrTypes = map[string]attr.Type{
	"vpc_id":            types.StringType,
	"v_switch_id":       types.StringType,
	"security_group_id": types.StringType,
}

var embeddingConfigAttrTypes = map[string]attr.Type{
	"batch_size":  types.Int64Type,
	"concurrency": types.Int64Type,
}

// --- scalar value helpers ---

func strVal(v interface{}) types.String {
	if s, ok := v.(string); ok && s != "" {
		return types.StringValue(s)
	}
	return types.StringNull()
}

func intVal(v interface{}) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	switch n := v.(type) {
	case int64:
		return types.Int64Value(n)
	case int:
		return types.Int64Value(int64(n))
	case float64:
		return types.Int64Value(int64(n))
	case json.Number:
		if i, err := n.Int64(); err == nil {
			return types.Int64Value(i)
		}
	case string:
		if i, err := strconv.ParseInt(n, 10, 64); err == nil {
			return types.Int64Value(i)
		}
	}
	return types.Int64Null()
}
