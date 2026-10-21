package pai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const promptTypeName = "alicloud_pai_workspace_prompt"

// The legacy SDKv2 resource declared a configurable timeouts block; plugin
// Framework offers no equivalent, so the documented defaults are fixed here.
const (
	promptCreateTimeout = 5 * time.Minute
	promptUpdateTimeout = 5 * time.Minute
	promptDeleteTimeout = 5 * time.Minute
	promptGetTimeout    = 1 * time.Minute
)

var (
	_ resource.Resource                = &promptResource{}
	_ resource.ResourceWithImportState = &promptResource{}
)

func NewPromptResource() resource.Resource {
	return &promptResource{}
}

type promptResource struct {
	fwadapt.ResourceBase
}

type promptResourceModel struct {
	Id               types.String `tfsdk:"id"`
	Accessibility    types.String `tfsdk:"accessibility"`
	CreateTime       types.String `tfsdk:"create_time"`
	Description      types.String `tfsdk:"description"`
	FrameworkContent types.String `tfsdk:"framework_content"`
	FrameworkType    types.String `tfsdk:"framework_type"`
	ModifyTime       types.String `tfsdk:"modify_time"`
	PromptId         types.String `tfsdk:"prompt_id"`
	PromptName       types.String `tfsdk:"prompt_name"`
	WorkspaceId      types.String `tfsdk:"workspace_id"`
}

func (r *promptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = promptTypeName
}

func (r *promptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI Workspace Prompt resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource supplied above. The value is formulated as `<workspace_id>:<prompt_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"accessibility": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("PRIVATE", "PUBLIC"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Represents the visibility property of the resource in the workspace.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Resource attribute field representing creation time.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Attributes that represent the context of the thread on the resource. Removing this argument keeps the current value because the API does not support clearing it.",
			},
			"framework_content": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Cue word frame content. Removing this argument keeps the current value because the API does not support clearing it.",
			},
			"framework_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Cue word frame type. Removing this argument keeps the current value because the API does not support clearing it.",
			},
			"modify_time": schema.StringAttribute{
				Computed:    true,
				Description: "Resource attribute representing modification time.",
			},
			"prompt_id": schema.StringAttribute{
				Computed:    true,
				Description: "Prompt word id.",
			},
			"prompt_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Resource attribute field representing the resource name.",
			},
			"workspace_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Workspace id.",
			},
		},
	}
}

func (r *promptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan promptResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"WorkspaceId": plan.WorkspaceId.ValueString(),
		"PromptName":  plan.PromptName.ValueString(),
	}
	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		request["Description"] = plan.Description.ValueString()
	}
	if !plan.Accessibility.IsNull() && plan.Accessibility.ValueString() != "" {
		request["Accessibility"] = plan.Accessibility.ValueString()
	}
	if !plan.FrameworkType.IsNull() && plan.FrameworkType.ValueString() != "" {
		request["FrameworkType"] = plan.FrameworkType.ValueString()
	}
	if !plan.FrameworkContent.IsNull() && plan.FrameworkContent.ValueString() != "" {
		request["FrameworkContent"] = plan.FrameworkContent.ValueString()
	}

	const action = "/api/v1/prompts"
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(promptCreateTimeout, func() *retry.RetryError {
		response, err = r.Client().RoaPost("AIWorkSpace", "2021-02-04", action, nil, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreatePrompt", promptTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, promptTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	rawPromptId, err := jsonpath.Get("$.PromptId", response)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading PromptId from the response", promptTypeName),
			errs.WrapErrorf(err, errs.FailedGetAttributeMsg, promptTypeName, "$.PromptId", response).Error(),
		)
		return
	}
	id := fmt.Sprintf("%v:%v", plan.WorkspaceId.ValueString(), rawPromptId)

	object, err := getPrompt(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling GetPrompt", promptTypeName),
			errs.WrapErrorf(err, "Creating %s: calling GetPrompt", promptTypeName).Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newPromptState(id, object))...)
}

func (r *promptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state promptResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getPrompt(r.Client(), state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetPrompt", promptTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetPrompt", promptTypeName).Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newPromptState(state.Id.ValueString(), object))...)
}

func (r *promptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state promptResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parts := strings.Split(state.Id.ValueString(), ":")
	promptId := ""
	workspaceId := ""
	if len(parts) == 2 {
		workspaceId, promptId = parts[0], parts[1]
	}
	action := fmt.Sprintf("/api/v1/prompts/%s", promptId)

	request := map[string]interface{}{
		"WorkspaceId": workspaceId,
	}
	// UpdatePrompt ignores both JSON null and empty string values, so optional
	// fields cannot be cleared through this API. The schema marks them Computed
	// and rejects empty strings, so removing an argument keeps the current value.
	if !plan.FrameworkType.Equal(state.FrameworkType) {
		request["FrameworkType"] = plan.FrameworkType.ValueString()
	}
	if !plan.FrameworkContent.Equal(state.FrameworkContent) {
		request["FrameworkContent"] = plan.FrameworkContent.ValueString()
	}
	if !plan.Description.Equal(state.Description) {
		request["Description"] = plan.Description.ValueString()
	}

	if len(request) > 1 {
		var err error
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(promptUpdateTimeout, func() *retry.RetryError {
			_, err = r.Client().RoaPut("AIWorkSpace", "2021-02-04", action, nil, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling UpdatePrompt", promptTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
	}

	object, err := getPrompt(r.Client(), state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: reading the updated Prompt", promptTypeName),
			errs.WrapErrorf(err, "Updating %s: reading the updated Prompt", promptTypeName).Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newPromptState(state.Id.ValueString(), object))...)
}

func (r *promptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state promptResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parts := strings.Split(state.Id.ValueString(), ":")
	promptId := ""
	if len(parts) == 2 {
		promptId = parts[1]
	}
	action := fmt.Sprintf("/api/v1/prompts/%s", promptId)
	query := map[string]*string{
		"WorkspaceId": tea.String(parts[0]),
	}

	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(promptDeleteTimeout, func() *retry.RetryError {
		_, err = r.Client().RoaDelete("AIWorkSpace", "2021-02-04", action, query, nil, nil, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil && !errs.NotFoundError(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeletePrompt", promptTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
}

func (r *promptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getPrompt lives in this package because the split service package cannot
// import the root alicloud package. A PromptNoPerm error is classified as
// not-found, matching the other PAI service helpers.
func getPrompt(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid Resource Id %s. Expected parts' length 2, got %d", id, len(parts))
	}
	action := fmt.Sprintf("/api/v1/prompts/%s", parts[1])
	query := map[string]*string{
		"WorkspaceId": tea.String(parts[0]),
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(promptGetTimeout, func() *retry.RetryError {
		response, err = client.RoaGet("AIWorkSpace", "2021-02-04", action, query, nil, nil)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if isPromptNoPerm(err) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("Prompt", id), errs.NotFoundMsg, response)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// isPromptNoPerm is the local form of the root package's
// IsExpectedErrors(err, []string{"PromptNoPerm"}): the only expected error
// code of the GetPrompt API.
func isPromptNoPerm(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil && *e.Code == "PromptNoPerm" {
		return true
	}
	return strings.Contains(err.Error(), "PromptNoPerm")
}

func newPromptState(id string, object map[string]interface{}) promptResourceModel {
	workspaceId := ""
	promptId := ""
	if parts := strings.Split(id, ":"); len(parts) == 2 {
		workspaceId, promptId = parts[0], parts[1]
	}
	return promptResourceModel{
		Id:               types.StringValue(id),
		Accessibility:    objectStringValue(object["Accessibility"]),
		CreateTime:       objectStringValue(object["CreateTime"]),
		Description:      optionalStringValue(object["Description"]),
		FrameworkContent: optionalStringValue(object["FrameworkContent"]),
		FrameworkType:    optionalStringValue(object["FrameworkType"]),
		ModifyTime:       objectStringValue(object["ModifyTime"]),
		PromptId:         types.StringValue(promptId),
		PromptName:       objectStringValue(object["PromptName"]),
		WorkspaceId:      types.StringValue(workspaceId),
	}
}

func objectStringValue(v interface{}) types.String {
	if s, ok := v.(string); ok {
		return types.StringValue(s)
	}
	return types.StringNull()
}

// optionalStringValue treats an absent or empty API value as null so that
// removing an optional attribute from the configuration stays converged.
func optionalStringValue(v interface{}) types.String {
	if s, ok := v.(string); ok && s != "" {
		return types.StringValue(s)
	}
	return types.StringNull()
}
