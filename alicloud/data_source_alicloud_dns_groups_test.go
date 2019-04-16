package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudDnsGroupsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100000, 999999)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDnsGroupsDataSource(rand, map[string]string{
			"name_regex": `"${alicloud_dns_group.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDnsGroupsDataSource(rand, map[string]string{
			"name_regex": `"${alicloud_dns_group.default.name}_fake"`,
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

const testAccCheckAlicloudDnsGroupsDataSourceNameRegexAll = `
data "alicloud_dns_groups" "default" {
  name_regex = "^ALL"
}`

func testAccCheckAlicloudDnsGroupsDataSource(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
resource "alicloud_dns_group" "default" {
	name = "tf-testacc-%d"
}

data "alicloud_dns_groups" "default" {
	%s
}`, rand, strings.Join(pairs, "\n	"))
}

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
