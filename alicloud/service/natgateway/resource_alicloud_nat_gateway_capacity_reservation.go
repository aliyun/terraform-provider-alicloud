package natgateway

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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	natGatewayCapacityReservationTypeName          = "alicloud_nat_gateway_capacity_reservation"
	natGatewayCapacityReservationTimeout           = 20 * time.Minute
	natGatewayCapacityReservationApiName           = "ModifyNatGatewayCapacityReservation"
	natGatewayCapacityReservationStatusProvisioned = "provisioned"
	natGatewayCapacityReservationStatusFailed      = "failed"
	natGatewayCapacityReservationStatusDeleted     = "deleted"
)

var (
	_ resource.Resource                = &natGatewayCapacityReservationResource{}
	_ resource.ResourceWithImportState = &natGatewayCapacityReservationResource{}
)

func NewNatGatewayCapacityReservationResource() resource.Resource {
	return &natGatewayCapacityReservationResource{}
}

// natGatewayCapacityReservationResource manages the CU reservation of one NAT
// gateway. The reservation is a singleton per gateway, so the gateway ID is the
// resource identity and changing it replaces the resource.
type natGatewayCapacityReservationResource struct {
	fwadapt.ResourceBase
}

type natGatewayCapacityReservationModel struct {
	Id                        types.String `tfsdk:"id"`
	NatGatewayId              types.String `tfsdk:"nat_gateway_id"`
	CapacityUnits             types.Int64  `tfsdk:"capacity_units"`
	Status                    types.String `tfsdk:"status"`
	StatusReason              types.String `tfsdk:"status_reason"`
	DecreaseRequestsRemaining types.String `tfsdk:"decrease_requests_remaining"`
	LastModifiedTime          types.String `tfsdk:"last_modified_time"`
	Bandwidth                 types.Int64  `tfsdk:"bandwidth"`
	Cps                       types.Int64  `tfsdk:"cps"`
	Pps                       types.Int64  `tfsdk:"pps"`
	SessionCount              types.Int64  `tfsdk:"session_count"`
}

func (r *natGatewayCapacityReservationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = natGatewayCapacityReservationTypeName
}

func (r *natGatewayCapacityReservationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"nat_gateway_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// CapacityUnits is the only write-side parameter of the reservation;
			// bandwidth, sessions, CPS and PPS are values the API derives from it.
			// The product minimum (9000 CU) is enforced by the API, because the
			// quota and the minimum are account-scoped.
			"capacity_units": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"status":                      schema.StringAttribute{Computed: true},
			"status_reason":               schema.StringAttribute{Computed: true},
			"decrease_requests_remaining": schema.StringAttribute{Computed: true},
			"last_modified_time":          schema.StringAttribute{Computed: true},
			"bandwidth":                   schema.Int64Attribute{Computed: true},
			"cps":                         schema.Int64Attribute{Computed: true},
			"pps":                         schema.Int64Attribute{Computed: true},
			"session_count":               schema.Int64Attribute{Computed: true},
		},
	}
}

func (r *natGatewayCapacityReservationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan natGatewayCapacityReservationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	natGatewayId := plan.NatGatewayId.ValueString()
	if err := r.modifyCapacityReservation(ctx, natGatewayId, plan.CapacityUnits.ValueInt64(), false); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", natGatewayCapacityReservationTypeName, natGatewayCapacityReservationApiName),
			errs.WrapErrorf(err, "Creating %s: calling %s", natGatewayCapacityReservationTypeName, natGatewayCapacityReservationApiName).Error(),
		)
		return
	}

	object, err := r.waitForCapacityReservationStatus(ctx, natGatewayId, plan.CapacityUnits.ValueInt64(), false)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for reservation to become ready", natGatewayCapacityReservationTypeName),
			errs.WrapErrorf(err, "Creating %s: waiting for reservation to become ready", natGatewayCapacityReservationTypeName).Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, natGatewayCapacityReservationState(object))...)
}

func (r *natGatewayCapacityReservationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state natGatewayCapacityReservationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	natGatewayId := state.Id.ValueString()
	if natGatewayId == "" {
		natGatewayId = state.NatGatewayId.ValueString()
	}

	object, err := DescribeNatGatewayCapacityReservation(ctx, r.Meta(), natGatewayId)
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeNatGatewayCapacityReservation", natGatewayCapacityReservationTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeNatGatewayCapacityReservation", natGatewayCapacityReservationTypeName).Error(),
		)
		return
	}

	// A cancelled reservation keeps answering with Status=deleted, and a gateway
	// without any reservation answers with no status at all; both mean the
	// Terraform resource no longer exists even though the gateway still does.
	switch capacityReservationStatus(object) {
	case "", natGatewayCapacityReservationStatusDeleted:
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, natGatewayCapacityReservationState(object))...)
}

func (r *natGatewayCapacityReservationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state natGatewayCapacityReservationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	natGatewayId := state.Id.ValueString()
	if !plan.CapacityUnits.Equal(state.CapacityUnits) {
		if err := r.modifyCapacityReservation(ctx, natGatewayId, plan.CapacityUnits.ValueInt64(), false); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling %s", natGatewayCapacityReservationTypeName, natGatewayCapacityReservationApiName),
				errs.WrapErrorf(err, "Updating %s: calling %s", natGatewayCapacityReservationTypeName, natGatewayCapacityReservationApiName).Error(),
			)
			return
		}

		// The write is asynchronous: the API can still report the previous
		// capacity and status until the new reservation takes effect, so wait
		// for the planned value instead of the first ready status.
		if _, err := r.waitForCapacityReservationStatus(ctx, natGatewayId, plan.CapacityUnits.ValueInt64(), false); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: waiting for reservation to become ready", natGatewayCapacityReservationTypeName),
				errs.WrapErrorf(err, "Updating %s: waiting for reservation to become ready", natGatewayCapacityReservationTypeName).Error(),
			)
			return
		}
	}

	object, err := DescribeNatGatewayCapacityReservation(ctx, r.Meta(), natGatewayId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeNatGatewayCapacityReservation", natGatewayCapacityReservationTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeNatGatewayCapacityReservation", natGatewayCapacityReservationTypeName).Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, natGatewayCapacityReservationState(object))...)
}

func (r *natGatewayCapacityReservationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state natGatewayCapacityReservationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	natGatewayId := state.Id.ValueString()
	if natGatewayId == "" {
		natGatewayId = state.NatGatewayId.ValueString()
	}

	// Cancelling the reservation counts against the daily decrease quota and
	// never touches the parent NAT gateway. An already cancelled reservation is
	// gone, so do not burn another decrease request on it.
	current, err := DescribeNatGatewayCapacityReservation(ctx, r.Meta(), natGatewayId)
	if err != nil {
		if errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DescribeNatGatewayCapacityReservation", natGatewayCapacityReservationTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DescribeNatGatewayCapacityReservation", natGatewayCapacityReservationTypeName).Error(),
		)
		return
	}
	if status := capacityReservationStatus(current); status == "" || status == natGatewayCapacityReservationStatusDeleted {
		return
	}

	if err := r.modifyCapacityReservation(ctx, natGatewayId, 0, true); err != nil {
		if errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling %s", natGatewayCapacityReservationTypeName, natGatewayCapacityReservationApiName),
			errs.WrapErrorf(err, "Deleting %s: calling %s", natGatewayCapacityReservationTypeName, natGatewayCapacityReservationApiName).Error(),
		)
		return
	}

	if _, err := r.waitForCapacityReservationStatus(ctx, natGatewayId, 0, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: waiting for the reservation to be cancelled", natGatewayCapacityReservationTypeName),
			errs.WrapErrorf(err, "Deleting %s: waiting for the reservation to be cancelled", natGatewayCapacityReservationTypeName).Error(),
		)
		return
	}
}

func (r *natGatewayCapacityReservationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// modifyCapacityReservation creates, updates or cancels the reservation through
// the single write API. NatGatewayCapacity is a JSON-styled parameter: the
// object is serialized as a whole instead of being flattened into query keys.
func (r *natGatewayCapacityReservationResource) modifyCapacityReservation(ctx context.Context, natGatewayId string, capacityUnits int64, reset bool) error {
	request := map[string]interface{}{
		"RegionId":     r.Meta().RegionId,
		"NatGatewayId": natGatewayId,
		"ClientToken":  helper.BuildClientToken(natGatewayCapacityReservationApiName),
	}
	if reset {
		request["ResetCapacityReservation"] = "true"
	} else {
		capacity, err := natGatewayCapacityReservationCapacityJSON(capacityUnits)
		if err != nil {
			return err
		}
		request["NatGatewayCapacity"] = capacity
	}

	response, err := helper.RetryRpcPost(ctx, natGatewayCapacityReservationApiName, func() (map[string]interface{}, error) {
		return r.Meta().RpcPost("Vpc", "2016-04-28", natGatewayCapacityReservationApiName, nil, request, true)
	}, natGatewayCapacityReservationTimeout, helper.NeedRetry)
	logs.AddDebug(natGatewayCapacityReservationApiName, response, request)
	return err
}

func natGatewayCapacityReservationCapacityJSON(capacityUnits int64) (string, error) {
	capacity, err := json.Marshal(map[string]interface{}{"CapacityUnits": capacityUnits})
	if err != nil {
		return "", err
	}
	return string(capacity), nil
}

// waitForCapacityReservationStatus polls the reservation until it reaches the
// provisioned status with the expected capacity. Waiting for the capacity value
// matters because the write is asynchronous and Describe can still answer with
// the previous snapshot. A failed reservation is reported with the API reason.
// With waitForGone set, an empty/deleted answer completes the wait instead.
func (r *natGatewayCapacityReservationResource) waitForCapacityReservationStatus(ctx context.Context, natGatewayId string, expectedCapacityUnits int64, waitForGone bool) (map[string]interface{}, error) {
	var object map[string]interface{}
	wait := helper.IncrementalWait(5*time.Second, 10*time.Second)
	err := retry.RetryContext(ctx, natGatewayCapacityReservationTimeout, func() *retry.RetryError {
		current, err := DescribeNatGatewayCapacityReservation(ctx, r.Meta(), natGatewayId)
		if err != nil {
			if errs.NotFoundError(err) && waitForGone {
				return nil
			}
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}

		status := capacityReservationStatus(current)
		if waitForGone && (status == "" || status == natGatewayCapacityReservationStatusDeleted) {
			object = current
			return nil
		}
		if status == natGatewayCapacityReservationStatusProvisioned {
			if units := capacityReservationCapacityUnits(current); units != nil && *units == expectedCapacityUnits {
				object = current
				return nil
			}
			wait()
			return retry.RetryableError(fmt.Errorf("capacity reservation of NAT gateway %s is provisioned with %s CU, waiting for %d CU", natGatewayId, capacityReservationCapacityUnitsString(current), expectedCapacityUnits))
		}
		// A failed status is only terminal when it describes the requested
		// capacity; otherwise Describe is still returning the snapshot taken
		// before this write. A cancellation is not aborted by it either: the
		// reservation may have failed before the cancel request.
		if status == natGatewayCapacityReservationStatusFailed && !waitForGone {
			if units := capacityReservationCapacityUnits(current); units != nil && *units == expectedCapacityUnits {
				return retry.NonRetryableError(fmt.Errorf("capacity reservation of NAT gateway %s entered status failed: %s", natGatewayId, capacityReservationStatusReason(current)))
			}
		}
		wait()
		return retry.RetryableError(fmt.Errorf("capacity reservation of NAT gateway %s is in status %q, waiting for %s", natGatewayId, status, natGatewayCapacityReservationStatusProvisioned))
	})
	if err != nil {
		return nil, err
	}
	return object, nil
}

func capacityReservationStatus(object map[string]interface{}) string {
	if v, ok := object["CapacityReservationStatus"].(map[string]interface{}); ok {
		if status, ok := v["Status"].(string); ok {
			return status
		}
	}
	return ""
}

func capacityReservationStatusReason(object map[string]interface{}) string {
	if v, ok := object["CapacityReservationStatus"].(map[string]interface{}); ok {
		if reason, ok := v["Reason"].(string); ok {
			return reason
		}
	}
	return ""
}

func capacityReservationCapacityUnits(object map[string]interface{}) *int64 {
	if capacity, ok := object["NatGatewayCapacity"].(map[string]interface{}); ok {
		if units, ok := capacityReservationInt64Pointer(capacity["CapacityUnits"]); ok {
			return &units
		}
	}
	return nil
}

func capacityReservationCapacityUnitsString(object map[string]interface{}) string {
	if units := capacityReservationCapacityUnits(object); units != nil {
		return strconv.FormatInt(*units, 10)
	}
	return "<none>"
}

// natGatewayCapacityReservationState maps a DescribeNatGatewayCapacityReservation
// response onto the resource model. Numbers arrive as float64, json.Number or
// int depending on the transport, so every encoding is normalized explicitly.
func natGatewayCapacityReservationState(object map[string]interface{}) natGatewayCapacityReservationModel {
	model := natGatewayCapacityReservationModel{
		NatGatewayId:              capacityReservationString(object["NatGatewayId"]),
		Status:                    capacityReservationString(capacityReservationStatus(object)),
		StatusReason:              capacityReservationString(capacityReservationStatusReason(object)),
		DecreaseRequestsRemaining: capacityReservationString(object["DecreaseRequestsRemaining"]),
		LastModifiedTime:          capacityReservationString(object["LastModifiedTime"]),
	}
	model.Id = model.NatGatewayId

	if units := capacityReservationCapacityUnits(object); units != nil {
		model.CapacityUnits = types.Int64Value(*units)
	}
	if values, ok := object["NatGatewayPerformanceValues"].(map[string]interface{}); ok {
		model.Bandwidth = capacityReservationInt64(values["Bandwidth"])
		model.Cps = capacityReservationInt64(values["Cps"])
		model.Pps = capacityReservationInt64(values["Pps"])
		model.SessionCount = capacityReservationInt64(values["SessionCount"])
	}
	return model
}

func capacityReservationString(value interface{}) types.String {
	switch v := value.(type) {
	case string:
		if v != "" {
			return types.StringValue(v)
		}
	case float64:
		return types.StringValue(strconv.FormatFloat(v, 'f', -1, 64))
	case json.Number:
		return types.StringValue(v.String())
	case int:
		return types.StringValue(strconv.Itoa(v))
	case int64:
		return types.StringValue(strconv.FormatInt(v, 10))
	}
	return types.StringNull()
}

func capacityReservationInt64(value interface{}) types.Int64 {
	if v, ok := capacityReservationInt64Pointer(value); ok {
		return types.Int64Value(v)
	}
	return types.Int64Null()
}

func capacityReservationInt64Pointer(value interface{}) (int64, bool) {
	switch v := value.(type) {
	case float64:
		return int64(v), true
	case json.Number:
		if n, err := v.Int64(); err == nil {
			return n, true
		}
	case int:
		return int64(v), true
	case int64:
		return v, true
	case string:
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n, true
		}
	}
	return 0, false
}
