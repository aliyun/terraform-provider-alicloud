package cloudmonitorservice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const cmsSubscriptionUnitServerId = "sub-unit-123"

type cmsSubscriptionUnitRecorder struct {
	mu       sync.Mutex
	methods  []string
	paths    []string
	bodies   []map[string]interface{}
	queries  []string
	response func(method, path string, body map[string]interface{}) (int, interface{})
}

func (r *cmsSubscriptionUnitRecorder) handle(w http.ResponseWriter, req *http.Request) {
	var body map[string]interface{}
	if req.Body != nil {
		_ = json.NewDecoder(req.Body).Decode(&body)
	}
	recorded := map[string]interface{}{}
	for key, value := range body {
		recorded[key] = value
	}
	r.mu.Lock()
	r.methods = append(r.methods, req.Method)
	r.paths = append(r.paths, req.URL.Path)
	r.bodies = append(r.bodies, recorded)
	r.queries = append(r.queries, req.URL.RawQuery)
	handler := r.response
	r.mu.Unlock()

	status, payload := handler(req.Method, req.URL.Path, body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

func (r *cmsSubscriptionUnitRecorder) requests(method string) ([]map[string]interface{}, []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	bodies := make([]map[string]interface{}, 0)
	paths := make([]string, 0)
	for i, recorded := range r.methods {
		if recorded == method {
			bodies = append(bodies, r.bodies[i])
			paths = append(paths, r.paths[i])
		}
	}
	return bodies, paths
}

// cmsSubscriptionUnitClient builds a real connectivity client whose CMS endpoint
// is a local HTTP server, so the assertions cover the signed request that the
// resource would send rather than a detached request map.
func cmsSubscriptionUnitClient(t *testing.T, response func(method, path string, body map[string]interface{}) (int, interface{})) (*connectivity.AliyunClient, *cmsSubscriptionUnitRecorder) {
	t.Helper()
	recorder := &cmsSubscriptionUnitRecorder{response: response}
	server := httptest.NewServer(http.HandlerFunc(recorder.handle))
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("cms", endpoint)
	return client, recorder
}

func cmsSubscriptionUnitResource(t *testing.T, apiClient *connectivity.AliyunClient) fwresource.Resource {
	t.Helper()
	r := NewCmsSubscriptionResource()
	configurable, ok := r.(fwresource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("%T does not implement resource.ResourceWithConfigure", r)
	}
	var resp fwresource.ConfigureResponse
	configurable.Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: apiClient}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("configuring the resource: %v", resp.Diagnostics)
	}
	return r
}

func cmsSubscriptionUnitSchema(t *testing.T, r fwresource.Resource) (fwschema.Schema, tftypes.Object) {
	t.Helper()
	ctx := context.Background()
	var schemaResp fwresource.SchemaResponse
	r.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building schema: %v", schemaResp.Diagnostics)
	}
	objType, ok := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatalf("schema type is %T, want tftypes.Object", schemaResp.Schema.Type().TerraformType(ctx))
	}
	return schemaResp.Schema, objType
}

type cmsSubscriptionUnitUnknown struct{}

func cmsSubscriptionUnitValue(t *testing.T, valueType tftypes.Type, raw interface{}) tftypes.Value {
	t.Helper()
	if _, ok := raw.(cmsSubscriptionUnitUnknown); ok {
		return tftypes.NewValue(valueType, tftypes.UnknownValue)
	}
	if raw == nil {
		return tftypes.NewValue(valueType, nil)
	}
	switch {
	case valueType.Is(tftypes.String):
		return tftypes.NewValue(valueType, raw.(string))
	case valueType.Is(tftypes.Bool):
		return tftypes.NewValue(valueType, raw.(bool))
	case valueType.Is(tftypes.Number):
		return tftypes.NewValue(valueType, raw.(float64))
	}
	switch typed := valueType.(type) {
	case tftypes.List:
		items := raw.([]interface{})
		values := make([]tftypes.Value, 0, len(items))
		for _, item := range items {
			values = append(values, cmsSubscriptionUnitValue(t, typed.ElementType, item))
		}
		return tftypes.NewValue(valueType, values)
	case tftypes.Object:
		attributes := raw.(map[string]interface{})
		values := make(map[string]tftypes.Value, len(typed.AttributeTypes))
		for name, attributeType := range typed.AttributeTypes {
			attributeRaw, present := attributes[name]
			if !present {
				values[name] = tftypes.NewValue(attributeType, nil)
				continue
			}
			values[name] = cmsSubscriptionUnitValue(t, attributeType, attributeRaw)
		}
		return tftypes.NewValue(valueType, values)
	}
	t.Fatalf("unsupported tftypes type %T", valueType)
	return tftypes.Value{}
}

func cmsSubscriptionUnitObject(t *testing.T, objType tftypes.Object, attributes map[string]interface{}) tftypes.Value {
	t.Helper()
	values := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, attributeType := range objType.AttributeTypes {
		raw, present := attributes[name]
		if !present {
			values[name] = tftypes.NewValue(attributeType, nil)
			continue
		}
		values[name] = cmsSubscriptionUnitValue(t, attributeType, raw)
	}
	return tftypes.NewValue(objType, values)
}

func cmsSubscriptionUnitObjectType(t *testing.T, objType tftypes.Object, name string) tftypes.Object {
	t.Helper()
	listType, ok := objType.AttributeTypes[name].(tftypes.List)
	if !ok {
		t.Fatalf("%s is %T, want a list", name, objType.AttributeTypes[name])
	}
	elementType, ok := listType.ElementType.(tftypes.Object)
	if !ok {
		t.Fatalf("%s element is %T, want an object", name, listType.ElementType)
	}
	return elementType
}

func cmsSubscriptionUnitStateAttribute(t *testing.T, state tfsdk.State, name string) types.String {
	t.Helper()
	var value types.String
	if diags := state.GetAttribute(context.Background(), path.Root(name), &value); diags.HasError() {
		t.Fatalf("reading state attribute %q: %v", name, diags)
	}
	return value
}

// TestUnitCmsSubscriptionCreateRequestShape proves that legal zero values are
// sent explicitly and that unset values stay absent on the wire.
func TestUnitCmsSubscriptionCreateRequestShape(t *testing.T) {
	created := map[string]interface{}{}
	apiClient, recorder := cmsSubscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		if method == http.MethodPost {
			for key, value := range body {
				created[key] = value
			}
			return http.StatusOK, map[string]interface{}{"requestId": "req-1", "subscriptionId": cmsSubscriptionUnitServerId}
		}
		subscription := map[string]interface{}{
			"subscriptionId":   cmsSubscriptionUnitServerId,
			"workspace":        "ws-1",
			"enable":           false,
			"userId":           "1",
			"subscriptionType": "NORMAL",
			"createTime":       "1",
			"updateTime":       "1",
		}
		for key, value := range created {
			subscription[key] = value
		}
		return http.StatusOK, map[string]interface{}{"requestId": "req-2", "subscription": subscription}
	})
	r := cmsSubscriptionUnitResource(t, apiClient)
	resourceSchema, objType := cmsSubscriptionUnitSchema(t, r)

	planValue := cmsSubscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                     cmsSubscriptionUnitUnknown{},
		"subscription_name":      "zero-values",
		"description":            "",
		"workspace":              "ws-1",
		"enable":                 false,
		"subscribe_legacy_event": false,
		"filter_setting": []interface{}{
			map[string]interface{}{
				"expression": "",
				"relation":   "",
				"conditions": []interface{}{
					map[string]interface{}{"field": "labels.alertname", "op": "EQ", "value": "CPU"},
				},
			},
		},
		"pushing_setting": []interface{}{
			map[string]interface{}{
				"alert_action_ids":   []interface{}{},
				"restore_action_ids": []interface{}{},
				"template_uuid":      "",
				"response_plan_id":   "",
			},
		},
	})
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	r.Create(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create returned diagnostics: %v", resp.Diagnostics)
	}

	creates, paths := recorder.requests(http.MethodPost)
	if len(creates) != 1 {
		t.Fatalf("expected exactly one CreateSubscription call, got %v", paths)
	}
	body := creates[0]
	if got, ok := body["description"]; !ok || got != "" {
		t.Fatalf("description must be sent explicitly as an empty string, got %v (present=%t)", got, ok)
	}
	if got, ok := body["subscribeLegacyEvent"]; !ok || got != false {
		t.Fatalf("subscribeLegacyEvent must be sent explicitly as false, got %v (present=%t)", got, ok)
	}
	filter, ok := body["filterSetting"].(map[string]interface{})
	if !ok {
		t.Fatalf("filterSetting missing from the request: %v", body)
	}
	if got, ok := filter["expression"]; !ok || got != "" {
		t.Fatalf("expression must be sent explicitly, got %v (present=%t)", got, ok)
	}
	conditions, ok := filter["conditions"].([]interface{})
	if !ok || len(conditions) != 1 {
		t.Fatalf("conditions must carry the configured member, got %v", filter["conditions"])
	}
	condition := conditions[0].(map[string]interface{})
	if condition["field"] != "labels.alertname" || condition["op"] != "EQ" || condition["value"] != "CPU" {
		t.Fatalf("unexpected condition payload: %v", condition)
	}
	pushing, ok := body["pushingSetting"].(map[string]interface{})
	if !ok {
		t.Fatalf("pushingSetting missing from the request: %v", body)
	}
	for _, key := range []string{"alertActionIds", "restoreActionIds"} {
		list, ok := pushing[key].([]interface{})
		if !ok || len(list) != 0 {
			t.Fatalf("%s must be sent as an explicit empty list, got %v", key, pushing[key])
		}
	}
	if got, ok := pushing["templateUuid"]; !ok || got != "" {
		t.Fatalf("templateUuid must be sent explicitly, got %v (present=%t)", got, ok)
	}

	if got := cmsSubscriptionUnitStateAttribute(t, resp.State, "id").ValueString(); got != cmsSubscriptionUnitServerId {
		t.Fatalf("id = %q, want %q", got, cmsSubscriptionUnitServerId)
	}
	if got := cmsSubscriptionUnitStateAttribute(t, resp.State, "region_id").ValueString(); got != "cn-hangzhou" {
		t.Fatalf("region_id = %q, want the provider region", got)
	}
}

// TestUnitCmsSubscriptionCreateKeepsIdentityOnWaitFailure proves that a failure
// after CreateSubscription returned an id still leaves a destroyable state.
func TestUnitCmsSubscriptionCreateKeepsIdentityOnWaitFailure(t *testing.T) {
	apiClient, recorder := cmsSubscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		if method == http.MethodPost {
			return http.StatusOK, map[string]interface{}{"requestId": "req-1", "subscriptionId": cmsSubscriptionUnitServerId}
		}
		return http.StatusForbidden, map[string]interface{}{"requestId": "req-2", "Code": "Forbidden", "Message": "denied"}
	})
	r := cmsSubscriptionUnitResource(t, apiClient)
	resourceSchema, objType := cmsSubscriptionUnitSchema(t, r)
	planValue := cmsSubscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                cmsSubscriptionUnitUnknown{},
		"subscription_name": "identity",
		"workspace":         "ws-1",
		"enable":            true,
	})
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	r.Create(context.Background(), req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Create to fail when the first read never succeeds")
	}
	if got := cmsSubscriptionUnitStateAttribute(t, resp.State, "id").ValueString(); got != cmsSubscriptionUnitServerId {
		t.Fatalf("state id = %q, want the created id to be persisted before the failing read", got)
	}
	if got := cmsSubscriptionUnitStateAttribute(t, resp.State, "workspace").ValueString(); got != "ws-1" {
		t.Fatalf("state workspace = %q, want ws-1", got)
	}
	if len(recorder.methods) == 0 {
		t.Fatal("no request was recorded")
	}
}

// TestUnitCmsSubscriptionUpdateRequiresRemovedValuesToClear proves the update
// waiter rejects a stale response that still carries a removed value.
func TestUnitCmsSubscriptionUpdateRequiresRemovedValuesToClear(t *testing.T) {
	plan := cmsSubscriptionResourceModel{
		SubscriptionName: types.StringValue("name"),
		Workspace:        types.StringValue("ws-1"),
	}
	state := cmsSubscriptionResourceModel{
		SubscriptionName: types.StringValue("name"),
		Workspace:        types.StringValue("ws-1"),
		Description:      types.StringValue("old description"),
	}
	removed := cmsSubscriptionRemovedPaths(plan, state)
	if len(removed) != 1 || strings.Join(removed[0], ".") != "description" {
		t.Fatalf("removed paths = %v, want [description]", removed)
	}
	current := map[string]interface{}{"subscriptionName": "name", "description": "old description"}
	if cmsSubscriptionPathsEmpty(current, removed) {
		t.Fatal("a stale response that still carries the removed description must not be accepted")
	}
	cleared := map[string]interface{}{"subscriptionName": "name"}
	if !cmsSubscriptionPathsEmpty(cleared, removed) {
		t.Fatal("a response without the removed description must be accepted")
	}
}

// TestUnitCmsSubscriptionRequestVisibleComparesNestedContent proves the waiter
// compares condition content, not only the number of conditions.
func TestUnitCmsSubscriptionRequestVisibleComparesNestedContent(t *testing.T) {
	request := map[string]interface{}{
		"subscriptionName": "name",
		"filterSetting": map[string]interface{}{
			"relation": "AND",
			"conditions": []map[string]interface{}{
				{"field": "labels.alertname", "op": "EQ", "value": "CPU"},
			},
		},
	}
	stale := map[string]interface{}{
		"subscriptionName": "name",
		"filterSetting": map[string]interface{}{
			"relation": "AND",
			"conditions": []interface{}{
				map[string]interface{}{"field": "labels.alertname", "op": "EQ", "value": "Memory"},
			},
		},
	}
	if cmsSubscriptionRequestVisible(stale, request) {
		t.Fatal("a response with the same condition count but different content must not be accepted")
	}
	applied := map[string]interface{}{
		"subscriptionName": "name",
		"filterSetting": map[string]interface{}{
			"relation": "AND",
			"conditions": []interface{}{
				map[string]interface{}{"field": "labels.alertname", "op": "EQ", "value": "CPU"},
			},
			"defaultsAddedByService": true,
		},
	}
	if !cmsSubscriptionRequestVisible(applied, request) {
		t.Fatal("a response that contains the requested configuration must be accepted")
	}
}

// TestUnitCmsSubscriptionStateKeepsExplicitZeroValues proves the read path does
// not turn configured empty values into nulls.
func TestUnitCmsSubscriptionStateKeepsExplicitZeroValues(t *testing.T) {
	plan := cmsSubscriptionResourceModel{
		SubscriptionName:     types.StringValue("name"),
		Workspace:            types.StringValue("ws-1"),
		Description:          types.StringValue(""),
		SubscribeLegacyEvent: types.BoolValue(false),
	}
	object := map[string]interface{}{
		"subscriptionId":   "sub-1",
		"subscriptionName": "name",
		"workspace":        "ws-1",
		"enable":           true,
	}
	state := cmsSubscriptionState(object, plan, "cn-hangzhou")
	if state.Description.IsNull() || state.Description.ValueString() != "" {
		t.Fatalf("description = %v, want an explicit empty string", state.Description)
	}
	if state.SubscribeLegacyEvent.IsNull() || state.SubscribeLegacyEvent.ValueBool() {
		t.Fatalf("subscribe_legacy_event = %v, want false", state.SubscribeLegacyEvent)
	}
	if state.RegionId.ValueString() != "cn-hangzhou" {
		t.Fatalf("region_id = %q, want cn-hangzhou", state.RegionId.ValueString())
	}
}

// TestUnitCmsSubscriptionChannelValidation proves CONTACT/GROUP/DUTY channels
// are rejected locally when enabled_sub_channels is missing.
func TestUnitCmsSubscriptionChannelValidation(t *testing.T) {
	r := NewCmsSubscriptionResource()
	resourceSchema, objType := cmsSubscriptionUnitSchema(t, r)
	channelType := cmsSubscriptionUnitObjectType(t, objType, "agent_config")
	routeType := cmsSubscriptionUnitObjectType(t, channelType, "routes")
	channelObjectType := cmsSubscriptionUnitObjectType(t, routeType, "channels")

	cases := map[string]struct {
		channel   map[string]interface{}
		wantError bool
	}{
		"contactWithoutSubChannels": {
			channel:   map[string]interface{}{"channel_type": "CONTACT", "receivers": []interface{}{"c1"}},
			wantError: true,
		},
		"groupWithEmptySubChannels": {
			channel:   map[string]interface{}{"channel_type": "GROUP", "enabled_sub_channels": []interface{}{}},
			wantError: true,
		},
		"dutyWithSubChannels": {
			channel: map[string]interface{}{"channel_type": "DUTY", "enabled_sub_channels": []interface{}{"SMS"}},
		},
		"webhookWithoutSubChannels": {
			channel: map[string]interface{}{"channel_type": "WEBHOOK", "receivers": []interface{}{"w1"}},
		},
	}
	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			channelValue := cmsSubscriptionUnitObject(t, channelObjectType, testCase.channel)
			routeValue := cmsSubscriptionUnitObject(t, routeType, map[string]interface{}{
				"channels": []interface{}{channelValue},
			})
			configValue := cmsSubscriptionUnitObject(t, objType, map[string]interface{}{
				"id":                cmsSubscriptionUnitUnknown{},
				"subscription_name": "name",
				"workspace":         "ws-1",
				"agent_config": []interface{}{
					map[string]interface{}{
						"agent_uuid": "agent-1",
						"routes":     []interface{}{routeValue},
					},
				},
			})
			req := fwresource.ValidateConfigRequest{Config: tfsdk.Config{Schema: resourceSchema, Raw: configValue}}
			resp := &fwresource.ValidateConfigResponse{}
			cmsSubscriptionChannelValidator{}.ValidateResource(context.Background(), req, resp)
			if testCase.wantError {
				if !resp.Diagnostics.HasError() {
					t.Fatal("expected a local validation error for a conditional channel type without enabled_sub_channels")
				}
				if !strings.Contains(fmt.Sprint(resp.Diagnostics), "enabled_sub_channels") {
					t.Fatalf("diagnostics %v do not point at enabled_sub_channels", resp.Diagnostics)
				}
				return
			}
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
		})
	}
}

// TestUnitCmsSubscriptionDescribeRejectsMalformedPayload proves a non-object
// $.subscription payload returns an error instead of panicking.
func TestUnitCmsSubscriptionDescribeRejectsMalformedPayload(t *testing.T) {
	cases := map[string]interface{}{
		"stringPayload": "not-an-object",
		"numberPayload": 42.0,
	}
	for name, subscription := range cases {
		t.Run(name, func(t *testing.T) {
			apiClient, _ := cmsSubscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
				return http.StatusOK, map[string]interface{}{"requestId": "req", "subscription": subscription}
			})
			_, err := NewCmsService(apiClient).DescribeCmsSubscription(context.Background(), "sub-1", "ws-1")
			if err == nil {
				t.Fatal("expected an error for a malformed subscription payload")
			}
		})
	}
}
