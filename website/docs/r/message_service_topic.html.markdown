---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_topic"
description: |-
  Provides a Alicloud Message Service Topic resource.
---

# alicloud_message_service_topic

Provides a Message Service Topic resource.



For information about Message Service Topic and how to use it, see [What is Topic](https://www.alibabacloud.com/help/en/message-service/latest/createtopic).

-> **NOTE:** Available since v1.188.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_message_service_topic&exampleId=f4078b56-da8a-868c-bec2-ce95370ce435d8bac445&activeTab=example&spm=docs.r.message_service_topic.0.f4078b56da&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_message_service_topic" "default" {
  topic_name       = var.name
  max_message_size = 16888
  enable_logging   = true
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_message_service_topic&spm=docs.r.message_service_topic.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `enable_logging` - (Optional, Bool, Available since v1.241.0) Specifies whether to enable the logging feature. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `enable_sse` - (Optional, Bool, Available since v1.299.0) Specifies whether to enable server-side encryption (SSE) for the topic. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `kms_key_id` - (Optional, Available since v1.299.0) The ID of the KMS key. Optional; recorded by the server when `enable_sse` is `true` (the MNS-open API always applies the `AES-256-GCM` algorithm).
* `max_message_size` - (Optional, Int) The maximum length of the message that is sent to the topic. Default value: `65536`. Valid values: `1024` to `65536`. Unit: bytes.
* `sse_algorithm` - (Optional, Available since v1.299.0) The server-side encryption algorithm. Valid value: `AES-256-GCM`.
* `sse_type` - (Optional, Available since v1.299.0) The server-side encryption type. Valid value: `SMQ`.
-> **NOTE:** The SSE attributes (`enable_sse`, `kms_key_id`, `sse_algorithm`, `sse_type`) are returned by the MNS-open `GetTopicAttributes` API only in some regions — for example, `cn-shanghai` returned them during testing — and are omitted from the response in other regions — for example, `cn-hangzhou` and `cn-beijing` did not return them during testing. This is a region deployment difference of the MNS-open API and does not affect `SetTopicAttributes`. In regions where these fields are not returned, the read-back value is empty, so a configured SSE block may show a plan diff; use `lifecycle.ignore_changes` on the SSE attributes if needed.
* `tags` - (Optional, Map, Available since v1.241.0) A mapping of tags to assign to the resource.
* `topic_name` - (Required, ForceNew) The name of the topic.
* `topic_type` - (Optional, ForceNew, Computed, Available since v1.283.0) The type of the topic. Default value: `normal`. Valid values:
  - `normal`: Standard topic.
  - `fifo`: FIFO queue.
The following arguments will be discarded. Please use new fields as soon as possible:
* `logging_enabled` - (Deprecated since v1.241.0). Field `logging_enabled` has been deprecated from provider version 1.241.0. New field `enable_logging` instead.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Topic.
* `create_time` - (Available since v1.241.0) The time when the topic was created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Topic.
* `delete` - (Defaults to 5 mins) Used when delete the Topic.
* `update` - (Defaults to 5 mins) Used when update the Topic.

## Import

Message Service Topic can be imported using the id, e.g.

```shell
$ terraform import alicloud_message_service_topic.example <id>
```
