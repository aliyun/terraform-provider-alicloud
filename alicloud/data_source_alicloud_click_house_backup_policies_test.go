package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlicloudClickHouseBackupPoliciesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.ClickHouseBackupPolicySupportRegions)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudClickHouseBackupPoliciesDataSourceName(rand, map[string]string{
			"db_cluster_id": `"${alicloud_click_house_backup_policy.default.db_cluster_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudClickHouseBackupPoliciesDataSourceName(rand, map[string]string{
			"db_cluster_id": `"${alicloud_click_house_backup_policy.default.db_cluster_id}_fake"`,
		}),
	}
	var existAlicloudClickHouseBackupPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"policies.#":                           "1",
			"policies.0.backup_retention_period":   "7",
			"policies.0.db_cluster_id":             CHECKSET,
			"policies.0.id":                        CHECKSET,
			"policies.0.preferred_backup_period.#": "2",
			"policies.0.preferred_backup_time":     "00:00Z-01:00Z",
			"policies.0.status":                    "true",
		}
	}
	var fakeAlicloudClickHouseBackupPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudClickHouseBackupPoliciesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_click_house_backup_policies.default",
		existMapFunc: existAlicloudClickHouseBackupPoliciesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudClickHouseBackupPoliciesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudClickHouseBackupPoliciesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
}

func testAccCheckAlicloudClickHouseBackupPoliciesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
  default = "tf-testAccBackupPolicy-%d"
}

data "alicloud_click_house_regions" "default" {
  current = true
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_click_house_regions.default.regions.0.zone_ids.0.zone_id
}

resource "alicloud_click_house_db_cluster" "default" {
  db_cluster_version      = "23.8"
  status                  = "Running"
  category                = "Basic"
  db_cluster_class        = "S8"
  db_cluster_network_type = "vpc"
  db_cluster_description  = var.name
  db_node_group_count     = "1"
  payment_type            = "PayAsYouGo"
  db_node_storage         = "500"
  storage_type            = "cloud_essd"
  vswitch_id              = alicloud_vswitch.default.id
  db_cluster_access_white_list {
    db_cluster_ip_array_name      = "test"
    security_ip_list              = "192.168.0.1"
  }
}

resource "alicloud_click_house_backup_policy" "default" {
	backup_retention_period = 7
	db_cluster_id = alicloud_click_house_db_cluster.default.id
	preferred_backup_period = ["Monday", "Tuesday"]
	preferred_backup_time = "00:00Z-01:00Z"
}

data "alicloud_click_house_backup_policies" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
