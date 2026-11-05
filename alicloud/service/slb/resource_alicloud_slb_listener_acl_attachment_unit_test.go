package slb

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
		{name: "empty acl suffix", id: "lb-123:http:80:", wantErr: true},
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

// The protocol part of the listener ID selects the protocol specific
// SetLoadBalancer*ListenerAttribute action, so two listeners that share the
// same load balancer and frontend port but use different protocols must never
// dispatch to the same action.
func TestUnitListenerAclAttachmentProtocolDispatch(t *testing.T) {
	cases := map[string]string{
		"http":  "SetLoadBalancerHTTPListenerAttribute",
		"https": "SetLoadBalancerHTTPSListenerAttribute",
		"tcp":   "SetLoadBalancerTCPListenerAttribute",
		"udp":   "SetLoadBalancerUDPListenerAttribute",
	}
	actions := map[string]string{}
	for protocol, want := range cases {
		listenerId, aclId, err := parseListenerAclAttachmentId("lb-shared:" + protocol + ":80:acl-1")
		if err != nil {
			t.Fatalf("parse %s listener id: %s", protocol, err)
		}
		if listenerId != "lb-shared:"+protocol+":80" || aclId != "acl-1" {
			t.Fatalf("unexpected identity for %s: (%q, %q)", protocol, listenerId, aclId)
		}
		action, err := listenerAttributeAction(protocol)
		if err != nil {
			t.Fatalf("action for %s: %s", protocol, err)
		}
		if action != want {
			t.Fatalf("action for %s = %q, want %q", protocol, action, want)
		}
		for otherProtocol, otherAction := range actions {
			if otherAction == action {
				t.Fatalf("%s and %s share action %q", protocol, otherProtocol, action)
			}
		}
		actions[protocol] = action
	}

	for _, protocol := range []string{"HTTP", "smtp", "", "http2"} {
		if _, err := listenerAttributeAction(protocol); err == nil {
			t.Fatalf("expected unsupported protocol %q to fail", protocol)
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

func TestUnitListenerAclAttachmentPlanState(t *testing.T) {
	plan := listenerAclAttachmentResourceModel{
		Id:         types.StringValue("lb-123:http:80"),
		ListenerId: types.StringValue("lb-123:http:80"),
		AclId:      types.StringValue("acl-456"),
		AclType:    types.StringValue("white"),
	}
	model := listenerAclAttachmentPlanState(plan)
	if model.Id.ValueString() != "lb-123:http:80" {
		t.Fatalf("unexpected id: %q", model.Id.ValueString())
	}
	if model.ListenerId.ValueString() != "lb-123:http:80" || model.AclId.ValueString() != "acl-456" {
		t.Fatalf("identity was not preserved: %+v", model)
	}
	if model.AclType.ValueString() != "white" || model.AclStatus.ValueString() != "on" {
		t.Fatalf("unexpected acl values: type=%q status=%q", model.AclType.ValueString(), model.AclStatus.ValueString())
	}
}
