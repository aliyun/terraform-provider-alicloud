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

			"instance_name": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMilvusInstanceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_milvus_instance.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,

			"instance_name": `"${var.name}_fake"`,
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
		"instances.0.kms_key_id":           CHECKSET,
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
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccMilvusInstance%d"
}


resource "alicloud_milvus_instance" "default" {
  ai_function       = false
  zone_id           = "cn-hangzhou-j"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  vswitch_ids {
    zone_id = "cn-hangzhou-j"
    vsw_id  = "vsw-bp1pommb2vygb0kzvf8i6"
  }
  vswitch_ids {
    zone_id = "cn-hangzhou-k"
    vsw_id  = "vsw-bp1tomony773mb6nlabw9"
  }
  encrypted             = false
  auto_renew            = false
  payment_duration_unit = "month"
  auto_pay              = true
  load_replicas         = "2"
  payment_duration      = "1"
  db_admin_password     = "@1234Test"
  instance_name         = "tf-parity-sub-month-{{function.randomIntString(100000,999999)}}"
  components {
    cu_type = "general"
    type    = "streaming"
    cu_num  = "4"
    replica = "2"
  }
  components {
    cu_type = "general"
    type    = "data"
    cu_num  = "4"
    replica = "1"
  }
  components {
    cu_type        = "general"
    type           = "query"
    cu_num         = "16"
    disk_size_type = "Normal"
    replica        = "2"
  }
  components {
    cu_type = "general"
    type    = "proxy"
    cu_num  = "2"
    replica = "2"
  }
  components {
    cu_type = "general"
    type    = "mix_coordinator"
    cu_num  = "4"
    replica = "2"
  }
  db_version          = "2.6"
  vpc_id              = "vpc-bp168d0ay5yft9aira762"
  is_multi_az_storage = true
  payment_type        = "Subscription"
  ha                  = true
  multi_zone_mode     = "single"
  auto_backup         = true
  promotion_no        = "youhuiquan_promotion_option_id_for_blank"
}

data "alicloud_milvus_instances" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
