package cms

import (
	"context"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const endpointConnectorTypeName = "alicloud_cms_endpoint_connector"

// endpointConnectorTimeout replaces the SDKv2 resource-level default timeouts
// (Create/Update/Delete each 5 minutes); terraform-plugin-framework does not
// support user-configurable resource timeouts yet.
const endpointConnectorTimeout = 5 * time.Minute

var (
	_ resource.Resource                = &endpointConnectorResource{}
	_ resource.ResourceWithImportState = &endpointConnectorResource{}
)

func NewEndpointConnectorResource() resource.Resource {
	return &endpointConnectorResource{}
}

type endpointConnectorResource struct {
	fwadapt.ResourceBase
}

type headerModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type endpointConnectorResourceModel struct {
	Id          types.String  `tfsdk:"id"`
	Workspace   types.String  `tfsdk:"workspace"`
	Type        types.String  `tfsdk:"type"`
	Name        types.String  `tfsdk:"name"`
	Endpoint    types.String  `tfsdk:"endpoint"`
	Alias       types.String  `tfsdk:"alias"`
	Description types.String  `tfsdk:"description"`
	Credential  types.Map     `tfsdk:"credential"`
	Headers     []headerModel `tfsdk:"headers"`
	Properties  types.Map     `tfsdk:"properties"`
	ConnectorId types.String  `tfsdk:"connector_id"`
	RegionId    types.String  `tfsdk:"region_id"`
	CreatedAt   types.Int64   `tfsdk:"created_at"`
	UpdatedAt   types.Int64   `tfsdk:"updated_at"`
}

func (r *endpointConnectorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = endpointConnectorTypeName
}

func (r *endpointConnectorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("model_service", "agent_app"),
				},
			},
			"name":        schema.StringAttribute{Required: true},
			"endpoint":    schema.StringAttribute{Required: true},
			"alias":       schema.StringAttribute{Optional: true},
			"description": schema.StringAttribute{Optional: true},
			"credential": schema.MapAttribute{
				Required:    true,
				Sensitive:   true,
				ElementType: types.StringType,
			},
			"properties": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"connector_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.Int64Attribute{Computed: true},
			"updated_at": schema.Int64Attribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			// headers uses ListNestedBlock (not ListNestedAttribute) because the
			// provider is served via tf5muxserver (protocol v5), which cannot
			// convert nested attributes with Attributes set.
			"headers": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"key":   schema.StringAttribute{Optional: true},
						"value": schema.StringAttribute{Optional: true},
					},
				},
			},
		},
	}
}

func (r *endpointConnectorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan endpointConnectorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspace := plan.Workspace.ValueString()
	action := fmt.Sprintf("/api/v1/endpoint-connectors/%s", workspace)
	request := map[string]interface{}{
		"type":     plan.Type.ValueString(),
		"name":     plan.Name.ValueString(),
		"endpoint": plan.Endpoint.ValueString(),
	}
	if !plan.Credential.IsNull() {
		request["credential"] = mapToInterface(plan.Credential)
	}
	if !plan.Alias.IsNull() {
		request["alias"] = plan.Alias.ValueString()
	}
	if !plan.Description.IsNull() {
		request["description"] = plan.Description.ValueString()
	}
	if len(plan.Headers) > 0 {
		request["headers"] = expandHeaders(plan.Headers)
	}
	if !plan.Properties.IsNull() {
		request["properties"] = mapToInterface(plan.Properties)
	}
	body := request
	query := make(map[string]*string)

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(endpointConnectorTimeout, func() *retry.RetryError {
		response, err = r.Client().RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateEndpointConnector", endpointConnectorTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, endpointConnectorTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	connectorId := fmt.Sprint(response["connectorId"])
	if connectorId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateEndpointConnector", endpointConnectorTypeName),
			"connector_id is empty in the CreateEndpointConnector response",
		)
		return
	}
	id := fmt.Sprintf("%s:%s", workspace, connectorId)

	object, err := NewService(r.Client()).DescribeEndpointConnector(id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetEndpointConnector", endpointConnectorTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetEndpointConnector", endpointConnectorTypeName).Error(),
		)
		return
	}
	state := endpointConnectorState(object, plan)
	state.Id = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *endpointConnectorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state endpointConnectorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := NewService(r.Client()).DescribeEndpointConnector(state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetEndpointConnector", endpointConnectorTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetEndpointConnector", endpointConnectorTypeName).Error(),
		)
		return
	}
	model := endpointConnectorState(object, state)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *endpointConnectorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state endpointConnectorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspace, connectorId, err := endpointConnectorID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing resource id", endpointConnectorTypeName),
			errs.WrapError(err).Error(),
		)
		return
	}
	action := fmt.Sprintf("/api/v1/endpoint-connectors/%s/%s", workspace, connectorId)
	request := map[string]interface{}{
		"name":     plan.Name.ValueString(),
		"endpoint": plan.Endpoint.ValueString(),
	}
	if !plan.Credential.IsNull() {
		request["credential"] = mapToInterface(plan.Credential)
	}
	if !plan.Alias.IsNull() {
		request["alias"] = plan.Alias.ValueString()
	}
	if !plan.Description.IsNull() {
		request["description"] = plan.Description.ValueString()
	}
	if len(plan.Headers) > 0 {
		request["headers"] = expandHeaders(plan.Headers)
	} else {
		request["headers"] = make([]map[string]interface{}, 0)
	}
	if !plan.Properties.IsNull() {
		request["properties"] = mapToInterface(plan.Properties)
	} else {
		request["properties"] = make(map[string]interface{})
	}
	body := request
	query := make(map[string]*string)

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(endpointConnectorTimeout, func() *retry.RetryError {
		response, err = r.Client().RoaPut("Cms", "2024-03-30", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateEndpointConnector", endpointConnectorTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	object, err := NewService(r.Client()).DescribeEndpointConnector(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetEndpointConnector", endpointConnectorTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetEndpointConnector", endpointConnectorTypeName).Error(),
		)
		return
	}
	model := endpointConnectorState(object, plan)
	model.Id = state.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *endpointConnectorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state endpointConnectorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspace, connectorId, err := endpointConnectorID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing resource id", endpointConnectorTypeName),
			errs.WrapError(err).Error(),
		)
		return
	}
	action := fmt.Sprintf("/api/v1/endpoint-connectors/%s/%s", workspace, connectorId)
	query := make(map[string]*string)

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(endpointConnectorTimeout, func() *retry.RetryError {
		response, err = r.Client().RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EndpointConnectorNotExist", "WorkspaceNotExist"}) || errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteEndpointConnector", endpointConnectorTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
}

func (r *endpointConnectorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// endpointConnectorState maps a GetEndpointConnector response onto the resource
// model. Credential, headers and properties are preserved from the plan/state
// when the API omits them: credential is write-sensitive (Get may return a
// desensitized shell), and headers/properties are optional and may be absent.
func endpointConnectorState(object map[string]interface{}, preserve endpointConnectorResourceModel) endpointConnectorResourceModel {
	reqStr := func(key string) types.String {
		if v, ok := object[key]; ok && v != nil {
			return types.StringValue(fmt.Sprint(v))
		}
		return types.StringNull()
	}
	// Optional strings preserve the configured value (null or "") when the API
	// does not return the key, so a cleared field stays cleared and an unset
	// field stays null instead of producing a plan drift.
	optStr := func(key string, fallback types.String) types.String {
		if v, ok := object[key]; ok && v != nil {
			return types.StringValue(fmt.Sprint(v))
		}
		return fallback
	}
	state := endpointConnectorResourceModel{
		Id:          preserve.Id,
		Workspace:   reqStr("workspace"),
		Type:        reqStr("type"),
		Name:        reqStr("name"),
		Endpoint:    reqStr("endpoint"),
		Alias:       optStr("alias", preserve.Alias),
		Description: optStr("description", preserve.Description),
		Credential:  preserve.Credential,
		Headers:     preserve.Headers,
		Properties:  preserve.Properties,
		ConnectorId: reqStr("connectorId"),
		RegionId:    reqStr("regionId"),
		CreatedAt:   mapInt64(object, "createdAt"),
		UpdatedAt:   mapInt64(object, "updatedAt"),
	}
	if v, ok := object["headers"]; ok && v != nil {
		state.Headers = flattenHeaders(v)
	}
	if v, ok := object["properties"]; ok && v != nil {
		state.Properties = mapToTypeMap(v)
	}
	return state
}

func expandHeaders(headers []headerModel) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(headers))
	for _, h := range headers {
		result = append(result, map[string]interface{}{
			"key":   h.Key.ValueString(),
			"value": h.Value.ValueString(),
		})
	}
	return result
}

func flattenHeaders(v interface{}) []headerModel {
	result := make([]headerModel, 0)
	list, ok := v.([]interface{})
	if !ok {
		return result
	}
	for _, h := range list {
		if h == nil {
			continue
		}
		hm, ok := h.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, headerModel{
			Key:   stringFromMap(hm, "key"),
			Value: stringFromMap(hm, "value"),
		})
	}
	return result
}

func stringFromMap(m map[string]interface{}, key string) types.String {
	if v, ok := m[key]; ok && v != nil {
		return types.StringValue(fmt.Sprint(v))
	}
	return types.StringNull()
}

// mapToTypeMap converts an API response map into a types.Map.
func mapToTypeMap(v interface{}) types.Map {
	if v == nil {
		return types.MapNull(types.StringType)
	}
	m := make(map[string]attr.Value)
	if raw, ok := v.(map[string]interface{}); ok {
		for k, val := range raw {
			if val == nil {
				continue
			}
			m[k] = types.StringValue(fmt.Sprint(val))
		}
	}
	result, d := types.MapValue(types.StringType, m)
	if d.HasError() {
		return types.MapNull(types.StringType)
	}
	return result
}

// mapToInterface extracts a types.Map into a map[string]interface{} for an API
// request. Iterating Elements directly avoids the reflection uncertainty of
// ElementsAs into map[string]interface{}.
func mapToInterface(m types.Map) map[string]interface{} {
	out := make(map[string]interface{})
	for k, v := range m.Elements() {
		if s, ok := v.(types.String); ok {
			out[k] = s.ValueString()
		} else {
			out[k] = fmt.Sprint(v)
		}
	}
	return out
}
