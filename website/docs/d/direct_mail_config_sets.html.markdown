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

-> **NOTE:** Available since v2.2.0.

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
* `all` - (Optional, Bool, Available since v2.2.0) Whether to query all records.
* `ids` - (Optional, Computed, List, Available since v2.2.0) A list of Config Set IDs.
* `keyword` - (Optional, String, Available since v2.2.0) The keyword used to search by name.
* `enable_details` - (Optional, Bool, Available since v2.2.0) Defaults to `false`. Set it to `true` to call `ConfigSetDetail` for each configuration set and output the detail attributes.
* `output_file` - (Optional, String, Available since v2.2.0) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - (Computed, List, Available since v2.2.0) A list of Config Set IDs.
* `sets` - (Computed, List, Available since v2.2.0) A list of Config Set Entries. Each element contains the following attributes:
  * `description` - (Computed, String, Available since v2.2.0) The description of the configuration set.
  * `id` - (Computed, String, Available since v2.2.0) The ID of the configuration set.
  * `ip_pool_id` - (Computed, String, Available since v2.2.0) The ID of the associated IP pool.
  * `ip_pool_name` - (Computed, String, Available since v2.2.0) The name of the associated IP pool.
  * `is_public_channel_backoff` - (Computed, Bool, Available since v2.2.0) Whether to fall back to the public channel when the associated dedicated IP pool has problems.
  * `name` - (Computed, String, Available since v2.2.0) The name of the configuration set.
  * `validation_option` - (Computed, List, Available since v2.2.0) Address validation options:
    * `enabled` - (Computed, Bool, Available since v2.2.0) Whether the address validation feature is enabled.
    * `forbidden_status_list` - (Computed, List, Available since v2.2.0) The list of address validation main statuses to block.
    * `forbidden_sub_status_list` - (Computed, List, Available since v2.2.0) The list of address validation sub-statuses to block, providing finer-grained interception control than the main statuses.
