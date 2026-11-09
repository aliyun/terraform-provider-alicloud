package privatezone

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestUnitInboundEndpointIpConfigRequest stays in the internal package because the request
// encoder is unexported; the ACC file uses package privatezone_test and cannot reach it.
func TestUnitInboundEndpointIpConfigRequest(t *testing.T) {
	ipConfigs := []inboundEndpointIpConfigModel{
		{
			AzId:      types.StringValue("cn-hangzhou-a"),
			CidrBlock: types.StringValue("172.16.0.0/20"),
			VSwitchId: types.StringValue("vsw-explicit"),
			Ip:        types.StringValue("172.16.0.10"),
		},
		{
			AzId:      types.StringValue("cn-hangzhou-b"),
			CidrBlock: types.StringValue("172.16.16.0/20"),
			VSwitchId: types.StringValue("vsw-null"),
			Ip:        types.StringNull(),
		},
		{
			AzId:      types.StringValue("cn-hangzhou-c"),
			CidrBlock: types.StringValue("172.16.32.0/20"),
			VSwitchId: types.StringValue("vsw-unknown"),
			Ip:        types.StringUnknown(),
		},
		{
			AzId:      types.StringValue("cn-hangzhou-d"),
			CidrBlock: types.StringValue("172.16.48.0/20"),
			VSwitchId: types.StringValue("vsw-empty"),
			Ip:        types.StringValue(""),
		},
	}

	request := inboundEndpointIpConfigRequest(ipConfigs)
	if got := request["IpConfig.1.Ip"]; got != "172.16.0.10" {
		t.Fatalf("explicit ip must be sent, got %#v", got)
	}
	if got := request["IpConfig.1.VSwitchId"]; got != "vsw-explicit" {
		t.Fatalf("vswitch id must be sent, got %#v", got)
	}
	if got := request["IpConfig.1.AzId"]; got != "cn-hangzhou-a" {
		t.Fatalf("az id must be sent as AzId, got %#v", got)
	}
	if got := request["IpConfig.1.CidrBlock"]; got != "172.16.0.0/20" {
		t.Fatalf("cidr block must be sent, got %#v", got)
	}
	for _, key := range []string{"IpConfig.2.Ip", "IpConfig.3.Ip", "IpConfig.4.Ip"} {
		if _, ok := request[key]; ok {
			t.Fatalf("%s must be omitted for null, unknown and empty ips", key)
		}
	}
	if got := request["IpConfig.2.VSwitchId"]; got != "vsw-null" {
		t.Fatalf("vswitch id must still be sent for an auto-assigned ip, got %#v", got)
	}
}

func inboundEndpointUnitPlan() (inboundEndpointResourceModel, []inboundEndpointIpConfigModel) {
	plan := inboundEndpointResourceModel{
		Id:              types.StringValue("ep-1"),
		EndpointName:    types.StringValue("tf-testAcc-example"),
		SecurityGroupId: types.StringValue("sg-1"),
		VpcId:           types.StringValue("vpc-1"),
		VpcRegionId:     types.StringValue("cn-hangzhou"),
		SpecVersion:     types.StringValue("standard"),
	}
	ipConfigs := []inboundEndpointIpConfigModel{
		{
			AzId:      types.StringValue("cn-hangzhou-a"),
			CidrBlock: types.StringValue("172.16.0.0/20"),
			VSwitchId: types.StringValue("vsw-1"),
			Ip:        types.StringValue("172.16.0.11"),
		},
		{
			AzId:      types.StringValue("cn-hangzhou-b"),
			CidrBlock: types.StringValue("172.16.16.0/20"),
			VSwitchId: types.StringValue("vsw-2"),
			Ip:        types.StringNull(),
		},
	}
	return plan, ipConfigs
}

// inboundEndpointTestClient builds a client whose pvtz endpoint is a local httptest server, so the
// tests read the request the provider actually sends instead of only the pre-serialization map.
func inboundEndpointTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parsing request form: %s", err)
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	client, err := (&connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}).Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("pvtz", endpoint)
	return client
}

// TestUnitInboundEndpointUpdateWireSendsExplicitIp proves on the wire (the HTTP request the
// provider serializes) that UpdateInboundEndpoint carries IpConfig.N.Ip for a configured ip and
// omits it for a service-assigned one. A request-map assertion alone cannot show this.
func TestUnitInboundEndpointUpdateWireSendsExplicitIp(t *testing.T) {
	ctx := context.Background()
	plan, newIpConfigs := inboundEndpointUnitPlan()
	planSet, diags := inboundEndpointIpConfigsToSet(ctx, newIpConfigs)
	if diags.HasError() {
		t.Fatalf("building plan ip_configs: %v", diags)
	}
	plan.IpConfigs = planSet

	state := plan
	state.EndpointName = types.StringValue("tf-testAcc-old")
	oldIpConfigs := []inboundEndpointIpConfigModel{
		{
			AzId:      types.StringValue("cn-hangzhou-a"),
			CidrBlock: types.StringValue("172.16.0.0/20"),
			VSwitchId: types.StringValue("vsw-1"),
			Ip:        types.StringValue("172.16.0.10"),
		},
		{
			AzId:      types.StringValue("cn-hangzhou-b"),
			CidrBlock: types.StringValue("172.16.16.0/20"),
			VSwitchId: types.StringValue("vsw-2"),
			Ip:        types.StringNull(),
		},
	}
	stateSet, diags := inboundEndpointIpConfigsToSet(ctx, oldIpConfigs)
	if diags.HasError() {
		t.Fatalf("building state ip_configs: %v", diags)
	}
	state.IpConfigs = stateSet

	request, update, diags := inboundEndpointUpdateRequest(ctx, plan, state)
	if diags.HasError() {
		t.Fatalf("building update request: %v", diags)
	}
	if !update {
		t.Fatal("expected a changed name and ip to require an update")
	}

	var mu sync.Mutex
	var form map[string][]string
	client := inboundEndpointTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		form = r.Form
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"RequestId":"req-update"}`)
	})

	service := NewService(client)
	if _, err := service.UpdateInboundEndpoint(ctx, request); err != nil {
		t.Fatalf("UpdateInboundEndpoint returned an error: %s", err)
	}

	mu.Lock()
	defer mu.Unlock()
	get := func(key string) string {
		if values, ok := form[key]; ok && len(values) > 0 {
			return values[0]
		}
		return ""
	}
	if got := get("Action"); got != "UpdateInboundEndpoint" {
		t.Fatalf("expected UpdateInboundEndpoint on the wire, got %q", got)
	}
	if got := get("EndpointId"); got != "ep-1" {
		t.Fatalf("expected EndpointId ep-1 on the wire, got %q", got)
	}
	if got := get("Name"); got != "tf-testAcc-example" {
		t.Fatalf("expected the changed name on the wire, got %q", got)
	}
	if got := get("IpConfig.1.Ip"); got != "172.16.0.11" {
		t.Fatalf("expected the explicit ip on the wire, got %q", got)
	}
	if got := get("IpConfig.1.AzId"); got != "cn-hangzhou-a" {
		t.Fatalf("expected the az id on the wire, got %q", got)
	}
	if _, ok := form["IpConfig.2.Ip"]; ok {
		t.Fatal("a service-assigned ip must be omitted from the update request")
	}
}

// TestUnitInboundEndpointCreateRetriesReuseClientToken records every CreateInboundEndpoint
// request: the first attempt is throttled and the retry must carry the same ClientToken. This is
// the direct evidence that one logical create cannot turn into two endpoints.
func TestUnitInboundEndpointCreateRetriesReuseClientToken(t *testing.T) {
	ctx := context.Background()
	plan, ipConfigs := inboundEndpointUnitPlan()
	request := inboundEndpointCreateRequest(plan, ipConfigs)

	var mu sync.Mutex
	var tokens []string
	client := inboundEndpointTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		tokens = append(tokens, r.Form.Get("ClientToken"))
		attempt := len(tokens)
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if attempt == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"Code":"Throttling.User","Message":"code: 500, request throttled","RequestId":"req-1"}`)
			return
		}
		fmt.Fprint(w, `{"EndpointId":"ep-1","RequestId":"req-2"}`)
	})

	service := NewService(client)
	response, err := service.CreateInboundEndpoint(ctx, request)
	if err != nil {
		t.Fatalf("CreateInboundEndpoint returned an error: %s", err)
	}
	if got := fmt.Sprint(response["EndpointId"]); got != "ep-1" {
		t.Fatalf("expected EndpointId ep-1, got %q", got)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(tokens) != 2 {
		t.Fatalf("expected the throttled create to be retried once, got %d attempts", len(tokens))
	}
	if tokens[0] == "" {
		t.Fatal("the create request must carry a ClientToken")
	}
	if tokens[0] != tokens[1] {
		t.Fatalf("every attempt of one create must reuse the same ClientToken, got %q and %q", tokens[0], tokens[1])
	}
}

// TestUnitInboundEndpointDeleteWaiterDistinguishesGoneFromAccepted covers the asynchronous delete
// terminal state: an accepted DeleteInboundEndpoint is not completion, and the waiter must only
// return once DescribeInboundEndpoint reports the endpoint is gone.
func TestUnitInboundEndpointDeleteWaiterDistinguishesGoneFromAccepted(t *testing.T) {
	ctx := context.Background()
	var mu sync.Mutex
	describes := 0
	client := inboundEndpointTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "DeleteInboundEndpoint":
			fmt.Fprint(w, `{"RequestId":"req-delete"}`)
		case "DescribeInboundEndpoint":
			mu.Lock()
			describes++
			attempt := describes
			mu.Unlock()
			if attempt == 1 {
				fmt.Fprint(w, `{"EndpointId":"ep-1","Status":"CHANGE_INIT","IpConfigs":[]}`)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"Code":"InboundEndpoint.NotExsts","Message":"code: 404, endpoint not exists","RequestId":"req-describe"}`)
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	service := NewService(client)
	if _, err := service.DeleteInboundEndpoint(ctx, map[string]interface{}{
		"EndpointId":  "ep-1",
		"Lang":        "en",
		"ClientToken": "token-1",
	}); err != nil {
		t.Fatalf("DeleteInboundEndpoint returned an error: %s", err)
	}
	if err := service.inboundEndpointWaitForDeleted(ctx, "ep-1", 30*time.Second); err != nil {
		t.Fatalf("waiting for the endpoint to disappear returned an error: %s", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if describes < 2 {
		t.Fatalf("the waiter must keep polling until DescribeInboundEndpoint reports not found, got %d describe calls", describes)
	}
}

// TestUnitInboundEndpointStateRejectsMalformedIpConfigs covers the response mapping: a malformed
// IpConfigs payload must produce a diagnostic instead of an empty or partially populated state.
func TestUnitInboundEndpointStateRejectsMalformedIpConfigs(t *testing.T) {
	ctx := context.Background()
	valid := map[string]interface{}{
		"Id":   "ep-1",
		"Name": "tf-testAcc-example",
		"IpConfigs": []interface{}{
			map[string]interface{}{
				"AzId":      "cn-hangzhou-a",
				"CidrBlock": "172.16.0.0/20",
				"Ip":        "172.16.0.10",
				"VSwitchId": "vsw-1",
			},
		},
	}
	model, diags := inboundEndpointState(ctx, valid)
	if diags.HasError() {
		t.Fatalf("a valid response must not produce diagnostics: %v", diags)
	}
	if got := model.IpConfigs.Elements(); len(got) != 1 {
		t.Fatalf("expected one ip_config, got %d", len(got))
	}

	cases := []struct {
		name   string
		mutate func(map[string]interface{})
	}{
		{name: "missing array", mutate: func(o map[string]interface{}) { delete(o, "IpConfigs") }},
		{name: "non array", mutate: func(o map[string]interface{}) { o["IpConfigs"] = "not-an-array" }},
		{name: "empty array", mutate: func(o map[string]interface{}) { o["IpConfigs"] = []interface{}{} }},
		{name: "non object element", mutate: func(o map[string]interface{}) {
			o["IpConfigs"] = []interface{}{"not-an-object"}
		}},
		{name: "missing az id", mutate: func(o map[string]interface{}) {
			o["IpConfigs"] = []interface{}{map[string]interface{}{
				"CidrBlock": "172.16.0.0/20",
				"VSwitchId": "vsw-1",
			}}
		}},
		{name: "non string cidr block", mutate: func(o map[string]interface{}) {
			o["IpConfigs"] = []interface{}{map[string]interface{}{
				"AzId":      "cn-hangzhou-a",
				"CidrBlock": 17216,
				"VSwitchId": "vsw-1",
			}}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			object := map[string]interface{}{
				"Id":   "ep-1",
				"Name": "tf-testAcc-example",
				"IpConfigs": []interface{}{
					map[string]interface{}{
						"AzId":      "cn-hangzhou-a",
						"CidrBlock": "172.16.0.0/20",
						"Ip":        "172.16.0.10",
						"VSwitchId": "vsw-1",
					},
				},
			}
			tc.mutate(object)
			if _, diags := inboundEndpointState(ctx, object); !diags.HasError() {
				t.Fatalf("expected a diagnostic for %s", tc.name)
			}
		})
	}
}
