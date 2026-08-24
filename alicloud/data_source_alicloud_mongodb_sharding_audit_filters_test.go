package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudMongodbShardingAuditFiltersDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.MongoDBSupportRegions)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongodbShardingAuditFiltersDataSourceName(rand, map[string]string{}),
		fakeConfig:  "",
	}
	var existAlicloudMongodbShardingAuditFiltersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"filters.#":                "1",
			"filters.0.id":             CHECKSET,
			"filters.0.db_instance_id": CHECKSET,
			"filters.0.audit_status":   "enable",
			"filters.0.filter":         "admin,slow",
			"filters.0.service_type":   CHECKSET,
			"filters.0.storage_period": CHECKSET,
		}
	}
	var fakeAlicloudMongodbShardingAuditFiltersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"filters.#": "0",
		}
	}
	var alicloudMongodbShardingAuditFiltersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mongodb_sharding_audit_filters.default",
		existMapFunc: existAlicloudMongodbShardingAuditFiltersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudMongodbShardingAuditFiltersDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudMongodbShardingAuditFiltersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
}

func testAccCheckAlicloudMongodbShardingAuditFiltersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccShardingAuditFilter-%d"
}

data "alicloud_mongodb_zones" "default" {}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
}

resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = data.alicloud_mongodb_zones.default.zones.0.id
  vswitch_id     = data.alicloud_vswitches.default.ids[0]
  engine_version = "4.2"
  name           = var.name
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
}

resource "alicloud_mongodb_sharding_audit_filter" "default" {
  db_instance_id     = alicloud_mongodb_sharding_instance.default.id
  audit_status       = "enable"
  filter             = "admin,slow"
  service_type       = "V2_Standard"
  storage_period     = 30
  hot_storage_period = 7
}

data "alicloud_mongodb_sharding_audit_filters" "default" {
  db_instance_id = alicloud_mongodb_sharding_audit_filter.default.db_instance_id
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
