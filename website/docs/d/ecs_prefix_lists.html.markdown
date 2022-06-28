---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_prefix_lists"
sidebar_current: "docs-alicloud-datasource-ecs-prefix-lists"
description: |-
  Provides a list of Ecs Prefix Lists to the user.
---

# alicloud\_ecs\_prefix\_lists

This data source provides the Ecs Prefix Lists of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.152.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecs_prefix_lists" "example" {
  ids        = ["E2RY53-xxxx"]
  name_regex = "tf-testAcc"
}

output "output_id" {
  value = data.alicloud_ecs_prefix_lists.example.lists.0.id
}
```

## Argument Reference

The following arguments are supported:

* `address_family` - (Optional, ForceNew) The address family of the prefix list. Valid values: `IPv4`,`IPv6`. This parameter is empty by default, which indicates that all prefix lists are to be queried.
* `ids` - (Optional, ForceNew, Computed)  A list of Prefix List IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by `prefix_list_name`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Prefix List names.
* `lists` - A list of Prefix Lists. Each element contains the following attributes:
    * `address_family` - The address family of the prefix list. Valid values:`IPv4`,`IPv6`.
    * `association_count` - The amount of associated resources.
    * `create_time` - The time when the prefix list was created.
    * `description` - The description of the prefix list.
    * `max_entries` - The maximum number of entries that the prefix list supports.
    * `prefix_list_id` - The ID of the prefix list.
    * `name` - The name of the prefix list.
    * `id` - The ID of the prefix list.
    * `prefix_list_name` - The name of the prefix list.

