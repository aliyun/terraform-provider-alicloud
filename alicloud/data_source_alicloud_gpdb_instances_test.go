package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudGpdbInstancesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_gpdb_instance.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_gpdb_instance.default.description}_fake"`,
		}),
	}
	availabilityZoneConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_gpdb_instance.default.description}"`,
			"availability_zone": `"${data.alicloud_zones.default.zones.0.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_gpdb_instance.default.description}"`,
			"availability_zone": `"${data.alicloud_zones.default.zones.1.id}"`,
		}),
	}
	fullConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_gpdb_instance.default.description}"`,
			"availability_zone": `"${data.alicloud_zones.default.zones.0.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_gpdb_instance.default.description}_fake"`,
			"availability_zone": `"${data.alicloud_zones.default.zones.0.id}"`,
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#":                       CHECKSET,
			"instances.0.id":                    CHECKSET,
			"instances.0.description":           fmt.Sprintf("tf-testAccGpdbInstance_datasource_%d", rand),
			"instances.0.engine":                "gpdb",
			"instances.0.engine_version":        "4.3",
			"instances.0.instance_class":        "gpdb.group.segsdx2",
			"instances.0.instance_group_count":  "2",
			"instances.0.region_id":             CHECKSET,
			"instances.0.status":                CHECKSET,
			"instances.0.creation_time":         CHECKSET,
			"instances.0.instance_network_type": CHECKSET,
			"instances.0.charge_type":           CHECKSET,
			"ids.#":                             "1",
			"ids.0":                             CHECKSET,
		}
	}
	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#": "0",
			"ids.#":       "0",
		}
	}
	var CheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_gpdb_instances.default",
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}

	CheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, availabilityZoneConf, fullConf)
}

func testAccCheckAlicloudGpdbDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
        data "alicloud_zones" "default" {
            available_resource_creation = "Gpdb"
        }
        resource "alicloud_vpc" "default" {
            description = "${var.name}"
            cidr_block  = "172.16.0.0/16"
        }
        resource "alicloud_vswitch" "default" {
            availability_zone = "${data.alicloud_zones.default.zones.0.id}"
            vpc_id            = "${alicloud_vpc.default.id}"
            cidr_block        = "172.16.0.0/24"
            description       = "${var.name}"
        }
        variable "name" {
            default = "tf-testAccGpdbInstance_datasource_%d"
        }
        resource "alicloud_gpdb_instance" "default" {
            vswitch_id           = "${alicloud_vswitch.default.id}"
            engine               = "gpdb"
            engine_version       = "4.3"
            instance_class       = "gpdb.group.segsdx2"
            instance_group_count = "2"
            description          = "${var.name}"
        }
        data "alicloud_gpdb_instances" "default" {
            %s
        }`, rand, strings.Join(pairs, "\n  "))
}
