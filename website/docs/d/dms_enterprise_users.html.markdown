---
subcategory: "DMS Enterprise"
layout: "alicloud"
page_title: "Alicloud: alicloud_dms_enterprise_users"
sidebar_current: "docs-alicloud-datasource-dms-enterprise-users"
description: |-
    Provides a list of available DMS Enterprise Users.
---

# alicloud\_dms\_enterprise\_users

This data source provides a list of DMS Enterprise Users in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.90.0+

## Example Usage

```terraform
# Declare the data source
data "alicloud_dms_enterprise_users" "dms_enterprise_users_ds" {
  ids    = ["uid"]
  role   = "USER"
  status = "NORMAL"
}

output "first_user_id" {
  value = "${data.alicloud_dms_enterprise_users.dms_enterprise_users_ds.users.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `role` - (Optional) The role of the user to query.
* `search_key` - (Optional) The keyword used to query users.
* `status` - (Optional) The status of the user.
* `tid` - (Optional) The ID of the tenant in DMS Enterprise.
* `ids` - (Optional)  A list of DMS Enterprise User IDs (UID).
* `name_regex` - (Optional, Available in 1.100.0+) A regex string to filter the results by the DMS Enterprise User nick_name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of DMS Enterprise User IDs (UID).
* `names` - A list of DMS Enterprise User names.
* `users` - A list of DMS Enterprise Users. Each element contains the following attributes:
  * `mobile` - The DingTalk number or mobile number of the user.
  * `nick_name` - The nickname of the user.
  * `user_name` - The nickname of the user.
  * `parent_uid` - The Alibaba Cloud unique ID (UID) of the parent account if the user corresponds to a Resource Access Management (RAM) user.
  * `role_ids` - The list ids of the role that the user plays.
  * `role_names` - The list names of the role that he user plays.
  * `status` - The status of the user.
  * `id` - The Alibaba Cloud unique ID (UID) of the user.
  * `user_id` - The ID of the user.
