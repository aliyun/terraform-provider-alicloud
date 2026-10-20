package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const snapshotLockTypeName = "alicloud_ecs_snapshot_lock"

// snapshotLockTimeout replaces the SDKv2 resource-level default timeouts
// (Create/Update/Delete each 5 minutes); terraform-plugin-framework does not
// support user-configurable resource timeouts yet.
const snapshotLockTimeout = 5 * time.Minute

var (
	_ resource.Resource                = &snapshotLockResource{}
	_ resource.ResourceWithImportState = &snapshotLockResource{}
)

func NewSnapshotLockResource() resource.Resource {
	return &snapshotLockResource{}
}

type snapshotLockResource struct {
	fwadapt.ResourceBase
}

type snapshotLockResourceModel struct {
	Id                       types.String `tfsdk:"id"`
	SnapshotLockId           types.String `tfsdk:"snapshot_lock_id"`
	CoolOffPeriod            types.Int64  `tfsdk:"cool_off_period"`
	CoolOffPeriodExpiredTime types.String `tfsdk:"cool_off_period_expired_time"`
	LockDuration             types.Int64  `tfsdk:"lock_duration"`
	LockDurationStartTime    types.String `tfsdk:"lock_duration_start_time"`
	LockExpiredTime          types.String `tfsdk:"lock_expired_time"`
	LockCreationTime         types.String `tfsdk:"lock_creation_time"`
	LockMode                 types.String `tfsdk:"lock_mode"`
	LockStatus               types.String `tfsdk:"lock_status"`
	DryRun                   types.Bool   `tfsdk:"dry_run"`
}

func (r *snapshotLockResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = snapshotLockTypeName
}

func (r *snapshotLockResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"snapshot_lock_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cool_off_period": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(0, 72),
				},
			},
			"lock_duration": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 36500),
				},
			},
			"lock_mode": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("compliance"),
				},
			},
			"dry_run": schema.BoolAttribute{
				Optional: true,
			},
			// The API-managed computed fields carry no UseStateForUnknown:
			// the expired/started times follow lock_duration and cool_off_period
			// updates and lock_status transitions server-side, so the planned
			// state value could never survive an apply.
			"cool_off_period_expired_time": schema.StringAttribute{Computed: true},
			"lock_creation_time":           schema.StringAttribute{Computed: true},
			"lock_duration_start_time":     schema.StringAttribute{Computed: true},
			"lock_expired_time":            schema.StringAttribute{Computed: true},
			"lock_status":                  schema.StringAttribute{Computed: true},
		},
	}
}

func (r *snapshotLockResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan snapshotLockResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "LockSnapshot"
	request := map[string]interface{}{
		"SnapshotId":    plan.SnapshotLockId.ValueString(),
		"RegionId":      r.Client().RegionId,
		"ClientToken":   helper.BuildClientToken(action),
		"CoolOffPeriod": plan.CoolOffPeriod.ValueInt64(),
		"DryRun":        plan.DryRun.ValueBool(),
		"LockDuration":  plan.LockDuration.ValueInt64(),
		"LockMode":      plan.LockMode.ValueString(),
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(snapshotLockTimeout, func() *retry.RetryError {
		response, err = r.Client().RpcPost("Ecs", "2014-05-26", action, nil, request, true)
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
			fmt.Sprintf("Creating %s: calling LockSnapshot", snapshotLockTypeName),
			errs.WrapErrorf(err, "Creating %s: calling LockSnapshot", snapshotLockTypeName).Error(),
		)
		return
	}

	object, err := NewService(r.Client()).DescribeSnapshotLock(plan.SnapshotLockId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeLockedSnapshots", snapshotLockTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeLockedSnapshots", snapshotLockTypeName).Error(),
		)
		return
	}
	state := snapshotLockState(object)
	state.DryRun = plan.DryRun
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *snapshotLockResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state snapshotLockResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := NewService(r.Client()).DescribeSnapshotLock(state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeLockedSnapshots", snapshotLockTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeLockedSnapshots", snapshotLockTypeName).Error(),
		)
		return
	}
	// DryRun is a write-only flag; DescribeLockedSnapshots never returns it,
	// so refresh must keep the previously configured value.
	model := snapshotLockState(object)
	model.DryRun = state.DryRun
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *snapshotLockResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state snapshotLockResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "LockSnapshot"
	request := map[string]interface{}{
		"SnapshotId":    state.Id.ValueString(),
		"RegionId":      r.Client().RegionId,
		"ClientToken":   helper.BuildClientToken(action),
		"CoolOffPeriod": plan.CoolOffPeriod.ValueInt64(),
		"DryRun":        plan.DryRun.ValueBool(),
		"LockDuration":  plan.LockDuration.ValueInt64(),
		"LockMode":      plan.LockMode.ValueString(),
	}
	update := false
	if !plan.CoolOffPeriod.Equal(state.CoolOffPeriod) {
		update = true
	}
	if !plan.LockDuration.Equal(state.LockDuration) {
		update = true
	}
	if !plan.LockMode.Equal(state.LockMode) {
		update = true
	}

	if update {
		var response map[string]interface{}
		var err error
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(snapshotLockTimeout, func() *retry.RetryError {
			response, err = r.Client().RpcPost("Ecs", "2014-05-26", action, nil, request, true)
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
				fmt.Sprintf("Updating %s: calling LockSnapshot", snapshotLockTypeName),
				errs.WrapErrorf(err, "Updating %s: calling LockSnapshot", snapshotLockTypeName).Error(),
			)
			return
		}
	}

	object, err := NewService(r.Client()).DescribeSnapshotLock(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeLockedSnapshots", snapshotLockTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeLockedSnapshots", snapshotLockTypeName).Error(),
		)
		return
	}
	model := snapshotLockState(object)
	model.DryRun = plan.DryRun
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *snapshotLockResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state snapshotLockResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "UnlockSnapshot"
	request := map[string]interface{}{
		"SnapshotId":  state.Id.ValueString(),
		"RegionId":    r.Client().RegionId,
		"ClientToken": helper.BuildClientToken(action),
		"DryRun":      state.DryRun.ValueBool(),
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(snapshotLockTimeout, func() *retry.RetryError {
		response, err = r.Client().RpcPost("Ecs", "2014-05-26", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"InvalidSnapshotLock.NotFound", "InvalidSnapshotId.NotFound"}) || errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling UnlockSnapshot", snapshotLockTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling UnlockSnapshot", snapshotLockTypeName).Error(),
		)
		return
	}
}

func (r *snapshotLockResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// snapshotLockState maps a DescribeLockedSnapshots item onto the resource
// model. DryRun is deliberately left null: it is write-only and every caller
// overlays the plan/state value before writing state.
func snapshotLockState(object map[string]interface{}) snapshotLockResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	// The tea-rpc client decodes LockSnapshot/DescribeLockedSnapshots numbers
	// as json.Number or int rather than float64, so every numeric encoding
	// the transport produces must map back to Int64 or planned values would
	// collapse to null ("inconsistent result after apply").
	int64Val := func(key string) types.Int64 {
		switch v := object[key].(type) {
		case float64:
			return types.Int64Value(int64(v))
		case json.Number:
			if n, err := v.Int64(); err == nil {
				return types.Int64Value(n)
			}
		case int:
			return types.Int64Value(int64(v))
		case int64:
			return types.Int64Value(v)
		case string:
			if n, err := strconv.ParseInt(v, 10, 64); err == nil {
				return types.Int64Value(n)
			}
		}
		return types.Int64Null()
	}
	return snapshotLockResourceModel{
		Id:                       str("SnapshotId"),
		SnapshotLockId:           str("SnapshotId"),
		CoolOffPeriod:            int64Val("CoolOffPeriod"),
		CoolOffPeriodExpiredTime: str("CoolOffPeriodExpiredTime"),
		LockDuration:             int64Val("LockDuration"),
		LockDurationStartTime:    str("LockDurationStartTime"),
		LockExpiredTime:          str("LockExpiredTime"),
		LockCreationTime:         str("LockCreationTime"),
		LockMode:                 str("LockMode"),
		LockStatus:               str("LockStatus"),
	}
}
