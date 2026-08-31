package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudPvtzUserVpcAuthorization_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pvtz_user_vpc_authorization.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivateZoneUserVpcAuthorizationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PvtzService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePvtzUserVpcAuthorization")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatezoneuservpcauthorization%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivateZoneUserVpcAuthorizationBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEnterpriseAccountEnabled(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"authorized_user_id": "1359528198477648",
					"auth_channel":       "RESOURCE_DIRECTORY",
					"auth_type":          "NORMAL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"authorized_user_id": CHECKSET,
						"auth_type":          "NORMAL",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"auth_channel"},
			},
		},
	})
}

// defaultAuthType exercises the bug-fix path: auth_type is omitted so the
// resource ID must default to "<authorized_user_id>:NORMAL" instead of the
// previous buggy "<authorized_user_id>:<nil>".
func TestAccAliCloudPvtzUserVpcAuthorization_defaultAuthType(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pvtz_user_vpc_authorization.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivateZoneUserVpcAuthorizationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PvtzService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePvtzUserVpcAuthorization")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatezoneuservpcauthorization%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivateZoneUserVpcAuthorizationBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEnterpriseAccountEnabled(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"authorized_user_id": "1359528198477648",
					"auth_channel":       "RESOURCE_DIRECTORY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"authorized_user_id": CHECKSET,
						"auth_type":          "NORMAL",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"auth_channel"},
			},
		},
	})
}

var AlicloudPrivateZoneUserVpcAuthorizationMap0 = map[string]string{}

func AlicloudPrivateZoneUserVpcAuthorizationBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudPvtzUserVpcAuthorization(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_pvtz_user_vpc_authorization"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_pvtz_user_vpc_authorization"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"authorized_user_id": "authorized_user_id",
		"auth_channel":       "RESOURCE_DIRECTORY",
		"auth_type":          "NORMAL",
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
		"Users": []interface{}{
			map[string]interface{}{
				"AuthType":         "NORMAL",
				"AuthorizedUserId": "authorized_user_id",
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
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_pvtz_user_vpc_authorization", "NORMAL"))
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
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewPvtzClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudPvtzUserVpcAuthorizationCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("System.Busy")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudPvtzUserVpcAuthorizationCreate(d, rawClient)
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
		err := resourceAlicloudPvtzUserVpcAuthorizationCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Create with auth_type omitted: the resource ID must default to
	// "<authorized_user_id>:NORMAL" instead of the previous buggy
	// "<authorized_user_id>:<nil>".
	t.Run("CreateNormalWithDefaultAuthType", func(t *testing.T) {
		dDefault, _ := schema.InternalMap(p["alicloud_pvtz_user_vpc_authorization"].Schema).Data(nil, nil)
		dDefault.MarkNewResource()
		err := dDefault.Set("authorized_user_id", "authorized_user_id")
		assert.Nil(t, err)
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["CreateNormal"]("")
		})
		err = resourceAlicloudPvtzUserVpcAuthorizationCreate(dDefault, rawClient)
		patches.Reset()
		assert.Nil(t, err)
		assert.Equal(t, "authorized_user_id:NORMAL", dDefault.Id())
	})

	// Read with a legacy "<nil>" AuthType in the ID: the migration must
	// rewrite it to "NORMAL" so the API receives the correct AuthType.
	t.Run("ReadMigrateLegacyNilAuthType", func(t *testing.T) {
		dLegacy, _ := schema.InternalMap(p["alicloud_pvtz_user_vpc_authorization"].Schema).Data(nil, nil)
		dLegacy.SetId("authorized_user_id:<nil>")
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudPvtzUserVpcAuthorizationRead(dLegacy, rawClient)
		patches.Reset()
		assert.Nil(t, err)
		assert.Equal(t, "authorized_user_id:NORMAL", dLegacy.Id())
		assert.Equal(t, "NORMAL", dLegacy.Get("auth_type"))
	})

	// Delete with a legacy "<nil>" AuthType in the ID: the migration must
	// send "NORMAL" to the DeleteUserVpcAuthorization API.
	t.Run("DeleteMigrateLegacyNilAuthType", func(t *testing.T) {
		dLegacy, _ := schema.InternalMap(p["alicloud_pvtz_user_vpc_authorization"].Schema).Data(nil, nil)
		dLegacy.SetId("authorized_user_id:<nil>")
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudPvtzUserVpcAuthorizationDelete(dLegacy, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("authorized_user_id", ":", "NORMAL"))

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewPvtzClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudPvtzUserVpcAuthorizationDelete(d, rawClient)
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
		err := resourceAlicloudPvtzUserVpcAuthorizationDelete(d, rawClient)
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
		err := resourceAlicloudPvtzUserVpcAuthorizationDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("DeleteIsExpectedErrors", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("ZoneVpc.Auth.NotExists")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudPvtzUserVpcAuthorizationDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribePvtzUserVpcAuthorizationNotFound", func(t *testing.T) {
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
		err := resourceAlicloudPvtzUserVpcAuthorizationRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribePvtzUserVpcAuthorizationAbnormal", func(t *testing.T) {
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
		err := resourceAlicloudPvtzUserVpcAuthorizationRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}
