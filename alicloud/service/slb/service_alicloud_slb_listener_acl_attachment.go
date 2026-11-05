package slb

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

// The operation budget, confirm window and poll interval are variables so unit
// tests can shorten them without changing production behaviour.
var (
	listenerAclAttachmentOperationTimeout = 300 * time.Second
	listenerAclAttachmentConfirmTimeout   = 30 * time.Second
	listenerAclAttachmentPollInterval     = 3 * time.Second
)

// errListenerAclAttachmentListeningNotVisible marks a Describe response that does
// not carry the listener port yet; it is retried like a transient error.
var errListenerAclAttachmentListeningNotVisible = errors.New("listener is not visible yet")

// errListenerAclAttachmentWriteUnknown marks a write whose transport level
// error cannot prove that the request never reached the service. The caller
// reconciles the listener before deciding whether the binding is managed.
var errListenerAclAttachmentWriteUnknown = errors.New("listener ACL write outcome unknown")

// listenerTransientErrors are the listener level busy codes the
// Set*ListenerAttribute actions return while another change is still being
// applied. Quota, invalid ACL and IP version conflicts are not listed and fail
// fast.
var listenerTransientErrors = append(append([]string{}, errs.SlbIsBusy...), "TLSPolicyConfiguring")

// listenerAclAttachmentRpcPost is the single transport entry for the two SLB
// actions used by this resource. The actual HTTP request carries the lifecycle
// context and the SDK level retry is disabled, so a cancelled operation aborts
// the round trip instead of leaving it running in the background.
var listenerAclAttachmentRpcPost = func(ctx context.Context, client *connectivity.AliyunClient, action string, request map[string]interface{}) (map[string]interface{}, error) {
	return client.RpcPostWithContext(ctx, "Slb", string(connectivity.ApiVersion20140515), action, nil, request, false)
}

// listenerAclAttachmentApi is the SLB surface of the attachment lifecycle. The
// concrete implementation below wraps the provider client; keeping the surface
// explicit lets the lifecycle be exercised without a live endpoint.
type listenerAclAttachmentApi interface {
	Describe(ctx context.Context, listenerId string) (map[string]interface{}, error)
	SetAcl(ctx context.Context, listenerId, aclId, aclType string, enable bool) error
}

type listenerAclAttachmentClient struct {
	client *connectivity.AliyunClient
}

func newListenerAclAttachmentApi(client *connectivity.AliyunClient) listenerAclAttachmentApi {
	return listenerAclAttachmentClient{client: client}
}

func (a listenerAclAttachmentClient) Describe(ctx context.Context, listenerId string) (map[string]interface{}, error) {
	return describeListenerWithContext(ctx, a.client, listenerId)
}

func (a listenerAclAttachmentClient) SetAcl(ctx context.Context, listenerId, aclId, aclType string, enable bool) error {
	return updateListenerAcl(ctx, a.client, listenerId, aclId, aclType, enable)
}

// updateListenerAcl sends only the ACL related parameters of the protocol
// specific SetLoadBalancer*ListenerAttribute API, so other listener settings
// (bandwidth, health check, certificates and so on) are left untouched.
func updateListenerAcl(ctx context.Context, client *connectivity.AliyunClient, listenerId, aclId, aclType string, enable bool) error {
	parts, err := parseListenerAclListenerId(listenerId)
	if err != nil {
		return errs.WrapError(err)
	}
	action, err := listenerAttributeAction(parts[1])
	if err != nil {
		return errs.WrapErrorf(err, "invalid listener ID %s", listenerId)
	}
	port, err := parseListenerPort(parts[2])
	if err != nil {
		return errs.WrapErrorf(err, "invalid listener port in listener ID %s", listenerId)
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
	response, err = helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return listenerAclAttachmentRpcPost(ctx, client, action, request)
	}, listenerAclAttachmentOperationTimeout, isTransientListenerError)
	logs.AddDebug(action, response, request)
	if err != nil {
		wrapped := errs.WrapErrorf(err, errs.DefaultErrorMsg, listenerId, action, errs.AlibabaCloudSdkGoERROR)
		if isUnknownListenerAclWriteOutcome(err) {
			return fmt.Errorf("%w: %v", errListenerAclAttachmentWriteUnknown, wrapped)
		}
		return wrapped
	}
	return nil
}

// describeListenerWithContext reads the protocol specific listener attributes
// with the CRUD context: the SDK level retry stays off and the shared helper
// owns both the backoff budget and cancellation.
func describeListenerWithContext(ctx context.Context, client *connectivity.AliyunClient, listenerId string) (map[string]interface{}, error) {
	parts, err := parseListenerAclListenerId(listenerId)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	action, err := listenerAttributeDescribeAction(parts[1])
	if err != nil {
		return nil, errs.WrapErrorf(err, "invalid listener ID %s", listenerId)
	}
	port, err := parseListenerPort(parts[2])
	if err != nil {
		return nil, errs.WrapErrorf(err, "invalid listener port in listener ID %s", listenerId)
	}
	request := map[string]interface{}{
		"RegionId":       client.RegionId,
		"LoadBalancerId": parts[0],
		"ListenerPort":   port,
	}

	var object map[string]interface{}
	object, err = helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		response, err := listenerAclAttachmentRpcPost(ctx, client, action, request)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"The specified resource does not exist"}) {
				return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
			}
			return nil, err
		}
		logs.AddDebug(action, response, request)
		if value, ok := response["ListenerPort"]; !ok || strings.TrimSpace(fmt.Sprint(value)) == "" || fmt.Sprint(value) == "0" {
			return response, fmt.Errorf("%w: %s", errListenerAclAttachmentListeningNotVisible, listenerId)
		}
		return response, nil
	}, listenerAclAttachmentOperationTimeout, isTransientListenerError)
	if err != nil {
		return nil, err
	}
	return object, nil
}

// isTransientListenerError classifies the errors that may be retried inside the
// operation budget: transport level failures, the shared SLB busy codes and a
// Describe response that does not expose the listener yet. A missing listener,
// a quota violation, an invalid ACL and an IP version conflict are terminal.
func isTransientListenerError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, errListenerAclAttachmentListeningNotVisible) {
		return true
	}
	if errs.NotFoundError(err) {
		return false
	}
	if helper.NeedRetry(err) {
		return true
	}
	return errs.IsExpectedErrors(err, listenerTransientErrors)
}

// isUnknownListenerAclWriteOutcome reports whether a failed write may still have
// been applied: only transport level failures leave the result unknown. A
// business error such as AclNotExist, IpVersionConflict or a quota code proves
// that the request did not change the listener.
func isUnknownListenerAclWriteOutcome(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, errListenerAclAttachmentWriteUnknown) {
		return true
	}
	return helper.NeedRetry(listenerAclErrorCause(err))
}

// isUnresolvedListenerAclWrite reports whether a failed write may still have
// been applied: a transport level failure or a write budget that expired while
// the parent operation is still alive. The caller reconciles those against the
// listener. A cancelled parent operation keeps the plain failure because no
// further request can be issued.
func isUnresolvedListenerAclWrite(parent, writeCtx context.Context, err error) bool {
	if parent.Err() != nil {
		return false
	}
	if writeCtx.Err() != nil {
		return true
	}
	if isUnknownListenerAclWriteOutcome(err) {
		return true
	}
	return errors.Is(listenerAclErrorCause(err), context.DeadlineExceeded)
}

// listenerAclErrorCause exposes the original error behind the provider's
// ComplexError wrapper so the transport checks keep working on wrapped errors.
func listenerAclErrorCause(err error) error {
	if complexErr, ok := err.(*errs.ComplexError); ok && complexErr.Cause != nil {
		return complexErr.Cause
	}
	return err
}

// listenerAclAttachmentConfirms reports whether the listener proves the planned
// binding: access control must be enabled, the policy group must be a member,
// and the reported type must match the plan when the response carries one.
func listenerAclAttachmentConfirms(object map[string]interface{}, aclId, aclType string) bool {
	if object == nil {
		return false
	}
	if status, _ := object["AclStatus"].(string); !strings.EqualFold(status, "on") {
		return false
	}
	if !listenerAclContains(object, aclId) {
		return false
	}
	// While access control is on the Describe contract reports AclType. A
	// missing, empty or non-string value cannot prove that the requested type
	// took effect, so the write stays unconfirmed.
	value, ok := object["AclType"].(string)
	if !ok || strings.TrimSpace(value) == "" {
		return false
	}
	return strings.EqualFold(value, aclType)
}

// waitListenerAclAttachmentConverged reads the listener until it reports the
// planned binding, so a stale response cannot store the previous type against a
// new plan. It polls on the caller goroutine and returns converged=false when
// the operation budget ends first.
func waitListenerAclAttachmentConverged(ctx context.Context, api listenerAclAttachmentApi, listenerId, aclId, aclType string) (map[string]interface{}, bool, error) {
	var last map[string]interface{}
	for {
		if err := ctx.Err(); err != nil {
			return last, false, nil
		}
		object, err := api.Describe(ctx, listenerId)
		if err != nil {
			return last, false, err
		}
		last = object
		if listenerAclAttachmentConfirms(object, aclId, aclType) {
			return last, true, nil
		}
		if err := waitListenerAclWithContext(ctx, listenerAclAttachmentPollInterval); err != nil {
			return last, false, nil
		}
	}
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
	if value, ok := object["AclId"].(string); ok && value != "" {
		for _, id := range strings.Split(value, ",") {
			appendID(id)
		}
	}
	if value, ok := object["AclIds"].(map[string]interface{}); ok {
		if items, ok := value["AclId"].([]interface{}); ok {
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

// parseListenerAclAttachmentId parses the resource and import ID. Both the three
// segment listener form and the four segment form that selects the managed
// policy group are accepted; the listener part always goes through the same
// protocol and port validation as the configured listener_id.
func parseListenerAclAttachmentId(id string) (listenerId, aclId string, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 3 && len(parts) != 4 {
		return "", "", fmt.Errorf("invalid %s ID %s, expected <load_balancer_id>:<protocol>:<listener_port>[:<acl_id>]", listenerAclAttachmentTypeName, id)
	}
	listenerId = strings.Join(parts[:3], ":")
	if _, err := parseListenerAclListenerId(listenerId); err != nil {
		return "", "", fmt.Errorf("invalid %s ID %s: %v", listenerAclAttachmentTypeName, id, err)
	}
	if len(parts) == 4 {
		if err := validateSingleListenerAclId(parts[3]); err != nil {
			return "", "", fmt.Errorf("invalid %s ID %s: %v", listenerAclAttachmentTypeName, id, err)
		}
		aclId = parts[3]
	}
	return listenerId, aclId, nil
}

// parseListenerAclListenerId validates the listener identity used by both the
// configuration and the import ID.
func parseListenerAclListenerId(id string) ([]string, error) {
	parts := strings.Split(id, ":")
	if len(parts) != 3 || parts[0] == "" {
		return nil, fmt.Errorf("invalid listener ID %s, expected <load_balancer_id>:<protocol>:<listener_port>", id)
	}
	if !isSupportedListenerProtocol(parts[1]) {
		return nil, fmt.Errorf("invalid listener ID %s, protocol must be one of http, https, tcp and udp", id)
	}
	if _, err := parseListenerPort(parts[2]); err != nil {
		return nil, fmt.Errorf("invalid listener ID %s, %v", id, err)
	}
	return parts, nil
}

// validateSingleListenerAclId rejects the combined ACL form: this resource
// manages exactly one policy group, and a comma separated value would bind
// several groups that the membership check can never match.
func validateSingleListenerAclId(aclId string) error {
	if strings.TrimSpace(aclId) == "" {
		return fmt.Errorf("acl_id must not be empty")
	}
	if strings.ContainsAny(aclId, ", \t\r\n") {
		return fmt.Errorf("acl_id %q must be a single access control policy group ID; comma separated IDs are not supported", aclId)
	}
	return nil
}

func parseListenerPort(port string) (int, error) {
	value, err := strconv.Atoi(port)
	if err != nil {
		return 0, fmt.Errorf("listener port %q must be an integer", port)
	}
	if value < 1 || value > 65535 {
		return 0, fmt.Errorf("listener port %d must be between 1 and 65535", value)
	}
	return value, nil
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

// listenerAttributeDescribeAction maps the protocol part of a listener ID to
// the protocol specific DescribeLoadBalancer*ListenerAttribute action. It is
// kept separate from the write action so both directions are covered.
func listenerAttributeDescribeAction(protocol string) (string, error) {
	if !isSupportedListenerProtocol(protocol) {
		return "", fmt.Errorf("unsupported listener protocol %q", protocol)
	}
	return fmt.Sprintf("DescribeLoadBalancer%sListenerAttribute", strings.ToUpper(protocol)), nil
}
