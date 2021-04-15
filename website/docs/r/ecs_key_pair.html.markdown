---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_key_pair"
sidebar_current: "docs-alicloud-resource-ecs-key-pair"
description: |-
  Provides a Alicloud ECS Key Pair resource.
---

# alicloud\_ecs\_key\_pair

Provides a ECS Key Pair resource.

For information about ECS Key Pair and how to use it, see [What is Key Pair](https://www.alibabacloud.com/help/en/doc-detail/51771.htm).

-> **NOTE:** Available in v1.121.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecs_key_pair" "example" {
  key_pair_name = "key_pair_name"
}

// Using name prefix to build key pair
resource "alicloud_ecs_key_pair" "prefix" {
  key_name_prefix = "terraform-test-key-pair-prefix"
}

// Import an existing public key to build a alicloud key pair
resource "alicloud_ecs_key_pair" "publickey" {
  key_pair_name = "my_public_key"
  public_key    = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
}

```

## Argument Reference

The following arguments are supported:

* `key_file` - (Optional, ForceNew) The key file.
* `key_name` - (Optional, ForceNew) The key pair's name. It is the only in one Alicloud account.
* `key_name_prefix` - (Optional, ForceNew) The key pair name's prefix. It is conflict with `key_pair_name`. If it is specified, terraform will using it to build the only key name.
* `public_key` - (Optional, ForceNew) You can import an existing public key and using Alicloud key pair to manage it. If this parameter is specified, `resource_group_id` is the key pair belongs.
* `resource_group_id` - (Optional) The Id of resource group which the key pair belongs.
* `tags` - (Optional) A mapping of tags to assign to the resource.

-> **NOTE:** If `key_pair_name` and `key_name_prefix` are not set, terraform will produce a specified ID to replace.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Key Pair. Value as `key_pair_name`.
* `finger_print` The finger print of the key pair.

## Import

ECS Key Pair can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_key_pair.example <key_name>
```