---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_key_pair"
sidebar_current: "docs-alicloud-resource-key-pair"
description: |-
  Provides a Alicloud key pair resource.
---

# alicloud\_key\_pair

-> **DEPRECATED:** This resource has been renamed to [alicloud_ecs_key_pair](https://www.terraform.io/docs/providers/alicloud/r/ecs_key_pair) from version 1.121.0.

Provides a key pair resource.

## Example Usage

Basic Usage

```
resource "alicloud_key_pair" "basic" {
  key_name = "terraform-test-key-pair"
}

// Using name prefix to build key pair
resource "alicloud_key_pair" "prefix" {
  key_name_prefix = "terraform-test-key-pair-prefix"
}

// Import an existing public key to build a alicloud key pair
resource "alicloud_key_pair" "publickey" {
  key_name   = "my_public_key"
  public_key = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
}
```
## Argument Reference

The following arguments are supported:

* `key_name` - (ForceNew) The key pair's name. It is the only in one Alicloud account.
* `key_name_prefix` - (ForceNew) The key pair name's prefix. It is conflict with `key_name`. If it is specified, terraform will using it to build the only key name.
* `public_key` - (ForceNew) You can import an existing public key and using Alicloud key pair to manage it. If this parameter is specified, `resource_group_id` is the key pair belongs.
* `key_file` - (ForceNew) The name of file to save your new key pair's private key. Strongly suggest you to specified it when you creating key pair, otherwise, you wouldn't get its private key ever.
* `resource_group_id` - (Optional, Available in 1.57.0+, Modifiable in 1.115.0+) The Id of resource group which the key pair belongs.
* `tags` - (Optional, Available in v1.66.0+) A mapping of tags to assign to the resource.
-> **NOTE:** If `key_name` and `key_name_prefix` are not set, terraform will produce a specified ID to replace.

## Attributes Reference

* `key_name` - The name of the key pair.
* `fingerprint` The finger print of the key pair.

## Import

Key pair can be imported using the name, e.g.

```
$ terraform import alicloud_key_pair.example my_public_key
```
