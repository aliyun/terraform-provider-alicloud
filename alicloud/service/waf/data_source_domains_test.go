package waf_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudWafDomainsDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.WAFSupportRegions)
	rand := tfacctest.RandIntRange(1000000, 99999999)
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccAlicloudWafDomainDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_waf_domain.default.instance_id}"`,
			"name_regex":  `"${alicloud_waf_domain.default.domain_name}"`,
		}),
		FakeConfig: testAccAlicloudWafDomainDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_waf_domain.default.instance_id}"`,
			"name_regex":  `"fake"`,
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccAlicloudWafDomainDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_waf_domain.default.instance_id}"`,
			"ids":         `["${alicloud_waf_domain.default.domain_name}"]`,
		}),
		FakeConfig: testAccAlicloudWafDomainDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_waf_domain.default.instance_id}"`,
			"ids":         `["fake"]`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccAlicloudWafDomainDataSourceConfig(rand, map[string]string{
			"instance_id":       `"${alicloud_waf_domain.default.instance_id}"`,
			"ids":               `["${alicloud_waf_domain.default.domain_name}"]`,
			"name_regex":        `"${alicloud_waf_domain.default.domain_name}"`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.groups.0.id}"`,
		}),
		FakeConfig: testAccAlicloudWafDomainDataSourceConfig(rand, map[string]string{
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
			"domains.0.id":                fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", acctest.DefaultRegionToTest, rand),
			"domains.0.domain":            fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", acctest.DefaultRegionToTest, rand),
			"domains.0.domain_name":       fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", acctest.DefaultRegionToTest, rand),
			"domains.0.cluster_type":      "PhysicalCluster",
			"domains.0.cname":             acctest.CHECKSET,
			"domains.0.connection_time":   acctest.CHECKSET,
			"domains.0.http2_port.#":      "1",
			"domains.0.http_port.#":       "1",
			"domains.0.http_to_user_ip":   "Off",
			"domains.0.https_port.#":      "1",
			"domains.0.https_redirect":    "Off",
			"domains.0.is_access_product": "Off",
			"domains.0.load_balancing":    "IpHash",
			"domains.0.log_headers.#":     "1",
			"domains.0.read_time":         acctest.CHECKSET,
			"domains.0.resource_group_id": acctest.CHECKSET,
			"domains.0.source_ips.#":      "1",
			"domains.0.version":           acctest.CHECKSET,
			"domains.0.write_time":        acctest.CHECKSET,
		}
	}

	var fakeWafDomainsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":   "0",
			"ids.#":     "0",
			"domains.#": "0",
		}
	}

	var WafDomainsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_waf_domains.default",
		ExistMapFunc: existWafDomainsMapFunc,
		FakeMapFunc:  fakeWafDomainsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	WafDomainsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
}

func testAccAlicloudWafDomainDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
		variable "name" {
      		default = "tf-testacc%s%d.alicloud-provider.cn"
		}

		data "alicloud_waf_instances" "default" {}

		resource "alicloud_waf_domain" "default" {
			domain_name = "${var.name}"
			instance_id = data.alicloud_waf_instances.default.ids.0
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
`, acctest.DefaultRegionToTest, rand, strings.Join(pairs, "\n  "))
	return config
}
