package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudBastionhostUsersDataSource(t *testing.T) {
	resourceId := "data.alicloud_bastionhost_users.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccBastionhostUsersTest%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceBastionhostUsersDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_user.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_user.default.user_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_user.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_user.default.id}-fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_user.default.instance_id}",
			"name_regex":  "${alicloud_bastionhost_user.default.user_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_user.default.instance_id}",
			"name_regex":  "${alicloud_bastionhost_user.default.user_name}" + "fake",
		}),
	}
	userNameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_user.default.instance_id}",
			"user_name":   "${alicloud_bastionhost_user.default.user_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_user.default.instance_id}",
			"user_name":   "${alicloud_bastionhost_user.default.user_name}" + "fake",
		}),
	}
	displayNameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":  "${alicloud_bastionhost_user.default.instance_id}",
			"display_name": "${alicloud_bastionhost_user.default.user_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":  "${alicloud_bastionhost_user.default.instance_id}",
			"display_name": "${alicloud_bastionhost_user.default.user_name}" + "fake",
		}),
	}
	mobileRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_user.default.instance_id}",
			"mobile":      "${alicloud_bastionhost_user.default.mobile}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_user.default.instance_id}",
			"mobile":      "${alicloud_bastionhost_user.default.mobile}" + "1",
		}),
	}
	sourceRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_user.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_user.default.user_id}"},
			"source":      "${alicloud_bastionhost_user.default.source}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_user.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_user.default.user_id}"},
			"source":      "Ram",
		}),
	}
	statusRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_user.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_user.default.user_id}"},
			"status":      "${alicloud_bastionhost_user.default.status}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_user.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_user.default.user_id}"},
			"status":      "Expired",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":  "${alicloud_bastionhost_user.default.instance_id}",
			"name_regex":   name,
			"user_name":    name,
			"ids":          []string{"${alicloud_bastionhost_user.default.user_id}"},
			"display_name": "${alicloud_bastionhost_user.default.user_name}",
			"mobile":       "${alicloud_bastionhost_user.default.mobile}",
			"source":       "${alicloud_bastionhost_user.default.source}",
			"status":       "${alicloud_bastionhost_user.default.status}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":  "${alicloud_bastionhost_user.default.instance_id}",
			"user_name":    name + "fake",
			"name_regex":   name + "fake",
			"ids":          []string{"${alicloud_bastionhost_user.default.id}-fake"},
			"display_name": "${alicloud_bastionhost_user.default.user_name}",
			"mobile":       "${alicloud_bastionhost_user.default.mobile}",
			"source":       "${alicloud_bastionhost_user.default.source}",
			"status":       "${alicloud_bastionhost_user.default.status}",
		}),
	}
	var existBastionhostUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"ids.0":                       CHECKSET,
			"names.#":                     "1",
			"names.0":                     name,
			"users.#":                     "1",
			"users.0.id":                  CHECKSET,
			"users.0.comment":             "",
			"users.0.instance_id":         CHECKSET,
			"users.0.user_id":             CHECKSET,
			"users.0.user_name":           name,
			"users.0.display_name":        name,
			"users.0.email":               "",
			"users.0.mobile":              CHECKSET,
			"users.0.mobile_country_code": CHECKSET,
			"users.0.source":              CHECKSET,
			"users.0.source_user_id":      "",
			"users.0.status":              "Normal",
		}
	}

	var fakeBastionhostUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"users.#": "0",
		}
	}

	var BastionhostUsersInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existBastionhostUsersMapFunc,
		fakeMapFunc:  fakeBastionhostUsersMapFunc,
	}

	BastionhostUsersInfo.dataSourceTestCheck(t, 0, idsConf, nameRegexConf, userNameRegexConf, displayNameRegexConf, mobileRegexConf, sourceRegexConf, statusRegexConf, allConf)
}

func dataSourceBastionhostUsersDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_bastionhost_instances" "default" {}

resource "alicloud_bastionhost_user" "default" {
  instance_id     = data.alicloud_bastionhost_instances.default.ids.0
  mobile         = "13312345678"
  mobile_country_code = "CN"
  password       = "YourPassword-123"
  source         = "Local"
  user_name      = var.name
}
`, name)
}
