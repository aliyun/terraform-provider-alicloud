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
func TestUnitAliCloudGPDBDBInstance(t *testing.T) {
	p := unitTestProviderResources()
	dInit, _ := schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"db_instance_category":  "CreateDBInstanceValue",
		"db_instance_class":     "CreateDBInstanceValue",
		"db_instance_mode":      "CreateDBInstanceValue",
		"description":           "CreateDBInstanceValue",
		"engine":                "CreateDBInstanceValue",
		"engine_version":        "CreateDBInstanceValue",
		"effective_time":        "Immediate",
		"zone_id":               "CreateDBInstanceValue",
		"instance_network_type": "CreateDBInstanceValue",
		"instance_spec":         "CreateDBInstanceValue",
		"payment_type":          "PayAsYouGo",
		"seg_storage_type":      "CreateDBInstanceValue",
		"seg_node_num":          4,
		"storage_size":          50,
		"vpc_id":                "CreateDBInstanceValue",
		"vswitch_id":            "CreateDBInstanceValue",
		"resource_group_id":     "CreateDBInstanceValue",
		"period":                "Month",
		"instance_group_count":  2,
		"create_sample_data":    false,
		"used_time":             "1",
		"ip_whitelist": []map[string]interface{}{
			{
				"security_ip_list": "127.0.0.1",
			},
		},
		"tags": map[string]string{
			"Created": "TF",
			"For":     "acceptance test",
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
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// DescribeDBInstanceAttribute
		"Items": map[string]interface{}{
			// DescribeDataShareInstances returns no data-share membership for
			// the mocked instance; an empty list keeps the Read path error-free.
			"DBInstance": []interface{}{},
			"DBInstanceAttribute": []interface{}{
				map[string]interface{}{
					"DBInstanceId":          "CreateDBInstanceValue",
					"DBInstanceCategory":    "CreateDBInstanceValue",
					"DBInstanceMode":        "CreateDBInstanceValue",
					"DBInstanceDescription": "CreateDBInstanceValue",
					"Engine":                "CreateDBInstanceValue",
					"EngineVersion":         "CreateDBInstanceValue",
					"InstanceNetworkType":   "CreateDBInstanceValue",
					"InstanceSpec":          "CreateDBInstanceValue",
					"MaintainEndTime":       "CreateDBInstanceValue",
					"MaintainStartTime":     "CreateDBInstanceValue",
					"MasterNodeNum":         1,
					"ResourceGroupId":       "CreateDBInstanceValue",
					"SegmentCounts":         0,
					"PayType":               "Postpaid",
					"SegNodeNum":            4,
					"DBInstanceStatus":      "Running",
					"StorageSize":           50,
					"StorageType":           "CreateDBInstanceValue",
					"VSwitchId":             "CreateDBInstanceValue",
					"VpcId":                 "CreateDBInstanceValue",
					"ZoneId":                "CreateDBInstanceValue",
					"TagResources": map[string]interface{}{
						"TagResource": []interface{}{
							map[string]interface{}{
								"TagKey":   "Created",
								"TagValue": "TF",
							},
							map[string]interface{}{
								"TagKey":   "For",
								"TagValue": "acceptance test123",
							},
						},
					},
				},
			},
			"DBInstanceIPArray": []interface{}{
				map[string]interface{}{
					"DBInstanceIPArrayAttribute": "",
					"DBInstanceIPArrayName":      "",
					"SecurityIPList":             "127.0.0.1",
				},
			},
		},
		"DBInstanceId": "CreateDBInstanceValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateDBInstance
		"DBInstanceId": "CreateDBInstanceValue",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_gpdb_instance", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	// The GPDB instance resource sends every request through AliyunClient.RpcPost,
	// so the loadEndpoint failure branch is covered by patching that real method
	// (there is no NewGpdbClient constructor on AliyunClient).
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, _, _, _ string, _, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:    tea.String("loadEndpoint error"),
			Data:    tea.String("loadEndpoint error"),
			Message: tea.String("loadEndpoint error"),
		}
	})
	err = resourceDbInstanceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeDBInstanceAttribute Response
		"DBInstanceId": "CreateDBInstanceValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateDBInstance" {
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
		err := resourceDbInstanceCreate(dInit, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, _, _, _ string, _, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:    tea.String("loadEndpoint error"),
			Data:    tea.String("loadEndpoint error"),
			Message: tea.String("loadEndpoint error"),
		}
	})
	err = resourceDbInstanceUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyDBInstanceDescription
	attributesDiff := map[string]interface{}{
		"description": "ModifyDBInstanceDescriptionValue",
	}
	diff, err := newInstanceDiff("alicloud_gpdb_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDBInstanceAttribute Response
		"Items": map[string]interface{}{
			"DBInstanceAttribute": []interface{}{
				map[string]interface{}{
					"DBInstanceDescription": "ModifyDBInstanceDescriptionValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDBInstanceDescription" {
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
		err := resourceDbInstanceUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ModifyDBInstanceMaintainTime
	attributesDiff = map[string]interface{}{
		"maintain_end_time":   "ModifyDBInstanceMaintainTimeValue",
		"maintain_start_time": "ModifyDBInstanceMaintainTimeValue",
	}
	diff, err = newInstanceDiff("alicloud_gpdb_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDBInstanceAttribute Response
		"Items": map[string]interface{}{
			"DBInstanceAttribute": []interface{}{
				map[string]interface{}{
					"MaintainEndTime":   "ModifyDBInstanceMaintainTimeValue",
					"MaintainStartTime": "ModifyDBInstanceMaintainTimeValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDBInstanceMaintainTime" {
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
		err := resourceDbInstanceUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ModifySecurityIps
	attributesDiff = map[string]interface{}{
		"ip_whitelist": []map[string]interface{}{
			{
				"security_ip_list": "1.1.1.1",
			},
		},
	}
	diff, err = newInstanceDiff("alicloud_gpdb_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDBInstanceAttribute Response
		"Items": map[string]interface{}{
			"DBInstanceIPArray": []interface{}{
				map[string]interface{}{
					"DBInstanceIPArrayAttribute": "",
					"DBInstanceIPArrayName":      "",
					"SecurityIPList":             "1.1.1.1",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifySecurityIps" {
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
		err := resourceDbInstanceUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// UpgradeDBInstance
	attributesDiff = map[string]interface{}{
		"seg_node_num": 8,
	}
	diff, err = newInstanceDiff("alicloud_gpdb_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDBInstanceAttribute Response
		"Items": map[string]interface{}{
			"DBInstanceAttribute": []interface{}{
				map[string]interface{}{
					"SegNodeNum": 8,
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpgradeDBInstance" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						assert.Equal(t, "Immediate", body["EffectiveTime"])
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceDbInstanceUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// UpgradeDBInstance
	attributesDiff = map[string]interface{}{
		"storage_size": 100,
	}
	diff, err = newInstanceDiff("alicloud_gpdb_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDBInstanceAttribute Response
		"Items": map[string]interface{}{
			"DBInstanceAttribute": []interface{}{
				map[string]interface{}{
					"StorageSize": 100,
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpgradeDBInstance" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						assert.Equal(t, "Immediate", body["EffectiveTime"])
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceDbInstanceUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// UpgradeDBVersion + ModifySQLCollectorPolicy
	// The last state-compared update block carries all three write-only changes at once:
	// effective_time (MaintainTime, covering the deferred-upgrade wait skip) and
	// sql_collector_status (Enabled) have no API readback, so they must be introduced in
	// the final diff block — earlier blocks re-derive their state from the mocked
	// DescribeDBInstanceAttribute response, which cannot report these attributes.
	attributesDiff = map[string]interface{}{
		"minor_version":        "6.6.2.21-202606241529",
		"effective_time":       "MaintainTime",
		"sql_collector_status": "Enabled",
	}
	diff, err = newInstanceDiff("alicloud_gpdb_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDBInstanceAttribute Response
		// DBInstanceStatus=IDLE locks the regression: Serverless instances settle in
		// IDLE, so the post-ModifySQLCollectorPolicy update wait must accept it instead
		// of polling for Running only until the update timeout.
		"Items": map[string]interface{}{
			"DBInstanceAttribute": []interface{}{
				map[string]interface{}{
					"MinorVersion":     "6.6.2.21-202606241529",
					"DBInstanceStatus": "IDLE",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling.User", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpgradeDBVersion" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						assert.Equal(t, "MaintainTime", body["EffectiveTime"])
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			if *action == "ModifySQLCollectorPolicy" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						assert.Equal(t, "Enable", body["SQLCollectorStatus"])
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceDbInstanceUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Updating the maintenance window and the minor version in one update must
	// commit ModifyDBInstanceMaintainTime before UpgradeDBVersion: an upgrade
	// scheduled with EffectiveTime = MaintainTime runs inside the window as it is
	// configured at submission time. The mocked describe response reports IDLE,
	// which also locks the IDLE acceptance of the minor-version wait for
	// serverless instances.
	attributesDiff = map[string]interface{}{
		"maintain_end_time":   "06:00Z",
		"maintain_start_time": "02:00Z",
		"minor_version":       "6.6.2.21-202606241530",
	}
	diff, err = newInstanceDiff("alicloud_gpdb_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDBInstanceAttribute Response
		"Items": map[string]interface{}{
			"DBInstanceAttribute": []interface{}{
				map[string]interface{}{
					"MinorVersion":     "6.6.2.21-202606241530",
					"DBInstanceStatus": "IDLE",
				},
			},
		},
	}
	var callLog []string
	lastBody := make(map[string]map[string]interface{})
	gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		switch *action {
		case "ModifyDBInstanceMaintainTime", "UpgradeDBVersion", "ModifySQLCollectorPolicy":
			if len(callLog) == 0 || callLog[len(callLog)-1] != *action {
				callLog = append(callLog, *action)
			}
			lastBody[*action] = body
			return ReadMockResponse, nil
		case "DescribeDBInstanceAttribute":
			return successResponseMock(ReadMockResponseDiff)
		default:
			return ReadMockResponse, nil
		}
	})
	err = resourceDbInstanceUpdate(dExisted, rawClient)
	assert.Nil(t, err)
	assert.Equal(t, []string{"ModifyDBInstanceMaintainTime", "UpgradeDBVersion"}, callLog)
	assert.Equal(t, "02:00Z", lastBody["ModifyDBInstanceMaintainTime"]["StartTime"])
	assert.Equal(t, "06:00Z", lastBody["ModifyDBInstanceMaintainTime"]["EndTime"])
	assert.Equal(t, "6.6.2.21-202606241530", lastBody["UpgradeDBVersion"]["MinorVersion"])
	assert.Equal(t, "Immediate", lastBody["UpgradeDBVersion"]["EffectiveTime"])

	// A maintenance-window update failure must abort the update before the
	// upgrade is submitted: an upgrade scheduled into the maintenance window
	// would otherwise run with a window the provider failed to apply.
	attributesDiff = map[string]interface{}{
		"maintain_end_time":   "07:00Z",
		"maintain_start_time": "03:00Z",
		"minor_version":       "6.6.2.21-202606241531",
	}
	diff, err = newInstanceDiff("alicloud_gpdb_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dInit.State(), diff)
	callLog = nil
	gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		switch *action {
		case "ModifyDBInstanceMaintainTime":
			callLog = append(callLog, *action)
			return failedResponseMock("NonRetryableError")
		case "UpgradeDBVersion", "ModifySQLCollectorPolicy":
			callLog = append(callLog, *action)
			return ReadMockResponse, nil
		default:
			return ReadMockResponse, nil
		}
	})
	err = resourceDbInstanceUpdate(dExisted, rawClient)
	assert.NotNil(t, err)
	assert.Equal(t, []string{"ModifyDBInstanceMaintainTime"}, callLog)

	// A sql_collector_status change is not persisted when a later block of the
	// same update fails: without a partial-state record the state keeps the
	// previous collector status ("Enabled" here), the next plan re-diffs the
	// configured value, and the re-apply re-sends the same idempotent
	// ModifySQLCollectorPolicy call (self-heal). The mocked ModifyParameters
	// failure reproduces a follow-up change failing after the collector switch
	// succeeded.
	attributesDiff = map[string]interface{}{
		"sql_collector_status": "Disabled",
		"parameters": []map[string]interface{}{
			{
				"name":  "optimizer",
				"value": "off",
			},
		},
	}
	dCollector, _ := schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dInit.State(), nil)
	err = dCollector.Set("sql_collector_status", "Enabled")
	assert.Nil(t, err)
	diff, err = newInstanceDiff("alicloud_gpdb_instance", attributes, attributesDiff, dCollector.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_instance"].Schema).Data(dCollector.State(), diff)
	callLog = nil
	gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		switch *action {
		case "ModifySQLCollectorPolicy":
			callLog = append(callLog, *action)
			assert.Equal(t, "Disabled", body["SQLCollectorStatus"])
			return ReadMockResponse, nil
		case "ModifyParameters":
			callLog = append(callLog, *action)
			return failedResponseMock("NonRetryableError")
		default:
			return ReadMockResponse, nil
		}
	})
	err = resourceDbInstanceUpdate(dExisted, rawClient)
	assert.NotNil(t, err)
	assert.Equal(t, []string{"ModifySQLCollectorPolicy", "ModifyParameters"}, callLog)
	assert.Equal(t, "Enabled", dExisted.State().Attributes["sql_collector_status"])
	assert.NotContains(t, dExisted.State().Attributes, "parameters.#")

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDBInstanceAttribute" {
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
		err := resourceDbInstanceRead(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, _, _, _ string, _, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:    tea.String("loadEndpoint error"),
			Data:    tea.String("loadEndpoint error"),
			Message: tea.String("loadEndpoint error"),
		}
	})
	err = resourceDbInstanceDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "AclNotExist", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteDBInstance" {
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
		err := resourceDbInstanceDelete(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}

// The whitelist read path maps remote groups onto the reference list: named
// groups by group name, anonymous groups by their sorted IP list. These cases
// lock that matching behavior without any cloud calls.
func TestUnitAliCloudGPDBDBInstanceIPWhitelistNormalization(t *testing.T) {
	reference := []interface{}{
		map[string]interface{}{"ip_group_name": "group1", "security_ip_list": "11.0.0.1"},
		map[string]interface{}{"ip_group_name": "default", "security_ip_list": "10.0.0.1"},
	}
	remote := []map[string]interface{}{
		{"ip_group_name": "default", "security_ip_list": "10.0.0.1"},
		{"ip_group_name": "group1", "security_ip_list": "11.0.0.1"},
	}
	normalized := normalizeIPWhitelistOrder(reference, remote)
	assert.Equal(t, "group1", normalized[0]["ip_group_name"])
	assert.Equal(t, "default", normalized[1]["ip_group_name"])

	// Anonymous groups match by their sorted IP list, so member order inside a
	// group is insignificant while the group keeps the remote member intact.
	// The read path always writes ip_group_name, so anonymous fixtures carry
	// an explicit empty name; a missing key would match via the named branch.
	reference = []interface{}{
		map[string]interface{}{"ip_group_name": "", "security_ip_list": "10.0.0.2,10.0.0.1"},
	}
	remote = []map[string]interface{}{
		{"ip_group_name": "", "security_ip_list": "10.0.0.1,10.0.0.2"},
	}
	normalized = normalizeIPWhitelistOrder(reference, remote)
	assert.Equal(t, 1, len(normalized))
	assert.Equal(t, "10.0.0.1,10.0.0.2", normalized[0]["security_ip_list"])
	assert.True(t, ipWhitelistMemberMatches(reference[0].(map[string]interface{}), remote[0]))
	assert.False(t, ipWhitelistMemberMatches(
		map[string]interface{}{"ip_group_name": "", "security_ip_list": "10.0.0.3"},
		map[string]interface{}{"ip_group_name": "", "security_ip_list": "10.0.0.1,10.0.0.2"}))

	// Named groups match by group name alone.
	assert.True(t, ipWhitelistMemberMatches(
		map[string]interface{}{"ip_group_name": "group1", "security_ip_list": "10.0.0.1"},
		map[string]interface{}{"ip_group_name": "group1", "security_ip_list": "11.0.0.9"}))
	assert.False(t, ipWhitelistMemberMatches(
		map[string]interface{}{"ip_group_name": "group1", "security_ip_list": "10.0.0.1"},
		map[string]interface{}{"ip_group_name": "group2", "security_ip_list": "10.0.0.1"}))

	// Remote members without a reference position are appended and reference
	// members the API no longer returns are dropped.
	reference = []interface{}{
		map[string]interface{}{"ip_group_name": "group1", "security_ip_list": "11.0.0.1"},
		map[string]interface{}{"ip_group_name": "", "security_ip_list": "10.0.0.3"},
	}
	remote = []map[string]interface{}{
		{"ip_group_name": "", "security_ip_list": "10.0.0.9"},
		{"ip_group_name": "group1", "security_ip_list": "11.0.0.1"},
	}
	normalized = normalizeIPWhitelistOrder(reference, remote)
	assert.Equal(t, 2, len(normalized))
	assert.Equal(t, "group1", normalized[0]["ip_group_name"])
	assert.Equal(t, "10.0.0.9", normalized[1]["security_ip_list"])

	// An empty reference (create, import) keeps the API order.
	remote = []map[string]interface{}{
		{"ip_group_name": "default", "security_ip_list": "10.0.0.1"},
		{"ip_group_name": "group1", "security_ip_list": "11.0.0.1"},
	}
	assert.Equal(t, remote, normalizeIPWhitelistOrder([]interface{}{}, remote))
}
