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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_smartag_flow_log",
		&resource.Sweeper{
			Name: "alicloud_smartag_flow_log",
			F:    testSweepSmartagFlowLog,
		})
}

func testSweepSmartagFlowLog(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.SmartagSupportedRegions) {
		log.Printf("[INFO] Skipping Smartag Flow Log unsupported region: %s", region)
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeFlowLogs"
	request := map[string]interface{}{}
	request["RegionId"] = aliyunClient.RegionId

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	conn, err := aliyunClient.NewSmartagClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-03-13"), StringPointer("AK"), nil, request, &runtime)
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

		resp, err := jsonpath.Get("$.FlowLogs.FlowLogSetType", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.FlowLogs.FlowLogSetType", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Name"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Smartag Flow Log: %s", item["Name"].(string))
				continue
			}
			action := "DeleteFlowLog"
			request := map[string]interface{}{
				"FlowLogId": item["FlowLogId"],
				"RegionId":  aliyunClient.RegionId,
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-03-13"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Smartag Flow Log (%s): %s", item["Name"].(string), err)
			}
			log.Printf("[INFO] Delete Smartag Flow Log success: %s ", item["Name"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudSmartagFlowLog_basic_netflow(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_smartag_flow_log.default"
	checkoutSupportedRegions(t, true, connectivity.SmartagSupportedRegions)
	ra := resourceAttrInit(resourceId, AlicloudSmartagFlowLogMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSmartagFlowLog")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssmartagflowlog%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSmartagFlowLogBasicDependence0)
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
					"netflow_server_port": "9995",
					"netflow_server_ip":   "192.168.0.2",
					"netflow_version":     "V9",
					"output_type":         "netflow",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"netflow_server_port": "9995",
						"netflow_server_ip":   "192.168.0.2",
						"netflow_version":     "V9",
						"output_type":         "netflow",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudSmartagFlowLog_basic_sls(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_smartag_flow_log.default"
	checkoutSupportedRegions(t, true, connectivity.SmartagSupportedRegions)
	ra := resourceAttrInit(resourceId, AlicloudSmartagFlowLogMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSmartagFlowLog")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssmartagflowlog%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSmartagFlowLogBasicDependence0)
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
					"logstore_name": "${var.name}",
					"project_name":  "${var.name}",
					"sls_region_id": defaultRegionToTest,
					"output_type":   "sls",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logstore_name": name,
						"project_name":  name,
						"sls_region_id": defaultRegionToTest,
						"output_type":   "sls",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudSmartagFlowLog_basic_all(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_smartag_flow_log.default"
	checkoutSupportedRegions(t, true, connectivity.SmartagSupportedRegions)
	ra := resourceAttrInit(resourceId, AlicloudSmartagFlowLogMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSmartagFlowLog")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssmartagflowlog%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSmartagFlowLogBasicDependence0)
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
					"netflow_server_port": "9995",
					"logstore_name":       "${var.name}",
					"description":         "${var.name}",
					"active_aging":        "300",
					"project_name":        "${var.name}",
					"netflow_server_ip":   "192.168.0.2",
					"netflow_version":     "V9",
					"inactive_aging":      "15",
					"flow_log_name":       "${var.name}",
					"sls_region_id":       defaultRegionToTest,
					"output_type":         "all",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"netflow_server_port": "9995",
						"logstore_name":       name,
						"description":         name,
						"active_aging":        "300",
						"project_name":        name,
						"netflow_server_ip":   "192.168.0.2",
						"netflow_version":     "V9",
						"inactive_aging":      "15",
						"flow_log_name":       name,
						"sls_region_id":       defaultRegionToTest,
						"output_type":         "all",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"netflow_server_port": "9999",
					"logstore_name":       "${var.name}_update",
					"description":         "${var.name}_update",
					"active_aging":        "400",
					"project_name":        "${var.name}_update",
					"netflow_server_ip":   "192.168.0.1",
					"netflow_version":     "V10",
					"inactive_aging":      "20",
					"flow_log_name":       "${var.name}_update",
					"sls_region_id":       defaultRegionToTest,
					"output_type":         "all",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"netflow_server_port": "9999",
						"logstore_name":       name + "_update",
						"description":         name + "_update",
						"active_aging":        "400",
						"project_name":        name + "_update",
						"netflow_server_ip":   "192.168.0.1",
						"netflow_version":     "V10",
						"inactive_aging":      "20",
						"flow_log_name":       name + "_update",
						"sls_region_id":       defaultRegionToTest,
						"output_type":         "all",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Inactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Inactive",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Active",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudSmartagFlowLogMap0 = map[string]string{
	"output_type": CHECKSET,
}

func AlicloudSmartagFlowLogBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

func TestUnitAlicloudSmartagFlowLog(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_smartag_flow_log"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_smartag_flow_log"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"netflow_server_port": 9995,
		"logstore_name":       "CreateSmartagFlowLogValue",
		"description":         "CreateSmartagFlowLogValue",
		"active_aging":        300,
		"project_name":        "CreateSmartagFlowLogValue",
		"netflow_server_ip":   "CreateSmartagFlowLogValue",
		"netflow_version":     "CreateSmartagFlowLogValue",
		"inactive_aging":      15,
		"flow_log_name":       "CreateSmartagFlowLogValue",
		"sls_region_id":       "CreateSmartagFlowLogValue",
		"output_type":         "CreateSmartagFlowLogValue",
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
		"TotalCount": 1,
		"FlowLogs": map[string]interface{}{
			"FlowLogSetType": []interface{}{
				map[string]interface{}{
					"Status":            "Active",
					"NetflowServerPort": 9995,
					"LogstoreName":      "CreateSmartagFlowLogValue",
					"ActiveAging":       300,
					"Description":       "CreateSmartagFlowLogValue",
					"ResourceGroupId":   "CreateSmartagFlowLogValue",
					"ProjectName":       "CreateSmartagFlowLogValue",
					"NetflowVersion":    "CreateSmartagFlowLogValue",
					"NetflowServerIp":   "CreateSmartagFlowLogValue",
					"InactiveAging":     15,
					"FlowLogId":         "CreateSmartagFlowLogValue",
					"Name":              "CreateSmartagFlowLogValue",
					"SlsRegionId":       "CreateSmartagFlowLogValue",
					"OutputType":        "CreateSmartagFlowLogValue",
				},
			},
		},
		"PageSize":   10,
		"PageNumber": 1,
	}
	CreateMockResponse := map[string]interface{}{
		"FlowLogId": "CreateSmartagFlowLogValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_smartag_flow_log", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewSmartagClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudSmartagFlowLogCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateFlowLog" {
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
		err := resourceAlicloudSmartagFlowLogCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_smartag_flow_log"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewSmartagClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudSmartagFlowLogUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"netflow_server_port": 9999,
		"logstore_name":       "UpdateSmartagFlowLogValue",
		"description":         "UpdateSmartagFlowLogValue",
		"active_aging":        400,
		"project_name":        "UpdateSmartagFlowLogValue",
		"netflow_server_ip":   "UpdateSmartagFlowLogValue",
		"netflow_version":     "UpdateSmartagFlowLogValue",
		"inactive_aging":      20,
		"flow_log_name":       "UpdateSmartagFlowLogValue",
		"sls_region_id":       "UpdateSmartagFlowLogValue",
		"output_type":         "UpdateSmartagFlowLogValue",
	}
	diff, err := newInstanceDiff("alicloud_smartag_flow_log", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_smartag_flow_log"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"FlowLogs": map[string]interface{}{
			"FlowLogSetType": []interface{}{
				map[string]interface{}{
					"NetflowServerPort": 9999,
					"LogstoreName":      "UpdateSmartagFlowLogValue",
					"ActiveAging":       400,
					"ResourceGroupId":   "UpdateSmartagFlowLogValue",
					"Description":       "UpdateSmartagFlowLogValue",
					"NetflowServerIp":   "UpdateSmartagFlowLogValue",
					"ProjectName":       "UpdateSmartagFlowLogValue",
					"NetflowVersion":    "UpdateSmartagFlowLogValue",
					"InactiveAging":     20,
					"Name":              "UpdateSmartagFlowLogValue",
					"SlsRegionId":       "UpdateSmartagFlowLogValue",
					"OutputType":        "UpdateSmartagFlowLogValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyFlowLogAttribute" {
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
		err := resourceAlicloudSmartagFlowLogUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_smartag_flow_log"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Inactive
	attributesDiff = map[string]interface{}{
		"status": "Inactive",
	}
	diff, err = newInstanceDiff("alicloud_smartag_flow_log", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_smartag_flow_log"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"FlowLogs": map[string]interface{}{
			"FlowLogSetType": []interface{}{
				map[string]interface{}{
					"Status": "Inactive",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeactiveFlowLog" {
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
		err := resourceAlicloudSmartagFlowLogUpdate(dExisted, rawClient)
		patches.Reset()

		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_smartag_flow_log"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Active
	attributesDiff = map[string]interface{}{
		"status": "Active",
	}
	diff, err = newInstanceDiff("alicloud_smartag_flow_log", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_smartag_flow_log"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"FlowLogs": map[string]interface{}{
			"FlowLogSetType": []interface{}{
				map[string]interface{}{
					"Status": "Active",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ActiveFlowLog" {
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
		err := resourceAlicloudSmartagFlowLogUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_smartag_flow_log"].Schema).Data(dExisted.State(), nil)
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
			if *action == "DescribeFlowLogs" {
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
		err := resourceAlicloudSmartagFlowLogRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewSmartagClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudSmartagFlowLogDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteFlowLog" {
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
		err := resourceAlicloudSmartagFlowLogDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
