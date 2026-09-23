package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func SkipTestAccAlicloudDMSEnterpriseLogicDatabaseDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDMSEnterpriseLogicDatabaseSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_dms_enterprise_logic_database.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDMSEnterpriseLogicDatabaseSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_dms_enterprise_logic_database.default.id}_fake"]`,
		}),
	}

	DMSEnterpriseLogicDatabaseCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

var existDMSEnterpriseLogicDatabaseMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"databases.#":    "1",
		"databases.0.id": CHECKSET,
	}
}

var fakeDMSEnterpriseLogicDatabaseMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"databases.#": "0",
	}
}

var DMSEnterpriseLogicDatabaseCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_dms_enterprise_logic_databases.default",
	existMapFunc: existDMSEnterpriseLogicDatabaseMapFunc,
	fakeMapFunc:  fakeDMSEnterpriseLogicDatabaseMapFunc,
}

func testAccCheckAlicloudDMSEnterpriseLogicDatabaseSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccDMSEnterpriseLogicDatabase%d"
}

data "alicloud_dms_enterprise_instances" "dms_enterprise_instances_ds" {
  instance_type = "mysql"
  search_key    = "tf-test-no-deleting"
}

data "alicloud_dms_enterprise_databases" "test2" {
  name_regex  = "test2"
  instance_id = data.alicloud_dms_enterprise_instances.dms_enterprise_instances_ds.instances.0.instance_id
}

data "alicloud_dms_enterprise_databases" "test3" {
  name_regex  = "test3"
  instance_id = data.alicloud_dms_enterprise_instances.dms_enterprise_instances_ds.instances.0.instance_id
}

resource "alicloud_dms_enterprise_logic_database" "default" {
  alias     = var.name
  database_ids = ["${data.alicloud_dms_enterprise_databases.test2.databases.0.id}", "${data.alicloud_dms_enterprise_databases.test3.databases.0.id}"]
}

data "alicloud_dms_enterprise_logic_databases" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
