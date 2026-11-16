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
		// Even before the context flips, the helper's own timeout text proves that
		// the attempt stopped being observed while the request may have been sent.
		{"retry helper timeout with a live budget", liveCtx, writeBudgetTimeout, true},
		// A cancellation returned before the request is sent proves that nothing
		// was transmitted, so it stays a refusal.
		{"cancellation before the request is sent", liveCtx, context.Canceled, false},
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

// A cancellation that ends the operation before any request was handed to the
// transport is a refusal: nothing was transmitted, so no state may be written for
// it. The attempt counter is what separates that case from a cancellation that
// interrupted a request already on the wire, which stays an unknown outcome.
func TestUnitListenerAclAttachmentPreSendCancellationIsRefusal(t *testing.T) {
	client := &connectivity.AliyunClient{RegionId: "cn-hangzhou"}

	t.Run("cancelled before the first request", func(t *testing.T) {
		calls := requestCapture(t, func(callCtx context.Context, action string, request map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{}, nil
		})
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := updateListenerAcl(ctx, client, "lb-123:http:80", "acl-1", "white", true)
		if err == nil {
			t.Fatal("expected the cancelled write to fail")
		}
		if *calls != 0 {
			t.Fatalf("a cancelled write must not reach the transport, got %d calls", *calls)
		}
		if !listenerAclErrorMatches(err, errListenerAclAttachmentNotSent) {
			t.Fatalf("the error must report that no request was sent: %v", err)
		}
		if got := strings.Count(err.Error(), errListenerAclAttachmentNotSent.Error()); got != 1 {
			t.Fatalf("the not-sent diagnostic must appear exactly once, got %d in %q", got, err.Error())
		}
		if !isDefinitiveListenerAclWriteOutcome(err) {
			t.Fatal("a write with no attempt is a definitive refusal")
		}
		if isUnconfirmedListenerAclWrite(ctx, err) {
			t.Fatal("a write with no attempt must not be treated as an unknown outcome")
		}
	})

	t.Run("cancelled after a request was on the wire", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		calls := requestCapture(t, func(callCtx context.Context, action string, request map[string]interface{}) (map[string]interface{}, error) {
			cancel()
			return nil, errors.New(`Post "https://slb.aliyuncs.com": i/o timeout`)
		})
		err := updateListenerAcl(ctx, client, "lb-123:http:80", "acl-1", "white", true)
		if err == nil {
			t.Fatal("expected the cancelled write to fail")
		}
		if *calls != 1 {
			t.Fatalf("expected exactly one attempt, got %d", *calls)
		}
		if listenerAclErrorMatches(err, errListenerAclAttachmentNotSent) {
			t.Fatal("a request that was on the wire must not be reported as not sent")
		}
		if !isUnconfirmedListenerAclWrite(ctx, err) {
			t.Fatal("a cancellation after an attempt must stay an unknown outcome")
		}
	})

	// The retry waiting loop can observe an already ended context without ever
	// running the callback. That path must get the same refusal classification.
	t.Run("ended context observed by the waiting loop", func(t *testing.T) {
		calls := requestCapture(t, func(callCtx context.Context, action string, request map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{}, nil
		})
		ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
		defer cancel()
		time.Sleep(time.Millisecond)
		err := updateListenerAcl(ctx, client, "lb-123:http:80", "acl-1", "white", true)
		if err == nil {
			t.Fatal("expected the ended context to fail the write")
		}
		if *calls != 0 {
			t.Fatalf("an ended context must not reach the transport, got %d calls", *calls)
		}
		if !listenerAclErrorMatches(err, errListenerAclAttachmentNotSent) {
			t.Fatalf("the error must report that no request was sent: %v", err)
		}
		if got := strings.Count(err.Error(), errListenerAclAttachmentNotSent.Error()); got != 1 {
			t.Fatalf("the not-sent diagnostic must appear exactly once, got %d in %q", got, err.Error())
		}
	})
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

// listenerAclAttachmentConfirms is the single proof used by the write wait, the
// pending marker and the destroy attribution. Every incomplete shape must fail
// it, including the shapes that still carry the requested group in the single
// AclId field: an empty or malformed member list cannot prove the binding.
func TestUnitListenerAclAttachmentConfirmationPredicate(t *testing.T) {
	cases := []struct {
		name   string
		object map[string]interface{}
		want   bool
	}{
		{
			name:   "planned binding",
			object: map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclId": "acl-1"},
			want:   true,
		},
		{
			name:   "planned binding behind a reported list",
			object: map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclIds": map[string]interface{}{"AclId": []interface{}{"acl-1", "acl-2"}}},
			want:   true,
		},
		{
			name:   "previous type",
			object: map[string]interface{}{"AclStatus": "on", "AclType": "black", "AclId": "acl-1"},
			want:   false,
		},
		{
			name:   "missing type",
			object: map[string]interface{}{"AclStatus": "on", "AclId": "acl-1"},
			want:   false,
		},
		{
			name:   "empty type",
			object: map[string]interface{}{"AclStatus": "on", "AclType": " ", "AclId": "acl-1"},
			want:   false,
		},
		{
			name:   "non string type",
			object: map[string]interface{}{"AclStatus": "on", "AclType": 5, "AclId": "acl-1"},
			want:   false,
		},
		{
			name:   "disabled listener",
			object: map[string]interface{}{"AclStatus": "off", "AclType": "white", "AclId": "acl-1"},
			want:   false,
		},
		{
			name:   "missing status",
			object: map[string]interface{}{"AclType": "white", "AclId": "acl-1"},
			want:   false,
		},
		{
			name:   "another group",
			object: map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclId": "acl-other"},
			want:   false,
		},
		{
			name:   "member list omitted while the single field names the group",
			object: map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclIds": map[string]interface{}{}},
			want:   false,
		},
		{
			name:   "member list reported empty",
			object: map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclIds": map[string]interface{}{"AclId": []interface{}{}}},
			want:   false,
		},
		{
			name:   "member list malformed beside the requested group",
			object: map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclId": "acl-1", "AclIds": map[string]interface{}{"AclId": []interface{}{5}}},
			want:   false,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if got := listenerAclAttachmentConfirms(tc.object, "acl-1", "white"); got != tc.want {
				t.Fatalf("listenerAclAttachmentConfirms(%v) = %v, want %v", tc.object, got, tc.want)
			}
		})
	}
	if listenerAclAttachmentConfirms(nil, "acl-1", "white") {
		t.Fatal("a nil response must never confirm a binding")
	}
}

// The TCP Describe contract documents LocationServiceTimeout as a transient
// error. It is retried on the read path only; the write path keeps failing fast
// because no write contract documents the code as transient.
func TestUnitListenerAclAttachmentReadTransientTimeout(t *testing.T) {
	timeoutErr := &tea.SDKError{
		StatusCode: tea.Int(400),
		Code:       tea.String("LocationServiceTimeout"),
		Message:    tea.String("Location service connection timeout. Please try again later."),
		Data:       tea.String("{}"),
	}
	if !isTransientListenerReadError(timeoutErr) {
		t.Fatal("the documented TCP describe timeout must be retried on the read path")
	}
	if isTransientListenerError(timeoutErr) {
		t.Fatal("the describe-only timeout code must not be extended to the write path")
	}
	terminal := &tea.SDKError{StatusCode: tea.Int(400), Code: tea.String("IpVersionConflict"), Message: tea.String("ip version conflict"), Data: tea.String("{}")}
	if isTransientListenerReadError(terminal) {
		t.Fatal("an IP version conflict stays terminal on the read path")
	}

	// Real transport shape: the first TCP describe answers with the documented
	// transient code, the second one succeeds.
	client := &connectivity.AliyunClient{RegionId: "cn-hangzhou"}
	calls := 0
	requestCapture(t, func(callCtx context.Context, action string, request map[string]interface{}) (map[string]interface{}, error) {
		calls++
		if calls == 1 {
			return nil, timeoutErr
		}
		return map[string]interface{}{"ListenerPort": 80, "AclStatus": "off"}, nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	object, err := describeListenerWithContext(ctx, client, "lb-123:tcp:80")
	if err != nil {
		t.Fatalf("the transient describe timeout must be retried: %s", err)
	}
	if calls != 2 {
		t.Fatalf("expected the describe to be retried once, got %d calls", calls)
	}
	if object["AclStatus"] != "off" {
		t.Fatalf("the retried describe did not return the listener: %v", object)
	}
}

// An undecodable 5xx keeps its HTTP status so the transient gate can retry it,
// while an undecodable 3xx or 4xx remains an explicit rejection. A write that
// ends with an undecodable 5xx is still an unknown outcome and must not be
// confirmed early.
func TestUnitListenerAclAttachmentUnreadableResponseClassification(t *testing.T) {
	unreadable5xx := func(status int) error {
		return &connectivity.UnreadableResponseError{
			StatusCode: status,
			Cause:      fmt.Errorf("%w: status %d: unexpected EOF", connectivity.ErrUnreadableResponse, status),
		}
	}
	unreadable4xx := fmt.Errorf("%w: status 400: unexpected EOF", connectivity.ErrRpcRejectedResponse)
	unreadable3xx := fmt.Errorf("%w: status 302: a signed rpc request is not followed across a redirect", connectivity.ErrRpcRejectedResponse)

	cases := []struct {
		name      string
		err       error
		transient bool
	}{
		{"undecodable 503", unreadable5xx(503), true},
		{"undecodable 500", unreadable5xx(500), true},
		{"undecodable 400", unreadable4xx, false},
		{"undecodable 302", unreadable3xx, false},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if got := isTransientListenerError(tc.err); got != tc.transient {
				t.Fatalf("isTransientListenerError = %v, want %v", got, tc.transient)
			}
			if got := isTransientListenerReadError(tc.err); got != tc.transient {
				t.Fatalf("isTransientListenerReadError = %v, want %v", got, tc.transient)
			}
		})
	}

	// The status survives the retry wrapper and the action wrapper.
	wrapped := errs.WrapErrorf(
		errs.WrapErrorf(unreadable5xx(502), "calling %s failed", "DescribeLoadBalancerHTTPListenerAttribute"),
		errs.DefaultErrorMsg, "lb-123:http:80", "DescribeLoadBalancerHTTPListenerAttribute", errs.AlibabaCloudSdkGoERROR,
	)
	if status, ok := listenerAclUnreadableStatus(wrapped); !ok || status != 502 {
		t.Fatalf("the unreadable status was lost in the wrapper chain: %d %v", status, ok)
	}
	if !isTransientListenerReadError(wrapped) {
		t.Fatal("a wrapped undecodable 5xx must stay retryable")
	}
	if !isUnconfirmedListenerAclWrite(context.Background(), unreadable5xx(503)) {
		t.Fatal("a write with an undecodable 5xx must stay an unknown outcome")
	}
	if isDefinitiveListenerAclWriteOutcome(unreadable5xx(503)) {
		t.Fatal("an undecodable 5xx must not be a definitive write refusal")
	}
	if !isDefinitiveListenerAclWriteOutcome(unreadable4xx) {
		t.Fatal("an undecodable 4xx must stay a definitive reject")
	}

	// The read path retries the transient 5xx and succeeds on the next answer.
	client := &connectivity.AliyunClient{RegionId: "cn-hangzhou"}
	calls := 0
	requestCapture(t, func(callCtx context.Context, action string, request map[string]interface{}) (map[string]interface{}, error) {
		calls++
		if calls == 1 {
			return nil, unreadable5xx(503)
		}
		return map[string]interface{}{"ListenerPort": 80, "AclStatus": "off"}, nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	object, err := describeListenerWithContext(ctx, client, "lb-123:http:80")
	if err != nil {
		t.Fatalf("the undecodable 5xx must be retried: %s", err)
	}
	if calls != 2 {
		t.Fatalf("expected the describe to be retried once, got %d calls", calls)
	}
	if object["AclStatus"] != "off" {
		t.Fatalf("the retried describe did not return the listener: %v", object)
	}

	// The 4xx control: exactly one attempt, no retry.
	calls = 0
	requestCapture(t, func(callCtx context.Context, action string, request map[string]interface{}) (map[string]interface{}, error) {
		calls++
		return nil, unreadable4xx
	})
	if _, err := describeListenerWithContext(ctx, client, "lb-123:http:80"); err == nil {
		t.Fatal("an undecodable 4xx must fail the describe")
	}
	if calls != 1 {
		t.Fatalf("an undecodable 4xx must not be retried, got %d calls", calls)
	}
}
