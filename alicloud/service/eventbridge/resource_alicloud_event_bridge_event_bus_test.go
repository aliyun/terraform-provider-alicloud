package eventbridge_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eventbridge"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	log "github.com/sirupsen/logrus"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_event_bridge_event_bus",
		&resource.Sweeper{
			Name:         "alicloud_event_bridge_event_bus",
			F:            testSweepEventBridgeEventBus,
			Dependencies: []string{"alicloud_event_bridge_event_source"},
		})
}

func testSweepEventBridgeEventBus(region string) error {
	prefixes := []string{
		"tf-testacc",
	}
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	var response map[string]interface{}
	request := make(map[string]interface{})
	action := "ListEventBuses"
	request["Limit"] = helper.PageSizeLarge
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
	if err != nil {
		log.Println(errs.WrapError(fmt.Errorf("ListEventBuses failed, response: %v", response)))
		return nil
	}
	resp, err := jsonpath.Get("$.Data.EventBuses", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data.EventBuses", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		eventBusName := fmt.Sprint(item["EventBusName"])
		skip := true
		if !sweep.SweepAll() {
			for _, prefix := range prefixes {
				if strings.Contains(strings.ToLower(eventBusName), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping EventBridge Bus: %s", eventBusName)
				continue
			}
		}

		action = "ListRules"
		request = map[string]interface{}{
			"Limit":        helper.PageSizeLarge,
			"EventBusName": eventBusName,
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
			log.Println(errs.WrapError(fmt.Errorf("ListRules failed, response: %v", response)))
			continue
		}
		resp, _ = jsonpath.Get("$.Data.Rules", response)
		ruleResult, _ := resp.([]interface{})
		for _, v := range ruleResult {
			ruleName := v.(map[string]interface{})["RuleName"]
			log.Printf("[INFO] Deleting EventBridge Rule: %s", ruleName)
			action = "DeleteRule"
			request = map[string]interface{}{
				"RuleName":     ruleName,
				"EventBusName": eventBusName,
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
				log.Printf("\n[ERROR] Deleting EventBridge rule %s failed. Response: %v. Error: %v.", ruleName, response, err)
			}
		}

		log.Printf("[INFO] Deleting EventBridge Bus: %s", eventBusName)
		action = "DeleteEventBus"
		request = map[string]interface{}{
			"EventBusName": eventBusName,
		}
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
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
			log.Printf("\n[ERROR] Deleting EventBridge bus %s failed. Response: %v. Error: %v.", eventBusName, response, err)
		}
	}
	return nil
}

func TestAccAliCloudEventBridgeEventBus_basic(t *testing.T) {
	var v map[string]interface{}
	ResourceId := "alicloud_event_bridge_event_bus.default"
	ra := acctest.ResourceAttrInit(ResourceId, AlicloudEventBridgeEventBusMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(ResourceId, &v, func() interface{} {
		return eventbridge.NewEventbridgeService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEventBridgeEventBus")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventbridgeeventbus%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(ResourceId, name, AlicloudEventBridgeEventBusBasicDependence0)
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
					"event_bus_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name": name,
					}),
				),
			},
			{
				ResourceName:      ResourceId,
				ImportState:       true,
				ImportStateVerify: true,
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
		},
	})
}

var AlicloudEventBridgeEventBusMap0 = map[string]string{
	"description":    "",
	"event_bus_name": acctest.CHECKSET,
}

func AlicloudEventBridgeEventBusBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}
`, name)
}
