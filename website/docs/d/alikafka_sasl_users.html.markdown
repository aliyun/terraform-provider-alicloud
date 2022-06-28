---
subcategory: "Alikafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_sasl_users"
sidebar_current: "docs-alicloud-datasource-alikafka-sasl-users"
description: |-
    Provides a list of alikafka sasl users available to the user.
---

# alicloud\_alikafka\_sasl\_users

This data source provides a list of ALIKAFKA Sasl users in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.66.0+

## Example Usage

```
data "alicloud_alikafka_sasl_users" "sasl_users_ds" {
  instance_id = "xxx"
  name_regex = "username"
  output_file = "saslUsers.txt"
}

output "first_sasl_username" {
  value = "${data.alicloud_alikafka_sasl_users.sasl_users_ds.users.0.username}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the ALIKAFKA Instance that owns the sasl users.
* `name_regex` - (Optional) A regex string to filter results by the username. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of sasl usernames.
* `users` - A list of sasl users. Each element contains the following attributes:
  * `username` - The username of the user.
  * `password` - The password of the user.
