package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Alb LoadBalancerAccessLogConfigAttachment. >>> Resource test cases, automatically generated.
// Case LoadBalancerAccessLogConfigAttachment_test250103_自动化 9818
func TestAccAliCloudAlbLoadBalancerAccessLogConfigAttachment_basic9818(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_load_balancer_access_log_config_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudAlbLoadBalancerAccessLogConfigAttachmentMap9818)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbLoadBalancerAccessLogConfigAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salb%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlbLoadBalancerAccessLogConfigAttachmentBasicDependence9818)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"log_store":        "${var.name}",
					"load_balancer_id": "${alicloud_alb_load_balancer.defaultDYswYo.id}",
					"log_project":      "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_store":        CHECKSET,
						"load_balancer_id": CHECKSET,
						"log_project":      CHECKSET,
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

var AlicloudAlbLoadBalancerAccessLogConfigAttachmentMap9818 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudAlbLoadBalancerAccessLogConfigAttachmentBasicDependence9818(name string) string {
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

resource "alicloud_alb_load_balancer" "defaultDYswYo" {
  load_balancer_name    = format("%%s4", var.name)
  load_balancer_edition = "Standard"
  vpc_id                = alicloud_vpc.alb_test_tf_vpc.id
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  address_type           = "Intranet"
  address_allocated_mode = "Fixed"
  zone_mappings {
    vswitch_id = alicloud_vswitch.defaultDSY0JJ.id
    zone_id    = alicloud_vswitch.defaultDSY0JJ.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.alb_test_tf_j.id
    zone_id    = alicloud_vswitch.alb_test_tf_j.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.alb_test_tf_k.id
    zone_id    = alicloud_vswitch.alb_test_tf_k.zone_id
  }
  lifecycle {
    ignore_changes = [access_log_config]
  }
}


`, name)
}

// Test Alb LoadBalancerAccessLogConfigAttachment. <<< Resource test cases, automatically generated.
