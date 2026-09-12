---
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_aigw_instances"
sidebar_current: "docs-alicloud-datasource-esa-aigw-instances"
description: |-
  This data source provides a list of ESA AI Gateway Instances.
---


# alicloud_esa_aigw_instances

This data source provides a list of ESA AI Gateway Instances in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available since v1.248.0.

## Example Usage

```terraform
data "alicloud_esa_aigw_instances" "default" {
  ids = ["<aigw_instance_id>"]
}

output "first_instance_id" {
  value = data.alicloud_esa_aigw_instances.default.instances.0.aigw_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of AI Gateway Instance IDs.
* `name_regex` - (Optional) A regex string to filter results by the AI Gateway Instance name.
* `fuzzy_search_key` - (Optional) The fuzzy search keyword, supports searching by name or ID.
* `page_number` - (Optional) The page number. Default is 1.
* `page_size` - (Optional) The number of items per page. Default is 20. Valid values: 1-500.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of AI Gateway Instance IDs.
* `names` - A list of AI Gateway Instance names.
* `instances` - A list of ESA AI Gateway Instances. Each element contains the following attributes:
  * `aigw_instance_id` - The ID of the AI Gateway instance.
  * `aigw_instance_name` - The name of the AI Gateway instance.
  * `auth_key` - The authentication key of the AI Gateway instance.
  * `comment` - The remark of the AI Gateway instance.
  * `create_time` - The creation time of the AI Gateway instance.
  * `enable_auth` - Whether AuthKey authentication is enabled.
  * `record_count` - The number of domains bound to the AI Gateway instance.
  * `status` - The status of the AI Gateway instance.
  * `update_time` - The update time of the AI Gateway instance.
