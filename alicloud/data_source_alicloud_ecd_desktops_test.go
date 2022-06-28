package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECDDesktopsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.EcdUserSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdDesktopsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_desktop.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcdDesktopsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_desktop.default.id}_fake"]`,
		}),
	}

	officeSiteIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdDesktopsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_ecd_desktop.default.id}"]`,
			"office_site_id": `"${alicloud_ecd_desktop.default.office_site_id}"`,
		}),
	}

	policyGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdDesktopsDataSourceName(rand, map[string]string{
			"ids":             `["${alicloud_ecd_desktop.default.id}"]`,
			"policy_group_id": `"${alicloud_ecd_desktop.default.policy_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdDesktopsDataSourceName(rand, map[string]string{
			"ids":             `["${alicloud_ecd_desktop.default.id}"]`,
			"policy_group_id": `"${alicloud_ecd_desktop.default.policy_group_id}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdDesktopsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_desktop.default.desktop_name}"`,
			"status":     `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdDesktopsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_desktop.default.desktop_name}"`,
			"status":     `"Stopped"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdDesktopsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_desktop.default.desktop_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdDesktopsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_desktop.default.desktop_name}_fake"`,
		}),
	}

	deskTopNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdDesktopsDataSourceName(rand, map[string]string{
			"desktop_name": `"${alicloud_ecd_desktop.default.desktop_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdDesktopsDataSourceName(rand, map[string]string{
			"desktop_name": `"${alicloud_ecd_desktop.default.desktop_name}_fake"`,
		}),
	}

	endUserIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdDesktopsDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_ecd_desktop.default.id}"]`,
			"end_user_ids": `["${alicloud_ecd_user.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcdDesktopsDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_ecd_desktop.default.id}"]`,
			"end_user_ids": `["${alicloud_ecd_user.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdDesktopsDataSourceName(rand, map[string]string{
			"ids":             `["${alicloud_ecd_desktop.default.id}"]`,
			"office_site_id":  `"${alicloud_ecd_desktop.default.office_site_id}"`,
			"policy_group_id": `"${alicloud_ecd_desktop.default.policy_group_id}"`,
			"status":          `"Running"`,
			"name_regex":      `"${alicloud_ecd_desktop.default.desktop_name}"`,
			"desktop_name":    `"${alicloud_ecd_desktop.default.desktop_name}"`,
			"end_user_ids":    `["${alicloud_ecd_user.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcdDesktopsDataSourceName(rand, map[string]string{
			"ids":             `["${alicloud_ecd_desktop.default.id}_fake"]`,
			"policy_group_id": `"${alicloud_ecd_desktop.default.policy_group_id}_fake"`,
			"status":          `"Stopped"`,
			"name_regex":      `"${alicloud_ecd_desktop.default.desktop_name}_fake"`,
			"desktop_name":    `"${alicloud_ecd_desktop.default.desktop_name}_fake"`,
			"end_user_ids":    `["${alicloud_ecd_user.default.id}_fake"]`,
		}),
	}
	var existAlicloudEcdDesktopsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"desktops.#":                      "1",
			"desktops.0.id":                   CHECKSET,
			"desktops.0.cpu":                  CHECKSET,
			"desktops.0.create_time":          CHECKSET,
			"desktops.0.desktop_id":           CHECKSET,
			"desktops.0.desktop_name":         fmt.Sprintf("tf-testaccdesktop%d", rand),
			"desktops.0.desktop_type":         CHECKSET,
			"desktops.0.directory_id":         CHECKSET,
			"desktops.0.status":               "Running",
			"desktops.0.expired_time":         CHECKSET,
			"desktops.0.image_id":             CHECKSET,
			"desktops.0.memory":               CHECKSET,
			"desktops.0.network_interface_id": CHECKSET,
			"desktops.0.policy_group_id":      CHECKSET,
			"desktops.0.system_disk_size":     CHECKSET,
			"desktops.0.payment_type":         "PayAsYouGo",
			"desktops.0.end_user_ids.#":       "1",
		}
	}
	var fakeAlicloudEcdDesktopsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"desktops.#": "0",
		}
	}
	var alicloudEcdDesktopsBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecd_desktops.default",
		existMapFunc: existAlicloudEcdDesktopsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcdDesktopsDataSourceNameMapFunc,
	}

	alicloudEcdDesktopsBusesCheckInfo.dataSourceTestCheck(t, rand, idsConf, officeSiteIdConf, policyGroupIdConf, statusConf, nameRegexConf, deskTopNameConf, endUserIdConf, allConf)
}
func testAccCheckAlicloudEcdDesktopsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tf-testaccdesktop%d"
}

data "alicloud_ecd_bundles" "default"{
	bundle_type = "SYSTEM"
}

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = var.name

}

resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = var.name
  clipboard = "readwrite"
  local_drive = "read"
  authorize_access_policy_rules{
    description= var.name
    cidr_ip=     "1.2.3.4/24"
  }
  authorize_security_policy_rules  {
    type=        "inflow"
    policy=      "accept"
    description=  var.name
    port_range= "80/80"
    ip_protocol= "TCP"
    priority=    "1"
    cidr_ip=     "0.0.0.0/0"
  }
}

resource "alicloud_ecd_desktop" "default" {
	office_site_id  = alicloud_ecd_simple_office_site.default.id
	policy_group_id = alicloud_ecd_policy_group.default.id
	bundle_id 		= data.alicloud_ecd_bundles.default.bundles.1.id
	desktop_name 	=  var.name
    end_user_ids = [alicloud_ecd_user.default.id]
}

resource "alicloud_ecd_user" "default" {
	end_user_id = "tf_testaccecduser%d"
	email       = "hello.%d@aaa.com"
	phone       = "158016%d"
	password    = "%d"
}

data "alicloud_ecd_desktops" "default" {	
	%s
}
`, rand, rand, rand, rand, rand, strings.Join(pairs, " \n "))
	return config
}
