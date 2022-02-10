package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/alibabacloud-go/tea-rpc/client"
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
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
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
			err = resource.Retry(1*time.Minute, func() *resource.RetryError {
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
				err = resource.Retry(1*time.Minute, func() *resource.RetryError {
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

func TestAccAlicloudAmqpBinding_unit(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_amqp_binding"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_amqp_binding"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"instance_id":       "instance_id",
		"virtual_host_name": "virtual_host_name",
		"argument":          "x-match:all",
		"binding_key":       "binding_key",
		"binding_type":      "EXCHANGE",
		"destination_name":  "destination_name",
		"source_exchange":   "MockSourceExchange",
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"Data": map[string]interface{}{
			"Bindings": []interface{}{
				map[string]interface{}{
					"DestinationName": "destination_name",
					"InstanceId":      "instance_id",
					"SourceExchange":  "MockSourceExchange",
					"VirtualHostName": "virtual_host_name",
					"Argument":        "x-match:all",
					"BindingKey":      "binding_key",
					"BindingType":     "EXCHANGE",
				},
			},
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:    String(errorCode),
				Data:    String(errorCode),
				Message: String(errorCode),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_amqp_binding", "MockSourceExchange"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:    String(errorCode),
				Data:    String(errorCode),
				Message: String(errorCode),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			result["SourceExchange"] = "MockSourceExchange"
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewOnsproxyClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudAmqpBindingCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudAmqpBindingCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudAmqpBindingCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("instance_id", ":", "virtual_host_name", ":", "MockSourceExchange", ":", "destination_name"))

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewOnsproxyClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudAmqpBindingDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudAmqpBindingDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudAmqpBindingDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeAmqpBindingNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudAmqpBindingRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeAmqpBindingAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := true
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudAmqpBindingRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}
