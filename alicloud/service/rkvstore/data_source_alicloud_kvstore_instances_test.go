package rkvstore_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudKvstoreInstancesDataSource(t *testing.T) {
	resourceId := "data.alicloud_kvstore_instances.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreInstance-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceKvstoreInstancesDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_kvstore_instance.default.db_instance_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_kvstore_instance.default.db_instance_name}-fake",
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_kvstore_instance.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_kvstore_instance.default.id}-fake"},
		}),
	}
	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_kvstore_instance.default.id}"},
			"enable_details": "true",
			"tags":           "${alicloud_kvstore_instance.default.tags}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_kvstore_instance.default.id}"},
			"enable_details": "true",
			"tags": map[string]string{
				"Created": "TF-fake",
				"For":     "update test fake",
			},
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_kvstore_instance.default.id}"},
			"status":         "Normal",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_kvstore_instance.default.id}"},
			"status":         "Changing",
			"enable_details": "true",
		}),
	}
	paramsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_kvstore_instance.default.id}"},
			"engine_version": "5.0",
			"instance_class": "${data.alicloud_kvstore_instance_classes.default.instance_classes.0}",
			"instance_type":  "Redis",
			"payment_type":   "PostPaid",
			"vpc_id":         "${data.alicloud_vpcs.default.ids.0}",
			"vswitch_id":     "${data.alicloud_vswitches.default.ids.0}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_kvstore_instance.default.id}-fake"},
			"engine_version": "5.0",
			"instance_class": "${data.alicloud_kvstore_instance_classes.default.instance_classes.0}",
			"instance_type":  "Redis",
			"payment_type":   "PostPaid",
			"vpc_id":         "${data.alicloud_vpcs.default.ids.0}",
			"vswitch_id":     "${data.alicloud_vswitches.default.ids.0}",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_kvstore_instance.default.id}"},
			"status":         "Normal",
			"name_regex":     name,
			"engine_version": "5.0",
			"instance_class": "${data.alicloud_kvstore_instance_classes.default.instance_classes.0}",
			"instance_type":  "Redis",
			"payment_type":   "PostPaid",
			"vpc_id":         "${data.alicloud_vpcs.default.ids.0}",
			"vswitch_id":     "${data.alicloud_vswitches.default.ids.0}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_kvstore_instance.default.id}-fake"},
			"status":         "Normal",
			"name_regex":     name,
			"engine_version": "5.0",
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
			"ids.0":                                   acctest.CHECKSET,
			"names.#":                                 "1",
			"names.0":                                 name,
			"instances.#":                             "1",
			"instances.0.architecture_type":           acctest.CHECKSET,
			"instances.0.bandwidth":                   acctest.CHECKSET,
			"instances.0.capacity":                    acctest.CHECKSET,
			"instances.0.config.%":                    acctest.CHECKSET,
			"instances.0.connection_mode":             "",
			"instances.0.connection_domain":           acctest.CHECKSET,
			"instances.0.id":                          acctest.CHECKSET,
			"instances.0.db_instance_id":              acctest.CHECKSET,
			"instances.0.db_instance_name":            name,
			"instances.0.destroy_time":                "",
			"instances.0.end_time":                    "",
			"instances.0.engine_version":              "5.0",
			"instances.0.has_renew_change_order":      acctest.CHECKSET,
			"instances.0.instance_class":              acctest.CHECKSET,
			"instances.0.instance_type":               "Redis",
			"instances.0.is_rds":                      acctest.CHECKSET,
			"instances.0.max_connections":             acctest.CHECKSET,
			"instances.0.network_type":                "VPC",
			"instances.0.node_type":                   acctest.CHECKSET,
			"instances.0.package_type":                acctest.CHECKSET,
			"instances.0.payment_type":                "PostPaid",
			"instances.0.port":                        acctest.CHECKSET,
			"instances.0.private_ip":                  acctest.CHECKSET,
			"instances.0.qps":                         acctest.CHECKSET,
			"instances.0.replacate_id":                "",
			"instances.0.resource_group_id":           acctest.CHECKSET,
			"instances.0.search_key":                  "",
			"instances.0.status":                      "Normal",
			"instances.0.vpc_cloud_instance_id":       "",
			"instances.0.vswitch_id":                  acctest.CHECKSET,
			"instances.0.vpc_id":                      acctest.CHECKSET,
			"instances.0.zone_id":                     acctest.CHECKSET,
			"instances.0.instance_release_protection": "false",
			"instances.0.maintain_end_time":           acctest.CHECKSET,
			"instances.0.maintain_start_time":         acctest.CHECKSET,
			"instances.0.vpc_auth_mode":               acctest.CHECKSET,
			"instances.0.auto_renew":                  acctest.CHECKSET,
			"instances.0.auto_renew_period":           acctest.CHECKSET,
			"instances.0.security_group_id":           "",
			"instances.0.security_ip_group_attribute": acctest.CHECKSET,
			"instances.0.security_ip_group_name":      acctest.CHECKSET,
			"instances.0.security_ips.#":              acctest.CHECKSET,
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

	var kvstoreInstancesInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existKvstoreInstanceMapFunc,
		FakeMapFunc:  fakeKvstoreInstanceMapFunc,
	}
	kvstoreInstancesInfo.DataSourceTestCheck(t, 0, nameRegexConf, idsConf, tagsConf, statusConf, paramsConf, allConf)
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
	  	zone_id               = data.alicloud_kvstore_zones.default.zones.0.id
	  	engine                = "Redis"
	  	engine_version        = "5.0"
	  	instance_charge_type  = "PostPaid"
	  	architecture          = "standard"
	}
	data "alicloud_resource_manager_resource_groups" "default" {
	}
	resource "alicloud_kvstore_instance" "default" {
		vswitch_id       = data.alicloud_vswitches.default.ids.0
	  	db_instance_name = "%s"
	  	security_ips = ["10.23.12.24"]
	  	instance_type  = "Redis"
	  	engine_version = "5.0"
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
