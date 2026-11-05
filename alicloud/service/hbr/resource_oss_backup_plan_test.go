package hbr_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
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

func TestAccAliCloudHbrOssBackupPlan_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_oss_backup_plan.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudHBROssBackupPlanMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return hbr.NewHbrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeHbrOssBackupPlan")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrossbackupplan%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudHBROssBackupPlanBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vault_id":                "${alicloud_hbr_vault.default.id}",
					"bucket":                  "${alicloud_oss_bucket.default.bucket}",
					"backup_type":             "COMPLETE",
					"schedule":                "I|1602673264|PT2H",
					"oss_backup_plan_name":    "tf-testAccHbrOss",
					"retention":               "1",
					"cross_account_type":      "SELF_ACCOUNT",
					"cross_account_user_id":   "${data.alicloud_account.default.id}",
					"cross_account_role_name": "${alicloud_ram_role.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_type":             "COMPLETE",
						"schedule":                "I|1602673264|PT2H",
						"oss_backup_plan_name":    "tf-testAccHbrOss",
						"retention":               "1",
						"cross_account_type":      "SELF_ACCOUNT",
						"cross_account_user_id":   acctest.CHECKSET,
						"cross_account_role_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"oss_backup_plan_name": "tf-testAccHbrOss2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_backup_plan_name": "tf-testAccHbrOss2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": "I|1602673264|P1D",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"schedule": "I|1602673264|P1D",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prefix": "var",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prefix": "var",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prefix": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prefix": "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"oss_backup_plan_name": "tf-testAccHbrOss3",
					"schedule":             "I|1602673264|PT2H",
					"retention":            "3",
					"prefix":               "root",
					"disabled":             "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_backup_plan_name": "tf-testAccHbrOss3",
						"schedule":             "I|1602673264|PT2H",
						"retention":            "3",
						"prefix":               "root",
						"disabled":             "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudHBROssBackupPlanMap0 = map[string]string{
	"path.#":                              acctest.NOSET,
	"disk_id":                             acctest.NOSET,
	"options":                             acctest.NOSET,
	"exclude":                             acctest.NOSET,
	"resource":                            acctest.NOSET,
	"rule":                                acctest.NOSET,
	"udm_region_id":                       acctest.NOSET,
	"speed_limit":                         acctest.NOSET,
	"include":                             acctest.NOSET,
	"detail":                              acctest.NOSET,
	"update_paths":                        acctest.NOSET,
	"instance_id":                         acctest.NOSET,
	"auto_remove_file_system_mount_point": acctest.NOSET,
	"file_system_id":                      acctest.NOSET,
	"create_time":                         acctest.NOSET,
	"bucket":                              acctest.CHECKSET,
	"vault_id":                            acctest.CHECKSET,
	"retention":                           "",
	"schedule":                            "I|1602673264|PT2H",
	"oss_backup_plan_name":                "tf-testAccHbrOss",
	"backup_type":                         "COMPLETE",
}

func AlicloudHBROssBackupPlanBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_account" "default" {
	}

	resource "alicloud_hbr_vault" "default" {
  		vault_name = var.name
	}

	resource "alicloud_oss_bucket" "default" {
  		bucket = var.name
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

// lintignore: R001
func TestUnitAlicloudHBROssBackupPlan(t *testing.T) {
	p := acctest.Provider.ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_hbr_oss_backup_plan"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_hbr_oss_backup_plan"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"backup_type":          "CreateBackupPlanValue",
		"bucket":               "CreateBackupPlanValue",
		"oss_backup_plan_name": "CreateBackupPlanValue",
		"prefix":               "CreateBackupPlanValue",
		"retention":            "1",
		"schedule":             "CreateBackupPlanValue",
		"vault_id":             "CreateBackupPlanValue",
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
		// DescribeBackupPlans
		"BackupPlans": map[string]interface{}{
			"BackupPlan": []interface{}{
				map[string]interface{}{
					"PlanId":     "CreateBackupPlanValue",
					"BackupType": "CreateBackupPlanValue",
					"Bucket":     "CreateBackupPlanValue",
					"PlanName":   "CreateBackupPlanValue",
					"Prefix":     "CreateBackupPlanValue",
					"Retention":  "1",
					"Schedule":   "CreateBackupPlanValue",
					"VaultId":    "CreateBackupPlanValue",
					"Disabled":   false,
				},
			},
		},
		"PlanId":  "CreateBackupPlanValue",
		"Success": "true",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateBackupPlan
		"PlanId": "CreateBackupPlanValue",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_hbr_oss_backup_plan", errorCode))
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
	err = hbr.ResourceOssBackupPlanCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeBackupPlans Response
		"BackupPlans": map[string]interface{}{
			"BackupPlan": []interface{}{
				map[string]interface{}{
					"PlanId": "CreateBackupPlanValue",
				},
			},
		},
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
		err := hbr.ResourceOssBackupPlanCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_hbr_oss_backup_plan"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewHbrClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = hbr.ResourceOssBackupPlanUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateBackupPlan
	attributesDiff := map[string]interface{}{
		"vault_id":             "UpdateBackupPlanValue",
		"prefix":               "UpdateBackupPlanValue",
		"oss_backup_plan_name": "UpdateBackupPlanValue",
		"retention":            "2",
		"schedule":             "UpdateBackupPlanValue",
	}
	diff, err := acctest.NewInstanceDiff("alicloud_hbr_oss_backup_plan", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_hbr_oss_backup_plan"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeBackupPlans Response
		"BackupPlans": map[string]interface{}{
			"BackupPlan": []interface{}{
				map[string]interface{}{
					"Prefix":    "UpdateBackupPlanValue",
					"PlanName":  "UpdateBackupPlanValue",
					"Retention": "2",
					"Schedule":  "UpdateBackupPlanValue",
					"VaultId":   "UpdateBackupPlanValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateBackupPlan" {
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
		err := hbr.ResourceOssBackupPlanUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_hbr_oss_backup_plan"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// EnableBackupPlan
	attributesDiff = map[string]interface{}{
		"disabled": true,
	}
	diff, err = acctest.NewInstanceDiff("alicloud_hbr_oss_backup_plan", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_hbr_oss_backup_plan"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeBackupPlans Response
		"BackupPlans": map[string]interface{}{
			"BackupPlan": []interface{}{
				map[string]interface{}{
					"Disabled": true,
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DisableBackupPlan" {
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
		err := hbr.ResourceOssBackupPlanUpdate(dExisted, rawClient)
		patches.Reset()

		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_hbr_oss_backup_plan"].Schema).Data(dExisted.State(), nil)
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
		err := hbr.ResourceOssBackupPlanRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewHbrClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = hbr.ResourceOssBackupPlanDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
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
						ReadMockResponse = map[string]interface{}{
							"Success": "true",
						}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := hbr.ResourceOssBackupPlanDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
