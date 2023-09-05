---
subcategory: "VPN Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_customer_gateway"
description: |-
  Provides a Alicloud VPN Gateway Customer Gateway resource.
---

# alicloud_vpn_customer_gateway

Provides a VPN Gateway Customer Gateway resource. By creating a customer gateway, you can register the information of the local gateway to the cloud, and then connect the customer gateway to the VPN gateway.

-> **NOTE:** Terraform will auto build vpn customer gateway instance  while it uses `alicloud_vpn_customer_gateway` to build a vpn customer gateway resource.

For information about VPN customer gateway and how to use it, see [What is VPN customer gateway](https://www.alibabacloud.com/help/en/doc-detail/120368.html).

-> **NOTE:** Available since v1.14.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_vpn_customer_gateway" "default" {
  description           = "defaultCustomerGateway"
  ip_address            = "1.1.1.1"
  asn                   = 1111
  customer_gateway_name = var.name
}
```

## Argument Reference

The following arguments are supported:
* `asn` - (Optional, ForceNew) The autonomous system number of the gateway device in the data center. The asn is a 4-byte number. You can enter the number in two segments and separate the first 16 bits from the following 16 bits with a period (.). Enter the number in each segment in the decimal format.
* `customer_gateway_name` - (Optional) The name of the customer gateway.
* `description` - (Optional) The description of the customer gateway.
* `ip_address` - (Required, ForceNew) The IP address of the customer gateway.
* `tags` - (Optional, Map) tag.

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.210.0). Field 'name' has been deprecated from provider version 1.210.0. New field 'customer_gateway_name' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the customer gateway was created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Customer Gateway.
* `delete` - (Defaults to 5 mins) Used when delete the Customer Gateway.
* `update` - (Defaults to 5 mins) Used when update the Customer Gateway.

## Import

VPN Gateway Customer Gateway can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpn_customer_gateway.example <id>
```