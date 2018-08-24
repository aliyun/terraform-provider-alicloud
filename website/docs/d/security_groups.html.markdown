---
layout: "alicloud"
page_title: "Alicloud: alicloud_security_groups"
sidebar_current: "docs-alicloud-datasource-security-groups"
description: |-
    Provides a list of Security Groups available to the user.
---

# alicloud\_security\_groups

This data source provides a list of Security Groups in an Alibaba Cloud account according to the specified filters.

## Example Usage

```
# Filter security groups and print the results into a file
data "alicloud_security_groups" "sec_groups_ds" {
  name_regex  = "^web-"
  output_file = "web_access.json"
}

# In conjunction with a VPC
resource "alicloud_vpc" "primary_vpc_ds" {
  # ...
}

data "alicloud_security_groups" "primary_sec_groups_ds" {
  vpc_id = "${alicloud_vpc.primary_vpc_ds.id}"
}

output "first_group_id" {
  value = "${data.alicloud_security_groups.primary_sec_groups_ds.groups.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter the resulting security groups by their names.
* `vpc_id` - (Optional) Used to retrieve security groups that belong to the specified VPC ID.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `groups` - A list of groups. Each element contains the following attributes:
  * `id` - The ID of the security group.
  * `name` - The name of the security group.
  * `description` - The description of the security group.
  * `vpc_id` - The ID of the VPC that owns the security group.
  * `inner_access` - Whether to allow inner network access.
  * `creation_time` - Creation time of the security group.
