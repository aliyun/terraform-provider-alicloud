package alicloud

import (
	"fmt"
	"testing"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudKVStoreAccount_basic(t *testing.T) {
	var v r_kvstore.Account
	resourceId := "alicloud_kvstore_account.default"
	ra := resourceAttrInit(resourceId, KvstoreAccountMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKvstoreAccounttftestnormal%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KvstoreAccountBasicdependence)
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

func TestAccAlicloudKVStoreAccount_basic_v5(t *testing.T) {
	var v r_kvstore.Account
	resourceId := "alicloud_kvstore_account.default"
	ra := resourceAttrInit(resourceId, KvstoreAccountMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKvstoreAccounttftestnormal%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KvstoreAccountBasicdependenceV5)
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

var KvstoreAccountMap = map[string]string{
	"account_privilege": "RoleReadWrite",
	"status":            CHECKSET,
}

func KvstoreAccountBasicdependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_kvstore_zones" "default"{
		instance_charge_type = "PostPaid"
	}
	variable "name" {
		default = "%v"
	}
	resource "alicloud_kvstore_instance" "default" {
		zone_id = data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id
		instance_class = "redis.master.small.default"
		instance_name  = var.name
		engine_version = "4.0"
	}
	`, name)
}

func KvstoreAccountBasicdependenceV5(name string) string {
	return fmt.Sprintf(`
	data "alicloud_kvstore_zones" "default"{
		instance_charge_type = "PostPaid"
	}
	variable "name" {
		default = "%v"
	}
	resource "alicloud_kvstore_instance" "default" {
		zone_id = data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id
		instance_class = "redis.master.small.default"
		instance_name  = var.name
		engine_version = "5.0"
	}
	`, name)
}
