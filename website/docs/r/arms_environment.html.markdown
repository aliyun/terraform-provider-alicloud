---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_environment"
description: |-
  Provides a Alicloud ARMS Environment resource.
---

# alicloud_arms_environment

Provides a ARMS Environment resource. The arms environment.

For information about ARMS Environment and how to use it, see [What is Environment](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_resource_manager_resource_group" "rg" {
  display_name        = "resource-group-for-create"
  resource_group_name = var.name

}

resource "alicloud_resource_manager_resource_group" "rg2" {
  display_name        = "resource-group-for-move"
  resource_group_name = var.name

}


resource "alicloud_arms_environment" "default" {
  environment_type = "Cloud"
  environment_name = var.name

  environment_sub_type = "Cloud"
  bind_resource_id     = "cn-hangzhou"
  resource_group_id    = alicloud_resource_manager_resource_group.rg.id
  aliyun_lang          = "zh"
}
```

## Argument Reference

The following arguments are supported:
* `aliyun_lang` - (Optional) The locale. The default is Chinese zh | en.
* `bind_resource_id` - (Optional, ForceNew) The id or vpcId of the bound container instance.
* `environment_id` - (Optional, ForceNew, Computed) The first ID of the resource.
* `environment_name` - (Optional) The name of the resource.
* `environment_sub_type` - (Required, ForceNew) Subtype of environment:
  - Type of CS: ACK is currently supported.
  - Type of ECS: currently supports ECS.
  - Type of Cloud: currently supports Cloud.
* `environment_type` - (Required, ForceNew) Type of environment.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `tags` - (Optional, Map) The tag of the resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Environment.
* `delete` - (Defaults to 5 mins) Used when delete the Environment.
* `update` - (Defaults to 5 mins) Used when update the Environment.

## Import

ARMS Environment can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_environment.example <id>
```