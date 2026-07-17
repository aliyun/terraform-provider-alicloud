// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlicloudApigDomainDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigDomainSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_domain.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudApigDomainSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_domain.default.id}_fake"]`,
		}),
	}

	ResourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigDomainSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_domain.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		fakeConfig: testAccCheckAlicloudApigDomainSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_domain.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigDomainSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_domain.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		fakeConfig: testAccCheckAlicloudApigDomainSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_domain.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}

	ApigDomainCheckInfo.dataSourceTestCheck(t, rand, idsConf, ResourceGroupIdConf, allConf)
}

var existApigDomainMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"domains.#":                   "1",
		"domains.0.domain_scope":      CHECKSET,
		"domains.0.resource_group_id": CHECKSET,
		"domains.0.domain_name":       CHECKSET,
		"domains.0.domain_id":         CHECKSET,
		"domains.0.protocol":          CHECKSET,
	}
}

var fakeApigDomainMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"domains.#": "0",
	}
}

var ApigDomainCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_apig_domains.default",
	existMapFunc: existApigDomainMapFunc,
	fakeMapFunc:  fakeApigDomainMapFunc,
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
