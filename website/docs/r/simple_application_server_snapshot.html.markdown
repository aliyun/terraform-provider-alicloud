---
subcategory: "Simple Application Server"
layout: "alicloud"
page_title: "Alicloud: alicloud_simple_application_server_snapshot"
sidebar_current: "docs-alicloud-resource-simple-application-server-snapshot"
description: |-
  Provides a Alicloud Simple Application Server Snapshot resource.
---

# alicloud_simple_application_server_snapshot

Provides a Simple Application Server Snapshot resource.

For information about Simple Application Server Snapshot and how to use it, see [What is Snapshot](https://www.alibabacloud.com/help/doc-detail/190452.htm).

-> **NOTE:** Available since v1.143.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf_example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_simple_application_server_images" "default" {
  platform = "Linux"
}
data "alicloud_simple_application_server_plans" "default" {
  platform = "Linux"
}

resource "alicloud_simple_application_server_instance" "default" {
  payment_type   = "Subscription"
  plan_id        = data.alicloud_simple_application_server_plans.default.plans.0.id
  instance_name  = var.name
  image_id       = data.alicloud_simple_application_server_images.default.images.0.id
  period         = 1
  data_disk_size = 100
}

data "alicloud_simple_application_server_disks" "default" {
  instance_id = alicloud_simple_application_server_instance.default.id
}

resource "alicloud_simple_application_server_snapshot" "default" {
  disk_id       = data.alicloud_simple_application_server_disks.default.ids.0
  snapshot_name = "${var.name}-${random_integer.default.result}"
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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when create the Snapshot.

## Import

Simple Application Server Snapshot can be imported using the id, e.g.

```shell
$ terraform import alicloud_simple_application_server_snapshot.example <id>
```