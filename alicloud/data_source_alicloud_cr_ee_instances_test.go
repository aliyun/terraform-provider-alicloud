// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCrInstanceDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCrInstanceSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cr_ee_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCrInstanceSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cr_ee_instance.default.id}_fake"]`,
		}),
	}

	InstanceNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCrInstanceSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_cr_ee_instance.default.id}"]`,
			"instance_name": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCrInstanceSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_cr_ee_instance.default.id}_fake"]`,
			"instance_name": `"${var.name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCrInstanceSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_cr_ee_instance.default.id}"]`,
			"instance_name": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCrInstanceSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_cr_ee_instance.default.id}_fake"]`,
			"instance_name": `"${var.name}_fake"`,
		}),
	}

	CrInstanceCheckInfo.dataSourceTestCheck(t, rand, idsConf, InstanceNameConf, allConf)
}

var existCrInstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#":                   "1",
		"instances.0.record_total":      CHECKSET,
		"instances.0.resource_group_id": CHECKSET,
		"instances.0.modified_time":     CHECKSET,
		"instances.0.instance_id":       CHECKSET,
		"instances.0.instance_issue":    CHECKSET,
		"instances.0.create_time":       CHECKSET,
		"instances.0.instance_name":     CHECKSET,
		"instances.0.region_id":         CHECKSET,
	}
}

var fakeCrInstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#": "0",
	}
}

var CrInstanceCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cr_ee_instances.default",
	existMapFunc: existCrInstanceMapFunc,
	fakeMapFunc:  fakeCrInstanceMapFunc,
}

func testAccCheckAlicloudCrInstanceSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCrInstance%d"
}


resource "alicloud_cr_ee_instance" "default" {
  default_oss_bucket = "true"
  instance_name      = "tf-test-eco-526"
  renewal_status     = "ManualRenewal"
  image_scanner      = "DISABLE"
  period             = 1
  payment_type       = "Subscription"
  instance_type      = "Economy"
}

data "alicloud_cr_ee_instances" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
