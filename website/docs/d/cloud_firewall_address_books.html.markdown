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

-> **NOTE:** Available in v1.152.0+.

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

* `contain_port` - (Optional, ForceNew) The contain port.
* `group_type` - (Optional, ForceNew) Address Book, it can be set to include:
  -**ip**: an IP address book
  -**domain**: domain name address book
  -**port**: port Address Book
  -**tag**:ECS tag address book. Valid values: `domain`, `ip`, `port`, `tag`.
* `ids` - (Optional, ForceNew, Computed)  A list of Address Book IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `books` - A list of Cloud Firewall Address Books. Each element contains the following attributes:
	* `address_list` - Returns the address book table.
	* `auto_add_tag_ecs` - Whether you want to automatically add new matching tags of the ECS IP address to the address book.
       -**1**: the automatically added
       -**0**: indicates that does not automatically add.
	* `description` - After the description of.
	* `group_name` - Address book name.
	* `group_type` - Address Book, it can be set to include:
-**ip**: an IP address book
-**domain**: domain name address book
-**port**: port Address Book
-**tag**:ECS tag address book.
	* `group_uuid` - Address book unique ID.
	* `id` - The ID of the Address Book.
	* `tag_relation` - One or more tags for the relationship between.
-**and**: the one or more tags for the Inter-as "and" relationship, that is to say, in the meantime matching the plurality of tags of the ECS IP address will be added to the address book.
-**or**: a plurality of inter-tag "or" relationship, that is, as long as the matching one of the tags of the ECS IP address will be added to the address book.
	* `ecs_tags` - ECS tags.
		* `tag_key` - Of the to-be-matched ECS tags Key.
		* `tag_value` - Of the to-be-matched ECS tag value.