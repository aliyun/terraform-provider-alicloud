---
layout: "alicloud"
page_title: "Alicloud: alicloud_key_pairs"
sidebar_current: "docs-alicloud-datasource-key-pairs"
description: |-
    Provides a list of available key pairs that can be used by an Alibaba Cloud account.
---

# alicloud\_key\_pairs

This data source provides a list of key pairs in an Alibaba Cloud account according to the specified filters.

## Example Usage

```
# Declare the data source
data "alicloud_key_pairs" "key_pairs_ds" {
	name_regex = "test"
	output_file = "my_key_pairs.json"
}

# Bind a key pair for several ECS instances by using the first matched key pair

resource "alicloud_key_pair_attachment" "attachment" {
  key_name = "${data.alicloud_key_pairs.key_pairs_ds.key_pairs.0.id}"
  instance_ids = [...]
}

```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the resulting key pairs.
* `finger_print` - (Optional) A finger print used to retrieve specified key pair.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `key_pairs` - A list of key pairs. Each element contains the following attributes:
  * `id` - ID of the key pair.
  * `key_name` - Name of the key pair.
  * `finger_print` - Finger print of the key pair.
  * `instances` - A list of ECS instances that has been bound this key pair.
    * `availability_zone` - The ID of the availability zone where the ECS instance is located.
    * `instance_id` - The ID of the ECS instance.
    * `instance_name` - The name of the ECS instance.
    * `vswitch_id` - The ID of the VSwitch attached to the ECS instance.
    * `public_ip` - The public IP address or EIP of the ECS instance.
    * `private_ip` - The private IP address of the ECS instance.