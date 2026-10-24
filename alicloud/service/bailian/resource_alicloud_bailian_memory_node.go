package bailian

import (
	"context"
	"fmt"
	"strings"

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

const memoryNodeTypeName = "alicloud_bailian_memory_node"

var (
	_ resource.Resource                = &memoryNodeResource{}
	_ resource.ResourceWithImportState = &memoryNodeResource{}
)

func NewMemoryNodeResource() resource.Resource {
	return &memoryNodeResource{}
}

type memoryNodeResource struct {
	fwadapt.ResourceBase
}

type memoryNodeResourceModel struct {
	Id           types.String `tfsdk:"id"`
	WorkspaceId  types.String `tfsdk:"workspace_id"`
	MemoryId     types.String `tfsdk:"memory_id"`
	MemoryNodeId types.String `tfsdk:"memory_node_id"`
	Content      types.String `tfsdk:"content"`
}

func (r *memoryNodeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = memoryNodeTypeName
}

func (r *memoryNodeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Bailian Memory Node resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, format as `<workspace_id>:<memory_id>:<memory_node_id>`.",
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
				Required:    true,
				Description: "The memory ID the node belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"memory_node_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the memory node.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"content": schema.StringAttribute{
				Required:    true,
				Description: "The content of the memory node.",
			},
		},
	}
}

func (r *memoryNodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan memoryNodeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "CreateMemoryNode"
	llmId := extractLlmId(plan.WorkspaceId.ValueString())
	clientToken := helper.BuildClientToken(action)
	pathName := fmt.Sprintf("/%s/memories/%s/memoryNodes", llmId, plan.MemoryId.ValueString())
	query := map[string]*string{
		"content":     strPtr(plan.Content.ValueString()),
		"ClientToken": strPtr(clientToken),
	}

	response, err := r.Client().RoaPost("bailian", "2023-12-29", pathName, query, nil, nil, true)
	logs.AddDebug(action, response, query)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", memoryNodeTypeName, action),
			errs.WrapErrorf(err, "Creating %s: calling %s", memoryNodeTypeName, action).Error(),
		)
		return
	}

	memoryNodeId, _ := response["memoryNodeId"].(string)
	if memoryNodeId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", memoryNodeTypeName, action),
			fmt.Sprintf("The response carried no memoryNodeId: %v", response),
		)
		return
	}

	object, err := getMemoryNode(r.Client(), plan.WorkspaceId.ValueString(), plan.MemoryId.ValueString(), memoryNodeId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetMemoryNode", memoryNodeTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetMemoryNode", memoryNodeTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", memoryNodeTypeName, action),
			fmt.Sprintf("The memory node %s is not visible via GetMemoryNode right after creation", memoryNodeId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, memoryNodeState(object, plan.WorkspaceId.ValueString(), plan.MemoryId.ValueString()))...)
}

func (r *memoryNodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state memoryNodeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId, memoryId, memoryNodeId, ok := parseTripleId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s", memoryNodeTypeName),
			fmt.Sprintf("Cannot parse resource ID %q: expected format `<workspace_id>:<memory_id>:<memory_node_id>`", state.Id.ValueString()),
		)
		return
	}

	object, err := getMemoryNode(r.Client(), workspaceId, memoryId, memoryNodeId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetMemoryNode", memoryNodeTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetMemoryNode", memoryNodeTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, memoryNodeState(object, workspaceId, memoryId))...)
}

func (r *memoryNodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state memoryNodeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId, memoryId, memoryNodeId, _ := parseTripleId(state.Id.ValueString())
	llmId := extractLlmId(workspaceId)
	pathName := fmt.Sprintf("/%s/memories/%s/memoryNodes/%s", llmId, memoryId, memoryNodeId)
	query := map[string]*string{
		"content": strPtr(plan.Content.ValueString()),
	}

	response, err := r.Client().RoaPut("bailian", "2023-12-29", pathName, query, nil, nil, true)
	logs.AddDebug("UpdateMemoryNode", response, query)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateMemoryNode", memoryNodeTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateMemoryNode", memoryNodeTypeName).Error(),
		)
		return
	}

	object, err := getMemoryNode(r.Client(), workspaceId, memoryId, memoryNodeId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetMemoryNode", memoryNodeTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetMemoryNode", memoryNodeTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, memoryNodeState(object, workspaceId, memoryId))...)
}

func (r *memoryNodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state memoryNodeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId, memoryId, memoryNodeId, _ := parseTripleId(state.Id.ValueString())
	llmId := extractLlmId(workspaceId)
	pathName := fmt.Sprintf("/%s/memories/%s/memoryNodes/%s", llmId, memoryId, memoryNodeId)

	response, err := r.Client().RoaDelete("bailian", "2023-12-29", pathName, nil, nil, nil, true)
	logs.AddDebug("DeleteMemoryNode", response, nil)
	if err != nil && !isBailianNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteMemoryNode", memoryNodeTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteMemoryNode", memoryNodeTypeName).Error(),
		)
		return
	}
}

func (r *memoryNodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getMemoryNode(client *connectivity.AliyunClient, workspaceId, memoryId, memoryNodeId string) (map[string]interface{}, error) {
	llmId := extractLlmId(workspaceId)
	pathName := fmt.Sprintf("/%s/memories/%s/memoryNodes/%s", llmId, memoryId, memoryNodeId)
	response, err := client.RoaGet("bailian", "2023-12-29", pathName, nil, nil, nil)
	if err != nil {
		if isBailianNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return response, nil
}

func memoryNodeState(object map[string]interface{}, workspaceId, memoryId string) memoryNodeResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	memoryNodeId := str("memoryNodeId")
	return memoryNodeResourceModel{
		Id:           types.StringValue(fmt.Sprintf("%s:%s:%s", workspaceId, memoryId, memoryNodeId.ValueString())),
		WorkspaceId:  types.StringValue(workspaceId),
		MemoryId:     types.StringValue(memoryId),
		MemoryNodeId: memoryNodeId,
		Content:      str("content"),
	}
}

// parseTripleId splits a "<workspaceId>:<memoryId>:<memoryNodeId>" ID into its parts.
func parseTripleId(id string) (string, string, string, bool) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) < 3 {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}
