package vpcipam_test

import (
	"fmt"
	"strings"
	"testing"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudVpcIpamIpamDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpc_ipam_ipam.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpc_ipam_ipam.default.id}_fake"]`,
		}),
	}

	IpamNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_vpc_ipam_ipam.default.id}"]`,
			"ipam_name": `"${var.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_vpc_ipam_ipam.default.id}_fake"]`,
			"ipam_name": `"${var.name}_fake"`,
		}),
	}
	ResourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_vpc_ipam_ipam.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_vpc_ipam_ipam.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_vpc_ipam_ipam.default.id}"]`,
			"ipam_name": `"${var.name}"`,

			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_vpc_ipam_ipam.default.id}_fake"]`,
			"ipam_name": `"${var.name}_fake"`,

			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}

	VpcIpamIpamCheckInfo.DataSourceTestCheck(t, rand, idsConf, IpamNameConf, ResourceGroupIdConf, allConf)
}

var existVpcIpamIpamMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ipams.#":        "1",
		"ipams.0.status": acctest.CHECKSET,
		"ipams.0.default_resource_discovery_association_id": acctest.CHECKSET,
		"ipams.0.ipam_name":                            acctest.CHECKSET,
		"ipams.0.resource_group_id":                    acctest.CHECKSET,
		"ipams.0.ipam_id":                              acctest.CHECKSET,
		"ipams.0.create_time":                          acctest.CHECKSET,
		"ipams.0.ipam_description":                     acctest.CHECKSET,
		"ipams.0.default_resource_discovery_id":        acctest.CHECKSET,
		"ipams.0.resource_discovery_association_count": acctest.CHECKSET,
		"ipams.0.region_id":                            acctest.CHECKSET,
		"ipams.0.private_default_scope_id":             acctest.CHECKSET,
		"ipams.0.public_default_scope_id":              acctest.CHECKSET,
		"ipams.0.tags.%":                               acctest.CHECKSET,
	}
}

var fakeVpcIpamIpamMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ipams.#": "0",
	}
}

var VpcIpamIpamCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_vpc_ipam_ipams.default",
	ExistMapFunc: existVpcIpamIpamMapFunc,
	FakeMapFunc:  fakeVpcIpamIpamMapFunc,
}

func testAccCheckAlicloudVpcIpamIpamSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccVpcIpamIpam%d"
}
data "alicloud_resource_manager_resource_groups" "default" {}



resource "alicloud_vpc_ipam_ipam" "default" {
  ipam_description  = "This is my first Ipam."
  ipam_name         = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  tags = {
    ipamKey = "ipamValue"
  }
  operating_region_list = ["cn-hangzhou"]
}

data "alicloud_vpc_ipam_ipams" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
