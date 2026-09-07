---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_backup"
description: |-
  Provides a Alicloud Mongodb Backup resource.
---

# alicloud_mongodb_backup

Provides a Mongodb Backup resource.

Instance-level or database-level backup objects.

For information about Mongodb Backup and how to use it, see [What is Backup](https://next.api.alibabacloud.com/document/Dds/2015-12-01/CreateBackup).

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}

variable "zone_id" {
  default = "cn-shanghai-b"
}

variable "cidr_block" {
  default = "10.0.0.0/24"
}

resource "alicloud_vpc" "default" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = "bgg-vpc-shanghai-b"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = var.zone_id
  cidr_block = var.cidr_block
}

resource "alicloud_mongodb_instance" "default" {
  engine_version      = "5.0"
  storage_type        = "cloud_essd1"
  vswitch_id          = alicloud_vswitch.default.id
  db_instance_storage = "20"
  vpc_id              = alicloud_vpc.default.id
  db_instance_class   = "mdb.shard.4x.large.d"
  storage_engine      = "WiredTiger"
  network_type        = "VPC"
  zone_id             = var.zone_id
  replication_factor  = "3"
  readonly_replicas   = "0"
}

resource "alicloud_mongodb_backup" "default" {
  backup_method           = "Snapshot"
  db_instance_id          = alicloud_mongodb_instance.default.id
  backup_retention_period = 7
}
```

## Argument Reference

The following arguments are supported:
* `backup_method` - (Optional, ForceNew) Backup Method. Valid values: `Snapshot`, `Physical`, `Logical`. Default value: `Physical`.
* `backup_retention_period` - (Optional, ForceNew, Int) Backup retention days. Valid values: `7` - `730`, `-1` (Long-term retention). Not passing means consistent with the default backup policy.
* `db_instance_id` - (Required, ForceNew) The ID of the MongoDB instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<db_instance_id>:<backup_id>`.
* `backup_db_names` - Backup DB Names.
* `backup_download_url` - Backup Download URL.
* `backup_id` - Backup Id.
* `backup_intranet_download_url` - Backup Intranet DownloadURL.
* `backup_job_id` - The backup task ID.
* `backup_mode` - Backup Mode. Valid values: `Automated`, `Manual`.
* `backup_size` - Backup Size in Bytes.
* `backup_start_time` - The start time of this backup, in the format of `yyyy-MM-ddTHH:mm:ssZ` (UTC time).
* `backup_end_time` - The end time of this backup, in the format of `yyyy-MM-ddTHH:mm:ssZ` (UTC time).
* `backup_type` - Backup Type. Valid values: `FullBackup`, `IncrementalBackup`.
* `status` - The status of the resource. Valid values: `Success`, `Failed`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 23 mins) Used when create the Backup.
* `delete` - (Defaults to 5 mins) Used when delete the Backup.

## Import

Mongodb Backup can be imported using the id, e.g.

```shell
$ terraform import alicloud_mongodb_backup.example <db_instance_id>:<backup_id>
```
