package rds_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudRdsBackup_MySQL(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_backup.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRdsBackupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRdsBackup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdsbackup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRdsBackupBasicDependenceMySQL)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      testAccCheckAlicloudRdsBackupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":    "${alicloud_db_instance.default.id}",
					"backup_method":     "Snapshot",
					"remove_from_state": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": acctest.CHECKSET,
						"backup_method":  "Snapshot",
						"backup_id":      acctest.CHECKSET,
						"backup_job_id":  acctest.CHECKSET,
						"store_status":   acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_from_state", "db_name", "backup_strategy", "backup_job_id"},
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRdsBackupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRdsBackup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdsbackup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRdsBackupBasicDependencePostgreSQL0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      testAccCheckAlicloudRdsBackupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":    "${alicloud_db_instance.default.id}",
					"backup_method":     "Snapshot",
					"remove_from_state": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": acctest.CHECKSET,
						"backup_method":  "Snapshot",
						"backup_id":      acctest.CHECKSET,
						"backup_job_id":  acctest.CHECKSET,
						"store_status":   acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_from_state", "db_name", "backup_strategy", "backup_job_id"},
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRdsBackupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRdsBackup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdsbackup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRdsBackupBasicDependenceSQLServer0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      testAccCheckAlicloudRdsBackupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":    "${alicloud_db_instance.default.id}",
					"backup_method":     "Physical",
					"remove_from_state": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": acctest.CHECKSET,
						"backup_method":  "Physical",
						"backup_id":      acctest.CHECKSET,
						"backup_job_id":  acctest.CHECKSET,
						"store_status":   acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_from_state", "db_name", "backup_strategy", "backup_job_id"},
			},
		},
	})
}

func TestAccAliCloudRdsBackup_SQLServer1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_backup.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRdsBackupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRdsBackup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdsbackup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRdsBackupBasicDependenceSQLServer0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      testAccCheckAlicloudRdsBackupDestroy,
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
						"db_instance_id": acctest.CHECKSET,
						"backup_method":  "Physical",
						"backup_type":    "FullBackup",
						"backup_id":      acctest.CHECKSET,
						"backup_job_id":  acctest.CHECKSET,
						"store_status":   acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remove_from_state", "db_name", "backup_strategy", "backup_retention_period", "backup_job_id"},
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
	rawClient, err := sweep.SharedClientForRegion(region)
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
		"errs.NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_rds_backup", "MockBackupId"))
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
		err := rds.ResourceAliCloudRdsBackupCreate(d, rawClient)
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
		err := rds.ResourceAliCloudRdsBackupCreate(d, rawClient)
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
		err := rds.ResourceAliCloudRdsBackupCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("DBInstanceId", ":", "MockBackupId"))
	// Update
	t.Run("UpdateNormal", func(t *testing.T) {
		patcheDescribeBackups := gomonkey.ApplyMethod(reflect.TypeOf(&rds.RdsService{}), "DescribeRdsBackup", func(*rds.RdsService, string) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		err := rds.ResourceAliCloudRdsBackupUpdate(d, rawClient)
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
		err := rds.ResourceAliCloudRdsBackupDelete(d, rawClient)
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
		err := rds.ResourceAliCloudRdsBackupDelete(d, rawClient)
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

		err := rds.ResourceAliCloudRdsBackupDelete(d, rawClient)
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
		err := rds.ResourceAliCloudRdsBackupDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("DeleteStoreStatusAbnormal", func(t *testing.T) {
		diff := sdktf.NewInstanceDiff()
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
		err := rds.ResourceAliCloudRdsBackupDelete(resourceData, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	//Read
	t.Run("ReadDescribeBackupsNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["errs.NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := rds.ResourceAliCloudRdsBackupRead(d, rawClient)
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
		err := rds.ResourceAliCloudRdsBackupRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

func testAccCheckAlicloudRdsBackupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_rds_backup" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		rdsServiceV2 := rds.NewRdsServiceV2(client)
		_, err := rdsServiceV2.DescribeRdsBackup(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("RDS backup %s still exists", rs.Primary.ID)
	}
	return nil
}

func rdsBackupDeleteTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("rds", endpoint)
	return client
}

// The resource ID is the composite "DBInstanceId:BackupId" but DeleteBackup only
// accepts the bare backup set ID, so these cases fail if Delete regresses to
// sending d.Id() as BackupId or calls the API without a deletable backup.
func TestUnitAlicloudRdsBackupDeleteCompositeId(t *testing.T) {
	for _, tc := range []struct {
		name            string
		id              string
		storeStatus     string
		removeFromState bool
		wantErr         bool
		wantCalls       int
	}{
		{name: "enabled_composite_id", id: "rm-dbtestinstance:2091052", storeStatus: "Enabled", wantCalls: 1},
		{name: "disabled_remove_from_state", id: "rm-dbtestinstance:2091052", storeStatus: "Disabled", removeFromState: true, wantCalls: 0},
		{name: "disabled_keep_in_state", id: "rm-dbtestinstance:2091052", storeStatus: "Disabled", removeFromState: false, wantErr: true, wantCalls: 0},
		{name: "invalid_id", id: "2091052", storeStatus: "Enabled", wantErr: true, wantCalls: 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var (
				mutex          sync.Mutex
				calls          int
				lastAction     string
				lastDBInstance string
				lastBackupId   string
			)
			client := rdsBackupDeleteTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Error(err)
				}
				mutex.Lock()
				calls++
				lastAction = r.Form.Get("Action")
				lastDBInstance = r.Form.Get("DBInstanceId")
				lastBackupId = r.Form.Get("BackupId")
				mutex.Unlock()
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"RequestId":"test-request-id","BackupId":"2091052"}`))
			})

			d, _ := schema.InternalMap(Provider().ResourcesMap["alicloud_rds_backup"].Schema).Data(nil, nil)
			assert.Nil(t, d.Set("db_instance_id", "rm-dbtestinstance"))
			assert.Nil(t, d.Set("store_status", tc.storeStatus))
			assert.Nil(t, d.Set("remove_from_state", tc.removeFromState))
			d.SetId(tc.id)

			err := rds.ResourceAliCloudRdsBackupDelete(d, client)
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			mutex.Lock()
			defer mutex.Unlock()
			assert.Equal(t, tc.wantCalls, calls)
			if tc.wantCalls > 0 {
				assert.Equal(t, "DeleteBackup", lastAction)
				assert.Equal(t, "rm-dbtestinstance", lastDBInstance)
				assert.Equal(t, "2091052", lastBackupId, "DeleteBackup must receive the bare backup set ID, not the composite resource ID")
			}
		})
	}
}
