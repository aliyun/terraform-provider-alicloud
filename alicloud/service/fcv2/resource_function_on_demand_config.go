package fcv2

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const functionOnDemandConfigTypeName = "alicloud_fcv2_function_on_demand_config"

// functionOnDemandConfigTimeout replaces the SDKv2 resource-level default
// timeouts (Create/Update/Delete each 5 minutes); terraform-plugin-framework
// does not support user-configurable resource timeouts yet.
const functionOnDemandConfigTimeout = 5 * time.Minute

var (
	_ resource.Resource                = &functionOnDemandConfigResource{}
	_ resource.ResourceWithImportState = &functionOnDemandConfigResource{}
)

// NewFunctionOnDemandConfigResource is the factory the ServicePackage
// registration calls once per RPC.
func NewFunctionOnDemandConfigResource() resource.Resource {
	return &functionOnDemandConfigResource{}
}

type functionOnDemandConfigResource struct {
	fwadapt.ResourceBase
}

type functionOnDemandConfigResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	ServiceName          types.String `tfsdk:"service_name"`
	FunctionName         types.String `tfsdk:"function_name"`
	Qualifier            types.String `tfsdk:"qualifier"`
	MaximumInstanceCount types.Int64  `tfsdk:"maximum_instance_count"`
	Resource             types.String `tfsdk:"resource"`
}

func (r *functionOnDemandConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = functionOnDemandConfigTypeName
}

func (r *functionOnDemandConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"function_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// qualifier is a service alias or "LATEST"; the FC-Open API rejects
			// bare version numbers, and the on-demand config is keyed by it, so a
			// change means a different resource.
			"qualifier": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"maximum_instance_count": schema.Int64Attribute{
				Required: true,
			},
			// resource is the FC-internal descriptor "accountId#serviceName#
			// qualifier#functionName"; Get returns it but Put does not accept it.
			"resource": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *functionOnDemandConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan functionOnDemandConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := plan.ServiceName.ValueString()
	functionName := plan.FunctionName.ValueString()
	qualifier := plan.Qualifier.ValueString()
	action, query, body := functionOnDemandConfigPutInput(serviceName, functionName, qualifier, plan.MaximumInstanceCount.ValueInt64())

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(functionOnDemandConfigTimeout, func() *retry.RetryError {
		response, err = r.Meta().RoaPut("FC-Open", "2021-04-06", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, body)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling PutFunctionOnDemandConfig", functionOnDemandConfigTypeName),
			errs.WrapErrorf(err, "Creating %s: calling PutFunctionOnDemandConfig", functionOnDemandConfigTypeName).Error(),
		)
		return
	}

	plan.Id = types.StringValue(functionOnDemandConfigId(serviceName, functionName, qualifier))

	object, err := NewService(r.Meta()).DescribeFunctionOnDemandConfig(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading after create", functionOnDemandConfigTypeName),
			errs.WrapErrorf(err, "Creating %s: reading after create", functionOnDemandConfigTypeName).Error(),
		)
		return
	}
	state := functionOnDemandConfigState(object)
	// Identity inputs are not echoed by Get; keep them from the plan so the
	// state matches the configuration and ImportState/refresh round-trip.
	state.Id = plan.Id
	state.ServiceName = plan.ServiceName
	state.FunctionName = plan.FunctionName
	state.Qualifier = plan.Qualifier
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *functionOnDemandConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state functionOnDemandConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := NewService(r.Meta()).DescribeFunctionOnDemandConfig(state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetFunctionOnDemandConfig", functionOnDemandConfigTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetFunctionOnDemandConfig", functionOnDemandConfigTypeName).Error(),
		)
		return
	}
	model := functionOnDemandConfigState(object)
	model.Id = state.Id
	model.ServiceName = state.ServiceName
	model.FunctionName = state.FunctionName
	model.Qualifier = state.Qualifier
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *functionOnDemandConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan functionOnDemandConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Identity fields are all ForceNew, so an Update only happens when
	// maximum_instance_count changes. PutFunctionOnDemandConfig is an upsert:
	// the same call path serves create and update.
	serviceName := plan.ServiceName.ValueString()
	functionName := plan.FunctionName.ValueString()
	qualifier := plan.Qualifier.ValueString()
	action, query, body := functionOnDemandConfigPutInput(serviceName, functionName, qualifier, plan.MaximumInstanceCount.ValueInt64())

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(functionOnDemandConfigTimeout, func() *retry.RetryError {
		response, err = r.Meta().RoaPut("FC-Open", "2021-04-06", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, body)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling PutFunctionOnDemandConfig", functionOnDemandConfigTypeName),
			errs.WrapErrorf(err, "Updating %s: calling PutFunctionOnDemandConfig", functionOnDemandConfigTypeName).Error(),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeFunctionOnDemandConfig(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: reading after update", functionOnDemandConfigTypeName),
			errs.WrapErrorf(err, "Updating %s: reading after update", functionOnDemandConfigTypeName).Error(),
		)
		return
	}
	state := functionOnDemandConfigState(object)
	state.Id = plan.Id
	state.ServiceName = plan.ServiceName
	state.FunctionName = plan.FunctionName
	state.Qualifier = plan.Qualifier
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *functionOnDemandConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state functionOnDemandConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := state.ServiceName.ValueString()
	functionName := state.FunctionName.ValueString()
	qualifier := state.Qualifier.ValueString()
	action := fmt.Sprintf("/2021-04-06/services/%s/functions/%s/on-demand-config", serviceName, functionName)
	query := map[string]*string{"qualifier": helper.StringPointer(qualifier)}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(functionOnDemandConfigTimeout, func() *retry.RetryError {
		response, err = r.Meta().RoaDelete("FC-Open", "2021-04-06", action, query, nil, nil, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, query)
	if err != nil {
		if errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteFunctionOnDemandConfig", functionOnDemandConfigTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteFunctionOnDemandConfig", functionOnDemandConfigTypeName).Error(),
		)
		return
	}
}

func (r *functionOnDemandConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Importing %s", functionOnDemandConfigTypeName),
			fmt.Sprintf("invalid import id %q, expected the format serviceName:functionName:qualifier (qualifier is a service alias or LATEST)", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), types.StringValue(parts[0]))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("function_name"), types.StringValue(parts[1]))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("qualifier"), types.StringValue(parts[2]))...)
}

// functionOnDemandConfigId is the Terraform resource id; serviceName and
// functionName cannot contain ":", so the three-part split is unambiguous.
func functionOnDemandConfigId(serviceName, functionName, qualifier string) string {
	return fmt.Sprintf("%s:%s:%s", serviceName, functionName, qualifier)
}

// functionOnDemandConfigPutInput builds the ROA path, query and body shared by
// Create and Update: PutFunctionOnDemandConfig is an upsert, so both use the
// same call.
func functionOnDemandConfigPutInput(serviceName, functionName, qualifier string, maximumInstanceCount int64) (string, map[string]*string, map[string]interface{}) {
	action := fmt.Sprintf("/2021-04-06/services/%s/functions/%s/on-demand-config", serviceName, functionName)
	query := map[string]*string{"qualifier": helper.StringPointer(qualifier)}
	body := map[string]interface{}{
		"maximumInstanceCount": maximumInstanceCount,
	}
	return action, query, body
}

// functionOnDemandConfigState maps the GetFunctionOnDemandConfig response onto
// the resource model. Identity fields are not part of the response, so callers
// overlay them from the plan/state before writing.
func functionOnDemandConfigState(object map[string]interface{}) functionOnDemandConfigResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	// The tea-rpc client decodes maximumInstanceCount as json.Number, int or
	// float64 depending on the transport; every encoding must map back to Int64
	// or the planned value collapses to null ("inconsistent result after apply").
	int64Val := func(key string) types.Int64 {
		switch v := object[key].(type) {
		case float64:
			return types.Int64Value(int64(v))
		case json.Number:
			if n, err := v.Int64(); err == nil {
				return types.Int64Value(n)
			}
		case int:
			return types.Int64Value(int64(v))
		case int64:
			return types.Int64Value(v)
		case string:
			if n, err := strconv.ParseInt(v, 10, 64); err == nil {
				return types.Int64Value(n)
			}
		}
		return types.Int64Null()
	}
	return functionOnDemandConfigResourceModel{
		MaximumInstanceCount: int64Val("maximumInstanceCount"),
		Resource:             str("resource"),
	}
}
