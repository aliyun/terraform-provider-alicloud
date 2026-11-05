package alb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/alb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Alb LoadBalancerSecurityGroupAttachment. >>> Resource test cases, automatically generated.
// Case LoadBalancerSecurityGroupAttachment 7120
func TestAccAliCloudAlbLoadBalancerSecurityGroupAttachment_basic7120(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_load_balancer_security_group_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAlbLoadBalancerSecurityGroupAttachmentMap7120)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alb.NewAlbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlbLoadBalancerSecurityGroupAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbloadbalancersecuritygroupattachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAlbLoadBalancerSecurityGroupAttachmentBasicDependence7120)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_id": "${alicloud_security_group.create_security_group.id}",
					"load_balancer_id":  "${alicloud_alb_load_balancer.create_alb.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": acctest.CHECKSET,
						"load_balancer_id":  acctest.CHECKSET,
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

var AlicloudAlbLoadBalancerSecurityGroupAttachmentMap7120 = map[string]string{}

func AlicloudAlbLoadBalancerSecurityGroupAttachmentBasicDependence7120(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_alb_zones" "default" {
	}

	resource "alicloud_vpc" "create_vpc" {
  		cidr_block = "192.168.0.0/16"
  		vpc_name   = var.name
	}

	resource "alicloud_vswitch" "create_vsw_1" {
  		vpc_id       = alicloud_vpc.create_vpc.id
  		zone_id      = data.alicloud_alb_zones.default.zones.0.id
  		cidr_block   = "192.168.1.0/24"
  		vswitch_name = var.name
	}

	resource "alicloud_vswitch" "create_vsw_2" {
  		vpc_id       = alicloud_vpc.create_vpc.id
  		zone_id      = data.alicloud_alb_zones.default.zones.1.id
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
