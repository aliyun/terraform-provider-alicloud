package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudBastionhostHostAccountsDataSource(t *testing.T) {
	resourceId := "data.alicloud_bastionhost_host_accounts.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccBastionhostHostAccountsTest%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceBastionhostHostAccountsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host_account.default.instance_id}",
			"host_id":     "${alicloud_bastionhost_host_account.default.host_id}",
			"ids":         []string{"${alicloud_bastionhost_host_account.default.host_account_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host_account.default.instance_id}",
			"host_id":     "${alicloud_bastionhost_host_account.default.host_id}",
			"ids":         []string{"${alicloud_bastionhost_host_account.default.host_account_id}-fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host_account.default.instance_id}",
			"host_id":     "${alicloud_bastionhost_host_account.default.host_id}",
			"name_regex":  "${alicloud_bastionhost_host_account.default.host_account_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host_account.default.instance_id}",
			"host_id":     "${alicloud_bastionhost_host_account.default.host_id}",
			"name_regex":  "${alicloud_bastionhost_host_account.default.host_account_name}" + "fake",
		}),
	}
	hostAccountNameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_bastionhost_host_account.default.instance_id}",
			"host_id":           "${alicloud_bastionhost_host_account.default.host_id}",
			"host_account_name": "${alicloud_bastionhost_host_account.default.host_account_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_bastionhost_host_account.default.instance_id}",
			"host_id":           "${alicloud_bastionhost_host_account.default.host_id}",
			"host_account_name": "${alicloud_bastionhost_host_account.default.host_account_name}" + "fake",
		}),
	}
	protocolNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":   "${alicloud_bastionhost_host_account.default.instance_id}",
			"ids":           []string{"${alicloud_bastionhost_host_account.default.host_account_id}"},
			"protocol_name": "${alicloud_bastionhost_host_account.default.protocol_name}",
			"host_id":       "${alicloud_bastionhost_host_account.default.host_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":   "${alicloud_bastionhost_host_account.default.instance_id}",
			"ids":           []string{"${alicloud_bastionhost_host_account.default.host_account_id}"},
			"protocol_name": "RDP",
			"host_id":       "${alicloud_bastionhost_host_account.default.host_id}",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_bastionhost_host_account.default.instance_id}",
			"ids":               []string{"${alicloud_bastionhost_host_account.default.host_account_id}"},
			"host_id":           "${alicloud_bastionhost_host_account.default.host_id}",
			"name_regex":        "${alicloud_bastionhost_host_account.default.host_account_name}",
			"host_account_name": "${alicloud_bastionhost_host_account.default.host_account_name}",
			"protocol_name":     "${alicloud_bastionhost_host_account.default.protocol_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_bastionhost_host_account.default.instance_id}",
			"ids":               []string{"${alicloud_bastionhost_host_account.default.host_account_id}"},
			"host_id":           "${alicloud_bastionhost_host_account.default.host_id}",
			"name_regex":        "${alicloud_bastionhost_host_account.default.host_account_name}",
			"host_account_name": "${alicloud_bastionhost_host_account.default.host_account_name}-fake",
		}),
	}
	var existBastionhostHostAccountsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"ids.0":                              CHECKSET,
			"names.#":                            "1",
			"names.0":                            name,
			"accounts.#":                         "1",
			"accounts.0.id":                      CHECKSET,
			"accounts.0.instance_id":             CHECKSET,
			"accounts.0.host_id":                 CHECKSET,
			"accounts.0.host_account_name":       CHECKSET,
			"accounts.0.has_password":            "true",
			"accounts.0.host_account_id":         CHECKSET,
			"accounts.0.private_key_fingerprint": "",
			"accounts.0.protocol_name":           CHECKSET,
		}
	}

	var fakeBastionhostHostAccountsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"accounts.#": "0",
		}
	}

	var BastionhostHostAccountsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existBastionhostHostAccountsMapFunc,
		fakeMapFunc:  fakeBastionhostHostAccountsMapFunc,
	}

	BastionhostHostAccountsInfo.dataSourceTestCheck(t, 0, idsConf, nameRegexConf, hostAccountNameRegexConf, protocolNameConf, allConf)
}

func dataSourceBastionhostHostAccountsDependence(name string) string {
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
resource "alicloud_bastionhost_host_account" "default" {
 instance_id          = alicloud_bastionhost_host.default.instance_id
 host_account_name = var.name
 host_id           = alicloud_bastionhost_host.default.host_id
 protocol_name     = "SSH"
 password          = "YourPassword12345"
}
`, name)
}
