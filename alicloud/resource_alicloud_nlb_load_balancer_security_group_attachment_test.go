package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAliCloudNlbLoadBalancerSecurityGroupAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_load_balancer_security_group_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudNlbLoadBalancerSecurityGroupAttachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNlbLoadBalancerSecurityGroupAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sNlbLoadBalancerSecurityGroupAttachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNlbLoadBalancerSecurityGroupAttachmentBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.NLBSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_id": "${alicloud_security_group.default.id}",
					"load_balancer_id":  "${alicloud_nlb_load_balancer.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": CHECKSET,
						"load_balancer_id":  CHECKSET,
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudNlbLoadBalancerSecurityGroupAttachmentMap = map[string]string{}

func AlicloudNlbLoadBalancerSecurityGroupAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_nlb_zones" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_vswitches" "default_1" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_nlb_zones.default.zones.0.id
	}

	data "alicloud_vswitches" "default_2" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_nlb_zones.default.zones.1.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	locals {
  		zone_id_1    = data.alicloud_nlb_zones.default.zones.0.id
  		vswitch_id_1 = data.alicloud_vswitches.default_1.ids[0]
  		zone_id_2    = data.alicloud_nlb_zones.default.zones.1.id
  		vswitch_id_2 = data.alicloud_vswitches.default_2.ids[0]
	}

	resource "alicloud_nlb_load_balancer" "default" {
  		load_balancer_name = var.name
  		resource_group_id  = data.alicloud_resource_manager_resource_groups.default.ids.0
  		load_balancer_type = "Network"
  		address_type       = "Internet"
  		address_ip_version = "Ipv4"
  		tags = {
    		Created = "tfTestAcc0"
    		For     = "Tftestacc 0"
  		}
  		vpc_id = data.alicloud_vpcs.default.ids.0
  		zone_mappings {
    		vswitch_id = local.vswitch_id_1
    		zone_id    = local.zone_id_1
  		}
  		zone_mappings {
    		vswitch_id = local.vswitch_id_2
    		zone_id    = local.zone_id_2
  		}
	}
`, name)
}

// Test Nlb LoadBalancerSecurityGroupAttachment. >>> Resource test cases, automatically generated.
// Case 4781
func TestAccAliCloudNlbLoadBalancerSecurityGroupAttachment_basic4781(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_load_balancer_security_group_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudNlbLoadBalancerSecurityGroupAttachmentMap4781)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNlbLoadBalancerSecurityGroupAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbloadbalancersecuritygroupattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNlbLoadBalancerSecurityGroupAttachmentBasicDependence4781)
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
					"load_balancer_id":  "${alicloud_nlb_load_balancer.nlb.id}",
					"security_group_id": "${alicloud_security_group.securityGroup.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":  CHECKSET,
						"security_group_id": CHECKSET,
					}),
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

var AlicloudNlbLoadBalancerSecurityGroupAttachmentMap4781 = map[string]string{
	"security_group_id": CHECKSET,
}

func AlicloudNlbLoadBalancerSecurityGroupAttachmentBasicDependence4781(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_nlb_zones" "default" {
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name

}

resource "alicloud_vswitch" "vswtich" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_nlb_zones.default.zones.0.id
  vswitch_name = var.name

  cidr_block = "192.168.10.0/24"
}

resource "alicloud_vswitch" "vswtich2" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_nlb_zones.default.zones.1.id
  vswitch_name = var.name

  cidr_block = "192.168.30.0/24"
}

resource "alicloud_nlb_load_balancer" "nlb" {
  zone_mappings {
    vswitch_id = alicloud_vswitch.vswtich2.id
    zone_id    = alicloud_vswitch.vswtich2.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vswtich.id
    zone_id    = alicloud_vswitch.vswtich.zone_id
  }
  load_balancer_type = "Network"
  vpc_id             = alicloud_vpc.vpc.id
  address_type       = "Internet"
  address_ip_version = "Ipv4"
}

resource "alicloud_security_group" "securityGroup" {
  name = var.name

  vpc_id = alicloud_vpc.vpc.id
}


`, name)
}

// Test Nlb LoadBalancerSecurityGroupAttachment. <<< Resource test cases, automatically generated.
