package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMongodbShardingNetworkPublicAddressesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.MongoDBSupportRegions)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongodbShardingNetworkPublicAddressesDataSourceName(rand, map[string]string{
			"node_id": `"${alicloud_mongodb_sharding_network_public_address.default.node_id}"`,
			"role":    `"Primary"`,
		}),
		fakeConfig: "",
	}
	var existAlicloudMongodbShardingNetworkPublicAddressesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"addresses.#":                 "1",
			"addresses.0.db_instance_id":  CHECKSET,
			"addresses.0.node_id":         CHECKSET,
			"addresses.0.expired_time":    "",
			"addresses.0.ip_address":      CHECKSET,
			"addresses.0.network_address": CHECKSET,
			"addresses.0.network_type":    CHECKSET,
			"addresses.0.node_type":       "",
			"addresses.0.port":            CHECKSET,
			"addresses.0.role":            CHECKSET,
			"addresses.0.vpc_id":          "",
			"addresses.0.vswitch_id":      "",
		}
	}
	var fakeAlicloudMongodbShardingNetworkPublicAddressesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"addresses.#": "0",
		}
	}
	var alicloudMongodbShardingNetworkPublicAddressesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mongodb_sharding_network_public_addresses.default",
		existMapFunc: existAlicloudMongodbShardingNetworkPublicAddressesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudMongodbShardingNetworkPublicAddressesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudMongodbShardingNetworkPublicAddressesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
}
func testAccCheckAlicloudMongodbShardingNetworkPublicAddressesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccShardingNetworkPublicAddress-%d"
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

resource "alicloud_mongodb_sharding_network_public_address" "default" {
 db_instance_id = alicloud_mongodb_sharding_instance.default.id
 node_id = alicloud_mongodb_sharding_instance.default.mongo_list.0.node_id
}

data "alicloud_mongodb_sharding_network_public_addresses" "default" {
	db_instance_id = alicloud_mongodb_sharding_network_public_address.default.db_instance_id
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
