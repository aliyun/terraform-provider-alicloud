// Package dataworks. The alicloud_dataworks_route_rule resource is a Framework-native,
// hand-written implementation of the dataworks-public 2018-06-01 RouteRule lifecycle.
// It is registered through ServicePackage.Resources (the Framework list), not the
// SDKv2 SDKResources list, and lives in the dataworks product package alongside the
// existing SDKv2 resources. The four target APIs are visibility=Private; development
// and ACC proceed, while external publication waits for the product owner to confirm
// the open/migration plan (a publish gate, not a development blocker).
package dataworks

import (
	"context"
	"fmt"
	"time"

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
)

const routeRuleTypeName = "alicloud_dataworks_route_rule"

// routeRuleOpTimeout bounds a single create/update/delete RPC; the underlying
// @returnMode is async:false, so no state-refresh wait follows the write.
const routeRuleOpTimeout = 5 * time.Minute

var (
	_ resource.Resource                = &routeRuleResource{}
	_ resource.ResourceWithImportState = &routeRuleResource{}
)

// NewRouteRuleResource is the Framework factory registered in register.go.
func NewRouteRuleResource() resource.Resource {
	return &routeRuleResource{}
}

type routeRuleResource struct {
	fwadapt.ResourceBase
}

// routeRuleResourceModel is the Terraform state shape. resource_group_id and
// resource_id are the DataWorks resource group and network resource the rule is
// attached to; both are Required and replace-only because UpdateRouteRule does not
// accept them. destination_addressing_type is Computed: ListRouteRules returns it
// but Create/Update do not accept it.
type routeRuleResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	RouteRuleId               types.String `tfsdk:"route_rule_id"`
	ResourceGroupId           types.String `tfsdk:"resource_group_id"`
	ResourceId                types.String `tfsdk:"resource_id"`
	DestinationCidr           types.String `tfsdk:"destination_cidr"`
	DestinationType           types.String `tfsdk:"destination_type"`
	DestinationAddressingType types.String `tfsdk:"destination_addressing_type"`
}

func (r *routeRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = routeRuleTypeName
}

func (r *routeRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The composite id of the route rule: <resource_group_id>:<resource_id>:<route_rule_id>. ListRouteRules is addressed by the resource group and network resource and has no single Get API, so the composite keeps Read self-contained.",
			},
			"route_rule_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The route rule id returned by CreateRouteRule and used to address Update/Delete.",
			},
			"resource_group_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The DataWorks resource group id the network resource (and its route rules) belong to. Replaces the resource on change because UpdateRouteRule does not accept it.",
			},
			"resource_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The network resource id (the alicloud_data_works_network id) the route rule is attached to. Replaces the resource on change because UpdateRouteRule does not accept it.",
			},
			"destination_cidr": schema.StringAttribute{
				Required:    true,
				Description: "The destination CIDR of the route rule, e.g. 192.168.10.0/24.",
			},
			"destination_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("VPC", "IDC"),
				},
				Description: "The destination type of the route rule. Valid values: VPC, IDC.",
			},
			"destination_addressing_type": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The destination addressing type returned by ListRouteRules. Valid values: CIDR, IP. Computed because Create/Update do not accept it.",
			},
		},
	}
}

func (r *routeRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan routeRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"ResourceGroupId": plan.ResourceGroupId.ValueString(),
		"ResourceId":      plan.ResourceId.ValueString(),
		"DestinationCidr": plan.DestinationCidr.ValueString(),
		"DestinationType": plan.DestinationType.ValueString(),
		"RegionId":        r.Meta().RegionId,
		"ClientToken":     helper.BuildClientToken(routeRuleCreateAPI),
	}

	opCtx, cancel := context.WithTimeout(ctx, routeRuleOpTimeout)
	defer cancel()
	response, err := routeRuleCallPost(opCtx, r.Meta(), routeRuleCreateAPI, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", routeRuleTypeName, routeRuleCreateAPI),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, routeRuleTypeName, routeRuleCreateAPI, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	routeRuleId, err := routeRuleIdFromResponse(routeRuleLookup(response, "data", "Data"))
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading route rule id from %s response", routeRuleTypeName, routeRuleCreateAPI),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, routeRuleTypeName, routeRuleCreateAPI, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	id := routeRuleCompositeId(plan.ResourceGroupId.ValueString(), plan.ResourceId.ValueString(), routeRuleId)

	object, err := NewService(r.Meta()).DescribeRouteRuleContext(ctx, id)
	if err != nil {
		// Persist the allocated identity so the rule stays trackable even when the
		// read-back fails; Terraform marks the resource as tainted.
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("route_rule_id"), types.StringValue(routeRuleId))...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: calling %s", routeRuleTypeName, routeRuleListAPI),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, routeRuleListAPI, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, routeRuleState(object, &plan))...)
}

func (r *routeRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state routeRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Id.ValueString() == "" {
		return
	}

	object, err := NewService(r.Meta()).DescribeRouteRuleContext(ctx, state.Id.ValueString())
	if err != nil {
		if errs.IsProviderNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", routeRuleTypeName, routeRuleListAPI),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), routeRuleListAPI, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, routeRuleState(object, &state))...)
}

// Update applies the only two fields UpdateRouteRule accepts: destination_cidr and
// destination_type. resource_group_id and resource_id are replacement-only, so a
// change to either replaces the resource instead of reaching this method.
func (r *routeRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state routeRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, routeRuleId, err := parseRouteRuleId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing id", routeRuleTypeName),
			err.Error(),
		)
		return
	}

	request := map[string]interface{}{
		"RouteRuleId":     routeRuleId,
		"DestinationCidr": plan.DestinationCidr.ValueString(),
		"DestinationType": plan.DestinationType.ValueString(),
		"RegionId":        r.Meta().RegionId,
		"ClientToken":     helper.BuildClientToken(routeRuleUpdateAPI),
	}

	opCtx, cancel := context.WithTimeout(ctx, routeRuleOpTimeout)
	defer cancel()
	if _, err := routeRuleCallPost(opCtx, r.Meta(), routeRuleUpdateAPI, request); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling %s", routeRuleTypeName, routeRuleUpdateAPI),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), routeRuleUpdateAPI, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeRouteRuleContext(ctx, state.Id.ValueString())
	if err != nil {
		if errs.IsProviderNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling %s after update", routeRuleTypeName, routeRuleListAPI),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), routeRuleListAPI, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, routeRuleState(object, &plan))...)
}

func (r *routeRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state routeRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opCtx, cancel := context.WithTimeout(ctx, routeRuleOpTimeout)
	defer cancel()
	if err := NewService(r.Meta()).DeleteRouteRuleContext(opCtx, state.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling %s", routeRuleTypeName, routeRuleDeleteAPI),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), routeRuleDeleteAPI, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
}

// ImportState accepts the composite id <resource_group_id>:<resource_id>:<route_rule_id>.
// The composite carries every identity Read needs because ListRouteRules has no single
// Get API; a bare route rule id cannot be re-read and is rejected by parseRouteRuleId.
func (r *routeRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// routeRuleState maps the Describe response object onto the resource model. Response
// fields are read by their attribute key (routeRuleNormalize fills them); the
// ForceNew values are preserved from prior state when the response omits them.
func routeRuleState(object map[string]interface{}, prior *routeRuleResourceModel) routeRuleResourceModel {
	str := func(key string) types.String {
		v, ok := object[key]
		if !ok || v == nil {
			return types.StringNull()
		}
		s := fmt.Sprint(v)
		if s == "" {
			return types.StringNull()
		}
		return types.StringValue(s)
	}
	m := routeRuleResourceModel{
		Id:                        str("Id"),
		RouteRuleId:               str("RouteRuleId"),
		ResourceGroupId:           str("ResourceGroupId"),
		ResourceId:                str("ResourceId"),
		DestinationCidr:           str("DestinationCidr"),
		DestinationType:           str("DestinationType"),
		DestinationAddressingType: str("DestinationAddressingType"),
	}
	// routeRuleNormalize always fills ResourceGroupId/ResourceId from the composite
	// id, so they are never missing here; keep the prior-guard pattern for parity with
	// the other Framework resources in case a future normalize variant omits a field.
	if prior != nil {
		if m.ResourceGroupId.IsNull() {
			m.ResourceGroupId = prior.ResourceGroupId
		}
		if m.ResourceId.IsNull() {
			m.ResourceId = prior.ResourceId
		}
	}
	return m
}
