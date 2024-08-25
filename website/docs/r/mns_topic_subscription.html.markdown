---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_mns_topic_subscription"
sidebar_current: "docs-alicloud-resource-mns-topic_subscription"
description: |-
  Provides a Alicloud MNS Topic Subscription resource.
---

# alicloud\_mns\_topic\_subscription

Provides a MNS topic subscription resource.

-> **NOTE:** Terraform will auto build a mns topic subscription  while it uses `alicloud_mns_topic_subscription` to build a mns topic subscription resource.

-> **DEPRECATED:**  This resource has been deprecated from version `1.188.0`. Please use new resource [message_service_subscription](https://www.terraform.io/docs/providers/alicloud/r/message_service_subscription).

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_mns_topic_subscription&exampleId=93f921bc-5cd5-9cc9-fc01-e9e7af91c910c21f48a5&activeTab=example&spm=docs.r.mns_topic_subscription.0.93f921bc5c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_mns_topic" "topic" {
  name                 = "tf-example-mnstopic"
  maximum_message_size = 65536
  logging_enabled      = false
}

resource "alicloud_mns_topic_subscription" "subscription" {
  topic_name            = "tf-example-mnstopic"
  name                  = "tf-example-mnstopic-sub"
  filter_tag            = "test"
  endpoint              = "http://www.xxx.com/xxx"
  notify_strategy       = "BACKOFF_RETRY"
  notify_content_format = "XML"
}
```

## Argument Reference

The following arguments are supported:

* `topic_name`- (Required, ForceNew) The topic which The subscription belongs to was named with the name.A topic name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 256 characters.
* `name` - (Required, ForceNew) Two topics subscription on a single account in the same topic cannot have the same name. A topic subscription name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 256 characters.
* `notify_strategy` - (Optional) The NotifyStrategy attribute of Subscription. This attribute specifies the retry strategy when message sending fails. The Valid values: `EXPONENTIAL_DECAY_RETRY` and `BACKOFF_RETRY`. Default value to `BACKOFF_RETRY` .
* `notify_content_format` - (Optional, ForceNew) The NotifyContentFormat attribute of Subscription. This attribute specifies the content format of the messages pushed to users. The valid values: `SIMPLIFIED`, `XML` and `JSON`. Default to `SIMPLIFIED`.
* `endpoint` - (Required, ForceNew) The endpoint has three format. Available values format:
- `HTTP Format`: http://xxx.com/xxx
- `Queue Format`: acs:mns:{REGION}:{AccountID}:queues/{QueueName}
- `Email Format`: mail:directmail:{MailAddress}

* `filter_tag` - (Optional, ForceNew) The length should be shorter than 16.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the topic subscription.Format to topic_name:name

## Import

MNS Topic subscription can be imported using the id, e.g.

```shell
$ terraform import alicloud_mns_topic_subscription.subscription tf-example-mnstopic:tf-example-mnstopic-sub
```
