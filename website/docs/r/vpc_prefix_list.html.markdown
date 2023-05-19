---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_prefix_list"
sidebar_current: "docs-alicloud-resource-vpc-prefix-list"
description: |-
  Provides a Alicloud Vpc Prefix List resource.
---

# alicloud_vpc_prefix_list

Provides a Vpc Prefix List resource. This resource is used to create a prefix list.

For information about Vpc Prefix List and how to use it, see [What is Prefix List](https://www.alibabacloud.com/help/zh/virtual-private-cloud/latest/creatvpcprefixlist).

-> **NOTE:** Available in v1.182.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testacc-example"
}

resource "alicloud_resource_manager_resource_group" "defaultRg" {
  display_name        = "tf-testacc-chenyi"
  resource_group_name = var.name
}

resource "alicloud_resource_manager_resource_group" "changeRg" {
  display_name        = "tf-testacc-chenyi-change"
  resource_group_name = "${var.name}1"
}


resource "alicloud_vpc_prefix_list" "default" {
  max_entries             = 50
  resource_group_id       = alicloud_resource_manager_resource_group.defaultRg.id
  prefix_list_description = "test"
  ip_version              = "IPV4"
  prefix_list_name        = var.name
  entries {
    cidr        = "192.168.0.0/16"
    description = "test"
  }
}
```

## Argument Reference

The following arguments are supported:
* `entries` - (Optional, Computed, Available in v1.205.0+) The CIDR address block list of the prefix list.See the following `Block Entries`.
* `ip_version` - (Optional, ForceNew, Computed) The IP version of the prefix list. Value:-**IPV4**:IPv4 version.-**IPV6**:IPv6 version.
* `max_entries` - (Optional, Computed) The maximum number of entries for CIDR address blocks in the prefix list.
* `prefix_list_description` - (Optional) The description of the prefix list.It must be 2 to 256 characters in length and must start with a letter or Chinese, but cannot start with `http://` or `https://`.
* `prefix_list_name` - (Optional) The name of the prefix list. The name must be 2 to 128 characters in length, and must start with a letter. It can contain digits, periods (.), underscores (_), and hyphens (-).
* `resource_group_id` - (Optional, Computed, Available in v1.205.0+) The ID of the resource group to which the PrefixList belongs.
* `tags` - (Optional, Map, Available in v1.205.0+) The tags of PrefixList.

The following arguments will be discarded. Please use new fields as soon as possible:
* `entrys` - (Deprecated from v1.205.0+) Field 'entrys' has been deprecated from provider version 1.205.0. New field 'entries' instead.

#### Block Entries

The Entries supports the following:
* `cidr` - (Optional) The CIDR address block of the prefix list.
* `description` - (Optional) The description of the cidr entry. It must be 2 to 256 characters in length and must start with a letter or Chinese, but cannot start with `http://` or `https://`.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the prefix list was created.
* `prefix_list_association` - The association list information of the prefix list.
  * `owner_id` - The ID of the Alibaba Cloud account (primary account) to which the prefix list belongs.
  * `prefix_list_id` - The instance ID of the prefix list.
  * `reason` - Reason when the association fails.
  * `region_id` - The region ID of the prefix list to be queried.
  * `resource_id` - The ID of the associated resource.
  * `resource_type` - The associated resource type. Value:-**vpcRouteTable**: The VPC route table.-**trRouteTable**: the routing table of the forwarding router.
  * `resource_uid` - The ID of the Alibaba Cloud account (primary account) to which the resource bound to the prefix list belongs.
  * `status` - The association status of the prefix list. Value:-**Created**: Success.-**ModifyFailed**: is not associated with the latest version.-**Creating**: Creating.-**Modifying**: Modifying.-**Deleting**: Deleting.-**Deleted**: Deleted.
* `prefix_list_id` - The ID of the query Prefix List.
* `share_type` - The share type of the prefix list. Value:-**Shared**: indicates that the prefix list is a Shared prefix list.-Null: indicates that the prefix list is not a shared prefix list.
* `status` - Resource attribute fields that represent the status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Prefix List.
* `delete` - (Defaults to 5 mins) Used when delete the Prefix List.
* `update` - (Defaults to 5 mins) Used when update the Prefix List.

## Import

Vpc Prefix List can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_prefix_list.example <id>
```