package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudHBREcsBackupClientsDataSource(t *testing.T) {
	defer checkoutAccount(t, false)
	checkoutAccount(t, true)
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	rand := acctest.RandIntRange(1000000, 9999999)

	ecsBackupIdsconf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrEcsBackupClientSourceConfig(rand, map[string]string{
			"ids": `[alicloud_hbr_ecs_backup_client.default.id]`,
		}),
		fakeConfig: testAccCheckAlicloudHbrEcsBackupClientSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_hbr_ecs_backup_client.default.id}_fake"]`,
		}),
	}

	ecsInstanceIdconf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrEcsBackupClientSourceConfig(rand, map[string]string{
			"ids":          `[alicloud_hbr_ecs_backup_client.default.id]`,
			"instance_ids": `[alicloud_hbr_ecs_backup_client.default.instance_id]`,
		}),
		fakeConfig: testAccCheckAlicloudHbrEcsBackupClientSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_hbr_ecs_backup_client.default.id}_fake"]`,
			"instance_ids": `["${alicloud_hbr_ecs_backup_client.default.instance_id}_fake"]`,
		}),
	}

	statusconf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrEcsBackupClientSourceConfig(rand, map[string]string{
			"ids":    `[alicloud_hbr_ecs_backup_client.default.id]`,
			"status": `"ACTIVATED"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrEcsBackupClientSourceConfig(rand, map[string]string{
			"ids":    `[alicloud_hbr_ecs_backup_client.default.id]`,
			"status": `"UNKNOWN"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrEcsBackupClientSourceConfig(rand, map[string]string{
			"ids":          `[alicloud_hbr_ecs_backup_client.default.id]`,
			"instance_ids": `[alicloud_hbr_ecs_backup_client.default.instance_id]`,
			"status":       `"ACTIVATED"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrEcsBackupClientSourceConfig(rand, map[string]string{
			"ids":          `[alicloud_hbr_ecs_backup_client.default.id]`,
			"instance_ids": `[alicloud_hbr_ecs_backup_client.default.instance_id]`,
			"status":       `"UNKNOWN"`,
		}),
	}

	HbrEcsBackupClientCheckInfo.dataSourceTestCheck(t, rand, ecsBackupIdsconf, ecsInstanceIdconf, statusconf, allConf)
}

func testAccCheckAlicloudHbrEcsBackupClientSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc%d"
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
  name_regex = "^default-NODELETING$"
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
resource "alicloud_hbr_ecs_backup_client" "default" {
  instance_id = alicloud_instance.default.id
}
data "alicloud_hbr_ecs_backup_clients" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

var existHbrEcsBackupClientMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"clients.#":             "1",
		"clients.0.id":          CHECKSET,
		"clients.0.instance_id": CHECKSET,
	}
}

var fakeHbrEcsBackupClientMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"clients.#": "0",
	}
}

var HbrEcsBackupClientCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_hbr_ecs_backup_clients.default",
	existMapFunc: existHbrEcsBackupClientMapFunc,
	fakeMapFunc:  fakeHbrEcsBackupClientMapFunc,
}
