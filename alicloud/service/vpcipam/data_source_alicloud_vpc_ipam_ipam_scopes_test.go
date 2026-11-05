package vpcipam_test

import (
	"fmt"
	"strings"
	"testing"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudVpcIpamIpamScopeDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpc_ipam_ipam_scope.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpc_ipam_ipam_scope.default.id}_fake"]`,
		}),
	}

	IpamScopeNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_vpc_ipam_ipam_scope.default.id}"]`,
			"ipam_scope_name": `"${var.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_vpc_ipam_ipam_scope.default.id}_fake"]`,
			"ipam_scope_name": `"${var.name}_fake"`,
		}),
	}
	ResourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_vpc_ipam_ipam_scope.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_vpc_ipam_ipam_scope.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}
	IpamScopeTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_vpc_ipam_ipam_scope.default.id}"]`,
			"ipam_scope_type": `"private"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_vpc_ipam_ipam_scope.default.id}_fake"]`,
			"ipam_scope_type": `"private_fake"`,
		}),
	}
	IpamIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_vpc_ipam_ipam_scope.default.id}"]`,
			"ipam_id": `"${alicloud_vpc_ipam_ipam.defaultIpam.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_vpc_ipam_ipam_scope.default.id}_fake"]`,
			"ipam_id": `"${alicloud_vpc_ipam_ipam.defaultIpam.id}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_vpc_ipam_ipam_scope.default.id}"]`,
			"ipam_scope_name": `"${var.name}"`,

			"ipam_scope_type": `"private"`,

			"ipam_id": `"${alicloud_vpc_ipam_ipam.defaultIpam.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamScopeSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_vpc_ipam_ipam_scope.default.id}_fake"]`,
			"ipam_scope_name": `"${var.name}_fake"`,

			"ipam_scope_type": `"private_fake"`,

			"ipam_id": `"${alicloud_vpc_ipam_ipam.defaultIpam.id}_fake"`,
		}),
	}

	VpcIpamIpamScopeCheckInfo.DataSourceTestCheck(t, rand, idsConf, IpamScopeNameConf, ResourceGroupIdConf, IpamScopeTypeConf, IpamIdConf, allConf)
}

var existVpcIpamIpamScopeMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"scopes.#":                        "1",
		"scopes.0.status":                 acctest.CHECKSET,
		"scopes.0.ipam_scope_name":        acctest.CHECKSET,
		"scopes.0.resource_group_id":      acctest.CHECKSET,
		"scopes.0.create_time":            acctest.CHECKSET,
		"scopes.0.ipam_id":                acctest.CHECKSET,
		"scopes.0.ipam_scope_type":        acctest.CHECKSET,
		"scopes.0.ipam_scope_id":          acctest.CHECKSET,
		"scopes.0.ipam_scope_description": acctest.CHECKSET,
		"scopes.0.tags.%":                 acctest.CHECKSET,
	}
}

var fakeVpcIpamIpamScopeMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"scopes.#": "0",
	}
}

var VpcIpamIpamScopeCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_vpc_ipam_ipam_scopes.default",
	ExistMapFunc: existVpcIpamIpamScopeMapFunc,
	FakeMapFunc:  fakeVpcIpamIpamScopeMapFunc,
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
