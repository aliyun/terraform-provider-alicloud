package alicloud

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
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudPolarDBGlobalDatabaseNetworkDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_polardb_global_database_network.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_polardb_global_database_network.default.id}_fake"]`,
		}),
	}
	gdnIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"gdn_id": `"${alicloud_polardb_global_database_network.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"gdn_id": `"${alicloud_polardb_global_database_network.default.id}_fake"`,
		}),
	}
	dbClusterIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_global_database_network.default.db_cluster_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_global_database_network.default.db_cluster_id}_fake"`,
		}),
	}
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"description": `"${alicloud_polardb_global_database_network.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"description": `"${alicloud_polardb_global_database_network.default.description}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_polardb_global_database_network.default.id}"]`,
			"status": `"active"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_polardb_global_database_network.default.id}_fake"]`,
			"status": `"locked"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_polardb_global_database_network.default.id}"]`,
			"gdn_id":        `"${alicloud_polardb_global_database_network.default.id}"`,
			"db_cluster_id": `"${alicloud_polardb_global_database_network.default.db_cluster_id}"`,
			"description":   `"${alicloud_polardb_global_database_network.default.description}"`,
			"status":        `"active"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_polardb_global_database_network.default.id}_fake"]`,
			"gdn_id":        `"${alicloud_polardb_global_database_network.default.id}_fake"`,
			"db_cluster_id": `"${alicloud_polardb_global_database_network.default.db_cluster_id}_fake"`,
			"description":   `"${alicloud_polardb_global_database_network.default.description}_fake"`,
			"status":        `"locked"`,
		}),
	}
	// DescribeGlobalDatabaseNetworks does not return ResourceGroupId, so the
	// value may be empty until the API contract is confirmed. The assertion
	// requires the key and, once populated, a resource group ID shape.
	var existAlicloudPolarDBGlobalDatabaseNetworkDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"networks.#":                             "1",
			"networks.0.id":                          CHECKSET,
			"networks.0.gdn_id":                      CHECKSET,
			"networks.0.description":                 CHECKSET,
			"networks.0.gdn_version":                 CHECKSET,
			"networks.0.resource_group_id":           REGEXMATCH + "^(|rg-[0-9a-z]+)$",
			"networks.0.db_type":                     CHECKSET,
			"networks.0.db_version":                  CHECKSET,
			"networks.0.create_time":                 CHECKSET,
			"networks.0.status":                      "active",
			"networks.0.db_clusters.#":               "1",
			"networks.0.db_clusters.0.db_cluster_id": CHECKSET,
			"networks.0.db_clusters.0.role":          CHECKSET,
			"networks.0.db_clusters.0.region_id":     CHECKSET,
		}
	}
	var fakeAlicloudPolarDBGlobalDatabaseNetworkDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"networks.#": "0",
		}
	}
	var alicloudPolarDBGlobalDatabaseNetworkCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_polardb_global_database_networks.default",
		existMapFunc: existAlicloudPolarDBGlobalDatabaseNetworkDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudPolarDBGlobalDatabaseNetworkDataSourceNameMapFunc,
	}
	alicloudPolarDBGlobalDatabaseNetworkCheckInfo.dataSourceTestCheck(t, rand, idsConf, gdnIdConf, dbClusterIdConf, descriptionConf, statusConf, allConf)
}

func testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
	variable "name" {
		default = "tf-testAcc-%d"
	}
	data "alicloud_resource_manager_resource_groups" "default" {}

	data "alicloud_polardb_zones" "default" {}

	resource "alicloud_vpc" "default" {
		vpc_name = var.name
	}

	resource "alicloud_vswitch" "default" {
		zone_id    = data.alicloud_polardb_zones.default.ids[length(data.alicloud_polardb_zones.default.ids)-1]
		vpc_id     = alicloud_vpc.default.id
		cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
	}

	resource "alicloud_polardb_cluster" "default" {
		db_type       = "MySQL"
		db_version    = "8.0"
		pay_type      = "PostPaid"
		db_node_class = "polar.mysql.x4.medium"
		vswitch_id    = alicloud_vswitch.default.id
		description   = "${var.name}"
		lifecycle {
			ignore_changes = [encrypt_new_tables]
		}
	}
	resource "alicloud_polardb_global_database_network" "default" {
		db_cluster_id     = "${alicloud_polardb_cluster.default.id}"
		description       = var.name
		resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.1
		gdn_version       = "1.0"
	}
	data "alicloud_polardb_global_database_networks" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}

// TestUnitAlicloudPolarDBGlobalDatabaseNetworksDataSourceResourceGroupId covers
// the resource_group_id output on the real Read path. The cloud response
// currently omits the field, so the non-empty branch is driven by controlled
// responses here; the acceptance checks still require the key to exist.
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

	d := schema.TestResourceDataRaw(t, dataSourceAlicloudPolarDBGlobalDatabaseNetworks().Schema, map[string]interface{}{})
	if err := dataSourceAlicloudPolarDBGlobalDatabaseNetworksRead(d, client); err != nil {
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
