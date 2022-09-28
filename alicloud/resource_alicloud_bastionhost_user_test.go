package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudBastionhostUser_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_user.default"
	ra := resourceAttrInit(resourceId, AlicloudBastionhostUserMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhostuser%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostUserBasicDependence0)
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
					"user_name":   "tf-testAccBastionHostUser-12345",
					"source":      "Local",
					"instance_id": "${data.alicloud_bastionhost_instances.default.ids.0}",
					"password":    "tf-testAcc-oAupFqRaH24MdOSrsIKsu3qw",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_name":   "tf-testAccBastionHostUser-12345",
						"source":      "Local",
						"instance_id": CHECKSET,
						"password":    "tf-testAcc-oAupFqRaH24MdOSrsIKsu3qw",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment": "tf-testAcc-mrPFCPi3MuIloLzTvVzQbUbs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment": "tf-testAcc-mrPFCPi3MuIloLzTvVzQbUbs",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"email": "tf-testAcc-5V8AgQKKw389irWIePb47aOq@163.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"email": "tf-testAcc-5V8AgQKKw389irWIePb47aOq@163.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "tf-testAcc-RZEdvPXF9A3w3ArhFwuAfUoY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "tf-testAcc-RZEdvPXF9A3w3ArhFwuAfUoY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "CN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "CN",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "HK",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "HK",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "MO",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "MO",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "TW",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "TW",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "RU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "RU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "SG",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "SG",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "MY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "MY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "ID",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "ID",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "DE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "DE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "AU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "AU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "US",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "US",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "AE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "AE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "JP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "JP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "GB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "GB",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "IN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "IN",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "KR",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "KR",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "PH",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "PH",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "CH",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "CH",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "SE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "SE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile": "702345672",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile": "702345672",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Frozen",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Frozen",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "tf-testAcc-Li6bvnYmD9ryuLUt2Wsdn4gy",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "tf-testAcc-Li6bvnYmD9ryuLUt2Wsdn4gy",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment":             "tf-testAcc-X23IfHiv8DnMoYChjEnb6X2h",
					"email":               "tf-testAcc-Rw5hfV8W1mkMO44chYBC07sq@163.com",
					"display_name":        "tf-testAcc-yAwB1akRJGW9RVMaTEdOHOHS",
					"mobile_country_code": "CN",
					"mobile":              "13312345678",
					"password":            "tf-testAcc-lBdrpSbUJ4Ddw9oSGCeI2u2p",
					"status":              "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment":             "tf-testAcc-X23IfHiv8DnMoYChjEnb6X2h",
						"email":               "tf-testAcc-Rw5hfV8W1mkMO44chYBC07sq@163.com",
						"display_name":        "tf-testAcc-yAwB1akRJGW9RVMaTEdOHOHS",
						"mobile_country_code": "CN",
						"mobile":              "13312345678",
						"password":            "tf-testAcc-lBdrpSbUJ4Ddw9oSGCeI2u2p",
						"status":              "Normal",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

var AlicloudBastionhostUserMap0 = map[string]string{
	"display_name":        CHECKSET,
	"status":              CHECKSET,
	"instance_id":         CHECKSET,
	"mobile_country_code": "",
	"user_id":             CHECKSET,
}

func AlicloudBastionhostUserBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_bastionhost_instances" "default" {}

`, name)
}
func TestAccAlicloudBastionhostUser_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_user.default"
	ra := resourceAttrInit(resourceId, AlicloudBastionhostUserMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhostuser%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostUserBasicDependence1)
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
					"user_name":      "tf-testAccBastionhostUserRam-123456",
					"source":         "Ram",
					"instance_id":    "${data.alicloud_bastionhost_instances.default.ids.0}",
					"source_user_id": "247823888127488180",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_name":      "tf-testAccBastionhostUserRam-123456",
						"source":         "Ram",
						"instance_id":    CHECKSET,
						"source_user_id": "247823888127488180",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment": "tf-testAccBastionhostUserRam-123456",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment": "tf-testAccBastionhostUserRam-123456",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"email": "tf-testAcc-LmwD6dS7fyO93I@163.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"email": "tf-testAcc-LmwD6dS7fyO93I@163.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "tf-testAccBastionhostUserRam-456789",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "tf-testAccBastionhostUserRam-456789",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "CN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "CN",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "HK",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "HK",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "MO",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "MO",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "TW",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "TW",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "RU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "RU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "SG",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "SG",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "MY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "MY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "ID",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "ID",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "DE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "DE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "AU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "AU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "US",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "US",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "AE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "AE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "JP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "JP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "GB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "GB",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "IN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "IN",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "KR",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "KR",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "PH",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "PH",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "CH",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "CH",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile_country_code": "SE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile_country_code": "SE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile": "702345672",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile": "702345672",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "tf-testAcc-2MeAHvjV3LvFsGfUSs73hXaI",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "tf-testAcc-2MeAHvjV3LvFsGfUSs73hXaI",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Frozen",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Frozen",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment":             "tf-testAcc-2MeAHvjV3LvFsGfUSs73hXaI",
					"email":               "tf-testAcc-75MYawy06OnL4zTD4xdi6n4T@163.com",
					"display_name":        "tf-testAcc-2MeAHvjV3LvFsGfUSs73hXaI",
					"mobile_country_code": "CN",
					"mobile":              "13312345678",
					"password":            "tf-testAcc-2MeAHvjV3LvFsGfUSs73hXaI",
					"status":              "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment":             "tf-testAcc-2MeAHvjV3LvFsGfUSs73hXaI",
						"email":               "tf-testAcc-75MYawy06OnL4zTD4xdi6n4T@163.com",
						"display_name":        "tf-testAcc-2MeAHvjV3LvFsGfUSs73hXaI",
						"mobile_country_code": "CN",
						"mobile":              "13312345678",
						"password":            "tf-testAcc-2MeAHvjV3LvFsGfUSs73hXaI",
						"status":              "Normal",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

var AlicloudBastionhostUserMap1 = map[string]string{
	"user_id":             CHECKSET,
	"display_name":        CHECKSET,
	"status":              CHECKSET,
	"instance_id":         CHECKSET,
	"mobile_country_code": "",
}

func AlicloudBastionhostUserBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_bastionhost_instances" "default" {}

`, name)
}

func TestUnitAlicloudBastionhostUser(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_bastionhost_user"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_bastionhost_user"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"comment":             "CreateBastionhostUserValue",
		"display_name":        "CreateBastionhostUserValue",
		"email":               "CreateBastionhostUserValue",
		"instance_id":         "CreateBastionhostUserValue",
		"mobile_country_code": "CreateBastionhostUserValue",
		"mobile":              "CreateBastionhostUserValue",
		"password":            "CreateBastionhostUserValue",
		"source":              "CreateBastionhostUserValue",
		"user_name":           "CreateBastionhostUserValue",
		"source_user_id":      "CreateBastionhostUserValue",
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
		// GetUser
		"User": map[string]interface{}{
			"InstanceId":        "CreateBastionhostUserValue",
			"UserId":            "CreateBastionhostUserValue",
			"Comment":           "CreateBastionhostUserValue",
			"DisplayName":       "CreateBastionhostUserValue",
			"Email":             "CreateBastionhostUserValue",
			"Mobile":            "CreateBastionhostUserValue",
			"MobileCountryCode": "CreateBastionhostUserValue",
			"Source":            "CreateBastionhostUserValue",
			"SourceUserId":      "CreateBastionhostUserValue",
			"UserState": []interface{}{
				"CreateBastionhostUserValue",
			},
			"UserName": "CreateBastionhostUserValue",
		},
		"UserId": "CreateBastionhostUserValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateUser
		"User": map[string]interface{}{
			"UserId": "CreateBastionhostUserValue",
		},
		"UserId": "CreateBastionhostUserValue",
	}
	ReadMockResponseDiff := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_bastionhost_user", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	t.Run("Create", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewBastionhostClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudBastionhostUserCreate(dInit, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		ReadMockResponseDiff = map[string]interface{}{
			// GetUser Response
			"User": map[string]interface{}{
				"UserId": "CreateBastionhostUserValue",
			},
		}
		errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1 // a counter used to cover retry scenario; the same below
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "CreateUser" {
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
			err := resourceAlicloudBastionhostUserCreate(dInit, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_bastionhost_user"].Schema).Data(dInit.State(), nil)
				for key, value := range attributes {
					_ = dCompare.Set(key, value)
				}
				assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
			}
			if retryIndex >= len(errorCodes)-1 {
				break
			}
		}
	})

	// Update
	t.Run("Update", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewBastionhostClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudBastionhostUserUpdate(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		// ModifyUser
		attributesDiff := map[string]interface{}{
			"comment":             "ModifyUserValue",
			"display_name":        "ModifyUserValue",
			"email":               "ModifyUserValue",
			"mobile":              "ModifyUserValue",
			"mobile_country_code": "ModifyUserValue",
			"password":            "ModifyUserValue",
			"source":              "ModifyUserValue",
		}
		diff, err := newInstanceDiff("alicloud_bastionhost_user", attributes, attributesDiff, dInit.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_bastionhost_user"].Schema).Data(dInit.State(), diff)
		ReadMockResponseDiff = map[string]interface{}{
			// GetUser Response
			"User": map[string]interface{}{
				"Comment":           "ModifyUserValue",
				"DisplayName":       "ModifyUserValue",
				"Email":             "ModifyUserValue",
				"Mobile":            "ModifyUserValue",
				"MobileCountryCode": "ModifyUserValue",
				"Password":          "ModifyUserValue",
				"Source":            "ModifyUserValue",
			},
		}
		errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "ModifyUser" {
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
			err := resourceAlicloudBastionhostUserUpdate(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_bastionhost_user"].Schema).Data(dExisted.State(), nil)
				for key, value := range attributes {
					_ = dCompare.Set(key, value)
				}
				assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
			}
			if retryIndex >= len(errorCodes)-1 {
				break
			}
		}

		// UnlockUsers
		attributesDiff = map[string]interface{}{
			"status": "Normal",
		}
		diff, err = newInstanceDiff("alicloud_bastionhost_user", attributes, attributesDiff, dExisted.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_bastionhost_user"].Schema).Data(dExisted.State(), diff)
		ReadMockResponseDiff = map[string]interface{}{
			// GetUser Response
			"User": map[string]interface{}{
				"UserState": []interface{}{
					"Normal",
				},
			},
		}
		errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "UnlockUsers" {
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
			err := resourceAlicloudBastionhostUserUpdate(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_bastionhost_user"].Schema).Data(dExisted.State(), nil)
				for key, value := range attributes {
					_ = dCompare.Set(key, value)
				}
				assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
			}
			if retryIndex >= len(errorCodes)-1 {
				break
			}
		}

		// LockUsers
		attributesDiff = map[string]interface{}{
			"status": "Frozen",
		}
		diff, err = newInstanceDiff("alicloud_bastionhost_user", attributes, attributesDiff, dExisted.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_bastionhost_user"].Schema).Data(dExisted.State(), diff)
		ReadMockResponseDiff = map[string]interface{}{
			// GetUser Response
			"User": map[string]interface{}{
				"UserState": []interface{}{
					"Frozen",
				},
			},
		}
		errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "LockUsers" {
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
			err := resourceAlicloudBastionhostUserUpdate(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_bastionhost_user"].Schema).Data(dExisted.State(), nil)
				for key, value := range attributes {
					_ = dCompare.Set(key, value)
				}
				assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
			}
			if retryIndex >= len(errorCodes)-1 {
				break
			}
		}
	})

	// Read
	t.Run("Read", func(t *testing.T) {
		errorCodes := []string{"NonRetryableError", "Throttling", "nil", "{}"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "GetUser" {
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
			err := resourceAlicloudBastionhostUserRead(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			case "{}":
				assert.Nil(t, err)
			}
		}
	})

	// Delete
	t.Run("Delete", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewBastionhostClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudBastionhostUserDelete(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		attributesDiff := map[string]interface{}{}
		diff, err := newInstanceDiff("alicloud_bastionhost_user", attributes, attributesDiff, dInit.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_bastionhost_user"].Schema).Data(dInit.State(), diff)
		errorCodes := []string{"NonRetryableError", "Throttling", "nil", "OBJECT_NOT_FOUND"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "DeleteUser" {
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
			err := resourceAlicloudBastionhostUserDelete(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			case "OBJECT_NOT_FOUND":
				assert.Nil(t, err)
			}
		}
	})
}
