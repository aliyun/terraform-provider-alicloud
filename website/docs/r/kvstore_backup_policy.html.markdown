---
subcategory: "Redis And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_backup_policy"
sidebar_current: "docs-alicloud-resource-kvstore-backup-policy"
description: |-
  Provides a backup policy for ApsaraDB Redis / Memcache instance resource.
---

# alicloud\_kvstore\_backup\_policy

-> **DEPRECATED:**  This resource  has been deprecated from version `1.104.0`. Please use resource [alicloud_kvstore_instance](https://www.terraform.io/docs/providers/alicloud/r/kvstore_instance).

Provides a backup policy for ApsaraDB Redis / Memcache instance resource. 

## Example Usage

Basic Usage

```
variable "creation" {
  default = "KVStore"
}

variable "multi_az" {
  default = "false"
}

variable "name" {
  default = "kvstorebackuppolicyvpc"
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

resource "alicloud_kvstore_instance" "default" {
  instance_class = "Memcache"
  instance_name  = var.name
  vswitch_id     = alicloud_vswitch.default.id
  private_ip     = "172.16.0.10"
  security_ips   = ["10.0.0.1"]
  instance_type  = "memcache.master.small.default"
  engine_version = "2.8"
}

resource "alicloud_kvstore_backup_policy" "default" {
  instance_id   = alicloud_kvstore_instance.default.id
  backup_period = ["Tuesday", "Wednesday"]
  backup_time   = "10:00Z-11:00Z"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The id of ApsaraDB for Redis or Memcache intance.
* `backup_time` - (Optional) Backup time, in the format of HH:mmZ- HH:mm Z
* `backup_period` - (Optional) Backup Cycle. Allowed values: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday

## Attributes Reference

The following attributes are exported:

* `id` - The id of the backup policy.
* `instance_id` - The id of ApsaraDB for Redis or Memcache intance.
* `backup_time` - Backup time, in the format of HH:mmZ- HH:mm Z
* `backup_period` - Backup Cycle. Allowed values: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday

## Import

KVStore backup policy can be imported using the id, e.g.

```
$ terraform import alicloud_kvstore_backup_policy.example r-abc12345678
```
