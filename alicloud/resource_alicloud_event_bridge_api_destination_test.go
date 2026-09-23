package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudEventBridgeApiDestination_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.EventBridgeConnectionSupportRegions)
	resourceId := "alicloud_event_bridge_api_destination.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeApiDestinationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeApiDestination")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventbridgeapidestination%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeApiDestinationBasicDependence0)
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
					"connection_name":      "${alicloud_event_bridge_connection.default.connection_name}",
					"api_destination_name": name,
					"http_api_parameters": []map[string]interface{}{
						{
							"endpoint": "http://127.0.0.1:8001",
							"method":   "POST",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_name":       CHECKSET,
						"api_destination_name":  name,
						"http_api_parameters.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-api-destination-connection",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-api-destination-connection",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"http_api_parameters": []map[string]interface{}{
						{
							"endpoint": "http://127.0.0.1:8002",
							"method":   "GET",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_api_parameters.#": "1",
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

func TestAccAliCloudEventBridgeApiDestination_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.EventBridgeConnectionSupportRegions)
	resourceId := "alicloud_event_bridge_api_destination.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeApiDestinationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeApiDestination")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventbridgeapidestination%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeApiDestinationBasicDependence0)
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
					"connection_name":      "${alicloud_event_bridge_connection.default.connection_name}",
					"api_destination_name": name,
					"description":          "test-api-destination-connection",
					"http_api_parameters": []map[string]interface{}{
						{
							"endpoint": "http://127.0.0.1:8001",
							"method":   "POST",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_name":       CHECKSET,
						"api_destination_name":  name,
						"description":           "test-api-destination-connection",
						"http_api_parameters.#": "1",
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

var AliCloudEventBridgeApiDestinationMap0 = map[string]string{
	"create_time": CHECKSET,
}

func AliCloudEventBridgeApiDestinationBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
    	default = "%s"
	}

	resource "alicloud_event_bridge_connection" "default" {
  		connection_name = var.name
  		network_parameters {
    		network_type = "PublicNetwork"
  		}
	}
`, name)
}
