package bastionhost_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudBastionhostHostsDataSource(t *testing.T) {
	resourceId := "data.alicloud_bastionhost_hosts.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccBastionhostHostsTest%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceBastionhostHostsDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_host.default.host_id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_host.default.id}-fake"},
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"name_regex":  "${alicloud_bastionhost_host.default.host_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"name_regex":  "${alicloud_bastionhost_host.default.host_name}" + "fake",
		}),
	}
	userNameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"host_name":   "${alicloud_bastionhost_host.default.host_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"host_name":   "${alicloud_bastionhost_host.default.host_name}" + "fake",
		}),
	}
	hostAddressConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id":  "${alicloud_bastionhost_host.default.instance_id}",
			"ids":          []string{"${alicloud_bastionhost_host.default.host_id}"},
			"host_address": "172.16.0.10",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":  "${alicloud_bastionhost_host.default.instance_id}",
			"ids":          []string{"${alicloud_bastionhost_host.default.host_id}"},
			"host_address": "172.16.0.1",
		}),
	}
	osTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_host.default.host_id}"},
			"os_type":     "Linux",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_host.default.host_id}"},
			"os_type":     "Windows",
		}),
	}
	sourceConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_host.default.host_id}"},
			"source":      "Local",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_host.default.host_id}"},
			"source":      "Rds",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id":  "${alicloud_bastionhost_host.default.instance_id}",
			"name_regex":   name,
			"host_name":    name,
			"ids":          []string{"${alicloud_bastionhost_host.default.host_id}"},
			"host_address": "172.16.0.10",
			"os_type":      "Linux",
			"source":       "Local",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host.default.instance_id}",
			"host_name":   name + "fake",
			"name_regex":  name + "fake",
			"ids":         []string{"${alicloud_bastionhost_host.default.id}-fake"},
		}),
	}
	var existBastionhostHostsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"ids.0":               acctest.CHECKSET,
			"names.#":             "1",
			"names.0":             name,
			"hosts.#":             "1",
			"hosts.0.id":          acctest.CHECKSET,
			"hosts.0.comment":     "",
			"hosts.0.instance_id": acctest.CHECKSET,
			"hosts.0.host_id":     acctest.CHECKSET,
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

	var BastionhostHostsInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existBastionhostHostsMapFunc,
		FakeMapFunc:  fakeBastionhostHostsMapFunc,
	}

	BastionhostHostsInfo.DataSourceTestCheck(t, 0, idsConf, nameRegexConf, userNameRegexConf, hostAddressConf, osTypeConf, sourceConf, allConf)
}

func dataSourceBastionhostHostsDependence(name string) string {
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
`, name)
}
