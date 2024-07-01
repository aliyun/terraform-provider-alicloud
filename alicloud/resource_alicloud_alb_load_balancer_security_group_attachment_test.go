package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Alb LoadBalancerSecurityGroupAttachment. >>> Resource test cases, automatically generated.
// Case LoadBalancerSecurityGroupAttachment 7120
func TestAccAliCloudAlbLoadBalancerSecurityGroupAttachment_basic7120(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_load_balancer_security_group_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudAlbLoadBalancerSecurityGroupAttachmentMap7120)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbLoadBalancerSecurityGroupAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbloadbalancersecuritygroupattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlbLoadBalancerSecurityGroupAttachmentBasicDependence7120)
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
					"security_group_id": "${alicloud_security_group.create_security_group.id}",
					"load_balancer_id":  "${alicloud_alb_load_balancer.create_alb.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": CHECKSET,
						"load_balancer_id":  CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudAlbLoadBalancerSecurityGroupAttachmentMap7120 = map[string]string{
	"security_group_id": CHECKSET,
}

func AlicloudAlbLoadBalancerSecurityGroupAttachmentBasicDependence7120(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "create_vpc" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "create_vsw_1" {
  vpc_id       = alicloud_vpc.create_vpc.id
  zone_id      = data.alicloud_zones.default.zones.0.id
  cidr_block   = "192.168.1.0/24"
  vswitch_name = var.name
}

resource "alicloud_vswitch" "create_vsw_2" {
  vpc_id       = alicloud_vpc.create_vpc.id
  zone_id      = data.alicloud_zones.default.zones.1.id
  cidr_block   = "192.168.2.0/24"
  vswitch_name = var.name
}

resource "alicloud_security_group" "create_security_group" {
  name   = var.name
  vpc_id = alicloud_vpc.create_vpc.id
}

resource "alicloud_alb_load_balancer" "create_alb" {
  load_balancer_name    = var.name
  load_balancer_edition = "Standard"
  vpc_id                = alicloud_vpc.create_vpc.id
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  address_type           = "Intranet"
  address_allocated_mode = "Fixed"
  zone_mappings {
    vswitch_id = alicloud_vswitch.create_vsw_2.id
    zone_id    = alicloud_vswitch.create_vsw_2.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.create_vsw_1.id
    zone_id    = alicloud_vswitch.create_vsw_1.zone_id
  }
}


`, name)
}

// Test Alb LoadBalancerSecurityGroupAttachment. <<< Resource test cases, automatically generated.
