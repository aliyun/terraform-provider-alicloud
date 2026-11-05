package slb

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
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

const listenerAclAttachmentTypeName = "alicloud_slb_listener_acl_attachment"

const (
	listenerAclAttachmentTimeoutSeconds = 300
	listenerIdPattern                   = `^[^:]+:(http|https|tcp|udp):[1-9][0-9]*$`
)

// listenerTransientErrors are the listener level transient codes that the
// Set*ListenerAttribute actions may return while another configuration change
// is still being applied. They are retryable; quota, invalid ACL and IP
// version conflicts are not and must fail fast.
var listenerTransientErrors = []string{
	"SystemBusy",
	"OperationBusy",
	"ServiceIsStopping",
	"BackendServer.configuring",
	"ServiceIsConfiguring",
	"TLSPolicyConfiguring",
}

var (
	_ resource.Resource                = &listenerAclAttachmentResource{}
	_ resource.ResourceWithImportState = &listenerAclAttachmentResource{}
)

func NewListenerAclAttachmentResource() resource.Resource {
	return &listenerAclAttachmentResource{}
}

type listenerAclAttachmentResource struct {
	fwadapt.ResourceBase
}

type listenerAclAttachmentResourceModel struct {
	Id         types.String `tfsdk:"id"`
	ListenerId types.String `tfsdk:"listener_id"`
	AclId      types.String `tfsdk:"acl_id"`
	AclType    types.String `tfsdk:"acl_type"`
	AclStatus  types.String `tfsdk:"acl_status"`
}

func (r *listenerAclAttachmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = listenerAclAttachmentTypeName
}

func (r *listenerAclAttachmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an Alicloud SLB listener access control (ACL) attachment resource. " +
			"It binds an existing access control policy group to an existing SLB listener and only manages the ACL settings of that listener.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"listener_id": schema.StringAttribute{
				Required: true,
				Description: "The ID of the SLB listener, in the format of `<load_balancer_id>:<protocol>:<listener_port>`. " +
					"It is the same as the `id` of `alicloud_slb_listener` and the `id` element of `alicloud_slb_listeners`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(listenerIdPattern),
						"must be in the format of <load_balancer_id>:<protocol>:<listener_port>, protocol is one of http, https, tcp and udp",
					),
				},
			},
			"acl_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the access control policy group bound to the listener.",
			},
			"acl_type": schema.StringAttribute{
				Required:    true,
				Description: "The access control type. Valid values: `white` and `black`.",
				Validators: []validator.String{
					stringvalidator.OneOf("white", "black"),
				},
			},
			"acl_status": schema.StringAttribute{
				Computed:    true,
				Description: "The access control status of the listener.",
			},
		},
	}
}

func (r *listenerAclAttachmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan listenerAclAttachmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	listenerId := plan.ListenerId.ValueString()
	aclId := plan.AclId.ValueString()

	if err := updateListenerAcl(ctx, r.Meta(), listenerId, aclId, plan.AclType.ValueString(), true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: binding access control policy group %s", listenerAclAttachmentTypeName, aclId),
			errs.WrapErrorf(err, "Creating %s: binding access control policy group %s", listenerAclAttachmentTypeName, aclId).Error(),
		)
		return
	}

	// The write already took effect, so the identity is persisted before the
	// first read. A transient read failure must not drop a binding that
	// Terraform now owns; the next refresh reconciles the server values.
	resp.Diagnostics.Append(resp.State.Set(ctx, listenerAclAttachmentPlanState(plan))...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := describeListenerWithContext(ctx, r.Meta(), listenerId)
	if err != nil {
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Created %s but the first read failed", listenerAclAttachmentTypeName),
			fmt.Sprintf("The access control policy group %s was bound to listener %s, but reading the listener failed: %s. "+
				"The binding is kept in state and will be reconciled on the next refresh.", aclId, listenerId, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, listenerAclAttachmentState(listenerId, aclId, plan.AclType, object))...)
}

func (r *listenerAclAttachmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state listenerAclAttachmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	listenerId, idAclId, err := parseListenerAclAttachmentId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing resource ID", listenerAclAttachmentTypeName),
			err.Error(),
		)
		return
	}
	targetAclId := state.AclId.ValueString()
	if targetAclId == "" {
		targetAclId = idAclId
	}

	object, err := describeListenerWithContext(ctx, r.Meta(), listenerId)
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeLoadBalancer*ListenerAttribute", listenerAclAttachmentTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeLoadBalancer*ListenerAttribute", listenerAclAttachmentTypeName).Error(),
		)
		return
	}

	// SLB has no dedicated de-association API: the binding is considered gone
	// once access control is disabled or the configured policy group is no
	// longer part of the listener ACLs.
	if status, _ := object["AclStatus"].(string); !strings.EqualFold(status, "on") {
		resp.State.RemoveResource(ctx)
		return
	}

	if targetAclId == "" {
		// Imported without the ACL part. A single bound policy group is
		// unambiguous; several groups need the four segment import form so the
		// target identity is never guessed.
		aclIds := listenerAclIds(object)
		switch len(aclIds) {
		case 0:
			resp.State.RemoveResource(ctx)
			return
		case 1:
			targetAclId = aclIds[0]
		default:
			resp.Diagnostics.AddError(
				fmt.Sprintf("Importing %s: the listener has multiple access control policy groups", listenerAclAttachmentTypeName),
				fmt.Sprintf("Listener %s is bound to several access control policy groups (%s). Import with the four segment ID "+
					"<load_balancer_id>:<protocol>:<listener_port>:<acl_id> to select the managed group.", listenerId, strings.Join(aclIds, ", ")),
			)
			return
		}
	} else if !listenerAclContains(object, targetAclId) {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, listenerAclAttachmentState(listenerId, targetAclId, state.AclType, object))...)
}

func (r *listenerAclAttachmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state listenerAclAttachmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	listenerId, _, err := parseListenerAclAttachmentId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing resource ID", listenerAclAttachmentTypeName),
			err.Error(),
		)
		return
	}

	aclId := plan.AclId.ValueString()
	if !plan.AclId.Equal(state.AclId) || !plan.AclType.Equal(state.AclType) {
		if err := updateListenerAcl(ctx, r.Meta(), listenerId, aclId, plan.AclType.ValueString(), true); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: binding access control policy group %s", listenerAclAttachmentTypeName, aclId),
				errs.WrapErrorf(err, "Updating %s: binding access control policy group %s", listenerAclAttachmentTypeName, aclId).Error(),
			)
			return
		}
	}

	// Persist the planned identity right after the successful write so a
	// transient read failure cannot lose the managed binding.
	resp.Diagnostics.Append(resp.State.Set(ctx, listenerAclAttachmentPlanState(plan))...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := describeListenerWithContext(ctx, r.Meta(), listenerId)
	if err != nil {
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Updated %s but the read failed", listenerAclAttachmentTypeName),
			fmt.Sprintf("The access control settings of listener %s were updated, but reading the listener failed: %s. "+
				"The planned state is kept and will be reconciled on the next refresh.", listenerId, err.Error()),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, listenerAclAttachmentState(listenerId, aclId, plan.AclType, object))...)
}

func (r *listenerAclAttachmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state listenerAclAttachmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	listenerId, _, err := parseListenerAclAttachmentId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing resource ID", listenerAclAttachmentTypeName),
			err.Error(),
		)
		return
	}

	if _, err := describeListenerWithContext(ctx, r.Meta(), listenerId); err != nil {
		if errs.NotFoundError(err) {
			// The listener is gone, so the ACL binding no longer exists.
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DescribeLoadBalancer*ListenerAttribute", listenerAclAttachmentTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DescribeLoadBalancer*ListenerAttribute", listenerAclAttachmentTypeName).Error(),
		)
		return
	}

	if err := updateListenerAcl(ctx, r.Meta(), listenerId, "", "", false); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: disabling listener access control", listenerAclAttachmentTypeName),
			errs.WrapErrorf(err, "Deleting %s: disabling listener access control", listenerAclAttachmentTypeName).Error(),
		)
		return
	}

	// Wait until the access control status is switched off. The server may keep
	// the previous AclId after disabling, so only AclStatus is asserted.
	err = retry.RetryContext(ctx, time.Duration(listenerAclAttachmentTimeoutSeconds)*time.Second, func() *retry.RetryError {
		if err := ctx.Err(); err != nil {
			return retry.NonRetryableError(err)
		}
		object, err := describeListenerWithContext(ctx, r.Meta(), listenerId)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil
			}
			return retry.NonRetryableError(err)
		}
		if status, _ := object["AclStatus"].(string); !strings.EqualFold(status, "on") {
			return nil
		}
		if err := waitListenerAclWithContext(ctx, 3*time.Second); err != nil {
			return retry.NonRetryableError(err)
		}
		return retry.RetryableError(fmt.Errorf("%s: access control of listener %s is still enabled", listenerAclAttachmentTypeName, listenerId))
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: waiting for listener access control to be disabled", listenerAclAttachmentTypeName),
			errs.WrapErrorf(err, "Deleting %s: waiting for listener access control to be disabled", listenerAclAttachmentTypeName).Error(),
		)
		return
	}
}

func (r *listenerAclAttachmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	listenerId, aclId, err := parseListenerAclAttachmentId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Importing %s: parsing import ID", listenerAclAttachmentTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), listenerId)...)
	if aclId != "" {
		// Keep the requested fourth segment as the managed policy group so the
		// following read never adopts a different group.
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("acl_id"), aclId)...)
	}
}

// listenerAclAttachmentPlanState stores the identity that is already known
// once the write succeeds. It is used before the first read so that a transient
// read failure can never leave an unmanaged binding behind; the next refresh
// replaces the assumed status with the server values.
func listenerAclAttachmentPlanState(plan listenerAclAttachmentResourceModel) listenerAclAttachmentResourceModel {
	return listenerAclAttachmentResourceModel{
		Id:         types.StringValue(plan.ListenerId.ValueString()),
		ListenerId: plan.ListenerId,
		AclId:      plan.AclId,
		AclType:    plan.AclType,
		AclStatus:  types.StringValue("on"),
	}
}

// updateListenerAcl only sends the ACL related parameters of the protocol
// specific SetLoadBalancer*ListenerAttribute API, so other listener settings
// (bandwidth, health check, certificates and so on) are left untouched.
func updateListenerAcl(ctx context.Context, client *connectivity.AliyunClient, listenerId, aclId, aclType string, enable bool) error {
	parts, err := ParseSlbListenerId(listenerId)
	if err != nil {
		return errs.WrapError(err)
	}
	if len(parts) != 3 {
		return errs.WrapError(fmt.Errorf("invalid listener ID %s, expected <load_balancer_id>:<protocol>:<listener_port>", listenerId))
	}
	protocol := parts[1]
	port, err := parseListenerPort(parts[2])
	if err != nil {
		return errs.WrapErrorf(err, "invalid listener port in listener ID %s", listenerId)
	}
	action, err := listenerAttributeAction(protocol)
	if err != nil {
		return errs.WrapErrorf(err, "invalid listener ID %s", listenerId)
	}
	request := map[string]interface{}{
		"RegionId":       client.RegionId,
		"LoadBalancerId": parts[0],
		"ListenerPort":   port,
	}
	if enable {
		request["AclStatus"] = "on"
		request["AclId"] = aclId
		request["AclType"] = aclType
	} else {
		request["AclStatus"] = "off"
	}

	var response map[string]interface{}
	err = retry.RetryContext(ctx, time.Duration(listenerAclAttachmentTimeoutSeconds)*time.Second, func() *retry.RetryError {
		if err := ctx.Err(); err != nil {
			return retry.NonRetryableError(err)
		}
		response, err = client.RpcPost("Slb", string(connectivity.ApiVersion20140515), action, nil, request, true)
		if err != nil {
			if isTransientListenerError(err) {
				if err := waitListenerAclWithContext(ctx, 3*time.Second); err != nil {
					return retry.NonRetryableError(err)
				}
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, listenerId, action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

// describeListenerWithContext reads the protocol specific listener attributes
// with the CRUD context, so a cancelled apply stops retrying and sleeping
// instead of inheriting an internal retry loop that ignores cancellation.
func describeListenerWithContext(ctx context.Context, client *connectivity.AliyunClient, listenerId string) (map[string]interface{}, error) {
	parts, err := ParseSlbListenerId(listenerId)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	if len(parts) != 3 {
		return nil, errs.WrapError(fmt.Errorf("invalid listener ID %s, expected <load_balancer_id>:<protocol>:<listener_port>", listenerId))
	}
	action := fmt.Sprintf("DescribeLoadBalancer%sListenerAttribute", strings.ToUpper(parts[1]))
	request := map[string]interface{}{
		"RegionId":       client.RegionId,
		"LoadBalancerId": parts[0],
		"ListenerPort":   parts[2],
	}

	var object map[string]interface{}
	err = retry.RetryContext(ctx, time.Duration(listenerAclAttachmentTimeoutSeconds)*time.Second, func() *retry.RetryError {
		if err := ctx.Err(); err != nil {
			return retry.NonRetryableError(err)
		}
		response, err := client.RpcPost("Slb", string(connectivity.ApiVersion20140515), action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"The specified resource does not exist"}) {
				return retry.NonRetryableError(errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR))
			}
			if isTransientListenerError(err) {
				if err := waitListenerAclWithContext(ctx, 3*time.Second); err != nil {
					return retry.NonRetryableError(err)
				}
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		object = response
		if port, ok := object["ListenerPort"]; ok && strings.TrimSpace(fmt.Sprint(port)) != "" && fmt.Sprint(port) != "0" {
			return nil
		}
		return retry.RetryableError(errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("Listener", listenerId)))
	})
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, listenerId, action, errs.AlibabaCloudSdkGoERROR)
	}
	return object, nil
}

// isTransientListenerError covers both the transport level retries and the
// listener level busy codes that the Set*ListenerAttribute actions return while
// another change is still in progress. Quota, invalid ACL and IP version
// conflicts are terminal and stay non retryable.
func isTransientListenerError(err error) bool {
	if helper.NeedRetry(err) {
		return true
	}
	return errs.IsExpectedErrors(err, listenerTransientErrors)
}

func waitListenerAclWithContext(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func listenerAclAttachmentState(listenerId, aclId string, aclType types.String, object map[string]interface{}) listenerAclAttachmentResourceModel {
	model := listenerAclAttachmentResourceModel{
		Id:         types.StringValue(listenerId),
		ListenerId: types.StringValue(listenerId),
		AclId:      types.StringValue(aclId),
		AclType:    aclType,
	}
	if v, ok := object["AclType"].(string); ok && v != "" {
		model.AclType = types.StringValue(v)
	}
	if v, ok := object["AclStatus"].(string); ok && v != "" {
		model.AclStatus = types.StringValue(v)
	} else {
		model.AclStatus = types.StringNull()
	}
	return model
}

// listenerAclIds returns every access control policy group reported by
// DescribeLoadBalancer*ListenerAttribute.
func listenerAclIds(object map[string]interface{}) []string {
	var ids []string
	seen := map[string]struct{}{}
	appendID := func(id string) {
		id = strings.TrimSpace(id)
		if id == "" {
			return
		}
		if _, ok := seen[id]; ok {
			return
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if v, ok := object["AclId"].(string); ok && v != "" {
		for _, id := range strings.Split(v, ",") {
			appendID(id)
		}
	}
	if v, ok := object["AclIds"].(map[string]interface{}); ok {
		if items, ok := v["AclId"].([]interface{}); ok {
			for _, item := range items {
				if id, ok := item.(string); ok {
					appendID(id)
				}
			}
		}
	}
	return ids
}

// listenerAclContains reports whether the configured policy group is part of
// the listener ACLs described by DescribeLoadBalancer*ListenerAttribute.
func listenerAclContains(object map[string]interface{}, aclId string) bool {
	for _, id := range listenerAclIds(object) {
		if id == aclId {
			return true
		}
	}
	return false
}

func parseListenerAclAttachmentId(id string) (listenerId, aclId string, err error) {
	parts := strings.Split(id, ":")
	if (len(parts) != 3 && len(parts) != 4) || parts[0] == "" || parts[1] == "" || parts[2] == "" || (len(parts) == 4 && parts[3] == "") {
		return "", "", fmt.Errorf("invalid %s ID %s, expected <load_balancer_id>:<protocol>:<listener_port>[:<acl_id>]", listenerAclAttachmentTypeName, id)
	}
	return strings.Join(parts[:3], ":"), strings.Join(parts[3:], ":"), nil
}

func parseListenerPort(port string) (int, error) {
	return strconv.Atoi(port)
}

func isSupportedListenerProtocol(protocol string) bool {
	switch helper.Protocol(protocol) {
	case helper.Http, helper.Https, helper.Tcp, helper.Udp:
		return true
	}
	return false
}

// listenerAttributeAction maps the protocol part of a listener ID to the
// protocol specific SetLoadBalancer*ListenerAttribute action.
func listenerAttributeAction(protocol string) (string, error) {
	if !isSupportedListenerProtocol(protocol) {
		return "", fmt.Errorf("unsupported listener protocol %q", protocol)
	}
	return fmt.Sprintf("SetLoadBalancer%sListenerAttribute", strings.ToUpper(protocol)), nil
}
