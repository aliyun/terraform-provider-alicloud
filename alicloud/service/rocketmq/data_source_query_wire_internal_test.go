package rocketmq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// rocketmqWireClient points the provider's ROA client at an isolated test
// server so the query string that actually leaves the data sources can be
// inspected instead of inferring it from the request map.
func rocketmqWireClient(t *testing.T, endpoint string) *connectivity.AliyunClient {
	t.Helper()
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("build credential: %v", err)
	}
	endpoints := new(sync.Map)
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("build client: %v", err)
	}
	endpoints.Store("rocketmq", endpoint)
	return client
}

func rocketmqWireServer(t *testing.T, captured *[]url.Values) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*captured = append(*captured, r.URL.Query())
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    map[string]interface{}{"list": []interface{}{}},
		}); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	return server
}

func rocketmqWireEndpoint(server *httptest.Server) string {
	return strings.TrimPrefix(server.URL, "http://")
}

// TestUnitRocketmqInstancesQueryWire proves the multi-value seriesCodes filter
// and the scalar filters reach the wire in the indexed form the gateway reads.
func TestUnitRocketmqInstancesQueryWire(t *testing.T) {
	var captured []url.Values
	server := rocketmqWireServer(t, &captured)

	ds := &instancesDataSource{}
	configure := &datasource.ConfigureResponse{}
	ds.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: rocketmqWireClient(t, rocketmqWireEndpoint(server)),
	}, configure)
	if configure.Diagnostics.HasError() {
		t.Fatalf("configure data source: %v", configure.Diagnostics)
	}

	seriesCodes := []string{"standard", "ultimate"}
	model := instancesModel{
		Filter:          types.StringValue("rmq-filter"),
		ResourceGroupId: types.StringValue("rg-test"),
		SeriesCodes: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("standard"),
			types.StringValue("ultimate"),
		}),
	}
	query := instancesQuery(&model, map[string]string{"env": "prod"}, seriesCodes)
	if _, err := ds.describeInstancesPage(context.Background(), query, 2); err != nil {
		t.Fatalf("describe instances page: %v", err)
	}

	if len(captured) != 1 {
		t.Fatalf("captured %d requests, want 1", len(captured))
	}
	want := map[string]string{
		"pageNumber":      "2",
		"pageSize":        "100",
		"filter":          "rmq-filter",
		"resourceGroupId": "rg-test",
		"tags":            `[{"key":"env","value":"prod"}]`,
		"seriesCodes.1":   "standard",
		"seriesCodes.2":   "ultimate",
	}
	for key, value := range want {
		if got := captured[0].Get(key); got != value {
			t.Errorf("wire query %s = %q, want %q (full query %q)", key, got, value, captured[0].Encode())
		}
	}
	if got := captured[0].Get("seriesCodes.3"); got != "" {
		t.Errorf("wire query carries an unexpected seriesCodes.3 = %q", got)
	}
}

// TestUnitRocketmqTopicsQueryWire proves the multi-value messageTypes filter
// and the scalar filter reach the wire in the indexed form the gateway reads.
func TestUnitRocketmqTopicsQueryWire(t *testing.T) {
	var captured []url.Values
	server := rocketmqWireServer(t, &captured)

	ds := &topicsDataSource{}
	configure := &datasource.ConfigureResponse{}
	ds.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: rocketmqWireClient(t, rocketmqWireEndpoint(server)),
	}, configure)
	if configure.Diagnostics.HasError() {
		t.Fatalf("configure data source: %v", configure.Diagnostics)
	}

	messageTypes := []string{"NORMAL", "LITE"}
	model := topicsModel{
		InstanceId: types.StringValue("rmq-test"),
		Filter:     types.StringValue("topic-filter"),
		MessageTypes: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("NORMAL"),
			types.StringValue("LITE"),
		}),
	}
	query := topicsQuery(&model, messageTypes)
	if _, err := ds.describeTopicsPage(context.Background(), "/instances/rmq-test/topics", query, 3); err != nil {
		t.Fatalf("describe topics page: %v", err)
	}

	if len(captured) != 1 {
		t.Fatalf("captured %d requests, want 1", len(captured))
	}
	want := map[string]string{
		"pageNumber":     "3",
		"pageSize":       "100",
		"filter":         "topic-filter",
		"messageTypes.1": "NORMAL",
		"messageTypes.2": "LITE",
	}
	for key, value := range want {
		if got := captured[0].Get(key); got != value {
			t.Errorf("wire query %s = %q, want %q (full query %q)", key, got, value, captured[0].Encode())
		}
	}
}
