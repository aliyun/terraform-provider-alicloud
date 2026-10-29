package amqpopen_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/amqpopen"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_amqp_virtual_host", &resource.Sweeper{
		Name: "alicloud_amqp_virtual_host",
		F:    testSweepAmqpVirtualHost,
	})
}

func testSweepAmqpVirtualHost(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	action := "ListInstances"
	request := make(map[string]interface{})
	request["MaxResults"] = helper.PageSizeXLarge
	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcGet("amqp-open", "2019-12-12", action, request, nil)
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
			log.Println(errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_amqp_virtual_hosts", action, errs.AlibabaCloudSdkGoERROR))
			return nil
		}
		resp, err := jsonpath.Get("$.Data.Instances", response)
		if err != nil {
			log.Println(errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data.Instances", response))
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
			request["MaxResults"] = helper.PageSizeLarge
			var response map[string]interface{}
			for {
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(5*time.Minute, func() *retry.RetryError {
					response, err = client.RpcGet("amqp-open", "2019-12-12", action, request, nil)
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
					log.Println(errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_amqp_virtual_hosts", action, errs.AlibabaCloudSdkGoERROR))
					return nil
				}
				resp, err := jsonpath.Get("$.Data.VirtualHosts", response)
				if err != nil {
					log.Println(errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data.VirtualHosts", response))
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

					wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
					err = retry.Retry(3*time.Minute, func() *retry.RetryError {
						response, err = client.RpcPost("amqp-open", "2019-12-12", action, nil, request, false)
						if err != nil {
							if helper.NeedRetry(err) {
								wait()
								return retry.RetryableError(err)
							}
							return retry.NonRetryableError(err)
						}
						return nil
					})
					log.Println(errs.WrapError(err))
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

// lintignore: R001
func TestUnitAliCloudAmqpVirtualHost(t *testing.T) {
	p := acctest.Provider.ResourcesMap
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
	rawClient, err := sweep.SharedClientForRegion(region)
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
				Code:       helper.String(errorCode),
				Data:       helper.String(errorCode),
				Message:    helper.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_amqp_virtual_host", "MockVirtualHostName"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       helper.String(errorCode),
				Data:       helper.String(errorCode),
				Message:    helper.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		// The transient failure amqp-open returns while a freshly created instance
		// is still initializing, even though it already reports SERVING.
		"InstanceNotReadyError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("InvalidParameter"),
				Data:       helper.String(`{"Code":"InvalidParameter","Message":"the internal error!","RequestId":"MockRequestId"}`),
				Message:    helper.String("code: 400, the internal error! request id: MockRequestId"),
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
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := p["alicloud_amqp_virtual_host"].Create(d, rawClient)
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
		err := p["alicloud_amqp_virtual_host"].Create(d, rawClient)
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
		err := p["alicloud_amqp_virtual_host"].Create(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	// Creating a virtual host in the same apply as its instance can hit the
	// instance initialization window; the first attempt must be retried rather
	// than aborting the apply.
	t.Run("CreateInstanceNotReadyThenNormal", func(t *testing.T) {
		dNotReady, _ := schema.InternalMap(p["alicloud_amqp_virtual_host"].Schema).Data(nil, nil)
		dNotReady.MarkNewResource()
		for key, value := range map[string]interface{}{
			"instance_id":       "instance_id",
			"virtual_host_name": "virtual_host_name",
		} {
			err := dNotReady.Set(key, value)
			assert.Nil(t, err)
		}
		notReadyFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if notReadyFlag {
				notReadyFlag = false
				return responseMock["InstanceNotReadyError"]("")
			}
			return responseMock["CreateNormal"]("")
		})
		err := p["alicloud_amqp_virtual_host"].Create(dNotReady, rawClient)
		patches.Reset()
		assert.Nil(t, err)
		assert.False(t, notReadyFlag, "the transient error should have been served and retried")
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("instance_id", ":", "MockVirtualHostName"))

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewOnsproxyClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := p["alicloud_amqp_virtual_host"].Delete(d, rawClient)
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
		err := p["alicloud_amqp_virtual_host"].Delete(d, rawClient)
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
		err := p["alicloud_amqp_virtual_host"].Delete(d, rawClient)
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
		err := p["alicloud_amqp_virtual_host"].Read(d, rawClient)
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
		err := p["alicloud_amqp_virtual_host"].Read(d, rawClient)
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
		patcheDescribeAmqpVirtualHost := gomonkey.ApplyMethod(reflect.TypeOf(&amqpopen.AmqpOpenService{}), "DescribeAmqpVirtualHost", func(*amqpopen.AmqpOpenService, string) (map[string]interface{}, error) {
			return responseMock["ReadNormal"]("")
		})
		err := p["alicloud_amqp_virtual_host"].Read(resourceData1, rawClient)
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAmqpVirtualHostMap7216)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return amqpopen.NewAmqpServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAmqpVirtualHost")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccamqp%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAmqpVirtualHostBasicDependence7216)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"virtual_host_name": name,
					"instance_id":       "${data.alicloud_amqp_instances.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_host_name": name,
						"instance_id":       acctest.CHECKSET,
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

	data "alicloud_amqp_instances" "default" {
  		status = "SERVING"
	}
`, name)
}

// Test Amqp VirtualHost. <<< Resource test cases, automatically generated.
