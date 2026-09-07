---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_mns_topics"
sidebar_current: "docs-alicloud-datasource-mns-topics"
description: |-
  Provides a list of mns topics available to the user.
---

# alicloud\_mns\_topics

This data source provides a list of MNS topics in an Alibaba Cloud account according to the specified parameters.

-> **DEPRECATED:**  This datasource has been deprecated from version `1.188.0`. Please use new datasource [message_service_topics](https://www.terraform.io/docs/providers/alicloud/d/message_service_topics).

## Example Usage

```terraform
data "alicloud_mns_topics" "topics" {
  name_prefix = "tf-"
}

output "first_topic_id" {
  value = "${data.alicloud_mns_topics.topics.topics.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_prefix` - (Optional) A string to filter resulting topics by their name prefixs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of topic names.
* `topics` - A list of topics. Each element contains the following attributes:
   * `id` - The id of the topic. The value is set to `name`.
   * `name` - The name of the topic.
   * `maximum_message_size` - This indicates the maximum length, in bytes, of any message body sent to the topic.
   * `logging_enabled` - Whether to enable logging.
