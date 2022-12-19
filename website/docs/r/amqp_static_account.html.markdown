---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_static_account"
sidebar_current: "docs-alicloud-resource-amqp-static-account"
description: |-
  Provides a Alicloud Amqp Static Account resource.
---

# alicloud_amqp_static_account

Provides a Amqp Static Account resource.

For information about Amqp Static Account and how to use it, see [What is Static Account](https://help.aliyun.com/document_detail/184399.html).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_amqp_static_account" "default" {
  instance_id = "amqp-cn-0ju2y01zs001"
  access_key  = "LTAI5t8beMmVM1eRZtEJ6vfo"
  secret_key  = "sample-secret-key"
}
```

## Argument Reference

The following arguments are supported:
* `access_key` - (Required,ForceNew) Access key.
* `instance_id` - (Required,ForceNew) Amqp instance ID.
* `secret_key` - (Required,ForceNew) Secret key.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.The value is formulated as `<instance_id>:<access_key>`.
* `username` - Static user name.
* `password` - Static password.
* `create_time` - Create time stamp. Unix timestamp, to millisecond level.
* `master_uid` - The ID of the user's primary account.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Static Account.
* `delete` - (Defaults to 5 mins) Used when delete the Static Account.

## Import

Amqp Static Account can be imported using the id, e.g.

```shell
$terraform import alicloud_amqp_static_account.example <instance_id>:<access_key>
```