---
subcategory: "VPN Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_gateway_enhanced_vpn_gateways"
sidebar_current: "docs-alicloud-datasource-vpn-gateway-enhanced-vpn-gateways"
description: |-
  Provides a list of Vpn Gateway Enhanced Vpn Gateway owned by an Alibaba Cloud account.
---

# alicloud_vpn_gateway_enhanced_vpn_gateways

This data source provides Vpn Gateway Enhanced Vpn Gateway available to the user.[What is Enhanced Vpn Gateway](https://next.api.alibabacloud.com/document/Vpc/2016-04-28/CreateEnhancedVpnGateway)

-> **NOTE:** Available since v1.280.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "ap-southeast-3"
}

variable "region" {
  default = "ap-southeast-3"
}

variable "zone2" {
  default = "ap-southeast-3a"
}

variable "zone1" {
  default = "ap-southeast-3b"
}

resource "alicloud_vpc" "defaulttYTx5F" {
  cidr_block = "192.168.0.0/16"
  is_default = false
}

resource "alicloud_vswitch" "defaultTRk7k3" {
  vpc_id     = alicloud_vpc.defaulttYTx5F.id
  zone_id    = var.zone1
  cidr_block = "192.168.10.0/24"
}

resource "alicloud_vswitch" "default23kGFr" {
  vpc_id     = alicloud_vpc.defaulttYTx5F.id
  zone_id    = var.zone2
  cidr_block = "192.168.20.0/24"
}


resource "alicloud_vpn_gateway_enhanced_vpn_gateway" "default" {
  vpn_type                     = "Normal"
  description                  = "default"
  disaster_recovery_vswitch_id = alicloud_vswitch.default23kGFr.id
  vpc_id                       = alicloud_vpc.defaulttYTx5F.id
  vpn_gateway_name             = "default"
  network_type                 = "public"
  vswitch_id                   = alicloud_vswitch.defaultTRk7k3.id
  gateway_type                 = "Enhanced.SiteToSite"
  auto_propagate               = false
}

data "alicloud_vpn_gateway_enhanced_vpn_gateways" "default" {
  ids    = ["${alicloud_vpn_gateway_enhanced_vpn_gateway.default.id}"]
  vpc_id = alicloud_vpc.defaulttYTx5F.id
}

output "alicloud_vpn_gateway_enhanced_vpn_gateway_example_id" {
  value = data.alicloud_vpn_gateway_enhanced_vpn_gateways.default.gateways.0.id
}
```

## Argument Reference

The following arguments are supported:
* `status` - (ForceNew, Optional) The status of the resource
* `vpc_id` - (ForceNew, Optional) The ID of the VPC to which the VPN gateway belongs.
* `vpn_instance_id` - (ForceNew, Optional) The ID of the VPN gateway.
* `ids` - (Optional, ForceNew, Computed) A list of Enhanced Vpn Gateway IDs.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).
* `enable_details` - (Optional) Default to `false`. Set it to `true` to query detailed attributes.


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Enhanced Vpn Gateway IDs.
* `gateways` - A list of Enhanced Vpn Gateway Entries. Each element contains the following attributes:
    * `auto_propagate` - Specifies whether to automatically propagate BGP routes to the VPC.
    * `create_time` - The time when the VPN gateway was created.
    * `description` - The description of the VPN gateway.
    * `disaster_recovery_vswitch_id` - The ID of the backup VSwitch to which the VPN gateway is attached.
    * `gateway_type` - VPN gateway type.
    * `network_type` - Type of Gateway.
    * `status` - The status of the resource.
    * `tags` - The Tag of.
    * `vswitch_id` - The ID of the VSwitch to which the VPN gateway is attached.
    * `vpc_id` - The ID of the VPC to which the VPN gateway belongs.
    * `vpn_gateway_name` - The name of the VPN gateway.
    * `vpn_instance_id` - The ID of the VPN gateway.
    * `vpn_type` - The Type of Vpn.
    * `id` - The ID of the resource supplied above.
