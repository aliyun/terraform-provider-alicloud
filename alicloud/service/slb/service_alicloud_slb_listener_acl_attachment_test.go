package slb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// These unit tests cover the listener ACL attachment helpers and the lifecycle
// decision functions the public Create/Read/Update/Delete entry points call.
// The acceptance scenarios live in
// resource_alicloud_slb_listener_acl_attachment_test.go (package slb_test).

func TestUnitListenerAclAttachmentParseId(t *testing.T) {
	cases := []struct {
		name       string
		id         string
		listenerId string
		aclId      string
		wantErr    bool
	}{
		{name: "three segments", id: "lb-123:http:80", listenerId: "lb-123:http:80"},
		{name: "four segments", id: "lb-123:https:443:acl-456", listenerId: "lb-123:https:443", aclId: "acl-456"},
		{name: "two segments", id: "lb-123:80", wantErr: true},
		{name: "five segments", id: "lb-123:http:80:acl-1:extra", wantErr: true},
		{name: "empty protocol", id: "lb-123::80", wantErr: true},
		{name: "unsupported protocol", id: "lb-123:smtp:80", wantErr: true},
		{name: "port zero", id: "lb-123:http:0", wantErr: true},
		{name: "port above range", id: "lb-123:http:65536", wantErr: true},
		{name: "non numeric port", id: "lb-123:http:http", wantErr: true},
		{name: "empty acl suffix", id: "lb-123:http:80:", wantErr: true},
		{name: "combined acl suffix", id: "lb-123:http:80:acl-1,acl-2", wantErr: true},
		{name: "port at upper bound", id: "lb-123:udp:65535", listenerId: "lb-123:udp:65535"},
		{name: "port at lower bound", id: "lb-123:tcp:1", listenerId: "lb-123:tcp:1"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			listenerId, aclId, err := parseListenerAclAttachmentId(tc.id)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error for %q", tc.id)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %s", tc.id, err)
			}
			if listenerId != tc.listenerId || aclId != tc.aclId {
				t.Fatalf("got (%q, %q), want (%q, %q)", listenerId, aclId, tc.listenerId, tc.aclId)
			}
		})
	}
}

// Two listeners that share the same load balancer and frontend port but use
// different protocols must dispatch both the write and the read to their own
// action, because the same port can be reused by a different protocol.
func TestUnitListenerAclAttachmentProtocolDispatch(t *testing.T) {
	wantSet := map[string]string{
		"http":  "SetLoadBalancerHTTPListenerAttribute",
		"https": "SetLoadBalancerHTTPSListenerAttribute",
		"tcp":   "SetLoadBalancerTCPListenerAttribute",
		"udp":   "SetLoadBalancerUDPListenerAttribute",
	}
	wantDescribe := map[string]string{
		"http":  "DescribeLoadBalancerHTTPListenerAttribute",
		"https": "DescribeLoadBalancerHTTPSListenerAttribute",
		"tcp":   "DescribeLoadBalancerTCPListenerAttribute",
		"udp":   "DescribeLoadBalancerUDPListenerAttribute",
	}

	seen := map[string]string{}
	for protocol, setAction := range wantSet {
		listenerId, aclId, err := parseListenerAclAttachmentId("lb-shared:" + protocol + ":80:acl-1")
		if err != nil {
			t.Fatalf("parse %s listener id: %s", protocol, err)
		}
		if listenerId != "lb-shared:"+protocol+":80" || aclId != "acl-1" {
			t.Fatalf("unexpected identity for %s: (%q, %q)", protocol, listenerId, aclId)
		}
		action, err := listenerAttributeAction(protocol)
		if err != nil {
			t.Fatalf("set action for %s: %s", protocol, err)
		}
		if action != setAction {
			t.Fatalf("set action for %s = %q, want %q", protocol, action, setAction)
		}
		describeAction, err := listenerAttributeDescribeAction(protocol)
		if err != nil {
			t.Fatalf("describe action for %s: %s", protocol, err)
		}
		if describeAction != wantDescribe[protocol] {
			t.Fatalf("describe action for %s = %q, want %q", protocol, describeAction, wantDescribe[protocol])
		}
		for otherProtocol, otherAction := range seen {
			if otherAction == action {
				t.Fatalf("%s and %s share set action %q", protocol, otherProtocol, action)
			}
		}
		seen[protocol] = action
	}

	for _, protocol := range []string{"HTTP", "smtp", "", "http2"} {
		if _, err := listenerAttributeAction(protocol); err == nil {
			t.Fatalf("expected unsupported protocol %q to fail the set action", protocol)
		}
		if _, err := listenerAttributeDescribeAction(protocol); err == nil {
			t.Fatalf("expected unsupported protocol %q to fail the describe action", protocol)
		}
	}
}

func TestUnitListenerAclIdsFromDescribe(t *testing.T) {
	object := map[string]interface{}{
		"AclId": "acl-a,acl-b",
		"AclIds": map[string]interface{}{
			"AclId": []interface{}{"acl-b", "acl-c"},
		},
	}
	ids := listenerAclIds(object)
	if strings.Join(ids, ",") != "acl-a,acl-b,acl-c" {
		t.Fatalf("unexpected acl list: %v", ids)
	}
	if !listenerAclContains(object, "acl-b") {
		t.Fatal("expected acl-b to be reported as bound")
	}
	if listenerAclContains(object, "acl-z") {
		t.Fatal("expected acl-z to be reported as not bound")
	}

	single := map[string]interface{}{"AclId": "acl-only"}
	if got := listenerAclIds(single); len(got) != 1 || got[0] != "acl-only" {
		t.Fatalf("unexpected single acl list: %v", got)
	}
}

// The schema rejects a combined ACL value and an invalid listener identity
// locally, before any request is built.
func TestUnitListenerAclAttachmentSchemaValidation(t *testing.T) {
	res := NewListenerAclAttachmentResource()
	var schemaResp resource.SchemaResponse
	res.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building the schema failed: %v", schemaResp.Diagnostics)
	}

	validate := func(attribute, value string) bool {
		attr, ok := schemaResp.Schema.Attributes[attribute].(schema.StringAttribute)
		if !ok {
			t.Fatalf("%s is not a string attribute", attribute)
		}
		var resp validator.StringResponse
		for _, v := range attr.Validators {
			v.ValidateString(context.Background(), validator.StringRequest{
				Path:        path.Root(attribute),
				ConfigValue: types.StringValue(value),
			}, &resp)
		}
		return resp.Diagnostics.HasError()
	}

	cases := []struct {
		attribute string
		value     string
		wantError bool
	}{
		{attribute: "acl_id", value: "acl-123", wantError: false},
		{attribute: "acl_id", value: "acl-123,acl-456", wantError: true},
		{attribute: "acl_id", value: "acl 123", wantError: true},
		{attribute: "listener_id", value: "lb-123:http:80", wantError: false},
		{attribute: "listener_id", value: "lb-123:https:65535", wantError: false},
		{attribute: "listener_id", value: "lb-123:http:0", wantError: true},
		{attribute: "listener_id", value: "lb-123:http:65536", wantError: true},
		{attribute: "listener_id", value: "lb-123:smtp:80", wantError: true},
		{attribute: "listener_id", value: "lb-123:http:http", wantError: true},
	}
	for _, tc := range cases {
		if got := validate(tc.attribute, tc.value); got != tc.wantError {
			t.Fatalf("validate %s=%q: got error %v, want %v", tc.attribute, tc.value, got, tc.wantError)
		}
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

// fakeListenerAclApi drives the lifecycle decision functions without a live
// endpoint so each failure mode is asserted on the exact call sequence.
type fakeListenerAclApi struct {
	describe  func(ctx context.Context, listenerId string) (map[string]interface{}, error)
	set       func(ctx context.Context, listenerId, aclId, aclType string, enable bool) error
	describes int
	sets      int
	setAclIds []string
}

func (f *fakeListenerAclApi) Describe(ctx context.Context, listenerId string) (map[string]interface{}, error) {
	f.describes++
	if f.describe == nil {
		return map[string]interface{}{"AclStatus": "on", "AclType": "white", "ListenerPort": 80}, nil
	}
	return f.describe(ctx, listenerId)
}

func (f *fakeListenerAclApi) SetAcl(ctx context.Context, listenerId, aclId, aclType string, enable bool) error {
	f.sets++
	f.setAclIds = append(f.setAclIds, aclId)
	if f.set == nil {
		return nil
	}
	return f.set(ctx, listenerId, aclId, aclType, enable)
}

func listenerAclAttachmentTestPlan(listenerId, aclId, aclType string) listenerAclAttachmentResourceModel {
	return listenerAclAttachmentResourceModel{
		ListenerId: types.StringValue(listenerId),
		AclId:      types.StringValue(aclId),
		AclType:    types.StringValue(aclType),
	}
}

func TestUnitListenerAclAttachmentCreateLifecycle(t *testing.T) {
	listenerId := "lb-123:http:80"
	plan := listenerAclAttachmentTestPlan(listenerId, "acl-1", "white")

	t.Run("write then read failure keeps identity", func(t *testing.T) {
		api := &fakeListenerAclApi{
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return nil, errors.New("read exploded")
			},
		}
		model, warnings, err := createListenerAclAttachment(context.Background(), api, plan)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(warnings) != 1 {
			t.Fatalf("expected one warning about the failed first read, got %v", warnings)
		}
		if model.Id.ValueString() != listenerId || model.ListenerId.ValueString() != listenerId || model.AclId.ValueString() != "acl-1" {
			t.Fatalf("identity was not kept after the failed first read: %+v", model)
		}
	})

	t.Run("unknown write outcome applied keeps identity", func(t *testing.T) {
		api := &fakeListenerAclApi{
			set: func(context.Context, string, string, string, bool) error {
				return errors.New(`Post "https://slb.aliyuncs.com": dial tcp 10.0.0.1:443: i/o timeout`)
			},
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclId": "acl-1", "ListenerPort": 80}, nil
			},
		}
		model, _, err := createListenerAclAttachment(context.Background(), api, plan)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if model.AclId.ValueString() != "acl-1" || model.Id.ValueString() != listenerId {
			t.Fatalf("applied binding was not kept in state: %+v", model)
		}
		if api.sets != 1 {
			t.Fatalf("expected exactly one write, got %d", api.sets)
		}
	})

	t.Run("unknown write outcome not applied returns error without state", func(t *testing.T) {
		api := &fakeListenerAclApi{
			set: func(context.Context, string, string, string, bool) error {
				return errors.New(`Post "https://slb.aliyuncs.com": dial tcp 10.0.0.1:443: i/o timeout`)
			},
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return map[string]interface{}{"AclStatus": "on", "AclType": "white", "ListenerPort": 80}, nil
			},
		}
		model, _, err := createListenerAclAttachment(context.Background(), api, plan)
		if err == nil {
			t.Fatal("expected an error when the binding is not present")
		}
		if model.Id.ValueString() != "" || model.AclId.ValueString() != "" {
			t.Fatalf("state must stay empty for an unapplied write: %+v", model)
		}
		if api.sets != 1 {
			t.Fatalf("expected exactly one write attempt, got %d", api.sets)
		}
	})

	t.Run("terminal write error makes a single call and keeps no state", func(t *testing.T) {
		api := &fakeListenerAclApi{
			set: func(context.Context, string, string, string, bool) error {
				return &tea.SDKError{Code: tea.String("QuotaLimitExceeds.AclListenerOverLimit"), Message: tea.String("quota exceeded"), Data: tea.String("{}")}
			},
		}
		model, _, err := createListenerAclAttachment(context.Background(), api, plan)
		if err == nil {
			t.Fatal("expected the quota error to be returned")
		}
		if api.sets != 1 || api.describes != 0 {
			t.Fatalf("terminal error must fail fast: writes=%d describes=%d", api.sets, api.describes)
		}
		if model.AclId.ValueString() != "" {
			t.Fatalf("no state may be written for a rejected request: %+v", model)
		}
	})

	t.Run("successful write and read uses server values", func(t *testing.T) {
		api := &fakeListenerAclApi{
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return map[string]interface{}{"AclStatus": "on", "AclType": "black", "AclId": "acl-1", "ListenerPort": 80}, nil
			},
		}
		model, warnings, err := createListenerAclAttachment(context.Background(), api, plan)
		if err != nil || len(warnings) != 0 {
			t.Fatalf("unexpected result: warnings=%v err=%v", warnings, err)
		}
		if model.AclType.ValueString() != "black" || model.AclStatus.ValueString() != "on" {
			t.Fatalf("server values were not adopted: %+v", model)
		}
	})
}

func TestUnitListenerAclAttachmentUpdateLifecycle(t *testing.T) {
	listenerId := "lb-123:http:80"
	state := listenerAclAttachmentTestPlan(listenerId, "acl-1", "white")
	state.Id = types.StringValue(listenerId)
	state.AclStatus = types.StringValue("on")

	t.Run("read failure keeps the planned identity", func(t *testing.T) {
		plan := listenerAclAttachmentTestPlan(listenerId, "acl-2", "black")
		api := &fakeListenerAclApi{
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return nil, errors.New("read exploded")
			},
		}
		model, warnings, err := updateListenerAclAttachment(context.Background(), api, plan, state)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(warnings) != 1 {
			t.Fatalf("expected one warning, got %v", warnings)
		}
		if model.AclId.ValueString() != "acl-2" || model.AclType.ValueString() != "black" {
			t.Fatalf("planned identity was not kept: %+v", model)
		}
		if api.sets != 1 {
			t.Fatalf("expected one write, got %d", api.sets)
		}
	})

	t.Run("unchanged identity skips the write", func(t *testing.T) {
		plan := listenerAclAttachmentTestPlan(listenerId, "acl-1", "white")
		api := &fakeListenerAclApi{}
		_, warnings, err := updateListenerAclAttachment(context.Background(), api, plan, state)
		if err != nil || len(warnings) != 0 {
			t.Fatalf("unexpected result: warnings=%v err=%v", warnings, err)
		}
		if api.sets != 0 {
			t.Fatalf("expected no write when only unrelated values change, got %d", api.sets)
		}
	})

	t.Run("unknown update outcome applied keeps the plan", func(t *testing.T) {
		plan := listenerAclAttachmentTestPlan(listenerId, "acl-2", "black")
		api := &fakeListenerAclApi{
			set: func(context.Context, string, string, string, bool) error {
				return errors.New(`Post "https://slb.aliyuncs.com": dial tcp 10.0.0.1:443: i/o timeout`)
			},
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return map[string]interface{}{"AclStatus": "on", "AclType": "black", "AclId": "acl-2", "ListenerPort": 80}, nil
			},
		}
		model, _, err := updateListenerAclAttachment(context.Background(), api, plan, state)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if model.AclId.ValueString() != "acl-2" {
			t.Fatalf("planned binding was not kept: %+v", model)
		}
	})
}

func TestUnitListenerAclAttachmentReadLifecycle(t *testing.T) {
	listenerId := "lb-123:http:80"
	baseState := listenerAclAttachmentResourceModel{
		Id:         types.StringValue(listenerId),
		ListenerId: types.StringValue(listenerId),
		AclId:      types.StringValue("acl-1"),
		AclType:    types.StringValue("white"),
	}

	cases := []struct {
		name        string
		object      map[string]interface{}
		describeErr error
		wantFound   bool
		wantErr     bool
		wantAclId   string
	}{
		{
			name:      "bound group keeps the resource",
			object:    map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclId": "acl-1,acl-2"},
			wantFound: true,
			wantAclId: "acl-1",
		},
		{
			name:      "disabled access control removes the resource",
			object:    map[string]interface{}{"AclStatus": "off", "AclId": "acl-1"},
			wantFound: false,
		},
		{
			name:      "missing group removes the resource",
			object:    map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclId": "acl-2"},
			wantFound: false,
		},
		{
			name:        "missing listener removes the resource",
			describeErr: listenerAclTestNotFoundError(),
			wantFound:   false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			api := &fakeListenerAclApi{
				describe: func(context.Context, string) (map[string]interface{}, error) {
					return tc.object, tc.describeErr
				},
			}
			model, found, err := readListenerAclAttachment(context.Background(), api, baseState)
			if (err != nil) != tc.wantErr {
				t.Fatalf("unexpected error state: %v", err)
			}
			if found != tc.wantFound {
				t.Fatalf("found = %v, want %v", found, tc.wantFound)
			}
			if tc.wantFound && model.AclId.ValueString() != tc.wantAclId {
				t.Fatalf("acl_id = %q, want %q", model.AclId.ValueString(), tc.wantAclId)
			}
		})
	}

	t.Run("import without acl suffix picks the only group", func(t *testing.T) {
		state := listenerAclAttachmentResourceModel{Id: types.StringValue(listenerId)}
		api := &fakeListenerAclApi{
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclId": "acl-9"}, nil
			},
		}
		model, found, err := readListenerAclAttachment(context.Background(), api, state)
		if err != nil || !found || model.AclId.ValueString() != "acl-9" {
			t.Fatalf("unexpected result: found=%v model=%+v err=%v", found, model, err)
		}
	})

	t.Run("import without acl suffix rejects several groups", func(t *testing.T) {
		state := listenerAclAttachmentResourceModel{Id: types.StringValue(listenerId)}
		api := &fakeListenerAclApi{
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclId": "acl-a,acl-b"}, nil
			},
		}
		if _, _, err := readListenerAclAttachment(context.Background(), api, state); err == nil {
			t.Fatal("expected an ambiguity error")
		}
	})
}

// requestCapture replaces the transport entry for the duration of one test.
func requestCapture(t *testing.T, handler func(action string, request map[string]interface{}) (map[string]interface{}, error)) *int {
	t.Helper()
	original := listenerAclAttachmentRpcPost
	calls := 0
	listenerAclAttachmentRpcPost = func(client *connectivity.AliyunClient, action string, request map[string]interface{}) (map[string]interface{}, error) {
		calls++
		return handler(action, request)
	}
	t.Cleanup(func() { listenerAclAttachmentRpcPost = original })
	return &calls
}

func TestUnitListenerAclAttachmentWriteRetry(t *testing.T) {
	client := &connectivity.AliyunClient{RegionId: "cn-hangzhou"}

	t.Run("busy then success retries and converges", func(t *testing.T) {
		var requests []map[string]interface{}
		calls := requestCapture(t, func(action string, request map[string]interface{}) (map[string]interface{}, error) {
			requests = append(requests, request)
			if len(requests) == 1 {
				return nil, &tea.SDKError{Code: tea.String("OperationBusy"), Message: tea.String("code: 400, OperationBusy"), Data: tea.String("{}")}
			}
			return map[string]interface{}{}, nil
		})
		if err := updateListenerAcl(context.Background(), client, "lb-123:http:80", "acl-1", "white", true); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if *calls != 2 {
			t.Fatalf("expected two attempts (busy then success), got %d", *calls)
		}
		if len(requests) != 2 {
			t.Fatalf("expected two captured requests, got %d", len(requests))
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
		calls := requestCapture(t, func(action string, request map[string]interface{}) (map[string]interface{}, error) {
			return nil, &tea.SDKError{Code: tea.String("IpVersionConflict"), Message: tea.String("ip version conflict"), Data: tea.String("{}")}
		})
		if err := updateListenerAcl(context.Background(), client, "lb-123:http:80", "acl-1", "white", true); err == nil {
			t.Fatal("expected the terminal error to be returned")
		}
		if *calls != 1 {
			t.Fatalf("terminal error must not be retried, got %d calls", *calls)
		}
	})

	t.Run("cancelled context sends no request", func(t *testing.T) {
		calls := requestCapture(t, func(action string, request map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{}, nil
		})
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if err := updateListenerAcl(ctx, client, "lb-123:http:80", "acl-1", "white", true); err == nil {
			t.Fatal("expected the cancelled context to fail the write")
		}
		if *calls != 0 {
			t.Fatalf("cancelled context must not send a request, got %d calls", *calls)
		}
	})

	t.Run("describe uses the protocol specific read action", func(t *testing.T) {
		var actions []string
		var request map[string]interface{}
		requestCapture(t, func(action string, req map[string]interface{}) (map[string]interface{}, error) {
			actions = append(actions, action)
			request = req
			return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": "acl-1"}, nil
		})
		if _, err := describeListenerWithContext(context.Background(), client, "lb-123:udp:7000"); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(actions) != 1 || actions[0] != "DescribeLoadBalancerUDPListenerAttribute" {
			t.Fatalf("unexpected describe action: %v", actions)
		}
		if request["LoadBalancerId"] != "lb-123" || fmt.Sprint(request["ListenerPort"]) != "7000" {
			t.Fatalf("describe request identity is wrong: %v", request)
		}
	})
}

func TestUnitListenerAclAttachmentDeleteLifecycle(t *testing.T) {
	listenerId := "lb-123:http:80"

	t.Run("missing listener is already gone", func(t *testing.T) {
		api := &fakeListenerAclApi{
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return nil, listenerAclTestNotFoundError()
			},
		}
		if err := deleteListenerAclAttachment(context.Background(), api, listenerId); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if api.sets != 0 {
			t.Fatalf("expected no write, got %d", api.sets)
		}
	})

	t.Run("disabled access control makes no write", func(t *testing.T) {
		api := &fakeListenerAclApi{
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return map[string]interface{}{"AclStatus": "off"}, nil
			},
		}
		if err := deleteListenerAclAttachment(context.Background(), api, listenerId); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if api.sets != 0 {
			t.Fatalf("expected no write, got %d", api.sets)
		}
	})

	t.Run("disable then wait until off", func(t *testing.T) {
		describeCalls := 0
		api := &fakeListenerAclApi{
			describe: func(context.Context, string) (map[string]interface{}, error) {
				describeCalls++
				if describeCalls == 1 {
					return map[string]interface{}{"AclStatus": "on"}, nil
				}
				return map[string]interface{}{"AclStatus": "off"}, nil
			},
		}
		if err := deleteListenerAclAttachment(context.Background(), api, listenerId); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if api.sets != 1 || api.setAclIds[0] != "" {
			t.Fatalf("expected one disable write, got %d (%v)", api.sets, api.setAclIds)
		}
	})

	t.Run("cancelled deadline stops the wait", func(t *testing.T) {
		api := &fakeListenerAclApi{
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return map[string]interface{}{"AclStatus": "on"}, nil
			},
		}
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		start := time.Now()
		if err := deleteListenerAclAttachment(ctx, api, listenerId); err == nil {
			t.Fatal("expected the deadline to fail the delete")
		}
		if elapsed := time.Since(start); elapsed > 5*time.Second {
			t.Fatalf("delete did not stop promptly: %s", elapsed)
		}
	})
}
