package bailian

import (
	"context"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
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

const promptTemplateTypeName = "alicloud_bailian_prompt_template"

var (
	_ resource.Resource                = &promptTemplateResource{}
	_ resource.ResourceWithImportState = &promptTemplateResource{}
)

func NewPromptTemplateResource() resource.Resource {
	return &promptTemplateResource{}
}

type promptTemplateResource struct {
	fwadapt.ResourceBase
}

type promptTemplateResourceModel struct {
	Id               types.String `tfsdk:"id"`
	WorkspaceId      types.String `tfsdk:"workspace_id"`
	PromptTemplateId types.String `tfsdk:"prompt_template_id"`
	Name             types.String `tfsdk:"name"`
	Content          types.String `tfsdk:"content"`
}

func (r *promptTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = promptTemplateTypeName
}

func (r *promptTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Bailian Prompt Template resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, format as `<workspace_id>:<prompt_template_id>`.",
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
			"prompt_template_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the prompt template.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the prompt template.",
			},
			"content": schema.StringAttribute{
				Required:    true,
				Description: "The content of the prompt template.",
			},
		},
	}
}

func (r *promptTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan promptTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "CreatePromptTemplate"
	llmId := extractLlmId(plan.WorkspaceId.ValueString())
	clientToken := helper.BuildClientToken(action)
	pathName := fmt.Sprintf("/%s/promptTemplates", llmId)
	query := map[string]*string{
		"name":        strPtr(plan.Name.ValueString()),
		"content":     strPtr(plan.Content.ValueString()),
		"ClientToken": strPtr(clientToken),
	}

	response, err := r.Client().RoaPost("bailian", "2023-12-29", pathName, query, nil, nil, true)
	logs.AddDebug(action, response, query)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", promptTemplateTypeName, action),
			errs.WrapErrorf(err, "Creating %s: calling %s", promptTemplateTypeName, action).Error(),
		)
		return
	}

	promptTemplateId, _ := response["promptTemplateId"].(string)
	if promptTemplateId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", promptTemplateTypeName, action),
			fmt.Sprintf("The response carried no promptTemplateId: %v", response),
		)
		return
	}

	object, err := getPromptTemplate(r.Client(), plan.WorkspaceId.ValueString(), promptTemplateId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetPromptTemplate", promptTemplateTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetPromptTemplate", promptTemplateTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", promptTemplateTypeName, action),
			fmt.Sprintf("The prompt template %s is not visible via GetPromptTemplate right after creation", promptTemplateId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, promptTemplateState(object, plan.WorkspaceId.ValueString()))...)
}

func (r *promptTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state promptTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId, promptTemplateId, ok := parseCompositeId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s", promptTemplateTypeName),
			fmt.Sprintf("Cannot parse resource ID %q: expected format `<workspace_id>:<prompt_template_id>`", state.Id.ValueString()),
		)
		return
	}

	object, err := getPromptTemplate(r.Client(), workspaceId, promptTemplateId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetPromptTemplate", promptTemplateTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetPromptTemplate", promptTemplateTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, promptTemplateState(object, workspaceId))...)
}

func (r *promptTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state promptTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId, promptTemplateId, _ := parseCompositeId(state.Id.ValueString())
	llmId := extractLlmId(workspaceId)
	pathName := fmt.Sprintf("/%s/promptTemplates/%s", llmId, promptTemplateId)
	query := map[string]*string{
		"name":    strPtr(plan.Name.ValueString()),
		"content": strPtr(plan.Content.ValueString()),
	}

	response, err := r.Client().RoaPatchWithApiName("bailian", "2023-12-29", "", pathName, query, nil, nil, true)
	logs.AddDebug("UpdatePromptTemplate", response, query)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdatePromptTemplate", promptTemplateTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdatePromptTemplate", promptTemplateTypeName).Error(),
		)
		return
	}

	object, err := getPromptTemplate(r.Client(), workspaceId, promptTemplateId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetPromptTemplate", promptTemplateTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetPromptTemplate", promptTemplateTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, promptTemplateState(object, workspaceId))...)
}

func (r *promptTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state promptTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId, promptTemplateId, _ := parseCompositeId(state.Id.ValueString())
	llmId := extractLlmId(workspaceId)
	pathName := fmt.Sprintf("/%s/promptTemplates/%s", llmId, promptTemplateId)

	response, err := r.Client().RoaDelete("bailian", "2023-12-29", pathName, nil, nil, nil, true)
	logs.AddDebug("DeletePromptTemplate", response, nil)
	if err != nil && !isBailianNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeletePromptTemplate", promptTemplateTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeletePromptTemplate", promptTemplateTypeName).Error(),
		)
		return
	}
}

func (r *promptTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getPromptTemplate queries GetPromptTemplate and returns the response object,
// or nil when the prompt template no longer exists.
func getPromptTemplate(client *connectivity.AliyunClient, workspaceId, promptTemplateId string) (map[string]interface{}, error) {
	llmId := extractLlmId(workspaceId)
	pathName := fmt.Sprintf("/%s/promptTemplates/%s", llmId, promptTemplateId)
	response, err := client.RoaGet("bailian", "2023-12-29", pathName, nil, nil, nil)
	if err != nil {
		if isBailianNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return response, nil
}

// promptTemplateState maps a GetPromptTemplate response onto the resource model.
func promptTemplateState(object map[string]interface{}, workspaceId string) promptTemplateResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	promptTemplateId := str("promptTemplateId")
	return promptTemplateResourceModel{
		Id:               types.StringValue(fmt.Sprintf("%s:%s", workspaceId, promptTemplateId.ValueString())),
		WorkspaceId:      types.StringValue(workspaceId),
		PromptTemplateId: promptTemplateId,
		Name:             str("name"),
		Content:          str("content"),
	}
}

// extractLlmId extracts the workspace identifier (llm-xxx) from the api host.
// The bailian ROA API uses the first segment of the workspace apiHost as the
// workspaceId path parameter, e.g. "llm-abc123" from "llm-abc123.cn-beijing.maas.aliyuncs.com".
func extractLlmId(apiHost string) string {
	if i := strings.Index(apiHost, "."); i > 0 {
		return apiHost[:i]
	}
	return apiHost
}

// strPtr returns a pointer to the given string, for use in map[string]*string query params.
func strPtr(s string) *string { return &s }

// isBailianNotExist reports whether err is a Bailian resource-not-found error.
func isBailianNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil {
		code := *e.Code
		if strings.Contains(code, "NotFound") || strings.Contains(code, "NotExist") {
			return true
		}
	}
	msg := err.Error()
	return strings.Contains(msg, "NotFound") || strings.Contains(msg, "NotExist")
}

// parseCompositeId splits a "<workspaceId>:<resourceId>" ID into its parts.
func parseCompositeId(id string) (string, string, bool) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) < 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}
