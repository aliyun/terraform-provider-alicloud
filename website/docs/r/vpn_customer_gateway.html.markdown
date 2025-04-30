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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpn_customer_gateway&exampleId=1865a011-1907-538a-6dc8-13c4f15e93ea5c03bddc&activeTab=example&spm=docs.r.vpn_customer_gateway.0.1865a01119&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Customer Gateway.
* `delete` - (Defaults to 5 mins) Used when delete the Customer Gateway.
* `update` - (Defaults to 5 mins) Used when update the Customer Gateway.

## Import

VPN customer gateway can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpn_customer_gateway.example <id>
```