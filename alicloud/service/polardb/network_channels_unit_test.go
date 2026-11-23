package polardb

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestUnitPolarDBNetworkChannelsSchema(t *testing.T) {
	p := &schema.Provider{DataSourcesMap: map[string]*schema.Resource{"alicloud_polardb_network_channels": DataSourceNetworkChannels()}}
	if err := p.InternalValidate(); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, ds := range ServicePackage().SDKDataSources {
		if ds.TypeName == "alicloud_polardb_network_channels" {
			found = true
		}
	}
	if !found {
		t.Fatal("data source is not registered")
	}
}

func TestUnitPolarDBNetworkChannelsEniID(t *testing.T) {
	ds := DataSourceNetworkChannels()
	field := ds.Schema["channels"].Elem.(*schema.Resource).Schema["eni_id"]
	if field == nil || field.Type != schema.TypeString || !field.Computed || field.Optional || field.Required {
		t.Fatal("channels.eni_id must be a computed-only string")
	}
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"db_cluster_id": "pc-source"})
	for _, tc := range []struct {
		name   string
		object map[string]interface{}
		want   string
	}{
		{"present", map[string]interface{}{"ChannelName": "ch-one", "EniId": "eni-unit"}, "eni-unit"},
		{"missing", map[string]interface{}{"ChannelName": "ch-one"}, ""},
		{"null", map[string]interface{}{"ChannelName": "ch-one", "EniId": nil}, ""},
		{"empty", map[string]interface{}{"ChannelName": "ch-one", "EniId": ""}, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			channels, _, err := flattenNetworkChannels([]interface{}{tc.object}, "pc-source", "cn-beijing", "")
			if err != nil {
				t.Fatal(err)
			}
			if err := d.Set("channels", channels); err != nil {
				t.Fatal(err)
			}
			if got := d.Get("channels.0.eni_id"); got != tc.want {
				t.Fatalf("eni_id in state = %#v, want %q", got, tc.want)
			}
		})
	}
}

func TestUnitPolarDBNetworkChannelsFlatten(t *testing.T) {
	item := map[string]interface{}{"ChannelName": " ch-one ", "DBClusterId": " pc-source ", "RegionId": " cn-beijing ", "TargetDBClusterId": "pc-target", "TargetIp": "172.16.0.10", "TargetPort": "5432", "Notes": "logical replication", "TargetType": "polardb_for_postgresql", "VpcId": "vpc-unit"}
	items := []interface{}{item}
	rawJSON, err := json.Marshal(items)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name  string
		value interface{}
	}{
		{"array", items}, {"json", string(rawJSON)}, {"wrapped", map[string]interface{}{"ChannelInfo": items}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, names, err := flattenNetworkChannels(tc.value, "pc-source", "cn-beijing", "ch-one")
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(names, []string{"ch-one"}) || len(got) != 1 {
				t.Fatalf("unexpected channels: %#v %#v", got, names)
			}
			for key, want := range map[string]interface{}{"channel_name": "ch-one", "db_cluster_id": "pc-source", "region_id": "cn-beijing", "target_db_cluster_id": "pc-target", "target_ip": "172.16.0.10", "target_port": "5432", "notes": "logical replication", "target_type": "polardb_for_postgresql", "vpc_id": "vpc-unit"} {
				if got[0][key] != want {
					t.Errorf("%s: got %#v, want %#v", key, got[0][key], want)
				}
			}
		})
	}
	got, names, err := flattenNetworkChannels([]interface{}{map[string]interface{}{"ChannelName": "ch-z"}, map[string]interface{}{"ChannelName": "ch-a"}}, "pc-source", "cn-beijing", "")
	if err != nil || len(got) != 2 || !reflect.DeepEqual(names, []string{"ch-a", "ch-z"}) {
		t.Fatalf("unstable ordering: %#v %#v %v", got, names, err)
	}
	if got[0]["db_cluster_id"] != "pc-source" || got[0]["region_id"] != "cn-beijing" {
		t.Fatalf("missing query identity fallback: %#v", got[0])
	}
	filtered, filteredNames, err := flattenNetworkChannels(items, "pc-source", "cn-beijing", "no-match")
	if err != nil || len(filtered) != 0 || len(filteredNames) != 0 {
		t.Fatalf("filter did not converge: %#v %#v %v", filtered, filteredNames, err)
	}
}

func TestUnitPolarDBNetworkChannelsRejectMalformedResponse(t *testing.T) {
	for _, value := range []interface{}{nil, 1, "not-json", "null", map[string]interface{}{"unexpected": []interface{}{}}, []interface{}{"invalid"}, []interface{}{map[string]interface{}{}}, []interface{}{map[string]interface{}{"ChannelName": 12}}} {
		if _, _, err := flattenNetworkChannels(value, "pc-source", "cn-beijing", ""); err == nil {
			t.Errorf("accepted malformed ChannelInfos: %#v", value)
		}
	}
	for _, value := range []interface{}{[]interface{}{}, "[]", map[string]interface{}{"ChannelInfo": []interface{}{}}} {
		items, names, err := flattenNetworkChannels(value, "pc-source", "cn-beijing", "")
		if err != nil || len(items) != 0 || len(names) != 0 {
			t.Errorf("empty result rejected: %#v %#v %v", items, names, err)
		}
	}
	for _, item := range []map[string]interface{}{
		{"ChannelName": "ch-one", "DBClusterId": "pc-other"},
		{"ChannelName": "ch-one", "RegionId": "cn-hangzhou"},
	} {
		if _, _, err := flattenNetworkChannels([]interface{}{item}, "pc-source", "cn-beijing", ""); err == nil {
			t.Fatalf("accepted a channel from a different query scope: %#v", item)
		}
	}
}

func networkChannelsUnitClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	endpoints := new(sync.Map)
	client, err := connectivity.NewClient(&connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential, RegionId: "cn-beijing",
		AccountType: "test", Protocol: "http", Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("polardb", endpoint)
	return client
}

func TestUnitPolarDBNetworkChannelsRead(t *testing.T) {
	var empty atomic.Bool
	var calls atomic.Int32
	client := networkChannelsUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		for key, want := range map[string]string{"Action": "DescribeNetworkChannel", "Version": "2017-08-01", "RegionId": "cn-beijing", "DBClusterId": "pc-source", "ChannelName": "ch-one", "ResourceGroupId": "rg-unit", "VpcId": "vpc-unit"} {
			if got := r.Form.Get(key); got != want {
				t.Errorf("%s: got %q, want %q", key, got, want)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		channels := []interface{}{}
		if !empty.Load() {
			channels = append(channels, map[string]interface{}{"ChannelName": "ch-one", "TargetPort": "5432", "EniId": "eni-unit"})
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"ChannelInfos": channels})
	})
	ds := DataSourceNetworkChannels()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"db_cluster_id": "pc-source", "channel_name": "ch-one", "resource_group_id": "rg-unit", "vpc_id": "vpc-unit"})
	if diags := ds.ReadContext(context.Background(), d, client); diags.HasError() {
		t.Fatal(diags)
	}
	if got := d.Get("names"); !reflect.DeepEqual(got, []interface{}{"ch-one"}) {
		t.Fatalf("unexpected names: %#v", got)
	}
	if got := d.Get("channels.0.eni_id"); got != "eni-unit" {
		t.Fatalf("API EniId was not exported to state: %#v", got)
	}
	id := d.Id()
	if id == "" {
		t.Fatal("read did not set an ID")
	}
	empty.Store(true)
	if diags := ds.ReadContext(context.Background(), d, client); diags.HasError() {
		t.Fatal(diags)
	}
	if d.Get("channels.#").(int) != 0 || d.Get("names.#").(int) != 0 {
		t.Fatal("refresh retained stale channel data")
	}
	if d.Id() != id {
		t.Fatal("query ID changed with returned contents")
	}
	if calls.Load() != 2 {
		t.Fatalf("unexpected request count: %d", calls.Load())
	}
}

func TestUnitPolarDBNetworkChannelsReadErrors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{"missing_cluster", 404, `{"Code":"InvalidDBClusterId.NotFound","Message":"cluster does not exist","RequestId":"unit-test"}`},
		{"denied", 403, `{"Code":"Forbidden.RAM","Message":"access denied","RequestId":"unit-test"}`},
		{"missing_channels", 200, `{"RequestId":"unit-test"}`},
		{"invalid_channels", 200, `{"ChannelInfos":{}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			client := networkChannelsUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			})
			ds := DataSourceNetworkChannels()
			d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"db_cluster_id": "pc-source"})
			diags := ds.ReadContext(context.Background(), d, client)
			if !diags.HasError() {
				t.Fatal("read silently accepted an error")
			}
			if tc.name == "missing_channels" && !strings.Contains(diags[0].Summary, "unit-test") {
				t.Fatal("malformed response error lost its RequestId")
			}
			if calls.Load() != 1 {
				t.Fatalf("non-retryable error was retried %d times", calls.Load())
			}
		})
	}
}
