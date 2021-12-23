---
subcategory: "Simple Application Server"
layout: "alicloud"
page_title: "Alicloud: alicloud_simple_application_server_snapshot"
sidebar_current: "docs-alicloud-resource-simple-application-server-snapshot"
description: |-
  Provides a Alicloud Simple Application Server Snapshot resource.
---

# alicloud\_simple\_application\_server\_snapshot

Provides a Simple Application Server Snapshot resource.

For information about Simple Application Server Snapshot and how to use it, see [What is Snapshot](https://www.alibabacloud.com/help/doc-detail/190452.htm).

-> **NOTE:** Available in v1.143.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_simple_application_server_instances" "example" {}

data "alicloud_simple_application_server_images" "example" {}

data "alicloud_simple_application_server_plans" "example" {}

resource "alicloud_simple_application_server_instance" "example" {
  count         = length(data.alicloud_simple_application_server_instances.example.ids) > 0 ? 0 : 1
  payment_type  = "Subscription"
  plan_id       = data.alicloud_simple_application_server_plans.example.plans.0.id
  instance_name = "example_value"
  image_id      = data.alicloud_simple_application_server_images.example.images.0.id
  period        = 1
}

data "alicloud_simple_application_server_disks" "example" {
  instance_id = length(data.alicloud_simple_application_server_instances.example.ids) > 0 ? data.alicloud_simple_application_server_instances.example.ids.0 : alicloud_simple_application_server_instance.example.0.id
}

resource "alicloud_simple_application_server_snapshot" "example" {
  disk_id       = data.alicloud_simple_application_server_disks.example.ids.0
  snapshot_name = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `disk_id` - (Required, ForceNew) The ID of the disk.
* `snapshot_name` - (Required, ForceNew) The name of the snapshot. The name must be `2` to `50` characters in length. It must start with a letter and cannot start with `http://` or `https://`. It can contain letters, digits, colons (:), underscores (_), periods (.),and hyphens (-).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Snapshot.
* `status` - The status of the snapshot. Valid values: `Progressing`, `Accomplished` and `Failed`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when create the Snapshot.

## Import

Simple Application Server Snapshot can be imported using the id, e.g.

```
$ terraform import alicloud_simple_application_server_snapshot.example <id>
```