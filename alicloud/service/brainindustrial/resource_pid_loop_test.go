package brainindustrial_test

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/brainindustrial"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudBrainIndustrialPidLoop_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_brain_industrial_pid_loop.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudBrainIndustrialPidLoopMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return brainindustrial.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePidLoop")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudBrainIndustrialPidLoopBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.BrainIndustrialRegions)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_loop_dcs_type":      "standard",
					"pid_loop_desc":          "Test For Terraform",
					"pid_loop_is_crucial":    "false",
					"pid_loop_name":          name,
					"pid_loop_type":          "0",
					"pid_project_id":         "${alicloud_brain_industrial_pid_project.default.id}",
					"pid_loop_configuration": `{\"baseParam\":{\"forwardController\":false,\"integral\":false,\"kd\":{\"tagValue\":\"20\"},\"kp\":{},\"op\":\"PIDBenchmark.FOPDT_OP\",\"opParam\":{\"increment\":{\"max\":10},\"operate\":{\"max\":115,\"min\":-15},\"range\":{\"max\":115,\"min\":-15},\"trend\":0},\"openLoopTime\":150,\"pv\":\"PIDBenchmark.FOPDT_PV\",\"pvRange\":{\"max\":100,\"min\":0},\"sampleTime\":5,\"sp\":\"PIDBenchmark.FOPDT_SP\",\"spOperate\":{\"max\":100,\"min\":0},\"splitRangeControl\":false,\"suitCtrlTime\":100,\"td\":{},\"ti\":{}},\"identParam\":{\"delay\":10,\"modelType\":3},\"resetParam\":{\"ctrlMode\":0,\"ctrlStuc\":1}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_loop_dcs_type":      "standard",
						"pid_loop_desc":          "Test For Terraform",
						"pid_loop_is_crucial":    "false",
						"pid_loop_name":          name,
						"pid_loop_type":          "0",
						"pid_project_id":         acctest.CHECKSET,
						"pid_loop_configuration": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pid_loop_configuration", "pid_loop_desc", "pid_project_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_loop_desc": "Test For Terraform Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_project_id": acctest.CHECKSET,
						"pid_loop_desc":  "Test For Terraform Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_loop_configuration": `{\"baseParam\":{\"forwardController\":false,\"integral\":false,\"kd\":{\"tagValue\":\"20\"},\"kp\":{},\"op\":\"PIDBenchmark.FOPDT_OP\",\"opParam\":{\"increment\":{\"max\":10},\"operate\":{\"max\":115,\"min\":-15},\"range\":{\"max\":115,\"min\":-15},\"trend\":0},\"openLoopTime\":150,\"pv\":\"PIDBenchmark.FOPDT_PV\",\"pvRange\":{\"max\":100,\"min\":0},\"sampleTime\":6,\"sp\":\"PIDBenchmark.FOPDT_SP\",\"spOperate\":{\"max\":100,\"min\":0},\"splitRangeControl\":false,\"suitCtrlTime\":100,\"td\":{},\"ti\":{}},\"identParam\":{\"delay\":10,\"modelType\":3},\"resetParam\":{\"ctrlMode\":0,\"ctrlStuc\":1}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_loop_configuration": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_loop_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_loop_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_loop_is_crucial": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_project_id":      acctest.CHECKSET,
						"pid_loop_is_crucial": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_loop_type": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_loop_type": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_loop_desc":          "Test For Terraform",
					"pid_loop_name":          name,
					"pid_loop_type":          "0",
					"pid_loop_is_crucial":    "false",
					"pid_loop_configuration": `{\"baseParam\":{\"forwardController\":false,\"integral\":false,\"kd\":{\"tagValue\":\"20\"},\"kp\":{},\"op\":\"PIDBenchmark.FOPDT_OP\",\"opParam\":{\"increment\":{\"max\":10},\"operate\":{\"max\":115,\"min\":-15},\"range\":{\"max\":115,\"min\":-15},\"trend\":0},\"openLoopTime\":150,\"pv\":\"PIDBenchmark.FOPDT_PV\",\"pvRange\":{\"max\":100,\"min\":0},\"sampleTime\":5,\"sp\":\"PIDBenchmark.FOPDT_SP\",\"spOperate\":{\"max\":100,\"min\":0},\"splitRangeControl\":false,\"suitCtrlTime\":100,\"td\":{},\"ti\":{}},\"identParam\":{\"delay\":10,\"modelType\":3},\"resetParam\":{\"ctrlMode\":0,\"ctrlStuc\":1}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_loop_desc":          "Test For Terraform",
						"pid_loop_name":          name,
						"pid_loop_type":          "0",
						"pid_loop_is_crucial":    "false",
						"pid_loop_configuration": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

var AliCloudBrainIndustrialPidLoopMap = map[string]string{}

func AliCloudBrainIndustrialPidLoopBasicDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_brain_industrial_pid_organization" "default" {
		pid_organization_name = "%[1]s"
	}
	resource "alicloud_brain_industrial_pid_project" "default" {
		pid_organization_id = alicloud_brain_industrial_pid_organization.default.id
		pid_project_name = "%[1]s"
	}`, name)
}

// lintignore: R001
func TestUnitAliCloudBrainIndustrialPidLoop(t *testing.T) {
	resourceSchema := brainindustrial.ResourcePidLoop()
	dInit, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"pid_loop_dcs_type":      "CreatePidLoopValue",
		"pid_loop_desc":          "CreatePidLoopValue",
		"pid_loop_is_crucial":    false,
		"pid_loop_name":          "CreatePidLoopValue",
		"pid_loop_type":          "CreatePidLoopValue",
		"pid_project_id":         "CreatePidLoopValue",
		"pid_loop_configuration": "CreatePidLoopValue",
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
	rawClient, err := brainIndustrialTestClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	ReadMockResponse := map[string]interface{}{
		// GetLoop
		"Data": map[string]interface{}{
			"PidLoopConfiguration": "CreatePidLoopValue",
			"PidLoopDcsType":       "CreatePidLoopValue",
			"PidLoopDesc":          "CreatePidLoopValue",
			"PidLoopIsCrucial":     false,
			"PidLoopName":          "CreatePidLoopValue",
			"PidLoopType":          "CreatePidLoopValue",
			"PidProjectId":         "CreatePidLoopValue",
			"Status":               "CreatePidLoopValue",
		},
		"PidLoopId": "CreatePidLoopValue",
		"Code":      "200",
	}
	CreateMockResponse := map[string]interface{}{
		// CreatePidLoop
		"PidLoopId": "CreatePidLoopValue",
		"Code":      "200",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_brain_industrial_pid_loop", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAistudioClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceSchema.Create(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// GetLoop Response
		"PidLoopId": "CreatePidLoopValue",
		"Code":      "200",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreatePidLoop" {
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
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAistudioClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceSchema.Update(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdatePidLoop
	attributesDiff := map[string]interface{}{
		"pid_loop_configuration": "UpdatePidLoopValue",
		"pid_loop_dcs_type":      "UpdatePidLoopValue",
		"pid_loop_is_crucial":    true,
		"pid_loop_name":          "UpdatePidLoopValue",
		"pid_loop_type":          "UpdatePidLoopValue",
		"pid_project_id":         "UpdatePidLoopValue",
		"pid_loop_desc":          "UpdatePidLoopValue",
	}
	diff, err := acctest.NewInstanceDiff("alicloud_brain_industrial_pid_loop", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetLoop Response
		"Data": map[string]interface{}{
			"PidLoopConfiguration": "UpdatePidLoopValue",
			"PidLoopDcsType":       "UpdatePidLoopValue",
			"PidLoopDesc":          "UpdatePidLoopValue",
			"PidLoopIsCrucial":     true,
			"PidLoopName":          "UpdatePidLoopValue",
			"PidLoopType":          "UpdatePidLoopValue",
			"PidProjectId":         "UpdatePidLoopValue",
		},
		"PidLoopId": "UpdatePidLoopValue",
		"Code":      "200",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdatePidLoop" {
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
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetLoop" {
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAistudioClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceSchema.Delete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeletePidLoop" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{
							"Code": "200",
						}
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
		case "nil":
			assert.Nil(t, err)
		}
	}
}
