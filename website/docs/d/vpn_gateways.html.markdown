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
  default = "terraform-example"
}

provider "alicloud" {
  region = "me-east-1"
}

variable "spec" {
  default = "20"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "me-east-1a"
}

resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = "me-east-1a"
  vswitch_name = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_vpn_gateway" "default" {
  vpn_type         = "Normal"
  vpn_gateway_name = var.name

  vswitch_id   = local.vswitch_id
  auto_pay     = true
  vpc_id       = data.alicloud_vpcs.default.ids.0
  network_type = "public"
  payment_type = "Subscription"
  enable_ipsec = true
  bandwidth    = var.spec
}

data "alicloud_vpn_gateways" "vpn_gateways" {
  ids                      = [alicloud_vpn_gateway.default.id]
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
* `ssl_vpn` - (Optional, ForceNew, Available since v1.243.0) Indicates whether the SSL-VPN feature is enabled. Valid value is `enable`, `disable`.
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
  * `enable_ssl` - Whether the ssl function is enabled. It has been deprecated from provider version 1.243.0, and using `ssl_vpn` instead.
  * `ssl_vpn` - Whether the ssl function is enabled.
  * `ssl_connections` - Total count of ssl vpn connections.
  * `network_type` - The network type of the VPN gateway.
  * `auto_propagate` - Whether to automatically propagate BGP routes to the VPC. Valid values: `true`, `false`.
  * `disaster_recovery_vswitch_id` - - The ID of the backup vSwitch to which the VPN gateway is attached.
  * `disaster_recovery_internet_ip` - The backup public IP address of the VPN gateway. The second IP address assigned by the system to create an IPsec-VPN connection. This parameter is returned only when the VPN gateway supports the dual-tunnel mode.
  * `vpn_type` - - The VPN gateway type. Value:  Normal (default): Normal type. NationalStandard: National Secret type.
  * `tags` - The Tag of.
  * `ssl_vpn_internet_ip` - The IP address of the SSL-VPN connection. This parameter is returned only when the VPN gateway is a public VPN gateway and supports only the single-tunnel mode. In addition, the VPN gateway must have the SSL-VPN feature enabled.
  * `vswitch_id` - - The ID of the vSwitch to which the VPN gateway is attached.
  * `resource_group_id` - The ID of the resource group.
