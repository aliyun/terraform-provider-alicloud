---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_key_pairs"
sidebar_current: "docs-alicloud-datasource-key-pairs"
description: |-
    Provides a list of available key pairs that can be used by an Alibaba Cloud account.
---

# alicloud\_key\_pairs

-> **DEPRECATED:** This datasource has been renamed to [alicloud_ecs_key_pairs](https://www.terraform.io/docs/providers/alicloud/d/ecs_key_pairs) from version 1.121.0.

This data source provides a list of key pairs in an Alibaba Cloud account according to the specified filters.

## Example Usage

```terraform
# Declare the data source
resource "alicloud_key_pair" "default" {
  key_name = "keyPairDatasource"
}
data "alicloud_key_pairs" "default" {
  name_regex = "${alicloud_key_pair.default.key_name}"
}

```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the resulting key pairs.
* `ids` - (Optional, Available 1.52.1+) A list of key pair IDs.
* `finger_print` - (Optional) A finger print used to retrieve specified key pair.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, Available since v1.57.0) The Id of resource group which the key pair belongs.
* `tags` - (Optional, ForceNew, Available since v1.66.0) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of key pair names.
* `key_pairs` - (Deprecated) A list of key pairs. Each element contains the following attributes:
  * `id` - ID of the key pair.
  * `key_name` - Name of the key pair.
  * `finger_print` - Finger print of the key pair.
  * `key_pair_name` - The name of the key pair.
  * `resource_group_id` - The Id of resource group.
  * `tags` - A map of tags assigned to the key pair.
  * `instances` - A list of ECS instances that has been bound this key pair.
    * `availability_zone` - The ID of the availability zone where the ECS instance is located.
    * `instance_id` - The ID of the ECS instance.
    * `instance_name` - The name of the ECS instance.
    * `vswitch_id` - The ID of the vSwitch attached to the ECS instance.
    * `public_ip` - The public IP address or EIP of the ECS instance.
    * `private_ip` - The private IP address of the ECS instance.
    * `description` - Description of the ECS instance.
    * `image_id` - The image ID.
    * `instance_type` - The instance type.
    * `key_name` - The key name.
    * `region_id` - The region ID.
    * `status` - The status of the ECS instance.
* `pairs` - A list of key pairs. Each element contains the following attributes:
  * `id` - ID of the key pair.
  * `key_pair_name` - Name of the key pair.
  * `key_name` - Name of the key pair.
  * `finger_print` - Finger print of the key pair.
  * `resource_group_id` - The Id of resource group.
  * `tags` - A map of tags assigned to the key pair.
  * `instances` - A list of ECS instances bound to this key pair.
    * `instance_id` - The ID of the ECS instance.
    * `instance_name` - The name of the ECS instance.
    * `description` - Description of the ECS instance.
    * `image_id` - The image ID.
    * `region_id` - The region ID.
    * `availability_zone` - The availability zone ID.
    * `instance_type` - The instance type.
    * `vswitch_id` - The vSwitch ID.
    * `public_ip` - The public IP address.
    * `private_ip` - The private IP address.
    * `key_name` - The key name.
    * `status` - The status of the ECS instance.
