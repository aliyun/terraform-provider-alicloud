package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAliCloudEcdDesktopGroupsDataSource_basic0(t *testing.T) {
	rand := 10000 + acctest.RandIntRange(0, 89999)
	resourceId := "data.alicloud_ecd_desktop_groups.default"
	name := fmt.Sprintf("tf-testaccdesktopgroupds%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEcdDesktopGroupsConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ecd_desktop_group.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecd_desktop_group.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ecd_desktop_group.default.desktop_group_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "tf-fake-nonexist",
		}),
	}

	officeSiteIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"office_site_id": "${alicloud_ecd_simple_office_site.default.id}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"office_site_id": "cn-shanghai+dir-0000000000",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                []string{"${alicloud_ecd_desktop_group.default.id}"},
			"name_regex":         "${alicloud_ecd_desktop_group.default.desktop_group_name}",
			"desktop_group_id":   "${alicloud_ecd_desktop_group.default.id}",
			"desktop_group_name": "${alicloud_ecd_desktop_group.default.desktop_group_name}",
			"office_site_id":     "${alicloud_ecd_simple_office_site.default.id}",
			"period_unit":        "Month",
			"enable_details":     "true",
			"output_file":        "desktop_groups_output.json",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alicloud_ecd_desktop_group.default.id}_fake"},
			"name_regex":       "tf-fake-nonexist",
			"desktop_group_id": "dg-0000000000fake",
			"enable_details":   "true",
		}),
	}

	var existAliCloudEcdDesktopGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"groups.#":                      "1",
			"groups.0.id":                   CHECKSET,
			"groups.0.desktop_group_id":     CHECKSET,
			"groups.0.desktop_group_name":   CHECKSET,
			"groups.0.bundle_id":            CHECKSET,
			"groups.0.comments":             CHECKSET,
			"groups.0.office_site_id":       CHECKSET,
			"groups.0.office_site_name":     CHECKSET,
			"groups.0.office_site_type":     CHECKSET,
			"groups.0.policy_group_id":      CHECKSET,
			"groups.0.policy_group_name":    CHECKSET,
			"groups.0.pay_type":             CHECKSET,
			"groups.0.create_time":          CHECKSET,
			"groups.0.cpu":                  CHECKSET,
			"groups.0.memory":               CHECKSET,
			"groups.0.system_disk_category": CHECKSET,
			"groups.0.system_disk_size":     CHECKSET,
			"groups.0.min_desktops_count":   CHECKSET,
			"groups.0.max_desktops_count":   CHECKSET,
			"groups.0.allow_auto_setup":     CHECKSET,
			"groups.0.allow_buffer_count":   CHECKSET,
			"groups.0.keep_duration":        CHECKSET,
			"groups.0.end_user_count":       CHECKSET,
			"groups.0.end_user_ids.#":       "1",
			"groups.0.end_user_ids.0":       CHECKSET,
			"groups.0.res_type":             CHECKSET,
			"groups.0.creator":              CHECKSET,
		}
	}

	var fakeAliCloudEcdDesktopGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"groups.#": "0",
		}
	}

	var aliCloudEcdDesktopGroupsInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecd_desktop_groups.default",
		existMapFunc: existAliCloudEcdDesktopGroupsMapFunc,
		fakeMapFunc:  fakeAliCloudEcdDesktopGroupsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.EcdUserSupportRegions)
	}

	aliCloudEcdDesktopGroupsInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, officeSiteIdConf, allConf)
}

func dataSourceEcdDesktopGroupsConfig(name string) string {
	rand := 10000 + acctest.RandIntRange(0, 89999)
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_ecd_bundles" "default" {
  bundle_type = "SYSTEM"
}

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block          = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = var.name
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
  end_user_id = "tf_testacc-dgds-u1-%d"
  email       = "hello.dgds.%d@aaa.com"
  phone       = "158016%d"
  password    = "%d"
}

resource "alicloud_ecd_desktop_group" "default" {
  office_site_id     = alicloud_ecd_simple_office_site.default.id
  policy_group_id    = alicloud_ecd_policy_group.default.id
  bundle_id          = data.alicloud_ecd_bundles.default.bundles.0.id
  end_user_ids       = [alicloud_ecd_user.default.id]
  desktop_group_name = var.name
  comments           = var.name
  allow_buffer_count = 0
}
`, name, rand, rand, rand, rand)
}
