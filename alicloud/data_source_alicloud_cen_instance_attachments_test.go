package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenInstanceAttachmentsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenInstanceAttachmentsSourceConfig(rand, map[string]string{
			"status": `"Attached"`,
		}),
		fakeConfig: testAccCheckAlicloudCenInstanceAttachmentsSourceConfig(rand, map[string]string{
			"status": `"Attaching"`,
		}),
	}

	childInstanceRegionIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenInstanceAttachmentsSourceConfig(rand, map[string]string{
			"child_instance_region_id": fmt.Sprintf(`"%s"`, defaultRegionToTest),
		}),
		fakeConfig: testAccCheckAlicloudCenInstanceAttachmentsSourceConfig(rand, map[string]string{
			"child_instance_region_id": `"fake"`,
		}),
	}

	childInstanceTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenInstanceAttachmentsSourceConfig(rand, map[string]string{
			"child_instance_type": `"VPC"`,
		}),
		fakeConfig: testAccCheckAlicloudCenInstanceAttachmentsSourceConfig(rand, map[string]string{
			"child_instance_type": `"CCN"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenInstanceAttachmentsSourceConfig(rand, map[string]string{
			"status":                   `"Attached"`,
			"child_instance_region_id": fmt.Sprintf(`"%s"`, defaultRegionToTest),
			"child_instance_type":      `"VPC"`,
		}),
		fakeConfig: testAccCheckAlicloudCenInstanceAttachmentsSourceConfig(rand, map[string]string{
			"status":                   `"Attaching"`,
			"child_instance_region_id": `"fake"`,
			"child_instance_type":      `"CCN"`,
		}),
	}

	var existCenInstanceAttachmentsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"attachments.#":    "1",
			"ids.#":            "1",
			"attachments.0.id": CHECKSET,
			"attachments.0.child_instance_attach_time": CHECKSET,
			"attachments.0.child_instance_id":          CHECKSET,
			"attachments.0.child_instance_owner_id":    CHECKSET,
			"attachments.0.child_instance_region_id":   fmt.Sprintf("%s", defaultRegionToTest),
			"attachments.0.child_instance_type":        "VPC",
			"attachments.0.instance_id":                CHECKSET,
			"attachments.0.status":                     "Attached",
		}
	}

	var fakeCenInstanceAttachmentsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"attachments.#": "0",
			"ids.#":         "0",
		}
	}

	var cenInstanceAttachmentsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_instance_attachments.default",
		existMapFunc: existCenInstanceAttachmentsRecordsMapFunc,
		fakeMapFunc:  fakeCenInstanceAttachmentsRecordsMapFunc,
	}

	cenInstanceAttachmentsCheckInfo.dataSourceTestCheck(t, rand, statusConf, childInstanceRegionIdConf, childInstanceTypeConf, allConf)

}

func testAccCheckAlicloudCenInstanceAttachmentsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCen%d"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  description = "for test"
}

resource "alicloud_cen_instance_attachment" "default" {
  child_instance_id = data.alicloud_vpcs.default.ids.0
  child_instance_region_id = "%s"
  instance_id = alicloud_cen_instance.default.id
  child_instance_type = "VPC"
}

data "alicloud_cen_instance_attachments" "default" {
  instance_id = alicloud_cen_instance_attachment.default.instance_id
	%s
}
`, rand, defaultRegionToTest, strings.Join(pairs, "\n   "))
	return config
}
