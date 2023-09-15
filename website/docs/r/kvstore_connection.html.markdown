---
subcategory: "Redis And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_connection"
sidebar_current: "docs-alicloud-resource-kvstore-connection"
description: |-
  Operate the public network ip of the specified resource.
---

# alicloud_kvstore_connection

Operate the public network ip of the specified resource. How to use it, see [What is Resource Alicloud KVStore Connection](https://www.alibabacloud.com/help/doc-detail/125795.htm).

-> **NOTE:** Available since v1.101.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_kvstore_zones" "default" {

}
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_kvstore_zones.default.zones.0.id
}

resource "alicloud_kvstore_instance" "default" {
  db_instance_name  = var.name
  vswitch_id        = alicloud_vswitch.default.id
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  zone_id           = data.alicloud_kvstore_zones.default.zones.0.id
  instance_class    = "redis.master.large.default"
  instance_type     = "Redis"
  engine_version    = "5.0"
  security_ips      = ["10.23.12.24"]
  config = {
    appendonly             = "yes"
    lazyfree-lazy-eviction = "yes"
  }
  tags = {
    Created = "TF",
    For     = "example",
  }
}

resource "alicloud_kvstore_connection" "default" {
  connection_string_prefix = "exampleconnection"
  instance_id              = alicloud_kvstore_instance.default.id
  port                     = "6370"
}
```

## Argument Reference

The following arguments are supported:
* `connection_string_prefix` - (Required) The prefix of the public endpoint. The prefix can be 8 to 64 characters in length, and can contain lowercase letters and digits. It must start with a lowercase letter.
* `instance_id`- (Required, ForceNew) The ID of the instance.
* `port` - (Required) The service port number of the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of KVStore DBInstance.
* `connection_string` - The public connection string of KVStore DBInstance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the KVStore connection (until it reaches the initial `Normal` status). 
* `update` - (Defaults to 5 mins) Used when updating the KVStore connection (until it reaches the initial `Normal` status). 
* `delete` - (Defaults to 30 mins) Used when deleting the KVStore connection (until it reaches the initial `Normal` status). 

## Import

KVStore connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_kvstore_connection.example r-abc12345678
```

