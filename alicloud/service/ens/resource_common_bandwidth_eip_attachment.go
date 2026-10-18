package ens

import (
	"context"
	"fmt"
	"strings"
	"time"

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

const commonBandwidthEipAttachmentTypeName = "alicloud_ens_common_bandwidth_eip_attachment"

// The backing API only exposes attach/detach semantics
// (AddCommonBandwidthPackageIps / DeleteCommonBandwidthPackageIps), so the
// resource is create-only and removed on delete; there is no in-place update.
var (
	_ resource.Resource                = &commonBandwidthEipAttachmentResource{}
	_ resource.ResourceWithImportState = &commonBandwidthEipAttachmentResource{}
)

func NewCommonBandwidthEipAttachmentResource() resource.Resource {
	return &commonBandwidthEipAttachmentResource{}
}

type commonBandwidthEipAttachmentResource struct {
	fwadapt.ResourceBase
}

type commonBandwidthEipAttachmentModel struct {
	Id                 types.String `tfsdk:"id"`
	BandwidthPackageId types.String `tfsdk:"bandwidth_package_id"`
	IpInstanceId       types.String `tfsdk:"ip_instance_id"`
}

func (r *commonBandwidthEipAttachmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = commonBandwidthEipAttachmentTypeName
}

func (r *commonBandwidthEipAttachmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a ENS Common Bandwidth Eip Attachment resource. " +
			"Bindings between ENS shared bandwidth and EIP.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. The value is formulated as `<ip_instance_id>:<bandwidth_package_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bandwidth_package_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the ENS common bandwidth package.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ip_instance_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the EIP (AllocationId) that is attached to the common bandwidth package.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *commonBandwidthEipAttachmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan commonBandwidthEipAttachmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	bandwidthPackageId := plan.BandwidthPackageId.ValueString()
	ipInstanceId := plan.IpInstanceId.ValueString()

	request := map[string]interface{}{
		"BandwidthPackageId": bandwidthPackageId,
		"IpInstanceIds.1":    ipInstanceId,
	}
	if _, err := client.RpcPost("Ens", "2017-11-10", "AddCommonBandwidthPackageIps", nil, request, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling AddCommonBandwidthPackageIps", commonBandwidthEipAttachmentTypeName),
			errs.WrapErrorf(err, "Creating %s: calling AddCommonBandwidthPackageIps", commonBandwidthEipAttachmentTypeName).Error(),
		)
		return
	}

	// AddCommonBandwidthPackageIps returns no object; the composite id is the
	// only handle Delete/Read can split back into its parts.
	plan.Id = types.StringValue(fmt.Sprintf("%s:%s", ipInstanceId, bandwidthPackageId))

	// The attachment is not always visible through DescribeEnsEipAddresses the
	// instant the attach call returns, so read back with a short bounded retry.
	// Falling back to the plan values keeps state populated when the describe
	// is still converging; the next refresh reconciles from the API.
	object := describeAttachmentWithRetry(ctx, client, plan.Id.ValueString())
	if object != nil {
		plan = attachmentState(object, plan)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *commonBandwidthEipAttachmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state commonBandwidthEipAttachmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := describeAttachment(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeEnsEipAddresses", commonBandwidthEipAttachmentTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeEnsEipAddresses", commonBandwidthEipAttachmentTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state = attachmentState(object, state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is required by the resource.Resource interface even though every
// argument is ForceNew: an attribute change routes to Delete+Create via
// RequiresReplace and never reaches here. The method reconciles any unknown
// Computed values by re-reading the attachment from the API; the backing API
// exposes no in-place update, so nothing is sent.
func (r *commonBandwidthEipAttachmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state commonBandwidthEipAttachmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := describeAttachment(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling DescribeEnsEipAddresses", commonBandwidthEipAttachmentTypeName),
			errs.WrapErrorf(err, "Updating %s: calling DescribeEnsEipAddresses", commonBandwidthEipAttachmentTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, attachmentState(object, state))...)
}

func (r *commonBandwidthEipAttachmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state commonBandwidthEipAttachmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parts := strings.Split(state.Id.ValueString(), ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing resource id", commonBandwidthEipAttachmentTypeName),
			fmt.Sprintf("invalid Resource Id %s. Expected 2 parts separated by ':', got %d", state.Id.ValueString(), len(parts)),
		)
		return
	}

	request := map[string]interface{}{
		"BandwidthPackageId": parts[1],
		"IpInstanceIds.1":    parts[0],
	}
	_, err := r.Client().RpcPost("Ens", "2017-11-10", "DeleteCommonBandwidthPackageIps", nil, request, true)
	if err != nil && !isAttachmentGone(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteCommonBandwidthPackageIps", commonBandwidthEipAttachmentTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteCommonBandwidthPackageIps", commonBandwidthEipAttachmentTypeName).Error(),
		)
		return
	}
}

func (r *commonBandwidthEipAttachmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// describeAttachment queries DescribeEnsEipAddresses by AllocationId and
// returns the EIP whose BandwidthPackageId matches the package half of id, or
// nil when the attachment no longer exists.
func describeAttachment(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid Resource Id %s. Expected 2 parts separated by ':', got %d", id, len(parts))
	}

	request := map[string]interface{}{
		"AllocationId": parts[0],
	}
	response, err := client.RpcPost("Ens", "2017-11-10", "DescribeEnsEipAddresses", nil, request, true)
	if err != nil {
		if isAttachmentGone(err) {
			return nil, nil
		}
		return nil, err
	}

	eips, _ := response["EipAddresses"].(map[string]interface{})
	list, _ := eips["EipAddress"].([]interface{})
	for _, v := range list {
		item, _ := v.(map[string]interface{})
		if fmt.Sprint(item["BandwidthPackageId"]) != parts[1] {
			continue
		}
		return item, nil
	}
	return nil, nil
}

// describeAttachmentWithRetry wraps describeAttachment with a short bounded
// retry for the create-time read-back, where the attachment is not yet visible.
func describeAttachmentWithRetry(ctx context.Context, client *connectivity.AliyunClient, id string) map[string]interface{} {
	// The context is honored between attempts; the API call itself is not
	// passed the context because RpcPost is not context-aware, matching the
	// rest of the provider's ENS service code.
	for attempt := 0; attempt < 5; attempt++ {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		object, err := describeAttachment(client, id)
		if err != nil || object != nil {
			return object
		}
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(3 * time.Second):
		}
	}
	return nil
}

// attachmentState maps a DescribeEnsEipAddresses EIP item onto the resource
// model, preserving the id and the user-supplied package id where the API
// response omits them.
func attachmentState(object map[string]interface{}, base commonBandwidthEipAttachmentModel) commonBandwidthEipAttachmentModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	out := base
	if v := str("BandwidthPackageId"); !v.IsNull() {
		out.BandwidthPackageId = v
	}
	if v := str("AllocationId"); !v.IsNull() {
		out.IpInstanceId = v
	}
	return out
}

// isAttachmentGone reports whether err means the attachment (or its EIP) is
// already gone, so Delete/Read can treat it as success. The code arrives on a
// *tea.SDKError from RpcPost; the substring fallback covers wrapped forms.
func isAttachmentGone(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil {
		switch *e.Code {
		case "IpInstanceId.NotInBandwidthPackage", "InvalidBandwidthPackageId.NotFound", "NotFound":
			return true
		}
	}
	msg := err.Error()
	return strings.Contains(msg, "IpInstanceId.NotInBandwidthPackage") ||
		strings.Contains(msg, "InvalidBandwidthPackageId.NotFound") ||
		strings.Contains(msg, "NotFound")
}
