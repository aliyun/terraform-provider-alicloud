---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_key_pair"
description: |-
  Provides a Alicloud ECS Key Pair resource.
---

# alicloud_ecs_key_pair

Provides a ECS Key Pair resource.

For information about ECS Key Pair and how to use it, see [What is Key Pair](https://www.alibabacloud.com/help/en/doc-detail/51771.htm).

-> **NOTE:** Available since v1.121.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_key_pair&exampleId=3dd74fa1-1377-6034-b4c4-c41fd49d223c74831f0e&activeTab=example&spm=docs.r.ecs_key_pair.0.3dd74fa113&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `key_pair_name` - (Optional, ForceNew) The name of the key pair. The name must be 2 to 128 characters in length. The name must start with a letter and cannot start with http:// or https://. The name can contain letters, digits, colons (:), underscores (_), and hyphens (-).
* `key_name_prefix` - (Optional, ForceNew) The key pair name's prefix. It is conflict with `key_pair_name`. If it is specified, terraform will using it to build the only key name.
* `public_key` - (Optional) The public key of the key pair.
* `resource_group_id` - (Optional) The ID of the resource group to which to add the key pair.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `key_file` - (Optional, ForceNew) The key file.
* `key_name` - (Optional, ForceNew, Deprecated since v1.121.0) Field `key_name` has been deprecated from provider version 1.121.0. New field `key_pair_name` instead.

-> **NOTE:** If `key_pair_name` and `key_name_prefix` are not set, terraform will produce a specified ID to replace.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Key Pair.
* `finger_print` The fingerprint of the key pair.
* `create_time` - (Available since v1.237.0) The time when the key pair was created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Key Pair.
* `delete` - (Defaults to 5 mins) Used when delete the Key Pair.
* `update` - (Defaults to 5 mins) Used when update the Key Pair.

## Import

ECS Key Pair can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_key_pair.example <id>
```
