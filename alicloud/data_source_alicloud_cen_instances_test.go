package alicloud

import (
	"strings"
	"testing"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/acctest"

	"fmt"
)

func TestAccAlicloudCenInstancesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 99999999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccAlicloudCenInstancesDataSourceConfig(rand, map[string]string{
			"ids": `[alicloud_cen_instance.default.id]`,
		}),
		fakeConfig: testAccAlicloudCenInstancesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cen_instance.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccAlicloudCenInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": `alicloud_cen_instance.default.name`,
		}),
		fakeConfig: testAccAlicloudCenInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_instance.default.name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccAlicloudCenInstancesDataSourceConfig(rand, map[string]string{
			"ids":        `[alicloud_cen_instance.default.id]`,
			"name_regex": `alicloud_cen_instance.default.name`,
		}),
		fakeConfig: testAccAlicloudCenInstancesDataSourceConfig(rand, map[string]string{
			"ids":        `[alicloud_cen_instance.default.id]`,
			"name_regex": `"${alicloud_cen_instance.default.name}_fake"`,
		}),
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.CenNoSkipRegions)
	}
	CenInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func testAccAlicloudCenInstancesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
		resource "alicloud_cen_instance" "default" {
			name = "tf-testAcc%sCenInstancesDataSourceCen-%d"
			description = "tf-testAccCenConfigDescription"
		}

		data "alicloud_cen_instances" "default" {
			%s
		}
`, defaultRegionToTest, rand, strings.Join(pairs, "\n  "))
	return config
}

var existCenInstancesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"names.#":                 "1",
		"instances.#":             "1",
		"instances.0.name":        fmt.Sprintf("tf-testAcc%sCenInstancesDataSourceCen-%d", defaultRegionToTest, rand),
		"instances.0.description": "tf-testAccCenConfigDescription",
	}
}

var fakeCenInstancesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"names.#":     "0",
		"instances.#": "0",
	}
}

var existCenInstancesMultiMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"names.#":                 "5",
		"instances.#":             "5",
		"instances.0.name":        fmt.Sprintf("tf-testAcc%sCenInstancesDataSourceCen-%d", defaultRegionToTest, rand),
		"instances.0.description": "tf-testAccCenConfigDescription",
	}
}

var CenInstancesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cen_instances.default",
	existMapFunc: existCenInstancesMapFunc,
	fakeMapFunc:  fakeCenInstancesMapFunc,
}
var CenInstancesCheckInfoMulti = dataSourceAttr{
	resourceId:   "data.alicloud_cen_instances.default",
	existMapFunc: existCenInstancesMultiMapFunc,
	fakeMapFunc:  fakeCenInstancesMapFunc,
}
