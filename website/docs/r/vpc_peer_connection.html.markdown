---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_peer_connection"
description: |-
  Provides a Alicloud VPC Peer Connection resource.
---

# alicloud_vpc_peer_connection

Provides a VPC Peer Connection resource.

For information about VPC Peer Connection and how to use it, see [What is Peer Connection](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/createvpcpeer).

-> **NOTE:** Available since v1.186.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_account" "default" {
}

variable "accepting_region" {
  default = "cn-beijing"
}

provider "alicloud" {
  alias  = "local"
  region = "cn-hangzhou"
}

provider "alicloud" {
  alias  = "accepting"
  region = var.accepting_region
}

resource "alicloud_vpc" "local_vpc" {
  provider   = alicloud.local
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vpc" "accepting_vpc" {
  provider   = alicloud.accepting
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vpc_peer_connection" "default" {
  provider             = alicloud.local
  peer_connection_name = "terraform-example"
  vpc_id               = alicloud_vpc.local_vpc.id
  accepting_ali_uid    = data.alicloud_account.default.id
  accepting_region_id  = var.accepting_region
  accepting_vpc_id     = alicloud_vpc.accepting_vpc.id
  description          = "terraform-example"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The ID of the requester VPC.
* `accepting_vpc_id` - (Required, ForceNew) The VPC ID of the receiving end of the VPC peer connection.
* `accepting_region_id` - (Required, ForceNew) The region ID of the recipient of the VPC peering connection to be created.
  - When creating a VPC peer-to-peer connection in the same region, enter the same region ID as the region ID of the initiator.
  - When creating a cross-region VPC peer-to-peer connection, enter a region ID that is different from the region ID of the initiator.
* `accepting_ali_uid` - (Optional, ForceNew, Int) The ID of the Alibaba Cloud account (primary account) of the receiving end of the VPC peering connection to be created.
  - Enter the ID of your Alibaba Cloud account to create a peer-to-peer connection to the VPC account.
  - Enter the ID of another Alibaba Cloud account to create a cross-account VPC peer-to-peer connection.
-> **NOTE:**  If the recipient account is a RAM user (sub-account), enter the ID of the Alibaba Cloud account corresponding to the RAM user.
* `bandwidth` - (Optional, Int) The bandwidth of the VPC peering connection to be modified. Unit: Mbps. The value range is an integer greater than 0.
* `resource_group_id` - (Optional, Available since v1.207.0) The ID of the resource group.
* `peer_connection_name` - (Optional) The name of the VPC peer connection. The name of the resource. The name must be 2 to 128 characters in length, and must start with a letter. It can contain digits, underscores (_), and hyphens (-).
* `description` - (Optional) The description of the VPC peer connection to be created.It must be 2 to 256 characters in length and must start with a letter or Chinese, but cannot start with `http://` or `https://`.
* `status` - (Optional) The status of the VPC peer connection.
* `tags` - (Optional, Map, Available since v1.207.0) A mapping of tags to assign to the resource.
* `dry_run` - (Optional, Bool) Whether to PreCheck only this request. Default value: `false`. Valid values:
  - `true`: The check request is sent without creating a VPC peer-to-peer connection. Check items include whether required parameters, request format, and business restrictions are filled in. If the check does not pass, the corresponding error is returned. If the check passes, the error code `DryRunOperation` is returned.
  - `false`: A normal request is sent. After checking, the HTTP 2xx status code is returned and the operation is performed directly.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the VPC peer connection. Use UTC time in the format `YYYY-MM-DDThh:mm:ssZ`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Peer Connection.
* `update` - (Defaults to 5 mins) Used when update the Peer Connection.
* `delete` - (Defaults to 5 mins) Used when delete the Peer Connection.

## Import

VPC Peer Connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_peer_connection.example <id>
```
