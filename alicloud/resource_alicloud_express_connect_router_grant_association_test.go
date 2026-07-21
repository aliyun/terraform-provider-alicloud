package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ExpressConnectRouter GrantAssociation. >>> Resource test cases, automatically generated.
// Case 预发测试用例 9639
func TestAccAliCloudExpressConnectRouterGrantAssociation_basic9639(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_grant_association.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectRouterGrantAssociationMap9639)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectRouterServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectRouterGrantAssociation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectroutergrantassociation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectRouterGrantAssociationBasicDependence9639)
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
					"instance_region_id": "cn-wulanchabu-test-5",
					"instance_id":        "vpc-7qbkkicdvc69hzhu48x98",
					"instance_type":      "VPC",
					"ecr_id":             "${var.ecr_id}",
					"ecr_owner_ali_uid":  "${var.ecr_uid}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecr_id":             CHECKSET,
						"instance_region_id": CHECKSET,
						"instance_id":        CHECKSET,
						"ecr_owner_ali_uid":  CHECKSET,
						"instance_type":      "VPC",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudExpressConnectRouterGrantAssociationMap9639 = map[string]string{
	"status": CHECKSET,
}

func AlicloudExpressConnectRouterGrantAssociationBasicDependence9639(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "ecr_id" {
  default = "ecr-g2g2f4ww1vduya9xcn"
}

variable "region" {
  default = "cn-hangzhou"
}

variable "ecr_uid" {
  default = 1891593620094065
}


`, name)
}

// Test ExpressConnectRouter GrantAssociation. <<< Resource test cases, automatically generated.
