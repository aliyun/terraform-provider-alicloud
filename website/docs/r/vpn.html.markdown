---
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn"
sidebar_current: "docs-alicloud-resource-vpn"
description: |-
  Provides a Alicloud VPN resource.
---

# alicloud\_vpn

Provides a VPN resource.

~> **NOTE:** Terraform will auto build vpn instance  while it uses `alicloud_vpn` to build a vpn resource.

## Example Usage

Basic Usage

```
resource "alicloud_vpn" "foo" {
        name = "testAccVpnConfig"
        vpc_id = "vpc-2ze9wy916mfwpwbf6hx4u"
		bandwidth = "10"
        enable_ssl = true
        instance_charge_type = "postpaid"
        auto_pay = true
		description = "test_create_description"
}
```
## Argument Reference

The following arguments are supported:
* `name` - (Optional) The name of the VPN. Defaults to null.
* `vpc_id` - (Required, Forces new resource) The VPN belongs the vpc_id, the field can't be changed.
* `instance_charge_type` - (Optional) The charge type for instance.
* `period` - (Optional) The filed is only required while the InstanceChargeType is prepaid.
* `auto_pay` - (Optional) The order will be automatically payed if ture.
* `bandwidth` - (Required) The value should be 10 or 100 if the user is postpaid, otherwise it can be 5,10,20,50,100.
                   It can't be changed by terraform.
* `enable_ipsec` - (Optional) Enable or Disable IPSec VPN. At least one type of VPN should be enabled.
* `enable_ssl`  - (Optional) Enable or Disable SSL VPN.  At least one type of VPN should be enabled.
* `ssl_connections` - (Optional) The max connections of SSL VPN.
* `description` - (Optional) The description of the VPN instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPN instance id.
* `internet_ip` - The internet ip of the VPN.
* `spec` - The specification of the VPN




