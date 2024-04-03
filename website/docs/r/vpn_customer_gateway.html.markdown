---
subcategory: "VPN Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_customer_gateway"
description: |-
  Provides a Alicloud VPN customer gateway resource.
---

# alicloud_vpn_customer_gateway

Provides a VPN customer gateway resource.

-> **NOTE:** Terraform will auto build vpn customer gateway instance  while it uses `alicloud_vpn_customer_gateway` to build a vpn customer gateway resource.

For information about VPN customer gateway and how to use it, see [What is VPN customer gateway](https://www.alibabacloud.com/help/en/doc-detail/120368.html).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "ap-southeast-1"
}

resource "alicloud_vpn_customer_gateway" "default" {
  description           = var.name
  ip_address            = "4.3.2.10"
  asn                   = "1219002"
  customer_gateway_name = var.name
}
```

## Argument Reference

The following arguments are supported:
* `asn` - (Optional, ForceNew) Asn.
* `customer_gateway_name` - (Optional) The name of the customer gateway.
* `description` - (Optional) The description of the customer gateway.
* `ip_address` - (Required, ForceNew) The IP address of the customer gateway.
* `tags` - (Optional, Map) tag.

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.216.0). Field 'name' has been deprecated from provider version 1.216.0. New field 'customer_gateway_name' instead.

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

VPN customer gateway can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpn_customer_gateway.example <id>
```