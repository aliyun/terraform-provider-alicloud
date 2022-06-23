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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudSchedulerxNamespace_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_schedulerx_namespace.default"
	ra := resourceAttrInit(resourceId, AlicloudSchedulerxNamespaceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Schedulerx2Service{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSchedulerxNamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sschedulerxnamespace%d", defaultRegionToTest, rand)
	checkoutSupportedRegions(t, true, connectivity.SchedulerxSupportRegions)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSchedulerxNamespaceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_name": name,
					"description":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name": name,
						"description":    name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_name": name,
					"description":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name": name,
						"description":    name,
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

var AlicloudSchedulerxNamespaceMap0 = map[string]string{}

func AlicloudSchedulerxNamespaceBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

func TestUnitAccAlicloudSchedulerxNamespace(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_schedulerx_namespace"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_schedulerx_namespace"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"description":    "CreateSchedulerxNamespaceValue",
		"namespace_name": "CreateSchedulerxNamespaceValue",
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
		"Data": map[string]interface{}{
			"Namespaces": []interface{}{
				map[string]interface{}{
					"UId":         "CreateSchedulerxNamespaceValue",
					"Description": "CreateSchedulerxNamespaceValue",
					"Name":        "CreateSchedulerxNamespaceValue",
				},
			},
		},
	}

	CreateMockResponse := map[string]interface{}{
		"Data": map[string]interface{}{
			"NamespaceUid": "CreateSchedulerxNamespaceValue",
		},
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_schedulerx_namespace", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEdasschedulerxClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudSchedulerxNamespaceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateNamespace" {
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
		err := resourceAlicloudSchedulerxNamespaceCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_schedulerx_namespace"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEdasschedulerxClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudSchedulerxNamespaceUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"description":    "UpdateSchedulerxNamespaceValue",
		"namespace_name": "UpdateSchedulerxNamespaceValue",
	}
	diff, err := newInstanceDiff("alicloud_schedulerx_namespace", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_schedulerx_namespace"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"Data": map[string]interface{}{
			"Namespaces": []interface{}{
				map[string]interface{}{
					"UId":         "CreateSchedulerxNamespaceValue",
					"Description": "UpdateSchedulerxNamespaceValue",
					"Name":        "UpdateSchedulerxNamespaceValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateNamespace" {
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
		err := resourceAlicloudSchedulerxNamespaceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_schedulerx_namespace"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_schedulerx_namespace", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_schedulerx_namespace"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ListNamespaces" {
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
		err := resourceAlicloudSchedulerxNamespaceRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

}
