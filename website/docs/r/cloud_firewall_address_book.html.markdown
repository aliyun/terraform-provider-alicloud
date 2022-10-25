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

For information about Cloud Firewall Address Book and how to use it, see [What is Address Book](https://www.alibabacloud.com/help/en/cloud-firewall/latest/describeaddressbook).

-> **NOTE:** Available in v1.178.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cloud_firewall_address_book" "example" {
  description      = "example_value"
  group_name       = "example_value"
  group_type       = "tag"
  tag_relation     = "and"
  auto_add_tag_ecs = 0
  ecs_tags {
    tag_key   = "created"
    tag_value = "tfTestAcc0"
  }
}
```

## Argument Reference

The following arguments are supported:

* `address_list` - (Optional) The list of addresses.
* `description` - (Required) The description of the Address Book.
* `group_name` - (Required) The name of the Address Book.
* `group_type` - (Required, ForceNew) The type of the Address Book. Valid values:  `ip`, `tag`.
* `ecs_tags` - (Optional) A list of ECS tags. See the following `Block ecs_tags`.
* `tag_relation` - (Optional) The logical relation among the ECS tags that to be matched. Valid values:
  - **and**: Only the public IP addresses of ECS instances that match all the specified tags can be added to the Address Book. This is the default value.
  - **or**: The public IP addresses of ECS instances that match one of the specified tags can be added to the Address Book.
* `auto_add_tag_ecs` - (Optional) Whether you want to automatically add new matching tags of the ECS IP address to the Address Book. Valid values: `0`, `1`.
* `lang` - (Optional) The language of the content within the request and response. Valid values: `en`, `zh`.

### Block ecs_tags

The ecs_tags supports the following:

* `tag_key` - (Optional) The key of ECS tag that to be matched.
* `tag_value` - (Optional) The value of ECS tag that to be matched.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Address Book.

## Import

Cloud Firewall Address Book can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_firewall_address_book.example <id>
```