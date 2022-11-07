package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGaBandwidthPackageAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_bandwidth_package_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudGaBandwidthPackageAttachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaBandwidthPackageAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", AlicloudGaBandwidthPackageAttachmentBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id":       "${local.accelerator_id}",
					"bandwidth_package_id": "${alicloud_ga_bandwidth_package.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":       CHECKSET,
						"bandwidth_package_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_id": "${alicloud_ga_bandwidth_package.update.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudGaBandwidthPackageAttachmentMap = map[string]string{
	"accelerators.#": CHECKSET,
	"status":         "active",
}

func AlicloudGaBandwidthPackageAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	locals {
  		accelerator_id = length(data.alicloud_ga_accelerators.default.accelerators) > 0 ? data.alicloud_ga_accelerators.default.accelerators.0.id : alicloud_ga_accelerator.default.0.id
	}

	data "alicloud_ga_accelerators" "default" {
  		status = "active"
	}

	resource "alicloud_ga_accelerator" "default" {
		count            = length(data.alicloud_ga_accelerators.default.accelerators) > 0 ? 0 : 1
		duration         = 1
		spec             = "1"
		accelerator_name = var.name
  		auto_use_coupon  = true
  		description      = var.name
	}

	resource "alicloud_ga_bandwidth_package" "default" {
  		bandwidth       = 100
  		type            = "Basic"
  		bandwidth_type  = "Basic"
  		payment_type    = "PayAsYouGo"
  		billing_type    = "PayBy95"
  		ratio           = 30
  		auto_pay        = true
  		auto_use_coupon = true
	}

	resource "alicloud_ga_bandwidth_package" "update" {
  		bandwidth       = 100
  		type            = "Basic"
  		bandwidth_type  = "Basic"
  		payment_type    = "PayAsYouGo"
  		billing_type    = "PayBy95"
  		ratio           = 30
  		auto_pay        = true
  		auto_use_coupon = true
	}
`, name)
}
