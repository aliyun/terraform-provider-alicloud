---
subcategory: "VPN Gateway"
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
data "alicloud_vpn_customer_gateways" "foo" {
  name_regex          = "testAcc*"
  ids                 = ["fake-id1", "fake-id2"] 
  output_file         = "/tmp/cgws"
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) ID of the VPN customer gateways.
* `name_regex` - (Optional) A regex string of VPN customer gateways name.
* `output_file` - (Optional) Save the result to the file.

## Attributes Reference

The following attributes are exported:

* `ids` IDs of VPN customer gateway.
* `names` names of VPN customer gateway.
* `gateways` - A list of VPN customer gateways. Each element contains the following attributes:
  * `id` - ID of the VPN customer gateway .
  * `name` - The name of the VPN customer gateway.
  * `description` - The description of the VPN customer gateway.
  * `ip_address` - The ip address of the VPN customer gateway.
  * `create_time` - The creation time of the VPN customer gateway.

