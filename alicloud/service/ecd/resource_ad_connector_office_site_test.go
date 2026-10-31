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
		"alicloud_ecd_ad_connector_office_site",
		&resource.Sweeper{
			Name: "alicloud_ecd_ad_connector_office_site",
			F:    testSweepEcdAdConnectorOfficeSite,
		})
}

func testSweepEcdAdConnectorOfficeSite(region string) error {
	if sweep.SweepPreCheckWithRegions(region, true, connectivity.EcdSupportRegions) {
		log.Printf("[INFO] Skipping Ecd Ad Connector Office Site unsupported region: %s", region)
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
	action := "DescribeOfficeSites"
	request := map[string]interface{}{}
	request["RegionId"] = aliyunClient.RegionId
	request["MaxResults"] = helper.PageSizeLarge

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
			_, err = aliyunClient.RpcPost("ecd", "2020-09-30", action, nil, request, false)
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

func TestAccAliCloudECDAdConnectorOfficeSite_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_ad_connector_office_site.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.EcdSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AliCloudECDAdConnectorOfficeSiteMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecd.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcdAdConnectorOfficeSite")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdadconnectorofficesite%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudECDAdConnectorOfficeSiteBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ad_hostname":                   "CreateEcdAdConnectorOfficeSiteValue",
					"cen_owner_id":                  "${data.alicloud_account.default.id}",
					"protocol_type":                 "ASP",
					"specification":                 "1",
					"verify_code":                   "CreateEcdAdConnectorOfficeSiteValue",
					"sub_domain_dns_address":        []string{"127.0.0.3"},
					"cidr_block":                    "10.0.0.0/12",
					"dns_address":                   []string{"127.0.0.2"},
					"sub_domain_name":               "child.example1234.com",
					"ad_connector_office_site_name": name,
					"bandwidth":                     "100",
					"enable_internet_access":        "true",
					"domain_name":                   "example1234.com",
					"enable_admin_access":           "true",
					"mfa_enabled":                   "true",
					"domain_password":               "YourPassword1234",
					"cen_id":                        "${data.alicloud_cen_instances.default.instances.0.id}",
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
						"cen_id":                        acctest.CHECKSET,
						"desktop_access_type":           "INTERNET",
						"domain_user_name":              "Administrator",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ad_connector_office_site_name": name + "Update",
					"sub_domain_dns_address":        []string{"127.0.0.5"},
					"cidr_block":                    "10.0.0.0/12",
					"dns_address":                   []string{"127.0.0.4"},
					"sub_domain_name":               "child.example5678.com",
					"bandwidth":                     "100",
					"enable_internet_access":        "true",
					"domain_name":                   "example5678.com",
					"enable_admin_access":           "true",
					"mfa_enabled":                   "false",
					"domain_password":               "YourPassword1234",
					"cen_id":                        "${data.alicloud_cen_instances.default.instances.0.id}",
					"desktop_access_type":           "VPC",
					"domain_user_name":              "Administrator2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sub_domain_dns_address.#":      acctest.CHECKSET,
						"cidr_block":                    "10.0.0.0/12",
						"dns_address.#":                 acctest.CHECKSET,
						"sub_domain_name":               acctest.CHECKSET,
						"ad_connector_office_site_name": acctest.CHECKSET,
						"bandwidth":                     "100",
						"enable_internet_access":        "true",
						"domain_name":                   acctest.CHECKSET,
						"enable_admin_access":           "true",
						"mfa_enabled":                   acctest.CHECKSET,
						"domain_password":               "YourPassword1234",
						"cen_id":                        acctest.CHECKSET,
						"desktop_access_type":           acctest.CHECKSET,
						"domain_user_name":              acctest.CHECKSET,
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

func TestAccAliCloudECDAdConnectorOfficeSite_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_ad_connector_office_site.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.EcdSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AliCloudECDAdConnectorOfficeSiteMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecd.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcdAdConnectorOfficeSite")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdadconnectorofficesite%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudECDAdConnectorOfficeSiteBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ad_connector_office_site_name": name,
					"cidr_block":                    "10.0.0.0/12",
					"domain_name":                   "example1234.com",
					"cen_id":                        "${data.alicloud_cen_instances.default.instances.0.id}",
					"dns_address":                   []string{"127.0.0.2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ad_connector_office_site_name": name,
						"cidr_block":                    "10.0.0.0/12",
						"domain_name":                   "example1234.com",
						"cen_id":                        acctest.CHECKSET,
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

var AliCloudECDAdConnectorOfficeSiteMap0 = map[string]string{
	"dns_address.#":       acctest.CHECKSET,
	"desktop_access_type": acctest.CHECKSET,
	"protocol_type":       acctest.NOSET,
	"verify_code":         acctest.NOSET,
	"status":              acctest.CHECKSET,
	"ad_hostname":         acctest.NOSET,
	"cen_owner_id":        acctest.NOSET,
}

var AliCloudECDAdConnectorOfficeSiteMap1 = map[string]string{
	"cen_owner_id":        acctest.NOSET,
	"dns_address.#":       acctest.CHECKSET,
	"desktop_access_type": acctest.CHECKSET,
	"protocol_type":       acctest.NOSET,
	"status":              acctest.CHECKSET,
	"verify_code":         acctest.NOSET,
	"ad_hostname":         acctest.NOSET,
}

func AliCloudECDAdConnectorOfficeSiteBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_cen_instances" "default" {
  		name_regex = "no-deleting-cen"
	}

	data "alicloud_account" "default" {}
`, name)
}

// lintignore: R001
func TestUnitAliCloudECDAdConnectorOfficeSite(t *testing.T) {
	p := ecd.ResourceAdConnectorOfficeSite()
	dInit, _ := schema.InternalMap(p.Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p.Schema).Data(nil, nil)
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
	rawClient, err := ecdTestClientForRegion(region)
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
			Code:       ecdString(errorCode),
			Data:       ecdString(errorCode),
			Message:    ecdString(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_ecd_ad_connector_office_site", errorCode))
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
	err = ecd.ResourceAdConnectorOfficeSite().Create(dInit, rawClient)
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
		err := ecd.ResourceAdConnectorOfficeSite().Create(dInit, rawClient)
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

	err = ecd.ResourceAdConnectorOfficeSite().Update(dExisted, rawClient)
	assert.Nil(t, err)

	// Read
	attributesDiff := map[string]interface{}{}
	diff, err := acctest.NewInstanceDiff("alicloud_ecd_ad_connector_office_site", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p.Schema).Data(dInit.State(), diff)
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
		err := ecd.ResourceAdConnectorOfficeSite().Read(dExisted, rawClient)
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
	err = ecd.ResourceAdConnectorOfficeSite().Delete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = acctest.NewInstanceDiff("alicloud_ecd_ad_connector_office_site", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p.Schema).Data(dInit.State(), diff)
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
		err := ecd.ResourceAdConnectorOfficeSite().Delete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
