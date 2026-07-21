---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_peer_connection"
description: |-
  Provides a Alicloud Vpc Peer Peer Connection resource.
---

# alicloud_vpc_peer_connection

Provides a Vpc Peer Peer Connection resource.

Vpc peer connection.

For information about Vpc Peer Peer Connection and how to use it, see [What is Peer Connection](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/createvpcpeer).

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_vpc_peer_connection&spm=docs.r.vpc_peer_connection.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `accepting_ali_uid` - (Optional, ForceNew, Int) The ID of the Alibaba Cloud account to which the accepter VPC belongs.

*   To create a VPC peering connection within your Alibaba Cloud account, enter the ID of your Alibaba Cloud account.
*   To create a VPC peering connection between your Alibaba Cloud account and another Alibaba Cloud account, enter the ID of the peer Alibaba Cloud account.

-> **NOTE:**   If the accepter is a RAM user, set `AcceptingAliUid` to the ID of the Alibaba Cloud account that created the RAM user.

* `accepting_region_id` - (Required, ForceNew) The region ID of the accepter VPC of the VPC peering connection that you want to create.

  - To create an intra-region VPC peering connection, enter a region ID that is the same as that of the requester VPC.
  - To create an inter-region VPC peering connection, enter a region ID that is different from that of the requester VPC.
* `accepting_vpc_id` - (Required, ForceNew) The ID of the accepter VPC.
* `bandwidth` - (Optional, Computed, Int) The bandwidth of the VPC peering connection. Unit: Mbit/s. The value must be an integer greater than 0. Before you specify this parameter, make sure that you create an inter-region VPC peering connection.
* `description` - (Optional, Computed) The description of the VPC peering connection.
The description must be 2 to 256 characters in length. The description must start with a letter but cannot start with `http://` or `https://`.
* `dry_run` - (Optional) Specifies whether to perform only a dry run, without performing the actual request. Valid values:

  - `true`: performs only a dry run. The system checks the request for potential issues, including missing parameter values, incorrect request syntax, and service limits. If the request fails the dry run, an error message is returned. If the request passes the dry run, the `DryRunOperation` error code is returned.
  - `false` (default): performs a dry run and performs the actual request. If the request passes the dry run, a 2xx HTTP status code is returned and the operation is performed.
* `force_delete` - (Optional, Available since v1.231.0) Specifies whether to forcefully delete the VPC peering connection. Valid values:

  - `false` (default): no.
  - `true`: yes. If you forcefully delete the VPC peering connection, the system deletes the routes that point to the VPC peering connection from the VPC route table.
* `link_type` - (Optional, Available since v1.240.0) The link type of the VPC peering connection that you want to create. Valid values:
  - Platinum.
  - Gold: default value.

-> **NOTE:**  

-> **NOTE:**  - If you need to specify this parameter, ensure that the VPC peering connection is an inter-region connection.

* `peer_connection_name` - (Optional, Computed) The name of the VPC peering connection.
The name must be 2 to 128 characters in length, and can contain digits, underscores (\_), and hyphens (-). It must start with a letter.
* `resource_group_id` - (Optional, Computed) The ID of the new resource group.

-> **NOTE:**   You can use resource groups to manage resources within your Alibaba Cloud account by group. This helps you resolve issues such as resource grouping and permission management for your Alibaba Cloud account. For more information, see [What is resource management?](https://www.alibabacloud.com/help/en/doc-detail/94475.html)

* `status` - (Optional, Computed) The status of the resource
* `tags` - (Optional, Map) The tags of VpcPeer.
* `vpc_id` - (Required, ForceNew) The ID of the requester VPC or accepter VPC of the VPC peering connection that you want to query.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the VPC peer connection. Use UTC time in the format' YYYY-MM-DDThh:mm:ssZ '.
* `region_id` - The region ID of the resource to which you want to create and add tags.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Peer Connection.
* `delete` - (Defaults to 5 mins) Used when delete the Peer Connection.
* `update` - (Defaults to 5 mins) Used when update the Peer Connection.

## Import

Vpc Peer Peer Connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_peer_connection.example <id>
```