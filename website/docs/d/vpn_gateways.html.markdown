---
subcategory: "VPN Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_gateways"
sidebar_current: "docs-alicloud-datasource-vpn-gateways"
description: |-
    Provides a list of VPNs which owned by an Alicloud account.
---

# alicloud_vpn_gateways

The VPNs data source lists a number of VPNs resource information owned by an Alicloud account.

-> **NOTE:** Available since v1.18.0.

## Example Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vpn_gateway" "default" {
  name                 = var.name
  vpc_id               = data.alicloud_vpcs.default.ids.0
  bandwidth            = "10"
  enable_ssl           = true
  enable_ipsec         = true
  instance_charge_type = "PrePaid"
  description          = var.name
  vswitch_id           = data.alicloud_vswitches.default.ids.0
  network_type         = "public"
}

data "alicloud_vpn_gateways" "vpn_gateways" {
  vpc_id                   = data.alicloud_vpcs.default.ids.0
  ids                      = [alicloud_vpn_gateway.default.id]
  status                   = "Active"
  business_status          = "Normal"
  name_regex               = "tf-example"
  include_reservation_data = true
  output_file              = "/tmp/vpns"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Optional, ForceNew) Use the VPC ID as the search key.
* `ids` - (Optional, ForceNew) IDs of the VPN.
* `status` - (Optional, ForceNew) Limit search to specific status - valid value is "Init", "Provisioning", "Active", "Updating", "Deleting".
* `business_status` - (Optional, ForceNew) Limit search to specific business status - valid value is "Normal", "FinancialLocked".
* `name_regex` - (Optional, ForceNew) A regex string of VPN name.
* `output_file` - (Optional) Save the result to the file.
* `enable_ipsec` - (Deprecated, Optional, Available 1.161.0+, has been deprecated from provider version 1.193.0, it will be removed in the future version.) Indicates whether the IPsec-VPN feature is enabled.
* `include_reservation_data` - (Optional, ForceNew, Available 1.193.0+) Include ineffective ordering data.

## Attributes Reference

The following attributes are exported:

* `ids` - IDs of the VPN.
* `names` - names of the VPN.
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
  * `instance_charge_type` - The charge type of the VPN gateway.
  * `enable_ipsec` - Whether the ipsec function is enabled.
  * `enable_ssl` - Whether the ssl function is enabled.
  * `ssl_connections` - Total count of ssl vpn connections.
  * `network_type` - The network type of the VPN gateway.
  * `auto_propagate` - Whether to automatically propagate BGP routes to the VPC. Valid values: `true`, `false`.
