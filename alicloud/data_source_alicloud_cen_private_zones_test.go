package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenPrivateZonesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenPrivateZonesSourceConfig(rand, map[string]string{
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
			"ids":    `[split(":",alicloud_cen_private_zone.default.id)[1]]`,
		}),
		fakeConfig: testAccCheckAlicloudCenPrivateZonesSourceConfig(rand, map[string]string{
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
			"ids":    `["fake"]`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenPrivateZonesSourceConfig(rand, map[string]string{
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
			"ids":    `[split(":",alicloud_cen_private_zone.default.id)[1]]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudCenPrivateZonesSourceConfig(rand, map[string]string{
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
			"ids":    `[split(":",alicloud_cen_private_zone.default.id)[1]]`,
			"status": `"Creating"`,
		}),
	}

	hostRegionIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenPrivateZonesSourceConfig(rand, map[string]string{
			"cen_id":         `"${alicloud_cen_instance.default.id}"`,
			"ids":            `[split(":",alicloud_cen_private_zone.default.id)[1]]`,
			"host_region_id": fmt.Sprintf(`"%s"`, defaultRegionToTest),
		}),
		fakeConfig: testAccCheckAlicloudCenPrivateZonesSourceConfig(rand, map[string]string{
			"cen_id":         `"${alicloud_cen_instance.default.id}"`,
			"ids":            `[split(":",alicloud_cen_private_zone.default.id)[1]]`,
			"host_region_id": `"fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenPrivateZonesSourceConfig(rand, map[string]string{
			"cen_id":         `"${alicloud_cen_instance.default.id}"`,
			"ids":            `[split(":",alicloud_cen_private_zone.default.id)[1]]`,
			"status":         `"Active"`,
			"host_region_id": fmt.Sprintf(`"%s"`, defaultRegionToTest),
		}),
		fakeConfig: testAccCheckAlicloudCenPrivateZonesSourceConfig(rand, map[string]string{
			"cen_id":         `"${alicloud_cen_instance.default.id}"`,
			"ids":            `[split(":",alicloud_cen_private_zone.default.id)[1]]`,
			"status":         `"Active"`,
			"host_region_id": `"fake"`,
		}),
	}

	var existCenPrivateZonesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"zones.#":                          "1",
			"zones.0.cen_id":                   CHECKSET,
			"zones.0.private_zone_dns_servers": CHECKSET,
			"zones.0.access_region_id":         defaultRegionToTest,
			"zones.0.host_region_id":           defaultRegionToTest,
			"zones.0.host_vpc_id":              CHECKSET,
			"zones.0.status":                   "Active",
		}
	}

	var fakeCenPrivateZonesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var cenPrivateZonesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_private_zones.default",
		existMapFunc: existCenPrivateZonesRecordsMapFunc,
		fakeMapFunc:  fakeCenPrivateZonesRecordsMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	cenPrivateZonesRecordsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, hostRegionIdConf, allConf)

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
resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/12"
}
resource "alicloud_cen_instance_attachment" "default" {
  instance_id              = alicloud_cen_instance.default.id
  child_instance_id        = alicloud_vpc.default.id
  child_instance_type      = "VPC"
  child_instance_region_id = "%[2]s"
}

resource "alicloud_cen_private_zone" "default" {
  access_region_id = "%[2]s"
  cen_id           = alicloud_cen_instance.default.id
  host_region_id   = "%[2]s"
  host_vpc_id      = alicloud_vpc.default.id
  depends_on       = [alicloud_cen_instance_attachment.default]
}

data "alicloud_cen_private_zones" "default" {
%s
}
`, rand, defaultRegionToTest, strings.Join(pairs, "\n   "))
	return config
}
