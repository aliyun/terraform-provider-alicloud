---
layout: "alicloud"
page_title: "Alicloud: alicloud_security_groups"
sidebar_current: "docs-alicloud-datasource-security-groups"
description: |-
    Provides a list of Security Groups available to the user.
---

# alicloud\_security\_groups

The Security Groups data source provides a list of Security Groups in an Alicloud account according to the specified filters.

## Example Usage

```
# Filter security groups and output to a file

data "alicloud_security_groups" "web" {
  name_regex  = "^web-"
  output_file = "web_access.json"
}

# in conjunction with vpc

resource "alicloud_vpc" "primary" {
  ...
}

data "alicloud_security_groups" "primary_groups" {
  vpc_id = "${alicloud_vpc.primary.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the security groups list returned by Alicloud.
* `vpc_id` - (Optional) Used to retrieve security groups belong to specified VPC ID.
* `output_file` - (Optional) The name of file that can save security groups data source after running `terraform plan`.

## Attributes Reference

A list of security groups will be exported and its every element contains the following attributes:

* `id` - The ID of the security group.
* `name` - The name of the security group.
* `description` - The description of the security group.
* `vpc_id` - The ID of the VPC.
* `inner_access` - Whether to allow inner network access.
* `creation_time` - Creation time of the security group.
