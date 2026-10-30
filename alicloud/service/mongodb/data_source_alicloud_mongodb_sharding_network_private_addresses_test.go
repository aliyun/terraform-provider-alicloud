package mongodb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudMongodbShardingNetworkPrivateAddressesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.MongoDBSupportRegions)
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMongodbShardingNetworkPrivateAddressesDataSourceName(rand, map[string]string{
			"node_id": `"${alicloud_mongodb_sharding_network_private_address.default.node_id}"`,
			"role":    `"Primary"`,
		}),
		FakeConfig: "",
	}
	var existAlicloudMongodbShardingNetworkPrivateAddressesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"addresses.#":                 "1",
			"addresses.0.db_instance_id":  acctest.CHECKSET,
			"addresses.0.node_id":         acctest.CHECKSET,
			"addresses.0.expired_time":    "",
			"addresses.0.ip_address":      acctest.CHECKSET,
			"addresses.0.network_address": acctest.CHECKSET,
			"addresses.0.network_type":    acctest.CHECKSET,
			"addresses.0.node_type":       "",
			"addresses.0.port":            acctest.CHECKSET,
			"addresses.0.role":            acctest.CHECKSET,
			"addresses.0.vpc_id":          acctest.CHECKSET,
			"addresses.0.vswitch_id":      acctest.CHECKSET,
		}
	}
	var fakeAlicloudMongodbShardingNetworkPrivateAddressesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"addresses.#": "0",
		}
	}
	var alicloudMongodbShardingNetworkPrivateAddressesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_mongodb_sharding_network_private_addresses.default",
		ExistMapFunc: existAlicloudMongodbShardingNetworkPrivateAddressesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudMongodbShardingNetworkPrivateAddressesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
	}
	alicloudMongodbShardingNetworkPrivateAddressesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
}
func testAccCheckAlicloudMongodbShardingNetworkPrivateAddressesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccShardingNetworkPrivateAddress-%d"
	}

	data "alicloud_mongodb_zones" "default" {
	}

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

	resource "alicloud_mongodb_sharding_network_private_address" "default" {
  		account_name     = "tf_test"
  		account_password = "YourPassword+12345"
  		db_instance_id   = alicloud_mongodb_sharding_instance.default.id
  		node_id          = alicloud_mongodb_sharding_instance.default.shard_list.0.node_id
  		zone_id          = alicloud_mongodb_sharding_instance.default.zone_id
	}

	data "alicloud_mongodb_sharding_network_private_addresses" "default" {
  		db_instance_id = alicloud_mongodb_sharding_instance.default.id
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
