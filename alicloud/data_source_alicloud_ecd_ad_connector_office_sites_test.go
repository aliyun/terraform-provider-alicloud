package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEcdAdConnectorOfficeSitesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.EcdSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdAdConnectorOfficeSitesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_ad_connector_office_site.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcdAdConnectorOfficeSitesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_ad_connector_office_site.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdAdConnectorOfficeSitesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_ad_connector_office_site.default.ad_connector_office_site_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdAdConnectorOfficeSitesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_ad_connector_office_site.default.ad_connector_office_site_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdAdConnectorOfficeSitesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecd_ad_connector_office_site.default.id}"]`,
			"status": `"${alicloud_ecd_ad_connector_office_site.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdAdConnectorOfficeSitesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecd_ad_connector_office_site.default.id}"]`,
			"status": `"REGISTERING"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdAdConnectorOfficeSitesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_ad_connector_office_site.default.id}"]`,
			"name_regex": `"${alicloud_ecd_ad_connector_office_site.default.ad_connector_office_site_name}"`,
			"status":     `"${alicloud_ecd_ad_connector_office_site.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdAdConnectorOfficeSitesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_ad_connector_office_site.default.id}_fake"]`,
			"name_regex": `"${alicloud_ecd_ad_connector_office_site.default.ad_connector_office_site_name}_fake"`,
			"status":     `"REGISTERING"`,
		}),
	}
	var existAlicloudEcdAdConnectorOfficeSitesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 "1",
			"names.#":                               "1",
			"sites.#":                               "1",
			"sites.0.ad_connector_office_site_name": fmt.Sprintf("tf-testAccAdConnectorOfficeSite-%d", rand),
			"sites.0.bandwidth":                     "100",
			"sites.0.cen_id":                        CHECKSET,
			"sites.0.id":                            CHECKSET,
			"sites.0.office_site_id":                CHECKSET,
			"sites.0.office_site_type":              CHECKSET,
			"sites.0.create_time":                   CHECKSET,
			"sites.0.custom_security_group_id":      "",
			"sites.0.desktop_vpc_endpoint":          "",
			"sites.0.dns_user_name":                 "",
			"sites.0.file_system_ids.#":             "0",
			"sites.0.ad_connectors.#":               CHECKSET,
			"sites.0.logs.#":                        CHECKSET,
			"sites.0.vswitch_ids.#":                 CHECKSET,
			"sites.0.network_package_id":            CHECKSET,
			"sites.0.status":                        CHECKSET,
			"sites.0.trust_password":                "",
			"sites.0.vpc_id":                        CHECKSET,
			"sites.0.cidr_block":                    "10.0.0.0/12",
			"sites.0.desktop_access_type":           "INTERNET",
			"sites.0.dns_address.#":                 "1",
			"sites.0.domain_name":                   "example1234.com",
			"sites.0.domain_user_name":              "Administrator",
			"sites.0.enable_admin_access":           "true",
			"sites.0.enable_cross_desktop_access":   "false",
			"sites.0.enable_internet_access":        "true",
			"sites.0.mfa_enabled":                   "false",
			"sites.0.sso_enabled":                   "false",
			"sites.0.sub_domain_dns_address.#":      "1",
			"sites.0.sub_domain_name":               "child.example1234.com",
		}
	}
	var fakeAlicloudEcdAdConnectorOfficeSitesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEcdAdConnectorOfficeSitesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecd_ad_connector_office_sites.default",
		existMapFunc: existAlicloudEcdAdConnectorOfficeSitesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcdAdConnectorOfficeSitesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEcdAdConnectorOfficeSitesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudEcdAdConnectorOfficeSitesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccAdConnectorOfficeSite-%d"
}
resource "alicloud_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	protection_level = "REDUCED"
}

resource "alicloud_ecd_ad_connector_office_site" "default" {
	ad_connector_office_site_name = var.name
	bandwidth = "100"
	cen_id = alicloud_cen_instance.default.id
	cidr_block = "10.0.0.0/12"
	desktop_access_type = "INTERNET"
	dns_address = ["127.0.0.2"]
	domain_name = "example1234.com"
	domain_password = "YourPassword1234"
	domain_user_name =  "Administrator"
	enable_admin_access = "true"
	enable_internet_access = "true"
	mfa_enabled = "false"
	sub_domain_dns_address = ["127.0.0.3"]
	sub_domain_name =  "child.example1234.com"
}

data "alicloud_ecd_ad_connector_office_sites" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
