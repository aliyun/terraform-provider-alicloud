package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudRdsInstanceCrossBackupPolicyMySql(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_instance_cross_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeInstanceCrossBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccRdsCrossBackupPolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceRdsInstanceCrossBackupPolicyMysqlConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.RDSInstanceClassesSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":         "${alicloud_db_instance.default.id}",
					"cross_backup_region": "${data.alicloud_rds_cross_regions.regions.ids[0]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":         CHECKSET,
						"cross_backup_region": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_backup_region": "${data.alicloud_rds_cross_regions.regions.ids[1]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_backup_region": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup_enabled": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_enabled": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention": "15",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention": "15",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":         "${alicloud_db_instance.default.id}",
					"cross_backup_region": "${data.alicloud_rds_cross_regions.regions.ids[2]}",
					"log_backup_enabled":  "Enable",
					"retention":           "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":         CHECKSET,
						"cross_backup_region": CHECKSET,
						"log_backup_enabled":  "Enable",
						"retention":           "30",
					}),
				),
			}},
	})
}

func resourceRdsInstanceCrossBackupPolicyMysqlConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "local_ssd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones[0].id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "local_ssd"
  instance_charge_type     = "PostPaid"
}

data "alicloud_rds_cross_regions" "regions" {
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_db_zones.default.ids[0]
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "local_ssd"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min
  vswitch_id               = data.alicloud_vswitches.default.ids[0]
  instance_name            = var.name
}
`, name)
}

func TestAccAliCloudRdsInstanceCrossBackupPolicyPostgreSQL(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_instance_cross_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeInstanceCrossBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccRdsCrossBackupPolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceRdsInstanceCrossBackupPolicyPostgreSQLConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.RDSInstanceClassesSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":         "${alicloud_db_instance.default.id}",
					"cross_backup_region": "${data.alicloud_rds_cross_regions.regions.ids[0]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":         CHECKSET,
						"cross_backup_region": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_backup_region": "${data.alicloud_rds_cross_regions.regions.ids[1]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_backup_region": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup_enabled": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_enabled": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention": "15",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention": "15",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":         "${alicloud_db_instance.default.id}",
					"cross_backup_region": "${data.alicloud_rds_cross_regions.regions.ids[2]}",
					"log_backup_enabled":  "Enable",
					"retention":           "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":         CHECKSET,
						"cross_backup_region": CHECKSET,
						"log_backup_enabled":  "Enable",
						"retention":           "30",
					}),
				),
			}},
	})
}

func resourceRdsInstanceCrossBackupPolicyPostgreSQLConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_db_zones" "default" {
  engine                   = "PostgreSQL"
  engine_version           = "18.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones[0].id
  engine                   = "PostgreSQL"
  engine_version           = "18.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

data "alicloud_rds_cross_regions" "regions" {
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_db_zones.default.zones[0].id
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_db_instance" "default" {
  engine                   = "PostgreSQL"
  engine_version           = "18.0"
  db_instance_storage_type = "cloud_essd"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min
  vswitch_id               = data.alicloud_vswitches.default.ids[0]
  instance_name            = var.name
}
`, name)
}

// TestUnitAlicloudRdsInstanceCrossBackupPolicy tests the cross backup policy resource
// with focus on the delete behavior that waits for DB instance to return to Running status
func TestUnitAlicloudRdsInstanceCrossBackupPolicy(t *testing.T) {
	p := Provider().ResourcesMap
	dCreate, _ := schema.InternalMap(p["alicloud_rds_instance_cross_backup_policy"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	dDelete, _ := schema.InternalMap(p["alicloud_rds_instance_cross_backup_policy"].Schema).Data(nil, nil)

	for key, value := range map[string]interface{}{
		"instance_id":         "pgm-test123456",
		"cross_backup_region": "eu-west-1",
	} {
		// lintignore:R001
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		// lintignore:R001
		err = dDelete.Set(key, value)
		assert.Nil(t, err)
	}
	dDelete.SetId("pgm-test123456")

	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}

	ReadMockResponse := map[string]interface{}{
		"BackupEnabled":     "Enable",
		"LogBackupEnabled":  "Enable",
		"CrossBackupRegion": "eu-west-1",
		"Retention":         7,
		"DBInstanceStatus":  "Running",
		"LockMode":          "Unlock",
	}

	ReadMockResponseDBInstance := map[string]interface{}{
		"DBInstanceStatus": "Running",
		"DBInstanceId":     "pgm-test123456",
	}

	ReadMockResponseDBInstanceModifying := map[string]interface{}{
		"DBInstanceStatus": "Modifying",
		"DBInstanceId":     "pgm-test123456",
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_rds_instance_cross_backup_policy", "pgm-test123456"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			return ReadMockResponse, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			return ReadMockResponse, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			return ReadMockResponse, nil
		},
		"ReadDBInstanceRunning": func(errorCode string) (map[string]interface{}, error) {
			return ReadMockResponseDBInstance, nil
		},
		"ReadDBInstanceModifying": func(errorCode string) (map[string]interface{}, error) {
			return ReadMockResponseDBInstanceModifying, nil
		},
	}

	// Test Delete - Normal case where DB instance is already Running
	t.Run("DeleteNormal", func(t *testing.T) {
		patchDescribeCrossBackupPolicy := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeInstanceCrossBackupPolicy",
			func(*RdsService, string) (map[string]interface{}, error) {
				return responseMock["ReadNormal"]("")
			})
		defer patchDescribeCrossBackupPolicy.Reset()

		patchRpcPost := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost",
			func(*connectivity.AliyunClient, string, string, string, map[string]string, map[string]interface{}, bool) (map[string]interface{}, error) {
				return responseMock["DeleteNormal"]("")
			})
		defer patchRpcPost.Reset()

		patchDescribeDBInstance := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeDBInstance",
			func(*RdsService, string) (map[string]interface{}, error) {
				return responseMock["ReadDBInstanceRunning"]("")
			})
		defer patchDescribeDBInstance.Reset()

		err := resourceAlicloudRdsInstanceCrossBackupPolicyDelete(dDelete, rawClient)
		assert.Nil(t, err)
	})

	// Test Delete - DB instance transitions from Modifying to Running BEFORE delete
	t.Run("DeleteWaitsForRunningStatusBeforeDelete", func(t *testing.T) {
		patchDescribeCrossBackupPolicy := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeInstanceCrossBackupPolicy",
			func(*RdsService, string) (map[string]interface{}, error) {
				return responseMock["ReadNormal"]("")
			})
		defer patchDescribeCrossBackupPolicy.Reset()

		patchRpcPost := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost",
			func(*connectivity.AliyunClient, string, string, string, map[string]string, map[string]interface{}, bool) (map[string]interface{}, error) {
				return responseMock["DeleteNormal"]("")
			})
		defer patchRpcPost.Reset()

		// Simulate DB instance transitioning from Modifying to Running BEFORE the delete API call
		callCount := 0
		patchDescribeDBInstance := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeDBInstance",
			func(*RdsService, string) (map[string]interface{}, error) {
				callCount++
				if callCount <= 2 {
					return responseMock["ReadDBInstanceModifying"]("")
				}
				return responseMock["ReadDBInstanceRunning"]("")
			})
		defer patchDescribeDBInstance.Reset()

		err := resourceAlicloudRdsInstanceCrossBackupPolicyDelete(dDelete, rawClient)
		assert.Nil(t, err)
		// Verify that DescribeDBInstance was called multiple times (waited for Running status before delete)
		assert.True(t, callCount >= 2, "Expected multiple calls to DescribeDBInstance while waiting for Running status before delete")
	})

	// Test Delete - DB instance not found during state wait (policy already gone)
	t.Run("DeleteDBInstanceNotFoundDuringWait", func(t *testing.T) {
		// RdsDBInstanceStateRefreshFunc returns nil, "", nil when instance not found
		patchDescribeDBInstance := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeDBInstance",
			func(*RdsService, string) (map[string]interface{}, error) {
				return responseMock["NotFoundError"]("")
			})
		defer patchDescribeDBInstance.Reset()

		err := resourceAlicloudRdsInstanceCrossBackupPolicyDelete(dDelete, rawClient)
		assert.Nil(t, err)
	})

	// Test Delete - Policy already not found (idempotent delete)
	t.Run("DeleteNotFound", func(t *testing.T) {
		patchDescribeDBInstance := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeDBInstance",
			func(*RdsService, string) (map[string]interface{}, error) {
				return responseMock["ReadDBInstanceRunning"]("")
			})
		defer patchDescribeDBInstance.Reset()

		patchDescribeCrossBackupPolicy := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeInstanceCrossBackupPolicy",
			func(*RdsService, string) (map[string]interface{}, error) {
				return responseMock["NotFoundError"]("")
			})
		defer patchDescribeCrossBackupPolicy.Reset()

		err := resourceAlicloudRdsInstanceCrossBackupPolicyDelete(dDelete, rawClient)
		assert.Nil(t, err)
	})

	// Test Delete - LockMode is null (nothing to delete)
	t.Run("DeleteLockModeNull", func(t *testing.T) {
		patchDescribeDBInstance := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeDBInstance",
			func(*RdsService, string) (map[string]interface{}, error) {
				return responseMock["ReadDBInstanceRunning"]("")
			})
		defer patchDescribeDBInstance.Reset()

		lockModeNullResponse := map[string]interface{}{
			"BackupEnabled": "Enable",
			"LockMode":      "null",
		}
		patchDescribeCrossBackupPolicy := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeInstanceCrossBackupPolicy",
			func(*RdsService, string) (map[string]interface{}, error) {
				return lockModeNullResponse, nil
			})
		defer patchDescribeCrossBackupPolicy.Reset()

		err := resourceAlicloudRdsInstanceCrossBackupPolicyDelete(dDelete, rawClient)
		assert.Nil(t, err)
	})

	// Test Delete - Policy already disabled at initial check (Ansatz 2 main case)
	// This is the key test - when we check upfront and policy is already disabled,
	// we should NOT call ModifyInstanceCrossBackupPolicy at all
	t.Run("DeletePolicyAlreadyDisabledAtStart", func(t *testing.T) {
		patchDescribeDBInstance := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeDBInstance",
			func(*RdsService, string) (map[string]interface{}, error) {
				return responseMock["ReadDBInstanceRunning"]("")
			})
		defer patchDescribeDBInstance.Reset()

		disabledPolicyResponse := map[string]interface{}{
			"BackupEnabled":     "Disabled",
			"LogBackupEnabled":  "Disabled",
			"CrossBackupRegion": "eu-west-1",
			"Retention":         7,
			"DBInstanceStatus":  "Running",
			"LockMode":          "Unlock",
		}
		patchDescribeCrossBackupPolicy := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeInstanceCrossBackupPolicy",
			func(*RdsService, string) (map[string]interface{}, error) {
				return disabledPolicyResponse, nil
			})
		defer patchDescribeCrossBackupPolicy.Reset()

		// RpcPost should NOT be called at all
		rpcCallCount := 0
		patchRpcPost := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost",
			func(*connectivity.AliyunClient, string, string, string, map[string]string, map[string]interface{}, bool) (map[string]interface{}, error) {
				rpcCallCount++
				return responseMock["DeleteNormal"]("")
			})
		defer patchRpcPost.Reset()

		err := resourceAlicloudRdsInstanceCrossBackupPolicyDelete(dDelete, rawClient)
		assert.Nil(t, err)
		// RpcPost should NOT have been called because policy was already disabled
		assert.Equal(t, 0, rpcCallCount, "RpcPost should not be called when policy is already disabled")
	})

	// Test Delete - DescribeInstanceCrossBackupPolicy returns InternalError (treat as not found)
	t.Run("DeleteDescribePolicyInternalError", func(t *testing.T) {
		patchDescribeDBInstance := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeDBInstance",
			func(*RdsService, string) (map[string]interface{}, error) {
				return responseMock["ReadDBInstanceRunning"]("")
			})
		defer patchDescribeDBInstance.Reset()

		patchDescribeCrossBackupPolicy := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeInstanceCrossBackupPolicy",
			func(*RdsService, string) (map[string]interface{}, error) {
				return responseMock["RetryError"]("InternalError")
			})
		defer patchDescribeCrossBackupPolicy.Reset()

		// RpcPost should NOT be called
		rpcCallCount := 0
		patchRpcPost := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost",
			func(*connectivity.AliyunClient, string, string, string, map[string]string, map[string]interface{}, bool) (map[string]interface{}, error) {
				rpcCallCount++
				return responseMock["DeleteNormal"]("")
			})
		defer patchRpcPost.Reset()

		err := resourceAlicloudRdsInstanceCrossBackupPolicyDelete(dDelete, rawClient)
		assert.Nil(t, err)
		assert.Equal(t, 0, rpcCallCount, "RpcPost should not be called when DescribeInstanceCrossBackupPolicy returns InternalError")
	})

	// Test Delete - API returns non-retryable error
	t.Run("DeleteAbnormal", func(t *testing.T) {
		patchDescribeDBInstance := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeDBInstance",
			func(*RdsService, string) (map[string]interface{}, error) {
				return responseMock["ReadDBInstanceRunning"]("")
			})
		defer patchDescribeDBInstance.Reset()

		patchDescribeCrossBackupPolicy := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeInstanceCrossBackupPolicy",
			func(*RdsService, string) (map[string]interface{}, error) {
				return responseMock["ReadNormal"]("")
			})
		defer patchDescribeCrossBackupPolicy.Reset()

		patchRpcPost := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost",
			func(*connectivity.AliyunClient, string, string, string, map[string]string, map[string]interface{}, bool) (map[string]interface{}, error) {
				return responseMock["NoRetryError"]("InvalidDBInstanceId.NotFound")
			})
		defer patchRpcPost.Reset()

		err := resourceAlicloudRdsInstanceCrossBackupPolicyDelete(dDelete, rawClient)
		assert.NotNil(t, err)
	})

	// Test Delete - InternalError is retried and eventually succeeds
	t.Run("DeleteRetriesInternalError", func(t *testing.T) {
		patchDescribeCrossBackupPolicy := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeInstanceCrossBackupPolicy",
			func(*RdsService, string) (map[string]interface{}, error) {
				return responseMock["ReadNormal"]("")
			})
		defer patchDescribeCrossBackupPolicy.Reset()

		// Simulate InternalError on first 2 calls, then success
		rpcCallCount := 0
		patchRpcPost := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost",
			func(*connectivity.AliyunClient, string, string, string, map[string]string, map[string]interface{}, bool) (map[string]interface{}, error) {
				rpcCallCount++
				if rpcCallCount <= 2 {
					return responseMock["RetryError"]("InternalError")
				}
				return responseMock["DeleteNormal"]("")
			})
		defer patchRpcPost.Reset()

		patchDescribeDBInstance := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeDBInstance",
			func(*RdsService, string) (map[string]interface{}, error) {
				return responseMock["ReadDBInstanceRunning"]("")
			})
		defer patchDescribeDBInstance.Reset()

		err := resourceAlicloudRdsInstanceCrossBackupPolicyDelete(dDelete, rawClient)
		assert.Nil(t, err)
		// Verify that RpcPost was called multiple times (retried InternalError)
		assert.True(t, rpcCallCount >= 2, "Expected multiple RpcPost calls due to InternalError retry")
	})

	// Test Delete - InternalError with policy becoming disabled after retries
	// This simulates the case where Alicloud API returns 500 but the policy becomes disabled
	t.Run("DeleteInternalErrorPolicyBecomesDisabled", func(t *testing.T) {
		patchDescribeDBInstance := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeDBInstance",
			func(*RdsService, string) (map[string]interface{}, error) {
				return responseMock["ReadDBInstanceRunning"]("")
			})
		defer patchDescribeDBInstance.Reset()

		// Simulate policy becoming disabled after first InternalError check
		describePolicyCallCount := 0
		patchDescribeCrossBackupPolicy := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeInstanceCrossBackupPolicy",
			func(*RdsService, string) (map[string]interface{}, error) {
				describePolicyCallCount++
				// First call (initial check) returns enabled, second call (after InternalError) returns NotFound
				if describePolicyCallCount <= 1 {
					return responseMock["ReadNormal"]("")
				}
				return responseMock["NotFoundError"]("")
			})
		defer patchDescribeCrossBackupPolicy.Reset()

		// Return InternalError on first call - the delete should succeed when policy check shows disabled
		rpcCallCount := 0
		patchRpcPost := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost",
			func(*connectivity.AliyunClient, string, string, string, map[string]string, map[string]interface{}, bool) (map[string]interface{}, error) {
				rpcCallCount++
				return responseMock["RetryError"]("InternalError")
			})
		defer patchRpcPost.Reset()

		err := resourceAlicloudRdsInstanceCrossBackupPolicyDelete(dDelete, rawClient)
		assert.Nil(t, err)
		// Verify that we checked the policy status after InternalError
		assert.True(t, describePolicyCallCount >= 2, "Expected at least 2 DescribeInstanceCrossBackupPolicy calls")
		assert.True(t, rpcCallCount >= 1, "Expected at least 1 RpcPost call")
	})

	// Test Delete - InternalError with policy showing BackupEnabled=Disabled after retry check
	t.Run("DeleteInternalErrorPolicyAlreadyDisabled", func(t *testing.T) {
		patchDescribeDBInstance := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeDBInstance",
			func(*RdsService, string) (map[string]interface{}, error) {
				return responseMock["ReadDBInstanceRunning"]("")
			})
		defer patchDescribeDBInstance.Reset()

		disabledPolicyResponse := map[string]interface{}{
			"BackupEnabled":     "Disabled",
			"LogBackupEnabled":  "Disabled",
			"CrossBackupRegion": "eu-west-1",
			"Retention":         7,
			"DBInstanceStatus":  "Running",
			"LockMode":          "Unlock",
		}

		describePolicyCallCount := 0
		patchDescribeCrossBackupPolicy := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeInstanceCrossBackupPolicy",
			func(*RdsService, string) (map[string]interface{}, error) {
				describePolicyCallCount++
				// First call (initial check) returns enabled, second call (after InternalError) returns disabled
				if describePolicyCallCount <= 1 {
					return responseMock["ReadNormal"]("")
				}
				return disabledPolicyResponse, nil
			})
		defer patchDescribeCrossBackupPolicy.Reset()

		// Return InternalError - should trigger re-check which shows disabled
		rpcCallCount := 0
		patchRpcPost := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost",
			func(*connectivity.AliyunClient, string, string, string, map[string]string, map[string]interface{}, bool) (map[string]interface{}, error) {
				rpcCallCount++
				return responseMock["RetryError"]("InternalError")
			})
		defer patchRpcPost.Reset()

		err := resourceAlicloudRdsInstanceCrossBackupPolicyDelete(dDelete, rawClient)
		assert.Nil(t, err)
		assert.True(t, describePolicyCallCount >= 2, "Expected at least 2 DescribeInstanceCrossBackupPolicy calls")
	})
}
