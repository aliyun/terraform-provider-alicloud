package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudPrivatelinkVpcEndpointConnection_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointConnectionMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivatelinkService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivatelinkVpcEndpointConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", AlicloudPrivatelinkVpcEndpointConnectionBasicDependence)
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
					"endpoint_id": "${alicloud_privatelink_vpc_endpoint.default.id}",
					"service_id":  "${alicloud_privatelink_vpc_endpoint_service.default.id}",
					"bandwidth":   "1024",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_id": CHECKSET,
						"service_id":  CHECKSET,
						"bandwidth":   "1024",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "1000",
					}),
				),
			},
		},
	})
}

var AlicloudPrivatelinkVpcEndpointConnectionMap = map[string]string{
	"status": CHECKSET,
}

func AlicloudPrivatelinkVpcEndpointConnectionBasicDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_vpcs" "default" {
	 name_regex = "default-NODELETING"
	}
	resource "alicloud_security_group" "default" {
	 name = "tf-testAcc-for-privatelink"
	 description = "privatelink test security group"
	 vpc_id = data.alicloud_vpcs.default.ids.0
	}
	resource "alicloud_privatelink_vpc_endpoint_service" "default" {
	service_description = "test for privatelink connection"
	connect_bandwidth = 103
	auto_accept_connection = false
	}
	resource "alicloud_privatelink_vpc_endpoint" "default" {
	 service_id = alicloud_privatelink_vpc_endpoint_service.default.id
	 vpc_id = data.alicloud_vpcs.default.ids.0
	 security_group_ids = [alicloud_security_group.default.id]
	 vpc_endpoint_name = "testformaintf"
	 depends_on = [alicloud_privatelink_vpc_endpoint_service.default]
	}
`)
}
