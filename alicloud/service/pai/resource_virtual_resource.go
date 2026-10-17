package pai

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const virtualResourceTypeName = "alicloud_pai_virtual_resource"

// EAS POP product code and API version for the VirtualResource lifecycle.
const (
	easProductCode = "EAS"
	easAPIVersion  = "2021-07-01"
)

var (
	_ resource.Resource                = &virtualResource{}
	_ resource.ResourceWithImportState = &virtualResource{}
)

func NewVirtualResourceResource() resource.Resource {
	return &virtualResource{}
}

type virtualResource struct {
	fwadapt.ResourceBase
}

// virtualResourceModel maps the Terraform state onto the DescribeVirtualResource
// response. The Resources list is a full-replace payload: Create and Update send
// the configured list verbatim, Read refreshes it from the Describe response.
type virtualResourceModel struct {
	Id                          types.String               `tfsdk:"id"`
	VirtualResourceId           types.String               `tfsdk:"virtual_resource_id"`
	VirtualResourceName         types.String               `tfsdk:"virtual_resource_name"`
	DisableSpotProtectionPeriod types.Bool                 `tfsdk:"disable_spot_protection_period"`
	RegionId                    types.String               `tfsdk:"region_id"`
	CreateTime                  types.String               `tfsdk:"create_time"`
	Resources                   []virtualResourceItemModel `tfsdk:"resources"`
}

// virtualResourceItemModel maps one entry of the Resources array. Per the
// CreateVirtualResource/OpenAPI contract, exactly one of instance_type,
// resource_id or quota_id identifies the backing resource group; the provider
// does not enforce the "exactly one" rule at plan time, the EAS API does.
type virtualResourceItemModel struct {
	Region         types.String  `tfsdk:"region"`
	Priority       types.Int64   `tfsdk:"priority"`
	InstanceType   types.String  `tfsdk:"instance_type"`
	QuotaId        types.String  `tfsdk:"quota_id"`
	SpotPriceLimit types.Float64 `tfsdk:"spot_price_limit"`
	ResourceId     types.String  `tfsdk:"resource_id"`
}

func (r *virtualResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = virtualResourceTypeName
}

func (r *virtualResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI EAS Virtual Resource resource. " +
			"A virtual resource group organizes compute resources (a public resource group instance type, " +
			"an EAS dedicated resource group, or a Lingjun quota) for EAS service scheduling.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `virtual_resource_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"virtual_resource_id": schema.StringAttribute{
				Computed:    true,
				Description: "The first-level ID of the virtual resource group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"virtual_resource_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the virtual resource group.",
			},
			"disable_spot_protection_period": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Description: "Whether to disable the spot protection period for spot instances in the group. " +
					"Defaults to `false`.",
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the resource, mapped from the EAS ClusterId.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the virtual resource group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"resources": schema.ListNestedBlock{
				Description: "The list of resources contained in the virtual resource group. " +
					"At least one entry is required, and each entry must specify exactly one of " +
					"`instance_type`, `resource_id` or `quota_id` to identify the backing resource " +
					"group; the EAS API enforces both rules.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The region of the resource entry.",
						},
						"priority": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "The scheduling priority of the resource. A larger number means a higher priority.",
						},
						"instance_type": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Description: "The instance type of the public resource group, e.g. `ecs.s6-c1m2.xlarge`. " +
								"Exactly one of `instance_type`, `resource_id` and `quota_id` must be set.",
						},
						"quota_id": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Description: "The Lingjun resource quota ID, for Lingjun智算 resource groups. " +
								"Exactly one of `instance_type`, `resource_id` and `quota_id` must be set.",
						},
						"spot_price_limit": schema.Float64Attribute{
							Optional: true,
							Computed: true,
							Description: "The spot price upper limit for public resource group spot instances. " +
								"Leave unset to not use spot instances.",
						},
						"resource_id": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Description: "The EAS dedicated resource group ID. " +
								"Exactly one of `instance_type`, `resource_id` and `quota_id` must be set.",
						},
					},
				},
			},
		},
	}
}

func (r *virtualResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan virtualResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := virtualResourceBody(&plan)
	response, err := r.Client().RoaPost(easProductCode, easAPIVersion, "/api/v2/virtualresources", nil, nil, body, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateVirtualResource", virtualResourceTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateVirtualResource", virtualResourceTypeName).Error(),
		)
		return
	}

	virtualResourceId, _ := response["VirtualResourceId"].(string)
	if virtualResourceId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateVirtualResource", virtualResourceTypeName),
			fmt.Sprintf("The CreateVirtualResource response carried no VirtualResourceId: %v", response),
		)
		return
	}

	object, err := describeVirtualResource(r.Client(), virtualResourceId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeVirtualResource", virtualResourceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeVirtualResource", virtualResourceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateVirtualResource", virtualResourceTypeName),
			fmt.Sprintf("The virtual resource %s is not visible via DescribeVirtualResource right after creation", virtualResourceId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, virtualResourceState(object, r.Client().RegionId, &plan))...)
}

func (r *virtualResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state virtualResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := describeVirtualResource(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeVirtualResource", virtualResourceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeVirtualResource", virtualResourceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, virtualResourceState(object, r.Client().RegionId, &state))...)
}

func (r *virtualResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state virtualResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := r.Client().RegionId
	action := fmt.Sprintf("/api/v2/virtualresources/%s/%s", clusterId, state.Id.ValueString())
	query := map[string]*string{"ClusterId": &clusterId}

	body := virtualResourceBody(&plan)
	if _, err := r.Client().RoaPut(easProductCode, easAPIVersion, action, query, nil, body, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateVirtualResource", virtualResourceTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateVirtualResource", virtualResourceTypeName).Error(),
		)
		return
	}

	object, err := describeVirtualResource(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeVirtualResource", virtualResourceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeVirtualResource", virtualResourceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, virtualResourceState(object, r.Client().RegionId, &plan))...)
}

func (r *virtualResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state virtualResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := r.Client().RegionId
	action := fmt.Sprintf("/api/v2/virtualresources/%s/%s", clusterId, state.Id.ValueString())
	query := map[string]*string{"ClusterId": &clusterId}
	_, err := r.Client().RoaDelete(easProductCode, easAPIVersion, action, query, nil, nil, true)
	if err != nil && !isVirtualResourceNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteVirtualResource", virtualResourceTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteVirtualResource", virtualResourceTypeName).Error(),
		)
		return
	}
}

func (r *virtualResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// describeVirtualResource queries DescribeVirtualResource by VirtualResourceId
// and returns the resource object, or nil when it no longer exists. ClusterId
// (the EAS cluster) is the provider region, matching the existing alicloud_pai_service
// resource.
func describeVirtualResource(client *connectivity.AliyunClient, virtualResourceId string) (map[string]interface{}, error) {
	action := fmt.Sprintf("/api/v2/virtualresources/%s/%s", client.RegionId, virtualResourceId)
	query := map[string]*string{"ClusterId": &client.RegionId}
	response, err := client.RoaGet(easProductCode, easAPIVersion, action, query, nil, nil)
	if err != nil {
		if isVirtualResourceNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return response, nil
}

// virtualResourceBody builds the Create/Update request body from the plan. Only
// set (non-null) attributes are sent, so an omitted optional field keeps its
// server default rather than being forced to a zero value.
func virtualResourceBody(plan *virtualResourceModel) map[string]interface{} {
	resources := make([]map[string]interface{}, 0, len(plan.Resources))
	for i := range plan.Resources {
		item := plan.Resources[i]
		entry := map[string]interface{}{}
		if !item.Region.IsNull() {
			entry["Region"] = item.Region.ValueString()
		}
		if !item.Priority.IsNull() && !item.Priority.IsUnknown() {
			entry["Priority"] = item.Priority.ValueInt64()
		}
		if !item.InstanceType.IsNull() {
			entry["InstanceType"] = item.InstanceType.ValueString()
		}
		if !item.QuotaId.IsNull() {
			entry["QuotaId"] = item.QuotaId.ValueString()
		}
		if !item.SpotPriceLimit.IsNull() && !item.SpotPriceLimit.IsUnknown() {
			entry["SpotPriceLimit"] = item.SpotPriceLimit.ValueFloat64()
		}
		if !item.ResourceId.IsNull() {
			entry["ResourceId"] = item.ResourceId.ValueString()
		}
		resources = append(resources, entry)
	}
	body := map[string]interface{}{
		"VirtualResourceName": plan.VirtualResourceName.ValueString(),
		"Resources":           resources,
	}
	if !plan.DisableSpotProtectionPeriod.IsNull() && !plan.DisableSpotProtectionPeriod.IsUnknown() {
		body["DisableSpotProtectionPeriod"] = plan.DisableSpotProtectionPeriod.ValueBool()
	}
	return body
}

// virtualResourceState maps a DescribeVirtualResource response onto the resource
// model. Empty strings become null so omitted optional fields read back as absent.
// RegionId is not echoed by DescribeVirtualResource (it is the request-side
// ClusterId), so it is populated from the provider region.
//
// Some optional fields (Priority, SpotPriceLimit, DisableSpotProtectionPeriod)
// are accepted by the Create/Update API but not echoed by DescribeVirtualResource
// for every resource type — e.g. a non-spot public resource group entry does not
// round-trip SpotPriceLimit/Priority. To keep the resource convergent, when the
// response omits such a field the previously configured value (from the plan on
// Create/Update, or the prior state on Read) is preserved. base may be nil when
// the caller has no prior value to fall back to.
func virtualResourceState(object map[string]interface{}, regionId string, base *virtualResourceModel) virtualResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	boolVal := func(key string) types.Bool {
		if v, ok := object[key].(bool); ok {
			return types.BoolValue(v)
		}
		// The EAS API serializes booleans as JSON booleans; keep a string fallback
		// for the "true"/"false" form some ROA responses use.
		if v, ok := object[key].(string); ok {
			return types.BoolValue(strings.EqualFold(v, "true"))
		}
		return types.BoolNull()
	}

	region := types.StringNull()
	if regionId != "" {
		region = types.StringValue(regionId)
	}
	model := virtualResourceModel{
		Id:                          str("VirtualResourceId"),
		VirtualResourceId:           str("VirtualResourceId"),
		VirtualResourceName:         str("VirtualResourceName"),
		DisableSpotProtectionPeriod: boolVal("DisableSpotProtectionPeriod"),
		RegionId:                    region,
		CreateTime:                  str("CreateTime"),
	}
	if base != nil && model.VirtualResourceName.IsNull() {
		model.VirtualResourceName = base.VirtualResourceName
	}
	if base != nil && model.DisableSpotProtectionPeriod.IsNull() && !base.DisableSpotProtectionPeriod.IsNull() && !base.DisableSpotProtectionPeriod.IsUnknown() {
		model.DisableSpotProtectionPeriod = base.DisableSpotProtectionPeriod
	}

	var baseItems []virtualResourceItemModel
	if base != nil {
		baseItems = base.Resources
	}
	if raw, ok := object["Resources"].([]interface{}); ok && len(raw) > 0 {
		model.Resources = make([]virtualResourceItemModel, 0, len(raw))
		for i, r := range raw {
			item, ok := r.(map[string]interface{})
			if !ok {
				if i < len(baseItems) {
					model.Resources = append(model.Resources, baseItems[i])
				}
				continue
			}
			entry := virtualResourceItemModel{
				Region:         itemString(item, "Region"),
				Priority:       itemInt64(item, "Priority"),
				InstanceType:   itemString(item, "InstanceType"),
				QuotaId:        itemString(item, "QuotaId"),
				SpotPriceLimit: itemFloat64(item, "SpotPriceLimit"),
				ResourceId:     itemString(item, "ResourceId"),
			}
			if i < len(baseItems) {
				b := baseItems[i]
				if entry.Region.IsNull() && !b.Region.IsNull() && !b.Region.IsUnknown() {
					entry.Region = b.Region
				}
				if entry.Priority.IsNull() && !b.Priority.IsNull() && !b.Priority.IsUnknown() {
					entry.Priority = b.Priority
				}
				if entry.InstanceType.IsNull() && !b.InstanceType.IsNull() && !b.InstanceType.IsUnknown() {
					entry.InstanceType = b.InstanceType
				}
				if entry.QuotaId.IsNull() && !b.QuotaId.IsNull() && !b.QuotaId.IsUnknown() {
					entry.QuotaId = b.QuotaId
				}
				if entry.SpotPriceLimit.IsNull() && !b.SpotPriceLimit.IsNull() && !b.SpotPriceLimit.IsUnknown() {
					entry.SpotPriceLimit = b.SpotPriceLimit
				}
				if entry.ResourceId.IsNull() && !b.ResourceId.IsNull() && !b.ResourceId.IsUnknown() {
					entry.ResourceId = b.ResourceId
				}
			}
			model.Resources = append(model.Resources, entry)
		}
	} else if base != nil {
		// The response carried no Resources list; keep the configured list so a
		// transiently-empty echo does not wipe the user's block configuration.
		model.Resources = baseItems
	}
	return model
}

func itemString(item map[string]interface{}, key string) types.String {
	if v, ok := item[key].(string); ok && v != "" {
		return types.StringValue(v)
	}
	return types.StringNull()
}

func itemInt64(item map[string]interface{}, key string) types.Int64 {
	switch v := item[key].(type) {
	case float64:
		return types.Int64Value(int64(v))
	case int:
		return types.Int64Value(int64(v))
	case int64:
		return types.Int64Value(v)
	case string:
		// EAS ROA serializes some numbers as strings ("6"); parse leniently.
		s := strings.TrimSpace(v)
		if s == "" {
			return types.Int64Null()
		}
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			return types.Int64Value(n)
		}
		// A decimal form like "6.0" is also accepted for an integer field.
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return types.Int64Value(int64(f))
		}
	}
	return types.Int64Null()
}

func itemFloat64(item map[string]interface{}, key string) types.Float64 {
	switch v := item[key].(type) {
	case float64:
		return types.Float64Value(v)
	case int:
		return types.Float64Value(float64(v))
	case int64:
		return types.Float64Value(float64(v))
	case string:
		// EAS ROA serializes some numbers as strings ("10.5"); parse leniently.
		s := strings.TrimSpace(v)
		if s == "" {
			return types.Float64Null()
		}
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return types.Float64Value(f)
		}
	}
	return types.Float64Null()
}

// isVirtualResourceNotExist reports whether err is the EAS virtual-resource
// not-found error. The code arrives on a *tea.SDKError from the ROA line; the
// substring fallback covers wrapped forms.
func isVirtualResourceNotExist(err error) bool {
	if err == nil {
		return false
	}
	if errs.NotFoundError(err) {
		return true
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil && *e.Code == "InvalidResource.NotFound" {
		return true
	}
	return strings.Contains(err.Error(), "InvalidResource.NotFound")
}
