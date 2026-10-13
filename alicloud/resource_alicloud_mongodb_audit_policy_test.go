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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// lintignore: AT001
func TestAccAliCloudMongoDBAuditPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_audit_policy.default"
	checkoutSupportedRegions(t, true, connectivity.MongoDBSupportRegions)
	ra := resourceAttrInit(resourceId, AliCloudMongoDBAuditPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongodbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongodbAuditPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-mongodbauditpolicy-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongoDBAuditPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":     "${alicloud_mongodb_instance.default.id}",
					"audit_status":       "enable",
					"service_type":       "V2_Standard",
					"storage_period":     "30",
					"hot_storage_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":     CHECKSET,
						"audit_status":       "enable",
						"service_type":       "V2_Standard",
						"storage_period":     "30",
						"hot_storage_period": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_type":       "V2_Standard",
					"storage_period":     "180",
					"hot_storage_period": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_period":     "180",
						"hot_storage_period": "5",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_type"},
			},
		},
	})
}

// lintignore: AT001
func TestAccAliCloudMongoDBAuditPolicy_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_audit_policy.default"
	checkoutSupportedRegions(t, true, connectivity.MongoDBSupportRegions)
	ra := resourceAttrInit(resourceId, AliCloudMongoDBAuditPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongodbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongodbAuditPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-mongodbauditpolicy-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongoDBAuditPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":     "${alicloud_mongodb_instance.default.id}",
					"audit_status":       "enable",
					"storage_period":     "30",
					"service_type":       "V2_Standard",
					"hot_storage_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":     CHECKSET,
						"audit_status":       "enable",
						"storage_period":     "30",
						"service_type":       "V2_Standard",
						"hot_storage_period": "7",
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

var AliCloudMongoDBAuditPolicyMap0 = map[string]string{
	"storage_period": CHECKSET,
	"filter":         CHECKSET,
	"service_type":   "V2_Standard",
}

func AliCloudMongoDBAuditPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_mongodb_zones" "default" {}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
}

resource "alicloud_mongodb_instance" "default" {
  engine_version      = "4.2"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  name                = var.name
  vswitch_id          = data.alicloud_vswitches.default.ids.0
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}
`, name)
}

// lintignore: R001
func TestUnitAliCloudMongoDBAuditPolicy(t *testing.T) {
	p := Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_mongodb_audit_policy"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_mongodb_audit_policy"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"db_instance_id": "CreateValue",
		"audit_status":   "enable",
		"storage_period": 20,
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
		// DescribeAuditPolicy
		"DBInstanceId":   "CreateValue",
		"LogAuditStatus": "Enable",
		"DBInstances": map[string]interface{}{
			"DBInstance": []interface{}{
				map[string]interface{}{
					"DBInstanceStatus": "Running",
					"DBInstanceId":     "CreateValue",
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// ResetAccountPassword
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_mongodb_audit_policy", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudMongodbAuditPolicyCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeAuditPolicy Response
		"DBInstanceId":   "CreateValue",
		"LogAuditStatus": "Enable",
		"DBInstances": map[string]interface{}{
			"DBInstance": []interface{}{
				map[string]interface{}{
					"DBInstanceStatus": "Running",
					"DBInstanceId":     "CreateValue",
				},
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyAuditPolicy" {
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
		err := resourceAliCloudMongodbAuditPolicyCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mongodb_audit_policy"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudMongodbAuditPolicyUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyAccountDescription
	attributesDiff := map[string]interface{}{
		"audit_status":   "disabled",
		"storage_period": 30,
	}
	diff, err := newInstanceDiff("alicloud_mongodb_audit_policy", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_mongodb_audit_policy"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeAuditPolicy Response
		"LogAuditStatus": "Disabled",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyAuditPolicy" {
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
		err := resourceAliCloudMongodbAuditPolicyUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mongodb_audit_policy"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeAuditPolicy" {
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
		err := resourceAliCloudMongodbAuditPolicyRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	err = resourceAliCloudMongodbAuditPolicyDelete(dExisted, rawClient)
	assert.Nil(t, err)

}

// Test Mongodb AuditPolicy. >>> Resource test cases, automatically generated.
// Case 审计日志TF覆盖 11599
// lintignore: AT001
func TestAccAliCloudMongodbAuditPolicy_basic11599(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_audit_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudMongodbAuditPolicyMap11599)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongodbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongodbAuditPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmongodb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongodbAuditPolicyBasicDependence11599)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id": "${alicloud_mongodb_instance.default.id}",
					"audit_status":   "enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": CHECKSET,
						"audit_status":   "enable",
						"service_type":   "Standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"audit_status": "disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_status": "disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"audit_status": "enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_status": "enable",
						"service_type": "Standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_period": "12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_period": "12",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"filter": "update,insert,delete",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"filter": "update,insert,delete",
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

// lintignore: AT001
func TestAccAliCloudMongodbAuditPolicy_basic11599_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_audit_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudMongodbAuditPolicyMap11599)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongodbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongodbAuditPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmongodb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongodbAuditPolicyBasicDependence11599)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id": "${alicloud_mongodb_instance.default.id}",
					"audit_status":   "enable",
					"storage_period": "12",
					"filter":         "update,insert,delete",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": CHECKSET,
						"audit_status":   "enable",
						"storage_period": "12",
						"filter":         "update,insert,delete",
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

var AliCloudMongodbAuditPolicyMap11599 = map[string]string{
	"storage_period": CHECKSET,
	"filter":         CHECKSET,
	"service_type":   "Standard",
}

func AliCloudMongodbAuditPolicyBasicDependence11599(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

# mongo.x8.large/local_ssd is not offered in every zone (e.g. cn-beijing-a is closed for it),
# so create a dedicated VPC/VSwitch in a zone that currently has stock for this spec.
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-h"
}

resource "alicloud_mongodb_instance" "default" {
  engine_version      = "4.2"
  db_instance_class   = "mongo.x8.large"
  db_instance_storage = 50
  storage_engine      = "WiredTiger"
  storage_type        = "local_ssd"
  name                = var.name
  vswitch_id          = alicloud_vswitch.default.id
}
`, name)
}

// TestAccAliCloudMongoDBAuditPolicy_disabledImport verifies that importing a resource
// whose audit_status is `disabled` does not create a perpetual diff on `service_type`.
// This is the core scenario the DiffSuppressFunc + disabled-preserve Read logic exist
// for: on import, state starts empty; if Read wrote back the API's actual value (which
// the server flips to the trial edition while disabled) the plan would forever want to
// change service_type back to the config's declared value. The step chain is:
//
//	create (audit_status = enable, service_type = Standard)
//	  → toggle audit_status to disabled
//	    → import
//	      → follow-up plan must be empty (verified implicitly by ImportStateVerify).
//
// lintignore: AT001
func TestAccAliCloudMongoDBAuditPolicy_disabledImport(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_audit_policy.default"
	checkoutSupportedRegions(t, true, connectivity.MongoDBSupportRegions)
	ra := resourceAttrInit(resourceId, map[string]string{
		"storage_period": CHECKSET,
		"filter":         CHECKSET,
		"service_type":   "Standard",
	})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongodbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongodbAuditPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-mongodbauditpolicy-di-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongoDBAuditPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id": "${alicloud_mongodb_instance.default.id}",
					"audit_status":   "enable",
					"service_type":   "Standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": CHECKSET,
						"audit_status":   "enable",
						"service_type":   "Standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id": "${alicloud_mongodb_instance.default.id}",
					"audit_status":   "disabled",
					"service_type":   "Standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_status": "disabled",
						"service_type": "Standard",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// When audit is disabled the resource Read intentionally skips writing service_type
				// (see Read's disabled branch), so a freshly imported state has an empty value while
				// the pre-import state carries "Standard". Skip the raw state comparison for this
				// field; the follow-up plan the test framework runs after import is what actually
				// verifies DiffSuppress prevents a perpetual diff.
				ImportStateVerifyIgnore: []string{"service_type"},
			},
		},
	})
}

// Test Mongodb AuditPolicy. <<< Resource test cases, automatically generated.
