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
  exchange_type     = "DIRECT"
  instance_id       = alicloud_amqp_instance.default.id
  internal          = false
  virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
}

resource "alicloud_amqp_queue" "example" {
  instance_id       = alicloud_amqp_instance.default.id
  queue_name        = "tf-example"
  virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
}
```

## Argument Reference

The following arguments are supported:

* `auto_delete_state` - (Optional, ForceNew) Specifies whether the Auto Delete attribute is configured. Valid values:
  * true: The Auto Delete attribute is configured. The queue is automatically deleted after the last subscription from consumers to this queue is canceled. 
  * false: The Auto Delete attribute is not configured.
  
* `auto_expire_state` - (Optional) The validity period after which the queue is automatically deleted.
  If the queue is not accessed within a specified period of time, it is automatically deleted.
* `dead_letter_exchange` - (Optional) The dead-letter exchange. A dead-letter exchange is used to receive rejected messages. 
  If a consumer rejects a message that cannot be retried, this message is routed to a specified dead-letter exchange.
  Then, the dead-letter exchange routes the message to the queue that is bound to the dead-letter exchange.
* `dead_letter_routing_key` - (Optional) The dead letter routing key.
* `exclusive_state` - (Optional, ForceNew) Specifies whether the queue is an exclusive queue. Valid values:
  * true: The queue is an exclusive queue. It can be used only for the connection that declares the exclusive queue. After the connection is closed, the exclusive queue is automatically deleted.
  * false: The queue is not an exclusive queue.
  
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `max_length` - (Optional) The maximum number of messages that can be stored in the queue.
  If this threshold is exceeded, the earliest messages that are routed to the queue are discarded.
* `maximum_priority` - (Optional) The highest priority supported by the queue. This parameter is set to a positive integer.
  Valid values: 0 to 255. Recommended values: 1 to 10
* `message_ttl` - (Optional) The message TTL of the queue.
  If the retention period of a message in the queue exceeds the message TTL of the queue, the message expires.
  Message TTL must be set to a non-negative integer, in milliseconds.
  For example, if the message TTL of the queue is 1000, messages survive for at most 1 second in the queue.
* `queue_name` - (Required, ForceNew) The name of the queue.
  The queue name must be 1 to 255 characters in length, and can contain only letters, digits, hyphens (-), underscores (_), periods (.), and at signs (@).
* `virtual_host_name` - (Required, ForceNew) The name of the virtual host.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Queue. The value formats as `<instance_id>:<virtual_host_name>:<queue_name>`.

## Import

RabbitMQ (AMQP) Queue can be imported using the id, e.g.

```shell
$ terraform import alicloud_amqp_queue.example <instance_id>:<virtual_host_name>:<queue_name>
```
