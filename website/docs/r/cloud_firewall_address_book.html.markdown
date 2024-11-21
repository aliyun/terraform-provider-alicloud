---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_address_book"
sidebar_current: "docs-alicloud-resource-cloud-firewall-address-book"
description: |-
  Provides a Alicloud Cloud Firewall Address Book resource.
---

# alicloud_cloud_firewall_address_book

Provides a Cloud Firewall Address Book resource.

For information about Cloud Firewall Address Book and how to use it, see [What is Address Book](https://www.alibabacloud.com/help/en/cloud-firewall/developer-reference/api-cloudfw-2017-12-07-addaddressbook).

-> **NOTE:** Available since v1.178.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_address_book&exampleId=25f27cdd-69ee-ab3a-9836-9ee11efea44fac907b7d&activeTab=example&spm=docs.r.cloud_firewall_address_book.0.25f27cdd69&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `group_name` - (Required) The name of the Address Book.
* `group_type` - (Required, ForceNew) The type of the Address Book. Valid values: `ip`, `ipv6`, `domain`, `port`, `tag`.
**NOTE:** From version 1.213.1, `group_type` can be set to `ipv6`, `domain`, `port`.
* `description` - (Required) The description of the Address Book.
* `auto_add_tag_ecs` - (Optional, Int) Whether you want to automatically add new matching tags of the ECS IP address to the Address Book. Valid values: `0`, `1`.
* `tag_relation` - (Optional) The logical relation among the ECS tags that to be matched. Default value: `and`. Valid values:
  - `and`: Only the public IP addresses of ECS instances that match all the specified tags can be added to the Address Book.
  - `or`: The public IP addresses of ECS instances that match one of the specified tags can be added to the Address Book.
* `lang` - (Optional) The language of the content within the request and response. Valid values: `zh`, `en`.
* `address_list` - (Optional, List) The list of addresses.
* `ecs_tags` - (Optional, Set) A list of ECS tags. See [`ecs_tags`](#ecs_tags) below.

### `ecs_tags`

The ecs_tags supports the following:

* `tag_key` - (Optional) The key of ECS tag that to be matched.
* `tag_value` - (Optional) The value of ECS tag that to be matched.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Address Book.

## Import

Cloud Firewall Address Book can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_address_book.example <id>
```
