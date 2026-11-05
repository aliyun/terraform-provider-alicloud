package ens

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const storageVolumeTypeName = "alicloud_ens_storage_volume"

// storageVolumeTimeout replaces the SDKv2 resource-level default timeouts
// (Create/Delete each 5 minutes); terraform-plugin-framework does not support
// user-configurable resource timeouts yet.
const (
	storageVolumeCreateTimeout = 5 * time.Minute
	storageVolumeDeleteTimeout = 5 * time.Minute
)

var (
	_ resource.Resource                = &storageVolumeResource{}
	_ resource.ResourceWithImportState = &storageVolumeResource{}
)

func NewEnsStorageVolumeResource() resource.Resource {
	return &storageVolumeResource{}
}

type storageVolumeResource struct {
	fwadapt.ResourceBase
}

type storageVolumeResourceModel struct {
	Id           types.String `tfsdk:"id"`
	GatewayId    types.String `tfsdk:"gateway_id"`
	StorageId    types.String `tfsdk:"storage_id"`
	EnsRegionId  types.String `tfsdk:"ens_region_id"`
	VolumeName   types.String `tfsdk:"volume_name"`
	Description  types.String `tfsdk:"description"`
	IsEnable     types.Bool   `tfsdk:"is_enable"`
	IsAuth       types.Bool   `tfsdk:"is_auth"`
	AuthProtocol types.String `tfsdk:"auth_protocol"`
	AuthUser     types.String `tfsdk:"auth_user"`
	AuthPassword types.String `tfsdk:"auth_password"`
	Status       types.String `tfsdk:"status"`
	CreationTime types.String `tfsdk:"creation_time"`
	TargetName   types.String `tfsdk:"target_name"`
}

func (r *storageVolumeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = storageVolumeTypeName
}

func (r *storageVolumeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// gateway_id, storage_id and ens_region_id identify the storage volume's
			// parent gateway and edge node; the ENS storage volume API exposes no
			// update operation, so every identifying argument forces replacement.
			"gateway_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The ID of the storage gateway that the storage volume belongs to. CreateStorageVolume accepts GatewayId, while DescribeStorageVolume returns it as StorageGatewayId.",
			},
			"storage_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The ID of the storage medium that backs the storage volume.",
			},
			"ens_region_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The ID of the ENS node (edge region) where the storage volume is created.",
			},
			"volume_name": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The name of the storage volume. CreateStorageVolume accepts VolumeName, while DescribeStorageVolume returns it as StorageVolumeName.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The description of the storage volume.",
			},
			// is_enable and is_auth are Optional + Computed because CreateStorageVolume
			// applies server defaults (IsEnable defaults to enabled, IsAuth to disabled)
			// when the argument is omitted, so the value read back must not drift a
			// planned null. The API exposes no update, so a configured change replaces.
			"is_enable": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
				Description: "Whether to enable the storage volume. Valid values: true (enabled, default), false (disabled). CreateStorageVolume encodes it as IsEnable (string \"1\"/\"0\"), while DescribeStorageVolume returns it as an integer.",
			},
			"is_auth": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
				Description: "Whether to enable CHAP authentication for the storage volume. Valid values: true (enabled), false (disabled, default). CreateStorageVolume encodes it as IsAuth (string \"1\"/\"0\"), while DescribeStorageVolume returns it as an integer.",
			},
			"auth_protocol": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("CHAP"),
				},
				Description: "The authentication protocol. Currently only CHAP is supported.",
			},
			// auth_user and auth_password are write-only: CreateStorageVolume accepts
			// them but DescribeStorageVolume never returns them, so refresh keeps the
			// previously configured value instead of collapsing it to null.
			"auth_user": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The CHAP authentication user name. Write-only: DescribeStorageVolume does not return it, so the configured value is preserved on refresh.",
			},
			"auth_password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The CHAP authentication password. Write-only and sensitive: DescribeStorageVolume does not return it, so the configured value is preserved on refresh.",
			},
			"status": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The status of the storage volume. Valid values: creating, available, deleting, deleted.",
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The creation time of the storage volume, in UTC format.",
			},
			"target_name": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The iSCSI target name of the storage volume.",
			},
		},
	}
}

func (r *storageVolumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan storageVolumeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "CreateStorageVolume"
	request := map[string]interface{}{
		"GatewayId":   plan.GatewayId.ValueString(),
		"StorageId":   plan.StorageId.ValueString(),
		"EnsRegionId": plan.EnsRegionId.ValueString(),
		"ClientToken": helper.BuildClientToken(action),
	}
	if !plan.VolumeName.IsNull() && plan.VolumeName.ValueString() != "" {
		request["VolumeName"] = plan.VolumeName.ValueString()
	}
	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		request["Description"] = plan.Description.ValueString()
	}
	// CreateStorageVolume encodes the booleans as string enum "1"/"0".
	if !plan.IsEnable.IsNull() {
		request["IsEnable"] = boolToStr(plan.IsEnable.ValueBool())
	}
	if !plan.IsAuth.IsNull() {
		request["IsAuth"] = boolToStr(plan.IsAuth.ValueBool())
	}
	if !plan.AuthProtocol.IsNull() && plan.AuthProtocol.ValueString() != "" {
		request["AuthProtocol"] = plan.AuthProtocol.ValueString()
	}
	if !plan.AuthUser.IsNull() && plan.AuthUser.ValueString() != "" {
		request["AuthUser"] = plan.AuthUser.ValueString()
	}
	if !plan.AuthPassword.IsNull() && plan.AuthPassword.ValueString() != "" {
		request["AuthPassword"] = plan.AuthPassword.ValueString()
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(storageVolumeCreateTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost("Ens", "2017-11-10", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateStorageVolume", storageVolumeTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateStorageVolume", storageVolumeTypeName).Error(),
		)
		return
	}

	// CreateStorageVolume returns VolumeId as an array of created volume IDs; a
	// single volume is created here, so the first element is the resource id
	// (StorageVolumeId, sv-***).
	id, err := firstStorageVolumeId(response)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading VolumeId from response", storageVolumeTypeName),
			err.Error(),
		)
		return
	}

	// The volume is created asynchronously; wait until it reaches the available
	// state before reading it back, so the computed status is not "creating".
	if err := storageVolumeWaitForState(r.Meta(), id, "available", []string{"deleted"}, storageVolumeCreateTimeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for available", storageVolumeTypeName),
			errs.WrapErrorf(err, "Creating %s: waiting for available", storageVolumeTypeName).Error(),
		)
		return
	}

	object, err := describeStorageVolume(r.Meta(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeStorageVolume", storageVolumeTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeStorageVolume", storageVolumeTypeName).Error(),
		)
		return
	}
	state := storageVolumeState(object, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *storageVolumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state storageVolumeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := describeStorageVolume(r.Meta(), state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeStorageVolume", storageVolumeTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeStorageVolume", storageVolumeTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, storageVolumeState(object, &state))...)
}

// Update is required by the resource.Resource interface but is never reached for
// a real change: every configurable argument is RequiresReplace, so Terraform
// recreates the volume instead of updating it. The only non-replace plan change
// is a no-op, so refresh reads the live volume back and carries the plan's
// write-only auth_user/auth_password forward.
func (r *storageVolumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state storageVolumeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := describeStorageVolume(r.Meta(), state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling DescribeStorageVolume", storageVolumeTypeName),
			errs.WrapErrorf(err, "Updating %s: calling DescribeStorageVolume", storageVolumeTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, storageVolumeState(object, &plan))...)
}

func (r *storageVolumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state storageVolumeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "DeleteStorageVolume"
	// DeleteStorageVolume takes VolumeId, which is the same value as the
	// DescribeStorageVolume primary key StorageVolumeId (sv-***).
	request := map[string]interface{}{
		"VolumeId": state.Id.ValueString(),
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(storageVolumeDeleteTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost("Ens", "2017-11-10", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		// A delete retry that reports the volume as already gone is success.
		if errs.IsExpectedErrors(err, []string{"InvalidStorageVolume.NotFound", "NotFound"}) || errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteStorageVolume", storageVolumeTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteStorageVolume", storageVolumeTypeName).Error(),
		)
		return
	}

	// Wait until the volume is no longer present so a re-create does not race a
	// lingering deleting/deleted entry. Any still-present status (deleting,
	// deleted, available) is transient and retried until NotFound.
	if err := storageVolumeWaitForState(r.Meta(), state.Id.ValueString(), "", nil, storageVolumeDeleteTimeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: waiting for removal", storageVolumeTypeName),
			errs.WrapErrorf(err, "Deleting %s: waiting for removal", storageVolumeTypeName).Error(),
		)
		return
	}
}

func (r *storageVolumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// describeStorageVolume delegates to the ENS Service Describe method, which
// calls DescribeStorageVolume with the VolumeId filter and returns the single
// matching storage volume object, or a NotFound error when the array is empty.
func describeStorageVolume(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	return NewService(client).DescribeStorageVolume(id)
}

// storageVolumeWaitForState polls DescribeStorageVolume until the status reaches
// target (empty target means "removed", i.e. NotFound) or one of the failStates
// is observed. It returns nil on success.
func storageVolumeWaitForState(client *connectivity.AliyunClient, id, target string, failStates []string, timeout time.Duration) error {
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	return retry.Retry(timeout, func() *retry.RetryError {
		object, err := describeStorageVolume(client, id)
		if err != nil {
			if errs.NotFoundError(err) {
				if target == "" {
					return nil
				}
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		status := ""
		if v, ok := object["Status"].(string); ok {
			status = v
		}
		for _, failState := range failStates {
			if status == failState {
				return retry.NonRetryableError(errs.WrapErrorf(errs.Error(errs.FailedToReachTargetStatus, status), "storage volume %s entered fail state %q", id, status))
			}
		}
		if status == target {
			return nil
		}
		wait()
		return retry.RetryableError(fmt.Errorf("storage volume %s is still %q", id, status))
	})
}

// storageVolumeState maps a DescribeStorageVolume item onto the resource model.
// AuthUser and AuthPassword are write-only: DescribeStorageVolume never returns
// them, so the previously configured value (from the plan/state carried as prior)
// is preserved instead of collapsing to null.
func storageVolumeState(object map[string]interface{}, prior *storageVolumeResourceModel) storageVolumeResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	// DescribeStorageVolume returns IsEnable/IsAuth as integer; the tea-rpc client
	// may decode numbers as json.Number, float64 or int, so every encoding the
	// transport produces must map back to Bool.
	boolVal := func(key string) types.Bool {
		switch v := object[key].(type) {
		case bool:
			return types.BoolValue(v)
		case float64:
			return types.BoolValue(int(v) != 0)
		case json.Number:
			if n, err := v.Int64(); err == nil {
				return types.BoolValue(n != 0)
			}
		case int:
			return types.BoolValue(v != 0)
		case int64:
			return types.BoolValue(v != 0)
		case string:
			if n, err := strconv.Atoi(v); err == nil {
				return types.BoolValue(n != 0)
			}
		}
		return types.BoolNull()
	}
	m := storageVolumeResourceModel{
		Id:           str("StorageVolumeId"),
		GatewayId:    str("StorageGatewayId"),
		StorageId:    str("StorageId"),
		EnsRegionId:  str("EnsRegionId"),
		VolumeName:   str("StorageVolumeName"),
		Description:  str("Description"),
		IsEnable:     boolVal("IsEnable"),
		IsAuth:       boolVal("IsAuth"),
		AuthProtocol: str("AuthProtocol"),
		Status:       str("Status"),
		CreationTime: str("CreationTime"),
		TargetName:   str("TargetName"),
	}
	// Write-only fields are not returned by DescribeStorageVolume; carry them
	// over from the plan (create) or prior state (refresh) verbatim.
	if prior != nil {
		m.AuthUser = prior.AuthUser
		m.AuthPassword = prior.AuthPassword
	}
	return m
}

// firstStorageVolumeId extracts the first created volume id from the
// CreateStorageVolume response, whose VolumeId is an array of strings.
func firstStorageVolumeId(response map[string]interface{}) (string, error) {
	raw, err := jsonpath.Get("$.VolumeId", response)
	if err != nil {
		return "", errs.WrapErrorf(err, "reading VolumeId from CreateStorageVolume response: %v", response)
	}
	switch v := raw.(type) {
	case []interface{}:
		for _, item := range v {
			if s, ok := item.(string); ok && s != "" {
				return s, nil
			}
		}
	case string:
		if v != "" {
			return v, nil
		}
	}
	return "", fmt.Errorf("CreateStorageVolume response did not contain a usable VolumeId: %v", response)
}

// boolToStr encodes a bool as the API string enum "1"/"0".
func boolToStr(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
