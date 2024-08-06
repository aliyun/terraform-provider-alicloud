---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_key_pairs"
sidebar_current: "docs-alicloud-datasource-ecs-key-pairs"
description: |-
  Provides a list of Ecs Key Pairs to the user.
---

# alicloud\_ecs\_key\_pairs

This data source provides the Ecs Key Pairs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.121.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecs_key_pairs" "example" {
  ids        = ["key_pair_name"]
  name_regex = "key_pair_name"
}

output "first_ecs_key_pair_id" {
  value = data.alicloud_ecs_key_pairs.example.pairs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `finger_print` - (Optional, ForceNew) The finger print of the key pair.
* `ids` - (Optional, ForceNew, Computed)  A list of Key Pair IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Key Pair name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The resource group Id.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Key Pair names.
* `pairs` - A list of Ecs Key Pairs. Each element contains the following attributes:
	* `finger_print` - The finger print of the key pair.
	* `id` - The ID of the Key Pair.
	* `key_name` - The Key Pair Name.
	* `resource_group_id` - The Resource Group Id.
	* `tags` - The tags.
		* `tag_key` - The tag key.
		* `tag_value` - The tag value.
  * `instances` - A list of ECS instances that has been bound this key pair.
    * `availability_zone` - The ID of the availability zone where the ECS instance is located.
    * `instance_id` - The ID of the ECS instance.
    * `instance_name` - The name of the ECS instance.
    * `vswitch_id` - The ID of the vSwitch attached to the ECS instance.
    * `public_ip` - The public IP address or EIP of the ECS instance.
    * `private_ip` - The private IP address of the ECS instance.
