package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudAmqpExchange_DIRECT(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AmqpExchangeBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpExchangebasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpExchangeConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alternate_exchange": name + "alternate_exchange",
					"instance_id":        "${alicloud_amqp_virtual_host.default.instance_id}",
					"virtual_host_name":  "${alicloud_amqp_virtual_host.default.virtual_host_name}",
					"auto_delete_state":  "true",
					"exchange_name":      name,
					"exchange_type":      "DIRECT",
					"internal":           "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":        CHECKSET,
						"virtual_host_name":  name,
						"alternate_exchange": name + "alternate_exchange",
						"auto_delete_state":  "true",
						"exchange_name":      name,
						"exchange_type":      "DIRECT",
						"internal":           "false",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal"},
			},
		},
	})
}

func TestAccAliCloudAmqpExchange_FANOUT(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AmqpExchangeBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpExchangebasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpExchangeConfigDependence)

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
					"alternate_exchange": name + "alternate_exchange",
					"instance_id":        "${alicloud_amqp_virtual_host.default.instance_id}",
					"virtual_host_name":  "${alicloud_amqp_virtual_host.default.virtual_host_name}",
					"auto_delete_state":  "true",
					"exchange_name":      name,
					"exchange_type":      "FANOUT",
					"internal":           "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":        CHECKSET,
						"virtual_host_name":  name,
						"alternate_exchange": name + "alternate_exchange",
						"auto_delete_state":  "true",
						"exchange_name":      name,
						"exchange_type":      "FANOUT",
						"internal":           "false",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal"},
			},
		},
	})
}

func TestAccAliCloudAmqpExchange_HEADERS(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AmqpExchangeBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpExchangebasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpExchangeConfigDependence)

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
					"alternate_exchange": name + "alternate_exchange",
					"instance_id":        "${alicloud_amqp_virtual_host.default.instance_id}",
					"virtual_host_name":  "${alicloud_amqp_virtual_host.default.virtual_host_name}",
					"auto_delete_state":  "true",
					"exchange_name":      name,
					"exchange_type":      "HEADERS",
					"internal":           "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":        CHECKSET,
						"virtual_host_name":  name,
						"alternate_exchange": name + "alternate_exchange",
						"auto_delete_state":  "true",
						"exchange_name":      name,
						"exchange_type":      "HEADERS",
						"internal":           "false",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal"},
			},
		},
	})
}

func TestAccAliCloudAmqpExchange_TOPIC(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AmqpExchangeBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpExchangebasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpExchangeConfigDependence)

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
					"alternate_exchange": name + "alternate_exchange",
					"instance_id":        "${alicloud_amqp_virtual_host.default.instance_id}",
					"virtual_host_name":  "${alicloud_amqp_virtual_host.default.virtual_host_name}",
					"auto_delete_state":  "true",
					"exchange_name":      name,
					"exchange_type":      "TOPIC",
					"internal":           "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":        CHECKSET,
						"virtual_host_name":  name,
						"alternate_exchange": name + "alternate_exchange",
						"auto_delete_state":  "true",
						"exchange_name":      name,
						"exchange_type":      "TOPIC",
						"internal":           "false",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal"},
			},
		},
	})
}

func resourceAmqpExchangeConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
 			default = "%v"
		}
		resource "alicloud_amqp_instance" "default" {
		  instance_name  = var.name
		  instance_type  = "enterprise"
		  max_tps        = "3000"
		  queue_capacity = "200"
		  period_cycle   = "Year"
		  support_eip    = "false"
		  period         = "1"
		  auto_renew     = "true"
		  payment_type   = "Subscription"
		}
		resource "alicloud_amqp_virtual_host" "default" {
		  instance_id       = alicloud_amqp_instance.default.id
		  virtual_host_name = var.name
		}
		`, name)
}

var AmqpExchangeBasicMap = map[string]string{}

func TestUnitAlicloudAmqpExchange(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_amqp_exchange"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_amqp_exchange"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"alternate_exchange": "alternate_exchange",
		"instance_id":        "instance_id",
		"virtual_host_name":  "virtual_host_name",
		"auto_delete_state":  true,
		"exchange_name":      "exchange_name",
		"exchange_type":      "DIRECT",
		"internal":           false,
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
			"Exchanges": []interface{}{
				map[string]interface{}{
					"ExchangeName":    "exchange_name",
					"InstanceId":      "instance_id",
					"VirtualHostName": "virtual_host_name",
					"ExclusiveState":  false,
					"ExchangeType":    "exchange_type",
					"Name":            "MockExchangeName",
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
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_amqp_exchange", "MockExchangeName"))
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
		err := resourceAliCloudAmqpExchangeCreate(d, rawClient)
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
		err := resourceAliCloudAmqpExchangeCreate(d, rawClient)
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
		err := resourceAliCloudAmqpExchangeCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("instance_id", ":", "virtual_host_name", ":", "MockExchangeName"))

	//// Update
	t.Run("UpdateNormal", func(t *testing.T) {
		patcheDescribeBackups := gomonkey.ApplyMethod(reflect.TypeOf(&AmqpOpenService{}), "DescribeAmqpExchange", func(*AmqpOpenService, string) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudAmqpExchangeUpdate(d, rawClient)
		patcheDescribeBackups.Reset()
		assert.Nil(t, err)
	})

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
		err := resourceAliCloudAmqpExchangeDelete(d, rawClient)
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
		err := resourceAliCloudAmqpExchangeDelete(d, rawClient)
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
		err := resourceAliCloudAmqpExchangeDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeAmqpExchangeNotFound", func(t *testing.T) {
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
		err := resourceAliCloudAmqpExchangeRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeAmqpExchangeAbnormal", func(t *testing.T) {
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
		err := resourceAliCloudAmqpExchangeRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

// Test Amqp Exchange. >>> Resource test cases, automatically generated.
// Case topicExchange_测试用例 7251
func TestAccAliCloudAmqpExchange_basic7251(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AlicloudAmqpExchangeMap7251)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpExchange")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccamqp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAmqpExchangeBasicDependence7251)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"virtual_host_name": "${var.vhost_name}",
					"instance_id":       "${alicloud_amqp_instance.defaultInstace.id}",
					"internal":          "false",
					"auto_delete_state": "false",
					"exchange_type":     "TOPIC",
					"exchange_name":     name,
					"x_delayed_type":    "TOPIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_host_name": CHECKSET,
						"instance_id":       CHECKSET,
						"internal":          "false",
						"auto_delete_state": "false",
						"exchange_type":     "TOPIC",
						"exchange_name":     name,
						"x_delayed_type":    "TOPIC",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal", "x_delayed_type"},
			},
		},
	})
}

var AlicloudAmqpExchangeMap7251 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudAmqpExchangeBasicDependence7251(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "vhost_name" {
  default = "/"
}

resource "alicloud_amqp_instance" "defaultInstace" {
  renewal_duration      = "1"
  max_tps               = "3000"
  period_cycle          = "Month"
  max_connections       = "2000"
  support_eip           = true
  auto_renew            = false
  renewal_status        = "AutoRenewal"
  period                = "1"
  instance_name         = "OpenAPI-TestCase"
  support_tracing       = false
  payment_type          = "Subscription"
  renewal_duration_unit = "Month"
  instance_type         = "enterprise"
  queue_capacity        = "200"
  max_eip_tps           = "128"
  storage_size          = "0"
}


`, name)
}

// Case headersExchange_测试用例 7237
func TestAccAliCloudAmqpExchange_basic7237(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AlicloudAmqpExchangeMap7237)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpExchange")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccamqp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAmqpExchangeBasicDependence7237)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"virtual_host_name":  "${var.vhost_name}",
					"instance_id":        "${alicloud_amqp_instance.defaultInstace.id}",
					"internal":           "false",
					"auto_delete_state":  "false",
					"exchange_type":      "HEADERS",
					"exchange_name":      name,
					"alternate_exchange": "amq.headers",
					"x_delayed_type":     "HEADERS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_host_name":  CHECKSET,
						"instance_id":        CHECKSET,
						"internal":           "false",
						"auto_delete_state":  "false",
						"exchange_type":      "HEADERS",
						"exchange_name":      name,
						"alternate_exchange": "amq.headers",
						"x_delayed_type":     "HEADERS",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal", "x_delayed_type"},
			},
		},
	})
}

var AlicloudAmqpExchangeMap7237 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudAmqpExchangeBasicDependence7237(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "vhost_name" {
  default = "/"
}

resource "alicloud_amqp_instance" "defaultInstace" {
  renewal_duration      = "1"
  max_tps               = "3000"
  period_cycle          = "Month"
  max_connections       = "2000"
  support_eip           = true
  auto_renew            = false
  renewal_status        = "AutoRenewal"
  period                = "1"
  instance_name         = "OpenAPI-TestCase"
  support_tracing       = false
  payment_type          = "Subscription"
  renewal_duration_unit = "Month"
  instance_type         = "enterprise"
  queue_capacity        = "200"
  max_eip_tps           = "128"
  storage_size          = "0"
}


`, name)
}

// Case directExchange_测试用例 7249
func TestAccAliCloudAmqpExchange_basic7249(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AlicloudAmqpExchangeMap7249)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpExchange")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccamqp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAmqpExchangeBasicDependence7249)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"virtual_host_name": "${var.vhost_name}",
					"instance_id":       "${alicloud_amqp_instance.defaultInstace.id}",
					"internal":          "false",
					"auto_delete_state": "false",
					"exchange_type":     "DIRECT",
					"exchange_name":     name,
					"x_delayed_type":    "DIRECT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_host_name": CHECKSET,
						"instance_id":       CHECKSET,
						"internal":          "false",
						"auto_delete_state": "false",
						"exchange_type":     "DIRECT",
						"exchange_name":     name,
						"x_delayed_type":    "DIRECT",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal", "x_delayed_type"},
			},
		},
	})
}

var AlicloudAmqpExchangeMap7249 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudAmqpExchangeBasicDependence7249(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "vhost_name" {
  default = "/"
}

resource "alicloud_amqp_instance" "defaultInstace" {
  renewal_duration      = "1"
  max_tps               = "3000"
  period_cycle          = "Month"
  max_connections       = "2000"
  support_eip           = true
  auto_renew            = false
  renewal_status        = "AutoRenewal"
  period                = "1"
  instance_name         = "OpenAPI-TestCase"
  support_tracing       = false
  payment_type          = "Subscription"
  renewal_duration_unit = "Month"
  instance_type         = "enterprise"
  queue_capacity        = "200"
  max_eip_tps           = "128"
  storage_size          = "0"
}


`, name)
}

// Case fanoutExchange_测试用例 7250
func TestAccAliCloudAmqpExchange_basic7250(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AlicloudAmqpExchangeMap7250)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpExchange")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccamqp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAmqpExchangeBasicDependence7250)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"virtual_host_name": "${var.vhost_name}",
					"instance_id":       "${alicloud_amqp_instance.defaultInstace.id}",
					"internal":          "true",
					"auto_delete_state": "true",
					"exchange_type":     "FANOUT",
					"exchange_name":     name,
					"x_delayed_type":    "FANOUT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_host_name": CHECKSET,
						"instance_id":       CHECKSET,
						"internal":          "true",
						"auto_delete_state": "true",
						"exchange_type":     "FANOUT",
						"exchange_name":     name,
						"x_delayed_type":    "FANOUT",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal", "x_delayed_type"},
			},
		},
	})
}

var AlicloudAmqpExchangeMap7250 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudAmqpExchangeBasicDependence7250(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "vhost_name" {
  default = "/"
}

resource "alicloud_amqp_instance" "defaultInstace" {
  renewal_duration      = "1"
  max_tps               = "3000"
  period_cycle          = "Month"
  max_connections       = "2000"
  support_eip           = true
  auto_renew            = false
  renewal_status        = "AutoRenewal"
  period                = "1"
  instance_name         = "OpenAPI-TestCase"
  support_tracing       = false
  payment_type          = "Subscription"
  renewal_duration_unit = "Month"
  instance_type         = "enterprise"
  queue_capacity        = "200"
  max_eip_tps           = "128"
  storage_size          = "0"
}


`, name)
}

// Case xJmsTopicExchange_测试用例 7261
func TestAccAliCloudAmqpExchange_basic7261(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AlicloudAmqpExchangeMap7261)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpExchange")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccamqp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAmqpExchangeBasicDependence7261)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"virtual_host_name": "${var.vhost_name}",
					"instance_id":       "${alicloud_amqp_instance.defaultInstace.id}",
					"internal":          "false",
					"auto_delete_state": "false",
					"exchange_type":     "DIRECT",
					"exchange_name":     name,
					"x_delayed_type":    "X_JMS_TOPIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_host_name": CHECKSET,
						"instance_id":       CHECKSET,
						"internal":          "false",
						"auto_delete_state": "false",
						"exchange_type":     "DIRECT",
						"exchange_name":     name,
						"x_delayed_type":    "X_JMS_TOPIC",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal", "x_delayed_type"},
			},
		},
	})
}

var AlicloudAmqpExchangeMap7261 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudAmqpExchangeBasicDependence7261(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "vhost_name" {
  default = "/"
}

resource "alicloud_amqp_instance" "defaultInstace" {
  renewal_duration      = "1"
  max_tps               = "3000"
  period_cycle          = "Month"
  max_connections       = "2000"
  support_eip           = true
  auto_renew            = false
  renewal_status        = "AutoRenewal"
  period                = "1"
  instance_name         = "OpenAPI-TestCase"
  support_tracing       = false
  payment_type          = "Subscription"
  renewal_duration_unit = "Month"
  instance_type         = "enterprise"
  queue_capacity        = "200"
  max_eip_tps           = "128"
  storage_size          = "0"
}


`, name)
}

// Case TF增加Exchange类型_HEADERS 10181
func TestAccAliCloudAmqpExchange_basic10181(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AlicloudAmqpExchangeMap10181)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpExchange")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccamqp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAmqpExchangeBasicDependence10181)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"virtual_host_name":  "${var.virtual_host_name}",
					"instance_id":        "${alicloud_amqp_instance.CreateInstace.id}",
					"internal":           "true",
					"auto_delete_state":  "false",
					"exchange_name":      name,
					"exchange_type":      "HEADERS",
					"alternate_exchange": "bakExchange",
					"x_delayed_type":     "HEADERS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_host_name":  CHECKSET,
						"instance_id":        CHECKSET,
						"internal":           "true",
						"auto_delete_state":  "false",
						"exchange_name":      name,
						"exchange_type":      "HEADERS",
						"alternate_exchange": "bakExchange",
						"x_delayed_type":     "HEADERS",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal", "x_delayed_type"},
			},
		},
	})
}

var AlicloudAmqpExchangeMap10181 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudAmqpExchangeBasicDependence10181(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "virtual_host_name" {
  default = "/"
}

resource "alicloud_amqp_instance" "CreateInstace" {
  renewal_duration      = "1"
  max_tps               = "3000"
  period_cycle          = "Month"
  max_connections       = "2000"
  support_eip           = true
  auto_renew            = false
  renewal_status        = "AutoRenewal"
  period                = "12"
  instance_name         = "OpenAPI-TestCase"
  support_tracing       = false
  payment_type          = "Subscription"
  renewal_duration_unit = "Month"
  instance_type         = "enterprise"
  queue_capacity        = "200"
  max_eip_tps           = "128"
  storage_size          = "0"
}


`, name)
}

// Case TF增加Exchange类型_DIRECT 10140
func TestAccAliCloudAmqpExchange_basic10140(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AlicloudAmqpExchangeMap10140)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpExchange")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccamqp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAmqpExchangeBasicDependence10140)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"virtual_host_name":  "${var.virtual_host_name}",
					"instance_id":        "${alicloud_amqp_instance.CreateInstance.id}",
					"internal":           "false",
					"auto_delete_state":  "true",
					"exchange_name":      name,
					"exchange_type":      "DIRECT",
					"alternate_exchange": "bakExchange",
					"x_delayed_type":     "DIRECT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_host_name":  CHECKSET,
						"instance_id":        CHECKSET,
						"internal":           "false",
						"auto_delete_state":  "true",
						"exchange_name":      name,
						"exchange_type":      "DIRECT",
						"alternate_exchange": "bakExchange",
						"x_delayed_type":     "DIRECT",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal", "x_delayed_type"},
			},
		},
	})
}

var AlicloudAmqpExchangeMap10140 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudAmqpExchangeBasicDependence10140(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "virtual_host_name" {
  default = "/"
}

resource "alicloud_amqp_instance" "CreateInstance" {
  renewal_duration      = "1"
  max_tps               = "3000"
  period_cycle          = "Month"
  max_connections       = "2000"
  support_eip           = true
  auto_renew            = false
  renewal_status        = "AutoRenewal"
  period                = "12"
  instance_name         = "OpenAPI-TestCase"
  support_tracing       = false
  payment_type          = "Subscription"
  renewal_duration_unit = "Month"
  instance_type         = "enterprise"
  queue_capacity        = "200"
  max_eip_tps           = "128"
  storage_size          = "0"
}


`, name)
}

// Case TF增加Exchange类型_FANOUT 10179
func TestAccAliCloudAmqpExchange_basic10179(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AlicloudAmqpExchangeMap10179)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpExchange")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccamqp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAmqpExchangeBasicDependence10179)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"virtual_host_name":  "${var.virtual_host_name}",
					"instance_id":        "${alicloud_amqp_instance.CreateInstance.id}",
					"internal":           "true",
					"auto_delete_state":  "false",
					"exchange_name":      name,
					"exchange_type":      "FANOUT",
					"alternate_exchange": "bakExchange",
					"x_delayed_type":     "FANOUT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_host_name":  CHECKSET,
						"instance_id":        CHECKSET,
						"internal":           "true",
						"auto_delete_state":  "false",
						"exchange_name":      name,
						"exchange_type":      "FANOUT",
						"alternate_exchange": "bakExchange",
						"x_delayed_type":     "FANOUT",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal", "x_delayed_type"},
			},
		},
	})
}

var AlicloudAmqpExchangeMap10179 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudAmqpExchangeBasicDependence10179(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "virtual_host_name" {
  default = "/"
}

resource "alicloud_amqp_instance" "CreateInstance" {
  renewal_duration      = "1"
  max_tps               = "3000"
  period_cycle          = "Month"
  max_connections       = "2000"
  support_eip           = true
  auto_renew            = false
  renewal_status        = "AutoRenewal"
  period                = "12"
  instance_name         = "OpenAPI-TestCase"
  support_tracing       = false
  payment_type          = "Subscription"
  renewal_duration_unit = "Month"
  instance_type         = "enterprise"
  queue_capacity        = "200"
  max_eip_tps           = "128"
  storage_size          = "0"
}


`, name)
}

// Case TF增加Exchange类型_TOPIC 10180
func TestAccAliCloudAmqpExchange_basic10180(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AlicloudAmqpExchangeMap10180)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpExchange")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccamqp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAmqpExchangeBasicDependence10180)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"virtual_host_name":  "${var.virtual_host_name}",
					"internal":           "true",
					"auto_delete_state":  "false",
					"exchange_name":      name,
					"exchange_type":      "TOPIC",
					"alternate_exchange": "bakExchange",
					"x_delayed_type":     "TOPIC",
					"instance_id":        "${alicloud_amqp_instance.CreateInstance.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_host_name":  CHECKSET,
						"internal":           "true",
						"auto_delete_state":  "false",
						"exchange_name":      name,
						"exchange_type":      "TOPIC",
						"alternate_exchange": "bakExchange",
						"x_delayed_type":     "TOPIC",
						"instance_id":        CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal", "x_delayed_type"},
			},
		},
	})
}

var AlicloudAmqpExchangeMap10180 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudAmqpExchangeBasicDependence10180(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "virtual_host_name" {
  default = "/"
}

resource "alicloud_amqp_instance" "CreateInstance" {
  renewal_duration      = "1"
  max_tps               = "3000"
  period_cycle          = "Month"
  max_connections       = "2000"
  support_eip           = true
  auto_renew            = false
  renewal_status        = "AutoRenewal"
  period                = "12"
  instance_name         = "OpenAPI-TestCase"
  support_tracing       = false
  payment_type          = "Subscription"
  renewal_duration_unit = "Month"
  instance_type         = "enterprise"
  queue_capacity        = "200"
  max_eip_tps           = "128"
  storage_size          = "0"
}


`, name)
}

// Case TF增加Exchange类型_X_DELAYED_MESSAGE 10182
func TestAccAliCloudAmqpExchange_basic10182(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AlicloudAmqpExchangeMap10182)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpExchange")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccamqp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAmqpExchangeBasicDependence10182)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"virtual_host_name":  "${var.virtual_host_name}",
					"instance_id":        "${alicloud_amqp_instance.CreateInstance.id}",
					"internal":           "true",
					"auto_delete_state":  "false",
					"exchange_name":      name,
					"exchange_type":      "X_DELAYED_MESSAGE",
					"alternate_exchange": "bakExchange",
					"x_delayed_type":     "X_JMS_TOPIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_host_name":  CHECKSET,
						"instance_id":        CHECKSET,
						"internal":           "true",
						"auto_delete_state":  "false",
						"exchange_name":      name,
						"exchange_type":      "X_DELAYED_MESSAGE",
						"alternate_exchange": "bakExchange",
						"x_delayed_type":     "X_JMS_TOPIC",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal", "x_delayed_type"},
			},
		},
	})
}

var AlicloudAmqpExchangeMap10182 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudAmqpExchangeBasicDependence10182(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "virtual_host_name" {
  default = "/"
}

resource "alicloud_amqp_instance" "CreateInstance" {
  renewal_duration      = "1"
  max_tps               = "3000"
  period_cycle          = "Month"
  max_connections       = "2000"
  support_eip           = true
  auto_renew            = false
  renewal_status        = "AutoRenewal"
  period                = "12"
  instance_name         = "OpenAPI-TestCase"
  support_tracing       = false
  payment_type          = "Subscription"
  renewal_duration_unit = "Month"
  instance_type         = "enterprise"
  queue_capacity        = "200"
  max_eip_tps           = "128"
  storage_size          = "0"
}


`, name)
}

// Case TF增加Exchange类型_X_CONSISTENT_HASH 10183
func TestAccAliCloudAmqpExchange_basic10183(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AlicloudAmqpExchangeMap10183)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpExchange")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccamqp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAmqpExchangeBasicDependence10183)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"virtual_host_name":  "${var.virtual_host_name}",
					"instance_id":        "${alicloud_amqp_instance.CreateInstance.id}",
					"internal":           "true",
					"auto_delete_state":  "false",
					"exchange_name":      name,
					"exchange_type":      "X_CONSISTENT_HASH",
					"alternate_exchange": "bakExchange",
					"x_delayed_type":     "DIRECT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_host_name":  CHECKSET,
						"instance_id":        CHECKSET,
						"internal":           "true",
						"auto_delete_state":  "false",
						"exchange_name":      name,
						"exchange_type":      "X_CONSISTENT_HASH",
						"alternate_exchange": "bakExchange",
						"x_delayed_type":     "DIRECT",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal", "x_delayed_type"},
			},
		},
	})
}

var AlicloudAmqpExchangeMap10183 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudAmqpExchangeBasicDependence10183(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "virtual_host_name" {
  default = "/"
}

resource "alicloud_amqp_instance" "CreateInstance" {
  renewal_duration      = "1"
  max_tps               = "3000"
  period_cycle          = "Month"
  max_connections       = "2000"
  support_eip           = true
  auto_renew            = false
  renewal_status        = "AutoRenewal"
  period                = "12"
  instance_name         = "OpenAPI-TestCase"
  support_tracing       = false
  payment_type          = "Subscription"
  renewal_duration_unit = "Month"
  instance_type         = "enterprise"
  queue_capacity        = "200"
  max_eip_tps           = "128"
  storage_size          = "0"
}


`, name)
}

// Test Amqp Exchange. <<< Resource test cases, automatically generated.
