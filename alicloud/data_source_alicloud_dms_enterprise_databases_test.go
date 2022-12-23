package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDMSEnterpriseDatabaseDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDMSEnterpriseDatabaseSourceConfig(rand, map[string]string{
			"name_regex":  `"test2"`,
			"instance_id": `"${data.alicloud_dms_enterprise_instances.dms_enterprise_instances_ds.instances.0.instance_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudDMSEnterpriseDatabaseSourceConfig(rand, map[string]string{
			"name_regex":  `"test2_fake"`,
			"instance_id": `"${data.alicloud_dms_enterprise_instances.dms_enterprise_instances_ds.instances.0.instance_id}"`,
		}),
	}

	DMSEnterpriseDatabaseCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

var existDMSEnterpriseDatabaseMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"databases.#":    "1",
		"databases.0.id": CHECKSET,
	}
}

var fakeDMSEnterpriseDatabaseMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"databases.#": "0",
	}
}

var DMSEnterpriseDatabaseCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_dms_enterprise_databases.default",
	existMapFunc: existDMSEnterpriseDatabaseMapFunc,
	fakeMapFunc:  fakeDMSEnterpriseDatabaseMapFunc,
}

func testAccCheckAlicloudDMSEnterpriseDatabaseSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccDMSEnterpriseDatabase%d"
}

data "alicloud_dms_enterprise_instances" "dms_enterprise_instances_ds" {
  instance_type = "mysql"
  search_key    = "tf-test-no-deleting"
}

data "alicloud_dms_enterprise_databases" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
