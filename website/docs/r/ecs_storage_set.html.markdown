---
subcategory: "unreleased"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_storage_set"
description: |-
  Provides a Alicloud ECS Storage Set resource.
---

# alicloud_ecs_storage_set

Provides a ECS Storage Set resource.



For information about ECS Storage Set and how to use it, see [What is Storage Set](https://next.api.alibabacloud.com/document/Ecs/2014-05-26/CreateStorageSet).

-> **NOTE:** Available since v1.295.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_ecs_storage_set" "default" {
  description          = "wTest"
  zone_id              = "cn-hangzhou-i"
  max_partition_number = 2
  storage_set_name     = "wexample"
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional, Computed) description
* `max_partition_number` - (Optional, Computed, ForceNew, Int) max partition number
* `resource_group_id` - (Optional, Computed, ForceNew) resource group id. Storage Set is not a supported resource type of JoinResourceGroup, so the resource group can only be set when the Storage Set is created; changing it recreates the resource.
* `storage_set_name` - (Optional, Computed) storage set name
* `tags` - (Optional, Computed, ForceNew, Map) tags. Storage Set is not a supported resource type of TagResources/UntagResources, so tags can only be set when the Storage Set is created; changing them recreates the resource. To clear all tags use an explicit `tags = {}` — omitting the `tags` block preserves the existing tags.
* `zone_id` - (Required, ForceNew) zone id

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource, same as `storage_set_id`.
* `storage_set_id` - The ID of the Storage Set.
* `region_id` - The region ID of the Storage Set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Storage Set.
* `delete` - (Defaults to 5 mins) Used when delete the Storage Set.
* `update` - (Defaults to 5 mins) Used when update the Storage Set.

## Import

ECS Storage Set can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_storage_set.example <storage_set_id>
```
