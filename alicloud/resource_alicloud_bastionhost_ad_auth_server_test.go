package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

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

func TestAccAlicloudBastionhostAdAuthServer_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_ad_auth_server.default"
	checkoutSupportedRegions(t, true, connectivity.BastionhostSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudBastionhostAdAuthServerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostAdAuthServer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhostadauthserver%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostAdAuthServerBasicDependence0)
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
					"instance_id": "${data.alicloud_bastionhost_instances.default.ids.0}",
					"server":      "192.168.1.1",
					"port":        "80",
					"is_ssl":      "false",
					"domain":      "domain",
					"account":     "cn=Manager,dc=test,dc=com",
					"password":    "YouPassword123",
					"base_dn":     "dc=test,dc=com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"server":      "192.168.1.1",
						"port":        "80",
						"is_ssl":      "false",
						"domain":      "domain",
						"account":     "cn=Manager,dc=test,dc=com",
						"base_dn":     "dc=test,dc=com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "81",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "81",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server": "192.168.1.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server": "192.168.1.2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"is_ssl": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_ssl": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"base_dn": "dc=test1,dc=com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"base_dn": "dc=test1,dc=com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account": "cn=Manager1,dc=test,dc=com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account": "cn=Manager1,dc=test,dc=com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"filter": "objectClass=person",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"filter": "objectClass=person",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"standby_server": "192.168.1.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"standby_server": "192.168.1.3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "YouPassword1234",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name_mapping": "nameAttr",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name_mapping": "nameAttr",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"email_mapping": "emailAttr",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"email_mapping": "emailAttr",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_mapping": "mobileAttr",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_mapping": "mobileAttr",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server":   "192.168.1.1",
					"port":     "80",
					"is_ssl":   "false",
					"domain":   "domain",
					"account":  "cn=Manager,dc=test,dc=com",
					"password": "YouPassword123",
					"base_dn":  "dc=test,dc=com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server":  "192.168.1.1",
						"port":    "80",
						"is_ssl":  "false",
						"domain":  "domain",
						"account": "cn=Manager,dc=test,dc=com",
						"base_dn": "dc=test,dc=com",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccAlicloudBastionhostAdAuthServer_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_ad_auth_server.default"
	checkoutSupportedRegions(t, true, connectivity.BastionhostSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudBastionhostAdAuthServerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostAdAuthServer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhostadauthserver%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostAdAuthServerBasicDependence0)
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
					"server":         "192.168.1.1",
					"standby_server": "192.168.1.3",
					"port":           "80",
					"domain":         "domain",
					"account":        "cn=Manager,dc=test,dc=com",
					"password":       "YouPassword123",
					"filter":         "objectClass=person",
					"name_mapping":   "nameAttr",
					"email_mapping":  "emailAttr",
					"mobile_mapping": "mobileAttr",
					"instance_id":    "${data.alicloud_bastionhost_instances.default.ids.0}",
					"is_ssl":         "true",
					"base_dn":        "dc=test,dc=com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server":         "192.168.1.1",
						"standby_server": "192.168.1.3",
						"port":           "80",
						"domain":         "domain",
						"account":        "cn=Manager,dc=test,dc=com",
						"filter":         "objectClass=person",
						"name_mapping":   "nameAttr",
						"email_mapping":  "emailAttr",
						"mobile_mapping": "mobileAttr",
						"instance_id":    CHECKSET,
						"is_ssl":         "true",
						"base_dn":        "dc=test,dc=com",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

var AlicloudBastionhostAdAuthServerMap0 = map[string]string{
	"instance_id": CHECKSET,
}

func AlicloudBastionhostAdAuthServerBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_bastionhost_instances" "default" {}
`, name)
}

func TestAccAlicloudBastionhostAdAuthServer_unit(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_bastionhost_ad_auth_server"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_bastionhost_ad_auth_server"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"server":         "CreateBastionhostAdAuthServerValue",
		"standby_server": "CreateBastionhostAdAuthServerValue",
		"port":           80,
		"domain":         "CreateBastionhostAdAuthServerValue",
		"account":        "CreateBastionhostAdAuthServerValue",
		"password":       "CreateBastionhostAdAuthServerValue",
		"filter":         "CreateBastionhostAdAuthServerValue",
		"name_mapping":   "CreateBastionhostAdAuthServerValue",
		"email_mapping":  "CreateBastionhostAdAuthServerValue",
		"mobile_mapping": "CreateBastionhostAdAuthServerValue",
		"instance_id":    "CreateBastionhostAdAuthServerValue",
		"is_ssl":         true,
		"base_dn":        "CreateBastionhostAdAuthServerValue",
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
		"AD": map[string]interface{}{
			"Account":       "CreateBastionhostAdAuthServerValue",
			"NameMapping":   "CreateBastionhostAdAuthServerValue",
			"Server":        "CreateBastionhostAdAuthServerValue",
			"Filter":        "CreateBastionhostAdAuthServerValue",
			"Port":          80,
			"BaseDN":        "CreateBastionhostAdAuthServerValue",
			"StandbyServer": "CreateBastionhostAdAuthServerValue",
			"EmailMapping":  "CreateBastionhostAdAuthServerValue",
			"IsSSL":         true,
			"MobileMapping": "CreateBastionhostAdAuthServerValue",
			"Domain":        "CreateBastionhostAdAuthServerValue",
			"HasPassword":   true,
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_bastionhost_ad_auth_server", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewBastionhostClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudBastionhostAdAuthServerCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyInstanceADAuthServer" {
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
		err := resourceAlicloudBastionhostAdAuthServerCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_bastionhost_ad_auth_server"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewBastionhostClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudBastionhostAdAuthServerUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"server":         "UpdateBastionhostAdAuthServerValue",
		"standby_server": "UpdateBastionhostAdAuthServerValue",
		"port":           81,
		"domain":         "UpdateBastionhostAdAuthServerValue",
		"account":        "UpdateBastionhostAdAuthServerValue",
		"password":       "UpdateBastionhostAdAuthServerValue",
		"filter":         "UpdateBastionhostAdAuthServerValue",
		"name_mapping":   "UpdateBastionhostAdAuthServerValue",
		"email_mapping":  "UpdateBastionhostAdAuthServerValue",
		"mobile_mapping": "UpdateBastionhostAdAuthServerValue",
		"is_ssl":         false,
		"base_dn":        "UpdateBastionhostAdAuthServerValue",
	}
	diff, err := newInstanceDiff("alicloud_bastionhost_ad_auth_server", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_bastionhost_ad_auth_server"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"AD": map[string]interface{}{
			"Account":       "UpdateBastionhostAdAuthServerValue",
			"NameMapping":   "UpdateBastionhostAdAuthServerValue",
			"Server":        "UpdateBastionhostAdAuthServerValue",
			"Filter":        "UpdateBastionhostAdAuthServerValue",
			"Port":          81,
			"BaseDN":        "UpdateBastionhostAdAuthServerValue",
			"StandbyServer": "UpdateBastionhostAdAuthServerValue",
			"EmailMapping":  "UpdateBastionhostAdAuthServerValue",
			"IsSSL":         false,
			"MobileMapping": "UpdateBastionhostAdAuthServerValue",
			"Domain":        "UpdateBastionhostAdAuthServerValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyInstanceADAuthServer" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudBastionhostAdAuthServerUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_bastionhost_ad_auth_server"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	diff, err = newInstanceDiff("alicloud_bastionhost_ad_auth_server", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_bastionhost_ad_auth_server"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetInstanceADAuthServer" {
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
		err := resourceAlicloudBastionhostAdAuthServerRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	err = resourceAlicloudBastionhostAdAuthServerDelete(dExisted, rawClient)
	assert.Nil(t, err)
}
