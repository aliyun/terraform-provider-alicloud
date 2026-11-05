package sslcertificateservice

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
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

const sslCertificatesServiceWarehouseTypeName = "alicloud_ssl_certificates_service_warehouse"

var (
	_ resource.Resource                = &sslCertificatesServiceWarehouseResource{}
	_ resource.ResourceWithImportState = &sslCertificatesServiceWarehouseResource{}
)

// NewSslCertificatesServiceWarehouseResource is the Framework factory for the Warehouse resource.
// The Warehouse resource is new (no SDKv2 counterpart) and is registered through ServicePackage.Resources.
func NewSslCertificatesServiceWarehouseResource() resource.Resource {
	return &sslCertificatesServiceWarehouseResource{}
}

type sslCertificatesServiceWarehouseResource struct {
	fwadapt.ResourceBase
}

type sslCertificatesServiceWarehouseResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	WarehouseInstanceId types.String `tfsdk:"warehouse_instance_id"`
	Name                types.String `tfsdk:"name"`
	Type                types.String `tfsdk:"type"`
	Biz                 types.String `tfsdk:"biz"`
}

func (r *sslCertificatesServiceWarehouseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = sslCertificatesServiceWarehouseTypeName
}

func (r *sslCertificatesServiceWarehouseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"warehouse_instance_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("pcaAuthorityCertificate", "pcaCertificate", "personal", "ssl", "uploadedAuthorityCertificate", "uploadedCertificate"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"biz": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *sslCertificatesServiceWarehouseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sslCertificatesServiceWarehouseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := NewSslCertificatesServiceServiceV2(r.Meta())
	id, err := client.CreateSslCertificatesServiceWarehouse(plan.Name.ValueString(), plan.Type.ValueString(), plan.Biz.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateWarehouse", sslCertificatesServiceWarehouseTypeName),
			err.Error(),
		)
		return
	}

	// Read back so the computed warehouse_instance_id and the refreshed name populate state.
	object, err := client.DescribeSslCertificatesServiceWarehouse(id)
	if err != nil {
		// Keep the allocated identity so the resource is tainted for the next apply instead of lost.
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("warehouse_instance_id"), types.StringValue(id))...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading back after create", sslCertificatesServiceWarehouseTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, sslCertificatesServiceWarehouseState(object, &plan))...)
}

func (r *sslCertificatesServiceWarehouseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sslCertificatesServiceWarehouseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if id == "" {
		id = state.WarehouseInstanceId.ValueString()
	}
	if id == "" {
		return
	}

	client := NewSslCertificatesServiceServiceV2(r.Meta())
	object, err := client.DescribeSslCertificatesServiceWarehouse(id)
	if err != nil {
		if errs.IsProviderNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListWarehouse", sslCertificatesServiceWarehouseTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, sslCertificatesServiceWarehouseState(object, &state))...)
}

// Update is required by the Framework resource interface. Every configurable argument is
// replacement-only (no UpdateWarehouse API exists), so Terraform replaces the resource instead of
// calling this method; it carries the new plan into state for a no-op in-place change such as a
// timeouts-only update.
func (r *sslCertificatesServiceWarehouseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sslCertificatesServiceWarehouseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *sslCertificatesServiceWarehouseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sslCertificatesServiceWarehouseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if id == "" {
		id = state.WarehouseInstanceId.ValueString()
	}
	client := NewSslCertificatesServiceServiceV2(r.Meta())
	if err := client.DeleteSslCertificatesServiceWarehouse(id); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteWarehouse", sslCertificatesServiceWarehouseTypeName),
			err.Error(),
		)
		return
	}
}

func (r *sslCertificatesServiceWarehouseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// sslCertificatesServiceWarehouseState maps a ListWarehouse item onto the resource model. The id
// and warehouse_instance_id come from WarehouseInstanceId; name is refreshed from WarehouseName.
// type is retained from prior state because ListWarehouse returns WarehouseType using a different
// enum value domain than CreateWarehouse Type (e.g. Create pcaCertificate -> List pcaCaCert), and
// refreshing it would create a perpetual diff; type is replacement-only. biz is write-only
// (CreateWarehouse input, absent from ListWarehouse) and is also retained from prior state.
func sslCertificatesServiceWarehouseState(object map[string]interface{}, prior *sslCertificatesServiceWarehouseResourceModel) sslCertificatesServiceWarehouseResourceModel {
	m := sslCertificatesServiceWarehouseResourceModel{
		Id:                  types.StringNull(),
		WarehouseInstanceId: types.StringNull(),
		Name:                types.StringNull(),
		Type:                types.StringNull(),
		Biz:                 types.StringNull(),
	}
	if wid, ok := object["WarehouseInstanceId"].(string); ok && wid != "" {
		m.Id = types.StringValue(wid)
		m.WarehouseInstanceId = types.StringValue(wid)
	}
	if name, ok := object["WarehouseName"].(string); ok {
		m.Name = types.StringValue(name)
	} else if prior != nil {
		m.Name = prior.Name
	}
	if prior != nil {
		m.Type = prior.Type
		m.Biz = prior.Biz
	}
	return m
}
