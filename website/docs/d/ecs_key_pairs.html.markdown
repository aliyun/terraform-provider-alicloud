---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_key_pairs"
sidebar_current: "docs-alicloud-datasource-ecs-key-pairs"
description: |-
  Provides a list of Ecs Key Pairs to the user.
---

# alicloud_ecs_key_pairs

This data source provides the Ecs Key Pairs of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.121.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_ecs_key_pair" "default" {
  key_pair_name     = var.name
  public_key        = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.1
  tags = {
    Created = "TF"
    For     = "KeyPair",
  }
}

data "alicloud_ecs_key_pairs" "ids" {
  ids = [alicloud_ecs_key_pair.default.id]
}

output "ecs_key_pair_id_0" {
  value = data.alicloud_ecs_key_pairs.ids.pairs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List)  A list of Key Pair IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Key Pair name.
* `finger_print` - (Optional, ForceNew) The fingerprint of the key pair.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Key Pair names.
* `pairs` - A list of Ecs Key Pairs. Each element contains the following attributes:
  * `id` - The ID of the Key Pair.
  * `key_pair_name` - The name of the Key Pair.
  * `key_name` - The name of the Key Pair.
  * `finger_print` - The fingerprint of the Key Pair.
  * `resource_group_id` - The ID of the resource group.
  * `tags` - The tags of the Key Pair.
  * `instances` - A list of ECS instances that has been bound this Key Pair.
    * `instance_id` - The ID of the ECS instance.
    * `instance_name` - The name of the ECS instance.
    * `description` - The description of the ECS instance.
    * `image_id` - The image ID of the instance.
    * `region_id` - The region ID of the instance.
    * `availability_zone` - The zone ID of the instance.
    * `instance_type` - The instance type of the instance.
    * `vswitch_id` - The ID of the vSwitch.
    * `public_ip` - The public IP address or EIP of the ECS instance.
    * `private_ip` - The private IP address of the ECS instance.
    * `key_name` - The name of the key pair.
    * `status` - The status of the instance.
* `key_pairs` - (Deprecated since v1.121.0) A list of Ecs Key Pairs. Each element contains the following attributes:
  * `id` - The ID of the Key Pair.
  * `key_pair_name` - The name of the Key Pair.
  * `key_name` - The name of the Key Pair.
  * `finger_print` - The fingerprint of the Key Pair.
  * `resource_group_id` - The ID of the resource group.
  * `tags` - The tags of the Key Pair.
  * `instances` - A list of ECS instances that has been bound this Key Pair.
    * `instance_id` - The ID of the ECS instance.
    * `instance_name` - The name of the ECS instance.
    * `description` - The description of the ECS instance.
    * `image_id` - The image ID of the instance.
    * `region_id` - The region ID of the instance.
    * `availability_zone` - The zone ID of the instance.
    * `instance_type` - The instance type of the instance.
    * `vswitch_id` - The ID of the vSwitch.
    * `public_ip` - The public IP address or EIP of the ECS instance.
    * `private_ip` - The private IP address of the ECS instance.
    * `key_name` - The name of the key pair.
    * `status` - The status of the instance.