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

// Test Alb LoadBalancerZoneShiftedAttachment. >>> Resource test cases, automatically generated.
// Case test 9795
func TestAccAliCloudAlbLoadBalancerZoneShiftedAttachment_basic9795(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_load_balancer_zone_shifted_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAlbLoadBalancerZoneShiftedAttachmentMap9795)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alb.NewAlbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlbLoadBalancerZoneShiftedAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbloadbalancerzoneshiftedattachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAlbLoadBalancerZoneShiftedAttachmentBasicDependence9795)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":          "${alicloud_vswitch.defaultDSY0JJ.zone_id}",
					"vswitch_id":       "${alicloud_vswitch.defaultDSY0JJ.id}",
					"load_balancer_id": "${alicloud_alb_load_balancer.default78TIYG.id}",
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

var AlicloudAlbLoadBalancerZoneShiftedAttachmentMap9795 = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudAlbLoadBalancerZoneShiftedAttachmentBasicDependence9795(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "alb_test_tf_vpc" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "alb_test_tf_j" {
  vpc_id       = alicloud_vpc.alb_test_tf_vpc.id
  zone_id      = "cn-beijing-j"
  cidr_block   = "192.168.1.0/24"
  vswitch_name = format("%%s1", var.name)
}

resource "alicloud_vswitch" "alb_test_tf_k" {
  vpc_id       = alicloud_vpc.alb_test_tf_vpc.id
  zone_id      = "cn-beijing-k"
  cidr_block   = "192.168.2.0/24"
  vswitch_name = format("%%s2", var.name)
}

resource "alicloud_vswitch" "defaultDSY0JJ" {
  vpc_id       = alicloud_vpc.alb_test_tf_vpc.id
  zone_id      = "cn-beijing-f"
  cidr_block   = "192.168.3.0/24"
  vswitch_name = format("%%s3", var.name)
}

resource "alicloud_alb_load_balancer" "default78TIYG" {
  load_balancer_edition = "Standard"
  vpc_id                = alicloud_vpc.alb_test_tf_vpc.id
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  address_type           = "Intranet"
  address_allocated_mode = "Fixed"
  zone_mappings {
    vswitch_id = alicloud_vswitch.alb_test_tf_j.id
    zone_id    = alicloud_vswitch.alb_test_tf_j.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.alb_test_tf_k.id
    zone_id    = alicloud_vswitch.alb_test_tf_k.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.defaultDSY0JJ.id
    zone_id    = alicloud_vswitch.defaultDSY0JJ.zone_id
  }
}


`, name)
}

// Test Alb LoadBalancerZoneShiftedAttachment. <<< Resource test cases, automatically generated.
