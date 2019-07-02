---
layout: "alicloud"
page_title: "Alicloud: alicloud_snapshot"
sidebar_current: "docs-alicloud-resource-snapshot"
description: |-
  Provides an ECS snapshot resource.
---

# alicloud\_snapshot

Provides an ECS snapshot resource.

For information about snapshot and how to use it, see [Snapshot](https://www.alibabacloud.com/help/doc-detail/25460.html).

## Example Usage

```
resource "alicloud_snapshot" "snapshot" {
  disk_id = "${alicloud_disk_attachment.instance-attachment.disk_id}"
  name = "test-snapshot"
  description = "this snapshot is created for testing"
  tags = {
    version = "1.2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `disk_id` - (Required, ForceNew) The source disk ID.
* `name` - (Optional, ForceNew) Name of the snapshot. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-", ".", "_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Default value is null.
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

```
$ terraform import alicloud_snapshot.snapshot s-abc1234567890000
```
