---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_cidr"
sidebar_current: "docs-alicloud-resource-vpc-cidr"
description: |-
 Provides a resource to call AssociateVpcCidrBlock to add additional IPv4 network segments to the VPC.
---

# alicloud\_vpc\_cidr

 Provides a resource to call AssociateVpcCidrBlock to add additional IPv4 network segments to the VPC.

-> **NOTE:** A VPC only supports adding one additional IPv4 network segment. If you need more quota, please submit a ticket.

-> **NOTE:**  You can use the three standard network segments 192.168.0.0/16, 172.16.0.0/12, and 10.0.0.0/8 and their subnets as additional IPv4 network segments for the VPC.

If you want to use the public network address segment as an additional IPv4 network segment of the VPC, please submit a work order.

When adding additional IPv4 network segments, the following principles should be followed:
- It cannot start with 0, and the valid range of mask length is 8~24 bits.
- The additional network segment cannot overlap with the VPC main network segment and the added additional network segment.
For example, in a VPC where the primary IPv4 network segment is 192.168.0.0/16, you cannot add the following network segments as additional IPv4 network segments.                                         
  - A network segment larger than 192.168.0.0/16, such as 192.168.0.0/8.
  - The same network segment as the 192.168.0.0/16 range.
  - A network segment smaller than the 192.168.0.0/16 range, such as 192.168.0.0/24.


## Example Usage

Basic usage

```
resource "alicloud_vpc" "default" {
  name       = "tf_testCidrCheckVpcName"
  cidr_block = "172.16.0.0/12"
}
resource "alicloud_vpc_cidr_block" "default" {
  vpc_id = alicloud_vpc.default.id
  secondary_cidr_block = "10.0.0.0/8"
}
data "alicloud_vpcs" "vpcs" {
  ids = [alicloud_vpc.default.id]
  is_default = false
  name_regex = alicloud_vpc.default.name
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The VPC ID.
* `secondary_cidr_block` -  Additional IPv4 network segment to be added.
                            
## Attributes Reference

The following attributes are exported:

* `vpc_id` - The VPC ID                          

