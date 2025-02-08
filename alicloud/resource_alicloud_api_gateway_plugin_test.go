package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
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

func init() {
	resource.AddTestSweepers(
		"alicloud_api_gateway_plugin",
		&resource.Sweeper{
			Name: "alicloud_api_gateway_plugin",
			F:    testSweepApiGatewayPlugin,
		})
}

func testSweepApiGatewayPlugin(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.ApiGatewaySupportRegions) {
		log.Printf("[INFO] Skipping Api Gateway Plugin unsupported region: %s", region)
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribePlugins"
	request := map[string]interface{}{}
	request["RegionId"] = aliyunClient.RegionId

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = aliyunClient.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.Plugins.PluginAttribute", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Plugins.PluginAttribute", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["PluginName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Api Gateway Plugin: %s", item["PluginName"].(string))
				continue
			}
			action := "DeletePlugin"
			request := map[string]interface{}{
				"PluginId": item["PluginId"],
			}
			_, err = aliyunClient.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Api Gateway Plugin (%s): %s", item["PluginName"].(string), err)
			}
			log.Printf("[INFO] Delete Api Gateway Plugin success: %s ", item["PluginName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudApiGatewayPlugin_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	checkoutSupportedRegions(t, true, connectivity.ApiGatewaySupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence0)
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
					"description": "${var.name}",
					"plugin_name": "${var.name}",
					"plugin_data": `{\"allowOrigins\": \"api.foo.com\",\"allowMethods\": \"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\",\"allowHeaders\": \"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\",\"exposeHeaders\": \"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\",\"maxAge\": 172800,\"allowCredentials\": true}`,
					"plugin_type": "cors",
					"tags": map[string]string{
						"Created": "tfTestAcc0",
						"For":     "Tftestacc 0",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  name,
						"plugin_name":  name,
						"plugin_data":  "{\"allowOrigins\": \"api.foo.com\",\"allowMethods\": \"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\",\"allowHeaders\": \"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\",\"exposeHeaders\": \"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\",\"maxAge\": 172800,\"allowCredentials\": true}",
						"plugin_type":  "cors",
						"tags.%":       "2",
						"tags.Created": "tfTestAcc0",
						"tags.For":     "Tftestacc 0",
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

var AlicloudApiGatewayPluginMap0 = map[string]string{
	"tags.%":      CHECKSET,
	"plugin_data": CHECKSET,
	"plugin_type": CHECKSET,
}

func AlicloudApiGatewayPluginBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

func TestAccAliCloudApiGatewayPlugin_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	checkoutSupportedRegions(t, true, connectivity.ApiGatewaySupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence1)
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
					"plugin_name": "${var.name}",
					"plugin_data": `{\"allowOrigins\": \"api.foo.com\",\"allowMethods\": \"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\",\"allowHeaders\": \"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\",\"exposeHeaders\": \"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\",\"maxAge\": 172800,\"allowCredentials\": true}`,
					"plugin_type": "cors",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_name": name,
						"plugin_data": "{\"allowOrigins\": \"api.foo.com\",\"allowMethods\": \"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\",\"allowHeaders\": \"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\",\"exposeHeaders\": \"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\",\"maxAge\": 172800,\"allowCredentials\": true}",
						"plugin_type": "cors",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plugin_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plugin_data": `{\"allowOrigins\": \"api.foo1.com\",\"allowMethods\": \"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\",\"allowHeaders\": \"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\",\"exposeHeaders\": \"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\",\"maxAge\": 172800,\"allowCredentials\": true}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_data": "{\"allowOrigins\": \"api.foo1.com\",\"allowMethods\": \"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\",\"allowHeaders\": \"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\",\"exposeHeaders\": \"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\",\"maxAge\": 172800,\"allowCredentials\": true}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "tfTestAcc4",
						"For":     "Tftestacc 4",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "tfTestAcc4",
						"tags.For":     "Tftestacc 4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}",
					"plugin_name": "${var.name}",
					"plugin_data": `{\"allowOrigins\": \"api.foo.com\",\"allowMethods\": \"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\",\"allowHeaders\": \"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\",\"exposeHeaders\": \"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\",\"maxAge\": 172800,\"allowCredentials\": true}`,
					"tags": map[string]string{
						"Created": "tfTestAcc5",
						"For":     "Tftestacc 5",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  name,
						"plugin_name":  name,
						"plugin_data":  "{\"allowOrigins\": \"api.foo.com\",\"allowMethods\": \"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\",\"allowHeaders\": \"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\",\"exposeHeaders\": \"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\",\"maxAge\": 172800,\"allowCredentials\": true}",
						"tags.%":       "2",
						"tags.Created": "tfTestAcc5",
						"tags.For":     "Tftestacc 5",
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

var AlicloudApiGatewayPluginMap1 = map[string]string{
	"plugin_data": CHECKSET,
	"plugin_type": CHECKSET,
}

func AlicloudApiGatewayPluginBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

func TestUnitAccAlicloudApiGatewayPlugin(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_api_gateway_plugin"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_api_gateway_plugin"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"description": "CreateApiGatewayPluginValue",
		"plugin_name": "CreateApiGatewayPluginValue",
		"plugin_data": "CreateApiGatewayPluginValue",
		"plugin_type": "CreateApiGatewayPluginValue",
		"tags": map[string]string{
			"Created": "CreateApiGatewayPluginValue",
		},
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}

	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"Plugins": map[string]interface{}{
			"PluginAttribute": []interface{}{
				map[string]interface{}{
					"Description": "CreateApiGatewayPluginValue",
					"PluginName":  "CreateApiGatewayPluginValue",
					"PluginData":  "CreateApiGatewayPluginValue",
					"PluginId":    "CreateApiGatewayPluginValue",
					"PluginType":  "CreateApiGatewayPluginValue",
					"Tags": map[string]interface{}{
						"TagInfo": []interface{}{
							map[string]interface{}{
								"Value": "Created",
								"Key":   "CreateApiGatewayPluginValue",
							},
						},
					},
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		"PluginId": "CreateApiGatewayPluginValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_api_gateway_plugin", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewApigatewayClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudApiGatewayPluginCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreatePlugin" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudApiGatewayPluginCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_api_gateway_plugin"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewApigatewayClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudApiGatewayPluginUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"description": "UpdateApiGatewayPluginValue",
		"plugin_name": "UpdateApiGatewayPluginValue",
		"plugin_data": "UpdateApiGatewayPluginValue",
		"tags": map[string]string{
			"Created": "UpdateApiGatewayPluginValue",
		},
	}
	diff, err := newInstanceDiff("alicloud_api_gateway_plugin", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_api_gateway_plugin"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"Plugins": map[string]interface{}{
			"PluginAttribute": []interface{}{
				map[string]interface{}{
					"Description": "UpdateApiGatewayPluginValue",
					"PluginName":  "UpdateApiGatewayPluginValue",
					"PluginData":  "UpdateApiGatewayPluginValue",
					"Tags": map[string]interface{}{
						"TagInfo": []interface{}{
							map[string]interface{}{
								"Value": "Created",
								"Key":   "UpdateApiGatewayPluginValue",
							},
						},
					},
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyPlugin" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudApiGatewayPluginUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_api_gateway_plugin"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	diff, err = newInstanceDiff("alicloud_api_gateway_plugin", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_api_gateway_plugin"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribePlugins" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudApiGatewayPluginRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewApigatewayClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudApiGatewayPluginDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_api_gateway_plugin", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_api_gateway_plugin"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeletePlugin" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudApiGatewayPluginDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}

// Test ApiGateway Plugin. >>> Resource test cases, automatically generated.
// Case 后端路由插件 6792
func TestAccAliCloudApiGatewayPlugin_basic6792(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6792)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6792)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"plugin_name": name,
					"plugin_data": "{   \\\"routes\\\": [     {       \\\"name\\\": \\\"Vip\\\",       \\\"condition\\\": \\\"$CaAppId = 123456\\\",       \\\"backend\\\": {         \\\"type\\\": \\\"HTTP-VPC\\\",         \\\"vpcAccessName\\\": \\\"slbAccessForVip\\\"       }     },     {       \\\"name\\\": \\\"MockForOldClient\\\",       \\\"condition\\\": \\\"$ClientVersion < '2.0.5'\\\",       \\\"backend\\\": {         \\\"type\\\": \\\"MOCK\\\",         \\\"statusCode\\\": 400,         \\\"mockBody\\\": \\\"This version is not supported!!!\\\"       }     },     {       \\\"name\\\": \\\"BlueGreenPercent05\\\",       \\\"condition\\\": \\\"1 = 1\\\",       \\\"backend\\\": {         \\\"type\\\": \\\"HTTP\\\",         \\\"address\\\": \\\"https://beta-version.api.foo.com\\\"       },       \\\"constant-parameters\\\": [         {           \\\"name\\\": \\\"x-route-blue-green\\\",           \\\"location\\\": \\\"header\\\",           \\\"value\\\": \\\"route-blue-green\\\"         }       ]     }   ] }",
					"plugin_type": "routing",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_name": name,
						"plugin_data": "{   \"routes\": [     {       \"name\": \"Vip\",       \"condition\": \"$CaAppId = 123456\",       \"backend\": {         \"type\": \"HTTP-VPC\",         \"vpcAccessName\": \"slbAccessForVip\"       }     },     {       \"name\": \"MockForOldClient\",       \"condition\": \"$ClientVersion < '2.0.5'\",       \"backend\": {         \"type\": \"MOCK\",         \"statusCode\": 400,         \"mockBody\": \"This version is not supported!!!\"       }     },     {       \"name\": \"BlueGreenPercent05\",       \"condition\": \"1 = 1\",       \"backend\": {         \"type\": \"HTTP\",         \"address\": \"https://beta-version.api.foo.com\"       },       \"constant-parameters\": [         {           \"name\": \"x-route-blue-green\",           \"location\": \"header\",           \"value\": \"route-blue-green\"         }       ]     }   ] }",
						"plugin_type": "routing",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
					"plugin_name": name + "_update",
					"plugin_data": "{   \\\"routes\\\": [     {       \\\"name\\\": \\\"Vip\\\",       \\\"condition\\\": \\\"$CaAppId = 123456\\\",       \\\"backend\\\": {         \\\"type\\\": \\\"HTTP-VPC\\\",         \\\"vpcAccessName\\\": \\\"slbAccessForVip\\\"       }     },     {       \\\"name\\\": \\\"MockForOldClient\\\",       \\\"condition\\\": \\\"$ClientVersion < '2.0.5'\\\",       \\\"backend\\\": {         \\\"type\\\": \\\"MOCK\\\",         \\\"statusCode\\\": 400,         \\\"mockBody\\\": \\\"This version is not supported!!!\\\"       }     },     {       \\\"name\\\": \\\"BlueGreenPercent05\\\",       \\\"condition\\\": \\\"1 = 1\\\",       \\\"backend\\\": {         \\\"type\\\": \\\"HTTP\\\",         \\\"address\\\": \\\"https://beta-version.api.foo.com\\\"       },       \\\"constant-parameters\\\": [         {           \\\"name\\\": \\\"x-route-blue-green\\\",           \\\"location\\\": \\\"header\\\",           \\\"value\\\": \\\"route-blue-green\\\"         }       ]     }   ] }",
					"plugin_type": "routing",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name + "_update",
						"plugin_data": "{   \"routes\": [     {       \"name\": \"Vip\",       \"condition\": \"$CaAppId = 123456\",       \"backend\": {         \"type\": \"HTTP-VPC\",         \"vpcAccessName\": \"slbAccessForVip\"       }     },     {       \"name\": \"MockForOldClient\",       \"condition\": \"$ClientVersion < '2.0.5'\",       \"backend\": {         \"type\": \"MOCK\",         \"statusCode\": 400,         \"mockBody\": \"This version is not supported!!!\"       }     },     {       \"name\": \"BlueGreenPercent05\",       \"condition\": \"1 = 1\",       \"backend\": {         \"type\": \"HTTP\",         \"address\": \"https://beta-version.api.foo.com\"       },       \"constant-parameters\": [         {           \"name\": \"x-route-blue-green\",           \"location\": \"header\",           \"value\": \"route-blue-green\"         }       ]     }   ] }",
						"plugin_type": "routing",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

var AlicloudApiGatewayPluginMap6792 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudApiGatewayPluginBasicDependence6792(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case BasicAuth插件 6789
func TestAccAliCloudApiGatewayPlugin_basic6789(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6789)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6789)
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
					"plugin_name": name,
					"plugin_data": "{   \\\"users\\\": [     {       \\\"username\\\": \\\"alice\\\",       \\\"password\\\": 123456     },     {       \\\"username\\\": \\\"bob\\\",       \\\"password\\\": 666666     },     {       \\\"username\\\": \\\"charlie\\\",       \\\"password\\\": 888888     },     {       \\\"username\\\": \\\"dave\\\",       \\\"password\\\": 111111     }   ] }",
					"plugin_type": "basicAuth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_name": name,
						"plugin_data": "{   \"users\": [     {       \"username\": \"alice\",       \"password\": 123456     },     {       \"username\": \"bob\",       \"password\": 666666     },     {       \"username\": \"charlie\",       \"password\": 888888     },     {       \"username\": \"dave\",       \"password\": 111111     }   ] }",
						"plugin_type": "basicAuth",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
					"plugin_name": name + "_update",
					"plugin_data": "{   \\\"users\\\": [     {       \\\"username\\\": \\\"alice\\\",       \\\"password\\\": 123456     },     {       \\\"username\\\": \\\"bob\\\",       \\\"password\\\": 666666     },     {       \\\"username\\\": \\\"charlie\\\",       \\\"password\\\": 888888     },     {       \\\"username\\\": \\\"dave\\\",       \\\"password\\\": 111111     }   ] }",
					"plugin_type": "basicAuth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name + "_update",
						"plugin_data": "{   \"users\": [     {       \"username\": \"alice\",       \"password\": 123456     },     {       \"username\": \"bob\",       \"password\": 666666     },     {       \"username\": \"charlie\",       \"password\": 888888     },     {       \"username\": \"dave\",       \"password\": 111111     }   ] }",
						"plugin_type": "basicAuth",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

var AlicloudApiGatewayPluginMap6789 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudApiGatewayPluginBasicDependence6789(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 参数访问控制插件 6793
func TestAccAliCloudApiGatewayPlugin_basic6793(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6793)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6793)
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
					"plugin_name": name,
					"plugin_data": "{   \\\"parameters\\\": {     \\\"userId\\\": \\\"Token:userId\\\",     \\\"userType\\\": \\\"Token:userType\\\",     \\\"pathUserId\\\": \\\"path:userId\\\"   },   \\\"rules\\\": [     {       \\\"name\\\": \\\"admin\\\",       \\\"condition\\\": \\\"$userType = 'admin'\\\",       \\\"ifTrue\\\": \\\"ALLOW\\\"     },     {       \\\"name\\\": \\\"user\\\",       \\\"condition\\\": \\\"$userId = $pathUserId\\\",       \\\"ifFalse\\\": \\\"DENY\\\",       \\\"statusCode\\\": 403,       \\\"errorMessage\\\": \\\"Path not match $${userId} vs /$${pathUserId}\\\",       \\\"responseHeaders\\\": {         \\\"Content-Type\\\": \\\"application/xml\\\"       },       \\\"responseBody\\\": \\\"<Reason>Path not match $${userId} vs /$${pathUserId}</Reason>\\\\n\\\"     }   ] }",
					"plugin_type": "accessControl",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_name": name,
						"plugin_data": "{   \"parameters\": {     \"userId\": \"Token:userId\",     \"userType\": \"Token:userType\",     \"pathUserId\": \"path:userId\"   },   \"rules\": [     {       \"name\": \"admin\",       \"condition\": \"$userType = 'admin'\",       \"ifTrue\": \"ALLOW\"     },     {       \"name\": \"user\",       \"condition\": \"$userId = $pathUserId\",       \"ifFalse\": \"DENY\",       \"statusCode\": 403,       \"errorMessage\": \"Path not match ${userId} vs /${pathUserId}\",       \"responseHeaders\": {         \"Content-Type\": \"application/xml\"       },       \"responseBody\": \"<Reason>Path not match ${userId} vs /${pathUserId}</Reason>\\n\"     }   ] }",
						"plugin_type": "accessControl",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
					"plugin_name": name + "_update",
					"plugin_data": "{   \\\"parameters\\\": {     \\\"userId\\\": \\\"Token:userId\\\",     \\\"userType\\\": \\\"Token:userType\\\",     \\\"pathUserId\\\": \\\"path:userId\\\"   },   \\\"rules\\\": [     {       \\\"name\\\": \\\"admin\\\",       \\\"condition\\\": \\\"$userType = 'admin'\\\",       \\\"ifTrue\\\": \\\"ALLOW\\\"     },     {       \\\"name\\\": \\\"user\\\",       \\\"condition\\\": \\\"$userId = $pathUserId\\\",       \\\"ifFalse\\\": \\\"DENY\\\",       \\\"statusCode\\\": 403,       \\\"errorMessage\\\": \\\"Path not match $${userId} vs /$${pathUserId}\\\",       \\\"responseHeaders\\\": {         \\\"Content-Type\\\": \\\"application/xml\\\"       },       \\\"responseBody\\\": \\\"<Reason>Path not match $${userId} vs /$${pathUserId}</Reason>\\\\n\\\"     }   ] }",
					"plugin_type": "accessControl",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name + "_update",
						"plugin_data": "{   \"parameters\": {     \"userId\": \"Token:userId\",     \"userType\": \"Token:userType\",     \"pathUserId\": \"path:userId\"   },   \"rules\": [     {       \"name\": \"admin\",       \"condition\": \"$userType = 'admin'\",       \"ifTrue\": \"ALLOW\"     },     {       \"name\": \"user\",       \"condition\": \"$userId = $pathUserId\",       \"ifFalse\": \"DENY\",       \"statusCode\": 403,       \"errorMessage\": \"Path not match ${userId} vs /${pathUserId}\",       \"responseHeaders\": {         \"Content-Type\": \"application/xml\"       },       \"responseBody\": \"<Reason>Path not match ${userId} vs /${pathUserId}</Reason>\\n\"     }   ] }",
						"plugin_type": "accessControl",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

var AlicloudApiGatewayPluginMap6793 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudApiGatewayPluginBasicDependence6793(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 第三方鉴权插件 6796
func TestAccAliCloudApiGatewayPlugin_basic6796(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6796)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6796)
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
					"plugin_name": name,
					"plugin_data": "{   \\\"parameters\\\": {     \\\"statusCode\\\": \\\"StatusCode\\\"   },   \\\"authUriType\\\": \\\"HTTP\\\",   \\\"authUri\\\": {     \\\"address\\\": \\\"http://your-auth-domain.com:8080\\\",     \\\"path\\\": \\\"/your/authPath\\\",     \\\"timeout\\\": 7000,     \\\"method\\\": \\\"POST\\\"   },   \\\"passThroughBody\\\": false,   \\\"cachedTimeBySecond\\\": 10,   \\\"authParameters\\\": [     {       \\\"targetParameterName\\\": \\\"x-userId\\\",       \\\"sourceParameterName\\\": \\\"userId\\\",       \\\"targetLocation\\\": \\\"form\\\",       \\\"sourceLocation\\\": \\\"query\\\"     },     {       \\\"targetParameterName\\\": \\\"x-passwoed\\\",       \\\"sourceParameterName\\\": \\\"password\\\",       \\\"targetLocation\\\": \\\"form\\\",       \\\"sourceLocation\\\": \\\"query\\\"     }   ],   \\\"successCondition\\\": \\\"$${statusCode} = 200\\\",   \\\"errorMessage\\\": \\\"auth failed\\\",   \\\"errorStatusCode\\\": 401 }",
					"plugin_type": "remoteAuth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_name": name,
						"plugin_data": "{   \"parameters\": {     \"statusCode\": \"StatusCode\"   },   \"authUriType\": \"HTTP\",   \"authUri\": {     \"address\": \"http://your-auth-domain.com:8080\",     \"path\": \"/your/authPath\",     \"timeout\": 7000,     \"method\": \"POST\"   },   \"passThroughBody\": false,   \"cachedTimeBySecond\": 10,   \"authParameters\": [     {       \"targetParameterName\": \"x-userId\",       \"sourceParameterName\": \"userId\",       \"targetLocation\": \"form\",       \"sourceLocation\": \"query\"     },     {       \"targetParameterName\": \"x-passwoed\",       \"sourceParameterName\": \"password\",       \"targetLocation\": \"form\",       \"sourceLocation\": \"query\"     }   ],   \"successCondition\": \"${statusCode} = 200\",   \"errorMessage\": \"auth failed\",   \"errorStatusCode\": 401 }",
						"plugin_type": "remoteAuth",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
					"plugin_name": name + "_update",
					"plugin_data": "{   \\\"parameters\\\": {     \\\"statusCode\\\": \\\"StatusCode\\\"   },   \\\"authUriType\\\": \\\"HTTP\\\",   \\\"authUri\\\": {     \\\"address\\\": \\\"http://your-auth-domain.com:8080\\\",     \\\"path\\\": \\\"/your/authPath\\\",     \\\"timeout\\\": 7000,     \\\"method\\\": \\\"POST\\\"   },   \\\"passThroughBody\\\": false,   \\\"cachedTimeBySecond\\\": 10,   \\\"authParameters\\\": [     {       \\\"targetParameterName\\\": \\\"x-userId\\\",       \\\"sourceParameterName\\\": \\\"userId\\\",       \\\"targetLocation\\\": \\\"form\\\",       \\\"sourceLocation\\\": \\\"query\\\"     },     {       \\\"targetParameterName\\\": \\\"x-passwoed\\\",       \\\"sourceParameterName\\\": \\\"password\\\",       \\\"targetLocation\\\": \\\"form\\\",       \\\"sourceLocation\\\": \\\"query\\\"     }   ],   \\\"successCondition\\\": \\\"$${statusCode} = 200\\\",   \\\"errorMessage\\\": \\\"auth failed\\\",   \\\"errorStatusCode\\\": 401 }",
					"plugin_type": "remoteAuth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name + "_update",
						"plugin_data": "{   \"parameters\": {     \"statusCode\": \"StatusCode\"   },   \"authUriType\": \"HTTP\",   \"authUri\": {     \"address\": \"http://your-auth-domain.com:8080\",     \"path\": \"/your/authPath\",     \"timeout\": 7000,     \"method\": \"POST\"   },   \"passThroughBody\": false,   \"cachedTimeBySecond\": 10,   \"authParameters\": [     {       \"targetParameterName\": \"x-userId\",       \"sourceParameterName\": \"userId\",       \"targetLocation\": \"form\",       \"sourceLocation\": \"query\"     },     {       \"targetParameterName\": \"x-passwoed\",       \"sourceParameterName\": \"password\",       \"targetLocation\": \"form\",       \"sourceLocation\": \"query\"     }   ],   \"successCondition\": \"${statusCode} = 200\",   \"errorMessage\": \"auth failed\",   \"errorStatusCode\": 401 }",
						"plugin_type": "remoteAuth",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

var AlicloudApiGatewayPluginMap6796 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudApiGatewayPluginBasicDependence6796(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 缓存插件 6791
func TestAccAliCloudApiGatewayPlugin_basic6791(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6791)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6791)
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
					"plugin_name": name,
					"plugin_data": "{   \\\"varyByApp\\\": false,   \\\"varyByParameters\\\": [     \\\"userId\\\"   ],   \\\"varyByHeaders\\\": [     \\\"Accept\\\"   ],   \\\"clientCacheControl\\\": {     \\\"mode\\\": \\\"app\\\",     \\\"apps\\\": [       1992323,       1239922     ]   },   \\\"cacheableHeaders\\\": [     \\\"X-Customer-Token\\\"   ] }",
					"plugin_type": "caching",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_name": name,
						"plugin_data": "{   \"varyByApp\": false,   \"varyByParameters\": [     \"userId\"   ],   \"varyByHeaders\": [     \"Accept\"   ],   \"clientCacheControl\": {     \"mode\": \"app\",     \"apps\": [       1992323,       1239922     ]   },   \"cacheableHeaders\": [     \"X-Customer-Token\"   ] }",
						"plugin_type": "caching",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
					"plugin_name": name + "_update",
					"plugin_data": "{   \\\"varyByApp\\\": false,   \\\"varyByParameters\\\": [     \\\"userId\\\"   ],   \\\"varyByHeaders\\\": [     \\\"Accept\\\"   ],   \\\"clientCacheControl\\\": {     \\\"mode\\\": \\\"app\\\",     \\\"apps\\\": [       1992323,       1239922     ]   },   \\\"cacheableHeaders\\\": [     \\\"X-Customer-Token\\\"   ] }",
					"plugin_type": "caching",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name + "_update",
						"plugin_data": "{   \"varyByApp\": false,   \"varyByParameters\": [     \"userId\"   ],   \"varyByHeaders\": [     \"Accept\"   ],   \"clientCacheControl\": {     \"mode\": \"app\",     \"apps\": [       1992323,       1239922     ]   },   \"cacheableHeaders\": [     \"X-Customer-Token\"   ] }",
						"plugin_type": "caching",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

var AlicloudApiGatewayPluginMap6791 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudApiGatewayPluginBasicDependence6791(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 后端签名插件 6799
func TestAccAliCloudApiGatewayPlugin_basic6799(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6799)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6799)
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
					"plugin_name": name,
					"plugin_data": "{   \\\"type\\\": \\\"APIGW_BACKEND\\\",   \\\"key\\\": \\\"SampleKey\\\",   \\\"secret\\\": \\\"SampleSecret\\\" }",
					"plugin_type": "backendSignature",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_name": name,
						"plugin_data": "{   \"type\": \"APIGW_BACKEND\",   \"key\": \"SampleKey\",   \"secret\": \"SampleSecret\" }",
						"plugin_type": "backendSignature",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
					"plugin_name": name + "_update",
					"plugin_data": "{   \\\"type\\\": \\\"APIGW_BACKEND\\\",   \\\"key\\\": \\\"SampleKey\\\",   \\\"secret\\\": \\\"SampleSecret\\\" }",
					"plugin_type": "backendSignature",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name + "_update",
						"plugin_data": "{   \"type\": \"APIGW_BACKEND\",   \"key\": \"SampleKey\",   \"secret\": \"SampleSecret\" }",
						"plugin_type": "backendSignature",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

var AlicloudApiGatewayPluginMap6799 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudApiGatewayPluginBasicDependence6799(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 错误码映射插件 6794
func TestAccAliCloudApiGatewayPlugin_basic6794(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6794)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6794)
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
					"plugin_name": name,
					"plugin_data": "{   \\\"parameters\\\": {     \\\"statusCode\\\": \\\"StatusCode\\\",     \\\"resultCode\\\": \\\"BodyJsonField:$.result_code\\\",     \\\"requestId\\\": \\\"BodyJsonField:$.req_msg_id\\\"   },   \\\"errorCondition\\\": \\\"$statusCode = 200 and $resultCode != null and $resultCode != 'OK'\\\",   \\\"errorCode\\\": \\\"resultCode\\\",   \\\"mappings\\\": [     {       \\\"code\\\": \\\"ROLE_NOT_EXISTS\\\",       \\\"statusCode\\\": 404,       \\\"errorMessage\\\": \\\"Role Not Exists, RequestId=$${requestId}\\\"     },     {       \\\"code\\\": \\\"INVALID_PARAMETER\\\",       \\\"statusCode\\\": 400,       \\\"errorMessage\\\": \\\"Invalid Parameter, RequestId=$${requestId}\\\"     }   ],   \\\"defaultMapping\\\": {     \\\"statusCode\\\": 500,     \\\"errorMessage\\\": \\\"Unknown Error, $${resultCode}, RequestId=$${requestId}\\\"   } }",
					"plugin_type": "errorMapping",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_name": name,
						"plugin_data": "{   \"parameters\": {     \"statusCode\": \"StatusCode\",     \"resultCode\": \"BodyJsonField:$.result_code\",     \"requestId\": \"BodyJsonField:$.req_msg_id\"   },   \"errorCondition\": \"$statusCode = 200 and $resultCode != null and $resultCode != 'OK'\",   \"errorCode\": \"resultCode\",   \"mappings\": [     {       \"code\": \"ROLE_NOT_EXISTS\",       \"statusCode\": 404,       \"errorMessage\": \"Role Not Exists, RequestId=${requestId}\"     },     {       \"code\": \"INVALID_PARAMETER\",       \"statusCode\": 400,       \"errorMessage\": \"Invalid Parameter, RequestId=${requestId}\"     }   ],   \"defaultMapping\": {     \"statusCode\": 500,     \"errorMessage\": \"Unknown Error, ${resultCode}, RequestId=${requestId}\"   } }",
						"plugin_type": "errorMapping",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
					"plugin_name": name + "_update",
					"plugin_data": "{   \\\"parameters\\\": {     \\\"statusCode\\\": \\\"StatusCode\\\",     \\\"resultCode\\\": \\\"BodyJsonField:$.result_code\\\",     \\\"requestId\\\": \\\"BodyJsonField:$.req_msg_id\\\"   },   \\\"errorCondition\\\": \\\"$statusCode = 200 and $resultCode != null and $resultCode != 'OK'\\\",   \\\"errorCode\\\": \\\"resultCode\\\",   \\\"mappings\\\": [     {       \\\"code\\\": \\\"ROLE_NOT_EXISTS\\\",       \\\"statusCode\\\": 404,       \\\"errorMessage\\\": \\\"Role Not Exists, RequestId=$${requestId}\\\"     },     {       \\\"code\\\": \\\"INVALID_PARAMETER\\\",       \\\"statusCode\\\": 400,       \\\"errorMessage\\\": \\\"Invalid Parameter, RequestId=$${requestId}\\\"     }   ],   \\\"defaultMapping\\\": {     \\\"statusCode\\\": 500,     \\\"errorMessage\\\": \\\"Unknown Error, $${resultCode}, RequestId=$${requestId}\\\"   } }",
					"plugin_type": "errorMapping",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name + "_update",
						"plugin_data": "{   \"parameters\": {     \"statusCode\": \"StatusCode\",     \"resultCode\": \"BodyJsonField:$.result_code\",     \"requestId\": \"BodyJsonField:$.req_msg_id\"   },   \"errorCondition\": \"$statusCode = 200 and $resultCode != null and $resultCode != 'OK'\",   \"errorCode\": \"resultCode\",   \"mappings\": [     {       \"code\": \"ROLE_NOT_EXISTS\",       \"statusCode\": 404,       \"errorMessage\": \"Role Not Exists, RequestId=${requestId}\"     },     {       \"code\": \"INVALID_PARAMETER\",       \"statusCode\": 400,       \"errorMessage\": \"Invalid Parameter, RequestId=${requestId}\"     }   ],   \"defaultMapping\": {     \"statusCode\": 500,     \"errorMessage\": \"Unknown Error, ${resultCode}, RequestId=${requestId}\"   } }",
						"plugin_type": "errorMapping",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

var AlicloudApiGatewayPluginMap6794 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudApiGatewayPluginBasicDependence6794(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 请求应答改写插件 6798
func TestAccAliCloudApiGatewayPlugin_basic6798(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6798)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6798)
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
					"plugin_name": name,
					"plugin_data": "{   \\\"addRequestParameterIfAbsent\\\": [     {       \\\"name\\\": \\\"userId\\\",       \\\"value\\\": 123456,       \\\"location\\\": \\\"query\\\"     }   ],   \\\"putRequestParameter\\\": [     {       \\\"name\\\": \\\"name\\\",       \\\"value\\\": null,       \\\"location\\\": \\\"header\\\"     }   ],   \\\"removeRequestParameter\\\": [     {       \\\"name\\\": \\\"address\\\",       \\\"location\\\": \\\"form\\\"     }   ],   \\\"setResponseStatusCode\\\": 200,   \\\"addResponseHeaderIfAbsent\\\": [     {       \\\"name\\\": \\\"age\\\",       \\\"value\\\": 18     }   ],   \\\"putResponseHeader\\\": [     {       \\\"name\\\": \\\"name\\\",       \\\"value\\\": \\\"Alice\\\"     }   ],   \\\"removeResponseHeader\\\": [     {       \\\"name\\\": \\\"phone\\\"     }   ] }",
					"plugin_type": "transformer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_name": name,
						"plugin_data": "{   \"addRequestParameterIfAbsent\": [     {       \"name\": \"userId\",       \"value\": 123456,       \"location\": \"query\"     }   ],   \"putRequestParameter\": [     {       \"name\": \"name\",       \"value\": null,       \"location\": \"header\"     }   ],   \"removeRequestParameter\": [     {       \"name\": \"address\",       \"location\": \"form\"     }   ],   \"setResponseStatusCode\": 200,   \"addResponseHeaderIfAbsent\": [     {       \"name\": \"age\",       \"value\": 18     }   ],   \"putResponseHeader\": [     {       \"name\": \"name\",       \"value\": \"Alice\"     }   ],   \"removeResponseHeader\": [     {       \"name\": \"phone\"     }   ] }",
						"plugin_type": "transformer",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
					"plugin_name": name + "_update",
					"plugin_data": "{   \\\"addRequestParameterIfAbsent\\\": [     {       \\\"name\\\": \\\"userId\\\",       \\\"value\\\": 123456,       \\\"location\\\": \\\"query\\\"     }   ],   \\\"putRequestParameter\\\": [     {       \\\"name\\\": \\\"name\\\",       \\\"value\\\": null,       \\\"location\\\": \\\"header\\\"     }   ],   \\\"removeRequestParameter\\\": [     {       \\\"name\\\": \\\"address\\\",       \\\"location\\\": \\\"form\\\"     }   ],   \\\"setResponseStatusCode\\\": 200,   \\\"addResponseHeaderIfAbsent\\\": [     {       \\\"name\\\": \\\"age\\\",       \\\"value\\\": 18     }   ],   \\\"putResponseHeader\\\": [     {       \\\"name\\\": \\\"name\\\",       \\\"value\\\": \\\"Alice\\\"     }   ],   \\\"removeResponseHeader\\\": [     {       \\\"name\\\": \\\"phone\\\"     }   ] }",
					"plugin_type": "transformer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name + "_update",
						"plugin_data": "{   \"addRequestParameterIfAbsent\": [     {       \"name\": \"userId\",       \"value\": 123456,       \"location\": \"query\"     }   ],   \"putRequestParameter\": [     {       \"name\": \"name\",       \"value\": null,       \"location\": \"header\"     }   ],   \"removeRequestParameter\": [     {       \"name\": \"address\",       \"location\": \"form\"     }   ],   \"setResponseStatusCode\": 200,   \"addResponseHeaderIfAbsent\": [     {       \"name\": \"age\",       \"value\": 18     }   ],   \"putResponseHeader\": [     {       \"name\": \"name\",       \"value\": \"Alice\"     }   ],   \"removeResponseHeader\": [     {       \"name\": \"phone\"     }   ] }",
						"plugin_type": "transformer",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

var AlicloudApiGatewayPluginMap6798 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudApiGatewayPluginBasicDependence6798(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 断路器插件 6795
func TestAccAliCloudApiGatewayPlugin_basic6795(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6795)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6795)
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
					"plugin_name": name,
					"plugin_data": "{   \\\"timeoutThreshold\\\": 15,   \\\"windowInSeconds\\\": 30,   \\\"openTimeoutSeconds\\\": 15,   \\\"downgradeBackend\\\": {     \\\"type\\\": \\\"mock\\\",     \\\"statusCode\\\": 302,     \\\"body\\\": \\\"<result>\\\\n  <errorCode>I's a teapot</errorCode>\\\\n</result>\\\\n\\\"   } }",
					"plugin_type": "circuitBreaker",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_name": name,
						"plugin_data": "{   \"timeoutThreshold\": 15,   \"windowInSeconds\": 30,   \"openTimeoutSeconds\": 15,   \"downgradeBackend\": {     \"type\": \"mock\",     \"statusCode\": 302,     \"body\": \"<result>\\n  <errorCode>I's a teapot</errorCode>\\n</result>\\n\"   } }",
						"plugin_type": "circuitBreaker",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
					"plugin_name": name + "_update",
					"plugin_data": "{   \\\"timeoutThreshold\\\": 15,   \\\"windowInSeconds\\\": 30,   \\\"openTimeoutSeconds\\\": 15,   \\\"downgradeBackend\\\": {     \\\"type\\\": \\\"mock\\\",     \\\"statusCode\\\": 302,     \\\"body\\\": \\\"<result>\\\\n  <errorCode>I's a teapot</errorCode>\\\\n</result>\\\\n\\\"   } }",
					"plugin_type": "circuitBreaker",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name + "_update",
						"plugin_data": "{   \"timeoutThreshold\": 15,   \"windowInSeconds\": 30,   \"openTimeoutSeconds\": 15,   \"downgradeBackend\": {     \"type\": \"mock\",     \"statusCode\": 302,     \"body\": \"<result>\\n  <errorCode>I's a teapot</errorCode>\\n</result>\\n\"   } }",
						"plugin_type": "circuitBreaker",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

var AlicloudApiGatewayPluginMap6795 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudApiGatewayPluginBasicDependence6795(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case ip控制插件 5922
func TestAccAliCloudApiGatewayPlugin_basic5922(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap5922)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence5922)
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
					"plugin_name": name,
					"plugin_data": "{\\\"type\\\":\\\"ALLOW\\\",\\\"items\\\":[{\\\"blocks\\\":[\\\"79.11.12.2\\\"]}]}",
					"plugin_type": "ipControl",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_name": name,
						"plugin_data": "{\"type\":\"ALLOW\",\"items\":[{\"blocks\":[\"79.11.12.2\"]}]}",
						"plugin_type": "ipControl",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tetetete",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tetetete",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plugin_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plugin_data": "{\\\"type\\\":\\\"ALLOW\\\",\\\"items\\\":[{\\\"blocks\\\":[\\\"79.11.12.2\\\",\\\"127.0.0.1\\\"]}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_data": "{\"type\":\"ALLOW\",\"items\":[{\"blocks\":[\"79.11.12.2\",\"127.0.0.1\"]}]}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
					"plugin_name": name + "_update",
					"plugin_data": "{\\\"type\\\":\\\"ALLOW\\\",\\\"items\\\":[{\\\"blocks\\\":[\\\"79.11.12.2\\\"]}]}",
					"plugin_type": "ipControl",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name + "_update",
						"plugin_data": "{\"type\":\"ALLOW\",\"items\":[{\"blocks\":[\"79.11.12.2\"]}]}",
						"plugin_type": "ipControl",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

var AlicloudApiGatewayPluginMap5922 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudApiGatewayPluginBasicDependence5922(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 日志脱敏插件 6797
func TestAccAliCloudApiGatewayPlugin_basic6797(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6797)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6797)
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
					"plugin_name": name,
					"plugin_data": "{   \\\"rules\\\": [     {       \\\"name\\\": \\\"request_query\\\",       \\\"location\\\": \\\"REQUEST_QUERY\\\",       \\\"parameters\\\": [         \\\"userid\\\",         \\\"name\\\"       ],       \\\"policy\\\": \\\"KEEP_LEFT:4\\\"     },     {       \\\"name\\\": \\\"request_header\\\",       \\\"location\\\": \\\"REQUEST_HEADER\\\",       \\\"parameters\\\": [         \\\"userid\\\",         \\\"name\\\"       ],       \\\"policy\\\": \\\"KEEP_CENTER:4,5\\\"     },     {       \\\"name\\\": \\\"request_body_HEX\\\",       \\\"location\\\": \\\"REQUEST_BODY\\\",       \\\"matchMode\\\": \\\"HEX:10\\\",       \\\"policy\\\": \\\"ALL\\\"     },     {       \\\"name\\\": \\\"response_body_EMAIL\\\",       \\\"location\\\": \\\"RESPONSE_BODY\\\",       \\\"matchMode\\\": \\\"EMAIL\\\",       \\\"policy\\\": \\\"KEEP_RIGHT:7\\\"     }   ] }",
					"plugin_type": "logMask",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_name": name,
						"plugin_data": "{   \"rules\": [     {       \"name\": \"request_query\",       \"location\": \"REQUEST_QUERY\",       \"parameters\": [         \"userid\",         \"name\"       ],       \"policy\": \"KEEP_LEFT:4\"     },     {       \"name\": \"request_header\",       \"location\": \"REQUEST_HEADER\",       \"parameters\": [         \"userid\",         \"name\"       ],       \"policy\": \"KEEP_CENTER:4,5\"     },     {       \"name\": \"request_body_HEX\",       \"location\": \"REQUEST_BODY\",       \"matchMode\": \"HEX:10\",       \"policy\": \"ALL\"     },     {       \"name\": \"response_body_EMAIL\",       \"location\": \"RESPONSE_BODY\",       \"matchMode\": \"EMAIL\",       \"policy\": \"KEEP_RIGHT:7\"     }   ] }",
						"plugin_type": "logMask",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
					"plugin_name": name + "_update",
					"plugin_data": "{   \\\"rules\\\": [     {       \\\"name\\\": \\\"request_query\\\",       \\\"location\\\": \\\"REQUEST_QUERY\\\",       \\\"parameters\\\": [         \\\"userid\\\",         \\\"name\\\"       ],       \\\"policy\\\": \\\"KEEP_LEFT:4\\\"     },     {       \\\"name\\\": \\\"request_header\\\",       \\\"location\\\": \\\"REQUEST_HEADER\\\",       \\\"parameters\\\": [         \\\"userid\\\",         \\\"name\\\"       ],       \\\"policy\\\": \\\"KEEP_CENTER:4,5\\\"     },     {       \\\"name\\\": \\\"request_body_HEX\\\",       \\\"location\\\": \\\"REQUEST_BODY\\\",       \\\"matchMode\\\": \\\"HEX:10\\\",       \\\"policy\\\": \\\"ALL\\\"     },     {       \\\"name\\\": \\\"response_body_EMAIL\\\",       \\\"location\\\": \\\"RESPONSE_BODY\\\",       \\\"matchMode\\\": \\\"EMAIL\\\",       \\\"policy\\\": \\\"KEEP_RIGHT:7\\\"     }   ] }",
					"plugin_type": "logMask",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name + "_update",
						"plugin_data": "{   \"rules\": [     {       \"name\": \"request_query\",       \"location\": \"REQUEST_QUERY\",       \"parameters\": [         \"userid\",         \"name\"       ],       \"policy\": \"KEEP_LEFT:4\"     },     {       \"name\": \"request_header\",       \"location\": \"REQUEST_HEADER\",       \"parameters\": [         \"userid\",         \"name\"       ],       \"policy\": \"KEEP_CENTER:4,5\"     },     {       \"name\": \"request_body_HEX\",       \"location\": \"REQUEST_BODY\",       \"matchMode\": \"HEX:10\",       \"policy\": \"ALL\"     },     {       \"name\": \"response_body_EMAIL\",       \"location\": \"RESPONSE_BODY\",       \"matchMode\": \"EMAIL\",       \"policy\": \"KEEP_RIGHT:7\"     }   ] }",
						"plugin_type": "logMask",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

var AlicloudApiGatewayPluginMap6797 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudApiGatewayPluginBasicDependence6797(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case jwt插件 6788
func TestAccAliCloudApiGatewayPlugin_basic6788(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6788)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6788)
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
					"plugin_name": name,
					"plugin_data": "{   \\\"parameter\\\": \\\"X-Token\\\",   \\\"parameterLocation\\\": \\\"header\\\",   \\\"claimParameters\\\": [     {       \\\"claimName\\\": \\\"aud\\\",       \\\"parameterName\\\": \\\"X-Aud\\\",       \\\"location\\\": \\\"header\\\"     },     {       \\\"claimName\\\": \\\"userId\\\",       \\\"parameterName\\\": \\\"userId\\\",       \\\"location\\\": \\\"query\\\"     }   ],   \\\"preventJtiReplay\\\": false,   \\\"jwk\\\": {     \\\"kty\\\": \\\"RSA\\\",     \\\"e\\\": \\\"AQAB\\\",     \\\"use\\\": \\\"sig\\\",     \\\"kid\\\": \\\"O8fpdhrViq2zaaaBEWZITz\\\",     \\\"alg\\\": \\\"RS256\\\",     \\\"n\\\": \\\"qSVxcknOm0uCq5vGsOmaorPDzHUubBmZZ4UXj-9do7w9X1uKFXAnqfto4TepSNuYU2bA_-tzSLAGBsR-BqvT6w9SjxakeiyQpVmexxnDw5WZwpWenUAcYrfSPEoNU-0hAQwFYgqZwJQMN8ptxkd0170PFauwACOx4Hfr-9FPGy8NCoIO4MfLXzJ3mJ7xqgIZp3NIOGXz-GIAbCf13ii7kSStpYqN3L_zzpvXUAos1FJ9IPXRV84tIZpFVh2lmRh0h8ImK-vI42dwlD_hOIzayL1Xno2R0T-d5AwTSdnep7g-Fwu8-sj4cCRWq3bd61Zs2QOJ8iustH0vSRMYdP5oYQ\\\"   } }",
					"plugin_type": "jwtAuth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_name": name,
						"plugin_data": "{   \"parameter\": \"X-Token\",   \"parameterLocation\": \"header\",   \"claimParameters\": [     {       \"claimName\": \"aud\",       \"parameterName\": \"X-Aud\",       \"location\": \"header\"     },     {       \"claimName\": \"userId\",       \"parameterName\": \"userId\",       \"location\": \"query\"     }   ],   \"preventJtiReplay\": false,   \"jwk\": {     \"kty\": \"RSA\",     \"e\": \"AQAB\",     \"use\": \"sig\",     \"kid\": \"O8fpdhrViq2zaaaBEWZITz\",     \"alg\": \"RS256\",     \"n\": \"qSVxcknOm0uCq5vGsOmaorPDzHUubBmZZ4UXj-9do7w9X1uKFXAnqfto4TepSNuYU2bA_-tzSLAGBsR-BqvT6w9SjxakeiyQpVmexxnDw5WZwpWenUAcYrfSPEoNU-0hAQwFYgqZwJQMN8ptxkd0170PFauwACOx4Hfr-9FPGy8NCoIO4MfLXzJ3mJ7xqgIZp3NIOGXz-GIAbCf13ii7kSStpYqN3L_zzpvXUAos1FJ9IPXRV84tIZpFVh2lmRh0h8ImK-vI42dwlD_hOIzayL1Xno2R0T-d5AwTSdnep7g-Fwu8-sj4cCRWq3bd61Zs2QOJ8iustH0vSRMYdP5oYQ\"   } }",
						"plugin_type": "jwtAuth",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
					"plugin_name": name + "_update",
					"plugin_data": "{   \\\"parameter\\\": \\\"X-Token\\\",   \\\"parameterLocation\\\": \\\"header\\\",   \\\"claimParameters\\\": [     {       \\\"claimName\\\": \\\"aud\\\",       \\\"parameterName\\\": \\\"X-Aud\\\",       \\\"location\\\": \\\"header\\\"     },     {       \\\"claimName\\\": \\\"userId\\\",       \\\"parameterName\\\": \\\"userId\\\",       \\\"location\\\": \\\"query\\\"     }   ],   \\\"preventJtiReplay\\\": false,   \\\"jwk\\\": {     \\\"kty\\\": \\\"RSA\\\",     \\\"e\\\": \\\"AQAB\\\",     \\\"use\\\": \\\"sig\\\",     \\\"kid\\\": \\\"O8fpdhrViq2zaaaBEWZITz\\\",     \\\"alg\\\": \\\"RS256\\\",     \\\"n\\\": \\\"qSVxcknOm0uCq5vGsOmaorPDzHUubBmZZ4UXj-9do7w9X1uKFXAnqfto4TepSNuYU2bA_-tzSLAGBsR-BqvT6w9SjxakeiyQpVmexxnDw5WZwpWenUAcYrfSPEoNU-0hAQwFYgqZwJQMN8ptxkd0170PFauwACOx4Hfr-9FPGy8NCoIO4MfLXzJ3mJ7xqgIZp3NIOGXz-GIAbCf13ii7kSStpYqN3L_zzpvXUAos1FJ9IPXRV84tIZpFVh2lmRh0h8ImK-vI42dwlD_hOIzayL1Xno2R0T-d5AwTSdnep7g-Fwu8-sj4cCRWq3bd61Zs2QOJ8iustH0vSRMYdP5oYQ\\\"   } }",
					"plugin_type": "jwtAuth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name + "_update",
						"plugin_data": "{   \"parameter\": \"X-Token\",   \"parameterLocation\": \"header\",   \"claimParameters\": [     {       \"claimName\": \"aud\",       \"parameterName\": \"X-Aud\",       \"location\": \"header\"     },     {       \"claimName\": \"userId\",       \"parameterName\": \"userId\",       \"location\": \"query\"     }   ],   \"preventJtiReplay\": false,   \"jwk\": {     \"kty\": \"RSA\",     \"e\": \"AQAB\",     \"use\": \"sig\",     \"kid\": \"O8fpdhrViq2zaaaBEWZITz\",     \"alg\": \"RS256\",     \"n\": \"qSVxcknOm0uCq5vGsOmaorPDzHUubBmZZ4UXj-9do7w9X1uKFXAnqfto4TepSNuYU2bA_-tzSLAGBsR-BqvT6w9SjxakeiyQpVmexxnDw5WZwpWenUAcYrfSPEoNU-0hAQwFYgqZwJQMN8ptxkd0170PFauwACOx4Hfr-9FPGy8NCoIO4MfLXzJ3mJ7xqgIZp3NIOGXz-GIAbCf13ii7kSStpYqN3L_zzpvXUAos1FJ9IPXRV84tIZpFVh2lmRh0h8ImK-vI42dwlD_hOIzayL1Xno2R0T-d5AwTSdnep7g-Fwu8-sj4cCRWq3bd61Zs2QOJ8iustH0vSRMYdP5oYQ\"   } }",
						"plugin_type": "jwtAuth",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

var AlicloudApiGatewayPluginMap6788 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudApiGatewayPluginBasicDependence6788(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 流控插件 6787
func TestAccAliCloudApiGatewayPlugin_basic6787(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6787)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6787)
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
					"plugin_name": name,
					"plugin_data": "{   \\\"unit\\\": \\\"SECOND\\\",   \\\"apiDefault\\\": 1000,   \\\"userDefault\\\": 30,   \\\"appDefault\\\": 30,   \\\"specials\\\": [     {       \\\"type\\\": \\\"APP\\\",       \\\"policies\\\": [         {           \\\"key\\\": 10123123,           \\\"value\\\": 10         },         {           \\\"key\\\": 10123123,           \\\"value\\\": 10         }       ]     },     {       \\\"type\\\": \\\"USER\\\",       \\\"policies\\\": [         {           \\\"key\\\": 123455,           \\\"value\\\": 100         }       ]     }   ] }",
					"plugin_type": "trafficControl",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_name": name,
						"plugin_data": "{   \"unit\": \"SECOND\",   \"apiDefault\": 1000,   \"userDefault\": 30,   \"appDefault\": 30,   \"specials\": [     {       \"type\": \"APP\",       \"policies\": [         {           \"key\": 10123123,           \"value\": 10         },         {           \"key\": 10123123,           \"value\": 10         }       ]     },     {       \"type\": \"USER\",       \"policies\": [         {           \"key\": 123455,           \"value\": 100         }       ]     }   ] }",
						"plugin_type": "trafficControl",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
					"plugin_name": name + "_update",
					"plugin_data": "{   \\\"unit\\\": \\\"SECOND\\\",   \\\"apiDefault\\\": 1000,   \\\"userDefault\\\": 30,   \\\"appDefault\\\": 30,   \\\"specials\\\": [     {       \\\"type\\\": \\\"APP\\\",       \\\"policies\\\": [         {           \\\"key\\\": 10123123,           \\\"value\\\": 10         },         {           \\\"key\\\": 10123123,           \\\"value\\\": 10         }       ]     },     {       \\\"type\\\": \\\"USER\\\",       \\\"policies\\\": [         {           \\\"key\\\": 123455,           \\\"value\\\": 100         }       ]     }   ] }",
					"plugin_type": "trafficControl",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name + "_update",
						"plugin_data": "{   \"unit\": \"SECOND\",   \"apiDefault\": 1000,   \"userDefault\": 30,   \"appDefault\": 30,   \"specials\": [     {       \"type\": \"APP\",       \"policies\": [         {           \"key\": 10123123,           \"value\": 10         },         {           \"key\": 10123123,           \"value\": 10         }       ]     },     {       \"type\": \"USER\",       \"policies\": [         {           \"key\": 123455,           \"value\": 100         }       ]     }   ] }",
						"plugin_type": "trafficControl",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

var AlicloudApiGatewayPluginMap6787 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudApiGatewayPluginBasicDependence6787(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 跨域插件 6790
func TestAccAliCloudApiGatewayPlugin_basic6790(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6790)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6790)
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
					"plugin_name": name,
					"plugin_data": "{   \\\"allowOrigins\\\": \\\"api.foo.com\\\",   \\\"allowMethods\\\": \\\"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\\\",   \\\"allowHeaders\\\": \\\"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\\\",   \\\"exposeHeaders\\\": \\\"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\\\",   \\\"maxAge\\\": 172800,   \\\"allowCredentials\\\": true }",
					"plugin_type": "cors",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_name": name,
						"plugin_data": "{   \"allowOrigins\": \"api.foo.com\",   \"allowMethods\": \"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\",   \"allowHeaders\": \"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\",   \"exposeHeaders\": \"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\",   \"maxAge\": 172800,   \"allowCredentials\": true }",
						"plugin_type": "cors",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
					"plugin_name": name + "_update",
					"plugin_data": "{   \\\"allowOrigins\\\": \\\"api.foo.com\\\",   \\\"allowMethods\\\": \\\"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\\\",   \\\"allowHeaders\\\": \\\"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\\\",   \\\"exposeHeaders\\\": \\\"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\\\",   \\\"maxAge\\\": 172800,   \\\"allowCredentials\\\": true }",
					"plugin_type": "cors",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name + "_update",
						"plugin_data": "{   \"allowOrigins\": \"api.foo.com\",   \"allowMethods\": \"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\",   \"allowHeaders\": \"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\",   \"exposeHeaders\": \"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\",   \"maxAge\": 172800,   \"allowCredentials\": true }",
						"plugin_type": "cors",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

var AlicloudApiGatewayPluginMap6790 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudApiGatewayPluginBasicDependence6790(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 后端路由插件 6792  twin
func TestAccAliCloudApiGatewayPlugin_basic6792_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6792)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6792)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"routes\\\": [     {       \\\"name\\\": \\\"Vip\\\",       \\\"condition\\\": \\\"$CaAppId = 123456\\\",       \\\"backend\\\": {         \\\"type\\\": \\\"HTTP-VPC\\\",         \\\"vpcAccessName\\\": \\\"slbAccessForVip\\\"       }     },     {       \\\"name\\\": \\\"MockForOldClient\\\",       \\\"condition\\\": \\\"$ClientVersion < '2.0.5'\\\",       \\\"backend\\\": {         \\\"type\\\": \\\"MOCK\\\",         \\\"statusCode\\\": 400,         \\\"mockBody\\\": \\\"This version is not supported!!!\\\"       }     },     {       \\\"name\\\": \\\"BlueGreenPercent05\\\",       \\\"condition\\\": \\\"1 = 1\\\",       \\\"backend\\\": {         \\\"type\\\": \\\"HTTP\\\",         \\\"address\\\": \\\"https://beta-version.api.foo.com\\\"       },       \\\"constant-parameters\\\": [         {           \\\"name\\\": \\\"x-route-blue-green\\\",           \\\"location\\\": \\\"header\\\",           \\\"value\\\": \\\"route-blue-green\\\"         }       ]     }   ] }",
					"plugin_type": "routing",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "test",
						"plugin_name":  name,
						"plugin_data":  "{   \"routes\": [     {       \"name\": \"Vip\",       \"condition\": \"$CaAppId = 123456\",       \"backend\": {         \"type\": \"HTTP-VPC\",         \"vpcAccessName\": \"slbAccessForVip\"       }     },     {       \"name\": \"MockForOldClient\",       \"condition\": \"$ClientVersion < '2.0.5'\",       \"backend\": {         \"type\": \"MOCK\",         \"statusCode\": 400,         \"mockBody\": \"This version is not supported!!!\"       }     },     {       \"name\": \"BlueGreenPercent05\",       \"condition\": \"1 = 1\",       \"backend\": {         \"type\": \"HTTP\",         \"address\": \"https://beta-version.api.foo.com\"       },       \"constant-parameters\": [         {           \"name\": \"x-route-blue-green\",           \"location\": \"header\",           \"value\": \"route-blue-green\"         }       ]     }   ] }",
						"plugin_type":  "routing",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case BasicAuth插件 6789  twin
func TestAccAliCloudApiGatewayPlugin_basic6789_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6789)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6789)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"users\\\": [     {       \\\"username\\\": \\\"alice\\\",       \\\"password\\\": 123456     },     {       \\\"username\\\": \\\"bob\\\",       \\\"password\\\": 666666     },     {       \\\"username\\\": \\\"charlie\\\",       \\\"password\\\": 888888     },     {       \\\"username\\\": \\\"dave\\\",       \\\"password\\\": 111111     }   ] }",
					"plugin_type": "basicAuth",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "test",
						"plugin_name":  name,
						"plugin_data":  "{   \"users\": [     {       \"username\": \"alice\",       \"password\": 123456     },     {       \"username\": \"bob\",       \"password\": 666666     },     {       \"username\": \"charlie\",       \"password\": 888888     },     {       \"username\": \"dave\",       \"password\": 111111     }   ] }",
						"plugin_type":  "basicAuth",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 参数访问控制插件 6793  twin
func TestAccAliCloudApiGatewayPlugin_basic6793_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6793)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6793)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"parameters\\\": {     \\\"userId\\\": \\\"Token:userId\\\",     \\\"userType\\\": \\\"Token:userType\\\",     \\\"pathUserId\\\": \\\"path:userId\\\"   },   \\\"rules\\\": [     {       \\\"name\\\": \\\"admin\\\",       \\\"condition\\\": \\\"$userType = 'admin'\\\",       \\\"ifTrue\\\": \\\"ALLOW\\\"     },     {       \\\"name\\\": \\\"user\\\",       \\\"condition\\\": \\\"$userId = $pathUserId\\\",       \\\"ifFalse\\\": \\\"DENY\\\",       \\\"statusCode\\\": 403,       \\\"errorMessage\\\": \\\"Path not match $${userId} vs /$${pathUserId}\\\",       \\\"responseHeaders\\\": {         \\\"Content-Type\\\": \\\"application/xml\\\"       },       \\\"responseBody\\\": \\\"<Reason>Path not match $${userId} vs /$${pathUserId}</Reason>\\\\n\\\"     }   ] }",
					"plugin_type": "accessControl",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "test",
						"plugin_name":  name,
						"plugin_data":  "{   \"parameters\": {     \"userId\": \"Token:userId\",     \"userType\": \"Token:userType\",     \"pathUserId\": \"path:userId\"   },   \"rules\": [     {       \"name\": \"admin\",       \"condition\": \"$userType = 'admin'\",       \"ifTrue\": \"ALLOW\"     },     {       \"name\": \"user\",       \"condition\": \"$userId = $pathUserId\",       \"ifFalse\": \"DENY\",       \"statusCode\": 403,       \"errorMessage\": \"Path not match ${userId} vs /${pathUserId}\",       \"responseHeaders\": {         \"Content-Type\": \"application/xml\"       },       \"responseBody\": \"<Reason>Path not match ${userId} vs /${pathUserId}</Reason>\\n\"     }   ] }",
						"plugin_type":  "accessControl",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 第三方鉴权插件 6796  twin
func TestAccAliCloudApiGatewayPlugin_basic6796_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6796)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6796)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"parameters\\\": {     \\\"statusCode\\\": \\\"StatusCode\\\"   },   \\\"authUriType\\\": \\\"HTTP\\\",   \\\"authUri\\\": {     \\\"address\\\": \\\"http://your-auth-domain.com:8080\\\",     \\\"path\\\": \\\"/your/authPath\\\",     \\\"timeout\\\": 7000,     \\\"method\\\": \\\"POST\\\"   },   \\\"passThroughBody\\\": false,   \\\"cachedTimeBySecond\\\": 10,   \\\"authParameters\\\": [     {       \\\"targetParameterName\\\": \\\"x-userId\\\",       \\\"sourceParameterName\\\": \\\"userId\\\",       \\\"targetLocation\\\": \\\"form\\\",       \\\"sourceLocation\\\": \\\"query\\\"     },     {       \\\"targetParameterName\\\": \\\"x-passwoed\\\",       \\\"sourceParameterName\\\": \\\"password\\\",       \\\"targetLocation\\\": \\\"form\\\",       \\\"sourceLocation\\\": \\\"query\\\"     }   ],   \\\"successCondition\\\": \\\"$${statusCode} = 200\\\",   \\\"errorMessage\\\": \\\"auth failed\\\",   \\\"errorStatusCode\\\": 401 }",
					"plugin_type": "remoteAuth",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "test",
						"plugin_name":  name,
						"plugin_data":  "{   \"parameters\": {     \"statusCode\": \"StatusCode\"   },   \"authUriType\": \"HTTP\",   \"authUri\": {     \"address\": \"http://your-auth-domain.com:8080\",     \"path\": \"/your/authPath\",     \"timeout\": 7000,     \"method\": \"POST\"   },   \"passThroughBody\": false,   \"cachedTimeBySecond\": 10,   \"authParameters\": [     {       \"targetParameterName\": \"x-userId\",       \"sourceParameterName\": \"userId\",       \"targetLocation\": \"form\",       \"sourceLocation\": \"query\"     },     {       \"targetParameterName\": \"x-passwoed\",       \"sourceParameterName\": \"password\",       \"targetLocation\": \"form\",       \"sourceLocation\": \"query\"     }   ],   \"successCondition\": \"${statusCode} = 200\",   \"errorMessage\": \"auth failed\",   \"errorStatusCode\": 401 }",
						"plugin_type":  "remoteAuth",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 缓存插件 6791  twin
func TestAccAliCloudApiGatewayPlugin_basic6791_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6791)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6791)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"varyByApp\\\": false,   \\\"varyByParameters\\\": [     \\\"userId\\\"   ],   \\\"varyByHeaders\\\": [     \\\"Accept\\\"   ],   \\\"clientCacheControl\\\": {     \\\"mode\\\": \\\"app\\\",     \\\"apps\\\": [       1992323,       1239922     ]   },   \\\"cacheableHeaders\\\": [     \\\"X-Customer-Token\\\"   ] }",
					"plugin_type": "caching",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "test",
						"plugin_name":  name,
						"plugin_data":  "{   \"varyByApp\": false,   \"varyByParameters\": [     \"userId\"   ],   \"varyByHeaders\": [     \"Accept\"   ],   \"clientCacheControl\": {     \"mode\": \"app\",     \"apps\": [       1992323,       1239922     ]   },   \"cacheableHeaders\": [     \"X-Customer-Token\"   ] }",
						"plugin_type":  "caching",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 后端签名插件 6799  twin
func TestAccAliCloudApiGatewayPlugin_basic6799_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6799)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6799)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"type\\\": \\\"APIGW_BACKEND\\\",   \\\"key\\\": \\\"SampleKey\\\",   \\\"secret\\\": \\\"SampleSecret\\\" }",
					"plugin_type": "backendSignature",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "test",
						"plugin_name":  name,
						"plugin_data":  "{   \"type\": \"APIGW_BACKEND\",   \"key\": \"SampleKey\",   \"secret\": \"SampleSecret\" }",
						"plugin_type":  "backendSignature",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 错误码映射插件 6794  twin
func TestAccAliCloudApiGatewayPlugin_basic6794_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6794)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6794)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"parameters\\\": {     \\\"statusCode\\\": \\\"StatusCode\\\",     \\\"resultCode\\\": \\\"BodyJsonField:$.result_code\\\",     \\\"requestId\\\": \\\"BodyJsonField:$.req_msg_id\\\"   },   \\\"errorCondition\\\": \\\"$statusCode = 200 and $resultCode != null and $resultCode != 'OK'\\\",   \\\"errorCode\\\": \\\"resultCode\\\",   \\\"mappings\\\": [     {       \\\"code\\\": \\\"ROLE_NOT_EXISTS\\\",       \\\"statusCode\\\": 404,       \\\"errorMessage\\\": \\\"Role Not Exists, RequestId=$${requestId}\\\"     },     {       \\\"code\\\": \\\"INVALID_PARAMETER\\\",       \\\"statusCode\\\": 400,       \\\"errorMessage\\\": \\\"Invalid Parameter, RequestId=$${requestId}\\\"     }   ],   \\\"defaultMapping\\\": {     \\\"statusCode\\\": 500,     \\\"errorMessage\\\": \\\"Unknown Error, $${resultCode}, RequestId=$${requestId}\\\"   } }",
					"plugin_type": "errorMapping",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "test",
						"plugin_name":  name,
						"plugin_data":  "{   \"parameters\": {     \"statusCode\": \"StatusCode\",     \"resultCode\": \"BodyJsonField:$.result_code\",     \"requestId\": \"BodyJsonField:$.req_msg_id\"   },   \"errorCondition\": \"$statusCode = 200 and $resultCode != null and $resultCode != 'OK'\",   \"errorCode\": \"resultCode\",   \"mappings\": [     {       \"code\": \"ROLE_NOT_EXISTS\",       \"statusCode\": 404,       \"errorMessage\": \"Role Not Exists, RequestId=${requestId}\"     },     {       \"code\": \"INVALID_PARAMETER\",       \"statusCode\": 400,       \"errorMessage\": \"Invalid Parameter, RequestId=${requestId}\"     }   ],   \"defaultMapping\": {     \"statusCode\": 500,     \"errorMessage\": \"Unknown Error, ${resultCode}, RequestId=${requestId}\"   } }",
						"plugin_type":  "errorMapping",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 请求应答改写插件 6798  twin
func TestAccAliCloudApiGatewayPlugin_basic6798_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6798)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6798)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"addRequestParameterIfAbsent\\\": [     {       \\\"name\\\": \\\"userId\\\",       \\\"value\\\": 123456,       \\\"location\\\": \\\"query\\\"     }   ],   \\\"putRequestParameter\\\": [     {       \\\"name\\\": \\\"name\\\",       \\\"value\\\": null,       \\\"location\\\": \\\"header\\\"     }   ],   \\\"removeRequestParameter\\\": [     {       \\\"name\\\": \\\"address\\\",       \\\"location\\\": \\\"form\\\"     }   ],   \\\"setResponseStatusCode\\\": 200,   \\\"addResponseHeaderIfAbsent\\\": [     {       \\\"name\\\": \\\"age\\\",       \\\"value\\\": 18     }   ],   \\\"putResponseHeader\\\": [     {       \\\"name\\\": \\\"name\\\",       \\\"value\\\": \\\"Alice\\\"     }   ],   \\\"removeResponseHeader\\\": [     {       \\\"name\\\": \\\"phone\\\"     }   ] }",
					"plugin_type": "transformer",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "test",
						"plugin_name":  name,
						"plugin_data":  "{   \"addRequestParameterIfAbsent\": [     {       \"name\": \"userId\",       \"value\": 123456,       \"location\": \"query\"     }   ],   \"putRequestParameter\": [     {       \"name\": \"name\",       \"value\": null,       \"location\": \"header\"     }   ],   \"removeRequestParameter\": [     {       \"name\": \"address\",       \"location\": \"form\"     }   ],   \"setResponseStatusCode\": 200,   \"addResponseHeaderIfAbsent\": [     {       \"name\": \"age\",       \"value\": 18     }   ],   \"putResponseHeader\": [     {       \"name\": \"name\",       \"value\": \"Alice\"     }   ],   \"removeResponseHeader\": [     {       \"name\": \"phone\"     }   ] }",
						"plugin_type":  "transformer",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 断路器插件 6795  twin
func TestAccAliCloudApiGatewayPlugin_basic6795_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6795)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6795)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"timeoutThreshold\\\": 15,   \\\"windowInSeconds\\\": 30,   \\\"openTimeoutSeconds\\\": 15,   \\\"downgradeBackend\\\": {     \\\"type\\\": \\\"mock\\\",     \\\"statusCode\\\": 302,     \\\"body\\\": \\\"<result>\\\\n  <errorCode>I's a teapot</errorCode>\\\\n</result>\\\\n\\\"   } }",
					"plugin_type": "circuitBreaker",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "test",
						"plugin_name":  name,
						"plugin_data":  "{   \"timeoutThreshold\": 15,   \"windowInSeconds\": 30,   \"openTimeoutSeconds\": 15,   \"downgradeBackend\": {     \"type\": \"mock\",     \"statusCode\": 302,     \"body\": \"<result>\\n  <errorCode>I's a teapot</errorCode>\\n</result>\\n\"   } }",
						"plugin_type":  "circuitBreaker",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case ip控制插件 5922  twin
func TestAccAliCloudApiGatewayPlugin_basic5922_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap5922)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence5922)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{\\\"type\\\":\\\"ALLOW\\\",\\\"items\\\":[{\\\"blocks\\\":[\\\"79.11.12.2\\\"]}]}",
					"plugin_type": "ipControl",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "test",
						"plugin_name":  name,
						"plugin_data":  "{\"type\":\"ALLOW\",\"items\":[{\"blocks\":[\"79.11.12.2\"]}]}",
						"plugin_type":  "ipControl",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 日志脱敏插件 6797  twin
func TestAccAliCloudApiGatewayPlugin_basic6797_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6797)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6797)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"rules\\\": [     {       \\\"name\\\": \\\"request_query\\\",       \\\"location\\\": \\\"REQUEST_QUERY\\\",       \\\"parameters\\\": [         \\\"userid\\\",         \\\"name\\\"       ],       \\\"policy\\\": \\\"KEEP_LEFT:4\\\"     },     {       \\\"name\\\": \\\"request_header\\\",       \\\"location\\\": \\\"REQUEST_HEADER\\\",       \\\"parameters\\\": [         \\\"userid\\\",         \\\"name\\\"       ],       \\\"policy\\\": \\\"KEEP_CENTER:4,5\\\"     },     {       \\\"name\\\": \\\"request_body_HEX\\\",       \\\"location\\\": \\\"REQUEST_BODY\\\",       \\\"matchMode\\\": \\\"HEX:10\\\",       \\\"policy\\\": \\\"ALL\\\"     },     {       \\\"name\\\": \\\"response_body_EMAIL\\\",       \\\"location\\\": \\\"RESPONSE_BODY\\\",       \\\"matchMode\\\": \\\"EMAIL\\\",       \\\"policy\\\": \\\"KEEP_RIGHT:7\\\"     }   ] }",
					"plugin_type": "logMask",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "test",
						"plugin_name":  name,
						"plugin_data":  "{   \"rules\": [     {       \"name\": \"request_query\",       \"location\": \"REQUEST_QUERY\",       \"parameters\": [         \"userid\",         \"name\"       ],       \"policy\": \"KEEP_LEFT:4\"     },     {       \"name\": \"request_header\",       \"location\": \"REQUEST_HEADER\",       \"parameters\": [         \"userid\",         \"name\"       ],       \"policy\": \"KEEP_CENTER:4,5\"     },     {       \"name\": \"request_body_HEX\",       \"location\": \"REQUEST_BODY\",       \"matchMode\": \"HEX:10\",       \"policy\": \"ALL\"     },     {       \"name\": \"response_body_EMAIL\",       \"location\": \"RESPONSE_BODY\",       \"matchMode\": \"EMAIL\",       \"policy\": \"KEEP_RIGHT:7\"     }   ] }",
						"plugin_type":  "logMask",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case jwt插件 6788  twin
func TestAccAliCloudApiGatewayPlugin_basic6788_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6788)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6788)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"parameter\\\": \\\"X-Token\\\",   \\\"parameterLocation\\\": \\\"header\\\",   \\\"claimParameters\\\": [     {       \\\"claimName\\\": \\\"aud\\\",       \\\"parameterName\\\": \\\"X-Aud\\\",       \\\"location\\\": \\\"header\\\"     },     {       \\\"claimName\\\": \\\"userId\\\",       \\\"parameterName\\\": \\\"userId\\\",       \\\"location\\\": \\\"query\\\"     }   ],   \\\"preventJtiReplay\\\": false,   \\\"jwk\\\": {     \\\"kty\\\": \\\"RSA\\\",     \\\"e\\\": \\\"AQAB\\\",     \\\"use\\\": \\\"sig\\\",     \\\"kid\\\": \\\"O8fpdhrViq2zaaaBEWZITz\\\",     \\\"alg\\\": \\\"RS256\\\",     \\\"n\\\": \\\"qSVxcknOm0uCq5vGsOmaorPDzHUubBmZZ4UXj-9do7w9X1uKFXAnqfto4TepSNuYU2bA_-tzSLAGBsR-BqvT6w9SjxakeiyQpVmexxnDw5WZwpWenUAcYrfSPEoNU-0hAQwFYgqZwJQMN8ptxkd0170PFauwACOx4Hfr-9FPGy8NCoIO4MfLXzJ3mJ7xqgIZp3NIOGXz-GIAbCf13ii7kSStpYqN3L_zzpvXUAos1FJ9IPXRV84tIZpFVh2lmRh0h8ImK-vI42dwlD_hOIzayL1Xno2R0T-d5AwTSdnep7g-Fwu8-sj4cCRWq3bd61Zs2QOJ8iustH0vSRMYdP5oYQ\\\"   } }",
					"plugin_type": "jwtAuth",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "test",
						"plugin_name":  name,
						"plugin_data":  "{   \"parameter\": \"X-Token\",   \"parameterLocation\": \"header\",   \"claimParameters\": [     {       \"claimName\": \"aud\",       \"parameterName\": \"X-Aud\",       \"location\": \"header\"     },     {       \"claimName\": \"userId\",       \"parameterName\": \"userId\",       \"location\": \"query\"     }   ],   \"preventJtiReplay\": false,   \"jwk\": {     \"kty\": \"RSA\",     \"e\": \"AQAB\",     \"use\": \"sig\",     \"kid\": \"O8fpdhrViq2zaaaBEWZITz\",     \"alg\": \"RS256\",     \"n\": \"qSVxcknOm0uCq5vGsOmaorPDzHUubBmZZ4UXj-9do7w9X1uKFXAnqfto4TepSNuYU2bA_-tzSLAGBsR-BqvT6w9SjxakeiyQpVmexxnDw5WZwpWenUAcYrfSPEoNU-0hAQwFYgqZwJQMN8ptxkd0170PFauwACOx4Hfr-9FPGy8NCoIO4MfLXzJ3mJ7xqgIZp3NIOGXz-GIAbCf13ii7kSStpYqN3L_zzpvXUAos1FJ9IPXRV84tIZpFVh2lmRh0h8ImK-vI42dwlD_hOIzayL1Xno2R0T-d5AwTSdnep7g-Fwu8-sj4cCRWq3bd61Zs2QOJ8iustH0vSRMYdP5oYQ\"   } }",
						"plugin_type":  "jwtAuth",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 流控插件 6787  twin
func TestAccAliCloudApiGatewayPlugin_basic6787_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6787)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6787)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"unit\\\": \\\"SECOND\\\",   \\\"apiDefault\\\": 1000,   \\\"userDefault\\\": 30,   \\\"appDefault\\\": 30,   \\\"specials\\\": [     {       \\\"type\\\": \\\"APP\\\",       \\\"policies\\\": [         {           \\\"key\\\": 10123123,           \\\"value\\\": 10         },         {           \\\"key\\\": 10123123,           \\\"value\\\": 10         }       ]     },     {       \\\"type\\\": \\\"USER\\\",       \\\"policies\\\": [         {           \\\"key\\\": 123455,           \\\"value\\\": 100         }       ]     }   ] }",
					"plugin_type": "trafficControl",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "test",
						"plugin_name":  name,
						"plugin_data":  "{   \"unit\": \"SECOND\",   \"apiDefault\": 1000,   \"userDefault\": 30,   \"appDefault\": 30,   \"specials\": [     {       \"type\": \"APP\",       \"policies\": [         {           \"key\": 10123123,           \"value\": 10         },         {           \"key\": 10123123,           \"value\": 10         }       ]     },     {       \"type\": \"USER\",       \"policies\": [         {           \"key\": 123455,           \"value\": 100         }       ]     }   ] }",
						"plugin_type":  "trafficControl",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 跨域插件 6790  twin
func TestAccAliCloudApiGatewayPlugin_basic6790_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6790)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6790)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"allowOrigins\\\": \\\"api.foo.com\\\",   \\\"allowMethods\\\": \\\"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\\\",   \\\"allowHeaders\\\": \\\"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\\\",   \\\"exposeHeaders\\\": \\\"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\\\",   \\\"maxAge\\\": 172800,   \\\"allowCredentials\\\": true }",
					"plugin_type": "cors",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "test",
						"plugin_name":  name,
						"plugin_data":  "{   \"allowOrigins\": \"api.foo.com\",   \"allowMethods\": \"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\",   \"allowHeaders\": \"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\",   \"exposeHeaders\": \"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\",   \"maxAge\": 172800,   \"allowCredentials\": true }",
						"plugin_type":  "cors",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 后端路由插件 6792  raw
func TestAccAliCloudApiGatewayPlugin_basic6792_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6792)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6792)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"routes\\\": [     {       \\\"name\\\": \\\"Vip\\\",       \\\"condition\\\": \\\"$CaAppId = 123456\\\",       \\\"backend\\\": {         \\\"type\\\": \\\"HTTP-VPC\\\",         \\\"vpcAccessName\\\": \\\"slbAccessForVip\\\"       }     },     {       \\\"name\\\": \\\"MockForOldClient\\\",       \\\"condition\\\": \\\"$ClientVersion < '2.0.5'\\\",       \\\"backend\\\": {         \\\"type\\\": \\\"MOCK\\\",         \\\"statusCode\\\": 400,         \\\"mockBody\\\": \\\"This version is not supported!!!\\\"       }     },     {       \\\"name\\\": \\\"BlueGreenPercent05\\\",       \\\"condition\\\": \\\"1 = 1\\\",       \\\"backend\\\": {         \\\"type\\\": \\\"HTTP\\\",         \\\"address\\\": \\\"https://beta-version.api.foo.com\\\"       },       \\\"constant-parameters\\\": [         {           \\\"name\\\": \\\"x-route-blue-green\\\",           \\\"location\\\": \\\"header\\\",           \\\"value\\\": \\\"route-blue-green\\\"         }       ]     }   ] }",
					"plugin_type": "routing",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name,
						"plugin_data": "{   \"routes\": [     {       \"name\": \"Vip\",       \"condition\": \"$CaAppId = 123456\",       \"backend\": {         \"type\": \"HTTP-VPC\",         \"vpcAccessName\": \"slbAccessForVip\"       }     },     {       \"name\": \"MockForOldClient\",       \"condition\": \"$ClientVersion < '2.0.5'\",       \"backend\": {         \"type\": \"MOCK\",         \"statusCode\": 400,         \"mockBody\": \"This version is not supported!!!\"       }     },     {       \"name\": \"BlueGreenPercent05\",       \"condition\": \"1 = 1\",       \"backend\": {         \"type\": \"HTTP\",         \"address\": \"https://beta-version.api.foo.com\"       },       \"constant-parameters\": [         {           \"name\": \"x-route-blue-green\",           \"location\": \"header\",           \"value\": \"route-blue-green\"         }       ]     }   ] }",
						"plugin_type": "routing",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case BasicAuth插件 6789  raw
func TestAccAliCloudApiGatewayPlugin_basic6789_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6789)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6789)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"users\\\": [     {       \\\"username\\\": \\\"alice\\\",       \\\"password\\\": 123456     },     {       \\\"username\\\": \\\"bob\\\",       \\\"password\\\": 666666     },     {       \\\"username\\\": \\\"charlie\\\",       \\\"password\\\": 888888     },     {       \\\"username\\\": \\\"dave\\\",       \\\"password\\\": 111111     }   ] }",
					"plugin_type": "basicAuth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name,
						"plugin_data": "{   \"users\": [     {       \"username\": \"alice\",       \"password\": 123456     },     {       \"username\": \"bob\",       \"password\": 666666     },     {       \"username\": \"charlie\",       \"password\": 888888     },     {       \"username\": \"dave\",       \"password\": 111111     }   ] }",
						"plugin_type": "basicAuth",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 参数访问控制插件 6793  raw
func TestAccAliCloudApiGatewayPlugin_basic6793_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6793)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6793)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"parameters\\\": {     \\\"userId\\\": \\\"Token:userId\\\",     \\\"userType\\\": \\\"Token:userType\\\",     \\\"pathUserId\\\": \\\"path:userId\\\"   },   \\\"rules\\\": [     {       \\\"name\\\": \\\"admin\\\",       \\\"condition\\\": \\\"$userType = 'admin'\\\",       \\\"ifTrue\\\": \\\"ALLOW\\\"     },     {       \\\"name\\\": \\\"user\\\",       \\\"condition\\\": \\\"$userId = $pathUserId\\\",       \\\"ifFalse\\\": \\\"DENY\\\",       \\\"statusCode\\\": 403,       \\\"errorMessage\\\": \\\"Path not match $${userId} vs /$${pathUserId}\\\",       \\\"responseHeaders\\\": {         \\\"Content-Type\\\": \\\"application/xml\\\"       },       \\\"responseBody\\\": \\\"<Reason>Path not match $${userId} vs /$${pathUserId}</Reason>\\\\n\\\"     }   ] }",
					"plugin_type": "accessControl",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name,
						"plugin_data": "{   \"parameters\": {     \"userId\": \"Token:userId\",     \"userType\": \"Token:userType\",     \"pathUserId\": \"path:userId\"   },   \"rules\": [     {       \"name\": \"admin\",       \"condition\": \"$userType = 'admin'\",       \"ifTrue\": \"ALLOW\"     },     {       \"name\": \"user\",       \"condition\": \"$userId = $pathUserId\",       \"ifFalse\": \"DENY\",       \"statusCode\": 403,       \"errorMessage\": \"Path not match ${userId} vs /${pathUserId}\",       \"responseHeaders\": {         \"Content-Type\": \"application/xml\"       },       \"responseBody\": \"<Reason>Path not match ${userId} vs /${pathUserId}</Reason>\\n\"     }   ] }",
						"plugin_type": "accessControl",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 第三方鉴权插件 6796  raw
func TestAccAliCloudApiGatewayPlugin_basic6796_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6796)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6796)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"parameters\\\": {     \\\"statusCode\\\": \\\"StatusCode\\\"   },   \\\"authUriType\\\": \\\"HTTP\\\",   \\\"authUri\\\": {     \\\"address\\\": \\\"http://your-auth-domain.com:8080\\\",     \\\"path\\\": \\\"/your/authPath\\\",     \\\"timeout\\\": 7000,     \\\"method\\\": \\\"POST\\\"   },   \\\"passThroughBody\\\": false,   \\\"cachedTimeBySecond\\\": 10,   \\\"authParameters\\\": [     {       \\\"targetParameterName\\\": \\\"x-userId\\\",       \\\"sourceParameterName\\\": \\\"userId\\\",       \\\"targetLocation\\\": \\\"form\\\",       \\\"sourceLocation\\\": \\\"query\\\"     },     {       \\\"targetParameterName\\\": \\\"x-passwoed\\\",       \\\"sourceParameterName\\\": \\\"password\\\",       \\\"targetLocation\\\": \\\"form\\\",       \\\"sourceLocation\\\": \\\"query\\\"     }   ],   \\\"successCondition\\\": \\\"$${statusCode} = 200\\\",   \\\"errorMessage\\\": \\\"auth failed\\\",   \\\"errorStatusCode\\\": 401 }",
					"plugin_type": "remoteAuth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name,
						"plugin_data": "{   \"parameters\": {     \"statusCode\": \"StatusCode\"   },   \"authUriType\": \"HTTP\",   \"authUri\": {     \"address\": \"http://your-auth-domain.com:8080\",     \"path\": \"/your/authPath\",     \"timeout\": 7000,     \"method\": \"POST\"   },   \"passThroughBody\": false,   \"cachedTimeBySecond\": 10,   \"authParameters\": [     {       \"targetParameterName\": \"x-userId\",       \"sourceParameterName\": \"userId\",       \"targetLocation\": \"form\",       \"sourceLocation\": \"query\"     },     {       \"targetParameterName\": \"x-passwoed\",       \"sourceParameterName\": \"password\",       \"targetLocation\": \"form\",       \"sourceLocation\": \"query\"     }   ],   \"successCondition\": \"${statusCode} = 200\",   \"errorMessage\": \"auth failed\",   \"errorStatusCode\": 401 }",
						"plugin_type": "remoteAuth",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 缓存插件 6791  raw
func TestAccAliCloudApiGatewayPlugin_basic6791_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6791)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6791)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"varyByApp\\\": false,   \\\"varyByParameters\\\": [     \\\"userId\\\"   ],   \\\"varyByHeaders\\\": [     \\\"Accept\\\"   ],   \\\"clientCacheControl\\\": {     \\\"mode\\\": \\\"app\\\",     \\\"apps\\\": [       1992323,       1239922     ]   },   \\\"cacheableHeaders\\\": [     \\\"X-Customer-Token\\\"   ] }",
					"plugin_type": "caching",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name,
						"plugin_data": "{   \"varyByApp\": false,   \"varyByParameters\": [     \"userId\"   ],   \"varyByHeaders\": [     \"Accept\"   ],   \"clientCacheControl\": {     \"mode\": \"app\",     \"apps\": [       1992323,       1239922     ]   },   \"cacheableHeaders\": [     \"X-Customer-Token\"   ] }",
						"plugin_type": "caching",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 后端签名插件 6799  raw
func TestAccAliCloudApiGatewayPlugin_basic6799_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6799)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6799)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"type\\\": \\\"APIGW_BACKEND\\\",   \\\"key\\\": \\\"SampleKey\\\",   \\\"secret\\\": \\\"SampleSecret\\\" }",
					"plugin_type": "backendSignature",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name,
						"plugin_data": "{   \"type\": \"APIGW_BACKEND\",   \"key\": \"SampleKey\",   \"secret\": \"SampleSecret\" }",
						"plugin_type": "backendSignature",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 错误码映射插件 6794  raw
func TestAccAliCloudApiGatewayPlugin_basic6794_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6794)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6794)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"parameters\\\": {     \\\"statusCode\\\": \\\"StatusCode\\\",     \\\"resultCode\\\": \\\"BodyJsonField:$.result_code\\\",     \\\"requestId\\\": \\\"BodyJsonField:$.req_msg_id\\\"   },   \\\"errorCondition\\\": \\\"$statusCode = 200 and $resultCode != null and $resultCode != 'OK'\\\",   \\\"errorCode\\\": \\\"resultCode\\\",   \\\"mappings\\\": [     {       \\\"code\\\": \\\"ROLE_NOT_EXISTS\\\",       \\\"statusCode\\\": 404,       \\\"errorMessage\\\": \\\"Role Not Exists, RequestId=$${requestId}\\\"     },     {       \\\"code\\\": \\\"INVALID_PARAMETER\\\",       \\\"statusCode\\\": 400,       \\\"errorMessage\\\": \\\"Invalid Parameter, RequestId=$${requestId}\\\"     }   ],   \\\"defaultMapping\\\": {     \\\"statusCode\\\": 500,     \\\"errorMessage\\\": \\\"Unknown Error, $${resultCode}, RequestId=$${requestId}\\\"   } }",
					"plugin_type": "errorMapping",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name,
						"plugin_data": "{   \"parameters\": {     \"statusCode\": \"StatusCode\",     \"resultCode\": \"BodyJsonField:$.result_code\",     \"requestId\": \"BodyJsonField:$.req_msg_id\"   },   \"errorCondition\": \"$statusCode = 200 and $resultCode != null and $resultCode != 'OK'\",   \"errorCode\": \"resultCode\",   \"mappings\": [     {       \"code\": \"ROLE_NOT_EXISTS\",       \"statusCode\": 404,       \"errorMessage\": \"Role Not Exists, RequestId=${requestId}\"     },     {       \"code\": \"INVALID_PARAMETER\",       \"statusCode\": 400,       \"errorMessage\": \"Invalid Parameter, RequestId=${requestId}\"     }   ],   \"defaultMapping\": {     \"statusCode\": 500,     \"errorMessage\": \"Unknown Error, ${resultCode}, RequestId=${requestId}\"   } }",
						"plugin_type": "errorMapping",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 请求应答改写插件 6798  raw
func TestAccAliCloudApiGatewayPlugin_basic6798_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6798)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6798)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"addRequestParameterIfAbsent\\\": [     {       \\\"name\\\": \\\"userId\\\",       \\\"value\\\": 123456,       \\\"location\\\": \\\"query\\\"     }   ],   \\\"putRequestParameter\\\": [     {       \\\"name\\\": \\\"name\\\",       \\\"value\\\": null,       \\\"location\\\": \\\"header\\\"     }   ],   \\\"removeRequestParameter\\\": [     {       \\\"name\\\": \\\"address\\\",       \\\"location\\\": \\\"form\\\"     }   ],   \\\"setResponseStatusCode\\\": 200,   \\\"addResponseHeaderIfAbsent\\\": [     {       \\\"name\\\": \\\"age\\\",       \\\"value\\\": 18     }   ],   \\\"putResponseHeader\\\": [     {       \\\"name\\\": \\\"name\\\",       \\\"value\\\": \\\"Alice\\\"     }   ],   \\\"removeResponseHeader\\\": [     {       \\\"name\\\": \\\"phone\\\"     }   ] }",
					"plugin_type": "transformer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name,
						"plugin_data": "{   \"addRequestParameterIfAbsent\": [     {       \"name\": \"userId\",       \"value\": 123456,       \"location\": \"query\"     }   ],   \"putRequestParameter\": [     {       \"name\": \"name\",       \"value\": null,       \"location\": \"header\"     }   ],   \"removeRequestParameter\": [     {       \"name\": \"address\",       \"location\": \"form\"     }   ],   \"setResponseStatusCode\": 200,   \"addResponseHeaderIfAbsent\": [     {       \"name\": \"age\",       \"value\": 18     }   ],   \"putResponseHeader\": [     {       \"name\": \"name\",       \"value\": \"Alice\"     }   ],   \"removeResponseHeader\": [     {       \"name\": \"phone\"     }   ] }",
						"plugin_type": "transformer",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 断路器插件 6795  raw
func TestAccAliCloudApiGatewayPlugin_basic6795_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6795)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6795)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"timeoutThreshold\\\": 15,   \\\"windowInSeconds\\\": 30,   \\\"openTimeoutSeconds\\\": 15,   \\\"downgradeBackend\\\": {     \\\"type\\\": \\\"mock\\\",     \\\"statusCode\\\": 302,     \\\"body\\\": \\\"<result>\\\\n  <errorCode>I's a teapot</errorCode>\\\\n</result>\\\\n\\\"   } }",
					"plugin_type": "circuitBreaker",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name,
						"plugin_data": "{   \"timeoutThreshold\": 15,   \"windowInSeconds\": 30,   \"openTimeoutSeconds\": 15,   \"downgradeBackend\": {     \"type\": \"mock\",     \"statusCode\": 302,     \"body\": \"<result>\\n  <errorCode>I's a teapot</errorCode>\\n</result>\\n\"   } }",
						"plugin_type": "circuitBreaker",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case ip控制插件 5922  raw
func TestAccAliCloudApiGatewayPlugin_basic5922_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap5922)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence5922)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{\\\"type\\\":\\\"ALLOW\\\",\\\"items\\\":[{\\\"blocks\\\":[\\\"79.11.12.2\\\"]}]}",
					"plugin_type": "ipControl",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name,
						"plugin_data": "{\"type\":\"ALLOW\",\"items\":[{\"blocks\":[\"79.11.12.2\"]}]}",
						"plugin_type": "ipControl",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tetetete",
					"plugin_name": name + "_update",
					"plugin_data": "{\\\"type\\\":\\\"ALLOW\\\",\\\"items\\\":[{\\\"blocks\\\":[\\\"79.11.12.2\\\",\\\"127.0.0.1\\\"]}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tetetete",
						"plugin_name": name + "_update",
						"plugin_data": "{\"type\":\"ALLOW\",\"items\":[{\"blocks\":[\"79.11.12.2\",\"127.0.0.1\"]}]}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 日志脱敏插件 6797  raw
func TestAccAliCloudApiGatewayPlugin_basic6797_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6797)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6797)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"rules\\\": [     {       \\\"name\\\": \\\"request_query\\\",       \\\"location\\\": \\\"REQUEST_QUERY\\\",       \\\"parameters\\\": [         \\\"userid\\\",         \\\"name\\\"       ],       \\\"policy\\\": \\\"KEEP_LEFT:4\\\"     },     {       \\\"name\\\": \\\"request_header\\\",       \\\"location\\\": \\\"REQUEST_HEADER\\\",       \\\"parameters\\\": [         \\\"userid\\\",         \\\"name\\\"       ],       \\\"policy\\\": \\\"KEEP_CENTER:4,5\\\"     },     {       \\\"name\\\": \\\"request_body_HEX\\\",       \\\"location\\\": \\\"REQUEST_BODY\\\",       \\\"matchMode\\\": \\\"HEX:10\\\",       \\\"policy\\\": \\\"ALL\\\"     },     {       \\\"name\\\": \\\"response_body_EMAIL\\\",       \\\"location\\\": \\\"RESPONSE_BODY\\\",       \\\"matchMode\\\": \\\"EMAIL\\\",       \\\"policy\\\": \\\"KEEP_RIGHT:7\\\"     }   ] }",
					"plugin_type": "logMask",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name,
						"plugin_data": "{   \"rules\": [     {       \"name\": \"request_query\",       \"location\": \"REQUEST_QUERY\",       \"parameters\": [         \"userid\",         \"name\"       ],       \"policy\": \"KEEP_LEFT:4\"     },     {       \"name\": \"request_header\",       \"location\": \"REQUEST_HEADER\",       \"parameters\": [         \"userid\",         \"name\"       ],       \"policy\": \"KEEP_CENTER:4,5\"     },     {       \"name\": \"request_body_HEX\",       \"location\": \"REQUEST_BODY\",       \"matchMode\": \"HEX:10\",       \"policy\": \"ALL\"     },     {       \"name\": \"response_body_EMAIL\",       \"location\": \"RESPONSE_BODY\",       \"matchMode\": \"EMAIL\",       \"policy\": \"KEEP_RIGHT:7\"     }   ] }",
						"plugin_type": "logMask",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case jwt插件 6788  raw
func TestAccAliCloudApiGatewayPlugin_basic6788_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6788)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6788)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"parameter\\\": \\\"X-Token\\\",   \\\"parameterLocation\\\": \\\"header\\\",   \\\"claimParameters\\\": [     {       \\\"claimName\\\": \\\"aud\\\",       \\\"parameterName\\\": \\\"X-Aud\\\",       \\\"location\\\": \\\"header\\\"     },     {       \\\"claimName\\\": \\\"userId\\\",       \\\"parameterName\\\": \\\"userId\\\",       \\\"location\\\": \\\"query\\\"     }   ],   \\\"preventJtiReplay\\\": false,   \\\"jwk\\\": {     \\\"kty\\\": \\\"RSA\\\",     \\\"e\\\": \\\"AQAB\\\",     \\\"use\\\": \\\"sig\\\",     \\\"kid\\\": \\\"O8fpdhrViq2zaaaBEWZITz\\\",     \\\"alg\\\": \\\"RS256\\\",     \\\"n\\\": \\\"qSVxcknOm0uCq5vGsOmaorPDzHUubBmZZ4UXj-9do7w9X1uKFXAnqfto4TepSNuYU2bA_-tzSLAGBsR-BqvT6w9SjxakeiyQpVmexxnDw5WZwpWenUAcYrfSPEoNU-0hAQwFYgqZwJQMN8ptxkd0170PFauwACOx4Hfr-9FPGy8NCoIO4MfLXzJ3mJ7xqgIZp3NIOGXz-GIAbCf13ii7kSStpYqN3L_zzpvXUAos1FJ9IPXRV84tIZpFVh2lmRh0h8ImK-vI42dwlD_hOIzayL1Xno2R0T-d5AwTSdnep7g-Fwu8-sj4cCRWq3bd61Zs2QOJ8iustH0vSRMYdP5oYQ\\\"   } }",
					"plugin_type": "jwtAuth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name,
						"plugin_data": "{   \"parameter\": \"X-Token\",   \"parameterLocation\": \"header\",   \"claimParameters\": [     {       \"claimName\": \"aud\",       \"parameterName\": \"X-Aud\",       \"location\": \"header\"     },     {       \"claimName\": \"userId\",       \"parameterName\": \"userId\",       \"location\": \"query\"     }   ],   \"preventJtiReplay\": false,   \"jwk\": {     \"kty\": \"RSA\",     \"e\": \"AQAB\",     \"use\": \"sig\",     \"kid\": \"O8fpdhrViq2zaaaBEWZITz\",     \"alg\": \"RS256\",     \"n\": \"qSVxcknOm0uCq5vGsOmaorPDzHUubBmZZ4UXj-9do7w9X1uKFXAnqfto4TepSNuYU2bA_-tzSLAGBsR-BqvT6w9SjxakeiyQpVmexxnDw5WZwpWenUAcYrfSPEoNU-0hAQwFYgqZwJQMN8ptxkd0170PFauwACOx4Hfr-9FPGy8NCoIO4MfLXzJ3mJ7xqgIZp3NIOGXz-GIAbCf13ii7kSStpYqN3L_zzpvXUAos1FJ9IPXRV84tIZpFVh2lmRh0h8ImK-vI42dwlD_hOIzayL1Xno2R0T-d5AwTSdnep7g-Fwu8-sj4cCRWq3bd61Zs2QOJ8iustH0vSRMYdP5oYQ\"   } }",
						"plugin_type": "jwtAuth",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 流控插件 6787  raw
func TestAccAliCloudApiGatewayPlugin_basic6787_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6787)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6787)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"unit\\\": \\\"SECOND\\\",   \\\"apiDefault\\\": 1000,   \\\"userDefault\\\": 30,   \\\"appDefault\\\": 30,   \\\"specials\\\": [     {       \\\"type\\\": \\\"APP\\\",       \\\"policies\\\": [         {           \\\"key\\\": 10123123,           \\\"value\\\": 10         },         {           \\\"key\\\": 10123123,           \\\"value\\\": 10         }       ]     },     {       \\\"type\\\": \\\"USER\\\",       \\\"policies\\\": [         {           \\\"key\\\": 123455,           \\\"value\\\": 100         }       ]     }   ] }",
					"plugin_type": "trafficControl",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name,
						"plugin_data": "{   \"unit\": \"SECOND\",   \"apiDefault\": 1000,   \"userDefault\": 30,   \"appDefault\": 30,   \"specials\": [     {       \"type\": \"APP\",       \"policies\": [         {           \"key\": 10123123,           \"value\": 10         },         {           \"key\": 10123123,           \"value\": 10         }       ]     },     {       \"type\": \"USER\",       \"policies\": [         {           \"key\": 123455,           \"value\": 100         }       ]     }   ] }",
						"plugin_type": "trafficControl",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Case 跨域插件 6790  raw
func TestAccAliCloudApiGatewayPlugin_basic6790_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayPluginMap6790)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewayplugin%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayPluginBasicDependence6790)
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
					"description": "test",
					"plugin_name": name,
					"plugin_data": "{   \\\"allowOrigins\\\": \\\"api.foo.com\\\",   \\\"allowMethods\\\": \\\"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\\\",   \\\"allowHeaders\\\": \\\"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\\\",   \\\"exposeHeaders\\\": \\\"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\\\",   \\\"maxAge\\\": 172800,   \\\"allowCredentials\\\": true }",
					"plugin_type": "cors",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
						"plugin_name": name,
						"plugin_data": "{   \"allowOrigins\": \"api.foo.com\",   \"allowMethods\": \"GET,POST,PUT,DELETE,HEAD,OPTIONS,PATCH\",   \"allowHeaders\": \"Authorization,Accept,Accept-Ranges,Cache-Control,Range,Date,Content-Type,Content-Length,Content-MD5,User-Agent,X-Ca-Signature,X-Ca-Signature-Headers,X-Ca-Signature-Method,X-Ca-Key,X-Ca-Timestamp,X-Ca-Nonce,X-Ca-Stage,X-Ca-Request-Mode,x-ca-deviceid\",   \"exposeHeaders\": \"Content-MD5,Server,Date,Latency,X-Ca-Request-Id,X-Ca-Error-Code,X-Ca-Error-Message\",   \"maxAge\": 172800,   \"allowCredentials\": true }",
						"plugin_type": "cors",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_time", "modified_time"},
			},
		},
	})
}

// Test ApiGateway Plugin. <<< Resource test cases, automatically generated.
