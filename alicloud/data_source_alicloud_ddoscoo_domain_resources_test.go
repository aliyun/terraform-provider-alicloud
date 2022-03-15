package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDdoscooDomainResourcesDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%d.qq.com", rand)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDdoscooDomainResourcesDataSourceName(name, map[string]string{
			"ids": `["${alicloud_ddoscoo_domain_resource.default.domain}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDdoscooDomainResourcesDataSourceName(name, map[string]string{
			"ids": `["${alicloud_ddoscoo_domain_resource.default.domain}_fake"]`,
		}),
	}
	instanceIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDdoscooDomainResourcesDataSourceName(name, map[string]string{
			"ids":          `["${alicloud_ddoscoo_domain_resource.default.domain}"]`,
			"instance_ids": `["${data.alicloud_ddoscoo_instances.default.ids.0}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDdoscooDomainResourcesDataSourceName(name, map[string]string{
			"ids":          `["${alicloud_ddoscoo_domain_resource.default.domain}"]`,
			"instance_ids": `["${data.alicloud_ddoscoo_instances.default.ids.1}"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDdoscooDomainResourcesDataSourceName(name, map[string]string{
			"ids":          `["${alicloud_ddoscoo_domain_resource.default.domain}"]`,
			"instance_ids": `["${data.alicloud_ddoscoo_instances.default.ids.0}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDdoscooDomainResourcesDataSourceName(name, map[string]string{
			"ids":          `["${alicloud_ddoscoo_domain_resource.default.domain}_fake"]`,
			"instance_ids": `["${data.alicloud_ddoscoo_instances.default.ids.1}"]`,
		}),
	}
	var existAlicloudDdoscooDomainResourcesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                   "1",
			"resources.#":                             "1",
			"resources.0.domain":                      name,
			"resources.0.instance_ids.#":              `1`,
			"resources.0.proxy_types.#":               `1`,
			"resources.0.proxy_types.0.proxy_ports.#": `1`,
			"resources.0.proxy_types.0.proxy_type":    `https`,
			"resources.0.real_servers.#":              `1`,
			"resources.0.rs_type":                     `0`,
		}
	}
	var fakeAlicloudDdoscooDomainResourcesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudDdoscooDomainResourcesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ddoscoo_domain_resources.default",
		existMapFunc: existAlicloudDdoscooDomainResourcesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudDdoscooDomainResourcesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.DdoscooSupportedRegions)
	}
	alicloudDdoscooDomainResourcesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, instanceIdsConf, allConf)
}
func testAccCheckAlicloudDdoscooDomainResourcesDataSourceName(name string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_ddoscoo_instances" "default" {}

resource "alicloud_ddoscoo_domain_resource" "default" {
	domain = "%s"
	instance_ids = [data.alicloud_ddoscoo_instances.default.ids.0]
  	proxy_types {   
    	proxy_ports = [443]
   	 	proxy_type = "https"
  	}
	real_servers = ["177.167.32.11"]
	rs_type = 0
}

data "alicloud_ddoscoo_domain_resources" "default" {	
	%s
}
`, name, strings.Join(pairs, " \n "))
	return config
}
