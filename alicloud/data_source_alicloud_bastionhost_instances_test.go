package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudBastionhostInstancesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_bastionhost_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, fmt.Sprintf("tf_testAcc%d", rand),
		dataSourceYundunBastionhostInstanceConfigDependency)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_bastionhost_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_bastionhost_instance.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_bastionhost_instance.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_bastionhost_instance.default.description}-fake",
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_bastionhost_instance.default.id}"},
			"tags": map[string]interface{}{
				"Created": "TF",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_bastionhost_instance.default.id}-fake"},
			"tags": map[string]interface{}{
				"Created": "TF-fake",
			},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_bastionhost_instance.default.description}",
			"ids":               []string{"${alicloud_bastionhost_instance.default.id}"},
			"tags": map[string]interface{}{
				"For": "acceptance test",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_bastionhost_instance.default.description}-fake",
			"ids":               []string{"${alicloud_bastionhost_instance.default.id}-fake"},
			"tags": map[string]interface{}{
				"For": "acceptance test-fake",
			},
		}),
	}

	var existYundunBastionhostInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"descriptions.#":                    "1",
			"ids.0":                             CHECKSET,
			"descriptions.0":                    fmt.Sprintf("tf_testAcc%d", rand),
			"instances.#":                       "1",
			"instances.0.description":           fmt.Sprintf("tf_testAcc%d", rand),
			"instances.0.license_code":          "bhah_ent_50_asset",
			"instances.0.user_vswitch_id":       CHECKSET,
			"instances.0.public_network_access": "true",
			"instances.0.private_domain":        CHECKSET,
			"instances.0.instance_status":       CHECKSET,
			"instances.0.security_group_ids.#":  "1",
		}
	}
	var fakeYundunBastionhostInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"descriptions.#": "0",
		}
	}
	var yundunBastionhostInstanceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_bastionhost_instances.default",
		existMapFunc: existYundunBastionhostInstanceMapFunc,
		fakeMapFunc:  fakeYundunBastionhostInstanceMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithAccountSiteType(t, DomesticSite)
	}

	yundunBastionhostInstanceCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, tagsConf, allConf)

}

func dataSourceYundunBastionhostInstanceConfigDependency(description string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_vswitch" "this" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vswitch_name = var.name
  vpc_id       = data.alicloud_vpcs.default.ids.0
  zone_id      = data.alicloud_zones.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  name   = var.name
}
locals {
  vswitch_id  = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id     = data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 1]
  instance_id = length(data.alicloud_bastionhost_instances.default.ids) > 0 ? data.alicloud_bastionhost_instances.default.ids.0 : concat(alicloud_bastionhost_instance.default.*.id, [""])[0]
}
				
resource "alicloud_bastionhost_instance" "default" {
  description        = "${var.name}"
  license_code       = "bhah_ent_50_asset"
  period             = "1"
  vswitch_id         = local.vswitch_id
  security_group_ids = ["${alicloud_security_group.default.id}"]
  tags 				 = {
		Created = "TF"
		For 	= "acceptance test"
  }
}`, description)
}
