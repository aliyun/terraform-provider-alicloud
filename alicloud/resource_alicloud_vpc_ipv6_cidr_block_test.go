package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudVPCIpv6CidrBlock_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipv6_cidr_block.default"
	checkoutSupportedRegions(t, true, connectivity.VPCSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudVPCIpv6CidrBlockMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpv6CidrBlock")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipv6cidrblock%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCIpv6CidrBlockBasicDependence0)
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
					"secondary_cidr_block": "tf-testAcc-Tp63T",
					"vpc_id":               "tf-testAcc-uiYus",
					"ipv6_isp":             "ChinaMobile",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_cidr_block": "tf-testAcc-Tp63T",
						"vpc_id":               "tf-testAcc-uiYus",
						"ipv6_isp":             "ChinaMobile",
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

var AlicloudVPCIpv6CidrBlockMap0 = map[string]string{
	"ipv6_isp":             CHECKSET,
	"vpc_id":               CHECKSET,
	"secondary_cidr_block": CHECKSET,
}

func AlicloudVPCIpv6CidrBlockBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
