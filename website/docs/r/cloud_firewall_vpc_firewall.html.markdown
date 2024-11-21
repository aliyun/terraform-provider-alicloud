---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_vpc_firewall"
sidebar_current: "docs-alicloud-resource-cloud-firewall-vpc-firewall"
description: |-
  Provides a Alicloud Cloud Firewall Vpc Firewall resource.
---

# alicloud_cloud_firewall_vpc_firewall

Provides a Cloud Firewall Vpc Firewall resource.

For information about Cloud Firewall Vpc Firewall and how to use it, see [What is Vpc Firewall](https://www.alibabacloud.com/help/en/cloud-firewall/developer-reference/api-cloudfw-2017-12-07-createvpcfirewallconfigure).

-> **NOTE:** Available since v1.194.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_vpc_firewall&exampleId=8fb2ff2c-239a-dc83-a522-95e488c8392a8702ef57&activeTab=example&spm=docs.r.cloud_firewall_vpc_firewall.0.8fb2ff2c23&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_account" "current" {
}

resource "alicloud_cloud_firewall_vpc_firewall" "default" {
  vpc_firewall_name = "tf-example"
  member_uid        = data.alicloud_account.current.id
  local_vpc {
    vpc_id    = "vpc-bp1d065m6hzn1xbw8ibfd"
    region_no = "cn-hangzhou"
    local_vpc_cidr_table_list {
      local_route_table_id = "vtb-bp1lj0ddg846856chpzrv"
      local_route_entry_list {
        local_next_hop_instance_id = "ri-bp1uobww3aputjlwwkyrh"
        local_destination_cidr     = "10.1.0.0/16"
      }
    }
  }
  peer_vpc {
    vpc_id    = "vpc-bp1gcmm64o3caox84v0nz"
    region_no = "cn-hangzhou"
    peer_vpc_cidr_table_list {
      peer_route_table_id = "vtb-bp1f516f2hh4sok1ig9b5"
      peer_route_entry_list {
        peer_destination_cidr     = "10.0.0.0/16"
        peer_next_hop_instance_id = "ri-bp1thhtgf6ydr2or52l3n"
      }
    }
  }
  status = "open"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_firewall_name` - (Required) The name of the VPC firewall instance.
* `status` - (Required) The status of the resource. Valid values:
  - `open`: protection is automatically enabled after the VPC boundary firewall is created.
  - `close`: Do not automatically enable protection after creating VPC boundary firewall.
* `member_uid` - (Optional, ForceNew) The UID of the Alibaba Cloud member account.
* `lang` - (Optional) The language type of the requested and received messages. Valid values:
  - `zh`: Chinese.
  - `en`: English.
* `local_vpc` - (Required, ForceNew, Set) The details of the local VPC. See [`local_vpc`](#local_vpc) below.
* `peer_vpc` - (Required, ForceNew, Set) The details of the peer VPC. See [`peer_vpc`](#peer_vpc) below.

### `local_vpc`

The local_vpc supports the following:

* `vpc_id` - (Required, ForceNew) The ID of the local VPC instance.
* `region_no` - (Required, ForceNew) The region ID of the local VPC.
* `local_vpc_cidr_table_list` - (Required, ForceNew, Set) The network segment list of the local VPC. See [`local_vpc_cidr_table_list`](#local_vpc-local_vpc_cidr_table_list) below.

### `local_vpc-local_vpc_cidr_table_list`

The local_vpc_cidr_table_list supports the following:

* `local_route_table_id` - (Required, ForceNew) The ID of the route table of the local VPC.
* `local_route_entry_list` - (Required, ForceNew, Set) The list of route entries of the local VPC. See [`local_route_entry_list`](#local_vpc-local_vpc_cidr_table_list-local_route_entry_list) below.

### `local_vpc-local_vpc_cidr_table_list-local_route_entry_list`

The local_route_entry_list supports the following:

* `local_next_hop_instance_id` - (Required, ForceNew) The ID of the next-hop instance in the local VPC.
* `local_destination_cidr` - (Required, ForceNew) The target network segment of the local VPC.

### `peer_vpc`

The peer_vpc supports the following:

* `vpc_id` - (Required, ForceNew) The ID of the peer VPC instance.
* `region_no` - (Required, ForceNew) The region ID of the peer VPC.
* `peer_vpc_cidr_table_list` - (Required, ForceNew, Set) The network segment list of the peer VPC. See [`peer_vpc_cidr_table_list`](#peer_vpc-peer_vpc_cidr_table_list) below.

### `peer_vpc-peer_vpc_cidr_table_list`

The peer_vpc_cidr_table_list supports the following:

* `peer_route_table_id` - (Required, ForceNew) The ID of the route table of the peer VPC.
* `peer_route_entry_list` - (Required, ForceNew, Set) Peer VPC route entry list information. See [`peer_route_entry_list`](#peer_vpc-peer_vpc_cidr_table_list-peer_route_entry_list) below.

### `peer_vpc-peer_vpc_cidr_table_list-peer_route_entry_list`

The peer_route_entry_list supports the following:

* `peer_next_hop_instance_id` - (Required, ForceNew) The ID of the next-hop instance in the peer VPC.
* `peer_destination_cidr` - (Required, ForceNew) The target network segment of the peer VPC.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID of the Vpc Firewall. The value formats as `vpc_firewall_id`.
* `vpc_firewall_id` - The ID of the VPC firewall instance.
* `connect_type` - The communication type of the VPC firewall.
* `bandwidth` - Bandwidth specifications for high-speed channels. Unit: Mbps.
* `region_status` - The region is open.
* `local_vpc` - The details of the Local VPC.
  * `vpc_name` - The instance name of the local VPC.
  * `eni_id` - The ID of the instance of the Eni in the local VPC.
  * `eni_private_ip_address` - The private IP address of the elastic network card in the local VPC.
  * `router_interface_id` - The ID of the router interface in the local VPC.
* `peer_vpc` - The details of the Peer VPC.
  * `vpc_name` - The instance name of the peer VPC.
  * `eni_id` - The ID of the instance of the ENI in the peer VPC.
  * `eni_private_ip_address` - The private IP address of the elastic network card in the peer VPC.
  * `router_interface_id` - The ID of the router interface in the peer VPC.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 31 mins) Used when create the Vpc Firewall.
* `update` - (Defaults to 31 mins) Used when update the Vpc Firewall.
* `delete` - (Defaults to 31 mins) Used when delete the Vpc Firewall.

## Import

Cloud Firewall Vpc Firewall can be imported using the id, e.g.

```shell
$terraform import alicloud_cloud_firewall_vpc_firewall.example <id>
```
