package fnf_test

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/fnf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_fnf_schedule", &resource.Sweeper{
		Name: "alicloud_fnf_schedule",
		F:    testSweepFnfSchedule,
	})
}

func testSweepFnfSchedule(region string) error {
	client, err := fnfTestClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	support := false
	for _, v := range connectivity.FnfSupportRegions {
		if v == connectivity.Region(region) {
			support = true
			break
		}
	}
	if !support {
		return nil
	}

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "ListFlows"
	request := map[string]interface{}{
		"Limit": 100,
	}
	var response map[string]interface{}
	response, err = client.RpcGet("fnf", "2019-03-15", action, request, nil)
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_fnf_flows", action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	resp, err := jsonpath.Get("$.Flows", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Flows", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		name := item["Name"].(string)

		action := "ListSchedules"
		request := map[string]interface{}{
			"FlowName": name,
			"Limit":    100,
		}
		var response map[string]interface{}
		response, err = client.RpcGet("fnf", "2019-03-15", action, request, nil)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_fnf_schedules", action, errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(action, response, request)

		resp, err := jsonpath.Get("$.Schedules", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Schedules", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := item["ScheduleName"].(string)
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(name, prefix) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Fnf Schedule: %s ", name)
				continue
			}
			log.Printf("[Info] Delete Fnf Schedule: %s", name)

			action := "DeleteSchedule"
			request := map[string]interface{}{
				"FlowName":     item["FlowName"],
				"ScheduleName": name,
			}
			_, err = client.RpcGet("fnf", "2019-03-15", action, request, nil)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Fnf Schedule (%s): %s", name, err)
			}
		}
	}
	return nil
}

const fnfScheduleResourceId = "alicloud_fnf_schedule.default"

func TestAccAlicloudFnfSchedule_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudFnfSchedule%d", fnfTestRegion(), rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.FnfSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkFnfScheduleDestroy(fnfScheduleResourceId),
		Steps: []resource.TestStep{
			{
				Config: fnfScheduleTestAccDependence(name) + fmt.Sprintf(`
resource "alicloud_fnf_schedule" "default" {
  cron_expression = "30 9 * * * *"
  flow_name       = alicloud_fnf_flow.default.name
  schedule_name   = var.name
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fnfScheduleResourceId, "cron_expression", "30 9 * * * *"),
					resource.TestCheckResourceAttrSet(fnfScheduleResourceId, "flow_name"),
					resource.TestCheckResourceAttr(fnfScheduleResourceId, "schedule_name", name),
				),
			},
			{
				ResourceName:      fnfScheduleResourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fnfScheduleTestAccDependence(name) + fmt.Sprintf(`
resource "alicloud_fnf_schedule" "default" {
  cron_expression = "30 18 * * * *"
  flow_name        = alicloud_fnf_flow.default.name
  schedule_name    = var.name
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fnfScheduleResourceId, "cron_expression", "30 18 * * * *"),
				),
			},
			{
				Config: fnfScheduleTestAccDependence(name) + fmt.Sprintf(`
resource "alicloud_fnf_schedule" "default" {
  cron_expression = "30 18 * * * *"
  description     = "tf-testaccFnFSchedule813242"
  flow_name       = alicloud_fnf_flow.default.name
  schedule_name   = var.name
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fnfScheduleResourceId, "description", "tf-testaccFnFSchedule813242"),
				),
			},
			{
				Config: fnfScheduleTestAccDependence(name) + fmt.Sprintf(`
resource "alicloud_fnf_schedule" "default" {
  cron_expression = "30 18 * * * *"
  description     = "tf-testaccFnFSchedule813242"
  enable          = false
  flow_name       = alicloud_fnf_flow.default.name
  schedule_name   = var.name
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fnfScheduleResourceId, "enable", "false"),
				),
			},
			{
				Config: fnfScheduleTestAccDependence(name) + fmt.Sprintf(`
resource "alicloud_fnf_schedule" "default" {
  cron_expression = "30 18 * * * *"
  description     = "tf-testaccFnFSchedule813242"
  enable          = false
  flow_name       = alicloud_fnf_flow.default.name
  payload         = "{\"tf-testchange\": \"test success\"}"
  schedule_name   = var.name
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fnfScheduleResourceId, "payload", `{"tf-testchange": "test success"}`),
				),
			},
			{
				Config: fnfScheduleTestAccDependence(name) + fmt.Sprintf(`
resource "alicloud_fnf_schedule" "default" {
  cron_expression = "30 9 * * * *"
  description     = "tf-testaccFnFSchedule983041"
  enable          = true
  flow_name       = alicloud_fnf_flow.default.name
  payload         = "{\"tf-test\": \"test success\"}"
  schedule_name   = var.name
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fnfScheduleResourceId, "cron_expression", "30 9 * * * *"),
					resource.TestCheckResourceAttr(fnfScheduleResourceId, "description", "tf-testaccFnFSchedule983041"),
					resource.TestCheckResourceAttr(fnfScheduleResourceId, "enable", "true"),
					resource.TestCheckResourceAttr(fnfScheduleResourceId, "payload", `{"tf-test": "test success"}`),
					resource.TestCheckResourceAttrSet(fnfScheduleResourceId, "last_modified_time"),
					resource.TestCheckResourceAttrSet(fnfScheduleResourceId, "schedule_id"),
				),
			},
		},
	})
}

func checkFnfScheduleDestroy(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		fnfService := fnf.NewService(client)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_fnf_schedule" || rs.Primary.ID != resourceId {
				continue
			}
			if _, err := fnfService.DescribeFnfSchedule(rs.Primary.ID); err != nil {
				if errs.NotFoundError(err) {
					continue
				}
				return err
			}
			return fmt.Errorf("the resource %s %s was not destroyed", resourceId, rs.Primary.ID)
		}
		return nil
	}
}

func fnfScheduleTestAccDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_fnf_flow" "default" {
  definition  = "version: v1beta1\ntype: flow\nsteps:\n  - type: pass\n    name: helloworld"
  description = "tf-testaccFnFFlow983041"
  name        = var.name
  type        = "FDL"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudFnfSchedule(t *testing.T) {
	resourceSchema := fnf.ResourceFnfSchedule()
	dInit, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"cron_expression": "CreateScheduleValue",
		"description":     "CreateScheduleValue",
		"flow_name":       "CreateScheduleValue",
		"payload":         "CreateScheduleValue",
		"enable":          false,
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
	region := fnfTestRegion()
	rawClient, err := fnfTestClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
	}
	ReadMockResponse := map[string]interface{}{
		// DescribeSchedule
		"FlowName":         "CreateScheduleValue",
		"ScheduleName":     "CreateScheduleValue",
		"CronExpression":   "CreateScheduleValue",
		"Description":      "CreateScheduleValue",
		"Enable":           false,
		"LastModifiedTime": "CreateScheduleValue",
		"Payload":          "CreateScheduleValue",
		"ScheduleId":       "CreateScheduleValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateSchedule
		"ScheduleName": "CreateScheduleValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String(errorCode),
			Data:       tea.String(errorCode),
			Message:    tea.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_fnf_schedule", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			for k, v := range operationMockResponse {
				ReadMockResponse[k] = v
			}
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, _ string, _ string, _ string, _ map[string]interface{}, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceSchema.Create(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeSchedule Response
		"ScheduleName": "CreateScheduleValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateSchedule" {
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
		err := resourceSchema.Create(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, _ string, _ string, _ string, _ map[string]interface{}, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceSchema.Update(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateSchedule
	attributesDiff := map[string]interface{}{
		"cron_expression": "UpdateScheduleValue",
		"description":     "UpdateScheduleValue",
		"enable":          true,
		"payload":         "UpdateScheduleValue",
	}
	diff := fnfStringAttrDiff(resourceSchema, attributes, attributesDiff)
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeSchedule Response
		"Description":    "UpdateScheduleValue",
		"Enable":         true,
		"Payload":        "UpdateScheduleValue",
		"CronExpression": "UpdateScheduleValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "ConcurrentUpdateError", "InternalServerError", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateSchedule" {
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
		err := resourceSchema.Update(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(resourceSchema.Schema).Data(dExisted.State(), nil)
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
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeSchedule" {
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
		err := resourceSchema.Read(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcGet", func(_ *connectivity.AliyunClient, _ string, _ string, _ string, _ map[string]interface{}, _ map[string]interface{}) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceSchema.Delete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "ConcurrentUpdateError", "InternalServerError", "nil", "FlowNotExists", "ExecutionNotExists", "ScheduleNotExists"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteSchedule" {
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
		err := resourceSchema.Delete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "FlowNotExists", "ExecutionNotExists", "ScheduleNotExists":
			assert.Nil(t, err)
		}
	}
}
