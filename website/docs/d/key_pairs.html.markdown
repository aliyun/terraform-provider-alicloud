---
layout: "alicloud"
page_title: "Alicloud: alicloud_key_pairs"
sidebar_current: "docs-alicloud-datasource-key-pairs"
description: |-
    Provides a list of Availability Key Pairs which can be used by an Alicloud account.
---

# alicloud\_key\_pairs

The Key Pairs data source provides a list of Alicloud Key Pairs in an Alicloud account according to the specified filters.

## Example Usage

```
# Declare the data source
data "alicloud_key_pairs" "name_regex" {
	name_regex = "test"
	output_file = "my_key_pairs.json"
}

# Bind a key pair for several ecs instances using the first matched key pair

resource "alicloud_key_pair_attachment" "attachment" {
  key_name = "${data.alicloud_key_pairs.default.key_pairs.0.id}"
  instance_ids = [...]
}

```

## Argument Reference

The following arguments are supported:

* `name_regex` - A regex string to apply to the key pair list returned by Alicloud.
* `finger_print` - A finger print used to retrieve specified key pair.
* `output_file` - (Optional) The name of file that can save key pairs data source after running `terraform plan`.

## Attributes Reference

A list of key pairs will be exported and its every element contains the following attributes:

* `id` - ID of the key pair.
* `key_name` - Name of the key pair.
* `finger_print` - Finger print of the key pair.
* `instances` - A List of ECS instances that has been bound a specified key pair.
    * `availability_zone` - The ID of availability zone that ECS instance launched.
    * `instance_id` - The ID of ECS instance.
    * `instance_name` - The name of ECS instance.
    * `vswitch_id` - The ID of VSwitch that ECS instance launched.
    * `public_ip` - The public IP address or EIP of the ECS instance.
    * `private_ip` - The private IP address of the ECS instance.