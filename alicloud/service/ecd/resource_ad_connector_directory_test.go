package ecd_test

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecd"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
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
	if sweep.SweepPreCheckWithRegions(region, true, connectivity.EcdSupportRegions) {
		log.Printf("[INFO] Skipping Ecd Ad Connector Directory unsupported region: %s", region)
		return nil
	}

	rawClient, err := ecdTestClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeDirectories"
	request := map[string]interface{}{}
	request["RegionId"] = aliyunClient.RegionId
	request["MaxResults"] = helper.PageSizeLarge
	request["DirectoryType"] = "AD_CONNECTOR"
	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = aliyunClient.RpcPost("ecd", "2020-09-30", action, nil, request, true)
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
			_, err = aliyunClient.RpcPost("ecd", "2020-09-30", action, nil, request, false)
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

func TestAccAliCloudECDAdConnectorDirectory_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_ad_connector_directory.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.EcdSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AliCloudECDAdConnectorDirectoryMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecd.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcdAdConnectorDirectory")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdadconnectordirectory%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudECDAdConnectorDirectoryBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":      fmt.Sprintf("corp.%d.example.com", rand),
					"vswitch_ids":      []string{"${alicloud_vswitch.default.id}"},
					"dns_address":      []string{"127.0.0.2"},
					"domain_password":  "YourPassword1234",
					"domain_user_name": "sAMAccountName",
					"directory_name":   name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":      acctest.CHECKSET,
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

func TestAccAliCloudECDAdConnectorDirectory_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_ad_connector_directory.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.EcdSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AliCloudECDAdConnectorDirectoryMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecd.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcdAdConnectorDirectory")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdadconnectordirectory%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudECDAdConnectorDirectoryBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name":            fmt.Sprintf("corp.%d.example.com", rand),
					"vswitch_ids":            []string{"${alicloud_vswitch.default.id}"},
					"dns_address":            []string{"127.0.0.2"},
					"desktop_access_type":    "INTERNET",
					"domain_password":        "YourPassword1234",
					"domain_user_name":       "sAMAccountName",
					"directory_name":         name,
					"specification":          "1",
					"sub_domain_dns_address": []string{"127.0.0.3"},
					"sub_domain_name":        "child.example.com",
					"enable_admin_access":    "true",
					"mfa_enabled":            "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":              acctest.CHECKSET,
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

var AliCloudECDAdConnectorDirectoryMap0 = map[string]string{
	"dns_address.#": acctest.CHECKSET,
	"status":        acctest.CHECKSET,
	"vswitch_ids.#": acctest.CHECKSET,
}

func AliCloudECDAdConnectorDirectoryBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "10.4.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		cidr_block   = "10.4.0.0/16"
  		vpc_id       = alicloud_vpc.default.id
  		zone_id      = "cn-hangzhou-j"
	}
`, name)
}

// lintignore: R001
func TestUnitAliCloudECDAdConnectorDirectory(t *testing.T) {
	p := ecd.ResourceAdConnectorDirectory()
	dInit, _ := schema.InternalMap(p.Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p.Schema).Data(nil, nil)
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
	rawClient, err := ecdTestClientForRegion(region)
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
			Code:       ecdString(errorCode),
			Data:       ecdString(errorCode),
			Message:    ecdString(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_ecd_ad_connector_directory", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			ecdMapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGwsecdClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       ecdString("loadEndpoint error"),
			Data:       ecdString("loadEndpoint error"),
			Message:    ecdString("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = ecd.ResourceAdConnectorDirectory().Create(dInit, rawClient)
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
		err := ecd.ResourceAdConnectorDirectory().Create(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p.Schema).Data(dInit.State(), nil)
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
	err = ecd.ResourceAdConnectorDirectory().Update(dExisted, rawClient)
	assert.Nil(t, err)

	// Read
	attributesDiff := map[string]interface{}{}
	diff, err := acctest.NewInstanceDiff("alicloud_ecd_ad_connector_directory", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p.Schema).Data(dInit.State(), diff)
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
		err := ecd.ResourceAdConnectorDirectory().Read(dExisted, rawClient)
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
			Code:       ecdString("loadEndpoint error"),
			Data:       ecdString("loadEndpoint error"),
			Message:    ecdString("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = ecd.ResourceAdConnectorDirectory().Delete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = acctest.NewInstanceDiff("alicloud_ecd_ad_connector_directory", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p.Schema).Data(dInit.State(), diff)
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
		err := ecd.ResourceAdConnectorDirectory().Delete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
