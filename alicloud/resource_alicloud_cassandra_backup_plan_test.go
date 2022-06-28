package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/alibabacloud-go/tea-rpc/client"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func SkipTestAccAlicloudCassandraBackupPlan_basic0(t *testing.T) {
	// Cassandra has been offline
	t.Skip("Cassandra has been offline")
	var v map[string]interface{}
	resourceId := "alicloud_cassandra_backup_plan.default"
	ra := resourceAttrInit(resourceId, AlicloudCassandraBackupPlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CassandraService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCassandraBackupPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scassandrabackupplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCassandraBackupPlanBasicDependence0)
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
					"data_center_id": "${alicloud_cassandra_cluster.default.zone_id}",
					"cluster_id":     "${alicloud_cassandra_cluster.default.id}",
					"backup_time":    "00:10Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_center_id": CHECKSET,
						"cluster_id":     CHECKSET,
						"backup_time":    "00:10Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"active": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"active": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period": "Tuesday",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period": "Tuesday",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_time": "00:30Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "00:30Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"active":           "false",
					"backup_period":    "Monday",
					"backup_time":      "00:10Z",
					"retention_period": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"active":           "false",
						"backup_period":    "Monday",
						"backup_time":      "00:10Z",
						"retention_period": "1",
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

var AlicloudCassandraBackupPlanMap0 = map[string]string{

	"retention_period": CHECKSET,
	"backup_period":    CHECKSET,
	"data_center_id":   CHECKSET,
	"cluster_id":       CHECKSET,
	"backup_time":      "00:10Z",
	"active":           CHECKSET,
}

func AlicloudCassandraBackupPlanBasicDependence0(name string) string {
	return fmt.Sprintf(` 
		variable "name" {
		  default = "%s"
		}

		data "alicloud_cassandra_zones" "default" {
		}
		
		data "alicloud_vpcs" "default" {
		  name_regex = "default-NODELETING"
		}
		
		data "alicloud_vswitches" "default_1" {
		  vpc_id = data.alicloud_vpcs.default.ids[0]
		  zone_id = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)+(-1)].id
		}
		
		resource "alicloud_vswitch" "this_1" {
		  count = length(data.alicloud_vswitches.default_1.ids) > 0 ? 0 : 1
		  vswitch_name = var.name
		  vpc_id = data.alicloud_vpcs.default.ids.0
		  zone_id = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)+(-1)].id
		  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
		}
		resource "alicloud_cassandra_cluster" "default" {
		  cluster_name = var.name
		  data_center_name = var.name
		  auto_renew = "false"
		  instance_type = "cassandra.c.large"
		  major_version = "3.11"
		  node_count = "2"
		  pay_type = "PayAsYouGo"
		  vswitch_id = length(data.alicloud_vswitches.default_1.ids) > 0 ? data.alicloud_vswitches.default_1.ids[0] : alicloud_vswitch.this_1[0].id
		  disk_size = "160"
		  disk_type = "cloud_ssd"
		  maintain_start_time = "18:00Z"
		  maintain_end_time = "20:00Z"
		  ip_white = "127.0.0.1"
		}
	`, name)
}

func TestUnitAlicloudCassandraBackupPlan(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_cassandra_backup_plan"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cassandra_backup_plan"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"data_center_id":   "CreateBackupPlanValue",
		"cluster_id":       "CreateBackupPlanValue",
		"backup_time":      "CreateBackupPlanValue",
		"active":           false,
		"backup_period":    "CreateBackupPlanValue",
		"retention_period": 1,
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
		// DescribeBackupPlans
		"BackupPlans": map[string]interface{}{
			"BackupPlan": []interface{}{
				map[string]interface{}{
					"ClusterId":       "CreateBackupPlanValue",
					"DataCenterId":    "CreateBackupPlanValue",
					"Active":          false,
					"BackupPeriod":    "CreateBackupPlanValue",
					"BackupTime":      "CreateBackupPlanValue",
					"RetentionPeriod": 1,
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateBackupPlan
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cassandra_backup_plan", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudCassandraBackupPlanCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeBackupPlans Response
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateBackupPlan" {
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
		err := resourceAlicloudCassandraBackupPlanCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cassandra_backup_plan"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudCassandraBackupPlanUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyBackupPlan
	attributesDiff := map[string]interface{}{
		"backup_time":      "ModifyBackupPlanValue",
		"active":           true,
		"backup_period":    "ModifyBackupPlanValue",
		"retention_period": 10,
	}
	diff, err := newInstanceDiff("alicloud_cassandra_backup_plan", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cassandra_backup_plan"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeBackupPlans Response
		"BackupPlans": map[string]interface{}{
			"BackupPlan": []interface{}{
				map[string]interface{}{
					"Active":          true,
					"BackupPeriod":    "ModifyBackupPlanValue",
					"BackupTime":      "ModifyBackupPlanValue",
					"RetentionPeriod": 10,
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyBackupPlan" {
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
		err := resourceAlicloudCassandraBackupPlanUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cassandra_backup_plan"].Schema).Data(dExisted.State(), nil)
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
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_cassandra_backup_plan", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cassandra_backup_plan"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeBackupPlans" {
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
		err := resourceAlicloudCassandraBackupPlanRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudCassandraBackupPlanDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_cassandra_backup_plan", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cassandra_backup_plan"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteBackupPlan" {
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
		err := resourceAlicloudCassandraBackupPlanDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
