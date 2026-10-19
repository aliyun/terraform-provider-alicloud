package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_kvstore_account", &resource.Sweeper{
		Name: "alicloud_kvstore_account",
		F:    testSweepKvstoreAccounts,
	})
}

// testSweepKvstoreAccounts sweeps KVStore (Redis) account sub-resources.
// Accounts are a sub-resource of a KVStore instance, so we enumerate instances
// first, then DescribeAccounts per instance and delete each non-system account.
func testSweepKvstoreAccounts(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	// 1. List all KVStore instances (Redis + Memcache) so we can enumerate accounts per instance.
	listAction := "DescribeInstances"
	listReq := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeXLarge,
		"PageNumber": 1,
	}
	instanceIds := make([]string, 0)
	for _, instanceType := range []string{string(KVStoreRedis), string(KVStoreMemcache)} {
		listReq["InstanceType"] = instanceType
		for {
			response, err := client.RpcPost("R-kvstore", "2015-01-01", listAction, nil, listReq, true)
			if err != nil {
				log.Printf("[ERROR] Failed to list KVStore instances for account sweep: %s", err)
				return nil
			}
			resp, err := jsonpath.Get("$.Instances.KVStoreInstance", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, listAction, "$.Instances.KVStoreInstance", response)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				instanceName := fmt.Sprint(item["InstanceName"])
				if protectedSweepVpcName(instanceName) {
					log.Printf("[INFO] Skipping KVStore Instance accounts (NODELETING): %s (%s)", instanceName, item["InstanceId"])
					continue
				}
				if !sweepAll() {
					match := false
					for _, prefix := range prefixes {
						if strings.HasPrefix(strings.ToLower(instanceName), strings.ToLower(prefix)) {
							match = true
							break
						}
					}
					if !match {
						continue
					}
				}
				instanceIds = append(instanceIds, fmt.Sprint(item["InstanceId"]))
			}
			if len(result) < PageSizeXLarge {
				break
			}
			listReq["PageNumber"] = listReq["PageNumber"].(int) + 1
		}
	}

	// 2. For each qualifying instance, list accounts and delete each (skip system/built-in accounts).
	for _, instanceId := range instanceIds {
		accReq := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": instanceId,
		}
		accResp, err := client.RpcPost("R-kvstore", "2015-01-01", "DescribeAccounts", nil, accReq, true)
		if err != nil {
			log.Printf("[ERROR] Failed to DescribeAccounts for instance %s: %s", instanceId, err)
			continue
		}
		accsObj, err := jsonpath.Get("$.Accounts.Account", accResp)
		if err != nil {
			log.Printf("[ERROR] Failed to parse accounts for instance %s: %s", instanceId, err)
			continue
		}
		accList, _ := accsObj.([]interface{})
		for _, a := range accList {
			acc, ok := a.(map[string]interface{})
			if !ok {
				continue
			}
			accName := fmt.Sprint(acc["AccountName"])
			// Skip built-in/system accounts — the cloud does not allow deleting them.
			if accType := fmt.Sprint(acc["AccountType"]); accType == "System" || accType == "Super" {
				log.Printf("[INFO] Skipping system KVStore account on instance %s: %s (%s)", instanceId, accName, accType)
				continue
			}
			if protectedSweepVpcName(accName) {
				log.Printf("[INFO] Skipping KVStore account (NODELETING): %s:%s", instanceId, accName)
				continue
			}
			log.Printf("[INFO] Deleting KVStore Account: %s:%s", instanceId, accName)
			delReq := map[string]interface{}{
				"RegionId":    client.RegionId,
				"InstanceId":  instanceId,
				"AccountName": accName,
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(1*time.Minute, func() *resource.RetryError {
				_, err = client.RpcPost("R-kvstore", "2015-01-01", "DeleteAccount", nil, delReq, false)
				if err != nil {
					if IsExpectedErrors(err, []string{"InvalidAccountName.NotFound", "InvalidInstanceId.NotFound"}) {
						return nil
					}
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete KVStore Account %s:%s: %s", instanceId, accName, err)
			}
		}
	}
	return nil
}

func TestAccAliCloudKvstoreAccount_basic(t *testing.T) {
	var v r_kvstore.Account
	resourceId := "alicloud_kvstore_account.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreAccountMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKvstoreAccounttftestnormal%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKVStoreAccountBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.KvStoreSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_name":     "tftest",
					"account_password": "YourPassword_123",
					"instance_id":      "${alicloud_kvstore_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name":     "tftest",
						"account_password": "YourPassword_123",
						"instance_id":      CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
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
					"account_privilege": "RoleReadOnly",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_privilege": "RoleReadOnly",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform_test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password":  "YourPassword_123",
					"account_privilege": "RoleReadWrite",
					"description":       "terraform_test_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password":  "YourPassword_123",
						"account_privilege": "RoleReadWrite",
						"description":       "terraform_test_update",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudKvstoreAccount_kms_encrypted_password(t *testing.T) {
	var v r_kvstore.Account
	resourceId := "alicloud_kvstore_account.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreAccountMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKvstoreAccounttftestnormal%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKVStoreAccountBasicDependence_kms_encrypted_password)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.KvStoreSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_name":           "tftest",
					"kms_encrypted_password": "${alicloud_kms_ciphertext.default.ciphertext_blob}",
					"kms_encryption_context": map[string]string{
						"name": name,
					},
					"instance_id": "${alicloud_kvstore_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name":             "tftest",
						"instance_id":              CHECKSET,
						"kms_encrypted_password":   CHECKSET,
						"kms_encryption_context.%": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kms_encrypted_password": "${alicloud_kms_ciphertext.update.ciphertext_blob}",
					"kms_encryption_context": map[string]string{
						"name": name + "update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kms_encrypted_password":   CHECKSET,
						"kms_encryption_context.%": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kms_encrypted_password", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAliCloudKvstoreAccount_basic_v5(t *testing.T) {
	var v r_kvstore.Account
	resourceId := "alicloud_kvstore_account.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreAccountMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKvstoreAccounttftestnormal%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKVStoreAccountBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.KvStoreSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_name":     "tftest",
					"account_password": "YourPassword_123",
					"instance_id":      "${alicloud_kvstore_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name":     "tftest",
						"account_password": "YourPassword_123",
						"instance_id":      CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
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
					"account_privilege": "RoleReadOnly",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_privilege": "RoleReadOnly",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform_test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password":  "YourPassword_123",
					"account_privilege": "RoleReadWrite",
					"description":       "terraform_test_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password":  "YourPassword_123",
						"account_privilege": "RoleReadWrite",
						"description":       "terraform_test_update",
					}),
				),
			},
		},
	})
}

var AliCloudKVStoreAccountMap0 = map[string]string{
	"account_privilege": "RoleReadWrite",
	"status":            CHECKSET,
}

func AliCloudKVStoreAccountBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_kvstore_zones" "default" {
  		instance_charge_type = "PostPaid"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		zone_id = data.alicloud_kvstore_zones.default.zones.0.id
  		vpc_id  = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_kvstore_instance" "default" {
  		zone_id          = data.alicloud_kvstore_zones.default.zones.0.id
  		instance_class   = "redis.master.small.default"
  		db_instance_name = var.name
  		engine_version   = "5.0"
  		vswitch_id       = data.alicloud_vswitches.default.ids.0
	}
`, name)
}

func AliCloudKVStoreAccountBasicDependence_kms_encrypted_password(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_kms_instances" "default" {
	}

	resource "alicloud_kms_key" "default" {
  		description            = var.name
  		status                 = "Enabled"
  		pending_window_in_days = 7
        dkms_instance_id       = data.alicloud_kms_instances.default.instances.0.instance_id
	}

	resource "alicloud_kms_ciphertext" "default" {
  		key_id    = alicloud_kms_key.default.id
  		plaintext = "YourPassword1234!"
  		encryption_context = {
    		"name" = var.name
  		}
	}

	resource "alicloud_kms_ciphertext" "update" {
  		key_id    = alicloud_kms_key.default.id
  		plaintext = "YourPassword1234!update"
  		encryption_context = {
    		"name" = "${var.name}update"
  		}
	}

	data "alicloud_kvstore_zones" "default" {
  		instance_charge_type = "PostPaid"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		zone_id = data.alicloud_kvstore_zones.default.zones.0.id
  		vpc_id  = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_kvstore_instance" "default" {
  		zone_id          = data.alicloud_kvstore_zones.default.zones.0.id
  		instance_class   = "redis.master.small.default"
  		db_instance_name = var.name
  		engine_version   = "5.0"
  		vswitch_id       = data.alicloud_vswitches.default.ids.0
	}
`, name)
}

func AliCloudKVStoreAccountBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_kvstore_zones" "default" {
  		instance_charge_type = "PostPaid"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		zone_id = data.alicloud_kvstore_zones.default.zones.0.id
  		vpc_id  = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_kvstore_instance" "default" {
  		zone_id          = data.alicloud_kvstore_zones.default.zones.0.id
  		instance_class   = "redis.master.small.default"
  		db_instance_name = var.name
  		engine_version   = "5.0"
  		vswitch_id       = data.alicloud_vswitches.default.ids.0
	}
`, name)
}

// Test Redis Account. >>> Resource test cases, automatically generated.
// Case Redis Account资源测试 11424
func TestAccAliCloudRedisAccount_basic11424(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kvstore_account.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisAccountMap11424)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccredis%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisAccountBasicDependence11424)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "tf测试 Account",
					"instance_id":       "${alicloud_kvstore_instance.defaultRedisInstance.id}",
					"account_privilege": "RoleReadWrite",
					"account_type":      "Normal",
					"account_name":      name,
					"account_password":  "123456_tf",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "tf测试 Account",
						"instance_id":       CHECKSET,
						"account_privilege": "RoleReadWrite",
						"account_type":      "Normal",
						"account_name":      name,
						"account_password":  "123456_tf",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "tf测试 Account 修改",
					"account_privilege": "RoleReadOnly",
					"account_password":  "tf_123456",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "tf测试 Account 修改",
						"account_privilege": "RoleReadOnly",
						"account_password":  "tf_123456",
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

var AlicloudRedisAccountMap11424 = map[string]string{
	"status": CHECKSET,
}

func AlicloudRedisAccountBasicDependence11424(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_kvstore_zones" "default" {
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
	zone_id = data.alicloud_kvstore_zones.default.zones.0.id
	vpc_id  = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_kvstore_instance" "defaultRedisInstance" {
	zone_id          = data.alicloud_kvstore_zones.default.zones.0.id
	instance_class   = "redis.master.small.default"
	db_instance_name = var.name
	engine_version   = "5.0"
	vswitch_id       = data.alicloud_vswitches.default.ids.0
}

`, name)
}

// Test Redis Account. <<< Resource test cases, automatically generated.

func TestAccAliCloudKvstoreAccount_passwordWo(t *testing.T) {
	var v r_kvstore.Account
	resourceId := "alicloud_kvstore_account.default"
	ra := resourceAttrInit(resourceId, AliCloudKVStoreAccountMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKvstoreAccounttftestnormal%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKVStoreAccountBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.KvStoreSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_name":                "tftest",
					"account_password_wo":         "YourPassword_123",
					"account_password_wo_version": 1,
					"instance_id":                 "${alicloud_kvstore_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name":                "tftest",
						"account_password_wo_version": "1",
						"instance_id":                 CHECKSET,
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
				ImportStateVerifyIgnore: []string{"account_password", "account_password_wo_version"},
			},
		},
	})
}
