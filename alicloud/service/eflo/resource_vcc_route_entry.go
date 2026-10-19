package eflo

import (
	"context"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const vccRouteEntryTypeName = "alicloud_eflo_vcc_route_entry"

// eflo product/version shared by the VccRouteEntry CRUD actions.
const (
	efloProduct = "eflo"
	efloVersion = "2022-05-30"
	idSeparator = ":"
)

var (
	_ resource.Resource                = (*vccRouteEntryResource)(nil)
	_ resource.ResourceWithImportState = (*vccRouteEntryResource)(nil)
)

func NewVccRouteEntryResource() resource.Resource {
	return &vccRouteEntryResource{}
}

type vccRouteEntryResource struct {
	fwadapt.ResourceBase
}

// vccRouteEntryModel is the Terraform state shape. The `id` attribute holds the
// composite identity "VccId:VccRouteEntryId": GetVccRouteEntry requires both, so a
// single id must encode them. Read parses it back; ImportState splits it explicitly.
type vccRouteEntryModel struct {
	Id                   types.String `tfsdk:"id"`
	VccId                types.String `tfsdk:"vcc_id"`
	DestinationCidrBlock types.String `tfsdk:"destination_cidr_block"`
	RegionId             types.String `tfsdk:"region_id"`
	VccRouteEntryId      types.String `tfsdk:"vcc_route_entry_id"`
	Status               types.String `tfsdk:"status"`
	RouteType            types.String `tfsdk:"route_type"`
	NextHopId            types.String `tfsdk:"next_hop_id"`
	NextHopType          types.String `tfsdk:"next_hop_type"`
	Message              types.String `tfsdk:"message"`
	GmtModified          types.String `tfsdk:"gmt_modified"`
	TenantId             types.String `tfsdk:"tenant_id"`
	ResourceGroupId      types.String `tfsdk:"resource_group_id"`
}

func (r *vccRouteEntryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = vccRouteEntryTypeName
}

func (r *vccRouteEntryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Eflo Vcc Route Entry resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, formatted as `<vcc_id>:<vcc_route_entry_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vcc_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Eflo VCC instance.",
				PlanModifiers: []planmodifier.String{
					// No Update API exists; changing identity-bearing fields forces replacement.
					stringplanmodifier.RequiresReplace(),
				},
			},
			"destination_cidr_block": schema.StringAttribute{
				Required:    true,
				Description: "The destination CIDR block of the route entry.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The region ID. Defaults to the provider region.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vcc_route_entry_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the route entry.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the route entry.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"route_type": schema.StringAttribute{
				Computed:    true,
				Description: "The route type, e.g. `System`, `BGP`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"next_hop_id": schema.StringAttribute{
				Computed:    true,
				Description: "The next hop instance ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"next_hop_type": schema.StringAttribute{
				Computed:    true,
				Description: "The next hop type, e.g. `VCC`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"message": schema.StringAttribute{
				Computed:    true,
				Description: "The message returned with the route entry.",
			},
			"gmt_modified": schema.StringAttribute{
				Computed:    true,
				Description: "The last modification time.",
			},
			"tenant_id": schema.StringAttribute{
				Computed:    true,
				Description: "The tenant ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_group_id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource group ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *vccRouteEntryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vccRouteEntryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	request := map[string]interface{}{
		"VccId":                plan.VccId.ValueString(),
		"DestinationCidrBlock": plan.DestinationCidrBlock.ValueString(),
	}
	regionId := regionOrDefault(client, plan.RegionId)
	request["RegionId"] = regionId

	response, err := client.RpcPost(efloProduct, efloVersion, "CreateVccRouteEntry", nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateVccRouteEntry", vccRouteEntryTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateVccRouteEntry", vccRouteEntryTypeName).Error(),
		)
		return
	}

	vccRouteEntryId, ok := contentString(response, "VccRouteEntryId")
	if !ok || vccRouteEntryId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateVccRouteEntry", vccRouteEntryTypeName),
			fmt.Sprintf("The response carried no Content.VccRouteEntryId: %v", response),
		)
		return
	}

	// The route entry is immediately queryable; read once to populate state.
	object, err := getVccRouteEntry(client, plan.VccId.ValueString(), vccRouteEntryId, regionId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetVccRouteEntry", vccRouteEntryTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetVccRouteEntry", vccRouteEntryTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateVccRouteEntry", vccRouteEntryTypeName),
			fmt.Sprintf("The route entry %s:%s is not visible via GetVccRouteEntry right after creation", plan.VccId.ValueString(), vccRouteEntryId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, vccRouteEntryState(object))...)
}

func (r *vccRouteEntryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vccRouteEntryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vccId, vccRouteEntryId, ok := parseCompositeId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing id", vccRouteEntryTypeName),
			fmt.Sprintf("Expected id in `<vcc_id>:<vcc_route_entry_id>` form, got %q", state.Id.ValueString()),
		)
		return
	}

	client := r.Client()
	object, err := getVccRouteEntry(client, vccId, vccRouteEntryId, regionOrDefault(client, state.RegionId))
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetVccRouteEntry", vccRouteEntryTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetVccRouteEntry", vccRouteEntryTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, vccRouteEntryState(object))...)
}

// Update implements resource.Resource. VccRouteEntry is immutable: there is no
// UpdateVccRouteEntry action, so every settable field is RequiresReplace and
// Terraform never reaches Update in practice. The method exists to satisfy the
// resource.Resource surface; if it is ever invoked it refreshes state from the API
// rather than silently dropping the plan.
func (r *vccRouteEntryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan vccRouteEntryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	vccId := plan.VccId.ValueString()
	vccRouteEntryId := plan.VccRouteEntryId.ValueString()
	if vccId == "" || vccRouteEntryId == "" {
		vccId, vccRouteEntryId, _ = parseCompositeId(plan.Id.ValueString())
	}
	if vccId == "" || vccRouteEntryId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: resolving identity", vccRouteEntryTypeName),
			fmt.Sprintf("Could not resolve vcc_id / vcc_route_entry_id from plan: %s", plan.Id.ValueString()),
		)
		return
	}

	object, err := getVccRouteEntry(client, vccId, vccRouteEntryId, regionOrDefault(client, plan.RegionId))
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling GetVccRouteEntry", vccRouteEntryTypeName),
			errs.WrapErrorf(err, "Updating %s: calling GetVccRouteEntry", vccRouteEntryTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, vccRouteEntryState(object))...)
}

// Delete implements resource.Resource. VccRouteEntry is immutable: there is no
// UpdateVccRouteEntry action, so every settable field is RequiresReplace.
func (r *vccRouteEntryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vccRouteEntryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	vccId, vccRouteEntryId, ok := parseCompositeId(state.Id.ValueString())
	if !ok {
		// Nothing to delete from an unparseable id; treat as already gone.
		return
	}
	request := map[string]interface{}{
		"VccId":           vccId,
		"VccRouteEntryId": vccRouteEntryId,
		"RegionId":        regionOrDefault(client, state.RegionId),
	}
	_, err := client.RpcPost(efloProduct, efloVersion, "DeleteVccRouteEntry", nil, request, true)
	if err != nil && !isVccRouteEntryNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteVccRouteEntry", vccRouteEntryTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteVccRouteEntry", vccRouteEntryTypeName).Error(),
		)
		return
	}
}

// ImportState accepts `<vcc_id>:<vcc_route_entry_id>` and seeds the id, vcc_id and
// vcc_route_entry_id attributes; the subsequent Read repopulates the rest.
func (r *vccRouteEntryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	vccId, vccRouteEntryId, ok := parseCompositeId(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Importing %s: parsing id", vccRouteEntryTypeName),
			fmt.Sprintf("Expected id in `<vcc_id>:<vcc_route_entry_id>` form, got %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &vccRouteEntryModel{
		Id:              types.StringValue(req.ID),
		VccId:           types.StringValue(vccId),
		VccRouteEntryId: types.StringValue(vccRouteEntryId),
	})...)
}

// getVccRouteEntry calls GetVccRouteEntry and returns the Content object, or nil when
// the route entry no longer exists. Eflo returns code "1003" for missing resources;
// errs.NotFoundError covers the wrapped 404/NotFound forms.
func getVccRouteEntry(client *connectivity.AliyunClient, vccId, vccRouteEntryId, regionId string) (map[string]interface{}, error) {
	response, err := client.RpcPost(efloProduct, efloVersion, "GetVccRouteEntry", nil, map[string]interface{}{
		"VccId":           vccId,
		"VccRouteEntryId": vccRouteEntryId,
		"RegionId":        regionId,
	}, true)
	if err != nil {
		if isVccRouteEntryNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	content, _ := response["Content"].(map[string]interface{})
	if content == nil {
		return nil, fmt.Errorf("the GetVccRouteEntry response carried no Content object: %v", response)
	}
	return content, nil
}

// vccRouteEntryState maps a GetVccRouteEntry Content object onto the resource model.
// Empty strings become null: optional/computed fields the API omits read back as
// absent rather than "", matching the IMS framework convention.
func vccRouteEntryState(object map[string]interface{}) vccRouteEntryModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	vccId := str("VccId")
	vccRouteEntryId := str("VccRouteEntryId")
	id := types.StringNull()
	if !vccId.IsNull() && !vccRouteEntryId.IsNull() {
		id = types.StringValue(vccId.ValueString() + idSeparator + vccRouteEntryId.ValueString())
	}
	return vccRouteEntryModel{
		Id:                   id,
		VccId:                vccId,
		DestinationCidrBlock: str("DestinationCidrBlock"),
		RegionId:             str("RegionId"),
		VccRouteEntryId:      vccRouteEntryId,
		Status:               str("Status"),
		RouteType:            str("RouteType"),
		NextHopId:            str("NextHopId"),
		NextHopType:          str("NextHopType"),
		Message:              str("Message"),
		GmtModified:          str("GmtModified"),
		TenantId:             str("TenantId"),
		ResourceGroupId:      str("ResourceGroupId"),
	}
}

// regionOrDefault returns the configured region_id or the provider region when unset.
func regionOrDefault(client *connectivity.AliyunClient, configured types.String) string {
	if !configured.IsNull() && configured.ValueString() != "" {
		return configured.ValueString()
	}
	return client.RegionId
}

// contentString reads a string field from the Content object of an RPC response.
func contentString(response map[string]interface{}, key string) (string, bool) {
	content, _ := response["Content"].(map[string]interface{})
	if content == nil {
		return "", false
	}
	v, ok := content[key].(string)
	return v, ok
}

// parseCompositeId splits "VccId:VccRouteEntryId" into its parts. Neither Eflo id
// format contains ':', so a two-way split is unambiguous.
func parseCompositeId(id string) (string, string, bool) {
	if id == "" {
		return "", "", false
	}
	parts := strings.SplitN(id, idSeparator, 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// isVccRouteEntryNotExist reports whether err means the route entry is gone. Eflo
// surfaces missing resources as SDK error code "1003"; errs.NotFoundError covers the
// wrapped 404 / "NotFound" / ResourceNotfound forms produced across SDKs.
func isVccRouteEntryNotExist(err error) bool {
	if err == nil {
		return false
	}
	if errs.NotFoundError(err) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "1003") ||
		strings.Contains(msg, "not found") ||
		strings.Contains(msg, "not exist")
}
