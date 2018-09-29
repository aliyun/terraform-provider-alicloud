---
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_gateways"
sidebar_current: "docs-alicloud-datasource-vpn-gateways"
description: |-
    Provides a list of VPNs which owned by an Alicloud account.
---

# alicloud\_vpn_gateways

The VPNs data source lists a number of VPNs resource information owned by an Alicloud account.

## Example Usage

```
data "alicloud_vpn_gateways" "vpn_gateways" {
	vpc_id = "fake-vpc-id"
	vpn_gateway_id = "fake-vpn-id"
	status = "active"
	business_status = "Normal"
	name_regex = "testAcc*"
	output_file = "/tmp/vpns"
}

```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Optional) Use the VPC ID as the search key.
* `ids` - (Optional) IDs of the VPN.
* `status` - (Optional) Limit search to specific status - valid value is "Init", "Provisioning", "Active", "Updating", "Deleting".
* `business_status` - (Optional) Limit search to specific business status - valid value is "Normal", "FinancialLocked".
* `name_regex` - (Optional) A regex string of VPN name.
* `output_file` - (Optional) Save the result to the file.

## Attributes Reference

The following attributes are exported:

* `gateways` - A list of VPN gateways. Each element contains the following attributes:
  * `id` - ID of the VPN.
  * `vpc_id` - ID of the VPC that the VPN belongs.
  * `internet_ip` - The internet ip of the VPN.
  * `create_time` - The creation time of the VPN gateway.
  * `end_time` - The expiration time of the VPN gateway.
  * `specification` - The Specification of the VPN
  * `name` - The name of the VPN.
  * `description` - The description of the VPN
  * `status` - The status of the VPN
  * `business_status` - The business status of the VPN gateway.
