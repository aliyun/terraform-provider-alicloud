package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_amqp_queue", &resource.Sweeper{
		Name: "alicloud_amqp_queue",
		F:    testSweepAmqpQueue,
	})
}

func testSweepAmqpQueue(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}
	action := "ListInstances"
	request := make(map[string]interface{})
	request["MaxResults"] = PageSizeLarge
	var response map[string]interface{}
	conn, err := client.NewOnsproxyClient()
	if err != nil {
		log.Println(WrapError(err))
		return nil
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-12-12"), StringPointer("AK"), request, nil, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Println(WrapErrorf(err, DataDefaultErrorMsg, "alicloud_amqp_queues", action, AlibabaCloudSdkGoERROR))
			return nil
		}
		resp, err := jsonpath.Get("$.Data.Instances", response)
		if err != nil {
			log.Println(WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.Queues", response))
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			instanceId := fmt.Sprint(item["InstanceId"])
			action := "ListQueues"
			request := make(map[string]interface{})
			request["InstanceId"] = instanceId
			request["MaxResults"] = PageSizeLarge
			var response map[string]interface{}
			conn, err := client.NewOnsproxyClient()
			if err != nil {
				log.Println(WrapError(err))
				return nil
			}
			for {
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(5*time.Minute, func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-12-12"), StringPointer("AK"), request, nil, &runtime)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					log.Println(WrapErrorf(err, DataDefaultErrorMsg, "alicloud_amqp_queues", action, AlibabaCloudSdkGoERROR))
					return nil
				}
				resp, err := jsonpath.Get("$.Data.Queues", response)
				if err != nil {
					log.Println(WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.Queues", response))
					return nil
				}
				result, _ := resp.([]interface{})
				for _, v := range result {
					item := v.(map[string]interface{})
					skip := true
					for _, prefixe := range prefixes {
						if strings.HasPrefix(fmt.Sprint(item["Name"]), prefixe) {
							skip = false
							break
						}
					}
					if skip {
						log.Printf("[DEBUG] Skipping the resource %s", item["Name"])
					}

					action := "DeleteQueue"
					request := map[string]interface{}{
						"InstanceId": instanceId,
						"Queue":      item["Name"],
					}

					wait := incrementalWait(3*time.Second, 3*time.Second)
					err = resource.Retry(3*time.Minute, func() *resource.RetryError {
						response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
						if err != nil {
							if NeedRetry(err) {
								wait()
								return resource.RetryableError(err)
							}
							return resource.NonRetryableError(err)
						}
						return nil
					})
					log.Println(WrapError(err))
				}
				if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
					request["NextToken"] = nextToken
				} else {
					break
				}
			}
		}
	}

	return nil
}

func TestAccAlicloudAmqpQueue_basic(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_queue.default"
	ra := resourceAttrInit(resourceId, AmqpQueueBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpQueuebasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpQueueConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":             "${alicloud_amqp_virtual_host.default.instance_id}",
					"virtual_host_name":       "${alicloud_amqp_virtual_host.default.virtual_host_name}",
					"auto_delete_state":       "true",
					"auto_expire_state":       "10000",
					"dead_letter_exchange":    "",
					"dead_letter_routing_key": "",
					"exclusive_state":         "false",
					"max_length":              "100",
					"maximum_priority":        "10",
					"message_ttl":             "100",
					"queue_name":              name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":             CHECKSET,
						"virtual_host_name":       name,
						"auto_delete_state":       "true",
						"auto_expire_state":       "10000",
						"dead_letter_exchange":    "",
						"dead_letter_routing_key": "",
						"exclusive_state":         "false",
						"max_length":              "100",
						"maximum_priority":        "10",
						"message_ttl":             "100",
						"queue_name":              name,
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_expire_state", "dead_letter_exchange", "dead_letter_routing_key", "max_length", "maximum_priority", "message_ttl"},
			},
		},
	})

}

func resourceAmqpQueueConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
 			default = "%v"
		}
		data "alicloud_amqp_instances" "default" {
			status = "SERVING"
		}
		resource "alicloud_amqp_virtual_host" "default" {
		  instance_id       = data.alicloud_amqp_instances.default.ids.0
		  virtual_host_name = var.name
		}
		`, name)
}

var AmqpQueueBasicMap = map[string]string{}
