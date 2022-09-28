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
		"alicloud_ecd_ad_connector_office_site",
		&resource.Sweeper{
			Name: "alicloud_ecd_ad_connector_office_site",
			F:    testSweepEcdAdConnectorOfficeSite,
		})
}

func testSweepEcdAdConnectorOfficeSite(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.EcdSupportRegions) {
		log.Printf("[INFO] Skipping Ecd Ad Connector Office Site unsupported region: %s", region)
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
	action := "DescribeOfficeSites"
	request := map[string]interface{}{}
	request["RegionId"] = aliyunClient.RegionId
	request["MaxResults"] = PageSizeLarge

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

		resp, err := jsonpath.Get("$.OfficeSites", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.OfficeSites", action, err)
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
				log.Printf("[INFO] Skipping Ecd Ad Connector Office Site: %s", item["Name"].(string))
				continue
			}
			action := "DeleteOfficeSites"
			request := map[string]interface{}{
				"OfficeSiteId": []string{"OfficeSiteId"},
				"RegionId":     aliyunClient.RegionId,
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ecd Ad Connector Office Site (%s): %s", item["Name"].(string), err)
			}
			log.Printf("[INFO] Delete Ecd Ad Connector Office Site success: %s ", item["Name"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudECDAdConnectorOfficeSite_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_ad_connector_office_site.default"
	checkoutSupportedRegions(t, true, connectivity.EcdSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudECDAdConnectorOfficeSiteMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcdService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcdAdConnectorOfficeSite")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdadconnectorofficesite%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECDAdConnectorOfficeSiteBasicDependence0)
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
					"sub_domain_dns_address":        []string{"127.0.0.3"},
					"cidr_block":                    "10.0.0.0/12",
					"dns_address":                   []string{"127.0.0.2"},
					"sub_domain_name":               "child.example1234.com",
					"ad_connector_office_site_name": "${var.name}",
					"bandwidth":                     "100",
					"enable_internet_access":        "true",
					"domain_name":                   "example1234.com",
					"enable_admin_access":           "true",
					"mfa_enabled":                   "true",
					"domain_password":               "YourPassword1234",
					"cen_id":                        "${alicloud_cen_instance.default.id}",
					"desktop_access_type":           "INTERNET",
					"domain_user_name":              "Administrator",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sub_domain_dns_address.#":      "1",
						"cidr_block":                    "10.0.0.0/12",
						"dns_address.#":                 "1",
						"sub_domain_name":               "child.example1234.com",
						"ad_connector_office_site_name": name,
						"bandwidth":                     "100",
						"enable_internet_access":        "true",
						"domain_name":                   "example1234.com",
						"enable_admin_access":           "true",
						"mfa_enabled":                   "true",
						"domain_password":               "YourPassword1234",
						"cen_id":                        CHECKSET,
						"desktop_access_type":           "INTERNET",
						"domain_user_name":              "Administrator",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domain_password", "protocol_type", "verify_code", "specification", "ad_hostname", "cen_owner_id"},
			},
		},
	})
}

var AlicloudECDAdConnectorOfficeSiteMap0 = map[string]string{
	"dns_address.#":       CHECKSET,
	"desktop_access_type": CHECKSET,
	"protocol_type":       NOSET,
	"verify_code":         NOSET,
	"status":              CHECKSET,
	"ad_hostname":         NOSET,
	"cen_owner_id":        NOSET,
}

func AlicloudECDAdConnectorOfficeSiteBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
resource "alicloud_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	protection_level = "REDUCED"
}
`, name)
}

func TestAccAlicloudECDAdConnectorOfficeSite_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_ad_connector_office_site.default"
	checkoutSupportedRegions(t, true, connectivity.EcdSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudECDAdConnectorOfficeSiteMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcdService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcdAdConnectorOfficeSite")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdadconnectorofficesite%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECDAdConnectorOfficeSiteBasicDependence0)
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
					"ad_connector_office_site_name": "${var.name}",
					"cidr_block":                    "10.0.0.0/12",
					"domain_name":                   "example1234.com",
					"cen_id":                        "${alicloud_cen_instance.default.id}",
					"dns_address":                   []string{"127.0.0.2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ad_connector_office_site_name": name,
						"cidr_block":                    "10.0.0.0/12",
						"domain_name":                   "example1234.com",
						"cen_id":                        CHECKSET,
						"dns_address.#":                 "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"protocol_type", "verify_code", "specification", "ad_hostname", "cen_owner_id"},
			},
		},
	})
}

var AlicloudECDAdConnectorOfficeSiteMap1 = map[string]string{
	"cen_owner_id":        NOSET,
	"dns_address.#":       CHECKSET,
	"desktop_access_type": CHECKSET,
	"protocol_type":       NOSET,
	"status":              CHECKSET,
	"verify_code":         NOSET,
	"ad_hostname":         NOSET,
}

func TestUnitAlicloudECDAdConnectorOfficeSite(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ecd_ad_connector_office_site"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ecd_ad_connector_office_site"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"sub_domain_dns_address":        []string{"CreateEcdAdConnectorOfficeSiteValue"},
		"cidr_block":                    "CreateEcdAdConnectorOfficeSiteValue",
		"dns_address":                   []string{"CreateEcdAdConnectorOfficeSiteValue"},
		"sub_domain_name":               "CreateEcdAdConnectorOfficeSiteValue",
		"ad_connector_office_site_name": "CreateEcdAdConnectorOfficeSiteValue",
		"bandwidth":                     100,
		"enable_internet_access":        true,
		"domain_name":                   "CreateEcdAdConnectorOfficeSiteValue",
		"enable_admin_access":           true,
		"mfa_enabled":                   false,
		"domain_password":               "CreateEcdAdConnectorOfficeSiteValue",
		"cen_id":                        "CreateEcdAdConnectorOfficeSiteValue",
		"desktop_access_type":           "CreateEcdAdConnectorOfficeSiteValue",
		"domain_user_name":              "CreateEcdAdConnectorOfficeSiteValue",
		"ad_hostname":                   "CreateEcdAdConnectorOfficeSiteValue",
		"cen_owner_id":                  "CreateEcdAdConnectorOfficeSiteValue",
		"protocol_type":                 "CreateEcdAdConnectorOfficeSiteValue",
		"specification":                 1,
		"verify_code":                   "CreateEcdAdConnectorOfficeSiteValue",
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
		"OfficeSites": []interface{}{
			map[string]interface{}{
				"Status":                   "REGISTERED",
				"NetworkPackageId":         "CreateEcdAdConnectorOfficeSiteValue",
				"CloudBoxOfficeSite":       false,
				"CidrBlock":                "CreateEcdAdConnectorOfficeSiteValue",
				"AdHostname":               "",
				"DnsAddress":               []string{"CreateEcdAdConnectorOfficeSiteValue"},
				"SubDomainName":            "CreateEcdAdConnectorOfficeSiteValue",
				"Name":                     "CreateEcdAdConnectorOfficeSiteValue",
				"DesktopCount":             0,
				"SubDnsAddress":            []string{"CreateEcdAdConnectorOfficeSiteValue"},
				"Bandwidth":                100,
				"EnableInternetAccess":     true,
				"DomainName":               "CreateEcdAdConnectorOfficeSiteValue",
				"EnableAdminAccess":        true,
				"OfficeSiteType":           "CreateEcdAdConnectorOfficeSiteValue",
				"SsoEnabled":               true,
				"MfaEnabled":               false,
				"OfficeSiteId":             "EcdAdConnectorOfficeSiteId",
				"VpcType":                  "CreateEcdAdConnectorOfficeSiteValue",
				"VpcId":                    "CreateEcdAdConnectorOfficeSiteValue",
				"EnableCrossDesktopAccess": true,
				"ProtocolType":             "CreateEcdAdConnectorOfficeSiteValue",
				"CreationTime":             "CreateEcdAdConnectorOfficeSiteValue",
				"CenId":                    "CreateEcdAdConnectorOfficeSiteValue",
				"DesktopAccessType":        "CreateEcdAdConnectorOfficeSiteValue",
				"DomainUserName":           "CreateEcdAdConnectorOfficeSiteValue",
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		"OfficeSiteId": "EcdAdConnectorOfficeSiteId",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ecd_ad_connector_office_site", errorCode))
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
	err = resourceAlicloudEcdAdConnectorOfficeSiteCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateADConnectorOfficeSite" {
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
		err := resourceAlicloudEcdAdConnectorOfficeSiteCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecd_ad_connector_office_site"].Schema).Data(dInit.State(), nil)
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

	err = resourceAlicloudEcdAdConnectorOfficeSiteUpdate(dExisted, rawClient)
	assert.Nil(t, err)

	// Read
	attributesDiff := map[string]interface{}{}
	diff, err := newInstanceDiff("alicloud_ecd_ad_connector_office_site", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecd_ad_connector_office_site"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeOfficeSites" {
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
		err := resourceAlicloudEcdAdConnectorOfficeSiteRead(dExisted, rawClient)
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
	err = resourceAlicloudEcdAdConnectorOfficeSiteDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_ecd_ad_connector_office_site", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecd_ad_connector_office_site"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteOfficeSites" {
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
			if *action == "DescribeOfficeSites" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEcdAdConnectorOfficeSiteDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
