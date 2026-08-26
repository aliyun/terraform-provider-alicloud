// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMilvusInstanceDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMilvusInstanceSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_milvus_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMilvusInstanceSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_milvus_instance.default.id}_fake"]`,
		}),
	}

	ResourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMilvusInstanceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_milvus_instance.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		fakeConfig: testAccCheckAlicloudMilvusInstanceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_milvus_instance.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}
	InstanceNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMilvusInstanceSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_milvus_instance.default.id}"]`,
			"instance_name": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMilvusInstanceSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_milvus_instance.default.id}_fake"]`,
			"instance_name": `"${var.name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMilvusInstanceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_milvus_instance.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
			"instance_name":     `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMilvusInstanceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_milvus_instance.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
			"instance_name":     `"${var.name}_fake"`,
		}),
	}

	MilvusInstanceCheckInfo.dataSourceTestCheck(t, rand, idsConf, ResourceGroupIdConf, InstanceNameConf, allConf)
}

var existMilvusInstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#":                      "1",
		"instances.0.resource_group_id":    CHECKSET,
		"instances.0.configuration":        CHECKSET,
		"instances.0.encrypted":            CHECKSET,
		"instances.0.components.#":         CHECKSET,
		"instances.0.db_version":           CHECKSET,
		"instances.0.ha":                   CHECKSET,
		"instances.0.payment_type":         CHECKSET,
		"instances.0.auto_backup":          CHECKSET,
		"instances.0.tags.%":               CHECKSET,
		"instances.0.status":               CHECKSET,
		"instances.0.zone_id":              CHECKSET,
		"instances.0.instance_id":          CHECKSET,
		"instances.0.vswitch_ids.#":        CHECKSET,
		"instances.0.create_time":          CHECKSET,
		"instances.0.order_id":             CHECKSET,
		"instances.0.security_group_ids.#": CHECKSET,
		"instances.0.instance_name":        CHECKSET,
		"instances.0.vpc_id":               CHECKSET,
		"instances.0.region_id":            CHECKSET,
		"instances.0.expire_time":          CHECKSET,
		"instances.0.multi_zone_mode":      CHECKSET,
		"instances.0.running_time":         CHECKSET,
	}
}

var fakeMilvusInstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#": "0",
	}
}

var MilvusInstanceCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_milvus_instances.default",
	existMapFunc: existMilvusInstanceMapFunc,
	fakeMapFunc:  fakeMilvusInstanceMapFunc,
}

func testAccCheckAlicloudMilvusInstanceSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	password := testAccMilvusInstancePassword(rand)
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccMilvusInstance%d"
}

variable "zone_id" {
  default = "cn-hangzhou-j"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultILXuit" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "defaultN80M7S" {
  vpc_id       = alicloud_vpc.defaultILXuit.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.1.0/24"
  vswitch_name = "milvus-test"
}

resource "alicloud_milvus_instance" "default" {
  zone_id               = var.zone_id
  vswitch_ids {
    vsw_id  = alicloud_vswitch.defaultN80M7S.id
    zone_id = alicloud_vswitch.defaultN80M7S.zone_id
  }
  db_admin_password     = "%s"
  components {
    type           = "data"
    cu_num         = 2
    replica        = 1
    cu_type        = "general"
    disk_size_type = "Normal"
  }
  components {
    type           = "index"
    cu_num         = 4
    replica        = 2
    cu_type        = "general"
    disk_size_type = "Normal"
  }
  components {
    type           = "query"
    cu_num         = 4
    replica        = 2
    cu_type        = "general"
    disk_size_type = "Normal"
  }
  components {
    type           = "proxy"
    cu_num         = 2
    replica        = 2
    cu_type        = "general"
    disk_size_type = "Normal"
  }
  components {
    type           = "mix_coordinator"
    cu_num         = 4
    replica        = 2
    cu_type        = "general"
    disk_size_type = "Normal"
  }
  instance_name         = var.name
  db_version            = "2.4"
  vpc_id                = alicloud_vpc.defaultILXuit.id
  ha                    = false
  payment_type          = "Subscription"
  multi_zone_mode       = "Single"
  payment_duration_unit = "year"
  payment_duration      = 1
  auto_pay              = true
  resource_group_id     = data.alicloud_resource_manager_resource_groups.default.ids.0
  configuration         = "rootCoord:\n    maxDatabaseNum: 64"
}

data "alicloud_milvus_instances" "default" {
%s
}
`, rand, password, strings.Join(pairs, "\n   "))
	return config
}
