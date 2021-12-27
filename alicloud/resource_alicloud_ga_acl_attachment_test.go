package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGaAclAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_acl_attachment.default"
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudGaAclAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAclAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaaclattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaAclAttachmentBasicDependence0)
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
					"acl_type":    "white",
					"listener_id": "${alicloud_ga_listener.default.id}",
					"acl_id":      "${alicloud_ga_acl.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_type":    "white",
						"listener_id": CHECKSET,
						"acl_id":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudGaAclAttachmentMap0 = map[string]string{}

func AlicloudGaAclAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_ga_accelerators" "default" {
  status = "active"
}

data "alicloud_ga_bandwidth_packages" "default" {
  status = "active"
}

resource "alicloud_ga_accelerator" "default" {
  count            = length(data.alicloud_ga_accelerators.default.accelerators) > 0 ? 0 : 1
  duration         = 1
  auto_use_coupon  = true
  spec             = "1"
}

resource "alicloud_ga_bandwidth_package" "default" {
  count                  = length(data.alicloud_ga_bandwidth_packages.default.packages) > 0 ? 0 : 1
  bandwidth              = 20
  type                   = "Basic"
  bandwidth_type         = "Basic"
  duration               = 1
  ratio                  = 30
  auto_pay               = true
  auto_use_coupon        = true
}

locals {
  accelerator_id       = length(data.alicloud_ga_accelerators.default.accelerators) > 0 ? data.alicloud_ga_accelerators.default.accelerators.0.id : alicloud_ga_accelerator.default.0.id
  bandwidth_package_id = length(data.alicloud_ga_bandwidth_packages.default.packages) > 0 ? data.alicloud_ga_bandwidth_packages.default.packages.0.id : alicloud_ga_bandwidth_package.default.0.id
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = local.accelerator_id
  bandwidth_package_id = local.bandwidth_package_id
}

resource "alicloud_ga_listener" "default" {
  depends_on     = [alicloud_ga_bandwidth_package_attachment.default]
  accelerator_id = local.accelerator_id
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
