package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECDNasFileSystemsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_nas_file_system.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_nas_file_system.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_nas_file_system.default.nas_file_system_name}"`,
			"status":     `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_nas_file_system.default.nas_file_system_name}"`,
			"status":     `"Stopped"`,
		}),
	}
	officeSiteIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_ecd_nas_file_system.default.id}"]`,
			"office_site_id": `"${alicloud_ecd_simple_office_site.default.id}"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_nas_file_system.default.nas_file_system_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_nas_file_system.default.nas_file_system_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_ecd_nas_file_system.default.id}"]`,
			"status":         `"Running"`,
			"office_site_id": `"${alicloud_ecd_simple_office_site.default.id}"`,
			"name_regex":     `"${alicloud_ecd_nas_file_system.default.nas_file_system_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_nas_file_system.default.id}_fake"]`,
			"status":     `"Stopped"`,
			"name_regex": `"${alicloud_ecd_nas_file_system.default.nas_file_system_name}_fake"`,
		}),
	}
	var existAlicloudEcdNasFileSystemsBusesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"systems.#":                      "1",
			"systems.0.description":          fmt.Sprintf("tf-testacc-ecdnasfilesystem%d", rand),
			"systems.0.file_system_type":     "standard",
			"systems.0.metered_size":         CHECKSET,
			"systems.0.capacity":             CHECKSET,
			"systems.0.mount_target_status":  CHECKSET,
			"systems.0.nas_file_system_name": fmt.Sprintf("tf-testacc-ecdnasfilesystem%d", rand),
			"systems.0.office_site_id":       CHECKSET,
			"systems.0.status":               "Running",
			"systems.0.storage_type":         "Capacity",
			"systems.0.support_acl":          "false",
		}
	}
	var fakeAlicloudEcdNasFileSystemsBusesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"systems.#": "0",
		}
	}
	var alicloudEcdNasFileSystemsBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecd_nas_file_systems.default",
		existMapFunc: existAlicloudEcdNasFileSystemsBusesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcdNasFileSystemsBusesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
	}
	alicloudEcdNasFileSystemsBusesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, officeSiteIdConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tf-testacc-ecdnasfilesystem%d"
}

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = var.name
  enable_internet_access = false
}

resource "alicloud_ecd_nas_file_system" "default" {
  description = var.name
  office_site_id = alicloud_ecd_simple_office_site.default.id
  nas_file_system_name =  var.name
}

data "alicloud_ecd_nas_file_systems" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
