---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_domain_groups"
sidebar_current: "docs-alicloud-datasource-alidns-domain-groups"
description: |-
    Provides a list of Domain Groups available to the Alidns.
---

# alicloud\_alidns\_domain\_groups

This data source provides a list of Alidns Domain Groups in an Alibaba Cloud account according to the specified filters.

-> **NOTE:**  Available in 1.85.0+.

## Example Usage

```terraform
data "alicloud_alidns_domain_groups" "example" {
  ids = ["c5ef2bc43064445787adf182af2****"]
}
output "first_domain_group_id" {
  value = "${data.alicloud_alidns_domain_groups.example.groups.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance IDs.
* `name_regex` - (Optional) A regex string to filter results by the domain group name. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of instance IDs. 
* `names` - A list of domain group names.
* `groups` - A list of instances. Each element contains the following attributes:
  * `id` - Id of the instance.
  * `domain_count` - Number of domain names in the group.
  * `group_id` - Id of the domain group.
  * `group_name` - The name of the domain group.
