---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_peer_connection_accepter"
sidebar_current: "docs-alicloud-resource-vpc-peer-connection-accepter"
description: |-
  Provides a Alicloud Vpc Peer Connection Accepter resource.
---

# alicloud_vpc_peer_connection_accepter

Provides a Vpc Peer Connection Accepter resource.

For information about Vpc Peer Connection Accepter and how to use it, see [What is Peer Connection Accepter](https://www.alibabacloud.com/help/en/vpc/developer-reference/api-vpcpeer-2022-01-01-acceptvpcpeerconnection).

-> **NOTE:** Available since v1.196.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_vpc_peer_connection_accepter&exampleId=8204ef10-2d10-b925-a0ac-68de784036d01a919786&activeTab=example&spm=docs.r.vpc_peer_connection_accepter.0.8204ef102d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
variable "accepting_region" {
  default = "cn-beijing"
}
variable "another_uid" {
  default = "xxxx"
}
# Method 1: Use assume_role to operate resources in the target account, detail see https://registry.terraform.io/providers/aliyun/alicloud/latest/docs#assume-role
provider "alicloud" {
  region = var.accepting_region
  alias  = "accepting"
  assume_role {
    role_arn = "acs:ram::${var.another_uid}:role/terraform-example-assume-role"
  }
}

# Method 2: Use the target account's access_key, secret_key
# provider "alicloud" {
#   region     = "cn-hangzhou"
#   access_key = "access_key"
#   secret_key = "secret_key"
#   alias      = "accepting"
# }

provider "alicloud" {
  alias  = "local"
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "local" {
  provider   = alicloud.local
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vpc" "accepting" {
  provider   = alicloud.accepting
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

data "alicloud_account" "accepting" {
  provider = alicloud.accepting
}
resource "alicloud_vpc_peer_connection" "default" {
  provider             = alicloud.local
  peer_connection_name = var.name
  vpc_id               = alicloud_vpc.local.id
  accepting_ali_uid    = data.alicloud_account.accepting.id
  accepting_region_id  = var.accepting_region
  accepting_vpc_id     = alicloud_vpc.accepting.id
  description          = var.name
}

resource "alicloud_vpc_peer_connection_accepter" "default" {
  provider    = alicloud.accepting
  instance_id = alicloud_vpc_peer_connection.default.id
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (Required, ForceNew) The ID of the instance of the created VPC peer connection.
* `dry_run` - (Optional) The dry run.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `accepting_owner_uid` - The ID of the Alibaba Cloud account (primary account) of the receiving end of the VPC peering connection to be created.-Enter the ID of your Alibaba Cloud account to create a peer-to-peer connection to the VPC account.-Enter the ID of another Alibaba Cloud account to create a cross-account VPC peer-to-peer connection.> If the recipient account is a RAM user (sub-account), enter the ID of the Alibaba Cloud account corresponding to the RAM user.
* `accepting_region_id` - The region ID of the recipient of the VPC peering connection to be created.-When creating a VPC peer-to-peer connection in the same region, enter the same region ID as the region ID of the initiator.-When creating a cross-region VPC peer-to-peer connection, enter a region ID that is different from the region ID of the initiator.
* `accepting_vpc_id` - The VPC ID of the receiving end of the VPC peer connection.
* `bandwidth` - The bandwidth of the VPC peering connection to be modified. Unit: Mbps. The value range is an integer greater than 0.
* `description` - The description of the VPC peer connection to be created.It must be 2 to 256 characters in length and must start with a letter or Chinese, but cannot start with http:// or https.
* `peer_connection_accepter_name` - The name of the resource
* `status` - The status of the resource
* `vpc_id` - You must create a VPC ID on the initiator of a VPC peer connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Peer Connection Accepter.
* `delete` - (Defaults to 5 mins) Used when delete the Peer Connection Accepter.

## Import

Vpc Peer Connection Accepter can be imported using the id, e.g.

```shell
$terraform import alicloud_vpc_peer_connection_accepter.example <id>
```