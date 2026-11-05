package ecd_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudECDNasFileSystemsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_nas_file_system.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_nas_file_system.default.id}_fake"]`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_nas_file_system.default.nas_file_system_name}"`,
			"status":     `"Running"`,
		}),
		FakeConfig: testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_nas_file_system.default.nas_file_system_name}"`,
			"status":     `"Stopped"`,
		}),
	}
	officeSiteIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_ecd_nas_file_system.default.id}"]`,
			"office_site_id": `"${alicloud_ecd_simple_office_site.default.id}"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_nas_file_system.default.nas_file_system_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_nas_file_system.default.nas_file_system_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_ecd_nas_file_system.default.id}"]`,
			"status":         `"Running"`,
			"office_site_id": `"${alicloud_ecd_simple_office_site.default.id}"`,
			"name_regex":     `"${alicloud_ecd_nas_file_system.default.nas_file_system_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudEcdNasFileSystemsDataSourceName(rand, map[string]string{
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
			"systems.0.metered_size":         acctest.CHECKSET,
			"systems.0.capacity":             acctest.CHECKSET,
			"systems.0.mount_target_status":  acctest.CHECKSET,
			"systems.0.nas_file_system_name": fmt.Sprintf("tf-testacc-ecdnasfilesystem%d", rand),
			"systems.0.office_site_id":       acctest.CHECKSET,
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
	var alicloudEcdNasFileSystemsBusesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ecd_nas_file_systems.default",
		ExistMapFunc: existAlicloudEcdNasFileSystemsBusesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudEcdNasFileSystemsBusesDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
	}
	alicloudEcdNasFileSystemsBusesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, officeSiteIdConf, nameRegexConf, allConf)
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
