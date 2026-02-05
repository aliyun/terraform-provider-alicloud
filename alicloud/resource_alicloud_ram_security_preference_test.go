package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudRamSecurityPreference_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_security_preference.default"
	ra := resourceAttrInit(resourceId, AlicloudRAMSecurityPreferenceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamSecurityPreference")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sramsecuritypreference%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRAMSecurityPreferenceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_save_mfa_ticket":           "true",
					"allow_user_to_change_password":    "true",
					"allow_user_to_manage_access_keys": "true",
					"allow_user_to_manage_mfa_devices": "true",
					"login_session_duration":           "7",
					"login_network_masks":              "42.120.66.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_save_mfa_ticket":           "true",
						"allow_user_to_change_password":    "true",
						"allow_user_to_manage_access_keys": "true",
						"allow_user_to_manage_mfa_devices": "true",
						"login_session_duration":           "7",
						"login_network_masks":              "42.120.66.0/24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mfa_operation_for_login": "independent",
					"login_network_masks":     REMOVEKEY,
					"login_session_duration":  REMOVEKEY,
					"enable_save_mfa_ticket":  REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mfa_operation_for_login": "independent",
						"login_network_masks":     REMOVEKEY,
						"login_session_duration":  REMOVEKEY,
						"enable_save_mfa_ticket":  REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_save_mfa_ticket": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_save_mfa_ticket": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allow_user_to_change_password": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_user_to_change_password": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allow_user_to_manage_access_keys": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_user_to_manage_access_keys": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allow_user_to_manage_mfa_devices": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_user_to_manage_mfa_devices": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"login_session_duration": "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_session_duration": "8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"login_network_masks":    "10.0.0.0/8",
					"login_session_duration": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_network_masks":    "10.0.0.0/8",
						"login_session_duration": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"login_network_masks": "192.168.0.0/16",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_network_masks": "192.168.0.0/16",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"login_network_masks":           "",
					"max_idle_days_for_access_keys": "90",
					"max_idle_days_for_users":       "90",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_network_masks":           "",
						"max_idle_days_for_access_keys": "90",
						"max_idle_days_for_users":       "90",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_save_mfa_ticket":           "true",
					"allow_user_to_change_password":    "true",
					"allow_user_to_manage_access_keys": "true",
					"allow_user_to_manage_mfa_devices": "true",
					"login_session_duration":           "7",
					"max_idle_days_for_access_keys":    "180",
					"max_idle_days_for_users":          "180",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_save_mfa_ticket":           "true",
						"allow_user_to_change_password":    "true",
						"allow_user_to_manage_access_keys": "true",
						"allow_user_to_manage_mfa_devices": "true",
						"login_session_duration":           "7",
						"max_idle_days_for_access_keys":    "180",
						"max_idle_days_for_users":          "180",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudRAMSecurityPreferenceMap0 = map[string]string{}

func AlicloudRAMSecurityPreferenceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
default = "%s"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudRAMSecurityPreference(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_ram_security_preference"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_ram_security_preference"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"enable_save_mfa_ticket":           false,
		"allow_user_to_change_password":    false,
		"allow_user_to_manage_access_keys": false,
		"allow_user_to_manage_mfa_devices": false,
		"login_session_duration":           7,
		"login_network_masks":              "42.120.66.0/24",
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
		//GetSecurityPreference
		"SecurityPreference": map[string]interface{}{
			"AccessKeyPreference": map[string]interface{}{"AllowUserToManageAccessKeys": false},
			"LoginProfilePreference": map[string]interface{}{
				"AllowUserToChangePassword": false,
				"EnableSaveMFATicket":       false,
				"LoginNetworkMasks":         "42.120.66.0/24",
				"LoginSessionDuration":      7,
				"EnforceMFAForLogin":        false,
			},
			"MFAPreference": map[string]interface{}{"AllowUserToManageMFADevices": false},
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
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ram_security_preference", "MockId"))
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
			// at present, The result that api returned does not contain id
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewImsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudRamSecurityPreferenceCreate(d, rawClient)
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
		err := resourceAliCloudRamSecurityPreferenceCreate(d, rawClient)
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
		err := resourceAliCloudRamSecurityPreferenceCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("RamSecurityPreference")
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewImsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAliCloudRamSecurityPreferenceUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateSetSecurityPreferenceAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"enable_save_mfa_ticket", "allow_user_to_change_password", "allow_user_to_manage_access_keys", "allow_user_to_manage_mfa_devices", "login_session_duration", "login_network_masks"} {
			switch p["alicloud_ram_security_preference"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ram_security_preference"].Schema).Data(nil, diff)
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
			return responseMock["Normal"]("")
		})
		err := resourceAliCloudRamSecurityPreferenceUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateSetSecurityPreferenceNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"enable_save_mfa_ticket", "allow_user_to_change_password", "allow_user_to_manage_access_keys", "allow_user_to_manage_mfa_devices", "login_session_duration", "login_network_masks"} {
			switch p["alicloud_ram_security_preference"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ram_security_preference"].Schema).Data(nil, diff)
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
		err := resourceAliCloudRamSecurityPreferenceUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		err := resourceAliCloudRamSecurityPreferenceDelete(d, rawClient)
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeRamSecurityPreferenceNotFound", func(t *testing.T) {
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
		err := resourceAliCloudRamSecurityPreferenceRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})
	t.Run("ReadDescribeRamSecurityPreferenceAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := false
			noRetryFlag := true
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudRamSecurityPreferenceRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})

}

// Test Ram SecurityPreference. >>> Resource test cases, automatically generated.
// Case SecurityPreference测试 9192
// lintignore: AT001
func TestAccAliCloudRamSecurityPreference_basic9192(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_security_preference.default"
	ra := resourceAttrInit(resourceId, AlicloudRamSecurityPreferenceMap9192)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamSecurityPreference")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRamSecurityPreferenceBasicDependence9192)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"login_session_duration":                  "10",
					"login_network_masks":                     "192.168.0.0/16",
					"allow_user_to_change_password":           "true",
					"allow_user_to_manage_access_keys":        "true",
					"operation_for_risk_login":                "enforceVerify",
					"enable_save_mfa_ticket":                  "true",
					"mfa_operation_for_login":                 "mandatory",
					"allow_user_to_manage_personal_ding_talk": "true",
					"allow_user_to_manage_mfa_devices":        "true",
					"verification_types":                      []string{"sms"},
					"allow_user_to_login_with_passkey":        "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_session_duration":                  "10",
						"login_network_masks":                     "192.168.0.0/16",
						"allow_user_to_change_password":           "true",
						"allow_user_to_manage_access_keys":        "true",
						"operation_for_risk_login":                "enforceVerify",
						"enable_save_mfa_ticket":                  "true",
						"mfa_operation_for_login":                 "mandatory",
						"allow_user_to_manage_personal_ding_talk": "true",
						"allow_user_to_manage_mfa_devices":        "true",
						"verification_types.#":                    "1",
						"allow_user_to_login_with_passkey":        "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"login_session_duration":                  "8",
					"login_network_masks":                     "192.168.0.0/15",
					"allow_user_to_change_password":           "false",
					"allow_user_to_manage_access_keys":        "false",
					"operation_for_risk_login":                "autonomous",
					"enable_save_mfa_ticket":                  "false",
					"mfa_operation_for_login":                 "independent",
					"allow_user_to_manage_personal_ding_talk": "false",
					"allow_user_to_manage_mfa_devices":        "false",
					"verification_types": []string{
						"sms", "email", "mfa"},
					"allow_user_to_login_with_passkey": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_session_duration":                  "8",
						"login_network_masks":                     "192.168.0.0/15",
						"allow_user_to_change_password":           "false",
						"allow_user_to_manage_access_keys":        "false",
						"operation_for_risk_login":                "autonomous",
						"enable_save_mfa_ticket":                  "false",
						"mfa_operation_for_login":                 "independent",
						"allow_user_to_manage_personal_ding_talk": "false",
						"allow_user_to_manage_mfa_devices":        "false",
						"verification_types.#":                    "3",
						"allow_user_to_login_with_passkey":        "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"login_session_duration": "6",
					"verification_types":     []string{},
					"login_network_masks":    REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_session_duration": "6",
						"verification_types.#":   "0",
						"login_network_masks":    REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"login_session_duration":                  "10",
					"login_network_masks":                     "192.168.0.0/16",
					"allow_user_to_change_password":           "true",
					"allow_user_to_manage_access_keys":        "true",
					"operation_for_risk_login":                "enforceVerify",
					"enable_save_mfa_ticket":                  "true",
					"mfa_operation_for_login":                 "mandatory",
					"allow_user_to_manage_personal_ding_talk": "true",
					"allow_user_to_manage_mfa_devices":        "true",
					"verification_types":                      []string{"sms"},
					"allow_user_to_login_with_passkey":        "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_session_duration":                  "10",
						"login_network_masks":                     "192.168.0.0/16",
						"allow_user_to_change_password":           "true",
						"allow_user_to_manage_access_keys":        "true",
						"operation_for_risk_login":                "enforceVerify",
						"enable_save_mfa_ticket":                  "true",
						"mfa_operation_for_login":                 "mandatory",
						"allow_user_to_manage_personal_ding_talk": "true",
						"allow_user_to_manage_mfa_devices":        "true",
						"verification_types.#":                    "1",
						"allow_user_to_login_with_passkey":        "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"login_session_duration":                  "8",
					"login_network_masks":                     "192.168.0.0/15",
					"allow_user_to_change_password":           "false",
					"allow_user_to_manage_access_keys":        "false",
					"operation_for_risk_login":                "autonomous",
					"enable_save_mfa_ticket":                  "false",
					"mfa_operation_for_login":                 "independent",
					"allow_user_to_manage_personal_ding_talk": "false",
					"allow_user_to_manage_mfa_devices":        "false",
					"allow_user_to_login_with_passkey":        "false",
					"verification_types":                      []interface{}{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_session_duration":                  "8",
						"login_network_masks":                     "192.168.0.0/15",
						"allow_user_to_change_password":           "false",
						"allow_user_to_manage_access_keys":        "false",
						"operation_for_risk_login":                "autonomous",
						"enable_save_mfa_ticket":                  "false",
						"mfa_operation_for_login":                 "independent",
						"allow_user_to_manage_personal_ding_talk": "false",
						"allow_user_to_manage_mfa_devices":        "false",
						"allow_user_to_login_with_passkey":        "false",
						"verification_types.#":                    "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"login_session_duration":                  "10",
					"login_network_masks":                     "192.168.0.0/16",
					"allow_user_to_change_password":           "true",
					"allow_user_to_manage_access_keys":        "true",
					"operation_for_risk_login":                "enforceVerify",
					"enable_save_mfa_ticket":                  "true",
					"mfa_operation_for_login":                 "mandatory",
					"allow_user_to_manage_personal_ding_talk": "true",
					"allow_user_to_manage_mfa_devices":        "true",
					"verification_types": []string{
						"sms"},
					"allow_user_to_login_with_passkey": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_session_duration":                  "10",
						"login_network_masks":                     "192.168.0.0/16",
						"allow_user_to_change_password":           "true",
						"allow_user_to_manage_access_keys":        "true",
						"operation_for_risk_login":                "enforceVerify",
						"enable_save_mfa_ticket":                  "true",
						"mfa_operation_for_login":                 "mandatory",
						"allow_user_to_manage_personal_ding_talk": "true",
						"allow_user_to_manage_mfa_devices":        "true",
						"verification_types.#":                    "1",
						"allow_user_to_login_with_passkey":        "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudRamSecurityPreferenceMap9192 = map[string]string{}

func AlicloudRamSecurityPreferenceBasicDependence9192(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case SecurityPreference测试_副本1745566136705 10751
// lintignore: AT001
func TestAccAliCloudRamSecurityPreference_basic10751(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_security_preference.default"
	ra := resourceAttrInit(resourceId, AlicloudRamSecurityPreferenceMap10751)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamSecurityPreference")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRamSecurityPreferenceBasicDependence10751)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"login_session_duration":                  "10",
					"login_network_masks":                     "192.168.0.0/16",
					"allow_user_to_change_password":           "true",
					"allow_user_to_manage_access_keys":        "true",
					"operation_for_risk_login":                "enforceVerify",
					"enable_save_mfa_ticket":                  "true",
					"mfa_operation_for_login":                 "mandatory",
					"allow_user_to_manage_personal_ding_talk": "true",
					"allow_user_to_manage_mfa_devices":        "true",
					"verification_types": []string{
						"sms"},
					"allow_user_to_login_with_passkey": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_session_duration":                  "10",
						"login_network_masks":                     "192.168.0.0/16",
						"allow_user_to_change_password":           "true",
						"allow_user_to_manage_access_keys":        "true",
						"operation_for_risk_login":                "enforceVerify",
						"enable_save_mfa_ticket":                  "true",
						"mfa_operation_for_login":                 "mandatory",
						"allow_user_to_manage_personal_ding_talk": "true",
						"allow_user_to_manage_mfa_devices":        "true",
						"verification_types.#":                    "1",
						"allow_user_to_login_with_passkey":        "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"login_session_duration":                  "8",
					"login_network_masks":                     "192.168.0.0/15",
					"allow_user_to_change_password":           "false",
					"allow_user_to_manage_access_keys":        "false",
					"operation_for_risk_login":                "autonomous",
					"enable_save_mfa_ticket":                  "false",
					"mfa_operation_for_login":                 "independent",
					"allow_user_to_manage_personal_ding_talk": "false",
					"allow_user_to_manage_mfa_devices":        "false",
					"verification_types": []string{
						"sms", "email", "mfa"},
					"allow_user_to_login_with_passkey": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_session_duration":                  "8",
						"login_network_masks":                     "192.168.0.0/15",
						"allow_user_to_change_password":           "false",
						"allow_user_to_manage_access_keys":        "false",
						"operation_for_risk_login":                "autonomous",
						"enable_save_mfa_ticket":                  "false",
						"mfa_operation_for_login":                 "independent",
						"allow_user_to_manage_personal_ding_talk": "false",
						"allow_user_to_manage_mfa_devices":        "false",
						"verification_types.#":                    "3",
						"allow_user_to_login_with_passkey":        "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"login_session_duration": "6",
					"verification_types":     []string{},
					"login_network_masks":    REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_session_duration": "6",
						"verification_types.#":   "0",
						"login_network_masks":    REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"login_session_duration":                  "10",
					"login_network_masks":                     "192.168.0.0/16",
					"allow_user_to_change_password":           "true",
					"allow_user_to_manage_access_keys":        "true",
					"operation_for_risk_login":                "enforceVerify",
					"enable_save_mfa_ticket":                  "true",
					"mfa_operation_for_login":                 "mandatory",
					"allow_user_to_manage_personal_ding_talk": "true",
					"allow_user_to_manage_mfa_devices":        "true",
					"verification_types": []string{
						"sms"},
					"allow_user_to_login_with_passkey": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_session_duration":                  "10",
						"login_network_masks":                     "192.168.0.0/16",
						"allow_user_to_change_password":           "true",
						"allow_user_to_manage_access_keys":        "true",
						"operation_for_risk_login":                "enforceVerify",
						"enable_save_mfa_ticket":                  "true",
						"mfa_operation_for_login":                 "mandatory",
						"allow_user_to_manage_personal_ding_talk": "true",
						"allow_user_to_manage_mfa_devices":        "true",
						"verification_types.#":                    "1",
						"allow_user_to_login_with_passkey":        "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"login_session_duration":                  "8",
					"login_network_masks":                     "192.168.0.0/15",
					"allow_user_to_change_password":           "false",
					"allow_user_to_manage_access_keys":        "false",
					"operation_for_risk_login":                "autonomous",
					"enable_save_mfa_ticket":                  "false",
					"mfa_operation_for_login":                 "independent",
					"allow_user_to_manage_personal_ding_talk": "false",
					"allow_user_to_manage_mfa_devices":        "false",
					"allow_user_to_login_with_passkey":        "false",
					"verification_types":                      []interface{}{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_session_duration":                  "8",
						"login_network_masks":                     "192.168.0.0/15",
						"allow_user_to_change_password":           "false",
						"allow_user_to_manage_access_keys":        "false",
						"operation_for_risk_login":                "autonomous",
						"enable_save_mfa_ticket":                  "false",
						"mfa_operation_for_login":                 "independent",
						"allow_user_to_manage_personal_ding_talk": "false",
						"allow_user_to_manage_mfa_devices":        "false",
						"allow_user_to_login_with_passkey":        "false",
						"verification_types.#":                    "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"login_session_duration":                  "10",
					"login_network_masks":                     "192.168.0.0/16",
					"allow_user_to_change_password":           "true",
					"allow_user_to_manage_access_keys":        "true",
					"operation_for_risk_login":                "enforceVerify",
					"enable_save_mfa_ticket":                  "true",
					"mfa_operation_for_login":                 "mandatory",
					"allow_user_to_manage_personal_ding_talk": "true",
					"allow_user_to_manage_mfa_devices":        "true",
					"verification_types": []string{
						"sms"},
					"allow_user_to_login_with_passkey": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_session_duration":                  "10",
						"login_network_masks":                     "192.168.0.0/16",
						"allow_user_to_change_password":           "true",
						"allow_user_to_manage_access_keys":        "true",
						"operation_for_risk_login":                "enforceVerify",
						"enable_save_mfa_ticket":                  "true",
						"mfa_operation_for_login":                 "mandatory",
						"allow_user_to_manage_personal_ding_talk": "true",
						"allow_user_to_manage_mfa_devices":        "true",
						"verification_types.#":                    "1",
						"allow_user_to_login_with_passkey":        "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudRamSecurityPreferenceMap10751 = map[string]string{}

func AlicloudRamSecurityPreferenceBasicDependence10751(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Ram SecurityPreference. <<< Resource test cases, automatically generated.
