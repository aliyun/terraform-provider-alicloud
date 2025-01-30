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
		"alicloud_cms_hybrid_monitor_sls_task",
		&resource.Sweeper{
			Name: "alicloud_cms_hybrid_monitor_sls_task",
			F:    testSweepCmsHybridMonitorSlsTask,
		})
}

func testSweepCmsHybridMonitorSlsTask(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.CloudMonitorServiceSupportRegions) {
		log.Printf("[INFO] Skipping Cms Hybrid Monitor Sls Task unsupported region: %s", region)
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeHybridMonitorTaskList"
	request := map[string]interface{}{}

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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

		resp, err := jsonpath.Get("$.TaskList", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.TaskList", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(item["TaskName"].(string)), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Cms Hybrid Monitor Sls Task: %s", item["TaskName"].(string))
					continue
				}
			}
			action := "DeleteHybridMonitorTask"
			request := map[string]interface{}{
				"TaskId": item["TaskId"],
			}
			_, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Cms Hybrid Monitor Sls Task (%s): %s", item["TaskName"].(string), err)
			}
			log.Printf("[INFO] Delete Cms Hybrid Monitor Sls Task success: %s ", item["TaskName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudCloudMonitorServiceHybridMonitorSlsTask_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_hybrid_monitor_sls_task.default"
	checkoutSupportedRegions(t, true, connectivity.CloudMonitorServiceSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceHybridMonitorSlsTaskMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsHybridMonitorSlsTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_cms_slstask%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceHybridMonitorSlsTaskBasicDependence0)
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
					"sls_process_config": []map[string]interface{}{
						{
							"filter": []map[string]interface{}{
								{
									"relation": "and",
									"filters": []map[string]interface{}{
										{
											"operator":     "=",
											"value":        "200",
											"sls_key_name": "code",
										},
									},
								},
							},
							"statistics": []map[string]interface{}{
								{
									"function":      "count",
									"alias":         "level_count",
									"sls_key_name":  "name",
									"parameter_one": "200",
									"parameter_two": "299",
								},
							},
							"group_by": []map[string]interface{}{
								{
									"alias":        "code",
									"sls_key_name": "ApiResult",
								},
							},
							"express": []map[string]interface{}{
								{
									"express": "success_count",
									"alias":   "SuccRate",
								},
							},
						},
					},
					"task_name":           "${var.name}",
					"namespace":           "${alicloud_cms_namespace.default.id}",
					"description":         "${var.name}",
					"collect_interval":    "60",
					"collect_target_type": "${alicloud_cms_sls_group.default.id}",
					"attach_labels": []map[string]interface{}{
						{
							"name":  "app_service",
							"value": "testValue",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sls_process_config.#": "1",
						"task_name":            name,
						"namespace":            CHECKSET,
						"description":          name,
						"collect_interval":     "60",
						"collect_target_type":  CHECKSET,
						"attach_labels.#":      "1",
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

var AlicloudCloudMonitorServiceHybridMonitorSlsTaskMap0 = map[string]string{
	"attach_labels.#":      CHECKSET,
	"sls_process_config.#": CHECKSET,
}

func AlicloudCloudMonitorServiceHybridMonitorSlsTaskBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_account" "this" {}

resource "alicloud_cms_sls_group" "default" {
	sls_group_config {
		sls_user_id = data.alicloud_account.this.id
		sls_logstore = "Logstore-ECS"
		sls_project = "aliyun-project"
		sls_region = "cn-hangzhou"
	}
	sls_group_description = var.name
	sls_group_name = var.name
}

resource "alicloud_cms_namespace" "default" {
	description = var.name
	namespace = "tf-testacc-cloudmonitorservicenamespace"
	specification = "cms.s1.large"
}
`, name)
}

func TestAccAlicloudCloudMonitorServiceHybridMonitorSlsTask_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_hybrid_monitor_sls_task.default"
	checkoutSupportedRegions(t, true, connectivity.CloudMonitorServiceSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceHybridMonitorSlsTaskMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsHybridMonitorSlsTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_cms_slstask%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceHybridMonitorSlsTaskBasicDependence0)
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
					"task_name":           "${var.name}",
					"namespace":           "${alicloud_cms_namespace.default.id}",
					"collect_target_type": "${alicloud_cms_sls_group.default.id}",
					"sls_process_config": []map[string]interface{}{
						{
							"statistics": []map[string]interface{}{
								{
									"function":      "count",
									"alias":         "level_count",
									"sls_key_name":  "name",
									"parameter_one": "200",
									"parameter_two": "299",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"collect_target_type":  CHECKSET,
						"task_name":            name,
						"namespace":            CHECKSET,
						"sls_process_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sls_process_config": []map[string]interface{}{
						{
							"statistics": []map[string]interface{}{
								{
									"function":      "count",
									"alias":         "level_count",
									"sls_key_name":  "name",
									"parameter_one": "100",
									"parameter_two": "299",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sls_process_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sls_process_config": []map[string]interface{}{
						{
							"filter": []map[string]interface{}{
								{
									"relation": "and",
									"filters": []map[string]interface{}{
										{
											"operator":     "=",
											"value":        "200",
											"sls_key_name": "code",
										},
									},
								},
							},
							"statistics": []map[string]interface{}{
								{
									"function":      "count",
									"alias":         "level_count",
									"sls_key_name":  "name",
									"parameter_one": "100",
									"parameter_two": "199",
								},
							},
							"group_by": []map[string]interface{}{
								{
									"alias":        "code",
									"sls_key_name": "ApiResult",
								},
							},
							"express": []map[string]interface{}{
								{
									"express": "success_count",
									"alias":   "SuccRate",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sls_process_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"collect_interval": "15",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"collect_interval": "15",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"attach_labels": []map[string]interface{}{
						{
							"name":  "app_service",
							"value": "testValue",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"attach_labels.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"attach_labels": []map[string]interface{}{
						{
							"name":  "app_service1",
							"value": "testValue1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"attach_labels.#": "1",
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

func TestUnitAccAlicloudCmsHybridMonitorSlsTask(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_cms_hybrid_monitor_sls_task"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cms_hybrid_monitor_sls_task"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"sls_process_config": []map[string]interface{}{
			{
				"filter": []map[string]interface{}{
					{
						"relation": "CreateCmsHybridMonitorSlsTaskValue",
						"filters": []map[string]interface{}{
							{
								"operator":     "CreateCmsHybridMonitorSlsTaskValue",
								"value":        "CreateCmsHybridMonitorSlsTaskValue",
								"sls_key_name": "CreateCmsHybridMonitorSlsTaskValue",
							},
						},
					},
				},
				"statistics": []map[string]interface{}{
					{
						"function":      "CreateCmsHybridMonitorSlsTaskValue",
						"alias":         "CreateCmsHybridMonitorSlsTaskValue",
						"sls_key_name":  "CreateCmsHybridMonitorSlsTaskValue",
						"parameter_one": "CreateCmsHybridMonitorSlsTaskValue",
						"parameter_two": "CreateCmsHybridMonitorSlsTaskValue",
					},
				},
				"group_by": []map[string]interface{}{
					{
						"alias":        "CreateCmsHybridMonitorSlsTaskValue",
						"sls_key_name": "CreateCmsHybridMonitorSlsTaskValue",
					},
				},
				"express": []map[string]interface{}{
					{
						"express": "CreateCmsHybridMonitorSlsTaskValue",
						"alias":   "CreateCmsHybridMonitorSlsTaskValue",
					},
				},
			},
		},
		"task_name":           "CreateCmsHybridMonitorSlsTaskValue",
		"namespace":           "CreateCmsHybridMonitorSlsTaskValue",
		"description":         "CreateCmsHybridMonitorSlsTaskValue",
		"collect_interval":    60,
		"collect_target_type": "CreateCmsHybridMonitorSlsTaskValue",
		"attach_labels": []map[string]interface{}{
			{
				"name":  "CreateCmsHybridMonitorSlsTaskValue",
				"value": "CreateCmsHybridMonitorSlsTaskValue",
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
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}

	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"TaskList": []interface{}{
			map[string]interface{}{
				"CollectInterval": 60,
				"Description":     "CreateCmsHybridMonitorSlsTaskValue",
				"TaskId":          "CmsHybridMonitorSlsTaskId",
				"AttachLabels": []interface{}{
					map[string]interface{}{
						"Name":  "CreateCmsHybridMonitorSlsTaskValue",
						"Value": "CreateCmsHybridMonitorSlsTaskValue",
					},
				},
				"CollectTargetPath":     "CreateCmsHybridMonitorSlsTaskValue",
				"Namespace":             "CreateCmsHybridMonitorSlsTaskValue",
				"GroupId":               0,
				"MatchExpressRelation":  "CreateCmsHybridMonitorSlsTaskValue",
				"CollectTargetEndpoint": "CreateCmsHybridMonitorSlsTaskValue",
				"TaskName":              "CreateCmsHybridMonitorSlsTaskValue",
				"UploadRegion":          "CreateCmsHybridMonitorSlsTaskValue",
				"NetworkType":           "CreateCmsHybridMonitorSlsTaskValue",
				"TaskType":              "CreateCmsHybridMonitorSlsTaskValue",
				"CollectTargetType":     "CreateCmsHybridMonitorSlsTaskValue",
				"SLSProcessConfig": map[string]interface{}{
					"GroupBy": []interface{}{
						map[string]interface{}{
							"Alias":      "CreateCmsHybridMonitorSlsTaskValue",
							"SLSKeyName": "CreateCmsHybridMonitorSlsTaskValue",
						},
					},
					"Express": []interface{}{
						map[string]interface{}{
							"Express": "CreateCmsHybridMonitorSlsTaskValue",
							"Alias":   "CreateCmsHybridMonitorSlsTaskValue",
						},
					},
					"Filter": map[string]interface{}{
						"Relation": "CreateCmsHybridMonitorSlsTaskValue",
						"Filters": []interface{}{
							map[string]interface{}{
								"Operator":   "CreateCmsHybridMonitorSlsTaskValue",
								"Value":      "CreateCmsHybridMonitorSlsTaskValue",
								"SLSKeyName": "CreateCmsHybridMonitorSlsTaskValue",
							},
						},
					},
					"Statistics": []interface{}{
						map[string]interface{}{
							"Parameter2": "CreateCmsHybridMonitorSlsTaskValue",
							"Function":   "CreateCmsHybridMonitorSlsTaskValue",
							"Parameter1": "CreateCmsHybridMonitorSlsTaskValue",
							"Alias":      "CreateCmsHybridMonitorSlsTaskValue",
							"SLSKeyName": "CreateCmsHybridMonitorSlsTaskValue",
						},
					},
				},
				"TargetUserId": "CreateCmsHybridMonitorSlsTaskValue",
			},
		},
		"Success": true,
	}
	CreateMockResponse := map[string]interface{}{
		"TaskId": "CmsHybridMonitorSlsTaskId",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cms_hybrid_monitor_sls_task", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudCmsHybridMonitorSlsTaskCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateHybridMonitorTask" {
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
		err := resourceAlicloudCmsHybridMonitorSlsTaskCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cms_hybrid_monitor_sls_task"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudCmsHybridMonitorSlsTaskUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{

		"sls_process_config": []map[string]interface{}{
			{
				"filter": []map[string]interface{}{
					{
						"relation": "UpdateCmsHybridMonitorSlsTaskValue",
						"filters": []map[string]interface{}{
							{
								"operator":     "UpdateCmsHybridMonitorSlsTaskValue",
								"value":        "UpdateCmsHybridMonitorSlsTaskValue",
								"sls_key_name": "UpdateCmsHybridMonitorSlsTaskValue",
							},
						},
					},
				},
				"statistics": []map[string]interface{}{
					{
						"function":      "UpdateCmsHybridMonitorSlsTaskValue",
						"alias":         "UpdateCmsHybridMonitorSlsTaskValue",
						"sls_key_name":  "UpdateCmsHybridMonitorSlsTaskValue",
						"parameter_one": "UpdateCmsHybridMonitorSlsTaskValue",
						"parameter_two": "UpdateCmsHybridMonitorSlsTaskValue",
					},
				},
				"group_by": []map[string]interface{}{
					{
						"alias":        "UpdateCmsHybridMonitorSlsTaskValue",
						"sls_key_name": "UpdateCmsHybridMonitorSlsTaskValue",
					},
				},
				"express": []map[string]interface{}{
					{
						"express": "UpdateCmsHybridMonitorSlsTaskValue",
						"alias":   "UpdateCmsHybridMonitorSlsTaskValue",
					},
				},
			},
		},
		"description":      "UpdateCmsHybridMonitorSlsTaskValue",
		"collect_interval": 15,
		"attach_labels": []map[string]interface{}{
			{
				"name":  "UpdateCmsHybridMonitorSlsTaskValue",
				"value": "UpdateCmsHybridMonitorSlsTaskValue",
			},
		},
	}
	diff, err := newInstanceDiff("alicloud_cms_hybrid_monitor_sls_task", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cms_hybrid_monitor_sls_task"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"TaskList": []interface{}{
			map[string]interface{}{
				"CollectInterval": 15,
				"Description":     "UpdateCmsHybridMonitorSlsTaskValue",
				"AttachLabels": []interface{}{
					map[string]interface{}{
						"Name":  "UpdateCmsHybridMonitorSlsTaskValue",
						"Value": "UpdateCmsHybridMonitorSlsTaskValue",
					},
				},
				"SLSProcessConfig": map[string]interface{}{
					"GroupBy": []interface{}{
						map[string]interface{}{
							"Alias":      "UpdateCmsHybridMonitorSlsTaskValue",
							"SLSKeyName": "UpdateCmsHybridMonitorSlsTaskValue",
						},
					},
					"Express": []interface{}{
						map[string]interface{}{
							"Express": "UpdateCmsHybridMonitorSlsTaskValue",
							"Alias":   "UpdateCmsHybridMonitorSlsTaskValue",
						},
					},
					"Filter": map[string]interface{}{
						"Relation": "UpdateCmsHybridMonitorSlsTaskValue",
						"Filters": []interface{}{
							map[string]interface{}{
								"Operator":   "UpdateCmsHybridMonitorSlsTaskValue",
								"Value":      "UpdateCmsHybridMonitorSlsTaskValue",
								"SLSKeyName": "UpdateCmsHybridMonitorSlsTaskValue",
							},
						},
					},
					"Statistics": []interface{}{
						map[string]interface{}{
							"Parameter2": "UpdateCmsHybridMonitorSlsTaskValue",
							"Function":   "UpdateCmsHybridMonitorSlsTaskValue",
							"Parameter1": "UpdateCmsHybridMonitorSlsTaskValue",
							"Alias":      "UpdateCmsHybridMonitorSlsTaskValue",
							"SLSKeyName": "UpdateCmsHybridMonitorSlsTaskValue",
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
			if *action == "ModifyHybridMonitorTask" {
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
		err := resourceAlicloudCmsHybridMonitorSlsTaskUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cms_hybrid_monitor_sls_task"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_cms_hybrid_monitor_sls_task", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cms_hybrid_monitor_sls_task"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeHybridMonitorTaskList" {
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
		err := resourceAlicloudCmsHybridMonitorSlsTaskRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudCmsHybridMonitorSlsTaskDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_cms_hybrid_monitor_sls_task", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cms_hybrid_monitor_sls_task"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteHybridMonitorTask" {
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
		err := resourceAlicloudCmsHybridMonitorSlsTaskDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
