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
	slbListenerIdPattern                = `^[^:]+:(http|https|tcp|udp):[1-9][0-9]*$`
)

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
						regexp.MustCompile(slbListenerIdPattern),
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
	if err := updateListenerAcl(r.Meta(), listenerId, aclId, plan.AclType.ValueString(), true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: binding access control policy group %s", listenerAclAttachmentTypeName, aclId),
			errs.WrapErrorf(err, "Creating %s: binding access control policy group %s", listenerAclAttachmentTypeName, aclId).Error(),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeSlbListener(listenerId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeLoadBalancer*ListenerAttribute", listenerAclAttachmentTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeLoadBalancer*ListenerAttribute", listenerAclAttachmentTypeName).Error(),
		)
		return
	}

	model := listenerAclAttachmentState(listenerId, aclId, object)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *listenerAclAttachmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state listenerAclAttachmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	listenerId, aclId, err := parseListenerAclAttachmentId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing resource ID", listenerAclAttachmentTypeName),
			err.Error(),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeSlbListener(listenerId)
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
	if aclId == "" {
		// Imported without the ACL part: adopt the ACL that is bound now.
		aclId = listenerFirstAclId(object)
		if aclId == "" {
			resp.State.RemoveResource(ctx)
			return
		}
	} else if !listenerAclContains(object, aclId) {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, listenerAclAttachmentState(listenerId, aclId, object))...)
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
		if err := updateListenerAcl(r.Meta(), listenerId, aclId, plan.AclType.ValueString(), true); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: binding access control policy group %s", listenerAclAttachmentTypeName, aclId),
				errs.WrapErrorf(err, "Updating %s: binding access control policy group %s", listenerAclAttachmentTypeName, aclId).Error(),
			)
			return
		}
	}

	object, err := NewService(r.Meta()).DescribeSlbListener(listenerId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeLoadBalancer*ListenerAttribute", listenerAclAttachmentTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeLoadBalancer*ListenerAttribute", listenerAclAttachmentTypeName).Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, listenerAclAttachmentState(listenerId, aclId, object))...)
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

	service := NewService(r.Meta())
	if _, err := service.DescribeSlbListener(listenerId); err != nil {
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

	if err := updateListenerAcl(r.Meta(), listenerId, "", "", false); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: disabling listener access control", listenerAclAttachmentTypeName),
			errs.WrapErrorf(err, "Deleting %s: disabling listener access control", listenerAclAttachmentTypeName).Error(),
		)
		return
	}

	// Wait until the access control status is switched off. The server may keep
	// the previous AclId after disabling, so only AclStatus is asserted.
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(time.Duration(listenerAclAttachmentTimeoutSeconds)*time.Second, func() *retry.RetryError {
		object, err := service.DescribeSlbListener(listenerId)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil
			}
			return retry.NonRetryableError(err)
		}
		if status, _ := object["AclStatus"].(string); !strings.EqualFold(status, "on") {
			return nil
		}
		wait()
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
	listenerId, _, err := parseListenerAclAttachmentId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Importing %s: parsing import ID", listenerAclAttachmentTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), listenerId)...)
}

// updateListenerAcl only sends the ACL related parameters of the protocol
// specific SetLoadBalancer*ListenerAttribute API, so other listener settings
// (bandwidth, health check, certificates and so on) are left untouched.
func updateListenerAcl(client *connectivity.AliyunClient, listenerId, aclId, aclType string, enable bool) error {
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
	if !isSupportedListenerProtocol(protocol) {
		return errs.WrapError(fmt.Errorf("unsupported listener protocol %s in listener ID %s", protocol, listenerId))
	}

	action := fmt.Sprintf("SetLoadBalancer%sListenerAttribute", strings.ToUpper(protocol))
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(time.Duration(listenerAclAttachmentTimeoutSeconds)*time.Second, func() *retry.RetryError {
		response, err = client.RpcPost("Slb", string(connectivity.ApiVersion20140515), action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, listenerId, action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

func listenerAclAttachmentState(listenerId, aclId string, object map[string]interface{}) listenerAclAttachmentResourceModel {
	model := listenerAclAttachmentResourceModel{
		Id:         types.StringValue(listenerId),
		ListenerId: types.StringValue(listenerId),
		AclId:      types.StringValue(aclId),
	}
	if v, ok := object["AclType"].(string); ok && v != "" {
		model.AclType = types.StringValue(v)
	} else {
		model.AclType = types.StringNull()
	}
	if v, ok := object["AclStatus"].(string); ok && v != "" {
		model.AclStatus = types.StringValue(v)
	} else {
		model.AclStatus = types.StringNull()
	}
	return model
}

// listenerAclContains reports whether the configured policy group is part of
// the listener ACLs described by DescribeLoadBalancer*ListenerAttribute.
func listenerAclContains(object map[string]interface{}, aclId string) bool {
	if v, ok := object["AclId"].(string); ok && v != "" {
		for _, id := range strings.Split(v, ",") {
			if strings.TrimSpace(id) == aclId {
				return true
			}
		}
	}
	if v, ok := object["AclIds"].(map[string]interface{}); ok {
		if items, ok := v["AclId"].([]interface{}); ok {
			for _, item := range items {
				if id, ok := item.(string); ok && id == aclId {
					return true
				}
			}
		}
	}
	return false
}

// listenerFirstAclId returns the first policy group reported by
// DescribeLoadBalancer*ListenerAttribute, used when the resource is imported
// without the optional ACL part.
func listenerFirstAclId(object map[string]interface{}) string {
	if v, ok := object["AclId"].(string); ok && v != "" {
		for _, id := range strings.Split(v, ",") {
			if id = strings.TrimSpace(id); id != "" {
				return id
			}
		}
	}
	if v, ok := object["AclIds"].(map[string]interface{}); ok {
		if items, ok := v["AclId"].([]interface{}); ok {
			for _, item := range items {
				if id, ok := item.(string); ok && id != "" {
					return id
				}
			}
		}
	}
	return ""
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
