package alidns_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudAlidnsGroupsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100000, 999999)

	testAccConfig := acctest.DataSourceTestAccConfigFunc("data.alicloud_dns_groups.default", fmt.Sprintf("tf-testacc-%d", rand), dataSourceDnsGroupsConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_dns_group.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_dns_group.default.name}_fake",
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_dns_group.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_dns_group.default.id}_fake"},
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_dns_group.default.name}",
			"ids":        []string{"${alicloud_dns_group.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_dns_group.default.name}_fake",
			"ids":        []string{"${alicloud_dns_group.default.id}"},
		}),
	}
	existChangeMap := map[string]string{
		"ids.#":               "1",
		"ids.0":               acctest.REMOVEKEY,
		"names.#":             "1",
		"names.0":             "ALL",
		"groups.#":            "1",
		"groups.0.group_id":   "",
		"groups.0.group_name": "ALL",
	}
	nameAllConf := acctest.DataSourceTestAccConfig{
		ExistConfig:   testAccCheckAlicloudDnsGroupsDataSourceNameRegexAll,
		ExistChangMap: existChangeMap,
	}

	dnsGroupsCheckInfo.DataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf, nameAllConf)
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
		"ids.0":               acctest.CHECKSET,
		"names.#":             "1",
		"names.0":             fmt.Sprintf("tf-testacc-%d", rand),
		"groups.0.group_id":   acctest.CHECKSET,
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

var dnsGroupsCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_dns_groups.default",
	ExistMapFunc: existDnsGroupsMapFunc,
	FakeMapFunc:  fakeDnsGroupsMapFunc,
}
