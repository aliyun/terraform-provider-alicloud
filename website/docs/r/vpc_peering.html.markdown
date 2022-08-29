---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_peering"
sidebar_current: "docs-alicloud-resource-vpc-peering"
description: |-
  Provides a Alicloud VPC Peering resource.
---

# alicloud\_vpc\_peering

Provides a VPC Peering resource.

For information about VPC Peering and how to use it, see [What is Peering](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/createvpcpeer).

-> **NOTE:** Available in v1.184.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_account" "default" {}

data "alicloud_vpcs" "default" {}

resource "alicloud_vpc_peering" "default" {
  peering_name        = var.name
  vpc_id              = data.alicloud_vpcs.default.ids.0
  accepting_ali_uid   = data.alicloud_account.default.id
  accepting_region_id = "cn-hangzhou"
  accepting_vpc_id    = data.alicloud_vpcs.default.ids.1
  description         = var.name
}
```

## Argument Reference

The following arguments are supported:

* `accepting_ali_uid` - (Optional, Computed, ForceNew) The ID of the Alibaba Cloud account (primary account) of the receiving end of the VPC peering connection to be created.
  - Enter the ID of your Alibaba Cloud account to create a peer-to-peer connection to the VPC account.
  - Enter the ID of another Alibaba Cloud account to create a cross-account VPC peer-to-peer connection.
  - If the recipient account is a RAM user (sub-account), enter the ID of the Alibaba Cloud account corresponding to the RAM user.
* `accepting_region_id` - (Required, ForceNew) The region ID of the recipient of the VPC peering connection to be created.
  - When creating a VPC peer-to-peer connection in the same region, enter the same region ID as the region ID of the initiator.
  - When creating a cross-region VPC peer-to-peer connection, enter a region ID that is different from the region ID of the initiator.
* `accepting_vpc_id` - (Required, ForceNew) The VPC ID of the receiving end of the VPC peer connection.
* `bandwidth` - (Optional, Computed) The bandwidth of the VPC peering connection to be modified. Unit: Mbps. The value range is an integer greater than 0.
* `description` - (Optional) The description of the VPC peer connection to be created. It must be 2 to 256 characters in length and must start with a letter or Chinese, but cannot start with `http://` or `https://`.
* `dry_run` - (Optional) The dry run.
* `peering_name` - (Optional) The name of the resource. The name must be 2 to 128 characters in length, and must start with a letter. It can contain digits, underscores (_), and hyphens (-).
* `vpc_id` - (Required, ForceNew) The ID of the requester VPC.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Peering.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Peering.
* `delete` - (Defaults to 1 mins) Used when delete the Peering.
* `update` - (Defaults to 1 mins) Used when update the Peering.

## Import

VPC Peering can be imported using the id, e.g.

```
$ terraform import alicloud_vpc_peering.example <id>
```