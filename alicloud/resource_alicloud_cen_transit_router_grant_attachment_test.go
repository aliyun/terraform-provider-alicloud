package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func SkipTestAccAlicloudCENTransitRouterGrantAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_grant_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterGrantAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterGrantAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%scentransitroutergrantattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenTransitRouterGrantAttachmentBasicDependence0)
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
					"order_type":    "PayByCenOwner",
					"instance_id":   "${data.alicloud_vpcs.default.ids.0}",
					"cen_owner_id":  "${var.cen_owner_id}",
					"cen_id":        "${alicloud_cen_instance.default.id}",
					"instance_type": "VPC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_type":    "PayByCenOwner",
						"instance_id":   CHECKSET,
						"cen_owner_id":  CHECKSET,
						"cen_id":        CHECKSET,
						"instance_type": "VPC",
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

var AlicloudCenTransitRouterGrantAttachmentMap0 = map[string]string{}

func AlicloudCenTransitRouterGrantAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "cen_owner_id" {
  default = "%s"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  description       = "test for transit router grant attachment"
}

`, name, os.Getenv("ALICLOUD_MAIN_ACCOUNT"))
}
