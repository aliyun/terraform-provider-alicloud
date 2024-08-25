---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_snapshot"
sidebar_current: "docs-alicloud-resource-snapshot"
description: |-
  Provides an ECS snapshot resource.
---

# alicloud\_snapshot

-> **DEPRECATED:** This resource has been renamed to [alicloud_ecs_snapshot](https://www.terraform.io/docs/providers/alicloud/r/ecs_snapshot) from version 1.120.0.

Provides an ECS snapshot resource.

For information about snapshot and how to use it, see [Snapshot](https://www.alibabacloud.com/help/doc-detail/25460.html).

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_snapshot&exampleId=968b3dc2-93e8-50df-633c-42947cb57c7ad4f11f3d&activeTab=example&spm=docs.r.snapshot.0.968b3dc293&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_snapshot" "snapshot" {
  disk_id     = alicloud_disk_attachment.instance-attachment.disk_id
  name        = "test-snapshot"
  description = "this snapshot is created for testing"
  tags = {
    version = "1.2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `disk_id` - (Required, ForceNew) The source disk ID.
* `name` - (Optional, ForceNew) The name of the snapshot to be created. The name must be 2 to 128 characters in length. It must start with a letter and cannot start with http:// or https://. It can contain letters, digits, colons (:), underscores (_), and hyphens (-).
It cannot start with auto, because snapshot names starting with auto are recognized as automatic snapshots.
* `resource_group_id` - (Optional, ForceNew, Available in 1.94.0+) The ID of the resource group.
* `description` - (Optional, ForceNew) Description of the snapshot. This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `tags` - (Optional) A mapping of tags to assign to the resource.

### Timeouts

-> **NOTE:** Available in 1.51.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when creating the snapshot (until it reaches the initial `SnapshotCreatingAccomplished` status). 
* `delete` - (Defaults to 2 mins) Used when terminating the snapshot. 

## Attributes Reference

The following attributes are exported:

* `id` - The snapshot ID.

## Import

Snapshot can be imported using the id, e.g.

```shell
$ terraform import alicloud_snapshot.snapshot s-abc1234567890000
```
