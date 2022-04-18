package alicloud

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"fmt"
)

func TestAccAlicloudWAFDomainsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 99999999)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccAlicloudWafDomainDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_waf_domain.default.instance_id}"`,
			"name_regex":  `"${alicloud_waf_domain.default.domain_name}"`,
		}),
		fakeConfig: testAccAlicloudWafDomainDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_waf_domain.default.instance_id}"`,
			"name_regex":  `"fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccAlicloudWafDomainDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_waf_domain.default.instance_id}"`,
			"ids":         `["${alicloud_waf_domain.default.domain_name}"]`,
		}),
		fakeConfig: testAccAlicloudWafDomainDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_waf_domain.default.instance_id}"`,
			"ids":         `["fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccAlicloudWafDomainDataSourceConfig(rand, map[string]string{
			"instance_id":       `"${alicloud_waf_domain.default.instance_id}"`,
			"ids":               `["${alicloud_waf_domain.default.domain_name}"]`,
			"name_regex":        `"${alicloud_waf_domain.default.domain_name}"`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.groups.0.id}"`,
		}),
		fakeConfig: testAccAlicloudWafDomainDataSourceConfig(rand, map[string]string{
			"instance_id":       `"${alicloud_waf_domain.default.instance_id}"`,
			"ids":               `["${alicloud_waf_domain.default.domain_name}"]`,
			"name_regex":        `"fake"`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.groups.0.id}"`,
		}),
	}
	var existWafDomainsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                     "1",
			"ids.#":                       "1",
			"domains.#":                   "1",
			"domains.0.id":                fmt.Sprintf("tf-testacc%s%d.wafqa3.com", defaultRegionToTest, rand),
			"domains.0.domain":            fmt.Sprintf("tf-testacc%s%d.wafqa3.com", defaultRegionToTest, rand),
			"domains.0.domain_name":       fmt.Sprintf("tf-testacc%s%d.wafqa3.com", defaultRegionToTest, rand),
			"domains.0.cluster_type":      "PhysicalCluster",
			"domains.0.cname":             CHECKSET,
			"domains.0.connection_time":   CHECKSET,
			"domains.0.http2_port.#":      "1",
			"domains.0.http_port.#":       "1",
			"domains.0.http_to_user_ip":   "Off",
			"domains.0.https_port.#":      "1",
			"domains.0.https_redirect":    "Off",
			"domains.0.is_access_product": "Off",
			"domains.0.load_balancing":    "IpHash",
			"domains.0.log_headers.#":     "1",
			"domains.0.read_time":         CHECKSET,
			"domains.0.resource_group_id": CHECKSET,
			"domains.0.source_ips.#":      "1",
			"domains.0.version":           CHECKSET,
			"domains.0.write_time":        CHECKSET,
		}
	}

	var fakeWafDomainsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":   "0",
			"ids.#":     "0",
			"domains.#": "0",
		}
	}

	var WafDomainsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_waf_domains.default",
		existMapFunc: existWafDomainsMapFunc,
		fakeMapFunc:  fakeWafDomainsMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithEnvVariable(t, "ALICLOUD_WAF_INSTANCE_ID")
	}
	WafDomainsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
}

func testAccAlicloudWafDomainDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
		variable "name" {
      		default = "tf-testacc%s%d.wafqa3.com"
		}

		resource "alicloud_waf_domain" "default" {
			domain_name = "${var.name}"
			instance_id = "%s"
			is_access_product = "Off"
			source_ips = ["1.1.1.1"]
			cluster_type = "PhysicalCluster"
			http2_port = ["443"]
			http_port = ["80"]
			https_port = ["443"]
			http_to_user_ip = "Off"
			https_redirect = "Off"
			load_balancing = "IpHash"
  			log_headers {
    			key = "tf"
    			value = "test"
  			}
		}
		
		data "alicloud_resource_manager_resource_groups" "default" {
			name_regex="^default$"
		}

		data "alicloud_waf_domains" "default" {
			%s
			  enable_details = true
		}
`, defaultRegionToTest, rand, os.Getenv("ALICLOUD_WAF_INSTANCE_ID"), strings.Join(pairs, "\n  "))
	return config
}
