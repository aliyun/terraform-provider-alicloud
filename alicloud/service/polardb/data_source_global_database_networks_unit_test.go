package polardb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestUnitAlicloudPolarDBGlobalDatabaseNetworksDataSourceResourceGroupId covers
// the resource_group_id output on the real Read path. The cloud response
// currently omits the field, so the non-empty branch is driven by controlled
// responses here; the acceptance checks still require the key to exist.
// lintignore: R001
func TestUnitAlicloudPolarDBGlobalDatabaseNetworksDataSourceResourceGroupId(t *testing.T) {
	details := map[string]map[string]interface{}{
		"gdn-with-rg": {
			"GDNId":           "gdn-with-rg",
			"ResourceGroupId": "rg-acfmzaq3ypaqkdy",
			"Labels":          map[string]interface{}{"GDNVersion": "1.0"},
			"DBClusters": []interface{}{
				map[string]interface{}{"DBClusterId": "pc-with-rg", "Role": "Primary", "RegionId": "cn-hangzhou"},
			},
		},
		"gdn-without-rg": {
			"GDNId":  "gdn-without-rg",
			"Labels": map[string]interface{}{"GDNVersion": "2.0"},
			"DBClusters": []interface{}{
				map[string]interface{}{"DBClusterId": "pc-without-rg", "Role": "Primary", "RegionId": "cn-hangzhou"},
			},
		},
	}
	var (
		mu        sync.Mutex
		detailIDs []string
	)
	client := polarDBGlobalDatabaseNetworksTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch action := r.Form.Get("Action"); action {
		case "DescribeGlobalDatabaseNetworks":
			writePolarDBGlobalDatabaseNetworksTestResponse(t, w, map[string]interface{}{
				"Items": []interface{}{
					map[string]interface{}{
						"GDNId": "gdn-with-rg", "GDNDescription": "with resource group", "DBType": "MySQL",
						"DBVersion": "8.0", "CreateTime": "2026-01-01T00:00:00Z", "GDNStatus": "active",
					},
					map[string]interface{}{
						"GDNId": "gdn-without-rg", "GDNDescription": "without resource group", "DBType": "MySQL",
						"DBVersion": "8.0", "CreateTime": "2026-01-02T00:00:00Z", "GDNStatus": "active",
					},
				},
			})
		case "DescribeGlobalDatabaseNetwork":
			id := r.Form.Get("GDNId")
			mu.Lock()
			detailIDs = append(detailIDs, id)
			mu.Unlock()
			detail, ok := details[id]
			if !ok {
				t.Errorf("unexpected DescribeGlobalDatabaseNetwork GDNId %q", id)
				return
			}
			writePolarDBGlobalDatabaseNetworksTestResponse(t, w, detail)
		default:
			t.Errorf("unexpected action %s", action)
		}
	})

	d := schema.TestResourceDataRaw(t, DataSourceGlobalDatabaseNetworks().Schema, map[string]interface{}{})
	if err := dataSourceGlobalDatabaseNetworksRead(d, client); err != nil {
		t.Fatalf("read failed: %v", err)
	}

	networks, ok := d.Get("networks").([]interface{})
	if !ok || len(networks) != 2 {
		t.Fatalf("expected 2 networks, got %#v", d.Get("networks"))
	}
	byID := make(map[string]map[string]interface{}, len(networks))
	for _, item := range networks {
		network, ok := item.(map[string]interface{})
		if !ok {
			t.Fatalf("unexpected network item %#v", item)
		}
		byID[fmt.Sprint(network["id"])] = network
	}

	if got := byID["gdn-with-rg"]["resource_group_id"]; got != "rg-acfmzaq3ypaqkdy" {
		t.Fatalf("gdn-with-rg resource_group_id = %#v, want rg-acfmzaq3ypaqkdy", got)
	}
	if got, exists := byID["gdn-without-rg"]["resource_group_id"]; !exists || got != "" {
		t.Fatalf("gdn-without-rg resource_group_id = %#v (exists=%t), want an empty value with the key present", got, exists)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(detailIDs) != 2 || detailIDs[0] != "gdn-with-rg" || detailIDs[1] != "gdn-without-rg" {
		t.Fatalf("detail queries = %v, want the resource group resolved for each listed GDN", detailIDs)
	}
}

func polarDBGlobalDatabaseNetworksTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	endpoints := new(sync.Map)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("polardb", endpoint)
	return client
}

func writePolarDBGlobalDatabaseNetworksTestResponse(t *testing.T, w http.ResponseWriter, body map[string]interface{}) {
	t.Helper()
	body["RequestId"] = "test-request"
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Error(err)
	}
}
