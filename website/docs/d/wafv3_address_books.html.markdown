---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_wafv3_address_books"
sidebar_current: "docs-alicloud-datasource-wafv3-address-books"
description: |-
  Provides a list of Wafv3 Address Book owned by an Alibaba Cloud account.
---

# alicloud_wafv3_address_books

This data source provides Wafv3 Address Book available to the user.[What is Address Book](https://next.api.alibabacloud.com/document/waf-openapi/2021-10-01/CreateDefenseRule)

-> **NOTE:** Available since v1.283.0.

## Example Usage

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
  description       = "test"
  instance_id       = data.alicloud_wafv3_instances.default.ids.0
  address_book_name = var.name
  address_list      = ["100.100.100.100/32", "101.101.101.101/32", "102.102.102.102/32"]
  address_book_type = "ip"
}

data "alicloud_wafv3_address_books" "default" {
  ids         = ["${alicloud_wafv3_address_book.default.id}"]
  name_regex  = alicloud_wafv3_address_book.default.address_book_name
  instance_id = data.alicloud_wafv3_instances.default.ids.0
}

output "alicloud_wafv3_address_book_example_id" {
  value = data.alicloud_wafv3_address_books.default.books.0.id
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (Required) The ID of the WAF instance.

-> **NOTE:**  You can call [DescribeInstance](~~ 433756 ~~) to view the ID of the current WAF instance.

* `ids` - (Optional, Computed) A list of Address Book IDs. The value is formulated as `<instance_id>:<address_book_id>`.
* `name_regex` - (Optional) A regex string to filter results by Address Book name.
* `enable_details` - (Optional) Default to `false`. Set it to `true` to fetch the `address_list` of each Address Book.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Address Book IDs.
* `names` - A list of name of Address Books.
* `books` - A list of Address Book Entries. Each element contains the following attributes:
    * `address_book_id` - The ID of the Address Book.
    * `address_book_name` - The name of the Address Book.
    * `address_book_type` - The type of the Address Book. Valid values: `ip`.
    * `address_list` - The address list of the Address Book. **NOTE:** This field is only available when `enable_details` is `true`.
    * `description` - The description of the Address Book.
    * `id` - The resource ID. It is formatted as `<instance_id>:<address_book_id>`.
