---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_key_pair"
sidebar_current: "docs-alicloud-resource-ecs-key-pair"
description: |-
  Provides a Alicloud ECS Key Pair resource.
---

# alicloud_ecs_key_pair

Provides a ECS Key Pair resource. Secret key pair.

For information about ECS Key Pair and how to use it, see [What is Key Pair](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available in v1.207.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_ecs_key_pair" "default" {
}
```


## Argument Reference

The following arguments are supported:
* `resource_group_id` - (Optional, Computed, Available in v1.121.0+) ResourceGroupId.
* `resource_type` - (Optional) Resource type.
* `tags` - (Optional, ForceNew, Map, Available in v1.121.0+) Tags.



## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Key Pair.
* `delete` - (Defaults to 5 mins) Used when delete the Key Pair.
* `update` - (Defaults to 5 mins) Used when update the Key Pair.

## Import

ECS Key Pair can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_key_pair.example <id>
```