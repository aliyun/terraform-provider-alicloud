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
		"alicloud_nlb_server_group",
		&resource.Sweeper{
			Name: "alicloud_nlb_server_group",
			F:    testSweepNlbServerGroup,
		})
}

func testSweepNlbServerGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListServerGroups"
	request := map[string]interface{}{}
	request["RegionId"] = aliyunClient.RegionId

	request["MaxResults"] = PageSizeLarge

	var response map[string]interface{}
	conn, err := aliyunClient.NewNlbClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &runtime)
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

		resp, err := jsonpath.Get("$.ServerGroups", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.ServerGroups", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["ServerGroupName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Nlb Server Group: %s", item["ServerGroupName"].(string))
				continue
			}
			action := "DeleteServerGroup"
			request := map[string]interface{}{
				"ServerGroupId": item["ServerGroupId"],
				"RegionId":      aliyunClient.RegionId,
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Nlb Server Group (%s): %s", item["ServerGroupName"].(string), err)
			}
			log.Printf("[INFO] Delete Nlb Server Group success: %s ", item["ServerGroupName"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudNLBServerGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := resourceAttrInit(resourceId, AlicloudNLBServerGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNlbServerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNLBServerGroupBasicDependence0)
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
					"resource_group_id":          "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"server_group_name":          "${var.name}",
					"server_group_type":          "Instance",
					"vpc_id":                     "${data.alicloud_vpcs.default.ids.0}",
					"scheduler":                  "Wrr",
					"preserve_client_ip_enabled": "true",
					"protocol":                   "TCP",
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":         "true",
							"health_check_type":            "HTTP",
							"health_check_connect_port":    "0",
							"healthy_threshold":            "2",
							"unhealthy_threshold":          "2",
							"health_check_connect_timeout": "5",
							"health_check_interval":        "10",
							"http_check_method":            "GET",
							"health_check_url":             "/test/index.html",
							"health_check_domain":          "tf-testAcc.com",
							"health_check_http_code":       []string{"http_2xx", "http_3xx", "http_4xx"},
						},
					},
					"connection_drain":         "true",
					"connection_drain_timeout": "60",
					"tags": map[string]string{
						"Created": "tfTestAcc0",
						"For":     "Tftestacc 0",
					},
					"address_ip_version": "Ipv4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id":          CHECKSET,
						"server_group_name":          name,
						"server_group_type":          "Instance",
						"vpc_id":                     CHECKSET,
						"scheduler":                  "Wrr",
						"preserve_client_ip_enabled": "true",
						"protocol":                   "TCP",
						"connection_drain":           "true",
						"connection_drain_timeout":   "60",
						"address_ip_version":         "Ipv4",
						"health_check.#":             "1",
						"tags.%":                     "2",
						"tags.Created":               "tfTestAcc0",
						"tags.For":                   "Tftestacc 0",
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

var AlicloudNLBServerGroupMap0 = map[string]string{
	"server_group_name":        CHECKSET,
	"address_ip_version":       CHECKSET,
	"health_check.#":           CHECKSET,
	"protocol":                 CHECKSET,
	"server_group_type":        CHECKSET,
	"status":                   CHECKSET,
	"connection_drain":         CHECKSET,
	"connection_drain_timeout": CHECKSET,
	"scheduler":                CHECKSET,
	"vpc_id":                   CHECKSET,
}

func AlicloudNLBServerGroupBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_resource_manager_resource_groups" "default" {}
`, name)
}

func TestAccAlicloudNLBServerGroup_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_server_group.default"
	ra := resourceAttrInit(resourceId, AlicloudNLBServerGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNlbServerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbservergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNLBServerGroupBasicDependence0)
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
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"server_group_name": "${var.name}",
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":            CHECKSET,
						"server_group_name": name,
						"health_check.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "Wrr",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": "Wrr",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check": []map[string]interface{}{
						{
							"health_check_enabled":         "true",
							"health_check_type":            "TCP",
							"health_check_connect_port":    "0",
							"healthy_threshold":            "2",
							"unhealthy_threshold":          "2",
							"health_check_connect_timeout": "5",
							"health_check_interval":        "10",
							"http_check_method":            "GET",
							"health_check_url":             "/test/index.html",
							"health_check_domain":          "tf-testAcc.com",
							"health_check_http_code":       []string{"http_2xx", "http_3xx", "http_4xx"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain":         "true",
					"connection_drain_timeout": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_drain":         "true",
						"connection_drain_timeout": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "tfTestAcc5",
						"For":     "Tftestacc 5",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "tfTestAcc5",
						"tags.For":     "Tftestacc 5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name":        "${var.name}",
					"scheduler":                "Rr",
					"connection_drain_timeout": "100",
					"tags": map[string]string{
						"Created": "tfTestAcc6",
						"For":     "Tftestacc 6",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name":        name,
						"scheduler":                "Rr",
						"connection_drain_timeout": "100",
						"tags.%":                   "2",
						"tags.Created":             "tfTestAcc6",
						"tags.For":                 "Tftestacc 6",
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

func TestUnitAccAlicloudNlbServerGroup(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_nlb_server_group"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_nlb_server_group"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"resource_group_id": "CreateNlbServerGroupValue",
		"server_group_name": "CreateNlbServerGroupValue",
		"server_group_type": "CreateNlbServerGroupValue",
		"vpc_id":            "CreateNlbServerGroupValue",
		"scheduler":         "CreateNlbServerGroupValue",
		"protocol":          "CreateNlbServerGroupValue",
		"health_check": []map[string]interface{}{
			{
				"health_check_enabled":         true,
				"health_check_type":            "CreateNlbServerGroupValue",
				"health_check_connect_port":    0,
				"healthy_threshold":            2,
				"health_check_url":             "CreateNlbServerGroupValue",
				"health_check_domain":          "CreateNlbServerGroupValue",
				"unhealthy_threshold":          2,
				"health_check_connect_timeout": 5,
				"health_check_interval":        10,
				"http_check_method":            "CreateNlbServerGroupValue",
				"health_check_http_code":       []string{"CreateNlbServerGroupValue"},
			},
		},
		"connection_drain":           true,
		"connection_drain_timeout":   60,
		"address_ip_version":         "CreateNlbServerGroupValue",
		"preserve_client_ip_enabled": true,
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
		"ServerGroups": []interface{}{
			map[string]interface{}{
				"RegionId":         "CreateNlbServerGroupValue",
				"ServerGroupId":    "CreateNlbServerGroupValue",
				"ServerGroupName":  "CreateNlbServerGroupValue",
				"ServerGroupType":  "CreateNlbServerGroupValue",
				"AddressIPVersion": "CreateNlbServerGroupValue",
				"VpcId":            "CreateNlbServerGroupValue",
				"Scheduler":        "CreateNlbServerGroupValue",
				"Protocol":         "CreateNlbServerGroupValue",
				"HealthCheck": map[string]interface{}{
					"HealthCheckEnabled":        true,
					"HealthCheckType":           "CreateNlbServerGroupValue",
					"HealthCheckConnectPort":    0,
					"HealthyThreshold":          2,
					"UnhealthyThreshold":        2,
					"HealthCheckConnectTimeout": 5,
					"HealthCheckInterval":       10,
					"HealthCheckDomain":         "CreateNlbServerGroupValue",
					"HealthCheckUrl":            "CreateNlbServerGroupValue",
					"HealthCheckHttpCode":       []string{"CreateNlbServerGroupValue"},
					"HttpCheckMethod":           "CreateNlbServerGroupValue",
				},
				"ConnectionDrainEnabled":  true,
				"ConnectionDrainTimeout":  60,
				"PreserveClientIpEnabled": true,
				"ResourceGroupId":         "CreateNlbServerGroupValue",
				"ServerGroupStatus":       "Available",
				"ServerCount":             2,
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		"ServerGroupId": "CreateNlbServerGroupValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_nlb_server_group", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewNlbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudNlbServerGroupCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateServerGroup" {
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
		err := resourceAlicloudNlbServerGroupCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_nlb_server_group"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewNlbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudNlbServerGroupUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"scheduler":                "UpdateNlbServerGroupValue",
		"connection_drain":         false,
		"connection_drain_timeout": 100,
		"health_check": []map[string]interface{}{
			{
				"health_check_enabled":         true,
				"health_check_type":            "UpdateNlbServerGroupValue",
				"health_check_connect_port":    0,
				"healthy_threshold":            2,
				"health_check_url":             "UpdateNlbServerGroupValue",
				"health_check_domain":          "UpdateNlbServerGroupValue",
				"unhealthy_threshold":          2,
				"health_check_connect_timeout": 5,
				"health_check_interval":        10,
				"http_check_method":            "UpdateNlbServerGroupValue",
				"health_check_http_code":       []string{"UpdateNlbServerGroupValue"},
			},
		},
	}
	diff, err := newInstanceDiff("alicloud_nlb_server_group", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_nlb_server_group"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"ServerGroups": []interface{}{
			map[string]interface{}{
				"Scheduler":              "UpdateNlbServerGroupValue",
				"ConnectionDrainEnabled": false,
				"ConnectionDrainTimeout": 100,
				"HealthCheck": map[string]interface{}{
					"HealthCheckEnabled":        true,
					"HealthCheckType":           "UpdateNlbServerGroupValue",
					"HealthCheckConnectPort":    0,
					"HealthyThreshold":          2,
					"UnhealthyThreshold":        2,
					"HealthCheckConnectTimeout": 5,
					"HealthCheckInterval":       10,
					"HealthCheckDomain":         "UpdateNlbServerGroupValue",
					"HealthCheckUrl":            "UpdateNlbServerGroupValue",
					"HealthCheckHttpCode":       []string{"UpdateNlbServerGroupValue"},
					"HttpCheckMethod":           "UpdateNlbServerGroupValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateServerGroupAttribute" {
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
		err := resourceAlicloudNlbServerGroupUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_nlb_server_group"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ServerGroupName
	attributesDiff = map[string]interface{}{
		"server_group_name": "UpdateNlbServerGroupValue",
	}
	diff, err = newInstanceDiff("alicloud_nlb_server_group", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_nlb_server_group"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"ServerGroups": []interface{}{
			map[string]interface{}{
				"ServerGroupName": "UpdateNlbServerGroupValue",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateServerGroupAttribute" {
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
		err := resourceAlicloudNlbServerGroupUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_nlb_server_group"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// TagResources
	attributesDiff = map[string]interface{}{
		"tags": map[string]interface{}{
			"TagResourcesValue_1": "TagResourcesValue_1",
			"TagResourcesValue_2": "TagResourcesValue_2",
		},
	}
	diff, err = newInstanceDiff("alicloud_nlb_server_group", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_nlb_server_group"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"ServerGroups": []interface{}{
			map[string]interface{}{
				"Tags": []interface{}{
					map[string]interface{}{
						"Key":   "TagResourcesValue_1",
						"Value": "TagResourcesValue_1",
					},
					map[string]interface{}{
						"Key":   "TagResourcesValue_2",
						"Value": "TagResourcesValue_2",
					},
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "TagResources" {
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
		err := resourceAlicloudNlbServerGroupUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_nlb_server_group"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// UntagResources
	attributesDiff = map[string]interface{}{
		"tags": map[string]interface{}{
			"UntagResourcesValue3_1": "UnTagResourcesValue3_1",
			"UntagResourcesValue3_2": "UnTagResourcesValue3_2",
		},
	}
	diff, err = newInstanceDiff("alicloud_nlb_server_group", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_nlb_server_group"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"ServerGroups": []interface{}{
			map[string]interface{}{
				"Tags": []interface{}{
					map[string]interface{}{
						"Key":   "UntagResourcesValue3_1",
						"Value": "UnTagResourcesValue3_1",
					},
					map[string]interface{}{
						"Key":   "UntagResourcesValue3_2",
						"Value": "UnTagResourcesValue3_2",
					},
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UntagResources" {
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
		err := resourceAlicloudNlbServerGroupUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_nlb_server_group"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_nlb_server_group", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_nlb_server_group"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ListServerGroups" {
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
		err := resourceAlicloudNlbServerGroupRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewNlbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudNlbServerGroupDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_nlb_server_group", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_nlb_server_group"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteServerGroup" {
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
			if *action == "ListServerGroups" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudNlbServerGroupDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
