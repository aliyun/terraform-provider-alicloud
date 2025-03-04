---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_peer_connection_accepter"
description: |-
  Provides a Alicloud Vpc Peer Peer Connection Accepter resource.
---

# alicloud_vpc_peer_connection_accepter

Provides a Vpc Peer Peer Connection Accepter resource.

Vpc peer connection receiver.

For information about Vpc Peer Peer Connection Accepter and how to use it, see [What is Peer Connection Accepter](https://www.alibabacloud.com/help/en/vpc/developer-reference/api-vpcpeer-2022-01-01-acceptvpcpeerconnection).

-> **NOTE:** Available since v1.196.0.

## Example Usage

Basic Usage

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
* `bandwidth` - (Optional, Computed, Int, Available since v1.231.0) The new bandwidth of the VPC peering connection. Unit: Mbit/s. The value must be an integer greater than 0.
* `description` - (Optional, Computed, Available since v1.231.0) The new description of the VPC peering connection.
The description must be 1 to 256 characters in length, and cannot start with `http://` or `https://`.
* `dry_run` - (Optional) Specifies whether to perform only a dry run, without performing the actual request. Valid values:

  - `true`: performs only a dry run. The system checks the request for potential issues, including missing parameter values, incorrect request syntax, and service limits. If the request fails the dry run, an error message is returned. If the request passes the dry run, the `DryRunOperation` error code is returned.
  - `false` (default): performs a dry run and performs the actual request. If the request passes the dry run, a 2xx HTTP status code is returned and the operation is performed.
* `force_delete` - (Optional, Available since v1.231.0) Specifies whether to forcefully delete the VPC peering connection. Valid values:

  - `false` (default): no.
  - `true`: yes. If you forcefully delete the VPC peering connection, the system deletes the routes that point to the VPC peering connection from the VPC route table.
* `instance_id` - (Required, ForceNew) The ID of the VPC peering connection whose name or description you want to modify.
* `link_type` - (Optional, Computed, Available since v1.240.0) Link Type
* `peer_connection_accepter_name` - (Optional, Computed, Available since v1.231.0) The new name of the VPC peering connection.
The name must be 1 to 128 characters in length, and cannot start with `http://` or `https://`.
* `resource_group_id` - (Optional, Computed) The ID of the new resource group.

-> **NOTE:**   You can use resource groups to manage resources within your Alibaba Cloud account by group. This helps you resolve issues such as resource grouping and permission management for your Alibaba Cloud account. For more information, see [What is resource management?](https://www.alibabacloud.com/help/en/doc-detail/94475.html)


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `accepting_owner_uid` - The ID of the Alibaba Cloud account (primary account) of the receiving end of the VPC peering connection to be created.-to-peer connection to the VPC account.-account VPC peer-to-peer connection.
* `accepting_region_id` - The region ID of the recipient of the VPC peering connection to be created.-to-peer connection in the same region, enter the same region ID as the region ID of the initiator.-region VPC peer-to-peer connection, enter a region ID that is different from the region ID of the initiator.
* `accepting_vpc_id` - The VPC ID of the receiving end of the VPC peer connection.
* `create_time` - The creation time of the VPC peer connection. Use UTC time in the format' YYYY-MM-DDThh:mm:ssZ '.
* `region_id` - The ID of the region where you want to query VPC peering connections.
* `status` - The status of the resource
* `vpc_id` - The VPC ID of the initiator of the VPC peering connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Peer Connection Accepter.
* `delete` - (Defaults to 5 mins) Used when delete the Peer Connection Accepter.
* `update` - (Defaults to 5 mins) Used when update the Peer Connection Accepter.

## Import

Vpc Peer Peer Connection Accepter can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_peer_connection_accepter.example <id>
```