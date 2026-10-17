package paiplugin

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const contentGenerationJobTypeName = "alicloud_pai_plugin_content_generation_job"

var (
	_ resource.Resource                = &contentGenerationJobResource{}
	_ resource.ResourceWithImportState = &contentGenerationJobResource{}
)

func NewContentGenerationJobResource() resource.Resource {
	return &contentGenerationJobResource{}
}

type contentGenerationJobResource struct {
	fwadapt.ResourceBase
}

type contentGenerationJobModel struct {
	Id                             types.String `tfsdk:"id"`
	ContentGenerationJobId         types.String `tfsdk:"content_generation_job_id"`
	ContentGenerationJobConfigId   types.String `tfsdk:"content_generation_job_config_id"`
	ContentGenerationJobConfigName types.String `tfsdk:"content_generation_job_config_name"`
	Status                         types.String `tfsdk:"status"`
	Results                        types.List   `tfsdk:"results"`
	Contents                       types.List   `tfsdk:"contents"`
	SignatureId                    types.String `tfsdk:"signature_id"`
	SignatureName                  types.String `tfsdk:"signature_name"`
	GmtCreateTime                  types.String `tfsdk:"gmt_create_time"`
	GmtModifiedTime                types.String `tfsdk:"gmt_modified_time"`
}

func (r *contentGenerationJobResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = contentGenerationJobTypeName
}

func (r *contentGenerationJobResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI Plugin Content Generation Job resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `content_generation_job_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"content_generation_job_id": schema.StringAttribute{
				Computed:    true,
				Description: "The first ID of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"content_generation_job_config_id": schema.StringAttribute{
				Required:    true,
				Description: "The config id of the content generation job, obtained from ListContentGenerationJobConfigs.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content_generation_job_config_name": schema.StringAttribute{
				Computed:    true,
				Description: "The config name of the content generation job.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the content generation job. Valid values: `Running`, `Failure`, `Success`.",
			},
			"results": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The results of the content generation job.",
			},
			"contents": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "The contents of the content generation job. Updatable via UpdateContentGenerationJob. The server echoes them back on Read, so it is also Computed.",
			},
			"signature_id": schema.StringAttribute{
				Computed:    true,
				Description: "The signature id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"signature_name": schema.StringAttribute{
				Computed:    true,
				Description: "The signature name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gmt_create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The create time.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gmt_modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "The modified time.",
			},
		},
	}
}

func (r *contentGenerationJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan contentGenerationJobModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"ContentGenerationJobConfigId": plan.ContentGenerationJobConfigId.ValueString(),
	}
	// Contents is an updatable input; send it at create so the server stores the
	// initial value and the first Read reflects what the user planned.
	if contents := strSliceFromList(plan.Contents); len(contents) > 0 {
		body["Contents"] = contents
	}

	response, err := r.Client().RoaPost(paiPluginProduct, paiPluginVersion, pathJobs, nil, nil, body, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateContentGenerationJob", contentGenerationJobTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateContentGenerationJob", contentGenerationJobTypeName).Error(),
		)
		return
	}

	data := dataObject(response)
	jobId := str(data, "ContentGenerationJobId")
	if jobId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateContentGenerationJob", contentGenerationJobTypeName),
			fmt.Sprintf("The response carried no Data.ContentGenerationJobId: %v", response),
		)
		return
	}

	// Carry forward the planned contents until the first Read populates the
	// server-side fields, so the apply result is consistent with the plan.
	state := contentGenerationJobStateFromModel(&plan)
	state.Id = types.StringValue(jobId)
	state.ContentGenerationJobId = types.StringValue(jobId)

	object, err := getContentGenerationJob(r.Client(), jobId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetContentGenerationJob", contentGenerationJobTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetContentGenerationJob", contentGenerationJobTypeName).Error(),
		)
		return
	}
	if object != nil {
		state = contentGenerationJobState(object)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *contentGenerationJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state contentGenerationJobModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getContentGenerationJob(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetContentGenerationJob", contentGenerationJobTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetContentGenerationJob", contentGenerationJobTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, contentGenerationJobState(object))...)
}

func (r *contentGenerationJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state contentGenerationJobModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"Contents": strSliceFromList(plan.Contents),
	}

	if _, err := r.Client().RoaPut(paiPluginProduct, paiPluginVersion, pathJobs+"/"+state.Id.ValueString(), nil, nil, body, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateContentGenerationJob", contentGenerationJobTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateContentGenerationJob", contentGenerationJobTypeName).Error(),
		)
		return
	}

	object, err := getContentGenerationJob(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetContentGenerationJob", contentGenerationJobTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetContentGenerationJob", contentGenerationJobTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, contentGenerationJobState(object))...)
}

func (r *contentGenerationJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// PAIPlugin exposes no DeleteContentGenerationJob: the job is an immutable
	// server-side record, so delete only removes the resource from state.
}

func (r *contentGenerationJobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getContentGenerationJob(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	response, err := roaGet(client, pathJobs+"/"+id, nil)
	if err != nil {
		if isNotFound(err, response) {
			return nil, nil
		}
		return nil, err
	}
	return dataObject(response), nil
}

func contentGenerationJobState(obj map[string]interface{}) contentGenerationJobModel {
	return contentGenerationJobModel{
		Id:                             strValue(obj, "ContentGenerationJobId"),
		ContentGenerationJobId:         strValue(obj, "ContentGenerationJobId"),
		ContentGenerationJobConfigId:   strValue(obj, "ContentGenerationJobConfigId"),
		ContentGenerationJobConfigName: strValue(obj, "ContentGenerationJobConfigName"),
		Status:                         strValue(obj, "Status"),
		Results:                        strSliceValue(obj, "Results"),
		Contents:                       strSliceValue(obj, "Contents"),
		SignatureId:                    strValue(obj, "SignatureId"),
		SignatureName:                  strValue(obj, "SignatureName"),
		GmtCreateTime:                  strValue(obj, "GmtCreateTime"),
		GmtModifiedTime:                strValue(obj, "GmtModifiedTime"),
	}
}

func contentGenerationJobStateFromModel(m *contentGenerationJobModel) contentGenerationJobModel {
	return *m
}
