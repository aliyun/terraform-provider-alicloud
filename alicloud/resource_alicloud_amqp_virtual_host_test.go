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
	resource.AddTestSweepers("alicloud_amqp_virtual_host", &resource.Sweeper{
		Name: "alicloud_amqp_virtual_host",
		F:    testSweepAmqpVirtualHost,
	})
}

func testSweepAmqpVirtualHost(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	action := "ListInstances"
	request := make(map[string]interface{})
	request["MaxResults"] = PageSizeXLarge
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RpcGet("amqp-open", "2019-12-12", action, request, nil)
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
			log.Println(WrapErrorf(err, DataDefaultErrorMsg, "alicloud_amqp_virtual_hosts", action, AlibabaCloudSdkGoERROR))
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
			if fmt.Sprint(item["Status"]) != "SERVING" {
				continue
			}
			action := "ListVirtualHosts"
			request := make(map[string]interface{})
			request["InstanceId"] = instanceId
			request["MaxResults"] = PageSizeLarge
			var response map[string]interface{}
			for {
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(5*time.Minute, func() *resource.RetryError {
					response, err = client.RpcGet("amqp-open", "2019-12-12", action, request, nil)
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
					log.Println(WrapErrorf(err, DataDefaultErrorMsg, "alicloud_amqp_virtual_hosts", action, AlibabaCloudSdkGoERROR))
					return nil
				}
				resp, err := jsonpath.Get("$.Data.VirtualHosts", response)
				if err != nil {
					log.Println(WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.VirtualHosts", response))
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

					action := "DeleteVirtualHost"
					request := map[string]interface{}{
						"InstanceId":  instanceId,
						"VirtualHost": item["Name"],
					}

					wait := incrementalWait(3*time.Second, 3*time.Second)
					err = resource.Retry(3*time.Minute, func() *resource.RetryError {
						response, err = client.RpcPost("amqp-open", "2019-12-12", action, nil, request, false)
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
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	return nil
}

func TestUnitAliCloudAmqpVirtualHost(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_amqp_virtual_host"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_amqp_virtual_host"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"instance_id":       "instance_id",
		"virtual_host_name": "virtual_host_name",
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
			"VirtualHosts": []interface{}{
				map[string]interface{}{
					"InstanceId":      "instance_id",
					"VirtualHostName": "virtual_host_name",
					"Name":            "MockVirtualHostName",
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
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_amqp_virtual_host", "MockVirtualHostName"))
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
		err := resourceAliCloudAmqpVirtualHostCreate(d, rawClient)
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
		err := resourceAliCloudAmqpVirtualHostCreate(d, rawClient)
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
		err := resourceAliCloudAmqpVirtualHostCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("instance_id", ":", "MockVirtualHostName"))

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
		err := resourceAliCloudAmqpVirtualHostDelete(d, rawClient)
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
		err := resourceAliCloudAmqpVirtualHostDelete(d, rawClient)
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
		err := resourceAliCloudAmqpVirtualHostDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeAmqpVirtualHostNotFound", func(t *testing.T) {
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
		err := resourceAliCloudAmqpVirtualHostRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeAmqpVirtualHostAbnormal", func(t *testing.T) {
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
		err := resourceAliCloudAmqpVirtualHostRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})

	t.Run("ReadMockParseResourceId", func(t *testing.T) {
		resourceData1, _ := schema.InternalMap(p["alicloud_amqp_virtual_host"].Schema).Data(nil, nil)
		resourceData1.SetId(fmt.Sprint("instance_id", ":", "MockVirtualHostName"))
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := false
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		patcheDescribeAmqpVirtualHost := gomonkey.ApplyMethod(reflect.TypeOf(&AmqpOpenService{}), "DescribeAmqpVirtualHost", func(*AmqpOpenService, string) (map[string]interface{}, error) {
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudAmqpVirtualHostRead(resourceData1, rawClient)
		patcheDorequest.Reset()
		patcheDescribeAmqpVirtualHost.Reset()
		assert.Nil(t, err)
	})
}

// Test Amqp VirtualHost. >>> Resource test cases, automatically generated.
// Case VirtualHost_测试用例 7216
func TestAccAliCloudAmqpVirtualHost_basic7216(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_virtual_host.default"
	ra := resourceAttrInit(resourceId, AliCloudAmqpVirtualHostMap7216)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpVirtualHost")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccamqp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAmqpVirtualHostBasicDependence7216)
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
					"virtual_host_name": name,
					"instance_id":       "${alicloud_amqp_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_host_name": name,
						"instance_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudAmqpVirtualHostMap7216 = map[string]string{}

func AliCloudAmqpVirtualHostBasicDependence7216(name string) string {
	return fmt.Sprintf(`
	variable "name" {
    	default = "%s"
	}

	resource "alicloud_amqp_instance" "default" {
  		instance_name         = var.name
  		instance_type         = "enterprise"
  		max_tps               = 3000
  		max_connections       = 2000
  		queue_capacity        = 200
  		payment_type          = "Subscription"
  		renewal_status        = "AutoRenewal"
  		renewal_duration      = 1
  		renewal_duration_unit = "Year"
  		support_eip           = true
	}
`, name)
}

// Test Amqp VirtualHost. <<< Resource test cases, automatically generated.
