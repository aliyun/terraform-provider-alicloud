---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_snapshot"
sidebar_current: "docs-alicloud-resource-ecs-snapshot"
description: |-
  Provides a Alicloud ECS Snapshot resource.
---

# alicloud\_ecs\_snapshot

Provides a ECS Snapshot resource.

For information about ECS Snapshot and how to use it, see [What is Snapshot](https://www.alibabacloud.com/help/en/doc-detail/25524.htm).

-> **NOTE:** Available in v1.120.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecs_snapshot" "default" {
  category       = "standard"
  description    = "Test For Terraform"
  disk_id        = "d-gw8csgxxxxxxxxx"
  retention_days = "20"
  snapshot_name  = "tf-test"
  tags = {
    Created = "TF"
    For     = "Acceptance-test"
  }
}

```

## Argument Reference

The following arguments are supported:

* `category` - (Optional, ForceNew) The category of the snapshot. Valid Values: `standard` and `flash`.
* `description` - (Optional) The description of the snapshot.
* `disk_id` - (Required, ForceNew) The ID of the disk.
* `force` - (Optional) Specifies whether to forcibly delete the snapshot that has been used to create disks.
* `instant_access` - (Optional) Specifies whether to enable the instant access feature.
* `instant_access_retention_days` - (Optional, ForceNew) Specifies the retention period of the instant access feature. After the retention period ends, the snapshot is automatically released.
* `resource_group_id` - (Optional, ForceNew) The resource group id.
* `retention_days` - (Optional, ForceNew) The retention period of the snapshot.
* `snapshot_name` - (Optional) The name of the snapshot.
* `name` - (Optional, Deprecated in v1.120.0+) Field `name` has been deprecated from provider version 1.120.0. New field `snapshot_name` instead. 
* `tags` - (Optional) A mapping of tags to assign to the snapshot.

-> **NOTE:** If `force` is true, After an snapshot is deleted, the disks created from this snapshot cannot be re-initialized.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Snapshot.
* `status` - The status of snapshot.

## Import

ECS Snapshot can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_snapshot.example <id>
```
