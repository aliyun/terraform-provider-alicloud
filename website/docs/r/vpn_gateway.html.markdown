---
subcategory: "VPN Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_gateway"
description: |-
  Provides a Alicloud VPN gateway resource.
---

# alicloud_vpn_gateway

Provides a VPN gateway resource.

-> **NOTE:** Terraform will auto build vpn instance  while it uses `alicloud_vpn_gateway` to build a vpn resource.

-> Currently International-Site account can open `PostPaid` VPN gateway and China-Site account can open `PrePaid` VPN gateway.

For information about VPN gateway and how to use it, see [What is VPN gateway](https://www.alibabacloud.com/help/en/doc-detail/120365.html).

-> **NOTE:** Available since v1.13.0.

## Example Usage

Basic Usage

[IPsec-VPN connections support the dual-tunnel mode](https://www.alibabacloud.com/help/en/vpn/product-overview/ipsec-vpn-connections-support-the-dual-tunnel-mode)

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpn_gateway&exampleId=e828140f-319c-9314-3b20-c45c1468b6f795539aa9&activeTab=example&spm=docs.r.vpn_gateway.0.e828140f31&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

variable "spec" {
  default = "20"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
  cidr_block = "172.16.0.0/16"
}

data "alicloud_vswitches" "default0" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids.0
}

data "alicloud_vswitches" "default1" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids.1
}

resource "alicloud_vpn_gateway" "default" {
  vpn_type         = "Normal"
  vpn_gateway_name = var.name

  vswitch_id                   = data.alicloud_vswitches.default0.ids.0
  disaster_recovery_vswitch_id = data.alicloud_vswitches.default1.ids.0
  auto_pay                     = true
  vpc_id                       = data.alicloud_vpcs.default.ids.0
  network_type                 = "public"
  payment_type                 = "Subscription"
  enable_ipsec                 = true
  bandwidth                    = var.spec
}
```

### Deleting `alicloud_vpn_gateway` or removing it from your configuration

The `alicloud_vpn_gateway` resource allows you to manage  `payment_type = "Subscription"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `auto_pay` - (Optional, Available since v1.160.0) Whether to pay automatically. Default value: `true`. Valid values:
  - `false`: If automatic payment is not enabled, you need to go to the order center to complete the payment after the order is generated.
  - `true`: Enable automatic payment, automatic payment order.
* `auto_propagate` - (Optional) Whether to automatically propagate the BGP route to the VPC. Value:  true: Propagate automatically.  false: does not propagate automatically.
* `bandwidth` - (Required, ForceNew) The Bandwidth specification of the VPN gateway. Unit: Mbps.  If you want to create a public VPN gateway, the value is 5, 10, 20, 50, 100, 200, 500, or 1000. If you want to create a private VPN gateway, the value is 200 or 1000.
* `description` - (Optional) The description of the VPN gateway.
* `disaster_recovery_vswitch_id` - (Optional, ForceNew) The ID of the backup VSwitch to which the VPN gateway is attached.
* `network_type` - (Optional, ForceNew) The network type of the VPN gateway. Value:  public (default): public VPN gateway. private: private network VPN gateway.
* `payment_type` - (Optional, ForceNew, Computed) Type of payment. Value: Subscription: prepaid PayAsYouGo: Post-paid.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `ssl_connections` - (Optional, ForceNew) Maximum number of clients.
* `tags` - (Optional, Map) The Tag of.
* `vswitch_id` - (Optional, ForceNew) The ID of the VSwitch to which the VPN gateway is attached.
* `vpc_id` - (Required, ForceNew) The ID of the VPC to which the VPN gateway belongs.
* `vpn_gateway_name` - (Optional) The name of the VPN gateway.
* `vpn_type` - (Optional, ForceNew) The VPN gateway type. Value:  Normal (default): Normal type. NationalStandard: National Secret type.
* `period` - (Optional) The filed is only required while the InstanceChargeType is PrePaid. Valid values: [1-9, 12, 24, 36]. Default to 1.
* `enable_ipsec` - (Optional) Enable or Disable IPSec VPN. At least one type of VPN should be enabled.
* `enable_ssl` - (Optional) Enable or Disable SSL VPN.  At least one type of VPN should be enabled.

The following arguments will be discarded. Please use new fields as soon as possible:
* `instance_charge_type` - (Deprecated since v1.216.0). Field 'instance_charge_type' has been deprecated from provider version 1.216.0. New field 'payment_type' instead.
* `name` - (Deprecated since v1.216.0). Field 'name' has been deprecated from provider version 1.216.0. New field 'vpn_gateway_name' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the VPN gateway was created.
* `status` - The status of the resource.
* `internet_ip` - The internet ip of the VPN.
* `business_status` - The business status of the VPN gateway.
* `ssl_vpn_internet_ip` - The IP address of the SSL-VPN connection. This parameter is returned only when the VPN gateway is a public VPN gateway and supports only the single-tunnel mode. In addition, the VPN gateway must have the SSL-VPN feature enabled.
* `disaster_recovery_internet_ip` - The backup public IP address of the VPN gateway. The second IP address assigned by the system to create an IPsec-VPN connection. This parameter is returned only when the VPN gateway supports the dual-tunnel mode.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the vpn gateway.
* `delete` - (Defaults to 10 mins) Used when delete the vpn gateway.

## Import

VPN gateway can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpn_gateway.example <id>
```