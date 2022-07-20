package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
		"alicloud_actiontrail_trail",
		&resource.Sweeper{
			Name: "alicloud_actiontrail_trail",
			F:    testSweepActiontrailTrail,
		})
}

func testSweepActiontrailTrail(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.ApiGatewayNoSupportedRegions) {
		log.Printf("[INFO] Skipping API Gateway unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}
	request := make(map[string]interface{})
	var response map[string]interface{}
	action := "DescribeTrails"
	conn, err := client.NewActiontrailClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-07-06"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_actiontrail_trails", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.TrailList", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TrailList", response)
	}
	sweeped := false
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
			log.Printf("[INFO] Skipping ActionTrail Trails: %s", item["Name"].(string))
			continue
		}
		sweeped = true
		action = "DeleteTrail"
		request := map[string]interface{}{
			"Name": item["Name"],
		}
		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-07-06"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to delete ActionTrail Trail (%s): %s", item["Name"].(string), err)
		}
		if sweeped {
			// Waiting 5 seconds to ensure these ActionTrail Trails have been deleted.
			time.Sleep(5 * time.Second)
		}
		log.Printf("[INFO] Delete ActionTrail Trail success: %s ", item["Name"].(string))
	}
	return nil
}

func TestAccAlicloudActiontrailTrail_basic(t *testing.T) {
	checkoutAccount(t, true)
	defer checkoutAccount(t, false)
	var v map[string]interface{}
	resourceId := "alicloud_actiontrail_trail.default"
	ra := resourceAttrInit(resourceId, ActiontrailTrailMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ActiontrailService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeActiontrailTrail")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccactiontrail%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ActiontrailTrailBasicdependence)
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
					"trail_name":         name,
					"oss_write_role_arn": "${data.alicloud_ram_roles.default.roles.0.arn}",
					"oss_bucket_name":    "${alicloud_oss_bucket.default.id}",
					"status":             "Disable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"trail_name":         name,
						"oss_write_role_arn": CHECKSET,
						"oss_bucket_name":    name,
						"status":             "Disable",
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
					"oss_write_role_arn": "${data.alicloud_ram_roles.update.roles.0.arn}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_write_role_arn": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"event_rw": "All",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_rw": "All",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"oss_bucket_name": "${alicloud_oss_bucket.default2.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_bucket_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"trail_region": "cn-beijing",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"trail_region": "cn-beijing",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"oss_bucket_name":    "${alicloud_oss_bucket.default.id}",
					"oss_write_role_arn": "${data.alicloud_ram_roles.default.roles.0.arn}",
					"trail_region":       "All",
					"event_rw":           "Write",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_bucket_name":    name,
						"oss_write_role_arn": CHECKSET,
						"trail_region":       "All",
						"event_rw":           "Write",
					}),
				),
			},
		},
	})
}

var ActiontrailTrailMap = map[string]string{}

func ActiontrailTrailBasicdependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	resource "alicloud_oss_bucket" "default" {
		bucket  = "${var.name}"
	}

	resource "alicloud_oss_bucket" "default2" {
		bucket  = "${var.name}-update"
	}

	data "alicloud_ram_roles" "default" {
		name_regex = "AliyunActionTrailDefaultRole"
	}

	data "alicloud_ram_roles" "update" {
		name_regex = "AliyunServiceRoleForActionTrail"
	}
`, name)
}

func TestUnitAlicloudActiontrailTrail(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_actiontrail_trail"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_actiontrail_trail"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"trail_name":            "trail_name",
		"oss_write_role_arn":    "oss_write_role_arn",
		"oss_bucket_name":       "oss_bucket_name",
		"status":                "Disable",
		"event_rw":              "event_rw",
		"is_organization_trail": true,
		"oss_key_prefix":        "oss_key_prefix",
		"sls_project_arn":       "sls_project_arn",
		"sls_write_role_arn":    "sls_write_role_arn",
		"trail_region":          "trail_region",
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"TrailList": []interface{}{
			map[string]interface{}{
				"EventRW":             "event_rw",
				"IsOrganizationTrail": true,
				"OssBucketName":       "MockName",
				"OssKeyPrefix":        "oss_key_prefix",
				"OssWriteRoleArn":     "oss_write_role_arn",
				"SlsProjectArn":       "sls_project_arn",
				"SlsWriteRoleArn":     "sls_write_role_arn",
				"Status":              "Fresh",
				"TrailRegion":         "trail_region",
				"Name":                "MockName",
			},
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_actiontrail_trail", "MockName"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			result["Name"] = "MockName"
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"UpdateStopLoggingNormal": func(errorCode string) (map[string]interface{}, error) {
			result := map[string]interface{}{
				"TrailList": []interface{}{
					map[string]interface{}{
						"EventRW":             "event_rw",
						"IsOrganizationTrail": true,
						"OssBucketName":       "MockName",
						"OssKeyPrefix":        "oss_key_prefix",
						"OssWriteRoleArn":     "oss_write_role_arn",
						"SlsProjectArn":       "sls_project_arn",
						"SlsWriteRoleArn":     "sls_write_role_arn",
						"Status":              "Disable",
						"TrailRegion":         "trail_region",
						"Name":                "MockName",
					},
				},
				"Status": "Enable",
			}
			result["Name"] = "MockName"
			return result, nil
		},
		"UpdateStatusDisable": func(errorCode string) (map[string]interface{}, error) {
			result := map[string]interface{}{
				"TrailList": []interface{}{
					map[string]interface{}{
						"EventRW":             "event_rw",
						"IsOrganizationTrail": true,
						"OssBucketName":       "MockName",
						"OssKeyPrefix":        "oss_key_prefix",
						"OssWriteRoleArn":     "oss_write_role_arn",
						"SlsProjectArn":       "sls_project_arn",
						"SlsWriteRoleArn":     "sls_write_role_arn",
						"Status":              "Disable",
						"TrailRegion":         "trail_region",
						"Name":                "MockName",
					},
				},
				"Status": "Disable",
			}
			result["Name"] = "MockName"
			return result, nil
		},
		"UpdateStartLoggingNormal": func(errorCode string) (map[string]interface{}, error) {
			result := map[string]interface{}{
				"TrailList": []interface{}{
					map[string]interface{}{
						"EventRW":             "event_rw",
						"IsOrganizationTrail": true,
						"OssBucketName":       "MockName",
						"OssKeyPrefix":        "oss_key_prefix",
						"OssWriteRoleArn":     "oss_write_role_arn",
						"SlsProjectArn":       "sls_project_arn",
						"SlsWriteRoleArn":     "sls_write_role_arn",
						"Status":              "Enable",
						"TrailRegion":         "trail_region",
						"Name":                "MockName",
					},
				},
				"Status": "Disable",
			}
			result["Name"] = "MockName"
			return result, nil
		},
		"UpdateStatusEnable": func(errorCode string) (map[string]interface{}, error) {
			result := map[string]interface{}{
				"TrailList": []interface{}{
					map[string]interface{}{
						"EventRW":             "event_rw",
						"IsOrganizationTrail": true,
						"OssBucketName":       "MockName",
						"OssKeyPrefix":        "oss_key_prefix",
						"OssWriteRoleArn":     "oss_write_role_arn",
						"SlsProjectArn":       "sls_project_arn",
						"SlsWriteRoleArn":     "sls_write_role_arn",
						"Status":              "Enable",
						"TrailRegion":         "trail_region",
						"Name":                "MockName",
					},
				},
				"Status": "Enable",
			}
			result["Name"] = "MockName"
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewActiontrailClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudActiontrailTrailCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("InsufficientBucketPolicyException")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudActiontrailTrailCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudActiontrailTrailCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("MockName")
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewActiontrailClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAlicloudActiontrailTrailUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateTrailAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"event_rw", "oss_bucket_name", "oss_key_prefix", "oss_write_role_arn", "sls_project_arn", "sls_write_role_arn", "trail_region", "sls_project_arn", "sls_write_role_arn", "oss_bucket_name", "oss_write_role_arn"} {
			switch p["alicloud_actiontrail_trail"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_actiontrail_trail"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("InsufficientBucketPolicyException")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAlicloudActiontrailTrailUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateTrailNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"event_rw", "oss_bucket_name", "oss_key_prefix", "oss_write_role_arn", "sls_project_arn", "sls_write_role_arn", "trail_region", "sls_project_arn", "sls_write_role_arn", "oss_bucket_name", "oss_write_role_arn"} {
			switch p["alicloud_actiontrail_trail"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_actiontrail_trail"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAlicloudActiontrailTrailUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("UpdateStopLoggingAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"status"} {
			switch p["alicloud_actiontrail_trail"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: "Disable"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_actiontrail_trail"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("InsufficientBucketPolicyException")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateStatusDisable"]("")
		})
		patcheDescribeActiontrailTrail := gomonkey.ApplyMethod(reflect.TypeOf(&ActiontrailService{}), "DescribeActiontrailTrail", func(*ActiontrailService, string) (map[string]interface{}, error) {
			return responseMock["UpdateStopLoggingNormal"]("")
		})
		err := resourceAlicloudActiontrailTrailUpdate(resourceData1, rawClient)
		patches.Reset()
		patcheDescribeActiontrailTrail.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateStopLoggingNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"status"} {
			switch p["alicloud_actiontrail_trail"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: "Disable"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_actiontrail_trail"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateStatusDisable"]("")
		})
		patcheDescribeActiontrailTrail := gomonkey.ApplyMethod(reflect.TypeOf(&ActiontrailService{}), "DescribeActiontrailTrail", func(*ActiontrailService, string) (map[string]interface{}, error) {
			return responseMock["UpdateStopLoggingNormal"]("")
		})
		patchActiontrailTrailStateRefreshFunc := gomonkey.ApplyMethod(reflect.TypeOf(&ActiontrailService{}), "ActiontrailTrailStateRefreshFunc", func(*ActiontrailService, string, []string) resource.StateRefreshFunc {
			return func() (interface{}, string, error) {
				object := map[string]interface{}{
					"EventRW":             "event_rw",
					"IsOrganizationTrail": true,
					"OssBucketName":       "MockName",
					"OssKeyPrefix":        "oss_key_prefix",
					"OssWriteRoleArn":     "oss_write_role_arn",
					"SlsProjectArn":       "sls_project_arn",
					"SlsWriteRoleArn":     "sls_write_role_arn",
					"Status":              "Disable",
					"TrailRegion":         "trail_region",
					"Name":                "MockName",
				}
				return object, "Disable", nil
			}
		})
		err := resourceAlicloudActiontrailTrailUpdate(resourceData1, rawClient)
		patches.Reset()
		patcheDescribeActiontrailTrail.Reset()
		patchActiontrailTrailStateRefreshFunc.Reset()
		assert.Nil(t, err)
	})

	t.Run("UpdateStartLoggingAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"status"} {
			switch p["alicloud_actiontrail_trail"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: "Enable"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_actiontrail_trail"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("InsufficientBucketPolicyException")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateStatusEnable"]("")
		})
		patcheDescribeActiontrailTrail := gomonkey.ApplyMethod(reflect.TypeOf(&ActiontrailService{}), "DescribeActiontrailTrail", func(*ActiontrailService, string) (map[string]interface{}, error) {
			return responseMock["UpdateStartLoggingNormal"]("")
		})
		err := resourceAlicloudActiontrailTrailUpdate(resourceData1, rawClient)
		patches.Reset()
		patcheDescribeActiontrailTrail.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateStartLoggingNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"status"} {
			switch p["alicloud_actiontrail_trail"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: "Enable"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_actiontrail_trail"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateStatusEnable"]("")
		})
		patcheDescribeActiontrailTrail := gomonkey.ApplyMethod(reflect.TypeOf(&ActiontrailService{}), "DescribeActiontrailTrail", func(*ActiontrailService, string) (map[string]interface{}, error) {
			return responseMock["UpdateStartLoggingNormal"]("")
		})
		patchActiontrailTrailStateRefreshFunc := gomonkey.ApplyMethod(reflect.TypeOf(&ActiontrailService{}), "ActiontrailTrailStateRefreshFunc", func(*ActiontrailService, string, []string) resource.StateRefreshFunc {
			return func() (interface{}, string, error) {
				object := map[string]interface{}{
					"EventRW":             "event_rw",
					"IsOrganizationTrail": true,
					"OssBucketName":       "MockName",
					"OssKeyPrefix":        "oss_key_prefix",
					"OssWriteRoleArn":     "oss_write_role_arn",
					"SlsProjectArn":       "sls_project_arn",
					"SlsWriteRoleArn":     "sls_write_role_arn",
					"Status":              "Enable",
					"TrailRegion":         "trail_region",
					"Name":                "MockName",
				}
				return object, "Enable", nil
			}
		})
		err := resourceAlicloudActiontrailTrailUpdate(resourceData1, rawClient)
		patches.Reset()
		patcheDescribeActiontrailTrail.Reset()
		patchActiontrailTrailStateRefreshFunc.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewActiontrailClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudActiontrailTrailDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudActiontrailTrailDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudActiontrailTrailDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("DeleteNonRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudActiontrailTrailDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	//Read
	t.Run("ReadDescribeActiontrailTrailNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudActiontrailTrailRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeActiontrailTrailAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := true
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudActiontrailTrailRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}
