---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_queue"
sidebar_current: "docs-alicloud-resource-amqp-queue"
description: |-
  Provides a Alicloud RabbitMQ (AMQP) Queue resource.
---

# alicloud_amqp_queue

Provides a RabbitMQ (AMQP) Queue resource.

For information about RabbitMQ (AMQP) Queue and how to use it, see [What is Queue](https://www.alibabacloud.com/help/en/message-queue-for-rabbitmq/latest/createqueue).

-> **NOTE:** Available since v1.127.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_amqp_queue&exampleId=160d3011-a66d-c43b-5206-4eaf762f291aa6ccc2c3&activeTab=example&spm=docs.r.amqp_queue.0.160d3011a6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_amqp_instance" "default" {
  instance_name         = "${var.name}-${random_integer.default.result}"
  instance_type         = "enterprise"
  max_tps               = 3000
  max_connections       = 2000
  queue_capacity        = 200
  payment_type          = "Subscription"
  renewal_status        = "AutoRenewal"
  renewal_duration      = 1
  renewal_duration_unit = "Year"
  support_eip           = true
}

resource "alicloud_amqp_virtual_host" "default" {
  instance_id       = alicloud_amqp_instance.default.id
  virtual_host_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_amqp_queue" "default" {
  instance_id       = alicloud_amqp_instance.default.id
  virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
  queue_name        = "${var.name}-${random_integer.default.result}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the ApsaraMQ for RabbitMQ instance to which the queue belongs.
* `virtual_host_name` - (Required, ForceNew) The name of the vhost to which the queue belongs. The name can contain only letters, digits, hyphens (-), underscores (_), periods (.), number signs (#), forward slashes (/), and at signs (@). The name must be 1 to 255 characters in length.
* `queue_name` - (Required, ForceNew) The name of the queue to create.
* `auto_delete_state` - (Optional, Bool, ForceNew) Specifies whether to automatically delete the queue. Valid values:
  - `true`: The queue is automatically deleted after the last consumer unsubscribes from it.
  - `false`: The queue is not automatically deleted.
* `max_length` - (Optional) The maximum number of messages that can be stored in the queue.
* `maximum_priority` - (Optional, Int) The priority of the queue.
* `message_ttl` - (Optional) The time to live (TTL) of a message in the queue.
* `dead_letter_exchange` - (Optional) The dead-letter exchange.
* `dead_letter_routing_key` - (Optional) The dead-letter routing key.
* `auto_expire_state` - (Optional) The auto-expiration time for the queue.
* `exclusive_state` - (Removed since v1.266.0) Field `exclusive_state` has been removed from provider version 1.266.0.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Queue. It formats as `<instance_id>:<virtual_host_name>:<queue_name>`.

## Import

RabbitMQ (AMQP) Queue can be imported using the id, e.g.

```shell
$ terraform import alicloud_amqp_queue.example <instance_id>:<virtual_host_name>:<queue_name>
```
