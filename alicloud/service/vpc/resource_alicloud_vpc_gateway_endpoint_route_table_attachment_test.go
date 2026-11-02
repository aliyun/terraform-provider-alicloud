package vpc_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Vpc GatewayEndpointRouteTableAttachment. >>> Resource test cases, automatically generated.
// Case 3634
func TestAccAliCloudVpcGatewayEndpointRouteTableAttachment_basic3634(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_gateway_endpoint_route_table_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudVpcGatewayEndpointRouteTableAttachmentMap3634)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return NewVpcServiceV2(testAccProvider.Meta().(*connectivity.AliyunClient))
	}, "DescribeVpcGatewayEndpointRouteTableAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcgatewayendpointroutetableattachment%d", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudVpcGatewayEndpointRouteTableAttachmentBasicDependence3634)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VPCGatewayEndpointSupportRegions)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_endpoint_id": "${alicloud_vpc_gateway_endpoint.defaultGE.id}",
					"route_table_id":      "${alicloud_route_table.defaultRT.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_endpoint_id": acctest.CHECKSET,
						"route_table_id":      acctest.CHECKSET,
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

var AlicloudVpcGatewayEndpointRouteTableAttachmentMap3634 = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudVpcGatewayEndpointRouteTableAttachmentBasicDependence3634(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaulteVpc" {
  description = "test"
}

resource "alicloud_vpc_gateway_endpoint" "defaultGE" {
  service_name                = "com.aliyun.cn-hangzhou.oss"
  policy_document             = "{ \"Version\" : \"1\", \"Statement\" : [ { \"Effect\" : \"Allow\", \"Resource\" : [ \"*\" ], \"Action\" : [ \"*\" ], \"Principal\" : [ \"*\" ] } ] }"
  vpc_id                      = alicloud_vpc.defaulteVpc.id
  gateway_endpoint_descrption = "test-gateway-endpoint"
  gateway_endpoint_name       = "${var.name}1"
}

resource "alicloud_route_table" "defaultRT" {
  vpc_id           = alicloud_vpc.defaulteVpc.id
  route_table_name = "${var.name}2"
}

resource "alicloud_route_table" "secondRT" {
  vpc_id           = alicloud_vpc.defaulteVpc.id
  route_table_name = "${var.name}3"
}

resource "alicloud_vpc_gateway_endpoint_route_table_attachment" "default0" {
  gateway_endpoint_id = alicloud_vpc_gateway_endpoint.defaultGE.id
  route_table_id      = alicloud_route_table.secondRT.id
}

`, name)
}

// Test Vpc GatewayEndpointRouteTableAttachment. <<< Resource test cases, automatically generated.
