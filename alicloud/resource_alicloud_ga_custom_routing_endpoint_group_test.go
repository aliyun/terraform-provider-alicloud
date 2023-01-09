package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGaCustomRoutingEndpointGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_custom_routing_endpoint_group.default"
	ra := resourceAttrInit(resourceId, AlicloudGaCustomRoutingEndpointGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaCustomRoutingEndpointGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sGaCustomRoutingEndpointGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaCustomRoutingEndpointGroupBasicDependence0)
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
					"accelerator_id":                     "${alicloud_ga_listener.default.accelerator_id}",
					"listener_id":                        "${alicloud_ga_listener.default.id}",
					"endpoint_group_region":              defaultRegionToTest,
					"custom_routing_endpoint_group_name": name,
					"description":                        name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":                     CHECKSET,
						"listener_id":                        CHECKSET,
						"endpoint_group_region":              CHECKSET,
						"custom_routing_endpoint_group_name": name,
						"description":                        name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_routing_endpoint_group_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_routing_endpoint_group_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "-update",
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

var AlicloudGaCustomRoutingEndpointGroupMap = map[string]string{
	"status": CHECKSET,
}

func AlicloudGaCustomRoutingEndpointGroupBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_ga_accelerators" "default" {
  		status = "active"
	}
	
	resource "alicloud_ga_bandwidth_package" "default" {
		bandwidth      = 100
  		type           = "Basic"
  		bandwidth_type = "Basic"
  		payment_type   = "PayAsYouGo"
  		billing_type   = "PayBy95"
  		ratio          = 30
	}
	
	resource "alicloud_ga_bandwidth_package_attachment" "default" {
  		accelerator_id       = data.alicloud_ga_accelerators.default.accelerators.0.id
  		bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
	}
	
	resource "alicloud_ga_listener" "default" {
  		accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
		listener_type  = "CustomRouting"
  		port_ranges {
    		from_port = 10000
    		to_port   = 16000
  		}
	}
`, name)
}
