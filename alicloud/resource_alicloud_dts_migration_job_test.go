package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_dts_migration_job",
		&resource.Sweeper{
			Name: "alicloud_dts_migration_job",
			F:    testSweepDTSMigrationJob,
		})
}

func testSweepDTSMigrationJob(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeDtsJobs"
	request := map[string]interface{}{}
	request["JobType"] = "MIGRATION"
	request["PageNumber"] = 1
	request["MaxResults"] = PageSizeXLarge

	var response map[string]interface{}
	conn, err := client.NewDtsClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

		resp, err := jsonpath.Get("$.DtsJobList", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.DtsJobList", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["DtsJobName"]; !ok {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["DtsJobName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping DTS Migration Job: %s", item["DtsJobName"].(string))
				continue
			}
			action := "DeleteDtsJob"
			request := map[string]interface{}{
				"DtsJobId": item["DtsJobId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete DTS Migration Job (%s): %s", item["DtsJobId"].(string), err)
			}
			log.Printf("[INFO] Delete DTS Migration Job success: %s ", item["DtsJobId"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudDTSMigrationJob_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_migration_job.default"
	checkoutSupportedRegions(t, true, connectivity.DTSSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudDTSMigrationJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsMigrationJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtsmigrationjob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDTSMigrationJobBasicDependence0)
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
					"dts_instance_id":                    "${alicloud_dts_migration_instance.default.id}",
					"dts_job_name":                       name,
					"source_endpoint_instance_type":      "RDS",
					"source_endpoint_instance_id":        "${alicloud_db_instance.default.0.id}",
					"source_endpoint_engine_name":        "MySQL",
					"source_endpoint_region":             "${var.region}",
					"source_endpoint_user_name":          "${alicloud_rds_account.default.0.name}",
					"source_endpoint_password":           "${var.password}",
					"destination_endpoint_instance_type": "RDS",
					"destination_endpoint_instance_id":   "${alicloud_db_instance.default.1.id}",
					"destination_endpoint_engine_name":   "MySQL",
					"destination_endpoint_region":        "${var.region}",
					"destination_endpoint_user_name":     "${alicloud_rds_account.default.1.name}",
					"destination_endpoint_password":      "${var.password}",
					"db_list":                            `{\"tftestdatabase\":{\"name\":\"tftestdatabase\",\"all\":true}}`,
					"structure_initialization":           "true",
					"data_initialization":                "true",
					"data_synchronization":               "true",
					"status":                             "Migrating",
					"depends_on":                         []string{"alicloud_db_account_privilege.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_instance_id":                    CHECKSET,
						"dts_job_name":                       name,
						"source_endpoint_instance_type":      "RDS",
						"source_endpoint_engine_name":        "MySQL",
						"source_endpoint_region":             CHECKSET,
						"source_endpoint_user_name":          CHECKSET,
						"destination_endpoint_instance_type": "RDS",
						"destination_endpoint_engine_name":   "MySQL",
						"destination_endpoint_region":        CHECKSET,
						"destination_endpoint_user_name":     CHECKSET,
						"db_list":                            CHECKSET,
						"structure_initialization":           "true",
						"data_initialization":                "true",
						"data_synchronization":               "true",
						"status":                             "Migrating",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Suspending",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Suspending",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Migrating",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Migrating",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"destination_endpoint_password", "source_endpoint_password"},
			},
		},
	})
}

var AlicloudDTSMigrationJobMap0 = map[string]string{}

func AlicloudDTSMigrationJobBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "region" {
  default = "%s"
}

variable "password" {
  default = "Test12345"
}

variable "database_name" {
  default = "tftestdatabase"
}

data "alicloud_db_zones" "default" {}

data "alicloud_db_instance_classes" "default" {
  engine               = "MySQL"
  engine_version       = "5.6"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_db_zones.default.zones[0].id
}

resource "alicloud_db_instance" "default" {
  count            = 2
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    =  data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
  instance_storage = "10"
  vswitch_id       = data.alicloud_vswitches.default.ids[0]
  instance_name    = join("", [var.name, count.index])
}

resource "alicloud_rds_account" "default" {
  count            = 2
  db_instance_id   = alicloud_db_instance.default[count.index].id
  account_name     = join("", [var.database_name, count.index])
  account_password = var.password
}

resource "alicloud_db_database" "default" {
  count       = 2
  instance_id = alicloud_db_instance.default[count.index].id
  name        = var.database_name
}

resource "alicloud_db_account_privilege" "default" {
  count        = 2
  instance_id  = alicloud_db_instance.default[count.index].id
  account_name = alicloud_rds_account.default[count.index].name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.default[count.index].name]
}

resource "alicloud_dts_migration_instance" "default" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "MySQL"
  source_endpoint_region           = var.region
  destination_endpoint_engine_name = "MySQL"
  destination_endpoint_region      = var.region
  instance_class                   = "small"
  sync_architecture                = "oneway"
}
`, name, defaultRegionToTest)
}

func TestUnitAlicloudDTSMigrationJob(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_dts_migration_job"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_dts_migration_job"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"dts_instance_id":                    "dts_instance_id",
		"dts_job_name":                       "dts_job_name",
		"source_endpoint_password":           "source_endpoint_password",
		"source_endpoint_database_name":      "source_endpoint_database_name",
		"source_endpoint_engine_name":        "MySQL",
		"source_endpoint_ip":                 "source_endpoint_ip",
		"source_endpoint_instance_id":        "source_endpoint_instance_id",
		"source_endpoint_instance_type":      "RDS",
		"source_endpoint_oracle_sid":         "source_endpoint_oracle_sid",
		"source_endpoint_port":               "source_endpoint_port",
		"source_endpoint_region":             "source_endpoint_region",
		"source_endpoint_user_name":          "source_endpoint_user_name",
		"source_endpoint_role":               "source_endpoint_role",
		"source_endpoint_owner_id":           "source_endpoint_owner_id",
		"destination_endpoint_database_name": "destination_endpoint_database_name",
		"destination_endpoint_engine_name":   "MySQL",
		"destination_endpoint_ip":            "destination_endpoint_ip",
		"destination_endpoint_instance_id":   "destination_endpoint_instance_id",
		"destination_endpoint_instance_type": "RDS",
		"destination_endpoint_oracle_sid":    "destination_endpoint_oracle_sid",
		"destination_endpoint_port":          "destination_endpoint_port",
		"destination_endpoint_region":        "destination_endpoint_region",
		"destination_endpoint_user_name":     "destination_endpoint_user_name",
		"destination_endpoint_password":      "destination_endpoint_password",
		"db_list":                            "db_list",
		"structure_initialization":           true,
		"data_initialization":                true,
		"data_synchronization":               true,
		"status":                             "Migrating",
		"checkpoint":                         "checkpoint",
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
		"DtsJobId": "dts_job_id",
		"Status":   "Migrating",
		"DtsJobList": []interface{}{
			map[string]interface{}{
				"Checkpoint": "checkpoint",
				"MigrationMode": map[string]interface{}{
					"DataInitialization":       true,
					"DataSynchronization":      true,
					"structure_initialization": true,
				},
				"DbObject": "db_list",
				"DestinationEndpoint": map[string]interface{}{
					"DatabaseName": "destination_endpoint_database_name",
					"EngineName":   "MySQL",
					"Ip":           "destination_endpoint_ip",
					"InstanceID":   "destination_endpoint_instance_id",
					"InstanceType": "RDS",
					"OracleSID":    "destination_endpoint_oracle_sid",
					"Port":         "destination_endpoint_port",
					"Region":       "destination_endpoint_region",
					"UserName":     "destination_endpoint_user_name",
				},
				"DtsInstanceID": "dts_instance_id",
				"DtsJobName":    "dts_job_name",
				"PayType":       "PostPaid",
				"SourceEndpoint": map[string]interface{}{
					"DatabaseName": "source_endpoint_database_name",
					"EngineName":   "MySQL",
					"Ip":           "source_endpoint_ip",
					"InstanceID":   "source_endpoint_instance_id",
					"InstanceType": "RDS",
					"OracleSID":    "source_endpoint_oracle_sid",
					"Port":         "source_endpoint_port",
					"Region":       "source_endpoint_region",
					"UserName":     "source_endpoint_user_name",
					"RoleName":     "source_endpoint_role",
					"AliyunUid":    "source_endpoint_owner_id",
				},
				"Status":   "Migrating",
				"DtsJobId": "dts_job_id",
			},
		},
	}

	SuspendingReadMockResponse := map[string]interface{}{
		"DtsJobId": "dts_job_id",
		"Status":   "Suspending",
		"DtsJobList": []interface{}{
			map[string]interface{}{
				"Checkpoint": "checkpoint",
				"MigrationMode": map[string]interface{}{
					"DataInitialization":       true,
					"DataSynchronization":      true,
					"structure_initialization": true,
				},
				"DbObject": "db_list",
				"DestinationEndpoint": map[string]interface{}{
					"DatabaseName": "destination_endpoint_database_name",
					"EngineName":   "MySQL",
					"Ip":           "destination_endpoint_ip",
					"InstanceID":   "destination_endpoint_instance_id",
					"InstanceType": "RDS",
					"OracleSID":    "destination_endpoint_oracle_sid",
					"Port":         "destination_endpoint_port",
					"Region":       "destination_endpoint_region",
					"UserName":     "destination_endpoint_user_name",
				},
				"DtsInstanceID": "dts_instance_id",
				"DtsJobName":    "dts_job_name",
				"PayType":       "PostPaid",
				"SourceEndpoint": map[string]interface{}{
					"DatabaseName": "source_endpoint_database_name",
					"EngineName":   "MySQL",
					"Ip":           "source_endpoint_ip",
					"InstanceID":   "source_endpoint_instance_id",
					"InstanceType": "RDS",
					"OracleSID":    "source_endpoint_oracle_sid",
					"Port":         "source_endpoint_port",
					"Region":       "source_endpoint_region",
					"UserName":     "source_endpoint_user_name",
					"RoleName":     "source_endpoint_role",
					"AliyunUid":    "source_endpoint_owner_id",
				},
				"Status":   "Suspending",
				"DtsJobId": "dts_job_id",
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
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_dts_migration_job", "dts_instance_id"))
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
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"UpdateSuspendingNormal": func(errorCode string) (map[string]interface{}, error) {
			result := SuspendingReadMockResponse
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDtsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudDtsMigrationJobCreate(d, rawClient)
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
		err := resourceAlicloudDtsMigrationJobCreate(d, rawClient)
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
		err := resourceAlicloudDtsMigrationJobCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("CreateRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudDtsMigrationJobCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDtsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAlicloudDtsMigrationJobUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("dts_job_id")

	t.Run("UpdateModifyDtsMigrationJobSuspendingStatusAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.SetAttribute("status", &terraform.ResourceAttrDiff{Old: "Migrating", New: "Suspending"})
		resourceData1, _ := schema.InternalMap(p["alicloud_dts_migration_job"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patcheDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DtsService{}), "DescribeDtsMigrationJob", func(*DtsService, string) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateSuspendingNormal"]("")
		})
		err := resourceAlicloudDtsMigrationJobUpdate(resourceData1, rawClient)
		patcheDescribe.Reset()
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyDTSMigrationJobSuspendingStatusNoRetryErrorAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.SetAttribute("status", &terraform.ResourceAttrDiff{Old: "Migrating", New: "Suspending"})
		resourceData, _ := schema.InternalMap(p["alicloud_dts_migration_job"].Schema).Data(nil, diff)
		resourceData.SetId(d.Id())
		patcheDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DtsService{}), "DescribeDtsMigrationJob", func(*DtsService, string) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("NoRetryError")
		})
		err := resourceAlicloudDtsMigrationJobUpdate(resourceData, rawClient)
		patcheDescribe.Reset()
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyDTSMigrationJobSuspendingStatusNoRetryErrorNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.SetAttribute("status", &terraform.ResourceAttrDiff{Old: "Migrating", New: "Suspending"})
		resourceData, _ := schema.InternalMap(p["alicloud_dts_migration_job"].Schema).Data(nil, diff)
		resourceData.SetId(d.Id())
		patcheDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DtsService{}), "DescribeDtsMigrationJob", func(*DtsService, string) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["UpdateSuspendingNormal"]("")
		})
		err := resourceAlicloudDtsMigrationJobUpdate(resourceData, rawClient)
		patcheDescribe.Reset()
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("UpdateModifyDtsMigrationJobMigratingStatusAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.SetAttribute("status", &terraform.ResourceAttrDiff{Old: "Suspending", New: "Migrating"})
		resourceData1, _ := schema.InternalMap(p["alicloud_dts_migration_job"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patcheDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DtsService{}), "DescribeDtsMigrationJob", func(*DtsService, string) (map[string]interface{}, error) {
			return responseMock["UpdateSuspendingNormal"]("")
		})
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAlicloudDtsMigrationJobUpdate(resourceData1, rawClient)
		patcheDescribe.Reset()
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyDTSMigrationJobMigratingStatusNoRetryErrorAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.SetAttribute("status", &terraform.ResourceAttrDiff{Old: "Suspending", New: "Migrating"})
		resourceData, _ := schema.InternalMap(p["alicloud_dts_migration_job"].Schema).Data(nil, diff)
		resourceData.SetId(d.Id())
		patcheDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DtsService{}), "DescribeDtsMigrationJob", func(*DtsService, string) (map[string]interface{}, error) {
			return responseMock["UpdateSuspendingNormal"]("")
		})
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("NoRetryError")
		})
		err := resourceAlicloudDtsMigrationJobUpdate(resourceData, rawClient)
		patcheDescribe.Reset()
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyDTSMigrationJobMigratingStatusNoRetryErrorNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.SetAttribute("status", &terraform.ResourceAttrDiff{Old: "Suspending", New: "Migrating"})
		resourceData, _ := schema.InternalMap(p["alicloud_dts_migration_job"].Schema).Data(nil, diff)
		resourceData.SetId(d.Id())
		patcheDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DtsService{}), "DescribeDtsMigrationJob", func(*DtsService, string) (map[string]interface{}, error) {
			return responseMock["UpdateSuspendingNormal"]("")
		})
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAlicloudDtsMigrationJobUpdate(resourceData, rawClient)
		patcheDescribe.Reset()
		patches.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDtsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudDtsMigrationJobDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				// retry until the timeout comes
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudDtsMigrationJobDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudDtsMigrationJobDelete(d, rawClient)
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
		err := resourceAlicloudDtsMigrationJobDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteDescribeVpcDTSMigrationJobExpectedError", func(t *testing.T) {
		patchRequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("Forbidden.InstanceNotFound")
		})
		err := resourceAlicloudDtsMigrationJobDelete(d, rawClient)
		patchRequest.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeVpcDTSMigrationJobNotFound", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		resourceData, _ := schema.InternalMap(p["alicloud_dts_migration_job"].Schema).Data(nil, diff)
		resourceData.SetId(d.Id())
		patcheDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DtsService{}), "DescribeDtsMigrationJob", func(*DtsService, string) (map[string]interface{}, error) {
			return responseMock["NotFoundError"]("NotFoundError")
		})
		err := resourceAlicloudDtsMigrationJobRead(resourceData, rawClient)
		patcheDescribe.Reset()
		assert.Nil(t, err)
	})
	t.Run("ReadDescribeDTSMigrationJobAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudDtsMigrationJobDelete(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})
}
