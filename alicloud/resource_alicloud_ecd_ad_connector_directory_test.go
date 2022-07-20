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
		"alicloud_ecd_ad_connector_directory",
		&resource.Sweeper{
			Name: "alicloud_ecd_ad_connector_directory",
			F:    testSweepEcdAdConnectorDirectory,
		})
}

func testSweepEcdAdConnectorDirectory(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.EcdSupportRegions) {
		log.Printf("[INFO] Skipping Ecd Ad Connector Directory unsupported region: %s", region)
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeDirectories"
	request := map[string]interface{}{}
	request["RegionId"] = aliyunClient.RegionId
	request["MaxResults"] = PageSizeLarge
	request["DirectoryType"] = "AD_CONNECTOR"
	var response map[string]interface{}
	conn, err := aliyunClient.NewGwsecdClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &runtime)
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

		resp, err := jsonpath.Get("$.Directories", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Directories", action, err)
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
				log.Printf("[INFO] Skipping Ecd Ad Connector Directory: %s", item["Name"].(string))
				continue
			}
			action := "DeleteDirectories"
			request := map[string]interface{}{
				"DirectoryId": []string{fmt.Sprint(item["DirectoryId"])},
				"RegionId":    aliyunClient.RegionId,
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ecd Ad Connector Directory (%s): %s", item["Name"].(string), err)
			}
			log.Printf("[INFO] Delete Ecd Ad Connector Directory success: %s ", item["Name"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudECDAdConnectorDirectory_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_ad_connector_directory.default"
	checkoutSupportedRegions(t, true, connectivity.EcdSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudECDAdConnectorDirectoryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcdService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcdAdConnectorDirectory")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdadconnectordirectory%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECDAdConnectorDirectoryBasicDependence0)
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
					"domain_name":      "corp.example.com",
					"vswitch_ids":      []string{"${data.alicloud_vswitches.default.ids.0}"},
					"dns_address":      []string{"127.0.0.2"},
					"domain_password":  "YourPassword1234",
					"domain_user_name": "sAMAccountName",
					"directory_name":   "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":      "corp.example.com",
						"vswitch_ids.#":    "1",
						"dns_address.#":    "1",
						"domain_user_name": "sAMAccountName",
						"directory_name":   name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"specification", "domain_password"},
			},
		},
	})
}

func TestAccAlicloudECDAdConnectorDirectory_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_ad_connector_directory.default"
	checkoutSupportedRegions(t, true, connectivity.EcdSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudECDAdConnectorDirectoryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcdService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcdAdConnectorDirectory")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdadconnectordirectory%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECDAdConnectorDirectoryBasicDependence0)
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
					"domain_name":            "corp.example.com",
					"vswitch_ids":            []string{"${data.alicloud_vswitches.default.ids.0}"},
					"dns_address":            []string{"127.0.0.2"},
					"desktop_access_type":    "INTERNET",
					"domain_password":        "YourPassword1234",
					"domain_user_name":       "sAMAccountName",
					"directory_name":         "${var.name}",
					"specification":          "1",
					"sub_domain_dns_address": []string{"127.0.0.3"},
					"sub_domain_name":        "child.example.com",
					"enable_admin_access":    "true",
					"mfa_enabled":            "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":              "corp.example.com",
						"vswitch_ids.#":            "1",
						"dns_address.#":            "1",
						"desktop_access_type":      "INTERNET",
						"domain_user_name":         "sAMAccountName",
						"directory_name":           name,
						"sub_domain_dns_address.#": "1",
						"sub_domain_name":          "child.example.com",
						"enable_admin_access":      "true",
						"mfa_enabled":              "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"specification", "domain_password"},
			},
		},
	})
}

var AlicloudECDAdConnectorDirectoryMap0 = map[string]string{
	"dns_address.#": CHECKSET,
	"status":        CHECKSET,
	"vswitch_ids.#": CHECKSET,
}

func AlicloudECDAdConnectorDirectoryBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_ecd_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_ecd_zones.default.ids.0
}
`, name)
}

func TestUnitAlicloudECDAdConnectorDirectory(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ecd_ad_connector_directory"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ecd_ad_connector_directory"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"domain_name":            "CreateEcdAdConnectorDirectoryValue",
		"dns_address":            []string{"CreateEcdAdConnectorDirectoryValue"},
		"domain_password":        "CreateEcdAdConnectorDirectoryValue",
		"domain_user_name":       "CreateEcdAdConnectorDirectoryValue",
		"directory_name":         "CreateEcdAdConnectorDirectoryValue",
		"specification":          1,
		"sub_domain_dns_address": []string{"CreateEcdAdConnectorDirectoryValue"},
		"sub_domain_name":        "CreateEcdAdConnectorDirectoryValue",
		"mfa_enabled":            false,
		"enable_admin_access":    true,
		"desktop_access_type":    "CreateEcdAdConnectorDirectoryValue",
		"vswitch_ids":            []string{"CreateEcdAdConnectorDirectoryValue"},
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
		"Directories": []interface{}{
			map[string]interface{}{
				"Status":                   "REGISTERED",
				"VSwitchIds":               []string{"CreateEcdAdConnectorDirectoryValue"},
				"SubDnsAddress":            []string{"CreateEcdAdConnectorDirectoryValue"},
				"DnsAddress":               []string{"CreateEcdAdConnectorDirectoryValue"},
				"EnableAdminAccess":        true,
				"MfaEnabled":               false,
				"DirectoryType":            "CreateEcdAdConnectorDirectoryValue",
				"SubDomainName":            "CreateEcdAdConnectorDirectoryValue",
				"Name":                     "CreateEcdAdConnectorDirectoryValue",
				"DirectoryId":              "CreateEcdAdConnectorDirectoryValue",
				"VpcId":                    "CreateEcdAdConnectorDirectoryValue",
				"EnableCrossDesktopAccess": false,
				"DesktopAccessType":        "CreateEcdAdConnectorDirectoryValue",
				"EnableInternetAccess":     true,
				"DomainName":               "CreateEcdAdConnectorDirectoryValue",
				"DomainUserName":           "CreateEcdAdConnectorDirectoryValue",
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		"DirectoryId": "CreateEcdAdConnectorDirectoryValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ecd_ad_connector_directory", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGwsecdClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEcdAdConnectorDirectoryCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateADConnectorDirectory" {
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
		err := resourceAlicloudEcdAdConnectorDirectoryCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecd_ad_connector_directory"].Schema).Data(dInit.State(), nil)
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
	err = resourceAlicloudEcdAdConnectorDirectoryUpdate(dExisted, rawClient)
	assert.Nil(t, err)

	// Read
	attributesDiff := map[string]interface{}{}
	diff, err := newInstanceDiff("alicloud_ecd_ad_connector_directory", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecd_ad_connector_directory"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDirectories" {
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
		err := resourceAlicloudEcdAdConnectorDirectoryRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGwsecdClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEcdAdConnectorDirectoryDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_ecd_ad_connector_directory", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecd_ad_connector_directory"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteDirectories" {
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
		err := resourceAlicloudEcdAdConnectorDirectoryDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
