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
  <a href="https://api.aliyun.com/terraform?resource=alicloud_message_service_queue&exampleId=d0dc31e6-0335-9ff7-e5d4-87536fb991ca86f6a8ac&activeTab=example&spm=docs.r.message_service_queue.0.d0dc31e603&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_message_service_queue" "default" {
  queue_name               = var.name
  delay_seconds            = "2"
  polling_wait_seconds     = "2"
  message_retention_period = "566"
  maximum_message_size     = "1126"
  visibility_timeout       = "30"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_message_service_queue&spm=docs.r.message_service_queue.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `delay_seconds` - (Optional, Int) The period after which all messages sent to the queue are consumed. Default value: `0`. Valid values: `0` to `604800`. Unit: seconds.
* `dlq_policy` - (Optional, Set, Available since v1.244.0) The dead-letter queue policy. See [`dlq_policy`](#dlq_policy) below.
* `enable_sse` - (Optional, Bool, Available since v1.291.0) Specifies whether to enable server-side encryption (SSE) for the messages in the queue. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `kms_key_id` - (Optional, Available since v1.291.0) The ID of the customer master key (CMK) in Key Management Service (KMS). This parameter is required when `sse_type` is set to `KMS`.
* `logging_enabled` - (Optional, Bool) Specifies whether to enable the logging feature. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `maximum_message_size` - (Optional, Int) The maximum length of the message that is sent to the queue. Valid values: `1024` to `65536`. Unit: bytes. Default value: `65536`.
* `message_retention_period` - (Optional, Int) The maximum duration for which a message is retained in the queue. After the specified retention period ends, the message is deleted regardless of whether the message is received. Valid values: `60` to `604800`. Unit: seconds. Default value: `345600`.
* `polling_wait_seconds` - (Optional, Int) The maximum duration for which long polling requests are held after the ReceiveMessage operation is called. Valid values: `0` to `30`. Unit: seconds. Default value: `0`.
* `queue_name` - (Required, ForceNew) The name of the queue.
* `queue_type` - (Optional, ForceNew, Available since v1.283.0) The type of the queue. Default value: `normal`. Valid values:
  - `normal`: Standard queue.
  - `fifo`: FIFO queue.
* `sse_algorithm` - (Optional, Available since v1.291.0) The encryption algorithm that is used to encrypt the messages in the queue. Valid value: `AES-256-GCM`.
* `sse_type` - (Optional, Available since v1.291.0) The type of server-side encryption (SSE). Valid values:
  - `SMQ`: The keys that are managed by Message Service are used to encrypt and decrypt messages.
  - `KMS`: The customer master key (CMK) that is managed by Key Management Service (KMS) is used to encrypt and decrypt messages.
* `tags` - (Optional, Map, Available since v1.241.0) A mapping of tags to assign to the resource.
* `visibility_timeout` - (Optional, Int) The duration for which a message stays in the Inactive state after the message is received from the queue. Valid values: `1` to `43200`. Unit: seconds. Default value: `30`.

-> **NOTE:** Server-side encryption (SSE) is available only after the feature is enabled for your account. To enable the feature, submit a ticket. Before you set `sse_type` to `KMS`, make sure that a symmetric key of the AES-256 key spec is created in KMS in the same region and within the same account as the queue, and that the service-linked role `AliyunMNSAccessingKMSRole` is authorized.

-> **NOTE:** To disable SSE, you must set `enable_sse` to `false` and remove `sse_type`, `sse_algorithm`, and `kms_key_id` from the configuration at the same time. SSE takes effect only on the messages that are sent after SSE is enabled or disabled, and existing messages are not encrypted or decrypted again.

### `dlq_policy`

The dlq_policy supports the following:
* `dead_letter_target_queue` - (Optional) The queue to which dead-letter messages are delivered.
* `enabled` - (Optional, Bool) Specifies whether to enable the dead-letter message delivery. Valid values: `true`, `false`.
* `max_receive_count` - (Optional, Int) The maximum number of retries.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Queue.
* `create_time` - (Available since v1.223.2) The time when the queue was created.
* `encryption_enabled` - (Available since v1.291.0) Indicates whether server-side encryption is applied to the queue. The value remains `true` after server-side encryption is disabled because existing messages are still stored as encrypted.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Queue.
* `delete` - (Defaults to 5 mins) Used when delete the Queue.
* `update` - (Defaults to 5 mins) Used when update the Queue.

## Import

Message Service Queue can be imported using the id, e.g.

```shell
$ terraform import alicloud_message_service_queue.example <id>
```
