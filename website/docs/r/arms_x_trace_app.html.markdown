---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_x_trace_app"
description: |-
  Provides a Alicloud ARMS X Trace App resource.
---

# alicloud_arms_x_trace_app

Provides a ARMS X Trace App resource. The Application of Trace.

For information about ARMS X Trace App and how to use it, see [What is X Trace App](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.215.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_resource_manager_resource_group" "defaultV8g7dc" {
  display_name        = "testg1"
  resource_group_name = var.name

  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}

resource "alicloud_resource_manager_resource_group" "defaultkkepEi" {
  display_name        = "testg2"
  resource_group_name = var.name

  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}


resource "alicloud_arms_x_trace_app" "default" {
  x_trace_app_name = var.name

  resource_group_id = alicloud_resource_manager_resource_group.defaultV8g7dc.id
}
```

## Argument Reference

The following arguments are supported:
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `tags` - (Optional, Map) The tag of the resource.
* `x_trace_app_name` - (Required, ForceNew) The name of the resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the X Trace App.
* `delete` - (Defaults to 5 mins) Used when delete the X Trace App.
* `update` - (Defaults to 5 mins) Used when update the X Trace App.

## Import

ARMS X Trace App can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_x_trace_app.example <id>
```