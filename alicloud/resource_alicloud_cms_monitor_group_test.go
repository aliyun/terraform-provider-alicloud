package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cms_monitor_group", &resource.Sweeper{
		Name: "alicloud_cms_monitor_group",
		F:    testSweepCmsMonitorgroup,
	})
}

func testSweepCmsMonitorgroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeMonitorGroups"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["Type"] = "custom"
	var response map[string]interface{}
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}

	for {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cms_monitor_groups", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Resources.Resource", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Resources.Resource", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := item["GroupName"].(string)
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Cms Monitor Group: %s ", name)
				continue
			}
			log.Printf("[INFO] Delete Cms Monitor Group: %s ", name)

			delAction := "DeleteMonitorGroup"
			conn, err := client.NewCmsClient()
			if err != nil {
				return WrapError(err)
			}
			delRequest := map[string]interface{}{
				"GroupId": fmt.Sprint(formatInt(item["GroupId"])),
			}
			_, err = conn.DoRequest(StringPointer(delAction), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, delRequest, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Cms Monitor Group (%s): %s", name, err)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudCmsMonitorGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_monitor_group.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsMonitorGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsMonitorGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudCmsMonitorGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsMonitorGroupBasicDependence)
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
					"monitor_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitor_group_name": name,
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
					"contact_groups": []string{"${alicloud_cms_alarm_contact_group.default.alarm_contact_group_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitor_group_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitor_group_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance-test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_groups":     []string{"${alicloud_cms_alarm_contact_group.default.alarm_contact_group_name}", "${alicloud_cms_alarm_contact_group.default1.alarm_contact_group_name}"},
					"monitor_group_name": "${var.name}",
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "acceptance-test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups.#":   "2",
						"monitor_group_name": name,
						"tags.%":             "2",
						"tags.Created":       "TF-update",
						"tags.For":           "acceptance-test-update",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCmsMonitorGroup_ByResourceGroupId(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_monitor_group.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsMonitorGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsMonitorGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudCmsMonitorGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsMonitorGroupBasicDependence1)
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
					"contact_groups":      []string{"${alicloud_cms_alarm_contact_group.default.alarm_contact_group_name}"},
					"resource_group_id":   "${alicloud_resource_manager_resource_group.default.id}",
					"resource_group_name": "${alicloud_resource_manager_resource_group.default.resource_group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups.#":    "1",
						"resource_group_id":   CHECKSET,
						"resource_group_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitor_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitor_group_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_group_id", "resource_group_name"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_groups": []string{"${alicloud_cms_alarm_contact_group.default.alarm_contact_group_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitor_group_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitor_group_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance-test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_groups":     []string{"${alicloud_cms_alarm_contact_group.default.alarm_contact_group_name}", "${alicloud_cms_alarm_contact_group.default1.alarm_contact_group_name}"},
					"monitor_group_name": "${var.name}",
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "acceptance-test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups.#":   "2",
						"monitor_group_name": name,
						"tags.%":             "2",
						"tags.Created":       "TF-update",
						"tags.For":           "acceptance-test-update",
					}),
				),
			},
		},
	})
}

var AlicloudCmsMonitorGroupMap = map[string]string{}

func AlicloudCmsMonitorGroupBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_cms_alarm_contact_group" "default" {
alarm_contact_group_name = var.name
}

resource "alicloud_cms_alarm_contact_group" "default1" {
alarm_contact_group_name = "${var.name}_update"
}
`, name)
}

func AlicloudCmsMonitorGroupBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_cms_alarm_contact_group" "default" {
	alarm_contact_group_name = var.name
}

resource "alicloud_cms_alarm_contact_group" "default1" {
	alarm_contact_group_name = "${var.name}_update"
}

resource "alicloud_resource_manager_resource_group" "default" {
	resource_group_name = var.name
	display_name        = var.name
}
`, name)
}

func TestUnitAlicloudCmsMonitorGroup(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_cms_monitor_group"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cms_monitor_group"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"monitor_group_name": "CreateMonitorGroupValue",
		"contact_groups":     []string{"CreateMonitorGroupValue"},
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
	dGroupInit, _ := schema.InternalMap(p["alicloud_cms_monitor_group"].Schema).Data(nil, nil)
	dGroupExisted, _ := schema.InternalMap(p["alicloud_cms_monitor_group"].Schema).Data(nil, nil)
	dGroupInit.MarkNewResource()
	groupAttributes := map[string]interface{}{
		"contact_groups":      []string{"CreateMonitorGroupValue"},
		"resource_group_name": "CreateMonitorGroupValue",
		"resource_group_id":   "CreateMonitorGroupValue",
	}
	for key, value := range groupAttributes {
		err := dGroupInit.Set(key, value)
		assert.Nil(t, err)
		err = dGroupExisted.Set(key, value)
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
		// DescribeMonitorGroups
		"Resources": map[string]interface{}{
			"Resource": []interface{}{
				map[string]interface{}{
					"GroupName": "CreateMonitorGroupValue",
					"GroupId":   1,
					"ContactGroups": map[string]interface{}{
						"ContactGroup": []interface{}{
							map[string]interface{}{
								"Name": "CreateMonitorGroupValue",
							},
						},
					},
					"Tags": map[string]interface{}{
						"Tag": []interface{}{
							map[string]interface{}{
								"Key":   "CreateMonitorGroupValue",
								"Value": "CreateMonitorGroupValue",
							},
						},
					},
				},
			},
		},
		"Success": true,
		"GroupId": 1,
		"Code":    "200",
		"Id":      "CreateMonitorGroupValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateMonitorGroup
		"Id":      "CreateMonitorGroupValue",
		"GroupId": 1,
		"Code":    "200",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cms_monitor_group", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudCmsMonitorGroupCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeMonitorGroups Response
		"GroupId": 1,
		"Code":    "200",
		"Id":      "CreateMonitorGroupValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateMonitorGroup" {
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
		err := resourceAlicloudCmsMonitorGroupCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cms_monitor_group"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	ReadMockResponseDiff = map[string]interface{}{
		// DescribeMonitorGroups Response
		"GroupId": 1,
		"Code":    "200",
		"Id":      "CreateMonitorGroupValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateMonitorGroupByResourceGroupId" {
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
		err := resourceAlicloudCmsMonitorGroupCreate(dGroupInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cms_monitor_group"].Schema).Data(dGroupInit.State(), nil)
			for key, value := range groupAttributes {
				_ = dCompare.Set(key, value)
			}
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudCmsMonitorGroupUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyMonitorGroup
	attributesDiff := map[string]interface{}{
		"monitor_group_name": "ModifyMonitorGroupValue",
		"contact_groups":     []string{"ModifyMonitorGroupValue"},
	}
	diff, err := newInstanceDiff("alicloud_cms_monitor_group", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cms_monitor_group"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeMonitorGroups Response
		"Resources": map[string]interface{}{
			"Resource": []interface{}{
				map[string]interface{}{
					"GroupName": "ModifyMonitorGroupValue",
					"ContactGroups": map[string]interface{}{
						"ContactGroup": []interface{}{
							map[string]interface{}{
								"Name": "ModifyMonitorGroupValue",
							},
						},
					},
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyMonitorGroup" {
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
		err := resourceAlicloudCmsMonitorGroupUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cms_monitor_group"].Schema).Data(dExisted.State(), nil)
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
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeMonitorGroups" {
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
		err := resourceAlicloudCmsMonitorGroupRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudCmsMonitorGroupDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteMonitorGroup" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{
							"Code": "ResourceNotFound",
						}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudCmsMonitorGroupDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
