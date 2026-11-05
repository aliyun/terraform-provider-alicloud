package mongodb_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mongodb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkTerraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	terraform "github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudMongoDBAccount_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_account.default"
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
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      testAccCheckAccountDependenciesDestroyed,
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

resource "alicloud_vswitch" "vswitch" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "10.0.0.0/24"
  zone_id      = local.zone_id
  vswitch_name = "subnet-for-local-test"
}

resource "alicloud_mongodb_instance" "default" {
  engine_version      = "4.2"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  name                = var.name
  vswitch_id          = alicloud_vswitch.vswitch.id
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

// TestAccAliCloudMongodbAccount_complete_sharding_account_roles covers the product
// decision for CreateAccount roles: shard direct-connection (db) accounts reject
// roles with InvalidRoles.NotSupported, config server (cs) accounts require and
// honor custom roles.
func TestAccAliCloudMongodbAccount_complete_sharding_account_roles(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-roles-%d", rand)
	dbAccountName := fmt.Sprintf("tf_db_%d", rand)
	dbRolesAccountName := fmt.Sprintf("tf_dbr_%d", rand)
	csAccountName := fmt.Sprintf("tf_cs_%d", rand)
	emptyRolesAccountName := fmt.Sprintf("tf_er_%d", rand)
	missingDbNameAccountName := fmt.Sprintf("tf_md_%d", rand)
	dbResourceId := "alicloud_mongodb_account.db"
	csResourceId := "alicloud_mongodb_account.cs"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", "cn-shanghai")
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAccountDependenciesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountShardingRolesConfig(name, dbAccountName, dbRolesAccountName, csAccountName, "qwer1234!!!", false, false, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dbResourceId, "account_name", dbAccountName),
					resource.TestCheckResourceAttr(dbResourceId, "character_type", "db"),
					resource.TestCheckResourceAttr(dbResourceId, "db_name", "admin"),
					resource.TestCheckResourceAttr(dbResourceId, "account_description", name),
					resource.TestCheckResourceAttrSet(dbResourceId, "instance_id"),
					resource.TestCheckResourceAttrSet(dbResourceId, "status"),
				),
			},
			{
				ResourceName:            dbResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "account_description"},
			},
			{
				// The legacy two-part id is accepted and upgraded in place.
				ResourceName:            dbResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccAccountLegacyImportId(dbResourceId),
				ImportStateVerifyIgnore: []string{"account_password", "account_description"},
			},
			{
				// A shard direct-connection account must not carry roles; the API
				// rejects the create with InvalidRoles.NotSupported. This is the
				// real service response, not a provider-side validation.
				Config:      testAccAccountShardingRolesConfig(name, dbAccountName, dbRolesAccountName, csAccountName, "qwer1234!!!", true, false, ""),
				ExpectError: regexp.MustCompile("InvalidRoles\\.NotSupported"),
			},
			{
				// Config server accounts require and honor custom roles.
				Config: testAccAccountShardingRolesConfig(name, dbAccountName, dbRolesAccountName, csAccountName, "qwer1234!!!", false, true, ""),
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
				ResourceName:            csResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
			{
				Config: testAccAccountShardingRolesConfig(name, dbAccountName, dbRolesAccountName, csAccountName, "qwer12345!!!", false, true, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(csResourceId, "roles.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(csResourceId, "roles.*", map[string]string{
						"role":    "read",
						"db_name": "admin",
					}),
				),
			},
			{
				ResourceName:            csResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
			{
				// An empty roles block is rejected at plan time.
				Config:      testAccAccountShardingRolesConfig(name, dbAccountName, dbRolesAccountName, csAccountName, "qwer12345!!!", false, true, testAccAccountRolesBlock(emptyRolesAccountName, "", "")),
				ExpectError: regexp.MustCompile(`every roles block requires a non-empty "role"`),
			},
			{
				// A role without its database is rejected at plan time.
				Config:      testAccAccountShardingRolesConfig(name, dbAccountName, dbRolesAccountName, csAccountName, "qwer12345!!!", false, true, testAccAccountRolesBlock(missingDbNameAccountName, "\n    role = \"read\"\n  ", "")),
				ExpectError: regexp.MustCompile(`every roles block requires a non-empty "db_name"`),
			},
		},
	})
}

// TestAccAliCloudMongodbAccount_complete_normal_account_roles covers replica set
// accounts with custom roles, including the description update that must carry
// the database name (C8).
func TestAccAliCloudMongodbAccount_complete_normal_account_roles(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-normal-%d", rand)
	accountName := fmt.Sprintf("tf_n_%d", rand)
	resourceId := "alicloud_mongodb_account.default"
	description := fmt.Sprintf("desc_%d", rand)
	updatedDescription := fmt.Sprintf("desc_upd_%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", "cn-shanghai")
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAccountDependenciesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountNormalRolesConfig(name, accountName, "qwer1234!!!", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "account_name", accountName),
					resource.TestCheckResourceAttr(resourceId, "character_type", "normal"),
					resource.TestCheckResourceAttr(resourceId, "db_name", "admin"),
					resource.TestCheckResourceAttr(resourceId, "roles.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "roles.*", map[string]string{
						"role":    "read",
						"db_name": "admin",
					}),
					resource.TestCheckResourceAttrSet(resourceId, "instance_id"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
			{
				Config: testAccAccountNormalRolesConfig(name, accountName, "qwer1234!!!", description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "account_description", description),
					resource.TestCheckResourceAttr(resourceId, "roles.#", "1"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
			{
				Config: testAccAccountNormalRolesConfig(name, accountName, "qwer12345!!!", updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "account_description", updatedDescription),
					resource.TestCheckResourceAttr(resourceId, "roles.#", "1"),
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

func testAccAccountShardingRolesConfig(name, dbAccount, dbRolesAccount, csAccount, password string, withDbRoles, withCs bool, extra string) string {
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
%s
`, testAccAccountShardingRolesDependence(name), dbAccount, password, name, dbRolesBlock, csBlock, extra)
}

func testAccAccountRolesBlock(accountName, roleBody, dbNameBody string) string {
	return fmt.Sprintf(`
resource "alicloud_mongodb_account" "invalid_roles" {
  account_name     = "%s"
  account_password = "qwer1234!!!"
  instance_id      = alicloud_mongodb_sharding_instance.default.id
  character_type   = "cs"
  db_name          = "admin"

  roles {%s%s}
}
`, accountName, roleBody, dbNameBody)
}

func testAccAccountShardingRolesDependence(name string) string {
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

func testAccAccountNormalRolesConfig(name, accountName, password, description string) string {
	descriptionLine := ""
	if description != "" {
		descriptionLine = fmt.Sprintf("\n  account_description = \"%s\"", description)
	}
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

resource "alicloud_mongodb_instance" "default" {
  engine_version      = "4.4"
  storage_type        = "cloud_essd1"
  vswitch_id          = alicloud_vswitch.default.id
  vpc_id              = alicloud_vpc.default.id
  storage_engine      = "WiredTiger"
  network_type        = "VPC"
  zone_id             = local.zone_id
  db_instance_storage = "20"
  db_instance_class   = "mdb.shard.4x.large.d"
  name                = var.name
}

resource "alicloud_mongodb_account" "default" {
  account_name     = "%s"
  account_password = "%s"
  instance_id      = alicloud_mongodb_instance.default.id
  character_type   = "normal"
  db_name          = "admin"%s

  roles {
    role    = "read"
    db_name = "admin"
  }
}
`, name, accountName, password, descriptionLine)
}

func testAccAccountLegacyImportId(resourceId string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceId]
		if !ok {
			return "", fmt.Errorf("resource %s not found in state", resourceId)
		}
		return fmt.Sprintf("%s:%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["account_name"]), nil
	}
}

// testAccCheckAccountDependenciesDestroyed verifies that the parent MongoDB
// instance and the VPC/vSwitch created by the case are really gone. The account
// resource itself cannot be deleted independently.
func testAccCheckAccountDependenciesDestroyed(state *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	mongodbService := mongodb.NewMongoDBService(client)
	vpcService := vpc.NewVpcService(client)

	instanceIds := make(map[string]struct{})
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "alicloud_mongodb_account" {
			continue
		}
		if instanceId := rs.Primary.Attributes["instance_id"]; instanceId != "" {
			instanceIds[instanceId] = struct{}{}
		}
	}
	for instanceId := range instanceIds {
		if _, err := mongodbService.DescribeMongoDBShardingInstance(instanceId); err == nil {
			return fmt.Errorf("parent MongoDB instance %s still exists after destroy", instanceId)
		} else if !errs.NotFoundError(err) {
			return fmt.Errorf("checking parent MongoDB instance %s: %v", instanceId, err)
		}
	}

	for _, rs := range state.RootModule().Resources {
		switch rs.Type {
		case "alicloud_vpc":
			if _, err := vpcService.DescribeVpc(rs.Primary.ID); err == nil {
				return fmt.Errorf("vpc %s still exists after destroy", rs.Primary.ID)
			} else if !errs.NotFoundError(err) {
				return fmt.Errorf("checking vpc %s: %v", rs.Primary.ID, err)
			}
		case "alicloud_vswitch":
			if _, err := vpcService.DescribeVSwitch(rs.Primary.ID); err == nil {
				return fmt.Errorf("vswitch %s still exists after destroy", rs.Primary.ID)
			} else if !errs.NotFoundError(err) {
				return fmt.Errorf("checking vswitch %s: %v", rs.Primary.ID, err)
			}
		}
	}

	return nil
}

func unitAccountResourceSchema(t *testing.T) map[string]*schema.Schema {
	t.Helper()
	return acctest.Provider.ResourcesMap["alicloud_mongodb_account"].Schema
}

func unitAccountResourceData(t *testing.T, id string, attrs map[string]interface{}, changes map[string]interface{}) *schema.ResourceData {
	t.Helper()
	state := &sdkTerraform.InstanceState{ID: id, Attributes: map[string]string{}}
	for k, v := range attrs {
		state.Attributes[k] = fmt.Sprint(v)
	}
	diff := &sdkTerraform.InstanceDiff{Attributes: map[string]*sdkTerraform.ResourceAttrDiff{}}
	for k, v := range changes {
		diff.Attributes[k] = &sdkTerraform.ResourceAttrDiff{Old: fmt.Sprint(attrs[k]), New: fmt.Sprint(v)}
	}
	d, err := schema.InternalMap(unitAccountResourceSchema(t)).Data(state, diff)
	if err != nil {
		t.Fatalf("failed to build resource data: %v", err)
	}
	return d
}

type accountApiCall struct {
	Action string
	Query  map[string]interface{}
	Body   map[string]interface{}
}

// accountUnitClient builds a client with a fake credential and a preloaded
// endpoint so the provider reaches the SDK request layer without cloud access.
func accountUnitClient(t *testing.T) *connectivity.AliyunClient {
	t.Helper()
	endpoints := &sync.Map{}
	endpoints.Store("dds", "dds.cn-shanghai.aliyuncs.com")
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("test-access-key").
		SetAccessKeySecret("test-secret-key"))
	if err != nil {
		t.Fatalf("failed to build the test credential: %v", err)
	}
	config := connectivity.Config{
		Region:               connectivity.Region("cn-shanghai"),
		RegionId:             "cn-shanghai",
		AccessKey:            "test-access-key",
		SecretKey:            "test-secret-key",
		Protocol:             "HTTPS",
		Endpoints:            endpoints,
		SignVersion:          &sync.Map{},
		ClientReadTimeout:    30000,
		ClientConnectTimeout: 60000,
		SkipRegionValidation: true,
		AccountType:          "Domestic",
		AccountId:            "1234567890123456",
		Credential:           credential,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("failed to build the test client: %v", err)
	}
	return client
}

// patchAccountDoRequest records the action, query and body the provider passes
// to the SDK, which is the last point before wire serialization.
func patchAccountDoRequest(t *testing.T, calls *[]accountApiCall, responses map[string]interface{}) *gomonkey.Patches {
	t.Helper()
	return gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, query map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		*calls = append(*calls, accountApiCall{Action: *action, Query: query, Body: body})
		if response, ok := responses[*action]; ok {
			if responseErr, ok := response.(error); ok {
				return nil, responseErr
			}
			return response.(map[string]interface{}), nil
		}
		return map[string]interface{}{}, nil
	})
}

func accountApiCallByName(calls []accountApiCall, action string) *accountApiCall {
	for i := range calls {
		if calls[i].Action == action {
			return &calls[i]
		}
	}
	return nil
}

func accountDescribeResponse(instanceId, dbName, accountName, characterType, description string, roles []map[string]string) map[string]interface{} {
	account := map[string]interface{}{
		"DBInstanceId":  instanceId,
		"AccountName":   accountName,
		"DBName":        dbName,
		"CharacterType": characterType,
		"AccountStatus": "Available",
	}
	if description != "" {
		account["AccountDescription"] = description
	}
	if len(roles) > 0 {
		list := make([]interface{}, 0, len(roles))
		for _, role := range roles {
			list = append(list, map[string]interface{}{"Role": role["role"], "DBName": role["db_name"]})
		}
		account["Roles"] = map[string]interface{}{"Role": list}
	}
	return map[string]interface{}{"Accounts": map[string]interface{}{"Account": []interface{}{account}}}
}

func TestUnitAccountSchemaContract(t *testing.T) {
	r := unitAccountResourceSchema(t)
	assert.NotNil(t, acctest.Provider.ResourcesMap["alicloud_mongodb_account"].CustomizeDiff, "roles members are validated at plan time")
	assert.True(t, r["db_name"].ForceNew, "db_name is part of the account identity")
	assert.True(t, r["db_name"].Computed)
	assert.True(t, r["roles"].Computed, "roles are returned by DescribeAccounts")
	assert.True(t, r["roles"].ForceNew)
	assert.True(t, r["character_type"].ForceNew)
	roleResource := r["roles"].Elem.(*schema.Resource)
	// The breaking-change gate rejects newly added required attributes, so the
	// members stay Optional and CustomizeDiff enforces them.
	assert.True(t, roleResource.Schema["role"].Optional)
	assert.True(t, roleResource.Schema["db_name"].Optional)
	_, enumErrs := r["character_type"].ValidateFunc("cs", "character_type")
	assert.Empty(t, enumErrs, "cs must be a valid character_type")
	_, enumErrs = r["character_type"].ValidateFunc("db", "character_type")
	assert.Empty(t, enumErrs, "db must be a valid character_type")
	_, enumErrs = r["character_type"].ValidateFunc("normal", "character_type")
	assert.Empty(t, enumErrs, "normal must be a valid character_type")
}

func TestUnitAccountDescribeRejectsInvalidId(t *testing.T) {
	service := &mongodb.MongodbServiceV2{}
	if _, err := service.DescribeMongodbAccount("dds-1"); err == nil {
		t.Fatal("expected an error for an id without an account name")
	}
	if _, err := service.DescribeMongodbAccount("dds-1:db:acc:extra"); err == nil {
		t.Fatal("expected an error for an id with too many parts")
	}
}

func TestUnitAccountStateRefreshChecksFieldError(t *testing.T) {
	service := &mongodb.MongodbServiceV2{}
	refresh := service.MongodbAccountStateRefreshFuncWithApi("dds-1:admin:acc", "AccountStatus", nil, func(id string) (map[string]interface{}, error) {
		return map[string]interface{}{"DBName": "admin"}, nil
	})
	if _, _, err := refresh(); err == nil {
		t.Fatal("expected an error when the watched field is missing from the response")
	}
}

func TestUnitAccountReadNotFoundClearsId(t *testing.T) {
	d := unitAccountResourceData(t, "dds-1:acc", map[string]interface{}{
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

func TestUnitAccountReadSetsDatabaseScopedId(t *testing.T) {
	d := unitAccountResourceData(t, "dds-1:acc1", map[string]interface{}{
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

func TestUnitAccountReadBuildsReversibleId(t *testing.T) {
	d := unitAccountResourceData(t, "dds-1:acc1", map[string]interface{}{
		"instance_id":  "dds-1",
		"account_name": "acc1",
	}, nil)
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&mongodb.MongodbServiceV2{}), "DescribeMongodbAccount", func(_ *mongodb.MongodbServiceV2, id string) (map[string]interface{}, error) {
		return map[string]interface{}{
			"DBInstanceId":  "dds-1",
			"AccountName":   "acc1",
			"DBName":        "app:logs",
			"CharacterType": "db",
			"AccountStatus": "Available",
		}, nil
	})
	defer patches.Reset()

	assert.Nil(t, mongodb.ResourceAliCloudMongodbAccountRead(d, &connectivity.AliyunClient{}))
	assert.Equal(t, "dds-1:app%3Alogs:acc1", d.Id())
}

func TestUnitAccountReadPropagatesError(t *testing.T) {
	d := unitAccountResourceData(t, "dds-1:acc", map[string]interface{}{
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

func TestUnitAccountUpdateDescriptionDescribeError(t *testing.T) {
	d := unitAccountResourceData(t, "dds-1:admin:acc", map[string]interface{}{
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

func TestUnitAccountUpdateShardDescriptionRejected(t *testing.T) {
	d := unitAccountResourceData(t, "dds-1:admin:acc", map[string]interface{}{
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

// TestUnitAccountCreateSendsRolesForReplicaSet proves that a replica set
// account with custom roles reaches CreateAccount with the roles, the database
// and the description in the request body (C1, C3-C5).
func TestUnitAccountCreateSendsRolesForReplicaSet(t *testing.T) {
	client := accountUnitClient(t)
	d := schema.TestResourceDataRaw(t, unitAccountResourceSchema(t), map[string]interface{}{
		"instance_id":         "dds-1",
		"account_name":        "normal_acc",
		"account_password":    "YourPassword+12345",
		"character_type":      "normal",
		"db_name":             "admin",
		"account_description": "normal description",
		"roles": []interface{}{
			map[string]interface{}{"role": "read", "db_name": "admin"},
		},
	})
	calls := make([]accountApiCall, 0)
	responses := map[string]interface{}{
		"CreateAccount":    map[string]interface{}{},
		"DescribeAccounts": accountDescribeResponse("dds-1", "admin", "normal_acc", "normal", "normal description", []map[string]string{{"role": "read", "db_name": "admin"}}),
	}
	patches := patchAccountDoRequest(t, &calls, responses)
	defer patches.Reset()

	assert.Nil(t, mongodb.ResourceAliCloudMongodbAccountCreate(d, client))
	call := accountApiCallByName(calls, "CreateAccount")
	if call == nil {
		t.Fatal("expected a replica set account with roles to reach CreateAccount")
	}
	assert.Equal(t, "normal", fmt.Sprint(call.Body["CharacterType"]))
	assert.Equal(t, "admin", fmt.Sprint(call.Body["DBName"]))
	assert.Equal(t, "normal description", fmt.Sprint(call.Body["AccountDescription"]))
	roles, ok := call.Body["Roles"].([]interface{})
	if !ok || len(roles) != 1 {
		t.Fatalf("expected one Role entry in the CreateAccount body, got %v", call.Body["Roles"])
	}
	role := roles[0].(map[string]interface{})
	assert.Equal(t, "read", fmt.Sprint(role["Role"]))
	assert.Equal(t, "admin", fmt.Sprint(role["DBName"]))
	assert.Equal(t, "dds-1:admin:normal_acc", d.Id())
}

// TestUnitAccountCreateDefaultsNormalForRoles covers the documented default:
// when character_type is omitted, roles are created as a replica set (normal)
// account instead of falling back to the API default db.
func TestUnitAccountCreateDefaultsNormalForRoles(t *testing.T) {
	client := accountUnitClient(t)
	d := schema.TestResourceDataRaw(t, unitAccountResourceSchema(t), map[string]interface{}{
		"instance_id":      "dds-1",
		"account_name":     "normal_def_acc",
		"account_password": "YourPassword+12345",
		"roles": []interface{}{
			map[string]interface{}{"role": "read", "db_name": "admin"},
		},
	})
	calls := make([]accountApiCall, 0)
	responses := map[string]interface{}{
		"CreateAccount":    map[string]interface{}{},
		"DescribeAccounts": accountDescribeResponse("dds-1", "admin", "normal_def_acc", "normal", "", []map[string]string{{"role": "read", "db_name": "admin"}}),
	}
	patches := patchAccountDoRequest(t, &calls, responses)
	defer patches.Reset()

	assert.Nil(t, mongodb.ResourceAliCloudMongodbAccountCreate(d, client))
	call := accountApiCallByName(calls, "CreateAccount")
	if call == nil {
		t.Fatal("expected an account with roles to reach CreateAccount")
	}
	assert.Equal(t, "normal", fmt.Sprint(call.Body["CharacterType"]))
	assert.Equal(t, "admin", fmt.Sprint(call.Body["DBName"]))
}

// TestUnitAccountCreateForwardsShardRolesToApi proves the provider no longer
// rejects a shard account with roles locally: the request reaches CreateAccount
// and the InvalidRoles.NotSupported response comes from the API.
func TestUnitAccountCreateForwardsShardRolesToApi(t *testing.T) {
	client := accountUnitClient(t)
	d := schema.TestResourceDataRaw(t, unitAccountResourceSchema(t), map[string]interface{}{
		"instance_id":      "dds-1",
		"account_name":     "shard_acc",
		"account_password": "YourPassword+12345",
		"character_type":   "db",
		"db_name":          "admin",
		"roles": []interface{}{
			map[string]interface{}{"role": "read", "db_name": "admin"},
		},
	})
	calls := make([]accountApiCall, 0)
	responses := map[string]interface{}{
		"CreateAccount": &tea.SDKError{
			Code:       helper.String("InvalidRoles.NotSupported"),
			Data:       helper.String("InvalidRoles.NotSupported"),
			Message:    helper.String("Roles is not supported for shard direct-connection account, whose privileges are fixed."),
			StatusCode: tea.Int(400),
		},
	}
	patches := patchAccountDoRequest(t, &calls, responses)
	defer patches.Reset()

	err := mongodb.ResourceAliCloudMongodbAccountCreate(d, client)
	if err == nil {
		t.Fatal("expected the CreateAccount rejection to be returned")
	}
	assert.Contains(t, err.Error(), "InvalidRoles.NotSupported")
	call := accountApiCallByName(calls, "CreateAccount")
	if call == nil {
		t.Fatal("expected the shard account with roles to reach CreateAccount")
	}
	assert.Equal(t, "db", fmt.Sprint(call.Body["CharacterType"]))
	roles, ok := call.Body["Roles"].([]interface{})
	if !ok || len(roles) != 1 {
		t.Fatalf("expected the roles to be forwarded to the API, got %v", call.Body["Roles"])
	}
}

// TestUnitAccountCreateSendsCreateAccountPayload covers the config server
// account payload: description, database, character type and roles (C1-C5).
func TestUnitAccountCreateSendsCreateAccountPayload(t *testing.T) {
	client := accountUnitClient(t)
	d := schema.TestResourceDataRaw(t, unitAccountResourceSchema(t), map[string]interface{}{
		"instance_id":         "dds-1",
		"account_name":        "cs_acc",
		"account_password":    "YourPassword+12345",
		"character_type":      "cs",
		"db_name":             "appdb",
		"account_description": "cs description",
		"roles": []interface{}{
			map[string]interface{}{"role": "readWrite", "db_name": "appdb"},
		},
	})
	calls := make([]accountApiCall, 0)
	responses := map[string]interface{}{
		"CreateAccount":    map[string]interface{}{},
		"DescribeAccounts": accountDescribeResponse("dds-1", "appdb", "cs_acc", "cs", "cs description", []map[string]string{{"role": "readWrite", "db_name": "appdb"}}),
	}
	patches := patchAccountDoRequest(t, &calls, responses)
	defer patches.Reset()

	assert.Nil(t, mongodb.ResourceAliCloudMongodbAccountCreate(d, client))
	call := accountApiCallByName(calls, "CreateAccount")
	if call == nil {
		t.Fatal("expected a config server account to reach CreateAccount")
	}
	assert.Equal(t, "cs", fmt.Sprint(call.Body["CharacterType"]))
	assert.Equal(t, "appdb", fmt.Sprint(call.Body["DBName"]))
	assert.Equal(t, "cs description", fmt.Sprint(call.Body["AccountDescription"]))
	roles, ok := call.Body["Roles"].([]interface{})
	if !ok || len(roles) != 1 {
		t.Fatalf("expected one Role entry in the CreateAccount body, got %v", call.Body["Roles"])
	}
	role := roles[0].(map[string]interface{})
	assert.Equal(t, "readWrite", fmt.Sprint(role["Role"]))
	assert.Equal(t, "appdb", fmt.Sprint(role["DBName"]))
	assert.Equal(t, "dds-1:appdb:cs_acc", d.Id())
}

// TestUnitAccountCreateRootSendsResetPasswordDbName covers C6: the default
// database is sent explicitly even when the configuration omits db_name.
func TestUnitAccountCreateRootSendsResetPasswordDbName(t *testing.T) {
	client := accountUnitClient(t)
	d := schema.TestResourceDataRaw(t, unitAccountResourceSchema(t), map[string]interface{}{
		"instance_id":      "dds-1",
		"account_name":     "root",
		"account_password": "YourPassword+12345",
	})
	calls := make([]accountApiCall, 0)
	responses := map[string]interface{}{
		"ResetAccountPassword": map[string]interface{}{},
		"DescribeAccounts":     accountDescribeResponse("dds-1", "admin", "root", "normal", "", nil),
	}
	patches := patchAccountDoRequest(t, &calls, responses)
	defer patches.Reset()

	assert.Nil(t, mongodb.ResourceAliCloudMongodbAccountCreate(d, client))
	call := accountApiCallByName(calls, "ResetAccountPassword")
	if call == nil {
		t.Fatal("expected the replica set root account to use ResetAccountPassword")
	}
	assert.Equal(t, "admin", fmt.Sprint(call.Body["DBName"]))
	_, hasRoles := call.Body["Roles"]
	assert.False(t, hasRoles, "ResetAccountPassword must not carry roles")
	assert.Equal(t, "dds-1:admin:root", d.Id())
}

// TestUnitAccountUpdatePasswordSendsDbName covers C7.
func TestUnitAccountUpdatePasswordSendsDbName(t *testing.T) {
	client := accountUnitClient(t)
	d := unitAccountResourceData(t, "dds-1:appdb:acc", map[string]interface{}{
		"instance_id":      "dds-1",
		"account_name":     "acc",
		"db_name":          "appdb",
		"character_type":   "db",
		"account_password": "YourPassword+12345",
	}, map[string]interface{}{
		"account_password": "YourPassword+123456",
	})
	calls := make([]accountApiCall, 0)
	responses := map[string]interface{}{
		"ResetAccountPassword": map[string]interface{}{},
		"DescribeAccounts":     accountDescribeResponse("dds-1", "appdb", "acc", "db", "", nil),
	}
	patches := patchAccountDoRequest(t, &calls, responses)
	defer patches.Reset()

	assert.Nil(t, mongodb.ResourceAliCloudMongodbAccountUpdate(d, client))
	call := accountApiCallByName(calls, "ResetAccountPassword")
	if call == nil {
		t.Fatal("expected the password update to call ResetAccountPassword")
	}
	assert.Equal(t, "appdb", fmt.Sprint(call.Body["DBName"]))
	assert.Equal(t, "db", fmt.Sprint(call.Body["CharacterType"]))
}

// TestUnitAccountUpdateDescriptionSendsDbName covers C8.
func TestUnitAccountUpdateDescriptionSendsDbName(t *testing.T) {
	client := accountUnitClient(t)
	d := unitAccountResourceData(t, "dds-1:appdb:acc", map[string]interface{}{
		"instance_id":         "dds-1",
		"account_name":        "acc",
		"db_name":             "appdb",
		"account_description": "before",
	}, map[string]interface{}{
		"account_description": "after",
	})
	calls := make([]accountApiCall, 0)
	responses := map[string]interface{}{
		"ModifyAccountDescription": map[string]interface{}{},
		"DescribeAccounts":         accountDescribeResponse("dds-1", "appdb", "acc", "normal", "after", nil),
	}
	patches := patchAccountDoRequest(t, &calls, responses)
	defer patches.Reset()

	assert.Nil(t, mongodb.ResourceAliCloudMongodbAccountUpdate(d, client))
	call := accountApiCallByName(calls, "ModifyAccountDescription")
	if call == nil {
		t.Fatal("expected the description update to call ModifyAccountDescription")
	}
	assert.Equal(t, "appdb", fmt.Sprint(call.Body["DBName"]))
	assert.Equal(t, "after", fmt.Sprint(call.Body["AccountDescription"]))
}

// TestUnitAccountDescribeFiltersByIdentity covers QA-PR-07: the collection read
// must return the requested account even when other accounts exist.
func TestUnitAccountDescribeFiltersByIdentity(t *testing.T) {
	client := accountUnitClient(t)
	calls := make([]accountApiCall, 0)
	responses := map[string]interface{}{
		"DescribeAccounts": map[string]interface{}{
			"Accounts": map[string]interface{}{
				"Account": []interface{}{
					map[string]interface{}{"AccountName": "shared", "DBName": "dbA", "DBInstanceId": "dds-1", "AccountStatus": "Available"},
					map[string]interface{}{"AccountName": "shared", "DBName": "dbB", "DBInstanceId": "dds-1", "AccountStatus": "Available"},
				},
			},
		},
	}
	patches := patchAccountDoRequest(t, &calls, responses)
	defer patches.Reset()
	service := mongodb.NewMongodbServiceV2(client)

	object, err := service.DescribeMongodbAccount("dds-1:dbB:shared")
	assert.NoError(t, err)
	assert.Equal(t, "dbB", fmt.Sprint(object["DBName"]), "the requested database must win over the first list entry")

	if _, err := service.DescribeMongodbAccount("dds-1:dbC:shared"); err == nil || !errs.NotFoundError(err) {
		t.Fatalf("expected NotFound when only other databases hold the account, got %v", err)
	}

	if _, err := service.DescribeMongodbAccount("dds-1:shared"); err == nil {
		t.Fatal("expected an ambiguity error for a two-part id that matches two databases")
	} else {
		assert.Contains(t, err.Error(), "db_name")
	}
}

// TestUnitAccountDescribeDecodesEscapedDatabase covers the reversible id for a
// legal database name that contains a colon (QA-PR-02).
func TestUnitAccountDescribeDecodesEscapedDatabase(t *testing.T) {
	client := accountUnitClient(t)
	calls := make([]accountApiCall, 0)
	responses := map[string]interface{}{
		"DescribeAccounts": accountDescribeResponse("dds-1", "app:logs", "acc1", "db", "", nil),
	}
	patches := patchAccountDoRequest(t, &calls, responses)
	defer patches.Reset()
	service := mongodb.NewMongodbServiceV2(client)

	object, err := service.DescribeMongodbAccount("dds-1:app%3Alogs:acc1")
	assert.NoError(t, err)
	assert.Equal(t, "acc1", fmt.Sprint(object["AccountName"]))
	call := accountApiCallByName(calls, "DescribeAccounts")
	if call == nil {
		t.Fatal("expected DescribeAccounts to be called")
	}
	assert.Equal(t, "app:logs", fmt.Sprint(call.Body["DBName"]), "the escaped database must be decoded before the request")
}

func TestUnitAccountCreateRequiresRolesForCsAccount(t *testing.T) {
	d := schema.TestResourceDataRaw(t, unitAccountResourceSchema(t), map[string]interface{}{
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

func TestUnitAccountDeleteNoop(t *testing.T) {
	d := unitAccountResourceData(t, "dds-1:acc", map[string]interface{}{}, nil)
	assert.Nil(t, mongodb.ResourceAliCloudMongodbAccountDelete(d, &connectivity.AliyunClient{}))
}

// accountWireCall records one serialized HTTP request as it left the SDK. The
// query and the urlencoded body are captured after signing, so they are the
// final wire values for the account API parameters.
type accountWireCall struct {
	Action string
	Method string
	Query  url.Values
	Form   url.Values
}

func (c accountWireCall) form() url.Values {
	if c.Form != nil {
		return c.Form
	}
	return url.Values{}
}

func accountWireCallsByName(calls []accountWireCall, action string) []accountWireCall {
	matched := make([]accountWireCall, 0, 1)
	for _, call := range calls {
		if call.Action == action {
			matched = append(matched, call)
		}
	}
	return matched
}

// accountWireClient points the DDS RPC client at a local HTTP server so the
// tests observe the serialized request instead of the pre-serialization
// DoRequest map.
func accountWireClient(t *testing.T, describe map[string]interface{}) (*connectivity.AliyunClient, *[]accountWireCall) {
	t.Helper()
	calls := make([]accountWireCall, 0)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read the serialized request body: %v", err)
		}
		form, err := url.ParseQuery(string(raw))
		if err != nil {
			t.Errorf("failed to parse the serialized request body: %v", err)
		}
		action := r.URL.Query().Get("Action")
		calls = append(calls, accountWireCall{Action: action, Method: r.Method, Query: r.URL.Query(), Form: form})
		w.Header().Set("Content-Type", "application/json")
		if action == "DescribeAccounts" {
			if err := json.NewEncoder(w).Encode(describe); err != nil {
				t.Errorf("failed to write the DescribeAccounts response: %v", err)
			}
			return
		}
		_, _ = io.WriteString(w, `{"RequestId":"mock"}`)
	}))
	t.Cleanup(server.Close)

	endpoints := &sync.Map{}
	endpoints.Store("dds", strings.TrimPrefix(server.URL, "http://"))
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("test-access-key").
		SetAccessKeySecret("test-secret-key"))
	if err != nil {
		t.Fatalf("failed to build the test credential: %v", err)
	}
	client, err := (&connectivity.Config{
		Region:               connectivity.Region("cn-shanghai"),
		RegionId:             "cn-shanghai",
		AccessKey:            "test-access-key",
		SecretKey:            "test-secret-key",
		Protocol:             "HTTP",
		Endpoints:            endpoints,
		SignVersion:          &sync.Map{},
		ClientReadTimeout:    30000,
		ClientConnectTimeout: 60000,
		SkipRegionValidation: true,
		AccountType:          "Domestic",
		AccountId:            "1234567890123456",
		Credential:           credential,
	}).Client()
	if err != nil {
		t.Fatalf("failed to build the wire test client: %v", err)
	}
	return client, &calls
}

func accountWireDescribe(instanceId, dbName, accountName, characterType, description string, roles []map[string]string) map[string]interface{} {
	account := map[string]interface{}{
		"DBInstanceId":  instanceId,
		"AccountName":   accountName,
		"DBName":        dbName,
		"AccountStatus": "Available",
		"CharacterType": characterType,
	}
	if description != "" {
		account["AccountDescription"] = description
	}
	if len(roles) > 0 {
		roleList := make([]interface{}, 0, len(roles))
		for _, role := range roles {
			roleList = append(roleList, map[string]interface{}{"Role": role["role"], "DBName": role["db_name"]})
		}
		account["Roles"] = map[string]interface{}{"Role": roleList}
	}
	return map[string]interface{}{
		"RequestId": "mock",
		"Accounts":  map[string]interface{}{"Account": []interface{}{account}},
	}
}

// TestUnitAccountWritesExposeSerializedCreateAccount asserts the final
// serialized CreateAccount request for a replica set account with roles, which
// covers C1-C5 at the wire boundary.
func TestUnitAccountWritesExposeSerializedCreateAccount(t *testing.T) {
	client, calls := accountWireClient(t, accountWireDescribe("dds-1", "admin", "normal_acc", "normal", "normal description", []map[string]string{{"role": "read", "db_name": "admin"}}))
	d := schema.TestResourceDataRaw(t, unitAccountResourceSchema(t), map[string]interface{}{
		"instance_id":         "dds-1",
		"account_name":        "normal_acc",
		"account_password":    "YourPassword+12345",
		"character_type":      "normal",
		"db_name":             "admin",
		"account_description": "normal description",
		"roles": []interface{}{
			map[string]interface{}{"role": "read", "db_name": "admin"},
		},
	})

	assert.NoError(t, mongodb.ResourceAliCloudMongodbAccountCreate(d, client))
	matched := accountWireCallsByName(*calls, "CreateAccount")
	if len(matched) != 1 {
		t.Fatalf("expected exactly one serialized CreateAccount request, got %d", len(matched))
	}
	form := matched[0].form()
	assert.Equal(t, "POST", matched[0].Method)
	assert.Equal(t, "normal_acc", form.Get("AccountName"))
	assert.Equal(t, "dds-1", form.Get("DBInstanceId"))
	assert.Equal(t, "admin", form.Get("DBName"))
	assert.Equal(t, "normal", form.Get("CharacterType"))
	assert.Equal(t, "normal description", form.Get("AccountDescription"))
	assert.Equal(t, "read", form.Get("Roles.1.Role"))
	assert.Equal(t, "admin", form.Get("Roles.1.DBName"))
	assert.Empty(t, form.Get("Roles.2.Role"), "only one role may be serialized")
}

// TestUnitAccountWritesExposeSerializedResetPassword asserts the final
// serialized ResetAccountPassword request for both the create-time root path and
// the update path, which covers C6 and C7 at the wire boundary.
func TestUnitAccountWritesExposeSerializedResetPassword(t *testing.T) {
	t.Run("create root account", func(t *testing.T) {
		client, calls := accountWireClient(t, accountWireDescribe("dds-1", "admin", "root", "normal", "", nil))
		d := schema.TestResourceDataRaw(t, unitAccountResourceSchema(t), map[string]interface{}{
			"instance_id":      "dds-1",
			"account_name":     "root",
			"account_password": "YourPassword+12345",
		})

		assert.NoError(t, mongodb.ResourceAliCloudMongodbAccountCreate(d, client))
		matched := accountWireCallsByName(*calls, "ResetAccountPassword")
		if len(matched) != 1 {
			t.Fatalf("expected exactly one serialized ResetAccountPassword request, got %d", len(matched))
		}
		form := matched[0].form()
		assert.Equal(t, "root", form.Get("AccountName"))
		assert.Equal(t, "dds-1", form.Get("DBInstanceId"))
		assert.Equal(t, "admin", form.Get("DBName"))
		assert.Empty(t, form.Get("Roles.1.Role"), "ResetAccountPassword must not serialize roles")
	})

	t.Run("update password", func(t *testing.T) {
		client, calls := accountWireClient(t, accountWireDescribe("dds-1", "appdb", "acc", "db", "", nil))
		d := unitAccountResourceData(t, "dds-1:appdb:acc", map[string]interface{}{
			"instance_id":      "dds-1",
			"account_name":     "acc",
			"db_name":          "appdb",
			"character_type":   "db",
			"account_password": "YourPassword+12345",
		}, map[string]interface{}{
			"account_password": "YourPassword+123456",
		})

		assert.NoError(t, mongodb.ResourceAliCloudMongodbAccountUpdate(d, client))
		matched := accountWireCallsByName(*calls, "ResetAccountPassword")
		if len(matched) != 1 {
			t.Fatalf("expected exactly one serialized ResetAccountPassword request, got %d", len(matched))
		}
		form := matched[0].form()
		assert.Equal(t, "acc", form.Get("AccountName"))
		assert.Equal(t, "appdb", form.Get("DBName"))
		assert.Equal(t, "db", form.Get("CharacterType"))
	})
}

// TestUnitAccountWritesExposeSerializedModifyDescription asserts the final
// serialized ModifyAccountDescription request, which covers C8 at the wire
// boundary.
func TestUnitAccountWritesExposeSerializedModifyDescription(t *testing.T) {
	client, calls := accountWireClient(t, accountWireDescribe("dds-1", "appdb", "acc", "normal", "after", nil))
	d := unitAccountResourceData(t, "dds-1:appdb:acc", map[string]interface{}{
		"instance_id":         "dds-1",
		"account_name":        "acc",
		"db_name":             "appdb",
		"account_description": "before",
	}, map[string]interface{}{
		"account_description": "after",
	})

	assert.NoError(t, mongodb.ResourceAliCloudMongodbAccountUpdate(d, client))
	matched := accountWireCallsByName(*calls, "ModifyAccountDescription")
	if len(matched) != 1 {
		t.Fatalf("expected exactly one serialized ModifyAccountDescription request, got %d", len(matched))
	}
	form := matched[0].form()
	assert.Equal(t, "acc", form.Get("AccountName"))
	assert.Equal(t, "dds-1", form.Get("DBInstanceId"))
	assert.Equal(t, "appdb", form.Get("DBName"))
	assert.Equal(t, "after", form.Get("AccountDescription"))
}
