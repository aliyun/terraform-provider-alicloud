package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudEventBridgeEventSource_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_source.default"
	ra := resourceAttrInit(resourceId, AlicloudEventBridgeEventSourceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventbridgeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventbridgeeventsource%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEventBridgeEventSourceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EventBridgeSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"event_bus_name":    "${var.name}",
					"event_source_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":    name,
						"event_source_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"linked_external_source": "true",
					"external_source_type":   "MNS",
					"external_source_config": map[string]interface{}{
						"QueueName": "${alicloud_mns_queue.queue2.name}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"linked_external_source":           "true",
						"external_source_type":             "MNS",
						"external_source_config.%":         "1",
						"external_source_config.QueueName": name + "change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"linked_external_source": "true",
					"external_source_type":   "MNS",
					"external_source_config": map[string]interface{}{
						"QueueName": "${alicloud_mns_queue.queue1.name}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_source_config.%":         "1",
						"external_source_config.QueueName": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
					"external_source_config": map[string]interface{}{
						"QueueName": "${alicloud_mns_queue.queue2.name}",
					},
					"external_source_type":   "MNS",
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                      name,
						"external_source_config.%":         "1",
						"external_source_config.QueueName": name + "change",
						"external_source_type":             "MNS",
						"linked_external_source":           "true",
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

var AlicloudEventBridgeEventSourceMap0 = map[string]string{
	"event_bus_name":           CHECKSET,
	"event_source_name":        CHECKSET,
	"external_source_config.%": "0",
	"external_source_type":     "",
	"linked_external_source":   CHECKSET,
	"description":              "",
}

func AlicloudEventBridgeEventSourceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

resource "alicloud_event_bridge_event_bus" "default" {
	event_bus_name = var.name
}

resource "alicloud_mns_queue" "queue1" {
  name = var.name
}

resource "alicloud_mns_queue" "queue2" {
  name = format("%%schange", var.name)
}
`, name)
}
