package sae_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudSAEGreyTagRouteDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 100)
	acctest.CheckoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSaeGreyTagRouteDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sae_grey_tag_route.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudSaeGreyTagRouteDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sae_grey_tag_route.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSaeGreyTagRouteDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_sae_grey_tag_route.default.grey_tag_route_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudSaeGreyTagRouteDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_sae_grey_tag_route.default.grey_tag_route_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSaeGreyTagRouteDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_sae_grey_tag_route.default.id}"]`,
			"name_regex": `"${alicloud_sae_grey_tag_route.default.grey_tag_route_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudSaeGreyTagRouteDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_sae_grey_tag_route.default.id}_fake"]`,
			"name_regex": `"${alicloud_sae_grey_tag_route.default.grey_tag_route_name}_fake"`,
		}),
	}
	var existAlicloudSaeGreyTagRouteDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                   "1",
			"names.#":                                 "1",
			"routes.#":                                "1",
			"routes.0.id":                             acctest.CHECKSET,
			"routes.0.description":                    fmt.Sprintf("tf-testAccsae-%d", rand),
			"routes.0.grey_tag_route_name":            fmt.Sprintf("tf-testAccsae-%d", rand),
			"routes.0.dubbo_rules.#":                  "1",
			"routes.0.dubbo_rules.0.method_name":      acctest.CHECKSET,
			"routes.0.dubbo_rules.0.service_name":     acctest.CHECKSET,
			"routes.0.dubbo_rules.0.version":          acctest.CHECKSET,
			"routes.0.dubbo_rules.0.condition":        acctest.CHECKSET,
			"routes.0.dubbo_rules.0.group":            acctest.CHECKSET,
			"routes.0.dubbo_rules.0.items.#":          "1",
			"routes.0.dubbo_rules.0.items.0.index":    acctest.CHECKSET,
			"routes.0.dubbo_rules.0.items.0.expr":     acctest.CHECKSET,
			"routes.0.dubbo_rules.0.items.0.cond":     acctest.CHECKSET,
			"routes.0.dubbo_rules.0.items.0.value":    acctest.CHECKSET,
			"routes.0.dubbo_rules.0.items.0.operator": acctest.CHECKSET,
			"routes.0.sc_rules.#":                     "1",
			"routes.0.sc_rules.0.path":                acctest.CHECKSET,
			"routes.0.sc_rules.0.condition":           acctest.CHECKSET,
			"routes.0.sc_rules.0.items.#":             "1",
			"routes.0.sc_rules.0.items.0.name":        acctest.CHECKSET,
			"routes.0.sc_rules.0.items.0.type":        acctest.CHECKSET,
			"routes.0.sc_rules.0.items.0.cond":        acctest.CHECKSET,
			"routes.0.sc_rules.0.items.0.value":       acctest.CHECKSET,
			"routes.0.sc_rules.0.items.0.operator":    acctest.CHECKSET,
		}
	}
	var fakeAlicloudSaeGreyTagRouteDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudSaeGreyTagRouteCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_sae_grey_tag_routes.default",
		ExistMapFunc: existAlicloudSaeGreyTagRouteDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudSaeGreyTagRouteDataSourceNameMapFunc,
	}

	alicloudSaeGreyTagRouteCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudSaeGreyTagRouteDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccsae-%d"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_zones" "default" {}

resource "alicloud_sae_namespace" "default" {
  namespace_description = var.name
  namespace_id          = join(":",["%s","%d"])
  namespace_name        = var.name
}

resource "alicloud_sae_application" "default" {
  app_description = var.name
  app_name        = var.name
  namespace_id    = alicloud_sae_namespace.default.id
  image_url       = "registry-vpc.cn-hangzhou.aliyuncs.com/lxepoo/apache-php5"
  package_type    = "Image"
  jdk             = "Open JDK 8"
  vswitch_id      = data.alicloud_vswitches.default.ids.0
  vpc_id          = data.alicloud_vpcs.default.ids.0
  timezone        = "Asia/Shanghai"
  replicas        = "5"
  cpu             = "500"
  memory          = "2048"
}

resource "alicloud_sae_grey_tag_route" "default" {
  grey_tag_route_name        = var.name
  description = var.name
  app_id             = alicloud_sae_application.default.id
  sc_rules {
    items {
      type     = "param"
      name     = "tftest"
      operator = "rawvalue"
      value    = "test"
      cond     = "=="
    }
    path      = "/tf/test"
    condition = "AND"
  }

  dubbo_rules {
    items {
      cond     = "=="
      expr     = ".key1"
      index    = "1"
      operator = "rawvalue"
      value    = "value1"
    }
    condition    = "OR"
    group        = "DUBBO"
    method_name  = "test"
    service_name = "com.test.service"
    version      = "1.0.0"
  }
}

data "alicloud_sae_grey_tag_routes" "default" {	
   app_id = alicloud_sae_application.default.id
	%s
}
`, rand, acctest.DefaultRegionToTest, rand, strings.Join(pairs, " \n "))
	return config
}
