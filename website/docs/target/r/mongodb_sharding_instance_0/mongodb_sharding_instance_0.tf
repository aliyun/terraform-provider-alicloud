variable "name" {
  default = "tf-example"
}

variable "shard" {
  default = {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
}

variable "mongo" {
  default = {
    node_class = "dds.mongos.mid"
  }
}

data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.0.0/24"
  zone_id    = data.alicloud_zones.default.zones[0].id
  name       = var.name
}

resource "alicloud_mongodb_sharding_instance" "foo" {
  zone_id        = data.alicloud_zones.default.zones[0].id
  vswitch_id     = alicloud_vswitch.default.id
  engine_version = "3.4"
  name           = var.name
  dynamic "shard_list" {
    for_each = [var.shard]
    content {
      # TF-UPGRADE-TODO: The automatic upgrade tool can't predict
      # which keys might be set in maps assigned here, so it has
      # produced a comprehensive set here. Consider simplifying
      # this after confirming which keys can be set in practice.

      node_class   = shard_list.value.node_class
      node_storage = shard_list.value.node_storage
    }
  }
  dynamic "shard_list" {
    for_each = [var.shard]
    content {
      # TF-UPGRADE-TODO: The automatic upgrade tool can't predict
      # which keys might be set in maps assigned here, so it has
      # produced a comprehensive set here. Consider simplifying
      # this after confirming which keys can be set in practice.

      node_class   = shard_list.value.node_class
      node_storage = shard_list.value.node_storage
    }
  }
  dynamic "mongo_list" {
    for_each = [var.mongo]
    content {
      # TF-UPGRADE-TODO: The automatic upgrade tool can't predict
      # which keys might be set in maps assigned here, so it has
      # produced a comprehensive set here. Consider simplifying
      # this after confirming which keys can be set in practice.

      node_class = mongo_list.value.node_class
    }
  }
  dynamic "mongo_list" {
    for_each = [var.mongo]
    content {
      # TF-UPGRADE-TODO: The automatic upgrade tool can't predict
      # which keys might be set in maps assigned here, so it has
      # produced a comprehensive set here. Consider simplifying
      # this after confirming which keys can be set in practice.

      node_class = mongo_list.value.node_class
    }
  }
}
