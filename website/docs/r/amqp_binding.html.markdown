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

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_amqp_instance" "default" {
  instance_type  = "enterprise"
  max_tps        = 3000
  queue_capacity = 200
  storage_size   = 700
  support_eip    = false
  max_eip_tps    = 128
  payment_type   = "Subscription"
}

resource "alicloud_amqp_virtual_host" "default" {
  instance_id       = alicloud_amqp_instance.default.id
  virtual_host_name = var.name
}

resource "alicloud_amqp_exchange" "default" {
  instance_id       = alicloud_amqp_instance.default.id
  virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
  exchange_name     = var.name
  exchange_type     = "HEADERS"
  auto_delete_state = false
  internal          = false
}

resource "alicloud_amqp_queue" "default" {
  instance_id       = alicloud_amqp_instance.default.id
  virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
  queue_name        = var.name
}

resource "alicloud_amqp_binding" "default" {
  instance_id       = alicloud_amqp_instance.default.id
  virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
  source_exchange   = alicloud_amqp_exchange.default.exchange_name
  destination_name  = var.name
  binding_type      = "QUEUE"
  binding_key       = alicloud_amqp_queue.default.queue_name
  argument          = "x-match:all"
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

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Binding. It formats as `<instance_id>:<virtual_host_name>:<source_exchange>:<destination_name>`.

## Import

RabbitMQ (AMQP) Binding can be imported using the id, e.g.

```shell
$ terraform import alicloud_amqp_binding.example <instance_id>:<virtual_host_name>:<source_exchange>:<destination_name>
```
