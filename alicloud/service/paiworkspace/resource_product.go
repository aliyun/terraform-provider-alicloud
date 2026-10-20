package paiworkspace

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const productTypeName = "alicloud_pai_workspace_product"

var (
	_ resource.Resource                = &productResource{}
	_ resource.ResourceWithImportState = &productResource{}
)

func NewProductResource() resource.Resource {
	return &productResource{}
}

// productResource manages a PAI Workspace Product. Creating the resource places
// a real, chargeable order via CreateProductOrders; there is no Update or Delete
// API, so changes force replacement and Delete is a no-op that leaves the
// purchased product in the cloud. The resource is therefore order-only and
// immutable.
type productResource struct {
	fwadapt.ResourceBase
}

// instancePropertyModel is one InstanceProperties entry as a framework object.
type instancePropertyModel struct {
	Code  types.String `tfsdk:"code"`
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type productResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	ProductId               types.String `tfsdk:"product_id"`
	ProductCode             types.String `tfsdk:"product_code"`
	ChargeType              types.String `tfsdk:"charge_type"`
	PricingCycle            types.String `tfsdk:"pricing_cycle"`
	Duration                types.Int64  `tfsdk:"duration"`
	AutoRenew               types.Bool   `tfsdk:"auto_renew"`
	OrderType               types.String `tfsdk:"order_type"`
	AutoPay                 types.Bool   `tfsdk:"auto_pay"`
	InstanceProperties      types.List   `tfsdk:"instance_properties"`
	IsPurchased             types.Bool   `tfsdk:"is_purchased"`
	HasPermissionToPurchase types.Bool   `tfsdk:"has_permission_to_purchase"`
	PurchaseUrl             types.String `tfsdk:"purchase_url"`
	RegionId                types.String `tfsdk:"region_id"`
}

func (r *productResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = productTypeName
}

func (r *productResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI Workspace Product resource. Creating the resource places a real product order via " +
			"the `CreateProductOrders` API. There is no Update or Delete API: the order is irreversible, so any " +
			"change to a configurable attribute forces replacement (a new order) and `terraform destroy` only " +
			"removes the resource from Terraform state without cancelling the purchased product.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `product_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"product_id": schema.StringAttribute{
				Computed: true,
				Description: "The product ID, the unique key of the purchased product. It is the first element of the " +
					"`ProductIds` returned by `CreateProductOrders`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"product_code": schema.StringAttribute{
				Required: true,
				Description: "The code of the product to purchase. The valid value is `PAI_share`. " +
					"Changing it forces a new order (replacement).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"charge_type": schema.StringAttribute{
				Optional: true,
				Description: "The charge type of the order, such as `POSTPAY`, `PayAsYouGo` or `Subscription`. " +
					"Only settable at creation. Changing it forces a new order.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"pricing_cycle": schema.StringAttribute{
				Optional: true,
				Description: "The pricing cycle, such as `Hour`, `Month` or `Year`. Only settable at creation. " +
					"Changing it forces a new order.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"duration": schema.Int64Attribute{
				Optional: true,
				Description: "The purchase duration; the unit depends on the pricing cycle. Only settable at creation. " +
					"Changing it forces a new order.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"auto_renew": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to auto renew. Only settable at creation. Changing it forces a new order.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"order_type": schema.StringAttribute{
				Optional: true,
				Description: "The order type, such as `BUY`, `New` or `Renew`. Only settable at creation. " +
					"Changing it forces a new order.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"auto_pay": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to auto pay for all products. Only settable at creation. Changing it forces a new order.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"is_purchased": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the product has been purchased.",
			},
			"has_permission_to_purchase": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the caller has permission to purchase the product.",
			},
			"purchase_url": schema.StringAttribute{
				Computed:    true,
				Description: "The URL to purchase the product.",
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the resource.",
			},
		},
		// instance_properties is a block, not an attribute: nested-attribute
		// objects (schema.NestedAttributeObject with Attributes) are a protocol
		// v6 concept and cannot be served under the v5 mux this provider uses.
		// schema.ListNestedBlock maps to the same types.List in the state model,
		// so the create/plan wiring is unchanged. The block itself has no list-
		// level plan modifier; RequiresReplace on each leaf attribute forces a
		// new order whenever any property changes (including add/remove of a
		// whole entry, since the leaf goes null<->value).
		Blocks: map[string]schema.Block{
			"instance_properties": schema.ListNestedBlock{
				Description: "The instance properties attached to the order. Only settable at creation; the API does " +
					"not read them back, so they are preserved from configuration. Changing the block forces a new order.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"code": schema.StringAttribute{
							Optional:    true,
							Description: "The property code.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"name": schema.StringAttribute{
							Optional:    true,
							Description: "The property name.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"value": schema.StringAttribute{
							Optional:    true,
							Description: "The property value.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
			},
		},
	}
}

func (r *productResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan productResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := productCreateInput{
		ProductCode: plan.ProductCode.ValueString(),
		AutoRenew:   plan.AutoRenew.ValueBool(),
		AutoPay:     plan.AutoPay.ValueBool(),
	}
	if !plan.ChargeType.IsNull() {
		in.ChargeType = plan.ChargeType.ValueString()
	}
	if !plan.PricingCycle.IsNull() {
		in.PricingCycle = plan.PricingCycle.ValueString()
	}
	if !plan.OrderType.IsNull() {
		in.OrderType = plan.OrderType.ValueString()
	}
	if !plan.Duration.IsNull() {
		in.Duration = plan.Duration.ValueInt64()
	}
	// AutoRenew is already read above; it is a boolean, so a null plan value is
	// treated as the false default and is sent explicitly.
	if !plan.InstanceProperties.IsNull() {
		var ips []instancePropertyModel
		resp.Diagnostics.Append(plan.InstanceProperties.ElementsAs(ctx, &ips, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, ip := range ips {
			entry := instancePropertyInput{}
			if !ip.Code.IsNull() {
				entry.Code = ip.Code.ValueString()
			}
			if !ip.Name.IsNull() {
				entry.Name = ip.Name.ValueString()
			}
			if !ip.Value.IsNull() {
				entry.Value = ip.Value.ValueString()
			}
			in.InstanceProperties = append(in.InstanceProperties, entry)
		}
	}

	productId, err := createProduct(r.Client(), buildCreateBody(in))
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateProductOrders", productTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateProductOrders", productTypeName).Error(),
		)
		return
	}

	// Carry the plan into state so the create-only order attributes and the
	// instance_properties block survive refresh (they are not read back).
	state := plan
	state.Id = types.StringValue(productId)
	state.ProductId = types.StringValue(productId)

	// Best-effort read for the catalog status fields. A freshly placed order may
	// not be consistent in ListProducts yet; in that case we record the product
	// id and mark it purchased, leaving the rest for the next refresh.
	obj, readErr := readProduct(r.Client(), productId, plan.ProductCode.ValueString())
	switch {
	case readErr != nil:
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Reading %s: calling ListProducts after create", productTypeName),
			errs.WrapErrorf(readErr, "Reading %s: calling ListProducts after create", productTypeName).Error(),
		)
		state.IsPurchased = types.BoolValue(true)
	case obj == nil:
		state.IsPurchased = types.BoolValue(true)
	default:
		applyProductObject(&state, obj, r.Client().RegionId)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *productResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state productResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj, err := readProduct(r.Client(), state.ProductId.ValueString(), state.ProductCode.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListProducts", productTypeName),
			errs.WrapErrorf(err, "Reading %s: calling ListProducts", productTypeName).Error(),
		)
		return
	}
	if obj == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Overlay the readable catalog fields; the create-only order attributes and
	// instance_properties are preserved from prior state since the API does not
	// return them.
	applyProductObject(&state, obj, r.Client().RegionId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *productResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No Delete API exists. Orders are irreversible: the purchased product is
	// not cancelled, and Terraform removes the resource from state only.
}

func (r *productResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No Update API exists and every configurable attribute is RequiresReplace,
	// so the framework never routes a real diff here. Keep the plan as state so
	// the interface is satisfied and any in-place drift the framework does pass
	// (none, by schema design) is carried through without a spurious cloud call.
	var plan productResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *productResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// applyProductObject maps a ListProducts entry onto the readable attributes of
// the state model. It does not touch the create-only order attributes or
// instance_properties, which are preserved from the prior state so refresh does
// not drop values the API never returns.
func applyProductObject(state *productResourceModel, obj map[string]interface{}, regionId string) {
	str := func(key string) types.String {
		if v, ok := obj[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	boolVal := func(key string) types.Bool {
		if v, ok := obj[key].(bool); ok {
			return types.BoolValue(v)
		}
		return types.BoolNull()
	}

	state.ProductId = str("ProductId")
	state.Id = state.ProductId
	if pc := str("ProductCode"); !pc.IsNull() {
		state.ProductCode = pc
	}
	state.IsPurchased = boolVal("IsPurchased")
	state.HasPermissionToPurchase = boolVal("HasPermissionToPurchase")
	state.PurchaseUrl = str("PurchaseUrl")
	if r := str("RegionId"); !r.IsNull() {
		state.RegionId = r
	} else if regionId != "" {
		state.RegionId = types.StringValue(regionId)
	}
}
