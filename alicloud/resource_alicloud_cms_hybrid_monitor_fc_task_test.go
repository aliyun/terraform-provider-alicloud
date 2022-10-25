package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudMonitorServiceHybridMonitorFcTask_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_hybrid_monitor_fc_task.default"
	checkoutSupportedRegions(t, true, connectivity.CloudMonitorServiceSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceHybridMonitorFcTaskMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsHybridMonitorFcTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacccmshybridmonitorfctask%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceHybridMonitorFcTaskBasicDependence0)
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
					"namespace":      "${alicloud_cms_namespace.default.id}",
					"target_user_id": "${data.alicloud_account.this.id}",
					"yarm_config":    "---\\nproducts:\\n- namespace: \\\"acs_ecs_dashboard1\\\"\\n  metric_info:\\n  - metric_list:\\n    - \\\"CPUUtilization\\\"\\n    - \\\"DiskReadBPS\\\"\\n    - \\\"InternetOut\\\"\\n    - \\\"IntranetOut\\\"\\n    - \\\"cpu_idle\\\"\\n    - \\\"cpu_system\\\"\\n    - \\\"cpu_total\\\"\\n    - \\\"diskusage_utilization\\\"\\n- namespace: \\\"acs_rds_dashboard\\\"\\n  metric_info:\\n  - metric_list:\\n    - \\\"MySQL_QPS\\\"\\n    - \\\"MySQL_TPS\\\"\\n- namespace: \\\"acs_ecs_dashboard\\\"\\n  metric_info:\\n  - metric_list:\\n    - \\\"cpu_total\\\"\\n    - \\\"diskusage_utilization\\\"\\n    - \\\"memory_usedutilization\\\"\\n    - \\\"net_tcpconnection\\\"\\n",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_user_id": CHECKSET,
						"namespace":      CHECKSET,
						"yarm_config":    "---\nproducts:\n- namespace: \"acs_ecs_dashboard\"\n  metric_info:\n  - metric_list:\n    - \"cpu_total\"\n    - \"diskusage_utilization\"\n    - \"memory_usedutilization\"\n    - \"net_tcpconnection\"\n- namespace: \"acs_ecs_dashboard1\"\n  metric_info:\n  - metric_list:\n    - \"CPUUtilization\"\n    - \"DiskReadBPS\"\n    - \"InternetOut\"\n    - \"IntranetOut\"\n    - \"cpu_idle\"\n    - \"cpu_system\"\n    - \"cpu_total\"\n    - \"diskusage_utilization\"\n- namespace: \"acs_rds_dashboard\"\n  metric_info:\n  - metric_list:\n    - \"MySQL_QPS\"\n    - \"MySQL_TPS\"\n",
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

var AlicloudCloudMonitorServiceHybridMonitorFcTaskMap0 = map[string]string{
	"namespace":                 CHECKSET,
	"hybrid_monitor_fc_task_id": CHECKSET,
}

func AlicloudCloudMonitorServiceHybridMonitorFcTaskBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_account" "this" {}
resource "alicloud_cms_namespace" "default" {
	description = var.name
	namespace = "tf-testacc-cloudmonitorservicenamespace"
	specification = "cms.s1.large"
}
`, name)
}

func TestAccAlicloudCloudMonitorServiceHybridMonitorFcTask_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_hybrid_monitor_fc_task.default"
	checkoutSupportedRegions(t, true, connectivity.CloudMonitorServiceSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceHybridMonitorFcTaskMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsHybridMonitorFcTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudmonitorservicehybridmonitorfctask%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceHybridMonitorFcTaskBasicDependence0)
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
					"namespace":   "${alicloud_cms_namespace.default.id}",
					"yarm_config": "---\\nproducts:\\n- namespace: \\\"acs_ecs_dashboard1\\\"\\n  metric_info:\\n  - metric_list:\\n    - \\\"CPUUtilization\\\"\\n    - \\\"DiskReadBPS\\\"\\n    - \\\"InternetOut\\\"\\n    - \\\"IntranetOut\\\"\\n    - \\\"cpu_idle\\\"\\n    - \\\"cpu_system\\\"\\n    - \\\"cpu_total\\\"\\n    - \\\"diskusage_utilization\\\"\\n- namespace: \\\"acs_rds_dashboard\\\"\\n  metric_info:\\n  - metric_list:\\n    - \\\"MySQL_QPS\\\"\\n    - \\\"MySQL_TPS\\\"\\n- namespace: \\\"acs_ecs_dashboard\\\"\\n  metric_info:\\n  - metric_list:\\n    - \\\"cpu_total\\\"\\n    - \\\"diskusage_utilization\\\"\\n    - \\\"memory_usedutilization\\\"\\n    - \\\"net_tcpconnection\\\"\\n",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace":      CHECKSET,
						"target_user_id": CHECKSET,
						"yarm_config":    "---\nproducts:\n- namespace: \"acs_ecs_dashboard\"\n  metric_info:\n  - metric_list:\n    - \"cpu_total\"\n    - \"diskusage_utilization\"\n    - \"memory_usedutilization\"\n    - \"net_tcpconnection\"\n- namespace: \"acs_ecs_dashboard1\"\n  metric_info:\n  - metric_list:\n    - \"CPUUtilization\"\n    - \"DiskReadBPS\"\n    - \"InternetOut\"\n    - \"IntranetOut\"\n    - \"cpu_idle\"\n    - \"cpu_system\"\n    - \"cpu_total\"\n    - \"diskusage_utilization\"\n- namespace: \"acs_rds_dashboard\"\n  metric_info:\n  - metric_list:\n    - \"MySQL_QPS\"\n    - \"MySQL_TPS\"\n",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"yarm_config": "---\\nproducts:\\n- namespace: \\\"acs_ecs_dashboard2\\\"\\n  metric_info:\\n  - metric_list:\\n    - \\\"CPUUtilization\\\"\\n    - \\\"DiskReadBPS\\\"\\n    - \\\"InternetOut\\\"\\n    - \\\"IntranetOut\\\"\\n    - \\\"cpu_idle\\\"\\n    - \\\"cpu_system\\\"\\n    - \\\"cpu_total\\\"\\n    - \\\"diskusage_utilization\\\"\\n- namespace: \\\"acs_rds_dashboard\\\"\\n  metric_info:\\n  - metric_list:\\n    - \\\"MySQL_QPS\\\"\\n    - \\\"MySQL_TPS\\\"\\n- namespace: \\\"acs_ecs_dashboard\\\"\\n  metric_info:\\n  - metric_list:\\n    - \\\"cpu_total\\\"\\n    - \\\"diskusage_utilization\\\"\\n    - \\\"memory_usedutilization\\\"\\n    - \\\"net_tcpconnection\\\"\\n",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"yarm_config": "---\nproducts:\n- namespace: \"acs_ecs_dashboard\"\n  metric_info:\n  - metric_list:\n    - \"cpu_total\"\n    - \"diskusage_utilization\"\n    - \"memory_usedutilization\"\n    - \"net_tcpconnection\"\n- namespace: \"acs_ecs_dashboard2\"\n  metric_info:\n  - metric_list:\n    - \"CPUUtilization\"\n    - \"DiskReadBPS\"\n    - \"InternetOut\"\n    - \"IntranetOut\"\n    - \"cpu_idle\"\n    - \"cpu_system\"\n    - \"cpu_total\"\n    - \"diskusage_utilization\"\n- namespace: \"acs_rds_dashboard\"\n  metric_info:\n  - metric_list:\n    - \"MySQL_QPS\"\n    - \"MySQL_TPS\"\n",
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

func TestUnitAccAlicloudCmsHybridMonitorFcTask(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_cms_hybrid_monitor_fc_task"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cms_hybrid_monitor_fc_task"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"yarm_config":    "CreateCmsHybridMonitorFcTaskValue",
		"namespace":      "CreateCmsHybridMonitorFcTaskValue",
		"target_user_id": "CreateCmsHybridMonitorFcTaskValue",
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
				"CollectInterval":       60,
				"Description":           "CreateCmsHybridMonitorFcTaskValue",
				"TaskId":                "CmsHybridMonitorFcTaskId",
				"CollectTargetPath":     "CreateCmsHybridMonitorFcTaskValue",
				"Namespace":             "CreateCmsHybridMonitorFcTaskValue",
				"GroupId":               0,
				"MatchExpressRelation":  "CreateCmsHybridMonitorFcTaskValue",
				"CollectTargetEndpoint": "CreateCmsHybridMonitorFcTaskValue",
				"TaskName":              "CreateCmsHybridMonitorFcTaskValue",
				"UploadRegion":          "CreateCmsHybridMonitorFcTaskValue",
				"NetworkType":           "CreateCmsHybridMonitorFcTaskValue",
				"TaskType":              "CreateCmsHybridMonitorFcTaskValue",
				"CollectTargetType":     "CreateCmsHybridMonitorFcTaskValue",
				"TargetUserId":          "CreateCmsHybridMonitorFcTaskValue",
				"YARMConfig":            "CreateCmsHybridMonitorFcTaskValue",
			},
		},
		"Success": true,
	}
	CreateMockResponse := map[string]interface{}{
		"TaskId": "CmsHybridMonitorFcTaskId",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cms_hybrid_monitor_fc_task", errorCode))
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
	err = resourceAlicloudCmsHybridMonitorFcTaskCreate(dInit, rawClient)
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
		err := resourceAlicloudCmsHybridMonitorFcTaskCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cms_hybrid_monitor_fc_task"].Schema).Data(dInit.State(), nil)
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
	err = resourceAlicloudCmsHybridMonitorFcTaskUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"yarm_config": "UpdateCmsHybridMonitorFcTaskValue",
	}
	diff, err := newInstanceDiff("alicloud_cms_hybrid_monitor_fc_task", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cms_hybrid_monitor_fc_task"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"TaskList": []interface{}{
			map[string]interface{}{
				"YARMConfig": "UpdateCmsHybridMonitorFcTaskValue",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateHybridMonitorTask" {
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
		err := resourceAlicloudCmsHybridMonitorFcTaskUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cms_hybrid_monitor_fc_task"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_cms_hybrid_monitor_fc_task", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cms_hybrid_monitor_fc_task"].Schema).Data(dInit.State(), diff)
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
		err := resourceAlicloudCmsHybridMonitorFcTaskRead(dExisted, rawClient)
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
	err = resourceAlicloudCmsHybridMonitorFcTaskDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_cms_hybrid_monitor_fc_task", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cms_hybrid_monitor_fc_task"].Schema).Data(dInit.State(), diff)
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
		err := resourceAlicloudCmsHybridMonitorFcTaskDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
