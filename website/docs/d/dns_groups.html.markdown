---
layout: "alicloud"
page_title: "Alicloud: alicloud_dns_groups"
sidebar_current: "docs-alicloud-datasource-dns-groups"
description: |-
    Provides a list of groups available to the dns.
---

# alicloud\_dns\_groups

This data source provides a list of DNS Domain Groups in an Alibaba Cloud account according to the specified filters.

## Example Usage

```
data "alicloud_dns_groups" "groups_ds" {
  name_regex = "^y[A-Za-z]+"
  output_file = "groups.txt"
}

output "first_group_name" {
  value = "${data.alicloud_dns_groups.groups_ds.groups.0.group_name}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by group name. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `groups` - A list of groups. Each element contains the following attributes:
  * `group_id` - Id of the group.
  * `group_name` - Name of the group.