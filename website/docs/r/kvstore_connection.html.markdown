---
subcategory: "Redis And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_connection"
sidebar_current: "docs-alicloud-resource-kvstore-connection"
description: |-
  Provides a kvstore instance connection resource.
---

# alicloud\_kvstore\_connection

Provides a kvstore connection resource to allocate an Internet connection string for kvstore instance.

-> **NOTE:** (Available in 1.70.1+)Each kvstore instance will allocate a intranet connnection string automatically and its prifix is kvstore instance ID.
 To avoid unnecessary conflict, please specified a internet connection prefix before applying the resource.

## Example Usage

```
variable "creation" {
  default = "KVStore"
}
variable "name" {
  default = "kvstoreinstancevpc"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "ap-southeast-1b"
  name              = "${var.name}"
}
resource "alicloud_kvstore_instance" "default" {
  instance_class = "redis.master.small.default"
  instance_name  = "${var.name}"
  availability_zone = "ap-southeast-1b"
  private_ip     = "172.16.0.10"
  security_ips   = ["10.0.0.1"]
  instance_type  = "Redis"
  engine_version = "4.0"
}
resource "alicloud_kvstore_connection" "default" {
  instance_id = "${alicloud_kvstore_instance.default.id}"
  connection_string_prefix = "redisterraform"
  port = "3308"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance that can run database.
* `connection_prefix` - (ForceNew) Prefix of an Internet connection string. It must be checked for uniqueness. It may consist of lowercase letters, numbers, and underlines, and must start with a letter and have no more than 30 characters. Default to <instance_id> + 'tf'.
* `port` - (Optional) Internet connection port. Valid value: [3001-3999]. Default to 3306.

## Attributes Reference

The following attributes are exported:

* `id` - The current instance connection resource ID. Composed of instance ID and connection string with format `<instance_id>:<connection_prefix>`.
* `connection_prefix` - Prefix of a connection string.
* `port` - Connection instance port.
* `connection_string` - Connection instance string.

## Import

kvstore connection can be imported using the id, e.g.

```
$ terraform import alicloud_kvstore_connection.example abc12345678
```
