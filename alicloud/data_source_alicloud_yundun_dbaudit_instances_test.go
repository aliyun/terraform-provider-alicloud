package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudYundunDbauditInstanceDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_yundun_dbaudit_instance.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, fmt.Sprintf("tf_testAcc%d", rand),
		dataSourceYundunDbauditInstanceConfigDependency)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_yundun_dbaudit_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_yundun_dbaudit_instance.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_yundun_dbaudit_instance.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_yundun_dbaudit_instance.default.description}-fake",
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_yundun_dbaudit_instance.default.id}"},
			"tags": map[string]interface{}{
				"Created": "TF",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_yundun_dbaudit_instance.default.id}-fake"},
			"tags": map[string]interface{}{
				"Created": "TF-fake",
			},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_yundun_dbaudit_instance.default.description}",
			"ids":               []string{"${alicloud_yundun_dbaudit_instance.default.id}"},
			"tags": map[string]interface{}{
				"For": "acceptance test",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
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
			"ids.0":                             CHECKSET,
			"descriptions.0":                    fmt.Sprintf("tf_testAcc%d", rand),
			"instances.#":                       "1",
			"instances.0.description":           fmt.Sprintf("tf_testAcc%d", rand),
			"instances.0.license_code":          "alpha.professional",
			"instances.0.user_vswitch_id":       CHECKSET,
			"instances.0.public_network_access": "false",
			"instances.0.private_domain":        CHECKSET,
			"instances.0.instance_status":       CHECKSET,
		}
	}
	var fakeYundunDbauditInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"descriptions.#": "0",
		}
	}
	var yundunDbauditInstanceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_yundun_dbaudit_instance.default",
		existMapFunc: existYundunDbauditInstanceMapFunc,
		fakeMapFunc:  fakeYundunDbauditInstanceMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.YundunDbauditSupportedRegions)
		testAccPreCheckWithAccountSiteType(t, DomesticSite)
	}

	yundunDbauditInstanceCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, tagsConf, allConf)

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
