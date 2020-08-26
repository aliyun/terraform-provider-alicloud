package alicloud

import (
	"fmt"
	"testing"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudKVStoreAccountUpdateV4(t *testing.T) {
	var v *r_kvstore.Account
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccKVstoreAccount-%d", rand)
	var basicMap = map[string]string{
		"instance_id":      CHECKSET,
		"account_name":     "tftestnormal",
		"account_password": "YourPassword_123",
		"account_type":     "Normal",
	}
	resourceId := "alicloud_kvstore_account.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeKVstoreAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceKVstoreAccountConfigDependenceV4)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":      "${alicloud_kvstore_instance.instance.id}",
					"account_name":     "tftestnormal",
					"account_password": "YourPassword_123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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
					"description": "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "from terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_privilege": "RoleRepl",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_privilege": "RoleRepl",
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
					"description":       "tf test",
					"account_password":  "YourPassword_123",
					"account_privilege": "RoleReadOnly",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "tf test",
						"account_password":  "YourPassword_123",
						"account_privilege": "RoleReadOnly",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudKVStoreAccountUpdateV5(t *testing.T) {
	var v *r_kvstore.Account
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccKVstoreAccount-%d", rand)
	var basicMap = map[string]string{
		"instance_id":      CHECKSET,
		"account_name":     "tftestnormal",
		"account_password": "YourPassword_123",
		"account_type":     "Normal",
	}
	resourceId := "alicloud_kvstore_account.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeKVstoreAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceKVstoreAccountConfigDependenceV5)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":      "${alicloud_kvstore_instance.instance.id}",
					"account_name":     "tftestnormal",
					"account_password": "YourPassword_123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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
					"description": "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "from terraform",
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
					"description":       "tf test",
					"account_password":  "YourPassword_123",
					"account_privilege": "RoleReadOnly",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "tf test",
						"account_password":  "YourPassword_123",
						"account_privilege": "RoleReadOnly",
					}),
				),
			},
		},
	})
}

func resourceKVstoreAccountConfigDependenceV4(name string) string {
	return fmt.Sprintf(`
	data "alicloud_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "%v"
	}
	resource "alicloud_kvstore_instance" "instance" {
		availability_zone = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%%length(data.alicloud_zones.default.zones)], "id")}"
		instance_class = "redis.master.small.default"
		instance_name  = "${var.name}"
		instance_charge_type = "PostPaid"
		engine_version = "4.0"
	}
	`, name)
}

func resourceKVstoreAccountConfigDependenceV5(name string) string {
	return fmt.Sprintf(`
	data "alicloud_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "%v"
	}
	resource "alicloud_kvstore_instance" "instance" {
		availability_zone = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%%length(data.alicloud_zones.default.zones)], "id")}"
		instance_class = "redis.master.small.default"
		instance_name  = "${var.name}"
		instance_charge_type = "PostPaid"
		engine_version = "5.0"
	}
	`, name)
}
