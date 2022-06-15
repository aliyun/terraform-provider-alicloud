package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudKVStoreInstancesDataSource(t *testing.T) {
	resourceId := "data.alicloud_kvstore_instances.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreInstance-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceKvstoreInstancesDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_kvstore_instance.default.db_instance_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_kvstore_instance.default.db_instance_name}-fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_kvstore_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_kvstore_instance.default.id}-fake"},
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_kvstore_instance.default.id}"},
			"enable_details": "true",
			"tags":           "${alicloud_kvstore_instance.default.tags}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_kvstore_instance.default.id}"},
			"enable_details": "true",
			"tags": map[string]string{
				"Created": "TF-fake",
				"For":     "update test fake",
			},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_kvstore_instance.default.id}"},
			"status":         "Normal",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_kvstore_instance.default.id}"},
			"status":         "Changing",
			"enable_details": "true",
		}),
	}
	paramsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_kvstore_instance.default.id}"},
			"engine_version": "4.0",
			"instance_class": "${data.alicloud_kvstore_instance_classes.default.instance_classes.0}",
			"instance_type":  "Redis",
			"payment_type":   "PostPaid",
			"vpc_id":         "${data.alicloud_vpcs.default.ids.0}",
			"vswitch_id":     "${data.alicloud_vswitches.default.ids.0}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_kvstore_instance.default.id}-fake"},
			"engine_version": "4.0",
			"instance_class": "${data.alicloud_kvstore_instance_classes.default.instance_classes.0}",
			"instance_type":  "Redis",
			"payment_type":   "PostPaid",
			"vpc_id":         "${data.alicloud_vpcs.default.ids.0}",
			"vswitch_id":     "${data.alicloud_vswitches.default.ids.0}",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_kvstore_instance.default.id}"},
			"status":         "Normal",
			"name_regex":     name,
			"engine_version": "4.0",
			"instance_class": "${data.alicloud_kvstore_instance_classes.default.instance_classes.0}",
			"instance_type":  "Redis",
			"payment_type":   "PostPaid",
			"vpc_id":         "${data.alicloud_vpcs.default.ids.0}",
			"vswitch_id":     "${data.alicloud_vswitches.default.ids.0}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_kvstore_instance.default.id}-fake"},
			"status":         "Normal",
			"name_regex":     name,
			"engine_version": "4.0",
			"instance_class": "${data.alicloud_kvstore_instance_classes.default.instance_classes.0}",
			"instance_type":  "Redis",
			"payment_type":   "PostPaid",
			"vpc_id":         "${data.alicloud_vpcs.default.ids.0}",
			"vswitch_id":     "${data.alicloud_vswitches.default.ids.0}",
		}),
	}
	var existKvstoreInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                   "1",
			"ids.0":                                   CHECKSET,
			"names.#":                                 "1",
			"names.0":                                 name,
			"instances.#":                             "1",
			"instances.0.architecture_type":           CHECKSET,
			"instances.0.bandwidth":                   CHECKSET,
			"instances.0.capacity":                    CHECKSET,
			"instances.0.config.%":                    CHECKSET,
			"instances.0.connection_mode":             "",
			"instances.0.connection_domain":           CHECKSET,
			"instances.0.id":                          CHECKSET,
			"instances.0.db_instance_id":              CHECKSET,
			"instances.0.db_instance_name":            name,
			"instances.0.destroy_time":                "",
			"instances.0.end_time":                    "",
			"instances.0.engine_version":              "4.0",
			"instances.0.has_renew_change_order":      CHECKSET,
			"instances.0.instance_class":              CHECKSET,
			"instances.0.instance_type":               "Redis",
			"instances.0.is_rds":                      CHECKSET,
			"instances.0.max_connections":             CHECKSET,
			"instances.0.network_type":                "VPC",
			"instances.0.node_type":                   CHECKSET,
			"instances.0.package_type":                CHECKSET,
			"instances.0.payment_type":                "PostPaid",
			"instances.0.port":                        CHECKSET,
			"instances.0.private_ip":                  CHECKSET,
			"instances.0.qps":                         CHECKSET,
			"instances.0.replacate_id":                "",
			"instances.0.resource_group_id":           CHECKSET,
			"instances.0.search_key":                  "",
			"instances.0.status":                      "Normal",
			"instances.0.vpc_cloud_instance_id":       "",
			"instances.0.vswitch_id":                  CHECKSET,
			"instances.0.vpc_id":                      CHECKSET,
			"instances.0.zone_id":                     CHECKSET,
			"instances.0.instance_release_protection": "false",
			"instances.0.maintain_end_time":           CHECKSET,
			"instances.0.maintain_start_time":         CHECKSET,
			"instances.0.vpc_auth_mode":               CHECKSET,
			"instances.0.auto_renew":                  CHECKSET,
			"instances.0.auto_renew_period":           CHECKSET,
			"instances.0.security_group_id":           "",
			"instances.0.security_ip_group_attribute": CHECKSET,
			"instances.0.security_ip_group_name":      CHECKSET,
			"instances.0.security_ips.#":              CHECKSET,
			"instances.0.secondary_zone_id":           "",
		}
	}

	var fakeKvstoreInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"instances.#": "0",
		}
	}

	var kvstoreInstancesInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKvstoreInstanceMapFunc,
		fakeMapFunc:  fakeKvstoreInstanceMapFunc,
	}
	kvstoreInstancesInfo.dataSourceTestCheck(t, 0, nameRegexConf, idsConf, tagsConf, statusConf, paramsConf, allConf)
}

func dataSourceKvstoreInstancesDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_kvstore_zones" "default" {
	  	instance_charge_type = "PostPaid"
	}
	data "alicloud_vpcs" "default" {
	  	name_regex = "default-NODELETING"
	}
	data "alicloud_vswitches" "default" {
	  	zone_id = data.alicloud_kvstore_zones.default.zones.0.id
	  	vpc_id  = data.alicloud_vpcs.default.ids.0
	}
	data "alicloud_kvstore_instance_classes" "default" {
	  	zone_id        = data.alicloud_kvstore_zones.default.zones.0.id
	  	engine         = "Redis"
	  	engine_version = "4.0"
	}
	data "alicloud_resource_manager_resource_groups" "default" {
	}
	resource "alicloud_kvstore_instance" "default" {
		vswitch_id       = data.alicloud_vswitches.default.ids.0
	  	db_instance_name = "%s"
	  	security_ips = ["10.23.12.24"]
	  	instance_type  = "Redis"
	  	engine_version = "4.0"
	  	config = {
			appendonly             = "yes",
			lazyfree-lazy-eviction = "yes",
		}
	  	tags = {
			Created = "TF",
			For     = "update test",
	  	}
		resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.1
	  	zone_id           = data.alicloud_kvstore_zones.default.zones.0.id
	  	instance_class    = data.alicloud_kvstore_instance_classes.default.instance_classes.0
	}
	`, name)
}
