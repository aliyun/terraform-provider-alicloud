---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipv6_internet_bandwidths"
sidebar_current: "docs-alicloud-datasource-vpc-ipv6-internet-bandwidths"
description: |-
  Provides a list of Vpc Ipv6 Internet Bandwidths to the user.
---

# alicloud\_vpc\_ipv6\_internet\_bandwidths

This data source provides the Vpc Ipv6 Internet Bandwidths of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.143.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_ipv6_internet_bandwidths" "ids" {
  ids = ["example_id"]
}
output "vpc_ipv6_internet_bandwidth_id_1" {
  value = data.alicloud_vpc_ipv6_internet_bandwidths.ids.bandwidths.0.id
}

data "alicloud_vpc_ipv6_internet_bandwidths" "ipv6InternetBandwidthId" {
  ipv6_internet_bandwidth_id = "example_value"
}
output "vpc_ipv6_internet_bandwidth_id_2" {
  value = data.alicloud_vpc_ipv6_internet_bandwidths.ipv6InternetBandwidthId.bandwidths.0.id
}

data "alicloud_vpc_ipv6_internet_bandwidths" "ipv6AddressId" {
  ipv6_address_id = "example_value"
}
output "vpc_ipv6_internet_bandwidth_id_3" {
  value = data.alicloud_vpc_ipv6_internet_bandwidths.ipv6AddressId.bandwidths.0.id
}

data "alicloud_vpc_ipv6_internet_bandwidths" "status" {
  status = "Normal"
}
output "vpc_ipv6_internet_bandwidth_id_4" {
  value = data.alicloud_vpc_ipv6_internet_bandwidths.status.bandwidths.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ipv6_internet_bandwidth_id` - (Optional, ForceNew) The ID of the Ipv6 Internet Bandwidth.
* `ipv6_address_id` - (Optional, ForceNew) The ID of the IPv6 address.
* `ids` - (Optional, ForceNew, Computed)  A list of Ipv6 Internet Bandwidth IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Normal`, `FinancialLocked` and `SecurityLocked`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Ipv6 Internet Bandwidth names.
* `bandwidths` - A list of Vpc Ipv6 Internet Bandwidths. Each element contains the following attributes:
	* `bandwidth` - The amount of Internet bandwidth resources of the IPv6 address, Unit: `Mbit/s`.
	* `id` - The ID of the Ipv6 Internet Bandwidth.
	* `internet_charge_type` - The metering method of the Internet bandwidth resources of the IPv6 gateway.
	* `ipv6_address_id` - The ID of the IPv6 address.
	* `ipv6_gateway_id` - The ID of the IPv6 gateway.
	* `ipv6_internet_bandwidth_id` - The ID of the Ipv6 Internet Bandwidth.
	* `payment_type` - The payment type of the resource.
	* `status` -  The status of the resource. Valid values: `Normal`, `FinancialLocked` and `SecurityLocked`.