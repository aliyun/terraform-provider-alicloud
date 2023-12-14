---
subcategory: "VPN Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_gateway"
sidebar_current: "docs-alicloud-resource-vpn-gateway"
description: |-
  Provides a Alicloud VPN gateway resource.
---

# alicloud_vpn_gateway

Provides a VPN gateway resource.

-> **NOTE:** Terraform will auto build vpn instance  while it uses `alicloud_vpn_gateway` to build a vpn resource.

-> Currently International-Site account can open `PostPaid` VPN gateway and China-Site account can open `PrePaid` VPN gateway.

For information about VPN gateway and how to use it, see [What is VPN gateway](https://www.alibabacloud.com/help/en/doc-detail/120365.html).

-> **NOTE:** Available since v1.13.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids.0
}

resource "alicloud_vpn_gateway" "default" {
  name                 = var.name
  vpc_id               = data.alicloud_vpcs.default.ids.0
  bandwidth            = "10"
  enable_ssl           = true
  description          = var.name
  instance_charge_type = "PrePaid"
  vswitch_id           = data.alicloud_vswitches.default.ids.0
}
```

### Deleting `alicloud_vpn_gateway` or removing it from your configuration

The `alicloud_vpn_gateway` resource allows you to manage `instance_charge_type = "Prepaid"` vpn gateway, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration
will remove it from your state file and management, but will not destroy the VPN Gateway.
You can resume managing the subscription vpn gateway via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the VPN. Defaults to null.
* `vpc_id` - (Required, ForceNew) The VPN belongs the vpc_id, the field can't be changed.
* `instance_charge_type` - (Optional, ForceNew) The charge type for instance. If it is an international site account, the valid value is PostPaid, otherwise PrePaid. 
                                Default to PostPaid. 
* `period` - (Optional) The filed is only required while the InstanceChargeType is PrePaid. Valid values: [1-9, 12, 24, 36]. Default to 1. 
* `bandwidth` - (Required, ForceNew) The value should be 10, 100, 200. if the user is postpaid, otherwise it can be 5, 10, 20, 50, 100, 200.
                   It can't be changed by terraform.
* `enable_ipsec` - (Optional) Enable or Disable IPSec VPN. At least one type of VPN should be enabled.
* `enable_ssl`  - (Optional, ForceNew) Enable or Disable SSL VPN.  At least one type of VPN should be enabled.
* `ssl_connections` - (Optional) The max connections of SSL VPN. Default to 5. The number of connections supported by each account is different. 
                        This field is ignored when enable_ssl is false.
* `description` - (Optional) The description of the VPN instance.
* `vswitch_id` - (Optional, ForceNew, Available in v1.56.0+) The VPN belongs the vswitch_id, the field can't be changed.
* `tags` - (Optional, Available in v1.160.0+) The tags of VPN gateway.
* `auto_pay` - (Optional, Available in v1.160.0+)  Whether to pay automatically. Default value: `true`. Valid values:
    - `false`: If automatic payment is not enabled, you need to go to the order center to complete the payment after the order is generated.
    - `true`: Enable automatic payment, automatic payment order.
* `auto_propagate` - (Optional, Available in v1.184.0+) Specifies whether to automatically advertise BGP routes to the virtual private cloud (VPC). Valid values:
    - `true`: Enable.
    - `false`: Disable.
* `network_type` - (Optional, ForceNew, Available in v1.193.0+) The network type of the VPN gateway. Value:
    - public (default): Public VPN gateway. 
    - private: Private VPN gateway.

  -> **NOTE:** Private VPN gateway can only be purchased by white list users, and the bandwidth only supports 200M or 1000M; In addition, SSL is not supported.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPN instance id.
* `internet_ip` - The internet ip of the VPN.
* `status` - The status of the VPN gateway.
* `business_status` - The business status of the VPN gateway.


## Timeouts

-> **NOTE:** Available in 1.160.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the vpn gateway.
* `delete` - (Defaults to 10 mins) Used when delete the vpn gateway.

## Import

VPN gateway can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpn_gateway.example vpn-abc123456
```


