package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
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
			"availability_zone": "${data.alicloud_zones.default.zones.0.id}F",
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

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_gpdb_instance.default.id}"},
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "acceptance test",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_gpdb_instance.default.id}"},
			"tags": map[string]interface{}{
				"Created": "TF-fake",
				"For":     "acceptance test",
			},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_gpdb_instance.default.id}"},
			"name_regex":        "${alicloud_gpdb_instance.default.description}",
			"availability_zone": "${data.alicloud_zones.default.zones.0.id}",
			"vswitch_id":        "${alicloud_gpdb_instance.default.vswitch_id}",
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "acceptance test",
			}}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_gpdb_instance.default.id}_fake"},
			"name_regex":        "${alicloud_gpdb_instance.default.description}_fake",
			"availability_zone": "${data.alicloud_zones.default.zones.0.id}",
			"vswitch_id":        "unknow",
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "acceptance test",
			}}),
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
	preCheck := func() {
		testAccPreCheckWithNoDefaultVpc(t)
	}

	CheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, availabilityZoneConf, vSwitchIdConf, tagsConf, allConf)
}

func dataSourceGpdbConfigDependence(name string) string {
	return fmt.Sprintf(`
        data "alicloud_zones" "default" {
            available_resource_creation = "Gpdb"
        }
        data "alicloud_vswitches" "default" {
		  zone_id = data.alicloud_zones.default.ids[0]
		  name_regex = "default-tf--testAcc-00"
		}
        variable "name" {
            default = "%s"
        }
        resource "alicloud_gpdb_instance" "default" {
            vswitch_id           = "${data.alicloud_vswitches.default.ids.0}"
            engine               = "gpdb"
            engine_version       = "4.3"
            instance_class       = "gpdb.group.segsdx2"
            instance_group_count = "2"
            description          = "${var.name}"
			tags 				 = {
				Created = "TF"
				For 	= "acceptance test"
			}
        }`, name)
}
