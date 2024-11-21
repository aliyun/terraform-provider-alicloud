---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_peer_connection"
description: |-
  Provides a Alicloud Vpc Peer Connection resource.
---

# alicloud_vpc_peer_connection

Provides a Vpc Peer Connection resource.

For information about VPC Peer Connection and how to use it, see [What is Peer Connection](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/createvpcpeer).

-> **NOTE:** Available since v1.186.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_peer_connection&exampleId=294fed06-9b0d-e5fe-a093-4ebb1a7b8fe9e29c352b&activeTab=example&spm=docs.r.vpc_peer_connection.0.294fed069b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_account" "default" {}

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
* `accepting_ali_uid` - (Optional, ForceNew, Int) The ID of the Alibaba Cloud account (primary account) of the receiving end of the VPC peering connection to be created.
  - Enter the ID of your Alibaba Cloud account to create a peer-to-peer connection to the VPC account.
  - Enter the ID of another Alibaba Cloud account to create a cross-account VPC peer-to-peer connection.

-> **NOTE:**  If the recipient account is a RAM user (sub-account), enter the ID of the Alibaba Cloud account corresponding to the RAM user.

* `accepting_region_id` - (Required, ForceNew) The region ID of the recipient of the VPC peering connection to be created.
  - When creating a VPC peer-to-peer connection in the same region, enter the same region ID as the region ID of the initiator.
  - When creating a cross-region VPC peer-to-peer connection, enter a region ID that is different from the region ID of the initiator.
* `accepting_vpc_id` - (Required, ForceNew) The VPC ID of the receiving end of the VPC peer connection.
* `bandwidth` - (Optional, Int) The bandwidth of the VPC peering connection to be modified. Unit: Mbps. The value range is an integer greater than 0.
* `description` - (Optional) The description of the VPC peer connection to be created.

  It must be 2 to 256 characters in length and must start with a letter or Chinese, but cannot start with http:// or https.
* `dry_run` - (Optional) Whether to PreCheck only this request. Value:
  - `true`: The check request is sent without creating a VPC peer-to-peer connection. Check items include whether required parameters, request format, and business restrictions are filled in. If the check does not pass, the corresponding error is returned. If the check passes, the error code 'DryRunOperation' is returned '.
  - `false` (default): A normal request is sent. After checking, the HTTP 2xx status code is returned and the operation is performed directly.
* `force_delete` - (Optional, Available since v1.231.0) Whether to forcibly delete the VPC peering connection. Value:
  - `false` (default): Does not forcibly delete the VPC peering connection.
  - `true`: forcibly deletes the VPC peering connection. During the forced deletion, the system deletes the route entries in the VPC routing table that point to the VPC peering connection.
* `peer_connection_name` - (Optional) The name of the resource.
* `resource_group_id` - (Optional, Computed) The ID of resource group.
* `status` - (Optional, Computed) The status of the resource.
* `tags` - (Optional, Map) The tags of the resource.
* `vpc_id` - (Required, ForceNew) You must create a VPC ID on the initiator of a VPC peer connection.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the VPC peer connection. Use UTC time in the format' YYYY-MM-DDThh:mm:ssZ '.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Peer Connection.
* `delete` - (Defaults to 5 mins) Used when delete the Peer Connection.
* `update` - (Defaults to 5 mins) Used when update the Peer Connection.

## Import

Vpc Peer Connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_peer_connection.example <id>
```