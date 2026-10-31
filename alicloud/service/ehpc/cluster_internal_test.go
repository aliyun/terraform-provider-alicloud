// Package ehpc. Internal test package (package ehpc) so the TestUnit
// gomonkey tests can call the unexported resourceAlicloudEhpc* CRUD
// functions directly, matching the migrated-subpackage convention
// (see alicloud/service/ecs/disk_unit_test.go).
package ehpc

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
func TestUnitAlicloudEhpcCluster(t *testing.T) {
	p := unitTestProviderResources()
	dInit, _ := schema.InternalMap(p["alicloud_ehpc_cluster"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ehpc_cluster"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"description":              "CreateEhpcClusterValue",
		"cluster_name":             "CreateEhpcClusterValue",
		"volume_protocol":          "CreateEhpcClusterValue",
		"volume_id":                "CreateEhpcClusterValue",
		"volume_mountpoint":        "CreateEhpcClusterValue",
		"deploy_mode":              "CreateEhpcClusterValue",
		"image_id":                 "CreateEhpcClusterValue",
		"image_owner_alias":        "CreateEhpcClusterValue",
		"cluster_version":          "CreateEhpcClusterValue",
		"volume_type":              "CreateEhpcClusterValue",
		"remote_directory":         "CreateEhpcClusterValue",
		"scheduler_type":           "CreateEhpcClusterValue",
		"account_type":             "CreateEhpcClusterValue",
		"ha_enable":                false,
		"os_tag":                   "CreateEhpcClusterValue",
		"manager_count":            2,
		"manager_instance_type":    "CreateEhpcClusterValue",
		"compute_count":            1,
		"compute_instance_type":    "CreateEhpcClusterValue",
		"login_count":              1,
		"login_instance_type":      "CreateEhpcClusterValue",
		"password":                 "CreateEhpcClusterValue",
		"vswitch_id":               "CreateEhpcClusterValue",
		"vpc_id":                   "CreateEhpcClusterValue",
		"zone_id":                  "CreateEhpcClusterValue",
		"auto_renew":               true,
		"auto_renew_period":        1,
		"compute_enable_ht":        true,
		"compute_spot_price_limit": "CreateEhpcClusterValue",
		"compute_spot_strategy":    "CreateEhpcClusterValue",
		"domain":                   "CreateEhpcClusterValue",
		"ecs_charge_type":          "CreateEhpcClusterValue",
		"ehpc_version":             "CreateEhpcClusterValue",
		"input_file_url":           "CreateEhpcClusterValue",
		"is_compute_ess":           true,
		"job_queue":                "CreateEhpcClusterValue",
		"period":                   1,
		"period_unit":              "CreateEhpcClusterValue",
		"plugin":                   "CreateEhpcClusterValue",
		"key_pair_name":            "CreateEhpcClusterValue",
		"ram_role_name":            "CreateEhpcClusterValue",
		"remote_vis_enable":        true,
		"resource_group_id":        "CreateEhpcClusterValue",
		"scc_cluster_id":           "CreateEhpcClusterValue",
		"security_group_id":        "CreateEhpcClusterValue",
		"security_group_name":      "CreateEhpcClusterValue",
		"system_disk_level":        "CreateEhpcClusterValue",
		"system_disk_size":         40,
		"without_agent":            true,
		"without_elastic_ip":       true,
		"system_disk_type":         "CreateEhpcClusterValue",
		"volume_mount_option":      "CreateEhpcClusterValue",
		"client_version":           "CreateEhpcClusterValue",
		"release_instance":         true,
		"additional_volumes": []map[string]interface{}{
			{
				"job_queue":        "CreateEhpcClusterValue",
				"local_directory":  "CreateEhpcClusterValue",
				"location":         "CreateEhpcClusterValue",
				"remote_directory": "CreateEhpcClusterValue",
				"roles": []map[string]interface{}{
					{
						"name": "CreateEhpcClusterValue",
					},
				},
				"volume_id":           "CreateEhpcClusterValue",
				"volume_mount_option": "CreateEhpcClusterValue",
				"volume_mountpoint":   "CreateEhpcClusterValue",
				"volume_protocol":     "CreateEhpcClusterValue",
				"volume_type":         "CreateEhpcClusterValue",
			},
		},
		"application": []map[string]interface{}{
			{
				"tag": "CreateEhpcClusterValue",
			},
		},
		"post_install_script": []map[string]interface{}{
			{
				"args": "CreateEhpcClusterValue",
				"url":  "CreateEhpcClusterValue",
			},
		},
		"ram_node_types": []string{"CreateEhpcClusterValue"},
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
		"ClusterInfo": map[string]interface{}{
			"Status":           "running",
			"VpcId":            "CreateEhpcClusterValue",
			"KeyPairName":      "CreateEhpcClusterValue",
			"EcsChargeType":    "CreateEhpcClusterValue",
			"SecurityGroupId":  "CreateEhpcClusterValue",
			"SccClusterId":     "CreateEhpcClusterValue",
			"CreateTime":       "CreateEhpcClusterValue",
			"AccountType":      "CreateEhpcClusterValue",
			"VolumeProtocol":   "CreateEhpcClusterValue",
			"Description":      "CreateEhpcClusterValue",
			"VolumeId":         "CreateEhpcClusterValue",
			"HaEnable":         false,
			"BaseOsTag":        "CreateEhpcClusterValue",
			"Name":             "CreateEhpcClusterValue",
			"ImageId":          "CreateEhpcClusterValue",
			"SchedulerType":    "CreateEhpcClusterValue",
			"DeployMode":       "CreateEhpcClusterValue",
			"ImageOwnerAlias":  "CreateEhpcClusterValue",
			"OsTag":            "CreateEhpcClusterValue",
			"VolumeMountpoint": "CreateEhpcClusterValue",
			"RemoteDirectory":  "CreateEhpcClusterValue",
			"RegionId":         "CreateEhpcClusterValue",
			"VSwitchId":        "CreateEhpcClusterValue",
			"ImageName":        "CreateEhpcClusterValue",
			"VolumeType":       "CreateEhpcClusterValue",
			"Location":         "CreateEhpcClusterValue",
			"Id":               "CreateEhpcClusterValue",
			"ClientVersion":    "CreateEhpcClusterValue",
			"Applications": map[string]interface{}{
				"ApplicationInfo": []interface{}{
					map[string]interface{}{
						"Tag":     "CreateEhpcClusterValue",
						"Name":    "CreateEhpcClusterValue",
						"Version": "CreateEhpcClusterValue",
					}},
			},
			"PostInstallScripts": map[string]interface{}{
				"PostInstallScriptInfo": []interface{}{
					map[string]interface{}{
						"Url":  "CreateEhpcClusterValue",
						"Args": "CreateEhpcClusterValue",
					},
				},
			},
			"EcsInfo": map[string]interface{}{
				"Manager": map[string]interface{}{
					"InstanceType": "CreateEhpcClusterValue",
					"Count":        2,
				},
				"Compute": map[string]interface{}{
					"InstanceType": "CreateEhpcClusterValue",
					"Count":        1,
				},
				"Login": map[string]interface{}{
					"InstanceType": "CreateEhpcClusterValue",
					"Count":        1,
				},
				"ProxyMgr": map[string]interface{}{
					"InstanceType": "CreateEhpcClusterValue",
					"Count":        1,
				},
			},
			"OnPremiseInfo": []interface{}{
				map[string]interface{}{
					"Type":     "CreateEhpcClusterValue",
					"HostName": "CreateEhpcClusterValue",
					"IP":       "CreateEhpcClusterValue",
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		"ClusterId": "CreateEhpcClusterValue",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_ehpc_cluster", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEhsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEhpcClusterCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateCluster" {
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
		err := resourceAlicloudEhpcClusterCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ehpc_cluster"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEhsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEhpcClusterUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"description":       "UpdateEhpcClusterValue",
		"cluster_name":      "UpdateEhpcClusterValue",
		"image_id":          "UpdateEhpcClusterValue",
		"image_owner_alias": "UpdateEhpcClusterValue",
	}
	diff, err := newInstanceDiff("alicloud_ehpc_cluster", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ehpc_cluster"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"ClusterInfo": map[string]interface{}{
			"Name":            "UpdateEhpcClusterValue",
			"Description":     "UpdateEhpcClusterValue",
			"ImageId":         "UpdateEhpcClusterValue",
			"ImageOwnerAlias": "UpdateEhpcClusterValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyClusterAttributes" {
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
		err := resourceAlicloudEhpcClusterUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ehpc_cluster"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_ehpc_cluster", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ehpc_cluster"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeCluster" {
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
		err := resourceAlicloudEhpcClusterRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEhsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEhpcClusterDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_ehpc_cluster", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ehpc_cluster"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "ClusterNotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteCluster" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					if errorCodes[retryIndex] == "ClusterNotFound" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			if *action == "DescribeCluster" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEhpcClusterDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
