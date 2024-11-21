---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_queue"
description: |-
  Provides a Alicloud Message Service Queue resource.
---

# alicloud_message_service_queue

Provides a Message Service Queue resource. 

For information about Message Service Queue and how to use it, see [What is Queue](https://www.alibabacloud.com/help/en/message-service/latest/createqueue).

-> **NOTE:** Available since v1.188.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_message_service_queue&exampleId=035c32fa-29c9-0589-56e8-a39bc719314c356e28b3&activeTab=example&spm=docs.r.message_service_queue.0.035c32fa29&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_message_service_queue" "default" {
  delay_seconds            = "2"
  polling_wait_seconds     = "2"
  message_retention_period = "566"
  maximum_message_size     = "1123"
  visibility_timeout       = "30"
  queue_name               = var.name
}
```

## Argument Reference

The following arguments are supported:
* `delay_seconds` - (Optional, Computed) This means that messages sent to the queue can only be consumed after the delay time set by this parameter, in seconds.
* `logging_enabled` - (Optional) Represents whether the log management function is enabled.
* `maximum_message_size` - (Optional, Computed) Represents the maximum length of the message body sent to the Queue, in Byte.
* `message_retention_period` - (Optional, Computed) Represents the longest life time of the message in the Queue.
* `polling_wait_seconds` - (Optional, Computed) The longest waiting time for a Queue request when the number of messages is empty, in seconds.
* `queue_name` - (Required, ForceNew) Representative resources.
* `visibility_timeout` - (Optional, Computed) Represents the duration after the message is removed from the Queue and changed from the Active state to the Inactive state.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Represents the time when the Queue was created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Queue.
* `delete` - (Defaults to 5 mins) Used when delete the Queue.
* `update` - (Defaults to 5 mins) Used when update the Queue.

## Import

Message Service Queue can be imported using the id, e.g.

```shell
$ terraform import alicloud_message_service_queue.example <id>
```