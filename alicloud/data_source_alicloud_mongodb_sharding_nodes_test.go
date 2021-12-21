package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMongodbShardingNodesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.MongoDBSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongodbShardingNodesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mongodb_sharding_node.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMongodbShardingNodesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mongodb_sharding_node.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongodbShardingNodesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_mongodb_sharding_node.default.id}"]`,
			"status": `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudMongodbShardingNodesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_mongodb_sharding_node.default.id}_fake"]`,
			"status": `"Deleting"`,
		}),
	}
	var existAlicloudMongodbShardingNodesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "1",
			"nodes.#":                "1",
			"nodes.0.db_instance_id": CHECKSET,
			"nodes.0.node_class":     "dds.shard.mid",
			"nodes.0.node_id":        CHECKSET,
			"nodes.0.id":             CHECKSET,
			"nodes.0.node_storage":   "10",
			"nodes.0.status":         "Running",
		}
	}
	var fakeAlicloudMongodbShardingNodesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudMongodbShardingNodesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mongodb_sharding_nodes.default",
		existMapFunc: existAlicloudMongodbShardingNodesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudMongodbShardingNodesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudMongodbShardingNodesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf)
}
func testAccCheckAlicloudMongodbShardingNodesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccShardingNode-%d"
}

data "alicloud_mongodb_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_mongodb_zones.default.zones.0.id
  vswitch_name = "subnet-for-local-test"
}

resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = data.alicloud_mongodb_zones.default.zones.0.id
  vswitch_id     = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  engine_version = "3.4"
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

resource "alicloud_mongodb_sharding_node" "default" {
	db_instance_id = alicloud_mongodb_sharding_instance.default.id
	node_class = "dds.shard.mid"
	node_storage = "10"
}

data "alicloud_mongodb_sharding_nodes" "default" {	
	db_instance_id = alicloud_mongodb_sharding_instance.default.id
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
