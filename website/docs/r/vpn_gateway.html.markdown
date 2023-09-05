---
subcategory: "VPN Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_gateway"
description: |-
  Provides a Alicloud VPN Gateway VPN Gateway resource.
---

# alicloud_vpn_gateway

Provides a VPN Gateway VPN Gateway resource. 

For information about VPN Gateway VPN Gateway and how to use it, see [What is VPN Gateway](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.13.0.
-> Currently International-Site account can open `PostPaid` VPN gateway and China-Site account can open `PrePaid` VPN gateway.
## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultVpc" {
  description = "test_wbw_vpc"
  vpc_name    = var.name
  cidr_block  = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultVswitch_1" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  description  = "defaultVswitch_1"
  cidr_block   = "192.168.1.0/24"
  vswitch_name = "${var.name}1"
  zone_id      = data.alicloud_zones.default.zones.0.id
}


resource "alicloud_vpn_gateway" "default" {
  ssl_connections  = 5
  enable_ipsec     = true
  enable_ssl       = true
  auto_pay         = true
  vpn_type         = "Normal"
  description      = "test_vpn"
  vpn_gateway_name = var.name
  network_type     = "public"
  bandwidth        = "5"
  vswitch_id       = alicloud_vswitch.defaultVswitch_1.id
  vpc_id           = alicloud_vpc.defaultVpc.id
  period           = 1
  payment_type     = "Subscription"
}
```

### Deleting `alicloud_vpn_gateway` or removing it from your configuration

The `alicloud_vpn_gateway` resource allows you to manage  `payment_type = "Subscription"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `auto_pay` - (Optional) Whether to automatically pay the bill of the VPN gateway. Value:true: Automatically pay the bill of the VPN gateway.false (default): does not automatically pay the bill of the VPN gateway.
* `auto_propagate` - (Optional) Whether to automatically propagate the BGP route to the VPC. Value:true: Propagate automatically.false: does not propagate automatically.
* `bandwidth` - (Required, ForceNew) The Bandwidth specification of the VPN gateway. Unit: Mbps.If you want to create a public VPN gateway, the value is 5, 10, 20, 50, 100, 200, 500, or 1000.If you want to create a private VPN gateway, the value is 200 or 1000.
* `description` - (Optional) The description of the VPN gateway.
* `enable_ipsec` - (Optional) Whether to enable the IPsec-VPN function. Value:true (default): Enables the IPsec-VPN function.false: IPsec-VPN function is not enabled.
* `enable_ssl` - (Optional) Whether to enable the SSL-VPN function. Value:true: Enable the SSL-VPN function.false (default): The SSL-VPN function is not enabled.
* `network_type` - (Optional, ForceNew) The network type of the VPN gateway. Value:public (default): public VPN gateway.private: private network VPN gateway.
* `payment_type` - (Optional, ForceNew, Computed) Type of payment. Value:Subscription: prepaidPayAsYouGo: Post-paid.
* `period` - (Optional) Duration of purchase. Unit: Month. Value: 1~9, 12, 24, or 36.
* `ssl_connections` - (Optional) Maximum number of clients.
* `tags` - (Optional, Map) The Tag of.
* `vswitch_id` - (Optional, ForceNew) The ID of the VSwitch to which the VPN gateway is attached.
* `vpc_id` - (Required, ForceNew) The ID of the VPC to which the VPN gateway belongs.
* `vpn_gateway_name` - (Optional) The name of the VPN gateway.
* `vpn_type` - (Optional, ForceNew) The VPN gateway type. Value:Normal (default): Normal type.NationalStandard: National Secret type.

The following arguments will be discarded. Please use new fields as soon as possible:
* `instance_charge_type` - (Deprecated since v1.210.0). Field 'instance_charge_type' has been deprecated from provider version 1.210.0. New field 'payment_type' instead.
* `name` - (Deprecated since v1.210.0). Field 'name' has been deprecated from provider version 1.210.0. New field 'vpn_gateway_name' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the VPN gateway was created.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the VPN Gateway.
* `delete` - (Defaults to 5 mins) Used when delete the VPN Gateway.
* `update` - (Defaults to 5 mins) Used when update the VPN Gateway.

## Import

VPN Gateway VPN Gateway can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpn_gateway.example <id>
```