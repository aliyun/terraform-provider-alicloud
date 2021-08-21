package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudBastionhostUserGroupsDataSource(t *testing.T) {
	resourceId := "data.alicloud_bastionhost_user_groups.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccBastionhostUserGroupsTest%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceBastionhostUserGroupsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_user_group.default.user_group_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_user_group.default.id}-fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.instance_id}",
			"name_regex":  name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.instance_id}",
			"name_regex":  name + "fake",
		}),
	}
	userGroupNameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":     "${local.instance_id}",
			"user_group_name": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":     "${local.instance_id}",
			"user_group_name": name + "fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":     "${local.instance_id}",
			"name_regex":      name,
			"user_group_name": name,
			"ids":             []string{"${alicloud_bastionhost_user_group.default.user_group_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":     "${local.instance_id}",
			"user_group_name": name + "fake",
			"name_regex":      name + "fake",
			"ids":             []string{"${alicloud_bastionhost_user_group.default.id}-fake"},
		}),
	}
	var existBastionhostUserGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"ids.0":                    CHECKSET,
			"names.#":                  "1",
			"names.0":                  name,
			"groups.#":                 "1",
			"groups.0.id":              CHECKSET,
			"groups.0.comment":         "",
			"groups.0.instance_id":     CHECKSET,
			"groups.0.user_group_id":   CHECKSET,
			"groups.0.user_group_name": name,
		}
	}

	var fakeBastionhostUserGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"groups.#": "0",
		}
	}

	var BastionhostUserGroupsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existBastionhostUserGroupsMapFunc,
		fakeMapFunc:  fakeBastionhostUserGroupsMapFunc,
	}

	BastionhostUserGroupsInfo.dataSourceTestCheck(t, 0, idsConf, nameRegexConf, userGroupNameRegexConf, allConf)
}

func dataSourceBastionhostUserGroupsDependence(name string) string {
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
  zone_id = local.zone_id
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
  instance_id = length(data.alicloud_yundun_bastionhost_instances.default.ids) > 0 ? data.alicloud_yundun_bastionhost_instances.default.ids.0 : concat(alicloud_yundun_bastionhost_instance.default.*.id, [""])[0]
}
data "alicloud_yundun_bastionhost_instances" "default" {}
resource "alicloud_yundun_bastionhost_instance" "default" {
  count              = length(data.alicloud_yundun_bastionhost_instances.default.ids) > 0 ? 0 : 1
  description        = var.name
  license_code       = "bhah_ent_50_asset"
  period             = "1"
  vswitch_id         = local.vswitch_id
  security_group_ids = [alicloud_security_group.default.id]
}
resource "alicloud_bastionhost_user_group" "default" {
  instance_id     = local.instance_id
  user_group_name = var.name
}
`, name)
}
