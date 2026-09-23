package alicloud

import (
	"fmt"
	"testing"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudKVStoreAccount_basic(t *testing.T) {
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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

func TestAccAliCloudKVStoreAccount_kms_encrypted_password(t *testing.T) {
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"password", "kms_encrypted_password", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAliCloudKVStoreAccount_basic_v5(t *testing.T) {
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
