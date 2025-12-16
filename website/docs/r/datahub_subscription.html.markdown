---
subcategory: "Datahub Service (DataHub)"
layout: "alicloud"
page_title: "Alicloud: alicloud_datahub_subscription"
sidebar_current: "docs-alicloud-resource-datahub-subscription"
description: |-
  Provides a Alicloud datahub subscription resource.
---

# alicloud_datahub_subscription

The subscription is the basic unit of resource usage in Datahub Service under Publish/Subscribe model. You can manage the relationships between user and topics by using subscriptions. [Refer to details](https://www.alibabacloud.com/help/en/datahub/latest/nerbcz).

-> **NOTE:** Available since v1.19.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_datahub_subscription&exampleId=f7816e86-5c20-c083-4f5e-3d106803e28671a12572&activeTab=example&spm=docs.r.datahub_subscription.0.f7816e865c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform_example"
}
resource "alicloud_datahub_project" "example" {
  name    = var.name
  comment = "created by terraform"
}

resource "alicloud_datahub_topic" "example" {
  name         = var.name
  project_name = alicloud_datahub_project.example.name
  record_type  = "BLOB"
  shard_count  = 3
  life_cycle   = 7
  comment      = "created by terraform"
}

resource "alicloud_datahub_subscription" "example" {
  project_name = alicloud_datahub_project.example.name
  topic_name   = alicloud_datahub_topic.example.name
  comment      = "created by terraform"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_datahub_subscription&spm=docs.r.datahub_subscription.example&intl_lang=EN_US)

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

```shell
$ terraform import alicloud_datahub_subscription.example tf_datahub_project:tf_datahub_topic:1539073399567UgCzY
```
