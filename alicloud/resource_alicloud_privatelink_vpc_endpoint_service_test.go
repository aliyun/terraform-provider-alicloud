package alicloud

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudPrivatelinkVpcEndpointService_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_service.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointServiceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivatelinkService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivatelinkVpcEndpointService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", AlicloudPrivatelinkVpcEndpointServiceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.PrivateLinkRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_description":    "ThisismyEndpointService",
					"connect_bandwidth":      "103",
					"auto_accept_connection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description":    "ThisismyEndpointService",
						"connect_bandwidth":      "103",
						"auto_accept_connection": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "payer"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_accept_connection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_accept_connection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_description": "ThisismyEndpointServiceUpdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description": "ThisismyEndpointServiceUpdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connect_bandwidth": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connect_bandwidth": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_accept_connection": "false",
					"service_description":    "ThisismyEndpointService",
					"connect_bandwidth":      "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_accept_connection": "false",
						"service_description":    "ThisismyEndpointService",
						"connect_bandwidth":      "200",
					}),
				),
			},
		},
	})
}

var AlicloudPrivatelinkVpcEndpointServiceMap = map[string]string{
	"service_business_status": "Normal",
	"service_domain":          CHECKSET,
	"status":                  CHECKSET,
}

func AlicloudPrivatelinkVpcEndpointServiceBasicDependence(name string) string {
	return ""
}
