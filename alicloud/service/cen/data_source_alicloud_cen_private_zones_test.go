package cen_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCenPrivateZonesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenPrivateZonesSourceConfig(rand, map[string]string{
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
			"ids":    `[alicloud_cen_private_zone.default.id]`,
		}),
		FakeConfig: testAccCheckAlicloudCenPrivateZonesSourceConfig(rand, map[string]string{
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
			"ids":    `["fake"]`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenPrivateZonesSourceConfig(rand, map[string]string{
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
			"ids":    `[alicloud_cen_private_zone.default.id]`,
			"status": `"Active"`,
		}),
		FakeConfig: testAccCheckAlicloudCenPrivateZonesSourceConfig(rand, map[string]string{
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
			"ids":    `[alicloud_cen_private_zone.default.id]`,
			"status": `"Creating"`,
		}),
	}

	hostRegionIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenPrivateZonesSourceConfig(rand, map[string]string{
			"cen_id":         `"${alicloud_cen_instance.default.id}"`,
			"ids":            `[alicloud_cen_private_zone.default.id]`,
			"host_region_id": fmt.Sprintf(`"%s"`, os.Getenv("ALICLOUD_REGION")),
		}),
		FakeConfig: testAccCheckAlicloudCenPrivateZonesSourceConfig(rand, map[string]string{
			"cen_id":         `"${alicloud_cen_instance.default.id}"`,
			"ids":            `[alicloud_cen_private_zone.default.id]`,
			"host_region_id": `"fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenPrivateZonesSourceConfig(rand, map[string]string{
			"cen_id":         `"${alicloud_cen_instance.default.id}"`,
			"ids":            `[alicloud_cen_private_zone.default.id]`,
			"status":         `"Active"`,
			"host_region_id": fmt.Sprintf(`"%s"`, os.Getenv("ALICLOUD_REGION")),
		}),
		FakeConfig: testAccCheckAlicloudCenPrivateZonesSourceConfig(rand, map[string]string{
			"cen_id":         `"${alicloud_cen_instance.default.id}"`,
			"ids":            `[alicloud_cen_private_zone.default.id]`,
			"status":         `"Active"`,
			"host_region_id": `"fake"`,
		}),
	}

	var existCenPrivateZonesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"zones.#":                          "1",
			"zones.0.id":                       acctest.CHECKSET,
			"zones.0.cen_id":                   acctest.CHECKSET,
			"zones.0.private_zone_dns_servers": acctest.CHECKSET,
			"zones.0.access_region_id":         os.Getenv("ALICLOUD_REGION"),
			"zones.0.host_region_id":           os.Getenv("ALICLOUD_REGION"),
			"zones.0.host_vpc_id":              acctest.CHECKSET,
			"zones.0.status":                   "Active",
		}
	}

	var fakeCenPrivateZonesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var cenPrivateZonesRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cen_private_zones.default",
		ExistMapFunc: existCenPrivateZonesRecordsMapFunc,
		FakeMapFunc:  fakeCenPrivateZonesRecordsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	cenPrivateZonesRecordsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, hostRegionIdConf, allConf)

}

func testAccCheckAlicloudCenPrivateZonesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccCenPrivateZones-%d"
}
resource "alicloud_cen_instance" "default" {
  name = var.name
}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
resource "alicloud_cen_instance_attachment" "default" {
  instance_id              = alicloud_cen_instance.default.id
  child_instance_id        = data.alicloud_vpcs.default.ids.0
  child_instance_type      = "VPC"
  child_instance_region_id = "%[2]s"
}

resource "alicloud_cen_private_zone" "default" {
  access_region_id = "%[2]s"
  cen_id           = alicloud_cen_instance.default.id
  host_region_id   = "%[2]s"
  host_vpc_id      = data.alicloud_vpcs.default.ids.0
  depends_on       = [alicloud_cen_instance_attachment.default]
}

data "alicloud_cen_private_zones" "default" {
%s
}
`, rand, os.Getenv("ALICLOUD_REGION"), strings.Join(pairs, "\n   "))
	return config
}
