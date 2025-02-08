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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_dbs_backup_plan",
		&resource.Sweeper{
			Name: "alicloud_dbs_backup_plan",
			F:    testSweepDbsBackupPlan,
		})
}

func testSweepDbsBackupPlan(region string) error {

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeBackupPlanList"
	request := map[string]interface{}{}
	request["RegionId"] = aliyunClient.RegionId

	request["PageSize"] = PageSizeLarge
	request["PageNum"] = 0

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = aliyunClient.RpcPost("Dbs", "2019-03-06", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.Items.BackupPlanDetail", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Items.BackupPlanDetail", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["BackupPlanName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Dbs Backup Plan: %s", item["BackupPlanName"].(string))
				continue
			}
			action := "ReleaseBackupPlan"
			request := map[string]interface{}{
				"BackupPlanId": item["BackupPlanId"],
			}
			_, err = aliyunClient.RpcPost("Dbs", "2019-03-06", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Dbs Backup Plan (%s): %s", item["BackupPlanName"].(string), err)
			}
			log.Printf("[INFO] Delete Dbs Backup Plan success: %s ", item["BackupPlanName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNum"] = request["PageNum"].(int) + 1
	}
	return nil
}

func TestAccAlicloudDBSBackupPlan_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dbs_backup_plan.default"
	ra := resourceAttrInit(resourceId, AlicloudDBSBackupPlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DbsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDbsBackupPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdbsbackupplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDBSBackupPlanBasicDependence0)
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
					"backup_plan_name":              "${var.name}",
					"payment_type":                  "PayAsYouGo",
					"instance_class":                "xlarge",
					"backup_method":                 "logical",
					"database_type":                 "MySQL",
					"database_region":               "${var.database_region}",
					"storage_region":                "${var.storage_region}",
					"instance_type":                 "RDS",
					"source_endpoint_instance_type": "RDS",
					"resource_group_id":             "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"source_endpoint_region":        "${var.source_endpoint_region}",
					"source_endpoint_instance_id":   "${alicloud_db_instance.default.id}",
					"source_endpoint_user_name":     "${alicloud_db_account_privilege.default.account_name}",
					"source_endpoint_password":      "${alicloud_rds_account.default.account_password}",
					"backup_objects":                `[{\"DBName\":\"${alicloud_db_database.default.name}\"}]`,
					"backup_period":                 "Monday",
					"backup_start_time":             "14:22",
					"backup_storage_type":           "system",
					"backup_retention_period":       "740",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_plan_name":        name,
						"payment_type":            "PayAsYouGo",
						"instance_class":          "xlarge",
						"backup_method":           "logical",
						"database_type":           "MySQL",
						"database_region":         CHECKSET,
						"storage_region":          CHECKSET,
						"instance_type":           "RDS",
						"backup_retention_period": "740",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_endpoint_password", "resource_group_id", "period", "database_region", "backup_strategy_type", "source_endpoint_port", "instance_type", "backup_speed_limit", "backup_log_interval_seconds", "source_endpoint_ip", "used_time", "backup_rate_limit", "storage_region"},
			},
		},
	})
}

var AlicloudDBSBackupPlanMap0 = map[string]string{
	"source_endpoint_port":                 NOSET,
	"payment_type":                         CHECKSET,
	"backup_period":                        CHECKSET,
	"backup_plan_name":                     CHECKSET,
	"backup_storage_type":                  CHECKSET,
	"enable_backup_log":                    CHECKSET,
	"instance_type":                        NOSET,
	"backup_objects":                       CHECKSET,
	"backup_speed_limit":                   NOSET,
	"backup_start_time":                    CHECKSET,
	"duplication_infrequent_access_period": CHECKSET,
	"source_endpoint_user_name":            CHECKSET,
	"backup_log_interval_seconds":          NOSET,
	"source_endpoint_ip":                   NOSET,
	"source_endpoint_instance_id":          CHECKSET,
	"used_time":                            NOSET,
	"backup_rate_limit":                    NOSET,
	"backup_retention_period":              CHECKSET,
	"duplication_archive_period":           CHECKSET,
	"storage_region":                       NOSET,
	"database_region":                      NOSET,
	"period":                               NOSET,
	"status":                               CHECKSET,
	"backup_strategy_type":                 NOSET,
	"source_endpoint_instance_type":        CHECKSET,
	"source_endpoint_region":               CHECKSET,
}

func AlicloudDBSBackupPlanBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
variable "database_region" {
  default = "%s"
}
variable "storage_region" {
  default = "%s"
}
variable "source_endpoint_region" {
  default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "local_ssd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "local_ssd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids.0
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_db_instance" "default" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "local_ssd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = local.vswitch_id
	instance_name = var.name
	security_group_ids = alicloud_security_group.default.*.id
}
resource "alicloud_db_database" "default" {
  instance_id = alicloud_db_instance.default.id
  name        = "tftestdatabase"
}
resource "alicloud_rds_account" "default" {
  db_instance_id = alicloud_db_instance.default.id
  account_name        = "tftestnormal000"
  account_password    = "Test12345"
}
resource "alicloud_db_account_privilege" "default" {
  instance_id  = alicloud_db_instance.default.id
  account_name = alicloud_rds_account.default.account_name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.default.name]
}
`, name, defaultRegionToTest, defaultRegionToTest, defaultRegionToTest)
}

func TestAccAlicloudDBSBackupPlan_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dbs_backup_plan.default"
	ra := resourceAttrInit(resourceId, AlicloudDBSBackupPlanMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DbsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDbsBackupPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdbsbackupplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDBSBackupPlanBasicDependence0)
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
					"backup_plan_name":              "${var.name}",
					"payment_type":                  "PayAsYouGo",
					"instance_class":                "xlarge",
					"backup_method":                 "logical",
					"database_type":                 "MySQL",
					"database_region":               "${var.database_region}",
					"storage_region":                "${var.storage_region}",
					"instance_type":                 "RDS",
					"source_endpoint_instance_type": "RDS",
					"source_endpoint_region":        "${var.database_region}",
					"source_endpoint_instance_id":   "${alicloud_db_instance.default.id}",
					"source_endpoint_user_name":     "${alicloud_db_account_privilege.default.account_name}",
					"source_endpoint_password":      "${alicloud_rds_account.default.account_password}",
					"backup_objects":                `[{\"DBName\":\"${alicloud_db_database.default.name}\"}]`,
					"backup_period":                 "Monday",
					"backup_start_time":             "14:22",
					"backup_storage_type":           "system",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_plan_name": name,
						"payment_type":     "PayAsYouGo",
						"instance_class":   "xlarge",
						"backup_method":    "logical",
						"database_type":    "MySQL",
						"database_region":  CHECKSET,
						"storage_region":   CHECKSET,
						"instance_type":    "RDS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "pause",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "pause",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "running",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_endpoint_password", "source_endpoint_port", "instance_type", "backup_speed_limit", "source_endpoint_ip", "used_time", "backup_log_interval_seconds", "backup_rate_limit", "storage_region", "database_region", "period", "backup_strategy_type"},
			},
		},
	})
}

var AlicloudDBSBackupPlanMap1 = map[string]string{
	"backup_rate_limit":                    NOSET,
	"backup_retention_period":              CHECKSET,
	"duplication_archive_period":           CHECKSET,
	"storage_region":                       NOSET,
	"database_region":                      NOSET,
	"period":                               NOSET,
	"status":                               CHECKSET,
	"backup_strategy_type":                 NOSET,
	"source_endpoint_instance_type":        CHECKSET,
	"source_endpoint_region":               CHECKSET,
	"source_endpoint_port":                 NOSET,
	"payment_type":                         CHECKSET,
	"backup_period":                        CHECKSET,
	"backup_plan_name":                     CHECKSET,
	"backup_storage_type":                  CHECKSET,
	"enable_backup_log":                    CHECKSET,
	"instance_type":                        NOSET,
	"backup_objects":                       CHECKSET,
	"backup_speed_limit":                   NOSET,
	"backup_start_time":                    CHECKSET,
	"duplication_infrequent_access_period": CHECKSET,
	"source_endpoint_user_name":            CHECKSET,
	"backup_log_interval_seconds":          NOSET,
	"source_endpoint_ip":                   NOSET,
	"source_endpoint_instance_id":          CHECKSET,
	"used_time":                            NOSET,
}

func TestUnitAccAlicloudDbsBackupPlan(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
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
	rawClient, err := sharedClientForRegion(region)
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
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_dbs_backup_plan", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCbsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
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
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
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
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
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
