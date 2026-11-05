package cen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCenInstancesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 99999999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccAlicloudCenInstancesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cen_instance.default.id}"]`,
		}),
		FakeConfig: testAccAlicloudCenInstancesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cen_instance.default.id}_fake"]`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccAlicloudCenInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_instance.default.name}"`,
			"status":     `"Active"`,
		}),
		FakeConfig: testAccAlicloudCenInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_instance.default.name}_fake"`,
			"status":     `"Deleting"`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccAlicloudCenInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_instance.default.name}"`,
		}),
		FakeConfig: testAccAlicloudCenInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_instance.default.name}_fake"`,
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccAlicloudCenInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_instance.default.name}"`,
			"tags": `{
							Created = "TF"
							For 	= "acceptance test"
					  }`,
		}),
		FakeConfig: testAccAlicloudCenInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_instance.default.name}"`,
			"tags": `{
							Created = "TF-fake"
							For 	= "acceptance test-fake"
					  }`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccAlicloudCenInstancesDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_cen_instance.default.id}"]`,
			"name_regex": `"${alicloud_cen_instance.default.name}"`,
		}),
		FakeConfig: testAccAlicloudCenInstancesDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_cen_instance.default.id}"]`,
			"name_regex": `"${alicloud_cen_instance.default.name}_fake"`,
		}),
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.CenNoSkipRegions)
	}
	CenInstancesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, nameRegexConf, tagsConf, allConf)
}

func testAccAlicloudCenInstancesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
		resource "alicloud_cen_instance" "default" {
			cen_instance_name = "tf-testAcc%sCenInstancesDataSourceCen-%d"
			description = "tf-testAccCenConfigDescription"
			tags 		= {
				Created = "TF"
				For 	= "acceptance test"
			}
		}

		data "alicloud_cen_instances" "default" {
			%s
		}
`, acctest.DefaultRegionToTest, rand, strings.Join(pairs, "\n  "))
	return config
}

var existCenInstancesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"names.#":     "1",
		"instances.#": "1",
		"instances.0.cen_bandwidth_package_ids.#": "0",
		"instances.0.id":                acctest.CHECKSET,
		"instances.0.cen_id":            acctest.CHECKSET,
		"instances.0.description":       "tf-testAccCenConfigDescription",
		"instances.0.name":              fmt.Sprintf("tf-testAcc%sCenInstancesDataSourceCen-%d", acctest.DefaultRegionToTest, rand),
		"instances.0.cen_instance_name": fmt.Sprintf("tf-testAcc%sCenInstancesDataSourceCen-%d", acctest.DefaultRegionToTest, rand),
		"instances.0.protection_level":  "REDUCED",
		"instances.0.status":            "Active",
		"instances.0.tags.%":            "2",
		"instances.0.tags.Created":      "TF",
		"instances.0.tags.For":          "acceptance test",
		"instances.0.create_time":       acctest.CHECKSET,
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
		"instances.0.name":        fmt.Sprintf("tf-testAcc%sCenInstancesDataSourceCen-%d", acctest.DefaultRegionToTest, rand),
		"instances.0.description": "tf-testAccCenConfigDescription",
	}
}

var CenInstancesCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_cen_instances.default",
	ExistMapFunc: existCenInstancesMapFunc,
	FakeMapFunc:  fakeCenInstancesMapFunc,
}
var CenInstancesCheckInfoMulti = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_cen_instances.default",
	ExistMapFunc: existCenInstancesMultiMapFunc,
	FakeMapFunc:  fakeCenInstancesMapFunc,
}
