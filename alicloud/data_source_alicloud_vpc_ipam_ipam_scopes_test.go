package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVpcIpamIpamScopeDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpc_ipam_ipam_scope.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpc_ipam_ipam_scope.default.id}_fake"]`,
		}),
	}

	IpamScopeNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_vpc_ipam_ipam_scope.default.id}"]`,
			"ipam_scope_name": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_vpc_ipam_ipam_scope.default.id}_fake"]`,
			"ipam_scope_name": `"${var.name}_fake"`,
		}),
	}
	ResourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_vpc_ipam_ipam_scope.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_vpc_ipam_ipam_scope.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}
	IpamScopeTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_vpc_ipam_ipam_scope.default.id}"]`,
			"ipam_scope_type": `"private"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_vpc_ipam_ipam_scope.default.id}_fake"]`,
			"ipam_scope_type": `"private_fake"`,
		}),
	}
	IpamIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_vpc_ipam_ipam_scope.default.id}"]`,
			"ipam_id": `"${alicloud_vpc_ipam_ipam.defaultIpam.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_vpc_ipam_ipam_scope.default.id}_fake"]`,
			"ipam_id": `"${alicloud_vpc_ipam_ipam.defaultIpam.id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_vpc_ipam_ipam_scope.default.id}"]`,
			"ipam_scope_name": `"${var.name}"`,

			"ipam_scope_type": `"private"`,

			"ipam_id": `"${alicloud_vpc_ipam_ipam.defaultIpam.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_vpc_ipam_ipam_scope.default.id}_fake"]`,
			"ipam_scope_name": `"${var.name}_fake"`,

			"ipam_scope_type": `"private_fake"`,

			"ipam_id": `"${alicloud_vpc_ipam_ipam.defaultIpam.id}_fake"`,
		}),
	}

	VpcIpamIpamScopeCheckInfo.dataSourceTestCheck(t, rand, idsConf, IpamScopeNameConf, ResourceGroupIdConf, IpamScopeTypeConf, IpamIdConf, allConf)
}

var existVpcIpamIpamScopeMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"scopes.#":                        "1",
		"scopes.0.status":                 CHECKSET,
		"scopes.0.ipam_scope_name":        CHECKSET,
		"scopes.0.resource_group_id":      CHECKSET,
		"scopes.0.create_time":            CHECKSET,
		"scopes.0.ipam_id":                CHECKSET,
		"scopes.0.ipam_scope_type":        CHECKSET,
		"scopes.0.ipam_scope_id":          CHECKSET,
		"scopes.0.ipam_scope_description": CHECKSET,
		"scopes.0.tags.%":                 CHECKSET,
	}
}

var fakeVpcIpamIpamScopeMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"scopes.#": "0",
	}
}

var VpcIpamIpamScopeCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_vpc_ipam_ipam_scopes.default",
	existMapFunc: existVpcIpamIpamScopeMapFunc,
	fakeMapFunc:  fakeVpcIpamIpamScopeMapFunc,
}

func testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccVpcIpamIpamScope%d"
}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
  ipam_name             = var.name
}



resource "alicloud_vpc_ipam_ipam_scope" "default" {
  ipam_scope_name        = var.name
  ipam_id                = alicloud_vpc_ipam_ipam.defaultIpam.id
  ipam_scope_description = "This is a ipam scope."
  ipam_scope_type        = "private"
  tags = {
    "k1": "v1"
  }
}

data "alicloud_vpc_ipam_ipam_scopes" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
