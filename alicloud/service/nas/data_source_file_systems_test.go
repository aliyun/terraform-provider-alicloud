package nas_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlicloudNASFileSystem_DataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100000, 999999)
	storageTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"storage_type":      `"${alicloud_nas_file_system.default.storage_type}"`,
			"description_regex": `"^${alicloud_nas_file_system.default.description}"`,
		}),
	}
	protocolTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"protocol_type":     `"${alicloud_nas_file_system.default.protocol_type}"`,
			"description_regex": `"^${alicloud_nas_file_system.default.description}"`,
		}),
	}
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"description_regex": `"^${alicloud_nas_file_system.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"description_regex": `"^${alicloud_nas_file_system.default.description}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_nas_file_system.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_nas_file_system.default.id}_fake"]`,
		}),
	}
	fileSystemTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"file_system_type":  `"standard"`,
			"description_regex": `"^${alicloud_nas_file_system.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"file_system_type":  `"extreme"`,
			"description_regex": `"^${alicloud_nas_file_system.default.description}"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"storage_type":      `"${alicloud_nas_file_system.default.storage_type}"`,
			"protocol_type":     `"${alicloud_nas_file_system.default.protocol_type}"`,
			"description_regex": `"^${alicloud_nas_file_system.default.description}"`,
			"ids":               `["${alicloud_nas_file_system.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudFileSystemDataSourceConfig(rand, map[string]string{
			"description_regex": `"^${alicloud_nas_file_system.default.description}_fake"`,
			"ids":               `["${alicloud_nas_file_system.default.id}_fake"]`,
		}),
	}

	fileSystemCheckInfo.dataSourceTestCheck(t, rand, storageTypeConf, protocolTypeConf,
		descriptionConf, idsConf, fileSystemTypeConf, allConf)
}

func testAccCheckAlicloudFileSystemDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "description" {
  default = "tf-testAccCheckAlicloudFileSystemsDataSource"
}
variable "storage_type" {
  default = "Capacity"
}
data "alicloud_nas_protocols" "default" {
        type = "${var.storage_type}"
}
resource "alicloud_nas_file_system" "default" {
  description = "${var.description}"
  storage_type = "${var.storage_type}"
  protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
}
data "alicloud_nas_file_systems" "default" {
	%s
}`, strings.Join(pairs, "\n  "))
	return config
}

var existFileSystemMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"systems.0.id":            CHECKSET,
		"systems.0.region_id":     CHECKSET,
		"systems.0.description":   "tf-testAccCheckAlicloudFileSystemsDataSource",
		"systems.0.protocol_type": CHECKSET,
		"systems.0.storage_type":  "Capacity",
		"systems.0.metered_size":  CHECKSET,
		"systems.0.create_time":   CHECKSET,
		"ids.#":                   "1",
		"ids.0":                   CHECKSET,
		"descriptions.#":          "1",
		"descriptions.0":          "tf-testAccCheckAlicloudFileSystemsDataSource",
	}
}

var fakeFileSystemMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"systems.#":      "0",
		"ids.#":          "0",
		"descriptions.#": "0",
	}
}

var fileSystemCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_nas_file_systems.default",
	existMapFunc: existFileSystemMapCheck,
	fakeMapFunc:  fakeFileSystemMapCheck,
}

// Case data source output for auto upgrade config on cpfsse file systems.
func TestAccAliCloudNasFileSystem_DataSource_AutoUpgradeConfig(t *testing.T) {
	rand := tfacctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf-testAccNasAutoUpgradeDs%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudFileSystemDataSourceAutoUpgradeConfigConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.alicloud_nas_file_systems.default", "systems.0.id", "alicloud_nas_file_system.default", "id"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.default", "systems.0.auto_upgrade_config.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.default", "systems.0.auto_upgrade_config.0.enabled", "true"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.default", "systems.0.auto_upgrade_config.0.capacity_used_ratio", "60"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.default", "systems.0.auto_upgrade_config.0.time", "90"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.default", "systems.0.auto_upgrade_config.0.step", "200"),
				),
			},
		},
	})
}

func testAccCheckAlicloudFileSystemDataSourceAutoUpgradeConfigConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

# This fixture pins three specific cn-beijing zones and the advance_100 CPFS SE
# specification: ZRS needs three distinct zones, and cpfsse advance_100 capacity
# is per-zone inventory. A region-wide zone query only proves API support, not
# inventory, so a dynamically selected triple can fail at apply time with
# Resource.OutOfStock. Update the fixed triple when these zones run out of stock.
resource "alicloud_vpc" "default" {
  is_default  = false
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "tf-testAccNasAutoUpgradeDsVpc"
  enable_ipv6 = true
}

resource "alicloud_vswitch" "vswitch_f" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-i"
  cidr_block   = "192.168.2.0/24"
  vswitch_name = "tf-testAccNasAutoUpgradeDsVswF"
}

resource "alicloud_vswitch" "vswitch_c" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-l"
  cidr_block   = "192.168.3.0/24"
  vswitch_name = "tf-testAccNasAutoUpgradeDsVswC"
}

resource "alicloud_vswitch" "vswitch_d" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-h"
  cidr_block   = "192.168.4.0/24"
  vswitch_name = "tf-testAccNasAutoUpgradeDsVswD"
}

resource "alicloud_nas_file_system" "default" {
  description      = var.name
  storage_type     = "advance_100"
  encrypt_type     = 0
  vpc_id           = alicloud_vpc.default.id
  redundancy_vswitch_ids = [
    alicloud_vswitch.vswitch_c.id,
    alicloud_vswitch.vswitch_d.id,
    alicloud_vswitch.vswitch_f.id,
  ]
  capacity         = 500
  protocol_type    = "cpfs"
  file_system_type = "cpfsse"
  redundancy_type  = "ZRS"

  auto_upgrade_config {
    enabled             = true
    capacity_used_ratio = 60
    time                = 90
    step                = 200
  }
}

data "alicloud_nas_file_systems" "default" {
  ids = [alicloud_nas_file_system.default.id]
}
`, name)
}
