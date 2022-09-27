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
		"alicloud_ecs_image_pipeline",
		&resource.Sweeper{
			Name: "alicloud_ecs_image_pipeline",
			F:    testSweepEcsImagePipeline,
		})
}

func testSweepEcsImagePipeline(region string) error {

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeImagePipelines"
	request := map[string]interface{}{}
	request["RegionId"] = aliyunClient.RegionId
	request["MaxResults"] = PageSizeLarge

	var response map[string]interface{}
	conn, err := aliyunClient.NewEcsClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
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

		resp, err := jsonpath.Get("$.ImagePipeline.ImagePipelineSet", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.ImagePipeline.ImagePipelineSet", action, err)
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
				log.Printf("[INFO] Skipping Ecs Image Pipeline: %s", item["Name"].(string))
				continue
			}
			action := "DeleteImagePipeline"
			request := map[string]interface{}{
				"ImagePipelineId": item["ImagePipelineId"],
				"RegionId":        aliyunClient.RegionId,
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ecs Image Pipeline (%s): %s", item["Name"].(string), err)
			}
			log.Printf("[INFO] Delete Ecs Image Pipeline success: %s ", item["Name"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudECSImagePipeline_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_image_pipeline.default"
	ra := resourceAttrInit(resourceId, AlicloudECSImagePipelineMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsImagePipeline")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsimagepipeline%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECSImagePipelineBasicDependence0)
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
					"base_image_type": "IMAGE",
					"base_image":      "${data.alicloud_images.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"base_image_type": "IMAGE",
						"base_image":      CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"add_account", "build_content"},
			},
		},
	})
}
func TestAccAlicloudECSImagePipeline_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_image_pipeline.default"
	ra := resourceAttrInit(resourceId, AlicloudECSImagePipelineMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsImagePipeline")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsimagepipeline%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECSImagePipelineBasicDependence0)
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
					"base_image_type":            "IMAGE",
					"name":                       "${var.name}",
					"description":                "${var.name}",
					"system_disk_size":           "40",
					"vswitch_id":                 "${data.alicloud_vswitches.default.ids.0}",
					"image_name":                 "${var.name}",
					"delete_instance_on_failure": "true",
					"internet_max_bandwidth_out": "20",
					"to_region_id":               []string{"cn-qingdao", "cn-zhangjiakou"},
					"add_account":                []string{"${data.alicloud_account.default.id}"},
					"instance_type":              "${data.alicloud_instance_types.default.ids.0}",
					"base_image":                 "${data.alicloud_images.default.ids.0}",
					"build_content":              "RUN yum update -y",
					"resource_group_id":          "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"tags":                       map[string]interface{}{"Created": "TF", "Foo": "Bar"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"base_image_type":            "IMAGE",
						"description":                name,
						"system_disk_size":           "40",
						"vswitch_id":                 CHECKSET,
						"image_name":                 name,
						"delete_instance_on_failure": "true",
						"internet_max_bandwidth_out": "20",
						"to_region_id.#":             "2",
						"add_account.#":              "1",
						"instance_type":              CHECKSET,
						"base_image":                 CHECKSET,
						"resource_group_id":          CHECKSET,
						"tags.%":                     "2",
						"tags.Created":               "TF",
						"tags.Foo":                   "Bar",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"add_account", "build_content"},
			},
		},
	})
}

var AlicloudECSImagePipelineMap0 = map[string]string{}

func AlicloudECSImagePipelineBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
data "alicloud_zones" "default" {}
data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones.0.id
}
data "alicloud_images" "default" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}
data "alicloud_instance_types" "default" {
  image_id = data.alicloud_images.default.ids.0
}
data "alicloud_account" "default" {}
`, name)
}

func TestUnitAlicloudEcsImagePipeline(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ecs_image_pipeline"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ecs_image_pipeline"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"base_image_type":            "CreateEcsImagePipelineValue",
		"name":                       "CreateEcsImagePipelineValue",
		"description":                "CreateEcsImagePipelineValue",
		"system_disk_size":           40,
		"vswitch_id":                 "CreateEcsImagePipelineValue",
		"image_name":                 "CreateEcsImagePipelineValue",
		"delete_instance_on_failure": true,
		"internet_max_bandwidth_out": 20,
		"to_region_id":               []string{"CreateEcsImagePipelineValue"},
		"add_account":                []string{"CreateEcsImagePipelineValue"},
		"instance_type":              "CreateEcsImagePipelineValue",
		"base_image":                 "CreateEcsImagePipelineValue",
		"build_content":              "CreateEcsImagePipelineValue",
		"resource_group_id":          "CreateEcsImagePipelineValue",
		"tags":                       map[string]interface{}{"TestKey": "TestValue"},
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
		"ImagePipeline": map[string]interface{}{
			"ImagePipelineSet": []interface{}{
				map[string]interface{}{
					"Description":     "CreateEcsImagePipelineValue",
					"BaseImageType":   "CreateEcsImagePipelineValue",
					"ResourceGroupId": "CreateEcsImagePipelineValue",
					"SystemDiskSize":  40,
					"ImagePipelineId": "CreateEcsImagePipelineValue",
					"VSwitchId":       "CreateEcsImagePipelineValue",
					"ToRegionIds": map[string]interface{}{
						"ToRegionId": []interface{}{
							"CreateEcsImagePipelineValue",
						},
					},
					"Name":                    "CreateEcsImagePipelineValue",
					"DeleteInstanceOnFailure": true,
					"ImageName":               "CreateEcsImagePipelineValue",
					"AddAccounts": map[string]interface{}{
						"AddAccount": []interface{}{
							"CreateEcsImagePipelineValue",
						},
					},
					"InternetMaxBandwidthOut": 20,
					"CreationTime":            "CreateEcsImagePipelineValue",
					"InstanceType":            "CreateEcsImagePipelineValue",
					"BaseImage":               "CreateEcsImagePipelineValue",
					"BuildContent":            "CreateEcsImagePipelineValue",
					"Tags": map[string]interface{}{
						"Tag": []interface{}{
							map[string]interface{}{
								"TagKey":   "TestKey",
								"TagValue": "TestValue",
							},
						},
					},
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ecs_image_pipeline", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEcsImagePipelineCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateImagePipeline" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						_, _ = successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEcsImagePipelineCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_image_pipeline"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	attributesDiff := map[string]interface{}{}
	diff, err := newInstanceDiff("alicloud_ecs_image_pipeline", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_image_pipeline"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeImagePipelines" {
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
		err := resourceAlicloudEcsImagePipelineRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEcsImagePipelineDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_ecs_image_pipeline", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_image_pipeline"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteImagePipeline" {
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
		err := resourceAlicloudEcsImagePipelineDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
