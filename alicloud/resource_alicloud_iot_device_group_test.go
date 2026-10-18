package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_iot_device_group", &resource.Sweeper{
		Name: "alicloud_iot_device_group",
		F:    testSweepDeviceGroup,
	})
}

func testSweepDeviceGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	listAction := "QueryDeviceGroupList"
	currentPage := 1
	for {
		request := map[string]interface{}{
			"CurrentPage": currentPage,
			"PageSize":    PageSizeSmall,
		}
		response, err := client.RpcPost("Iot", "2018-01-20", listAction, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] %s got an error: %#v", listAction, err)
			return nil
		}
		addDebug(listAction, response, request)

		// Handle a missing Data.GroupInfo key as an empty list — the IoT
		// QueryDeviceGroupList API omits the list key entirely when the
		// account has no device groups, and jsonpath.Get returns an error in
		// that case. Treat the error as "no items" rather than failing the
		// sweep (see rocketmq-sub #10925 for the same defect pattern).
		result := []interface{}{}
		if resp, jsonErr := jsonpath.Get("$.Data.GroupInfo", response); jsonErr == nil {
			if list, ok := resp.([]interface{}); ok {
				result = list
			}
		}

		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := fmt.Sprint(item["GroupName"])

			// NODELETING / do-not-delete name guard at the TOP of the per-instance
			// loop. Runs in BOTH sweepAll + normal mode (red-line #2; mirrors the
			// hbase-sub #10923 / rds-sub #10924 / bastionhost-sub #10926 pattern).
			if protectedSweepVpcName(name) {
				log.Printf("[INFO] Skipping Iot Device Group with do-not-delete marker: %s", name)
				continue
			}

			// Standard prefix filter — only restrict by name when SWEEP_ALL is
			// not set. When ALICLOUD_SWEEP_ALL_RESOURCES=true, bypass the prefix
			// check (the NODELETING guard above still protects by name).
			if !sweepAll() {
				skip := true
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Iot Device Group (prefix mismatch): %s", name)
					continue
				}
			}

			deleteAction := "DeleteDeviceGroup"
			deleteRequest := map[string]interface{}{
				"GroupId": item["GroupId"],
			}
			if _, err := client.RpcPost("Iot", "2018-01-20", deleteAction, nil, deleteRequest, false); err != nil {
				log.Printf("[ERROR] Failed to delete Iot Device Group (%s): %s", name, err)
				continue
			}

			log.Printf("[INFO] Delete Iot Device Group Success: %s", name)
		}

		if len(result) < PageSizeSmall {
			break
		}
		currentPage++
	}
	return nil
}

func TestAccAlicloudIotDeviceGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_iot_device_group.default"
	ra := resourceAttrInit(resourceId, AlicloudIotDeviceGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &IotService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeIotDeviceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 1000)
	name := fmt.Sprintf("tf_testacciot%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudIotDeviceGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", "cn-shanghai")
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_desc": name + "group_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_desc": name + "group_desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_desc": name + "group_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_desc": name + "group_desc",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"iot_instance_id", "super_group_id"},
			},
		},
	})
}

var AlicloudIotDeviceGroupMap0 = map[string]string{}

func AlicloudIotDeviceGroupBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudIotDeviceGroup(t *testing.T) {
	p := Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_iot_device_group"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_iot_device_group"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"group_desc":      "CreateDeviceGroupValue",
		"group_name":      "CreateDeviceGroupValue",
		"iot_instance_id": "CreateDeviceGroupValue",
		"super_group_id":  "CreateDeviceGroupValue",
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
		// QueryDeviceGroupInfo
		"Data": map[string]interface{}{
			"GroupDesc": "CreateDeviceGroupValue",
			"GroupName": "CreateDeviceGroupValue",
			"GroupId":   "CreateDeviceGroupValue",
		},
		"Success": "true",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateDeviceGroup
		"Data": map[string]interface{}{
			"GroupId": "CreateDeviceGroupValue",
		},
		"Success": "true",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_iot_device_group", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewIotClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudIotDeviceGroupCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// QueryDeviceGroupInfo Response
		"Data": map[string]interface{}{
			"GroupId": "CreateDeviceGroupValue",
		},
		"Success": "true",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateDeviceGroup" {
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
		err := resourceAlicloudIotDeviceGroupCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_iot_device_group"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewIotClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudIotDeviceGroupUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateDeviceGroup
	attributesDiff := map[string]interface{}{
		"group_desc":      "UpdateDeviceGroupValue",
		"iot_instance_id": "UpdateDeviceGroupValue",
	}
	diff, err := newInstanceDiff("alicloud_iot_device_group", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_iot_device_group"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// QueryDeviceGroupInfo Response
		"Data": map[string]interface{}{
			"GroupDesc": "UpdateDeviceGroupValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateDeviceGroup" {
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
		err := resourceAlicloudIotDeviceGroupUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_iot_device_group"].Schema).Data(dExisted.State(), nil)
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
			if *action == "QueryDeviceGroupInfo" {
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
		err := resourceAlicloudIotDeviceGroupRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewIotClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudIotDeviceGroupDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteDeviceGroup" {
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
		err := resourceAlicloudIotDeviceGroupDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
