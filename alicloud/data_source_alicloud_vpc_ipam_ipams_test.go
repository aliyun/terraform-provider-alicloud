package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVpcIpamIpamDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpamIpamSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpc_ipam_ipam.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpamIpamSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpc_ipam_ipam.default.id}_fake"]`,
		}),
	}

	IpamNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpamIpamSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_vpc_ipam_ipam.default.id}"]`,
			"ipam_name": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpamIpamSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_vpc_ipam_ipam.default.id}_fake"]`,
			"ipam_name": `"${var.name}_fake"`,
		}),
	}
	ResourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpamIpamSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_vpc_ipam_ipam.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpamIpamSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_vpc_ipam_ipam.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpamIpamSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_vpc_ipam_ipam.default.id}"]`,
			"ipam_name": `"${var.name}"`,

			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpamIpamSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_vpc_ipam_ipam.default.id}_fake"]`,
			"ipam_name": `"${var.name}_fake"`,

			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}

	VpcIpamIpamCheckInfo.dataSourceTestCheck(t, rand, idsConf, IpamNameConf, ResourceGroupIdConf, allConf)
}

var existVpcIpamIpamMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ipams.#":        "1",
		"ipams.0.status": CHECKSET,
		"ipams.0.default_resource_discovery_association_id": CHECKSET,
		"ipams.0.ipam_name":                            CHECKSET,
		"ipams.0.resource_group_id":                    CHECKSET,
		"ipams.0.ipam_id":                              CHECKSET,
		"ipams.0.create_time":                          CHECKSET,
		"ipams.0.ipam_description":                     CHECKSET,
		"ipams.0.default_resource_discovery_id":        CHECKSET,
		"ipams.0.resource_discovery_association_count": CHECKSET,
		"ipams.0.region_id":                            CHECKSET,
		"ipams.0.private_default_scope_id":             CHECKSET,
		"ipams.0.public_default_scope_id":              CHECKSET,
		"ipams.0.tags.%":                               CHECKSET,
	}
}

var fakeVpcIpamIpamMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ipams.#": "0",
	}
}

var VpcIpamIpamCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_vpc_ipam_ipams.default",
	existMapFunc: existVpcIpamIpamMapFunc,
	fakeMapFunc:  fakeVpcIpamIpamMapFunc,
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
