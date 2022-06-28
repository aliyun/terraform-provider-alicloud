---
subcategory: "Simple Application Server"
layout: "alicloud"
page_title: "Alicloud: alicloud_simple_application_server_custom_image"
sidebar_current: "docs-alicloud-resource-simple-application-server-custom-image"
description: |-
  Provides a Alicloud Simple Application Server Custom Image resource.
---

# alicloud\_simple\_application\_server\_custom\_image

Provides a Simple Application Server Custom Image resource.

For information about Simple Application Server Custom Image and how to use it, see [What is Custom Image](https://www.alibabacloud.com/help/zh/doc-detail/333535.htm).

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
  disk_type   = "System"
  instance_id = length(data.alicloud_simple_application_server_instances.example.ids) > 0 ? data.alicloud_simple_application_server_instances.example.ids.0 : alicloud_simple_application_server_instance.example.0.id
}

resource "alicloud_simple_application_server_snapshot" "example" {
  disk_id       = data.alicloud_simple_application_server_disks.example.ids.0
  snapshot_name = "example_value"
}

resource "alicloud_simple_application_server_custom_image" "example" {
  custom_image_name  = "example_value"
  instance_id        = data.alicloud_simple_application_server_disks.example.disks.0.instance_id
  system_snapshot_id = alicloud_simple_application_server_snapshot.example.id
  status             = "Share"
  description        = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `custom_image_name` - (Required, ForceNew) The name of the resource. The name must be `2` to `128` characters in length. It must start with a letter or a number. It can contain letters, digits, colons (:), underscores (_) and hyphens (-).
* `description` - (Optional, ForceNew) The description of the Custom Image.
* `instance_id` - (Required) The ID of the instance.
* `system_snapshot_id` - (Required) The ID of the system snapshot.
* `status` - (Optional) The Shared status of the Custom Image. Valid values: `Share`, `UnShare`.

 **NOTE:** The `status` will be automatically change to `UnShare` when the resource is deleted, please operate with caution.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Custom Image.

## Import

Simple Application Server Custom Image can be imported using the id, e.g.

```
$ terraform import alicloud_simple_application_server_custom_image.example <id>
```