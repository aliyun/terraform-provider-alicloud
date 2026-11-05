package slb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
)

// These unit tests cover the transport layer of the listener ACL attachment:
// request construction, the retry classification, same-port protocol dispatch
// and the not-found behaviour. The resource lifecycle is exercised through the
// public resource factory in resource_alicloud_slb_listener_acl_attachment_test.go.

// requestCapture replaces the transport entry for the duration of one test.
func requestCapture(t *testing.T, handler func(ctx context.Context, action string, request map[string]interface{}) (map[string]interface{}, error)) *int {
	t.Helper()
	original := listenerAclAttachmentRpcPost
	calls := 0
	listenerAclAttachmentRpcPost = func(ctx context.Context, client *connectivity.AliyunClient, action string, request map[string]interface{}) (map[string]interface{}, error) {
		calls++
		return handler(ctx, action, request)
	}
	t.Cleanup(func() { listenerAclAttachmentRpcPost = original })
	return &calls
}

func TestUnitListenerAclAttachmentWriteRetry(t *testing.T) {
	client := &connectivity.AliyunClient{RegionId: "cn-hangzhou"}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	t.Run("busy then success retries and converges", func(t *testing.T) {
		var requests []map[string]interface{}
		sawDeadline := false
		calls := requestCapture(t, func(callCtx context.Context, action string, request map[string]interface{}) (map[string]interface{}, error) {
			if _, ok := callCtx.Deadline(); ok {
				sawDeadline = true
			}
			requests = append(requests, request)
			if len(requests) == 1 {
				return nil, &tea.SDKError{Code: tea.String("OperationBusy"), Message: tea.String("code: 400, OperationBusy"), Data: tea.String("{}")}
			}
			return map[string]interface{}{}, nil
		})
		if err := updateListenerAcl(ctx, client, "lb-123:http:80", "acl-1", "white", true); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if *calls != 2 {
			t.Fatalf("expected two attempts (busy then success), got %d", *calls)
		}
		if !sawDeadline {
			t.Fatal("the transport must receive the lifecycle context with its deadline")
		}
		request := requests[0]
		if request["LoadBalancerId"] != "lb-123" || fmt.Sprint(request["ListenerPort"]) != "80" {
			t.Fatalf("request identity is wrong: %v", request)
		}
		if request["AclId"] != "acl-1" || request["AclType"] != "white" || request["AclStatus"] != "on" {
			t.Fatalf("request ACL parameters are wrong: %v", request)
		}
	})

	t.Run("terminal error makes exactly one request", func(t *testing.T) {
		calls := requestCapture(t, func(callCtx context.Context, action string, request map[string]interface{}) (map[string]interface{}, error) {
			return nil, &tea.SDKError{Code: tea.String("IpVersionConflict"), Message: tea.String("ip version conflict"), Data: tea.String("{}")}
		})
		if err := updateListenerAcl(ctx, client, "lb-123:http:80", "acl-1", "white", true); err == nil {
			t.Fatal("expected the terminal error to be returned")
		}
		if *calls != 1 {
			t.Fatalf("terminal error must not be retried, got %d calls", *calls)
		}
	})

	t.Run("quota error goes through the retry path once", func(t *testing.T) {
		calls := requestCapture(t, func(callCtx context.Context, action string, request map[string]interface{}) (map[string]interface{}, error) {
			return nil, &tea.SDKError{Code: tea.String("QuotaLimitExceeds.AclListenerOverLimit"), Message: tea.String("quota exceeded"), Data: tea.String("{}")}
		})
		if err := updateListenerAcl(ctx, client, "lb-123:http:80", "acl-1", "white", true); err == nil {
			t.Fatal("expected the quota error to be returned")
		}
		if *calls != 1 {
			t.Fatalf("quota error must make exactly one request, got %d", *calls)
		}
	})

	t.Run("cancelled context sends no request", func(t *testing.T) {
		calls := requestCapture(t, func(callCtx context.Context, action string, request map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{}, nil
		})
		cancelled, cancelCancelled := context.WithCancel(context.Background())
		cancelCancelled()
		if err := updateListenerAcl(cancelled, client, "lb-123:http:80", "acl-1", "white", true); err == nil {
			t.Fatal("expected the cancelled context to fail the write")
		}
		if *calls != 0 {
			t.Fatalf("cancelled context must not send a request, got %d calls", *calls)
		}
	})
}

// sameListenerTransport records the actions and requests issued for the same
// load balancer and port through two different protocols.
type sameListenerTransport struct {
	actions  []string
	requests []map[string]interface{}
}

func (s *sameListenerTransport) install(t *testing.T) {
	t.Helper()
	requestCapture(t, func(callCtx context.Context, action string, request map[string]interface{}) (map[string]interface{}, error) {
		s.actions = append(s.actions, action)
		s.requests = append(s.requests, request)
		return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": "acl-1"}, nil
	})
}

func TestUnitListenerAclAttachmentSamePortProtocolDispatch(t *testing.T) {
	client := &connectivity.AliyunClient{RegionId: "cn-hangzhou"}
	transport := &sameListenerTransport{}
	transport.install(t)

	for _, protocol := range []string{"http", "https"} {
		if err := updateListenerAcl(context.Background(), client, "lb-shared:"+protocol+":80", "acl-1", "white", true); err != nil {
			t.Fatalf("set %s listener: %s", protocol, err)
		}
		if _, err := describeListenerWithContext(context.Background(), client, "lb-shared:"+protocol+":80"); err != nil {
			t.Fatalf("describe %s listener: %s", protocol, err)
		}
	}

	want := []string{
		"SetLoadBalancerHTTPListenerAttribute",
		"DescribeLoadBalancerHTTPListenerAttribute",
		"SetLoadBalancerHTTPSListenerAttribute",
		"DescribeLoadBalancerHTTPSListenerAttribute",
	}
	if len(transport.actions) != len(want) {
		t.Fatalf("expected %d requests, got %v", len(want), transport.actions)
	}
	for i, action := range want {
		if transport.actions[i] != action {
			t.Fatalf("request %d action = %q, want %q", i, transport.actions[i], action)
		}
	}
	for _, request := range transport.requests {
		if request["LoadBalancerId"] != "lb-shared" || fmt.Sprint(request["ListenerPort"]) != "80" {
			t.Fatalf("same port request identity is wrong: %v", request)
		}
	}
}

// The busy codes and the terminal codes come from the shared SLB contract; the
// table keeps the retry classification auditable when the shared list changes.
func TestUnitListenerAclAttachmentErrorClassification(t *testing.T) {
	for _, code := range listenerTransientErrors {
		err := &tea.SDKError{Code: tea.String(code), Message: tea.String(""), Data: tea.String("{}")}
		if !isTransientListenerError(err) {
			t.Fatalf("%s should be retryable", code)
		}
	}

	terminal := []string{
		"AclListenerOverLimit",
		"QuotaLimitExceeds.AclAttachedToListener",
		"QuotaLimitExceeds.AclListenerOverLimit",
		"IpVersionConflict",
		"AclNotExist",
		"InvalidParameter",
	}
	for _, code := range terminal {
		err := &tea.SDKError{Code: tea.String(code), Message: tea.String(""), Data: tea.String("{}")}
		if isTransientListenerError(err) {
			t.Fatalf("%s must not be retried", code)
		}
	}

	notFound := listenerAclTestNotFoundError()
	if isTransientListenerError(notFound) {
		t.Fatal("a missing listener must not be retried")
	}
	if !listenerAclListenerGone(notFound) {
		t.Fatal("a 404 with the SLB business code must stay recognizable as a missing listener")
	}

	// A bare status, an empty business code and a plain transport failure never
	// prove that the listener is gone. Only the transport failure stays retryable;
	// the others are terminal SLB rejections.
	notGone := []struct {
		name      string
		err       error
		retryable bool
	}{
		{name: "empty business code with a 403 status", err: &tea.SDKError{StatusCode: tea.Int(403), Code: tea.String(""), Message: tea.String("The specified resource does not exist."), Data: tea.String("{}")}},
		{name: "empty code and message with a 404 status", err: &tea.SDKError{StatusCode: tea.Int(404), Code: tea.String(""), Message: tea.String(""), Data: tea.String("{}")}},
		{name: "gateway 404 without a business body", err: fmt.Errorf("%w: status 404: unexpected EOF", connectivity.ErrRpcRejectedResponse)},
		{name: "plain transport failure", err: errors.New(`Post "https://slb.aliyuncs.com/": i/o timeout`), retryable: true},
	}
	for _, tc := range notGone {
		if listenerAclListenerGone(tc.err) {
			t.Fatalf("%s must not be classified as a missing listener", tc.name)
		}
		if got := isTransientListenerError(tc.err); got != tc.retryable {
			t.Fatalf("%s retryable = %v, want %v", tc.name, got, tc.retryable)
		}
	}
}

// listenerAclTestNotFoundError is the error shape that proves the SLB entity is
// gone: a decodable business error, not only a status.
func listenerAclTestNotFoundError() error {
	return &tea.SDKError{
		StatusCode: tea.Int(404),
		Code:       tea.String("InvalidLoadBalancerId.NotFound"),
		Message:    tea.String("The specified resource does not exist."),
		Data:       tea.String("{}"),
	}
}

// The provider debug log receives the response payload. A gateway that echoes the
// request in its response must not be able to write the signed URL into it.
func TestUnitListenerAclAttachmentDebugLogIsRedacted(t *testing.T) {
	var buffer bytes.Buffer
	previousWriter := log.Writer()
	log.SetOutput(&buffer)
	t.Cleanup(func() { log.SetOutput(previousWriter) })
	t.Setenv("DEBUG", "terraform")

	requestCapture(t, func(ctx context.Context, action string, request map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"RequestId":    "r-1",
			"ListenerPort": 80,
			"AclStatus":    "on",
			"AclType":      "white",
			"AclId":        "acl-1",
			"Echo":         "https://slb.aliyuncs.com/?AccessKeyId=unit-test-access-key&Signature=synthetic-signature&SecurityToken=synthetic-token",
			"QueryEcho":    "AccessKeyId=unit-test-access-key&Signature=synthetic-signature&SignatureNonce=synthetic-nonce&SecurityToken=synthetic-token",
		}, nil
	})
	client := &connectivity.AliyunClient{RegionId: "cn-hangzhou"}
	if _, err := describeListenerWithContext(context.Background(), client, "lb-123:http:80"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	output := buffer.String()
	if !strings.Contains(output, "DescribeLoadBalancerHTTPListenerAttribute") {
		t.Fatalf("the debug log did not record the call: %q", output)
	}
	// The parameter names may stay in the log; only their values are secret.
	for _, secret := range []string{"unit-test-access-key", "synthetic-signature", "synthetic-nonce", "synthetic-token"} {
		if strings.Contains(output, secret) {
			// Do not print the leaking line: it contains the protected value.
			t.Fatalf("the debug log leaked %q", secret)
		}
	}
	if !strings.Contains(output, "slb.aliyuncs.com") {
		t.Fatalf("the redaction dropped the usable response: %q", output)
	}
}

// A decoded service answer is terminal even when the write budget expires at the
// same moment: reconciling it could adopt a pre-existing matching binding as the
// result of a refused write.
func TestUnitListenerAclAttachmentWriteOutcomeClassification(t *testing.T) {
	terminal := &tea.SDKError{
		StatusCode: tea.Int(400),
		Code:       tea.String("QuotaLimitExceeds.AclListenerOverLimit"),
		Message:    tea.String("quota exceeded"),
		Data:       tea.String("{}"),
	}
	serverFailure := &tea.SDKError{
		StatusCode: tea.Int(500),
		Code:       tea.String("InternalError"),
		Message:    tea.String("code: 500, internal error"),
		Data:       tea.String("{}"),
	}
	rejection := fmt.Errorf("%w: status 400: unexpected EOF", connectivity.ErrRpcRejectedResponse)
	transport := errors.New(`Post "https://slb.aliyuncs.com/": i/o timeout`)
	// The write error reaches the resource wrapped exactly like production: the
	// retry helper and the action wrapper each add a ComplexError layer.
	productionRejection := errs.WrapErrorf(
		errs.WrapErrorf(rejection, "calling %s failed", "SetLoadBalancerHTTPListenerAttribute"),
		errs.DefaultErrorMsg, "lb-123:http:80", "SetLoadBalancerHTTPListenerAttribute", errs.AlibabaCloudSdkGoERROR,
	)
	// A transport deadline that is not the write budget keeps the outcome unknown.
	productionDeadline := errs.WrapErrorf(
		errs.WrapErrorf(fmt.Errorf("request failed: %w", context.DeadlineExceeded), "calling %s failed", "SetLoadBalancerHTTPListenerAttribute"),
		errs.DefaultErrorMsg, "lb-123:http:80", "SetLoadBalancerHTTPListenerAttribute", errs.AlibabaCloudSdkGoERROR,
	)
	writeBudgetTimeout := errors.New("timeout while waiting for state to become 'success' (timeout: 200ms)")
	liveCtx := context.Background()
	expiredCtx, cancelExpired := context.WithCancel(context.Background())
	cancelExpired()

	cases := []struct {
		name     string
		writeCtx context.Context
		err      error
		want     bool
	}{
		{"decoded terminal rejection", liveCtx, terminal, false},
		{"undecodable 4xx rejection", liveCtx, rejection, false},
		{"production wrapped rejection", liveCtx, productionRejection, false},
		{"production wrapped deadline", liveCtx, productionDeadline, true},
		{"transport failure", liveCtx, transport, true},
		{"5xx answer", liveCtx, serverFailure, true},
		// The retry helper reports an exhausted write budget as a generic timeout;
		// the expired budget itself proves that the outcome is unknown.
		{"expired write budget with a generic timeout", expiredCtx, writeBudgetTimeout, true},
		// The parent deadline ends the write context as well; the request may have
		// been in flight, so the planned identity is still kept.
		{"cancelled parent leaves the write unconfirmed", expiredCtx, transport, true},
	}
	for _, tc := range cases {
		if got := isUnconfirmedListenerAclWrite(tc.writeCtx, tc.err); got != tc.want {
			t.Fatalf("%s: isUnconfirmedListenerAclWrite = %v, want %v", tc.name, got, tc.want)
		}
	}
	// The classification must not depend on the parent operation context: a
	// cancelled parent stops further requests but does not prove that an in-flight
	// write was not applied, so the caller still keeps the planned identity.
	if !isUnconfirmedListenerAclWrite(expiredCtx, transport) {
		t.Fatal("a transport failure must stay unconfirmed even after the parent stopped")
	}
	if !isDefinitiveListenerAclWriteOutcome(productionRejection) {
		t.Fatal("an explicitly rejected write must be definitive even when it is wrapped twice")
	}
	if isDefinitiveListenerAclWriteOutcome(transport) {
		t.Fatal("a transport failure must stay unknown")
	}
}

// The transport layer must keep the not-found and unreadable-response classes
// distinguishable after the shared retry wrapper and the URL sanitizer.
func TestUnitListenerAclAttachmentTransportClassification(t *testing.T) {
	client := &connectivity.AliyunClient{RegionId: "cn-hangzhou"}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	t.Run("not found stays a not found error", func(t *testing.T) {
		requestCapture(t, func(callCtx context.Context, action string, request map[string]interface{}) (map[string]interface{}, error) {
			return nil, listenerAclTestNotFoundError()
		})
		err := updateListenerAcl(ctx, client, "lb-123:http:80", "acl-1", "white", true)
		if err == nil || !errs.NotFoundError(err) {
			t.Fatalf("expected a not found error, got %v", err)
		}
		if isUnknownListenerAclWriteOutcome(err) {
			t.Fatal("a missing listener is a terminal result, not an unknown write")
		}
	})

	t.Run("unreadable response is an unknown write", func(t *testing.T) {
		requestCapture(t, func(callCtx context.Context, action string, request map[string]interface{}) (map[string]interface{}, error) {
			return nil, fmt.Errorf("%w: unexpected EOF", connectivity.ErrUnreadableResponse)
		})
		err := updateListenerAcl(ctx, client, "lb-123:http:80", "acl-1", "white", true)
		if err == nil {
			t.Fatal("expected an error")
		}
		if !isUnknownListenerAclWriteOutcome(err) {
			t.Fatalf("an unreadable response must be reconciled, got %v", err)
		}
	})

	t.Run("cancel during an attempt stops further attempts", func(t *testing.T) {
		calls := 0
		cancelCtx, cancelAttempt := context.WithCancel(context.Background())
		original := listenerAclAttachmentRpcPost
		listenerAclAttachmentRpcPost = func(callCtx context.Context, c *connectivity.AliyunClient, action string, request map[string]interface{}) (map[string]interface{}, error) {
			calls++
			cancelAttempt()
			return nil, errors.New(`Post "https://slb.aliyuncs.com/?Signature=abc": i/o timeout`)
		}
		t.Cleanup(func() { listenerAclAttachmentRpcPost = original })
		if err := updateListenerAcl(cancelCtx, client, "lb-123:http:80", "acl-1", "white", true); err == nil {
			t.Fatal("expected the cancelled operation to fail")
		}
		if calls != 1 {
			t.Fatalf("a cancelled operation must not issue another attempt, got %d calls", calls)
		}
	})
}
