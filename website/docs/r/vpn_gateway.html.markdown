---
subcategory: "VPN"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_gateway"
sidebar_current: "docs-alicloud-resource-vpn-gateway"
description: |-
  Provides a Alicloud VPN gateway resource.
---

# alicloud\_vpn_gateway

Provides a VPN gateway resource.

-> **NOTE:** Terraform will auto build vpn instance  while it uses `alicloud_vpn_gateway` to build a vpn resource.

-> Currently International-Site account can open `PostPaid` VPN gateway and China-Site account can open `PrePaid` VPN gateway.

For information about VPN gateway and how to use it, see [What is VPN gateway](https://www.alibabacloud.com/help/en/doc-detail/120365.html).

## Example Usage

Basic Usage

```
resource "alicloud_vpc" "vpc" {
  name       = "tf_test_foo"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vsw" {
  vpc_id            = alicloud_vpc.vpc.id
  cidr_block        = "172.16.0.0/21"
  availability_zone = "cn-beijing-b"
}

resource "alicloud_vpn_gateway" "foo" {
  name                 = "vpnGatewayConfig"
  vpc_id               = alicloud_vpc.vpc.id
  bandwidth            = "10"
  enable_ssl           = true
  instance_charge_type = "PostPaid"
  description          = "test_create_description"
  vswitch_id           = alicloud_vswitch.vsw.id
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
* `instance_charge_type` - (ForceNew) The charge type for instance. If it is an international site account, the valid value is PostPaid, otherwise PrePaid. 
                                Default to PostPaid. 
* `period` - (Optional) The filed is only required while the InstanceChargeType is PrePaid. Valid values: [1-9, 12, 24, 36]. Default to 1. 
* `bandwidth` - (Required) The value should be 10, 100, 200. if the user is postpaid, otherwise it can be 5, 10, 20, 50, 100, 200.
                   It can't be changed by terraform.
* `enable_ipsec` - (Optional) Enable or Disable IPSec VPN. At least one type of VPN should be enabled.
* `enable_ssl`  - (Optional) Enable or Disable SSL VPN.  At least one type of VPN should be enabled.
* `ssl_connections` - (Optional) The max connections of SSL VPN. Default to 5. The number of connections supported by each account is different. 
                        This field is ignored when enable_ssl is false.
* `description` - (Optional) The description of the VPN instance.
* `vswitch_id` - (Optional, ForceNew, Available in v1.56.0+) The VPN belongs the vswitch_id, the field can't be changed.
* `tags` - (Optional, Available in v1.160.0+) The tags of VPN gateway.
* `auto_pay` - (Optional, Available in v1.160.0+)  Whether to pay automatically. Default value: `true`. Valid values:
  `false`: If automatic payment is not enabled, you need to go to the order center to complete the payment after the order is generated.
  `true`: Enable automatic payment, automatic payment order.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPN instance id.
* `internet_ip` - The internet ip of the VPN.
* `status` - The status of the VPN gateway.
* `business_status` - The business status of the VPN gateway.


#### Timeouts

-> **NOTE:** Available in 1.160.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the vpn gateway.
* `delete` - (Defaults to 10 mins) Used when delete the vpn gateway.

## Import

VPN gateway can be imported using the id, e.g.

```
$ terraform import alicloud_vpn_gateway.example vpn-abc123456
```


