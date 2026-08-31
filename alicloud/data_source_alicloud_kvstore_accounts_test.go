package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudKVStoreAccountsDataSource(t *testing.T) {
	resourceId := "data.alicloud_kvstore_accounts.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreAccount-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceKvstoreAccountsDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":  "${alicloud_kvstore_account.default.account_name}",
			"instance_id": "${alicloud_kvstore_instance.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":  "${alicloud_kvstore_account.default.account_name}_fake",
			"instance_id": "${alicloud_kvstore_instance.default.id}",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":  "${alicloud_kvstore_instance.default.id}",
			"account_name": "${alicloud_kvstore_account.default.account_name}",
			"status":       "Available",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":  "${alicloud_kvstore_instance.default.id}",
			"account_name": "${alicloud_kvstore_account.default.account_name}",
			"status":       "Unavailable",
		}),
	}

	var existKvstoreAccountMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"ids.0":                   CHECKSET,
			"names.#":                 "1",
			"names.0":                 "tftest",
			"accounts.#":              "1",
			"accounts.0.id":           CHECKSET,
			"accounts.0.account_name": "tftest",
			"accounts.0.account_type": CHECKSET,
			"accounts.0.instance_id":  CHECKSET,
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

	var kvstoreAccountsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKvstoreAccountMapFunc,
		fakeMapFunc:  fakeKvstoreAccountMapFunc,
	}
	kvstoreAccountsInfo.dataSourceTestCheck(t, rand, nameRegexConf, statusConf)
}

func dataSourceKvstoreAccountsDependence(name string) string {
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
	resource "alicloud_kvstore_account" "default" {
		account_name     =   "tftest"
		account_password =  "YourPassword_123"
		instance_id      =   alicloud_kvstore_instance.default.id
	}
	`, name)
}
