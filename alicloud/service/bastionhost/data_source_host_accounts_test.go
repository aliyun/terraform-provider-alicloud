package bastionhost_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudBastionhostHostAccountsDataSource(t *testing.T) {
	resourceId := "data.alicloud_bastionhost_host_accounts.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccBastionhostHostAccountsTest%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceBastionhostHostAccountsDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host_account.default.instance_id}",
			"host_id":     "${alicloud_bastionhost_host_account.default.host_id}",
			"ids":         []string{"${alicloud_bastionhost_host_account.default.host_account_id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host_account.default.instance_id}",
			"host_id":     "${alicloud_bastionhost_host_account.default.host_id}",
			"ids":         []string{"${alicloud_bastionhost_host_account.default.host_account_id}-fake"},
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host_account.default.instance_id}",
			"host_id":     "${alicloud_bastionhost_host_account.default.host_id}",
			"name_regex":  "${alicloud_bastionhost_host_account.default.host_account_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host_account.default.instance_id}",
			"host_id":     "${alicloud_bastionhost_host_account.default.host_id}",
			"name_regex":  "${alicloud_bastionhost_host_account.default.host_account_name}" + "fake",
		}),
	}
	hostAccountNameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_bastionhost_host_account.default.instance_id}",
			"host_id":           "${alicloud_bastionhost_host_account.default.host_id}",
			"host_account_name": "${alicloud_bastionhost_host_account.default.host_account_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_bastionhost_host_account.default.instance_id}",
			"host_id":           "${alicloud_bastionhost_host_account.default.host_id}",
			"host_account_name": "${alicloud_bastionhost_host_account.default.host_account_name}" + "fake",
		}),
	}
	protocolNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id":   "${alicloud_bastionhost_host_account.default.instance_id}",
			"ids":           []string{"${alicloud_bastionhost_host_account.default.host_account_id}"},
			"protocol_name": "${alicloud_bastionhost_host_account.default.protocol_name}",
			"host_id":       "${alicloud_bastionhost_host_account.default.host_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":   "${alicloud_bastionhost_host_account.default.instance_id}",
			"ids":           []string{"${alicloud_bastionhost_host_account.default.host_account_id}"},
			"protocol_name": "RDP",
			"host_id":       "${alicloud_bastionhost_host_account.default.host_id}",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_bastionhost_host_account.default.instance_id}",
			"ids":               []string{"${alicloud_bastionhost_host_account.default.host_account_id}"},
			"host_id":           "${alicloud_bastionhost_host_account.default.host_id}",
			"name_regex":        "${alicloud_bastionhost_host_account.default.host_account_name}",
			"host_account_name": "${alicloud_bastionhost_host_account.default.host_account_name}",
			"protocol_name":     "${alicloud_bastionhost_host_account.default.protocol_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"ids.0":                              acctest.CHECKSET,
			"names.#":                            "1",
			"names.0":                            name,
			"accounts.#":                         "1",
			"accounts.0.id":                      acctest.CHECKSET,
			"accounts.0.instance_id":             acctest.CHECKSET,
			"accounts.0.host_id":                 acctest.CHECKSET,
			"accounts.0.host_account_name":       acctest.CHECKSET,
			"accounts.0.has_password":            "true",
			"accounts.0.host_account_id":         acctest.CHECKSET,
			"accounts.0.private_key_fingerprint": "",
			"accounts.0.protocol_name":           acctest.CHECKSET,
		}
	}

	var fakeBastionhostHostAccountsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"accounts.#": "0",
		}
	}

	var BastionhostHostAccountsInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existBastionhostHostAccountsMapFunc,
		FakeMapFunc:  fakeBastionhostHostAccountsMapFunc,
	}

	BastionhostHostAccountsInfo.DataSourceTestCheck(t, 0, idsConf, nameRegexConf, hostAccountNameRegexConf, protocolNameConf, allConf)
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
