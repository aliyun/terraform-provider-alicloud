package ims

import (
	"context"
	"fmt"
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

const mfaDeviceTypeName = "alicloud_ram_mfa_device"

var (
	_ resource.Resource                = &mfaDeviceResource{}
	_ resource.ResourceWithImportState = &mfaDeviceResource{}
)

func NewMFADeviceResource() resource.Resource {
	return &mfaDeviceResource{}
}

type mfaDeviceResource struct {
	fwadapt.ResourceBase
}

type mfaDeviceResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	VirtualMFADeviceName types.String `tfsdk:"virtual_mfa_device_name"`
	ActivateDate         types.String `tfsdk:"activate_date"`
}

func (r *mfaDeviceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = mfaDeviceTypeName
}

func (r *mfaDeviceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a RAM virtual MFA Device resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, the serial number of the virtual MFA device.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"virtual_mfa_device_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the virtual MFA device. The name must be 1 to 64 characters in length and can contain letters, digits, periods (.) and hyphens (-).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"activate_date": schema.StringAttribute{
				Computed:    true,
				Description: "The activation time of the MFA device. This value is empty if the device has not been activated.",
			},
		},
	}
}

func (r *mfaDeviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan mfaDeviceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	const action = "CreateVirtualMFADevice"
	response, err := r.Client().RpcPost("Ims", "2019-08-15", action, nil, map[string]interface{}{
		"VirtualMFADeviceName": plan.VirtualMFADeviceName.ValueString(),
	}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", mfaDeviceTypeName, action),
			errs.WrapErrorf(err, "Creating %s: calling %s", mfaDeviceTypeName, action).Error(),
		)
		return
	}

	serialNumber := serialNumberFromCreateResponse(response)
	if serialNumber == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", mfaDeviceTypeName, action),
			fmt.Sprintf("The response carried no VirtualMFADevice.SerialNumber: %v", response),
		)
		return
	}

	object, err := getMFADevice(r.Client(), serialNumber)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListVirtualMFADevices", mfaDeviceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling ListVirtualMFADevices", mfaDeviceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", mfaDeviceTypeName, action),
			fmt.Sprintf("The virtual MFA device %s is not visible via ListVirtualMFADevices right after creation", serialNumber),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, mfaDeviceState(object))...)
}

func (r *mfaDeviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state mfaDeviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getMFADevice(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListVirtualMFADevices", mfaDeviceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling ListVirtualMFADevices", mfaDeviceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, mfaDeviceState(object))...)
}

// Update is required to satisfy the resource.Resource interface. It is never
// reached in practice: virtual_mfa_device_name carries RequiresReplace, and the
// remaining attributes are Computed, so any change forces replacement rather
// than an in-place update. The plan is passed through unchanged.
func (r *mfaDeviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan mfaDeviceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *mfaDeviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state mfaDeviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	const action = "DeleteVirtualMFADevice"
	_, err := r.Client().RpcPost("Ims", "2019-08-15", action, nil, map[string]interface{}{
		"SerialNumber": state.Id.ValueString(),
	}, true)
	if err != nil && !isMFADeviceNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling %s", mfaDeviceTypeName, action),
			errs.WrapErrorf(err, "Deleting %s: calling %s", mfaDeviceTypeName, action).Error(),
		)
		return
	}
}

func (r *mfaDeviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mfaDeviceStr reads a string field from a device map; absent or empty values
// become null, matching the IMS convention of treating absent attributes as
// unset rather than "".
func mfaDeviceStr(device map[string]interface{}, key string) types.String {
	if v, ok := device[key].(string); ok && v != "" {
		return types.StringValue(v)
	}
	return types.StringNull()
}

// mfaDeviceState maps a ListVirtualMFADevices entry onto the resource model.
// The API does not return VirtualMFADeviceName directly; it is the last path
// segment of the SerialNumber, e.g. acs:ram::<accountId>:mfa/<device-name>.
func mfaDeviceState(device map[string]interface{}) mfaDeviceResourceModel {
	serialNumber := fmt.Sprint(device["SerialNumber"])
	return mfaDeviceResourceModel{
		Id:                   types.StringValue(serialNumber),
		VirtualMFADeviceName: types.StringValue(mfaDeviceName(serialNumber)),
		ActivateDate:         mfaDeviceStr(device, "ActivateDate"),
	}
}

// serialNumberFromCreateResponse pulls VirtualMFADevice.SerialNumber out of a
// CreateVirtualMFADevice response.
func serialNumberFromCreateResponse(response map[string]interface{}) string {
	if v, ok := response["VirtualMFADevice"].(map[string]interface{}); ok {
		return fmt.Sprint(v["SerialNumber"])
	}
	return ""
}

// isMFADeviceNotExist reports whether err is the IMS virtual-MFA-device
// not-found error. The code arrives on a *tea.SDKError from RpcPost; the
// substring fallback covers wrapped forms.
func isMFADeviceNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil && *e.Code == "EntityNotExist.VirtualMFADevice" {
		return true
	}
	return strings.Contains(err.Error(), "EntityNotExist.VirtualMFADevice")
}

// getMFADevice lists virtual MFA devices and returns the one whose SerialNumber
// matches, or nil when it no longer exists.
func getMFADevice(client *connectivity.AliyunClient, serialNumber string) (map[string]interface{}, error) {
	devices, err := listVirtualMFADevices(client)
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		if fmt.Sprint(device["SerialNumber"]) == serialNumber {
			return device, nil
		}
	}
	return nil, nil
}
