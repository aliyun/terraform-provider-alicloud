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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const contentGenerationJobConfigTypeName = "alicloud_pai_plugin_content_generation_job_config"

var (
	_ resource.Resource                = &contentGenerationJobConfigResource{}
	_ resource.ResourceWithImportState = &contentGenerationJobConfigResource{}
)

func NewContentGenerationJobConfigResource() resource.Resource {
	return &contentGenerationJobConfigResource{}
}

type contentGenerationJobConfigResource struct {
	fwadapt.ResourceBase
}

type contentGenerationJobConfigModel struct {
	Id                           types.String `tfsdk:"id"`
	ContentGenerationJobConfigId types.String `tfsdk:"content_generation_job_config_id"`
	Name                         types.String `tfsdk:"name"`
	Purpose                      types.String `tfsdk:"purpose"`
	Type                         types.String `tfsdk:"type"`
	Model                        types.String `tfsdk:"model"`
	IndustryCode                 types.String `tfsdk:"industry_code"`
	IndustryName                 types.String `tfsdk:"industry_name"`
	SignatureId                  types.String `tfsdk:"signature_id"`
	SignatureName                types.String `tfsdk:"signature_name"`
	RecommendedKeywords          types.List   `tfsdk:"recommended_keywords"`
	CustomKeywords               types.List   `tfsdk:"custom_keywords"`
	JobCount                     types.Int64  `tfsdk:"job_count"`
	GmtCreateTime                types.String `tfsdk:"gmt_create_time"`
	GmtModifiedTime              types.String `tfsdk:"gmt_modified_time"`
}

func (r *contentGenerationJobConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = contentGenerationJobConfigTypeName
}

func (r *contentGenerationJobConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI Plugin Content Generation Job Config resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `content_generation_job_config_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"content_generation_job_config_id": schema.StringAttribute{
				Computed:    true,
				Description: "The first ID of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the config.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"purpose": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The purpose of the config. Valid values: `UserTouch`, `SmsRecall`. When omitted, the server applies its default, so the field is also Computed and populated by Read.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The content generation type. Currently only `Template` is supported.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"model": schema.StringAttribute{
				Required:    true,
				Description: "The model used by the config. Currently only `QWen` is supported.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"industry_code": schema.StringAttribute{
				Required:    true,
				Description: "The industry code.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"industry_name": schema.StringAttribute{
				Computed:    true,
				Description: "The industry name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"signature_id": schema.StringAttribute{
				Required:    true,
				Description: "The signature id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"signature_name": schema.StringAttribute{
				Computed:    true,
				Description: "The signature name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"recommended_keywords": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "The recommended keywords list.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"custom_keywords": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "The custom keywords list.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"job_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The job count.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
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

func (r *contentGenerationJobConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan contentGenerationJobConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"Name":                plan.Name.ValueString(),
		"Type":                plan.Type.ValueString(),
		"Model":               plan.Model.ValueString(),
		"IndustryCode":        plan.IndustryCode.ValueString(),
		"SignatureId":         plan.SignatureId.ValueString(),
		"RecommendedKeywords": strSliceFromList(plan.RecommendedKeywords),
	}
	if !plan.Purpose.IsNull() {
		body["Purpose"] = plan.Purpose.ValueString()
	}
	if custom := strSliceFromList(plan.CustomKeywords); len(custom) > 0 {
		body["CustomKeywords"] = custom
	}

	response, err := r.Client().RoaPost(paiPluginProduct, paiPluginVersion, pathJobConfigs, nil, nil, body, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateContentGenerationJobConfig", contentGenerationJobConfigTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateContentGenerationJobConfig", contentGenerationJobConfigTypeName).Error(),
		)
		return
	}

	data := dataObject(response)
	configId := str(data, "ContentGenerationJobConfigId")
	if configId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateContentGenerationJobConfig", contentGenerationJobConfigTypeName),
			fmt.Sprintf("The response carried no Data.ContentGenerationJobConfigId: %v", response),
		)
		return
	}

	object, err := getContentGenerationJobConfig(r.Client(), configId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetContentGenerationJobConfig", contentGenerationJobConfigTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetContentGenerationJobConfig", contentGenerationJobConfigTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateContentGenerationJobConfig", contentGenerationJobConfigTypeName),
			fmt.Sprintf("The config %s is not visible via GetContentGenerationJobConfig right after creation", configId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, contentGenerationJobConfigState(object))...)
}

func (r *contentGenerationJobConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state contentGenerationJobConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getContentGenerationJobConfig(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetContentGenerationJobConfig", contentGenerationJobConfigTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetContentGenerationJobConfig", contentGenerationJobConfigTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, contentGenerationJobConfigState(object))...)
}

func (r *contentGenerationJobConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// PAIPlugin exposes no UpdateContentGenerationJobConfig: every input change
	// RequiresReplace, so an in-place update never reaches the provider.
	var plan, state contentGenerationJobConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, contentGenerationJobConfigStateFromModel(&state))...)
}

func (r *contentGenerationJobConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state contentGenerationJobConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.Client().RoaDelete(paiPluginProduct, paiPluginVersion, pathJobConfigs+"/"+state.Id.ValueString(), nil, nil, nil, true)
	if err != nil && !isNotFound(err, nil) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteContentGenerationJobConfig", contentGenerationJobConfigTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteContentGenerationJobConfig", contentGenerationJobConfigTypeName).Error(),
		)
		return
	}
}

func (r *contentGenerationJobConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getContentGenerationJobConfig queries GetContentGenerationJobConfig by id and
// returns the Data object, or nil when the config no longer exists.
func getContentGenerationJobConfig(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	response, err := roaGet(client, pathJobConfigs+"/"+id, nil)
	if err != nil {
		if isNotFound(err, response) {
			return nil, nil
		}
		return nil, err
	}
	return dataObject(response), nil
}

// contentGenerationJobConfigState maps a Get response Data object onto the
// resource model. The config id is the identity; id mirrors it.
func contentGenerationJobConfigState(obj map[string]interface{}) contentGenerationJobConfigModel {
	m := contentGenerationJobConfigModel{
		Id:                           strValue(obj, "ContentGenerationJobConfigId"),
		ContentGenerationJobConfigId: strValue(obj, "ContentGenerationJobConfigId"),
		Name:                         strValue(obj, "Name"),
		Purpose:                      strValue(obj, "Purpose"),
		Type:                         strValue(obj, "Type"),
		Model:                        strValue(obj, "Model"),
		IndustryCode:                 strValue(obj, "IndustryCode"),
		IndustryName:                 strValue(obj, "IndustryName"),
		SignatureId:                  strValue(obj, "SignatureId"),
		SignatureName:                strValue(obj, "SignatureName"),
		RecommendedKeywords:          strSliceValue(obj, "RecommendedKeywords"),
		CustomKeywords:               strSliceValue(obj, "CustomKeywords"),
		GmtCreateTime:                strValue(obj, "GmtCreateTime"),
		GmtModifiedTime:              strValue(obj, "GmtModifiedTime"),
	}
	if jc := intValue(obj, "JobCount"); jc != 0 {
		m.JobCount = types.Int64Value(jc)
	} else {
		m.JobCount = types.Int64Null()
	}
	return m
}

// contentGenerationJobConfigStateFromModel preserves state when no API read is
// possible (the only Update path is a no-op prior to replace).
func contentGenerationJobConfigStateFromModel(m *contentGenerationJobConfigModel) contentGenerationJobConfigModel {
	return *m
}
