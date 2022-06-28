---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_snapshots"
sidebar_current: "docs-alicloud-datasource-ecs-snapshots"
description: |-
  Provides a list of Ecs Snapshots to the user.
---

# alicloud\_ecs\_snapshots

This data source provides the Ecs Snapshots of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.120.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecs_snapshots" "example" {
  ids        = ["s-bp1fvuxxxxxxxx"]
  name_regex = "tf-test"
}

output "first_ecs_snapshot_id" {
  value = data.alicloud_ecs_snapshots.example.snapshots.0.id
}
```

## Argument Reference

The following arguments are supported:

* `category` - (Optional, ForceNew) The category of the snapshot. Valid Values: `flash` and `standard`.
* `dry_run` - (Optional, ForceNew) Specifies whether to check the validity of the request without actually making the request.
* `encrypted` - (Optional, ForceNew) Specifies whether the snapshot is encrypted.
* `ids` - (Optional, ForceNew, Computed)  A list of Snapshot IDs.
* `kms_key_id` - (Optional, ForceNew) The kms key id.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Snapshot name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The resource group id.
* `snapshot_link_id` - (Optional, ForceNew) The snapshot link id.
* `snapshot_name` - (Optional, ForceNew) The name of the snapshot.
* `snapshot_type` - (Optional, ForceNew) The type of the snapshot. Valid Values: `auto`, `user` and `all`. Default to: `all`.
* `source_disk_type` - (Optional, ForceNew) The type of the disk for which the snapshot was created. Valid Values: `System`, `Data`.
* `status` - (Optional, ForceNew) The status of the snapshot. Valid Values: `accomplished`, `failed`, `progressing` and `all`.
* `usage` - (Optional, ForceNew) A resource type that has a reference relationship. Valid Values: `image`, `disk`, `image_disk` and `none`.
* `tags` - (Optional) A mapping of tags to assign to the snapshot.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Snapshot names.
* `snapshots` - A list of Ecs Snapshots. Each element contains the following attributes:
	* `category` - The category of the snapshot.
	* `description` - The description of the snapshot.
	* `disk_id` - The source disk id.
	* `encrypted` - Whether the snapshot is encrypted.
	* `id` - The ID of the Snapshot.
	* `instant_access` - Whether snapshot speed availability is enabled.
	* `instant_access_retention_days` - Specifies the retention period of the instant access feature. After the retention period ends, the snapshot is automatically released.
	* `product_code` - The product number inherited from the mirror market.
	* `progress` - Snapshot creation progress, in percentage.
	* `remain_time` - Remaining completion time for the snapshot being created.
	* `resource_group_id` - The resource group id.
	* `retention_days` - Automatic snapshot retention days.
	* `snapshot_id` - The snapshot id.
	* `snapshot_name` - Snapshot Display Name.
	* `snapshot_type` - Snapshot creation type.
	* `snapshot_sn` - The serial number of the snapshot.
	* `source_disk_size` - Source disk capacity.
	* `source_disk_type` - Source disk attributes.
	* `source_storage_type` - Original disk type.
	* `status` - The status of the snapshot.
	* `tags` - The tags.
		* `tag_key` - The tag key.
		* `tag_value` - The tag value.
	* `usage` - A resource type that has a reference relationship.
