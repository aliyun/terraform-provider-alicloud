package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudHbrHanaBackupClientsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrHanaBackupClientsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_hbr_hana_backup_client.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudHbrHanaBackupClientsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_hbr_hana_backup_client.default.id}_fake"]`,
		}),
	}
	clientIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrHanaBackupClientsDataSourceName(rand, map[string]string{
			"client_id": `"${alicloud_hbr_hana_backup_client.default.client_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrHanaBackupClientsDataSourceName(rand, map[string]string{
			"client_id": `"${alicloud_hbr_hana_backup_client.default.client_id}_fake"`,
		}),
	}
	clusterIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrHanaBackupClientsDataSourceName(rand, map[string]string{
			"cluster_id": `"${alicloud_hbr_hana_backup_client.default.cluster_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrHanaBackupClientsDataSourceName(rand, map[string]string{
			"cluster_id": `"${alicloud_hbr_hana_backup_client.default.cluster_id}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrHanaBackupClientsDataSourceName(rand, map[string]string{
			"status": `"ACTIVATED"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrHanaBackupClientsDataSourceName(rand, map[string]string{
			"status": `"REGISTERED"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrHanaBackupClientsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_hbr_hana_backup_client.default.id}"]`,
			"client_id":  `"${alicloud_hbr_hana_backup_client.default.client_id}"`,
			"cluster_id": `"${alicloud_hbr_hana_backup_client.default.cluster_id}"`,
			"status":     `"ACTIVATED"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrHanaBackupClientsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_hbr_hana_backup_client.default.id}_fake"]`,
			"client_id":  `"${alicloud_hbr_hana_backup_client.default.client_id}_fake"`,
			"cluster_id": `"${alicloud_hbr_hana_backup_client.default.cluster_id}_fake"`,
			"status":     `"REGISTERED"`,
		}),
	}
	var existAlicloudHbrHanaBackupClientsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"hana_backup_clients.#":                "1",
			"hana_backup_clients.0.id":             CHECKSET,
			"hana_backup_clients.0.vault_id":       CHECKSET,
			"hana_backup_clients.0.client_id":      CHECKSET,
			"hana_backup_clients.0.client_name":    CHECKSET,
			"hana_backup_clients.0.client_type":    "ECS_AGENT",
			"hana_backup_clients.0.client_version": CHECKSET,
			"hana_backup_clients.0.max_version":    CHECKSET,
			"hana_backup_clients.0.cluster_id":     CHECKSET,
			"hana_backup_clients.0.instance_id":    CHECKSET,
			"hana_backup_clients.0.instance_name":  CHECKSET,
			"hana_backup_clients.0.alert_setting":  "INHERITED",
			"hana_backup_clients.0.use_https":      "true",
			"hana_backup_clients.0.network_type":   "VPC",
			"hana_backup_clients.0.status_message": "",
			"hana_backup_clients.0.status":         "ACTIVATED",
		}
	}
	var fakeAlicloudHbrHanaBackupClientsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "0",
			"hana_backup_clients.#": "0",
		}
	}
	var alicloudHbrHanaBackupClientsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_hbr_hana_backup_clients.default",
		existMapFunc: existAlicloudHbrHanaBackupClientsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudHbrHanaBackupClientsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudHbrHanaBackupClientsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, clientIdConf, clusterIdConf, statusConf, allConf)
}

func testAccCheckAlicloudHbrHanaBackupClientsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccHbrHanaBackupClient-%d"
	}

	data "alicloud_hbr_vaults" "default" {
  		name_regex = "tf-test-hbr-hana-client"
	}

	resource "alicloud_hbr_hana_backup_client" "default" {
  		vault_id      = data.alicloud_hbr_vaults.default.vaults.0.id
  		client_info   = "[ { \"instanceId\": \"i-bp1dpl8hfbkh5rvvcmsg\", \"clusterId\": \"cl-000cnu7ti2rmj23dhp77\", \"sourceTypes\": [ \"HANA\" ]  }]"
  		alert_setting = "INHERITED"
  		use_https     = true
	}

	data "alicloud_hbr_hana_backup_clients" "default" {
  		vault_id = alicloud_hbr_hana_backup_client.default.vault_id
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
