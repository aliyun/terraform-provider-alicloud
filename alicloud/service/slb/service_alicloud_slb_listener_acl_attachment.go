package slb

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
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
	return client.RpcPostWithContext(ctx, "Slb", string(connectivity.ApiVersion20140515), action, nil, request)
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
	// The debug log receives redacted copies: a gateway that echoes the request in
	// its response must not be able to write the signed URL (access key, signature
	// and security token) into the provider log.
	logs.AddDebug(action, connectivity.RedactLogValue(response), connectivity.RedactLogValue(request))
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
			if listenerAclListenerGone(err) {
				return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
			}
			return nil, err
		}
		logs.AddDebug(action, connectivity.RedactLogValue(response), connectivity.RedactLogValue(request))
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
	if listenerAclListenerGone(err) {
		return false
	}
	if helper.NeedRetry(err) {
		return true
	}
	code := listenerAclBusinessCode(err)
	if code == "" {
		return false
	}
	for _, busy := range listenerTransientErrors {
		if code == busy {
			return true
		}
	}
	return false
}

// listenerAclSdkError returns the tea SDK error behind the provider's wrappers
// (the retry helper wraps once, the lifecycle wraps again) so the business code
// and message survive every layer.
func listenerAclSdkError(err error) *tea.SDKError {
	for err != nil {
		if sdkErr, ok := err.(*tea.SDKError); ok {
			return sdkErr
		}
		switch typed := err.(type) {
		case *errs.ComplexError:
			// ComplexError predates errors.Unwrap and carries its cause directly.
			err = typed.Cause
		case interface{ Unwrap() error }:
			err = typed.Unwrap()
		default:
			return nil
		}
	}
	return nil
}

// listenerAclBusinessCode returns the SLB business code, or an empty string when
// the response carried none. An empty code must never be matched against a code
// list: the shared errs.IsExpectedErrors helper treats an empty code as a
// wildcard, which would classify an unrelated gateway rejection (for example an
// empty-bodied 403) as a known SLB result.
func listenerAclBusinessCode(err error) string {
	sdkErr := listenerAclSdkError(err)
	if sdkErr == nil {
		return ""
	}
	return strings.TrimSpace(tea.StringValue(sdkErr.Code))
}

// listenerAclListenerGone reports whether err proves that the addressed SLB load
// balancer or listener no longer exists. Only a decodable business error may
// prove it: a bare HTTP status (a gateway 404, an undecodable body) or a
// response without a business code does not show that the SLB entity is gone,
// and treating it as gone would drop live state on Read or report a successful
// destroy. The two Describe actions of this resource address only the load
// balancer and the listener, so a not-found result refers to one of them.
func listenerAclListenerGone(err error) bool {
	sdkErr := listenerAclSdkError(err)
	if sdkErr == nil {
		return false
	}
	code := strings.TrimSpace(tea.StringValue(sdkErr.Code))
	if code == "" {
		return false
	}
	if strings.Contains(strings.ToLower(code), "notfound") {
		return true
	}
	message := strings.ToLower(strings.TrimSpace(tea.StringValue(sdkErr.Message)))
	return strings.Contains(message, "does not exist") || strings.Contains(message, "not found")
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
	cause := listenerAclErrorCause(err)
	// A response that reached the service but could not be decoded cannot prove
	// whether the write applied, so it is reconciled like a transport failure.
	// A 4xx status is the opposite: the service refused the request even though
	// its body was unreadable, so the outcome is known.
	if errors.Is(cause, connectivity.ErrUnreadableResponse) {
		return true
	}
	if errors.Is(cause, connectivity.ErrRpcRejectedResponse) {
		return false
	}
	// A transport level failure (dial, TLS, reset or timeout) never proves that
	// the request was not applied. It is classified by type rather than by the
	// scheme dependent error message, so an http test endpoint behaves like the
	// https endpoint used in production.
	var urlErr *url.Error
	if errors.As(cause, &urlErr) {
		return true
	}
	return helper.NeedRetry(cause)
}

// isUnresolvedListenerAclWrite reports whether a failed write may still have
// been applied: a transport level failure or a write budget that expired while
// the parent operation is still alive. The caller reconciles those against the
// listener. A cancelled parent operation keeps the plain failure because no
// further request can be issued.
//
// A decoded service answer is checked first: it proves the request was answered,
// so the terminal classification wins over a write budget that expired at the
// same moment. Reconciling such a rejection could otherwise adopt a pre-existing
// matching binding as the result of a refused write.
func isUnresolvedListenerAclWrite(parent, writeCtx context.Context, err error) bool {
	if parent.Err() != nil {
		return false
	}
	if isDefinitiveListenerAclWriteOutcome(err) {
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

// isDefinitiveListenerAclWriteOutcome reports whether the service answered the
// write with a decoded error: a 4xx business error, or an explicit rejection
// whose body could not be decoded. Both prove that the request was not applied.
// A 5xx answer, an unreadable 2xx body and a transport level failure stay
// unknown and are reconciled instead.
func isDefinitiveListenerAclWriteOutcome(err error) bool {
	if err == nil {
		return false
	}
	// The write error reaches this package wrapped twice (the retry helper and the
	// action wrapper), and ComplexError does not implement errors.Unwrap, so the
	// whole chain has to be walked explicitly.
	if listenerAclErrorMatches(err, connectivity.ErrRpcRejectedResponse) {
		return true
	}
	sdkErr := listenerAclSdkError(err)
	if sdkErr == nil {
		return false
	}
	status := tea.IntValue(sdkErr.StatusCode)
	if status >= 400 && status < 500 {
		return true
	}
	// A decoded error without an HTTP status still carries the service verdict.
	return status == 0 && strings.TrimSpace(tea.StringValue(sdkErr.Code)) != ""
}

// listenerAclErrorMatches reports whether target appears in the wrapped error
// chain, including the provider's ComplexError layers that carry their cause in a
// field instead of implementing errors.Unwrap.
func listenerAclErrorMatches(err, target error) bool {
	for err != nil {
		if errors.Is(err, target) {
			return true
		}
		switch typed := err.(type) {
		case *errs.ComplexError:
			err = typed.Cause
		case interface{ Unwrap() error }:
			err = typed.Unwrap()
		default:
			return false
		}
	}
	return false
}

// listenerAclErrorCause exposes the original error behind the provider's
// ComplexError wrapper so the transport checks keep working on wrapped errors.
func listenerAclErrorCause(err error) error {
	if complexErr, ok := err.(*errs.ComplexError); ok && complexErr.Cause != nil {
		return complexErr.Cause
	}
	return err
}

// listenerAclStatus reads the AclStatus field of a Describe response. The
// Describe contract always reports it as `on` or `off`; an absent or non string
// value means the response is incomplete and can neither confirm a binding nor
// prove that access control is off.
func listenerAclStatus(object map[string]interface{}) (status string, ok bool) {
	if object == nil {
		return "", false
	}
	value, ok := object["AclStatus"].(string)
	if !ok || strings.TrimSpace(value) == "" {
		return "", false
	}
	return value, true
}

// listenerAclAttachmentConfirms reports whether the listener proves the planned
// binding: access control must be enabled, the policy group must be a member,
// and the reported type must match the plan when the response carries one.
func listenerAclAttachmentConfirms(object map[string]interface{}, aclId, aclType string) bool {
	if object == nil {
		return false
	}
	status, ok := listenerAclStatus(object)
	if !ok || !strings.EqualFold(status, "on") {
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

// listenerAclAttachmentObserved reports whether a Describe response carries
// enough usable fields to either confirm or contradict a planned binding. A
// response without a usable status, or one that reports an enabled listener
// without a usable type, without reporting the membership at all, or with
// malformed members, is incomplete: it can neither prove nor disprove the
// planned change.
func listenerAclAttachmentObserved(object map[string]interface{}) bool {
	if object == nil {
		return false
	}
	status, ok := listenerAclStatus(object)
	if !ok {
		return false
	}
	if !strings.EqualFold(status, "on") {
		return true
	}
	if value, ok := object["AclType"].(string); !ok || strings.TrimSpace(value) == "" {
		return false
	}
	_, membersReported, membersComplete := listenerAclIds(object)
	return membersReported && membersComplete
}

// waitListenerAclAttachmentConverged reads the listener until it reports the
// planned binding, so a stale response cannot store the previous type against a
// new plan. It polls on the caller goroutine and returns converged=false when
// the operation budget ends first. The complete flag then reports whether the
// last response was usable enough to prove that the binding differs from the
// plan: an incomplete response is a pending confirmation, a complete one is a
// contradiction that must not be reported as success.
func waitListenerAclAttachmentConverged(ctx context.Context, api listenerAclAttachmentApi, listenerId, aclId, aclType string) (map[string]interface{}, bool, bool, error) {
	var last map[string]interface{}
	for {
		if err := ctx.Err(); err != nil {
			return last, false, listenerAclAttachmentObserved(last), nil
		}
		object, err := api.Describe(ctx, listenerId)
		if err != nil {
			return last, false, listenerAclAttachmentObserved(last), err
		}
		last = object
		if listenerAclAttachmentConfirms(object, aclId, aclType) {
			return last, true, true, nil
		}
		if err := waitListenerAclWithContext(ctx, listenerAclAttachmentPollInterval); err != nil {
			return last, false, listenerAclAttachmentObserved(last), nil
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
// DescribeLoadBalancer*ListenerAttribute. The reported flag says whether the
// response carried the member field at all, and complete says whether every
// reported member could be read.
//
// An enabled listener must report at least one member, so a response that omits
// the member fields (or serialises them as JSON null) has not reported the
// membership at all: it can neither prove that the planned group is bound nor
// that it is missing. Only a reported and readable list is usable as evidence.
func listenerAclIds(object map[string]interface{}) (ids []string, reported, complete bool) {
	seen := map[string]struct{}{}
	complete = true
	appendID := func(id string) {
		id = strings.TrimSpace(id)
		if id == "" {
			complete = false
			return
		}
		if _, ok := seen[id]; ok {
			return
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if value, ok := object["AclId"]; ok && value != nil {
		if raw, ok := value.(string); ok {
			reported = true
			for _, id := range strings.Split(raw, ",") {
				appendID(id)
			}
		} else {
			complete = false
		}
	}
	if value, ok := object["AclIds"]; ok && value != nil {
		container, ok := value.(map[string]interface{})
		if !ok {
			// A present but unusable container was reported and is incomplete.
			reported = true
			complete = false
		} else if members, ok := container["AclId"]; ok {
			reported = true
			switch typed := members.(type) {
			case []interface{}:
				for _, item := range typed {
					id, ok := item.(string)
					if !ok {
						complete = false
						continue
					}
					appendID(id)
				}
			case string:
				// A single member may be reported without its one element array.
				appendID(typed)
			default:
				complete = false
			}
		}
	}
	return ids, reported, complete
}

// listenerAclContains reports whether the configured policy group is part of
// the listener ACLs described by DescribeLoadBalancer*ListenerAttribute.
func listenerAclContains(object map[string]interface{}, aclId string) bool {
	ids, _, _ := listenerAclIds(object)
	for _, id := range ids {
		if id == aclId {
			return true
		}
	}
	return false
}

// listenerAclIdsForMessage lists the readable policy groups for a diagnostic;
// the completeness of the member list is irrelevant there.
func listenerAclIdsForMessage(object map[string]interface{}) []string {
	ids, _, _ := listenerAclIds(object)
	return ids
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
