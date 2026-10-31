package config_test

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/config"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
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
	rawClient, err := sweep.SharedClientForRegion(region)
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
		"MaxResults": helper.PageSizeLarge,
	}
	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Println("List Config Aggregator Failed!", err)
			return nil
		}
		resp, err := jsonpath.Get("$.AggregatorsResult.Aggregators", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.AggregatorsResult.Aggregators", response)
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err = client.RpcPost("Config", "2020-09-07", action, nil, deleteRequest, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[INFO] Delete Aggregate Failed:  (%s)", strings.Join(aggregatorIds, ","))
	}
	return nil
}

// lintignore: R001
func TestUnitAliCloudConfigAggregator(t *testing.T) {
	p := acctest.Provider.ResourcesMap
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
	rawClient, err := sweep.SharedClientForRegion(region)
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
			Code:       helper.String(errorCode),
			Data:       helper.String(errorCode),
			Message:    helper.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_config_aggregator", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create

	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewConfigClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    helper.String("loadEndpoint error"),
			Data:    helper.String("loadEndpoint error"),
			Message: helper.String("loadEndpoint error"),
		}
	})
	err = config.ResourceAggregatorCreate(dInit, rawClient)
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
		err := config.ResourceAggregatorCreate(dInit, rawClient)
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
			Code:    helper.String("loadEndpoint error"),
			Data:    helper.String("loadEndpoint error"),
			Message: helper.String("loadEndpoint error"),
		}
	})
	err = config.ResourceAggregatorUpdate(dExisted, rawClient)
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
	diff, err := acctest.NewInstanceDiff("alicloud_config_aggregator", attributes, attributesDiff, dInit.State())
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
		err := config.ResourceAggregatorUpdate(dExisted, rawClient)
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
		err := config.ResourceAggregatorRead(dExisted, rawClient)
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
			Code:    helper.String("loadEndpoint error"),
			Data:    helper.String("loadEndpoint error"),
			Message: helper.String("loadEndpoint error"),
		}
	})
	err = config.ResourceAggregatorDelete(dExisted, rawClient)
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
		err := config.ResourceAggregatorDelete(dExisted, rawClient)
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudConfigAggregatorMap11746)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return config.NewConfigServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeConfigAggregator")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccconfig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudConfigAggregatorBasicDependence11746)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudConfigAggregatorMap11746)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return config.NewConfigServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeConfigAggregator")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccconfig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudConfigAggregatorBasicDependence11746)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

// Case 账号组-资源测试-CUSTOM 11747
func TestAccAliCloudConfigAggregator_basic11747(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregator.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudConfigAggregatorMap11746)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return config.NewConfigServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeConfigAggregator")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccconfig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudConfigAggregatorBasicDependence11747)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudConfigAggregatorMap11746)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return config.NewConfigServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeConfigAggregator")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccconfig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudConfigAggregatorBasicDependence11747)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudConfigAggregatorMap11746)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return config.NewConfigServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeConfigAggregator")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccconfig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudConfigAggregatorBasicDependence11748)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"folder_id":       acctest.CHECKSET,
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
						"folder_id": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudConfigAggregatorMap11746)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return config.NewConfigServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeConfigAggregator")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccconfig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudConfigAggregatorBasicDependence11748)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"folder_id":       acctest.CHECKSET,
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
