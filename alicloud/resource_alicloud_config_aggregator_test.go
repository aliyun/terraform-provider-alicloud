package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// collection-order: bypass aggregator_accounts -- pre-existing attribute outside this change

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
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
			if err != nil {
				if NeedRetry(err) {
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
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err = client.RpcPost("Config", "2020-09-07", action, nil, deleteRequest, false)
		if err != nil {
			if NeedRetry(err) {
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

// testUnitConfigAggregatorMockMeta builds a credential-free *connectivity.AliyunClient by
// constructing a real client from a stub access-key Config (same approach as
// TestUnitAliCloudCdnServiceRead). The "config" endpoint is pre-seeded in the
// Endpoints sync.Map so AliyunClient.loadApiEndpoint returns immediately via
// formatEndpointWithAccountID and never falls into loadEndpoint (which would
// require real credentials/network). The real rpcRequest therefore reaches
// conn.DoRequest, which the test patches per action.
func testUnitConfigAggregatorMockMeta(t *testing.T) *connectivity.AliyunClient {
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoints.Store("config", "config.aliyuncs.com")
	endpoints.Store("Config", "config.aliyuncs.com")
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	meta, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	return meta
}

// lintignore: R001
func TestUnitAliCloudConfigAggregator(t *testing.T) {
	p := Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_config_aggregator"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_config_aggregator"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"aggregator_accounts": []interface{}{
			map[string]interface{}{
				"account_id":   "CreateAggregatorValue",
				"account_name": "CreateAggregatorValue",
				"account_type": "ResourceDirectory",
			},
		},
		"aggregator_name": "CreateAggregatorValue",
		"aggregator_type": "CUSTOM",
		"description":     "CreateAggregatorValue",
		"tags": []interface{}{
			map[string]interface{}{
				"key":   "tf-key1",
				"value": "tf-value1",
			},
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
	mockMeta := testUnitConfigAggregatorMockMeta(t)

	ReadMockResponse := map[string]interface{}{
		// GetAggregator
		"Aggregator": map[string]interface{}{
			"AggregatorAccounts": []interface{}{
				map[string]interface{}{
					"AccountId":   "CreateAggregatorValue",
					"AccountName": "CreateAggregatorValue",
					"AccountType": "ResourceDirectory",
				},
			},
			"AggregatorName":            "CreateAggregatorValue",
			"AggregatorType":            "CUSTOM",
			"Description":               "CreateAggregatorValue",
			"AggregatorStatus":          "1",
			"AggregatorCreateTimestamp": json.Number("1609459200000"),
			"Tags": []interface{}{
				map[string]interface{}{
					"TagKey":   "tf-key1",
					"TagValue": "tf-value1",
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateAggregator
		"AggregatorId": "MockCreateAggregatorValue",
	}
	UpdateReadMockResponse := map[string]interface{}{
		// GetAggregator after UpdateAggregator
		"Aggregator": map[string]interface{}{
			"AggregatorAccounts": []interface{}{
				map[string]interface{}{
					"AccountId":   "UpdateAggregatorValue",
					"AccountName": "UpdateAggregatorValue",
					"AccountType": "ResourceDirectory",
				},
			},
			"AggregatorName":            "UpdateAggregatorValue",
			"AggregatorType":            "CUSTOM",
			"Description":               "UpdateAggregatorValue",
			"AggregatorStatus":          "1",
			"AggregatorCreateTimestamp": json.Number("1609459200000"),
			"Tags": []interface{}{
				map[string]interface{}{
					"TagKey":   "tf-key2",
					"TagValue": "tf-value2",
				},
			},
		},
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String(errorCode),
			Data:       tea.String(errorCode),
			Message:    tea.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}

	// Create: success (CreateAggregator -> CreateMockResponse; state poll + Read -> ReadMockResponse with tags)
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		switch *action {
		case "CreateAggregator":
			return CreateMockResponse, nil
		default:
			return ReadMockResponse, nil
		}
	})
	err := resourceAliCloudConfigAggregatorCreate(dInit, mockMeta)
	patches.Reset()
	assert.Nil(t, err)
	assert.Equal(t, "MockCreateAggregatorValue", dInit.Id())
	assert.Equal(t, "tf-key1", dInit.Get("tags.0.key"))
	assert.Equal(t, "tf-value1", dInit.Get("tags.0.value"))

	// Create: non-retryable error
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		if *action == "CreateAggregator" {
			return failedResponseMock("NonRetryableError")
		}
		return ReadMockResponse, nil
	})
	err = resourceAliCloudConfigAggregatorCreate(dInit, mockMeta)
	patches.Reset()
	assert.NotNil(t, err)

	// Update: success (tags change tf-key1 -> tf-key2; aggregator_accounts carries a
	// different account_id/account_name than the Create step.)
	attributesDiff := map[string]interface{}{
		"aggregator_name": "UpdateAggregatorValue",
		"description":     "UpdateAggregatorValue",
		"aggregator_accounts": []interface{}{
			map[string]interface{}{
				"account_id":   "UpdateAggregatorValue",
				"account_name": "UpdateAggregatorValue",
				"account_type": "ResourceDirectory",
			},
		},
		"tags": []interface{}{
			map[string]interface{}{
				"key":   "tf-key2",
				"value": "tf-value2",
			},
		},
	}
	diff, err := newInstanceDiff("alicloud_config_aggregator", attributes, attributesDiff, dInit.State())
	assert.Nil(t, err)
	dExisted, _ = schema.InternalMap(p["alicloud_config_aggregator"].Schema).Data(dInit.State(), diff)
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		return UpdateReadMockResponse, nil
	})
	err = resourceAliCloudConfigAggregatorUpdate(dExisted, mockMeta)
	patches.Reset()
	assert.Nil(t, err)
	assert.Equal(t, "tf-key2", dExisted.Get("tags.0.key"))
	assert.Equal(t, "tf-value2", dExisted.Get("tags.0.value"))

	// Update: non-retryable error
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		if *action == "UpdateAggregator" {
			return failedResponseMock("NonRetryableError")
		}
		return UpdateReadMockResponse, nil
	})
	err = resourceAliCloudConfigAggregatorUpdate(dExisted, mockMeta)
	patches.Reset()
	assert.NotNil(t, err)

	// Delete: success
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		if *action == "DeleteAggregators" {
			return map[string]interface{}{}, nil
		}
		return ReadMockResponse, nil
	})
	err = resourceAliCloudConfigAggregatorDelete(dExisted, mockMeta)
	patches.Reset()
	assert.Nil(t, err)

	// Delete: AccountNotExisted swallowed
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		if *action == "DeleteAggregators" {
			return failedResponseMock("AccountNotExisted")
		}
		return ReadMockResponse, nil
	})
	err = resourceAliCloudConfigAggregatorDelete(dExisted, mockMeta)
	patches.Reset()
	assert.Nil(t, err)

	// Delete: non-retryable error
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		if *action == "DeleteAggregators" {
			return failedResponseMock("NonRetryableError")
		}
		return ReadMockResponse, nil
	})
	err = resourceAliCloudConfigAggregatorDelete(dExisted, mockMeta)
	patches.Reset()
	assert.NotNil(t, err)
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregator_name": name,
					"description":     name,
					"aggregator_type": "CUSTOM",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_name": name,
						"description":     name,
						"aggregator_type": "CUSTOM",
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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

// Case 账号组-资源测试-TAGS 11749
func TestAccAliCloudConfigAggregator_tags11749(t *testing.T) {
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregator_name": name,
					"description":     name,
					"aggregator_type": "RD",
					"tags": []map[string]interface{}{
						{"key": "tf-key1", "value": "tf-value1"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_name": name,
						"description":     name,
						"aggregator_type": "RD",
						"tags.0.key":      "tf-key1",
						"tags.0.value":    "tf-value1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregator_name": name,
					"description":     name,
					"aggregator_type": "RD",
					"tags": []map[string]interface{}{
						{"key": "tf-key2", "value": "tf-value2"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.0.key":   "tf-key2",
						"tags.0.value": "tf-value2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregator_name": name,
					"description":     name,
					"aggregator_type": "RD",
					"tags":            REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.#": "0",
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

// TestAccAliCloudConfigAggregator_tagsOrder establishes TypeList Order Coverage
// for the tags attribute: an apply of the original member order followed by a
// convergent apply with the members reordered. Config echoes the submitted Tag
// order, so the reordered apply persists and the post-apply plan is empty.
func TestAccAliCloudConfigAggregator_tagsOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregator.default"
	ra := resourceAttrInit(resourceId, AliCloudConfigAggregatorMap11746)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregator")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudConfigAggregatorBasicDependence11746)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregator_name": name,
					"description":     name,
					"aggregator_type": "RD",
					"tags": []map[string]interface{}{
						{"key": "tf-key1", "value": "tf-value1"},
						{"key": "tf-key2", "value": "tf-value2"},
					},
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregator_name": name,
					"description":     name,
					"aggregator_type": "RD",
					"tags": []map[string]interface{}{
						{"key": "tf-key2", "value": "tf-value2"},
						{"key": "tf-key1", "value": "tf-value1"},
					},
				}),
			},
		},
	})
}

// Test Config Aggregator. <<< Resource test cases, automatically generated.
