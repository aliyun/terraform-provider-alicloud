package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudBastionhostHostsDataSource(t *testing.T) {
	resourceId := "data.alicloud_bastionhost_hosts.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccBastionhostHostsTest%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceBastionhostHostsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_host.default.host_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_host.default.id}-fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"name_regex":  "${alicloud_bastionhost_host.default.host_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"name_regex":  "${alicloud_bastionhost_host.default.host_name}" + "fake",
		}),
	}
	userNameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"host_name":   "${alicloud_bastionhost_host.default.host_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"host_name":   "${alicloud_bastionhost_host.default.host_name}" + "fake",
		}),
	}
	hostAddressConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":  "${alicloud_bastionhost_host.default.instance_id}",
			"ids":          []string{"${alicloud_bastionhost_host.default.host_id}"},
			"host_address": "172.16.0.10",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":  "${alicloud_bastionhost_host.default.instance_id}",
			"ids":          []string{"${alicloud_bastionhost_host.default.host_id}"},
			"host_address": "172.16.0.1",
		}),
	}
	osTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_host.default.host_id}"},
			"os_type":     "Linux",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_host.default.host_id}"},
			"os_type":     "Windows",
		}),
	}
	sourceConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_host.default.host_id}"},
			"source":      "Local",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_host.default.host_id}"},
			"source":      "Rds",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":  "${alicloud_bastionhost_host.default.instance_id}",
			"name_regex":   name,
			"host_name":    name,
			"ids":          []string{"${alicloud_bastionhost_host.default.host_id}"},
			"host_address": "172.16.0.10",
			"os_type":      "Linux",
			"source":       "Local",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"host_name":   name + "fake",
			"name_regex":  name + "fake",
			"ids":         []string{"${alicloud_bastionhost_host.default.id}-fake"},
		}),
	}
	var existBastionhostHostsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"ids.0":               CHECKSET,
			"names.#":             "1",
			"names.0":             name,
			"hosts.#":             "1",
			"hosts.0.id":          CHECKSET,
			"hosts.0.comment":     "",
			"hosts.0.instance_id": CHECKSET,
			"hosts.0.host_id":     CHECKSET,
			"hosts.0.host_name":   name,
		}
	}

	var fakeBastionhostHostsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"hosts.#": "0",
		}
	}

	var BastionhostHostsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existBastionhostHostsMapFunc,
		fakeMapFunc:  fakeBastionhostHostsMapFunc,
	}

	BastionhostHostsInfo.dataSourceTestCheck(t, 0, idsConf, nameRegexConf, userNameRegexConf, hostAddressConf, osTypeConf, sourceConf, allConf)
}

func dataSourceBastionhostHostsDependence(name string) string {
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
}
resource "alicloud_bastionhost_instance" "default" {
 description        = var.name
 license_code       = "bhah_ent_50_asset"
 period             = "1"
 vswitch_id         = local.vswitch_id
 security_group_ids = [alicloud_security_group.default.id]
}
resource "alicloud_bastionhost_host" "default" {
  instance_id          = alicloud_bastionhost_instance.default.id
  host_name            = var.name
  active_address_type  = "Private"
  host_private_address = "172.16.0.10"
  os_type              = "Linux"
  source               = "Local"
}
`, name)
}
