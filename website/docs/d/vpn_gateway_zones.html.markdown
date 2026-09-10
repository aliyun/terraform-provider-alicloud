---
subcategory: "VPN Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_gateway_zones"
sidebar_current: "docs-alicloud-datasource-vpn-gateway-zones"
description: |-
  Provides a list of VPN Gateway Zone owned by an Alibaba Cloud account.
---

# alicloud_vpn_gateway_zones

This data source provides VPN Gateway Zone available to the user.[What is Zone](https://next.api.alibabacloud.com/api/Vpc/2016-04-28/DescribeVpnGatewayAvailableZones?lang=JAVA)

-> **NOTE:** Available since v1.216.0.

## Example Usage

```terraform
data "alicloud_vpn_gateway_zones" "default" {
  spec = "5M"
}
```

## Argument Reference

The following arguments are supported:
* `spec` - (Required, ForceNew) Bandwidth specification.-If an IPsec connection is bound to a VPN gateway instance, this parameter indicates the Bandwidth specification of the VPN gateway instance.-If an IPsec connection is bound to a forwarding router, this parameter indicates the bandwidth that you expect the IPsec connection to support.Different bandwidth specifications may affect the zone information that is found. Value:
  - `5M`
  - `10M`
  - `20M`
  - `50M`
  - `100M`
  - `200M`
  - `500M`
  - `1000M`
* `ids` - (Optional, ForceNew, Computed) A list of Zone IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Zone IDs.
* `zones` - A list of Zone Entries. Each element contains the following attributes:
  * `zone_id` - The zone ID.
  * `zone_name` - The zone name.
