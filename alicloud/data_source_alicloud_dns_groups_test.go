package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudDnsGroupsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100000, 999999)

	testAccConfig := dataSourceTestAccConfigFunc("data.alicloud_dns_groups.default", fmt.Sprintf("tf-testacc-%d", rand), dataSourceDnsGroupsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_dns_group.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_dns_group.default.name}_fake",
		}),
	}
	existChangeMap := map[string]string{
		"groups.#":            "1",
		"groups.0.group_id":   "",
		"groups.0.group_name": "ALL",
	}
	nameAllConf := dataSourceTestAccConfig{
		existConfig:   testAccCheckAlicloudDnsGroupsDataSourceNameRegexAll,
		existChangMap: existChangeMap,
	}

	dnsGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, nameAllConf)
}

func dataSourceDnsGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_dns_group" "default" {
	name = "%s"
}
`, name)
}

const testAccCheckAlicloudDnsGroupsDataSourceNameRegexAll = `
data "alicloud_dns_groups" "default" {
  name_regex = "^ALL"
}`

var existDnsGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#":            "1",
		"groups.0.group_id":   CHECKSET,
		"groups.0.group_name": fmt.Sprintf("tf-testacc-%d", rand),
	}
}

var fakeDnsGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#": "0",
	}
}

var dnsGroupsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_dns_groups.default",
	existMapFunc: existDnsGroupsMapFunc,
	fakeMapFunc:  fakeDnsGroupsMapFunc,
}
