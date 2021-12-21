package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudMongoDBShardingNetworkPublicAddress_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_sharding_network_public_address.default"
	checkoutSupportedRegions(t, true, connectivity.MongoDBSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudMongoDBShardingNetworkPublicAddressMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongodbShardingNetworkPublicAddress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smongodbshardingnetworkpublicaddress%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMongoDBShardingNetworkPublicAddressBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id": "${alicloud_mongodb_sharding_instance.default.id}",
					"node_id":        "${alicloud_mongodb_sharding_instance.default.mongo_list.0.node_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":                    CHECKSET,
						"node_id":                           CHECKSET,
						"network_address.#":                 "1",
						"network_address.0.node_id":         CHECKSET,
						"network_address.0.expired_time":    "",
						"network_address.0.ip_address":      CHECKSET,
						"network_address.0.network_address": CHECKSET,
						"network_address.0.network_type":    CHECKSET,
						"network_address.0.node_type":       "",
						"network_address.0.port":            CHECKSET,
						"network_address.0.role":            CHECKSET,
						"network_address.0.vpc_id":          "",
						"network_address.0.vswitch_id":      "",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudMongoDBShardingNetworkPublicAddressMap0 = map[string]string{
	"node_id":        CHECKSET,
	"db_instance_id": CHECKSET,
}

func AlicloudMongoDBShardingNetworkPublicAddressBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
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
`, name)
}
