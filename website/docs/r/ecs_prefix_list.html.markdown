---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_prefix_list"
sidebar_current: "docs-alicloud-resource-ecs-prefix-list"
description: |-
  Provides a Alicloud ECS Prefix List resource.
---

# alicloud\_ecs\_prefix\_list\

Provides a ECS Prefix List resource.

For information about ECS Prefix List and how to use it, see [What is Prefix List.](https://www.alibabacloud.com/help/en/doc-detail/207969.html).

-> **NOTE:** Available in v1.152.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecs_prefix_list" "default" {
  address_family   = "IPv4"
  max_entries      = 2
  prefix_list_name = "tftest"
  description      = "description"
  entry {
    cidr        = "192.168.0.0/24"
    description = "description"
  }
}
```

## Argument Reference

The following arguments are supported:

* `address_family` - (Required, ForceNew) The IP address family. Valid values: `IPv4`,`IPv6`.
* `max_entries` - (Required, ForceNew) The maximum number of entries that the prefix list can contain.  Valid values: 1 to 200.
* `prefix_list_name` - (Required) The name of the prefix. The name must be 2 to 128 characters in length. It must start with a letter and cannot start with `http://`, `https://`, `com.aliyun`, or `com.alibabacloud`. It can contain letters, digits, colons (:), underscores (_), periods (.), and hyphens (-).
* `description` - (Optional) The description of the prefix list. The description must be 2 to 256 characters in length and cannot start with `http://` or `https://`.
* `entry` - (Required) The Entry. The details see Block `entry`. 



#### entry
The entry supports the following:

* `cidr` - (Optional) The CIDR block in entry. This parameter is empty by default.  Take note of the following items:
  * The total number of entries must not exceed the `max_entries` value.
  * CIDR block types are determined by the IP address family. You cannot combine `IPv4` and `IPv6` CIDR blocks in a single entry.
  * CIDR blocks must be unique across all entries in a prefix list. For example, you cannot specify 192.168.1.0/24 twice in the entries of the prefix list.
  * IP addresses are supported. The system converts IP addresses into CIDR blocks. For example, if you specify 192.168.1.100, the system converts it into the 192.168.1.100/32 CIDR block.
  * If an IPv6 CIDR block is used, the system converts it to the zero compression format and changes uppercase letters into lowercase ones. For example, if you specify 2001:0DB8:0000:0000:0000:0000:0000:0000/32, the system converts it into 2001:db8::/32.
  * For more information about CIDR blocks, see the "What is CIDR block?" section of the [Network FAQ](https://www.alibabacloud.com/help/doc-detail/40637.htm) topic.  * The total number of entries must not exceed the `max_entries` value.
* `description` - (Optional) The description in entry. The description must be 2 to 32 characters in length and cannot start with `http://` or `https://`.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Prefix List.


## Import

ECS Prefix List can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_prefix_list.example <id>
```
