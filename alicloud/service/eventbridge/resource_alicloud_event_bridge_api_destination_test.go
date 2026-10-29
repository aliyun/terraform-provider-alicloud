package eventbridge_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eventbridge"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudEventBridgeApiDestination_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.EventBridgeConnectionSupportRegions)
	ResourceId := "alicloud_event_bridge_api_destination.default"
	ra := acctest.ResourceAttrInit(ResourceId, AliCloudEventBridgeApiDestinationMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(ResourceId, &v, func() interface{} {
		return eventbridge.NewEventBridgeServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEventBridgeApiDestination")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventbridgeapidestination%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(ResourceId, name, AliCloudEventBridgeApiDestinationBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"connection_name":       acctest.CHECKSET,
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
				ResourceName:      ResourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudEventBridgeApiDestination_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.EventBridgeConnectionSupportRegions)
	ResourceId := "alicloud_event_bridge_api_destination.default"
	ra := acctest.ResourceAttrInit(ResourceId, AliCloudEventBridgeApiDestinationMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(ResourceId, &v, func() interface{} {
		return eventbridge.NewEventBridgeServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEventBridgeApiDestination")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventbridgeapidestination%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(ResourceId, name, AliCloudEventBridgeApiDestinationBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"connection_name":       acctest.CHECKSET,
						"api_destination_name":  name,
						"description":           "test-api-destination-connection",
						"http_api_parameters.#": "1",
					}),
				),
			},
			{
				ResourceName:      ResourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AliCloudEventBridgeApiDestinationMap0 = map[string]string{
	"create_time": acctest.CHECKSET,
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
