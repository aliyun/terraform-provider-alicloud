package ens

import (
	"context"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const ensVpdTypeName = "alicloud_ens_vpd"

// ensVpdTimeout replaces the SDKv2 resource-level default timeouts
// (Create/Update/Delete each 5 minutes); terraform-plugin-framework does not
// support user-configurable resource timeouts yet.
const ensVpdTimeout = 5 * time.Minute

var (
	_ resource.Resource                = &ensVpdResource{}
	_ resource.ResourceWithImportState = &ensVpdResource{}
)

func NewEnsVpdResource() resource.Resource {
	return &ensVpdResource{}
}

type ensVpdResource struct {
	fwadapt.ResourceBase
}

type ensVpdResourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnsRegionId   types.String `tfsdk:"ens_region_id"`
	NetworkId     types.String `tfsdk:"network_id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	CidrBlock     types.String `tfsdk:"cidr_block"`
	Status        types.String `tfsdk:"status"`
	CreationTime  types.String `tfsdk:"creation_time"`
	VpdVSwitchIds types.List   `tfsdk:"vpd_vswitch_ids"`
	InstanceIds   types.List   `tfsdk:"instance_ids"`
}

func (r *ensVpdResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ensVpdTypeName
}

func (r *ensVpdResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ens_region_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"network_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"cidr_block":    schema.StringAttribute{Computed: true},
			"status":        schema.StringAttribute{Computed: true},
			"creation_time": schema.StringAttribute{Computed: true},
			"vpd_vswitch_ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"instance_ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *ensVpdResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ensVpdResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "CreateVpd"
	request := map[string]interface{}{
		"EnsRegionId": plan.EnsRegionId.ValueString(),
		"NetworkId":   plan.NetworkId.ValueString(),
	}
	if v := plan.Name.ValueString(); v != "" {
		request["Name"] = v
	}
	if v := plan.Description.ValueString(); v != "" {
		request["Description"] = v
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(ensVpdTimeout, func() *retry.RetryError {
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
			fmt.Sprintf("Creating %s: calling CreateVpd", ensVpdTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, ensVpdTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	id := fmt.Sprint(response["VpdId"])
	if id == "" || id == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: empty VpdId in CreateVpd response", ensVpdTypeName),
			fmt.Sprintf("CreateVpd response did not contain a VpdId: %v", response),
		)
		return
	}

	// Creation is asynchronous (Status moves creating -> Available); wait for
	// the target state before refreshing.
	ensServiceV2 := NewService(r.Meta())
	stateConf := helper.BuildStateConf([]string{"Creating"}, []string{"Available"}, ensVpdTimeout, 30*time.Second, ensServiceV2.EnsVpdStateRefreshFunc(id, "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for Available", ensVpdTypeName),
			errs.WrapErrorf(err, errs.IdMsg, id).Error(),
		)
		return
	}

	// Populate Computed attributes by reading back the freshly created VPD,
	// mirroring the SDKv2 "Create ends with Read" convention (R4.16).
	object, err := NewService(r.Meta()).DescribeEnsVpd(id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeVpdAttribute", ensVpdTypeName),
			errs.WrapErrorf(err, errs.IdMsg, id).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, ensVpdState(object))...)
}

func (r *ensVpdResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ensVpdResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := NewService(r.Meta()).DescribeEnsVpd(state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeVpdAttribute", ensVpdTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), "DescribeVpdAttribute", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, ensVpdState(object))...)
}

func (r *ensVpdResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ensVpdResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// UpdateVpd only accepts Name and Description; both have a minLength of 1,
	// so empty values are invalid and must not be sent.
	request := map[string]interface{}{
		"VpdId": state.Id.ValueString(),
	}
	update := false
	if !plan.Name.Equal(state.Name) {
		if v := plan.Name.ValueString(); v != "" {
			request["Name"] = v
			update = true
		}
	}
	if !plan.Description.Equal(state.Description) {
		if v := plan.Description.ValueString(); v != "" {
			request["Description"] = v
			update = true
		}
	}

	if update {
		action := "UpdateVpd"
		var response map[string]interface{}
		var err error
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(ensVpdTimeout, func() *retry.RetryError {
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
				fmt.Sprintf("Updating %s: calling UpdateVpd", ensVpdTypeName),
				errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
			)
			return
		}
	}

	object, err := NewService(r.Meta()).DescribeEnsVpd(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeVpdAttribute", ensVpdTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), "DescribeVpdAttribute", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, ensVpdState(object))...)
}

func (r *ensVpdResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ensVpdResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "DeleteVpd"
	request := map[string]interface{}{
		"VpdId": state.Id.ValueString(),
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(ensVpdTimeout, func() *retry.RetryError {
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
		// Idempotent delete: a missing VPD is already gone.
		if errs.IsExpectedErrors(err, []string{"Vpd.NotFound"}) || errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteVpd", ensVpdTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	// Deletion is asynchronous (Status moves to Deleting then the resource
	// disappears); wait until DescribeVpdAttribute reports NotFound.
	ensServiceV2 := NewService(r.Meta())
	stateConf := helper.BuildStateConf([]string{}, []string{}, ensVpdTimeout, 30*time.Second, ensServiceV2.EnsVpdStateRefreshFunc(state.Id.ValueString(), "VpdId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: waiting for disappearance", ensVpdTypeName),
			errs.WrapErrorf(err, errs.IdMsg, state.Id.ValueString()).Error(),
		)
		return
	}
}

func (r *ensVpdResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// ensVpdState maps a DescribeVpdAttribute response onto the resource model.
// DescribeVpdAttribute returns all schema fields (including user-settable
// ones), so no prior-state preservation is needed.
func ensVpdState(object map[string]interface{}) ensVpdResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	return ensVpdResourceModel{
		Id:            str("VpdId"),
		EnsRegionId:   str("EnsRegionId"),
		NetworkId:     str("NetworkId"),
		Name:          str("Name"),
		Description:   str("Description"),
		CidrBlock:     str("CidrBlock"),
		Status:        str("Status"),
		CreationTime:  str("CreationTime"),
		VpdVSwitchIds: flattenEnsVpdStringList(object["VpdVSwitchIds"]),
		InstanceIds:   flattenEnsVpdStringList(object["InstanceIds"]),
	}
}

// flattenEnsVpdStringList converts a response []interface{} into a
// types.List of strings, tolerating nil or unexpected element types.
func flattenEnsVpdStringList(raw interface{}) types.List {
	if raw == nil {
		return types.ListNull(types.StringType)
	}
	items, ok := raw.([]interface{})
	if !ok || len(items) == 0 {
		return types.ListNull(types.StringType)
	}
	values := make([]attr.Value, 0, len(items))
	for _, v := range items {
		values = append(values, types.StringValue(fmt.Sprint(v)))
	}
	list, diags := types.ListValue(types.StringType, values)
	if diags.HasError() {
		return types.ListNull(types.StringType)
	}
	return list
}
