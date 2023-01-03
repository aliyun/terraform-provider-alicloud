---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_prefix_list"
sidebar_current: "docs-alicloud-resource-vpc-prefix-list"
description: |-
  Provides a Alicloud VPC Prefix List resource.
---

# alicloud\_vpc\_prefix\_list

Provides a VPC Prefix List resource.

For information about VPC Prefix List and how to use it, see [What is Prefix List](https://www.alibabacloud.com/help/zh/virtual-private-cloud/latest/creatvpcprefixlist).

-> **NOTE:** Available in v1.182.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vpc_prefix_list" "default" {
  entrys {
    cidr        = "192.168.0.0/16"
    description = "description"
  }
  ip_version              = "IPV4"
  max_entries             = 50
  prefix_list_name        = var.name
  prefix_list_description = "description"
}
```

## Argument Reference

The following arguments are supported:

* `entrys` - (Optional) The CIDR address block list of the prefix list. See the following `Block entrys`.
* `ip_version` - (Optional, Computed, ForceNew) The IP version of the prefix list. Valid values: `IPV4`, `IPV6`.
* `max_entries` - (Optional, Computed) The maximum number of entries for CIDR address blocks in the prefix list.
* `prefix_list_description` - (Optional) The description of the prefix list. It must be 2 to 256 characters in length and must start with a letter or Chinese, but cannot start with `http://` or `https://`.
* `prefix_list_name` - (Optional) The name of the prefix list. The name must be 2 to 128 characters in length and must start with a letter. It can contain digits, periods (.), underscores (_), and hyphens (-).

#### Block entrys

The entrys supports the following: 

* `cidr` - (Optional) The CIDR address block of the prefix list.
* `description` - (Optional) The description of the cidr entry. It must be 2 to 256 characters in length and must start with a letter or Chinese, but cannot start with `http://` or `https://`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Prefix List.
* `status` - (Available in v1.196.0+) The status of the Prefix List.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when creating the Prefix List.
* `update` - (Defaults to 3 mins) Used when updating the Prefix List.
* `delete` - (Defaults to 3 mins) Used when deleting the Prefix List.


## Import

VPC Prefix List can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_prefix_list.example <id>
```
