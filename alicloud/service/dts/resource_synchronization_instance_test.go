package dts_test

import (
	"fmt"
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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

const synchronizationInstanceResourceId = "alicloud_dts_synchronization_instance.default"

func TestAccAliCloudDTSSynchronizationInstance_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssynchronizationinstance%d", dtsTestRegion(), rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-beijing")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkSynchronizationInstanceDestroy(synchronizationInstanceResourceId),
		Steps: []resource.TestStep{
			{
				Config: synchronizationInstanceBasicDependence0(name) + `
resource "alicloud_dts_synchronization_instance" "default" {
  auto_pay                         = "false"
  auto_start                       = "false"
  compute_unit                     = "2"
  database_count                   = "1"
  destination_endpoint_engine_name = "MySQL"
  destination_endpoint_region      = "${data.alicloud_regions.default.regions.0.id}"
  instance_class                   = "small"
  payment_duration                 = "1"
  payment_duration_unit            = "Month"
  payment_type                     = "PayAsYouGo"
  quantity                          = "1"
  source_endpoint_engine_name      = "MySQL"
  source_endpoint_region            = "${data.alicloud_regions.default.regions.0.id}"
  sync_architecture                 = "oneway"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(synchronizationInstanceResourceId, "payment_type", "PayAsYouGo"),
					resource.TestCheckResourceAttr(synchronizationInstanceResourceId, "source_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttrSet(synchronizationInstanceResourceId, "source_endpoint_region"),
					resource.TestCheckResourceAttr(synchronizationInstanceResourceId, "destination_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttrSet(synchronizationInstanceResourceId, "destination_endpoint_region"),
					resource.TestCheckNoResourceAttr(synchronizationInstanceResourceId, "period"),
					resource.TestCheckNoResourceAttr(synchronizationInstanceResourceId, "used_time"),
					resource.TestCheckNoResourceAttr(synchronizationInstanceResourceId, "order_type"),
					resource.TestCheckNoResourceAttr(synchronizationInstanceResourceId, "synchronization_direction"),
					resource.TestCheckNoResourceAttr(synchronizationInstanceResourceId, "status"),
				),
			},
			{
				ResourceName:            synchronizationInstanceResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database_count", "status", "quantity", "sync_architecture", "auto_pay", "auto_start", "compute_unit", "payment_duration", "payment_duration_unit"},
			},
		},
	})
}

func checkSynchronizationInstanceDestroy(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_dts_synchronization_instance" {
				continue
			}
			client := acctest.Provider.Meta().(*connectivity.AliyunClient)
			if _, err := dts.NewService(client).DescribeDtsSynchronizationInstance(rs.Primary.ID); err != nil {
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

func synchronizationInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

// The region has to come from the provider rather than from ALICLOUD_REGION: that
// variable is not always set, and an empty region_id makes CreateDtsInstance fail
// with MissingDestinationRegion before the instance is created.
data "alicloud_regions" "default" {
  current = true
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudDTSSynchronizationInstance(t *testing.T) {
	resourceSchema := dts.ResourceSynchronizationInstance()
	dInit, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"auto_pay":                         "false",
		"auto_start":                       "false",
		"payment_type":                     "PayAsYouGo",
		"source_endpoint_engine_name":      "CreateDtsInstanceValue",
		"source_endpoint_region":           "CreateDtsInstanceValue",
		"destination_endpoint_engine_name": "CreateDtsInstanceValue",
		"destination_endpoint_region":      "CreateDtsInstanceValue",
		"instance_class":                   "CreateDtsInstanceValue",
		"compute_unit":                     1,
		"database_count":                   1,
		"payment_duration_unit":            "CreateDtsInstanceValue",
		"quantity":                         1,
		"sync_architecture":                "CreateDtsInstanceValue",
		"payment_duration":                 1,
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
	}
	region := "cn-beijing"
	rawClient, err := dtsTestClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	ReadMockResponse := map[string]interface{}{
		// DescribeDtsJobDetail
		"SourceEndpoint": map[string]interface{}{
			"EngineName": "CreateDtsInstanceValue",
			"Region":     "CreateDtsInstanceValue",
		},
		"DestinationEndpoint": map[string]interface{}{
			"EngineName": "CreateDtsInstanceValue",
			"Region":     "CreateDtsInstanceValue",
		},
		"PayType":     "PostPaid",
		"DtsJobClass": "CreateDtsInstanceValue",
		"InstanceId":  "CreateDtsInstanceValue",
	}
	CreateMockResponse := map[string]interface{}{
		//CreateDtsInstance
		"InstanceId": "CreateDtsInstanceValue",
	}
	ReadMockResponseDiff := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String(errorCode),
			Data:       tea.String(errorCode),
			Message:    tea.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_dts_synchronization_instance", errorCode))
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
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateDtsInstance" {
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
	err = resourceSchema.Update(dExisted, rawClient)
	assert.Nil(t, err)

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDtsJobDetail" {
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, _ string, _ string, _ string, _ map[string]interface{}, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
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
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "InvalidJobId"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteSynchronizationJob" {
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
		case "InvalidJobId":
			assert.Nil(t, err)
		}
	}
}
