package ga

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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

const resourceAttachmentTypeName = "alicloud_ga_resource_attachment"

var (
	_ resource.Resource                = &resourceAttachmentResource{}
	_ resource.ResourceWithImportState = &resourceAttachmentResource{}
)

func NewResourceAttachmentResource() resource.Resource {
	return &resourceAttachmentResource{}
}

type resourceAttachmentResource struct {
	fwadapt.ResourceBase
}

type resourceAttachmentResourceModel struct {
	Id                         types.String `tfsdk:"id"`
	AcceleratorId              types.String `tfsdk:"accelerator_id"`
	AssociatedResourceType     types.String `tfsdk:"associated_resource_type"`
	AssociatedResourceRegionId types.String `tfsdk:"associated_resource_region_id"`
	AssociatedResourceId       types.String `tfsdk:"associated_resource_id"`
	AssociatedMode             types.String `tfsdk:"associated_mode"`
	State                      types.String `tfsdk:"state"`
}

func (r *resourceAttachmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = resourceAttachmentTypeName
}

func (r *resourceAttachmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Global Accelerator (GA) Resource Attachment resource. It associates a resource with a GA instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource attachment, in the format `<accelerator_id>:<associated_resource_type>:<associated_resource_region_id>:<associated_resource_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"accelerator_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Global Accelerator instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"associated_resource_type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the associated resource. Valid value: `WAF`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("WAF"),
				},
			},
			"associated_resource_region_id": schema.StringAttribute{
				Required:    true,
				Description: "The region ID of the associated resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"associated_resource_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the associated resource. For `WAF`, it is the ID of the Web Application Firewall instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"associated_mode": schema.StringAttribute{
				Computed:    true,
				Description: "The association mode. Valid values: `Associated` and `Managed`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The association state. Valid values: `active`, `binding`, `unBinding` and `wafRelease`.",
			},
		},
	}
}

func (r *resourceAttachmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceAttachmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	acceleratorId := plan.AcceleratorId.ValueString()
	resourceType := plan.AssociatedResourceType.ValueString()
	resourceRegionId := plan.AssociatedResourceRegionId.ValueString()
	resourceId := plan.AssociatedResourceId.ValueString()

	client := r.Client()
	request := map[string]interface{}{
		"AcceleratorId":              acceleratorId,
		"AssociatedResourceType":     resourceType,
		"AssociatedResourceRegionId": resourceRegionId,
		"AssociatedResourceId":       resourceId,
		"RegionId":                   client.RegionId,
		"ClientToken":                helper.BuildClientToken("AssociateResources"),
	}
	err := retry.RetryContext(ctx, resourceAttachmentCreateTimeout, func() *retry.RetryError {
		_, callErr := client.RpcPost("Ga", "2019-11-20", "AssociateResources", nil, request, true)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				if sleepErr := sleepWithContext(ctx, resourceAttachmentPollInterval); sleepErr != nil {
					return retry.NonRetryableError(sleepErr)
				}
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling AssociateResources", resourceAttachmentTypeName),
			errs.WrapErrorf(err, "Creating %s: calling AssociateResources", resourceAttachmentTypeName).Error(),
		)
		return
	}

	plan.Id = types.StringValue(resourceAttachmentId(acceleratorId, resourceType, resourceRegionId, resourceId))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	record, err := waitResourceAttachmentActive(ctx, client, acceleratorId, resourceType, resourceRegionId, resourceId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for the association to become active", resourceAttachmentTypeName),
			errs.WrapErrorf(err, "Creating %s: waiting for the association to become active", resourceAttachmentTypeName).Error(),
		)
		return
	}
	plan.AssociatedMode = stringValue(record["AssociatedMode"])
	plan.State = stringValue(record["State"])
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *resourceAttachmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceAttachmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	found, diags := r.refresh(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update cannot be reached: every configurable attribute requires replacement.
// It refreshes state so an unexpected entry still converges instead of failing.
func (r *resourceAttachmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state resourceAttachmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	found, diags := r.refresh(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *resourceAttachmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceAttachmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	acceleratorId, resourceType, resourceRegionId, resourceId, err := parseResourceAttachmentId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing the resource ID", resourceAttachmentTypeName),
			err.Error(),
		)
		return
	}

	client := r.Client()
	existing, err := getResourceAttachment(ctx, client, acceleratorId, resourceType, resourceRegionId, resourceId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling GetGlobalAcceleratorResources", resourceAttachmentTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling GetGlobalAcceleratorResources", resourceAttachmentTypeName).Error(),
		)
		return
	}
	if existing == nil {
		return
	}

	request := map[string]interface{}{
		"AcceleratorId":              acceleratorId,
		"AssociatedResourceType":     resourceType,
		"AssociatedResourceRegionId": resourceRegionId,
		"AssociatedResourceId":       resourceId,
		"RegionId":                   client.RegionId,
		"ClientToken":                helper.BuildClientToken("DisassociateResources"),
	}
	err = retry.RetryContext(ctx, resourceAttachmentDeleteTimeout, func() *retry.RetryError {
		_, callErr := client.RpcPost("Ga", "2019-11-20", "DisassociateResources", nil, request, true)
		if callErr != nil {
			if isAcceleratorNotExist(callErr) {
				return nil
			}
			if helper.NeedRetry(callErr) {
				if sleepErr := sleepWithContext(ctx, resourceAttachmentPollInterval); sleepErr != nil {
					return retry.NonRetryableError(sleepErr)
				}
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DisassociateResources", resourceAttachmentTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DisassociateResources", resourceAttachmentTypeName).Error(),
		)
		return
	}

	if err := waitResourceAttachmentDeleted(ctx, client, acceleratorId, resourceType, resourceRegionId, resourceId); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: waiting for the association to be removed", resourceAttachmentTypeName),
			errs.WrapErrorf(err, "Deleting %s: waiting for the association to be removed", resourceAttachmentTypeName).Error(),
		)
	}
}

func (r *resourceAttachmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// refresh reads the association described by state.Id into state. The boolean
// result is false when the association no longer exists.
func (r *resourceAttachmentResource) refresh(ctx context.Context, state *resourceAttachmentResourceModel) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	acceleratorId, resourceType, resourceRegionId, resourceId, err := parseResourceAttachmentId(state.Id.ValueString())
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Reading %s: parsing the resource ID", resourceAttachmentTypeName),
			err.Error(),
		)
		return false, diags
	}

	record, err := getResourceAttachment(ctx, r.Client(), acceleratorId, resourceType, resourceRegionId, resourceId)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Reading %s: calling GetGlobalAcceleratorResources", resourceAttachmentTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetGlobalAcceleratorResources", resourceAttachmentTypeName).Error(),
		)
		return false, diags
	}
	if record == nil {
		return false, diags
	}

	state.AcceleratorId = types.StringValue(acceleratorId)
	state.AssociatedResourceType = types.StringValue(resourceType)
	state.AssociatedResourceRegionId = types.StringValue(resourceRegionId)
	state.AssociatedResourceId = types.StringValue(resourceId)
	if v := stringValue(record["AssociatedMode"]); !v.IsNull() {
		state.AssociatedMode = v
	}
	if v := stringValue(record["State"]); !v.IsNull() {
		state.State = v
	}
	return true, diags
}

func stringValue(v interface{}) types.String {
	if s, ok := v.(string); ok && s != "" {
		return types.StringValue(s)
	}
	return types.StringNull()
}
