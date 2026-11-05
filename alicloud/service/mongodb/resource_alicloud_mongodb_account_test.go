package mongodb_test

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mongodb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkTerraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudMongoDBAccount_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_account.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.MongoDBSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMongoDBAccountMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mongodb.NewMongodbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMongodbAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smongodbaccount%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMongoDBAccountBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":         "${alicloud_mongodb_instance.default.id}",
					"account_name":        "root",
					"account_password":    "YourPassword+12345",
					"account_description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":         acctest.CHECKSET,
						"account_name":        "root",
						"account_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword+12345678",
				}),
				Check: resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "${var.name}",
					"account_password":    "YourPassword+12345",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
		},
	})
}

var AlicloudMongoDBAccountMap0 = map[string]string{
	"account_name":   acctest.CHECKSET,
	"instance_id":    acctest.CHECKSET,
	"status":         acctest.CHECKSET,
	"character_type": "normal",
}

func AlicloudMongoDBAccountBasicDependence0(name string) string {
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
resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_mongodb_zones.default.zones.0.id
  vswitch_name      = "subnet-for-local-test"
}
resource "alicloud_mongodb_instance" "default" {
  engine_version      = "4.2"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  name                = var.name
  vswitch_id          = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudMongoDBAccount(t *testing.T) {
	p := acctest.Provider.ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_mongodb_account"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_mongodb_account"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"instance_id":      "ResetAccountPasswordValue",
		"account_name":     "ResetAccountPasswordValue",
		"account_password": "ResetAccountPasswordValue",
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
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// DescribeAccounts
		"Accounts": map[string]interface{}{
			"Account": []interface{}{
				map[string]interface{}{
					"AccountName":        "ResetAccountPasswordValue",
					"DBInstanceId":       "ResetAccountPasswordValue",
					"AccountDescription": "ResetAccountPasswordValue",
					"AccountStatus":      "Available",
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// ResetAccountPassword
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       helper.String(errorCode),
			Data:       helper.String(errorCode),
			Message:    helper.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_mongodb_account", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		// The account lookup is database scoped and requires the account name,
		// so partial mocks inherit the identity fields from the base response.
		if accounts, ok := ReadMockResponse["Accounts"].(map[string]interface{}); ok {
			if list, ok := accounts["Account"].([]interface{}); ok {
				for _, item := range list {
					account, ok := item.(map[string]interface{})
					if !ok {
						continue
					}
					if _, ok := account["AccountName"]; !ok {
						account["AccountName"] = "ResetAccountPasswordValue"
					}
					if _, ok := account["DBName"]; !ok {
						account["DBName"] = "admin"
					}
					if _, ok := account["CharacterType"]; !ok {
						account["CharacterType"] = "normal"
					}
				}
			}
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = mongodb.ResourceAliCloudMongodbAccountCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeAccounts Response
		"Accounts": map[string]interface{}{
			"Account": []interface{}{
				map[string]interface{}{
					"AccountName":        "ResetAccountPasswordValue",
					"DBInstanceId":       "ResetAccountPasswordValue",
					"AccountDescription": "ResetAccountPasswordValue",
					"AccountStatus":      "Available",
				},
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ResetAccountPassword" {
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
		err := mongodb.ResourceAliCloudMongodbAccountCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mongodb_account"].Schema).Data(dInit.State(), nil)
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
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = mongodb.ResourceAliCloudMongodbAccountUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyAccountDescription
	attributesDiff := map[string]interface{}{
		"account_description": "ModifyAccountDescriptionValue",
	}
	diff, err := acctest.NewInstanceDiff("alicloud_mongodb_account", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_mongodb_account"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeAccounts Response
		"Accounts": map[string]interface{}{
			"Account": []interface{}{
				map[string]interface{}{
					"AccountDescription": "ModifyAccountDescriptionValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyAccountDescription" {
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
		err := mongodb.ResourceAliCloudMongodbAccountUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mongodb_account"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ModifyAccountDescription
	attributesDiff = map[string]interface{}{
		"account_password": "UpdateValue",
	}
	diff, err = acctest.NewInstanceDiff("alicloud_mongodb_account", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_mongodb_account"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeAccounts Response
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ResetAccountPassword" {
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
		err := mongodb.ResourceAliCloudMongodbAccountUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mongodb_account"].Schema).Data(dExisted.State(), nil)
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
			if *action == "DescribeAccounts" {
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
		err := mongodb.ResourceAliCloudMongodbAccountRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	err = mongodb.ResourceAliCloudMongodbAccountDelete(dExisted, rawClient)
	assert.Nil(t, err)

}

// Test Mongodb Account. >>> Resource test cases, automatically generated.
// Case 账号测试用例初始化_分片集群_create账号覆盖度 9243
func TestAccAliCloudMongodbAccount_basic9243(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMongodbAccountMap9243)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mongodb.NewMongodbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMongodbAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMongodbAccountBasicDependence9243)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_name":     name,
					"account_password": "qwer1234!!!",
					"instance_id":      "${alicloud_mongodb_sharding_instance.default7eOftZ.id}",
					"character_type":   "db",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name":     name,
						"account_password": "qwer1234!!!",
						"instance_id":      acctest.CHECKSET,
						"character_type":   "db",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
		},
	})
}

var AlicloudMongodbAccountMap9243 = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudMongodbAccountBasicDependence9243(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-shanghai-b"
}

variable "region_id" {
  default = "cn-shanghai"
}

variable "other_zone_id" {
  default = "cn-shanghai-g"
}

resource "alicloud_vpc" "defaultIzCngN" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultLWEsPM" {
  vpc_id     = alicloud_vpc.defaultIzCngN.id
  zone_id    = var.zone_id
  cidr_block = "10.0.0.0/24"
}

resource "alicloud_mongodb_sharding_instance" "default7eOftZ" {
  engine_version = "4.4"
  storage_type   = "cloud_essd1"
  vswitch_id     = alicloud_vswitch.defaultLWEsPM.id
  vpc_id         = alicloud_vpc.defaultIzCngN.id
  storage_engine = "WiredTiger"
  network_type   = "VPC"
  zone_id        = var.zone_id
  name           = var.name
  mongo_list {
    node_class = "mdb.shard.4x.large.d"
  }
  mongo_list {
    node_class = "mdb.shard.4x.large.d"
  }
  shard_list {
    node_class   = "mdb.shard.4x.large.d"
    node_storage = "20"
  }
  shard_list {
    node_class   = "mdb.shard.4x.large.d"
    node_storage = "20"
  }
  config_server_list {
    node_class        = "mdb.shard.2x.xlarge.d"
    node_storage      = "80"
  }
}


`, name)
}

// Case 账号测试用例初始化_副本集_2024年11月7日_写入 8763
func TestAccAliCloudMongodbAccount_basic8763(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMongodbAccountMap8763)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mongodb.NewMongodbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMongodbAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMongodbAccountBasicDependence8763)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_name":        "root",
					"account_password":    "qwer1234!!!",
					"instance_id":         "${alicloud_mongodb_instance.default7eOftZ.id}",
					"character_type":      "normal",
					"account_description": "bgg-test-power",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name":        "root",
						"account_password":    "qwer1234!!!",
						"instance_id":         acctest.CHECKSET,
						"character_type":      "normal",
						"account_description": "bgg-test-power",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
		},
	})
}

var AlicloudMongodbAccountMap8763 = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudMongodbAccountBasicDependence8763(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-shanghai-b"
}

variable "region_id" {
  default = "cn-shanghai"
}

variable "other_zone_id" {
  default = "cn-shanghai-g"
}

resource "alicloud_vpc" "defaultIzCngN" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultLWEsPM" {
  vpc_id     = alicloud_vpc.defaultIzCngN.id
  zone_id    = var.zone_id
  cidr_block = "10.0.0.0/24"
}

resource "alicloud_mongodb_instance" "default7eOftZ" {
  engine_version      = "4.4"
  storage_type        = "cloud_essd1"
  vswitch_id          = alicloud_vswitch.defaultLWEsPM.id
  vpc_id              = alicloud_vpc.defaultIzCngN.id
  storage_engine      = "WiredTiger"
  network_type        = "VPC"
  zone_id             = var.zone_id
  db_instance_storage = "20"
  db_instance_class   = "mdb.shard.4x.large.d"
}


`, name)
}

// Test Mongodb Account. <<< Resource test cases, automatically generated.

// TestAccAliCloudMongodbAccount_shardingAccountRoles covers the product decision
// for CreateAccount roles: shard direct-connection (db) accounts reject roles
// with InvalidRoles.NotSupported, config server (cs) accounts require and honor
// custom roles.
func TestAccAliCloudMongodbAccount_shardingAccountRoles(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_roles_%d", rand)
	dbAccountName := fmt.Sprintf("tf_db_%d", rand)
	dbRolesAccountName := fmt.Sprintf("tf_dbr_%d", rand)
	csAccountName := fmt.Sprintf("tf_cs_%d", rand)
	dbResourceId := "alicloud_mongodb_account.db"
	csResourceId := "alicloud_mongodb_account.cs"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", "cn-shanghai")
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbAccountShardingRolesConfig(name, dbAccountName, dbRolesAccountName, csAccountName, "qwer1234!!!", false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dbResourceId, "account_name", dbAccountName),
					resource.TestCheckResourceAttr(dbResourceId, "character_type", "db"),
					resource.TestCheckResourceAttr(dbResourceId, "db_name", "admin"),
					resource.TestCheckResourceAttrSet(dbResourceId, "instance_id"),
					resource.TestCheckResourceAttrSet(dbResourceId, "status"),
				),
			},
			{
				// A shard direct-connection account must not carry roles; the API
				// rejects the create with InvalidRoles.NotSupported.
				Config:      testAccMongodbAccountShardingRolesConfig(name, dbAccountName, dbRolesAccountName, csAccountName, "qwer1234!!!", true, false),
				ExpectError: regexp.MustCompile("InvalidRoles.NotSupported"),
			},
			{
				// Config server accounts require and honor custom roles.
				Config: testAccMongodbAccountShardingRolesConfig(name, dbAccountName, dbRolesAccountName, csAccountName, "qwer1234!!!", false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(csResourceId, "account_name", csAccountName),
					resource.TestCheckResourceAttr(csResourceId, "character_type", "cs"),
					resource.TestCheckResourceAttr(csResourceId, "db_name", "admin"),
					resource.TestCheckResourceAttr(csResourceId, "roles.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(csResourceId, "roles.*", map[string]string{
						"role":    "read",
						"db_name": "admin",
					}),
				),
			},
			{
				Config: testAccMongodbAccountShardingRolesConfig(name, dbAccountName, dbRolesAccountName, csAccountName, "qwer12345!!!", false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(csResourceId, "roles.#", "1"),
				),
			},
			{
				ResourceName:            csResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
		},
	})
}

func testAccMongodbAccountShardingRolesConfig(name, dbAccount, dbRolesAccount, csAccount, password string, withDbRoles, withCs bool) string {
	dbRolesBlock := ""
	if withDbRoles {
		dbRolesBlock = fmt.Sprintf(`
resource "alicloud_mongodb_account" "db_roles" {
  account_name     = "%s"
  account_password = "%s"
  instance_id      = alicloud_mongodb_sharding_instance.default.id
  character_type   = "db"
  db_name          = "admin"

  roles {
    role    = "read"
    db_name = "admin"
  }
}
`, dbRolesAccount, password)
	}
	csBlock := ""
	if withCs {
		csBlock = fmt.Sprintf(`
resource "alicloud_mongodb_account" "cs" {
  account_name     = "%s"
  account_password = "%s"
  instance_id      = alicloud_mongodb_sharding_instance.default.id
  character_type   = "cs"
  db_name          = "admin"

  roles {
    role    = "read"
    db_name = "admin"
  }
}
`, csAccount, password)
	}
	return fmt.Sprintf(`
%s

resource "alicloud_mongodb_account" "db" {
  account_name        = "%s"
  account_password    = "%s"
  instance_id         = alicloud_mongodb_sharding_instance.default.id
  character_type      = "db"
  db_name             = "admin"
  account_description = "%s"
}
%s
%s
`, testAccMongodbAccountShardingRolesDependence(name), dbAccount, password, name, dbRolesBlock, csBlock)
}

func testAccMongodbAccountShardingRolesDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_mongodb_zones" "default" {}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

locals {
  vswitch_zone_ids = [for zone in data.alicloud_zones.default.zones : zone.id]
  mongodb_zone_ids = [for zone in data.alicloud_mongodb_zones.default.zones : zone.id if contains(local.vswitch_zone_ids, zone.id)]
  zone_id          = local.mongodb_zone_ids[0]
}

resource "alicloud_vpc" "default" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = local.zone_id
  cidr_block = "10.0.0.0/24"
}

resource "alicloud_mongodb_sharding_instance" "default" {
  engine_version = "4.2"
  storage_type   = "cloud_essd1"
  vswitch_id     = alicloud_vswitch.default.id
  vpc_id         = alicloud_vpc.default.id
  storage_engine = "WiredTiger"
  network_type   = "VPC"
  zone_id        = local.zone_id
  name           = var.name

  mongo_list {
    node_class = "mdb.shard.4x.large.d"
  }
  mongo_list {
    node_class = "mdb.shard.4x.large.d"
  }
  shard_list {
    node_class   = "mdb.shard.4x.large.d"
    node_storage = "20"
  }
  shard_list {
    node_class   = "mdb.shard.4x.large.d"
    node_storage = "20"
  }
  config_server_list {
    node_class   = "mdb.shard.2x.xlarge.d"
    node_storage = "80"
  }
}
`, name)
}

func unitMongodbAccountResourceSchema(t *testing.T) map[string]*schema.Schema {
	t.Helper()
	return acctest.Provider.ResourcesMap["alicloud_mongodb_account"].Schema
}

func unitMongodbAccountResourceData(t *testing.T, id string, attrs map[string]interface{}, changes map[string]interface{}) *schema.ResourceData {
	t.Helper()
	state := &sdkTerraform.InstanceState{ID: id, Attributes: map[string]string{}}
	for k, v := range attrs {
		state.Attributes[k] = fmt.Sprint(v)
	}
	diff := &sdkTerraform.InstanceDiff{Attributes: map[string]*sdkTerraform.ResourceAttrDiff{}}
	for k, v := range changes {
		diff.Attributes[k] = &sdkTerraform.ResourceAttrDiff{Old: fmt.Sprint(attrs[k]), New: fmt.Sprint(v)}
	}
	d, err := schema.InternalMap(unitMongodbAccountResourceSchema(t)).Data(state, diff)
	if err != nil {
		t.Fatalf("failed to build resource data: %v", err)
	}
	return d
}

func TestUnitMongodbAccountSchemaContract(t *testing.T) {
	r := unitMongodbAccountResourceSchema(t)
	assert.True(t, r["db_name"].ForceNew, "db_name is part of the account identity")
	assert.True(t, r["db_name"].Computed)
	assert.True(t, r["roles"].Computed, "roles are returned by DescribeAccounts")
	assert.True(t, r["roles"].ForceNew)
	assert.True(t, r["character_type"].ForceNew)
	_, enumErrs := r["character_type"].ValidateFunc("cs", "character_type")
	assert.Empty(t, enumErrs, "cs must be a valid character_type")
	_, enumErrs = r["character_type"].ValidateFunc("db", "character_type")
	assert.Empty(t, enumErrs, "db must be a valid character_type")
}

func TestUnitMongodbAccountDescribeRejectsInvalidId(t *testing.T) {
	service := &mongodb.MongodbServiceV2{}
	if _, err := service.DescribeMongodbAccount("dds-1"); err == nil {
		t.Fatal("expected an error for an id without an account name")
	}
	if _, err := service.DescribeMongodbAccount("dds-1:db:acc:extra"); err == nil {
		t.Fatal("expected an error for an id with too many parts")
	}
}

func TestUnitMongodbAccountStateRefreshChecksFieldError(t *testing.T) {
	service := &mongodb.MongodbServiceV2{}
	refresh := service.MongodbAccountStateRefreshFuncWithApi("dds-1:admin:acc", "AccountStatus", nil, func(id string) (map[string]interface{}, error) {
		return map[string]interface{}{"DBName": "admin"}, nil
	})
	if _, _, err := refresh(); err == nil {
		t.Fatal("expected an error when the watched field is missing from the response")
	}
}

func TestUnitMongodbAccountReadNotFoundClearsId(t *testing.T) {
	d := unitMongodbAccountResourceData(t, "dds-1:acc", map[string]interface{}{
		"instance_id":  "dds-1",
		"account_name": "acc",
	}, nil)
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&mongodb.MongodbServiceV2{}), "DescribeMongodbAccount", func(_ *mongodb.MongodbServiceV2, id string) (map[string]interface{}, error) {
		return nil, errs.WrapErrorf(errs.NotFoundErr("Account", id), errs.NotFoundMsg, "mock")
	})
	defer patches.Reset()

	assert.Nil(t, mongodb.ResourceAliCloudMongodbAccountRead(d, &connectivity.AliyunClient{}))
	assert.Equal(t, "", d.Id())
}

func TestUnitMongodbAccountReadSetsDatabaseScopedId(t *testing.T) {
	d := unitMongodbAccountResourceData(t, "dds-1:acc1", map[string]interface{}{
		"instance_id":  "dds-1",
		"account_name": "acc1",
	}, nil)
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&mongodb.MongodbServiceV2{}), "DescribeMongodbAccount", func(_ *mongodb.MongodbServiceV2, id string) (map[string]interface{}, error) {
		return map[string]interface{}{
			"DBInstanceId":  "dds-1",
			"AccountName":   "acc1",
			"DBName":        "admin",
			"CharacterType": "db",
			"AccountStatus": "Available",
			"Roles": map[string]interface{}{
				"Role": []interface{}{map[string]interface{}{"DBName": "admin", "Role": "read"}},
			},
		}, nil
	})
	defer patches.Reset()

	assert.Nil(t, mongodb.ResourceAliCloudMongodbAccountRead(d, &connectivity.AliyunClient{}))
	assert.Equal(t, "dds-1:admin:acc1", d.Id())
	assert.Equal(t, "admin", d.Get("db_name"))
	assert.Equal(t, 1, d.Get("roles").(*schema.Set).Len())
}

func TestUnitMongodbAccountReadPropagatesError(t *testing.T) {
	d := unitMongodbAccountResourceData(t, "dds-1:acc", map[string]interface{}{
		"instance_id":  "dds-1",
		"account_name": "acc",
	}, nil)
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&mongodb.MongodbServiceV2{}), "DescribeMongodbAccount", func(_ *mongodb.MongodbServiceV2, id string) (map[string]interface{}, error) {
		return nil, errors.New("Forbidden.RAM")
	})
	defer patches.Reset()

	err := mongodb.ResourceAliCloudMongodbAccountRead(d, &connectivity.AliyunClient{})
	if err == nil {
		t.Fatal("expected the DescribeAccounts error to be propagated")
	}
	assert.Equal(t, "dds-1:acc", d.Id())
}

func TestUnitMongodbAccountUpdateDescriptionDescribeError(t *testing.T) {
	d := unitMongodbAccountResourceData(t, "dds-1:admin:acc", map[string]interface{}{
		"instance_id":  "dds-1",
		"account_name": "acc",
		"db_name":      "admin",
	}, map[string]interface{}{
		"account_description": "updated",
	})
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&mongodb.MongodbServiceV2{}), "DescribeMongodbAccount", func(_ *mongodb.MongodbServiceV2, id string) (map[string]interface{}, error) {
		return nil, errors.New("Forbidden.RAM")
	})
	defer patches.Reset()

	if err := mongodb.ResourceAliCloudMongodbAccountUpdate(d, &connectivity.AliyunClient{}); err == nil {
		t.Fatal("expected the pre-update DescribeAccounts error to be propagated")
	}
}

func TestUnitMongodbAccountUpdateShardDescriptionRejected(t *testing.T) {
	d := unitMongodbAccountResourceData(t, "dds-1:admin:acc", map[string]interface{}{
		"instance_id":  "dds-1",
		"account_name": "acc",
		"db_name":      "admin",
	}, map[string]interface{}{
		"account_description": "updated",
	})
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&mongodb.MongodbServiceV2{}), "DescribeMongodbAccount", func(_ *mongodb.MongodbServiceV2, id string) (map[string]interface{}, error) {
		return map[string]interface{}{"DBInstanceId": "dds-1", "AccountName": "acc", "DBName": "admin", "CharacterType": "db", "AccountStatus": "Available"}, nil
	})
	defer patches.Reset()

	err := mongodb.ResourceAliCloudMongodbAccountUpdate(d, &connectivity.AliyunClient{})
	if err == nil {
		t.Fatal("expected a shard account description update to be rejected")
	}
	assert.Contains(t, err.Error(), "replica set")
}

func TestUnitMongodbAccountCreateRejectsRolesForReplicaSet(t *testing.T) {
	d := schema.TestResourceDataRaw(t, unitMongodbAccountResourceSchema(t), map[string]interface{}{
		"instance_id":      "dds-1",
		"account_name":     "root",
		"account_password": "YourPassword+12345",
		"character_type":   "normal",
		"db_name":          "admin",
		"roles": []interface{}{
			map[string]interface{}{"role": "read", "db_name": "admin"},
		},
	})
	err := mongodb.ResourceAliCloudMongodbAccountCreate(d, &connectivity.AliyunClient{})
	if err == nil {
		t.Fatal("expected roles on a replica set account to be rejected")
	}
	assert.Contains(t, err.Error(), "config server account")
}

func TestUnitMongodbAccountCreateRejectsRolesForShardAccount(t *testing.T) {
	d := schema.TestResourceDataRaw(t, unitMongodbAccountResourceSchema(t), map[string]interface{}{
		"instance_id":      "dds-1",
		"account_name":     "shard_acc",
		"account_password": "YourPassword+12345",
		"character_type":   "db",
		"db_name":          "admin",
		"roles": []interface{}{
			map[string]interface{}{"role": "read", "db_name": "admin"},
		},
	})
	err := mongodb.ResourceAliCloudMongodbAccountCreate(d, &connectivity.AliyunClient{})
	if err == nil {
		t.Fatal("expected roles on a shard direct-connection account to be rejected")
	}
	assert.Contains(t, err.Error(), "InvalidRoles.NotSupported")
}

func TestUnitMongodbAccountCreateRequiresRolesForCsAccount(t *testing.T) {
	d := schema.TestResourceDataRaw(t, unitMongodbAccountResourceSchema(t), map[string]interface{}{
		"instance_id":      "dds-1",
		"account_name":     "cs_acc",
		"account_password": "YourPassword+12345",
		"character_type":   "cs",
		"db_name":          "admin",
	})
	err := mongodb.ResourceAliCloudMongodbAccountCreate(d, &connectivity.AliyunClient{})
	if err == nil {
		t.Fatal("expected a config server account without roles to be rejected")
	}
	assert.Contains(t, err.Error(), "required")
}

func TestUnitMongodbAccountDeleteNoop(t *testing.T) {
	d := unitMongodbAccountResourceData(t, "dds-1:acc", map[string]interface{}{}, nil)
	assert.Nil(t, mongodb.ResourceAliCloudMongodbAccountDelete(d, &connectivity.AliyunClient{}))
}
