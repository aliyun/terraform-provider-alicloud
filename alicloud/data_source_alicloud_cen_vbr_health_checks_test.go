package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenVbrHealthCheckDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 2999)
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	cenIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenVbrHealthCheckSourceConfig(rand, map[string]string{
			"cen_id":          `"${alicloud_cen_instance.default.id}"`,
			"vbr_instance_id": `"${alicloud_cen_instance_attachment.vbr.child_instance_id}"`,
		}),
	}
	vbrInstanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenVbrHealthCheckSourceConfig(rand, map[string]string{
			"vbr_instance_id": `"${alicloud_cen_instance_attachment.vbr.child_instance_id}"`,
		}),
	}

	vbrInstanceOwnerIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenVbrHealthCheckSourceConfig(rand, map[string]string{
			"vbr_instance_owner_id": `"${alicloud_cen_instance_attachment.vbr.child_instance_owner_id}"`,
			"vbr_instance_id":       `"${alicloud_cen_instance_attachment.vbr.child_instance_id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenVbrHealthCheckSourceConfig(rand, map[string]string{
			"cen_id":                `"${alicloud_cen_instance.default.id}"`,
			"vbr_instance_id":       `"${alicloud_cen_instance_attachment.vbr.child_instance_id}"`,
			"vbr_instance_owner_id": `"${alicloud_cen_instance_attachment.vbr.child_instance_owner_id}"`,
		}),
	}

	var existCenVbrHealthCheckRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"checks.#":                        "1",
			"checks.0.id":                     CHECKSET,
			"checks.0.cen_id":                 CHECKSET,
			"checks.0.health_check_interval":  "2",
			"checks.0.health_check_source_ip": "192.168.1.2",
			"checks.0.health_check_target_ip": "10.0.0.2",
			"checks.0.healthy_threshold":      "8",
			"checks.0.vbr_instance_id":        CHECKSET,
			"checks.0.vbr_instance_region_id": defaultRegionToTest,
		}
	}

	var fakeCenVbrHealthCheckRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"checks.#": "0",
		}
	}

	var CenVbrHealthCheckRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_vbr_health_checks.default",
		existMapFunc: existCenVbrHealthCheckRecordsMapFunc,
		fakeMapFunc:  fakeCenVbrHealthCheckRecordsMapFunc,
	}

	CenVbrHealthCheckRecordsCheckInfo.dataSourceTestCheck(t, rand, cenIdConf, vbrInstanceIdConf, vbrInstanceOwnerIdConf, allConf)

}

// Because of the VBR instance requires a physical dedicated line，get it form the Environment variable.
func testAccCheckAlicloudCenVbrHealthCheckSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCenVbrHealthCheckDataSource%d"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = "${var.name}"
}

data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^preserved-NODELETING"
}

resource "alicloud_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = %d
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}
resource "alicloud_cen_instance_attachment" "vbr" {
  instance_id = "${alicloud_cen_instance.default.id}"
  child_instance_id = alicloud_express_connect_virtual_border_router.default.id
  child_instance_region_id = "%s"
  child_instance_type = "VBR"
}

resource "alicloud_cen_vbr_health_check" "default" {
	cen_id = "${alicloud_cen_instance.default.id}"
	health_check_source_ip = "192.168.1.2"
	health_check_target_ip = "10.0.0.2"
	vbr_instance_id = "${alicloud_cen_instance_attachment.vbr.child_instance_id}"
	vbr_instance_region_id = "${alicloud_cen_instance_attachment.vbr.child_instance_region_id}"
	health_check_interval = 2
	healthy_threshold = 8
}

data "alicloud_cen_vbr_health_checks" "default" {
  vbr_instance_region_id = "${alicloud_cen_vbr_health_check.default.vbr_instance_region_id}"
%s
}
`, rand, rand, defaultRegionToTest, strings.Join(pairs, "\n   "))
	return config
}
