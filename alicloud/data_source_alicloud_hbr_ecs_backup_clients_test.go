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
data "alicloud_instances" "default" {
  name_regex = "hbr-ecs-backup-plan"
  status     = "Running"
}
resource "alicloud_hbr_ecs_backup_client" "default" {
  instance_id = data.alicloud_instances.default.instances.0.id
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
