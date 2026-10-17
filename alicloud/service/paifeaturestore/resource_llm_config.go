package paifeaturestore

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	llmConfigTypeName = "alicloud_pai_feature_store_llm_config"
	// apiProduct / apiVersion of the PaiFeatureStore ROA OpenAPI.
	paiFeatureStoreProduct = "PaiFeatureStore"
	paiFeatureStoreVersion = "2023-06-21"
	// idSeparator mirrors the CloudSpec uniqueKeyFields delimiter between
	// InstanceId and LLMConfigId.
	idSeparator = ":"
)

var (
	_ resource.Resource                = &llmConfigResource{}
	_ resource.ResourceWithImportState = &llmConfigResource{}
)

func NewLLMConfigResource() resource.Resource {
	return &llmConfigResource{}
}

type llmConfigResource struct {
	fwadapt.ResourceBase
}

// llmConfigResourceModel is the tfsdk model. Id carries the composite
// "<InstanceId>:<LLMConfigId>" key; the API never echoes InstanceId on Read,
// so InstanceId is recovered from Id rather than from the response body.
type llmConfigResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	InstanceId         types.String `tfsdk:"instance_id"`
	WorkspaceId        types.String `tfsdk:"workspace_id"`
	Name               types.String `tfsdk:"name"`
	ApiKey             types.String `tfsdk:"api_key"`
	Model              types.String `tfsdk:"model"`
	BaseUrl            types.String `tfsdk:"base_url"`
	Rps                types.Int64  `tfsdk:"rps"`
	MaxTokens          types.Int64  `tfsdk:"max_tokens"`
	BatchSize          types.Int64  `tfsdk:"batch_size"`
	EmbeddingDimension types.Int64  `tfsdk:"embedding_dimension"`
	EnableFusion       types.Bool   `tfsdk:"enable_fusion"`
	ModelType          types.String `tfsdk:"model_type"`
	LLMConfigId        types.String `tfsdk:"llm_config_id"`
	GmtCreateTime      types.String `tfsdk:"gmt_create_time"`
	GmtModifiedTime    types.String `tfsdk:"gmt_modified_time"`
}

func (r *llmConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = llmConfigTypeName
}

func (r *llmConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PaiFeatureStore LLM Config resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, formatted as `<InstanceId>:<LLMConfigId>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the PaiFeatureStore instance.可通过接口 ListInstances 获取.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"workspace_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the workspace.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the LLM config.",
			},
			"api_key": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The API key used to call the LLM.",
			},
			"model": schema.StringAttribute{
				Required:    true,
				Description: "The model name, e.g. `text-embedding-v1`.",
			},
			"base_url": schema.StringAttribute{
				Required:    true,
				Description: "The base URL used to call the LLM.",
			},
			"rps": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The number of concurrent requests per second.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"max_tokens": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The maximum number of input tokens per row.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"batch_size": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The batch size.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"embedding_dimension": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The vector dimension. Omit to let the model use its own default dimension; the API treats 0 as the model default and echoes the resolved value, so 0 is rejected to avoid a perpetual diff.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"enable_fusion": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable fusion.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"model_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The model type. Valid values: `TEXT_EMBEDDING`, `MULTI_MODAL_EMBEDDING`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("TEXT_EMBEDDING", "MULTI_MODAL_EMBEDDING"),
				},
			},
			"llm_config_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the LLM config.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				// No UseStateForUnknown: UpdateLLMConfig bumps the timestamp,
				// so a planned (state) value can never survive an apply.
			},
		},
	}
}

func (r *llmConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan llmConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := plan.InstanceId.ValueString()
	request := map[string]interface{}{
		"ApiKey":      plan.ApiKey.ValueString(),
		"Model":       plan.Model.ValueString(),
		"Name":        plan.Name.ValueString(),
		"WorkspaceId": plan.WorkspaceId.ValueString(),
	}
	// Optional body fields: only send what the user actually configured. The
	// Optional+Computed attributes are null when omitted, so the null check
	// keeps them out of the create payload and lets the API pick defaults.
	if !plan.BaseUrl.IsNull() {
		request["BaseUrl"] = plan.BaseUrl.ValueString()
	}
	if !plan.Rps.IsNull() {
		request["Rps"] = plan.Rps.ValueInt64()
	}
	if !plan.MaxTokens.IsNull() {
		request["MaxTokens"] = plan.MaxTokens.ValueInt64()
	}
	if !plan.BatchSize.IsNull() {
		request["BatchSize"] = plan.BatchSize.ValueInt64()
	}
	if !plan.EmbeddingDimension.IsNull() {
		request["EmbeddingDimension"] = plan.EmbeddingDimension.ValueInt64()
	}
	if !plan.EnableFusion.IsNull() {
		request["EnableFusion"] = plan.EnableFusion.ValueBool()
	}
	if !plan.ModelType.IsNull() {
		request["ModelType"] = plan.ModelType.ValueString()
	}

	action := fmt.Sprintf("/api/v1/instances/%s/llmconfigs", instanceId)
	response, err := r.Client().RoaPost(paiFeatureStoreProduct, paiFeatureStoreVersion, action, nil, nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateLLMConfig", llmConfigTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateLLMConfig", llmConfigTypeName).Error(),
		)
		return
	}

	llmConfigId, _ := response["LLMConfigId"].(string)
	if llmConfigId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateLLMConfig", llmConfigTypeName),
			fmt.Sprintf("The response carried no LLMConfigId: %v", response),
		)
		return
	}

	object, err := getLLMConfig(r.Client(), instanceId, llmConfigId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetLLMConfig", llmConfigTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetLLMConfig", llmConfigTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateLLMConfig", llmConfigTypeName),
			fmt.Sprintf("The LLM config %s/%s is not visible via GetLLMConfig right after creation", instanceId, llmConfigId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, llmConfigStateWithApiKey(object, instanceId, plan.ApiKey))...)
}

func (r *llmConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state llmConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, llmConfigId, ok := parseLLMConfigID(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing resource id", llmConfigTypeName),
			fmt.Sprintf("The id %q is not in the expected <InstanceId>%s<LLMConfigId> format", state.Id.ValueString(), idSeparator),
		)
		return
	}

	object, err := getLLMConfig(r.Client(), instanceId, llmConfigId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetLLMConfig", llmConfigTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetLLMConfig", llmConfigTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, llmConfigStateWithApiKey(object, instanceId, state.ApiKey))...)
}

func (r *llmConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state llmConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := state.InstanceId.ValueString()
	if instanceId == "" {
		instanceId = parseLLMConfigIDLeft(state.Id.ValueString())
	}
	llmConfigId := state.LLMConfigId.ValueString()
	if llmConfigId == "" {
		_, llmConfigId, _ = parseLLMConfigID(state.Id.ValueString())
	}

	// UpdateLLMConfig is a full PUT: ApiKey, BaseUrl, MaxTokens, Model, Name and
	// Rps are required by the API, so they are always sent. Optional fields are
	// forwarded when the user has configured them; the Optional+Computed plan
	// value falls back to the prior state when unset, so nothing is cleared by
	// accident.
	request := map[string]interface{}{
		"ApiKey":    plan.ApiKey.ValueString(),
		"Model":     plan.Model.ValueString(),
		"Name":      plan.Name.ValueString(),
		"BaseUrl":   strOr(plan.BaseUrl, state.BaseUrl),
		"MaxTokens": int64Or(plan.MaxTokens, state.MaxTokens),
		"Rps":       int64Or(plan.Rps, state.Rps),
	}
	if !plan.BatchSize.IsNull() {
		request["BatchSize"] = plan.BatchSize.ValueInt64()
	} else if !state.BatchSize.IsNull() {
		request["BatchSize"] = state.BatchSize.ValueInt64()
	}
	if !plan.EmbeddingDimension.IsNull() {
		request["EmbeddingDimension"] = plan.EmbeddingDimension.ValueInt64()
	} else if !state.EmbeddingDimension.IsNull() {
		request["EmbeddingDimension"] = state.EmbeddingDimension.ValueInt64()
	}
	if !plan.EnableFusion.IsNull() {
		request["EnableFusion"] = plan.EnableFusion.ValueBool()
	} else if !state.EnableFusion.IsNull() {
		request["EnableFusion"] = state.EnableFusion.ValueBool()
	}
	if !plan.ModelType.IsNull() {
		request["ModelType"] = plan.ModelType.ValueString()
	} else if !state.ModelType.IsNull() {
		request["ModelType"] = state.ModelType.ValueString()
	}

	action := fmt.Sprintf("/api/v1/instances/%s/llmconfigs/%s", instanceId, llmConfigId)
	if _, err := r.Client().RoaPut(paiFeatureStoreProduct, paiFeatureStoreVersion, action, nil, nil, request, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateLLMConfig", llmConfigTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateLLMConfig", llmConfigTypeName).Error(),
		)
		return
	}

	object, err := getLLMConfig(r.Client(), instanceId, llmConfigId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetLLMConfig", llmConfigTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetLLMConfig", llmConfigTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, llmConfigStateWithApiKey(object, instanceId, plan.ApiKey))...)
}

func (r *llmConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state llmConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := state.InstanceId.ValueString()
	if instanceId == "" {
		instanceId = parseLLMConfigIDLeft(state.Id.ValueString())
	}
	llmConfigId := state.LLMConfigId.ValueString()
	if llmConfigId == "" {
		_, llmConfigId, _ = parseLLMConfigID(state.Id.ValueString())
	}

	action := fmt.Sprintf("/api/v1/instances/%s/llmconfigs/%s", instanceId, llmConfigId)
	_, err := r.Client().RoaDelete(paiFeatureStoreProduct, paiFeatureStoreVersion, action, nil, nil, nil, true)
	if err != nil && !isLLMConfigNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteLLMConfig", llmConfigTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteLLMConfig", llmConfigTypeName).Error(),
		)
		return
	}
}

func (r *llmConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The composite id "<InstanceId>:<LLMConfigId>" is stored on the id
	// attribute; Read splits it back into its parts.
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getLLMConfig queries GetLLMConfig and returns the flat response object, or
// nil when the LLM config no longer exists. The PaiFeatureStore ROA gateway
// surfaces a 404 that errs.NotFoundError recognises on a *tea.SDKError.
func getLLMConfig(client *connectivity.AliyunClient, instanceId, llmConfigId string) (map[string]interface{}, error) {
	action := fmt.Sprintf("/api/v1/instances/%s/llmconfigs/%s", instanceId, llmConfigId)
	response, err := client.RoaGet(paiFeatureStoreProduct, paiFeatureStoreVersion, action, nil, nil, nil)
	if err != nil {
		if isLLMConfigNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return response, nil
}

// llmConfigState maps a GetLLMConfig response onto the resource model. The
// response is flat (no wrapper). Empty strings become null; InstanceId comes
// from the composite id because the API never echoes it.
func llmConfigState(object map[string]interface{}, instanceId string) llmConfigResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	i64 := func(key string) types.Int64 {
		switch v := object[key].(type) {
		case nil:
			return types.Int64Null()
		case float64:
			return types.Int64Value(int64(v))
		case int:
			return types.Int64Value(int64(v))
		case int64:
			return types.Int64Value(v)
		case json.Number:
			if i, err := v.Int64(); err == nil {
				return types.Int64Value(i)
			}
			return types.Int64Null()
		default:
			return types.Int64Null()
		}
	}
	b := func(key string) types.Bool {
		if v, ok := object[key].(bool); ok {
			return types.BoolValue(v)
		}
		return types.BoolNull()
	}
	llmConfigId := str("LLMConfigId")
	return llmConfigResourceModel{
		Id:                 types.StringValue(fmt.Sprintf("%s%s%s", instanceId, idSeparator, llmConfigId.ValueString())),
		InstanceId:         types.StringValue(instanceId),
		WorkspaceId:        str("WorkspaceId"),
		Name:               str("Name"),
		ApiKey:             str("ApiKey"),
		Model:              str("Model"),
		BaseUrl:            str("BaseUrl"),
		Rps:                i64("Rps"),
		MaxTokens:          i64("MaxTokens"),
		BatchSize:          i64("BatchSize"),
		EmbeddingDimension: i64("EmbeddingDimension"),
		EnableFusion:       b("EnableFusion"),
		ModelType:          str("ModelType"),
		LLMConfigId:        llmConfigId,
		GmtCreateTime:      str("GmtCreateTime"),
		GmtModifiedTime:    str("GmtModifiedTime"),
	}
}

// llmConfigStateWithApiKey builds the state and, when the API does not echo the
// sensitive ApiKey, keeps the configured/prior value so the state is stable
// across refresh and import verification.
func llmConfigStateWithApiKey(object map[string]interface{}, instanceId string, fallback types.String) llmConfigResourceModel {
	s := llmConfigState(object, instanceId)
	if s.ApiKey.IsNull() && !fallback.IsNull() {
		s.ApiKey = fallback
	}
	return s
}

// parseLLMConfigID splits "<InstanceId>:<LLMConfigId>" on the first separator.
func parseLLMConfigID(id string) (string, string, bool) {
	parts := strings.SplitN(id, idSeparator, 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func parseLLMConfigIDLeft(id string) string {
	parts := strings.SplitN(id, idSeparator, 2)
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

// strOr returns the plan value when non-null, else the state value. Used for
// Optional+Computed attributes that the API requires on update.
func strOr(plan, state types.String) string {
	if !plan.IsNull() {
		return plan.ValueString()
	}
	return state.ValueString()
}

func int64Or(plan, state types.Int64) int64 {
	if !plan.IsNull() {
		return plan.ValueInt64()
	}
	return state.ValueInt64()
}

// isLLMConfigNotExist reports whether err is a not-found error. errs.NotFoundError
// recognises the ROA 404 on a *tea.SDKError; the substring fallback covers
// wrapped forms, matching IsExpectedErrors' behavior.
func isLLMConfigNotExist(err error) bool {
	if err == nil {
		return false
	}
	if errs.NotFoundError(err) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "notfound") || strings.Contains(msg, "not exist")
}
