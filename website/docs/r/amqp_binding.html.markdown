---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_binding"
sidebar_current: "docs-alicloud-resource-amqp-binding"
description: |-
  Provides a Alicloud RabbitMQ (AMQP) Binding resource.
---

# alicloud_amqp_binding

Provides a RabbitMQ (AMQP) Binding resource.

For information about RabbitMQ (AMQP) Binding and how to use it, see [What is Binding](https://www.alibabacloud.com/help/en/message-queue-for-rabbitmq/latest/createbinding).

-> **NOTE:** Available since v1.135.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_amqp_binding&exampleId=92808aec-04b2-679a-d6a1-186f9a7b69aee582f3d5&activeTab=example&spm=docs.r.amqp_binding.0.92808aec04&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-shanghai"
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

resource "alicloud_amqp_virtual_host" "default" {
  instance_id       = alicloud_amqp_instance.default.id
  virtual_host_name = "tf-example"
}

resource "alicloud_amqp_exchange" "default" {
  auto_delete_state = false
  exchange_name     = "tf-example"
  exchange_type     = "HEADERS"
  instance_id       = alicloud_amqp_instance.default.id
  internal          = false
  virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
}

resource "alicloud_amqp_queue" "default" {
  instance_id       = alicloud_amqp_instance.default.id
  queue_name        = "tf-example"
  virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
}

resource "alicloud_amqp_binding" "default" {
  argument          = "x-match:all"
  binding_key       = alicloud_amqp_queue.default.queue_name
  binding_type      = "QUEUE"
  destination_name  = "tf-example"
  instance_id       = alicloud_amqp_instance.default.id
  source_exchange   = alicloud_amqp_exchange.default.exchange_name
  virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the instance.
* `virtual_host_name` - (Required, ForceNew) The name of the vhost.
* `source_exchange` - (Required, ForceNew) The name of the source exchange.
* `destination_name` - (Required, ForceNew) The name of the object that you want to bind to the source exchange.
* `binding_type` - (Required, ForceNew) The type of the object that you want to bind to the source exchange. Valid values: `EXCHANGE`, `QUEUE`.
* `binding_key` - (Required, ForceNew) The Binding Key.
  * For a non-topic source exchange: The binding key can contain only letters, digits, hyphens (-), underscores (_), periods (.), and at signs (@).
    The binding key must be 1 to 255 characters in length.
  * For a topic source exchange: The binding key can contain letters, digits, hyphens (-), underscores (_), periods (.), and at signs (@).
    If the binding key contains a number sign (#), the binding key must start with a number sign (#) followed by a period (.) or end with a number sign (#) that follows a period (.).
    The binding key must be 1 to 255 characters in length.
* `argument` - (Optional, ForceNew) The key-value pairs that are configured for the headers attributes of a message. Default value: `x-match:all`. Valid values:
  - `x-match:all`: A headers exchange routes a message to a queue only if all binding attributes of the queue except for x-match match the headers attributes of the message.
  - `x-match:any`: A headers exchange routes a message to a queue if one or more binding attributes of the queue except for x-match match the headers attributes of the message. 
  - -> **NOTE:** If the exchange type is not 'HEADERS', the `argument` should not been set, otherwise, there are always "forces replacement" changes.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Binding. It formats as `<instance_id>:<virtual_host_name>:<source_exchange>:<destination_name>`.

## Import

RabbitMQ (AMQP) Binding can be imported using the id, e.g.

```shell
$ terraform import alicloud_amqp_binding.example <instance_id>:<virtual_host_name>:<source_exchange>:<destination_name>
```
