package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCddcDedicatedHostsDataSource(t *testing.T) {
	resourceId := "data.alicloud_cddc_dedicated_hosts.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cddcdedicatedhost-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceCddcDedicatedHostsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"enable_details":          "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}-fake"},
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"enable_details":          "true",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"status":                  "1",
			"enable_details":          "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"status":                  "2",
			"enable_details":          "true",
		}),
	}
	zoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"zone_id":                 "${alicloud_cddc_dedicated_host.default.zone_id}",
			"enable_details":          "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"zone_id":                 "${alicloud_cddc_dedicated_host.default.zone_id}-fake",
			"enable_details":          "true",
		}),
	}
	allocationStatusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"allocation_status":       "Allocatable",
			"enable_details":          "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"allocation_status":       "Suspended",
			"enable_details":          "true",
		}),
	}
	hostTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"host_type":               "dhg_cloud_ssd",
			"enable_details":          "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"host_type":               "dhg_local_ssd",
			"enable_details":          "true",
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"enable_details":          "true",
			"tags": map[string]string{
				"Create": "TF",
				"For":    "CDDC_DEDICATED",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"enable_details":          "true",
			"tags": map[string]string{
				"Create": "CDDC_DEDICATED",
				"For":    "TF",
			},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"status":                  "1",
			"allocation_status":       "Allocatable",
			"zone_id":                 "${alicloud_cddc_dedicated_host.default.zone_id}",
			"enable_details":          "true",
			"tags": map[string]string{
				"Create": "TF",
				"For":    "CDDC_DEDICATED",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}-fake"},
			"status":                  "2",
			"allocation_status":       "Suspended",
			"zone_id":                 "${alicloud_cddc_dedicated_host.default.zone_id}-fake",
			"enable_details":          "true",
			"tags": map[string]string{
				"Create": "CDDC_DEDICATED",
				"For":    "TF",
			},
		}),
	}
	var existCddcDedicatedHostMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"ids.0":                           CHECKSET,
			"hosts.#":                         "1",
			"hosts.0.bastion_instance_id":     "",
			"hosts.0.cpu_allocation_ratio":    CHECKSET,
			"hosts.0.cpu_used":                CHECKSET,
			"hosts.0.create_time":             CHECKSET,
			"hosts.0.dedicated_host_group_id": CHECKSET,
			"hosts.0.id":                      CHECKSET,
			"hosts.0.dedicated_host_id":       CHECKSET,
			"hosts.0.disk_allocation_ratio":   CHECKSET,
			"hosts.0.ecs_class_code":          CHECKSET,
			"hosts.0.end_time":                CHECKSET,
			"hosts.0.engine":                  CHECKSET,
			"hosts.0.expired_time":            CHECKSET,
			"hosts.0.host_class":              CHECKSET,
			"hosts.0.host_cpu":                CHECKSET,
			"hosts.0.host_mem":                CHECKSET,
			"hosts.0.host_name":               fmt.Sprintf("tf-testacc-cddcdedicatedhost-%d", rand),
			"hosts.0.host_storage":            CHECKSET,
			"hosts.0.host_type":               "dhg_cloud_ssd",
			"hosts.0.image_category":          "",
			"hosts.0.ip_address":              CHECKSET,
			"hosts.0.mem_allocation_ratio":    CHECKSET,
			"hosts.0.memory_used":             CHECKSET,
			"hosts.0.open_permission":         CHECKSET,
			"hosts.0.allocation_status":       "Allocatable",
			"hosts.0.status":                  "1",
			"hosts.0.storage_used":            CHECKSET,
			"hosts.0.tags.%":                  "2",
			"hosts.0.tags.Created":            "TF",
			"hosts.0.tags.For":                "CDDC_DEDICATED",
			"hosts.0.vswitch_id":              CHECKSET,
			"hosts.0.vpc_id":                  CHECKSET,
			"hosts.0.zone_id":                 CHECKSET,
		}
	}

	var fakeCddcDedicatedHostMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"hosts.#": "0",
		}
	}

	var CddcDedicatedHostCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCddcDedicatedHostMapFunc,
		fakeMapFunc:  fakeCddcDedicatedHostMapFunc,
	}

	CddcDedicatedHostCheckInfo.dataSourceTestCheck(t, rand, idsConf, statusConf, zoneIdConf, allocationStatusConf, hostTypeConf, tagsConf, allConf)
}

func dataSourceCddcDedicatedHostsDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_cddc_zones" "default" {}

data "alicloud_cddc_host_ecs_level_infos" "default" {
  db_type      = "mysql"
  zone_id      = data.alicloud_cddc_zones.default.ids.0
  storage_type = "cloud_essd"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_cddc_zones.default.ids.0
}

data "alicloud_cddc_dedicated_host_groups" "default" {
  engine     = "MySQL"
}

resource "alicloud_cddc_dedicated_host_group" "default" {
	count = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? 0 : 1
	engine = "MySQL"
	vpc_id = data.alicloud_vpcs.default.ids.0
	cpu_allocation_ratio = 101
	mem_allocation_ratio = 50
	disk_allocation_ratio = 200
	allocation_policy = "Evenly"
	host_replace_policy = "Manual"
	dedicated_host_group_desc = var.name
	open_permission = true
}
locals {
	dedicated_host_group_id = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? data.alicloud_cddc_dedicated_host_groups.default.ids.0 : concat(alicloud_cddc_dedicated_host_group.default[*].id, [""])[0]
}

resource "alicloud_cddc_dedicated_host" "default" {
  host_name               = var.name
  dedicated_host_group_id = local.dedicated_host_group_id
  host_class              = data.alicloud_cddc_host_ecs_level_infos.default.infos.0.res_class_code
  zone_id                 = data.alicloud_cddc_zones.default.ids.0
  vswitch_id              = data.alicloud_vswitches.default.ids.0
  payment_type            = "Subscription"
  tags = {
	Created = "TF"
	For = "CDDC_DEDICATED"
  }
}`, name)
}
