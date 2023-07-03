---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_virtual_host"
sidebar_current: "docs-alicloud-resource-amqp-virtual-host"
description: |-
  Provides a Alicloud RabbitMQ (AMQP) Virtual Host resource.
---

# alicloud_amqp_virtual_host

Provides a RabbitMQ (AMQP) Virtual Host resource.

For information about RabbitMQ (AMQP) Virtual Host and how to use it, see [What is Virtual Host](https://www.alibabacloud.com/help/en/message-queue-for-rabbitmq/latest/createvirtualhost).

-> **NOTE:** Available since v1.126.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_amqp_instance" "default" {
  instance_type  = "professional"
  max_tps        = 1000
  queue_capacity = 50
  support_eip    = true
  max_eip_tps    = 128
  payment_type   = "Subscription"
  period         = 1
}

resource "alicloud_amqp_virtual_host" "default" {
  instance_id       = alicloud_amqp_instance.default.id
  virtual_host_name = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) InstanceId.
* `virtual_host_name` - (Required, ForceNew) VirtualHostName.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Virtual Host. The value is formatted `<instance_id>:<virtual_host_name>`.

## Import

RabbitMQ (AMQP) Virtual Host can be imported using the id, e.g.

```shell
$ terraform import alicloud_amqp_virtual_host.example <instance_id>:<virtual_host_name>
```
