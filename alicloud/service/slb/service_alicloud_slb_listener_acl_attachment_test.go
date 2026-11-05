package slb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	if !errs.NotFoundError(notFound) {
		t.Fatal("a 404 response must stay recognizable as not found")
	}
}

// listenerAclTestNotFoundError is the error shape errs.NotFoundError detects
// for a listener that no longer exists.
func listenerAclTestNotFoundError() error {
	return &tea.SDKError{
		StatusCode: tea.Int(404),
		Code:       tea.String("InvalidLoadBalancerId.NotFound"),
		Message:    tea.String("The specified resource does not exist."),
		Data:       tea.String("{}"),
	}
}

// The not-found conversion has to survive the shared retry wrapper: read must
// clear the state and delete must succeed when the listener is already gone.
func TestUnitListenerAclAttachmentTransportNotFound(t *testing.T) {
	client := &connectivity.AliyunClient{RegionId: "cn-hangzhou"}
	requestCapture(t, func(callCtx context.Context, action string, request map[string]interface{}) (map[string]interface{}, error) {
		return nil, listenerAclTestNotFoundError()
	})
	api := newListenerAclAttachmentApi(client)
	state := listenerAclAttachmentResourceModel{
		Id:         types.StringValue("lb-123:http:80"),
		ListenerId: types.StringValue("lb-123:http:80"),
		AclId:      types.StringValue("acl-1"),
		AclType:    types.StringValue("white"),
	}
	if _, found, err := readListenerAclAttachment(context.Background(), api, state); err != nil || found {
		t.Fatalf("read must report a missing listener without an error: found=%v err=%v", found, err)
	}
	if err := deleteListenerAclAttachment(context.Background(), api, "lb-123:http:80"); err != nil {
		t.Fatalf("delete must tolerate a missing listener: %s", err)
	}
}
