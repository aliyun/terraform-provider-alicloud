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
	resource.AddTestSweepers("alicloud_amqp_binding", &resource.Sweeper{
		Name: "alicloud_amqp_binding",
		F:    testSweepAmqpBinding,
	})
}

func testSweepAmqpBinding(region string) error {
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
			log.Println(WrapErrorf(err, DataDefaultErrorMsg, "alicloud_amqp_bindings", action, AlibabaCloudSdkGoERROR))
			return nil
		}
		resp, err := jsonpath.Get("$.Data.Instances", response)
		if err != nil {
			log.Println(WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.Instances", response))
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			instanceId := fmt.Sprint(item["InstanceId"])
			action := "ListExchanges"
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
					log.Println(WrapErrorf(err, DataDefaultErrorMsg, "alicloud_amqp_bindings", action, AlibabaCloudSdkGoERROR))
					return nil
				}
				resp, err := jsonpath.Get("$.Data.Exchanges", response)
				if err != nil {
					log.Println(WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.Exchanges", response))
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

					action := "DeleteExchange"
					request := map[string]interface{}{
						"InstanceId": instanceId,
						"Exchange":   item["Name"],
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

func TestAccAlicloudAmqpBinding_all_EXCHANGE(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_binding.default"
	ra := resourceAttrInit(resourceId, AmqpBindingBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpBindingbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpBindingConfigDependence)

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
					"instance_id":       "${alicloud_amqp_exchange.default.instance_id}",
					"virtual_host_name": "${alicloud_amqp_exchange.default.virtual_host_name}",
					"argument":          "x-match:all",
					"binding_key":       "${alicloud_amqp_exchange.default2.exchange_name}",
					"binding_type":      "EXCHANGE",
					"destination_name":  name,
					"source_exchange":   "${alicloud_amqp_exchange.default.exchange_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":       CHECKSET,
						"virtual_host_name": name,
						"argument":          "x-match:all",
						"binding_key":       name + "-2",
						"binding_type":      "EXCHANGE",
						"destination_name":  name,
						"source_exchange":   name,
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

func TestAccAlicloudAmqpBinding_all_QUEUE(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_binding.default"
	ra := resourceAttrInit(resourceId, AmqpBindingBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpBindingbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpBindingConfigDependence)

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
					"instance_id":       "${alicloud_amqp_exchange.default.instance_id}",
					"virtual_host_name": "${alicloud_amqp_exchange.default.virtual_host_name}",
					"argument":          "x-match:all",
					"binding_key":       "${alicloud_amqp_queue.default.queue_name}",
					"binding_type":      "QUEUE",
					"destination_name":  name,
					"source_exchange":   "${alicloud_amqp_exchange.default.exchange_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":       CHECKSET,
						"virtual_host_name": name,
						"argument":          "x-match:all",
						"binding_key":       name,
						"binding_type":      "QUEUE",
						"destination_name":  name,
						"source_exchange":   name,
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

func TestAccAlicloudAmqpBinding_any_EXCHANGE(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_binding.default"
	ra := resourceAttrInit(resourceId, AmqpBindingBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpBindingbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpBindingConfigDependence)

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
					"instance_id":       "${alicloud_amqp_exchange.default.instance_id}",
					"virtual_host_name": "${alicloud_amqp_exchange.default.virtual_host_name}",
					"argument":          "x-match:any",
					"binding_key":       "${alicloud_amqp_exchange.default2.exchange_name}",
					"binding_type":      "EXCHANGE",
					"destination_name":  name,
					"source_exchange":   "${alicloud_amqp_exchange.default.exchange_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":       CHECKSET,
						"virtual_host_name": name,
						"argument":          "x-match:any",
						"binding_key":       name + "-2",
						"binding_type":      "EXCHANGE",
						"destination_name":  name,
						"source_exchange":   name,
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

func TestAccAlicloudAmqpBinding_any_QUEUE(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_binding.default"
	ra := resourceAttrInit(resourceId, AmqpBindingBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpBindingbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpBindingConfigDependence)

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
					"instance_id":       "${alicloud_amqp_exchange.default.instance_id}",
					"virtual_host_name": "${alicloud_amqp_exchange.default.virtual_host_name}",
					"argument":          "x-match:any",
					"binding_key":       "${alicloud_amqp_queue.default.queue_name}",
					"binding_type":      "QUEUE",
					"destination_name":  name,
					"source_exchange":   "${alicloud_amqp_exchange.default.exchange_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":       CHECKSET,
						"virtual_host_name": name,
						"argument":          "x-match:any",
						"binding_key":       name,
						"binding_type":      "QUEUE",
						"destination_name":  name,
						"source_exchange":   name,
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

func TestAccAlicloudAmqpBinding_empty_EXCHANGE(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_binding.default"
	ra := resourceAttrInit(resourceId, AmqpBindingBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpBindingbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpBindingTopicConfigDependence)

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
					"instance_id":       "${alicloud_amqp_exchange.default.instance_id}",
					"virtual_host_name": "${alicloud_amqp_exchange.default.virtual_host_name}",
					"binding_key":       "${alicloud_amqp_exchange.default2.exchange_name}",
					"binding_type":      "EXCHANGE",
					"destination_name":  name,
					"source_exchange":   "${alicloud_amqp_exchange.default.exchange_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":       CHECKSET,
						"virtual_host_name": name,
						"argument":          "",
						"binding_key":       name + "-2",
						"binding_type":      "EXCHANGE",
						"destination_name":  name,
						"source_exchange":   name,
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

func TestAccAlicloudAmqpBinding_empty_QUEUE(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_binding.default"
	ra := resourceAttrInit(resourceId, AmqpBindingBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpBindingbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpBindingTopicConfigDependence)

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
					"instance_id":       "${alicloud_amqp_exchange.default.instance_id}",
					"virtual_host_name": "${alicloud_amqp_exchange.default.virtual_host_name}",
					"binding_key":       "${alicloud_amqp_queue.default.queue_name}",
					"binding_type":      "QUEUE",
					"destination_name":  name,
					"source_exchange":   "${alicloud_amqp_exchange.default.exchange_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":       CHECKSET,
						"virtual_host_name": name,
						"argument":          "",
						"binding_key":       name,
						"binding_type":      "QUEUE",
						"destination_name":  name,
						"source_exchange":   name,
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

func resourceAmqpBindingConfigDependence(name string) string {
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
		resource "alicloud_amqp_exchange" "default" {
			instance_id = alicloud_amqp_virtual_host.default.instance_id
			virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
			auto_delete_state = true
			exchange_name = var.name
			exchange_type = "HEADERS"
			internal = false
		}
		resource "alicloud_amqp_exchange" "default2" {
			instance_id = alicloud_amqp_virtual_host.default.instance_id
			virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
			auto_delete_state = true
			exchange_name = "${var.name}-2"
			exchange_type = "HEADERS"
			internal = false
		}
		resource "alicloud_amqp_queue" "default" {
		  instance_id = alicloud_amqp_virtual_host.default.instance_id
          virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
		  queue_name = var.name
		  auto_delete_state = true
		}
		`, name)
}

func resourceAmqpBindingTopicConfigDependence(name string) string {
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
		resource "alicloud_amqp_exchange" "default" {
			instance_id = alicloud_amqp_virtual_host.default.instance_id
			virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
			auto_delete_state = true
			exchange_name = var.name
			exchange_type = "TOPIC"
			internal = false
		}
		resource "alicloud_amqp_exchange" "default2" {
			instance_id = alicloud_amqp_virtual_host.default.instance_id
			virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
			auto_delete_state = true
			exchange_name = "${var.name}-2"
			exchange_type = "TOPIC"
			internal = false
		}
		resource "alicloud_amqp_queue" "default" {
		  instance_id = alicloud_amqp_virtual_host.default.instance_id
          virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
		  queue_name = var.name
		  auto_delete_state = true
		}
		`, name)
}

var AmqpBindingBasicMap = map[string]string{}
