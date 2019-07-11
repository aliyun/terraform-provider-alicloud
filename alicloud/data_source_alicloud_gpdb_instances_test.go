package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudGpdbInstancesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_gpdb_instances.default"
	name := fmt.Sprintf("tf-testAccGpdbInstance_datasource-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceGpdbConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_gpdb_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_gpdb_instance.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_gpdb_instance.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_gpdb_instance.default.description}_fake",
		}),
	}

	availabilityZoneConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alicloud_gpdb_instance.default.description}",
			"availability_zone": "${data.alicloud_zones.default.zones.0.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alicloud_gpdb_instance.default.description}",
			"availability_zone": "${data.alicloud_zones.default.zones.1.id}",
		}),
	}

	vSwitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_gpdb_instance.default.description}",
			"vswitch_id": "${alicloud_gpdb_instance.default.vswitch_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_gpdb_instance.default.description}",
			"vswitch_id": "unknow",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_gpdb_instance.default.id}"},
			"name_regex":        "${alicloud_gpdb_instance.default.description}",
			"availability_zone": "${data.alicloud_zones.default.zones.0.id}",
			"vswitch_id":        "${alicloud_gpdb_instance.default.vswitch_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_gpdb_instance.default.id}_fake"},
			"name_regex":        "${alicloud_gpdb_instance.default.description}_fake",
			"availability_zone": "${data.alicloud_zones.default.zones.0.id}",
			"vswitch_id":        "unknow",
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"instances.#":                       CHECKSET,
			"instances.0.id":                    CHECKSET,
			"instances.0.description":           fmt.Sprintf("tf-testAccGpdbInstance_datasource-%d", rand),
			"instances.0.engine":                "gpdb",
			"instances.0.engine_version":        "4.3",
			"instances.0.instance_class":        "gpdb.group.segsdx2",
			"instances.0.instance_group_count":  "2",
			"instances.0.region_id":             CHECKSET,
			"instances.0.status":                CHECKSET,
			"instances.0.creation_time":         CHECKSET,
			"instances.0.instance_network_type": CHECKSET,
			"instances.0.charge_type":           CHECKSET,
		}
	}

	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#": "0",
			"ids.#":       "0",
		}
	}

	var CheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}

	CheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, availabilityZoneConf, vSwitchIdConf, allConf)
}

func dataSourceGpdbConfigDependence(name string) string {
	return fmt.Sprintf(`
        data "alicloud_zones" "default" {
            available_resource_creation = "Gpdb"
        }
        resource "alicloud_vpc" "default" {
            name = "${var.name}"
            cidr_block  = "172.16.0.0/16"
        }
        resource "alicloud_vswitch" "default" {
            availability_zone = "${data.alicloud_zones.default.zones.0.id}"
            vpc_id            = "${alicloud_vpc.default.id}"
            cidr_block        = "172.16.0.0/24"
            name       = "${var.name}"
        }
        variable "name" {
            default = "%s"
        }
        resource "alicloud_gpdb_instance" "default" {
            vswitch_id           = "${alicloud_vswitch.default.id}"
            engine               = "gpdb"
            engine_version       = "4.3"
            instance_class       = "gpdb.group.segsdx2"
            instance_group_count = "2"
            description          = "${var.name}"
        }`, name)
}
