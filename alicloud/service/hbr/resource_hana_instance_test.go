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

// collection-order: bypass ecs_instance_ids -- pre-existing attribute outside this change

func TestAccAliCloudHbrHanaInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_hana_instance.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.HBRSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudHBRHanaInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return hbr.NewHbrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeHbrHanaInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrhanainstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudHBRHanaInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_number":      "1",
					"user_name":            "admin",
					"host":                 "1.1.1.1",
					"use_ssl":              "true",
					"vault_id":             "${alicloud_hbr_vault.default.id}",
					"hana_name":            "${var.name}",
					"alert_setting":        "INHERITED",
					"validate_certificate": "false",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"sid":                  "HXE",
					"password":             "YouPassword123",
					"ecs_instance_ids":     []string{"${data.alicloud_instances.default.ids.0}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_number":      "1",
						"user_name":            "admin",
						"host":                 "1.1.1.1",
						"use_ssl":              "true",
						"vault_id":             acctest.CHECKSET,
						"hana_name":            name,
						"alert_setting":        "INHERITED",
						"validate_certificate": "false",
						"resource_group_id":    acctest.CHECKSET,
						"sid":                  "HXE",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "sid", "ecs_instance_ids"},
			},
		},
	})
}

var AlicloudHBRHanaInstanceMap0 = map[string]string{
	"hana_instance_id": acctest.CHECKSET,
	"vault_id":         acctest.CHECKSET,
	"status":           acctest.CHECKSET,
}

func AlicloudHBRHanaInstanceBasicDependence0(name string) string {
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
data "alicloud_instances" "default" {
  status     = "Running"
}
`, name)
}

func TestAccAliCloudHbrHanaInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbr_hana_instance.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.HBRSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudHBRHanaInstanceMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return hbr.NewHbrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeHbrHanaInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%shbrhanainstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudHBRHanaInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vault_id":  "${alicloud_hbr_vault.default.id}",
					"host":      "1.1.1.1",
					"hana_name": "${var.name}",
					"sid":       "HXE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vault_id":  acctest.CHECKSET,
						"host":      "1.1.1.1",
						"hana_name": name,
						"sid":       "HXE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"instance_number":   "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_number": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_name": "admin",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_name": "admin",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host": "1.1.1.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host": "1.1.1.2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"use_ssl": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"use_ssl": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hana_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hana_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_setting": "INHERITED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_setting": "INHERITED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"validate_certificate": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"validate_certificate": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_number":      "1",
					"user_name":            "admin1",
					"host":                 "1.1.1.1",
					"use_ssl":              "true",
					"hana_name":            "${var.name}",
					"alert_setting":        "INHERITED",
					"validate_certificate": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_number":      "1",
						"user_name":            "admin1",
						"host":                 "1.1.1.1",
						"use_ssl":              "true",
						"hana_name":            name,
						"alert_setting":        "INHERITED",
						"validate_certificate": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "sid"},
			},
		},
	})
}

var AlicloudHBRHanaInstanceMap1 = map[string]string{
	"status":           acctest.CHECKSET,
	"hana_instance_id": acctest.CHECKSET,
	"vault_id":         acctest.CHECKSET,
}

// lintignore: R001
func TestUnitAccAlicloudHbrHanaInstance(t *testing.T) {
	p := acctest.Provider.ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_hbr_hana_instance"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_hbr_hana_instance"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"instance_number":      1,
		"user_name":            "CreateHbrHanaInstanceValue",
		"host":                 "CreateHbrHanaInstanceValue",
		"use_ssl":              true,
		"vault_id":             "CreateHbrHanaInstanceValue",
		"hana_name":            "CreateHbrHanaInstanceValue",
		"alert_setting":        "CreateHbrHanaInstanceValue",
		"validate_certificate": true,
		"resource_group_id":    "CreateHbrHanaInstanceValue",
		"sid":                  "CreateHbrHanaInstanceValue",
		"password":             "CreateHbrHanaInstanceValue",
		"ecs_instance_ids":     []string{"CreateHbrHanaInstanceValue"},
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
		"Hanas": map[string]interface{}{
			"Hana": []interface{}{
				map[string]interface{}{
					"Status":              0,
					"Host":                "CreateHbrHanaInstanceValue",
					"VaultId":             "CreateHbrHanaInstanceValue",
					"UseSsl":              true,
					"HanaName":            "CreateHbrHanaInstanceValue",
					"InstanceNumber":      1,
					"ValidateCertificate": true,
					"AlertSetting":        "CreateHbrHanaInstanceValue",
					"UserName":            "CreateHbrHanaInstanceValue",
					"StatusMessage":       "CreateHbrHanaInstanceValue",
					"ClusterId":           "HbrHanaInstanceId",
					"ResourceGroupId":     "CreateHbrHanaInstanceValue",
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		"ClusterId": "HbrHanaInstanceId",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_hbr_hana_instance", errorCode))
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
	err = hbr.ResourceHanaInstanceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateHanaInstance" {
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
		err := hbr.ResourceHanaInstanceCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_hbr_hana_instance"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
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
	err = hbr.ResourceHanaInstanceUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"instance_number":      2,
		"user_name":            "UpdateHbrHanaInstanceValue",
		"host":                 "UpdateHbrHanaInstanceValue",
		"use_ssl":              false,
		"hana_name":            "UpdateHbrHanaInstanceValue",
		"alert_setting":        "UpdateHbrHanaInstanceValue",
		"validate_certificate": false,
		"resource_group_id":    "UpdateHbrHanaInstanceValue",
	}
	diff, err := acctest.NewInstanceDiff("alicloud_hbr_hana_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_hbr_hana_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"Hanas": map[string]interface{}{
			"Hana": []interface{}{
				map[string]interface{}{
					"Host":                "UpdateHbrHanaInstanceValue",
					"UseSsl":              false,
					"HanaName":            "UpdateHbrHanaInstanceValue",
					"InstanceNumber":      2,
					"ValidateCertificate": false,
					"AlertSetting":        "UpdateHbrHanaInstanceValue",
					"UserName":            "UpdateHbrHanaInstanceValue",
					"StatusMessage":       "UpdateHbrHanaInstanceValue",
					"ResourceGroupId":     "UpdateHbrHanaInstanceValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateHanaInstance" {
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
		err := hbr.ResourceHanaInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_hbr_hana_instance"].Schema).Data(dExisted.State(), nil)
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
	diff, err = acctest.NewInstanceDiff("alicloud_hbr_hana_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_hbr_hana_instance"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeHanaInstances" {
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
		err := hbr.ResourceHanaInstanceRead(dExisted, rawClient)
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
	err = hbr.ResourceHanaInstanceDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = acctest.NewInstanceDiff("alicloud_hbr_hana_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_hbr_hana_instance"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteHanaInstance" {
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
		err := hbr.ResourceHanaInstanceDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
