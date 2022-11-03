---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_connection"
sidebar_current: "docs-alicloud-resource-gpdb-connection"
description: |-
  Provides an AnalyticDB for PostgreSQL instance connection resource.
---

# alicloud\_gpdb\_connection

Provides a connection resource to allocate an Internet connection string for instance.

-> **NOTE:**  Available in 1.48.0+

-> **NOTE:** Each instance will allocate a intranet connection string automatically and its prefix is instance ID.
 To avoid unnecessary conflict, please specified a internet connection prefix before applying the resource.

## Example Usage

```terraform
variable "creation" {
  default = "Gpdb"
}

variable "name" {
  default = "gpdbConnectionBasic"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
  name              = var.name
}

resource "alicloud_gpdb_instance" "default" {
  vswitch_id           = alicloud_vswitch.default.id
  engine               = "gpdb"
  engine_version       = "4.3"
  instance_class       = "gpdb.group.segsdx2"
  instance_group_count = "2"
  description          = var.name
}

resource "alicloud_gpdb_connection" "default" {
  instance_id       = alicloud_gpdb_instance.default.id
  connection_prefix = "testAbc"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance that can run database.
* `connection_prefix` - (ForceNew) Prefix of an Internet connection string. It must be checked for uniqueness. It may consist of lowercase letters, numbers, and underlines, and must start with a letter and have no more than 30 characters. Default to <instance_id> + '-tf'.
* `port` - (Optional) Internet connection port. Valid value: [3200-3999]. Default to 3306.

### Timeouts

-> **NOTE:** Available in 1.53.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the Internet connection (until DB instance reaches the initial `Running` status). 
* `update` - (Defaults to 10 mins) Used when activating the DB instance during update.
* `delete` - (Defaults to 10 mins) Used when terminating the DB instance. 

## Attributes Reference

The following attributes are exported:

* `id` - The current instance connection resource ID. Composed of instance ID and connection string with format `<instance_id>:<connection_prefix>`.
* `connection_string` - Connection instance string.
* `ip_address` - The ip address of connection string.

## Import

AnalyticDB for PostgreSQL's connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_connection.example abc12345678
```
