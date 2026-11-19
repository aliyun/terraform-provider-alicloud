package ecd_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEcdDesktopsDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_ecd_desktops.default"
	name := fmt.Sprintf("tf_testaccdesktop%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceEcdDesktopsConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecd_desktop.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecd_desktop.default.id}_fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ecd_desktop.default.desktop_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ecd_desktop.default.desktop_name}_fake",
		}),
	}

	officeSiteIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"office_site_id": "${alicloud_ecd_desktop.default.office_site_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"office_site_id": "${alicloud_ecd_desktop.default.office_site_id}_fake",
		}),
	}

	policyGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"policy_group_id": "${alicloud_ecd_desktop.default.policy_group_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"policy_group_id": "${alicloud_ecd_desktop.default.policy_group_id}_fake",
		}),
	}

	deskTopNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"desktop_name": "${alicloud_ecd_desktop.default.desktop_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"desktop_name": "${alicloud_ecd_desktop.default.desktop_name}_fake",
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"status": "Running",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"status": "Stopped",
		}),
	}

	endUserIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"end_user_ids": []string{"${alicloud_ecd_user.default.end_user_id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"end_user_ids": []string{"${alicloud_ecd_user.default.end_user_id}_fake"},
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":             []string{"${alicloud_ecd_desktop.default.id}"},
			"name_regex":      "${alicloud_ecd_desktop.default.desktop_name}",
			"office_site_id":  "${alicloud_ecd_desktop.default.office_site_id}",
			"policy_group_id": "${alicloud_ecd_desktop.default.policy_group_id}",
			"desktop_name":    "${alicloud_ecd_desktop.default.desktop_name}",
			"status":          "Running",
			"end_user_ids":    []string{"${alicloud_ecd_user.default.end_user_id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":             []string{"${alicloud_ecd_desktop.default.id}_fake"},
			"name_regex":      "${alicloud_ecd_desktop.default.desktop_name}_fake",
			"office_site_id":  "${alicloud_ecd_desktop.default.office_site_id}_fake",
			"policy_group_id": "${alicloud_ecd_desktop.default.policy_group_id}_fake",
			"desktop_name":    "${alicloud_ecd_desktop.default.desktop_name}_fake",
			"status":          "Stopped",
			"end_user_ids":    []string{"${alicloud_ecd_user.default.end_user_id}_fake"},
		}),
	}

	var existAliCloudEcdDesktopsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"names.#":                              "1",
			"desktops.#":                           "1",
			"desktops.0.id":                        acctest.CHECKSET,
			"desktops.0.desktop_id":                acctest.CHECKSET,
			"desktops.0.desktop_type":              acctest.CHECKSET,
			"desktops.0.desktop_name":              acctest.CHECKSET,
			"desktops.0.policy_group_id":           acctest.CHECKSET,
			"desktops.0.directory_id":              acctest.CHECKSET,
			"desktops.0.network_interface_id":      acctest.CHECKSET,
			"desktops.0.cpu":                       acctest.CHECKSET,
			"desktops.0.image_id":                  acctest.CHECKSET,
			"desktops.0.memory":                    acctest.CHECKSET,
			"desktops.0.system_disk_size":          acctest.CHECKSET,
			"desktops.0.payment_type":              acctest.CHECKSET,
			"desktops.0.status":                    acctest.CHECKSET,
			"desktops.0.create_time":               acctest.CHECKSET,
			"desktops.0.expired_time":              acctest.CHECKSET,
			"desktops.0.end_user_ids.#":            acctest.CHECKSET,
			"desktops.0.disks.#":                   acctest.CHECKSET,
			"desktops.0.disks.0.disk_id":           acctest.CHECKSET,
			"desktops.0.disks.0.disk_type":         acctest.CHECKSET,
			"desktops.0.disks.0.disk_category":     acctest.CHECKSET,
			"desktops.0.disks.0.disk_size":         acctest.CHECKSET,
			"desktops.0.disks.0.performance_level": acctest.CHECKSET,
		}
	}

	var fakeAliCloudEcdDesktopsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"desktops.#": "0",
		}
	}

	var aliCloudEcdDesktopsBusesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ecd_desktops.default",
		ExistMapFunc: existAliCloudEcdDesktopsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudEcdDesktopsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.EcdUserSupportRegions)
	}

	aliCloudEcdDesktopsBusesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, officeSiteIdConf, policyGroupIdConf, deskTopNameConf, statusConf, endUserIdConf, allConf)
}

func dataSourceEcdDesktopsConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_ecd_images" "default" {
  image_type            = "SYSTEM"
  os_type               = "Windows"
  desktop_instance_type = "eds.hf.4c8g"
}

data "alicloud_ecd_desktop_types" "default" {
  instance_type_family = "eds.hf"
  cpu_count            = 4
  memory_size          = 8192
}

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block          = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = var.name
}

resource "alicloud_ecd_bundle" "default" {
  desktop_type                = data.alicloud_ecd_desktop_types.default.ids.0
  image_id                    = data.alicloud_ecd_images.default.ids.0
  user_disk_size_gib          = [110]
  root_disk_size_gib          = 120
  root_disk_performance_level = "PL1"
  user_disk_performance_level = "PL1"
  bundle_name                 = var.name
  description                 = var.name
}

resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = var.name
  clipboard         = "readwrite"
  local_drive       = "read"
  authorize_access_policy_rules {
    description = var.name
    cidr_ip     = "1.2.3.4/24"
  }
  authorize_security_policy_rules {
    type        = "inflow"
    policy      = "accept"
    description = var.name
    port_range  = "80/80"
    ip_protocol = "TCP"
    priority    = "1"
    cidr_ip     = "0.0.0.0/0"
  }
}

resource "alicloud_ecd_user" "default" {
  end_user_id = var.name
  email       = "${var.name}@aaa.com"
}

resource "alicloud_ecd_desktop" "default" {
  office_site_id  = alicloud_ecd_simple_office_site.default.id
  policy_group_id = alicloud_ecd_policy_group.default.id
  bundle_id       = alicloud_ecd_bundle.default.id
  desktop_name    = var.name
  end_user_ids    = [alicloud_ecd_user.default.id]
}
`, name)
}
