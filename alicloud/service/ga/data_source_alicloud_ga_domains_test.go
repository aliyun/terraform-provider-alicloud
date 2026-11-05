package ga_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudGaDomainDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGaDomainSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ga_domain.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudGaDomainSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ga_domain.default.id}_fake"]`,
		}),
	}

	AcceleratorIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGaDomainSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ga_domain.default.id}"]`,
			"accelerator_id": `"${alicloud_ga_domain.default.accelerator_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudGaDomainSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ga_domain.default.id}_fake"]`,
			"accelerator_id": `"${alicloud_ga_domain.default.accelerator_id}_fake"`,
		}),
	}
	DomainConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGaDomainSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_ga_domain.default.id}"]`,
			"domain": `"${alicloud_ga_domain.default.domain}"`,
		}),
		FakeConfig: testAccCheckAlicloudGaDomainSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_ga_domain.default.id}_fake"]`,
			"domain": `"${alicloud_ga_domain.default.domain}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGaDomainSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_ga_domain.default.id}"]`,
			"status": `"active"`,
		}),
		FakeConfig: testAccCheckAlicloudGaDomainSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_ga_domain.default.id}_fake"]`,
			"status": `"illegal"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGaDomainSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ga_domain.default.id}"]`,
			"accelerator_id": `"${alicloud_ga_domain.default.accelerator_id}"`,
			"domain":         `"${alicloud_ga_domain.default.domain}"`,
			"status":         `"active"`,
		}),
		FakeConfig: testAccCheckAlicloudGaDomainSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ga_domain.default.id}_fake"]`,
			"accelerator_id": `"${alicloud_ga_domain.default.accelerator_id}_fake"`,
			"domain":         `"${alicloud_ga_domain.default.domain}_fake"`,
			"status":         `"illegal"`,
		}),
	}

	GaDomainCheckInfo.DataSourceTestCheck(t, rand, idsConf, AcceleratorIdConf, DomainConf, statusConf, allConf)
}

var existGaDomainMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                    "1",
		"domains.#":                "1",
		"domains.0.id":             acctest.CHECKSET,
		"domains.0.accelerator_id": acctest.CHECKSET,
		"domains.0.domain":         acctest.CHECKSET,
		"domains.0.status":         acctest.CHECKSET,
	}
}

var fakeGaDomainMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":     "0",
		"domains.#": "0",
	}
}

var GaDomainCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_ga_domains.default",
	ExistMapFunc: existGaDomainMapFunc,
	FakeMapFunc:  fakeGaDomainMapFunc,
}

func testAccCheckAlicloudGaDomainSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccGaDomain%d"
}

data "alicloud_ga_accelerators" "default" {
  status = "active"
}

resource "alicloud_ga_accelerator" "default" {
  count           = length(data.alicloud_ga_accelerators.default.accelerators) > 0 ? 0 : 1
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}

locals {
  accelerator_id = length(data.alicloud_ga_accelerators.default.accelerators) > 0 ? data.alicloud_ga_accelerators.default.accelerators.0.id : alicloud_ga_accelerator.default.0.id
}

resource "alicloud_ga_domain" "default" {
  domain         = "changes.com.cn"
  accelerator_id = local.accelerator_id
}

data "alicloud_ga_domains" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
