package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenVbrHealthCheckDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	cenIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenVbrHealthCheckSourceConfig(rand, map[string]string{
			"cen_id":          `"${alicloud_cen_instance.default.id}"`,
			"vbr_instance_id": `"${alicloud_cen_instance_attachment.vbr.child_instance_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenVbrHealthCheckSourceConfig(rand, map[string]string{
			"cen_id":          `"${alicloud_cen_instance.default.id}"`,
			"vbr_instance_id": `"fake"`,
		}),
	}
	vbrInstanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenVbrHealthCheckSourceConfig(rand, map[string]string{
			"vbr_instance_id": `"${alicloud_cen_instance_attachment.vbr.child_instance_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenVbrHealthCheckSourceConfig(rand, map[string]string{
			"vbr_instance_id": `"fake"`,
		}),
	}

	vbrInstanceOwnerIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenVbrHealthCheckSourceConfig(rand, map[string]string{
			"vbr_instance_owner_id": `"${alicloud_cen_instance_attachment.vbr.child_instance_owner_id}"`,
			"vbr_instance_id":       `"${alicloud_cen_instance_attachment.vbr.child_instance_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenVbrHealthCheckSourceConfig(rand, map[string]string{
			"vbr_instance_owner_id": `123456`,
			"vbr_instance_id":       `"${alicloud_cen_instance_attachment.vbr.child_instance_id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenVbrHealthCheckSourceConfig(rand, map[string]string{
			"cen_id":                `"${alicloud_cen_instance.default.id}"`,
			"vbr_instance_id":       `"${alicloud_cen_instance_attachment.vbr.child_instance_id}"`,
			"vbr_instance_owner_id": `"${alicloud_cen_instance_attachment.vbr.child_instance_owner_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenVbrHealthCheckSourceConfig(rand, map[string]string{
			"cen_id":                `"${alicloud_cen_instance.default.id}"`,
			"vbr_instance_id":       `"fake"`,
			"vbr_instance_owner_id": `123456`,
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
			"checks.0.vbr_instance_id":        os.Getenv("VBR_INSTANCE_ID"),
			"checks.0.vbr_instance_region_id": os.Getenv("VBR_INSTANCE_REGION_ID"),
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

	preCheck := func() {
		testAccPreCheckWithCenVbrHealthCheckSetting(t)
	}

	CenVbrHealthCheckRecordsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, cenIdConf, vbrInstanceIdConf, vbrInstanceOwnerIdConf, allConf)

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

resource "alicloud_cen_instance_attachment" "vbr" {
  instance_id = "${alicloud_cen_instance.default.id}"
  child_instance_id = "%[2]s"
  child_instance_region_id = "%[3]s"
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
`, rand, os.Getenv("VBR_INSTANCE_ID"), os.Getenv("VBR_INSTANCE_REGION_ID"), strings.Join(pairs, "\n   "))
	return config
}
