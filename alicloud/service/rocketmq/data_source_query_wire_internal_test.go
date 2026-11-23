package rocketmq

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
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
	client, err := connectivity.NewClient(config)
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
			"data": map[string]interface{}{
				"pageNumber": 1, "pageSize": 100, "totalCount": 0, "list": []interface{}{},
			},
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

func rocketmqReadWithIds(t *testing.T, ds datasource.DataSource, server *httptest.Server, ids *[]string) *datasource.ReadResponse {
	t.Helper()
	configurable, ok := ds.(datasource.DataSourceWithConfigure)
	if !ok {
		t.Fatalf("data source %T does not support Configure", ds)
	}
	configure := &datasource.ConfigureResponse{}
	configurable.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: rocketmqWireClient(t, rocketmqWireEndpoint(server)),
	}, configure)
	if configure.Diagnostics.HasError() {
		t.Fatalf("configure data source: %v", configure.Diagnostics)
	}

	schemaResponse := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("build schema: %v", schemaResponse.Diagnostics)
	}
	objectType, ok := schemaResponse.Schema.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatalf("schema type is %T, want tftypes.Object", schemaResponse.Schema.Type().TerraformType(context.Background()))
	}
	config := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, attributeType := range objectType.AttributeTypes {
		config[name] = tftypes.NewValue(attributeType, nil)
	}
	if _, ok := objectType.AttributeTypes["instance_id"]; ok {
		config["instance_id"] = tftypes.NewValue(tftypes.String, "rmq-test")
	}
	if ids != nil {
		values := make([]tftypes.Value, 0, len(*ids))
		for _, id := range *ids {
			values = append(values, tftypes.NewValue(tftypes.String, id))
		}
		config["ids"] = tftypes.NewValue(objectType.AttributeTypes["ids"], values)
	}

	response := &datasource.ReadResponse{State: tfsdk.State{
		Schema: schemaResponse.Schema,
		Raw:    tftypes.NewValue(objectType, nil),
	}}
	ds.Read(context.Background(), datasource.ReadRequest{Config: tfsdk.Config{
		Schema: schemaResponse.Schema,
		Raw:    tftypes.NewValue(objectType, config),
	}}, response)
	return response
}

func rocketmqDecodeReadState(t *testing.T, response *datasource.ReadResponse, target interface{}) {
	t.Helper()
	if diags := response.State.Get(context.Background(), target); diags.HasError() {
		t.Fatalf("decode state: %v", diags)
	}
}

func rocketmqDecodeStringList(t *testing.T, list types.List) []string {
	t.Helper()
	var values []string
	if diags := list.ElementsAs(context.Background(), &values, false); diags.HasError() {
		t.Fatalf("decode string list: %v", diags)
	}
	return values
}

func TestUnitRocketmqInstancesReadPreservesConfiguredIdsOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"pageNumber": 1,
				"pageSize":   100,
				"totalCount": 2,
				"list": []interface{}{
					map[string]interface{}{"instanceId": "rmq-b", "instanceName": "name-b", "topicCount": 0, "groupCount": 0},
					map[string]interface{}{"instanceId": "rmq-a", "instanceName": "name-a", "topicCount": 0, "groupCount": 0},
				},
			},
		}); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	response := rocketmqReadWithIds(t, &instancesDataSource{}, server, &[]string{"rmq-b", "rmq-a"})
	if response.Diagnostics.HasError() {
		t.Fatalf("read data source: %v", response.Diagnostics)
	}
	var state instancesModel
	rocketmqDecodeReadState(t, response, &state)
	gotIds := rocketmqDecodeStringList(t, state.Ids)
	if want := []string{"rmq-b", "rmq-a"}; !reflect.DeepEqual(gotIds, want) {
		t.Fatalf("configured ids = %v, want preserved order %v", gotIds, want)
	}
	gotNames := rocketmqDecodeStringList(t, state.Names)
	if want := []string{"name-a", "name-b"}; !reflect.DeepEqual(gotNames, want) {
		t.Fatalf("matched names = %v, want sorted names %v", gotNames, want)
	}
	if len(state.Instances) != 2 || state.Instances[0].InstanceId.ValueString() != "rmq-a" || state.Instances[1].InstanceId.ValueString() != "rmq-b" {
		t.Fatalf("matched instances are not sorted by instance ID: %+v", state.Instances)
	}
}

func TestUnitRocketmqReadPreservesExplicitEmptyIds(t *testing.T) {
	tests := []struct {
		name       string
		dataSource datasource.DataSource
		items      []interface{}
		matchedIds []string
		wantNames  []string
		decode     func(*testing.T, *datasource.ReadResponse) ([]string, []string, []string)
	}{
		{
			name:       "instances",
			dataSource: &instancesDataSource{},
			items: []interface{}{
				map[string]interface{}{"instanceId": "rmq-b", "instanceName": "name-b", "topicCount": 0, "groupCount": 0},
				map[string]interface{}{"instanceId": "rmq-a", "instanceName": "name-a", "topicCount": 0, "groupCount": 0},
			},
			matchedIds: []string{"rmq-a", "rmq-b"},
			wantNames:  []string{"name-a", "name-b"},
			decode: func(t *testing.T, response *datasource.ReadResponse) ([]string, []string, []string) {
				var state instancesModel
				rocketmqDecodeReadState(t, response, &state)
				objectIds := make([]string, 0, len(state.Instances))
				for _, item := range state.Instances {
					objectIds = append(objectIds, item.InstanceId.ValueString())
				}
				return rocketmqDecodeStringList(t, state.Ids), rocketmqDecodeStringList(t, state.Names), objectIds
			},
		},
		{
			name:       "topics",
			dataSource: &topicsDataSource{},
			items: []interface{}{
				map[string]interface{}{"topicName": "topic-b"},
				map[string]interface{}{"topicName": "topic-a", "maxSendTps": 1500},
			},
			matchedIds: []string{"topic-a", "topic-b"},
			wantNames:  []string{"topic-a", "topic-b"},
			decode: func(t *testing.T, response *datasource.ReadResponse) ([]string, []string, []string) {
				var state topicsModel
				rocketmqDecodeReadState(t, response, &state)
				objectIds := make([]string, 0, len(state.Topics))
				for _, item := range state.Topics {
					objectIds = append(objectIds, item.TopicName.ValueString())
				}
				return rocketmqDecodeStringList(t, state.Ids), rocketmqDecodeStringList(t, state.Names), objectIds
			},
		},
		{
			name:       "consumer_groups",
			dataSource: &consumerGroupsDataSource{},
			items: []interface{}{
				map[string]interface{}{"consumerGroupId": "group-b"},
				map[string]interface{}{"consumerGroupId": "group-a", "maxReceiveTps": 1000},
			},
			matchedIds: []string{"group-a", "group-b"},
			wantNames:  []string{"group-a", "group-b"},
			decode: func(t *testing.T, response *datasource.ReadResponse) ([]string, []string, []string) {
				var state consumerGroupsModel
				rocketmqDecodeReadState(t, response, &state)
				objectIds := make([]string, 0, len(state.Groups))
				for _, item := range state.Groups {
					objectIds = append(objectIds, item.ConsumerGroupId.ValueString())
				}
				return rocketmqDecodeStringList(t, state.Ids), rocketmqDecodeStringList(t, state.Names), objectIds
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]interface{}{
					"success": true,
					"data": map[string]interface{}{
						"pageNumber": 1, "pageSize": 100, "totalCount": len(test.items), "list": test.items,
					},
				}); err != nil {
					t.Errorf("encode response: %v", err)
				}
			}))
			t.Cleanup(server.Close)

			for _, testCase := range []struct {
				name    string
				ids     *[]string
				wantIds []string
			}{
				{name: "explicit_empty", ids: &[]string{}, wantIds: []string{}},
				{name: "omitted", ids: nil, wantIds: test.matchedIds},
			} {
				t.Run(testCase.name, func(t *testing.T) {
					response := rocketmqReadWithIds(t, test.dataSource, server, testCase.ids)
					if response.Diagnostics.HasError() {
						t.Fatalf("read data source: %v", response.Diagnostics)
					}
					gotIds, gotNames, gotObjectIds := test.decode(t, response)
					if len(gotIds) != len(testCase.wantIds) || (len(gotIds) > 0 && !reflect.DeepEqual(gotIds, testCase.wantIds)) {
						t.Fatalf("ids = %v, want %v", gotIds, testCase.wantIds)
					}
					if !reflect.DeepEqual(gotNames, test.wantNames) {
						t.Fatalf("matched names = %v, want %v", gotNames, test.wantNames)
					}
					if !reflect.DeepEqual(gotObjectIds, test.matchedIds) {
						t.Fatalf("matched object IDs = %v, want sorted IDs %v", gotObjectIds, test.matchedIds)
					}
				})
			}
		})
	}
}

func TestUnitRocketmqReadRejectsInvalidNumericFields(t *testing.T) {
	tests := []struct {
		name       string
		dataSource datasource.DataSource
		item       map[string]interface{}
		wantField  string
	}{
		{
			name:       "instances",
			dataSource: &instancesDataSource{},
			item:       map[string]interface{}{"instanceId": "rmq-a", "topicCount": 0},
			wantField:  "groupCount",
		},
		{
			name:       "topics",
			dataSource: &topicsDataSource{},
			item:       map[string]interface{}{"topicName": "topic-a", "maxSendTps": "invalid"},
			wantField:  "maxSendTps",
		},
		{
			name:       "consumer_groups",
			dataSource: &consumerGroupsDataSource{},
			item:       map[string]interface{}{"consumerGroupId": "group-a", "maxReceiveTps": "invalid"},
			wantField:  "maxReceiveTps",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]interface{}{
					"success": true,
					"data": map[string]interface{}{
						"pageNumber": 1, "pageSize": 100, "totalCount": 1, "list": []interface{}{test.item},
					},
				}); err != nil {
					t.Errorf("encode response: %v", err)
				}
			}))
			t.Cleanup(server.Close)

			response := rocketmqReadWithIds(t, test.dataSource, server, nil)
			if !response.Diagnostics.HasError() || !strings.Contains(fmt.Sprint(response.Diagnostics), test.wantField) {
				t.Fatalf("expected a diagnostic naming %s, got %v", test.wantField, response.Diagnostics)
			}
		})
	}
}

func TestUnitRocketmqTags(t *testing.T) {
	if got := rocketmqTagsRequestValue(map[string]string{"b": "2", "a": "1"}); got != `[{"key":"a","value":"1"},{"key":"b","value":"2"}]` {
		t.Fatalf("unexpected request value: %s", got)
	}

	item := map[string]interface{}{
		"tags": []interface{}{
			map[string]interface{}{"key": "a", "value": "1"},
			map[string]interface{}{"key": "b", "value": "2"},
		},
	}
	if !rocketmqTagsMatch(item, map[string]string{"a": "1"}) {
		t.Fatal("expected the configured tag to match")
	}
	if rocketmqTagsMatch(item, map[string]string{"a": "2"}) {
		t.Fatal("expected a different tag value to be rejected")
	}
	if rocketmqTagsMatch(item, map[string]string{"c": "1"}) {
		t.Fatal("expected a missing tag to be rejected")
	}
	if !rocketmqTagsMatch(map[string]interface{}{}, nil) {
		t.Fatal("expected an empty filter to match everything")
	}

	tags := rocketmqTagsFromAPI(item["tags"])
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if !reflect.DeepEqual([]string{tags[0].Key.ValueString(), tags[1].Key.ValueString()}, []string{"a", "b"}) {
		t.Fatalf("expected tags sorted by key, got %v", tags)
	}
}

func TestUnitRocketmqInstanceMapping(t *testing.T) {
	item := map[string]interface{}{
		"instanceId":      "rmq-1",
		"instanceName":    "name-1",
		"status":          "RUNNING",
		"regionId":        "cn-hangzhou",
		"seriesCode":      "ultimate",
		"subSeriesCode":   "cluster_ha",
		"paymentType":     "PayAsYouGo",
		"commodityCode":   "ons_rmqsub_public_cn",
		"serviceCode":     "rmq",
		"resourceGroupId": "rg-1",
		"userId":          "1234567890",
		"remark":          "example",
		"topicCount":      json.Number("3"),
		"groupCount":      json.Number("2"),
		"createTime":      "2026-01-01 00:00:00",
		"updateTime":      "2026-01-02 00:00:00",
		"startTime":       "2026-01-01 00:00:00",
		"expireTime":      "2026-02-01 00:00:00",
		"releaseTime":     "2026-02-02 00:00:00",
		"tags": []interface{}{
			map[string]interface{}{"key": "Created", "value": "TF"},
		},
		"productInfo": map[string]interface{}{
			"traceOn":      true,
			"capacityType": "provisioned",
		},
	}

	got, err := instanceFromAPI(item)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checks := []struct {
		name string
		got  string
		want string
	}{
		{"instance_id", got.InstanceId.ValueString(), "rmq-1"},
		{"user_id", got.UserId.ValueString(), "1234567890"},
		{"commodity_code", got.CommodityCode.ValueString(), "ons_rmqsub_public_cn"},
		{"remark", got.Remark.ValueString(), "example"},
		{"update_time", got.UpdateTime.ValueString(), "2026-01-02 00:00:00"},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Fatalf("%s: expected %q, got %q", c.name, c.want, c.got)
		}
	}
	if got.TopicCount.ValueInt64() != 3 || got.GroupCount.ValueInt64() != 2 {
		t.Fatalf("expected numeric counts 3/2, got %d/%d", got.TopicCount.ValueInt64(), got.GroupCount.ValueInt64())
	}
	if len(got.ProductInfo) != 1 {
		t.Fatalf("expected one product_info element, got %d", len(got.ProductInfo))
	}
	if !got.ProductInfo[0].TraceOn.ValueBool() || got.ProductInfo[0].CapacityType.ValueString() != "provisioned" {
		t.Fatalf("unexpected product_info mapping: %+v", got.ProductInfo[0])
	}
	if len(got.Tags) != 1 || got.Tags[0].Key.ValueString() != "Created" {
		t.Fatalf("unexpected tags mapping: %+v", got.Tags)
	}
}

func TestUnitRocketmqConsumerGroupMapping(t *testing.T) {
	item := map[string]interface{}{
		"consumerGroupId": "GID-1",
		"status":          "RUNNING",
		"remark":          "example",
		"instanceId":      "rmq-1",
		"regionId":        "cn-hangzhou",
		"topicName":       "lite-topic",
		"createTime":      "2026-01-01 00:00:00",
		"updateTime":      "2026-01-02 00:00:00",
		"maxReceiveTps":   json.Number("1000"),
		"messageModel":    "LITE_SELECTIVE",
	}

	got, err := consumerGroupFromAPI(item)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.MessageModel.ValueString() != "LITE_SELECTIVE" {
		t.Fatalf("expected message_model to be mapped from messageModel, got %q", got.MessageModel.ValueString())
	}
	if got.TopicName.ValueString() != "lite-topic" {
		t.Fatalf("expected topic_name to be mapped from topicName, got %q", got.TopicName.ValueString())
	}
	if got.MaxReceiveTps.ValueInt64() != 1000 {
		t.Fatalf("expected max_receive_tps 1000, got %d", got.MaxReceiveTps.ValueInt64())
	}

	normal, err := consumerGroupFromAPI(map[string]interface{}{"consumerGroupId": "GID-2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if normal.MessageModel.ValueString() != "" {
		t.Fatalf("expected an empty message_model for a normal group, got %q", normal.MessageModel.ValueString())
	}
	if !normal.MaxReceiveTps.IsNull() {
		t.Fatalf("expected a missing max_receive_tps to remain null, got %v", normal.MaxReceiveTps)
	}
}

func TestUnitRocketmqTopicMapping(t *testing.T) {
	got, err := topicFromAPI(map[string]interface{}{
		"topicName":           "topic-1",
		"maxSendTps":          json.Number("1500"),
		"liteTopicExpiration": json.Number("20"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.MaxSendTps.ValueInt64() != 1500 || got.LiteTopicExpiration.ValueInt64() != 20 {
		t.Fatalf("unexpected numeric mapping: %+v", got)
	}

	normal, err := topicFromAPI(map[string]interface{}{"topicName": "topic-2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !normal.MaxSendTps.IsNull() || !normal.LiteTopicExpiration.IsNull() {
		t.Fatalf("expected missing optional topic numerics to remain null, got %+v", normal)
	}

	t.Run("optional missing field is null without an error", func(t *testing.T) {
		got, err := rocketmqOptionalInt64(map[string]interface{}{}, "liteTopicExpiration")
		if err != nil || !got.IsNull() {
			t.Fatalf("expected null without an error, got value=%v err=%v", got, err)
		}
	})

	t.Run("optional invalid field is an error", func(t *testing.T) {
		_, err := rocketmqOptionalInt64(map[string]interface{}{"liteTopicExpiration": "20"}, "liteTopicExpiration")
		if err == nil || !strings.Contains(err.Error(), "liteTopicExpiration") {
			t.Fatalf("expected a field-specific error, got %v", err)
		}
	})
}
