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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_rds_parameter_group", &resource.Sweeper{
		Name: "alicloud_rds_parameter_group",
		F:    testSweepRdsParameterGroups,
	})
}

// testSweepRdsParameterGroups sweeps leftover test RDS parameter groups.
// ParameterGroupType: 0 = system default, 1 = user custom (deletable),
// 2 = system auto-backup. Only custom (type 1) groups created by tests are
// deleted; system groups are always preserved.
func testSweepRdsParameterGroups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{"tf-testAcc", "tf_testAcc"}

	action := "DescribeParameterGroups"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	var response map[string]interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		log.Printf("[ERROR] %s got an error: %v", action, err)
		return nil
	}

	resp, err := jsonpath.Get("$.ParameterGroups.ParameterGroup", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ParameterGroups.ParameterGroup", response)
	}
	objects, _ := resp.([]interface{})

	for _, v := range objects {
		item := v.(map[string]interface{})
		name := fmt.Sprint(item["ParameterGroupName"])
		id := fmt.Sprint(item["ParameterGroupId"])

		// Only custom parameter groups (type 1) are deletable; system (0) and
		// auto-backup (2) templates are preserved regardless of SWEEP_ALL.
		if pgType := fmt.Sprint(item["ParameterGroupType"]); pgType != "1" {
			log.Printf("[INFO] Skipping RDS Parameter Group (system type %s): %s (%s)", pgType, name, id)
			continue
		}
		if protectedSweepVpcName(name) {
			log.Printf("[INFO] Skipping RDS Parameter Group (do-not-delete marker): %s (%s)", name, id)
			continue
		}
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
		} else {
			skip = false
		}
		if skip {
			log.Printf("[INFO] Skipping RDS Parameter Group: %s (%s)", name, id)
			continue
		}

		log.Printf("[INFO] Deleting RDS Parameter Group: %s (%s)", name, id)
		deleteAction := "DeleteParameterGroup"
		deleteRequest := map[string]interface{}{
			"ParameterGroupId": id,
			"SourceIp":         client.SourceIp,
			"RegionId":         client.RegionId,
		}
		if e := retry.Retry(2*time.Minute, func() *retry.RetryError {
			resp, e := client.RpcPost("Rds", "2014-08-15", deleteAction, nil, deleteRequest, false)
			if e != nil {
				if NeedRetry(e) || IsExpectedErrors(e, []string{"InternalError"}) {
					return retry.RetryableError(e)
				}
				return retry.NonRetryableError(e)
			}
			addDebug(deleteAction, resp, deleteRequest)
			return nil
		}); e != nil {
			log.Printf("[ERROR] Failed to delete RDS Parameter Group (%s (%s)): %s", name, id, e)
		}
	}

	return nil
}

func TestAccAliCloudRdsParameterGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_parameter_group.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsParameterGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsParameterGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testAccAlicloudRdsParameterGroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsParameterGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":         "mysql",
					"engine_version": `5.7`,
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "back_log",
							"param_value": `3000`,
						},
						{
							"param_name":  "wait_timeout",
							"param_value": `86400`,
						},
					},
					"parameter_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "mysql",
						"engine_version":       "5.7",
						"param_detail.#":       "2",
						"parameter_group_name": name,
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
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "back_log",
							"param_value": `4000`,
						},
						{
							"param_name":  "wait_timeout",
							"param_value": `86460`,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_detail.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameter_group_desc": "update_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameter_group_desc": "update_test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameter_group_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameter_group_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "back_log",
							"param_value": `3000`,
						},
						{
							"param_name":  "wait_timeout",
							"param_value": `86400`,
						},
					},
					"parameter_group_desc": "test",
					"parameter_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_detail.#":       "2",
						"parameter_group_desc": "test",
						"parameter_group_name": name,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudRdsParameterGroup_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_parameter_group.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsParameterGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsParameterGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testAccAlicloudRdsParameterGroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsParameterGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":         "mysql",
					"engine_version": `5.7`,
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "back_log",
							"param_value": `3000`,
						},
						{
							"param_name":  "wait_timeout",
							"param_value": `86400`,
						},
					},
					"parameter_group_name": "${var.name}",
					"parameter_group_desc": "update_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "mysql",
						"engine_version":       "5.7",
						"param_detail.#":       "2",
						"parameter_group_name": name,
						"parameter_group_desc": "update_test",
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

var AlicloudRdsParameterGroupMap0 = map[string]string{}

func AlicloudRdsParameterGroupBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}
`, name)
}

// lintignore: R001
func TestUnitAlicloudRdsParameterGroup(t *testing.T) {
	p := Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_rds_parameter_group"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_rds_parameter_group"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"engine":         "CreateParameterGroupValue",
		"engine_version": "CreateParameterGroupValue",
		"param_detail": []map[string]interface{}{
			{
				"param_name":  "CreateParameterGroupValue",
				"param_value": "CreateParameterGroupValue",
			},
		},
		"parameter_group_name": "CreateParameterGroupValue",
		"parameter_group_desc": "CreateParameterGroupValue",
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
		// DescribeParameterGroup
		"ParamGroup": map[string]interface{}{
			"ParameterGroup": []interface{}{
				map[string]interface{}{
					"ParameterGroupId": "CreateParameterGroupValue",
					"Engine":           "CreateParameterGroupValue",
					"EngineVersion":    "CreateParameterGroupValue",
					"ParamDetail": map[string]interface{}{
						"ParameterDetail": []interface{}{
							map[string]interface{}{
								"ParamName":  "CreateParameterGroupValue",
								"ParamValue": "CreateParameterGroupValue",
							},
						},
					},
					"ParameterGroupName": "CreateParameterGroupValue",
					"ParameterGroupDesc": "CreateParameterGroupValue",
				},
			},
		},
		"ParameterGroupId": "CreateParameterGroupValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateParameterGroup
		"ParameterGroupId": "CreateParameterGroupValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_rds_parameter_group", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewRdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudRdsParameterGroupCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeParameterGroup Response
		"ParameterGroupId": "CreateParameterGroupValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateParameterGroup" {
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
		err := resourceAlicloudRdsParameterGroupCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_rds_parameter_group"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewRdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudRdsParameterGroupUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyParameterGroup
	attributesDiff := map[string]interface{}{
		"param_detail": []map[string]interface{}{
			{
				"param_name":  "ModifyParameterGroupValue",
				"param_value": "ModifyParameterGroupValue",
			},
		},
		"parameter_group_desc": "ModifyParameterGroupValue",
		"parameter_group_name": "ModifyParameterGroupValue",
	}
	diff, err := newInstanceDiff("alicloud_rds_parameter_group", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_rds_parameter_group"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeParameterGroup Response
		"ParamGroup": map[string]interface{}{
			"ParameterGroup": []interface{}{
				map[string]interface{}{
					"ParamDetail": map[string]interface{}{
						"ParameterDetail": []interface{}{
							map[string]interface{}{
								"ParamName":  "ModifyParameterGroupValue",
								"ParamValue": "ModifyParameterGroupValue",
							},
						},
					},
					"ParameterGroupName": "ModifyParameterGroupValue",
					"ParameterGroupDesc": "ModifyParameterGroupValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyParameterGroup" {
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
		err := resourceAlicloudRdsParameterGroupUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_rds_parameter_group"].Schema).Data(dExisted.State(), nil)
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
			if *action == "DescribeParameterGroup" {
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
		err := resourceAlicloudRdsParameterGroupRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewRdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudRdsParameterGroupDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteParameterGroup" {
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
		err := resourceAlicloudRdsParameterGroupDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}

func TestAccAliCloudRdsParameterGroupPostgreSQL(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_parameter_group.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsParameterGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsParameterGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testAccAlicloudRdsParameterGroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsParameterGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":         "PostgreSQL",
					"engine_version": "11.0",
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "enable_sort",
							"param_value": "off",
						},
						{
							"param_name":  "geqo_seed",
							"param_value": "0.1",
						},
					},
					"parameter_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "PostgreSQL",
						"engine_version":       "11.0",
						"param_detail.#":       "2",
						"parameter_group_name": name,
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
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "enable_sort",
							"param_value": "on",
						},
						{
							"param_name":  "geqo_seed",
							"param_value": "0.2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_detail.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameter_group_desc": "update_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameter_group_desc": "update_test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameter_group_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameter_group_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "enable_sort",
							"param_value": "off",
						},
						{
							"param_name":  "geqo_seed",
							"param_value": "0.1",
						},
					},
					"parameter_group_desc": "test",
					"parameter_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_detail.#":       "2",
						"parameter_group_desc": "test",
						"parameter_group_name": name,
					}),
				),
			},
		},
	})
}
