// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package apig_test

import (
	"fmt"
	"strings"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudApigDomainDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudApigDomainSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_domain.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudApigDomainSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_domain.default.id}_fake"]`,
		}),
	}

	ResourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudApigDomainSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_domain.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		FakeConfig: testAccCheckAlicloudApigDomainSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_domain.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudApigDomainSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_domain.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		FakeConfig: testAccCheckAlicloudApigDomainSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_domain.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}

	ApigDomainCheckInfo.DataSourceTestCheck(t, rand, idsConf, ResourceGroupIdConf, allConf)
}

var existApigDomainMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"domains.#":                   "1",
		"domains.0.domain_scope":      acctest.CHECKSET,
		"domains.0.resource_group_id": acctest.CHECKSET,
		"domains.0.domain_name":       acctest.CHECKSET,
		"domains.0.domain_id":         acctest.CHECKSET,
		"domains.0.protocol":          acctest.CHECKSET,
	}
}

var fakeApigDomainMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"domains.#": "0",
	}
}

var ApigDomainCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_apig_domains.default",
	ExistMapFunc: existApigDomainMapFunc,
	FakeMapFunc:  fakeApigDomainMapFunc,
}

func testAccCheckAlicloudApigDomainSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf.apig%d.com"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_apig_domain" "default" {
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  domain_name       = var.name
  gateway_type      = "API"
  protocol          = "HTTP"
}

data "alicloud_apig_domains" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
