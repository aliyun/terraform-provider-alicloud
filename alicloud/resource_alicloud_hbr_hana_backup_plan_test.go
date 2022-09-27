package alicloud

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestAccAlicloudHBRHanaBackupPlan_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_hana_backup_plan.default"
	checkoutSupportedRegions(t, true, connectivity.HBRSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudHBRHanaBackupPlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrHanaBackupPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrhanabackupplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHBRHanaBackupPlanBasicDependence0)
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
					"vault_id":          "${alicloud_hbr_vault.default.id}",
					"cluster_id":        "${alicloud_hbr_hana_instance.default.hana_instance_id}",
					"plan_name":         "${var.name}",
					"backup_prefix":     "DIFF_DATA_BACKUP",
					"schedule":          "I|1602673264|P1D",
					"database_name":     "SYSTEMDB",
					"backup_type":       "COMPLETE",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vault_id":          CHECKSET,
						"cluster_id":        CHECKSET,
						"plan_name":         name,
						"backup_prefix":     "DIFF_DATA_BACKUP",
						"schedule":          "I|1602673264|P1D",
						"database_name":     "SYSTEMDB",
						"backup_type":       "COMPLETE",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_group_id"},
			},
		},
	})
}

var AlicloudHBRHanaBackupPlanMap0 = map[string]string{
	"cluster_id": CHECKSET,
	"vault_id":   CHECKSET,
	"plan_id":    CHECKSET,
}

func AlicloudHBRHanaBackupPlanBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
}
data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_hbr_hana_instance" "default" {
	alert_setting = "INHERITED"
	hana_name = var.name
	host = "1.1.1.1"
	instance_number = "1"
	password = "YouPassword123"
	resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
	sid = "HXE"
	use_ssl = "false"
	user_name = "admin"
	validate_certificate = "false"
	vault_id = alicloud_hbr_vault.default.id
}
`, name)
}

func TestAccAlicloudHBRHanaBackupPlan_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_hana_backup_plan.default"
	checkoutSupportedRegions(t, true, connectivity.HBRSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudHBRHanaBackupPlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HbrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHbrHanaBackupPlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrhanabackupplan%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudHBRHanaBackupPlanBasicDependence0)
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
					"vault_id":      "${alicloud_hbr_vault.default.id}",
					"cluster_id":    "${alicloud_hbr_hana_instance.default.hana_instance_id}",
					"plan_name":     "${var.name}",
					"schedule":      "I|1602673264|P1D",
					"database_name": "SYSTEMDB",
					"backup_type":   "COMPLETE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vault_id":      CHECKSET,
						"cluster_id":    CHECKSET,
						"plan_name":     name,
						"schedule":      "I|1602673264|P1D",
						"database_name": "SYSTEMDB",
						"backup_type":   "COMPLETE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_prefix": "DIFF_DATA_BACKUP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_prefix": "DIFF_DATA_BACKUP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": "I|1602673265|P1D",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"schedule": "I|1602673265|P1D",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_name":     "${var.name}",
					"backup_prefix": "DIFF_DATA_BACKUP1",
					"schedule":      "I|1602673264|P1D",
					"status":        "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_name":     name,
						"backup_prefix": "DIFF_DATA_BACKUP1",
						"schedule":      "I|1602673264|P1D",
						"status":        "Enabled",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_group_id"},
			},
		},
	})
}

func TestUnitAccAlicloudHbrHanaBackupPlan(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_hbr_hana_backup_plan"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_hbr_hana_backup_plan"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"vault_id":          "CreateHbrHanaBackupPlanValue",
		"cluster_id":        "CreateHbrHanaBackupPlanValue",
		"plan_name":         "CreateHbrHanaBackupPlanValue",
		"backup_prefix":     "CreateHbrHanaBackupPlanValue",
		"schedule":          "CreateHbrHanaBackupPlanValue",
		"database_name":     "CreateHbrHanaBackupPlanValue",
		"backup_type":       "CreateHbrHanaBackupPlanValue",
		"resource_group_id": "CreateHbrHanaBackupPlanValue",
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
		"HanaBackupPlans": map[string]interface{}{
			"HanaBackupPlan": []interface{}{
				map[string]interface{}{
					"ClusterId":    "CreateHbrHanaBackupPlanValue",
					"BackupPrefix": "CreateHbrHanaBackupPlanValue",
					"Schedule":     "CreateHbrHanaBackupPlanValue",
					"PlanId":       "CreateHbrHanaBackupPlanValue",
					"DatabaseName": "CreateHbrHanaBackupPlanValue",
					"BackupType":   "CreateHbrHanaBackupPlanValue",
					"VaultId":      "CreateHbrHanaBackupPlanValue",
					"PlanName":     "CreateHbrHanaBackupPlanValue",
					"Disabled":     false,
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		"PlanId": "CreateHbrHanaBackupPlanValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_hbr_hana_backup_plan", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewHbrClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudHbrHanaBackupPlanCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateHanaBackupPlan" {
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
		err := resourceAlicloudHbrHanaBackupPlanCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_hbr_hana_backup_plan"].Schema).Data(dInit.State(), nil)
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
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudHbrHanaBackupPlanUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"plan_name":     "UpdateHbrHanaBackupPlanValue",
		"backup_prefix": "UpdateHbrHanaBackupPlanValue",
		"schedule":      "UpdateHbrHanaBackupPlanValue",
		"database_name": "UpdateHbrHanaBackupPlanValue",
	}
	diff, err := newInstanceDiff("alicloud_hbr_hana_backup_plan", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_hbr_hana_backup_plan"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"HanaBackupPlans": map[string]interface{}{
			"HanaBackupPlan": []interface{}{
				map[string]interface{}{
					"BackupPrefix": "UpdateHbrHanaBackupPlanValue",
					"Schedule":     "UpdateHbrHanaBackupPlanValue",
					"DatabaseName": "UpdateHbrHanaBackupPlanValue",
					"PlanName":     "UpdateHbrHanaBackupPlanValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateHanaBackupPlan" {
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
		err := resourceAlicloudHbrHanaBackupPlanUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_hbr_hana_backup_plan"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Disabled
	attributesDiff = map[string]interface{}{
		"status": "Disabled",
	}
	diff, err = newInstanceDiff("alicloud_hbr_hana_backup_plan", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_hbr_hana_backup_plan"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"HanaBackupPlans": map[string]interface{}{
			"HanaBackupPlan": []interface{}{
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
			if *action == "DisableHanaBackupPlan" {
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
		err := resourceAlicloudHbrHanaBackupPlanUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_hbr_hana_backup_plan"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Enabled
	attributesDiff = map[string]interface{}{
		"status": "Enabled",
	}
	diff, err = newInstanceDiff("alicloud_hbr_hana_backup_plan", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_hbr_hana_backup_plan"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"HanaBackupPlans": map[string]interface{}{
			"HanaBackupPlan": []interface{}{
				map[string]interface{}{
					"Disabled": false,
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "EnableHanaBackupPlan" {
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
		err := resourceAlicloudHbrHanaBackupPlanUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_hbr_hana_backup_plan"].Schema).Data(dExisted.State(), nil)
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
			if *action == "DescribeHanaBackupPlans" {
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
		err := resourceAlicloudHbrHanaBackupPlanRead(dExisted, rawClient)
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
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudHbrHanaBackupPlanDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_hbr_hana_backup_plan", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_hbr_hana_backup_plan"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteHanaBackupPlan" {
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
		err := resourceAlicloudHbrHanaBackupPlanDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
