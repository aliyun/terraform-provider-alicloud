package brainindustrial_test

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/brainindustrial"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

// brainIndustrialTestClientForRegion delegates to the shared sweeper client
// (alicloud/sweep.SharedClientForRegion) so the package tests reuse the same
// credential, ALICLOUD_* / ALIBABA_CLOUD_* account-id/type fallback and the
// region-isolated endpoint cache as the root suite, instead of a lossy local
// copy that only handled AK/STS with a single global endpoint map.
func brainIndustrialTestClientForRegion(region string) (*connectivity.AliyunClient, error) {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return nil, err
	}
	return rawClient.(*connectivity.AliyunClient), nil
}

// brainIndustrialSweepSupportedRegion mirrors the root-suite
// testSweepPreCheckWithRegions(region, false, BrainIndustrialRegions): sweep
// only runs in the regions that actually support Brain Industrial.
func brainIndustrialSweepSupportedRegion(region string) bool {
	for _, r := range connectivity.BrainIndustrialRegions {
		if region == string(r) {
			return true
		}
	}
	return false
}

func init() {
	resource.AddTestSweepers(
		"alicloud_brain_industrial_pid_organization",
		&resource.Sweeper{
			Name: "alicloud_brain_industrial_pid_organization",
			F:    testSweepBrainIndustrialPidOrganization,
			Dependencies: []string{
				"alicloud_brain_industrial_pid_project",
			},
		})

}

func testSweepBrainIndustrialPidOrganization(region string) error {
	if !brainIndustrialSweepSupportedRegion(region) {
		log.Printf("[INFO] Skipping Brain Industrial unsupported region: %s", region)
		return nil
	}
	rawClient, err := brainIndustrialTestClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}
	request := map[string]interface{}{}
	var response map[string]interface{}
	action := "ListPidOrganizations"
	response, err = client.RpcPost("brain-industrial", "2020-09-20", action, nil, request, true)
	if err != nil {
		return errs.WrapError(err)
	}
	resp, err := jsonpath.Get("$.OrganizationList", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.OrganizationList", response)
	}
	sweeped := false
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(item["OrganizationName"].(string)), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Brain Industrial Organization: %s", item["OrganizationName"].(string))
			continue
		}
		sweeped = true
		action = "DeletePidOrganization"
		request := map[string]interface{}{
			"OrganizationId": item["OrganizationId"],
		}
		responseDelete, _ := client.RpcPost("brain-industrial", "2020-09-20", action, nil, request, false)
		if fmt.Sprintf(`%v`, responseDelete["Code"]) != "200" {
			log.Printf("[ERROR] Failed to delete Brain Industrial Organization (%s): %s", item["OrganizationName"].(string), responseDelete["Message"])
		}
		if sweeped {
			// Waiting 5 seconds to ensure  Brain Industrial Organization have been deleted.
			time.Sleep(5 * time.Second)
		}
		log.Printf("[INFO] Delete Brain Industrial Organization success: %s ", item["OrganizationName"].(string))
	}
	return nil
}

func TestAccAliCloudBrainIndustrialPidOrganization_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_brain_industrial_pid_organization.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudBrainIndustrialPidOrganizationMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return brainindustrial.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePidOrganization")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudBrainIndustrialPidOrganizationBasicDependence)
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
					"pid_organization_name":      name,
					"parent_pid_organization_id": "${alicloud_brain_industrial_pid_organization.parent.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_organization_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent_pid_organization_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_organization_name":      "tf-testAccUp",
					"parent_pid_organization_id": "${alicloud_brain_industrial_pid_organization.parent.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_organization_name": "tf-testAccUp",
					}),
				),
			},
		},
	})
}

var AliCloudBrainIndustrialPidOrganizationMap = map[string]string{}

func AliCloudBrainIndustrialPidOrganizationBasicDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_brain_industrial_pid_organization" "parent" {
	pid_organization_name = "%[1]s-parent"
}`, name)
}

// lintignore: R001
func TestUnitAliCloudBrainIndustrialPidOrganization(t *testing.T) {
	resourceSchema := brainindustrial.ResourcePidOrganization()
	dInit, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"pid_organization_name":      "CreatePidOrganizationValue",
		"parent_pid_organization_id": "CreatePidOrganizationValue",
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
		// ListPidOrganizations
		"OrganizationList": []interface{}{
			map[string]interface{}{
				"OrganizationName": "CreatePidOrganizationValue",
				"OrganizationId":   "CreatePidOrganizationValue",
			},
		},
		"OrganizationId": "CreatePidOrganizationValue",
		"Message":        "CreatePidOrganizationValue",
		"Code":           "200",
	}
	CreateMockResponse := map[string]interface{}{
		// CreatePidOrganization
		"OrganizationId": "CreatePidOrganizationValue",
		"Message":        "CreatePidOrganizationValue",
		"Code":           "200",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_brain_industrial_pid_organization", errorCode))
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
		// ListPidOrganizations Response
		"OrganizationId": "CreatePidOrganizationValue",
		"Message":        "CreatePidOrganizationValue",
		"Code":           "200",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreatePidOrganization" {
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
	// UpdatePidOrganization
	attributesDiff := map[string]interface{}{
		"pid_organization_name": "UpdatePidOrganizationValue",
	}
	diff, err := acctest.NewInstanceDiff("alicloud_brain_industrial_pid_organization", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// ListPidOrganizations Response
		"OrganizationList": []interface{}{
			map[string]interface{}{
				"OrganizationName": "UpdatePidOrganizationValue",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdatePidOrganization" {
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
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ListPidOrganizations" {
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeletePidOrganization" {
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
