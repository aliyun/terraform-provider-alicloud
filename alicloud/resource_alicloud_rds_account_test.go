package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_rds_account", &resource.Sweeper{
		Name: "alicloud_rds_account",
		F:    testSweepRdsAccounts,
	})
}

// listSweepRdsDbInstances enumerates RDS DB instances in the region. It is a
// shared helper for RDS sub-resource sweepers (accounts, backups, ...) that
// must walk the parent DB instance list before listing their child resources.
func listSweepRdsDbInstances(client *connectivity.AliyunClient) ([]map[string]interface{}, error) {
	action := "DescribeDBInstances"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	objects := make([]map[string]interface{}, 0)
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			resp, e := client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
			response = resp
			if e != nil {
				if NeedRetry(e) {
					wait()
					return retry.RetryableError(e)
				}
				return retry.NonRetryableError(e)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return objects, err
		}
		resp, e := jsonpath.Get("$.Items.DBInstance", response)
		if e != nil {
			return objects, WrapErrorf(e, FailedGetAttributeMsg, action, "$.Items.DBInstance", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			objects = append(objects, v.(map[string]interface{}))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return objects, nil
}

// testSweepRdsAccounts sweeps leftover test RDS accounts. RDS accounts are
// sub-resources that live on a DB instance, so the sweeper enumerates DB
// instances first, then the accounts on each. Accounts whose name (or parent
// instance name) carries a do-not-delete marker are preserved. Accounts on
// blocked instances fail deletion and are recorded as leftover (non-fatal).
func testSweepRdsAccounts(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{"tf-testAcc", "tf_testAcc"}

	instances, err := listSweepRdsDbInstances(client)
	if err != nil {
		log.Printf("[ERROR] DescribeDBInstances got an error: %v", err)
		return nil
	}

	for _, inst := range instances {
		instanceId := fmt.Sprint(inst["DBInstanceId"])
		instanceName := fmt.Sprint(inst["DBInstanceDescription"])
		if protectedSweepVpcName(instanceName) {
			log.Printf("[INFO] Skipping RDS Accounts on do-not-delete instance: %s (%s)", instanceName, instanceId)
			continue
		}

		action := "DescribeAccounts"
		req := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": instanceId,
			"SourceIp":     client.SourceIp,
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		if e := retry.Retry(5*time.Minute, func() *retry.RetryError {
			resp, e := client.RpcPost("Rds", "2014-08-15", action, nil, req, true)
			response = resp
			if e != nil {
				if rdsErrorHasCode(e, "InvalidDBInstanceId.NotFound", "InvalidDBInstanceName.NotFound") || NotFoundError(e) {
					return nil
				}
				if NeedRetry(e) || IsExpectedErrors(e, dbInstanceStatusErrorCodes) {
					wait()
					return retry.RetryableError(e)
				}
				return retry.NonRetryableError(e)
			}
			return nil
		}); e != nil {
			log.Printf("[ERROR] DescribeAccounts for %s got an error: %v", instanceId, e)
			continue
		}
		addDebug(action, response, req)

		accResp, e := jsonpath.Get("$.Accounts.DBInstanceAccount", response)
		if e != nil {
			log.Printf("[ERROR] DescribeAccounts parse for %s: %v", instanceId, e)
			continue
		}
		accounts, _ := accResp.([]interface{})
		for _, a := range accounts {
			acc := a.(map[string]interface{})
			accName := fmt.Sprint(acc["AccountName"])
			if protectedSweepVpcName(accName) {
				log.Printf("[INFO] Skipping RDS Account (do-not-delete marker): %s/%s", instanceId, accName)
				continue
			}
			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(accName), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
			} else {
				skip = false
			}
			if skip {
				log.Printf("[INFO] Skipping RDS Account: %s/%s", instanceId, accName)
				continue
			}

			log.Printf("[INFO] Deleting RDS Account: %s/%s", instanceId, accName)
			delAction := "DeleteAccount"
			delReq := map[string]interface{}{
				"DBInstanceId": instanceId,
				"AccountName":  accName,
				"SourceIp":     client.SourceIp,
			}
			if e := retry.Retry(2*time.Minute, func() *retry.RetryError {
				resp, e := client.RpcPost("Rds", "2014-08-15", delAction, nil, delReq, false)
				if e != nil {
					if rdsErrorHasCode(e, "InvalidDBInstanceId.NotFound", "InvalidDBInstanceName.NotFound", "InvalidAccountName.NotFound") {
						return nil
					}
					if NeedRetry(e) || IsExpectedErrors(e, dbInstanceStatusErrorCodes) || IsExpectedErrors(e, []string{"InternalError", "AccountActionForbidden", "IncorrectDBInstanceState"}) {
						return retry.RetryableError(e)
					}
					return retry.NonRetryableError(e)
				}
				addDebug(delAction, resp, delReq)
				return nil
			}); e != nil {
				log.Printf("[ERROR] Failed to delete RDS Account %s/%s: %s", instanceId, accName, e)
			}
		}
	}

	return nil
}

func TestAccAliCloudRdsAccount_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_account.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsAccountMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdsaccount%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsAccountBasicDependenceBasic)
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
					"db_instance_id":   "${alicloud_db_instance.default.id}",
					"account_name":     "tftestnormal999",
					"account_password": "YourPassword_123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":   CHECKSET,
						"account_name":     "tftestnormal999",
						"account_password": "YourPassword_123",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "reset_permission_flag"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "测试账号A",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": "测试账号A",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "测试账号",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": "测试账号",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword_1234",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "YourPassword_1234",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "tf test",
					"account_password":    "YourPassword_123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": "tf test",
						"account_password":    "YourPassword_123",
					}),
				),
			},
		},
	})
}

func AlicloudRdsAccountBasicDependenceBasic(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "local_ssd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "local_ssd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_security_group" "default" {
	security_group_name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_db_instance" "default" {
    engine = "MySQL"
	engine_version = "8.0"
 	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = local.vswitch_id
	instance_name = var.name
	instance_charge_type = "Postpaid"
	monitoring_period =     60
	db_instance_storage_type =  "local_ssd"
	db_is_ignore_case =  false
}
`, name)

}

func TestAccAliCloudRdsAccount_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_account.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsAccountMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-basic-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsAccountBasicDependenceNormal)
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
					"db_instance_id":         "${alicloud_db_instance.default.id}",
					"name":                   "tftestnormal999",
					"description":            "测试账号A",
					"type":                   "Normal",
					"kms_encrypted_password": "${alicloud_kms_ciphertext.default.ciphertext_blob}",
					"kms_encryption_context": map[string]string{
						"name": name,
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": CHECKSET,
						"name":           "tftestnormal999",
						"description":    "测试账号A",
						"type":           "Normal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kms_encrypted_password", "kms_encryption_context", "reset_permission_flag"},
			},
		},
	})
}

func TestAccAliCloudRdsAccount_normal(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_account.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdsaccount%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsAccountBasicDependenceNormal)
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
					"db_instance_id":      "${alicloud_db_instance.default.id}",
					"account_name":        "tf_test_normal",
					"account_password":    "!Q2w3e4r",
					"account_description": "TF测试普通账号",
					"account_type":        "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":      CHECKSET,
						"account_name":        "tf_test_normal",
						"account_password":    "!Q2w3e4r",
						"account_description": "TF测试普通账号",
						"account_type":        "Normal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password":      "!Q2w3e4r5t",
					"reset_permission_flag": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password":      "!Q2w3e4r5t",
						"reset_permission_flag": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password":      "!Q2w3e4r",
					"reset_permission_flag": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password":      "!Q2w3e4r",
						"reset_permission_flag": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "reset_permission_flag"},
			},
		},
	})
}

func AlicloudRdsAccountBasicDependenceNormal(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "local_ssd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "local_ssd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_security_group" "default" {
	security_group_name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_db_instance" "default" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "local_ssd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = local.vswitch_id
	instance_name = var.name
	instance_charge_type = "Postpaid"
	monitoring_period =     60
	db_is_ignore_case =  false
}

data "alicloud_kms_keys" "default" {
	  status = "Enabled"
	}
	resource "alicloud_kms_key" "default" {
	  count = length(data.alicloud_kms_keys.default.ids) > 0 ? 0 : 1
	  description = var.name
	  status = "Enabled"
	  pending_window_in_days = 7
	}
	
	resource "alicloud_kms_ciphertext" "default" {
	  key_id = length(data.alicloud_kms_keys.default.ids) > 0 ? data.alicloud_kms_keys.default.ids.0 : concat(alicloud_kms_key.default.*.id, [""])[0]
	  plaintext = "YourPassword1234"
	  encryption_context = {
		"name" = var.name
	  }
	}
`, name)

}

func TestAccAliCloudRdsAccount_super(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_account.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdsaccount%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsAccountBasicDependenceNormal)
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
					"db_instance_id":      "${alicloud_db_instance.default.id}",
					"account_name":        "tf_test_super",
					"account_password":    "!Q2w3e4r",
					"account_description": "TF测试高权限账号",
					"account_type":        "Super",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":      CHECKSET,
						"account_name":        "tf_test_super",
						"account_password":    "!Q2w3e4r",
						"account_description": "TF测试高权限账号",
						"account_type":        "Super",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password":      "!Q2w3e4r5t",
					"reset_permission_flag": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password":      "!Q2w3e4r5t",
						"reset_permission_flag": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password":      "!Q2w3e4r",
					"reset_permission_flag": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password":      "!Q2w3e4r",
						"reset_permission_flag": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "reset_permission_flag"},
			},
		},
	})
}

var AlicloudRdsAccountMap0 = map[string]string{
	"account_type": "Normal",
	"status":       "Available",
}

// Test Rds Account. >>> Resource test cases, automatically generated.
// Case RDS_ACCOUNT_TEST_SQLSERVICE 11761
func TestAccAliCloudRdsAccount_basic11761(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_account.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsAccountMap11761)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsAccountBasicDependence11761)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":      "${alicloud_db_instance.default.id}",
					"account_name":        name,
					"account_password":    "1qaz@4321",
					"account_description": "test001",
					"account_type":        "Super",
					"check_policy":        "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":      CHECKSET,
						"account_name":        name,
						"account_password":    "1qaz@4321",
						"account_description": "test001",
						"account_type":        "Super",
						"check_policy":        "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"check_policy": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"check_policy": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"check_policy": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"check_policy": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "1qaz@4312",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "1qaz@4312",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "reset_permission_flag"},
			},
		},
	})
}

// Case RDS_ACCOUNT_TEST_SQLSERVICE
func TestAccAliCloudRdsAccount_basic11761_2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_account.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsAccountMap11761)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsAccountBasicDependence11761)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":      "${alicloud_db_instance.default.id}",
					"account_name":        name,
					"account_password":    "1qaz@4321",
					"account_description": "test001",
					"account_type":        "Sysadmin",
					"check_policy":        "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":      CHECKSET,
						"account_name":        name,
						"account_password":    "1qaz@4321",
						"account_description": "test001",
						"account_type":        "Sysadmin",
						"check_policy":        "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "1qaz@4312",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "1qaz@4312",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "reset_permission_flag"},
			},
		},
	})
}

var AlicloudRdsAccountMap11761 = map[string]string{}

func AlicloudRdsAccountBasicDependence11761(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_db_zones" "default"{
  engine = "SQLServer"
  engine_version = "2012_std_ha"
  instance_charge_type = "PostPaid"
  category = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id = data.alicloud_db_zones.default.zones.0.id
  engine = "SQLServer"
  engine_version = "2012_std_ha"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type = "PostPaid"
  category = "HighAvailability"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

resource "alicloud_db_instance" "default" {
  engine         = "SQLServer"
  engine_version = "2012_std_ha"
  vswitch_id     = local.vswitch_id
  instance_type  = "mssql.x4.medium.s2"
  instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
  db_instance_storage_type = "cloud_essd"
  instance_charge_type =  "Postpaid"
  monitoring_period = "60"
  category = "HighAvailability"
  instance_name  = var.name
}

`, name)
}

// Case RDS_ACCOUNT_TEST_PGSQL 11752
func TestAccAliCloudRdsAccount_basic11752(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_account.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsAccountMap11752)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsAccountBasicDependence11752)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":      "${alicloud_db_instance.default.id}",
					"account_name":        name,
					"account_password":    "1qaz@4321",
					"account_description": "test001",
					"status":              "Available",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":      CHECKSET,
						"account_name":        name,
						"account_password":    "1qaz@4321",
						"account_description": "test001",
						"status":              "Available",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "test002",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": "test002",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "1qaz@432111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "1qaz@432111",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Lock",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Lock",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Available",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Available",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "reset_permission_flag"},
			},
		},
	})
}

var AlicloudRdsAccountMap11752 = map[string]string{}

func AlicloudRdsAccountBasicDependence11752(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_db_zones" "default"{
  	engine               = "PostgreSQL"
  	engine_version       = "12.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
  	engine               = "PostgreSQL"
  	engine_version       = "12.0"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

// RDS PostgreSQL on ECS requires this service-linked role before CreateDBInstance
resource "alicloud_rds_service_linked_role" "default" {
	service_name = "AliyunServiceRoleForRdsPgsqlOnEcs"
}

resource "alicloud_db_instance" "default" {
  engine         	= "PostgreSQL"
  engine_version 	= "12.0"
  instance_type 	=  data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage	=  data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
  db_instance_storage_type =  "cloud_essd"
  zone_id			=      data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.0.id
  instance_charge_type  =  "Postpaid"
  instance_name			=  var.name
  vswitch_id			=  local.vswitch_id
  monitoring_period 	=  "60"
  category				=  "HighAvailability"
  target_minor_version	=  "rds_postgres_1200_20231030"
	depends_on = [alicloud_rds_service_linked_role.default]
}


`, name)
}

// Test Rds Account. <<< Resource test cases, automatically generated.

func TestAccAliCloudRdsAccount_passwordWo(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_account.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsAccountMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdsaccount%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsAccountBasicDependenceBasic)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id": "${alicloud_db_instance.default.id}",
					"account_name":   "tftestnormal999",
				}),
				ExpectError: regexp.MustCompile("one\\s+of\\s+`account_password,account_password_wo,kms_encrypted_password,password`\\s+must\\s+be\\s+specified"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":              "${alicloud_db_instance.default.id}",
					"account_name":                "tftestnormal999",
					"account_password_wo":         "YourPassword_123",
					"account_password_wo_version": 1,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":              CHECKSET,
						"account_name":                "tftestnormal999",
						"account_password_wo_version": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password_wo":         "YourPassword_1234",
					"account_password_wo_version": 2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password_wo_version": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password":            "YourPassword_123",
					"account_password_wo":         REMOVEKEY,
					"account_password_wo_version": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password":            "YourPassword_123",
						"account_password_wo_version": REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password", "account_password_wo_version", "reset_permission_flag"},
			},
		},
	})
}

// rdsAccountTestClient stands up a credential-free AliyunClient whose RDS
// transport is pointed at an in-process httptest server, mirroring the
// ecs-snapshot unit-test construction. This avoids gomonkey binary patching
// (which the macOS arm64 kernel refuses) and runs under a plain `go test`.
func rdsAccountTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("rds", endpoint)
	return client
}

func writeGone403(t *testing.T, w http.ResponseWriter) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"Code":    "OperationDenied.DBInstanceStatus",
		"Message": "the request was not permitted in the current state",
	}); err != nil {
		t.Error(err)
	}
}

func rdsAccountResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()
	p := Provider().ResourcesMap
	d, err := schema.InternalMap(p["alicloud_rds_account"].Schema).Data(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	d.SetId("rm-test:test-account")
	return d
}

func TestUnitRdsAccountReadParentGone(t *testing.T) {
	client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "DescribeDBInstanceAttribute":
			writeRdsError(t, w, 404, "InvalidDBInstanceId.NotFound")
			return
		case "DescribeAccounts":
			writeGone403(t, w)
			return
		}
		t.Errorf("unexpected RDS action in Read: %s", r.Form.Get("Action"))
		w.WriteHeader(http.StatusBadRequest)
	})
	d := rdsAccountResourceData(t)

	err := resourceAliCloudRdsAccountRead(d, client)
	assert.Nil(t, err, "Read must not hard-fail when the parent instance is gone")
	assert.Equal(t, "", d.Id(), "account state must be cleared when the parent instance is gone")
}

func TestUnitRdsAccountDeleteParentGone(t *testing.T) {
	var deleteAccountCalls int32
	client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "DescribeDBInstanceAttribute":
			writeRdsError(t, w, 404, "InvalidDBInstanceId.NotFound")
			return
		case "DeleteAccount":
			atomic.AddInt32(&deleteAccountCalls, 1)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{})
			return
		}
		t.Errorf("unexpected RDS action in Delete: %s", r.Form.Get("Action"))
		w.WriteHeader(http.StatusBadRequest)
	})
	d := rdsAccountResourceData(t)

	err := resourceAliCloudRdsAccountDelete(d, client)
	assert.Nil(t, err, "Delete must be idempotent when the parent instance is gone")
	assert.Equal(t, int32(0), atomic.LoadInt32(&deleteAccountCalls),
		"DeleteAccount must not be called when the parent instance is already gone")
}

func TestUnitRdsAccountDeleteDeleteAccountTerminal403(t *testing.T) {
	var describeDBInstanceCalls int32
	client := rdsAccountTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "DescribeDBInstanceAttribute":
			n := atomic.AddInt32(&describeDBInstanceCalls, 1)
			if n == 1 {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"Items": map[string]interface{}{
						"DBInstanceAttribute": []map[string]interface{}{
							{"DBInstanceStatus": "Running", "DBInstanceId": "rm-test"},
						},
					},
				})
				return
			}
			writeRdsError(t, w, 404, "InvalidDBInstanceId.NotFound")
			return
		case "DeleteAccount":
			writeGone403(t, w)
			return
		case "DescribeAccounts":
			writeGone403(t, w)
			return
		}
		t.Errorf("unexpected RDS action in Delete: %s", r.Form.Get("Action"))
		w.WriteHeader(http.StatusBadRequest)
	})
	d := rdsAccountResourceData(t)

	err := resourceAliCloudRdsAccountDelete(d, client)
	assert.Nil(t, err, "Delete must finish idempotently when DeleteAccount returns the terminal 403")
	assert.Equal(t, int32(2), atomic.LoadInt32(&describeDBInstanceCalls),
		"the follow-up DescribeDBInstance must run exactly once after the DeleteAccount 403 (pre-check + follow-up)")
}
