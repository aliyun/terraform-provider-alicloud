---
layout: "alicloud"
page_title: "Alicloud: alicloud_dns_groups"
sidebar_current: "docs-alicloud-datasource-dns-groups"
description: |-
    Provides a list of groups available to the dns.
---

# alicloud\_dns\_groups

The Dns Domain Groups data source provides a list of Alicloud Dns Domain Groups in an Alicloud account according to the specified filters.

## Example Usage

```
data "alicloud_dns_groups" "group" {
  name_regex = "^y[A-Za-z]+"
  output_file = "groups.txt"
}

```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the group list returned by Alicloud. 
* `output_file` - (Optional) The name of file that can save groups data source after running `terraform plan`.

## Attributes Reference

A list of groups will be exported and its every element contains the following attributes:

* `group_id` - Id of the group .
* `group_name` - Name of the group .