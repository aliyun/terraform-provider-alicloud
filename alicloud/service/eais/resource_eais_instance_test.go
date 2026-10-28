// Package eais_test.
package eais_test

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eais"
	sweep "github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

// collection-order: bypass environment_var -- pre-existing EAIS attribute unchanged by this subpackage migration; product deprecated since v1.278.0 with all acceptance tests SkipTest-prefixed, so the required two-member reorder convergence cannot be exercised against the live API.
func init() {
	resource.AddTestSweepers(
		"alicloud_eais_instance",
		&resource.Sweeper{
			Name: "alicloud_eais_instance",
			F:    testSweepEaisInstance,
		})
}

func testSweepEaisInstance(region string) error {

	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeEais"
	request := map[string]interface{}{}

	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
	request["RegionId"] = client.RegionId

	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("eais", "2019-06-24", action, nil, request, true)
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
		v, err := jsonpath.Get("$.Instances.Instance", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Instances.Instance", action, err)
			return nil
		}
		result, _ := v.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["InstanceName"]; !ok {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["InstanceName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Eais Instance: %s", item["InstanceName"].(string))
				continue
			}
			action := "DeleteEai"
			request := map[string]interface{}{
				"ElasticAcceleratedInstanceId": item["ElasticAcceleratedInstanceId"],
				"Force":                        true,
			}
			request["ClientToken"] = helper.BuildClientToken("DeleteEai")
			_, err = client.RpcPost("eais", "2019-06-24", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Eais Instance (%s): %s", item["ElasticAcceleratedInstanceId"].(string), err)
			}
			log.Printf("[INFO] Delete Eais Instance success: %s ", item["ElasticAcceleratedInstanceId"].(string))

		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func SkipTestAccAliCloudEaisInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.EAISSystemSupportRegions)
	resourceId := "alicloud_eais_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEaisInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eais.NewEaisService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEaisInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seaisinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEaisInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type":     "eais.ei-a6.2xlarge",
					"vswitch_id":        "${alicloud_vswitch.default.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"category":          "eais",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":     "eais.ei-a6.2xlarge",
						"vswitch_id":        acctest.CHECKSET,
						"security_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"category"},
			},
		},
	})
}

func SkipTestAccAliCloudEaisInstance_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.EAISSystemSupportRegions)
	resourceId := "alicloud_eais_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEaisInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eais.NewEaisService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEaisInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seaisinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEaisInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type":     "eais.ei-a6.2xlarge",
					"vswitch_id":        "${alicloud_vswitch.default.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"instance_name":     name,
					"force":             "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":     "eais.ei-a6.2xlarge",
						"vswitch_id":        acctest.CHECKSET,
						"security_group_id": acctest.CHECKSET,
						"instance_name":     name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

var AliCloudEaisInstanceMap0 = map[string]string{
	"instance_name": acctest.CHECKSET,
	"status":        acctest.CHECKSET,
}

func AliCloudEaisInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	locals {
  		zone_id = "cn-hangzhou-h"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = local.zone_id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}
`, name)
}

// lintignore: R001
func SkipTestUnitAliCloudEaisInstance(t *testing.T) {
	p := acctest.Provider.ResourcesMap
	acctest.CheckoutSupportedRegions(t, true, connectivity.EAISSystemSupportRegions)
	dInit, _ := schema.InternalMap(p["alicloud_eais_instance"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_eais_instance"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"instance_name":     "CreateEaiValue",
		"instance_type":     "CreateEaiValue",
		"security_group_id": "CreateEaiValue",
		"vswitch_id":        "CreateEaiValue",
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
		// DescribeEais
		"Instances": map[string]interface{}{
			"Instance": []interface{}{
				map[string]interface{}{
					"InstanceName":                 "CreateEaiValue",
					"InstanceType":                 "CreateEaiValue",
					"Status":                       "Available",
					"ElasticAcceleratedInstanceId": "CreateEaiValue",
				},
			},
		},
		"ElasticAcceleratedInstanceId": "CreateEaiValue",
	}
	CreateMockResponse := map[string]interface{}{
		//CreateEaiValue
		"ElasticAcceleratedInstanceId": "CreateEaiValue",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_eais_instance", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEaisClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = eais.ResourceEaisInstanceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateEai" {
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
		err := eais.ResourceEaisInstanceCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_eais_instance"].Schema).Data(dInit.State(), nil)
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
	err = eais.ResourceEaisInstanceUpdate(dExisted, rawClient)
	assert.NotNil(t, err)

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeEais" {
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
		err := eais.ResourceEaisInstanceRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEaisClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = eais.ResourceEaisInstanceDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteEai" {
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
		err := eais.ResourceEaisInstanceDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}

// Test Eais Instance. >>> Resource test cases, automatically generated.
// Case cc_jupyter_pro 10128
func SkipTestAccAliCloudEaisInstance_basic10128(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eais_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEaisInstanceMap10128)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eais.NewEaisServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEaisInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceais%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEaisInstanceBasicDependence10128)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"security_group_id": "${alicloud_security_group.sg.id}",
					"vswitch_id":        "${alicloud_vswitch.vsw.id}",
					"instance_type":     "eais.ei-a6.2xlarge",
					"image":             "registry-vpc.cn-hangzhou.aliyuncs.com/eai_beijing/eais-triton:v3.2.5-server",
					"environment_var": []map[string]interface{}{
						{
							"key":   "testKey",
							"value": "testValue",
						},
					},
					"status":   "InUse",
					"category": "jupyter",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name,
						"resource_group_id": acctest.CHECKSET,
						"security_group_id": acctest.CHECKSET,
						"vswitch_id":        acctest.CHECKSET,
						"instance_type":     "eais.ei-a6.2xlarge",
						"environment_var.#": "1",
						"status":            "InUse",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Stopped",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Stopped",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "InUse",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "InUse",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"environment_var": []map[string]interface{}{
						{
							"key":   "testKey-update",
							"value": "testValue-update",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"environment_var", "image", "category"},
			},
		},
	})
}

var AlicloudEaisInstanceMap10128 = map[string]string{
	"create_time": acctest.CHECKSET,
	"region_id":   "cn-hangzhou",
}

func AlicloudEaisInstanceBasicDependence10128(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone" {
  default = "cn-hangzhou-i"
}

variable "region" {
  default = "cn-hangzhou"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vsw" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = alicloud_vpc.vpc.id
  zone_id      = var.zone
}

resource "alicloud_security_group" "sg" {
  security_group_name = alicloud_vpc.vpc.id
  vpc_id              = alicloud_vpc.vpc.id
}


`, name)
}

// Test Eais Instance. <<< Resource test cases, automatically generated.
