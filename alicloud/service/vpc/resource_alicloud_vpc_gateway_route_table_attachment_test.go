package vpc_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudVPCVpcGatewayRouteTableAttachment_basic0(t *testing.T) {
	t.Skip("Skipping: depends on pre-existing default VPC (default-NODELETING) which may not exist in all test regions")
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.VPCSupportRegions)
	resourceId := "alicloud_vpc_gateway_route_table_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudVpcGatewayRouteTableAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return NewVpcServiceV2(testAccProvider.Meta().(*connectivity.AliyunClient))
	}, "DescribeVpcGatewayRouteTableAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcgatewayroutetable%d", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudVpcGatewayRouteTableAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":  "${alicloud_route_table.default.id}",
					"ipv4_gateway_id": "${alicloud_vpc_ipv4_gateway.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":  acctest.CHECKSET,
						"ipv4_gateway_id": acctest.CHECKSET,
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

var AlicloudVpcGatewayRouteTableAttachmentMap0 = map[string]string{}

func AlicloudVpcGatewayRouteTableAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alicloud_route_table" "default" {
  vpc_id           = "${data.alicloud_vpcs.default.ids.0}"
  route_table_name = "${var.name}"
  description      = "${var.name}_description"
  associate_type   = "Gateway"
}

resource "alicloud_vpc_ipv4_gateway" "default" {
  ipv4_gateway_description = var.name
  ipv4_gateway_name        = var.name
  vpc_id                   = "${data.alicloud_vpcs.default.ids.0}"
  enabled                  = "true"
}

`, name)
}
