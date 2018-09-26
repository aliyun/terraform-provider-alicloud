---
layout: "alicloud"
page_title: "Alicloud: alicloud_mns_topics"
sidebar_current: "docs-alicloud-datasource-mns-topics"
description: |-
    Provides a list of mns topics available to the user.
---

# alicloud\_msn\_topics

This data source provides a list of MNS topics  in an Alibaba Cloud account according to the specified parameters.

## Example Usage

```
data "alicloud_mns_topics" "topics" {
  name_prefix = "tf-"
}

output "first_topic_id" {
  value = "${data.alicloud_mns_topics.topics.topics.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_prefix` - (Optional) A  string to filter resulting topics by their name prefixs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `topics` - A list of users. Each element contains the following attributes:
   * `name` - Two topics on a single account in the same region cannot have the same name. A topic name must start with an English letter or a digit, and can contain English letters, digits, and hyphens, with the length not exceeding 256 characters.
   * `maximum_message_size` - This indicates the maximum length, in bytes, of any message body sent to the topic. Valid value range: 1024-65536, i.e., 1K to 64K.
   * `logging_enabled` - is logging enabled? true or false default value is false
