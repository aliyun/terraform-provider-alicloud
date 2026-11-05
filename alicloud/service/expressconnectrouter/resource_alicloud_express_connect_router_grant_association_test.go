package expressconnectrouter_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/expressconnectrouter"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ExpressConnectRouter GrantAssociation. >>> Resource test cases, automatically generated.
// Case 预发测试用例 9639
func TestAccAliCloudExpressConnectRouterGrantAssociation_basic9639(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_grant_association.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterGrantAssociationMap9639)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterGrantAssociation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectroutergrantassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterGrantAssociationBasicDependence9639)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// The pre-provisioned ECR (ecr-g2g2f4ww1vduya9xcn) lives in cn-hangzhou.
			// Force cn-hangzhou via TestSalveRegions so GrantInstanceToExpressConnectRouter
			// resolves the ECR in its home region, matching the data source test.
			acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"ecr_id":             acctest.CHECKSET,
						"instance_region_id": acctest.CHECKSET,
						"instance_id":        acctest.CHECKSET,
						"ecr_owner_ali_uid":  acctest.CHECKSET,
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
	"status": acctest.CHECKSET,
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
