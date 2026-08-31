---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_snapshot_groups"
sidebar_current: "docs-alicloud-datasource-ecs-snapshot-groups"
description: |-
  Provides a list of Ecs Snapshot Groups to the user.
---

# alicloud\_ecs\_snapshot\_groups

This data source provides the Ecs Snapshot Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.160.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecs_snapshot_groups" "ids" {
  ids = ["example-id"]
}
output "ecs_snapshot_group_id_1" {
  value = data.alicloud_ecs_snapshot_groups.ids.groups.0.id
}

data "alicloud_ecs_snapshot_groups" "nameRegex" {
  name_regex = "^my-SnapshotGroup"
}
output "ecs_snapshot_group_id_2" {
  value = data.alicloud_ecs_snapshot_groups.nameRegex.groups.0.id
}

data "alicloud_ecs_snapshot_groups" "status" {
  status = "accomplished"
}
output "ecs_snapshot_group_id_3" {
  value = data.alicloud_ecs_snapshot_groups.status.groups.0.id
}

data "alicloud_ecs_snapshot_groups" "instanceId" {
  instance_id = "example-instance_id"
}
output "ecs_snapshot_group_id_4" {
  value = data.alicloud_ecs_snapshot_groups.instanceId.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Snapshot Group IDs.
* `instance_id` - (Optional, ForceNew) The ID of the instance.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Snapshot Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `snapshot_group_name` - (Optional, ForceNew) The name of the snapshot-consistent group.
* `status` - (Optional, ForceNew) The state of snapshot-consistent group. Valid Values: `accomplished`, `failed` and `progressing`.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the snapshot group.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Snapshot Group names.
* `groups` - A list of Ecs Snapshot Groups. Each element contains the following attributes:
  * `description` - The description of the snapshot-consistent group.
  * `id` - The ID of the Snapshot Group.
  * `instance_id` - The ID of the instance.
  * `resource_group_id` - The ID of the resource group to which the snapshot consistency group belongs.
  * `snapshot_group_id` - The first ID of the resource.
  * `snapshot_group_name` - The name of the snapshot-consistent group.
  * `status` - The status of the resource.
  * `tags` - List of label key-value pairs.
    * `tag_key` - The key of the tag.
    * `tag_value` - The value of the tag.