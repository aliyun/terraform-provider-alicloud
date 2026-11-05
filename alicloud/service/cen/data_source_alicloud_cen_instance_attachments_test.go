package cen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCenInstanceAttachmentsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenInstanceAttachmentsSourceConfig(rand, map[string]string{
			"status": `"Attached"`,
		}),
		FakeConfig: testAccCheckAlicloudCenInstanceAttachmentsSourceConfig(rand, map[string]string{
			"status": `"Attaching"`,
		}),
	}

	childInstanceRegionIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenInstanceAttachmentsSourceConfig(rand, map[string]string{
			"child_instance_region_id": fmt.Sprintf(`"%s"`, acctest.DefaultRegionToTest),
		}),
		FakeConfig: testAccCheckAlicloudCenInstanceAttachmentsSourceConfig(rand, map[string]string{
			"child_instance_region_id": `"fake"`,
		}),
	}

	childInstanceTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenInstanceAttachmentsSourceConfig(rand, map[string]string{
			"child_instance_type": `"VPC"`,
		}),
		FakeConfig: testAccCheckAlicloudCenInstanceAttachmentsSourceConfig(rand, map[string]string{
			"child_instance_type": `"CCN"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenInstanceAttachmentsSourceConfig(rand, map[string]string{
			"status":                   `"Attached"`,
			"child_instance_region_id": fmt.Sprintf(`"%s"`, acctest.DefaultRegionToTest),
			"child_instance_type":      `"VPC"`,
		}),
		FakeConfig: testAccCheckAlicloudCenInstanceAttachmentsSourceConfig(rand, map[string]string{
			"status":                   `"Attaching"`,
			"child_instance_region_id": `"fake"`,
			"child_instance_type":      `"CCN"`,
		}),
	}

	var existCenInstanceAttachmentsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"attachments.#":    "1",
			"ids.#":            "1",
			"attachments.0.id": acctest.CHECKSET,
			"attachments.0.child_instance_attach_time": acctest.CHECKSET,
			"attachments.0.child_instance_id":          acctest.CHECKSET,
			"attachments.0.child_instance_owner_id":    acctest.CHECKSET,
			"attachments.0.child_instance_region_id":   fmt.Sprintf("%s", acctest.DefaultRegionToTest),
			"attachments.0.child_instance_type":        "VPC",
			"attachments.0.instance_id":                acctest.CHECKSET,
			"attachments.0.status":                     "Attached",
		}
	}

	var fakeCenInstanceAttachmentsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"attachments.#": "0",
			"ids.#":         "0",
		}
	}

	var cenInstanceAttachmentsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cen_instance_attachments.default",
		ExistMapFunc: existCenInstanceAttachmentsRecordsMapFunc,
		FakeMapFunc:  fakeCenInstanceAttachmentsRecordsMapFunc,
	}

	cenInstanceAttachmentsCheckInfo.DataSourceTestCheck(t, rand, statusConf, childInstanceRegionIdConf, childInstanceTypeConf, allConf)

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
  name_regex = "^default-NODELETING$"
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
`, rand, acctest.DefaultRegionToTest, strings.Join(pairs, "\n   "))
	return config
}
