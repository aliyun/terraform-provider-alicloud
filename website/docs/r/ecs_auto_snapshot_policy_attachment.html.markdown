---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_auto_snapshot_policy_attachment"
sidebar_current: "docs-alicloud-resource-ecs-auto-snapshot-policy-attachment"
description: |-
  Provides a Alicloud ECS Auto Snapshot Policy Attachment resource.
---

# alicloud\_ecs\_auto\_snapshot\_policy\_attachment

Provides a ECS Auto Snapshot Policy Attachment resource.

For information about ECS Auto Snapshot Policy Attachment and how to use it, see [What is Auto Snapshot Policy Attachment](https://www.alibabacloud.com/help/en/doc-detail/25531.htm).

-> **NOTE:** Available in v1.122.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecs_auto_snapshot_policy_attachment" "example" {
  auto_snapshot_policy_id = "s-ge465xxxx"
  disk_id                 = "d-gw835xxxx"
}

```

## Argument Reference

The following arguments are supported:

* `auto_snapshot_policy_id` - (Required, ForceNew) The auto snapshot policy id.
* `disk_id` - (Required, ForceNew) The disk id.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Auto Snapshot Policy Attachment. The value is formatted `<auto_snapshot_policy_id>:<disk_id>`.

## Import

ECS Auto Snapshot Policy Attachment can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_auto_snapshot_policy_attachment.example s-abcd12345:d-abcd12345
```
