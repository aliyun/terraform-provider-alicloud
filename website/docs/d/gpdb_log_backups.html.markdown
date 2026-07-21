---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_log_backups"
sidebar_current: "docs-alicloud-datasource-gpdb-logbackups"
description: |-
  Provides a list of Gpdb Log Backup owned by an Alibaba Cloud account.
---

# alicloud_gpdb_log_backups

This data source provides Gpdb Logbackup available to the user.[What is Log Backup](https://www.alibabacloud.com/help/en/)

-> **NOTE:** Available since v1.231.0.

## Example Usage

```terraform

data "alicloud_gpdb_instances" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_gpdb_log_backups" "default" {
  start_time     = "2022-12-12T02:00Z"
  end_time       = "2024-12-12T02:00Z"
  db_instance_id = data.alicloud_gpdb_instances.default.ids.0
  ids            = ["${data.alicloud_gpdb_instances.default.ids.0}"]
}

output "alicloud_gpdb_logbackup_example_id" {
  value = data.alicloud_gpdb_log_backups.default.logbackups.0.db_instance_id
}
```

## Argument Reference

The following arguments are supported:
* `db_instance_id` - (Required, ForceNew) The ID of the Master node of the instance.
* `end_time` - (ForceNew, Optional) The query end time, which must be greater than the query start time. Format: yyyy-MM-ddTHH:mmZ(UTC time).
* `page_number` - (ForceNew, Optional) Current page number.
* `page_size` - (ForceNew, Optional) Number of records per page.
* `start_time` - (ForceNew, Optional) The query start time. Format: yyyy-MM-ddTHH:mmZ(UTC time).
* `ids` - (Optional, ForceNew, Computed) A list of Logbackup IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Logbackup IDs.
* `logbackups` - A list of Logbackup Entries. Each element contains the following attributes:
  * `db_instance_id` - The ID of the Master node of the instance.
  * `log_backup_id` - The first ID of the resource
  * `log_file_name` - Log file name (OSS path).
  * `log_file_size` - Size of the backup log file. Unit: Byte.
  * `log_time` - The log timestamp.
  * `record_total` - Total number of records.
  * `segment_name` - The node name.
