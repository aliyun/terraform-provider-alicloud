---
subcategory: "Direct Mail"
layout: "alicloud"
page_title: "Alicloud: alicloud_direct_mail_config_sets"
sidebar_current: "docs-alicloud-datasource-direct-mail-config-sets"
description: |-
  Provides a list of Direct Mail Config Set owned by an Alibaba Cloud account.
---

# alicloud_direct_mail_config_sets

This data source provides Direct Mail Config Set available to the user. [What is Config Set](https://next.api.alibabacloud.com/document/Dm/2015-11-23/ConfigSetCreate)

-> **NOTE:** Available since v1.293.0.

## Example Usage

```terraform
data "alicloud_direct_mail_config_sets" "default" {}

output "first_config_set_id" {
  value = data.alicloud_direct_mail_config_sets.default.sets.0.id
}

data "alicloud_direct_mail_config_sets" "by_name" {
  keyword = "terraform-example"
}
```

## Argument Reference

The following arguments are supported:
* `all` - (Optional) Whether to query all records.
* `ids` - (Optional, Computed) A list of Config Set IDs.
* `keyword` - (Optional) The keyword used to search by name.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Config Set IDs.
* `sets` - A list of Config Set Entries. Each element contains the following attributes:
  * `description` - The description of the configuration set.
  * `id` - The ID of the configuration set.
  * `ip_pool_id` - The ID of the associated IP pool.
  * `ip_pool_name` - The name of the associated IP pool.
  * `is_public_channel_backoff` - Whether to fall back to the public channel when the associated dedicated IP pool has problems.
  * `name` - The name of the configuration set.
  * `validation_option` - Address validation options:
    * `enabled` - Whether the address validation feature is enabled.
    * `forbidden_status_list` - The list of address validation main statuses to block.
    * `forbidden_sub_status_list` - The list of address validation sub-statuses to block, providing finer-grained interception control than the main statuses.
