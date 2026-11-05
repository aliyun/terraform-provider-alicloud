package ens

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
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	ensNetworkInterfaceTypeName = "alicloud_ens_network_interface"
	// terraform-plugin-framework does not support user-configurable resource
	// timeouts; mirror the SDKv2 default of 5 minutes per operation.
	ensNetworkInterfaceTimeout = 5 * time.Minute
)

var (
	_ resource.Resource                = &ensNetworkInterfaceResource{}
	_ resource.ResourceWithImportState = &ensNetworkInterfaceResource{}
)

func NewEnsNetworkInterfaceResource() resource.Resource {
	return &ensNetworkInterfaceResource{}
}

type ensNetworkInterfaceResource struct {
	fwadapt.ResourceBase
}

type ensNetworkInterfaceResourceModel struct {
	Id                    types.String   `tfsdk:"id"`
	NetworkInterfaceId    types.String   `tfsdk:"network_interface_id"`
	VSwitchId             types.String   `tfsdk:"vswitch_id"`
	SecurityGroupIds      []types.String `tfsdk:"security_group_ids"`
	NetworkInterfaceName  types.String   `tfsdk:"network_interface_name"`
	Description           types.String   `tfsdk:"description"`
	VmncLearn             types.Bool     `tfsdk:"vmnc_learn"`
	InstanceId            types.String   `tfsdk:"instance_id"`
	PrivateIpAddressCount types.Int64    `tfsdk:"private_ip_address_count"`
	Status                types.String   `tfsdk:"status"`
	PrimaryIp             types.String   `tfsdk:"primary_ip"`
	PrimaryIpType         types.String   `tfsdk:"primary_ip_type"`
	MacAddress            types.String   `tfsdk:"mac_address"`
	NetworkId             types.String   `tfsdk:"network_id"`
	EnsRegionId           types.String   `tfsdk:"ens_region_id"`
	CreateTime            types.String   `tfsdk:"create_time"`
	PrivateIpAddresses    []types.String `tfsdk:"private_ip_addresses"`
	Ipv6Addresses         []types.String `tfsdk:"ipv6_addresses"`
	Type                  types.String   `tfsdk:"type"`
}

func (r *ensNetworkInterfaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ensNetworkInterfaceTypeName
}

func (r *ensNetworkInterfaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network_interface_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// VSwitchId is read-only on the CloudSpec model and cannot be
			// modified after creation (no update API mutates it).
			"vswitch_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// SecurityGroupIds is required at creation; ModifyNetworkInterfaceAttribute
			// does not touch it, so changes force replacement.
			"security_group_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.SizeAtMost(5),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"network_interface_name": schema.StringAttribute{
				Optional: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"vmnc_learn": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			// InstanceId drives AttachNetworkInterface/DetachNetworkInterface:
			// setting it attaches the ENI to the instance, clearing detaches.
			"instance_id": schema.StringAttribute{
				Optional: true,
			},
			// AssignPrivateIpAddresses assigns one server-chosen secondary
			// private IP per call; the provider reconciles the desired count.
			"private_ip_address_count": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int64{
					int64validator.Between(0, 10),
				},
			},
			"status":          schema.StringAttribute{Computed: true},
			"primary_ip":      schema.StringAttribute{Computed: true},
			"primary_ip_type": schema.StringAttribute{Computed: true},
			"mac_address":     schema.StringAttribute{Computed: true},
			"network_id":      schema.StringAttribute{Computed: true},
			"ens_region_id":   schema.StringAttribute{Computed: true},
			"create_time":     schema.StringAttribute{Computed: true},
			"private_ip_addresses": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"ipv6_addresses": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"type": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *ensNetworkInterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ensNetworkInterfaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "CreateNetworkInterface"
	request := map[string]interface{}{
		"VSwitchId":        plan.VSwitchId.ValueString(),
		"SecurityGroupIds": ensNetworkInterfaceSecurityGroupIds(plan.SecurityGroupIds),
		"ClientToken":      helper.BuildClientToken(action),
	}
	if !plan.NetworkInterfaceName.IsNull() && !plan.NetworkInterfaceName.IsUnknown() {
		request["Name"] = plan.NetworkInterfaceName.ValueString()
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		request["Description"] = plan.Description.ValueString()
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(ensNetworkInterfaceTimeout, func() *retry.RetryError {
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
			fmt.Sprintf("Creating %s: calling CreateNetworkInterface", ensNetworkInterfaceTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateNetworkInterface", ensNetworkInterfaceTypeName).Error(),
		)
		return
	}

	id := ensNetworkInterfaceExtractId(response)
	if id == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: parsing NetworkInterfaceIds from response", ensNetworkInterfaceTypeName),
			fmt.Sprintf("CreateNetworkInterface response did not contain a NetworkInterfaceId: %#v", response),
		)
		return
	}

	// VmncLearn has a server default; apply the configured value when set.
	if !plan.VmncLearn.IsNull() && !plan.VmncLearn.IsUnknown() {
		if err := r.modifyVmncLearn(ctx, id, plan.VmncLearn.ValueBool()); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: calling ModifyNetworkInterfaceVmncLearn", ensNetworkInterfaceTypeName),
				err.Error(),
			)
			return
		}
	}

	if plan.InstanceId.ValueString() != "" {
		if err := r.attachNetworkInterface(id, plan.InstanceId.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: calling AttachNetworkInterface", ensNetworkInterfaceTypeName),
				err.Error(),
			)
			return
		}
	}

	if !plan.PrivateIpAddressCount.IsNull() && !plan.PrivateIpAddressCount.IsUnknown() && plan.PrivateIpAddressCount.ValueInt64() > 0 {
		if err := r.reconcilePrivateIpAddresses(id, plan.VSwitchId.ValueString(), 0, int(plan.PrivateIpAddressCount.ValueInt64())); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: calling AssignPrivateIpAddresses", ensNetworkInterfaceTypeName),
				err.Error(),
			)
			return
		}
	}

	object, err := r.describeNetworkInterface(id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeNetworkInterfaces", ensNetworkInterfaceTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, ensNetworkInterfaceState(object))...)
}

func (r *ensNetworkInterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ensNetworkInterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := r.describeNetworkInterface(state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeNetworkInterfaces", ensNetworkInterfaceTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, ensNetworkInterfaceState(object))...)
}

func (r *ensNetworkInterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ensNetworkInterfaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()

	// ModifyNetworkInterfaceAttribute covers name and description.
	if !plan.NetworkInterfaceName.Equal(state.NetworkInterfaceName) || !plan.Description.Equal(state.Description) {
		if err := r.modifyNetworkInterfaceAttribute(id, plan.NetworkInterfaceName.ValueString(), plan.Description.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling ModifyNetworkInterfaceAttribute", ensNetworkInterfaceTypeName),
				err.Error(),
			)
			return
		}
	}

	// VmncLearn is updated through a dedicated API.
	if !plan.VmncLearn.Equal(state.VmncLearn) {
		if err := r.modifyVmncLearn(ctx, id, plan.VmncLearn.ValueBool()); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling ModifyNetworkInterfaceVmncLearn", ensNetworkInterfaceTypeName),
				err.Error(),
			)
			return
		}
	}

	// Attach/detach on instance_id change.
	if !plan.InstanceId.Equal(state.InstanceId) {
		if state.InstanceId.ValueString() != "" {
			if err := r.detachNetworkInterface(id); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Updating %s: calling DetachNetworkInterface", ensNetworkInterfaceTypeName),
					err.Error(),
				)
				return
			}
		}
		if plan.InstanceId.ValueString() != "" {
			if err := r.attachNetworkInterface(id, plan.InstanceId.ValueString()); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Updating %s: calling AttachNetworkInterface", ensNetworkInterfaceTypeName),
					err.Error(),
				)
				return
			}
		}
	}

	// Reconcile secondary private IP count.
	currentCount := 0
	for _, v := range state.PrivateIpAddresses {
		if !v.IsNull() {
			currentCount++
		}
	}
	desired := 0
	if !plan.PrivateIpAddressCount.IsNull() && !plan.PrivateIpAddressCount.IsUnknown() {
		desired = int(plan.PrivateIpAddressCount.ValueInt64())
	}
	if desired != currentCount {
		var current []string
		for _, v := range state.PrivateIpAddresses {
			if !v.IsNull() {
				current = append(current, v.ValueString())
			}
		}
		if err := r.reconcilePrivateIpAddresses(id, state.VSwitchId.ValueString(), currentCount, desired, current...); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: reconciling private IP addresses", ensNetworkInterfaceTypeName),
				err.Error(),
			)
			return
		}
	}

	object, err := r.describeNetworkInterface(id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeNetworkInterfaces", ensNetworkInterfaceTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, ensNetworkInterfaceState(object))...)
}

func (r *ensNetworkInterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ensNetworkInterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()

	// A bound ENI must be detached before it can be deleted.
	if state.InstanceId.ValueString() != "" {
		if err := r.detachNetworkInterface(id); err != nil {
			if !errs.NotFoundError(err) && !errs.IsExpectedErrors(err, []string{"InvalidNetworkInterface.NotFound", "NotFound"}) {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Deleting %s: calling DetachNetworkInterface", ensNetworkInterfaceTypeName),
					err.Error(),
				)
				return
			}
		}
	}

	action := "DeleteNetworkInterfaces"
	request := map[string]interface{}{
		"NetworkInterfaceIds": []string{id},
	}
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(ensNetworkInterfaceTimeout, func() *retry.RetryError {
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
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"InvalidNetworkInterface.NotFound", "NotFound"}) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteNetworkInterfaces", ensNetworkInterfaceTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteNetworkInterfaces", ensNetworkInterfaceTypeName).Error(),
		)
		return
	}

	// DeleteNetworkInterfaces is asynchronous: poll until the ENI disappears.
	stateConf := helper.BuildStateConf([]string{"Available", "Attaching", "InUse", "Detaching", "Deleting"}, []string{}, ensNetworkInterfaceTimeout, 30*time.Second, r.ensNetworkInterfaceStateRefreshFunc(id))
	if _, err := stateConf.WaitForState(); err != nil {
		if errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: waiting for deletion", ensNetworkInterfaceTypeName),
			errs.WrapErrorf(err, "Deleting %s: waiting for deletion", ensNetworkInterfaceTypeName).Error(),
		)
		return
	}
}

func (r *ensNetworkInterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// --- service helpers ---

func (r *ensNetworkInterfaceResource) describeNetworkInterface(id string) (map[string]interface{}, error) {
	client := r.Meta()
	action := "DescribeNetworkInterfaces"
	query := map[string]interface{}{"NetworkInterfaceId": id}
	request := map[string]interface{}{}
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(ensNetworkInterfaceTimeout, func() *retry.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, ok := response["NetworkInterfaceSets"]
	if !ok || v == nil {
		return nil, errs.WrapErrorf(errs.NotFoundErr("NetworkInterface", id), errs.NotFoundMsg, response)
	}
	sets, ok := v.(map[string]interface{})
	if !ok {
		return nil, errs.WrapErrorf(fmt.Errorf("unexpected NetworkInterfaceSets type: %T", v), errs.FailedGetAttributeMsg, id, "$.NetworkInterfaceSets", response)
	}
	items, ok := sets["NetworkInterfaceSet"].([]interface{})
	if !ok || len(items) == 0 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("NetworkInterface", id), errs.NotFoundMsg, response)
	}
	object, ok := items[0].(map[string]interface{})
	if !ok {
		return nil, errs.WrapErrorf(fmt.Errorf("unexpected NetworkInterfaceSet item type: %T", items[0]), errs.FailedGetAttributeMsg, id, "$.NetworkInterfaceSets.NetworkInterfaceSet", response)
	}
	return object, nil
}

func (r *ensNetworkInterfaceResource) ensNetworkInterfaceStateRefreshFunc(id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := r.describeNetworkInterface(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", err
		}
		status, _ := object["Status"].(string)
		return object, status, nil
	}
}

func (r *ensNetworkInterfaceResource) modifyNetworkInterfaceAttribute(id, name, description string) error {
	action := "ModifyNetworkInterfaceAttribute"
	request := map[string]interface{}{
		"NetworkInterfaceId": id,
	}
	update := false
	if name != "" {
		request["NetworkInterfaceName"] = name
		update = true
	}
	if description != "" {
		request["Description"] = description
		update = true
	}
	if !update {
		return nil
	}
	return r.callEns(action, nil, request)
}

func (r *ensNetworkInterfaceResource) modifyVmncLearn(ctx context.Context, id string, vmncLearn bool) error {
	action := "ModifyNetworkInterfaceVmncLearn"
	request := map[string]interface{}{
		"NetworkInterfaceId": id,
		"VmncLearn":          vmncLearn,
	}
	return r.callEns(action, nil, request)
}

func (r *ensNetworkInterfaceResource) attachNetworkInterface(id, instanceId string) error {
	action := "AttachNetworkInterface"
	request := map[string]interface{}{
		"NetworkInterfaceId": id,
		"InstanceId":         instanceId,
	}
	return r.callEns(action, nil, request)
}

func (r *ensNetworkInterfaceResource) detachNetworkInterface(id string) error {
	action := "DetachNetworkInterface"
	request := map[string]interface{}{
		"NetworkInterfaceId": id,
	}
	return r.callEns(action, nil, request)
}

// reconcilePrivateIpAddresses brings the number of secondary private IPs in
// line with the desired count. AssignPrivateIpAddresses assigns one
// server-chosen IP per call; UnassignPrivateIpAddresses removes a list.
func (r *ensNetworkInterfaceResource) reconcilePrivateIpAddresses(id, vswitchId string, current, desired int, existing ...string) error {
	if desired > current {
		for i := 0; i < desired-current; i++ {
			if err := r.assignPrivateIpAddress(id, vswitchId); err != nil {
				return err
			}
		}
		return nil
	}
	if desired < current {
		extras := existing[desired:]
		if len(extras) == 0 {
			return nil
		}
		return r.unassignPrivateIpAddresses(id, extras)
	}
	return nil
}

func (r *ensNetworkInterfaceResource) assignPrivateIpAddress(id, vswitchId string) error {
	action := "AssignPrivateIpAddresses"
	request := map[string]interface{}{
		"NetworkInterfaceId": id,
		"VSwitchId":          vswitchId,
	}
	return r.callEns(action, nil, request)
}

func (r *ensNetworkInterfaceResource) unassignPrivateIpAddresses(id string, ips []string) error {
	action := "UnassignPrivateIpAddresses"
	request := map[string]interface{}{
		"NetworkInterfaceId": id,
		"PrivateIpAddress":   ips,
	}
	return r.callEns(action, nil, request)
}

func (r *ensNetworkInterfaceResource) callEns(action string, query, request map[string]interface{}) error {
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(ensNetworkInterfaceTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost("Ens", "2017-11-10", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "", action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

// ensNetworkInterfaceExtractId reads the created NetworkInterfaceId from the
// CreateNetworkInterface response, where it is returned as the first element
// of the NetworkInterfaceIds array (or, defensively, as a bare string).
func ensNetworkInterfaceExtractId(response map[string]interface{}) string {
	raw, ok := response["NetworkInterfaceIds"]
	if !ok || raw == nil {
		return ""
	}
	switch v := raw.(type) {
	case string:
		return v
	case []interface{}:
		if len(v) > 0 {
			if s, ok := v[0].(string); ok {
				return s
			}
		}
	}
	return ""
}

func ensNetworkInterfaceSecurityGroupIds(values []types.String) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		if !v.IsNull() {
			out = append(out, v.ValueString())
		}
	}
	return out
}

// ensNetworkInterfaceState maps a DescribeNetworkInterfaces item onto the
// resource model.
func ensNetworkInterfaceState(object map[string]interface{}) ensNetworkInterfaceResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key]; ok {
			switch s := v.(type) {
			case string:
				if s != "" {
					return types.StringValue(s)
				}
			case json.Number:
				return types.StringValue(s.String())
			}
		}
		return types.StringNull()
	}
	boolVal := func(key string) types.Bool {
		if v, ok := object[key]; ok {
			switch b := v.(type) {
			case bool:
				return types.BoolValue(b)
			case string:
				if parsed, err := strconv.ParseBool(b); err == nil {
					return types.BoolValue(parsed)
				}
			case json.Number:
				if parsed, err := strconv.ParseBool(b.String()); err == nil {
					return types.BoolValue(parsed)
				}
			}
		}
		return types.BoolNull()
	}
	id := str("NetworkInterfaceId")
	var sgs []types.String
	if raw, ok := object["SecurityGroupIds"]; ok {
		if groupMap, ok := raw.(map[string]interface{}); ok {
			if list, ok := groupMap["SecurityGroup"].([]interface{}); ok {
				for _, item := range list {
					if s, ok := item.(string); ok && s != "" {
						sgs = append(sgs, types.StringValue(s))
					}
				}
			}
		}
	}
	if sgs == nil {
		sgs = []types.String{}
	}
	var privateIps []types.String
	if raw, ok := object["PrivateIpSets"]; ok {
		if setMap, ok := raw.(map[string]interface{}); ok {
			if list, ok := setMap["PrivateIpSet"].([]interface{}); ok {
				for _, item := range list {
					if m, ok := item.(map[string]interface{}); ok {
						// Secondary private IPs only; the primary is reported
						// separately via PrimaryIp.
						if primary, ok := m["Primary"].(bool); ok && primary {
							continue
						}
						if ip, ok := m["PrivateIpAddress"].(string); ok && ip != "" {
							privateIps = append(privateIps, types.StringValue(ip))
						}
					}
				}
			}
		}
	}
	if privateIps == nil {
		privateIps = []types.String{}
	}
	var ipv6s []types.String
	if raw, ok := object["Ipv6Sets"]; ok {
		if setMap, ok := raw.(map[string]interface{}); ok {
			if list, ok := setMap["Ipv6Set"].([]interface{}); ok {
				for _, item := range list {
					if m, ok := item.(map[string]interface{}); ok {
						if ip, ok := m["Ipv6Address"].(string); ok && ip != "" {
							ipv6s = append(ipv6s, types.StringValue(ip))
						}
					}
				}
			}
		}
	}
	if ipv6s == nil {
		ipv6s = []types.String{}
	}
	return ensNetworkInterfaceResourceModel{
		Id:                    id,
		NetworkInterfaceId:    id,
		VSwitchId:             str("VSwitchId"),
		SecurityGroupIds:      sgs,
		NetworkInterfaceName:  str("NetworkInterfaceName"),
		Description:           str("Description"),
		VmncLearn:             boolVal("VmncLearn"),
		InstanceId:            str("InstanceId"),
		PrivateIpAddressCount: types.Int64Value(int64(len(privateIps))),
		Status:                str("Status"),
		PrimaryIp:             str("PrimaryIp"),
		PrimaryIpType:         str("PrimaryIpType"),
		MacAddress:            str("MacAddress"),
		NetworkId:             str("NetworkId"),
		EnsRegionId:           str("EnsRegionId"),
		CreateTime:            str("CreationTime"),
		PrivateIpAddresses:    privateIps,
		Ipv6Addresses:         ipv6s,
		Type:                  str("Type"),
	}
}
