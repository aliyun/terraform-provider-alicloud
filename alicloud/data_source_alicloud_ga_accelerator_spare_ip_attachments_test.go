package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGaAcceleratorSpareIpAttachmentsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaAcceleratorSpareIpAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_accelerator_spare_ip_attachment.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudGaAcceleratorSpareIpAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_accelerator_spare_ip_attachment.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaAcceleratorSpareIpAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ga_accelerator_spare_ip_attachment.default.id}"]`,
			"status": `"inuse"`,
		}),
		fakeConfig: testAccCheckAlicloudGaAcceleratorSpareIpAttachmentsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ga_accelerator_spare_ip_attachment.default.id}"]`,
			"status": `"active"`,
		}),
	}
	var existAlicloudGaAcceleratorSpareIpAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"attachments.#":                "1",
			"attachments.0.accelerator_id": CHECKSET,
			"attachments.0.id":             CHECKSET,
			"attachments.0.spare_ip":       "127.0.0.1",
			"attachments.0.status":         "inuse",
		}
	}
	var fakeAlicloudGaAcceleratorSpareIpAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudGaAcceleratorSpareIpAttachmentsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ga_accelerator_spare_ip_attachments.default",
		existMapFunc: existAlicloudGaAcceleratorSpareIpAttachmentsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudGaAcceleratorSpareIpAttachmentsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudGaAcceleratorSpareIpAttachmentsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf)
}
func testAccCheckAlicloudGaAcceleratorSpareIpAttachmentsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccAcceleratorSpareIpAttachment-%d"
}

resource "alicloud_ga_accelerator" "default" {
  duration         = 1
  spec             = "1"
  accelerator_name = var.name
  auto_use_coupon  = true
  description      = var.name
}
resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth              = 100
  type                   = "Basic"
  bandwidth_type         = "Basic"
  payment_type           = "PayAsYouGo"
  billing_type           = "PayBy95"
  ratio                  = 30
  bandwidth_package_name = var.name
  auto_pay               = true
  auto_use_coupon        = true
}
resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = alicloud_ga_accelerator.default.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_accelerator_spare_ip_attachment" "default" {
  accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  spare_ip       = "127.0.0.1"
}

data "alicloud_ga_accelerator_spare_ip_attachments" "default" {	
	accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
