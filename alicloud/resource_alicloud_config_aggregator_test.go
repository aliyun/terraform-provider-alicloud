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
	resource.AddTestSweepers("alicloud_config_aggregator", &resource.Sweeper{
		Name: "alicloud_config_aggregator",
		F:    testSweepConfigAggregator,
		Dependencies: []string{
			"alicloud_config_aggregate_compliance_pack",
			"alicloud_config_aggregate_config_rule",
		},
	})
}

func testSweepConfigAggregator(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	// Get all AggregatorId
	aggregatorIds := make([]string, 0)
	action := "ListAggregators"
	request := map[string]interface{}{
		"MaxResults": PageSizeLarge,
	}
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Println("List Config Aggregator Failed!", err)
			return nil
		}
		resp, err := jsonpath.Get("$.AggregatorsResult.Aggregators", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AggregatorsResult.Aggregators", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["AggregatorName"])), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Aggregate: %v (%v)", item["AggregatorName"], item["AggregatorId"])
				continue
			}
			aggregatorIds = append(aggregatorIds, fmt.Sprint(item["AggregatorId"]))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	log.Printf("[INFO] Deleting Aggregate:  (%s)", strings.Join(aggregatorIds, ","))
	action = "DeleteAggregators"
	deleteRequest := map[string]interface{}{
		"AggregatorIds": strings.Join(aggregatorIds, ","),
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = client.RpcPost("Config", "2020-09-07", action, nil, deleteRequest, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[INFO] Delete Aggregate Failed:  (%s)", strings.Join(aggregatorIds, ","))
	}
	return nil
}

func TestUnitAliCloudConfigAggregator(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_config_aggregator"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_config_aggregator"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"aggregator_accounts": []interface{}{
			map[string]interface{}{
				"account_id":   "CreateAggregatorValue",
				"account_name": "CreateAggregatorValue",
				"account_type": "CreateAggregatorValue",
			},
		},
		"aggregator_name": "CreateAggregatorValue",
		"aggregator_type": "CreateAggregatorValue",
		"description":     "CreateAggregatorValue",
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
		// GetAggregator
		"Aggregator": map[string]interface{}{
			"AggregatorAccounts": []interface{}{
				map[string]interface{}{
					"AccountId":   "CreateAggregatorValue",
					"AccountName": "CreateAggregatorValue",
					"AccountType": "CreateAggregatorValue",
				},
			},
			"AggregatorName":   "CreateAggregatorValue",
			"AggregatorType":   "CreateAggregatorValue",
			"Description":      "CreateAggregatorValue",
			"AggregatorStatus": "1",
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateAggregator
		"AggregatorId": "MockCreateAggregatorValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_config_aggregator", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create

	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewConfigClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudConfigAggregatorCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// GetAggregator Response
		"Aggregator": map[string]interface{}{
			"AggregatorId": "MockCreateAggregatorValue",
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateAggregator" {
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
		err := resourceAliCloudConfigAggregatorCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_aggregator"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewConfigClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudConfigAggregatorUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateAggregator
	attributesDiff := map[string]interface{}{
		"aggregator_accounts": []interface{}{
			map[string]interface{}{
				"account_id":   "UpdateAggregatorValue",
				"account_name": "UpdateAggregatorValue",
				"account_type": "UpdateAggregatorValue",
			},
		},
		"aggregator_name": "UpdateAggregatorValue",
		"description":     "UpdateAggregatorValue",
	}
	diff, err := newInstanceDiff("alicloud_config_aggregator", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_aggregator"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetAggregator Response
		"Aggregator": map[string]interface{}{
			"AggregatorAccounts": []interface{}{
				map[string]interface{}{
					"AccountId":   "UpdateAggregatorValue",
					"AccountName": "UpdateAggregatorValue",
					"AccountType": "UpdateAggregatorValue",
				},
			},
			"AggregatorName": "UpdateAggregatorValue",
			"Description":    "UpdateAggregatorValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateAggregator" {
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
		err := resourceAliCloudConfigAggregatorUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_aggregator"].Schema).Data(dExisted.State(), nil)
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
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "Invalid.AggregatorId.Value", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetAggregator" {
				switch errorCode {
				case "{}", "Invalid.AggregatorId.Value":
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
		err := resourceAliCloudConfigAggregatorRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}", "Invalid.AggregatorId.Value":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewConfigClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudConfigAggregatorDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "AccountNotExisted", "Invalid.AggregatorIds.Empty"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteAggregators" {
				switch errorCode {
				case "NonRetryableError", "AccountNotExisted", "Invalid.AggregatorIds.Empty":
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
		err := resourceAliCloudConfigAggregatorDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "AccountNotExisted", "Invalid.AggregatorIds.Empty":
			assert.Nil(t, err)
		}
	}
}

// Test Config Aggregator. >>> Resource test cases, automatically generated.
// Case 账号组-资源测试-CUSTOM 11746
func TestAccAliCloudConfigAggregator_basic11746(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregator.default"
	ra := resourceAttrInit(resourceId, AliCloudConfigAggregatorMap11746)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudConfigAggregatorBasicDependence11746)
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
					"aggregator_name": name,
					"description":     name,
					"aggregator_type": "RD",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_name": name,
						"description":     name,
						"aggregator_type": "RD",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregator_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
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

func TestAccAliCloudConfigAggregator_basic11746_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregator.default"
	ra := resourceAttrInit(resourceId, AliCloudConfigAggregatorMap11746)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudConfigAggregatorBasicDependence11746)
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
					"aggregator_name": name,
					"description":     name,
					"aggregator_type": "RD",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_name": name,
						"description":     name,
						"aggregator_type": "RD",
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

func AliCloudConfigAggregatorBasicDependence11746(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}

var AliCloudConfigAggregatorMap11746 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

// Case 账号组-资源测试-CUSTOM 11747
func TestAccAliCloudConfigAggregator_basic11747(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregator.default"
	ra := resourceAttrInit(resourceId, AliCloudConfigAggregatorMap11746)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudConfigAggregatorBasicDependence11747)
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
					"aggregator_name": name,
					"description":     name,
					"aggregator_accounts": []map[string]interface{}{
						{
							"account_id":   "${data.alicloud_resource_manager_accounts.default.accounts.0.account_id}",
							"account_name": "${data.alicloud_resource_manager_accounts.default.accounts.0.display_name}",
							"account_type": "ResourceDirectory",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_name":       name,
						"description":           name,
						"aggregator_accounts.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregator_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregator_accounts": []map[string]interface{}{
						{
							"account_id":   "${data.alicloud_resource_manager_accounts.default.accounts.1.account_id}",
							"account_name": "${data.alicloud_resource_manager_accounts.default.accounts.1.display_name}",
							"account_type": "ResourceDirectory",
						},
						{
							"account_id":   "${data.alicloud_resource_manager_accounts.default.accounts.2.account_id}",
							"account_name": "${data.alicloud_resource_manager_accounts.default.accounts.2.display_name}",
							"account_type": "ResourceDirectory",
						},
						{
							"account_id":   "${data.alicloud_resource_manager_accounts.default.accounts.3.account_id}",
							"account_name": "${data.alicloud_resource_manager_accounts.default.accounts.3.display_name}",
							"account_type": "ResourceDirectory",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_accounts.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregator_accounts": []map[string]interface{}{
						{
							"account_id":   "${data.alicloud_resource_manager_accounts.default.accounts.0.account_id}",
							"account_name": "${data.alicloud_resource_manager_accounts.default.accounts.0.display_name}",
							"account_type": "ResourceDirectory",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_accounts.#": "1",
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

func TestAccAliCloudConfigAggregator_basic11747_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregator.default"
	ra := resourceAttrInit(resourceId, AliCloudConfigAggregatorMap11746)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudConfigAggregatorBasicDependence11747)
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
					"aggregator_name": name,
					"description":     name,
					"aggregator_type": "CUSTOM",
					"aggregator_accounts": []map[string]interface{}{
						{
							"account_id":   "${data.alicloud_resource_manager_accounts.default.accounts.1.account_id}",
							"account_name": "${data.alicloud_resource_manager_accounts.default.accounts.1.display_name}",
							"account_type": "ResourceDirectory",
						},
						{
							"account_id":   "${data.alicloud_resource_manager_accounts.default.accounts.2.account_id}",
							"account_name": "${data.alicloud_resource_manager_accounts.default.accounts.2.display_name}",
							"account_type": "ResourceDirectory",
						},
						{
							"account_id":   "${data.alicloud_resource_manager_accounts.default.accounts.3.account_id}",
							"account_name": "${data.alicloud_resource_manager_accounts.default.accounts.3.display_name}",
							"account_type": "ResourceDirectory",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_name":       name,
						"description":           name,
						"aggregator_type":       "CUSTOM",
						"aggregator_accounts.#": "3",
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

func AliCloudConfigAggregatorBasicDependence11747(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_accounts" "default" {
  status  = "CreateSuccess"
}
`, name)
}

// Case 账号组-资源测试-基础用例-FOLDER 10748
func TestAccAliCloudConfigAggregator_basic11748(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregator.default"
	ra := resourceAttrInit(resourceId, AliCloudConfigAggregatorMap11746)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudConfigAggregatorBasicDependence11748)
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
					"aggregator_name": name,
					"description":     name,
					"aggregator_type": "FOLDER",
					"folder_id":       "${data.alicloud_resource_manager_folders.default.folders.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_name": name,
						"description":     name,
						"aggregator_type": "FOLDER",
						"folder_id":       CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregator_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"folder_id": fmt.Sprintf("%s,%s,%s", "${data.alicloud_resource_manager_folders.default.folders.1.id}", "${data.alicloud_resource_manager_folders.default.folders.2.id}", "${data.alicloud_resource_manager_folders.default.folders.3.id}"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"folder_id": CHECKSET,
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

func TestAccAliCloudConfigAggregator_basic11748_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregator.default"
	ra := resourceAttrInit(resourceId, AliCloudConfigAggregatorMap11746)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudConfigAggregatorBasicDependence11748)
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
					"aggregator_name": name,
					"description":     name,
					"aggregator_type": "FOLDER",
					"folder_id":       fmt.Sprintf("%s,%s,%s", "${data.alicloud_resource_manager_folders.default.folders.1.id}", "${data.alicloud_resource_manager_folders.default.folders.2.id}", "${data.alicloud_resource_manager_folders.default.folders.3.id}"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_name": name,
						"description":     name,
						"aggregator_type": "FOLDER",
						"folder_id":       CHECKSET,
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

func AliCloudConfigAggregatorBasicDependence11748(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_folders" "default" {
}
`, name)
}

// Test Config Aggregator. <<< Resource test cases, automatically generated.
