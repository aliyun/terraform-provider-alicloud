package privatezone

import (
	"context"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
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
)

const (
	inboundEndpointTypeName = "alicloud_pvtz_inbound_endpoint"
	// inboundEndpointTimeout mirrors the SDKv2 alicloud_pvtz_endpoint create/update default
	// timeouts (10 minutes each). This resource does not expose user-configurable operation
	// timeouts, so the default is applied as a fixed internal retry/wait budget. A user-facing
	// timeouts block would be added through terraform-plugin-framework-timeouts and kept in
	// state; the fixed budget below is not a substitute for that capability.
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
	AzId      types.String `tfsdk:"az_id"`
	CidrBlock types.String `tfsdk:"cidr_block"`
	Ip        types.String `tfsdk:"ip"`
	VSwitchId types.String `tfsdk:"vswitch_id"`
}

// inboundEndpointIpConfigsRequiredValidator rejects a missing ip_configs block. SetNestedBlock
// cannot declare Required and the size validators skip a null value, so without this a config
// that omits the block reaches Create without the contract-required IpConfig.
type inboundEndpointIpConfigsRequiredValidator struct{}

func (v inboundEndpointIpConfigsRequiredValidator) Description(context.Context) string {
	return "block must be configured with at least 2 elements"
}

func (v inboundEndpointIpConfigsRequiredValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v inboundEndpointIpConfigsRequiredValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() {
		resp.Diagnostics.AddAttributeError(req.Path, "Missing required block", "ip_configs must be configured with at least 2 elements")
	}
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
						"az_id": schema.StringAttribute{
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
					// SetNestedBlock has no Required field and the size validators skip a null
					// value, so enforce the contract-required block explicitly.
					inboundEndpointIpConfigsRequiredValidator{},
					// The PD openapi-pvtz-inbound-endpoint.json evidence records
					// ipConfig.maxItems=10 and notes that the field description's
					// "at least 2, at most 6" wording is copied from the outbound endpoint, so
					// the upper bound follows the structure. The minimum stays at 2 per the
					// high-availability requirement documented for both endpoint types.
					setvalidator.SizeAtLeast(2),
					setvalidator.SizeAtMost(10),
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

// ModifyPlan keeps an explicitly configured ip and clears only the service-assigned value in
// ip_config elements when ip_configs is changing. UpdateInboundEndpoint accepts IpConfig.Ip and
// documents that an omitted Ip is assigned by the system, so a user-configured ip must reach the
// update request and must not be replaced by an unknown planned value. Elements whose ip is not
// configured are set to unknown so the Framework correlates them by vswitch_id/az_id/cidr_block
// and accepts whatever IP the API assigns; without this the prior IP carried by
// UseStateForUnknown would no longer match the API result and produce an inconsistent result
// after apply. When ip_configs is unchanged and the resource is updated in place, the prior IPs
// are preserved and nothing is modified. A planned replacement counts as a create: the plan then
// carries the prior state for Optional+Computed values, so an omitted ip would otherwise be
// replayed onto the new endpoint instead of being assigned by the service.
func (r *inboundEndpointResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// On create there is no prior state; a configured ip is planned as-is and an omitted ip is
	// unknown, and Create sends only the configured value.
	if req.State.Raw.IsNull() {
		return
	}
	// On destroy the plan is null; there is nothing to modify and attempting to read a null
	// plan into the model raises a Value Conversion Error, so return early.
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state, config inboundEndpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If ip_configs hasn't changed and the resource is not being replaced, UseStateForUnknown
	// correctly preserves the service-assigned ip values from state so they don't appear as
	// planned changes.
	replacement := plan.SecurityGroupId != state.SecurityGroupId ||
		plan.VpcId != state.VpcId ||
		plan.VpcRegionId != state.VpcRegionId ||
		plan.SpecVersion != state.SpecVersion
	if plan.IpConfigs.Equal(state.IpConfigs) && !replacement {
		return
	}

	configured, diags := inboundEndpointIpConfigsFromSet(ctx, config.IpConfigs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	explicitIps := make(map[string]struct{}, len(configured))
	for _, ic := range configured {
		if !ic.Ip.IsNull() && !ic.Ip.IsUnknown() && ic.Ip.ValueString() != "" {
			explicitIps[inboundEndpointIpConfigKey(ic)] = struct{}{}
		}
	}

	ipConfigs, diags := inboundEndpointIpConfigsFromSet(ctx, plan.IpConfigs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	for i := range ipConfigs {
		if _, ok := explicitIps[inboundEndpointIpConfigKey(ipConfigs[i])]; !ok {
			ipConfigs[i].Ip = types.StringUnknown()
		}
	}
	newSet, diags := inboundEndpointIpConfigsToSet(ctx, ipConfigs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.IpConfigs = newSet
	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
}

// inboundEndpointCreateRequest builds the AddInboundEndpoint body. The ClientToken is created
// once here, before the retry loop starts, so every attempt of the same logical create reuses the
// same token: a replay after a lost response is an idempotent no-op instead of a second endpoint.
func inboundEndpointCreateRequest(plan inboundEndpointResourceModel, ipConfigs []inboundEndpointIpConfigModel) map[string]interface{} {
	action := "AddInboundEndpoint"
	// AddInboundEndpoint encodes IpConfig as repeatList with the same AzId/CidrBlock/Ip/VSwitchId
	// sub-fields as AddResolverEndpoint, so the encoding mirrors the SDKv2 outbound endpoint.
	request := map[string]interface{}{
		"Name":            plan.EndpointName.ValueString(),
		"VpcId":           plan.VpcId.ValueString(),
		"VpcRegionId":     plan.VpcRegionId.ValueString(),
		"SecurityGroupId": plan.SecurityGroupId.ValueString(),
		"SpecVersion":     plan.SpecVersion.ValueString(),
		"Lang":            "en",
		"ClientToken":     helper.BuildClientToken(action),
	}
	for k, v := range inboundEndpointIpConfigRequest(ipConfigs) {
		request[k] = v
	}
	return request
}

func (r *inboundEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan inboundEndpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Build the ip_config request from the configuration, not the plan: the plan carries prior
	// state for Optional+Computed values during a replacement, so an ip the user omitted would
	// otherwise be replayed onto the new endpoint instead of being assigned by the service.
	var config inboundEndpointResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ipConfigs, diags := inboundEndpointIpConfigsFromSet(ctx, config.IpConfigs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No published InboundEndpoint errorCodes table exists (PD E5), so the gate relies on
	// helper.NeedRetry (5xx/throttling/timeouts); per-code gates are added when reports identify them.
	pvtzService := Service{r.Meta()}
	response, err := pvtzService.CreateInboundEndpoint(ctx, inboundEndpointCreateRequest(plan, ipConfigs))
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling AddInboundEndpoint", inboundEndpointTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, inboundEndpointTypeName, "AddInboundEndpoint", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	endpointId, ok := response["EndpointId"].(string)
	if !ok || endpointId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: AddInboundEndpoint returned no EndpointId", inboundEndpointTypeName),
			"The response did not contain a non-empty string EndpointId, so the created endpoint cannot be tracked.",
		)
		return
	}
	id := endpointId
	// Persist the ID as soon as it is validated, before the status wait and the response mapping.
	// Every later failure path must leave the created endpoint trackable so it can be imported or
	// deleted; Terraform marks the resource as tainted when Create returns an error with an ID.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, waitErr := pvtzService.inboundEndpointWaitForState(ctx, id, []string{"SUCCESS"}, inboundEndpointFailStates, inboundEndpointTimeout); waitErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for SUCCESS status", inboundEndpointTypeName),
			errs.WrapErrorf(waitErr, errs.IdMsg, id).Error(),
		)
		return
	}

	object, err := pvtzService.DescribeInboundEndpointContext(ctx, id)
	if err != nil {
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

// inboundEndpointUpdateRequest builds the UpdateInboundEndpoint body from the plan and the prior
// state. UpdateInboundEndpoint accepts only Name and IpConfig, so a field that did not change is
// omitted and keeps its current service value. The service-assigned ip is Optional+Computed and
// carries the prior value in the plan when ip_configs is unchanged, so set equality distinguishes
// a real change from an unchanged block.
func inboundEndpointUpdateRequest(ctx context.Context, plan, state inboundEndpointResourceModel) (map[string]interface{}, bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	request := map[string]interface{}{
		"EndpointId": plan.Id.ValueString(),
		"Lang":       "en",
	}
	update := false
	if plan.EndpointName != state.EndpointName {
		update = true
		request["Name"] = plan.EndpointName.ValueString()
	}
	if !plan.IpConfigs.Equal(state.IpConfigs) {
		update = true
		ipConfigs, d := inboundEndpointIpConfigsFromSet(ctx, plan.IpConfigs)
		diags.Append(d...)
		if diags.HasError() {
			return request, false, diags
		}
		// ModifyPlan keeps a user-configured ip in the plan and marks an omitted
		// (service-assigned) ip unknown, so this sends explicit values and omits the rest;
		// UpdateInboundEndpoint assigns an IP when Ip is absent.
		for k, v := range inboundEndpointIpConfigRequest(ipConfigs) {
			request[k] = v
		}
	}
	return request, update, diags
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
	request, update, diags := inboundEndpointUpdateRequest(ctx, plan, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pvtzService := Service{r.Meta()}
	if update {
		_, err := pvtzService.UpdateInboundEndpoint(ctx, request)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling UpdateInboundEndpoint", inboundEndpointTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "UpdateInboundEndpoint", errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
		if _, waitErr := pvtzService.inboundEndpointWaitForState(ctx, id, []string{"SUCCESS"}, inboundEndpointFailStates, inboundEndpointTimeout); waitErr != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: waiting for SUCCESS status", inboundEndpointTypeName),
				errs.WrapErrorf(waitErr, errs.IdMsg, id).Error(),
			)
			return
		}
	}

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
	pvtzService := Service{r.Meta()}
	request := map[string]interface{}{
		"EndpointId":  id,
		"Lang":        "en",
		"ClientToken": helper.BuildClientToken("DeleteInboundEndpoint"),
	}
	_, err := pvtzService.DeleteInboundEndpoint(ctx, request)
	if err != nil {
		// A typed 404 or either NotExists spelling means the endpoint is already gone.
		if IsInboundEndpointNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteInboundEndpoint", inboundEndpointTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DeleteInboundEndpoint", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	// DeleteInboundEndpoint only accepts the request; the endpoint keeps releasing its listeners
	// and VSwitch bindings afterwards. Wait for the object to be observably gone so Terraform does
	// not release dependent resources while the endpoint still exists.
	if waitErr := pvtzService.inboundEndpointWaitForDeleted(ctx, id, inboundEndpointTimeout); waitErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: waiting for the endpoint to disappear", inboundEndpointTypeName),
			errs.WrapErrorf(waitErr, errs.IdMsg, id).Error(),
		)
		return
	}
}

func (r *inboundEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// inboundEndpointIpConfigKey identifies an ip_config element by its required fields.
func inboundEndpointIpConfigKey(ic inboundEndpointIpConfigModel) string {
	return fmt.Sprintf("%s|%s|%s", ic.AzId.ValueString(), ic.CidrBlock.ValueString(), ic.VSwitchId.ValueString())
}

// inboundEndpointIpConfigRequest encodes ip_config as IpConfig.N.* repeatList entries. A known,
// configured ip is sent; a null or unknown ip is omitted so the service assigns it (Create) or
// re-assigns it from VSwitchId/CidrBlock (Update).
func inboundEndpointIpConfigRequest(ipConfigs []inboundEndpointIpConfigModel) map[string]interface{} {
	request := make(map[string]interface{}, len(ipConfigs)*4)
	for k, ic := range ipConfigs {
		request[fmt.Sprintf("IpConfig.%d.VSwitchId", k+1)] = ic.VSwitchId.ValueString()
		request[fmt.Sprintf("IpConfig.%d.AzId", k+1)] = ic.AzId.ValueString()
		request[fmt.Sprintf("IpConfig.%d.CidrBlock", k+1)] = ic.CidrBlock.ValueString()
		if !ic.Ip.IsNull() && !ic.Ip.IsUnknown() && ic.Ip.ValueString() != "" {
			request[fmt.Sprintf("IpConfig.%d.Ip", k+1)] = ic.Ip.ValueString()
		}
	}
	return request
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
		"az_id":      types.StringType,
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

// inboundEndpointRequiredString reads a required IpConfig sub-field. A missing, empty or
// non-string value is an error: writing null instead would let a contract drift or a malformed
// response masquerade as a legal, partially populated configuration.
func inboundEndpointRequiredString(object map[string]interface{}, key string) (types.String, error) {
	raw, ok := object[key]
	if !ok || raw == nil {
		return types.StringNull(), fmt.Errorf("%s is missing from the response", key)
	}
	s, ok := raw.(string)
	if !ok {
		return types.StringNull(), fmt.Errorf("%s is %T, want a string", key, raw)
	}
	if s == "" {
		return types.StringNull(), fmt.Errorf("%s is empty", key)
	}
	return types.StringValue(s), nil
}

// inboundEndpointOptionalString reads an optional IpConfig sub-field. Absent and empty values map
// to null; a present non-string value is an error.
func inboundEndpointOptionalString(object map[string]interface{}, key string) (types.String, error) {
	raw, ok := object[key]
	if !ok || raw == nil {
		return types.StringNull(), nil
	}
	s, ok := raw.(string)
	if !ok {
		return types.StringNull(), fmt.Errorf("%s is %T, want a string", key, raw)
	}
	if s == "" {
		return types.StringNull(), nil
	}
	return types.StringValue(s), nil
}

// inboundEndpointState maps a DescribeInboundEndpoint response object to the resource model. A
// malformed IpConfigs payload is reported as a diagnostic instead of being turned into an empty or
// partially populated state: silently dropping listeners would make Terraform plan against a
// configuration the service never reported.
func inboundEndpointState(ctx context.Context, object map[string]interface{}) (inboundEndpointResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	model := inboundEndpointResourceModel{
		Id:              inboundEndpointString(object["Id"]),
		EndpointName:    inboundEndpointString(object["Name"]),
		SecurityGroupId: inboundEndpointString(object["SecurityGroupId"]),
		VpcId:           inboundEndpointString(object["VpcId"]),
		VpcRegionId:     inboundEndpointString(object["VpcRegionId"]),
		SpecVersion:     inboundEndpointString(object["SpecVersion"]),
		Status:          inboundEndpointString(object["Status"]),
		EndpointType:    inboundEndpointString(object["EndpointType"]),
	}

	raw, ok := object["IpConfigs"]
	if !ok || raw == nil {
		diags.AddError(
			fmt.Sprintf("Reading %s: DescribeInboundEndpoint returned no IpConfigs", inboundEndpointTypeName),
			"The response did not contain the IpConfigs array, so the endpoint listeners cannot be represented in state.",
		)
		return model, diags
	}
	ipConfigList, ok := raw.([]interface{})
	if !ok {
		diags.AddError(
			fmt.Sprintf("Reading %s: unexpected IpConfigs type", inboundEndpointTypeName),
			fmt.Sprintf("DescribeInboundEndpoint returned IpConfigs as %T, want an array.", raw),
		)
		return model, diags
	}
	if len(ipConfigList) == 0 {
		diags.AddError(
			fmt.Sprintf("Reading %s: empty IpConfigs array", inboundEndpointTypeName),
			"The contract requires at least two ip_config elements, so an empty array cannot be written to state.",
		)
		return model, diags
	}

	ipConfigs := make([]inboundEndpointIpConfigModel, 0, len(ipConfigList))
	for i, item := range ipConfigList {
		ic, ok := item.(map[string]interface{})
		if !ok {
			diags.AddError(
				fmt.Sprintf("Reading %s: unexpected ip_config element", inboundEndpointTypeName),
				fmt.Sprintf("IpConfigs[%d] is %T, want an object.", i, item),
			)
			return model, diags
		}
		azId, err := inboundEndpointRequiredString(ic, "AzId")
		if err != nil {
			diags.AddError(fmt.Sprintf("Reading %s: invalid ip_config element", inboundEndpointTypeName), fmt.Sprintf("IpConfigs[%d].AzId: %s", i, err))
			return model, diags
		}
		cidrBlock, err := inboundEndpointRequiredString(ic, "CidrBlock")
		if err != nil {
			diags.AddError(fmt.Sprintf("Reading %s: invalid ip_config element", inboundEndpointTypeName), fmt.Sprintf("IpConfigs[%d].CidrBlock: %s", i, err))
			return model, diags
		}
		vSwitchId, err := inboundEndpointRequiredString(ic, "VSwitchId")
		if err != nil {
			diags.AddError(fmt.Sprintf("Reading %s: invalid ip_config element", inboundEndpointTypeName), fmt.Sprintf("IpConfigs[%d].VSwitchId: %s", i, err))
			return model, diags
		}
		ip, err := inboundEndpointOptionalString(ic, "Ip")
		if err != nil {
			diags.AddError(fmt.Sprintf("Reading %s: invalid ip_config element", inboundEndpointTypeName), fmt.Sprintf("IpConfigs[%d].Ip: %s", i, err))
			return model, diags
		}
		ipConfigs = append(ipConfigs, inboundEndpointIpConfigModel{
			AzId:      azId,
			CidrBlock: cidrBlock,
			Ip:        ip,
			VSwitchId: vSwitchId,
		})
	}

	ipConfigsSet, d := inboundEndpointIpConfigsToSet(ctx, ipConfigs)
	diags.Append(d...)
	if diags.HasError() {
		return model, diags
	}
	model.IpConfigs = ipConfigsSet
	return model, diags
}
