---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_prefix_lists"
sidebar_current: "docs-alicloud-datasource-vpc-prefix-lists"
description: |-
  Provides a list of Vpc Prefix Lists to the user.
---

# alicloud_vpc_prefix_lists

This data source provides the Vpc Prefix Lists of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.182.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_prefix_lists" "ids" {}
output "vpc_prefix_list_id_1" {
  value = data.alicloud_vpc_prefix_lists.ids.lists.0.id
}

data "alicloud_vpc_prefix_lists" "nameRegex" {
  name_regex = "^my-PrefixList"
}
output "vpc_prefix_list_id_2" {
  value = data.alicloud_vpc_prefix_lists.nameRegex.lists.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `true`. Set it to `false` can hide the `entrys` to output.
* `ids` - (Optional, ForceNew, Computed) A list of Prefix List IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Prefix List name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `prefix_list_name` - (Optional, ForceNew) The name of the prefix list.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Prefix List names.
* `lists` - A list of Vpc Prefix Lists. Each element contains the following attributes:
  * `create_time` - The time when the prefix list was created.
  * `entrys` - The CIDR address block list of the prefix list.
    * `cidr` - The CIDR address block of the prefix list.
    * `description` - The description of the cidr entry.
  * `id` - The ID of the Prefix List.
  * `ip_version` - The IP version of the prefix list.
  * `max_entries` - The maximum number of entries for CIDR address blocks in the prefix list.
  * `prefix_list_id` - The ID of the query Prefix List.
  * `prefix_list_name` - The name of the prefix list.
  * `share_type` - The share type of the prefix list.
  * `prefix_list_description` - The description of the prefix list.