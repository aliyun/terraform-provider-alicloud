---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_snapshots"
sidebar_current: "docs-alicloud-datasource-snapshots"
description: |-
  Provides a data source to get a list of snapshot according to the specified filters.
---

# alicloud\_snapshots

-> **DEPRECATED:** This datasource has been renamed to [alicloud_ecs_snapshots](https://www.terraform.io/docs/providers/alicloud/d/ecs_snapshots) from version 1.120.0.

Use this data source to get a list of snapshot according to the specified filters in an Alibaba Cloud account.

For information about snapshot and how to use it, see [Snapshot](https://www.alibabacloud.com/help/doc-detail/25460.html).

-> **NOTE:**  Available in 1.40.0+.

## Example Usage

```
data "alicloud_snapshots" "snapshots" {
  ids        = ["s-123456890abcdef"]
  name_regex = "tf-testAcc-snapshot"
}
```

##  Argument Reference

The following arguments are supported:

* `instance_id` - (Optional) The specified instance ID.
* `disk_id` - (Optional) The specified disk ID.
* `encrypted` - (Optional) Queries the encrypted snapshots. Optional values: `true`: Encrypted snapshots. `false`: No encryption attribute limit. Default value: `false`.
* `ids` - (Optional)  A list of snapshot IDs.
* `name_regex` - (Optional) A regex string to filter results by snapshot name.
* `status` - (Optional) The specified snapshot status. Default value: `all`. Optional values:
  * progressing: The snapshots are being created.
  * accomplished: The snapshots are ready to use.
  * failed: The snapshot creation failed.
  * all: All status.
* `type` - (Optional) The snapshot category. Default value: `all`. Optional values:
  * auto: Auto snapshots.
  * user: Manual snapshots.
  * all: Auto and manual snapshots.
* `source_disk_type` - (Optional) The type of source disk:
  * System: The snapshots are created for system disks.
  * Data: The snapshots are created for data disks.
* `usage` - (Optional) The usage of the snapshot:
  * image: The snapshots are used to create custom images.
  * disk: The snapshots are used to CreateDisk.
  * mage_disk: The snapshots are used to create custom images and data disks.
  * none: The snapshots are not used yet.
* `tags` - (Optional) A map of tags assigned to snapshots.
* `output_file` - (Optional) The name of output file that saves the filter results.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of snapshot IDs.
* `names` - A list of snapshots names.
* `snapshots` - A list of snapshots. Each element contains the following attributes:
    * `id` - ID of the snapshot.
    * `name` - Name of the snapshot.
    * `description` - Description of the snapshot.
    * `encrypted` - Whether the snapshot is encrypted or not.
    * `progress` - Progress of snapshot creation, presented in percentage.
    * `source_disk_id` - Source disk ID, which is retained after the source disk of the snapshot is deleted.
    * `source_disk_size` - Size of the source disk, measured in GB.
    * `source_disk_type` - Source disk attribute. Value range: `System`,`Data`.
    * `product_code` - Product code on the image market place.
    * `retention_days` - The number of days that an automatic snapshot retains in the console for your instance.
    * `remain_time` - The remaining time of a snapshot creation task, in seconds.
    * `creation_time` - Creation time. Time of creation. It is represented according to ISO8601, and UTC time is used. Format: YYYY-MM-DDThh:mmZ.
    * `status` - The snapshot status. Value range: `progressing`, `accomplished` and `failed`.
    * `usage` - Whether the snapshots are used to create resources or not. Value range: `image`, `disk`, `image_disk` and `none`.
    * `tags` - A map of tags assigned to the snapshot.
