package privatezone

import (
	"context"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	inboundEndpointTypeName = "alicloud_pvtz_inbound_endpoint"
	// inboundEndpointTimeout mirrors the SDKv2 alicloud_pvtz_endpoint create/update default
	// timeouts (10 minutes each). terraform-plugin-framework does not support user-configurable
	// resource timeouts, so the default is applied as a fixed budget.
	inboundEndpointTimeout = 10 * time.Minute
)

var (
	_ resource.Resource                = &inboundEndpointResource{}
	_ resource.ResourceWithImportState = &inboundEndpointResource{}
	_ resource.ResourceWithModifyPlan  = &inboundEndpointResource{}
)

func NewInboundEndpointResource() resource.Resource {
	return &inboundEndpointResource{}
}

type inboundEndpointResource struct {
	fwadapt.ResourceBase
}

type inboundEndpointResourceModel struct {
	Id              types.String `tfsdk:"id"`
	EndpointName    types.String `tfsdk:"endpoint_name"`
	IpConfigs       types.Set    `tfsdk:"ip_configs"`
	SecurityGroupId types.String `tfsdk:"security_group_id"`
	VpcId           types.String `tfsdk:"vpc_id"`
	VpcRegionId     types.String `tfsdk:"vpc_region_id"`
	SpecVersion     types.String `tfsdk:"spec_version"`
	Status          types.String `tfsdk:"status"`
	EndpointType    types.String `tfsdk:"endpoint_type"`
}

type inboundEndpointIpConfigModel struct {
	ZoneId    types.String `tfsdk:"zone_id"`
	CidrBlock types.String `tfsdk:"cidr_block"`
	Ip        types.String `tfsdk:"ip"`
	VSwitchId types.String `tfsdk:"vswitch_id"`
}

func (r *inboundEndpointResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = inboundEndpointTypeName
}

func (r *inboundEndpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"ip_configs": schema.SetNestedBlock{
				// SetNestedBlock implements schema.Block (not schema.Attribute), so it
				// lives in Blocks, not Attributes. fwadapt serves the Framework resource
				// through the SDKv2 provider at protocol v5, where nested attributes are
				// rejected but blocks carry nested Attributes fine, matching the SDKv2
				// outbound endpoint's SchemaSetNestedBlock. SetNestedBlock has no Required
				// field; SizeAtLeast(2) enforces the high-availability minimum.
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"zone_id": schema.StringAttribute{
							Required: true,
						},
						"cidr_block": schema.StringAttribute{
							Required: true,
						},
						"ip": schema.StringAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"vswitch_id": schema.StringAttribute{
							Required: true,
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(2),
					setvalidator.SizeAtMost(6),
				},
			},
		},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"endpoint_name": schema.StringAttribute{
				Required: true,
			},
			"security_group_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpc_region_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"spec_version": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"endpoint_type": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// ModifyPlan clears the service-assigned ip in ip_config elements when ip_configs is changing.
// UpdateInboundEndpoint re-assigns every IP from VSwitchId/CidrBlock (it rejects stale IPs with
// InboundEndpoint.Ipconfig.InvalidIp), so any ip value carried into plan from state by
// UseStateForUnknown would no longer match the actual IP after the API re-assigns it. The
// Framework's set correlation would then fail with "Provider produced inconsistent result after
// apply". Setting ip to unknown for all elements lets the Framework correlate them by
// vswitch_id/zone_id/cidr_block and accept whatever IP the API assigns. When ip_configs is
// unchanged, the prior IPs are preserved and no modification is needed.
func (r *inboundEndpointResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// On create there is no prior state; the ip is unknown in plan and assigned by the API.
	if req.State.Raw.IsNull() {
		return
	}
	// On destroy the plan is null; there is nothing to modify and attempting to read a null
	// plan into the model raises a Value Conversion Error, so return early.
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state inboundEndpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If ip_configs hasn't changed, UseStateForUnknown correctly preserves the service-assigned
	// ip values from state so they don't appear as planned changes.
	if plan.IpConfigs.Equal(state.IpConfigs) {
		return
	}

	// ip_configs is changing: clear all ip values to unknown so the Framework does not expect
	// specific IP values that the API will re-assign from the new VSwitchId/CidrBlock.
	ipConfigs, diags := inboundEndpointIpConfigsFromSet(ctx, plan.IpConfigs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	for i := range ipConfigs {
		ipConfigs[i].Ip = types.StringUnknown()
	}
	newSet, diags := inboundEndpointIpConfigsToSet(ctx, ipConfigs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.IpConfigs = newSet
	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
}

func (r *inboundEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan inboundEndpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ipConfigs, diags := inboundEndpointIpConfigsFromSet(ctx, plan.IpConfigs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// AddInboundEndpoint encodes IpConfig as repeatList with the same AzId/CidrBlock/Ip/VSwitchId
	// sub-fields as AddResolverEndpoint, so the encoding mirrors the SDKv2 outbound endpoint.
	request := map[string]interface{}{
		"Name":            plan.EndpointName.ValueString(),
		"VpcId":           plan.VpcId.ValueString(),
		"VpcRegionId":     plan.VpcRegionId.ValueString(),
		"SecurityGroupId": plan.SecurityGroupId.ValueString(),
		"SpecVersion":     plan.SpecVersion.ValueString(),
		"Lang":            "en",
	}
	for k, ic := range ipConfigs {
		request[fmt.Sprintf("IpConfig.%d.VSwitchId", k+1)] = ic.VSwitchId.ValueString()
		request[fmt.Sprintf("IpConfig.%d.AzId", k+1)] = ic.ZoneId.ValueString()
		request[fmt.Sprintf("IpConfig.%d.CidrBlock", k+1)] = ic.CidrBlock.ValueString()
		if !ic.Ip.IsNull() && !ic.Ip.IsUnknown() && ic.Ip.ValueString() != "" {
			request[fmt.Sprintf("IpConfig.%d.Ip", k+1)] = ic.Ip.ValueString()
		}
	}

	action := "AddInboundEndpoint"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.RetryContext(ctx, inboundEndpointTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("pvtz", "2018-01-01", action, nil, request, false)
		if callErr != nil {
			// ResolverEndpoint.Operation.ExceedsFrequency is the documented throttle for the
			// outbound endpoint; InboundEndpoint may reuse the same pattern, so both are retried.
			if errs.IsExpectedErrors(callErr, []string{"InboundEndpoint.Operation.ExceedsFrequency", "ResolverEndpoint.Operation.ExceedsFrequency"}) || helper.NeedRetry(callErr) {
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
			fmt.Sprintf("Creating %s: calling AddInboundEndpoint", inboundEndpointTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, inboundEndpointTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	id := fmt.Sprint(response["EndpointId"])
	pvtzService := Service{r.Meta()}
	if _, waitErr := pvtzService.inboundEndpointWaitForState(ctx, id, []string{"SUCCESS"}, inboundEndpointFailStates, inboundEndpointTimeout); waitErr != nil {
		// Keep the created ID in state so the endpoint stays trackable even when the status wait
		// fails; Terraform marks the resource as tainted.
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for SUCCESS status", inboundEndpointTypeName),
			errs.WrapErrorf(waitErr, errs.IdMsg, id).Error(),
		)
		return
	}

	object, err := pvtzService.DescribeInboundEndpointContext(ctx, id)
	if err != nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: calling DescribeInboundEndpoint", inboundEndpointTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeInboundEndpoint", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, d := inboundEndpointState(ctx, object)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Preserve the exact EndpointId returned by AddInboundEndpoint rather than object["Id"]
	// from the Describe response. AddInboundEndpoint returns EndpointId as a string, but
	// DescribeInboundEndpoint may return Id as a JSON number whose 18-digit value loses
	// precision through float64 when the SDK decodes into map[string]interface{}. Using the
	// corrupted object["Id"] would make subsequent Update/Delete calls fail with
	// InboundEndpoint.InvalidId. The outbound SDKv2 endpoint keeps d.Id() from Create for
	// the same reason and never overwrites it in Read.
	model.Id = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *inboundEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state inboundEndpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if id == "" {
		return
	}

	pvtzService := Service{r.Meta()}
	object, err := pvtzService.DescribeInboundEndpointContext(ctx, id)
	if err != nil {
		if errs.IsProviderNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeInboundEndpoint", inboundEndpointTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeInboundEndpoint", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, diags := inboundEndpointState(ctx, object)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Keep the original state Id instead of object["Id"] from the Describe response; see the
	// comment in Create for why object["Id"] cannot be trusted to round-trip the identifier.
	model.Id = state.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *inboundEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan inboundEndpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state inboundEndpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.Id.ValueString()
	request := map[string]interface{}{
		"EndpointId": id,
		"Lang":       "en",
	}
	update := false
	if plan.EndpointName != state.EndpointName {
		update = true
		request["Name"] = plan.EndpointName.ValueString()
	}
	// The service-assigned ip is Optional+Computed and carries the prior value in plan when
	// ip_configs is unchanged, so set equality distinguishes a real change from an unchanged
	// block; UpdateInboundEndpoint only accepts Name and IpConfig, which matches this scope.
	if !plan.IpConfigs.Equal(state.IpConfigs) {
		update = true
		ipConfigs, diags := inboundEndpointIpConfigsFromSet(ctx, plan.IpConfigs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		// The service-assigned ip is Optional+Computed. UseStateForUnknown carries the
		// prior ip value from state into plan for elements that persist across the change.
		// Unlike UpdateResolverEndpoint (outbound), UpdateInboundEndpoint rejects that stale
		// ip with InboundEndpoint.Ipconfig.InvalidIp, so ip is never sent on update; the API
		// re-assigns IPs from VSwitchId/CidrBlock exactly as in Create (where ip is omitted).
		for k, ic := range ipConfigs {
			request[fmt.Sprintf("IpConfig.%d.VSwitchId", k+1)] = ic.VSwitchId.ValueString()
			request[fmt.Sprintf("IpConfig.%d.AzId", k+1)] = ic.ZoneId.ValueString()
			request[fmt.Sprintf("IpConfig.%d.CidrBlock", k+1)] = ic.CidrBlock.ValueString()
		}
	}

	if update {
		action := "UpdateInboundEndpoint"
		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err := retry.RetryContext(ctx, inboundEndpointTimeout, func() *retry.RetryError {
			var callErr error
			response, callErr = r.Meta().RpcPost("pvtz", "2018-01-01", action, nil, request, false)
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
				fmt.Sprintf("Updating %s: calling UpdateInboundEndpoint", inboundEndpointTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
		pvtzService := Service{r.Meta()}
		if _, waitErr := pvtzService.inboundEndpointWaitForState(ctx, id, []string{"SUCCESS"}, inboundEndpointFailStates, inboundEndpointTimeout); waitErr != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: waiting for SUCCESS status", inboundEndpointTypeName),
				errs.WrapErrorf(waitErr, errs.IdMsg, id).Error(),
			)
			return
		}
	}

	pvtzService := Service{r.Meta()}
	object, err := pvtzService.DescribeInboundEndpointContext(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update: calling DescribeInboundEndpoint", inboundEndpointTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeInboundEndpoint", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	model, diags := inboundEndpointState(ctx, object)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Keep the plan Id instead of object["Id"] from the Describe response; see the comment
	// in Create for why object["Id"] cannot be trusted to round-trip the identifier.
	model.Id = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *inboundEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state inboundEndpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	action := "DeleteInboundEndpoint"
	request := map[string]interface{}{
		"EndpointId": id,
		"Lang":       "en",
	}
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.RetryContext(ctx, inboundEndpointTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("pvtz", "2018-01-01", action, nil, request, false)
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
		// The not-exists code is not declared in the published DeleteInboundEndpoint contract,
		// and the live API returns it with a spelling error ("InboundEndpoint.NotExsits" rather
		// than "InboundEndpoint.NotExists"); a typed 404 or either spelling means the endpoint
		// is already gone.
		if errs.IsTypedNotFound(err) || errs.IsExpectedErrors(err, []string{"InboundEndpoint.NotExists", "InboundEndpoint.NotExsits"}) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteInboundEndpoint", inboundEndpointTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
}

func (r *inboundEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// inboundEndpointIpConfigsFromSet unwraps a types.Set of ip_config objects into a slice of
// inboundEndpointIpConfigModel. ElementsAs validates the element type against the nested object
// schema.
func inboundEndpointIpConfigsFromSet(ctx context.Context, set types.Set) ([]inboundEndpointIpConfigModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	ipConfigs := make([]inboundEndpointIpConfigModel, 0, len(set.Elements()))
	d := set.ElementsAs(ctx, &ipConfigs, false)
	diags.Append(d...)
	return ipConfigs, diags
}

// inboundEndpointIpConfigAttrTypes returns the attribute types of the ip_config nested object,
// used to build types.Object values when mapping the response back to state.
func inboundEndpointIpConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"zone_id":    types.StringType,
		"cidr_block": types.StringType,
		"ip":         types.StringType,
		"vswitch_id": types.StringType,
	}
}

// inboundEndpointIpConfigsToSet builds a types.Set from a slice of ip_config models, mapping the
// DescribeInboundEndpoint IpConfigs array back to state.
func inboundEndpointIpConfigsToSet(ctx context.Context, ipConfigs []inboundEndpointIpConfigModel) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	objs := make([]attr.Value, 0, len(ipConfigs))
	attrTypes := inboundEndpointIpConfigAttrTypes()
	for _, ic := range ipConfigs {
		objVal, d := types.ObjectValueFrom(ctx, attrTypes, ic)
		diags.Append(d...)
		objs = append(objs, objVal)
	}
	set, d := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: attrTypes}, objs)
	diags.Append(d...)
	return set, diags
}

// inboundEndpointString maps a response value to a types.String, treating nil and empty values as
// null so a missing optional field does not overwrite configured state.
func inboundEndpointString(raw interface{}) types.String {
	if raw == nil {
		return types.StringNull()
	}
	if s, ok := raw.(string); ok {
		if s == "" {
			return types.StringNull()
		}
		return types.StringValue(s)
	}
	return types.StringValue(fmt.Sprint(raw))
}

// inboundEndpointState maps a DescribeInboundEndpoint response object to the resource model.
func inboundEndpointState(ctx context.Context, object map[string]interface{}) (inboundEndpointResourceModel, diag.Diagnostics) {
	ipConfigsSli := make([]inboundEndpointIpConfigModel, 0)
	if ipConfigs, ok := object["IpConfigs"].([]interface{}); ok {
		for _, icArg := range ipConfigs {
			ic, ok := icArg.(map[string]interface{})
			if !ok {
				continue
			}
			ipConfigsSli = append(ipConfigsSli, inboundEndpointIpConfigModel{
				ZoneId:    inboundEndpointString(ic["AzId"]),
				CidrBlock: inboundEndpointString(ic["CidrBlock"]),
				Ip:        inboundEndpointString(ic["Ip"]),
				VSwitchId: inboundEndpointString(ic["VSwitchId"]),
			})
		}
	}
	ipConfigsSet, diags := inboundEndpointIpConfigsToSet(ctx, ipConfigsSli)

	return inboundEndpointResourceModel{
		Id:              inboundEndpointString(object["Id"]),
		EndpointName:    inboundEndpointString(object["Name"]),
		IpConfigs:       ipConfigsSet,
		SecurityGroupId: inboundEndpointString(object["SecurityGroupId"]),
		VpcId:           inboundEndpointString(object["VpcId"]),
		VpcRegionId:     inboundEndpointString(object["VpcRegionId"]),
		SpecVersion:     inboundEndpointString(object["SpecVersion"]),
		Status:          inboundEndpointString(object["Status"]),
		EndpointType:    inboundEndpointString(object["EndpointType"]),
	}, diags
}
