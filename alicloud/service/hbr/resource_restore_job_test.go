package hbr_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/hbr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudHBRRestoreJob_basic0(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	acctest.CheckoutAccount(t, true)
	defer acctest.CheckoutAccount(t, false)
	var v map[string]interface{}
	resourceId := "alicloud_hbr_restore_job.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudHBRRestoreJobMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return hbr.NewHbrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeHbrRestoreJob")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrrestorejob%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudHBRRestoreJobBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.HbrSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"snapshot_hash":           "${data.alicloud_hbr_snapshots.nas_snapshots.snapshots.0.snapshot_hash}",
					"vault_id":                "${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}",
					"source_type":             "NAS",
					"restore_type":            "NAS",
					"snapshot_id":             "${data.alicloud_hbr_snapshots.nas_snapshots.snapshots.0.snapshot_id}",
					"target_file_system_id":   "${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}",
					"target_create_time":      "${data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time}",
					"target_path":             "/",
					"options":                 "{\\\"includes\\\":[],\\\"excludes\\\":[]}",
					"include":                 "[\\\"/proc\\\"]",
					"exclude":                 "[\\\"/home\\\", \\\"/var/\\\"]",
					"cross_account_type":      "SELF_ACCOUNT",
					"cross_account_user_id":   "${data.alicloud_account.default.id}",
					"cross_account_role_name": "${alicloud_ram_role.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":             "NAS",
						"restore_type":            "NAS",
						"target_path":             "/",
						"options":                 "{\"includes\":[],\"excludes\":[]}",
						"include":                 "[\"/proc\"]",
						"exclude":                 "[\"/home\", \"/var/\"]",
						"cross_account_type":      "SELF_ACCOUNT",
						"cross_account_user_id":   acctest.CHECKSET,
						"cross_account_role_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"target_client_id":      "cl-tfacc-placeholder",
					"target_data_source_id": "ds-tfacc-placeholder",
					"udm_detail":            "{}",
					"vault_id":              acctest.REMOVEKEY,
				}),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"include", "exclude"},
			},
		},
	})
}

func TestAccAliCloudHBRRestoreJob_basic1(t *testing.T) {
	acctest.CheckoutAccount(t, true)
	defer acctest.CheckoutAccount(t, false)
	var v map[string]interface{}
	resourceId := "alicloud_hbr_restore_job.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudHBRRestoreJobMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return hbr.NewHbrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeHbrRestoreJob")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrrestorejob%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudHBRRestoreJobBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.HbrSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"snapshot_hash": "${data.alicloud_hbr_snapshots.oss_snapshots.snapshots.0.snapshot_hash}",
					"vault_id":      "${data.alicloud_hbr_oss_backup_plans.default.plans.0.vault_id}",
					"source_type":   "OSS",
					"restore_type":  "OSS",
					"snapshot_id":   "${data.alicloud_hbr_snapshots.oss_snapshots.snapshots.0.snapshot_id}",
					"target_bucket": "${data.alicloud_hbr_oss_backup_plans.default.plans.0.bucket}",
					"target_prefix": "",
					"options":       "{\\\"includes\\\":[],\\\"excludes\\\":[]}",
					"include":       "[\\\"/proc\\\"]",
					"exclude":       "[\\\"/home\\\", \\\"/var/\\\"]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":   "OSS",
						"restore_type":  "OSS",
						"target_prefix": "",
						"options":       "{\"includes\":[],\"excludes\":[]}",
						"include":       "[\"/proc\"]",
						"exclude":       "[\"/home\", \"/var/\"]",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"include", "exclude"},
			},
		},
	})
}

func TestAccAliCloudHBRRestoreJob_basic2(t *testing.T) {
	acctest.CheckoutAccount(t, true)
	defer acctest.CheckoutAccount(t, false)
	var v map[string]interface{}
	resourceId := "alicloud_hbr_restore_job.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudHBRRestoreJobMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return hbr.NewHbrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeHbrRestoreJob")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrrestorejob%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudHBRRestoreJobBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.HbrSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"snapshot_hash":      "${data.alicloud_hbr_snapshots.ecs_snapshots.snapshots.0.snapshot_hash}",
					"vault_id":           "${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}",
					"source_type":        "ECS_FILE",
					"restore_type":       "ECS_FILE",
					"snapshot_id":        "${data.alicloud_hbr_snapshots.ecs_snapshots.snapshots.0.snapshot_id}",
					"target_instance_id": "${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}",
					"target_path":        "/",
					"include":            "[\\\"/proc\\\"]",
					"exclude":            "[\\\"/home\\\", \\\"/var/\\\"]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":  "ECS_FILE",
						"restore_type": "ECS_FILE",
						"target_path":  "/",
						"include":      "[\"/proc\"]",
						"exclude":      "[\"/home\", \"/var/\"]",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"include", "exclude"},
			},
		},
	})
}

func TestAccAliCloudHBRRestoreJob_basic3(t *testing.T) {
	acctest.CheckoutAccount(t, true)
	defer acctest.CheckoutAccount(t, false)
	var v map[string]interface{}
	resourceId := "alicloud_hbr_restore_job.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudHBRRestoreJobMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return hbr.NewHbrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeHbrRestoreJob")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrrestorejob%d", acctest.DefaultRegionToTest, rand)
	ecsId := fmt.Sprintf("tf-testacc%d", rand+3)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudHBRRestoreJobBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.HbrSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"restore_job_id":     ecsId,
					"snapshot_hash":      "${data.alicloud_hbr_snapshots.ecs_snapshots.snapshots.0.snapshot_hash}",
					"vault_id":           "${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}",
					"source_type":        "ECS_FILE",
					"restore_type":       "ECS_FILE",
					"snapshot_id":        "${data.alicloud_hbr_snapshots.ecs_snapshots.snapshots.0.snapshot_id}",
					"target_instance_id": "${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}",
					"target_path":        "/",
					"include":            "[\\\"/proc\\\"]",
					"exclude":            "[\\\"/home\\\", \\\"/var/\\\"]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"restore_job_id": ecsId,
						"source_type":    "ECS_FILE",
						"restore_type":   "ECS_FILE",
						"target_path":    "/",
						"include":        "[\"/proc\"]",
						"exclude":        "[\"/home\", \"/var/\"]",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"include", "exclude"},
			},
		},
	})
}

func TestAccAliCloudHBRRestoreJob_basic4(t *testing.T) {
	acctest.CheckoutAccount(t, true)
	defer acctest.CheckoutAccount(t, false)
	var v map[string]interface{}
	resourceId := "alicloud_hbr_restore_job.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudHBRRestoreJobMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return hbr.NewHbrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeHbrRestoreJob")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrrestorejob%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudHBRRestoreJobBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.HbrSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vault_id":             "${data.alicloud_hbr_ots_snapshots.default.snapshots.0.vault_id}",
					"snapshot_hash":        "${data.alicloud_hbr_ots_snapshots.default.snapshots.0.snapshot_hash}",
					"snapshot_id":          "${data.alicloud_hbr_ots_snapshots.default.snapshots.0.snapshot_id}",
					"source_type":          "OTS_TABLE",
					"restore_type":         "OTS_TABLE",
					"include":              "${data.alicloud_hbr_ots_snapshots.default.snapshots.0.table_name}",
					"target_time":          "${data.alicloud_hbr_ots_snapshots.default.snapshots.0.start_time}",
					"target_instance_name": "${data.alicloud_hbr_ots_snapshots.default.snapshots.0.instance_name}",
					"target_table_name":    "${data.alicloud_hbr_ots_snapshots.default.snapshots.0.table_name}",
					"ots_detail": []map[string]interface{}{
						{
							"overwrite_existing": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vault_id":             acctest.CHECKSET,
						"snapshot_hash":        acctest.CHECKSET,
						"snapshot_id":          acctest.CHECKSET,
						"source_type":          "OTS_TABLE",
						"restore_type":         "OTS_TABLE",
						"target_time":          acctest.CHECKSET,
						"target_instance_name": acctest.CHECKSET,
						"target_table_name":    acctest.CHECKSET,
						"ots_detail.#":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ots_detail": []map[string]interface{}{
						{
							"overwrite_existing": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ots_detail.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"include", "exclude"},
			},
		},
	})
}

var AlicloudHBRRestoreJobMap0 = map[string]string{
	"status":        acctest.CHECKSET,
	"exclude":       acctest.NOSET,
	"udm_region_id": acctest.NOSET,
}

func AlicloudHBRRestoreJobBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_account" "default" {
	}

	data "alicloud_hbr_ecs_backup_plans" "default" {
    	name_regex = "plan-tf-used-dont-delete"
	}

	data "alicloud_hbr_oss_backup_plans" "default" {
		name_regex = "plan-tf-used-dont-delete"
	}

	data "alicloud_hbr_nas_backup_plans" "default" {
		name_regex = "plan-tf-used-dont-delete"
	}

	data "alicloud_hbr_snapshots" "ecs_snapshots" {
    	source_type  = "ECS_FILE"
		vault_id     =  data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id
		instance_id  =  data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id
	}

	data "alicloud_hbr_snapshots" "oss_snapshots" {
    	source_type  = "OSS"
		vault_id     =  data.alicloud_hbr_oss_backup_plans.default.plans.0.vault_id
		bucket       =  data.alicloud_hbr_oss_backup_plans.default.plans.0.bucket
	}

	data "alicloud_hbr_snapshots" "nas_snapshots" {
    	source_type     = "NAS"
		vault_id        =  data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id
		file_system_id  =  data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id
    	create_time     =  data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time	
	}

	resource "alicloud_ram_role" "default" {
  		name     = var.name
  		document = <<EOF
		{
			"Statement": [
			{
				"Action": "sts:AssumeRole",
				"Effect": "Allow",
				"Principal": {
					"Service": [
						"crossbackup.hbr.aliyuncs.com"
					]
				}
			}
			],
  			"Version": "1"
		}
  		EOF
  		force    = true
	}
`, name)
}

func AlicloudHBRRestoreJobBasicDependence1(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_hbr_ots_snapshots" "default" {
	}
`, name)
}

// lintignore: R001
func TestUnitAlicloudHBRRestoreJob(t *testing.T) {
	p := acctest.Provider.ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_hbr_restore_job"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_hbr_restore_job"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"snapshot_hash":         "CreateRestoreJobValue",
		"vault_id":              "CreateRestoreJobValue",
		"source_type":           "NAS",
		"restore_type":          "NAS",
		"snapshot_id":           "CreateRestoreJobValue",
		"target_file_system_id": "CreateRestoreJobValue",
		"target_create_time":    "2019-04-04T11:08:33CST",
		"target_path":           "/",
		"options":               "CreateRestoreJobValue",
		"include":               "CreateRestoreJobValue",
		"exclude":               "CreateRestoreJobValue",
		"restore_job_id":        "CreateRestoreJobValue",
		"target_bucket":         "CreateRestoreJobValue",
		"target_client_id":      "CreateRestoreJobValue",
		"target_data_source_id": "CreateRestoreJobValue",
		"target_instance_id":    "CreateRestoreJobValue",
		"target_prefix":         "CreateRestoreJobValue",
		"target_instance_name":  "CreateRestoreJobValue",
		"target_table_name":     "CreateRestoreJobValue",
		"target_time":           "1",
		"udm_detail":            "CreateRestoreJobValue",
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
	rawClient, err := hbrTestClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// DescribeRestoreJobs2
		"RestoreJobs": map[string]interface{}{
			"RestoreJob": []interface{}{
				map[string]interface{}{
					"RestoreId":          "CreateRestoreJobValue",
					"RestoreType":        "NAS",
					"Options":            "CreateRestoreJobValue",
					"SnapshotHash":       "CreateRestoreJobValue",
					"SnapshotId":         "CreateRestoreJobValue",
					"Schedule":           "CreateRestoreJobValue",
					"SourceType":         "NAS",
					"Status":             "RUNNING",
					"TargetBucket":       "CreateRestoreJobValue",
					"TargetClientId":     "CreateRestoreJobValue",
					"TargetCreateTime":   "1554347313",
					"TargetDataSourceId": "CreateRestoreJobValue",
					"TargetFileSystemId": "CreateRestoreJobValue",
					"TargetInstanceId":   "CreateRestoreJobValue",
					"TargetPath":         "/",
					"TargetPrefix":       "CreateRestoreJobValue",
					"VaultId":            "CreateRestoreJobValue",
					"TargetInstanceName": "CreateRestoreJobValue",
					"TargetTableName":    "CreateRestoreJobValue",
					"TargetTime":         "1",
					"UdmDetail":          "CreateRestoreJobValue",
				},
			},
		},
		"RestoreId": "CreateRestoreJobValue",
		"Success":   "true",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateRestoreJob
		"RestoreJobs": map[string]interface{}{
			"RestoreJob": []interface{}{
				map[string]interface{}{
					"RestoreId": "CreateRestoreJobValue",
				},
			},
		},
		"RestoreId": "CreateRestoreJobValue",
		"Success":   "true",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_hbr_restore_job", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewHbrClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = hbr.ResourceRestoreJobCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeRestoreJobs2 Response
		"RestoreJobs": map[string]interface{}{
			"RestoreJob": []interface{}{
				map[string]interface{}{
					"RestoreId": "CreateRestoreJobValue",
				},
			},
		},
		"RestoreId": "CreateRestoreJobValue",
		"Success":   "true",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateRestoreJob" {
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
		err := hbr.ResourceRestoreJobCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_hbr_restore_job"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	attributesDiff := map[string]interface{}{}
	diff, err := acctest.NewInstanceDiff("alicloud_hbr_restore_job", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_hbr_restore_job"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeRestoreJobs2`" {
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
		err := hbr.ResourceRestoreJobRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.Nil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	err = hbr.ResourceRestoreJobDelete(dExisted, rawClient)
	assert.Nil(t, err)
}
