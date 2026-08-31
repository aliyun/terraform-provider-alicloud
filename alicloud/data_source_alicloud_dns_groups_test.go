package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAlidnsGroupsDataSource(t *testing.T) {
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
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_dns_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_dns_group.default.id}_fake"},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_dns_group.default.name}",
			"ids":        []string{"${alicloud_dns_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_dns_group.default.name}_fake",
			"ids":        []string{"${alicloud_dns_group.default.id}"},
		}),
	}
	existChangeMap := map[string]string{
		"ids.#":               "1",
		"ids.0":               REMOVEKEY,
		"names.#":             "1",
		"names.0":             "ALL",
		"groups.#":            "1",
		"groups.0.group_id":   "",
		"groups.0.group_name": "ALL",
	}
	nameAllConf := dataSourceTestAccConfig{
		existConfig:   testAccCheckAlicloudDnsGroupsDataSourceNameRegexAll,
		existChangMap: existChangeMap,
	}

	dnsGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf, nameAllConf)
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
		"ids.#":               "1",
		"ids.0":               CHECKSET,
		"names.#":             "1",
		"names.0":             fmt.Sprintf("tf-testacc-%d", rand),
		"groups.0.group_id":   CHECKSET,
		"groups.0.group_name": fmt.Sprintf("tf-testacc-%d", rand),
	}
}

var fakeDnsGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":    "0",
		"names.#":  "0",
		"groups.#": "0",
	}
}

var dnsGroupsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_dns_groups.default",
	existMapFunc: existDnsGroupsMapFunc,
	fakeMapFunc:  fakeDnsGroupsMapFunc,
}
