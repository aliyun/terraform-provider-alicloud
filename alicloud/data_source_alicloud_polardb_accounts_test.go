package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudPolarDBClusterAccountsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.EcsClassicSupportedRegions)
	idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterAccountsDataSourceConfig(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_account.account.db_cluster_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterAccountsDataSourceConfig(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_account.account.db_cluster_id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterAccountsDataSourceConfig(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_account.account.db_cluster_id}"`,
			"name_regex":    `"${alicloud_polardb_account.account.account_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterAccountsDataSourceConfig(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_account.account.db_cluster_id}"`,
			"name_regex":    `"^test1234"`,
		}),
	}

	var existPolarClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                          CHECKSET,
			"accounts.0.account_description":   CHECKSET,
			"accounts.0.account_lock_state":    CHECKSET,
			"accounts.0.account_name":          CHECKSET,
			"accounts.0.account_status":        CHECKSET,
			"accounts.0.account_type":          CHECKSET,
			"accounts.0.database_privileges.#": CHECKSET,
		}
	}

	var fakePolarClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"accounts.#": CHECKSET,
		}
	}

	var PolarClusterCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_polardb_accounts.default",
		existMapFunc: existPolarClusterMapFunc,
		fakeMapFunc:  fakePolarClusterMapFunc,
	}
	PolarClusterCheckInfo.dataSourceTestCheck(t, rand, idConf, allConf)
}

func testAccCheckAlicloudPolarClusterAccountsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
	%s
	variable "name" {
	  default = "tf-testAccPolarClusterConfig_%d"
	}
	data "alicloud_polardb_node_classes" "this" {
	  db_type    = "MySQL"
	  db_version = "8.0"
      pay_type   = "PostPaid"
	  zone_id    = local.zone_id
	}
	resource "alicloud_polardb_cluster" "default" {
	  db_type           = "MySQL"
	  db_version        = "8.0"
      pay_type          = "PostPaid"
	  db_node_class     = data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class
	  vswitch_id        = local.vswitch_id
	  description       = var.name
	}

	resource "alicloud_polardb_account" "account" {
	  db_cluster_id        = alicloud_polardb_cluster.default.id
	  account_name         = "tftestnormal"
	  account_password     = "Test12345"
      account_description  = var.name
      account_type         = "Normal"
	}

	data "alicloud_polardb_accounts" "default" {
	  %s
	}
`, PolarDBCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
