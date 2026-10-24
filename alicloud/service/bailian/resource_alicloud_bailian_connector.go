package bailian

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
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

const connectorTypeName = "alicloud_bailian_connector"

var (
	_ resource.Resource                = &connectorResource{}
	_ resource.ResourceWithImportState = &connectorResource{}
)

func NewConnectorResource() resource.Resource {
	return &connectorResource{}
}

type connectorResource struct {
	fwadapt.ResourceBase
}

type connectorResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	WorkspaceId         types.String `tfsdk:"workspace_id"`
	ConnectorId         types.String `tfsdk:"connector_id"`
	ConnectorName       types.String `tfsdk:"connector_name"`
	ConnectorType       types.String `tfsdk:"connector_type"`
	Description         types.String `tfsdk:"description"`
	FileConnectorConfig types.String `tfsdk:"file_connector_config"`
}

func (r *connectorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = connectorTypeName
}

func (r *connectorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Bailian Connector resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, format as `<workspace_id>:<connector_id>`.",
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
			"connector_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the connector.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connector_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the connector.",
			},
			"connector_type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the connector. Valid value: `FILE`.",
				Validators: []validator.String{
					stringvalidator.OneOf("FILE"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the connector.",
			},
			"file_connector_config": schema.StringAttribute{
				Optional:    true,
				Description: "The file connector config, a JSON string with fields `bucket_name`, `region_id`, `storage_type`.",
			},
		},
	}
}

func (r *connectorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan connectorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "AddConnector"
	llmId := extractLlmId(plan.WorkspaceId.ValueString())
	clientToken := helper.BuildClientToken(action)
	pathName := fmt.Sprintf("/%s/datacenter/connector", llmId)
	query := map[string]*string{
		"ClientToken": strPtr(clientToken),
	}
	// ConnectorName, ConnectorType, Description and FileConnectorConfig are
	// formData parameters per the bailian OpenAPI definition, and AddConnector
	// consumes application/json; send them in the request body as a JSON
	// object (via RoaPost), not as form-urlencoded data or URL query parameters.
	body := map[string]interface{}{
		"ConnectorName": plan.ConnectorName.ValueString(),
		"ConnectorType": plan.ConnectorType.ValueString(),
	}
	if !plan.Description.IsNull() {
		body["Description"] = plan.Description.ValueString()
	}
	if !plan.FileConnectorConfig.IsNull() {
		body["FileConnectorConfig"] = plan.FileConnectorConfig.ValueString()
	}

	response, err := r.Client().RoaPost("bailian", "2023-12-29", pathName, query, nil, body, true)
	logs.AddDebug(action, response, body)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", connectorTypeName, action),
			errs.WrapErrorf(err, "Creating %s: calling %s", connectorTypeName, action).Error(),
		)
		return
	}

	data, _ := response["Data"].(map[string]interface{})
	if data == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", connectorTypeName, action),
			fmt.Sprintf("The response carried no Data object: %v", response),
		)
		return
	}
	connectorId, _ := data["ConnectorId"].(string)
	if connectorId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", connectorTypeName, action),
			fmt.Sprintf("The response Data carried no ConnectorId: %v", response),
		)
		return
	}

	object, err := getConnector(r.Client(), plan.WorkspaceId.ValueString(), connectorId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetConnector", connectorTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetConnector", connectorTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", connectorTypeName, action),
			fmt.Sprintf("The connector %s is not visible via GetConnector right after creation", connectorId),
		)
		return
	}
	state := connectorState(object, plan.WorkspaceId.ValueString())
	// connector_type is write-only via AddConnector but not returned by GetConnector;
	// preserve the planned value so refresh does not drop it.
	state.ConnectorType = plan.ConnectorType
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *connectorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state connectorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId, connectorId, ok := parseCompositeId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s", connectorTypeName),
			fmt.Sprintf("Cannot parse resource ID %q: expected format `<workspace_id>:<connector_id>`", state.Id.ValueString()),
		)
		return
	}

	object, err := getConnector(r.Client(), workspaceId, connectorId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetConnector", connectorTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetConnector", connectorTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	model := connectorState(object, workspaceId)
	// connector_type is write-only; keep the previously configured value.
	model.ConnectorType = state.ConnectorType
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *connectorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state connectorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId, connectorId, _ := parseCompositeId(state.Id.ValueString())
	llmId := extractLlmId(workspaceId)
	pathName := fmt.Sprintf("/%s/datacenter/connector/%s", llmId, connectorId)
	// ConnectorName and Description are formData parameters per the bailian
	// OpenAPI definition, and UpdateConnector consumes application/json; send
	// them in the request body as a JSON object (via RoaPut), not form-encoded.
	body := map[string]interface{}{
		"ConnectorName": plan.ConnectorName.ValueString(),
	}
	if !plan.Description.IsNull() {
		body["Description"] = plan.Description.ValueString()
	}

	response, err := r.Client().RoaPut("bailian", "2023-12-29", pathName, nil, nil, body, true)
	logs.AddDebug("UpdateConnector", response, body)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateConnector", connectorTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateConnector", connectorTypeName).Error(),
		)
		return
	}

	object, err := getConnector(r.Client(), workspaceId, connectorId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetConnector", connectorTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetConnector", connectorTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	model := connectorState(object, workspaceId)
	model.ConnectorType = plan.ConnectorType
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *connectorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state connectorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceId, connectorId, _ := parseCompositeId(state.Id.ValueString())
	llmId := extractLlmId(workspaceId)
	pathName := fmt.Sprintf("/%s/datacenter/connector/%s", llmId, connectorId)

	response, err := r.Client().RoaDelete("bailian", "2023-12-29", pathName, nil, nil, nil, true)
	logs.AddDebug("DeleteConnector", response, nil)
	if err != nil && !isBailianNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteConnector", connectorTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteConnector", connectorTypeName).Error(),
		)
		return
	}
}

func (r *connectorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getConnector(client *connectivity.AliyunClient, workspaceId, connectorId string) (map[string]interface{}, error) {
	llmId := extractLlmId(workspaceId)
	pathName := fmt.Sprintf("/%s/datacenter/connector", llmId)
	query := map[string]*string{
		"ConnectorId": strPtr(connectorId),
	}
	response, err := client.RoaGet("bailian", "2023-12-29", pathName, query, nil, nil)
	if err != nil {
		if isBailianNotExist(err) || errs.IsExpectedErrors(err, []string{"ConnectorNotFound"}) {
			return nil, nil
		}
		return nil, err
	}
	data, _ := response["Data"].(map[string]interface{})
	if data == nil {
		return nil, fmt.Errorf("the GetConnector response carried no Data object: %v", response)
	}
	return data, nil
}

func connectorState(object map[string]interface{}, workspaceId string) connectorResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	connectorId := str("ConnectorId")
	return connectorResourceModel{
		Id:            types.StringValue(fmt.Sprintf("%s:%s", workspaceId, connectorId.ValueString())),
		WorkspaceId:   types.StringValue(workspaceId),
		ConnectorId:   connectorId,
		ConnectorName: str("ConnectorName"),
		Description:   str("Description"),
	}
}
