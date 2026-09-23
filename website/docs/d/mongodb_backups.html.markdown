---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_backups"
sidebar_current: "docs-alicloud-datasource-mongodb-backups"
description: |-
  Provides a list of Mongodb Backup owned by an Alibaba Cloud account.
---

# alicloud_mongodb_backups

This data source provides Mongodb Backup available to the user. [What is Backup](https://next.api.alibabacloud.com/document/Dds/2015-12-01/CreateBackup)

-> **NOTE:** Available since v1.292.0.

## Example Usage

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

data "alicloud_mongodb_backups" "default" {
  ids            = [alicloud_mongodb_backup.default.id]
  db_instance_id = alicloud_mongodb_instance.default.id
}

output "alicloud_mongodb_backup_example_id" {
  value = data.alicloud_mongodb_backups.default.backups.0.id
}
```

## Argument Reference

The following arguments are supported:
* `backup_id` - (Optional) Backup Id.
* `db_instance_id` - (Required) The ID of the MongoDB instance.
* `end_time` - (Optional) The end time of the backup query range, format: `yyyy-MM-ddTHH:mmZ` (UTC time).
* `start_time` - (Optional) The start time of the backup query range, format: `yyyy-MM-ddTHH:mmZ` (UTC time).
* `ids` - (Optional, Computed) A list of Backup IDs. The value is formulated as `<db_instance_id>:<backup_id>`.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Backup IDs.
* `backups` - A list of Backup Entries. Each element contains the following attributes:
  * `backup_db_names` - Backup DB Names.
  * `backup_download_url` - Backup Download URL.
  * `backup_id` - Backup Id.
  * `backup_intranet_download_url` - Backup Intranet DownloadURL.
  * `backup_job_id` - **NOTE:** This field is only available when `enable_details` is `true`. The backup task ID.
  * `backup_method` - Backup Method.
  * `backup_mode` - Backup Mode.
  * `backup_size` - Backup Size.
  * `backup_start_time` - The start time of this backup, in the format of `yyyy-MM-ddTHH:mm:ssZ` (UTC time).
  * `backup_end_time` - The end time of this backup, in the format of `yyyy-MM-ddTHH:mm:ssZ` (UTC time).
  * `backup_type` - Backup Type.
  * `status` - The status of the resource.
  * `id` - The ID of the resource supplied above.
