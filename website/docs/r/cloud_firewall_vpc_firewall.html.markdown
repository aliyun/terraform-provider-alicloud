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

For information about Cloud Firewall Vpc Firewall and how to use it, see [What is Vpc Firewall](https://help.aliyun.com/document_detail/342893.html).

-> **NOTE:** Available in v1.194.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cloud_firewall_vpc_firewall" "default" {
  vpc_firewall_name = "tf-test"
  member_uid        = "1415189284827022"
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
* `lang` - (Optional) The language type of the requested and received messages. Value:**zh** (default): Chinese.**en**: English.
* `local_vpc` - (Required) The details of the local VPC. See the following `Block LocalVpc`.
* `member_uid` - (Optional) The UID of the Alibaba Cloud member account.
* `peer_vpc` - (Required) The details of the peer VPC. See the following `Block PeerVpc`.
* `status` - (Required) The status of the resource
  - `open` (default): protection is automatically enabled after the VPC boundary firewall is created.
  - `close`: Do not automatically enable protection after creating VPC boundary firewall
* `vpc_firewall_name` - (Required) The name of the VPC firewall instance.


#### Block LocalVpc

The LocalVpc supports the following:
* `eni_id` - (Computed) The ID of the instance of the Eni in the local VPC.
* `eni_private_ip_address` - (Computed) The private IP address of the elastic network card in the local VPC.
* `local_vpc_cidr_table_list` - (Required) The network segment list of the local VPC.See the following `Block LocalVpcCidrTableList`.
* `region_no` - (Required,ForceNew) The region ID of the local VPC.
* `router_interface_id` - (Computed) The ID of the router interface in the local VPC.
* `vpc_id` - (Required,ForceNew) The ID of the local VPC instance.
* `vpc_name` - (Computed) The instance name of the local VPC.

#### Block LocalVpcCidrTableList

The LocalVpcCidrTableList supports the following:
* `local_route_entry_list` - (Required) The list of route entries of the local VPC.See the following `Block LocalRouteEntryList`.
* `local_route_table_id` - (Required) The ID of the route table of the local VPC.

#### Block LocalRouteEntryList

The LocalRouteEntryList supports the following:
* `local_destination_cidr` - (Required) The target network segment of the local VPC.
* `local_next_hop_instance_id` - (Required) The ID of the next-hop instance in the local VPC.

#### Block PeerVpc

The PeerVpc supports the following:
* `eni_id` - (Computed) The ID of the instance of the ENI in the peer VPC.
* `eni_private_ip_address` - (Computed) The private IP address of the elastic network card in the peer VPC.
* `peer_vpc_cidr_table_list` - (Required) The network segment list of the peer VPC.See the following `Block PeerVpcCidrTableList`.
* `region_no` - (Required,ForceNew) The region ID of the peer VPC.
* `router_interface_id` - (Computed) The ID of the router interface in the peer VPC.
* `vpc_id` - (Required,ForceNew) The ID of the peer VPC instance.
* `vpc_name` - (Computed) The instance name of the peer VPC.

#### Block PeerVpcCidrTableList

The PeerVpcCidrTableList supports the following:
* `peer_route_entry_list` - (Required) Peer VPC route entry list information.See the following `Block PeerRouteEntryList`.
* `peer_route_table_id` - (Required) The ID of the route table of the peer VPC.

#### Block PeerRouteEntryList

The PeerRouteEntryList supports the following:
* `peer_destination_cidr` - (Required) The target network segment of the peer VPC.
* `peer_next_hop_instance_id` - (Required) The ID of the next-hop instance in the peer VPC.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the VPC firewall instance and the value same as `vpc_firewall_id`.
* `vpc_firewall_id` - The ID of the VPC firewall instance.
* `bandwidth` - Bandwidth specifications for high-speed channels. Unit: Mbps.
* `connect_type` - The communication type of the VPC firewall. Valid value: **expressconnect**, which indicates Express Connect.
* `region_status` - The region is open. Value:-**enable**: is enabled, indicating that VPC firewall can be configured in this region.-**disable**: indicates that VPC firewall cannot be configured in this region.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 31 mins) Used when create the Vpc Firewall.
* `delete` - (Defaults to 31 mins) Used when delete the Vpc Firewall.
* `update` - (Defaults to 31 mins) Used when update the Vpc Firewall.

## Import

Cloud Firewall Vpc Firewall can be imported using the id, e.g.

```shell
$terraform import alicloud_cloud_firewall_vpc_firewall.example <id>
```