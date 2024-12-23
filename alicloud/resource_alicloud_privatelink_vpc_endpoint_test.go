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
	resource.AddTestSweepers(
		"alicloud_privatelink_vpc_endpoint",
		&resource.Sweeper{
			Name: "alicloud_privatelink_vpc_endpoint",
			F:    testSweepPrivatelinkVpcEndpoint,
		})
}

func testSweepPrivatelinkVpcEndpoint(region string) error {
	if !testSweepPreCheckWithRegions(region, false, connectivity.PrivateLinkRegions) {
		log.Printf("[INFO] Skipping privatelink unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}
	request := map[string]interface{}{
		"MaxResults": PageSizeLarge,
	}
	var response map[string]interface{}
	action := "ListVpcEndpoints"
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = client.RpcPost("Privatelink", "2020-04-15", action, nil, request, true)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_privatelink_vpc_endpoint", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Endpoints", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Endpoints", response)
		}
		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			if item["EndpointName"] == nil {
				continue
			}
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["EndpointName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Privatelink VpcEndpoint: %s", item["EndpointName"].(string))
				continue
			}
			sweeped = true
			action = "DeleteVpcEndpoint"
			request := map[string]interface{}{
				"EndpointId": item["EndpointId"],
			}
			_, err = client.RpcPost("Privatelink", "2020-04-15", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Privatelink VpcEndpoint (%s): %s", item["EndpointName"].(string), err)
			}
			if sweeped {
				// Waiting 5 seconds to ensure Privatelink VpcEndpoint have been deleted.
				time.Sleep(5 * time.Second)
			}
			log.Printf("[INFO] Delete Privatelink VpcEndpoint success: %s ", item["EndpointName"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAliCloudPrivatelinkVpcEndpoint_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivatelinkService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivatelinkVpcEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivatelinkVpcEndpointBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.PrivateLinkRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_id":         "${alicloud_privatelink_vpc_endpoint_service.default.id}",
					"vpc_id":             "${data.alicloud_vpcs.default.ids.0}",
					"security_group_ids": []string{"${alicloud_security_group.default.id}"},
					"vpc_endpoint_name":  name,
					"depends_on":         []string{"alicloud_privatelink_vpc_endpoint_service.default"},
					"address_ip_version": "IPv4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_id":           CHECKSET,
						"vpc_id":               CHECKSET,
						"security_group_ids.#": "1",
						"vpc_endpoint_name":    name,
						"address_ip_version":   "IPv4",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_description": "Terraform Test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_description": "Terraform Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_endpoint_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_endpoint_name": name + "update",
					}),
				),
			},
			// TODO：There is a bug with the API here, which means calling the API will not result in an error,
			// but the modified value will not become what is expected.
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"address_ip_version": "DualStack",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"address_ip_version": "DualStack",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{"${alicloud_security_group.default.id}", "${alicloud_security_group.default2.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{"${alicloud_security_group.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids":   []string{"${alicloud_security_group.default.id}", "${alicloud_security_group.default2.id}"},
					"endpoint_description": "Terraform Test Update",
					"vpc_endpoint_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "2",
						"endpoint_description": "Terraform Test Update",
						"vpc_endpoint_name":    name,
					}),
				),
			},
		},
	})
}

var AlicloudPrivatelinkVpcEndpointMap = map[string]string{
	"bandwidth":                CHECKSET,
	"connection_status":        CHECKSET,
	"endpoint_business_status": CHECKSET,
	"endpoint_domain":          CHECKSET,
	"service_name":             CHECKSET,
	"status":                   CHECKSET,
}

func AlicloudPrivatelinkVpcEndpointBasicDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_vpcs" "default" {
	  name_regex = "^default-NODELETING$"
	}
	resource "alicloud_security_group" "default" {
	  name        = "tf-testAcc-for-privatelink"
	  description = "privatelink test security group"
	  vpc_id      = data.alicloud_vpcs.default.ids.0
	}
	resource "alicloud_security_group" "default2" {
	  name        = "%[1]s"
	  description = "privatelink test security group2"
	  vpc_id      = data.alicloud_vpcs.default.ids.0
	}
	resource "alicloud_privatelink_vpc_endpoint_service" "default" {
	 service_description = "%[1]s"
	 connect_bandwidth = 103
     auto_accept_connection = false
	}
`, name)
}

func TestUnitAlicloudPrivatelinkVpcEndpoint(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"dry_run":              false,
		"service_id":           "CreateVpcEndpointValue",
		"vpc_id":               "CreateVpcEndpointValue",
		"security_group_ids":   []string{"CreateVpcEndpointValue"},
		"vpc_endpoint_name":    "CreateVpcEndpointValue",
		"endpoint_description": "CreateVpcEndpointValue",
		"service_name":         "CreateVpcEndpointValue",
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
		// GetVpcEndpointAttribute
		"EndpointId":             "CreateVpcEndpointValue",
		"Bandwidth":              100,
		"ConnectionStatus":       "CreateVpcEndpointValue",
		"EndpointBusinessStatus": "CreateVpcEndpointValue",
		"EndpointDescription":    "CreateVpcEndpointValue",
		"EndpointDomain":         "CreateVpcEndpointValue",
		"ServiceId":              "CreateVpcEndpointValue",
		"ServiceName":            "CreateVpcEndpointValue",
		"EndpointStatus":         "Active",
		"EndpointName":           "CreateVpcEndpointValue",
		"VpcId":                  "CreateVpcEndpointValue",
		"SecurityGroups": []interface{}{
			map[string]interface{}{
				"SecurityGroupId": "CreateVpcEndpointValue",
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateVpcEndpoint
		"EndpointId": "CreateVpcEndpointValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_privatelink_vpc_endpoint", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewPrivatelinkClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudPrivateLinkVpcEndpointCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// GetVpcEndpointAttribute Response
		"EndpointId": "CreateVpcEndpointValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateVpcEndpoint" {
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
		err := resourceAliCloudPrivateLinkVpcEndpointCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewPrivatelinkClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudPrivateLinkVpcEndpointUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateVpcEndpointAttribute
	attributesDiff := map[string]interface{}{
		"endpoint_description": "UpdateVpcEndpointAttributeValue",
		"vpc_endpoint_name":    "UpdateVpcEndpointAttributeValue",
		"dry_run":              true,
	}
	diff, err := newInstanceDiff("alicloud_privatelink_vpc_endpoint", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_privatelink_vpc_endpoint"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetVpcEndpointAttribute Response
		"EndpointDescription": "UpdateVpcEndpointAttributeValue",
		"EndpointName":        "UpdateVpcEndpointAttributeValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateVpcEndpointAttribute" {
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
		err := resourceAliCloudPrivateLinkVpcEndpointUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// AttachSecurityGroupToVpcEndpoint
	attributesDiff = map[string]interface{}{
		"security_group_ids": []string{"AttachSecurityGroupToVpcEndpointValue"},
		"dry_run":            true,
	}
	diff, err = newInstanceDiff("alicloud_privatelink_vpc_endpoint", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_privatelink_vpc_endpoint"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetVpcEndpointAttribute Response
		"SecurityGroups": []interface{}{
			map[string]interface{}{
				"SecurityGroupId": "AttachSecurityGroupToVpcEndpointValue",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "AttachSecurityGroupToVpcEndpoint" {
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
		err := resourceAliCloudPrivateLinkVpcEndpointUpdate(dExisted, rawClient)
		patches.Reset()

		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// DetachSecurityGroupFromVpcEndpoint
	attributesDiff = map[string]interface{}{
		"security_group_ids": []string{"DetachSecurityGroupFromVpcEndpointValue"},
		"dry_run":            true,
	}
	diff, err = newInstanceDiff("alicloud_privatelink_vpc_endpoint", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_privatelink_vpc_endpoint"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetVpcEndpointAttribute Response
		"SecurityGroups": []interface{}{
			map[string]interface{}{
				"SecurityGroupId": "DetachSecurityGroupFromVpcEndpointValue",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DetachSecurityGroupFromVpcEndpoint" {
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
		err := resourceAliCloudPrivateLinkVpcEndpointUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetVpcEndpointAttribute" {
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
		err := resourceAliCloudPrivateLinkVpcEndpointRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewPrivatelinkClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudPrivateLinkVpcEndpointDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "EndpointOperationDenied", "nil", "EndpointNotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteVpcEndpoint" {
				switch errorCode {
				case "NonRetryableError", "EndpointNotFound":
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
		err := resourceAliCloudPrivateLinkVpcEndpointDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "EndpointNotFound":
			assert.Nil(t, err)
		}
	}

}

func TestAccAliCloudPrivateLinkVpcEndpoint_basic4793(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointMap4793)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivateLinkServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivateLinkVpcEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpoint%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivateLinkVpcEndpointBasicDependence4793)
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
					"service_id":        "${alicloud_privatelink_vpc_endpoint_service.defaultr0WBYX.id}",
					"vpc_id":            "${alicloud_vpc.defaultbFzA4a.id}",
					"vpc_endpoint_name": name,
					"security_group_ids": []string{
						"${alicloud_security_group.default1FTFrP.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_id":        CHECKSET,
						"vpc_id":            CHECKSET,
						"vpc_endpoint_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_description": "test-endpoint-zejun",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_description": "test-endpoint-zejun",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{
						"${alicloud_security_group.default1FTFrP.id}", "${alicloud_security_group.defaultHtejEL.id}", "${alicloud_security_group.default97JOJ3.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_endpoint_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_endpoint_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_endpoint_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_endpoint_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_description": "test-endpoint-zejun1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_description": "test-endpoint-zejun1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{
						"${alicloud_security_group.defaultHtejEL.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_endpoint_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_endpoint_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_description": "test-endpoint-zejun",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_description": "test-endpoint-zejun",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{
						"${alicloud_security_group.default1FTFrP.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_endpoint_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_endpoint_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_description": "test-endpoint-zejun1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_description": "test-endpoint-zejun1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{
						"${alicloud_security_group.defaultHtejEL.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_endpoint_name":             name + "_update",
					"vpc_id":                        "${alicloud_vpc.defaultbFzA4a.id}",
					"endpoint_description":          "test-endpoint-zejun",
					"service_id":                    "${alicloud_privatelink_vpc_endpoint_service.defaultr0WBYX.id}",
					"service_name":                  "${alicloud_privatelink_vpc_endpoint_service.defaultr0WBYX.vpc_endpoint_service_name}",
					"resource_group_id":             "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"endpoint_type":                 "Interface",
					"zone_private_ip_address_count": "1",
					"security_group_ids": []string{
						"${alicloud_security_group.default1FTFrP.id}", "${alicloud_security_group.defaultHtejEL.id}", "${alicloud_security_group.default97JOJ3.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_endpoint_name":             name + "_update",
						"vpc_id":                        CHECKSET,
						"endpoint_description":          "test-endpoint-zejun",
						"service_id":                    CHECKSET,
						"service_name":                  CHECKSET,
						"resource_group_id":             CHECKSET,
						"endpoint_type":                 "Interface",
						"zone_private_ip_address_count": "1",
						"security_group_ids.#":          "3",
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
				ImportStateVerifyIgnore: []string{"dry_run", "protected_enabled"},
			},
		},
	})
}

var AlicloudPrivateLinkVpcEndpointMap4793 = map[string]string{
	"endpoint_domain":          CHECKSET,
	"bandwidth":                CHECKSET,
	"connection_status":        CHECKSET,
	"status":                   CHECKSET,
	"create_time":              CHECKSET,
	"endpoint_business_status": CHECKSET,
}

func AlicloudPrivateLinkVpcEndpointBasicDependence4793(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultbFzA4a" {
  description = "test-terraform"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name

}

resource "alicloud_security_group" "default1FTFrP" {
  name = var.name

  vpc_id = alicloud_vpc.defaultbFzA4a.id
}

resource "alicloud_security_group" "defaultHtejEL" {
  name = var.name

  vpc_id = alicloud_vpc.defaultbFzA4a.id
}

resource "alicloud_security_group" "default97JOJ3" {
  name = var.name

  vpc_id = alicloud_vpc.defaultbFzA4a.id
}

resource "alicloud_privatelink_vpc_endpoint_service" "defaultr0WBYX" {
  service_description   = "test-zejun-service"
  connect_bandwidth     = "3072"
  service_resource_type = "nlb"
}


`, name)
}

// Case 4793  twin
func TestAccAliCloudPrivateLinkVpcEndpoint_basic4793_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointMap4793)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivateLinkServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivateLinkVpcEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpoint%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivateLinkVpcEndpointBasicDependence4793)
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
					"vpc_endpoint_name":             name,
					"vpc_id":                        "${alicloud_vpc.defaultbFzA4a.id}",
					"endpoint_description":          "test-endpoint-zejun1",
					"service_id":                    "${alicloud_privatelink_vpc_endpoint_service.defaultr0WBYX.id}",
					"service_name":                  "${alicloud_privatelink_vpc_endpoint_service.defaultr0WBYX.vpc_endpoint_service_name}",
					"resource_group_id":             "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"endpoint_type":                 "Interface",
					"zone_private_ip_address_count": "1",
					"security_group_ids": []string{
						"${alicloud_security_group.defaultHtejEL.id}", "${alicloud_security_group.default97JOJ3.id}"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_endpoint_name":             name,
						"vpc_id":                        CHECKSET,
						"endpoint_description":          "test-endpoint-zejun1",
						"service_id":                    CHECKSET,
						"service_name":                  CHECKSET,
						"resource_group_id":             CHECKSET,
						"endpoint_type":                 "Interface",
						"zone_private_ip_address_count": "1",
						"security_group_ids.#":          "2",
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.For":                      "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "protected_enabled"},
			},
		},
	})
}

// Test PrivateLink VpcEndpoint. >>> Resource test cases, automatically generated.
// Case 生命周期测试（PolicyDocument发布terraform） 6705
func TestAccAliCloudPrivateLinkVpcEndpoint_basic6705(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointMap6705)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivateLinkServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivateLinkVpcEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpoint%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivateLinkVpcEndpointBasicDependence6705)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"ap-southeast-5"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":            "${alicloud_vpc.defaultbFzA4a.id}",
					"vpc_endpoint_name": name,
					"service_id":        "epsrv-k1apjysze8u1l9t6uyg9",
					"service_name":      "com.aliyuncs.privatelink.ap-southeast-5.oss",
					"security_group_ids": []string{
						"${alicloud_security_group.default1FTFrP.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":            CHECKSET,
						"vpc_endpoint_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_description": "test-endpoint-zejun",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_description": "test-endpoint-zejun",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_document": "{   \\\"Version\\\": \\\"1\\\",   \\\"Statement\\\": [     {       \\\"Effect\\\": \\\"Allow\\\",       \\\"Action\\\": [         \\\"*\\\"       ],       \\\"Resource\\\": [         \\\"*\\\"       ],       \\\"Principal\\\": \\\"*\\\"     }   ] }",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_document": "{   \"Version\": \"1\",   \"Statement\": [     {       \"Effect\": \"Allow\",       \"Action\": [         \"*\"       ],       \"Resource\": [         \"*\"       ],       \"Principal\": \"*\"     }   ] }",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{
						"${alicloud_security_group.default1FTFrP.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_description": "test-endpoint-zejun1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_description": "test-endpoint-zejun1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_endpoint_name":             name + "_update",
					"vpc_id":                        "${alicloud_vpc.defaultbFzA4a.id}",
					"endpoint_description":          "test-endpoint-zejun",
					"service_id":                    "epsrv-k1apjysze8u1l9t6uyg9",
					"service_name":                  "com.aliyuncs.privatelink.ap-southeast-5.oss",
					"protected_enabled":             "false",
					"resource_group_id":             "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"dry_run":                       "false",
					"endpoint_type":                 "Interface",
					"zone_private_ip_address_count": "1",
					"policy_document":               "{   \\\"Version\\\": \\\"1\\\",   \\\"Statement\\\": [     {       \\\"Effect\\\": \\\"Allow\\\",       \\\"Action\\\": [         \\\"*\\\"       ],       \\\"Resource\\\": [         \\\"*\\\"       ],       \\\"Principal\\\": \\\"*\\\"     }   ] }",
					"security_group_ids": []string{
						"${alicloud_security_group.default1FTFrP.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_endpoint_name":             name + "_update",
						"vpc_id":                        CHECKSET,
						"endpoint_description":          "test-endpoint-zejun",
						"service_id":                    "epsrv-k1apjysze8u1l9t6uyg9",
						"service_name":                  "com.aliyuncs.privatelink.ap-southeast-5.oss",
						"protected_enabled":             "false",
						"resource_group_id":             CHECKSET,
						"dry_run":                       "false",
						"endpoint_type":                 "Interface",
						"zone_private_ip_address_count": "1",
						"policy_document":               "{   \"Version\": \"1\",   \"Statement\": [     {       \"Effect\": \"Allow\",       \"Action\": [         \"*\"       ],       \"Resource\": [         \"*\"       ],       \"Principal\": \"*\"     }   ] }",
						"security_group_ids.#":          "1",
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
				ImportStateVerifyIgnore: []string{"dry_run", "protected_enabled"},
			},
		},
	})
}

var AlicloudPrivateLinkVpcEndpointMap6705 = map[string]string{
	"endpoint_domain":               CHECKSET,
	"bandwidth":                     CHECKSET,
	"endpoint_type":                 CHECKSET,
	"connection_status":             CHECKSET,
	"status":                        CHECKSET,
	"create_time":                   CHECKSET,
	"zone_private_ip_address_count": CHECKSET,
	"endpoint_business_status":      CHECKSET,
	"service_name":                  CHECKSET,
}

func AlicloudPrivateLinkVpcEndpointBasicDependence6705(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultbFzA4a" {
  description = "test-terraform"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_security_group" "default1FTFrP" {
  name = var.name
  vpc_id              = alicloud_vpc.defaultbFzA4a.id
}

resource "alicloud_security_group" "defaultjljY5S" {
  name = var.name
  vpc_id              = alicloud_vpc.defaultbFzA4a.id
}


`, name)
}

// Case 生命周期测试（PolicyDocument发布terraform） 6705  twin
func TestAccAliCloudPrivateLinkVpcEndpoint_basic6705_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointMap6705)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivateLinkServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivateLinkVpcEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpoint%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivateLinkVpcEndpointBasicDependence6705)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"ap-southeast-5"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_endpoint_name":             name,
					"vpc_id":                        "${alicloud_vpc.defaultbFzA4a.id}",
					"endpoint_description":          "test-endpoint-zejun",
					"service_id":                    "epsrv-k1apjysze8u1l9t6uyg9",
					"service_name":                  "com.aliyuncs.privatelink.ap-southeast-5.oss",
					"protected_enabled":             "false",
					"resource_group_id":             "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"dry_run":                       "false",
					"endpoint_type":                 "Interface",
					"zone_private_ip_address_count": "1",
					"policy_document":               "{   \\\"Version\\\": \\\"1\\\",   \\\"Statement\\\": [     {       \\\"Effect\\\": \\\"Allow\\\",       \\\"Action\\\": [         \\\"*\\\"       ],       \\\"Resource\\\": [         \\\"*\\\"       ],       \\\"Principal\\\": \\\"*\\\"     }   ] }",
					"security_group_ids": []string{
						"${alicloud_security_group.default1FTFrP.id}"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_endpoint_name":             name,
						"vpc_id":                        CHECKSET,
						"endpoint_description":          "test-endpoint-zejun",
						"service_id":                    "epsrv-k1apjysze8u1l9t6uyg9",
						"service_name":                  "com.aliyuncs.privatelink.ap-southeast-5.oss",
						"protected_enabled":             "false",
						"resource_group_id":             CHECKSET,
						"dry_run":                       "false",
						"endpoint_type":                 "Interface",
						"zone_private_ip_address_count": "1",
						"policy_document":               "{   \"Version\": \"1\",   \"Statement\": [     {       \"Effect\": \"Allow\",       \"Action\": [         \"*\"       ],       \"Resource\": [         \"*\"       ],       \"Principal\": \"*\"     }   ] }",
						"security_group_ids.#":          "1",
						"tags.%":                        "2",
						"tags.Created":                  "TF",
						"tags.For":                      "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "protected_enabled"},
			},
		},
	})
}

// Case 生命周期测试（PolicyDocument发布terraform） 6705  raw
func TestAccAliCloudPrivateLinkVpcEndpoint_basic6705_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointMap6705)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivateLinkServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivateLinkVpcEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpoint%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivateLinkVpcEndpointBasicDependence6705)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"ap-southeast-5"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_endpoint_name":             name,
					"vpc_id":                        "${alicloud_vpc.defaultbFzA4a.id}",
					"endpoint_description":          "test-endpoint-zejun",
					"service_id":                    "epsrv-k1apjysze8u1l9t6uyg9",
					"service_name":                  "com.aliyuncs.privatelink.ap-southeast-5.oss",
					"protected_enabled":             "false",
					"resource_group_id":             "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"dry_run":                       "false",
					"endpoint_type":                 "Interface",
					"zone_private_ip_address_count": "1",
					"policy_document":               "{   \\\"Version\\\": \\\"1\\\",   \\\"Statement\\\": [     {       \\\"Effect\\\": \\\"Allow\\\",       \\\"Action\\\": [         \\\"*\\\"       ],       \\\"Resource\\\": [         \\\"*\\\"       ],       \\\"Principal\\\": \\\"*\\\"     }   ] }",
					"security_group_ids": []string{
						"${alicloud_security_group.default1FTFrP.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_endpoint_name":             name,
						"vpc_id":                        CHECKSET,
						"endpoint_description":          "test-endpoint-zejun",
						"service_id":                    "epsrv-k1apjysze8u1l9t6uyg9",
						"service_name":                  "com.aliyuncs.privatelink.ap-southeast-5.oss",
						"protected_enabled":             "false",
						"resource_group_id":             CHECKSET,
						"dry_run":                       "false",
						"endpoint_type":                 "Interface",
						"zone_private_ip_address_count": "1",
						"policy_document":               "{   \"Version\": \"1\",   \"Statement\": [     {       \"Effect\": \"Allow\",       \"Action\": [         \"*\"       ],       \"Resource\": [         \"*\"       ],       \"Principal\": \"*\"     }   ] }",
						"security_group_ids.#":          "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_description": "test-endpoint-zejun1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_description": "test-endpoint-zejun1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "protected_enabled"},
			},
		},
	})
}

// Test PrivateLink VpcEndpoint. <<< Resource test cases, automatically generated.
