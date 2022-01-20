---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_address_book"
sidebar_current: "docs-alicloud-resource-cloud-firewall-address-book"
description: |-
  Provides a Alicloud Cloud Firewall Address Book resource.
---

# alicloud\_cloud\_firewall\_address\_book

Provides a Cloud Firewall Address Book resource.

For information about Cloud Firewall Address Book and how to use it, see [What is Address Book](https://help.aliyun.com/).

-> **NOTE:** Available in v1.152.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cloud_firewall_address_book" "example" {
  description = "example_value"
  group_name  = "example_value"
  group_type =       "ip"
  address_list =     ["10.21.0.0/16", "10.168.0.0/16"]
  
}

```

```terraform
resource "alicloud_cloud_firewall_address_book" "example" {
  description = "example_value"
  group_name  = "example_value"
  group_type = "tag"
  tag_relation = "and"
  auto_add_tag_ecs = 0
  ecs_tags = [{
    "tag_key":   "created",
    "tag_value": "tfTestAcc0",
  },{
    "tag_key":   "for",
    "tag_value": "Tftestacc1",
  }]
  
}

```

## Argument Reference

The following arguments are supported:

* `address_list` - (Optional) Returns the address book table.

* `description` - (Required) After the description of.
* `group_name` - (Required) Address book name.
* `group_type` - (Required, ForceNew) Address Book, it can be set to include:
  - **ip**: an IP address book
  - **domain**: domain name address book
  - **port**: port Address Book
  - **tag**:ECS tag address book. Valid values: `domain`, `ip`, `port`, `tag`.
* `ecs_tags` - (Optional) ECS tags.
  * `tag_key` - (Optional) Of the to-be-matched ECS tags Key.
  * `tag_value` - (Optional) Of the to-be-matched ECS tags value.
* `tag_relation` - (Optional) One or more tags for the relationship between.
  - **and**: the one or more tags for the Inter-as `and` relationship, that is to say, in the meantime matching the plurality of tags of the ECS IP address will be added to the address book.
  - **or**: a plurality of inter-tag `or` relationship, that is, as long as the matching one of the tags of the ECS IP address will be added to the address book.
* `auto_add_tag_ecs` - (Optional) Whether you want to automatically add new matching tags of the ECS IP address to the address book.
  -**1**: the automatically added
  -**0**: indicates that does not automatically add.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Address Book.

## Import

Cloud Firewall Address Book can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_firewall_address_book.example <id>
```