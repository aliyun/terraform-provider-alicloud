package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"strings"
	"testing"
)

func TestAccAlicloudGpdbInstancesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.GPDBSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_gpdb_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_gpdb_instance.default.id}_fake"]`,
		}),
	}
	dBInstanceCategoriesConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbInstancesDataSourceName(rand, map[string]string{
			"ids":                    `["${alicloud_gpdb_instance.default.id}"]`,
			"db_instance_categories": `"highavailability"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbInstancesDataSourceName(rand, map[string]string{
			"ids":                    `["${alicloud_gpdb_instance.default.id}"]`,
			"db_instance_categories": `"highavailability_fake"`,
		}),
	}
	dBInstanceModesConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_gpdb_instance.default.id}"]`,
			"db_instance_modes": `"storageelastic"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_gpdb_instance.default.id}"]`,
			"db_instance_modes": `"storageelastic_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_gpdb_instance.default.id}"]`,
			"status": `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_gpdb_instance.default.id}"]`,
			"status": `"DBInstanceClassChanging"`,
		}),
	}
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbInstancesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_gpdb_instance.default.id}"]`,
			"description": `"${alicloud_gpdb_instance.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbInstancesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_gpdb_instance.default.id}"]`,
			"description": `"${alicloud_gpdb_instance.default.description}_fake"`,
		}),
	}
	instanceNetworkTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbInstancesDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_gpdb_instance.default.id}"]`,
			"instance_network_type": `"${alicloud_gpdb_instance.default.instance_network_type}"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbInstancesDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_gpdb_instance.default.id}"]`,
			"instance_network_type": `"${alicloud_gpdb_instance.default.instance_network_type}_fake"`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_gpdb_instance.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_gpdb_instance.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_gpdb_instance.default.id}"]`,
			"tags": `{
				"Created" = "TF"
				"For" = "acceptance test"
		}`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_gpdb_instance.default.id}"]`,
			"tags": `{
				"Created" = "TF-fake"
				"For" = "acceptance test fake"
			}`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbInstancesDataSourceName(rand, map[string]string{
			"db_instance_categories": `"highavailability"`,
			"db_instance_modes":      `"storageelastic"`,
			"description":            `"${alicloud_gpdb_instance.default.description}"`,
			"ids":                    `["${alicloud_gpdb_instance.default.id}"]`,
			"instance_network_type":  `"${alicloud_gpdb_instance.default.instance_network_type}"`,
			"resource_group_id":      `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
			"status":                 `"Running"`,
			"tags": `{
				"Created" = "TF"
				"For" = "acceptance test"
		}`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbInstancesDataSourceName(rand, map[string]string{
			"db_instance_categories": `"highavailability_fake"`,
			"db_instance_modes":      `"storageelastic_fake"`,
			"description":            `"${alicloud_gpdb_instance.default.description}_fake"`,
			"ids":                    `["${alicloud_gpdb_instance.default.id}_fake"]`,
			"instance_network_type":  `"${alicloud_gpdb_instance.default.instance_network_type}_fake"`,
			"resource_group_id":      `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
			"status":                 `"DBInstanceClassChanging"`,
			"tags": `{
				"Created" = "TF-fake"
				"For" = "acceptance test fake"
			}`,
		}),
	}
	var existAlicloudGpdbInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"instances.#":                       "1",
			"instances.0.status":                "Running",
			"instances.0.connection_string":     CHECKSET,
			"instances.0.cpu_cores":             CHECKSET,
			"instances.0.create_time":           CHECKSET,
			"instances.0.db_instance_category":  "HighAvailability",
			"instances.0.id":                    CHECKSET,
			"instances.0.db_instance_id":        CHECKSET,
			"instances.0.db_instance_mode":      "StorageElastic",
			"instances.0.description":           CHECKSET,
			"instances.0.engine":                "gpdb",
			"instances.0.engine_version":        "6.0",
			"instances.0.ip_whitelist.#":        "1",
			"instances.0.instance_network_type": "VPC",
			"instances.0.maintain_end_time":     CHECKSET,
			"instances.0.maintain_start_time":   CHECKSET,
			"instances.0.master_node_num":       "1",
			"instances.0.memory_size":           CHECKSET,
			"instances.0.payment_type":          "PayAsYouGo",
			"instances.0.seg_node_num":          CHECKSET,
			"instances.0.storage_size":          "50",
			"instances.0.storage_type":          "cloud_essd",
			"instances.0.vswitch_id":            CHECKSET,
			"instances.0.vpc_id":                CHECKSET,
			"instances.0.zone_id":               CHECKSET,
			"instances.0.region_id":             CHECKSET,
			"instances.0.availability_zone":     CHECKSET,
			"instances.0.creation_time":         CHECKSET,
			"instances.0.charge_type":           CHECKSET,
		}
	}
	var fakeAlicloudGpdbInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}
	var alicloudGpdbInstancesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_gpdb_instances.default",
		existMapFunc: existAlicloudGpdbInstancesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudGpdbInstancesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudGpdbInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, dBInstanceCategoriesConf, dBInstanceModesConf, statusConf, descriptionConf, instanceNetworkTypeConf, resourceGroupIdConf, tagsConf, allConf)
}
func testAccCheckAlicloudGpdbInstancesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccDBInstance-%d"
}

data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_gpdb_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.ids.0
}
resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_gpdb_zones.default.ids.0
  vswitch_name = var.name
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}
resource "alicloud_gpdb_instance" "default" {
  db_instance_category  = "HighAvailability"
  db_instance_class     = "gpdb.group.segsdx1"
  db_instance_mode      = "StorageElastic"
  description           = var.name
  engine                = "gpdb"
  engine_version        = "6.0"
  zone_id               = data.alicloud_gpdb_zones.default.ids.0
  instance_network_type = "VPC"
  instance_spec         = "2C16G"
  master_node_num       = 1
  payment_type          = "PayAsYouGo"
  private_ip_address    = "1.1.1.1"
  seg_storage_type      = "cloud_essd"
  seg_node_num          = 4
  storage_size          = 50
  vpc_id                = data.alicloud_vpcs.default.ids.0
  vswitch_id            = local.vswitch_id
  ip_whitelist {
    security_ip_list = "127.0.0.1"
  }
  tags = {
    Created = "TF"
    For =     "acceptance test"
  }
}
data "alicloud_gpdb_instances" "default" {	
	enable_details = true
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
