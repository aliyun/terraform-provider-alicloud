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

## Example Usage

Basic Usage

```
resource "alicloud_vpc" "vpc" {
  name       = "tf_test_foo"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vsw" {
  vpc_id            = "${alicloud_vpc.vpc.id}"
  cidr_block        = "172.16.0.0/21"
  availability_zone = "cn-beijing-b"
}

resource "alicloud_vpn_gateway" "foo" {
  name                 = "vpnGatewayConfig"
  vpc_id               = "${alicloud_vpc.vpc.id}"
  bandwidth            = "10"
  enable_ssl           = true
  instance_charge_type = "PostPaid"
  description          = "test_create_description"
  vswitch_id           = "${alicloud_vswitch.vsw.id}"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the VPN. Defaults to null.
* `vpc_id` - (Required, ForceNew) The VPN belongs the vpc_id, the field can't be changed.
* `instance_charge_type` - (ForceNew) The charge type for instance. Valid value: PostPaid, PrePaid. Default to PostPaid.
* `period` - (Optional) The filed is only required while the InstanceChargeType is prepaid.
* `bandwidth` - (Required) The value should be 10, 100, 200, 500, 1000 if the user is postpaid, otherwise it can be 5, 10, 20, 50, 100, 200, 500, 1000.
                   It can't be changed by terraform.
* `enable_ipsec` - (Optional) Enable or Disable IPSec VPN. At least one type of VPN should be enabled.
* `enable_ssl`  - (Optional) Enable or Disable SSL VPN.  At least one type of VPN should be enabled.
* `ssl_connections` - (Optional) The max connections of SSL VPN. Default to 5. This field is ignored when enable_ssl is false.
* `description` - (Optional) The description of the VPN instance.
* `vswitch_id` - (Optional, ForceNew, Available in v1.56.0+) The VPN belongs the vswitch_id, the field can't be changed.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPN instance id.
* `internet_ip` - The internet ip of the VPN.
* `status` - The status of the VPN gateway.
* `business_status` - The business status of the VPN gateway.

## Import

VPN gateway can be imported using the id, e.g.

```
$ terraform import alicloud_vpn_gateway.example vpn-abc123456
```


