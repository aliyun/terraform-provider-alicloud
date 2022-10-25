package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEcdAdConnectorDirectoriesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.EcdSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdAdConnectorDirectoriesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_ad_connector_directory.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcdAdConnectorDirectoriesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_ad_connector_directory.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdAdConnectorDirectoriesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_ad_connector_directory.default.directory_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdAdConnectorDirectoriesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_ad_connector_directory.default.directory_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdAdConnectorDirectoriesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_ad_connector_directory.default.id}"]`,
			"name_regex": `"${alicloud_ecd_ad_connector_directory.default.directory_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdAdConnectorDirectoriesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_ad_connector_directory.default.id}_fake"]`,
			"name_regex": `"${alicloud_ecd_ad_connector_directory.default.directory_name}_fake"`,
		}),
	}
	var existAlicloudEcdAdConnectorDirectoriesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                                 "1",
			"ids.#":                                   "1",
			"directories.#":                           "1",
			"directories.0.directory_name":            fmt.Sprintf("tf-testAccAdConnectorDirectory-%d", rand),
			"directories.0.directory_type":            "AD_CONNECTOR",
			"directories.0.dns_address.#":             "1",
			"directories.0.domain_name":               "corp.example.com",
			"directories.0.domain_user_name":          "sAMAccountName",
			"directories.0.enable_admin_access":       "false",
			"directories.0.mfa_enabled":               "false",
			"directories.0.sub_dns_address.#":         "1",
			"directories.0.sub_domain_name":           "child.example.com",
			"directories.0.vswitch_ids.#":             "1",
			"directories.0.ad_connectors.#":           CHECKSET,
			"directories.0.status":                    CHECKSET,
			"directories.0.id":                        CHECKSET,
			"directories.0.ad_connector_directory_id": CHECKSET,
			"directories.0.create_time":               CHECKSET,
			"directories.0.custom_security_group_id":  "",
			"directories.0.dns_user_name":             "",
			"directories.0.trust_password":            "",
			"directories.0.vpc_id":                    CHECKSET,
		}
	}
	var fakeAlicloudEcdAdConnectorDirectoriesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEcdAdConnectorDirectoriesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecd_ad_connector_directories.default",
		existMapFunc: existAlicloudEcdAdConnectorDirectoriesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcdAdConnectorDirectoriesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEcdAdConnectorDirectoriesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudEcdAdConnectorDirectoriesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccAdConnectorDirectory-%d"
}
data "alicloud_ecd_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_ecd_zones.default.ids.0
}
resource "alicloud_ecd_ad_connector_directory" "default" {
	directory_name = var.name
	desktop_access_type = "INTERNET"
	dns_address =  ["127.0.0.2"]
	domain_name = "corp.example.com"
	domain_password = "YourPassword1234"
	domain_user_name = "sAMAccountName"
	enable_admin_access = false
	mfa_enabled = false
	specification = 1
	sub_domain_dns_address = ["127.0.0.3"]
	sub_domain_name = "child.example.com"
	vswitch_ids = [data.alicloud_vswitches.default.ids.0]
}

data "alicloud_ecd_ad_connector_directories" "default" {	
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
