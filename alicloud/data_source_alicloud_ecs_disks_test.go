package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECSDisksDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_disk.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_disk.default.id}_fake"]`,
		}),
	}
	ZoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_ecs_disk.default.id}"]`,
			"zone_id": `"${alicloud_ecs_disk.default.zone_id}"`,
		}),
	}
	categoryConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_ecs_disk.default.id}"]`,
			"category": `"cloud_efficiency"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_ecs_disk.default.id}"]`,
			"category": `"cloud"`,
		}),
	}
	deleteAutoSnapshotConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_ecs_disk.default.id}"]`,
			"delete_auto_snapshot": `"true"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_ecs_disk.default.id}"]`,
			"delete_auto_snapshot": `"false"`,
		}),
	}
	deleteWithInstanceConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_ecs_disk.default.id}"]`,
			"delete_with_instance": `"true"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_ecs_disk.default.id}"]`,
			"delete_with_instance": `"false"`,
		}),
	}
	diskNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_ecs_disk.default.id}"]`,
			"disk_name": `"${alicloud_ecs_disk.default.disk_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_ecs_disk.default.id}"]`,
			"disk_name": `"${alicloud_ecs_disk.default.disk_name}_fake"`,
		}),
	}
	enableAutoSnapshotConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_ecs_disk.default.id}"]`,
			"enable_auto_snapshot": `"true"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_ecs_disk.default.id}"]`,
			"enable_auto_snapshot": `"false"`,
		}),
	}
	encryptedConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_ecs_disk.default.id}"]`,
			"encrypted": `"on"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_ecs_disk.default.id}"]`,
			"encrypted": `"off"`,
		}),
	}
	paymentTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_ecs_disk.default.id}"]`,
			"payment_type": `"PayAsYouGo"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_ecs_disk.default.id}"]`,
			"payment_type": `"Subscription"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_disk.default.id}"]`,
			"tags": `{
				"Created" = "TF"
				"Environment" = "Acceptance-test"
		}`,
		}),
		fakeConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_disk.default.id}"]`,
			"tags": `{
				"Created" = "TF-fake"
				"Environment" = "Acceptance-test"
			}`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecs_disk.default.disk_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecs_disk.default.disk_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecs_disk.default.id}"]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecs_disk.default.id}"]`,
			"status": `"Creating"`,
		}),
	}

	pagingConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"disk_name":   `"${alicloud_ecs_disk.default.disk_name}"`,
			"page_number": `1`,
		}),
		fakeConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"disk_name":   `"${alicloud_ecs_disk.default.disk_name}"`,
			"page_number": `2`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"zone_id":              `"${alicloud_ecs_disk.default.zone_id}"`,
			"category":             `"cloud_efficiency"`,
			"delete_auto_snapshot": `"true"`,
			"delete_with_instance": `"true"`,
			"disk_name":            `"${alicloud_ecs_disk.default.disk_name}"`,
			"enable_auto_snapshot": `"true"`,
			"encrypted":            `"on"`,
			"ids":                  `["${alicloud_ecs_disk.default.id}"]`,
			"name_regex":           `"${alicloud_ecs_disk.default.disk_name}"`,
			"payment_type":         `"PayAsYouGo"`,
			"status":               `"Available"`,
			"tags": `{
				"Created" = "TF"
				"Environment" = "Acceptance-test"
		}`,
			"page_number": `1`,
		}),
		fakeConfig: testAccCheckAlicloudEcsDisksDataSourceName(rand, map[string]string{
			"category":             `"cloud"`,
			"delete_auto_snapshot": `"false"`,
			"delete_with_instance": `"false"`,
			"disk_name":            `"${alicloud_ecs_disk.default.disk_name}_fake"`,
			"enable_auto_snapshot": `"false"`,
			"encrypted":            `"off"`,
			"ids":                  `["${alicloud_ecs_disk.default.id}_fake"]`,
			"name_regex":           `"${alicloud_ecs_disk.default.disk_name}_fake"`,
			"payment_type":         `"Subscription"`,
			"status":               `"Creating"`,
			"tags": `{
				"Created" = "TF-fake"
				"Environment" = "Acceptance-test"
			}`,
			"page_number": `2`,
		}),
	}
	var existAlicloudEcsDisksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"disks.#":                         "1",
			"total_count":                     CHECKSET,
			"disks.0.zone_id":                 CHECKSET,
			"disks.0.category":                `cloud_efficiency`,
			"disks.0.delete_auto_snapshot":    `true`,
			"disks.0.delete_with_instance":    `true`,
			"disks.0.description":             `Test For Terraform`,
			"disks.0.disk_name":               CHECKSET,
			"disks.0.enable_auto_snapshot":    `true`,
			"disks.0.encrypted":               "on",
			"disks.0.instance_id":             "",
			"disks.0.kms_key_id":              "",
			"disks.0.payment_type":            "PayAsYouGo",
			"disks.0.performance_level":       "",
			"disks.0.size":                    `500`,
			"disks.0.image_id":                "",
			"disks.0.device":                  "",
			"disks.0.auto_snapshot_policy_id": "",
			"disks.0.attached_time":           "",
			"disks.0.mount_instance_num":      "0",
			"disks.0.product_code":            "",
			"disks.0.snapshot_id":             "",
			"disks.0.tags.%":                  "2",
			"disks.0.tags.Created":            "TF",
			"disks.0.tags.Environment":        "Acceptance-test",
		}
	}
	var fakeAlicloudEcsDisksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"disks.#": "0",
		}
	}
	var alicloudEcsDisksCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecs_disks.default",
		existMapFunc: existAlicloudEcsDisksDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcsDisksDataSourceNameMapFunc,
	}
	alicloudEcsDisksCheckInfo.dataSourceTestCheck(t, rand, idsConf, ZoneIdConf, categoryConf, deleteAutoSnapshotConf, deleteWithInstanceConf, diskNameConf, enableAutoSnapshotConf, encryptedConf, paymentTypeConf, tagsConf, nameRegexConf, statusConf, pagingConf, allConf)
}
func testAccCheckAlicloudEcsDisksDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccDisk-%d"
}
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
resource "alicloud_ecs_disk" "default" {
	zone_id = "${data.alicloud_zones.default.zones.0.id}"
	category = "cloud_efficiency"
	delete_auto_snapshot = "true"
	delete_with_instance = "true"
	description = "Test For Terraform"
	disk_name = var.name
	enable_auto_snapshot = "true"
	encrypted = "true"
	size = "500"
  	tags = {
    	Created     = "TF"
    	Environment = "Acceptance-test"
  	}
}

data "alicloud_ecs_disks" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
