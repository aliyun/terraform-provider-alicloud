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

For information about Amqp Static Account and how to use it, see [What is Static Account](https://www.alibabacloud.com/help/en/message-queue-for-rabbitmq/latest/create-a-pair-of-static-username-and-password).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_amqp_static_account&exampleId=919dcffb-0ea4-f6df-823a-d8cfd79ddce79a715064&activeTab=example&spm=docs.r.amqp_static_account.0.919dcffb0e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}
variable "access_key" {
  default = "access_key"
}
variable "secret_key" {
  default = "secret_key"
}
resource "alicloud_amqp_instance" "default" {
  instance_type  = "enterprise"
  max_tps        = 3000
  queue_capacity = 200
  storage_size   = 700
  support_eip    = false
  max_eip_tps    = 128
  payment_type   = "Subscription"
  period         = 1
}
resource "alicloud_amqp_static_account" "default" {
  instance_id = alicloud_amqp_instance.default.id
  access_key  = var.access_key
  secret_key  = var.secret_key
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) Amqp instance ID.
* `access_key` - (Required, ForceNew) Access key.
* `secret_key` - (Required, ForceNew) Secret key.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Static Account. It formats as `<instance_id>:<access_key>`.
* `user_name` - The static username.
* `password` - The static password.
* `master_uid` - The ID of the user's primary account.
* `create_time` - The timestamp that indicates when the pair of static username and password was created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Static Account.
* `delete` - (Defaults to 5 mins) Used when delete the Static Account.

## Import

Amqp Static Account can be imported using the id, e.g.

```shell
$terraform import alicloud_amqp_static_account.example <instance_id>:<access_key>
```
