---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_dns_domain_groups"
sidebar_current: "docs-alicloud-datasource-dns-domain-groups"
description: |-
    Provides a list of groups available to the domain.
---

# alicloud_dns_domain_groups

-> **DEPRECATED:** This data source has been deprecated from [v1.3.2](https://github.com/alibaba/terraform-provider/releases/tag/V1.3.2). Please use the data source `alicloud_dns_groups` instead.

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by group name.
* `ids` - (Optional) A list of group IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of group names.
* `groups` - A list of groups.

---

### `groups`

The groups supports the following attributes:

* `group_id` - The ID of the group.
* `group_name` - The name of the group.
