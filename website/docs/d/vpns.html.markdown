---
layout: "alicloud"
page_title: "Alicloud: alicloud_vpns"
sidebar_current: "docs-alicloud-datasource-vpns"
description: |-
    Provides a list of VPNs which owned by an Alicloud account.
---

# alicloud\_vpns

The VPNs data source lists a number of VPNs resource information owned by an Alicloud account.

## Example Usage

```
data "alicloud_vpns" "foo" {
	vpc_id = "vpc-2ze9wy916mfwpwbf6hx4u"
	output_file = "/tmp/vpns"
}

```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Optional) Use the VPC ID as the search key.
* `status` - (Optional) Limit search to specific status - valid value is "init","provisioning","active","updating","deleting".
* `name_regex` - (Optional) A regex string of VPN name.
* `output_file` - (Optional) Save the result to the file.

## Attributes Reference

The following attributes are exported:

* `vpn_gateway_id` - ID of the VPN.
* `vpc_id` - ID of the VPC that the VPN belongs.
* `vswitch_id` - ID of the Switch that the VPN belongs.
* `internet_ip` - The internet ip of the VPN.
* `spec` - The Specification of the VPN
* `name` - The name of the VPN.
* `description` - The description of the VPN
* `status` - The status of the VPN
