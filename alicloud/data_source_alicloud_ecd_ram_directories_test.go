package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEcdRamDirectoriesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.EcdSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdRamDirectoriesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_ram_directory.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcdRamDirectoriesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_ram_directory.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdRamDirectoriesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_ram_directory.default.ram_directory_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdRamDirectoriesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_ram_directory.default.ram_directory_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdRamDirectoriesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecd_ram_directory.default.id}"]`,
			"status": `"REGISTERED"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdRamDirectoriesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecd_ram_directory.default.id}"]`,
			"status": `"REGISTERING"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdRamDirectoriesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_ram_directory.default.id}"]`,
			"name_regex": `"${alicloud_ecd_ram_directory.default.ram_directory_name}"`,
			"status":     `"REGISTERED"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdRamDirectoriesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_ram_directory.default.id}_fake"]`,
			"name_regex": `"${alicloud_ecd_ram_directory.default.ram_directory_name}_fake"`,
			"status":     `"REGISTERING"`,
		}),
	}
	var existAlicloudEcdRamDirectoriesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                     "1",
			"names.#":                                   "1",
			"directories.#":                             "1",
			"directories.0.desktop_access_type":         "INTERNET",
			"directories.0.enable_admin_access":         "true",
			"directories.0.enable_internet_access":      "true",
			"directories.0.ram_directory_name":          CHECKSET,
			"directories.0.vswitch_ids.#":               "1",
			"directories.0.vswitch_ids.0":               CHECKSET,
			"directories.0.status":                      "REGISTERED",
			"directories.0.id":                          CHECKSET,
			"directories.0.ad_connectors.#":             CHECKSET,
			"directories.0.create_time":                 CHECKSET,
			"directories.0.custom_security_group_id":    "",
			"directories.0.desktop_vpc_endpoint":        "",
			"directories.0.directory_type":              "RAM",
			"directories.0.dns_address.#":               "0",
			"directories.0.dns_user_name":               "",
			"directories.0.domain_name":                 "",
			"directories.0.domain_password":             "",
			"directories.0.domain_user_name":            "",
			"directories.0.enable_cross_desktop_access": CHECKSET,
			"directories.0.file_system_ids.#":           CHECKSET,
			"directories.0.logs.#":                      CHECKSET,
			"directories.0.mfa_enabled":                 CHECKSET,
			"directories.0.ram_directory_id":            CHECKSET,
			"directories.0.sso_enabled":                 CHECKSET,
			"directories.0.sub_dns_address.#":           CHECKSET,
			"directories.0.sub_domain_name":             "",
			"directories.0.trust_password":              "",
			"directories.0.vpc_id":                      CHECKSET,
		}
	}
	var fakeAlicloudEcdRamDirectoriesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEcdRamDirectoriesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecd_ram_directories.default",
		existMapFunc: existAlicloudEcdRamDirectoriesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcdRamDirectoriesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEcdRamDirectoriesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudEcdRamDirectoriesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccRamDirectory-%d"
}
data "alicloud_ecd_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_ecd_zones.default.ids.0
}
resource "alicloud_ecd_ram_directory" "default" {
	desktop_access_type = "INTERNET"
	enable_admin_access = "true"
	enable_internet_access = "true"
	ram_directory_name = var.name
	vswitch_ids = [data.alicloud_vswitches.default.ids.0]
}

data "alicloud_ecd_ram_directories" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
