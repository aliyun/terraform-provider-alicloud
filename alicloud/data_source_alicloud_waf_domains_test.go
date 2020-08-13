package alicloud

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"fmt"
)

func TestAccAlicloudWafDomainsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 99999999)
	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccAlicloudWafDomainDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_waf_domain.default.instance_id}"`,
		}),
		//fakeConfig: testAccAlicloudWafDomainDataSourceConfig(rand, map[string]string{
		//	"instance_id": `"${alicloud_waf_domain.default.instance_id}_fake"`,
		//}),
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	WafDomainsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, instanceIdConf)
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
		resource "alicloud_waf_instance" "default" {
			big_screen = "0"
			exclusive_ip_package = "1"
			ext_bandwidth = "50"
			ext_domain_package = "1"
			package_code = "version_3"
			prefessional_service = "false"
			subscription_type = "Subscription"
			period = "1"
			waf_log = "false"
			log_storage = "3"
			log_time = "180"
		}
		resource "alicloud_waf_domain" "default" {
			domain = "${var.name}"
			instance_id = "${alicloud_waf_instance.default.id}"
			is_access_product = "Off"
			source_ips = ["1.1.1.1"]
			cluster_type = "PhysicalCluster"
			http2_port = ["443"]
			http_port = ["80"]
			https_port = ["443"]
			http_to_user_ip = "Off"
			https_redirect = "Off"
			load_balancing = "IpHash"
		}
		data "alicloud_waf_domains" "default" {
			%s
		}
`, defaultRegionToTest, rand, strings.Join(pairs, "\n  "))
	return config
}

var existWafDomainsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"domains.#":        "1",
		"domains.0.domain": CHECKSET,
	}
}

var fakeWafDomainsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"names.#":     "0",
		"instances.#": "0",
	}
}

var existWafDomainsMultiMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"names.#":                 "5",
		"instances.#":             "5",
		"instances.0.name":        fmt.Sprintf("tf-testacc%s%d.wafqa3.com", defaultRegionToTest, rand),
		"instances.0.description": "tf-testAccCenConfigDescription",
	}
}

var WafDomainsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_waf_domains.default",
	existMapFunc: existWafDomainsMapFunc,
	fakeMapFunc:  fakeWafDomainsMapFunc,
}
var WafDomainsCheckInfoMulti = dataSourceAttr{
	resourceId:   "data.alicloud_waf_domains.default",
	existMapFunc: existWafDomainsMultiMapFunc,
	fakeMapFunc:  fakeWafDomainsMapFunc,
}
