---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_replica_set_role"
description: |-
  Provides an Alicloud MongoDB replica set role resource.
---

# alicloud_mongodb_replica_set_role

Provides an Alicloud MongoDB replica set role resource to modify the connection string of the replica set.

For information about how to modify connection string of MongoDB, see [Modify Connection String](https://alibabacloud.com/help/en/mongodb/user-guide/change-the-endpoint-and-port-of-an-instance).

-> **NOTE:** Available since v1.248.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_mongodb_replica_set_role&exampleId=0597f55b-14cb-17fd-804d-1044e9bd084990f8495a&activeTab=example&spm=docs.r.mongodb_replica_set_role.0.0597f55b14&intl_lang=EN_US" target="_blank">
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

# modify private network address.
resource "alicloud_mongodb_replica_set_role" "private" {
  db_instance_id    = alicloud_mongodb_instance.default.id
  role_id           = alicloud_mongodb_instance.default.replica_sets[0].role_id
  connection_prefix = "test-tf-private-change"
  connection_port   = 3718
  network_type      = "VPC"
}

# modify public network address.
resource "alicloud_mongodb_replica_set_role" "public" {
  db_instance_id    = alicloud_mongodb_instance.default.id
  role_id           = alicloud_mongodb_public_network_address.default.replica_sets[0].role_id
  connection_prefix = "test-tf-public-0"
  connection_port   = 3719
  network_type      = "Public"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_mongodb_replica_set_role&spm=docs.r.mongodb_replica_set_role.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `db_instance_id` - (Required, ForceNew) The instance ID.
* `role_id` - (Required, ForceNew) The role ID in the replica set.
* `connection_prefix` - (Optional) The prefix of the connection string, will be computed if not specified.
* `connection_port` - (Optional) The port of the connection string, will be computed if not specified.`
* `network_type` - (Required, ForceNew) The network type of the connection string. Valid values:
    - `VPC`: private network address.
    - `Public`: public network address.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. Composed of instance ID, network type and role ID with format `<db_instance_id>:<network_type>:<role_id>`.
* `replica_set_role` - The role of the related connection string.
* `connection_domain` - The connection address of the role.


## Import

MongoDB replica set role can be imported using the id, e.g. Composed of instance ID, network type and role ID with format `<db_instance_id>:<network_type>:<role_id>`.

```shell
$ terraform import alicloud_mongodb_replica_set_role.example <id>
```
