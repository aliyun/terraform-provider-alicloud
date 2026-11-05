package dataworks

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	dataQualityScanTypeName = "alicloud_data_works_data_quality_scan"
	// Framework resources cannot expose user-configurable per-operation
	// timeouts; dataQualityScanTimeout mirrors the SDKv2 5-minute default for
	// create/update/delete.
	dataQualityScanTimeout     = 5 * time.Minute
	dataQualityScanReadTimeout = 1 * time.Minute
)

var (
	_ resource.Resource                = &dataQualityScanResource{}
	_ resource.ResourceWithImportState = &dataQualityScanResource{}
)

func NewDataQualityScanResource() resource.Resource {
	return &dataQualityScanResource{}
}

type dataQualityScanResource struct {
	fwadapt.ResourceBase
}

// dataQualityScanResourceModel is the Terraform state shape. The nested
// configuration blocks (Parameters, Trigger, ComputeResource, RuntimeResource,
// Hooks) are typed so the API's structured responses do not create permanent
// diffs against user-supplied JSON strings; Spec stays a plain string because
// the API stores and returns it verbatim.
type dataQualityScanResourceModel struct {
	Id              types.String                    `tfsdk:"id"`
	Name            types.String                    `tfsdk:"name"`
	ProjectId       types.String                    `tfsdk:"project_id"`
	Spec            types.String                    `tfsdk:"spec"`
	Description     types.String                    `tfsdk:"description"`
	Owner           types.String                    `tfsdk:"owner"`
	Parameters      []dataQualityScanParameter      `tfsdk:"parameters"`
	Trigger         *dataQualityScanTrigger         `tfsdk:"trigger"`
	ComputeResource *dataQualityScanComputeResource `tfsdk:"compute_resource"`
	RuntimeResource *dataQualityScanRuntimeResource `tfsdk:"runtime_resource"`
	Hooks           []dataQualityScanHook           `tfsdk:"hooks"`
	CreateTime      types.String                    `tfsdk:"create_time"`
	ModifyTime      types.String                    `tfsdk:"modify_time"`
	CreateUser      types.String                    `tfsdk:"create_user"`
	ModifyUser      types.String                    `tfsdk:"modify_user"`
}

type dataQualityScanParameter struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type dataQualityScanTrigger struct {
	Type    types.String `tfsdk:"type"`
	TaskIds types.List   `tfsdk:"task_ids"`
}

type dataQualityScanComputeResource struct {
	Name    types.String                   `tfsdk:"name"`
	EnvType types.String                   `tfsdk:"env_type"`
	Runtime *dataQualityScanComputeRuntime `tfsdk:"runtime"`
}

type dataQualityScanComputeRuntime struct {
	Engine    types.String `tfsdk:"engine"`
	HiveConf  types.Map    `tfsdk:"hive_conf"`
	SparkConf types.Map    `tfsdk:"spark_conf"`
}

type dataQualityScanRuntimeResource struct {
	Id    types.String  `tfsdk:"id"`
	Cu    types.Float64 `tfsdk:"cu"`
	Image types.String  `tfsdk:"image"`
}

type dataQualityScanHook struct {
	Condition types.String `tfsdk:"condition"`
	Type      types.String `tfsdk:"type"`
}

func (r *dataQualityScanResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = dataQualityScanTypeName
}

func (r *dataQualityScanResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"spec": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"owner": schema.StringAttribute{
				Optional: true,
			},
			"parameters": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"value": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"trigger": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
					},
					"task_ids": schema.ListAttribute{
						ElementType: types.Int64Type,
						Optional:    true,
					},
				},
			},
			"compute_resource": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Optional: true,
					},
					"env_type": schema.StringAttribute{
						Optional: true,
					},
					"runtime": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"engine": schema.StringAttribute{
								Optional: true,
							},
							"hive_conf": schema.MapAttribute{
								ElementType: types.StringType,
								Optional:    true,
								PlanModifiers: []planmodifier.Map{
									mapplanmodifier.RequiresReplace(),
								},
							},
							"spark_conf": schema.MapAttribute{
								ElementType: types.StringType,
								Optional:    true,
								PlanModifiers: []planmodifier.Map{
									mapplanmodifier.RequiresReplace(),
								},
							},
						},
					},
				},
			},
			"runtime_resource": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Optional: true,
					},
					"cu": schema.Float64Attribute{
						Optional: true,
					},
					"image": schema.StringAttribute{
						Optional: true,
					},
				},
			},
			"hooks": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"condition": schema.StringAttribute{
							Required: true,
						},
						"type": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"create_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modify_time": schema.StringAttribute{
				Computed: true,
			},
			"create_user": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modify_user": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *dataQualityScanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dataQualityScanResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request, err := dataQualityScanRequest(plan, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: building CreateDataQualityScan request", dataQualityScanTypeName),
			err.Error(),
		)
		return
	}
	request["RegionId"] = r.Meta().RegionId

	action := "CreateDataQualityScan"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, dataQualityScanTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("dataworks-public", "2024-05-18", action, nil, request, true)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDataQualityScan", dataQualityScanTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, dataQualityScanTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	rawId, err := jsonpath.Get("$.Id", response)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading $.Id", dataQualityScanTypeName),
			errs.WrapErrorf(err, errs.FailedGetAttributeMsg, dataQualityScanTypeName, "$.Id", response).Error(),
		)
		return
	}
	id, err := dataQualityScanId(rawId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading $.Id", dataQualityScanTypeName),
			err.Error(),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeDataQualityScanContext(ctx, id)
	if err != nil {
		// Keep the created ID in state so the scan stays trackable when the
		// read-back fails; Terraform marks the resource as tainted.
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: calling GetDataQualityScan", dataQualityScanTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetDataQualityScan", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, err := dataQualityScanState(object)
	if err != nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: mapping GetDataQualityScan response", dataQualityScanTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *dataQualityScanResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dataQualityScanResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if id == "" {
		return
	}

	object, err := NewService(r.Meta()).DescribeDataQualityScanContext(ctx, id)
	if err != nil {
		if errs.IsProviderNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetDataQualityScan", dataQualityScanTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetDataQualityScan", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, err := dataQualityScanState(object)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: mapping GetDataQualityScan response", dataQualityScanTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

// Update sends the full configured object (CreateDataQualityScan's style=json
// parameters form the whole editable shape, so every configured field is
// resent; the API has no partial-patch) and then re-reads the scan.
func (r *dataQualityScanResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dataQualityScanResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The Id is not part of the plan (it is Computed); carry it from state.
	var state dataQualityScanResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Id = state.Id

	request, err := dataQualityScanRequest(plan, false)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: building UpdateDataQualityScan request", dataQualityScanTypeName),
			err.Error(),
		)
		return
	}
	request["RegionId"] = r.Meta().RegionId

	action := "UpdateDataQualityScan"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, dataQualityScanTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("dataworks-public", "2024-05-18", action, nil, request, true)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateDataQualityScan", dataQualityScanTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeDataQualityScanContext(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update: calling GetDataQualityScan", dataQualityScanTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), "GetDataQualityScan", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, err := dataQualityScanState(object)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update: mapping GetDataQualityScan response", dataQualityScanTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *dataQualityScanResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dataQualityScanResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := NewService(r.Meta()).DeleteDataQualityScanContext(ctx, state.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteDataQualityScan", dataQualityScanTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), "DeleteDataQualityScan", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
}

func (r *dataQualityScanResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// dataQualityScanId converts the numeric Id returned by CreateDataQualityScan
// (declared as number, sometimes returned as a numeric string) into the
// canonical decimal identity used in state, read-back and delete requests.
func dataQualityScanId(raw interface{}) (string, error) {
	switch v := raw.(type) {
	case string:
		if v == "" {
			return "", fmt.Errorf("Id is empty")
		}
		return v, nil
	case json.Number:
		return fmt.Sprint(v), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case nil:
		return "", fmt.Errorf("Id is missing")
	default:
		return "", fmt.Errorf("Id has unexpected type %T", raw)
	}
}

// dataQualityScanRequest builds the request map shared by Create and Update.
// The style=json parameters (Parameters, Trigger, ComputeResource,
// RuntimeResource, Hooks) are serialized to JSON strings, matching the
// dataworks-public encoding for structured form parameters (see di_job).
// Configured (non-null) optional blocks are sent; null blocks are omitted.
func dataQualityScanRequest(plan dataQualityScanResourceModel, isCreate bool) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"Name":      plan.Name.ValueString(),
		"ProjectId": plan.ProjectId.ValueString(),
		"Spec":      plan.Spec.ValueString(),
	}
	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		request["Description"] = plan.Description.ValueString()
	}
	if !plan.Owner.IsNull() && plan.Owner.ValueString() != "" {
		request["Owner"] = plan.Owner.ValueString()
	}
	if isCreate {
		request["ClientToken"] = helper.BuildClientToken("CreateDataQualityScan")
	} else {
		request["Id"] = plan.Id.ValueString()
	}
	if len(plan.Parameters) > 0 {
		params := make([]map[string]interface{}, 0, len(plan.Parameters))
		for _, p := range plan.Parameters {
			params = append(params, map[string]interface{}{
				"Name":  p.Name.ValueString(),
				"Value": p.Value.ValueString(),
			})
		}
		b, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("serializing Parameters: %w", err)
		}
		request["Parameters"] = string(b)
	}
	if plan.Trigger != nil {
		trigger := map[string]interface{}{
			"Type": plan.Trigger.Type.ValueString(),
		}
		if !plan.Trigger.TaskIds.IsNull() && len(plan.Trigger.TaskIds.Elements()) > 0 {
			ids := make([]int64, 0, len(plan.Trigger.TaskIds.Elements()))
			for _, e := range plan.Trigger.TaskIds.Elements() {
				if s, ok := e.(types.Int64); ok && !s.IsNull() && !s.IsUnknown() {
					ids = append(ids, s.ValueInt64())
				}
			}
			if len(ids) > 0 {
				trigger["TaskIds"] = ids
			}
		}
		b, err := json.Marshal(trigger)
		if err != nil {
			return nil, fmt.Errorf("serializing Trigger: %w", err)
		}
		request["Trigger"] = string(b)
	}
	if plan.ComputeResource != nil {
		cr := map[string]interface{}{}
		if !plan.ComputeResource.Name.IsNull() && plan.ComputeResource.Name.ValueString() != "" {
			cr["Name"] = plan.ComputeResource.Name.ValueString()
		}
		if !plan.ComputeResource.EnvType.IsNull() && plan.ComputeResource.EnvType.ValueString() != "" {
			cr["EnvType"] = plan.ComputeResource.EnvType.ValueString()
		}
		if plan.ComputeResource.Runtime != nil {
			rt := map[string]interface{}{}
			if !plan.ComputeResource.Runtime.Engine.IsNull() && plan.ComputeResource.Runtime.Engine.ValueString() != "" {
				rt["Engine"] = plan.ComputeResource.Runtime.Engine.ValueString()
			}
			if hive, ok := dataQualityScanConfToMap(plan.ComputeResource.Runtime.HiveConf); ok {
				rt["HiveConf"] = hive
			}
			if spark, ok := dataQualityScanConfToMap(plan.ComputeResource.Runtime.SparkConf); ok {
				rt["SparkConf"] = spark
			}
			cr["Runtime"] = rt
		}
		b, err := json.Marshal(cr)
		if err != nil {
			return nil, fmt.Errorf("serializing ComputeResource: %w", err)
		}
		request["ComputeResource"] = string(b)
	}
	if plan.RuntimeResource != nil {
		rr := map[string]interface{}{}
		if !plan.RuntimeResource.Id.IsNull() && plan.RuntimeResource.Id.ValueString() != "" {
			rr["Id"] = plan.RuntimeResource.Id.ValueString()
		}
		if !plan.RuntimeResource.Cu.IsNull() && !plan.RuntimeResource.Cu.IsUnknown() {
			rr["Cu"] = plan.RuntimeResource.Cu.ValueFloat64()
		}
		if !plan.RuntimeResource.Image.IsNull() && plan.RuntimeResource.Image.ValueString() != "" {
			rr["Image"] = plan.RuntimeResource.Image.ValueString()
		}
		b, err := json.Marshal(rr)
		if err != nil {
			return nil, fmt.Errorf("serializing RuntimeResource: %w", err)
		}
		request["RuntimeResource"] = string(b)
	}
	if len(plan.Hooks) > 0 {
		hooks := make([]map[string]interface{}, 0, len(plan.Hooks))
		for _, h := range plan.Hooks {
			hooks = append(hooks, map[string]interface{}{
				"Condition": h.Condition.ValueString(),
				"Type":      h.Type.ValueString(),
			})
		}
		b, err := json.Marshal(hooks)
		if err != nil {
			return nil, fmt.Errorf("serializing Hooks: %w", err)
		}
		request["Hooks"] = string(b)
	}
	return request, nil
}

// dataQualityScanConfToMap converts a configured map attribute into a plain
// map[string]interface{} of string values for JSON serialization. Returns
// ok=false when the map is null or empty so the caller omits the field.
func dataQualityScanConfToMap(value types.Map) (map[string]interface{}, bool) {
	if value.IsNull() {
		return nil, false
	}
	elems := value.Elements()
	if len(elems) == 0 {
		return nil, false
	}
	result := make(map[string]interface{}, len(elems))
	for k, v := range elems {
		if s, ok := v.(types.String); ok && !s.IsNull() {
			result[k] = s.ValueString()
		}
	}
	if len(result) == 0 {
		return nil, false
	}
	return result, true
}

// dataQualityScanState maps the GetDataQualityScan response object into the
// Framework model. Contract drift (missing required scalar, unexpected nested
// type) fails the read instead of overwriting configured state with null.
func dataQualityScanState(object map[string]interface{}) (dataQualityScanResourceModel, error) {
	id, err := dataQualityScanId(object["Id"])
	if err != nil {
		return dataQualityScanResourceModel{}, err
	}
	name, err := dataQualityScanRequiredString(object["Name"], "Name")
	if err != nil {
		return dataQualityScanResourceModel{}, err
	}
	projectId, err := dataQualityScanRequiredString(object["ProjectId"], "ProjectId")
	if err != nil {
		return dataQualityScanResourceModel{}, err
	}
	spec, err := dataQualityScanRequiredString(object["Spec"], "Spec")
	if err != nil {
		return dataQualityScanResourceModel{}, err
	}

	model := dataQualityScanResourceModel{
		Id:          types.StringValue(id),
		Name:        name,
		ProjectId:   projectId,
		Spec:        spec,
		Description: dataQualityScanString(object["Description"]),
		Owner:       dataQualityScanString(object["Owner"]),
		CreateTime:  dataQualityScanString(object["CreateTime"]),
		ModifyTime:  dataQualityScanString(object["ModifyTime"]),
		CreateUser:  dataQualityScanString(object["CreateUser"]),
		ModifyUser:  dataQualityScanString(object["ModifyUser"]),
	}

	if params, ok := object["Parameters"].([]interface{}); ok {
		model.Parameters = dataQualityScanParametersFromResponse(params)
	}
	if triggerRaw, ok := object["Trigger"].(map[string]interface{}); ok {
		model.Trigger = dataQualityScanTriggerFromResponse(triggerRaw)
	}
	if crRaw, ok := object["ComputeResource"].(map[string]interface{}); ok {
		model.ComputeResource = dataQualityScanComputeResourceFromResponse(crRaw)
	}
	if rrRaw, ok := object["RuntimeResource"].(map[string]interface{}); ok {
		model.RuntimeResource = dataQualityScanRuntimeResourceFromResponse(rrRaw)
	}
	if hooks, ok := object["Hooks"].([]interface{}); ok {
		model.Hooks = dataQualityScanHooksFromResponse(hooks)
	}

	return model, nil
}

func dataQualityScanParametersFromResponse(raw []interface{}) []dataQualityScanParameter {
	if len(raw) == 0 {
		return nil
	}
	out := make([]dataQualityScanParameter, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, dataQualityScanParameter{
			Name:  dataQualityScanString(m["Name"]),
			Value: dataQualityScanString(m["Value"]),
		})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func dataQualityScanTriggerFromResponse(raw map[string]interface{}) *dataQualityScanTrigger {
	t := &dataQualityScanTrigger{
		Type: dataQualityScanString(raw["Type"]),
	}
	if ids, ok := raw["TaskIds"].([]interface{}); ok && len(ids) > 0 {
		elems := make([]attr.Value, 0, len(ids))
		for _, id := range ids {
			if v, ok := id.(float64); ok {
				elems = append(elems, types.Int64Value(int64(v)))
			} else if s, ok := id.(string); ok {
				if n, err := strconv.ParseInt(s, 10, 64); err == nil {
					elems = append(elems, types.Int64Value(n))
				}
			} else if n, ok := id.(json.Number); ok {
				if v, err := n.Int64(); err == nil {
					elems = append(elems, types.Int64Value(v))
				}
			}
		}
		if len(elems) > 0 {
			t.TaskIds = types.ListValueMust(types.Int64Type, elems)
		}
	}
	if t.TaskIds.IsNull() {
		t.TaskIds = types.ListNull(types.Int64Type)
	}
	return t
}

func dataQualityScanComputeResourceFromResponse(raw map[string]interface{}) *dataQualityScanComputeResource {
	cr := &dataQualityScanComputeResource{
		Name:    dataQualityScanString(raw["Name"]),
		EnvType: dataQualityScanString(raw["EnvType"]),
	}
	if runtimeRaw, ok := raw["Runtime"].(map[string]interface{}); ok {
		cr.Runtime = &dataQualityScanComputeRuntime{
			Engine:    dataQualityScanString(runtimeRaw["Engine"]),
			HiveConf:  dataQualityScanConfFromResponse(runtimeRaw["HiveConf"]),
			SparkConf: dataQualityScanConfFromResponse(runtimeRaw["SparkConf"]),
		}
	}
	return cr
}

func dataQualityScanRuntimeResourceFromResponse(raw map[string]interface{}) *dataQualityScanRuntimeResource {
	rr := &dataQualityScanRuntimeResource{
		Id:    dataQualityScanString(raw["Id"]),
		Image: dataQualityScanString(raw["Image"]),
	}
	switch v := raw["Cu"].(type) {
	case float64:
		rr.Cu = types.Float64Value(v)
	case json.Number:
		if f, err := v.Float64(); err == nil {
			rr.Cu = types.Float64Value(f)
		}
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			rr.Cu = types.Float64Value(f)
		}
	}
	return rr
}

func dataQualityScanHooksFromResponse(raw []interface{}) []dataQualityScanHook {
	if len(raw) == 0 {
		return nil
	}
	out := make([]dataQualityScanHook, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, dataQualityScanHook{
			Condition: dataQualityScanString(m["Condition"]),
			Type:      dataQualityScanString(m["Type"]),
		})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// dataQualityScanConfFromResponse normalizes a HiveConf/SparkConf value (a map
// whose values may be strings or numbers) into a Map of strings. A null map
// is returned when the field is absent.
func dataQualityScanConfFromResponse(raw interface{}) types.Map {
	m, ok := raw.(map[string]interface{})
	if !ok || len(m) == 0 {
		return types.MapNull(types.StringType)
	}
	elems := make(map[string]attr.Value, len(m))
	for k, v := range m {
		elems[k] = types.StringValue(fmt.Sprint(v))
	}
	return types.MapValueMust(types.StringType, elems)
}

// dataQualityScanRequiredString rejects a missing or non-string required
// attribute instead of silently writing null over configured state.
func dataQualityScanRequiredString(raw interface{}, field string) (types.String, error) {
	switch v := raw.(type) {
	case string:
		if v == "" {
			return types.StringNull(), fmt.Errorf("%s is required in the GetDataQualityScan response, got empty string", field)
		}
		return types.StringValue(v), nil
	case float64:
		return types.StringValue(strconv.FormatFloat(v, 'f', -1, 64)), nil
	case json.Number:
		return types.StringValue(v.String()), nil
	case int:
		return types.StringValue(strconv.Itoa(v)), nil
	case int64:
		return types.StringValue(strconv.FormatInt(v, 10)), nil
	default:
		return types.StringNull(), fmt.Errorf("%s is required in the GetDataQualityScan response, got %T %v", field, raw, raw)
	}
}

// dataQualityScanString normalizes an optional scalar (string or number) into
// a String; missing/non-scalar values become null.
func dataQualityScanString(raw interface{}) types.String {
	switch v := raw.(type) {
	case string:
		if v == "" {
			return types.StringNull()
		}
		return types.StringValue(v)
	case float64:
		return types.StringValue(strconv.FormatFloat(v, 'f', -1, 64))
	case json.Number:
		return types.StringValue(v.String())
	case int:
		return types.StringValue(strconv.Itoa(v))
	case int64:
		return types.StringValue(strconv.FormatInt(v, 10))
	default:
		return types.StringNull()
	}
}
