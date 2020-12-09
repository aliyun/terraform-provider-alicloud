package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudKvstoreConnectionsDataSource(t *testing.T) {
	resourceId := "data.alicloud_kvstore_connections.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreConnection-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceKvstoreConnectionsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_kvstore_instance.default.id}"},
		}),
	}
	var existKvstoreConnectionMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"ids.0":                              CHECKSET,
			"connections.#":                      "1",
			"connections.0.connection_string":    "allocatetest.redis.rds.aliyuncs.com",
			"connections.0.db_instance_net_type": "0",
			"connections.0.expired_time":         "",
			"connections.0.ip_address":           CHECKSET,
			"connections.0.port":                 "6370",
			"connections.0.upgradeable":          "0",
			"connections.0.vpc_id":               "",
			"connections.0.vpc_instance_id":      "",
			"connections.0.vswitch_id":           "",
		}
	}

	var fakeKvstoreConnectionMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "0",
			"connections.#": "0",
		}
	}

	var KvstoreInstancesInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKvstoreConnectionMapFunc,
		fakeMapFunc:  fakeKvstoreConnectionMapFunc,
	}

	KvstoreInstancesInfo.dataSourceTestCheck(t, 0, idsConf)
}

func dataSourceKvstoreConnectionsDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_kvstore_zones" "default"{
		instance_charge_type = "PostPaid"
	}
	data "alicloud_vpcs" "default" {
	  is_default = true
	}
	data "alicloud_vswitches" "default" {
	  zone_id = data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id
	  vpc_id = data.alicloud_vpcs.default.ids.0
	}
	data "alicloud_resource_manager_resource_groups" "default" {
	}
	resource "alicloud_kvstore_instance" "default" {
		db_instance_name = "%s"
  		vswitch_id = data.alicloud_vswitches.default.ids.0
		instance_type = "Redis"
		engine_version = "4.0"
		tags = {
			Created = "TF",
			For = "update test",
		}
		resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.1
		zone_id = data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id
		instance_class="redis.master.large.default"
	}
	resource "alicloud_kvstore_connection" "default" {
	  connection_string_prefix = "allocatetest"
	  instance_id = alicloud_kvstore_instance.default.id
	  port = "6370"
	}
	`, name)
}
