---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipv6_cidr_block"
sidebar_current: "docs-alicloud-resource-vpc-ipv6-cidr-block"
description: |-
  Provides a Alicloud VPC Ipv6 Cidr Block resource.
---

# alicloud\_vpc\_ipv6\_cidr\_block

Provides a VPC Ipv6 Cidr Block resource.

For information about VPC Ipv6 Cidr Block and how to use it, see [What is Ipv6 Cidr Block](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/associatevpccidrblock).

-> **NOTE:** Available in v1.185.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vpc_ipv6_cidr_block" "example" {
  vpc_id = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `secondary_cidr_block` - (Required,ForceNew) The secondary IPv6 CIDR block.
* `ipv6_isp` - (Optional, Computed, ForceNew) The IPv6 address segment type of the VPC. Valid values: `BGP`, `ChinaMobile`, `ChinaUnicom`, `ChinaTelecom`, `ChinaTelecom`.
  - BGP (default): Alibaba Cloud BGP IPv6.
  - ChinaMobile: China Mobile (single line).
  - ChinaUnicom: China Unicom (single line).
  - ChinaTelecom: China Telecom (single line).
  - If a single-line bandwidth whitelist is enabled, the field can be set to `ChinaTelecom` (China Telecom), `ChinaUnicom` (China Unicom), and `ChinaMobile` (China Mobile). 
* `vpc_id` - (Required, ForceNew) The ID of the VPC.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Ipv6 Cidr Block. The value formats as `<vpc_id>:<secondary_cidr_block>`.

## Import

VPC Ipv6 Cidr Block can be imported using the id, e.g.

```
$ terraform import alicloud_vpc_ipv6_cidr_block.example <vpc_id>:<secondary_cidr_block>
```