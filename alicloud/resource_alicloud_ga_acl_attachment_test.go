package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudGaAclAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_acl_attachment.default"
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	ra := resourceAttrInit(resourceId, AliCloudGaAclAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAclAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaaclattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaAclAttachmentBasicDependence0)
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
					"listener_id": "${alicloud_ga_listener.default.id}",
					"acl_id":      "${alicloud_ga_acl.default.id}",
					"acl_type":    "white",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"listener_id": CHECKSET,
						"acl_id":      CHECKSET,
						"acl_type":    "white",
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

var AliCloudGaAclAttachmentMap0 = map[string]string{}

func AliCloudGaAclAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_ga_accelerators" "default" {
  		status = "active"
		bandwidth_billing_type = "BandwidthPackage"
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
  		// Please run resource ga_accelerator test case to ensure this account has at least one accelerator before run this case.
  		accelerator_id       = data.alicloud_ga_accelerators.default.ids.0
  		bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
	}

	resource "alicloud_ga_listener" "default" {
  		accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  		port_ranges {
    		from_port = 60
    		to_port   = 70
  		}
	}

	resource "alicloud_ga_acl" "default" {
  		acl_name           = var.name
  		address_ip_version = "IPv4"
  		acl_entries {
    		entry             = "192.168.1.0/24"
    		entry_description = "tf-test1"
  		}
	}
`, name)
}
