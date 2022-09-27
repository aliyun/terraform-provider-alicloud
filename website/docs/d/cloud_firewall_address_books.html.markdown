---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_address_books"
sidebar_current: "docs-alicloud-datasource-cloud-firewall-address-books"
description: |-
  Provides a list of Cloud Firewall Address Books to the user.
---

# alicloud\_cloud\_firewall\_address\_books

This data source provides the Cloud Firewall Address Books of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.178.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_firewall_address_books" "ids" {}
output "cloud_firewall_address_book_id_1" {
  value = data.alicloud_cloud_firewall_address_books.ids.books.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter results Address Book name.
* `group_type` - (Optional, ForceNew) The type of the Address Book. Valid values: `ip`, `tag`.
* `ids` - (Optional, ForceNew, Computed)  A list of Address Book IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Address Book names.
* `books` - A list of Cloud Firewall Address Books. Each element contains the following attributes:
    * `address_list` - The addresses in the Address Book.
    * `auto_add_tag_ecs` - Whether you want to automatically add new matching tags of the ECS IP address to the Address Book.
    * `description` - The description of the Address Book.
    * `group_name` - The name of the Address Book.
    * `group_type` - The type of the Address Book.
    * `group_uuid` - The ID of the Address Book.
    * `id` - The ID of the Address Book.
    * `tag_relation` - One or more tags for the relationship between.
    * `ecs_tags` - The logical relation among the ECS tags that to be matchedh.
		* `tag_key` - The key of ECS tag that to be matched.
        * `tag_value` - The value of ECS tag that to be matched.