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
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

const consumerChannelResourceId = "alicloud_dts_consumer_channel.default"

func TestAccAlicloudDTSConsumerChannel_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tftest%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkConsumerChannelDestroy(consumerChannelResourceId),
		Steps: []resource.TestStep{
			{
				Config: consumerChannelBasicDependence0(name) + fmt.Sprintf(`
resource "alicloud_dts_consumer_channel" "default" {
  consumer_group_password  = "tftestAcc123"
  consumer_group_user_name = "%s"
  consumer_group_name      = "%s"
  dts_instance_id          = "${alicloud_dts_subscription_job.default.dts_instance_id}"
}
`, name, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(consumerChannelResourceId, "consumer_group_user_name", name),
					resource.TestCheckResourceAttr(consumerChannelResourceId, "consumer_group_name", name),
					resource.TestCheckResourceAttrSet(consumerChannelResourceId, "dts_instance_id"),
				),
			},
			{
				Config: consumerChannelBasicDependence0(name) + fmt.Sprintf(`
resource "alicloud_dts_consumer_channel" "default" {
  consumer_group_password  = "tftestAcc123update"
  consumer_group_user_name = "%s"
  consumer_group_name      = "%s"
  dts_instance_id          = "${alicloud_dts_subscription_job.default.dts_instance_id}"
}
`, name, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(consumerChannelResourceId, "consumer_group_user_name", name),
					resource.TestCheckResourceAttr(consumerChannelResourceId, "consumer_group_name", name),
					resource.TestCheckResourceAttrSet(consumerChannelResourceId, "dts_instance_id"),
				),
			},
			{
				ResourceName:            consumerChannelResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"consumer_group_password"},
			},
		},
	})
}

func checkConsumerChannelDestroy(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_dts_consumer_channel" {
				continue
			}
			client := acctest.Provider.Meta().(*connectivity.AliyunClient)
			if _, err := dts.NewService(client).DescribeDtsConsumerChannel(rs.Primary.ID); err != nil {
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

func consumerChannelBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

## RDS MySQL Source
resource "alicloud_db_instance" "source" {
  engine           = "MySQL"
  engine_version   = "8.0"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  vswitch_id       = data.alicloud_vswitches.default.ids.0
  instance_name    = "rds-mysql-source"
}

resource "alicloud_db_database" "source_db" {
  instance_id = alicloud_db_instance.source.id
  name        = "test_database"
}

resource "alicloud_rds_account" "source_account" {
  db_instance_id   = alicloud_db_instance.source.id
  account_name     = "test_mysql"
  account_password = "N1cetest"
}

resource "alicloud_db_account_privilege" "source_privilege" {
  instance_id  = alicloud_db_instance.source.id
  account_name = alicloud_rds_account.source_account.name
  privilege    = "ReadWrite"
  db_names     = alicloud_db_database.source_db.*.name
}

## RDS MySQL Target
resource "alicloud_db_instance" "target" {
  engine           = "MySQL"
  engine_version   = "8.0"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  vswitch_id       = data.alicloud_vswitches.default.ids.0
  instance_name    = "rds-mysql-target"
}

resource "alicloud_rds_account" "target_account" {
  db_instance_id   = alicloud_db_instance.target.id
  account_name     = "test_mysql"
  account_password = "N1cetest"
}

resource "alicloud_dts_subscription_job" "default" {
    dts_job_name                        = var.name
    payment_type                        = "PayAsYouGo"
    source_endpoint_engine_name         = "MySQL"
    source_endpoint_region              = data.alicloud_regions.default.regions.0.id
    source_endpoint_instance_type       = "RDS"
    source_endpoint_instance_id         = "${alicloud_db_instance.source.id}"
    source_endpoint_database_name       = "${alicloud_rds_account.source_account.account_password}"
    source_endpoint_user_name           = "${alicloud_rds_account.source_account.account_name}"
    source_endpoint_password            = "${alicloud_rds_account.source_account.account_password}"
    db_list = "{\"test_database\":{\"name\":\"test_database\",\"all\":true,\"state\":\"normal\"}}"
    subscription_instance_network_type  = "vpc"
    subscription_instance_vpc_id        = data.alicloud_vpcs.default.ids[0]
    subscription_instance_vswitch_id    = data.alicloud_vswitches.default.ids[0]
    status                              = "Normal"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudDTSConsumerChannel(t *testing.T) {
	resourceSchema := dts.ResourceConsumerChannel()
	dInit, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"consumer_group_password":  "CreateConsumerChannelValue",
		"consumer_group_user_name": "CreateConsumerChannelValue",
		"consumer_group_name":      "CreateConsumerChannelValue",
		"dts_instance_id":          "CreateConsumerChannelValue",
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := dtsTestClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	ReadMockResponse := map[string]interface{}{
		// DescribeConsumerChannel
		"ConsumerChannels": []interface{}{
			map[string]interface{}{
				"ConsumerGroupID":       "CreateConsumerChannelValue",
				"DtsInstanceId":         "CreateConsumerChannelValue",
				"ConsumerGroupName":     "CreateConsumerChannelValue",
				"ConsumerGroupUserName": "CreateConsumerChannelValue",
				"ConsumerGroupId":       "CreateConsumerChannelValue",
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		//CreateConsumerChannel
		"ConsumerGroupID": "CreateConsumerChannelValue",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_dts_consumer_channel", errorCode))
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
		"ConsumerGroupID": "CreateConsumerChannelValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "OperationDenied.JobStatus", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateConsumerChannel" {
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
	// ModifyConsumerChannel
	attributesDiff := map[string]interface{}{
		"consumer_group_password": "ModifyConsumerChannelValue",
	}
	// newInstanceDiff merged attributesDiff into attributes in-place; replicate that merge so the
	// dCompare loop below compares against the updated value.
	for k, v := range attributesDiff {
		attributes[k] = v
	}
	diff := sdkterraform.NewInstanceDiff()
	diff.Attributes["consumer_group_password"] = &sdkterraform.ResourceAttrDiff{Old: dInit.Get("consumer_group_password").(string), New: "ModifyConsumerChannelValue"}
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeConsumerChannel Response

	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyConsumerChannel" {
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
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeConsumerChannel" {
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
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "OperationDenied.JobStatus", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteConsumerChannel" {
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
		case "nil":
			assert.Nil(t, err)
		}
	}
}
