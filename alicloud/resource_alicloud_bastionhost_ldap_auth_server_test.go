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

func TestAccAlicloudBastionhostLdapAuthServer_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_ldap_auth_server.default"
	checkoutSupportedRegions(t, true, connectivity.BastionhostSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudBastionhostLdapAuthServerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostLdapAuthServer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhostldapauthserver%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostLdapAuthServerBasicDependence0)
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
					"password":    "YouPassword123",
					"account":     "cn=Manager,dc=test,dc=com",
					"base_dn":     "dc=test,dc=com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"server":      "192.168.1.1",
						"port":        "80",
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
					"login_name_mapping": "uid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_name_mapping": "uid",
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
					"account":  "cn=Manager,dc=test,dc=com",
					"password": "YouPassword123",
					"base_dn":  "dc=test,dc=com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server":  "192.168.1.1",
						"port":    "80",
						"is_ssl":  "false",
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

func TestAccAlicloudBastionhostLdapAuthServer_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_ldap_auth_server.default"
	checkoutSupportedRegions(t, true, connectivity.BastionhostSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudBastionhostLdapAuthServerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostLdapAuthServer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhostldapauthserver%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostLdapAuthServerBasicDependence0)
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
					"server":             "192.168.1.1",
					"standby_server":     "192.168.1.3",
					"port":               "80",
					"login_name_mapping": "uid",
					"account":            "cn=Manager,dc=test,dc=com",
					"password":           "YouPassword123",
					"filter":             "objectClass=person",
					"name_mapping":       "nameAttr",
					"email_mapping":      "emailAttr",
					"mobile_mapping":     "mobileAttr",
					"instance_id":        "${data.alicloud_bastionhost_instances.default.ids.0}",
					"is_ssl":             "true",
					"base_dn":            "dc=test,dc=com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server":             "192.168.1.1",
						"standby_server":     "192.168.1.3",
						"port":               "80",
						"login_name_mapping": "uid",
						"account":            "cn=Manager,dc=test,dc=com",
						"filter":             "objectClass=person",
						"name_mapping":       "nameAttr",
						"email_mapping":      "emailAttr",
						"mobile_mapping":     "mobileAttr",
						"instance_id":        CHECKSET,
						"is_ssl":             "true",
						"base_dn":            "dc=test,dc=com",
					}),
				),
			},
		},
	})
}

var AlicloudBastionhostLdapAuthServerMap0 = map[string]string{
	"instance_id": CHECKSET,
}

func AlicloudBastionhostLdapAuthServerBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_bastionhost_instances" "default" {}
`, name)
}

func TestAccAlicloudBastionhostLdapAuthServer_unit(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_bastionhost_ldap_auth_server"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_bastionhost_ldap_auth_server"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"server":             "CreateBastionhostLdapAuthServerValue",
		"standby_server":     "CreateBastionhostLdapAuthServerValue",
		"port":               80,
		"login_name_mapping": "CreateBastionhostLdapAuthServerValue",
		"account":            "CreateBastionhostLdapAuthServerValue",
		"password":           "CreateBastionhostLdapAuthServerValue",
		"filter":             "CreateBastionhostLdapAuthServerValue",
		"name_mapping":       "CreateBastionhostLdapAuthServerValue",
		"email_mapping":      "CreateBastionhostLdapAuthServerValue",
		"mobile_mapping":     "CreateBastionhostLdapAuthServerValue",
		"instance_id":        "CreateBastionhostLdapAuthServerValue",
		"is_ssl":             true,
		"base_dn":            "CreateBastionhostLdapAuthServerValue",
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
		"LDAP": map[string]interface{}{
			"Account":          "CreateBastionhostLdapAuthServerValue",
			"NameMapping":      "CreateBastionhostLdapAuthServerValue",
			"Server":           "CreateBastionhostLdapAuthServerValue",
			"Filter":           "CreateBastionhostLdapAuthServerValue",
			"Port":             80,
			"BaseDN":           "CreateBastionhostLdapAuthServerValue",
			"StandbyServer":    "CreateBastionhostLdapAuthServerValue",
			"EmailMapping":     "CreateBastionhostLdapAuthServerValue",
			"IsSSL":            true,
			"MobileMapping":    "CreateBastionhostLdapAuthServerValue",
			"LoginNameMapping": "CreateBastionhostLdapAuthServerValue",
			"HasPassword":      true,
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_bastionhost_ldap_auth_server", errorCode))
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
	err = resourceAlicloudBastionhostLdapAuthServerCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyInstanceLDAPAuthServer" {
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
		err := resourceAlicloudBastionhostLdapAuthServerCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_bastionhost_ldap_auth_server"].Schema).Data(dInit.State(), nil)
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
	err = resourceAlicloudBastionhostLdapAuthServerUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"server":             "UpdateBastionhostLdapAuthServerValue",
		"standby_server":     "UpdateBastionhostLdapAuthServerValue",
		"port":               81,
		"login_name_mapping": "UpdateBastionhostLdapAuthServerValue",
		"account":            "UpdateBastionhostLdapAuthServerValue",
		"password":           "UpdateBastionhostLdapAuthServerValue",
		"filter":             "UpdateBastionhostLdapAuthServerValue",
		"name_mapping":       "UpdateBastionhostLdapAuthServerValue",
		"email_mapping":      "UpdateBastionhostLdapAuthServerValue",
		"mobile_mapping":     "UpdateBastionhostLdapAuthServerValue",
		"is_ssl":             false,
		"base_dn":            "UpdateBastionhostLdapAuthServerValue",
	}
	diff, err := newInstanceDiff("alicloud_bastionhost_ldap_auth_server", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_bastionhost_ldap_auth_server"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"LDAP": map[string]interface{}{
			"Account":          "UpdateBastionhostLdapAuthServerValue",
			"NameMapping":      "UpdateBastionhostLdapAuthServerValue",
			"Server":           "UpdateBastionhostLdapAuthServerValue",
			"Filter":           "UpdateBastionhostLdapAuthServerValue",
			"Port":             81,
			"BaseDN":           "UpdateBastionhostLdapAuthServerValue",
			"StandbyServer":    "UpdateBastionhostLdapAuthServerValue",
			"EmailMapping":     "UpdateBastionhostLdapAuthServerValue",
			"IsSSL":            false,
			"MobileMapping":    "UpdateBastionhostLdapAuthServerValue",
			"LoginNameMapping": "UpdateBastionhostLdapAuthServerValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyInstanceLDAPAuthServer" {
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
		err := resourceAlicloudBastionhostLdapAuthServerUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_bastionhost_ldap_auth_server"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_bastionhost_ldap_auth_server", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_bastionhost_ldap_auth_server"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetInstanceLDAPAuthServer" {
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
		err := resourceAlicloudBastionhostLdapAuthServerRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	err = resourceAlicloudBastionhostLdapAuthServerDelete(dExisted, rawClient)
	assert.Nil(t, err)
}
