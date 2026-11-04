package nlb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/nlb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Nlb LoadBalancerZoneShiftedAttachment. >>> Resource test cases, automatically generated.
// Case shiftzone 9761
func TestAccAliCloudNlbLoadBalancerZoneShiftedAttachment_basic9761(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_load_balancer_zone_shifted_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNlbLoadBalancerZoneShiftedAttachmentMap9761)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbLoadBalancerZoneShiftedAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnlb%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNlbLoadBalancerZoneShiftedAttachmentBasicDependence9761)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":          "${alicloud_vswitch.vsw1.zone_id}",
					"vswitch_id":       "${alicloud_vswitch.vsw1.id}",
					"load_balancer_id": "${alicloud_nlb_load_balancer.nlb.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":          acctest.CHECKSET,
						"vswitch_id":       acctest.CHECKSET,
						"load_balancer_id": acctest.CHECKSET,
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

var AlicloudNlbLoadBalancerZoneShiftedAttachmentMap9761 = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudNlbLoadBalancerZoneShiftedAttachmentBasicDependence9761(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "vpc" {
  description = "test"
  cidr_block  = "10.0.0.0/8"
  enable_ipv6 = true
  vpc_name    = "tf-testacc-878"
}

resource "alicloud_vswitch" "vsw1" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = "cn-beijing-l"
  cidr_block   = "10.0.1.0/24"
  vswitch_name = "tf-testacc-949"
}

resource "alicloud_vswitch" "vsw2" {
  vpc_id               = alicloud_vpc.vpc.id
  zone_id              = "cn-beijing-k"
  cidr_block           = "10.0.2.0/24"
  vswitch_name         = "tf-testacc-884"
  ipv6_cidr_block_mask = "8"
}

resource "alicloud_nlb_load_balancer" "nlb" {
  zone_mappings {
    vswitch_id = alicloud_vswitch.vsw1.id
    zone_id    = alicloud_vswitch.vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vsw2.id
    zone_id    = alicloud_vswitch.vsw2.zone_id
  }
  vpc_id       = alicloud_vpc.vpc.id
  address_type = "Intranet"
}


`, name)
}

// Test Nlb LoadBalancerZoneShiftedAttachment. <<< Resource test cases, automatically generated.
