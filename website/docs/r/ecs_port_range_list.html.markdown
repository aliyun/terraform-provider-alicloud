---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_port_range_list"
sidebar_current: "docs-alicloud-resource-ecs-port-range-list"
description: |-
  Provides a ECS Port Range List resource.
---

# alicloud_ecs_port_range_list

Provides a ECS Port Range List resource.

For information about ECS Port Range List and how to use it, see [What is Port Range List](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/createportrangelist).

-> **NOTE:** Available since v1.295.0.

## Example Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_ecs_port_range_list" "example" {
  port_range_list_name = "tf-example"
  description          = "terraform-example"
  max_entries          = 10
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids[0]

  entry {
    port_range  = "80/80"
    description = "http"
  }

  entry {
    port_range  = "443/443"
    description = "https"
  }

  tags = {
    Created = "TF"
    For     = "Example"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the port range list. The description must be `2` to `256` characters in length and cannot start with `http://` or `https://`.
* `entry` - (Optional) The entries of the port range list. See [`entry`](#entry) below.
* `max_entries` - (Required) The maximum number of entries that the port range list can contain. Valid values: `1` to `2000`. This value can be updated.
* `port_range_list_name` - (Required) The name of the port range list. The name must be `2` to `128` characters in length, must start with an English letter or a Chinese character, and can contain Chinese characters, English letters, digits, colons (:), underscores (_), periods (.), and hyphens (-). It cannot start with `http://`, `https://`, `com.aliyun`, or `com.alibabacloud`.
* `resource_group_id` - (Optional, ForceNew, Computed) The ID of the resource group to which the port range list belongs. The resource group can only be specified at creation: PortRangeList does not support JoinResourceGroup and ModifyPortRangeList does not accept ResourceGroupId.
* `tags` - (Optional) A mapping of tags to assign to the resource.

### `entry`

The `entry` block supports:

* `port_range` - (Optional) The port range. Format: `<start>/<end>`, for example `80/80`. Port ranges in a list cannot be duplicated.
* `description` - (Optional) The description of the port range. The description must be `2` to `32` characters in length and cannot start with `http://` or `https://`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Port Range List. It is the same as `port_range_list_id`.
* `association_count` - The number of resources associated with the port range list.
* `create_time` - The time when the port range list was created.

## Import

ECS Port Range List can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_port_range_list.example <id>
```
