package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudHBRServerBackupPlansDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.HbrSupportRegions)
	resourceId := "data.alicloud_hbr_server_backup_plans.default"
	name := fmt.Sprintf("tf-testAccHbrServerBackupPlanTest%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceHbrServerBackupPlanSourceConfig)

	backupIdsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_hbr_server_backup_plan.example.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_hbr_server_backup_plan.example.id}_fake"},
		}),
	}

	planIdBackupConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"filters": []map[string]interface{}{
				{
					"key":    "planId",
					"values": []string{"${alicloud_hbr_server_backup_plan.example.id}"},
				},
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"filters": []map[string]interface{}{
				{
					"key":    "planId",
					"values": []string{"${alicloud_hbr_server_backup_plan.example.id}_fake"},
				},
			},
		}),
	}

	instanceIdBackupConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"filters": []map[string]interface{}{
				{
					"key":    "instanceId",
					"values": []string{"${alicloud_instance.default.id}"},
				},
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"filters": []map[string]interface{}{
				{
					"key":    "instanceId",
					"values": []string{"${alicloud_instance.default.id}_fake"},
				},
			},
		}),
	}

	planNameBackupConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"filters": []map[string]interface{}{
				{
					"key":    "planName",
					"values": []string{"${alicloud_hbr_server_backup_plan.example.ecs_server_backup_plan_name}"},
				},
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"filters": []map[string]interface{}{
				{
					"key":    "planName",
					"values": []string{"${alicloud_hbr_server_backup_plan.example.ecs_server_backup_plan_name}_fake"},
				},
			},
		}),
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	HbrServerBackupPlanCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, backupIdsConf, planIdBackupConf, instanceIdBackupConf, planNameBackupConf)
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
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  	cpu_core_count    = 1
	memory_size       = 2
}
data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_security_group" "default" {
  name = "${var.name}"
  description = "New security group"
  vpc_id = data.alicloud_vpcs.default.ids.0
}
data "alicloud_images" "default" {
  owners      = "system"
  name_regex  = "^centos_8"
  most_recent = true
}
resource "alicloud_instance" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_name   = "${var.name}"
  host_name       = "tf-testAcc"
  image_id        = data.alicloud_images.default.images.0.id
  instance_type   = data.alicloud_instance_types.default.instance_types.0.id
  security_groups = [alicloud_security_group.default.id]
  vswitch_id      = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_hbr_server_backup_plan" "example" {
  ecs_server_backup_plan_name = "server_backup_plan"
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
		"plans.#":                                CHECKSET,
		"plans.0.id":                             CHECKSET,
		"plans.0.create_time":                    CHECKSET,
		"plans.0.detail.#":                       CHECKSET,
		"plans.0.detail.0.app_consistent":        CHECKSET,
		"plans.0.detail.0.snapshot_group":        CHECKSET,
		"plans.0.detail.0.enable_fs_freeze":      CHECKSET,
		"plans.0.detail.0.destination_region_id": "",
		"plans.0.detail.0.pre_script_path":       "",
		"plans.0.detail.0.post_script_path":      "",
		"plans.0.detail.0.timeout_in_seconds":    CHECKSET,
		"plans.0.detail.0.disk_id_list.#":        CHECKSET,
		"plans.0.detail.0.do_copy":               CHECKSET,
		"plans.0.detail.0.destination_retention": CHECKSET,
		"plans.0.disabled":                       CHECKSET,
		"plans.0.ecs_server_backup_plan_id":      CHECKSET,
		"plans.0.ecs_server_backup_plan_name":    CHECKSET,
		"plans.0.instance_id":                    CHECKSET,
		"plans.0.retention":                      CHECKSET,
		"plans.0.schedule":                       CHECKSET,
	}
}

var fakeHbrServerBackupPlanMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"plans.#": "0",
	}
}

var HbrServerBackupPlanCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_hbr_server_backup_plans.default",
	existMapFunc: existHbrServerBackupPlanMapFunc,
	fakeMapFunc:  fakeHbrServerBackupPlanMapFunc,
}
