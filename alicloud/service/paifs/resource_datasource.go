package paifs

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	datasourceTypeName = "alicloud_pai_feature_store_datasource"
	paifsProduct       = "PaiFeatureStore"
	paifsVersion       = "2023-06-21"
	idSeparator        = ":"
	datasourceNotExist = "DoConfigCenterRequestErrorProblem"
)

var (
	_ resource.Resource                = &datasourceResource{}
	_ resource.ResourceWithImportState = &datasourceResource{}

	// datasourceNamePattern mirrors the CloudSpec contract: start with a letter
	// or digit, only letters, digits and underscores, at most 63 characters.
	datasourceNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_]{0,62}$`)
)

func NewDatasourceResource() resource.Resource {
	return &datasourceResource{}
}

type datasourceResource struct {
	fwadapt.ResourceBase
}

type datasourceResourceModel struct {
	Id             types.String `tfsdk:"id"`
	InstanceId     types.String `tfsdk:"instance_id"`
	DatasourceId   types.String `tfsdk:"datasource_id"`
	Type           types.String `tfsdk:"type"`
	DatasourceName types.String `tfsdk:"datasource_name"`
	WorkspaceId    types.String `tfsdk:"workspace_id"`
	Uri            types.String `tfsdk:"uri"`
	Config         types.String `tfsdk:"config"`
	CreateTime     types.String `tfsdk:"create_time"`
}

func (r *datasourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = datasourceTypeName
}

func (r *datasourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI Feature Store Datasource resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, in the format `<InstanceId>:<DatasourceId>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the FeatureStore instance, which can be obtained through the ListInstances interface.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"datasource_id": schema.StringAttribute{
				Computed:    true,
				Description: "The first ID of the resource, allocated by the server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The data source type. Valid values: `Hologres`, `GraphCompute`, `FeatureDB`, `MaxCompute`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("Hologres", "GraphCompute", "FeatureDB", "MaxCompute"),
				},
			},
			"datasource_name": schema.StringAttribute{
				Required:    true,
				Description: "The datasource name. Start with a letter or number and contain only letters, numbers, and underscores. The maximum length cannot exceed 63.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					stringvalidator.RegexMatches(datasourceNamePattern,
						"must start with a letter or digit and contain only letters, digits and underscores"),
				},
			},
			"workspace_id": schema.StringAttribute{
				Required:    true,
				Description: "The PAI workspace ID, which can be obtained through the ListWorkspaces interface.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"uri": schema.StringAttribute{
				Required:    true,
				Description: "The resource address of the datasource. The maximum length cannot exceed 64.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"config": schema.StringAttribute{
				Optional:    true,
				Description: "The resource configuration information. The value is a JSON string, and the expected contents depend on `type`.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *datasourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan datasourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"InstanceId":  plan.InstanceId.ValueString(),
		"Type":        plan.Type.ValueString(),
		"Name":        plan.DatasourceName.ValueString(),
		"WorkspaceId": plan.WorkspaceId.ValueString(),
		"Uri":         plan.Uri.ValueString(),
	}
	if !plan.Config.IsNull() {
		body["Config"] = plan.Config.ValueString()
	}

	response, err := r.Client().RpcPost(paifsProduct, paifsVersion, "CreateDatasource", nil, body, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDatasource", datasourceTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateDatasource", datasourceTypeName).Error(),
		)
		return
	}

	datasourceId, _ := response["DatasourceId"].(string)
	if datasourceId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDatasource", datasourceTypeName),
			fmt.Sprintf("The response carried no DatasourceId: %v", response),
		)
		return
	}

	object, err := getDatasource(r.Client(), plan.InstanceId.ValueString(), datasourceId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetDatasource", datasourceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetDatasource", datasourceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDatasource", datasourceTypeName),
			fmt.Sprintf("The datasource %s is not visible via GetDatasource right after creation", datasourceId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, datasourceState(object, plan.InstanceId.ValueString(), datasourceId))...)
}

func (r *datasourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state datasourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, datasourceId := datasourceIdentifiers(state)
	if instanceId == "" || datasourceId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: missing identifiers", datasourceTypeName),
			fmt.Sprintf("The state carried no usable InstanceId:DatasourceId pair (id=%q)", state.Id.ValueString()),
		)
		return
	}

	object, err := getDatasource(r.Client(), instanceId, datasourceId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetDatasource", datasourceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetDatasource", datasourceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, datasourceState(object, instanceId, datasourceId))...)
}

func (r *datasourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state datasourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, datasourceId := datasourceIdentifiers(state)
	body := map[string]interface{}{
		"InstanceId":   instanceId,
		"DatasourceId": datasourceId,
		"Name":         plan.DatasourceName.ValueString(),
		"Uri":          plan.Uri.ValueString(),
	}
	if !plan.Config.IsNull() {
		body["Config"] = plan.Config.ValueString()
	}

	if _, err := r.Client().RpcPost(paifsProduct, paifsVersion, "UpdateDatasource", nil, body, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateDatasource", datasourceTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateDatasource", datasourceTypeName).Error(),
		)
		return
	}

	object, err := getDatasource(r.Client(), instanceId, datasourceId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetDatasource", datasourceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetDatasource", datasourceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, datasourceState(object, instanceId, datasourceId))...)
}

func (r *datasourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state datasourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, datasourceId := datasourceIdentifiers(state)
	if instanceId == "" || datasourceId == "" {
		return
	}

	_, err := r.Client().RpcPost(paifsProduct, paifsVersion, "DeleteDatasource", nil, map[string]interface{}{
		"InstanceId":   instanceId,
		"DatasourceId": datasourceId,
	}, true)
	if err != nil && !isDatasourceNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteDatasource", datasourceTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteDatasource", datasourceTypeName).Error(),
		)
		return
	}
}

func (r *datasourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, idSeparator)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Importing %s", datasourceTypeName),
			fmt.Sprintf("The import ID must be `<InstanceId>:<DatasourceId>`, got %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_id"), types.StringValue(parts[0]))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("datasource_id"), types.StringValue(parts[1]))...)
}

// getDatasource queries GetDatasource by InstanceId and DatasourceId and returns
// the Datasource object, or nil when it no longer exists.
func getDatasource(client *connectivity.AliyunClient, instanceId, datasourceId string) (map[string]interface{}, error) {
	response, err := client.RpcPost(paifsProduct, paifsVersion, "GetDatasource", nil, map[string]interface{}{
		"InstanceId":   instanceId,
		"DatasourceId": datasourceId,
	}, true)
	if err != nil {
		if isDatasourceNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	object, _ := response["Datasource"].(map[string]interface{})
	if object == nil {
		// GetDatasource returns the fields at the root of the response; fall back
		// to the response itself when the API does not wrap them under Datasource.
		if _, ok := response["DatasourceId"]; ok {
			return response, nil
		}
		return nil, fmt.Errorf("the GetDatasource response carried no Datasource object: %v", response)
	}
	return object, nil
}

// datasourceState maps a GetDatasource response onto the resource model.
// InstanceId is not returned by GetDatasource, so it is carried in from the
// state that triggered the read. Empty strings become null.
func datasourceState(object map[string]interface{}, instanceId, datasourceId string) datasourceResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	return datasourceResourceModel{
		Id:             types.StringValue(fmt.Sprintf("%s%s%s", instanceId, idSeparator, datasourceId)),
		InstanceId:     nonEmptyString(instanceId),
		DatasourceId:   nonEmptyString(datasourceId),
		Type:           str("Type"),
		DatasourceName: str("Name"),
		WorkspaceId:    str("WorkspaceId"),
		Uri:            str("Uri"),
		Config:         str("Config"),
		CreateTime:     str("GmtCreateTime"),
	}
}

// datasourceIdentifiers resolves the InstanceId and DatasourceId from the
// resource state, falling back to splitting the composite id when the
// individual attributes are absent.
func datasourceIdentifiers(state datasourceResourceModel) (instanceId, datasourceId string) {
	instanceId = state.InstanceId.ValueString()
	datasourceId = state.DatasourceId.ValueString()
	if instanceId == "" || datasourceId == "" {
		parts := strings.Split(state.Id.ValueString(), idSeparator)
		if len(parts) == 2 {
			if instanceId == "" {
				instanceId = parts[0]
			}
			if datasourceId == "" {
				datasourceId = parts[1]
			}
		}
	}
	return instanceId, datasourceId
}

func nonEmptyString(v string) types.String {
	if v == "" {
		return types.StringNull()
	}
	return types.StringValue(v)
}

// isDatasourceNotExist reports whether err is the PaiFeatureStore
// datasource-not-found error. The code arrives on a *tea.SDKError from RpcPost;
// the substring fallback covers wrapped forms.
func isDatasourceNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil && *e.Code == datasourceNotExist {
		return true
	}
	return strings.Contains(err.Error(), datasourceNotExist)
}
