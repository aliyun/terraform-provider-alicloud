package bailian

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
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
)

const memoryTypeName = "alicloud_bailian_memory"

var (
	_ resource.Resource                = &memoryResource{}
	_ resource.ResourceWithImportState = &memoryResource{}
)

func NewMemoryResource() resource.Resource {
	return &memoryResource{}
}

type memoryResource struct {
	fwadapt.ResourceBase
}

type memoryResourceModel struct {
	Id          types.String `tfsdk:"id"`
	WorkspaceId types.String `tfsdk:"workspace_id"`
	MemoryId    types.String `tfsdk:"memory_id"`
	Description types.String `tfsdk:"description"`
}

func (r *memoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = memoryTypeName
}

func (r *memoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Bailian Memory resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, format as `<workspace_id>:<memory_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.StringAttribute{
				Required:    true,
				Description: "The api host of the Bailian workspace (e.g. llm-xxx.cn-beijing.maas.aliyuncs.com).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"memory_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the memory.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Required:    true,
				Description: "The description of the memory.",
			},
		},
	}
}

func (r *memoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan memoryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "CreateMemory"
	llmId := extractLlmId(plan.WorkspaceId.ValueString())
	clientToken := helper.BuildClientToken(action)
	pathName := fmt.Sprintf("/%s/memories", llmId)
	query := map[string]*string{
		"description": strPtr(plan.Description.ValueString()),
		"ClientToken": strPtr(clientToken),
	}

	response, err := r.Client().RoaPost("bailian", "2023-12-29", pathName, query, nil, nil, true)
	logs.AddDebug(action, response, query)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", memoryTypeName, action),
			errs.WrapErrorf(err, "Creating %s: calling %s", memoryTypeName, action).Error(),
		)
		return
	}

	memoryId, _ := response["memoryId"].(string)
	if memoryId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", memoryTypeName, action),
			fmt.Sprintf("The response carried no memoryId: %v", response),
		)
		return
	}

	object, err := getMemory(r.Client(), plan.WorkspaceId.ValueString(), memoryId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetMemory", memoryTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetMemory", memoryTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", memoryTypeName, action),
			fmt.Sprintf("The memory %s is not visible via GetMemory right after creation", memoryId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, memoryState(object, plan.WorkspaceId.ValueString()))...)
}

func (r *memoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state memoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId, memoryId, ok := parseCompositeId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s", memoryTypeName),
			fmt.Sprintf("Cannot parse resource ID %q: expected format `<workspace_id>:<memory_id>`", state.Id.ValueString()),
		)
		return
	}

	object, err := getMemory(r.Client(), workspaceId, memoryId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetMemory", memoryTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetMemory", memoryTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, memoryState(object, workspaceId))...)
}

func (r *memoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state memoryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId, memoryId, _ := parseCompositeId(state.Id.ValueString())
	llmId := extractLlmId(workspaceId)
	pathName := fmt.Sprintf("/%s/memories/%s", llmId, memoryId)
	query := map[string]*string{
		"description": strPtr(plan.Description.ValueString()),
	}

	response, err := r.Client().RoaPut("bailian", "2023-12-29", pathName, query, nil, nil, true)
	logs.AddDebug("UpdateMemory", response, query)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateMemory", memoryTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateMemory", memoryTypeName).Error(),
		)
		return
	}

	object, err := getMemory(r.Client(), workspaceId, memoryId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetMemory", memoryTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetMemory", memoryTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, memoryState(object, workspaceId))...)
}

func (r *memoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state memoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId, memoryId, _ := parseCompositeId(state.Id.ValueString())
	llmId := extractLlmId(workspaceId)
	pathName := fmt.Sprintf("/%s/memories/%s", llmId, memoryId)

	response, err := r.Client().RoaDelete("bailian", "2023-12-29", pathName, nil, nil, nil, true)
	logs.AddDebug("DeleteMemory", response, nil)
	if err != nil && !isBailianNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteMemory", memoryTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteMemory", memoryTypeName).Error(),
		)
		return
	}
}

func (r *memoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getMemory(client *connectivity.AliyunClient, workspaceId, memoryId string) (map[string]interface{}, error) {
	llmId := extractLlmId(workspaceId)
	pathName := fmt.Sprintf("/%s/memories/%s", llmId, memoryId)
	response, err := client.RoaGet("bailian", "2023-12-29", pathName, nil, nil, nil)
	if err != nil {
		if isBailianNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return response, nil
}

func memoryState(object map[string]interface{}, workspaceId string) memoryResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	memoryId := str("memoryId")
	return memoryResourceModel{
		Id:          types.StringValue(fmt.Sprintf("%s:%s", workspaceId, memoryId.ValueString())),
		WorkspaceId: types.StringValue(workspaceId),
		MemoryId:    memoryId,
		Description: str("description"),
	}
}
