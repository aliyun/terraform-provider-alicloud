package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDtsInstanceDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDtsInstanceSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_dts_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDtsInstanceSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_dts_instance.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDtsInstanceSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_dts_instance.default.instance_name}"`,
			"ids":        `["${alicloud_dts_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDtsInstanceSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_dts_instance.default.id}"]`,
			"name_regex": `"${alicloud_dts_instance.default.instance_name}_fake"`,
		}),
	}

	ResourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDtsInstanceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_dts_instance.default.id}"]`,
			"resource_group_id": `"${alicloud_dts_instance.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudDtsInstanceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_dts_instance.default.id}_fake"]`,
			"resource_group_id": `"${alicloud_dts_instance.default.resource_group_id}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDtsInstanceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_dts_instance.default.id}"]`,
			"resource_group_id": `"${alicloud_dts_instance.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudDtsInstanceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_dts_instance.default.id}_fake"]`,
			"resource_group_id": `"${alicloud_dts_instance.default.resource_group_id}_fake"`,
		}),
	}

	DtsInstanceCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, ResourceGroupIdConf, allConf)
}

var existDtsInstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                   "1",
		"names.#":                 "1",
		"instances.#":             "1",
		"instances.0.id":          CHECKSET,
		"instances.0.create_time": CHECKSET,
		"instances.0.destination_endpoint_engine_name": CHECKSET,
		"instances.0.dts_instance_id":                  CHECKSET,
		"instances.0.instance_name":                    CHECKSET,
		"instances.0.instance_class":                   CHECKSET,
		"instances.0.payment_type":                     CHECKSET,
		"instances.0.resource_group_id":                CHECKSET,
		"instances.0.source_endpoint_engine_name":      CHECKSET,
		"instances.0.source_region":                    CHECKSET,
		"instances.0.destination_region":               CHECKSET,
		"instances.0.status":                           CHECKSET,
		"instances.0.type":                             CHECKSET,
		"instances.0.tags.%":                           "2",
		"instances.0.tags.Create":                      "TF",
		"instances.0.tags.For":                         "acceptance test",
	}
}

var fakeDtsInstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":       "0",
		"names.#":     "0",
		"instances.#": "0",
	}
}

var DtsInstanceCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_dts_instances.default",
	existMapFunc: existDtsInstanceMapFunc,
	fakeMapFunc:  fakeDtsInstanceMapFunc,
}

func testAccCheckAlicloudDtsInstanceSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccDtsInstance%d"
}

variable "region" {
	default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_dts_instance" "default" {
  type                             = "sync"
  resource_group_id                = data.alicloud_resource_manager_resource_groups.default.ids.0
  payment_type                     = "PayAsYouGo"
  instance_class                   = "large"
  source_endpoint_engine_name      = "MySQL"
  source_region                    = var.region
  destination_endpoint_engine_name = "MySQL"
  destination_region               = var.region
  tags = {
	  Create = "TF"
	  For = "acceptance test",
	}
}

data "alicloud_dts_instances" "default" {
%s
}
`, rand, os.Getenv("ALICLOUD_REGION"), strings.Join(pairs, "\n   "))
	return config
}
