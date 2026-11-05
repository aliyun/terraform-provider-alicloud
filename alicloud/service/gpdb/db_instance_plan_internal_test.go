package gpdb

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// lintignore: R001
func TestUnitAccAliCloudGpdbDbInstancePlan(t *testing.T) {
	p := unitTestProviderResources()
	dInit, _ := schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"db_instance_plan_name": "CreateGpdbDbInstancePlanValue",
		"plan_desc":             "CreateGpdbDbInstancePlanValue",
		"plan_type":             "CreateGpdbDbInstancePlanValue",
		"plan_schedule_type":    "CreateGpdbDbInstancePlanValue",
		"plan_start_date":       "CreateGpdbDbInstancePlanValue",
		"plan_end_date":         "CreateGpdbDbInstancePlanValue",
		"plan_config": []interface{}{
			map[string]interface{}{
				"resume": []interface{}{
					map[string]interface{}{
						"plan_cron_time": "CreateGpdbDbInstancePlanValue",
						"execute_time":   "CreateGpdbDbInstancePlanValue",
					},
				},
				"pause": []interface{}{
					map[string]interface{}{
						"plan_cron_time": "CreateGpdbDbInstancePlanValue",
						"execute_time":   "CreateGpdbDbInstancePlanValue",
					},
				},
				"scale_in": []interface{}{
					map[string]interface{}{
						"plan_cron_time":   "CreateGpdbDbInstancePlanValue",
						"execute_time":     "CreateGpdbDbInstancePlanValue",
						"segment_node_num": "2",
					},
				},
				"scale_out": []interface{}{
					map[string]interface{}{
						"plan_cron_time":   "CreateGpdbDbInstancePlanValue",
						"execute_time":     "CreateGpdbDbInstancePlanValue",
						"segment_node_num": "2",
					},
				},
			},
		},
		"db_instance_id": "CreateGpdbDbInstancePlanValue",
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
		"Items": map[string]interface{}{
			"PlanList": []interface{}{
				map[string]interface{}{
					"DBInstanceId":     "CreateGpdbDbInstancePlanValue",
					"PlanId":           "CreateGpdbDbInstancePlanValue",
					"PlanScheduleType": "CreateGpdbDbInstancePlanValue",
					"PlanType":         "CreateGpdbDbInstancePlanValue",
					"PlanDesc":         "CreateGpdbDbInstancePlanValue",
					"PlanName":         "CreateGpdbDbInstancePlanValue",
					"PlanStartDate":    "CreateGpdbDbInstancePlanValue",
					"PlanEndDate":      "CreateGpdbDbInstancePlanValue",
					"PlanConfig":       "{\"resume\":{\"planCronTime\":\"CreateGpdbDbInstancePlanValue\",\"executeTime\":\"CreateGpdbDbInstancePlanValue\"},\"pause\":{\"planCronTime\":\"CreateGpdbDbInstancePlanValue\",\"executeTime\":\"CreateGpdbDbInstancePlanValue\"},\"scaleOut\":{\"planCronTime\":\"CreateGpdbDbInstancePlanValue\",\"executeTime\":\"CreateGpdbDbInstancePlanValue\",\"segmentNodeNum\":\"2\"},\"scaleIn\":{\"planCronTime\":\"CreateGpdbDbInstancePlanValue\",\"executeTime\":\"CreateGpdbDbInstancePlanValue\",\"segmentNodeNum\":\"2\"}}",
					"PlanStatus":       "active",
				},
			},
		},
		"Status": "success",
	}
	CreateMockResponse := map[string]interface{}{
		"PlanId": "CreateGpdbDbInstancePlanValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String(errorCode),
			Data:       tea.String(errorCode),
			Message:    tea.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_gpdb_db_instance_plan", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGpdbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceDbInstancePlanCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateDBInstancePlan" {
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
		err := resourceDbInstancePlanCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGpdbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceDbInstancePlanUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"db_instance_plan_name": "UpdateGpdbDbInstancePlanValue",
		"plan_desc":             "UpdateGpdbDbInstancePlanValue",
		"plan_start_date":       "UpdateGpdbDbInstancePlanValue",
		"plan_end_date":         "UpdateGpdbDbInstancePlanValue",
		"status":                "cancel",
		"plan_config": []interface{}{
			map[string]interface{}{
				"resume": []interface{}{
					map[string]interface{}{
						"plan_cron_time": "UpdateGpdbDbInstancePlanValue",
						"execute_time":   "UpdateGpdbDbInstancePlanValue",
					},
				},
				"pause": []interface{}{
					map[string]interface{}{
						"plan_cron_time": "UpdateGpdbDbInstancePlanValue",
						"execute_time":   "UpdateGpdbDbInstancePlanValue",
					},
				},
				"scale_in": []interface{}{
					map[string]interface{}{
						"plan_cron_time":   "UpdateGpdbDbInstancePlanValue",
						"execute_time":     "UpdateGpdbDbInstancePlanValue",
						"segment_node_num": "2",
					},
				},
				"scale_out": []interface{}{
					map[string]interface{}{
						"plan_cron_time":   "UpdateGpdbDbInstancePlanValue",
						"execute_time":     "UpdateGpdbDbInstancePlanValue",
						"segment_node_num": "2",
					},
				},
			},
		},
	}
	diff, err := newInstanceDiff("alicloud_gpdb_db_instance_plan", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"Items": map[string]interface{}{
			"PlanList": []interface{}{
				map[string]interface{}{
					"PlanDesc":      "UpdateGpdbDbInstancePlanValue",
					"PlanName":      "UpdateGpdbDbInstancePlanValue",
					"PlanStartDate": "UpdateGpdbDbInstancePlanValue",
					"PlanEndDate":   "UpdateGpdbDbInstancePlanValue",
					"PlanConfig":    "{\"resume\":{\"planCronTime\":\"UpdateGpdbDbInstancePlanValue\",\"executeTime\":\"UpdateGpdbDbInstancePlanValue\"},\"pause\":{\"planCronTime\":\"UpdateGpdbDbInstancePlanValue\",\"executeTime\":\"UpdateGpdbDbInstancePlanValue\"},\"scaleOut\":{\"planCronTime\":\"UpdateGpdbDbInstancePlanValue\",\"executeTime\":\"UpdateGpdbDbInstancePlanValue\",\"segmentNodeNum\":\"2\"},\"scaleIn\":{\"planCronTime\":\"UpdateGpdbDbInstancePlanValue\",\"executeTime\":\"UpdateGpdbDbInstancePlanValue\",\"segmentNodeNum\":\"2\"}}",
					"PlanStatus":    "cancel",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateDBInstancePlan" {
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
			if *action == "SetDBInstancePlanStatus" {
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
		err := resourceDbInstancePlanUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_gpdb_db_instance_plan", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDBInstancePlans" {
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
		err := resourceDbInstancePlanRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGpdbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceDbInstancePlanDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_gpdb_db_instance_plan", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteDBInstancePlan" {
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
		err := resourceDbInstancePlanDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
