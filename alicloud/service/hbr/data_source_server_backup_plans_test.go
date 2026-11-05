package hbr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudHbrServerBackupPlansDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.HbrSupportRegions)
	resourceId := "data.alicloud_hbr_server_backup_plans.default"
	name := fmt.Sprintf("tf-testAccHbrServerBackupPlanTest%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceHbrServerBackupPlanSourceConfig)

	backupIdsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_hbr_server_backup_plan.example.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_hbr_server_backup_plan.example.id}_fake"},
		}),
	}

	planIdBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"filters": []map[string]interface{}{
				{
					"key":    "planId",
					"values": []string{"${alicloud_hbr_server_backup_plan.example.id}"},
				},
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"filters": []map[string]interface{}{
				{
					"key":    "planId",
					"values": []string{"${alicloud_hbr_server_backup_plan.example.id}_fake"},
				},
			},
		}),
	}

	instanceIdBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"filters": []map[string]interface{}{
				{
					"key":    "instanceId",
					"values": []string{"${alicloud_instance.default.id}"},
				},
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"filters": []map[string]interface{}{
				{
					"key":    "instanceId",
					"values": []string{"${alicloud_instance.default.id}_fake"},
				},
			},
		}),
	}

	planNameBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"filters": []map[string]interface{}{
				{
					"key":    "planName",
					"values": []string{"${alicloud_hbr_server_backup_plan.example.ecs_server_backup_plan_name}"},
				},
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"filters": []map[string]interface{}{
				{
					"key":    "planName",
					"values": []string{"${alicloud_hbr_server_backup_plan.example.ecs_server_backup_plan_name}_fake"},
				},
			},
		}),
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	HbrServerBackupPlanCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, backupIdsConf, planIdBackupConf, instanceIdBackupConf, planNameBackupConf)
}

func dataSourceHbrServerBackupPlanSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" default {
  available_resource_creation = "Instance"
}
data "alicloud_instance_types" "default" {
	availability_zone = data.alicloud_zones.default.zones.0.id
  	cpu_core_count    = 1
	memory_size       = 2
}
data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_security_group" "default" {
  security_group_name = var.name
  description = "New security group"
  vpc_id = data.alicloud_vpcs.default.ids.0
}
data "alicloud_images" "default" {
  owners      = "system"
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
}
resource "alicloud_instance" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  instance_name   = var.name
  host_name       = var.name
  image_id        = data.alicloud_images.default.images.0.id
  instance_type   = data.alicloud_instance_types.default.instance_types.0.id
  security_groups = [alicloud_security_group.default.id]
  vswitch_id      = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_hbr_server_backup_plan" "example" {
  ecs_server_backup_plan_name = var.name
  instance_id                 = alicloud_instance.default.id
  schedule                    = "I|1602673264|PT2H"
  retention                   = 1
  detail {
    app_consistent     = false
    snapshot_group     = true
    enable_fs_freeze   = true
    pre_script_path    = ""
    post_script_path   = ""
    timeout_in_seconds = 60
    disk_id_list       = ["/home"]
  }
  disabled = false
}`, name)
}

var existHbrServerBackupPlanMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"plans.#":                                acctest.CHECKSET,
		"plans.0.id":                             acctest.CHECKSET,
		"plans.0.create_time":                    acctest.CHECKSET,
		"plans.0.detail.#":                       acctest.CHECKSET,
		"plans.0.detail.0.app_consistent":        acctest.CHECKSET,
		"plans.0.detail.0.snapshot_group":        acctest.CHECKSET,
		"plans.0.detail.0.enable_fs_freeze":      acctest.CHECKSET,
		"plans.0.detail.0.destination_region_id": "",
		"plans.0.detail.0.pre_script_path":       "",
		"plans.0.detail.0.post_script_path":      "",
		"plans.0.detail.0.timeout_in_seconds":    acctest.CHECKSET,
		"plans.0.detail.0.disk_id_list.#":        acctest.CHECKSET,
		"plans.0.detail.0.do_copy":               acctest.CHECKSET,
		"plans.0.detail.0.destination_retention": acctest.CHECKSET,
		"plans.0.disabled":                       acctest.CHECKSET,
		"plans.0.ecs_server_backup_plan_id":      acctest.CHECKSET,
		"plans.0.ecs_server_backup_plan_name":    acctest.CHECKSET,
		"plans.0.instance_id":                    acctest.CHECKSET,
		"plans.0.retention":                      acctest.CHECKSET,
		"plans.0.schedule":                       acctest.CHECKSET,
	}
}

var fakeHbrServerBackupPlanMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"plans.#": "0",
	}
}

var HbrServerBackupPlanCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_hbr_server_backup_plans.default",
	ExistMapFunc: existHbrServerBackupPlanMapFunc,
	FakeMapFunc:  fakeHbrServerBackupPlanMapFunc,
}
