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
data "alicloud_bastionhost_instances" "default" {}

resource "alicloud_bastionhost_host" "default" {
 instance_id          = data.alicloud_bastionhost_instances.default.ids.0
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
