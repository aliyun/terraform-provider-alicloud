---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_wafv3_address_book"
description: |-
  Provides a Alicloud WAFV3 Address Book resource.
---

# alicloud_wafv3_address_book

Provides a WAFV3 Address Book resource.

An Address Book is a named collection of IP/CIDR entries that can be referenced from WAFV3 protection rules.

For information about WAFV3 Address Book and how to use it, see [What is Address Book](https://next.api.alibabacloud.com/document/waf-openapi/2021-10-01/CreateDefenseRule).

-> **NOTE:** Available since v1.283.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_wafv3_instances" "default" {
}

resource "alicloud_wafv3_address_book" "default" {
  description       = "example"
  instance_id       = data.alicloud_wafv3_instances.default.ids.0
  address_book_name = var.name
  address_list      = ["100.100.100.100/32", "101.101.101.101/32", "102.102.102.102/32"]
  address_book_type = "ip"
}
```

## Argument Reference

The following arguments are supported:
* `address_book_name` - (Optional) The name of the Address Book.
* `address_book_type` - (Required, ForceNew) The type of the Address Book. Valid values: `ip`.
* `address_list` - (Optional, Set) The address list of the Address Book. Each entry is a single IP address or a CIDR block, IPv4 or IPv6.
* `description` - (Optional) The description of the Address Book.
* `instance_id` - (Required, ForceNew) The ID of the WAF instance.

-> **NOTE:**  You can call [DescribeInstance](https://next.api.alibabacloud.com/document/waf-openapi/2021-10-01/DescribeInstance) to view the ID of the current WAF instance.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<instance_id>:<address_book_id>`.
* `address_book_id` - The ID of the Address Book.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Address Book.
* `delete` - (Defaults to 5 mins) Used when delete the Address Book.
* `update` - (Defaults to 5 mins) Used when update the Address Book.

## Import

WAFV3 Address Book can be imported using the id, e.g.

```shell
$ terraform import alicloud_wafv3_address_book.example <instance_id>:<address_book_id>
```
