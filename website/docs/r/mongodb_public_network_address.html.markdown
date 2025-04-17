---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_public_network_address"
description: |-
  Provides an Alicloud MongoDB public network address resource.
---

# alicloud_mongodb_public_network_address

Provides an Alicloud MongoDB public network address resource.

For information about MongoDB public network address and how to use it, see [Allocate Public Network Address for MongoDB](https://www.alibabacloud.com/help/en/mongodb/getting-started/apply-for-a-public-endpoint-for-an-apsaradb-for-mongodb-instance-optional).

-> **NOTE:** Available since v1.248.0.

## Example Usage


<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_mongodb_public_network_address&exampleId=65adeb2d-5510-4f42-102d-f310f9f6c98ecbeafc5c&activeTab=example&spm=docs.r.mongodb_public_network_address.0.65adeb2d55&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_mongodb_zones" "default" {}

locals {
  index   = length(data.alicloud_mongodb_zones.default.zones) - 1
  zone_id = data.alicloud_mongodb_zones.default.zones[local.index].id
}

resource "alicloud_vpc" "default" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = local.zone_id
  cidr_block = "10.0.0.0/24"
}

resource "alicloud_mongodb_instance" "default" {
  engine_version      = "4.4"
  storage_type        = "cloud_essd1"
  vswitch_id          = alicloud_vswitch.default.id
  db_instance_storage = "20"
  vpc_id              = alicloud_vpc.default.id
  db_instance_class   = "mdb.shard.4x.large.d"
  storage_engine      = "WiredTiger"
  network_type        = "VPC"
  zone_id             = local.zone_id
}

resource "alicloud_mongodb_public_network_address" "default" {
  db_instance_id = alicloud_mongodb_instance.default.id
}
```

## Argument Reference

The following arguments are supported:
* `db_instance_id` - (Required, ForceNew) The instance ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource. Equal to the `db_instance_id`.
* `replica_sets` - Replica set instance information.
 * `connection_port` - The connection port of the node.
 * `replica_set_role` - The role of the node.
 * `connection_domain` - The connection address of the node.
 * `network_type` - The network type, should be always "Public".
 * `role_id` - The id of the role.
 * `connection_type` - The connection type.

## Import

MongoDB public network address can be imported using the id, e.g.

```shell
$ terraform import alicloud_mongodb_public_network_address.example <id>
```