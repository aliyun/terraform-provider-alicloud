package rkvstore_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudKvstoreAccountsDataSource(t *testing.T) {
	resourceId := "data.alicloud_kvstore_accounts.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreAccount-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceKvstoreAccountsDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":  "${alicloud_kvstore_account.default.account_name}",
			"instance_id": "${alicloud_kvstore_instance.default.id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":  "${alicloud_kvstore_account.default.account_name}_fake",
			"instance_id": "${alicloud_kvstore_instance.default.id}",
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id":  "${alicloud_kvstore_instance.default.id}",
			"account_name": "${alicloud_kvstore_account.default.account_name}",
			"status":       "Available",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":  "${alicloud_kvstore_instance.default.id}",
			"account_name": "${alicloud_kvstore_account.default.account_name}",
			"status":       "Unavailable",
		}),
	}

	var existKvstoreAccountMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"ids.0":                   acctest.CHECKSET,
			"names.#":                 "1",
			"names.0":                 "tftest",
			"accounts.#":              "1",
			"accounts.0.id":           acctest.CHECKSET,
			"accounts.0.account_name": "tftest",
			"accounts.0.account_type": acctest.CHECKSET,
			"accounts.0.instance_id":  acctest.CHECKSET,
			"accounts.0.status":       "Available",
		}
	}

	var fakeKvstoreAccountMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"accounts.#": "0",
		}
	}

	var kvstoreAccountsInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existKvstoreAccountMapFunc,
		FakeMapFunc:  fakeKvstoreAccountMapFunc,
	}
	kvstoreAccountsInfo.DataSourceTestCheck(t, rand, nameRegexConf, statusConf)
}

func dataSourceKvstoreAccountsDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%v"
	}
	variable "zone_id" {
		default = "cn-beijing-k"
	}
	resource "alicloud_vpc" "default" {
		vpc_name   = var.name
		cidr_block = "172.16.0.0/16"
	}
	resource "alicloud_vswitch" "default" {
		vpc_id       = alicloud_vpc.default.id
		cidr_block   = "172.16.0.0/24"
		zone_id      = var.zone_id
		vswitch_name = var.name
	}
	resource "alicloud_kvstore_instance" "default" {
		zone_id           = var.zone_id
		instance_class    = "redis.master.small.default"
		instance_name     = var.name
		engine_version    = "5.0"
		vswitch_id        = alicloud_vswitch.default.id
	}
	resource "alicloud_kvstore_account" "default" {
		account_name     =   "tftest"
		account_password =  "YourPassword_123"
		instance_id      =   alicloud_kvstore_instance.default.id
	}
	`, name)
}
