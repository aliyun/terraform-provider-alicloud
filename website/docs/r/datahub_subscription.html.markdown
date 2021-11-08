---
subcategory: "Datahub Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_datahub_subscription"
sidebar_current: "docs-alicloud-resource-datahub-subscription"
description: |-
  Provides a Alicloud datahub subscription resource.
---

# alicloud\_datahub\_subscription

The subscription is the basic unit of resource usage in Datahub Service under Publish/Subscribe model. You can manage the relationships between user and topics by using subscriptions. [Refer to details](https://help.aliyun.com/document_detail/47440.html).

## Example Usage

Basic Usage

```
resource "alicloud_datahub_subscription" "example" {
  project_name = "tf_datahub_project"
  topic_name   = "tf_datahub_topic"
  comment      = "created by terraform"
}
```

## Argument Reference

The following arguments are supported:

* `project_name` - (Required, ForceNew) The name of the datahub project that the subscription belongs to. Its length is limited to 3-32 and only characters such as letters, digits and '_' are allowed. It is case-insensitive.
* `topic_name` - (Required, ForceNew) The name of the datahub topic that the subscription belongs to. Its length is limited to 1-128 and only characters such as letters, digits and '_' are allowed. It is case-insensitive.
* `comment` - (Optional) Comment of the datahub subscription. It cannot be longer than 255 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the datahub subscription as terraform resource. It was composed of project name, topic name and practical subscription ID generated from server side. Format to `<project_name>:<topic_name>:<sub_id>`.
* `sub_id` - The identidy of the subscription, generate from server side.
* `create_time` - Create time of the datahub subscription. It is a human-readable string rather than 64-bits UTC.
* `last_modify_time` - Last modify time of the datahub subscription. It is the same as *create_time* at the beginning. It is also a human-readable string rather than 64-bits UTC.

## Import

Datahub subscription can be imported using the ID, e.g.

```
$ terraform import alicloud_datahub_subscription.example tf_datahub_project:tf_datahub_topic:1539073399567UgCzY
```
