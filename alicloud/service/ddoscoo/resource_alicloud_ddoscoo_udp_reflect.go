package ddoscoo

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	udpReflectTypeName      = "alicloud_ddoscoo_udp_reflect"
	udpReflectAPIVersion    = "2020-01-01"
	udpReflectProduct       = "ddoscoo"
	udpReflectCreateTimeout = 5 * time.Minute
	udpReflectUpdateTimeout = 5 * time.Minute
	udpReflectDeleteTimeout = 5 * time.Minute
)

var (
	_ resource.Resource                = &udpReflectResource{}
	_ resource.ResourceWithImportState = &udpReflectResource{}
)

// NewDdoscooUdpReflectResource is the native Framework factory registered
// through ServicePackage.Resources.
func NewDdoscooUdpReflectResource() resource.Resource {
	return &udpReflectResource{}
}

// udpReflectResource manages the UDP reflect attack protection policy of a
// DdosCoo instance. The policy is a single per-instance configuration: there is
// no dedicated Create/Read/Update/Delete API set, only ConfigUdpReflect (which
// applies the full UdpSports list) and DescribeUdpReflect (which reads it).
// Each instance owns exactly one policy, so InstanceId is the resource id.
type udpReflectResource struct {
	fwadapt.ResourceBase
}

type udpReflectResourceModel struct {
	Id         types.String `tfsdk:"id"`
	InstanceId types.String `tfsdk:"instance_id"`
	UdpSports  types.List   `tfsdk:"udp_sports"`
}

func (r *udpReflectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = udpReflectTypeName
}

func (r *udpReflectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// InstanceId identifies both the owning instance and the single
			// policy it carries, so a change means a different policy and a
			// replacement.
			"instance_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// UdpSports is the full set of UDP reflect source ports to block.
			// ConfigUdpReflect applies it as a full replacement, so the list
			// is the source of truth and an empty list clears the policy. The
			// API returns the ports in arbitrary order, so the read path
			// normalizes them to an ascending list; the plan modifier below
			// applies the same canonical order to the configured value, which
			// keeps plan-vs-state comparison stable and avoids a perpetual
			// diff when the user supplies the ports in a different order.
			"udp_sports": schema.ListAttribute{
				Required:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listSortPlanModifier{},
				},
			},
		},
	}
}

func (r *udpReflectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan udpReflectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := plan.InstanceId.ValueString()
	if instanceId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s", udpReflectTypeName),
			"instance_id is required",
		)
		return
	}

	if err := r.configUdpReflect(instanceId, udpReflectPlanPorts(plan.UdpSports), udpReflectCreateTimeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling ConfigUdpReflect", udpReflectTypeName),
			err.Error(),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeDdoscooUdpReflect(instanceId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading the policy after configuration", udpReflectTypeName),
			err.Error(),
		)
		return
	}

	model := udpReflectModelFromAPI(ctx, object, plan)
	model.Id = types.StringValue(instanceId)
	model.InstanceId = types.StringValue(instanceId)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *udpReflectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state udpReflectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := state.Id.ValueString()
	if instanceId == "" {
		instanceId = state.InstanceId.ValueString()
	}

	object, err := NewService(r.Meta()).DescribeDdoscooUdpReflect(instanceId)
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeUdpReflect", udpReflectTypeName),
			err.Error(),
		)
		return
	}

	model := udpReflectModelFromAPI(ctx, object, state)
	model.Id = types.StringValue(instanceId)
	model.InstanceId = types.StringValue(instanceId)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *udpReflectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state udpReflectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := state.Id.ValueString()
	if instanceId == "" {
		instanceId = state.InstanceId.ValueString()
	}

	// ConfigUdpReflect applies the full UdpSports list, so any change (added,
	// removed or reordered ports) is pushed as the new full configuration.
	if !plan.UdpSports.Equal(state.UdpSports) {
		if err := r.configUdpReflect(instanceId, udpReflectPlanPorts(plan.UdpSports), udpReflectUpdateTimeout); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling ConfigUdpReflect", udpReflectTypeName),
				err.Error(),
			)
			return
		}
	}

	object, err := NewService(r.Meta()).DescribeDdoscooUdpReflect(instanceId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update", udpReflectTypeName),
			err.Error(),
		)
		return
	}

	model := udpReflectModelFromAPI(ctx, object, plan)
	model.Id = types.StringValue(instanceId)
	model.InstanceId = types.StringValue(instanceId)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *udpReflectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state udpReflectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := state.Id.ValueString()
	if instanceId == "" {
		instanceId = state.InstanceId.ValueString()
	}

	// There is no dedicated removal API: clearing the policy is done by
	// configuring an empty UdpSports list. If the instance itself is gone the
	// policy is gone with it, so the documented not-found codes are treated as
	// a successful delete.
	err := r.configUdpReflect(instanceId, []int64{}, udpReflectDeleteTimeout)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InstanceNotFound", "ddos_coop3301"}) || errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: clearing the policy via ConfigUdpReflect", udpReflectTypeName),
			err.Error(),
		)
		return
	}
}

func (r *udpReflectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// listSortPlanModifier normalizes a planned udp_sports list to ascending
// order. The API returns the ports in arbitrary order and the read path stores
// them sorted, so applying the same canonical order to the configured value
// keeps the plan stable and prevents a perpetual diff when the user lists the
// ports in a different order than the server happens to return.
type listSortPlanModifier struct{}

func (listSortPlanModifier) Description(_ context.Context) string {
	return "Sorts the planned UDP sports list ascending so the configured order does not cause a perpetual diff against the API-normalized state."
}

func (m listSortPlanModifier) MarkdownDescription(_ context.Context) string {
	return m.Description(nil)
}

func (listSortPlanModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	elems := req.ConfigValue.Elements()
	ports := make([]int64, 0, len(elems))
	for _, e := range elems {
		v, ok := e.(types.Int64)
		if !ok || v.IsNull() || v.IsUnknown() {
			// Unexpected element shape; leave the plan unchanged rather than
			// risk masking a real error.
			resp.PlanValue = req.ConfigValue
			return
		}
		ports = append(ports, v.ValueInt64())
	}
	sort.Slice(ports, func(i, j int) bool { return ports[i] < ports[j] })
	list, diags := types.ListValueFrom(ctx, types.Int64Type, ports)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.PlanValue = list
}

// configUdpReflect applies the full UdpSports list to a DdosCoo instance. The
// Config parameter is a JSON body whose only field is UdpSports; an empty slice
// clears the policy. The call is retried for the usual throttling and transient
// errors; the QPS limit of ConfigUdpReflect is 10/s.
func (r *udpReflectResource) configUdpReflect(instanceId string, ports []int64, timeout time.Duration) error {
	configBody, err := json.Marshal(map[string]interface{}{"UdpSports": ports})
	if err != nil {
		return errs.WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId": instanceId,
		"Config":     string(configBody),
		"RegionId":   r.Meta().RegionId,
	}
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(timeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost(udpReflectProduct, udpReflectAPIVersion, "ConfigUdpReflect", nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug("ConfigUdpReflect", response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, instanceId, "ConfigUdpReflect", errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

// udpReflectPlanPorts extracts the configured ports from a planned list. An
// unknown or null list yields an empty slice, which is also the value pushed
// to clear the policy on delete.
func udpReflectPlanPorts(list types.List) []int64 {
	ports := []int64{}
	if list.IsNull() || list.IsUnknown() {
		return ports
	}
	for _, e := range list.Elements() {
		v, ok := e.(types.Int64)
		if !ok || v.IsNull() || v.IsUnknown() {
			continue
		}
		ports = append(ports, v.ValueInt64())
	}
	return ports
}

// udpReflectModelFromAPI overlays the DescribeUdpReflect response onto base,
// which is the plan on create/update and the prior state on read. The API
// returns UdpSports as a list of strings, normalized here to a sorted int list
// so the state is deterministic regardless of the server-side order.
func udpReflectModelFromAPI(ctx context.Context, object map[string]interface{}, base udpReflectResourceModel) udpReflectResourceModel {
	model := base
	ports := UdpReflectResponsePorts(object["UdpSports"])
	list, diags := types.ListValueFrom(ctx, types.Int64Type, ports)
	if diags.HasError() {
		return model
	}
	model.UdpSports = list
	return model
}

// UdpReflectResponsePorts normalizes the UdpSports field of a
// DescribeUdpReflect response into a sorted int64 slice. Exported so the
// acceptance test can assert the server-side policy independently of state.
func UdpReflectResponsePorts(value interface{}) []int64 {
	ports := []int64{}
	arr, ok := value.([]interface{})
	if !ok {
		return ports
	}
	for _, v := range arr {
		n, err := strconv.ParseInt(fmt.Sprint(v), 10, 64)
		if err != nil {
			continue
		}
		ports = append(ports, n)
	}
	sort.Slice(ports, func(i, j int) bool { return ports[i] < ports[j] })
	return ports
}
