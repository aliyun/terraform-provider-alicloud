package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

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
	resource.AddTestSweepers("alicloud_kms_key", &resource.Sweeper{
		Name: "alicloud_kms_key",
		F:    testSweepKmsKey,
		Dependencies: []string{
			"alicloud_kms_alias",
		},
	})
}

func testSweepKmsKey(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	request := map[string]interface{}{
		"PageSize":   PageSizeXLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
		"Filters":    "[ {\"Key\":\"KeyState\", \"Values\":[\"Enabled\",\"Disabled\",\"PendingImport\"]} ]",
	}
	action := "ListKeys"

	var response map[string]interface{}
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kms_key", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Keys.Key", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Keys.Key", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			if !sweepAll() {
				if _, ok := item["Description"]; !ok {
					continue
				}
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(item["Description"].(string)), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Kms Key: %s", item["Description"].(string))
					continue
				}
			}
			action = "SetDeletionProtection"
			request := map[string]interface{}{
				"ProtectedResourceArn":     item["KeyArn"],
				"EnableDeletionProtection": false,
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to cancel Kms Key DeletionProtection %s (%s): %s", item["Description"], item["KeyId"], err)
			}

			action = "ScheduleKeyDeletion"
			request = map[string]interface{}{
				"KeyId":               item["KeyId"],
				"PendingWindowInDays": 7,
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Kms Key %s (%s): %s", item["Description"], item["KeyId"], err)
			}
		}
		if len(result) < PageSizeXLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudKMSKey_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, KmsKeyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKmsKey%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KmsKeyBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.KmsKeyHSMSupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":            name,
					"key_spec":               "Aliyun_AES_256",
					"protection_level":       "HSM",
					"pending_window_in_days": "7",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Key",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            name,
						"key_spec":               "Aliyun_AES_256",
						"protection_level":       "HSM",
						"pending_window_in_days": "7",
						"tags.%":                 "2",
						"tags.Created":           "TF",
						"tags.For":               "Key",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF_Update",
						"For":     "Key_Update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF_Update",
						"tags.For":     "Key_Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "from_terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "from_terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "Enabled",
					"rotation_interval":  "2678400s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Enabled",
						"rotation_interval":  "2678400s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":        name,
					"automatic_rotation": "Disabled",
					"rotation_interval":  REMOVEKEY,
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":        name,
						"automatic_rotation": "Disabled",
						"rotation_interval":  REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pending_window_in_days", "deletion_window_in_days", "is_enabled"},
			},
		},
	})
}

var KmsKeyMap = map[string]string{
	"arn":                 CHECKSET,
	"automatic_rotation":  "Disabled",
	"creation_date":       CHECKSET,
	"creator":             CHECKSET,
	"status":              "Enabled",
	"key_usage":           "ENCRYPT/DECRYPT",
	"last_rotation_date":  CHECKSET,
	"origin":              "Aliyun_KMS",
	"primary_key_version": CHECKSET,
	"protection_level":    "SOFTWARE",
}

func SkipTestAccAlicloudKMSKey_DKMS(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, KmsKeyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKmsKey%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KmsKeyBasicdependence)
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
					"description":            name,
					"key_spec":               "Aliyun_AES_256",
					"protection_level":       "HSM",
					"pending_window_in_days": "7",
					"dkms_instance_id":       os.Getenv("DKMS_INSTANCE_ID"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            name,
						"key_spec":               "Aliyun_AES_256",
						"protection_level":       "HSM",
						"pending_window_in_days": "7",
						"dkms_instance_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pending_window_in_days", "deletion_window_in_days", "is_enabled"},
			},
		},
	})
}

func KmsKeyBasicdependence(name string) string {
	return ""
}

func TestUnitAlicloudKMSKey(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"description":            "description",
		"key_spec":               "Aliyun_AES_256",
		"protection_level":       "HSM",
		"pending_window_in_days": 7,
		"key_usage":              "key_usage",
		"origin":                 "origin",
		"rotation_interval":      "rotation_interval",
		"automatic_rotation":     "Disabled",
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
		"KeyMetadata": map[string]interface{}{
			"Arn":                "arn",
			"AutomaticRotation":  "automatic_rotation",
			"Creator":            "creator",
			"CreationDate":       "creation_date",
			"DeleteDate":         "delete_date",
			"Description":        "description",
			"KeySpec":            "key_spec",
			"KeyUsage":           "key_usage",
			"LastRotationDate":   "last_rotation_date",
			"MaterialExpireTime": "material_expire_time",
			"NextRotationDate":   "next_rotation_date",
			"Origin":             "origin",
			"PrimaryKeyVersion":  "primary_key_version",
			"ProtectionLevel":    "protection_level",
			"RotationInterval":   "rotation_interval",
			"KeyState":           "status",
			"KeyId":              "MockKeyId",
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
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_kms_key", "MockKeyId"))
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
		"UpdateDisableKey": func(errorCode string) (map[string]interface{}, error) {
			result := map[string]interface{}{
				"KeyState": "Disable",
				"KeyId":    "MockKeyId",
			}
			return result, nil
		},
		"UpdateEnableKey": func(errorCode string) (map[string]interface{}, error) {
			result := map[string]interface{}{
				"KeyState": "Enable",
				"KeyId":    "MockKeyId",
			}
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewKmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudKmsKeyCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
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
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudKmsKeyCreate(d, rawClient)
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
		err := resourceAliCloudKmsKeyCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("MockKeyId")
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewKmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAliCloudKmsKeyUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateKeyDescriptionAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"description"} {
			switch p["alicloud_kms_key"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
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
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudKmsKeyUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateKeyDescriptionNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"description"} {
			switch p["alicloud_kms_key"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, diff)
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
		err := resourceAliCloudKmsKeyUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("UpdateRotationPolicyAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"automatic_rotation", "rotation_interval"} {
			switch p["alicloud_kms_key"].Schema[key].Type {
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
		resourceData1, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
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
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudKmsKeyUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateRotationPolicyNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"automatic_rotation", "rotation_interval"} {
			switch p["alicloud_kms_key"].Schema[key].Type {
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
		resourceData1, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, diff)
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
		err := resourceAliCloudKmsKeyUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("UpdateDisableKeyAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"status"} {
			switch p["alicloud_kms_key"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: "Enabled"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: "false"})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
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
			return responseMock["UpdateNormal"]("")
		})
		patcheDescribeKmsKey := gomonkey.ApplyMethod(reflect.TypeOf(&KmsService{}), "DescribeKmsKey", func(*KmsService, string) (map[string]interface{}, error) {
			return responseMock["UpdateDisableKey"]("")
		})
		err := resourceAliCloudKmsKeyUpdate(resourceData1, rawClient)
		patches.Reset()
		patcheDescribeKmsKey.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateDisableKeyNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"status"} {
			switch p["alicloud_kms_key"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: "Enabled"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: "false"})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, diff)
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
		patcheDescribeKmsKey := gomonkey.ApplyMethod(reflect.TypeOf(&KmsService{}), "DescribeKmsKey", func(*KmsService, string) (map[string]interface{}, error) {
			return responseMock["UpdateDisableKey"]("")
		})
		err := resourceAliCloudKmsKeyUpdate(resourceData1, rawClient)
		patches.Reset()
		patcheDescribeKmsKey.Reset()
		assert.Nil(t, err)
	})

	t.Run("UpdateEnableKeyAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"status"} {
			switch p["alicloud_kms_key"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: "Disabled"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: "true"})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
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
			return responseMock["UpdateNormal"]("")
		})
		patcheDescribeKmsKey := gomonkey.ApplyMethod(reflect.TypeOf(&KmsService{}), "DescribeKmsKey", func(*KmsService, string) (map[string]interface{}, error) {
			return responseMock["UpdateEnableKey"]("")
		})
		err := resourceAliCloudKmsKeyUpdate(resourceData1, rawClient)
		patches.Reset()
		patcheDescribeKmsKey.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateEnableKeyNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"status"} {
			switch p["alicloud_kms_key"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: "Disabled"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: "true"})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, diff)
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
		patcheDescribeKmsKey := gomonkey.ApplyMethod(reflect.TypeOf(&KmsService{}), "DescribeKmsKey", func(*KmsService, string) (map[string]interface{}, error) {
			return responseMock["UpdateEnableKey"]("")
		})
		err := resourceAliCloudKmsKeyUpdate(resourceData1, rawClient)
		patches.Reset()
		patcheDescribeKmsKey.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewKmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudKmsKeyDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		resourceData, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, nil)
		resourceData.SetId(d.Id())
		_ = resourceData.Set("pending_window_in_days", 7)
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
		err := resourceAliCloudKmsKeyDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		resourceData, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, nil)
		resourceData.SetId(d.Id())
		_ = resourceData.Set("deletion_window_in_days", 7)
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
		err := resourceAliCloudKmsKeyDelete(resourceData, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeKmsKeyNotFound", func(t *testing.T) {
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
		err := resourceAliCloudKmsKeyRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeKmsKeyAbnormal", func(t *testing.T) {
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
		err := resourceAliCloudKmsKeyRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

// Test Kms Key. >>> Resource test cases, automatically generated.
// Case KeyPolicy 6388
func TestAccAliCloudKmsKey_basic6388(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsKeyMap6388)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmskey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsKeyBasicDependence6388)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level": "SOFTWARE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level": "SOFTWARE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "key description example",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "key description example",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rotation_interval": "604800s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rotation_interval": "604800s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": "{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"kms:*\\\"],\\\"Resource\\\":[\\\"*\\\"],\\\"Effect\\\":\\\"Allow\\\",\\\"Principal\\\":{\\\"RAM\\\":[\\\"acs:ram::1192853035118460:*\\\"]},\\\"Sid\\\":\\\"kms default key policy\\\"}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"kms:*\"],\"Resource\":[\"*\"],\"Effect\":\"Allow\",\"Principal\":{\"RAM\":[\"acs:ram::1192853035118460:*\"]},\"Sid\":\"kms default key policy\"}]}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test111",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rotation_interval": "1209600s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rotation_interval": "1209600s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": "{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"kms:*\\\"],\\\"Resource\\\":[\\\"*\\\"],\\\"Effect\\\":\\\"Allow\\\",\\\"Principal\\\":{\\\"RAM\\\":[\\\"acs:ram::1192853035118460:*\\\"]},\\\"Sid\\\":\\\"kms key policy\\\"}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"kms:*\"],\"Resource\":[\"*\"],\"Effect\":\"Allow\",\"Principal\":{\"RAM\":[\"acs:ram::1192853035118460:*\"]},\"Sid\":\"kms key policy\"}]}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": "{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"kms:*\\\"],\\\"Resource\\\":[\\\"*\\\"],\\\"Effect\\\":\\\"Allow\\\",\\\"Principal\\\":{\\\"RAM\\\":[\\\"acs:ram::1192853035118460:*\\\"]},\\\"Sid\\\":\\\"kms policy\\\"}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"kms:*\"],\"Resource\":[\"*\"],\"Effect\":\"Allow\",\"Principal\":{\"RAM\":[\"acs:ram::1192853035118460:*\"]},\"Sid\":\"kms policy\"}]}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "PendingDeletion",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "PendingDeletion",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level":          "SOFTWARE",
					"description":               "key description example",
					"key_spec":                  "Aliyun_AES_256",
					"key_usage":                 "ENCRYPT/DECRYPT",
					"rotation_interval":         "604800s",
					"dkms_instance_id":          "${alicloud_kms_instance.create-instance.id}",
					"origin":                    "Aliyun_KMS",
					"policy":                    "{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"kms:*\\\"],\\\"Resource\\\":[\\\"*\\\"],\\\"Effect\\\":\\\"Allow\\\",\\\"Principal\\\":{\\\"RAM\\\":[\\\"acs:ram::1192853035118460:*\\\"]},\\\"Sid\\\":\\\"kms default key policy\\\"}]}",
					"enable_automatic_rotation": "true",
					"automatic_rotation":        "Enabled",
					"status":                    "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level":          "SOFTWARE",
						"description":               "key description example",
						"key_spec":                  "Aliyun_AES_256",
						"key_usage":                 "ENCRYPT/DECRYPT",
						"rotation_interval":         "604800s",
						"dkms_instance_id":          CHECKSET,
						"origin":                    "Aliyun_KMS",
						"policy":                    "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"kms:*\"],\"Resource\":[\"*\"],\"Effect\":\"Allow\",\"Principal\":{\"RAM\":[\"acs:ram::1192853035118460:*\"]},\"Sid\":\"kms default key policy\"}]}",
						"enable_automatic_rotation": "true",
						"automatic_rotation":        "Enabled",
						"status":                    "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
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
				ImportStateVerifyIgnore: []string{"certificate_id", "enable_automatic_rotation", "pending_window_in_days", "policy_name", "secret_name"},
			},
		},
	})
}

var AlicloudKmsKeyMap6388 = map[string]string{
	"status": CHECKSET,
}

func AlicloudKmsKeyBasicDependence6388(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "create-vpc" {
  vpc_name = var.name
}

resource "alicloud_vswitch" "vswitch-k" {
  vpc_id     = alicloud_vpc.create-vpc.id
  cidr_block = "172.16.1.0/24"
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch-j" {
  vpc_id     = alicloud_vpc.create-vpc.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_kms_instance" "create-instance" {
  vpc_num         = "1"
  key_num         = "1000"
  product_type    = "kms_ddi_public_cn"
  secret_num      = "100"
  product_version = "3"
  renew_status    = "AutoRenewal"
  vpc_id          = alicloud_vpc.create-vpc.id
  vswitch_ids     = ["${alicloud_vswitch.vswitch-j.id}"]
  zone_ids        = ["${alicloud_vswitch.vswitch-j.zone_id}", "${alicloud_vswitch.vswitch-k.zone_id}"]
  spec            = "1000"
  renew_period    = "1"
}


`, name)
}

// Case KeyPolicy_副本1712829083140——ceshi 6495
func TestAccAliCloudKmsKey_basic6495(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsKeyMap6495)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmskey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsKeyBasicDependence6495)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level": "SOFTWARE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level": "SOFTWARE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "key description example",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "key description example",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rotation_interval": "604800s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rotation_interval": "604800s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": "{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"kms:*\\\"],\\\"Resource\\\":[\\\"*\\\"],\\\"Effect\\\":\\\"Allow\\\",\\\"Principal\\\":{\\\"RAM\\\":[\\\"acs:ram::1117600963847258:*\\\"]},\\\"Sid\\\":\\\"kms default key policy\\\"}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"kms:*\"],\"Resource\":[\"*\"],\"Effect\":\"Allow\",\"Principal\":{\"RAM\":[\"acs:ram::1117600963847258:*\"]},\"Sid\":\"kms default key policy\"}]}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test111",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": "{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"kms:*\\\"],\\\"Resource\\\":[\\\"*\\\"],\\\"Effect\\\":\\\"Allow\\\",\\\"Principal\\\":{\\\"RAM\\\":[\\\"acs:ram::1117600963847258:*\\\"]},\\\"Sid\\\":\\\"kms default policy\\\"}]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"kms:*\"],\"Resource\":[\"*\"],\"Effect\":\"Allow\",\"Principal\":{\"RAM\":[\"acs:ram::1117600963847258:*\"]},\"Sid\":\"kms default policy\"}]}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "PendingDeletion",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "PendingDeletion",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level":          "SOFTWARE",
					"description":               "key description example",
					"key_spec":                  "Aliyun_AES_256",
					"key_usage":                 "ENCRYPT/DECRYPT",
					"rotation_interval":         "604800s",
					"dkms_instance_id":          "kst-phzz61dbabacquz826y29f",
					"origin":                    "Aliyun_KMS",
					"policy":                    "{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"kms:*\\\"],\\\"Resource\\\":[\\\"*\\\"],\\\"Effect\\\":\\\"Allow\\\",\\\"Principal\\\":{\\\"RAM\\\":[\\\"acs:ram::1117600963847258:*\\\"]},\\\"Sid\\\":\\\"kms default key policy\\\"}]}",
					"enable_automatic_rotation": "true",
					"status":                    "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level":          "SOFTWARE",
						"description":               "key description example",
						"key_spec":                  "Aliyun_AES_256",
						"key_usage":                 "ENCRYPT/DECRYPT",
						"rotation_interval":         "604800s",
						"dkms_instance_id":          "kst-phzz61dbabacquz826y29f",
						"origin":                    "Aliyun_KMS",
						"policy":                    "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"kms:*\"],\"Resource\":[\"*\"],\"Effect\":\"Allow\",\"Principal\":{\"RAM\":[\"acs:ram::1117600963847258:*\"]},\"Sid\":\"kms default key policy\"}]}",
						"enable_automatic_rotation": "true",
						"status":                    "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
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
				ImportStateVerifyIgnore: []string{"certificate_id", "enable_automatic_rotation", "pending_window_in_days", "policy_name", "secret_name"},
			},
		},
	})
}

var AlicloudKmsKeyMap6495 = map[string]string{
	"status": CHECKSET,
}

func AlicloudKmsKeyBasicDependence6495(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case KeyPolicy_副本1711981694735 6410
func TestAccAliCloudKmsKey_basic6410(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsKeyMap6410)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmskey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsKeyBasicDependence6410)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level": "SOFTWARE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level": "SOFTWARE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "key description example",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "key description example",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rotation_interval": "604800s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rotation_interval": "604800s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": "{     \\\"Version\\\": \\\"1\\\",     \\\"Statement\\\": [         {             \\\"Action\\\": [                 \\\"kms:*\\\"             ],             \\\"Effect\\\": \\\"Allow\\\",             \\\"Principal\\\": {                 \\\"RAM\\\": [                     \\\"acs:ram::1117600963847258:*\\\"                 ]             },             \\\"Resource\\\": [                 \\\"*\\\"             ],             \\\"Sid\\\": \\\"kms default key policy\\\"         }     ] }",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": "{     \"Version\": \"1\",     \"Statement\": [         {             \"Action\": [                 \"kms:*\"             ],             \"Effect\": \"Allow\",             \"Principal\": {                 \"RAM\": [                     \"acs:ram::1117600963847258:*\"                 ]             },             \"Resource\": [                 \"*\"             ],             \"Sid\": \"kms default key policy\"         }     ] }",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level":  "SOFTWARE",
					"description":       "key description example",
					"key_spec":          "Aliyun_AES_256",
					"key_usage":         "ENCRYPT/DECRYPT",
					"rotation_interval": "604800s",
					"dkms_instance_id":  "kst-phzz61dbabacquz826y29f",
					"origin":            "Aliyun_KMS",
					"policy":            "{     \\\"Version\\\": \\\"1\\\",     \\\"Statement\\\": [         {             \\\"Action\\\": [                 \\\"kms:*\\\"             ],             \\\"Effect\\\": \\\"Allow\\\",             \\\"Principal\\\": {                 \\\"RAM\\\": [                     \\\"acs:ram::1117600963847258:*\\\"                 ]             },             \\\"Resource\\\": [                 \\\"*\\\"             ],             \\\"Sid\\\": \\\"kms default key policy\\\"         }     ] }",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level":  "SOFTWARE",
						"description":       "key description example",
						"key_spec":          "Aliyun_AES_256",
						"key_usage":         "ENCRYPT/DECRYPT",
						"rotation_interval": "604800s",
						"dkms_instance_id":  "kst-phzz61dbabacquz826y29f",
						"origin":            "Aliyun_KMS",
						"policy":            "{     \"Version\": \"1\",     \"Statement\": [         {             \"Action\": [                 \"kms:*\"             ],             \"Effect\": \"Allow\",             \"Principal\": {                 \"RAM\": [                     \"acs:ram::1117600963847258:*\"                 ]             },             \"Resource\": [                 \"*\"             ],             \"Sid\": \"kms default key policy\"         }     ] }",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
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
				ImportStateVerifyIgnore: []string{"certificate_id", "enable_automatic_rotation", "pending_window_in_days", "policy_name", "secret_name"},
			},
		},
	})
}

var AlicloudKmsKeyMap6410 = map[string]string{
	"status": CHECKSET,
}

func AlicloudKmsKeyBasicDependence6410(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 全生命周期Key 4863
func TestAccAliCloudKmsKey_basic4863(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsKeyMap4863)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmskey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsKeyBasicDependence4863)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level": "SOFTWARE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level": "SOFTWARE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "key description example",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "key description example",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "key description example11111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "key description example11111",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "key description example",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "key description example",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "key description example11111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "key description example11111",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level":   "SOFTWARE",
					"description":        "key description example",
					"key_spec":           "Aliyun_AES_256",
					"key_usage":          "ENCRYPT/DECRYPT",
					"automatic_rotation": "Disabled",
					"dkms_instance_id":   "${alicloud_kms_instance.kms.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level":   "SOFTWARE",
						"description":        "key description example",
						"key_spec":           "Aliyun_AES_256",
						"key_usage":          "ENCRYPT/DECRYPT",
						"automatic_rotation": "Disabled",
						"dkms_instance_id":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
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
				ImportStateVerifyIgnore: []string{"certificate_id", "enable_automatic_rotation", "pending_window_in_days", "policy_name", "secret_name"},
			},
		},
	})
}

var AlicloudKmsKeyMap4863 = map[string]string{
	"status": CHECKSET,
}

func AlicloudKmsKeyBasicDependence4863(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vs1" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_zones.default.zones.0.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_vswitch" "vs2" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_zones.default.zones.1.id
  cidr_block   = "172.18.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_kms_instance" "kms" {
  vpc_id          = alicloud_vpc.vpc.id
  vpc_num         = "1"
  key_num         = "1000"
  vswitch_ids     = ["${alicloud_vswitch.vs1.id}", "${alicloud_vswitch.vs2.id}"]
  product_type    = "kms_ddi_public_cn"
  zone_ids        = ["${alicloud_vswitch.vs1.zone_id}", "${alicloud_vswitch.vs2.zone_id}"]
  secret_num      = "0"
  product_version = "3"
  spec            = "1000"
  renew_status    = "AutoRenewal"
}


`, name)
}

// Case Key 3386
func TestAccAliCloudKmsKey_basic3386(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsKeyMap3386)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmskey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsKeyBasicDependence3386)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level": "SOFTWARE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level": "SOFTWARE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "key description example",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "key description example",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rotation_interval": "604800s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rotation_interval": "604800s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "key description example11111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "key description example11111",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rotation_interval": "6048000s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rotation_interval": "6048000s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "key description example",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "key description example",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rotation_interval": "604800s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rotation_interval": "604800s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "key description example11111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "key description example11111",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rotation_interval": "6048000s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rotation_interval": "6048000s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level":   "SOFTWARE",
					"description":        "key description example",
					"key_spec":           "Aliyun_AES_256",
					"key_usage":          "ENCRYPT/DECRYPT",
					"automatic_rotation": "Enabled",
					"rotation_interval":  "604800s",
					"dkms_instance_id":   "kst-phzz61dbabacquz826y29f",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level":   "SOFTWARE",
						"description":        "key description example",
						"key_spec":           "Aliyun_AES_256",
						"key_usage":          "ENCRYPT/DECRYPT",
						"automatic_rotation": "Enabled",
						"rotation_interval":  "604800s",
						"dkms_instance_id":   "kst-phzz61dbabacquz826y29f",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
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
				ImportStateVerifyIgnore: []string{"certificate_id", "enable_automatic_rotation", "pending_window_in_days", "policy_name", "secret_name"},
			},
		},
	})
}

var AlicloudKmsKeyMap3386 = map[string]string{
	"status": CHECKSET,
}

func AlicloudKmsKeyBasicDependence3386(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 全生命周期_副本1676539573013 2532
func TestAccAliCloudKmsKey_basic2532(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsKeyMap2532)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmskey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsKeyBasicDependence2532)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"origin": "Aliyun_KMS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"origin": "Aliyun_KMS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"origin":             "Aliyun_KMS",
					"protection_level":   "SOFTWARE",
					"description":        "test",
					"key_spec":           "Aliyun_AES_256",
					"key_usage":          "ENCRYPT/DECRYPT",
					"automatic_rotation": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"origin":             "Aliyun_KMS",
						"protection_level":   "SOFTWARE",
						"description":        "test",
						"key_spec":           "Aliyun_AES_256",
						"key_usage":          "ENCRYPT/DECRYPT",
						"automatic_rotation": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
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
				ImportStateVerifyIgnore: []string{"certificate_id", "enable_automatic_rotation", "pending_window_in_days", "policy_name", "secret_name"},
			},
		},
	})
}

var AlicloudKmsKeyMap2532 = map[string]string{
	"status": CHECKSET,
}

func AlicloudKmsKeyBasicDependence2532(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 全生命周期 1980
func TestAccAliCloudKmsKey_basic1980(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsKeyMap1980)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmskey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsKeyBasicDependence1980)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"origin": "Aliyun_KMS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"origin": "Aliyun_KMS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"origin":             "Aliyun_KMS",
					"protection_level":   "SOFTWARE",
					"description":        "test",
					"key_spec":           "Aliyun_AES_256",
					"key_usage":          "ENCRYPT/DECRYPT",
					"automatic_rotation": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"origin":             "Aliyun_KMS",
						"protection_level":   "SOFTWARE",
						"description":        "test",
						"key_spec":           "Aliyun_AES_256",
						"key_usage":          "ENCRYPT/DECRYPT",
						"automatic_rotation": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
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
				ImportStateVerifyIgnore: []string{"certificate_id", "enable_automatic_rotation", "pending_window_in_days", "policy_name", "secret_name"},
			},
		},
	})
}

var AlicloudKmsKeyMap1980 = map[string]string{
	"status": CHECKSET,
}

func AlicloudKmsKeyBasicDependence1980(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case RDK接入 1323
func TestAccAliCloudKmsKey_basic1323(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsKeyMap1323)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmskey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsKeyBasicDependence1323)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"origin": "Aliyun_KMS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"origin": "Aliyun_KMS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "rdk_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "rdk_test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"origin":             "Aliyun_KMS",
					"protection_level":   "SOFTWARE",
					"description":        "rdk_test",
					"key_spec":           "Aliyun_AES_256",
					"key_usage":          "ENCRYPT/DECRYPT",
					"automatic_rotation": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"origin":             "Aliyun_KMS",
						"protection_level":   "SOFTWARE",
						"description":        "rdk_test",
						"key_spec":           "Aliyun_AES_256",
						"key_usage":          "ENCRYPT/DECRYPT",
						"automatic_rotation": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
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
				ImportStateVerifyIgnore: []string{"certificate_id", "enable_automatic_rotation", "pending_window_in_days", "policy_name", "secret_name"},
			},
		},
	})
}

var AlicloudKmsKeyMap1323 = map[string]string{
	"status": CHECKSET,
}

func AlicloudKmsKeyBasicDependence1323(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case KeyPolicy 6388  twin
func TestAccAliCloudKmsKey_basic6388_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsKeyMap6388)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmskey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsKeyBasicDependence6388)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level":          "SOFTWARE",
					"description":               "key description example",
					"key_spec":                  "Aliyun_AES_256",
					"key_usage":                 "ENCRYPT/DECRYPT",
					"rotation_interval":         "604800s",
					"dkms_instance_id":          "${alicloud_kms_instance.create-instance.id}",
					"origin":                    "Aliyun_KMS",
					"policy":                    "{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"kms:*\\\"],\\\"Resource\\\":[\\\"*\\\"],\\\"Effect\\\":\\\"Allow\\\",\\\"Principal\\\":{\\\"RAM\\\":[\\\"acs:ram::1192853035118460:*\\\"]},\\\"Sid\\\":\\\"kms default key policy\\\"}]}",
					"enable_automatic_rotation": "true",
					"automatic_rotation":        "Enabled",
					"status":                    "Enabled",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level":          "SOFTWARE",
						"description":               "key description example",
						"key_spec":                  "Aliyun_AES_256",
						"key_usage":                 "ENCRYPT/DECRYPT",
						"rotation_interval":         "604800s",
						"dkms_instance_id":          CHECKSET,
						"origin":                    "Aliyun_KMS",
						"policy":                    "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"kms:*\"],\"Resource\":[\"*\"],\"Effect\":\"Allow\",\"Principal\":{\"RAM\":[\"acs:ram::1192853035118460:*\"]},\"Sid\":\"kms default key policy\"}]}",
						"enable_automatic_rotation": "true",
						"automatic_rotation":        "Enabled",
						"status":                    "Enabled",
						"tags.%":                    "2",
						"tags.Created":              "TF",
						"tags.For":                  "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate_id", "enable_automatic_rotation", "pending_window_in_days", "policy_name", "secret_name"},
			},
		},
	})
}

// Case KeyPolicy_副本1712829083140——ceshi 6495  twin
func TestAccAliCloudKmsKey_basic6495_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsKeyMap6495)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmskey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsKeyBasicDependence6495)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level":          "SOFTWARE",
					"description":               "key description example",
					"key_spec":                  "Aliyun_AES_256",
					"key_usage":                 "ENCRYPT/DECRYPT",
					"rotation_interval":         "604800s",
					"dkms_instance_id":          "kst-phzz61dbabacquz826y29f",
					"origin":                    "Aliyun_KMS",
					"policy":                    "{\\\"Version\\\":\\\"1\\\",\\\"Statement\\\":[{\\\"Action\\\":[\\\"kms:*\\\"],\\\"Resource\\\":[\\\"*\\\"],\\\"Effect\\\":\\\"Allow\\\",\\\"Principal\\\":{\\\"RAM\\\":[\\\"acs:ram::1117600963847258:*\\\"]},\\\"Sid\\\":\\\"kms default key policy\\\"}]}",
					"enable_automatic_rotation": "true",
					"status":                    "Enabled",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level":          "SOFTWARE",
						"description":               "key description example",
						"key_spec":                  "Aliyun_AES_256",
						"key_usage":                 "ENCRYPT/DECRYPT",
						"rotation_interval":         "604800s",
						"dkms_instance_id":          "kst-phzz61dbabacquz826y29f",
						"origin":                    "Aliyun_KMS",
						"policy":                    "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"kms:*\"],\"Resource\":[\"*\"],\"Effect\":\"Allow\",\"Principal\":{\"RAM\":[\"acs:ram::1117600963847258:*\"]},\"Sid\":\"kms default key policy\"}]}",
						"enable_automatic_rotation": "true",
						"status":                    "Enabled",
						"tags.%":                    "2",
						"tags.Created":              "TF",
						"tags.For":                  "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate_id", "enable_automatic_rotation", "pending_window_in_days", "policy_name", "secret_name"},
			},
		},
	})
}

// Case KeyPolicy_副本1711981694735 6410  twin
func TestAccAliCloudKmsKey_basic6410_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsKeyMap6410)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmskey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsKeyBasicDependence6410)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level":  "SOFTWARE",
					"description":       "key description example",
					"key_spec":          "Aliyun_AES_256",
					"key_usage":         "ENCRYPT/DECRYPT",
					"rotation_interval": "604800s",
					"dkms_instance_id":  "kst-phzz61dbabacquz826y29f",
					"origin":            "Aliyun_KMS",
					"policy":            "{     \\\"Version\\\": \\\"1\\\",     \\\"Statement\\\": [         {             \\\"Action\\\": [                 \\\"kms:*\\\"             ],             \\\"Effect\\\": \\\"Allow\\\",             \\\"Principal\\\": {                 \\\"RAM\\\": [                     \\\"acs:ram::1117600963847258:*\\\"                 ]             },             \\\"Resource\\\": [                 \\\"*\\\"             ],             \\\"Sid\\\": \\\"kms default key policy\\\"         }     ] }",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level":  "SOFTWARE",
						"description":       "key description example",
						"key_spec":          "Aliyun_AES_256",
						"key_usage":         "ENCRYPT/DECRYPT",
						"rotation_interval": "604800s",
						"dkms_instance_id":  "kst-phzz61dbabacquz826y29f",
						"origin":            "Aliyun_KMS",
						"policy":            "{     \"Version\": \"1\",     \"Statement\": [         {             \"Action\": [                 \"kms:*\"             ],             \"Effect\": \"Allow\",             \"Principal\": {                 \"RAM\": [                     \"acs:ram::1117600963847258:*\"                 ]             },             \"Resource\": [                 \"*\"             ],             \"Sid\": \"kms default key policy\"         }     ] }",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate_id", "enable_automatic_rotation", "pending_window_in_days", "policy_name", "secret_name"},
			},
		},
	})
}

// Case 全生命周期Key 4863  twin
func TestAccAliCloudKmsKey_basic4863_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsKeyMap4863)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmskey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsKeyBasicDependence4863)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level":   "SOFTWARE",
					"description":        "key description example",
					"key_spec":           "Aliyun_AES_256",
					"key_usage":          "ENCRYPT/DECRYPT",
					"automatic_rotation": "Disabled",
					"dkms_instance_id":   "${alicloud_kms_instance.kms.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level":   "SOFTWARE",
						"description":        "key description example",
						"key_spec":           "Aliyun_AES_256",
						"key_usage":          "ENCRYPT/DECRYPT",
						"automatic_rotation": "Disabled",
						"dkms_instance_id":   CHECKSET,
						"tags.%":             "2",
						"tags.Created":       "TF",
						"tags.For":           "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate_id", "enable_automatic_rotation", "pending_window_in_days", "policy_name", "secret_name"},
			},
		},
	})
}

// Case Key 3386  twin
func TestAccAliCloudKmsKey_basic3386_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsKeyMap3386)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmskey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsKeyBasicDependence3386)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"protection_level":   "SOFTWARE",
					"description":        "key description example",
					"key_spec":           "Aliyun_AES_256",
					"key_usage":          "ENCRYPT/DECRYPT",
					"automatic_rotation": "Enabled",
					"rotation_interval":  "604800s",
					"dkms_instance_id":   "kst-phzz61dbabacquz826y29f",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protection_level":   "SOFTWARE",
						"description":        "key description example",
						"key_spec":           "Aliyun_AES_256",
						"key_usage":          "ENCRYPT/DECRYPT",
						"automatic_rotation": "Enabled",
						"rotation_interval":  "604800s",
						"dkms_instance_id":   "kst-phzz61dbabacquz826y29f",
						"tags.%":             "2",
						"tags.Created":       "TF",
						"tags.For":           "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate_id", "enable_automatic_rotation", "pending_window_in_days", "policy_name", "secret_name"},
			},
		},
	})
}

// Case 全生命周期_副本1676539573013 2532  twin
func TestAccAliCloudKmsKey_basic2532_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsKeyMap2532)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmskey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsKeyBasicDependence2532)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"origin":             "Aliyun_KMS",
					"protection_level":   "SOFTWARE",
					"description":        "test",
					"key_spec":           "Aliyun_AES_256",
					"key_usage":          "ENCRYPT/DECRYPT",
					"automatic_rotation": "false",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"origin":             "Aliyun_KMS",
						"protection_level":   "SOFTWARE",
						"description":        "test",
						"key_spec":           "Aliyun_AES_256",
						"key_usage":          "ENCRYPT/DECRYPT",
						"automatic_rotation": "false",
						"tags.%":             "2",
						"tags.Created":       "TF",
						"tags.For":           "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate_id", "enable_automatic_rotation", "pending_window_in_days", "policy_name", "secret_name"},
			},
		},
	})
}

// Case 全生命周期 1980  twin
func TestAccAliCloudKmsKey_basic1980_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsKeyMap1980)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmskey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsKeyBasicDependence1980)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"origin":             "Aliyun_KMS",
					"protection_level":   "SOFTWARE",
					"description":        "test",
					"key_spec":           "Aliyun_AES_256",
					"key_usage":          "ENCRYPT/DECRYPT",
					"automatic_rotation": "false",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"origin":             "Aliyun_KMS",
						"protection_level":   "SOFTWARE",
						"description":        "test",
						"key_spec":           "Aliyun_AES_256",
						"key_usage":          "ENCRYPT/DECRYPT",
						"automatic_rotation": "false",
						"tags.%":             "2",
						"tags.Created":       "TF",
						"tags.For":           "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate_id", "enable_automatic_rotation", "pending_window_in_days", "policy_name", "secret_name"},
			},
		},
	})
}

// Case RDK接入 1323  twin
func TestAccAliCloudKmsKey_basic1323_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsKeyMap1323)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmskey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsKeyBasicDependence1323)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"origin":             "Aliyun_KMS",
					"protection_level":   "SOFTWARE",
					"description":        "rdk_test",
					"key_spec":           "Aliyun_AES_256",
					"key_usage":          "ENCRYPT/DECRYPT",
					"automatic_rotation": "false",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"origin":             "Aliyun_KMS",
						"protection_level":   "SOFTWARE",
						"description":        "rdk_test",
						"key_spec":           "Aliyun_AES_256",
						"key_usage":          "ENCRYPT/DECRYPT",
						"automatic_rotation": "false",
						"tags.%":             "2",
						"tags.Created":       "TF",
						"tags.For":           "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate_id", "enable_automatic_rotation", "pending_window_in_days", "policy_name", "secret_name"},
			},
		},
	})
}

// Test Kms Key. <<< Resource test cases, automatically generated.
