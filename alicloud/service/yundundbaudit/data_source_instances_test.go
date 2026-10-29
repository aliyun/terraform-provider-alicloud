package yundundbaudit_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudYundunDbauditInstanceDataSource_basic(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_yundun_dbaudit_instance.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, fmt.Sprintf("tf_testAcc%d", rand),
		dataSourceYundunDbauditInstanceConfigDependency)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_yundun_dbaudit_instance.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_yundun_dbaudit_instance.default.id}-fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_yundun_dbaudit_instance.default.description}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_yundun_dbaudit_instance.default.description}-fake",
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_yundun_dbaudit_instance.default.id}"},
			"tags": map[string]interface{}{
				"Created": "TF",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_yundun_dbaudit_instance.default.id}-fake"},
			"tags": map[string]interface{}{
				"Created": "TF-fake",
			},
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_yundun_dbaudit_instance.default.description}",
			"ids":               []string{"${alicloud_yundun_dbaudit_instance.default.id}"},
			"tags": map[string]interface{}{
				"For": "acceptance test",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_yundun_dbaudit_instance.default.description}-fake",
			"ids":               []string{"${alicloud_yundun_dbaudit_instance.default.id}-fake"},
			"tags": map[string]interface{}{
				"For": "acceptance test-fake",
			},
		}),
	}

	var existYundunDbauditInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"descriptions.#":                    "1",
			"ids.0":                             acctest.CHECKSET,
			"descriptions.0":                    fmt.Sprintf("tf_testAcc%d", rand),
			"instances.#":                       "1",
			"instances.0.description":           fmt.Sprintf("tf_testAcc%d", rand),
			"instances.0.license_code":          "alpha.professional",
			"instances.0.user_vswitch_id":       acctest.CHECKSET,
			"instances.0.public_network_access": "false",
			"instances.0.private_domain":        acctest.CHECKSET,
			"instances.0.instance_status":       acctest.CHECKSET,
		}
	}
	var fakeYundunDbauditInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"descriptions.#": "0",
		}
	}
	var yundunDbauditInstanceCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_yundun_dbaudit_instance.default",
		ExistMapFunc: existYundunDbauditInstanceMapFunc,
		FakeMapFunc:  fakeYundunDbauditInstanceMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckPrePaidResources(t)
		acctest.PreCheckWithRegions(t, true, connectivity.YundunDbauditSupportedRegions)
		acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
	}

	yundunDbauditInstanceCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, tagsConf, allConf)

}

func dataSourceYundunDbauditInstanceConfigDependency(description string) string {
	return fmt.Sprintf(`
data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

variable "name" {
	default = "%s"
}

data "alicloud_vpcs" "default"{
	name_regex = "default-NODELETING"
}
data "alicloud_slb_zones" "default" {
	available_slb_address_type = "vpc"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_yundun_dbaudit_instance" "default" {
	description       = "${var.name}"
	plan_code         = "alpha.professional"
	period            = "1"
	vswitch_id        = data.alicloud_vswitches.default.ids.0
	tags 				 = {
		Created = "TF"
		For 	= "acceptance test"
  }
}`, description)
}
