---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_queues"
sidebar_current: "docs-alicloud-datasource-amqp-queues"
description: |-
  Provides a list of Amqp Queues to the user.
---

# alicloud\_amqp\_queues

This data source provides the Amqp Queues of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.127.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_amqp_queues" "ids" {
  instance_id       = "amqp-abc12345"
  virtual_host_name = "my-VirtualHost"
  ids               = ["my-Queue-1", "my-Queue-2"]
}
output "amqp_queue_id_1" {
  value = data.alicloud_amqp_queues.ids.queues.0.id
}

data "alicloud_amqp_queues" "nameRegex" {
  instance_id       = "amqp-abc12345"
  virtual_host_name = "my-VirtualHost"
  name_regex        = "^my-Queue"
}
output "amqp_queue_id_2" {
  value = data.alicloud_amqp_queues.nameRegex.queues.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Queue IDs. Its element value is same as Queue Name.
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Queue name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `virtual_host_name` - (Required, ForceNew) The name of the virtual host.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Queue names.
* `queues` - A list of Amqp Queues. Each element contains the following attributes:
	* `attributes` - The attributes for the Queue.
	* `auto_delete_state` - Specifies whether the Auto Delete attribute is configured.
	* `create_time` - CreateTime.
	* `exclusive_state` - Specifies whether the queue is an exclusive queue.
	* `id` - The ID of the Queue. Its value is same as Queue Name.
	* `instance_id` - The ID of the instance.
	* `last_consume_time` - The last consume time.
	* `queue_name` - The queue name.
	* `virtual_host_name` - The name of the virtual host.
