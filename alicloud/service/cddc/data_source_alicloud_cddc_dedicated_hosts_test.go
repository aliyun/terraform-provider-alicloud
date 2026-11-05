package cddc_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCddcDedicatedHostsDataSource(t *testing.T) {
	resourceId := "data.alicloud_cddc_dedicated_hosts.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cddcdedicatedhost-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceCddcDedicatedHostsDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"enable_details":          "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}-fake"},
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"enable_details":          "true",
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"status":                  "1",
			"enable_details":          "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"status":                  "2",
			"enable_details":          "true",
		}),
	}
	zoneIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"zone_id":                 "${alicloud_cddc_dedicated_host.default.zone_id}",
			"enable_details":          "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"zone_id":                 "${alicloud_cddc_dedicated_host.default.zone_id}-fake",
			"enable_details":          "true",
		}),
	}
	allocationStatusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"allocation_status":       "Allocatable",
			"enable_details":          "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"allocation_status":       "Suspended",
			"enable_details":          "true",
		}),
	}
	hostTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"host_type":               "dhg_cloud_ssd",
			"enable_details":          "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"host_type":               "dhg_local_ssd",
			"enable_details":          "true",
		}),
	}
	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"enable_details":          "true",
			"tags": map[string]string{
				"Create": "TF",
				"For":    "CDDC_DEDICATED",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_group_id": "${alicloud_cddc_dedicated_host.default.dedicated_host_group_id}",
			"ids":                     []string{"${alicloud_cddc_dedicated_host.default.id}"},
			"enable_details":          "true",
			"tags": map[string]string{
				"Create": "CDDC_DEDICATED",
				"For":    "TF",
			},
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
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
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"ids.0":                           acctest.CHECKSET,
			"hosts.#":                         "1",
			"hosts.0.bastion_instance_id":     "",
			"hosts.0.cpu_allocation_ratio":    acctest.CHECKSET,
			"hosts.0.cpu_used":                acctest.CHECKSET,
			"hosts.0.create_time":             acctest.CHECKSET,
			"hosts.0.dedicated_host_group_id": acctest.CHECKSET,
			"hosts.0.id":                      acctest.CHECKSET,
			"hosts.0.dedicated_host_id":       acctest.CHECKSET,
			"hosts.0.disk_allocation_ratio":   acctest.CHECKSET,
			"hosts.0.ecs_class_code":          acctest.CHECKSET,
			"hosts.0.end_time":                acctest.CHECKSET,
			"hosts.0.engine":                  acctest.CHECKSET,
			"hosts.0.expired_time":            acctest.CHECKSET,
			"hosts.0.host_class":              acctest.CHECKSET,
			"hosts.0.host_cpu":                acctest.CHECKSET,
			"hosts.0.host_mem":                acctest.CHECKSET,
			"hosts.0.host_name":               fmt.Sprintf("tf-testacc-cddcdedicatedhost-%d", rand),
			"hosts.0.host_storage":            acctest.CHECKSET,
			"hosts.0.host_type":               "dhg_cloud_ssd",
			"hosts.0.image_category":          "",
			"hosts.0.ip_address":              acctest.CHECKSET,
			"hosts.0.mem_allocation_ratio":    acctest.CHECKSET,
			"hosts.0.memory_used":             acctest.CHECKSET,
			"hosts.0.open_permission":         acctest.CHECKSET,
			"hosts.0.allocation_status":       "Allocatable",
			"hosts.0.status":                  "1",
			"hosts.0.storage_used":            acctest.CHECKSET,
			"hosts.0.tags.%":                  "2",
			"hosts.0.tags.Created":            "TF",
			"hosts.0.tags.For":                "CDDC_DEDICATED",
			"hosts.0.vswitch_id":              acctest.CHECKSET,
			"hosts.0.vpc_id":                  acctest.CHECKSET,
			"hosts.0.zone_id":                 acctest.CHECKSET,
		}
	}

	var fakeCddcDedicatedHostMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"hosts.#": "0",
		}
	}

	var CddcDedicatedHostCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existCddcDedicatedHostMapFunc,
		FakeMapFunc:  fakeCddcDedicatedHostMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, connectivity.CDDCSupportRegions)
	}

	CddcDedicatedHostCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, zoneIdConf, allocationStatusConf, hostTypeConf, tagsConf, allConf)
}

func dataSourceCddcDedicatedHostsDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_cddc_zones" "default" {}

data "alicloud_cddc_host_ecs_level_infos" "default" {
  db_type      = "mysql"
  zone_id      = data.alicloud_cddc_zones.default.ids.0
  storage_type = "cloud_essd"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_cddc_dedicated_host_groups.default.groups.0.vpc_id
  zone_id = data.alicloud_cddc_zones.default.ids.0
}

data "alicloud_cddc_dedicated_host_groups" "default" {
  name_regex = "^NO-DELETING"
  engine     = "MySQL"
}


resource "alicloud_cddc_dedicated_host" "default" {
  host_name               = var.name
  dedicated_host_group_id = data.alicloud_cddc_dedicated_host_groups.default.ids.0
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
