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

* `ids` - (Optional, Available 1.52.0+) A list of Security Group IDs.
* `name_regex` - (Optional) A regex string to filter the resulting security groups by their names.
* `vpc_id` - (Optional) Used to retrieve security groups that belong to the specified VPC ID.
* `resource_group_id` - (Optional, Available in 1.58.0+) The Id of resource group which the security_group belongs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `tags` - (Optional) A map of tags assigned to the ECS instances. It must be in the format:
  ```
  data "alicloud_security_groups" "taggedSecurityGroups" {
    tags = {
      tagKey1 = "tagValue1",
      tagKey2 = "tagValue2"
    }
  }
  ```

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Security Group IDs.
* `names` - A list of Security Group names.
* `groups` - A list of Security Groups. Each element contains the following attributes:
  * `id` - The ID of the security group.
  * `name` - The name of the security group.
  * `description` - The description of the security group.
  * `vpc_id` - The ID of the VPC that owns the security group.
  * `resource_group_id` - The Id of resource group which the security_group belongs.
  * `security_group_type` - The type of the security group.
  * `inner_access` - Whether to allow inner network access.
  * `creation_time` - Creation time of the security group.
  * `tags` - A map of tags assigned to the ECS instance.