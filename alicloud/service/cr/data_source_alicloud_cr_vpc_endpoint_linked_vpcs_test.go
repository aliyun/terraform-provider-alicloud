package cr_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCrVpcEndpointLinkedVpcsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCrVpcEndpointLinkedVpcsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cr_vpc_endpoint_linked_vpc.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCrVpcEndpointLinkedVpcsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cr_vpc_endpoint_linked_vpc.default.id}_fake"]`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCrVpcEndpointLinkedVpcsDataSourceName(rand, map[string]string{
			"status": `"RUNNING"`,
		}),
		FakeConfig: testAccCheckAlicloudCrVpcEndpointLinkedVpcsDataSourceName(rand, map[string]string{
			"status": `"CREATING"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCrVpcEndpointLinkedVpcsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cr_vpc_endpoint_linked_vpc.default.id}"]`,
			"status": `"RUNNING"`,
		}),
		FakeConfig: testAccCheckAlicloudCrVpcEndpointLinkedVpcsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cr_vpc_endpoint_linked_vpc.default.id}_fake"]`,
			"status": `"CREATING"`,
		}),
	}
	var existAlicloudCrVpcEndpointLinkedVpcsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                     "1",
			"vpc_endpoint_linked_vpcs.#":                "1",
			"vpc_endpoint_linked_vpcs.0.id":             acctest.CHECKSET,
			"vpc_endpoint_linked_vpcs.0.instance_id":    acctest.CHECKSET,
			"vpc_endpoint_linked_vpcs.0.vpc_id":         acctest.CHECKSET,
			"vpc_endpoint_linked_vpcs.0.vswitch_id":     acctest.CHECKSET,
			"vpc_endpoint_linked_vpcs.0.module_name":    "Registry",
			"vpc_endpoint_linked_vpcs.0.ip":             acctest.CHECKSET,
			"vpc_endpoint_linked_vpcs.0.default_access": acctest.CHECKSET,
			"vpc_endpoint_linked_vpcs.0.status":         "RUNNING",
		}
	}
	var fakeAlicloudCrVpcEndpointLinkedVpcsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "0",
			"vpc_endpoint_linked_vpcs.#": "0",
		}
	}
	var alicloudCrVpcEndpointLinkedVpcsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cr_vpc_endpoint_linked_vpcs.default",
		ExistMapFunc: existAlicloudCrVpcEndpointLinkedVpcsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCrVpcEndpointLinkedVpcsDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudCrVpcEndpointLinkedVpcsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, allConf)
}

func testAccCheckAlicloudCrVpcEndpointLinkedVpcsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-CrVpcEndpointLinkedVpc-%d"
	}

	data "alicloud_cr_ee_instances" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_cr_vpc_endpoint_linked_vpc" "default" {
  		instance_id                      = data.alicloud_cr_ee_instances.default.ids.0
  		vpc_id                           = data.alicloud_vpcs.default.ids.0
  		vswitch_id                       = data.alicloud_vswitches.default.ids.0
  		module_name                      = "Registry"
  		enable_create_dns_record_in_pvzt = true
	}

	data "alicloud_cr_vpc_endpoint_linked_vpcs" "default" {
  		instance_id = alicloud_cr_vpc_endpoint_linked_vpc.default.instance_id
  		module_name = alicloud_cr_vpc_endpoint_linked_vpc.default.module_name
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
