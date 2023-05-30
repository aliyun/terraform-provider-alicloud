---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipv4_cidr_block"
sidebar_current: "docs-alicloud-resource-vpc-ipv4-cidr-block"
description: |-
  Provides a Alicloud VPC Ipv4 Cidr Block resource.
---

# alicloud\_vpc\_ipv4\_cidr\_block

Provides a VPC Ipv4 Cidr Block resource.

For information about VPC Ipv4 Cidr Block and how to use it, see [What is Ipv4 Cidr Block](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/associatevpccidrblock).

-> **NOTE:** Available in v1.185.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/24"
  vpc_name   = "terraform-example"
}
resource "alicloud_vpc_ipv4_cidr_block" "example" {
  vpc_id               = alicloud_vpc.default.id
  secondary_cidr_block = "192.163.0.0/16"
}
```

## Argument Reference

The following arguments are supported:

* `secondary_cidr_block` - (Required, ForceNew) The IPv4 CIDR block. Take note of the following requirements:
  * You can specify one of the following standard IPv4 CIDR blocks or their subnets as the secondary IPv4 CIDR block: 192.168.0.0/16, 172.16.0.0/12, and 10.0.0.0/8.
  * You can also use a custom CIDR block other than 100.64.0.0/10, 224.0.0.0/4, 127.0.0.0/8, 169.254.0.0/16, or their subnets as the secondary IPv4 CIDR block of the VPC.
  * The CIDR block cannot start with 0. The subnet mask must be 8 to 28 bits in length.
  * The secondary CIDR block cannot overlap with the primary CIDR block or an existing secondary CIDR block.
* `vpc_id` - (Required, ForceNew) The ID of the VPC.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Ipv4 Cidr Block. The value formats as `<vpc_id>:<secondary_cidr_block>`.

## Import

VPC Ipv4 Cidr Block can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_ipv4_cidr_block.example <vpc_id>:<secondary_cidr_block>
```