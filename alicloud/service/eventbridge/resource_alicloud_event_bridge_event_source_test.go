package eventbridge_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eventbridge"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_event_bridge_event_source",
		&resource.Sweeper{
			Name: "alicloud_event_bridge_event_source",
			F:    testSweepEventBridgeEventSource,
		})
}

func testSweepEventBridgeEventSource(region string) error {
	prefixes := []string{
		"tf-testacc",
	}
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "Error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ListUserDefinedEventSources"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"Limit":    helper.PageSizeLarge,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapError(fmt.Errorf("ListUserDefinedEventSources failed. Response: %v. Error: %v", response, err))
	}
	resp, err := jsonpath.Get("$.Data.EventSourceList", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data.EventSourceList", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		skip := true
		if !sweep.SweepAll() {
			for _, prefix := range prefixes {
				EventBusName := ""
				if val, exist := item["Name"]; exist {
					EventBusName = val.(string)
				}
				if strings.Contains(strings.ToLower(EventBusName), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping EventBridge Source: %s", item["Name"])
				continue
			}
		}
		log.Printf("[INFO] Deleting EventBridge Source: %s", item["Name"])
		action = "DeleteEventSource"
		request = map[string]interface{}{
			"EventSourceName": item["Name"],
		}
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("\n[ERROR] Deleting EventBridge source %s failed. Response: %v. Error: %v.", item["Name"], response, err)
		}
	}
	return nil
}

func TestAccAliCloudEventBridgeEventSource_basic0(t *testing.T) {
	var v map[string]interface{}
	ResourceId := "alicloud_event_bridge_event_source.default"
	ra := acctest.ResourceAttrInit(ResourceId, AliCloudEventBridgeEventSourceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(ResourceId, &v, func() interface{} {
		return eventbridge.NewEventbridgeService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEventBridgeEventSource")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventbridgeeventsource%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(ResourceId, name, AliCloudEventBridgeEventSourceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EventBridgeSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"event_bus_name":    "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":    acctest.CHECKSET,
						"event_source_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"external_source_type": "MNS",
					"external_source_config": map[string]interface{}{
						"QueueName": "${alicloud_mns_queue.default.name}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_source_type":             "MNS",
						"external_source_config.%":         "1",
						"external_source_config.QueueName": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"linked_external_source": "true",
					}),
				),
			},
			{
				ResourceName:            ResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

func TestAccAliCloudEventBridgeEventSource_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	ResourceId := "alicloud_event_bridge_event_source.default"
	ra := acctest.ResourceAttrInit(ResourceId, AliCloudEventBridgeEventSourceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(ResourceId, &v, func() interface{} {
		return eventbridge.NewEventbridgeService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEventBridgeEventSource")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventbridgeeventsource%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(ResourceId, name, AliCloudEventBridgeEventSourceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.EventBridgeSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"event_bus_name":       "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name":    name,
					"external_source_type": "MNS",
					"external_source_config": map[string]interface{}{
						"QueueName": "${alicloud_mns_queue.default.name}",
					},
					"description":            name,
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":                   acctest.CHECKSET,
						"event_source_name":                name,
						"external_source_type":             "MNS",
						"external_source_config.%":         "1",
						"external_source_config.QueueName": name,
						"description":                      name,
					}),
				),
			},
			{
				ResourceName:            ResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

var AliCloudEventBridgeEventSourceMap0 = map[string]string{}

func AliCloudEventBridgeEventSourceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_event_bridge_event_bus" "default" {
  		event_bus_name = var.name
	}

	resource "alicloud_mns_queue" "default" {
  		name = var.name
	}
`, name)
}
