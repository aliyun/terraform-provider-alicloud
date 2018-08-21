---
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_customer_gateways"
sidebar_current: "docs-alicloud-datasource-vpn-customer-gateways"
description: |-
    Provides a list of VPN customer gateways which owned by an Alicloud account.
---

# alicloud\_vpn_customer_gateways

The VPN customers gateways data source lists a number of VPN customer gateways resource information owned by an Alicloud account.

## Example Usage

```
data "alicloud_vpn_customer_gateways" "cgws" {
	name_regex = "test_cgw"
	output_file = "/tmp/cgws"
}

```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string of VPN customer gateways name.
* `output_file` - (Optional) Save the result to the file.

## Attributes Reference

The following attributes are exported:

* `customer_gateways_id` - ID of the customer gateway .
* `name` - The name of the VPN customer gateway.
* `description` - The description of the VPN customer gateway.
* `ip_address` - The ip address of the customer gateway.

