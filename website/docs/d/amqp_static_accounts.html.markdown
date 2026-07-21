---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ons_static_accounts"
sidebar_current: "docs-alicloud-datasource-ons-static-accounts"
description: |-
  Provides a list of Amqp Static Account owned by an Alibaba Cloud account.
---

# alicloud_amqp_static_accounts

This data source provides Amqp Static Account available to the user.[What is Static Account](https://help.aliyun.com/document_detail/184399.html)

-> **NOTE:** Available in 1.195.0+

## Example Usage

```
data "alicloud_amqp_static_accounts" "default" {
  instance_id = "amqp-cn-0ju2y01zs001"
}

output "alicloud_amqp_static_account_example_id" {
  value = data.alicloud_amqp_static_accounts.default.accounts.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed)  The `key` of the resource supplied above.The value is formulated as `<instance_id>:<access_key>`.
* `instance_id` - (ForceNew, Optional) InstanceId
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `accounts` - A list of Static Account Entries. Each element contains the following attributes:
    * `id` - The `key` of the resource supplied above.The value is formulated as `<instance_id>:<access_key>`.
    * `access_key` - Access key.
    * `create_time` - Create time stamp. Unix timestamp, to millisecond level.
    * `instance_id` - Amqp instance ID.
    * `master_uid` - The ID of the user's primary account.
    * `password` - Static password.
    * `user_name` - Static username.