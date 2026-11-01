package cen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCenRouteServicesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenRouteServicesSourceConfig(rand, map[string]string{
			"cen_id":           `"${alicloud_cen_route_service.this.cen_id}"`,
			"access_region_id": `"${alicloud_cen_route_service.this.access_region_id}"`,
			"host":             `"${alicloud_cen_route_service.this.host}"`,
			"host_region_id":   `"${alicloud_cen_route_service.this.host_region_id}"`,
			"status":           `"Active"`,
		}),
		FakeConfig: testAccCheckAlicloudCenRouteServicesSourceConfig(rand, map[string]string{
			"cen_id":           `"${alicloud_cen_route_service.this.cen_id}"`,
			"access_region_id": `"${alicloud_cen_route_service.this.access_region_id}"`,
			"host":             `"${alicloud_cen_route_service.this.host}"`,
			"host_region_id":   `"${alicloud_cen_route_service.this.host_region_id}"`,
			"status":           `"Creating"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenRouteServicesSourceConfig(rand, map[string]string{
			"cen_id":           `"${alicloud_cen_route_service.this.cen_id}"`,
			"access_region_id": `"${alicloud_cen_route_service.this.access_region_id}"`,
			"host":             `"${alicloud_cen_route_service.this.host}"`,
			"host_region_id":   `"${alicloud_cen_route_service.this.host_region_id}"`,
			"host_vpc_id":      `"${alicloud_cen_route_service.this.host_vpc_id}"`,
			"status":           `"Active"`,
		}),
		FakeConfig: testAccCheckAlicloudCenRouteServicesSourceConfig(rand, map[string]string{
			"cen_id":           `"${alicloud_cen_route_service.this.cen_id}"`,
			"access_region_id": `"${alicloud_cen_route_service.this.access_region_id}"`,
			"host":             `"${alicloud_cen_route_service.this.host}"`,
			"host_region_id":   `"${alicloud_cen_route_service.this.host_region_id}"`,
			"host_vpc_id":      `"${alicloud_cen_route_service.this.host_vpc_id}"`,
			"status":           `"Creating"`,
		}),
	}

	var existCenRouteServicesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"services.#":                  "1",
			"services.0.id":               acctest.CHECKSET,
			"services.0.access_region_id": acctest.DefaultRegionToTest,
			"services.0.cen_id":           acctest.CHECKSET,
			"services.0.cidrs.#":          "1",
			"services.0.description":      fmt.Sprintf("tf-testAccCenRouteServicesDataSource%d", rand),
			"services.0.host":             "100.118.28.52/32",
			"services.0.host_region_id":   acctest.DefaultRegionToTest,
			"services.0.host_vpc_id":      acctest.CHECKSET,
			"services.0.status":           "Active",
			"services.0.update_interval":  "",
		}
	}

	var fakeCenRouteServicesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"services.#": "0",
		}
	}

	var cenRouteServicesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cen_route_services.default",
		ExistMapFunc: existCenRouteServicesMapFunc,
		FakeMapFunc:  fakeCenRouteServicesMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
	}

	cenRouteServicesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, statusConf, allConf)

}

func testAccCheckAlicloudCenRouteServicesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCenRouteServicesDataSource%d"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = var.name
}

resource "alicloud_cen_instance_attachment" "vpc" {
  instance_id              = alicloud_cen_instance.example.id
  child_instance_id        = data.alicloud_vpcs.default.ids.0
  child_instance_type      = "VPC"
  child_instance_region_id = "%s"
}

resource "alicloud_cen_route_service" "this" {
  access_region_id = alicloud_cen_instance_attachment.vpc.child_instance_region_id
  cen_id           = alicloud_cen_instance_attachment.vpc.instance_id
  host             = "100.118.28.52/32"
  host_region_id   = alicloud_cen_instance_attachment.vpc.child_instance_region_id
  host_vpc_id      = alicloud_cen_instance_attachment.vpc.child_instance_id
  description      = var.name
}

data "alicloud_cen_route_services" "default" {
%s
}
`, rand, acctest.DefaultRegionToTest, strings.Join(pairs, "\n   "))
	return config
}
