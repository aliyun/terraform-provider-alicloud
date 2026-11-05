package dbs

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

var unitTestEndpoints sync.Map
var unitTestSignVersion sync.Map

func unitTestProviderResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"alicloud_dbs_backup_plan": ResourceDbsBackupPlan(),
	}
}

// unitTestSharedClientForRegion mirrors the root test helper so the unit test
// can be skipped when no cloud credentials are configured.
func unitTestSharedClientForRegion(region string) (interface{}, error) {
	var accessKey, secretKey, securityToken string
	if accessKey = os.Getenv("ALICLOUD_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}

	if secretKey = os.Getenv("ALICLOUD_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}

	securityToken = os.Getenv("ALICLOUD_SECURITY_TOKEN")

	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   &unitTestEndpoints,
		SignVersion: &unitTestSignVersion,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
	}
	accountId := os.Getenv("ALICLOUD_ACCOUNT_ID")
	if accountId == "" {
		accountId = os.Getenv("ALIBABA_CLOUD_ACCOUNT_ID")
	}
	if accountId != "" {
		conf.AccountId = accountId
	}
	accountType := os.Getenv("ALICLOUD_ACCOUNT_TYPE")
	if accountType == "" {
		accountType = os.Getenv("ALIBABA_CLOUD_ACCOUNT_TYPE")
	}
	if accountType != "" {
		conf.AccountType = accountType
	}
	credentialConfig := new(credentials.Config).SetType("access_key").SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
	if v := strings.TrimSpace(securityToken); v != "" {
		credentialConfig.SetType("sts").SetSecurityToken(v)
	}
	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		return nil, err
	}
	conf.Credential = credential

	client, err := conf.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// newInstanceDiff mirrors the root helper for the migrated resource.
// lintignore: R001
func newInstanceDiff(resourceName string, attributes, attributesDiff map[string]interface{}, state *terraform.InstanceState) (*terraform.InstanceDiff, error) {
	p := unitTestProviderResources()
	dOld, _ := schema.InternalMap(p[resourceName].Schema).Data(state, nil)
	dNew, _ := schema.InternalMap(p[resourceName].Schema).Data(state, nil)
	for key, value := range attributes {
		err := dOld.Set(key, value)
		if err != nil {
			return nil, errs.WrapErrorf(err, "[ERROR] the field %s setting error.", key)
		}
	}
	for key, value := range attributesDiff {
		attributes[key] = value
	}

	for key, value := range attributes {
		err := dNew.Set(key, value)
		if err != nil {
			return nil, errs.WrapErrorf(err, "[ERROR] the field %s setting error.", key)
		}
	}

	diff := terraform.NewInstanceDiff()
	objectKey := ""
	for _, key := range helper.MapSort(dNew.State().Attributes) {
		newValue := dNew.State().Attributes[key]
		if objectKey != "" && !strings.HasPrefix(key, objectKey) {
			objectKey = ""
		}
		if objectKey == "" {
			for _, suffix := range []string{"#", "%"} {
				if strings.HasSuffix(key, suffix) {
					objectKey = strings.TrimSuffix(key, suffix)
					break
				}
			}
		}
		oldValue, ok := dOld.State().Attributes[key]
		if ok && oldValue == newValue {
			continue
		}
		if oldValue == "" {
			for _, suffix := range []string{"#", "%"} {
				if strings.HasSuffix(key, suffix) {
					oldValue = "0"
				}
			}
		}
		diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: oldValue, New: newValue}

		if objectKey != "" {
			for removeKey, removeValue := range dOld.State().Attributes {
				if strings.HasPrefix(removeKey, objectKey) {
					if _, ok := dNew.State().Attributes[removeKey]; !ok {
						if len(strings.Split(removeKey, ".")) > 2 {
							delete(diff.Attributes, removeKey)
						} else {
							diff.Attributes[removeKey] = &terraform.ResourceAttrDiff{Old: removeValue, New: ""}
						}
					}
				}
			}
			objectKey = ""
		}
	}
	return diff, nil
}

// lintignore: R001
func TestUnitAccAlicloudDbsBackupPlan(t *testing.T) {
	p := unitTestProviderResources()
	dInit, _ := schema.InternalMap(p["alicloud_dbs_backup_plan"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_dbs_backup_plan"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"backup_plan_name":                     "CreateDbsBackupPlanValue",
		"payment_type":                         "CreateDbsBackupPlanValue",
		"instance_class":                       "CreateDbsBackupPlanValue",
		"backup_method":                        "CreateDbsBackupPlanValue",
		"database_type":                        "CreateDbsBackupPlanValue",
		"database_region":                      "CreateDbsBackupPlanValue",
		"storage_region":                       "CreateDbsBackupPlanValue",
		"instance_type":                        "CreateDbsBackupPlanValue",
		"source_endpoint_instance_type":        "CreateDbsBackupPlanValue",
		"resource_group_id":                    "CreateDbsBackupPlanValue",
		"source_endpoint_region":               "CreateDbsBackupPlanValue",
		"source_endpoint_instance_id":          "CreateDbsBackupPlanValue",
		"source_endpoint_user_name":            "CreateDbsBackupPlanValue",
		"source_endpoint_password":             "CreateDbsBackupPlanValue",
		"backup_objects":                       "CreateDbsBackupPlanValue",
		"backup_period":                        "CreateDbsBackupPlanValue",
		"backup_start_time":                    "CreateDbsBackupPlanValue",
		"backup_storage_type":                  "CreateDbsBackupPlanValue",
		"backup_gateway_id":                    "123",
		"backup_log_interval_seconds":          1000,
		"backup_rate_limit":                    "CreateDbsBackupPlanValue",
		"backup_speed_limit":                   "CreateDbsBackupPlanValue",
		"backup_strategy_type":                 "CreateDbsBackupPlanValue",
		"cross_aliyun_id":                      "CreateDbsBackupPlanValue",
		"cross_role_name":                      "CreateDbsBackupPlanValue",
		"duplication_archive_period":           100,
		"duplication_infrequent_access_period": 100,
		"enable_backup_log":                    true,
		"oss_bucket_name":                      "CreateDbsBackupPlanValue",
		"source_endpoint_oracle_sid":           "CreateDbsBackupPlanValue",
		"period":                               "CreateDbsBackupPlanValue",
		"source_endpoint_database_name":        "CreateDbsBackupPlanValue",
		"source_endpoint_ip":                   "CreateDbsBackupPlanValue",
		"source_endpoint_port":                 80,
		"source_endpoint_sid":                  "CreateDbsBackupPlanValue",
		"used_time":                            1,
		"backup_retention_period":              740,
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
	rawClient, err := unitTestSharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}

	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"Items": map[string]interface{}{
			"BackupPlanDetail": []interface{}{
				map[string]interface{}{
					"BackupMethod":                      "CreateDbsBackupPlanValue",
					"BackupPlanCreateTime":              1661331208000,
					"BackupPlanName":                    "CreateDbsBackupPlanValue",
					"BackupPeriod":                      "CreateDbsBackupPlanValue",
					"BackupObjects":                     "CreateDbsBackupPlanValue",
					"CrossAliyunId":                     "CreateDbsBackupPlanValue",
					"DatabaseType":                      "CreateDbsBackupPlanValue",
					"SourceEndpointInstanceID":          "CreateDbsBackupPlanValue",
					"InstanceClass":                     "CreateDbsBackupPlanValue",
					"OSSBucketRegion":                   "CreateDbsBackupPlanValue",
					"SourceEndpointRegion":              "CreateDbsBackupPlanValue",
					"OpenBackupSetAutoDownload":         false,
					"SourceEndpointIpPort":              "CreateDbsBackupPlanValue",
					"DuplicationArchivePeriod":          100,
					"BackupPlanStatus":                  "running",
					"OSSBucketName":                     "CreateDbsBackupPlanValue",
					"CrossRoleName":                     "CreateDbsBackupPlanValue",
					"BackupStartTime":                   "CreateDbsBackupPlanValue",
					"EnableBackupLog":                   true,
					"BackupPlanId":                      "CreateDbsBackupPlanValue",
					"BackupRetentionPeriod":             740,
					"SourceEndpointInstanceType":        "CreateDbsBackupPlanValue",
					"BackupGatewayId":                   123,
					"BackupStorageType":                 "CreateDbsBackupPlanValue",
					"SourceEndpointUserName":            "CreateDbsBackupPlanValue",
					"SourceEndpointDatabaseName":        "CreateDbsBackupPlanValue",
					"SourceEndpointOracleSID":           "CreateDbsBackupPlanValue",
					"DuplicationInfrequentAccessPeriod": 100,
				},
			},
		},
	}
	DescribeBackupPlanBillingMockResponse := map[string]interface{}{
		"Item": map[string]interface{}{
			"BuyChargeType": "CreateDbsBackupPlanValue",
		},
	}
	CreateMockResponse := map[string]interface{}{
		"BackupPlanId": "CreateDbsBackupPlanValue",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_dbs_backup_plan", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCbsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudDbsBackupPlanCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateAndStartBackupPlan" {
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
			if *action == "DescribeBackupPlanBilling" {
				return DescribeBackupPlanBillingMockResponse, nil
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudDbsBackupPlanCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_dbs_backup_plan"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update pause
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCbsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudDbsBackupPlanUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"status": "pause",
	}
	diff, err := newInstanceDiff("alicloud_dbs_backup_plan", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_dbs_backup_plan"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"Items": map[string]interface{}{
			"BackupPlanDetail": []interface{}{
				map[string]interface{}{
					"BackupPlanStatus": "pause",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "StopBackupPlan" {
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
			if *action == "DescribeBackupPlanBilling" {
				return DescribeBackupPlanBillingMockResponse, nil
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudDbsBackupPlanUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_dbs_backup_plan"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update running
	attributesDiff = map[string]interface{}{
		"status": "running",
	}
	diff, err = newInstanceDiff("alicloud_dbs_backup_plan", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_dbs_backup_plan"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"Items": map[string]interface{}{
			"BackupPlanDetail": []interface{}{
				map[string]interface{}{
					"BackupPlanStatus": "running",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "StartBackupPlan" {
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
			if *action == "DescribeBackupPlanBilling" {
				return DescribeBackupPlanBillingMockResponse, nil
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudDbsBackupPlanUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_dbs_backup_plan"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_dbs_backup_plan", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_dbs_backup_plan"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeBackupPlanList" {
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
			if *action == "DescribeBackupPlanBilling" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return DescribeBackupPlanBillingMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudDbsBackupPlanRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCbsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudDbsBackupPlanDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_dbs_backup_plan", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_dbs_backup_plan"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ReleaseBackupPlan" {
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
			if *action == "DescribeBackupPlanList" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudDbsBackupPlanDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
