package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_eais_instance",
		&resource.Sweeper{
			Name: "alicloud_eais_instance",
			F:    testSweepEaisInstance,
		})
}

func testSweepEaisInstance(region string) error {

	rawClient, err := sharedClientForRegion(region)
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

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["RegionId"] = client.RegionId

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("eais", "2019-06-24", action, nil, request, true)
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
			request["ClientToken"] = buildClientToken("DeleteEai")
			_, err = client.RpcPost("eais", "2019-06-24", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Eais Instance (%s): %s", item["ElasticAcceleratedInstanceId"].(string), err)
			}
			log.Printf("[INFO] Delete Eais Instance success: %s ", item["ElasticAcceleratedInstanceId"].(string))

		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudEaisInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.EAISSystemSupportRegions)
	resourceId := "alicloud_eais_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudEaisInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EaisService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEaisInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seaisinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEaisInstanceBasicDependence0)
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
					"instance_type":     "eais.ei-a6.2xlarge",
					"vswitch_id":        "${alicloud_vswitch.default.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"category":          "eais",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":     "eais.ei-a6.2xlarge",
						"vswitch_id":        CHECKSET,
						"security_group_id": CHECKSET,
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

func TestAccAliCloudEaisInstance_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.EAISSystemSupportRegions)
	resourceId := "alicloud_eais_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudEaisInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EaisService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEaisInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seaisinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEaisInstanceBasicDependence0)
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
					"instance_type":     "eais.ei-a6.2xlarge",
					"vswitch_id":        "${alicloud_vswitch.default.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"instance_name":     name,
					"force":             "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":     "eais.ei-a6.2xlarge",
						"vswitch_id":        CHECKSET,
						"security_group_id": CHECKSET,
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
	"instance_name": CHECKSET,
	"status":        CHECKSET,
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

func TestUnitAliCloudEaisInstance(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	checkoutSupportedRegions(t, true, connectivity.EAISSystemSupportRegions)
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
	rawClient, err := sharedClientForRegion(region)
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
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_eais_instance", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEaisClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudEaisInstanceCreate(dInit, rawClient)
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
		err := resourceAliCloudEaisInstanceCreate(dInit, rawClient)
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
	err = resourceAliCloudEaisInstanceUpdate(dExisted, rawClient)
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
		err := resourceAliCloudEaisInstanceRead(dExisted, rawClient)
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
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudEaisInstanceDelete(dExisted, rawClient)
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
		err := resourceAliCloudEaisInstanceDelete(dExisted, rawClient)
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
func TestAccAliCloudEaisInstance_basic10128(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eais_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudEaisInstanceMap10128)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EaisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEaisInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceais%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEaisInstanceBasicDependence10128)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"resource_group_id": CHECKSET,
						"security_group_id": CHECKSET,
						"vswitch_id":        CHECKSET,
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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
	"create_time": CHECKSET,
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
