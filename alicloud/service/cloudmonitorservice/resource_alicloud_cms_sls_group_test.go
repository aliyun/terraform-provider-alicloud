package cloudmonitorservice_test

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudmonitorservice"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_cms_sls_group",
		&resource.Sweeper{
			Name: "alicloud_cms_sls_group",
			F:    testSweepCmsSlsGroup,
		})
}

func testSweepCmsSlsGroup(region string) error {

	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeHybridMonitorSLSGroup"
	request := map[string]interface{}{}

	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.List", response)

		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.List", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			if !sweep.SweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(item["SLSGroupName"].(string)), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Cms Sls Group: %s", item["SLSGroupName"].(string))
					continue
				}
			}
			action := "DeleteHybridMonitorSLSGroup"
			request := map[string]interface{}{
				"SLSGroupName": item["SLSGroupName"],
			}
			_, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Cms Sls Group (%s): %s", item["SLSGroupName"].(string), err)
			}
			log.Printf("[INFO] Delete Cms Sls Group success: %s ", item["SLSGroupName"].(string))
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudCloudMonitorServiceSlsGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_sls_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudMonitorServiceSlsGroupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsSlsGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_cms_slsgroup%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudMonitorServiceSlsGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"sls_group_name": "${var.name}",
					"sls_group_config": []map[string]interface{}{
						{
							"sls_user_id":  "${data.alicloud_account.this.id}",
							"sls_logstore": "Logstore-ECS",
							"sls_project":  "aliyun-project",
							"sls_region":   "cn-hangzhou",
						},
					},
					"sls_group_description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sls_group_name":        name,
						"sls_group_config.#":    "1",
						"sls_group_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sls_group_config": []map[string]interface{}{
						{
							"sls_user_id":  "${data.alicloud_account.this.id}-update",
							"sls_logstore": "Logstore-ECS-Update",
							"sls_project":  "aliyun-project-update",
							"sls_region":   "cn-beijing",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sls_group_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"sls_group_description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sls_group_description": name + "_update",
					}),
				),
			},
			{
				ResourceName: resourceId,
				ImportState:  true,
			},
		},
	})
}

var AlicloudCloudMonitorServiceSlsGroupMap0 = map[string]string{
	"sls_group_name":     acctest.CHECKSET,
	"sls_group_config.#": acctest.CHECKSET,
}

func AlicloudCloudMonitorServiceSlsGroupBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_account" "this" {}
`, name)
}

// lintignore: R001
func TestUnitAlicloudCloudMonitorServiceSlsGroup(t *testing.T) {
	p := unitTestProviderResources()
	dInit, _ := schema.InternalMap(p["alicloud_cms_sls_group"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cms_sls_group"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"sls_group_name": "CreateCloudMonitorServiceSlsGroupValue",
		"sls_group_config": []map[string]interface{}{
			{
				"sls_user_id":  "CreateCloudMonitorServiceSlsGroupValue",
				"sls_logstore": "CreateCloudMonitorServiceSlsGroupValue",
				"sls_project":  "CreateCloudMonitorServiceSlsGroupValue",
				"sls_region":   "CreateCloudMonitorServiceSlsGroupValue",
			},
		},
		"sls_group_description": "CreateCloudMonitorServiceSlsGroupValue",
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
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}

	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"List": []interface{}{
			map[string]interface{}{
				"SLSGroupName":        "CreateCloudMonitorServiceSlsGroupValue",
				"SLSGroupDescription": "CreateCloudMonitorServiceSlsGroupValue",
				"SLSGroupConfig": []interface{}{
					map[string]interface{}{
						"SLSProject":  "CreateCloudMonitorServiceSlsGroupValue",
						"SLSRegion":   "CreateCloudMonitorServiceSlsGroupValue",
						"SLSLogstore": "CreateCloudMonitorServiceSlsGroupValue",
						"SLSUserId":   "CreateCloudMonitorServiceSlsGroupValue",
					},
				},
			},
		},
		"Code":    200,
		"Success": true,
	}
	CreateMockResponse := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       helper.String(errorCode),
			Data:       helper.String(errorCode),
			Message:    helper.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_cms_sls_group", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = p["alicloud_cms_sls_group"].Create(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateHybridMonitorSLSGroup" {
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
		err := p["alicloud_cms_sls_group"].Create(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cms_sls_group"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = p["alicloud_cms_sls_group"].Update(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"sls_group_config": []map[string]interface{}{
			{
				"sls_user_id":  "UpdateCloudMonitorServiceSlsGroupValue",
				"sls_logstore": "UpdateCloudMonitorServiceSlsGroupValue",
				"sls_project":  "UpdateCloudMonitorServiceSlsGroupValue",
				"sls_region":   "UpdateCloudMonitorServiceSlsGroupValue",
			},
		},
		"sls_group_description": "UpdateCloudMonitorServiceSlsGroupValue",
	}
	diff, err := acctest.NewInstanceDiff("alicloud_cms_sls_group", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cms_sls_group"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"List": []interface{}{
			map[string]interface{}{
				"SLSGroupDescription": "UpdateCloudMonitorServiceSlsGroupValue",
				"SLSGroupConfig": []interface{}{
					map[string]interface{}{
						"SLSProject":  "UpdateCloudMonitorServiceSlsGroupValue",
						"SLSRegion":   "UpdateCloudMonitorServiceSlsGroupValue",
						"SLSLogstore": "UpdateCloudMonitorServiceSlsGroupValue",
						"SLSUserId":   "UpdateCloudMonitorServiceSlsGroupValue",
					},
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyHybridMonitorSLSGroup" {
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
		err := p["alicloud_cms_sls_group"].Update(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cms_sls_group"].Schema).Data(dExisted.State(), nil)
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
	diff, err = acctest.NewInstanceDiff("alicloud_cms_sls_group", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cms_sls_group"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeHybridMonitorSLSGroup" {
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
		err := p["alicloud_cms_sls_group"].Read(dExisted, rawClient)
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
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = p["alicloud_cms_sls_group"].Delete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = acctest.NewInstanceDiff("alicloud_cms_sls_group", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cms_sls_group"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteHybridMonitorSLSGroup" {
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
		err := p["alicloud_cms_sls_group"].Delete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
