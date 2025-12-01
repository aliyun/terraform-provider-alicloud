package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudAmqpBinding_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudAmqpBindingMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-%s-AmqpBinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAmqpBindingBasicDependence0)
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
					"instance_id":       "${alicloud_amqp_exchange.source.instance_id}",
					"virtual_host_name": "${alicloud_amqp_exchange.source.virtual_host_name}",
					"source_exchange":   "${alicloud_amqp_exchange.source.exchange_name}",
					"destination_name":  "${alicloud_amqp_exchange.destination.exchange_name}",
					"binding_type":      "EXCHANGE",
					"binding_key":       "${alicloud_amqp_exchange.destination.exchange_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":       CHECKSET,
						"virtual_host_name": CHECKSET,
						"source_exchange":   CHECKSET,
						"destination_name":  CHECKSET,
						"binding_type":      "EXCHANGE",
						"binding_key":       CHECKSET,
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

func TestAccAliCloudAmqpBinding_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudAmqpBindingMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-%s-AmqpBinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAmqpBindingBasicDependence0)
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
					"instance_id":       "${alicloud_amqp_exchange.source.instance_id}",
					"virtual_host_name": "${alicloud_amqp_exchange.source.virtual_host_name}",
					"source_exchange":   "${alicloud_amqp_exchange.source.exchange_name}",
					"destination_name":  "${alicloud_amqp_exchange.destination.exchange_name}",
					"binding_type":      "EXCHANGE",
					"binding_key":       "${alicloud_amqp_exchange.destination.exchange_name}",
					"argument":          "x-match:any",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":       CHECKSET,
						"virtual_host_name": CHECKSET,
						"source_exchange":   CHECKSET,
						"destination_name":  CHECKSET,
						"binding_type":      "EXCHANGE",
						"binding_key":       CHECKSET,
						"argument":          "x-match:any",
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

func TestAccAliCloudAmqpBinding_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudAmqpBindingMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-%s-AmqpBinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAmqpBindingBasicDependence0)
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
					"instance_id":       "${alicloud_amqp_exchange.source.instance_id}",
					"virtual_host_name": "${alicloud_amqp_exchange.source.virtual_host_name}",
					"source_exchange":   "${alicloud_amqp_exchange.source.exchange_name}",
					"destination_name":  "${alicloud_amqp_queue.default.queue_name}",
					"binding_type":      "QUEUE",
					"binding_key":       "${alicloud_amqp_queue.default.queue_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":       CHECKSET,
						"virtual_host_name": CHECKSET,
						"source_exchange":   CHECKSET,
						"destination_name":  CHECKSET,
						"binding_type":      "QUEUE",
						"binding_key":       CHECKSET,
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

func TestAccAliCloudAmqpBinding_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_binding.default"
	ra := resourceAttrInit(resourceId, AliCloudAmqpBindingMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-%s-AmqpBinding%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAmqpBindingBasicDependence0)
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
					"instance_id":       "${alicloud_amqp_exchange.source.instance_id}",
					"virtual_host_name": "${alicloud_amqp_exchange.source.virtual_host_name}",
					"source_exchange":   "${alicloud_amqp_exchange.source.exchange_name}",
					"destination_name":  "${alicloud_amqp_queue.default.queue_name}",
					"binding_type":      "QUEUE",
					"binding_key":       "${alicloud_amqp_queue.default.queue_name}",
					"argument":          "x-match:any",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":       CHECKSET,
						"virtual_host_name": CHECKSET,
						"source_exchange":   CHECKSET,
						"destination_name":  CHECKSET,
						"binding_type":      "QUEUE",
						"binding_key":       CHECKSET,
						"argument":          "x-match:any",
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

var AliCloudAmqpBindingMap0 = map[string]string{
	"argument": CHECKSET,
}

func AliCloudAmqpBindingBasicDependence0(name string) string {
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

	resource "alicloud_amqp_exchange" "source" {
  		instance_id       = alicloud_amqp_virtual_host.default.instance_id
  		virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
  		exchange_name     = var.name
  		exchange_type     = "HEADERS"
  		auto_delete_state = true
  		internal          = false
	}

	resource "alicloud_amqp_exchange" "destination" {
  		instance_id       = alicloud_amqp_virtual_host.default.instance_id
  		virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
  		exchange_name     = "${var.name}-destination"
  		exchange_type     = "HEADERS"
  		auto_delete_state = true
  		internal          = false
	}

	resource "alicloud_amqp_queue" "default" {
  		instance_id       = alicloud_amqp_virtual_host.default.instance_id
  		virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
  		queue_name        = var.name
  		auto_delete_state = true
	}
`, name)
}

func TestUnitAliCloudAmqpBinding(t *testing.T) {
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
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_amqp_binding", "MockSourceExchange"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
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
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudAmqpBindingCreate(d, rawClient)
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
		err := resourceAliCloudAmqpBindingCreate(d, rawClient)
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
		err := resourceAliCloudAmqpBindingCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("instance_id", ":", "virtual_host_name", ":", "MockSourceExchange", ":", "destination_name"))

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewOnsproxyClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudAmqpBindingDelete(d, rawClient)
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
		err := resourceAliCloudAmqpBindingDelete(d, rawClient)
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
		err := resourceAliCloudAmqpBindingDelete(d, rawClient)
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
		err := resourceAliCloudAmqpBindingRead(d, rawClient)
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
		err := resourceAliCloudAmqpBindingRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}
