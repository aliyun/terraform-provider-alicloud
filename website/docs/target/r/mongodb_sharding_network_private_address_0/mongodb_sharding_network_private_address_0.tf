data "alicloud_mongodb_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
}

resource "alicloud_mongodb_sharding_instance" "default" {
  zone_id        = data.alicloud_mongodb_zones.default.zones.0.id
  vswitch_id     = data.alicloud_vswitches.default.ids[0]
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
resource "alicloud_mongodb_sharding_network_private_address" "example" {
  db_instance_id   = alicloud_mongodb_sharding_instance.default.id
  node_id          = alicloud_mongodb_sharding_instance.default.shard_list.0.node_id
  zone_id          = alicloud_mongodb_sharding_instance.default.zone_id
  account_name     = "example_value"
  account_password = "YourPassword+12345"
}
