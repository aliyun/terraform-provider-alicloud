package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudServiceMeshExtensionProvider_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_mesh_extension_provider.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudServiceMeshExtensionProviderMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServicemeshService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceMeshExtensionProvider")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudServiceMeshExtensionProviderBasicDependence)
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
					"service_mesh_id":         "${alicloud_service_mesh_service_mesh.default.id}",
					"extension_provider_name": "httpextauth-" + name,
					"type":                    "httpextauth",
					"config":                  `{\"headersToDownstreamOnDeny\":[\"content-type\",\"set-cookie\"],\"headersToUpstreamOnAllow\":[\"authorization\",\"cookie\",\"path\",\"x-auth-request-access-token\",\"x-forwarded-access-token\"],\"includeRequestHeadersInCheck\":[\"cookie\",\"x-forward-access-token\"],\"oidc\":{\"clientID\":\"qweqweqwewqeqwe\",\"clientSecret\":\"asdasdasdasdsadas\",\"cookieExpire\":\"1000\",\"cookieRefresh\":\"500\",\"cookieSecret\":\"scxzcxzcxzcxzcxz\",\"issuerURI\":\"qweqwewqeqweqweqwe\",\"redirectDomain\":\"www.baidu.com\",\"redirectProtocol\":\"http\",\"scopes\":[\"profile\"]},\"port\":4180,\"service\":\"asm-oauth2proxy-httpextauth-` + name + `.istio-system.svc.cluster.local\",\"timeout\":\"10s\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_mesh_id":         CHECKSET,
						"extension_provider_name": "httpextauth-" + name,
						"type":                    "httpextauth",
						"config":                  "{\"headersToDownstreamOnDeny\":[\"content-type\",\"set-cookie\"],\"headersToUpstreamOnAllow\":[\"authorization\",\"cookie\",\"path\",\"x-auth-request-access-token\",\"x-forwarded-access-token\"],\"includeRequestHeadersInCheck\":[\"cookie\",\"x-forward-access-token\"],\"oidc\":{\"clientID\":\"qweqweqwewqeqwe\",\"clientSecret\":\"asdasdasdasdsadas\",\"cookieExpire\":\"1000\",\"cookieRefresh\":\"500\",\"cookieSecret\":\"scxzcxzcxzcxzcxz\",\"issuerURI\":\"qweqwewqeqweqweqwe\",\"redirectDomain\":\"www.baidu.com\",\"redirectProtocol\":\"http\",\"scopes\":[\"profile\"]},\"port\":4180,\"service\":\"asm-oauth2proxy-httpextauth-" + name + ".istio-system.svc.cluster.local\",\"timeout\":\"10s\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": `{\"headersToDownstreamOnDeny\":[\"content-type\",\"set-cookie\"],\"headersToUpstreamOnAllow\":[\"authorization\",\"cookie\",\"path\",\"x-auth-request-access-token\",\"x-forwarded-access-token\"],\"includeRequestHeadersInCheck\":[\"cookie\",\"x-forward-access-token\"],\"oidc\":{\"clientID\":\"qweqweqwewqeqwe\",\"clientSecret\":\"asdasdasdasdsadas\",\"cookieExpire\":\"1000\",\"cookieRefresh\":\"500\",\"cookieSecret\":\"scxzcxzcxzcxzcxz\",\"issuerURI\":\"qweqwewqeqweqweqwe\",\"redirectDomain\":\"www.baidu.com\",\"redirectProtocol\":\"http\",\"scopes\":[\"profile\"]},\"port\":4180,\"service\":\"asm-oauth2proxy-httpextauth-` + name + `.istio-system.svc.cluster.local\",\"timeout\":\"20s\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config": "{\"headersToDownstreamOnDeny\":[\"content-type\",\"set-cookie\"],\"headersToUpstreamOnAllow\":[\"authorization\",\"cookie\",\"path\",\"x-auth-request-access-token\",\"x-forwarded-access-token\"],\"includeRequestHeadersInCheck\":[\"cookie\",\"x-forward-access-token\"],\"oidc\":{\"clientID\":\"qweqweqwewqeqwe\",\"clientSecret\":\"asdasdasdasdsadas\",\"cookieExpire\":\"1000\",\"cookieRefresh\":\"500\",\"cookieSecret\":\"scxzcxzcxzcxzcxz\",\"issuerURI\":\"qweqwewqeqweqweqwe\",\"redirectDomain\":\"www.baidu.com\",\"redirectProtocol\":\"http\",\"scopes\":[\"profile\"]},\"port\":4180,\"service\":\"asm-oauth2proxy-httpextauth-" + name + ".istio-system.svc.cluster.local\",\"timeout\":\"20s\"}",
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

var resourceAlicloudServiceMeshExtensionProviderMap = map[string]string{}

func resourceAlicloudServiceMeshExtensionProviderBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}

	data "alicloud_vpcs" "default" {
		name_regex = "default-NODELETING"
	}

	resource "alicloud_vpc" "default" {
		count    = length(data.alicloud_vpcs.default.ids) > 0 ? 0 : 1
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
	}

	resource "alicloud_vswitch" "default" {
  		count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  		vpc_id       = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
  		cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_service_mesh_service_mesh" "default" {
		network {
			vpc_id        = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
			vswitche_list = [length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : alicloud_vswitch.default[0].id]
		}
	}
`, name)
}

func TestUnitAlicloudServiceMeshExtensionProvider(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_service_mesh_extension_provider"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_service_mesh_extension_provider"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"service_mesh_id":         "CreateServiceMeshExtensionProviderValue",
		"extension_provider_name": "CreateServiceMeshExtensionProviderValue",
		"type":                    "CreateServiceMeshExtensionProviderValue",
		"config":                  "CreateServiceMeshExtensionProviderValue",
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
		// DescribeExtensionProvider
		"ExtensionProviderArray": []interface{}{
			map[string]interface{}{
				"Type":   "CreateServiceMeshExtensionProviderValue",
				"Name":   "CreateServiceMeshExtensionProviderValue",
				"Config": "CreateServiceMeshExtensionProviderValue",
			},
		},
	}
	CreateMockResponse := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_service_mesh_extension_provider", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewServicemeshClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudServiceMeshExtensionProviderCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateExtensionProvider" {
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
		err := resourceAlicloudServiceMeshExtensionProviderCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_service_mesh_extension_provider"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewServicemeshClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudServiceMeshExtensionProviderUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"service_mesh_id":         "CreateServiceMeshExtensionProviderValue",
		"extension_provider_name": "CreateServiceMeshExtensionProviderValue",
		"type":                    "CreateServiceMeshExtensionProviderValue",
		"config":                  "PutServiceMeshExtensionProviderValue",
	}
	diff, err := newInstanceDiff("alicloud_service_mesh_extension_provider", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_service_mesh_extension_provider"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeExtensionProvider Response
		"ExtensionProviderArray": []interface{}{
			map[string]interface{}{
				"Type":   "CreateServiceMeshExtensionProviderValue",
				"Name":   "CreateServiceMeshExtensionProviderValue",
				"Config": "PutServiceMeshExtensionProviderValue",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateExtensionProvider" {
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
		err := resourceAlicloudServiceMeshExtensionProviderUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_service_mesh_extension_provider"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_service_mesh_extension_provider", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_service_mesh_extension_provider"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeExtensionProvider" {
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
		err := resourceAlicloudServiceMeshExtensionProviderRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewServicemeshClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudServiceMeshExtensionProviderDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_service_mesh_extension_provider", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_service_mesh_extension_provider"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteExtensionProvider" {
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
		err := resourceAlicloudServiceMeshExtensionProviderDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
