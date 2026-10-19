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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_rds_backup", &resource.Sweeper{
		Name: "alicloud_rds_backup",
		F:    testSweepRdsBackups,
	})
}

// testSweepRdsBackups sweeps leftover test RDS manual backups. Backups are
// sub-resources that live on a DB instance, so the sweeper enumerates DB
// instances first, then the backups on each. Only manual backups
// (BackupMode=Manual, created by tests) are deleted; automated backups are
// managed by the backup policy and auto-expire. Instances carrying a
// do-not-delete marker (or, in non-SWEEP_ALL mode, an instance name that
// does not match test prefixes) are skipped.
func testSweepRdsBackups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{"tf-testAcc", "tf_testAcc"}

	instances, err := listSweepRdsDbInstances(client)
	if err != nil {
		log.Printf("[ERROR] DescribeDBInstances got an error: %v", err)
		return nil
	}

	for _, inst := range instances {
		instanceId := fmt.Sprint(inst["DBInstanceId"])
		instanceName := fmt.Sprint(inst["DBInstanceDescription"])
		if protectedSweepVpcName(instanceName) {
			log.Printf("[INFO] Skipping RDS Backups on do-not-delete instance: %s (%s)", instanceName, instanceId)
			continue
		}
		if !sweepAll() {
			match := false
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(instanceName), strings.ToLower(prefix)) {
					match = true
					break
				}
			}
			if !match {
				log.Printf("[INFO] Skipping RDS Backups (instance %s does not match test prefixes): %s", instanceName, instanceId)
				continue
			}
		}

		action := "DescribeBackups"
		req := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": instanceId,
			"SourceIp":     client.SourceIp,
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		if e := retry.Retry(5*time.Minute, func() *retry.RetryError {
			resp, e := client.RpcPost("Rds", "2014-08-15", action, nil, req, true)
			response = resp
			if e != nil {
				if rdsErrorHasCode(e, "InvalidDBInstanceId.NotFound", "InvalidDBInstanceName.NotFound") || NotFoundError(e) {
					return nil
				}
				if NeedRetry(e) || IsExpectedErrors(e, dbInstanceStatusErrorCodes) {
					wait()
					return retry.RetryableError(e)
				}
				return retry.NonRetryableError(e)
			}
			return nil
		}); e != nil {
			log.Printf("[ERROR] DescribeBackups for %s got an error: %v", instanceId, e)
			continue
		}
		addDebug(action, response, req)

		bkResp, e := jsonpath.Get("$.Items.Backup", response)
		if e != nil {
			log.Printf("[ERROR] DescribeBackups parse for %s: %v", instanceId, e)
			continue
		}
		backups, _ := bkResp.([]interface{})
		for _, b := range backups {
			bk := b.(map[string]interface{})
			backupId := fmt.Sprint(bk["BackupId"])
			backupMode := fmt.Sprint(bk["BackupMode"])
			backupStatus := fmt.Sprint(bk["BackupStatus"])
			// Only manual backups created by tests are deletable; automated
			// backups are managed by the backup policy and auto-expire.
			if !strings.EqualFold(backupMode, "manual") {
				log.Printf("[INFO] Skipping RDS Backup (non-manual mode %s): %s/%s", backupMode, instanceId, backupId)
				continue
			}
			// Only delete finished backups; an in-progress backup cannot be removed.
			if backupStatus != "" && !strings.EqualFold(backupStatus, "success") && !strings.EqualFold(backupStatus, "finished") {
				log.Printf("[INFO] Skipping RDS Backup (status %s): %s/%s", backupStatus, instanceId, backupId)
				continue
			}

			log.Printf("[INFO] Deleting RDS Backup: %s/%s", instanceId, backupId)
			delAction := "DeleteBackup"
			delReq := map[string]interface{}{
				"DBInstanceId": instanceId,
				"BackupId":     backupId,
				"RegionId":     client.RegionId,
			}
			if e := retry.Retry(2*time.Minute, func() *retry.RetryError {
				resp, e := client.RpcPost("Rds", "2014-08-15", delAction, nil, delReq, true)
				if e != nil {
					if NotFoundError(e) || rdsErrorHasCode(e, "InvalidDBInstanceId.NotFound", "InvalidDBInstanceName.NotFound", "InvalidBackupId.NotFound") {
						return nil
					}
					if NeedRetry(e) || IsExpectedErrors(e, dbInstanceStatusErrorCodes) {
						return retry.RetryableError(e)
					}
					return retry.NonRetryableError(e)
				}
				addDebug(delAction, resp, delReq)
				return nil
			}); e != nil {
				log.Printf("[ERROR] Failed to delete RDS Backup %s/%s: %s", instanceId, backupId, e)
			}
		}
	}

	return nil
}

func TestAccAliCloudRdsBackup_MySQL(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_backup.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsBackupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsBackup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdsbackup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsBackupBasicDependenceMySQL)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":    "${alicloud_db_instance.default.id}",
					"backup_method":     "Snapshot",
					"remove_from_state": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": CHECKSET,
						"backup_method":  "Snapshot",
						"backup_id":      CHECKSET,
						"store_status":   CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_from_state", "db_name", "backup_strategy"},
			},
		},
	})
}

var AlicloudRdsBackupMap0 = map[string]string{}

func AlicloudRdsBackupBasicDependenceMySQL(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids[0]
	zone_id      = data.alicloud_db_zones.default.ids[0]
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids[0]
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_db_zones.default.ids[0]
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_resource_manager_resource_group" "default" {
	resource_group_name = var.name
	display_name        = var.name
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "cloud_essd"
  instance_type            = "mysql.x2.large.2c"
  instance_storage         = 20
  vswitch_id               = local.vswitch_id
  instance_name            = var.name
}

resource "alicloud_db_database" "default" {
  instance_id = alicloud_db_instance.default.id
  name        = "tftestdatabase"
}
`, name)
}

func TestAccAliCloudRdsBackup_PostgreSQL(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_backup.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsBackupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsBackup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdsbackup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsBackupBasicDependencePostgreSQL0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":    "${alicloud_db_instance.default.id}",
					"backup_method":     "Snapshot",
					"remove_from_state": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": CHECKSET,
						"backup_method":  "Snapshot",
						"backup_id":      CHECKSET,
						"store_status":   CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_from_state", "db_name", "backup_strategy"},
			},
		},
	})
}

func AlicloudRdsBackupBasicDependencePostgreSQL0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_db_zones" "default" {
  engine                   = "PostgreSQL"
  engine_version           = "14.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones[0].id
  engine                   = "PostgreSQL"
  engine_version           = "14.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids[0]
	zone_id      = data.alicloud_db_zones.default.ids[0]
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids[0]
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_db_zones.default.ids[0]
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_db_instance" "default" {
  engine                   = "PostgreSQL"
  engine_version           = "14.0"
  db_instance_storage_type = "cloud_essd"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes[0].storage_range[0].min
  vswitch_id               = local.vswitch_id
  instance_name            = var.name
}
`, name)
}

func TestAccAliCloudRdsBackup_SQLServer(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_backup.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsBackupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsBackup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdsbackup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsBackupBasicDependenceSQLServer0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":    "${alicloud_db_instance.default.id}",
					"backup_method":     "Physical",
					"remove_from_state": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": CHECKSET,
						"backup_method":  "Physical",
						"backup_id":      CHECKSET,
						"store_status":   CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_from_state", "db_name", "backup_strategy"},
			},
		},
	})
}

func TestAccAliCloudRdsBackup_SQLServer1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_backup.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsBackupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsBackup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdsbackup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsBackupBasicDependenceSQLServer0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":          "${alicloud_db_instance.default.id}",
					"backup_strategy":         "db",
					"backup_method":           "Physical",
					"backup_type":             "FullBackup",
					"backup_retention_period": "7",
					"db_name":                 "${alicloud_db_database.default.name}",
					"remove_from_state":       "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": CHECKSET,
						"backup_method":  "Physical",
						"backup_type":    "FullBackup",
						"backup_id":      CHECKSET,
						"store_status":   CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_from_state", "db_name", "backup_strategy", "backup_retention_period"},
			},
		},
	})
}

func AlicloudRdsBackupBasicDependenceSQLServer0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_db_zones" "default" {
  engine                   = "SQLServer"
  engine_version           = "2019_ent"
  instance_charge_type     = "PostPaid"
  category                 = "AlwaysOn"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids[0]
	zone_id      = data.alicloud_db_zones.default.ids[0]
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids[0]
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_db_zones.default.ids[0]
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_db_instance" "default" {
  engine                   = "SQLServer"
  engine_version           = "2019_ent"
  db_instance_storage_type = "cloud_essd"
  instance_type            = "mssql.x4.medium.e2"
  instance_storage         = 20
  vswitch_id               = local.vswitch_id
  instance_name            = var.name
  category                 = "AlwaysOn"
}

resource "alicloud_db_database" "default" {
  instance_id   = alicloud_db_instance.default.id
  name          = "tftestdatabase"
  character_set = "Chinese_PRC_CI_AS"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudRdsBackup(t *testing.T) {
	p := Provider().ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_rds_backup"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_rds_backup"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	dCreateCompletion, _ := schema.InternalMap(p["alicloud_rds_backup"].Schema).Data(nil, nil)
	dCreateCompletion.MarkNewResource()
	for key, value := range map[string]interface{}{
		"db_instance_id":    "db_instance_id",
		"db_name":           "db_name",
		"backup_strategy":   "instance",
		"backup_method":     "Snapshot",
		"backup_type":       "FullBackup",
		"remove_from_state": true,
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
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
			"Backup": []interface{}{
				map[string]interface{}{
					"BackupMethod": "Snapshot",
					"BackupType":   "FullBackup",
					"DBInstanceId": "DBInstanceId",
					"BackupId":     "MockBackupId",
					"StoreStatus":  "store_status",
				},
			},
			"BackupJob": []interface{}{
				map[string]interface{}{
					"BackupStatus": "Finished",
					"BackupJobId":  "MockBackupJobId",
					"BackupId":     "MockBackupId",
				},
			},
		},
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
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_rds_backup", "MockBackupId"))
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
			result := ReadMockResponse
			result["BackupJobId"] = "MockBackupJobId"
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewRdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudRdsBackupCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudRdsBackupCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudRdsBackupCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("DBInstanceId", ":", "MockBackupId"))
	// Update
	t.Run("UpdateNormal", func(t *testing.T) {
		patcheDescribeBackups := gomonkey.ApplyMethod(reflect.TypeOf(&RdsService{}), "DescribeRdsBackup", func(*RdsService, string) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudRdsBackupUpdate(d, rawClient)
		patcheDescribeBackups.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewRdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudRdsBackupDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudRdsBackupDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})

		err := resourceAliCloudRdsBackupDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("DeleteNonRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudRdsBackupDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("DeleteStoreStatusAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		resourceData, _ := schema.InternalMap(p["alicloud_rds_backup"].Schema).Data(nil, diff)
		err = resourceData.Set("store_status", "Disabled")
		err = resourceData.Set("remove_from_state", false)
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudRdsBackupDelete(resourceData, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	//Read
	t.Run("ReadDescribeBackupsNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudRdsBackupRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeBackupsAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := true
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudRdsBackupRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}
